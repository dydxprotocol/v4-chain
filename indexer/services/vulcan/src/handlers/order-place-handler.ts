import {
  logger, getInstanceId, runFuncWithTimingStat, stats,
} from '@dydxprotocol-indexer/base';
import { createSubaccountWebsocketMessage, KafkaTopics } from '@dydxprotocol-indexer/kafka';
import {
  blockHeightRefresher,
  OrderFromDatabase,
  OrderTable,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import {
  CanceledOrdersCache,
  convertToRedisOrder,
  placeOrder,
  PlaceOrderResult,
  StatefulOrderUpdatesCache,
} from '@dydxprotocol-indexer/redis';
import { getOrderIdHash, isStatefulOrder, ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerSubaccountId,
  OffChainUpdateV1,
  OrderPlaceV1,
  OrderPlaceV1_OrderPlacementStatus,
  OrderUpdateV1,
  RedisOrder,
} from '@dydxprotocol-indexer/v4-protos';
import { IHeaders, Message } from 'kafkajs';

import config from '../config';
import { isVaultOrder } from '../helpers/orders';
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

    if (placeOrderResult.replaced) {
      stats.increment(
        `${config.SERVICE_NAME}.place_order_handler.replaced_order`,
        1,
        { instance: getInstanceId() },
      );
    }

    // TODO(CLOB-597): Remove this logic and log erorrs once best-effort-open is not sent for
    // stateful orders in the protocol
    if (this.shouldSendSubaccountMessage(
      orderPlace,
      placeOrderResult,
      placementStatus,
    )) {
      // TODO(IND-171): Determine whether we should always be sending a message, even when the cache
      // isn't updated.
      // For stateful and conditional orders, look the order up in the db for the createdAtHeight
      // and send any cached order updates for the stateful or conditional order
      let dbOrder: OrderFromDatabase | undefined;
      if (isStatefulOrder(redisOrder.order!.orderId!.orderFlags)) {
        const orderUuid: string = OrderTable.orderIdToUuid(redisOrder.order!.orderId!);
        // Since vault orders are not persisted by ender (to improve processing latency), skip
        // looking them up in db. However, we should still send corresponding cached order update.
        if (!isVaultOrder(redisOrder.order!.orderId!)) {
          dbOrder = await OrderTable.findById(orderUuid);
          if (dbOrder === undefined) {
            logger.crit({
              at: 'OrderPlaceHandler#createSubaccountWebsocketMessage',
              message: 'Stateful order not found in database',
            });
            throw new Error(`Stateful order not found in database: ${orderUuid}`);
          }
        }
        await this.sendCachedOrderUpdate(orderUuid, headers);
      }
      const subaccountMessage: Message = {
        key: Buffer.from(
          IndexerSubaccountId.encode(redisOrder.order!.orderId!.subaccountId!).finish(),
        ),
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
  ): boolean {
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
        OffChainUpdateV1.encode({ orderUpdate: cachedOrderUpdate }).finish(),
      ),
      headers,
    };
    sendMessageWrapper(orderUpdateMessage, KafkaTopics.TO_VULCAN);
  }
}
