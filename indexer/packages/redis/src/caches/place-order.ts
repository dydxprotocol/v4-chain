import { OrderTable } from '@dydxprotocol-indexer/postgres';
import { ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import { RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import { Callback, RedisClient } from 'redis';

import { PlaceOrderResult } from '../types';
import { getOrderExpiry, validateRedisOrder } from './helpers';
import { ORDER_EXPIRY_CACHE_KEY } from './order-expiry-cache';
import { getOrderCacheKey } from './orders-cache';
import { getOrderDataCacheKey } from './orders-data-cache';
import { placeOrderScript } from './scripts';
import { getSubaccountOrderIdsCacheKey } from './subaccount-order-ids-cache';

// Number of keys that the lua script will access.
const numKeys: number = 4;

/**
 * Updates order caches in Redis for an order being placed/replaced. Evaluates the `place_order.lua`
 * script to update the caches.
 * The behavior of this should be:
 * - if the order is new:
 *   - the encoded `RedisOrder` should be saved to the `ORDERS_CACHE`, mapped to the order uuid
 *   - {good-til-block/sequence-number}_{totalFilled}_false should be saved to the
 *     `ORDERS_DATA_CACHE,` mapped to the order uuid, `totalFilled` being set to 0
 *   - the order uuid should be added to the `SUBACCOUNT_ORDERS_CACHE` for the subaccount that
 *     placed the order
 *   - if the order is a short-term order:
 *     - the order uuid should be added to the `ORDER_EXPIRY_CACHE` with its score set to its expiry
 * - if the order exists in the cache, and the placed order has a lower or equal expiry
 *   (good-til-block/seq. number)
 *   - no updates
 * - if the order exists in the cache, and the placed order has a greater expiry
 *   - the encoded `RedisOrder` should be saved to the `ORDERS_CACHE`, replacing the existing order
 *   - {good-til-block/sequence-number}_{totalFilled}_false should be saved to the
 *     `ORDERS_DATA_CACHE`, replacing the existing order, with `totalFilled` being the value for the
 *     existing order
 *   - if the order is a short-term order:
 *     - the score for the order uuid in the `ORDER_EXPIRY_CACHE` should be updated to the new
 *       expiry
 * See the `place_order.lua` script for more context.
 * @param param0 Contains the `RedisOrder` of the order being placed.
 * @returns `PlaceOrderResult` for the result of placing the order.
 */
export async function placeOrder({
  redisOrder,
  client,
}: {
  redisOrder: RedisOrder,
  client: RedisClient,
}): Promise<PlaceOrderResult> {
  validateRedisOrder(redisOrder);

  let evalAsync: (
    orderKey: string,
    orderDataKey: string,
    subaccountOrderIdsKey: string,
    orderExpiryKey: string,
    encodedOrder: string,
    orderExpiry: string,
    orderId: string,
    isShortTermOrder: boolean,
  ) => Promise<PlaceOrderResult> = (
    orderKey,
    orderDataKey,
    subaccountOrderIdsKey,
    orderExpiryKey,
    encodedOrder,
    orderExpiry,
    orderId,
    isShortTermOrder,
  ) => {
    return new Promise((resolve, reject) => {
      const callback: Callback<[number, number, string, string, string]> = (
        err: Error | null,
        results: [number, number, string, string, string],
      ) => {
        if (err) {
          return reject(err);
        }
        const [
          placed,
          replaced,
          filledStr,
          bookStr,
          oldOrder,
        ]:[
          number,
          number,
          string,
          string,
          string,
        ] = results;
        return resolve({
          placed: placed === 1,
          replaced: replaced === 1,
          oldTotalFilledQuantums: replaced === 1 ? parseInt(filledStr, 10) : undefined,
          restingOnBook: replaced === 1 ? bookStr === 'true' : undefined,
          oldOrder: replaced === 1 ? RedisOrder.decode(Buffer.from(oldOrder, 'binary')) : undefined,
        });
      };
      client.evalsha(
        placeOrderScript.hash,
        numKeys,
        orderKey,
        orderDataKey,
        subaccountOrderIdsKey,
        orderExpiryKey,
        encodedOrder,
        orderExpiry,
        orderId,
        isShortTermOrder.toString(),
        callback,
      );
    });
  };
  evalAsync = evalAsync.bind(client);

  return evalAsync(
    getOrderCacheKey(redisOrder.order!.orderId!),
    getOrderDataCacheKey(redisOrder.order!.orderId!),
    getSubaccountOrderIdsCacheKey(redisOrder.order!.orderId!.subaccountId!),
    ORDER_EXPIRY_CACHE_KEY,
    // TODO: use String to directly convert the UInt8Array to a string
    Buffer.from(RedisOrder.encode(redisOrder).finish()).toString('binary'),
    getOrderExpiry(redisOrder.order!).toString(),
    OrderTable.orderIdToUuid(redisOrder.order!.orderId!),
    redisOrder.order!.orderId!.orderFlags === ORDER_FLAG_SHORT_TERM,
  );
}
