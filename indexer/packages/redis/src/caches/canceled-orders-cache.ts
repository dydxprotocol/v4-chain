import { Callback, RedisClient } from 'redis';

import { zRemAsync, zScoreAsync } from '../helpers/redis';
import { CanceledOrderStatus } from '../types';
import { addCanceledOrderIdScript } from './scripts';
// Cache of cancelled orders
export const CANCELED_ORDERS_CACHE_KEY: string = 'v4/cancelled_orders';
export const BEST_EFFORT_CANCELED_ORDERS_CACHE_KEY: string = 'v4/best_effort_cancelled_orders';
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
  const [
    canceledScore,
    bestEffortCanceledScore,
  ]: (string | null)[] = await Promise.all([
    zScoreAsync({ hash: CANCELED_ORDERS_CACHE_KEY, key: orderId }, client),
    zScoreAsync({ hash: BEST_EFFORT_CANCELED_ORDERS_CACHE_KEY, key: orderId }, client),
  ]);
  return canceledScore !== null || bestEffortCanceledScore !== null;
}

export async function getOrderCanceledStatus(
  orderId: string,
  client: RedisClient,
): Promise<CanceledOrderStatus> {
  const [
    canceledScore,
    bestEffortCanceledScore,
  ]: (string | null)[] = await Promise.all([
    zScoreAsync({ hash: CANCELED_ORDERS_CACHE_KEY, key: orderId }, client),
    zScoreAsync({ hash: BEST_EFFORT_CANCELED_ORDERS_CACHE_KEY, key: orderId }, client),
  ]);

  if (canceledScore !== null) {
    return CanceledOrderStatus.CANCELED;
  }

  if (bestEffortCanceledScore !== null) {
    return CanceledOrderStatus.BEST_EFFORT_CANCELED;
  }

  return CanceledOrderStatus.NOT_CANCELED;
}

export async function removeOrderFromCaches(
  orderId: string,
  client: RedisClient,
): Promise<void> {
  await Promise.all([
    zRemAsync({ hash: CANCELED_ORDERS_CACHE_KEY, key: orderId }, client),
    zRemAsync({ hash: BEST_EFFORT_CANCELED_ORDERS_CACHE_KEY, key: orderId }, client),
  ]);
}

/**
 * addCanceledOrderId adds the order id to the best effort canceled orders cache.
 *
 * @param orderId
 * @param timestamp
 * @param client
 */
export async function addBestEffortCanceledOrderId(
  orderId: string,
  timestamp: number,
  client: RedisClient,
): Promise<number> {
  return addOrderIdtoCache(orderId, timestamp, client, BEST_EFFORT_CANCELED_ORDERS_CACHE_KEY);
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
  return addOrderIdtoCache(orderId, timestamp, client, CANCELED_ORDERS_CACHE_KEY);
}

/**
 * addCanceledOrderId adds the order id to the cacheKey's cache.
 *
 * @param orderId
 * @param timestamp
 * @param client
 */
export async function addOrderIdtoCache(
  orderId: string,
  timestamp: number,
  client: RedisClient,
  cacheKey: string,
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
        cacheKey,
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
