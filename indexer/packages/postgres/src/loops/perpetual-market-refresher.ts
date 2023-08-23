import { stats, delay, NodeEnv } from '@dydxprotocol-indexer/base';

import config from '../config';
import * as PerpetualMarketTable from '../stores/perpetual-market-table';
import { Options, PerpetualMarketFromDatabase, PerpetualMarketsMap } from '../types';

let clobPairIdToPerpetualMarket: Record<string, PerpetualMarketFromDatabase> = {};
let tickerToPerpetualMarket: Record<string, PerpetualMarketFromDatabase> = {};
let idToPerpetualMarket: Record<string, PerpetualMarketFromDatabase> = {};

// TODO(DEC-642): Update the in-memory mapping of perpetual market tickers to ids from websocket
// messages from `v4_markets` topic rather than periodically refereshing from the database.

/**
 * Refresh loop to cache the list of all perpetual markets from the database in-memory.
 */
export async function start(): Promise<void> {
  for (;;) {
    await updatePerpetualMarkets();
    await delay(config.PERPETUAL_MARKETS_REFRESHER_INTERVAL_MS);
  }
}

/**
 * Clears the in-memory map of perpetual market clob pair ids to tickers and tickers to clob pair.
 * Used for testing.
 */
export function clear(): void {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw new Error('resetCache cannot be used in non-test env');
  }
  clobPairIdToPerpetualMarket = {};
  tickerToPerpetualMarket = {};
  idToPerpetualMarket = {};
}

/**
 * Updates in-memory map of perpetual market clob pair ids to tickers and tickers to clob pair ids.
 */
export async function updatePerpetualMarkets(options?: Options): Promise<void> {
  const startTime: number = Date.now();
  const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
    {},
    [],
    options || { readReplica: true },
  );

  const tmpClobPairIdToPerpetualMarket: Record<string, PerpetualMarketFromDatabase> = {};
  const tmpTickerToPerpetualMarket: Record<string, PerpetualMarketFromDatabase> = {};
  const tmpIdToPerpetualMarket: Record<string, PerpetualMarketFromDatabase> = {};
  perpetualMarkets.forEach(
    (market: PerpetualMarketFromDatabase) => {
      tmpClobPairIdToPerpetualMarket[market.clobPairId] = market;
      tmpTickerToPerpetualMarket[market.ticker] = market;
      tmpIdToPerpetualMarket[market.id] = market;
    },
  );

  clobPairIdToPerpetualMarket = tmpClobPairIdToPerpetualMarket;
  tickerToPerpetualMarket = tmpTickerToPerpetualMarket;
  idToPerpetualMarket = tmpIdToPerpetualMarket;
  stats.timing(`${config.SERVICE_NAME}.loops.update_perpetual_markets`, Date.now() - startTime);
}

/**
 * Validates a ticker references a perpetual market.
 * @param ticker Ticker to validate.
 * @returns true if ticker matches a perpetual market ticker, false otherwise.
 */
export function isValidPerpetualMarketTicker(ticker: string): boolean {
  return tickerToPerpetualMarket[ticker] !== undefined && tickerToPerpetualMarket[ticker] !== null;
}

/**
 * Gets the clob pair id for a given perpetual market ticker.
 * @param ticker Ticker for the perpetual market.
 * @returns Clob pair id to get perpetual market with the ticker, if no perpetual market exists
 * with the given clob pair id, undefined is returned.
 */
export function getClobPairIdFromTicker(ticker: string): string | undefined {
  return tickerToPerpetualMarket[ticker]?.clobPairId;
}

/**
 * Gets the perpetual market ticker given a clob pair id.
 * @param clobPairId Clob pair id to get perpetual market for.
 * @returns Ticker for the perpetual market with the clob pair id, if no perpetual market exists
 * with the given clob pair id, undefined is returned.
 */
export function getPerpetualMarketTicker(clobPairId: string): string | undefined {
  return clobPairIdToPerpetualMarket[clobPairId]?.ticker;
}

/**
 * Gets the perpetual market for a given ticker.
 */
export function getPerpetualMarketFromTicker(
  ticker: string,
): PerpetualMarketFromDatabase | undefined {
  return tickerToPerpetualMarket[ticker];
}

/**
 * Gets the perpetual market for a given clob pair id.
 */
export function getPerpetualMarketFromClobPairId(
  clobPairId: string,
): PerpetualMarketFromDatabase | undefined {
  const ticker: string | undefined = getPerpetualMarketTicker(clobPairId);
  if (ticker === undefined) {
    return undefined;
  }
  return getPerpetualMarketFromTicker(ticker);
}

/**
 * Gets the perpetual market for a given id.
 */
export function getPerpetualMarketFromId(id: string): PerpetualMarketFromDatabase | undefined {
  return idToPerpetualMarket[id];
}

export function getClobPairIdToPerpetualMarket(): Record<string, PerpetualMarketFromDatabase> {
  return clobPairIdToPerpetualMarket;
}

export function getTickerToPerpetualMarketForTest(): Record<string, PerpetualMarketFromDatabase> {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw Error(
      `getTickerToPerpetualMarketForTest cannot be used in env ${config.NODE_ENV}`);
  }

  return tickerToPerpetualMarket;
}

export function getPerpetualMarketsMap(): PerpetualMarketsMap {
  return idToPerpetualMarket;
}
