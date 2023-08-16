import { IndexerOrderId, RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import { Callback, RedisClient } from 'redis';

import { InvalidTotalFilledQuantumsError } from '../errors';
import { UpdateOrderResult } from '../types';
import { validateOrderId } from './helpers';
import { getOrderCacheKey } from './orders-cache';
import { getOrderDataCacheKey } from './orders-data-cache';
import { updateOrderScript } from './scripts';

// Number of keys that the lua script will access.
const numKeys: number = 2;

/**
 * Updates order caches in Redis for an order being updated. Evaluates the `update_order.lua`
 * script to update the caches.
 * The behavior of this should be:
 * - if the order exists in the caches:
 *   - the order data cache should be updated with the new total filled quantums of the order, and
 *     `resting_on_book` should be set to "true" for the order
 * - if the order does not exist in the caches:
 *   - no-op
 * See the `update_order.lua` script for more context.
 * @param param0 Contains the `OrderId` of the order to update, and the new total filled quantums.
 * @returns `UpdateOrderResult` for the result of removing the order.
 */
export async function updateOrder({
  updatedOrderId,
  newTotalFilledQuantums,
  client,
}: {
  updatedOrderId: IndexerOrderId,
  newTotalFilledQuantums: number,
  client: RedisClient,
}): Promise<UpdateOrderResult> {
  validateOrderId(updatedOrderId);

  if (newTotalFilledQuantums < 0) {
    throw new InvalidTotalFilledQuantumsError(
      `Total filled cannot be negative, but was ${newTotalFilledQuantums}`,
    );
  }

  let evalAsync: (
    orderKey: string,
    orderDataKey: string,
    totalFilledQuantums: number,
  ) => Promise<UpdateOrderResult> = (
    orderKey,
    orderDataKey,
    totalFilledQuantums,
  ) => {
    return new Promise((resolve, reject) => {
      const callback: Callback<[number, string, string, string]> = (
        err: Error | null,
        results: [number, string, string, string],
      ) => {
        if (err) {
          return reject(err);
        }
        const [
          updated,
          filledStr,
          restingOnBookStr,
          order,
        ]:[
          number,
          string,
          string,
          string,
        ] = results;
        return resolve({
          updated: updated === 1,
          oldTotalFilledQuantums: updated === 1 ? parseInt(filledStr, 10) : undefined,
          oldRestingOnBook: updated === 1 ? restingOnBookStr === 'true' : undefined,
          order: updated === 1 ? RedisOrder.decode(Buffer.from(order, 'binary')) : undefined,
        });
      };
      client.evalsha(
        updateOrderScript.hash,
        numKeys,
        orderKey,
        orderDataKey,
        totalFilledQuantums,
        callback,
      );
    });
  };
  evalAsync = evalAsync.bind(client);

  return evalAsync(
    getOrderCacheKey(updatedOrderId),
    getOrderDataCacheKey(updatedOrderId),
    newTotalFilledQuantums,
  );
}
