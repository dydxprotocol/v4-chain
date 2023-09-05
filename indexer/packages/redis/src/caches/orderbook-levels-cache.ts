import { logger } from '@dydxprotocol-indexer/base';
import { OrderSide } from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import _ from 'lodash';
import { Callback, RedisClient } from 'redis';

import { InvalidOptionsError, InvalidPriceLevelUpdateError } from '../errors';
import { hGetAsync } from '../helpers/redis';
import { OrderbookLevels, PriceLevel } from '../types';
import { deleteZeroPriceLevelScript, getOrderbookSideScript, incrementOrderbookLevelScript } from './scripts';

// Cache of orderbook levels for each clob pair
// Each side of each exchange pair is an HSET with the hash = price, and value = total size of
// orders at the price in quantums
// TODO(CORE-512): add info/resources around caches. Doc:
// https://www.notion.so/dydx/Indexer-Technical-Spec-a6b15644502048f994c98dee35b96e96#61d5f8ca5117476caab78b3f0691b1d0
export const ORDERS_CACHE_KEY_PREFIX: string = 'v4/orderbookLevels/';

/**
 * Update the total size of orders at a price level for a specific ticker/side with a delta. The
 * delta is in quantums (integers) rather than a human-readable value as floating point math in
 * Redis is inexact.
 * @param param0 Ticker of the exchange pair, side, human readable, and the delta to apply to the
 * total size in quantums.
 * @returns The updated total size in quantums.
 */
export async function updatePriceLevel({
  ticker,
  side,
  humanPrice,
  sizeDeltaInQuantums,
  client,
}: {
  ticker: string,
  side: OrderSide,
  humanPrice: string,
  sizeDeltaInQuantums: string,
  client: RedisClient,
// TODO(DEC-1314): Return a string once redis client is updated to use `stringNumbers` option.
}): Promise<number> {
  const updatedQuantums: number = await incrementOrderbookLevel(
    ticker,
    side,
    humanPrice,
    sizeDeltaInQuantums,
    client,
  );

  // This case should never happen while the protcol and indexer are working correctly, and emits a
  // critical error log. As updates to price levels come from updates to orders, and updates to
  // are procesed in order of place -> update -> remove, there should never be quantums removed from
  // a price level in excess of the quantums added by an order.
  // NOTE: If this happens from a single price level update, it's possible for multiple subsequent
  // price level updates to fail with the same error due to interleaved price level updates.
  if (updatedQuantums < 0) {
    // Undo the update. This can't be done in a Lua script as Redis runs Lua 5.1, which only
    // uses doubles which support up to 53-bit integers. Race-condition where it's possible for a
    // price-level to have negative quantums handled in `getOrderbookLevels` where price-levels with
    // negative quantums are filtered out. Note: even though we are reverting this information, each
    // call to incrementOrderbookLevel updates the lastUpdated key in the cache.
    await incrementOrderbookLevel(
      ticker,
      side,
      humanPrice,
      // Needs to be an integer
      Big(sizeDeltaInQuantums).mul(-1).toFixed(0),
      client,
    );
    logger.crit({
      at: 'orderbookLevelsCache#updatePriceLevel',
      message: 'Price level updated to negative quantums',
      ticker,
      side,
      humanPrice,
      updatedQuantums,
      sizeDeltaInQuantums,
    });
    throw new InvalidPriceLevelUpdateError(
      '#updatePriceLevel: Resulting price level has negative quantums, quantums = ' +
      `${updatedQuantums}`,
    );
  }

  return updatedQuantums;
}

/**
 * Update the orderbooks level cache and its lastUpdated cache values.
 * @param ticker Ticker of the exchange pair.
 * @param side OrderSide of the orderbook.
 * @param humanPrice Human readable price key in the HSET.
 * @param sizeDeltaInQuantums Delta to apply to the total size in quantums.
 * @param client Redis client.
 */
async function incrementOrderbookLevel(
  ticker: string,
  side: OrderSide,
  humanPrice: string,
  sizeDeltaInQuantums: string,
  client: RedisClient,
): Promise<number> {
  // Number of keys for the lua script.
  const numKeys: number = 2;

  let evalAsync: (
    orderbookKey: string,
    lastUpdatedKey: string,
    priceLevel: string,
    delta: string,
  ) => Promise<number> = (
    orderbookKey,
    lastUpdatedKey,
    priceLevel,
    delta,
  ) => {
    return new Promise((resolve, reject) => {
      const callback: Callback<number> = (
        err: Error | null,
        results: number,
      ) => {
        if (err) {
          return reject(err);
        }
        return resolve(results);
      };
      client.evalsha(
        incrementOrderbookLevelScript.hash,
        numKeys,
        orderbookKey,
        lastUpdatedKey,
        priceLevel,
        delta,
        callback,
      );
    });
  };
  evalAsync = evalAsync.bind(client);

  return evalAsync(
    getKey(ticker, side),
    getLastUpdatedKey(ticker, side),
    humanPrice,
    sizeDeltaInQuantums,
  );
}

/**
 * Get the orderbook level at a specific ticker, side, and human price.
 */
export async function getOrderbookLevel(
  ticker: string,
  side: OrderSide,
  humanPrice: string,
  client: RedisClient,
): Promise<string> {
  const result: string | null = await hGetAsync(
    {
      hash: getKey(ticker, side),
      key: humanPrice,
    },
    client,
  );

  return result ?? '0';
}

/**
 * Get the order book price levels for a specific exchange pair.
 * @param ticker Ticker of the exchange pair.
 * @param client Redis client
 * @param options Additional options to apply transformations to the orderbook before returning
 *  - removeZeros: Remove orderbook levels with zero quantums before returning, if this option
 *    is not defined, behavior is to default to remove all orderbook levels with zero quantums
 *  - sortSides: Sorts bids in descending order and asks in ascending order
 *  - uncrossBook: Uncross the book before returning, such that the best bid < best ask, the order
 *    book may be crossed due to receiving updates to orders out-of-order, as ordering is only
 *    guaranteed per-order, and not per price-level. Used to return an uncrossed orderbook in the
 *    REST API.
 *  - limitPerSide: Returns a maximum number of levels per side of the book.
 *    Considered to be infinite if undefined or a non-positive number.
 *    Can only be set if sortSides is true.
 * @returns Order book price levels for the exchange pair
 */
export async function getOrderBookLevels(
  ticker: string,
  client: RedisClient,
  options: {
    removeZeros?: boolean,
    sortSides?: boolean,
    uncrossBook?: boolean,
    limitPerSide?: number,
  } = {},
): Promise<OrderbookLevels> {
  // Sanity-check the options.
  if (options.sortSides !== true) {
    if (options.uncrossBook === true) {
      throw new InvalidOptionsError(
        '#getOrderbookLevels: uncrossBook cannot be true if sortSides is not true',
      );
    }
    if (options.limitPerSide !== undefined) {
      throw new InvalidOptionsError(
        '#getOrderbookLevels: limitPerSide cannot be defined if sortSides is not true ',
      );
    }
  }

  // Default to removing zeros unless `false` is passed in
  const removeZeros: boolean = options.removeZeros ?? true;

  let [
    bids,
    asks,
  ]: [
    PriceLevel[],
    PriceLevel[],
  ] = await Promise.all([
    getOrderbookSide(ticker, OrderSide.BUY, client, removeZeros),
    getOrderbookSide(ticker, OrderSide.SELL, client, removeZeros),
  ]);

  // Sort bids in descending order. Sort asks in ascending order.
  if (options.sortSides === true) {
    bids.sort((a, b) => Number(b.humanPrice) - Number(a.humanPrice));
    asks.sort((a, b) => Number(a.humanPrice) - Number(b.humanPrice));
  }

  // Prevent the bid/ask sides from crossing. Sides are sorted, as an error is thrown above if this
  // option is true while the `sortSides` option is false.
  if (options.uncrossBook) {
    // Keep track of index pointers for bids and asks.
    let ai = 0;
    let bi = 0;

    // While the books are crossing...
    while (
      ai < asks.length &&
      bi < bids.length &&
      Number(bids[bi].humanPrice) >= Number(asks[ai].humanPrice)
    ) {
      // With ordering:
      // 1. Give precedence to newer price level over older price level.
      // 2. Give precedence to the side with the larger size in quantums.
      // 3. If both sides have the same recency and size, give precedence to the ask.
      // This is an arbitrary choice to remove crossing levels in the orderbook.
      if (Number(bids[bi].lastUpdated) > Number(asks[ai].lastUpdated) ||
        Number(bids[bi].quantums) > Number(asks[ai].quantums)) {
        ai += 1;
      } else {
        bi += 1;
      }
    }

    // Remove any price levels that are crossing.
    if (ai > 0) {
      asks = asks.slice(ai);
    }
    if (bi > 0) {
      bids = bids.slice(bi);
    }
  }

  // Limit the number of levels reported per side. Non-positive is considered infinite.
  const limitPerSide: number = options.limitPerSide ?? 0;
  if (limitPerSide > 0) {
    // Only run the costly `slice()` operations if necessary.
    if (asks.length > limitPerSide) {
      asks = asks.slice(0, limitPerSide);
    }
    if (bids.length > limitPerSide) {
      bids = bids.slice(0, limitPerSide);
    }
  }

  return { bids, asks };
}

/**
 * Deletes a zero size price level from the orderbook levels cache idempotently using a Lua script.
 * @param param0 Ticker of the exchange pair, side, human readable price level to delete.
 * @returns `boolean`, true/false for whether the level was deleted.
 */
export async function deleteZeroPriceLevel({
  ticker,
  side,
  humanPrice,
  client,
}: {
  ticker: string,
  side: OrderSide,
  humanPrice: string,
  client: RedisClient,
}): Promise<boolean> {
  // Number of keys for the lua script.
  const numKeys: number = 2;

  let evalAsync: (
    orderbookKey: string,
    lastUpdatedKey: string,
    priceLevel: string,
  ) => Promise<boolean> = (
    orderbookKey,
    lastUpdatedKey,
    priceLevel,
  ) => {
    return new Promise((resolve, reject) => {
      const callback: Callback<number> = (
        err: Error | null,
        results: number,
      ) => {
        if (err) {
          return reject(err);
        }
        const deleted: number = results;
        return resolve(deleted === 1);
      };
      client.evalsha(
        deleteZeroPriceLevelScript.hash,
        numKeys,
        orderbookKey,
        lastUpdatedKey,
        priceLevel,
        callback,
      );
    });
  };
  evalAsync = evalAsync.bind(client);

  return evalAsync(
    getKey(ticker, side),
    getLastUpdatedKey(ticker, side),
    humanPrice,
  );
}

/**
 * Gets the quantums and lastUpdated data from the cache for the given orderbook side.
 * @param param0 Ticker of the exchange pair, side, Redis client.
 * @returns An mapping from human-readable price to objects containing data for the price point.
 * {
 *   "<human-readable-price>": {
 *     "quantums": "<quantums>",
 *     "lastUpdated": "<timestamp>",
 *   },
 *   ...
 * }
 */
export async function getOrderbookSideData({
  ticker,
  side,
  client,
}: {
  ticker: string,
  side: OrderSide,
  client: RedisClient,
}): Promise<PriceLevel[]> {
  // Number of keys for the lua script.
  const numKeys: number = 2;

  let evalAsync: (
    orderbookKey: string,
    lastUpdatedKey: string,
  ) => Promise<string[][]> = (
    orderbookKey,
    lastUpdatedKey,
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
        getOrderbookSideScript.hash,
        numKeys,
        orderbookKey,
        lastUpdatedKey,
        callback,
      );
    });
  };
  evalAsync = evalAsync.bind(client);

  const rawRedisResults: string[][] = await evalAsync(
    getKey(ticker, side),
    getLastUpdatedKey(ticker, side),
  );

  // The Lua script returns a list of lists of strings.
  //   rawRedisResults = [<list of quantums data>, <list of lastUpdated data>].
  //   each subarray = ['key', 'value', 'key', 'value', ...]
  // The 1st list is a flat array of alternating key-value pairs representing prices and quantums.
  // The 2nd is a flat array of alternating key-value pairs representing prices and lastUpdated
  // values.
  const quantumsMapping: {[field: string]: string} = _.fromPairs(_.chunk(rawRedisResults[0], 2));
  const lastUpdatedMapping: {[field: string]: string} = _.fromPairs(_.chunk(rawRedisResults[1], 2));

  return convertToPriceLevels(quantumsMapping, lastUpdatedMapping);

}

export function getKey(ticker: string, side: OrderSide): string {
  return `${ORDERS_CACHE_KEY_PREFIX}${ticker}/${side}`;
}

export function getLastUpdatedKey(ticker: string, side: OrderSide): string {
  return `${getKey(ticker, side)}/lastUpdated`;
}

async function getOrderbookSide(
  ticker: string,
  side: OrderSide,
  client: RedisClient,
  removeZeros: boolean,
): Promise<PriceLevel[]> {
  let sideLevels: PriceLevel[] = await getOrderbookSideData({ ticker, side, client });

  // Remove any negative levels - possible due to race condition in updatePriceLevel
  sideLevels = sideLevels.filter((level: PriceLevel) => Big(level.quantums).gte(Big(0)));
  if (removeZeros) {
    // Remove all zero levels
    sideLevels = sideLevels.filter((level: PriceLevel) => level.quantums !== '0');
  }

  return sideLevels;
}

function convertToPriceLevels(
  price2QuantumsMapping: {[field: string]: string},
  price2LastUpdatedMapping: {[field: string]: string},
): PriceLevel[] {
  const quantumsKeys: string[] = _.keys(price2QuantumsMapping);
  const lastUpdatedKeys: string[] = _.keys(price2LastUpdatedMapping);
  const pricesMissingData: string[] = _.xor(quantumsKeys, lastUpdatedKeys);
  // If the cache is behaving correctly, this should never occur. Price keys should be added and
  // deleted atomically, so any missing keys here could signify a larger problem.
  if (!_.isEmpty(pricesMissingData)) {
    logger.error({
      at: 'orderbook-levels-cache#getOrderbookSideData',
      message: 'Key mismatch detected amongst orderbook levels caches.',
      quantumsKeysWithoutMatchingData: _.intersection(quantumsKeys, pricesMissingData),
      lastUpdatedKeysWithoutMatchingData: _.intersection(lastUpdatedKeys, pricesMissingData),
    });
  }

  const priceKeys: string[] = _.without(quantumsKeys, ...pricesMissingData);
  return _.map(priceKeys, (price: string) => {
    return {
      humanPrice: price,
      quantums: price2QuantumsMapping[price],
      lastUpdated: price2LastUpdatedMapping[price],
    };
  });
}
