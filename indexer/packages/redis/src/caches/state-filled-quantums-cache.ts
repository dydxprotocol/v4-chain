import { RedisClient } from 'redis';

import { getAsync, setexAsync } from '../helpers/redis';

export const STATE_FILLED_QUANTUMS_CACHE_KEY_PREFIX: string = 'v4/state_filled_quantums/';
export const STATE_FILLED_QUANTUMS_TTL_SECONDS: number = 300; // 5 minutes

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
