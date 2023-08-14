import _ from 'lodash';
import { RedisClient } from 'redis';

import { zRangeByScoreAsync } from '../helpers/redis';

// Cache of order expiry to order uuid
export const ORDER_EXPIRY_CACHE_KEY: string = 'v4/orderExpiry';

/**
 * Grabs and returns a mapping from order UUID to expiry value for expiries on/before the given
 * expiry.
 * @param {Object} args
 * * **latestExpiry** `number` - The most recent expiry value to retrieve.
 * * **latestExpiryIsInclusive** `boolean[=true]` - Whether to include entries with an expiry of
 * `latestExpiry` or only earlier expiries.
 * @param {RedisClient} client - Redis client.
 * @returns {Promise<Record<string, Number>>} A mapping from order uuid to expiry value.
 */
export async function getOrdersAndExpiries(
  {
    latestExpiry,
    latestExpiryIsInclusive = true,
  }: {
    latestExpiry: number,
    latestExpiryIsInclusive?: boolean,
  },
  client: RedisClient,
): Promise<Record<string, Number>> {

  const rawResults: string[] = await zRangeByScoreAsync({
    key: ORDER_EXPIRY_CACHE_KEY,
    start: -Infinity,
    end: latestExpiry,
    endIsInclusive: latestExpiryIsInclusive,
    withScores: true,
  }, client);
  return _.fromPairs(
    _.map(
      _.chunk(rawResults, 2),
      (keyValuePair) => [keyValuePair[0], Number(keyValuePair[1])],
    ),
  );
}

/**
 * Grabs all Order UUIDs with expiries on/before the given expiry.
 * @param {Object} args
 * * **latestExpiry** `number` - The most recent expiry value to retrieve.
 * * **latestExpiryIsInclusive** `boolean[=true]` - Whether to include entries with an expiry of
 * `latestExpiry` or only earlier expiries.
 * @param {RedisClient} client - Redis client.
 * @returns {Promise<string[]>} An array of order UUIDs.
 */
export async function getOrderExpiries(
  {
    latestExpiry,
    latestExpiryIsInclusive = true,
  }: {
    latestExpiry: number,
    latestExpiryIsInclusive?: boolean,
  },
  client: RedisClient,
): Promise<string[]> {
  return zRangeByScoreAsync({
    key: ORDER_EXPIRY_CACHE_KEY,
    start: -Infinity,
    end: latestExpiry,
    endIsInclusive: latestExpiryIsInclusive,
    withScores: false,
  }, client);
}
