import { logger, runFuncWithTimingStat, stats } from '@dydxprotocol-indexer/base';
import { createSubaccountWebsocketMessage, KafkaTopics } from '@dydxprotocol-indexer/kafka';
import {
  OrderFromDatabase,
  OrderTable,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import {
  CanceledOrdersCache,
  OpenOrdersCache,
  OrderbookLevelsCache,
  placeOrder,
  PlaceOrderResult,
  StatefulOrderUpdatesCache,
  convertToRedisOrder,
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
  OffChainUpdateV1,
  OrderPlaceV1,
  OrderPlaceV1_OrderPlacementStatus,
  OrderUpdateV1,
  RedisOrder,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import { IHeaders, Message } from 'kafkajs';

import config from '../config';
import { redisClient } from '../helpers/redis/redis-controller';
import { sendMessageWrapper } from '../lib/send-message-helper';
import { Handler } from './handler';

/**
 * Handler for OrderPlace messages.
 * The behavior is as follows:
 * - Add the order to the OrdersCache, OrdersDataCache, and SubaccountOrderIdsCache
 *  - this is done using the `placeOrder` function from the `redis` package
 *  - Remove the order from the CanceledOrdersCache if it exists
 * - If the order is a stateful order, attempt to remove any cached order update from the
 *   StatefulOrderUpdatesCache, and then queue the order update to be re-sent and re-processed
 * - If the order doesn't already exist in the caches, return
 * - If the order exists in the caches, but was not replaced due to the expiry of the existing order
 *   being greater than or equal to the expiry of the order in the OrderPlace message, return
 */
export class OrderPlaceHandler extends Handler {
  protected async handle(update: OffChainUpdateV1, headers: IHeaders): Promise<void> {
    logger.info({
      at: 'OrderPlaceHandler#handle',
      message: 'Received OffChainUpdate with OrderPlace.',
      update,
      txHash: this.txHash,
    });
    const orderPlace: OrderPlaceV1 = update.orderPlace!;
    this.validateOrderPlace(update.orderPlace!);
    const order: IndexerOrder = orderPlace.order!;
    const placementStatus: OrderPlaceV1_OrderPlacementStatus = orderPlace.placementStatus;

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(order.orderId!.clobPairId.toString());

    if (perpetualMarket === undefined) {
      this.logAndThrowParseMessageError(
        'Order in OrderPlace has invalid clobPairId',
        {
          order,
        },
      );
      // Needed so TS compiler knows `perpetualMarket` is defined
      return;
    }

    const redisOrder: RedisOrder = convertToRedisOrder(order, perpetualMarket);
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
      at: 'OrderPlaceHandler#handle',
      message: 'OrderPlace processed',
      order,
      placeOrderResult,
    });

    // TODO(IND-68): Remove once order replacement flow in V4 protocol removes the old order and
    // places the updated order.
    const updatedQuantums: number | undefined = await this.updatePriceLevel(
      placeOrderResult,
      perpetualMarket,
      update,
    );

    // TODO(IND-68): Error on this case once replacements are done by first removing the order, then
    // placing a new order.
    if (placeOrderResult.replaced) {
      // Replaced orders are no longer counted as resting on the book until an order update message
      // is received, so remove the order from the set of open orders when replaced.
      const clobPairId: string = order.orderId!.clobPairId.toString();
      await OpenOrdersCache.removeOpenOrder(
        redisOrder.id,
        clobPairId,
        redisClient,
      );
      // TODO(IND-172): Replace this with a logger.error call
      stats.increment(`${config.SERVICE_NAME}.place_order_handler.replaced_order`, 1);
    }

    // TODO(CLOB-597): Remove this logic and log erorrs once best-effort-open is not sent for
    // stateful orders in the protocol
    if (this.shouldSendSubaccountMessage(
      orderPlace,
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
            at: 'OrderPlaceHandler#createSubaccountWebsocketMessage',
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
        ),
        headers: {
          message_received_timestamp: headers.message_received_timestamp,
        },
      };
      sendMessageWrapper(subaccountMessage, KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS);
    }

    // TODO(IND-68): Remove once order replacement flow in V4 protocol removes the old order and
    // places the updated order.
    if (updatedQuantums !== undefined) {
      const orderbookMessage: Message = {
        value: this.createOrderbookWebsocketMessage(
          placeOrderResult.oldOrder!,
          perpetualMarket,
          updatedQuantums,
        ),
        headers: {
          message_received_timestamp: headers.message_received_timestamp,
        },
      };
      sendMessageWrapper(orderbookMessage, KafkaTopics.TO_WEBSOCKETS_ORDERBOOKS);
    }
  }

  /**
   * Updates the price level given the result of calling `placeOrder`.
   * @param result `PlaceOrderResult` from calling `placeOrder`
   * @param perpetualMarket Perpetual market object corresponding to the perpetual market of the
   * order
   * @param update Off-chain update
   * @returns
   */
  // eslint-disable-next-line @typescript-eslint/require-await
  protected async updatePriceLevel(
    result: PlaceOrderResult,
    perpetualMarket: PerpetualMarketFromDatabase,
    update: OffChainUpdateV1,
  ): Promise<number | undefined> {
    // TODO(DEC-1339): Update price levels based on if the order is reduce-only and if the replaced
    // order is reduce-only.
    if (
      result.replaced !== true ||
      result.restingOnBook !== true ||
      requiresImmediateExecution(result.oldOrder!.order!.timeInForce)
    ) {
      return undefined;
    }

    const remainingSizeDeltaInQuantums: Big = this.getRemainingSizeDeltaInQuantums(result);

    if (remainingSizeDeltaInQuantums.eq(0)) {
      // No update to the price level if remaining quantums = 0
      // An order could have remaining quantums = 0 intra-block, as an order is only considered
      // filled once the fills are committed in a block
      return undefined;
    }

    if (remainingSizeDeltaInQuantums.lt(0)) {
      // Log error and skip updating orderbook levels if old order had negative remaining
      // quantums
      logger.info({
        at: 'OrderPlaceHandler#handle',
        message: 'Total filled of order in Redis exceeds order quantums.',
        placeOrderResult: result,
        update,
      });
      stats.increment(`${config.SERVICE_NAME}.order_place_total_filled_exceeds_size`, 1);
      return undefined;
    }

    // If the remaining size is not equal or less than 0, it must be greater than 0.
    // Remove the remaining size of the replaced order from the orderbook, by decrementing
    // the total size of orders at the price of the replaced order
    return runFuncWithTimingStat(
      OrderbookLevelsCache.updatePriceLevel({
        ticker: perpetualMarket.ticker,
        side: protocolTranslations.protocolOrderSideToOrderSide(result.oldOrder!.order!.side),
        humanPrice: result.oldOrder!.price,
        // Delta should be -1 * remaining size of order in quantums and an integer
        sizeDeltaInQuantums: remainingSizeDeltaInQuantums.mul(-1).toFixed(0),
        client: redisClient,
      }),
      this.generateTimingStatsOptions('update_price_level'),
    );
  }

  /**
   * Gets the remaining size of the old order if the order was replaced.
   * @param result Result of placing an order, should be for a replaced order so both `oldOrder` and
   * `oldTotalFilledQuantums` properties should exist on the place order result.
   * @returns Remaining size of the old order that was replaced.
   */
  protected getRemainingSizeDeltaInQuantums(result: PlaceOrderResult): Big {
    const sizeDeltaInQuantums: Big = Big(
      result.oldOrder!.order!.quantums.toString(),
    ).minus(
      result.oldTotalFilledQuantums!,
    );
    return sizeDeltaInQuantums;
  }

  protected validateOrderPlace(orderPlace: OrderPlaceV1): void {
    if (orderPlace.order === undefined) {
      this.logAndThrowParseMessageError('Invalid OrderPlace, order is undefined');
      return;
    }

    if (orderPlace.order.orderId === undefined) {
      this.logAndThrowParseMessageError('Invalid OrderPlace, order id is undefined');
      return;
    }

    if (orderPlace.order.orderId.subaccountId === undefined) {
      this.logAndThrowParseMessageError('Invalid OrderPlace, subaccount id is undefined');
    }

    if (
      orderPlace
        .placementStatus === OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_UNSPECIFIED
    ) {
      this.logAndThrowParseMessageError('Invalid OrderPlace, placement status is UNSPECIFIED');
    }
  }

  /**
   * Determine whether to send a subaccount websocket message given the order place.
   * @param orderPlace
   * @returns TODO(CLOB-597): Remove once best-effort-opened messages are not sent for stateful
   * orders.
   */
  protected shouldSendSubaccountMessage(
    orderPlace: OrderPlaceV1,
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

    const orderFlags: number = orderPlace.order!.orderId!.orderFlags;
    const status: OrderPlaceV1_OrderPlacementStatus = orderPlace.placementStatus;
    // Best-effort-opened status should only be sent for short-term orders
    if (
      orderFlags !== ORDER_FLAG_SHORT_TERM &&
      status === OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED
    ) {
      return false;
    }

    // In the case where a stateful orderPlace is opened with a more recent expiry than an
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
   * @param blockHeight
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
      headers: {
        message_received_timestamp: headers.message_received_timestamp,
        event_type: String(headers.event_type),
      },
    };
    sendMessageWrapper(orderUpdateMessage, KafkaTopics.TO_VULCAN);
  }
}
