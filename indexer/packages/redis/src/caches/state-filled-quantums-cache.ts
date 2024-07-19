import { RedisClient } from 'redis';

import { getAsync, setexAsync } from '../helpers/redis';

export const STATE_FILLED_QUANTUMS_CACHE_KEY_PREFIX: string = 'v4/state_filled_quantums/';
export const STATE_FILLED_QUANTUMS_TTL_SECONDS: number = 300; // 5 minutes

/**
 * Updates the state-filled quantums for an order id. This is the total filled quantums of the order
 * in the state of the network.
 * @param orderId
 * @param filledQuantums
 * @param client
 */
export async function updateStateFilledQuantums(
  orderId: string,
  filledQuantums: string,
  client: RedisClient,
): Promise<void> {
  await setexAsync({
    key: getCacheKey(orderId),
    value: filledQuantums,
    timeToLiveSeconds: STATE_FILLED_QUANTUMS_TTL_SECONDS,
  }, client);
}

/**
 * Gets the state-filled quantums for an order id. This is the total filled quantums of the order
 * in the state of the network.
 * @param orderId
 * @param client
 * @returns
 */
export async function getStateFilledQuantums(
  orderId: string,
  client: RedisClient,
): Promise<string | undefined> {
  const filledQuantums: string | null = await getAsync(
    getCacheKey(orderId),
    client,
  );

  if (filledQuantums === null) {
    return undefined;
  }

  return filledQuantums;
}

export function getCacheKey(orderId: string): string {
  return `${STATE_FILLED_QUANTUMS_CACHE_KEY_PREFIX}${orderId}`;
}
