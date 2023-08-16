import { NextFundingCache } from '@dydxprotocol-indexer/redis';
import Big from 'big.js';

import { redisClient } from '../../src/helpers/redis/redis-controller';

export async function expectNextFundingRate(
  ticker: string,
  rate: Big | undefined,
): Promise<void> {
  const rates: { [ticker: string]: Big | undefined } = await NextFundingCache.getNextFunding(
    redisClient,
    [ticker],
  );
  expect(rates[ticker]).toEqual(rate);
}
