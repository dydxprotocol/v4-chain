import Big from 'big.js';
import { RedisClient } from 'redis';

import { deleteAsync, lRangeAsync, rPushAsync } from '../helpers/redis';

const KEY_PREFIX: string = 'v4/nextFunding/';

function getKey(ticker: string): string {
  return `${KEY_PREFIX}${ticker}`;
}

/**
 * Get the next funding rate for given tickers.
 *
 * If the actual funding rate was just published by the protocol, the next funding
 * rate will be undefined.
 *
 * @param client
 * @param tickers
 */
export async function getNextFunding(
  client: RedisClient,
  tickers: string[],
): Promise<{ [ticker: string]: Big | undefined }> {
  const fundingRates: { [ticker: string]: Big | undefined } = {};
  await Promise.all(
    tickers.map(async (ticker: string) => {
      const rates: string[] = await lRangeAsync(
        getKey(ticker),
        client,
      );
      // get average of rates
      if (rates.length > 0) {
        const sum: Big = rates.reduce(
          (acc: Big, val: string) => acc.plus(new Big(val)),
          new Big(0),
        );
        const avg: Big = sum.div(rates.length);
        fundingRates[ticker] = avg;
      } else {
        fundingRates[ticker] = undefined;
      }
    }),
  );
  return fundingRates;
}

export async function addFundingSample(
  ticker: string,
  rate: Big,
  client: RedisClient,
): Promise<number> {
  return rPushAsync({
    key: getKey(ticker),
    value: rate.toString(),
  }, client);
}

export async function clearFundingSamples(
  ticker: string,
  client: RedisClient,
): Promise<number> {
  return deleteAsync(
    getKey(ticker),
    client,
  );
}
