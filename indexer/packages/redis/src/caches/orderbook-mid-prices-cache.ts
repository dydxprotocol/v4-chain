import Big from 'big.js';
import { Callback, RedisClient } from 'redis';

import {
  addMarketPriceScript,
  getMarketMedianScript,
} from './scripts';

// Cache of orderbook prices for each clob pair
// Each price is cached for a 5 second window and in a ZSET
export const ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX: string = 'v4/orderbook_mid_prices/';

/**
 * Generates a cache key for a given ticker's orderbook mid price.
 * @param ticker The ticker symbol
 * @returns The cache key string
 */
function getOrderbookMidPriceCacheKey(ticker: string): string {
  return `${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${ticker}`;
}

/**
 * Adds a price to the market prices cache for a given ticker.
 * Uses a Lua script to add the price with a timestamp to a sorted set in Redis.
 * @param client The Redis client
 * @param ticker The ticker symbol
 * @param price The price to be added
 * @returns A promise that resolves when the operation is complete
 */
export async function setPrice(
  client: RedisClient,
  ticker: string,
  price: string,
): Promise<void> {
  // Number of keys for the lua script.
  const numKeys: number = 1;

  let evalAsync: (
    marketCacheKey: string,
  ) => Promise<void> = (marketCacheKey) => {

    return new Promise<void>((resolve, reject) => {
      const callback: Callback<void> = (
        err: Error | null,
      ) => {
        if (err) {
          return reject(err);
        }
        return resolve();
      };

      const nowSeconds = Math.floor(Date.now() / 1000); // Current time in seconds
      client.evalsha(
        addMarketPriceScript.hash,
        numKeys,
        marketCacheKey,
        price,
        nowSeconds,
        callback,
      );

    });
  };
  evalAsync = evalAsync.bind(client);

  return evalAsync(
    getOrderbookMidPriceCacheKey(ticker),
  );
}

/**
 * Retrieves the median price for a given ticker from the cache.
 * Uses a Lua script to fetch either the middle element (for odd number of prices)
 * or the two middle elements (for even number of prices) from a sorted set in Redis.
 * If two middle elements are returned, their average is calculated in JavaScript.
 * @param client The Redis client
 * @param ticker The ticker symbol
 * @returns A promise that resolves with the median price as a string, or null if not found
 */
export async function getMedianPrice(client: RedisClient, ticker: string): Promise<string | null> {
  let evalAsync: (
    marketCacheKey: string,
  ) => Promise<string[]> = (
    marketCacheKey,
  ) => {
    return new Promise((resolve, reject) => {
      const callback: Callback<string[]> = (
        err: Error | null,
        results: string[],
      ) => {
        if (err) {
          return reject(err);
        }
        return resolve(results);
      };

      client.evalsha(
        getMarketMedianScript.hash,
        1,
        marketCacheKey,
        callback,
      );
    });
  };
  evalAsync = evalAsync.bind(client);

  const prices = await evalAsync(
    getOrderbookMidPriceCacheKey(ticker),
  );

  if (!prices || prices.length === 0) {
    return null;
  }

  if (prices.length === 1) {
    return Big(prices[0]).toFixed();
  }

  if (prices.length === 2) {
    const [price1, price2] = prices.map((price) => {
      return Big(price);
    });
    return price1.plus(price2).div(2).toFixed();
  }

  return null;
}
