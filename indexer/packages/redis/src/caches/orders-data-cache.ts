import { OrderTable } from '@dydxprotocol-indexer/postgres';
import { IndexerOrderId } from '@dydxprotocol-indexer/v4-protos';
import { RedisClient } from 'redis';

import { getAsync } from '../helpers/redis';
import { OrderData } from '../types';

// Cache of order uuid to [good-til-block or sequence number]_[totalFilled]_[resting on book]
// where `totalFilled` = total quantums filled for an order,
// and `resting on book` = true/false, indicating if the order is resting on the orderbook.
// These values are for use in Lua scripts.
export const ORDERS_DATA_CACHE_KEY_PREFIX: string = 'v4/orderData/';

export async function getOrderData(
  orderId: IndexerOrderId,
  redisClient: RedisClient,
): Promise<OrderData | null> {
  const key: string = getOrderDataCacheKey(orderId);
  return getOrderDataFromCacheKey(key, redisClient);
}

export async function getOrderDataWithUUID(
  orderUuid: string,
  redisClient: RedisClient,
): Promise<OrderData | null> {
  const key: string = getOrderDataCacheKeyWithUUID(orderUuid);
  return getOrderDataFromCacheKey(key, redisClient);
}

async function getOrderDataFromCacheKey(
  cacheKey: string,
  redisClient: RedisClient,
): Promise<OrderData | null> {
  const orderDataString: string | null = await getAsync(cacheKey, redisClient);
  if (orderDataString === null) {
    return null;
  }
  return orderDataStringToOrderData(orderDataString);
}

function orderDataStringToOrderData(orderDataString: string): OrderData {
  const [
    goodTilBlock,
    totalFilledQuantums,
    restingOnBook,
  ]: [string, string, string] = orderDataString.split('_') as [string, string, string];

  return {
    goodTilBlock,
    totalFilledQuantums,
    restingOnBook: restingOnBook === 'true',
  };
}

export function getOrderDataCacheKey(orderId: IndexerOrderId): string {
  return getOrderDataCacheKeyWithUUID(OrderTable.orderIdToUuid(orderId));
}

function getOrderDataCacheKeyWithUUID(orderUuid: string): string {
  return `${ORDERS_DATA_CACHE_KEY_PREFIX}${orderUuid}`;
}
