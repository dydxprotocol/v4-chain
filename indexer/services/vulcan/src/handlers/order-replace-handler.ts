import {
  ParseMessageError, logger, runFuncWithTimingStat, stats,
} from '@dydxprotocol-indexer/base';
import {
  createSubaccountWebsocketMessage,
  getTriggerPrice,
  KafkaTopics,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import {
  IsoString,
  OrderFromDatabase,
  OrderStatus,
  OrderTable,
  PerpetualMarketFromDatabase,
  SubaccountMessageContents,
  SubaccountTable,
  TimeInForce,
  apiTranslations,
  blockHeightRefresher,
  perpetualMarketRefresher,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import {
  CanceledOrdersCache,
  OrderbookLevelsCache,
  placeOrder,
  PlaceOrderResult,
  StatefulOrderUpdatesCache,
  convertToRedisOrder,
  removeOrder,
  RemoveOrderResult,
} from '@dydxprotocol-indexer/redis';
import {
  getOrderIdHash,
  isLongTermOrder,
  isStatefulOrder,
  ORDER_FLAG_SHORT_TERM,
  requiresImmediateExecution,
} from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerOrderId,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  OrderRemovalReason,
  OrderReplaceV1,
  OrderUpdateV1,
  RedisOrder,
  SubaccountMessage,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import { IHeaders, Message } from 'kafkajs';

import config from '../config';
import { redisClient } from '../helpers/redis/redis-controller';
import { sendMessageWrapper } from '../lib/send-message-helper';
import { Handler } from './handler';
import { getStateRemainingQuantums } from './helpers';

/**
 * Handler for OrderReplace messages. This is currently only expected for stateful vault orders.
 * The behavior is as follows:
 * - Remove the old order
 *  - this is done using the `removeOrder` function from the `redis` package
 * - Add the new order to the OrdersCache, OrdersDataCache, and SubaccountOrderIdsCache
 *  - this is done using the `placeOrder` function from the `redis` package
 *  - Remove the order from the CanceledOrdersCache if it exists
 * - If the order is a stateful order, attempt to remove any cached order update from the
 *   StatefulOrderUpdatesCache, and then queue the order update to be re-sent and re-processed
 * - If the order doesn't already exist in the caches, return
 * - If the order exists in the caches, but was not replaced due to the expiry of the existing order
 *   being greater than or equal to the expiry of the order in the OrderReplace message, return
 */
export class OrderReplaceHandler extends Handler {
  protected async handle(update: OffChainUpdateV1, headers: IHeaders): Promise<void> {
    logger.info({
      at: 'OrderReplaceHandler#handle',
      message: 'Received OffChainUpdate with OrderReplace.',
      update,
      txHash: this.txHash,
    });
    const orderReplace: OrderReplaceV1 = update.orderReplace!;
    this.validateOrderReplace(orderReplace);
    const oldOrderId: IndexerOrderId = orderReplace.oldOrderId!;

    /* Remove old order */
    let removeOrderResult: RemoveOrderResult = { removed: false };
    try {
      removeOrderResult = await this.removeOldOrder(oldOrderId, headers);
    } catch (error) {
      if (error instanceof ParseMessageError) {
        return;
      }
    }

    /* We don't want to fail if old order is not found (new order should still be placed),
    so log and track metric */
    if (!removeOrderResult.removed) {
      logger.info({
        at: 'OrderReplaceHandler#handle',
        message: 'Old order not found in cache',
        oldOrderId,
      });
      stats.increment(`${config.SERVICE_NAME}.replace_order_handler.old_order_not_found_in_cache`, 1);
    }

    /* Place new order */
    const order: IndexerOrder = orderReplace.order!;
    const placementStatus: OrderPlaceV1_OrderPlacementStatus = orderReplace.placementStatus;
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(order.orderId!.clobPairId.toString());
    if (perpetualMarket === undefined) {
      this.logAndThrowParseMessageError(
        'Order in OrderReplace has invalid clobPairId',
        {
          order,
        },
      );
      // Needed so TS compiler knows `perpetualMarket` is defined
      return;
    }
    const redisOrder: RedisOrder = convertToRedisOrder(order, perpetualMarket);
    const placeOrderResult: PlaceOrderResult = await this.placeNewOrder(redisOrder);
    await this.removeOrderFromCanceledOrdersCache(
      OrderTable.orderIdToUuid(redisOrder.order?.orderId!),
    );

    if (placeOrderResult.replaced) {
      // This is not expected because the replaced orders either have different order IDs or
      // should have been removed before being placed again
      stats.increment(`${config.SERVICE_NAME}.replace_order_handler.place_order_result_replaced_order`, 1);
    }

    // If an order was removed from the Orders cache and was resting on the book, update the
    // orderbook levels cache
    // Orders that require immediate execution do not rest on the book, and also should not lead
    // to an update to the orderbook levels cache
    if (
      removeOrderResult.removed &&
      removeOrderResult.restingOnBook === true &&
      !requiresImmediateExecution(removeOrderResult.removedOrder!.order!.timeInForce)
    ) {
      // Don't send orderbook message if price is the same to prevent flickering because
      // the order update will send the correct size update
      const sendOrderbookMessage: boolean = (
        redisOrder.order!.subticks.neq(removeOrderResult.removedOrder!.order!.subticks)
      );
      if (sendOrderbookMessage) {
        logger.info({
          at: 'OrderReplaceHandler#handle',
          message: 'Sending orderbook message because price is different',
          redisOrder,
          removedOrder: removeOrderResult.removedOrder!.order,
        });
      }
      await this.removeOldOrderFromOrderbook(
        removeOrderResult,
        perpetualMarket,
        headers,
        sendOrderbookMessage,
      );
    }

    // TODO(CLOB-597): Remove this logic and log erorrs once best-effort-open is not sent for
    // stateful orders in the protocol
    if (this.shouldSendSubaccountMessage(
      orderReplace,
      placeOrderResult,
      placementStatus,
      redisOrder,
    )) {
      // TODO(IND-171): Determine whether we should always be sending a message, even when the cache
      // isn't updated.
      // For stateful and conditional orders, look the order up in the db for the createdAtHeight
      // and send any cached order updates for the stateful or conditional order
      let dbOrder: OrderFromDatabase | undefined;
      if (isStatefulOrder(redisOrder.order!.orderId!.orderFlags)) {
        const orderUuid: string = OrderTable.orderIdToUuid(redisOrder.order!.orderId!);
        dbOrder = await OrderTable.findById(orderUuid);
        if (dbOrder === undefined) {
          logger.crit({
            at: 'OrderReplaceHandler#createSubaccountWebsocketMessage',
            message: 'Stateful order not found in database',
          });
          throw new Error(`Stateful order not found in database: ${orderUuid}`);
        }
        await this.sendCachedOrderUpdate(orderUuid, headers);
      }
      const subaccountMessage: Message = {
        value: createSubaccountWebsocketMessage(
          redisOrder,
          dbOrder,
          perpetualMarket,
          placementStatus,
          blockHeightRefresher.getLatestBlockHeight(),
        ),
        headers,
      };
      sendMessageWrapper(subaccountMessage, KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS);
    }

    await this.addOrderToCanceledOrdersCache(
      oldOrderId,
      Date.now(),
    );
  }

  protected validateOrderReplace(orderReplace: OrderReplaceV1): void {
    if (orderReplace.order === undefined) {
      this.logAndThrowParseMessageError('Invalid OrderReplace, order is undefined');
      return;
    }

    if (orderReplace.order.orderId === undefined) {
      this.logAndThrowParseMessageError('Invalid OrderReplace, order id is undefined');
      return;
    }

    if (orderReplace.order.orderId.subaccountId === undefined) {
      this.logAndThrowParseMessageError('Invalid OrderReplace, subaccount id is undefined');
    }

    if (
      orderReplace
        .placementStatus === OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_UNSPECIFIED
    ) {
      this.logAndThrowParseMessageError('Invalid OrderReplace, placement status is UNSPECIFIED');
    }
  }

  protected async removeOldOrder(
    oldOrderId: IndexerOrderId,
    headers: IHeaders,
  ): Promise<RemoveOrderResult> {
    const order: OrderFromDatabase | undefined = await runFuncWithTimingStat(
      OrderTable.findById(
        OrderTable.orderIdToUuid(oldOrderId),
      ),
      this.generateTimingStatsOptions('find_order_for_stateful_cancelation'),
    );
    if (order === undefined) {
      logger.error({
        at: 'OrderReplaceHandler#removeOldOrder',
        message: 'Could not find old order ID to remove in DB',
        orderId: oldOrderId,
      });
      throw new ParseMessageError(`Could not find old order ID to remove in DB: ${oldOrderId}`);
    }

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(
        order.clobPairId,
      );
    if (perpetualMarket === undefined) {
      logger.error({
        at: 'OrderReplaceHandler#removeOldOrder',
        message: `Unable to find the perpetual market with clobPairId: ${order.clobPairId}`,
        order,
      });
      throw new ParseMessageError(`Unable to find the perpetual market with clobPairId: ${order.clobPairId}`);
    }

    const removeOrderResult: RemoveOrderResult = await runFuncWithTimingStat(
      removeOrder({
        removedOrderId: oldOrderId,
        client: redisClient,
      }),
      this.generateTimingStatsOptions('remove_order'),
    );
    logger.info({
      at: 'OrderReplaceHandler#handle',
      message: 'removeOrder processed',
      oldOrderId,
      removeOrderResult,
    });

    const stateRemainingQuantums: Big = await getStateRemainingQuantums(
      removeOrderResult.removedOrder!,
    );

    // If the remaining amount of the order in state is <= 0, the order is filled and
    // does not need to have it's status updated
    let canceledOrder: OrderFromDatabase | undefined;
    if (stateRemainingQuantums.gt(0)) {
      canceledOrder = await runFuncWithTimingStat(
        this.cancelOrderInPostgres(oldOrderId),
        this.generateTimingStatsOptions('cancel_order_in_postgres'),
      );
    } else {
      canceledOrder = await runFuncWithTimingStat(
        OrderTable.findById(OrderTable.orderIdToUuid(oldOrderId)),
        this.generateTimingStatsOptions('find_order'),
      );
    }

    const subaccountMessage: Message = {
      value: this.createSubaccountWebsocketMessageFromRemoveOrderResult(
        removeOrderResult,
        canceledOrder,
        oldOrderId,
        perpetualMarket,
        blockHeightRefresher.getLatestBlockHeight(),
      ),
      headers,
    };

    sendMessageWrapper(subaccountMessage, KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS);

    return removeOrderResult;
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  protected async cancelOrderInPostgres(
    oldOrderId: IndexerOrderId,
  ): Promise<OrderFromDatabase | undefined> {
    return OrderTable.update({
      id: OrderTable.orderIdToUuid(oldOrderId),
      status: OrderStatus.CANCELED,
    });
  }

  protected createSubaccountWebsocketMessageFromRemoveOrderResult(
    removeOrderResult: RemoveOrderResult,
    canceledOrder: OrderFromDatabase | undefined,
    oldOrderId: IndexerOrderId,
    perpetualMarket: PerpetualMarketFromDatabase,
    blockHeight: string | undefined,
  ): Buffer {
    const redisOrder: RedisOrder = removeOrderResult.removedOrder!;
    const orderTIF: TimeInForce = protocolTranslations.protocolOrderTIFToTIF(
      redisOrder.order!.timeInForce,
    );
    const createdAtHeight: string | undefined = canceledOrder?.createdAtHeight;
    const updatedAt: IsoString | undefined = canceledOrder?.updatedAt;
    const updatedAtHeight: string | undefined = canceledOrder?.updatedAtHeight;
    const contents: SubaccountMessageContents = {
      orders: [
        {
          id: OrderTable.orderIdToUuid(redisOrder.order!.orderId!),
          subaccountId: SubaccountTable.subaccountIdToUuid(
            oldOrderId.subaccountId!,
          ),
          clientId: oldOrderId.clientId.toString(),
          clobPairId: perpetualMarket.clobPairId,
          side: protocolTranslations.protocolOrderSideToOrderSide(redisOrder.order!.side),
          size: redisOrder.size,
          totalOptimisticFilled: protocolTranslations.quantumsToHumanFixedString(
            removeOrderResult.totalFilledQuantums!.toString(),
            perpetualMarket.atomicResolution,
          ),
          price: redisOrder.price,
          type: protocolTranslations.protocolConditionTypeToOrderType(
            redisOrder.order!.conditionType,
          ),
          status: OrderStatus.CANCELED,
          timeInForce: apiTranslations.orderTIFToAPITIF(orderTIF),
          postOnly: apiTranslations.isOrderTIFPostOnly(orderTIF),
          reduceOnly: redisOrder.order!.reduceOnly,
          orderFlags: redisOrder.order!.orderId!.orderFlags.toString(),
          goodTilBlock: protocolTranslations.getGoodTilBlock(redisOrder.order!)
            ?.toString(),
          goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(redisOrder.order!),
          ticker: redisOrder.ticker,
          removalReason: OrderRemovalReason[OrderRemovalReason.ORDER_REMOVAL_REASON_USER_CANCELED],
          ...(createdAtHeight && { createdAtHeight }),
          ...(updatedAt && { updatedAt }),
          ...(updatedAtHeight && { updatedAtHeight }),
          clientMetadata: redisOrder.order!.clientMetadata.toString(),
          triggerPrice: getTriggerPrice(redisOrder.order!, perpetualMarket),
        },
      ],
      ...(blockHeight && { blockHeight }),
    };

    const subaccountMessage: SubaccountMessage = SubaccountMessage.fromPartial({
      contents: JSON.stringify(contents),
      subaccountId: oldOrderId.subaccountId!,
      version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
    });

    return Buffer.from(Uint8Array.from(SubaccountMessage.encode(subaccountMessage).finish()));
  }

  protected async placeNewOrder(redisOrder: RedisOrder): Promise<PlaceOrderResult> {
    const placeOrderResult: PlaceOrderResult = await runFuncWithTimingStat(
      placeOrder({
        redisOrder,
        client: redisClient,
      }),
      this.generateTimingStatsOptions('place_order_cache_update'),
    );
    await this.removeOrderFromCanceledOrdersCache(
      OrderTable.orderIdToUuid(redisOrder.order?.orderId!),
    );
    logger.info({
      at: 'OrderReplaceHandler#handle',
      message: 'placeOrder processed',
      redisOrder,
      placeOrderResult,
    });
    return placeOrderResult;
  }

  /**
   * Determine whether to send a subaccount websocket message given the order place.
   * @param orderReplace
   * @returns TODO(CLOB-597): Remove once best-effort-opened messages are not sent for stateful
   * orders.
   */
  protected shouldSendSubaccountMessage(
    orderReplace: OrderReplaceV1,
    placeOrderResult: PlaceOrderResult,
    placementStatus: OrderPlaceV1_OrderPlacementStatus,
    redisOrder: RedisOrder,
  ): boolean {
    if (
      isLongTermOrder(redisOrder.order!.orderId!.orderFlags) &&
      !config.SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_STATEFUL_ORDERS
    ) {
      return false;
    }

    const orderFlags: number = orderReplace.order!.orderId!.orderFlags;
    const status: OrderPlaceV1_OrderPlacementStatus = orderReplace.placementStatus;
    // Best-effort-opened status should only be sent for short-term orders
    if (
      orderFlags !== ORDER_FLAG_SHORT_TERM &&
      status === OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED
    ) {
      return false;
    }

    // In the case where a stateful orderReplace is opened with a more recent expiry than an
    // existing order on the indexer, then the order will not have been placed or replaced and
    // no subaccount message should be sent.
    if (placeOrderResult.placed === false &&
      placeOrderResult.replaced === false &&
      placementStatus ===
      OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED) {
      return false;
    }
    return true;
  }

  /**
   * Removes the order from the cancelled orders cache in Redis.
   *
   * @param orderId
   * @protected
   */
  protected async removeOrderFromCanceledOrdersCache(
    orderId: string,
  ): Promise<void> {
    await runFuncWithTimingStat(
      CanceledOrdersCache.removeOrderFromCaches(orderId, redisClient),
      this.generateTimingStatsOptions('remove_order_from_cancel_cache'),
    );
  }

  /**
   * Removes and sends the cached order update for the given order id if it exists.
   *
   * @param orderId
   * @returns
   */
  protected async sendCachedOrderUpdate(
    orderId: string,
    headers: IHeaders,
  ): Promise<void> {
    const cachedOrderUpdate: OrderUpdateV1 | undefined = await StatefulOrderUpdatesCache
      .removeStatefulOrderUpdate(
        orderId,
        Date.now(),
        redisClient,
      );

    if (cachedOrderUpdate === undefined) {
      return;
    }

    const orderUpdateMessage: Message = {
      key: getOrderIdHash(cachedOrderUpdate.orderId!),
      value: Buffer.from(
        Uint8Array.from(OffChainUpdateV1.encode({ orderUpdate: cachedOrderUpdate }).finish()),
      ),
      headers,
    };
    sendMessageWrapper(orderUpdateMessage, KafkaTopics.TO_VULCAN);
  }

  /**
   * Updates the orderbook and sends a message to socks for the change in the orderbook.
   * @param removeOrderResult
   * @param perpetualMarket
   */
  protected async removeOldOrderFromOrderbook(
    removeOrderResult: RemoveOrderResult,
    perpetualMarket: PerpetualMarketFromDatabase,
    headers: IHeaders,
    sendWebsocketMessage: boolean,
  ): Promise<void> {
    const updatedQuantums: number = await runFuncWithTimingStat(
      this.updatePriceLevelsCache(
        removeOrderResult,
      ),
      this.generateTimingStatsOptions('update_price_level_cache'),
    );
    logger.info({
      at: 'OrderReplaceHandler#removeOldOrderFromOrderbook',
      message: 'Removed old order from orderbook price levels cache',
      updatedQuantums,
    });

    if (sendWebsocketMessage) {
      const orderbookMessage: Message = {
        value: this.createOrderbookWebsocketMessage(
          removeOrderResult.removedOrder!,
          perpetualMarket,
          updatedQuantums,
        ),
        headers,
      };
      sendMessageWrapper(orderbookMessage, KafkaTopics.TO_WEBSOCKETS_ORDERBOOKS);
    }
  }

  /**
   * Update orderbookLevelsCache, and assumes that the order is resting on the book
   * @param removeOrderResult
   * @returns
   */
  // eslint-disable-next-line @typescript-eslint/require-await
  protected async updatePriceLevelsCache(
    removeOrderResult: RemoveOrderResult,
  ): Promise<number> {
    const redisOrder: RedisOrder = removeOrderResult.removedOrder!;
    logger.info({
      at: 'OrderReplaceHandler#updatePriceLevelsCache',
      message: 'Orderbook price level before removing old order',
      orderbookLevel: OrderbookLevelsCache.getOrderbookLevel(
        redisOrder.ticker,
        protocolTranslations.protocolOrderSideToOrderSide(redisOrder.order!.side),
        redisOrder.price,
        redisClient,
      ),
    });
    return OrderbookLevelsCache.updatePriceLevel({
      ticker: redisOrder.ticker,
      side: protocolTranslations.protocolOrderSideToOrderSide(redisOrder.order!.side),
      humanPrice: redisOrder.price,
      sizeDeltaInQuantums: this.getSizeDeltaInQuantums(
        removeOrderResult,
        redisOrder,
      ),
      client: redisClient,
    });
  }

  protected getSizeDeltaInQuantums(
    removeOrderResult: RemoveOrderResult,
    redisOrder: RedisOrder,
  ): string {
    const sizeDelta: Big = Big(
      removeOrderResult.totalFilledQuantums!.toString(),
    ).minus(
      redisOrder.order!.quantums.toString(),
    );

    // TODO(IND-147): This should not be happening once `ender` updates orderbook for filled orders
    // rather than having off-chain updates sent from the protocol. Change to error once it's
    // confirmed this case no longer happens normally.
    if (sizeDelta.gt(0)) {
      logger.info({
        at: 'OrderReplaceHandler#getSizeDeltaInQuantums',
        message: 'Total filled of order exceeds quantums of order',
        totalFilled: removeOrderResult.totalFilledQuantums!.toString(),
        quantums: redisOrder.order!.quantums.toString(),
        removeOrderResult,
        redisOrder,
      });
      return '0';
    }

    return sizeDelta.toFixed();
  }

  /**
   * Adds the removed order to the cancelled orders cache in Redis.
   *
   * @param orderId
   * @param timestampMs
   * @protected
   */
  protected async addOrderToCanceledOrdersCache(
    oldOrderId: IndexerOrderId,
    timestampMs: number,
  ): Promise<void> {
    const orderId: string = OrderTable.orderIdToUuid(oldOrderId);

    await runFuncWithTimingStat(
      CanceledOrdersCache.addCanceledOrderId(orderId, timestampMs, redisClient),
      this.generateTimingStatsOptions('add_order_to_canceled_order_cache'),
    );
  }
}
