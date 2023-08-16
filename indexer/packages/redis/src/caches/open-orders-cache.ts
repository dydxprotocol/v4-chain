import { RedisClient } from 'redis';

import {
  hDelAsync,
  hGetAllAsync,
  hSetAsync,
} from '../helpers/redis';

// Cache of open orders in a market
export const OPEN_ORDERS_CACHE_KEY_PREFIX: string = 'v4/open_orders/';

export async function addOpenOrder(
  orderUuid: string,
  clobPairId: string,
  client: RedisClient,
): Promise<void> {
  await hSetAsync(
    {
      hash: getOpenOrderCacheKey(clobPairId),
      pairs: { [orderUuid]: '1' },
    },
    client,
  );
}

export async function removeOpenOrder(
  orderUuid: string,
  clobPairId: string,
  client: RedisClient,
): Promise<void> {
  await hDelAsync(
    {
      hash: getOpenOrderCacheKey(clobPairId),
      keys: [orderUuid],
    },
    client,
  );
}

export async function getOpenOrderIds(
  clobPairId: string,
  client: RedisClient,
): Promise<string[]> {
  const pairs: {[key: string]: string} = await hGetAllAsync(
    getOpenOrderCacheKey(clobPairId),
    client,
  );
  return Object.keys(pairs);
}

function getOpenOrderCacheKey(clobPairId: string): string {
  return `${OPEN_ORDERS_CACHE_KEY_PREFIX}${clobPairId}`;
}
