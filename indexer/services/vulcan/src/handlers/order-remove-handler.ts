import {
  logger, getInstanceId, runFuncWithTimingStat, stats,
} from '@dydxprotocol-indexer/base';
import { KafkaTopics, SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION, getTriggerPrice } from '@dydxprotocol-indexer/kafka';
import {
  blockHeightRefresher,
  BlockTable,
  BlockFromDatabase,
  OrderFromDatabase,
  OrderStatus,
  OrderTable,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  protocolTranslations,
  SubaccountMessageContents,
  SubaccountTable,
  apiTranslations,
  TimeInForce,
  IsoString,
  OrderSide,
  APITimeInForce,
  OrderType,
} from '@dydxprotocol-indexer/postgres';
import {
  OrderbookLevelsCache,
  OrdersCache,
  RemoveOrderResult,
  removeOrder,
  CanceledOrdersCache,
} from '@dydxprotocol-indexer/redis';
import {
  ORDER_FLAG_SHORT_TERM,
  isStatefulOrder,
  requiresImmediateExecution,
} from '@dydxprotocol-indexer/v4-proto-parser';
import {
  OffChainUpdateV1,
  IndexerOrder,
  IndexerSubaccountId,
  OrderRemoveV1,
  OrderRemovalReason,
  OrderRemoveV1_OrderRemovalStatus,
  RedisOrder,
  SubaccountMessage,
} from '@dydxprotocol-indexer/v4-protos';
import { Big } from 'big.js';
import { IHeaders, Message } from 'kafkajs';

import config from '../config';
import { isVaultOrder } from '../helpers/orders';
import { redisClient } from '../helpers/redis/redis-controller';
import { sendMessageWrapper } from '../lib/send-message-helper';
import { Handler } from './handler';
import { getStateRemainingQuantums } from './helpers';

/**
 * Handler for OrderRemove messages.
 * The behavior is as follows:
 * - Cancel the order in redis from the SubaccountOrderIdsCache, OrdersCache, and OrderDataCache.
 * - Add the order id to the CanceledOrdersCache. This is used to properly set the status of the
 *   order in Postgres to BEST_EFFORT_CANCELED when a fill is received for the canceled order.
 * - If the order is a stateful cancelation indicated by the reason USER_CANCELED, status being
 *   CANCELED and the order being a LONG_TERM order (should only be sent from `ender`)
 *   - send a message in socks that the order was removed with the status CANCELED, using the order
 *     in Postgres to populate the message
 *   - if the order existed in redis and was resting on the book
 *     - update the orderbookLevels cache, reducing the size of the price level of the order by
 *       orderSize - totalFilled
 * - If the order is not a stateful cancelation
 *   - If the order did not exist in redis, ignore the removal
 *   - Update the status of the order in postgres to the removal status of the message
 *   - If the order being removed was on the book, update the orderbookLevelsCache
 *     - orderbookSize -= orderSize - totalFilled
 *   - Send a message to socks that an order was removed along with reason and order_status
 */
export class OrderRemoveHandler extends Handler {
  // eslint-disable-next-line @typescript-eslint/require-await
  protected async handle(update: OffChainUpdateV1, headers: IHeaders): Promise<void> {
    logger.info({
      at: 'OrderRemoveHandler#handle',
      message: 'Received OffChainUpdate with OrderRemove.',
      update,
      txHash: this.txHash,
    });
    const orderRemove: OrderRemoveV1 = update.orderRemove!;
    const reason: OrderRemovalReason = orderRemove.reason;

    this.validateOrderRemove(orderRemove);

    // If the Indexer sent this expire message, check to verify it's still relevant. Updates may
    // have come in between its send and receipt.
    if (
      reason === OrderRemovalReason.ORDER_REMOVAL_REASON_INDEXER_EXPIRED &&
      !(await this.isOrderExpired(orderRemove))
    ) {
      stats.increment(
        `${config.SERVICE_NAME}.order_remove_reason_indexer_temp_expired`,
        1,
        { instance: getInstanceId() },
      );
      logger.info({
        at: 'OrderRemoveHandler#handle',
        message: 'Order was expired by Indexer but is no longer expired. Ignoring.',
        orderRemove,
      });
      return;
    }

    const removeOrderResult: RemoveOrderResult = await runFuncWithTimingStat(
      removeOrder({
        removedOrderId: orderRemove.removedOrderId!,
        client: redisClient,
      }),
      this.generateTimingStatsOptions('remove_order'),
    );
    logger.info({
      at: 'OrderRemoveHandler#handle',
      message: 'OrderRemove processed',
      orderRemove,
      removeOrderResult,
    });

    if (
      orderRemove.reason === OrderRemovalReason.ORDER_REMOVAL_REASON_INDEXER_EXPIRED
    ) {
      stats.increment(
        `${config.SERVICE_NAME}.order_remove_reason_indexer_expired`,
        1,
        { instance: getInstanceId() },
      );
      logger.info({
        at: 'OrderRemoveHandler#handle',
        message: 'Order was expired by Indexer',
        orderRemove,
        removeOrderResult,
      });
    }

    // Since vault orders are not persisted by ender (to improve processing latency), we
    // should handle their cancellation as non-stateful.
    if (this.isStatefulOrderCancelation(orderRemove) &&
      !isVaultOrder(orderRemove.removedOrderId!)) {
      await this.handleStatefulOrderCancelation(orderRemove, removeOrderResult, headers);
      return;
    }

    await this.handleOrderRemoval(orderRemove, removeOrderResult, headers);
  }

  protected validateOrderRemove(orderRemove: OrderRemoveV1): void {
    if (orderRemove.removedOrderId === undefined) {
      return this.logAndThrowParseMessageError(
        'OrderRemove must contain a removedOrderId',
        { orderRemove },
      );
    }

    if (orderRemove.removedOrderId.subaccountId === undefined) {
      return this.logAndThrowParseMessageError(
        'OrderRemove must contain a removedOrderId.subaccountId',
        { orderRemove },
      );
    }

    if (orderRemove.removedOrderId.clientId === undefined) {
      return this.logAndThrowParseMessageError(
        'OrderRemove must contain a removedOrderId.clientId',
        { orderRemove },
      );
    }

    if (orderRemove.removalStatus ===
      OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_UNSPECIFIED) {
      return this.logAndThrowParseMessageError(
        'OrderRemove removalStatus cannot be unspecified',
        { orderRemove },
      );
    }

    if (orderRemove.reason === OrderRemovalReason.ORDER_REMOVAL_REASON_UNSPECIFIED) {
      return this.logAndThrowParseMessageError(
        'OrderRemove reason cannot be unspecified',
        { orderRemove },
      );
    }
  }

  /**
   * Handles an order removal that is a stateful cancelation
   * - sends a message to socks indicating the order has status CANCELED and reason USER_CANCELED
   *   with the details from the Postgres order
   * - updates the orderbook if the removed order was in Redis and was resting on the book
   * Note: It's possible for there to be a race condition where another stateful order with the same
   * id is placed after the cancelation, in which case the message sent to socks will have the
   * incorrect details. This is acceptable as a user cannot be certain an order was canceled until
   * receiving the status CANCELED for the order, and so re-placing the order without receiving the
   * CANCELED message can lead to invalid responses.
   * @param orderRemove
   * @param removeOrderResult
   * @returns
   */
  protected async handleStatefulOrderCancelation(
    orderRemove: OrderRemoveV1,
    removeOrderResult: RemoveOrderResult,
    headers: IHeaders,
  ): Promise<void> {
    const order: OrderFromDatabase | undefined = await runFuncWithTimingStat(
      OrderTable.findById(
        OrderTable.orderIdToUuid(orderRemove.removedOrderId!),
      ),
      this.generateTimingStatsOptions('find_order_for_stateful_cancelation'),
    );
    if (order === undefined) {
      logger.error({
        at: 'orderRemoveHandler#handleStatefulOrderCancelation',
        message: 'Could not find order for stateful order cancelation',
        orderId: orderRemove.removedOrderId,
        orderRemove,
      });
      return;
    }

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(
        order.clobPairId,
      );
    if (perpetualMarket === undefined) {
      logger.error({
        at: 'orderRemoveHandler#handleStatefulOrderCancelation',
        message: `Unable to find the perpetual market with clobPairId: ${order.clobPairId}`,
        order,
        orderRemove,
      });
      return;
    }

    const subaccountMessage: Message = {
      key: Buffer.from(
        IndexerSubaccountId.encode(orderRemove.removedOrderId!.subaccountId!).finish(),
      ),
      value: this.createSubaccountWebsocketMessageFromPostgresOrder(
        order,
        orderRemove,
        perpetualMarket.ticker,
        blockHeightRefresher.getLatestBlockHeight(),
      ),
      headers,
    };
    sendMessageWrapper(subaccountMessage, KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS);

    // If an order was removed from the Orders cache and was resting on the book, update the
    // orderbook levels cache
    // Orders that require immediate execution do not rest on the book, and also should not lead
    // to an update to the orderbook levels cache
    if (
      removeOrderResult.removed &&
      removeOrderResult.restingOnBook === true &&
      !requiresImmediateExecution(removeOrderResult.removedOrder!.order!.timeInForce)) {
      await this.updateOrderbook(removeOrderResult, perpetualMarket, headers);
    }

  }

  /**
   * Handles an order removal that is not a stateful cancelation.
   * - if an order was not removed from redis, ignore the removal
   * - send a message to the subaccount indicating the order was removed
   * - update the status of the order in Postgres
   * - update the orderbook if the order was resting on the book
   * @param orderRemove
   * @param removeOrderResult
   * @returns
   */
  protected async handleOrderRemoval(
    orderRemove: OrderRemoveV1,
    removeOrderResult: RemoveOrderResult,
    headers: IHeaders,
  ): Promise<void> {
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(orderRemove.removedOrderId!.clobPairId.toString());
    if (perpetualMarket === undefined) {
      const clobPairId: string = orderRemove.removedOrderId!.clobPairId.toString();
      logger.error({
        at: 'orderRemoveHandler#handle',
        message: `Unable to find perpetual market with clobPairId: ${clobPairId}`,
      });
      return;
    }
    // This can happen for short term orders if the order place message was not received.
    if (!removeOrderResult.removed) {
      logger.info({
        at: 'orderRemoveHandler#handleOrderRemoval',
        message: 'Unable to find order',
        orderId: orderRemove.removedOrderId,
        orderRemove,
      });
      if (config.SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_CANCELS_MISSING_ORDERS) {
        const canceledOrder: OrderFromDatabase | undefined = await runFuncWithTimingStat(
          OrderTable.findById(OrderTable.orderIdToUuid(orderRemove.removedOrderId!)),
          this.generateTimingStatsOptions('find_order'),
        );
        const subaccountMessage: Message = {
          key: Buffer.from(
            IndexerSubaccountId.encode(orderRemove.removedOrderId!.subaccountId!).finish(),
          ),
          value: this.createSubaccountWebsocketMessageFromOrderRemoveMessage(
            canceledOrder,
            orderRemove,
            perpetualMarket.ticker,
            blockHeightRefresher.getLatestBlockHeight(),
          ),
          headers,
        };
        const reason: OrderRemovalReason = orderRemove.reason;
        if (!(
          reason === OrderRemovalReason.ORDER_REMOVAL_REASON_INDEXER_EXPIRED ||
          reason === OrderRemovalReason.ORDER_REMOVAL_REASON_FULLY_FILLED
        )) {
          sendMessageWrapper(subaccountMessage, KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS);
        }
      }
      return;
    }

    const stateRemainingQuantums: Big = await getStateRemainingQuantums(
      removeOrderResult.removedOrder!,
    );

    // If the remaining amount of the order in state is <= 0, the order is filled and
    // does not need to have it's status updated
    let canceledOrder: OrderFromDatabase | undefined;
    if (stateRemainingQuantums.gt(0)) {
      canceledOrder = await runFuncWithTimingStat(
        this.cancelOrderInPostgres(orderRemove),
        this.generateTimingStatsOptions('cancel_order_in_postgres'),
      );
    } else {
      canceledOrder = await runFuncWithTimingStat(
        OrderTable.findById(OrderTable.orderIdToUuid(orderRemove.removedOrderId!)),
        this.generateTimingStatsOptions('find_order'),
      );
    }

    const subaccountMessage: Message = {
      key: Buffer.from(
        IndexerSubaccountId.encode(orderRemove.removedOrderId!.subaccountId!).finish(),
      ),
      value: this.createSubaccountWebsocketMessageFromRemoveOrderResult(
        removeOrderResult,
        canceledOrder,
        orderRemove,
        perpetualMarket,
        blockHeightRefresher.getLatestBlockHeight(),
      ),
      headers,
    };

    if (this.shouldSendSubaccountMessage(orderRemove, removeOrderResult, stateRemainingQuantums)) {
      sendMessageWrapper(subaccountMessage, KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS);
    }

    const remainingQuantums: Big = Big(this.getSizeDeltaInQuantums(
      removeOrderResult,
      removeOrderResult.removedOrder!,
    ));
    // Do not update orderbook if order being cancelled has no remaining quantums or is
    // resting on book, or requires immediate execution and will not rest on the book
    if (
      !remainingQuantums.eq('0') &&
      removeOrderResult.restingOnBook !== false &&
      !requiresImmediateExecution(removeOrderResult.removedOrder!.order!.timeInForce)) {
      await this.updateOrderbook(removeOrderResult, perpetualMarket, headers);
    }
    // TODO: consolidate remove handler logic into a single lua script.
    await this.addOrderToCanceledOrdersCache(
      orderRemove,
      Date.now(),
    );
  }

  /**
   * Updates the orderbook and sends a message to socks for the change in the orderbook.
   * @param removeOrderResult
   * @param perpetualMarket
   */
  protected async updateOrderbook(
    removeOrderResult: RemoveOrderResult,
    perpetualMarket: PerpetualMarketFromDatabase,
    headers: IHeaders,
  ): Promise<void> {
    const updatedQuantums: number = await runFuncWithTimingStat(
      this.updatePriceLevelsCache(
        removeOrderResult,
      ),
      this.generateTimingStatsOptions('update_price_level_cache'),
    );
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

  /**
   * Adds the removed order to the cancelled orders cache in Redis.
   *
   * @param orderId
   * @param blockHeight
   * @protected
   */
  protected async addOrderToCanceledOrdersCache(
    orderRemove: OrderRemoveV1,
    timestampMs: number,
  ): Promise<void> {
    const orderId: string = OrderTable.orderIdToUuid(orderRemove.removedOrderId!);

    if (
      orderRemove.removalStatus ===
      OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED
    ) {
      await runFuncWithTimingStat(
        CanceledOrdersCache.addBestEffortCanceledOrderId(orderId, timestampMs, redisClient),
        this.generateTimingStatsOptions('add_order_to_canceled_order_cache'),
      );
    } else if (
      orderRemove.removalStatus ===
      OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED
    ) {
      await runFuncWithTimingStat(
        CanceledOrdersCache.addCanceledOrderId(orderId, timestampMs, redisClient),
        this.generateTimingStatsOptions('add_order_to_canceled_order_cache'),
      );
    }
  }

  /**
   * When the Indexer sends an expire message, we want to verify that the order hasn't received an
   * update since it occurred which would invalidate the message.
   */
  protected async isOrderExpired(orderRemove: OrderRemoveV1): Promise<boolean> {
    let block: BlockFromDatabase;
    try {
      block = await runFuncWithTimingStat(
        BlockTable.getLatest({ readReplica: true }),
        this.generateTimingStatsOptions('get_latest_block_for_indexer_expired_expiry_verification'),
      );
    } catch {
      logger.error({
        at: 'orderRemoveHandler#isOrderExpired',
        message: 'Unable to find latest block',
        orderRemove,
      });
      // We can't say with certainty that this order is expired, so return false
      return false;
    }

    const redisOrder: RedisOrder | null = await runFuncWithTimingStat(
      OrdersCache.getOrder(OrderTable.orderIdToUuid(orderRemove.removedOrderId!), redisClient),
      this.generateTimingStatsOptions('find_order_for_indexer_expired_expiry_verification'),
    );
    if (redisOrder === null) {
      stats.increment(
        `${config.SERVICE_NAME}.indexer_expired_order_not_found`,
        1,
        { instance: getInstanceId() },
      );
      logger.info({
        at: 'orderRemoveHandler#isOrderExpired',
        message: 'Could not find order for Indexer-expired expiry verification',
        orderRemove,
      });
      // We can't say with certainty that this order is expired, if it still exists, so return false
      return false;
    }
    const order: IndexerOrder = redisOrder.order!;

    // Indexer should only ever send expiration messages for short-term orders
    if (order.orderId!.orderFlags !== ORDER_FLAG_SHORT_TERM) {
      logger.error({
        at: 'orderRemoveHandler#isOrderExpired',
        message: 'Long-term order retrieved during Indexer-expired expiry verification',
        orderRemove,
        redisOrder,
      });
      return false;
    }

    // We know the order is short-term, so the goodTilBlock must exist.
    if (order.goodTilBlock! >= +block.blockHeight) {
      stats.increment(
        `${config.SERVICE_NAME}.indexer_expired_order_is_not_expired`,
        1,
        { instance: getInstanceId() },
      );
      logger.info({
        at: 'orderRemoveHandler#isOrderExpired',
        message: 'Indexer marked order that is not yet expired as expired',
        orderRemove,
        redisOrder,
        block,
      });
      return false;
    }
    return true;
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  protected async cancelOrderInPostgres(
    orderRemove: OrderRemoveV1,
  ): Promise<OrderFromDatabase | undefined> {
    return OrderTable.update({
      id: OrderTable.orderIdToUuid(orderRemove.removedOrderId!),
      status: this.orderRemovalStatusToOrderStatus(orderRemove.removalStatus),
    });
  }

  protected orderRemovalStatusToOrderStatus(
    orderRemovalStatus: OrderRemoveV1_OrderRemovalStatus,
  ): OrderStatus {
    switch (orderRemovalStatus) {
      case OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED:
        return OrderStatus.CANCELED;
      case OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_FILLED:
        return OrderStatus.FILLED;
      case OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED:
      default:
        return OrderStatus.BEST_EFFORT_CANCELED;
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
    return OrderbookLevelsCache.updatePriceLevel(
      redisOrder.ticker,
      protocolTranslations.protocolOrderSideToOrderSide(redisOrder.order!.side),
      redisOrder.price,
      this.getSizeDeltaInQuantums(
        removeOrderResult,
        redisOrder,
      ),
      redisClient,
    );
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
        at: 'orderRemoveHandler#getSizeDeltaInQuantums',
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
   * Should be called when an OrderRemove message is received for a non-existent order.
   * This can happen when the order was not found in redis because the initial order
   * placement message wasn't received.
   *
   * @param canceledOrder
   * @param orderRemove
   * @param perpetualMarket
   * @param blockHeight: latest block height processed by Indexer
   * @protected
   */
  protected createSubaccountWebsocketMessageFromOrderRemoveMessage(
    canceledOrder: OrderFromDatabase | undefined,
    orderRemove: OrderRemoveV1,
    ticker: string,
    blockHeight: string,
  ): Buffer {
    const createdAtHeight: string | undefined = canceledOrder?.createdAtHeight;
    const updatedAt: IsoString | undefined = canceledOrder?.updatedAt;
    const updatedAtHeight: string | undefined = canceledOrder?.updatedAtHeight;
    const price: string | undefined = canceledOrder?.price;
    const size: string | undefined = canceledOrder?.size;
    const clientMetadata: string | undefined = canceledOrder?.clientMetadata;
    const reduceOnly: boolean | undefined = canceledOrder?.reduceOnly;
    const side: OrderSide | undefined = canceledOrder?.side;
    const timeInForce: APITimeInForce | undefined = canceledOrder
      ? apiTranslations.orderTIFToAPITIF(canceledOrder.timeInForce) : undefined;
    const totalFilled: string | undefined = canceledOrder?.totalFilled;
    const goodTilBlock: string | undefined = canceledOrder?.goodTilBlock;
    const goodTilBlockTime: string | undefined = canceledOrder?.goodTilBlockTime;
    const triggerPrice: string | undefined = canceledOrder?.triggerPrice;
    const type: OrderType | undefined = canceledOrder?.type;

    const contents: SubaccountMessageContents = {
      orders: [
        {
          id: OrderTable.orderIdToUuid(orderRemove.removedOrderId!),
          subaccountId: SubaccountTable.subaccountIdToUuid(
            orderRemove.removedOrderId!.subaccountId!,
          ),
          clientId: orderRemove.removedOrderId!.clientId.toString(),
          clobPairId: orderRemove.removedOrderId!.clobPairId.toString(),
          status: this.orderRemovalStatusToOrderStatus(orderRemove.removalStatus),
          orderFlags: orderRemove.removedOrderId!.orderFlags.toString(),
          ticker,
          removalReason: OrderRemovalReason[orderRemove.reason],
          ...(createdAtHeight && { createdAtHeight }),
          ...(updatedAt && { updatedAt }),
          ...(updatedAtHeight && { updatedAtHeight }),
          ...(price && { price }),
          ...(size && { size }),
          ...(clientMetadata && { clientMetadata }),
          ...(reduceOnly && { reduceOnly }),
          ...(side && { side }),
          ...(timeInForce && { timeInForce }),
          ...(totalFilled && { totalFilled }),
          ...(goodTilBlock && { goodTilBlock }),
          ...(goodTilBlockTime && { goodTilBlockTime }),
          ...(triggerPrice && { triggerPrice }),
          ...(type && { type }),
        },
      ],
      blockHeight,
    };

    const subaccountMessage: SubaccountMessage = SubaccountMessage.fromPartial({
      contents: JSON.stringify(contents),
      subaccountId: orderRemove.removedOrderId!.subaccountId!,
      version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
    });

    return Buffer.from(SubaccountMessage.encode(subaccountMessage).finish());
  }

  protected createSubaccountWebsocketMessageFromRemoveOrderResult(
    removeOrderResult: RemoveOrderResult,
    canceledOrder: OrderFromDatabase | undefined,
    orderRemove: OrderRemoveV1,
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
            orderRemove.removedOrderId!.subaccountId!,
          ),
          clientId: orderRemove.removedOrderId!.clientId.toString(),
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
          status: this.orderRemovalStatusToOrderStatus(orderRemove.removalStatus),
          timeInForce: apiTranslations.orderTIFToAPITIF(orderTIF),
          postOnly: apiTranslations.isOrderTIFPostOnly(orderTIF),
          reduceOnly: redisOrder.order!.reduceOnly,
          orderFlags: redisOrder.order!.orderId!.orderFlags.toString(),
          goodTilBlock: protocolTranslations.getGoodTilBlock(redisOrder.order!)
            ?.toString(),
          goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(redisOrder.order!),
          ticker: redisOrder.ticker,
          removalReason: OrderRemovalReason[orderRemove.reason],
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
      subaccountId: orderRemove.removedOrderId!.subaccountId!,
      version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
    });

    return Buffer.from(SubaccountMessage.encode(subaccountMessage).finish());
  }

  protected createSubaccountWebsocketMessageFromPostgresOrder(
    order: OrderFromDatabase,
    orderRemove: OrderRemoveV1,
    orderTicker: string,
    blockHeight: string | undefined,
  ): Buffer {
    const contents: SubaccountMessageContents = {
      orders: [
        {
          id: order.id,
          subaccountId: SubaccountTable.subaccountIdToUuid(
            orderRemove.removedOrderId!.subaccountId!,
          ),
          clientId: orderRemove.removedOrderId!.clientId.toString(),
          clobPairId: order.clobPairId,
          side: order.side,
          size: order.size,
          totalFilled: order.totalFilled,
          price: order.price,
          type: order.type,
          status: this.orderRemovalStatusToOrderStatus(orderRemove.removalStatus),
          timeInForce: apiTranslations.orderTIFToAPITIF(order.timeInForce),
          postOnly: apiTranslations.isOrderTIFPostOnly(order.timeInForce),
          reduceOnly: order.reduceOnly,
          orderFlags: order.orderFlags,
          goodTilBlock: order.goodTilBlock ?? undefined,
          goodTilBlockTime: order.goodTilBlockTime ?? undefined,
          ticker: orderTicker,
          removalReason: OrderRemovalReason[orderRemove.reason],
          createdAtHeight: order.createdAtHeight,
          updatedAt: order.updatedAt,
          updatedAtHeight: order.updatedAtHeight,
          clientMetadata: order.clientMetadata,
          triggerPrice: order.triggerPrice ?? undefined,
        },
      ],
      ...(blockHeight && { blockHeight }),
    };

    const subaccountMessage: SubaccountMessage = SubaccountMessage.fromPartial({
      contents: JSON.stringify(contents),
      subaccountId: orderRemove.removedOrderId!.subaccountId!,
      version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
    });

    return Buffer.from(SubaccountMessage.encode(subaccountMessage).finish());
  }

  protected isStatefulOrderCancelation(
    orderRemove: OrderRemoveV1,
  ): boolean {
    return (
      isStatefulOrder(orderRemove.removedOrderId!.orderFlags) &&
      orderRemove.reason === OrderRemovalReason.ORDER_REMOVAL_REASON_USER_CANCELED &&
      orderRemove.removalStatus === OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED
    );
  }

  /**
   * Determine if a subaccount message should be sent for an order removal. Do not send messages if:
   * - best effort cancelling orders that are optimistically fully filled for user cancelations
   * - indexer expired cancelations
   * - orders that are removed due to being fully filled
   * @param orderRemove
   * @param removeOrderResult
   * @param redisOrder
   * @returns TODO(IND-147): Remove this logic once we remove orders from redis when filled in ender
   */
  protected shouldSendSubaccountMessage(
    orderRemove: OrderRemoveV1,
    removeOrderResult: RemoveOrderResult,
    stateRemainingQuantums: Big,
  ): boolean {
    const status: OrderRemoveV1_OrderRemovalStatus = orderRemove.removalStatus;
    const reason: OrderRemovalReason = orderRemove.reason;

    logger.info({
      at: 'orderRemoveHandler#shouldSendSubaccountMessage',
      message: 'Compared state filled quantums and size',
      stateRemainingQuantums: stateRemainingQuantums.toFixed(),
      removeOrderResult,
    });

    if (
      stateRemainingQuantums.lte(0) &&
      status === OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED &&
      reason === OrderRemovalReason.ORDER_REMOVAL_REASON_USER_CANCELED
    ) {
      return false;
    } else if (
      stateRemainingQuantums.lte(0) &&
      status === OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED &&
      reason === OrderRemovalReason.ORDER_REMOVAL_REASON_INDEXER_EXPIRED
    ) {
      return false;
    } else if (reason === OrderRemovalReason.ORDER_REMOVAL_REASON_FULLY_FILLED) {
      return false;
    }
    return true;
  }
}
