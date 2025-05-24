import {
  logger,
  getInstanceId,
  runFuncWithTimingStat,
  stats,
} from '@dydxprotocol-indexer/base';
import { KafkaTopics } from '@dydxprotocol-indexer/kafka';
import {
  PerpetualMarketFromDatabase,
  protocolTranslations,
  perpetualMarketRefresher,
  OrderTable,
} from '@dydxprotocol-indexer/postgres';
import {
  updateOrder,
  UpdateOrderResult,
  OrderbookLevelsCache,
  StatefulOrderUpdatesCache,
} from '@dydxprotocol-indexer/redis';
import { isStatefulOrder, requiresImmediateExecution } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  OffChainUpdateV1,
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
 * Handler for OrderUpdate messages.
 * The behavior is as follows:
 * - Update the total filled quantums of the order in the OrdersDataCache
 * - If the order did not exist in the OrdersDataCache, log an error and return
 * - Update the orderbook levels cache
 *  - if the order was previously not on the book, increase the total size of orders at the price
 *    level of the updated order by the remaining size of the order in quantums (order size in
 *    quantums - total filled quantums for the order)
 *  - if the order was previously on the book, update the total size of orders at the price level
 *    of the updated order by the delta between the old total filled quantums of the order and
 *    the new total filled quantums of the order (in the OrderUpdate message)
 *  NOTE: If the new total filled quantums of the order in the OrderUpdate exceeds the order size
 *  in quantums, an error is logged and the total filled quantums of the order used in calculating
 *  the delta is capped to the size of the order in quantums
 *  NOTE: If the old total filled quantums of the order in the OrderUpdate excceeds the order size
 *  in quantums, the old total filled quantums of the order used to calculate the delta to the
 *  price level is capped to the size of the order in quantums
 */
export class OrderUpdateHandler extends Handler {
  protected async handle(update: OffChainUpdateV1, headers: IHeaders): Promise<void> {
    logger.info({
      at: 'OrderUpdateHandler#handle',
      message: 'Received OffChainUpdate with OrderUpdate.',
      update,
      txHash: this.txHash,
    });

    this.validateOrderUpdate(update.orderUpdate!);
    const orderUpdate: OrderUpdateV1 = update.orderUpdate!;

    const updateResult: UpdateOrderResult = await runFuncWithTimingStat(
      updateOrder({
        updatedOrderId: orderUpdate.orderId!,
        newTotalFilledQuantums: orderUpdate.totalFilledQuantums.toNumber(),
        client: redisClient,
      }),
      this.generateTimingStatsOptions('update_order_cache_update'),
    );

    logger.info({
      at: 'OrderUpdateHandler#handle',
      message: 'OrderUpdate processed',
      orderUpdate,
      updateResult,
    });

    if (updateResult.updated !== true) {
      const orderFlags: number = orderUpdate.orderId!.orderFlags;
      if (isStatefulOrder(orderFlags)) {
        // If the order update was for a stateful order, add it to a cache of order updates
        // for stateful orders, so it can be re-sent after `ender` processes the on-chain
        // event for the stateful order placement
        await StatefulOrderUpdatesCache.addStatefulOrderUpdate(
          OrderTable.orderIdToUuid(orderUpdate.orderId!),
          orderUpdate,
          Date.now(),
          redisClient,
        );
      }
      logger.info({
        at: 'OrderUpdateHandler#handle',
        message: 'Received order update for order that does not exist, order id ' +
                 `${JSON.stringify(orderUpdate.orderId!)}`,
        update,
        updateResult,
      });
      stats.increment(
        `${config.SERVICE_NAME}.order_update_order_does_not_exist`,
        1,
        {
          orderFlags: String(orderFlags),
          instance: getInstanceId(),
        },
      );
      return;
    }

    const sizeDeltaInQuantums: Big = this.getSizeDeltaInQuantums(updateResult, orderUpdate);

    if (sizeDeltaInQuantums.eq(0)) {
      stats.increment(
        `${config.SERVICE_NAME}.order_update_with_zero_delta.count`,
        1,
        { instance: getInstanceId() },
      );
      return;
    }

    // Orders that require immediate execution do not rest on the order book and will not lead to
    // a change in the order book level for the order's price
    if (!requiresImmediateExecution(updateResult.order!.order!.timeInForce)) {
      const updatedQuantums: number = await this.updatePriceLevel(
        updateResult,
        sizeDeltaInQuantums,
      );

      const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
        .getPerpetualMarketFromTicker(updateResult.order!.ticker);
      if (perpetualMarket === undefined) {
        logger.error({
          at: 'OrderUpdateHandler#handle',
          message: `Received order update for order with unknown perpetual market, ticker ${
            updateResult.order!.ticker}`,
        });
        return;
      }

      const orderbookMessage: Message = {
        value: this.createOrderbookWebsocketMessage(
          updateResult.order!,
          perpetualMarket,
          updatedQuantums,
        ),
        headers,
      };
      sendMessageWrapper(orderbookMessage, KafkaTopics.TO_WEBSOCKETS_ORDERBOOKS);
    }
  }

  /**
   * Updates the price level given the result of updating an order. The `UpdateResult` passed into
   * this function should be from a successful order update.
   * @param updateResult Result of updating the orders caches.
   * @param sizeDeltaInQuantums Size delta in quantums to update the price level by.
   * @returns
   */
  // eslint-disable-next-line @typescript-eslint/require-await
  protected async updatePriceLevel(
    updateResult: UpdateOrderResult,
    sizeDeltaInQuantums: Big,
  ): Promise<number> {
    const redisOrder: RedisOrder = updateResult.order!;

    return runFuncWithTimingStat(
      OrderbookLevelsCache.updatePriceLevel(
        redisOrder.ticker,
        protocolTranslations.protocolOrderSideToOrderSide(redisOrder.order!.side),
        redisOrder.price,
        sizeDeltaInQuantums.toFixed(0),
        redisClient,
      ),
      this.generateTimingStatsOptions('update_price_level'),
    );
  }

  /**
   * Gets the old total filled amount of an order from the result of updating the order. If the
   * old total filled amount exceeds the size of the order, return the size of the order instead.
   * @param updateResult Result from updating the order.
   * @returns
   */
  protected getCappedOldTotalFilledQuantums(
    updateResult: UpdateOrderResult,
  ): Big {
    if (updateResult.oldTotalFilledQuantums! <= updateResult.order!.order!.quantums.toNumber()) {
      return Big(updateResult.oldTotalFilledQuantums!);
    }

    // Cap old total filled quantums for an order to the size of the order in quantums, log an error
    // if the old total filled quantums for an order exceeds it's size in quantums
    logger.info({
      at: 'OrderUpdateHandler#getCappedOldTotalFilledQuantums',
      message: 'Old total filled quantums of order exceeds order size in quantums.',
      updateResult,
    });
    stats.increment(
      `${config.SERVICE_NAME}.order_update_old_total_filled_exceeds_size`,
      1,
      { instance: getInstanceId() },
    );
    return Big(updateResult.order!.order!.quantums.toNumber().toString());
  }

  /**
   * Gets the new total filled amount of an order. If the new total filled amount exceeds the size
   * of the order, return the size of the order instead.
   * @param orderUpdate
   * @param updateResult
   * @returns
   */
  protected getCappedNewTotalFilledQuantums(
    orderUpdate: OrderUpdateV1,
    updateResult: UpdateOrderResult,
  ): Big {
    if (orderUpdate.totalFilledQuantums.lte(updateResult.order!.order!.quantums)) {
      return Big(orderUpdate.totalFilledQuantums!.toNumber().toString());
    }

    // Cap new total filled quantums for an order to the size of the order in quantums, log an error
    // if the new total filled quantums for an order exceeds it's size in quantums
    logger.info({
      at: 'OrderUpdateHandler#getCappedNewTotalFilledQuantums',
      message: 'New total filled quantums of order exceeds order size in quantums.',
      orderUpdate,
      updateResult,
    });
    stats.increment(
      `${config.SERVICE_NAME}.order_update_total_filled_exceeds_size`,
      1,
      { instance: getInstanceId() },
    );

    return Big(updateResult.order!.order!.quantums.toNumber().toString());
  }

  /**
   * Get the delta to the price level for an order update.
   * If the updated order was previously not resting on the book
   * - delta should be increasing the total size at the price level of the order by the remaining
   *   size of the order (order size - new total filled quantums)
   * If the updated order was previously resting on the book
   * - delta should be changing the total size at the price level of the order by the difference
   *   between the old total filled quantums of the order and the new total filled quantums
   * @param updateResult Result from updating the order
   * @param orderUpdate Order update message
   * @returns Size delta in quantums to the update the price level
   */
  protected getSizeDeltaInQuantums(
    updateResult: UpdateOrderResult,
    orderUpdate: OrderUpdateV1,
  ): Big {
    const cappedNewTotalFilledQuantums: Big = this.getCappedNewTotalFilledQuantums(
      orderUpdate,
      updateResult,
    );
    const cappedOldTotalFilledQuantums: Big = this.getCappedOldTotalFilledQuantums(updateResult);
    const orderQuantums: Big = Big(updateResult.order!.order!.quantums.toNumber().toString());

    // If the updated order was not resting on the book before being updated, the size delta for
    // the price level update should be the remaining size of the order
    if (updateResult.oldRestingOnBook === false) {
      // Size delta should be the remaining size of the order in quantums
      // Size of order in quantums - total filled in quantums
      return orderQuantums.minus(cappedNewTotalFilledQuantums);
    }

    // If the updated order was resting on the book before being updated, the size detla for the
    // price level update should be the difference between the old total filled of the order and
    // the new total filled of the order
    return cappedOldTotalFilledQuantums.minus(cappedNewTotalFilledQuantums);
  }

  protected validateOrderUpdate(orderUpdate: OrderUpdateV1): void {
    if (orderUpdate.orderId === undefined) {
      this.logAndThrowParseMessageError('Invalid OrderUpdate, order id is undefined');
      return;
    }

    if (orderUpdate.orderId.subaccountId === undefined) {
      this.logAndThrowParseMessageError('Invalid OrderUpdate, subaccount id is undefined');
    }
  }
}
