import { Callback, RedisClient } from 'redis';

import { zRemAsync, zScoreAsync } from '../helpers/redis';
import { addCanceledOrderIdScript } from './scripts';
// Cache of cancelled orders
export const CANCELED_ORDERS_CACHE_KEY: string = 'v4/cancelled_orders';
// 10 seconds in milliseconds
export const CANCELED_ORDER_WINDOW_SIZE: number = 30 * 1000;

/**
 * isOrderCanceled returns true if the order is canceled.
 *
 * In Vulcan, we cancel orders by adding them to the sorted set (by timestamp).
 * For replace orders, we cancel the old order and place the new order. When placing the new order,
 * we remove the order from the canceled orders cache.
 *
 * In Ender, when handling order fills, we check if the order is in the sorted set.
 * If so, we set the status to be BEST_EFFORT_CANCELED.
 *
 * @param orderId
 * @param client
 */
export async function isOrderCanceled(
  orderId: string,
  client: RedisClient,
): Promise<boolean> {
  const score: string | null = await
  zScoreAsync({ hash: CANCELED_ORDERS_CACHE_KEY, key: orderId }, client);
  return score !== null;
}

export async function removeOrderFromCache(
  orderId: string,
  client: RedisClient,
): Promise<number> {
  return zRemAsync({ hash: CANCELED_ORDERS_CACHE_KEY, key: orderId }, client);
}

/**
 * addCanceledOrderId adds the order id to the canceled orders cache.
 *
 * @param orderId
 * @param timestamp
 * @param client
 */
export async function addCanceledOrderId(
  orderId: string,
  timestamp: number,
  client: RedisClient,
): Promise<number> {
  const numKeys: number = 2;
  let evalAsync: (
    canceledOrderId: string,
    currentTimestampMs: number,
  ) => Promise<number> = (
    canceledOrderId,
    currentTimestampMs,
  ) => {
    return new Promise((resolve, reject) => {
      const callback: Callback<number> = (
        err: Error | null,
        results: number,
      ) => {
        if (err) {
          return reject(err);
        }
        return resolve(results);
      };
      client.evalsha(
        addCanceledOrderIdScript.hash,
        numKeys,
        CANCELED_ORDERS_CACHE_KEY,
        CANCELED_ORDER_WINDOW_SIZE,
        canceledOrderId,
        currentTimestampMs,
        callback,
      );
    });
  };

  evalAsync = evalAsync.bind(client);

  return evalAsync(
    orderId,
    timestamp,
  );
}
