import Big from 'big.js';
import { Callback, RedisClient } from 'redis';

import { getOrderBookMidPrice } from './orderbook-levels-cache';
import {
  addOrderbookMidPricesScript,
  getOrderbookMidPricesScript,
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
 * Fetches and caches mid prices for multiple tickers.
 * @param client The Redis client
 * @param tickers An array of ticker symbols
 * @returns A promise that resolves when all prices are fetched and cached
 */
export async function fetchAndCacheOrderbookMidPrices(
  client: RedisClient,
  tickers: string[],
): Promise<void> {
  // Fetch midPrices and filter out undefined values
  const cacheKeyPricePairs = await Promise.all(
    tickers.map(async (ticker) => {
      const cacheKey = getOrderbookMidPriceCacheKey(ticker);
      const midPrice = await getOrderBookMidPrice(cacheKey, client);
      if (midPrice !== undefined) {
        return { cacheKey, midPrice };
      }
      return null; // Return null for undefined midPrice
    }),
  );

  // Filter out null values
  const validPairs = cacheKeyPricePairs.filter(
    (pair): pair is { cacheKey: string, midPrice: string } => pair !== null,
  );
  if (validPairs.length === 0) {
    // No valid midPrices to cache
    return;
  }

  const nowSeconds = Math.floor(Date.now() / 1000); // Current time in seconds
  // Extract cache keys and prices
  const priceCacheKeys = validPairs.map((pair) => pair.cacheKey);
  const priceValues = validPairs.map((pair) => pair.midPrice);

  return new Promise<void>((resolve, reject) => {
    client.evalsha(
      addOrderbookMidPricesScript.hash,
      priceCacheKeys.length,
      ...priceCacheKeys,
      ...priceValues,
      nowSeconds,
      (err: Error | null) => {
        if (err) {
          reject(err);
        } else {
          resolve();
        }
      },
    );
  });
}

/**
 * Retrieves the median prices for a given array of tickers from the cache.
 * @param client The Redis client
 * @param tickers Array of ticker symbols
 * @returns A promise that resolves with an object mapping tickers
 *  to their median prices (as strings) or undefined if not found
 */
export async function getMedianPrices(
  client: RedisClient,
  tickers: string[],
): Promise<{ [ticker: string]: string | undefined }> {

  let evalAsync: (
    marketCacheKeys: string[],
  ) => Promise<string[][]> = (
    marketCacheKeys,
  ) => {
    return new Promise((resolve, reject) => {
      const callback: Callback<string[][]> = (
        err: Error | null,
        results: string[][],
      ) => {
        if (err) {
          return reject(err);
        }
        return resolve(results);
      };

      client.evalsha(
        getOrderbookMidPricesScript.hash,  // The Lua script to get cached prices
        marketCacheKeys.length,
        ...marketCacheKeys,
        callback,
      );
    });
  };
  evalAsync = evalAsync.bind(client);

  // Map tickers to cache keys
  const marketCacheKeys = tickers.map(getOrderbookMidPriceCacheKey);
  // Fetch the prices arrays from Redis (without scores)
  const pricesArrays = await evalAsync(marketCacheKeys);

  const result: { [ticker: string]: string | undefined } = {};
  tickers.forEach((ticker, index) => {
    const prices = pricesArrays[index];

    // Check if there are any prices
    if (!prices || prices.length === 0) {
      result[ticker] = undefined;
      return;
    }

    // Convert the prices to Big.js objects for precision
    const bigPrices = prices.map((price) => Big(price));

    // Sort the prices in ascending order
    bigPrices.sort((a, b) => a.cmp(b));

    // Calculate the median
    const mid = Math.floor(bigPrices.length / 2);
    if (bigPrices.length % 2 === 1) {
      // Odd number of prices: the middle one is the median
      result[ticker] = bigPrices[mid].toFixed();
    } else {
      // Even number of prices: average the two middle ones
      result[ticker] = bigPrices[mid - 1].plus(bigPrices[mid]).div(2).toFixed();
    }
  });

  return result;
}
