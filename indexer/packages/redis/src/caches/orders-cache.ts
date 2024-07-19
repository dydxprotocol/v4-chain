import { OrderTable } from '@dydxprotocol-indexer/postgres';
import {
  IndexerOrderId,
  RedisOrder,
} from '@dydxprotocol-indexer/v4-protos';
import { RedisClient } from 'redis';

import { getAsync } from '../helpers/redis';

// Cache of order uuid to encoded `RedisOrder`
export const ORDERS_CACHE_KEY_PREFIX: string = 'v4/orders/';

/**
 * Get an order by the UUID of the order.
 * @param orderUuid Indexer assigned UUID of the order to get.
 * @param client Redis client.
 * @returns `RedisOrder` if the order exists in the cache, otherwise `null`.
 */
export async function getOrder(orderUuid: string, client: RedisClient): Promise<RedisOrder | null> {
  const orderString: string | null = await getAsync(getOrderCacheKeyWithUUID(orderUuid), client);

  if (orderString === null) {
    return null;
  }

  return RedisOrder.decode(Buffer.from(orderString, 'binary'));
}

export function getOrderCacheKey(orderId: IndexerOrderId): string {
  return getOrderCacheKeyWithUUID(OrderTable.orderIdToUuid(orderId));
}

function getOrderCacheKeyWithUUID(orderUuid: string): string {
  return `${ORDERS_CACHE_KEY_PREFIX}${orderUuid}`;
}
