import { OrderTable } from '@dydxprotocol-indexer/postgres';
import { IndexerOrderId, RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import { Callback, RedisClient } from 'redis';

import { RemoveOrderResult } from '../types';
import { validateOrderId } from './helpers';
import { ORDER_EXPIRY_CACHE_KEY } from './order-expiry-cache';
import { getOrderCacheKey } from './orders-cache';
import { getOrderDataCacheKey } from './orders-data-cache';
import { removeOrderScript } from './scripts';
import { getSubaccountOrderIdsCacheKey } from './subaccount-order-ids-cache';

// Number of keys that the lua script will access.
const numKeys: number = 4;

/**
 * Updates order caches in Redis for an order being removed. Evaluates the `remove_order.lua`
 * script to update the caches.
 * The behavior of this should be:
 * - if the order exists in the caches:
 *   - the order should removed from the order and order data cache, and the order id should be
 *     removed from the subaccount order ids cache for the subaccount the order belongs to
 * - if the order does not exist in the caches:
 *   - no-op
 * See the `remove_order.lua` script for more context.
 * @param param0 Contains the `OrderId` of the order to remove.
 * @returns `RemoveOrderResult` for the result of removing the order.
 */
export async function removeOrder({
  removedOrderId,
  client,
}: {
  removedOrderId: IndexerOrderId,
  client: RedisClient,
}): Promise<RemoveOrderResult> {
  validateOrderId(removedOrderId);

  let evalAsync: (
    orderKey: string,
    orderDataKey: string,
    subaccountOrderIdsKey: string,
    orderExpiryKey: string,
    orderId: string,
  ) => Promise<RemoveOrderResult> = (
    orderKey,
    orderDataKey,
    subaccountOrderIdsKey,
    orderExpiryKey,
    orderId,
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
          removed,
          filledStr,
          bookStr,
          removedOrder,
        ]:[
          number,
          string,
          string,
          string,
        ] = results;
        return resolve({
          removed: removed === 1,
          totalFilledQuantums: removed === 1 ? parseInt(filledStr, 10) : undefined,
          restingOnBook: removed === 1 ? bookStr === 'true' : undefined,
          removedOrder: removed === 1 ? RedisOrder.decode(Buffer.from(removedOrder, 'binary')) : undefined,
        });
      };
      client.evalsha(
        removeOrderScript.hash,
        numKeys,
        orderKey,
        orderDataKey,
        subaccountOrderIdsKey,
        orderExpiryKey,
        orderId,
        callback,
      );
    });
  };
  evalAsync = evalAsync.bind(client);

  return evalAsync(
    getOrderCacheKey(removedOrderId),
    getOrderDataCacheKey(removedOrderId),
    getSubaccountOrderIdsCacheKey(removedOrderId.subaccountId!),
    ORDER_EXPIRY_CACHE_KEY,
    OrderTable.orderIdToUuid(removedOrderId),
  );
}
