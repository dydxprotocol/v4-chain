import {
  stats,
  NodeEnv,
} from '@dydxprotocol-indexer/base';
import _ from 'lodash';

import config from '../config';
import * as PerpetualMarketTable from '../stores/perpetual-market-table';
import {
  Options, PerpetualMarketColumns, PerpetualMarketFromDatabase, PerpetualMarketsMap,
} from '../types';
import { startUpdateLoop } from './loopHelper';

let idToPerpetualMarket: Record<string, PerpetualMarketFromDatabase> = {};

// TODO(DEC-642): Update the in-memory mapping of perpetual market tickers to ids from websocket
// messages from `v4_markets` topic rather than periodically refereshing from the database.

/**
 * Refresh loop to cache the list of all perpetual markets from the database in-memory.
 */
export async function start(): Promise<void> {
  await startUpdateLoop(
    updatePerpetualMarkets,
    config.PERPETUAL_MARKETS_REFRESHER_INTERVAL_MS,
    'updatePerpetualMarkets',
  );
}

/**
 * Clears the in-memory map of perpetual market clob pair ids to tickers and tickers to clob pair.
 * Used for testing.
 */
export function clear(): void {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw new Error('clear cannot be used in non-test env');
  }
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

  const tmpIdToPerpetualMarket: Record<string, PerpetualMarketFromDatabase> = {};
  perpetualMarkets.forEach(
    (market: PerpetualMarketFromDatabase) => {
      tmpIdToPerpetualMarket[market.id] = market;
    },
  );

  idToPerpetualMarket = tmpIdToPerpetualMarket;
  stats.timing(`${config.SERVICE_NAME}.loops.update_perpetual_markets`, Date.now() - startTime);
}

/**
 * Validates a ticker references a perpetual market.
 * @param ticker Ticker to validate.
 * @returns true if ticker matches a perpetual market ticker, false otherwise.
 */
export function isValidPerpetualMarketTicker(ticker: string): boolean {
  return _.some(idToPerpetualMarket, (perpetualMarket: PerpetualMarketFromDatabase) => {
    return perpetualMarket.ticker === ticker;
  });
}

export function getPerpetualMarketsList(): PerpetualMarketFromDatabase[] {
  return Object.values(idToPerpetualMarket);
}

/**
 * Gets the clob pair id for a given perpetual market ticker.
 * @param ticker Ticker for the perpetual market.
 * @returns Clob pair id to get perpetual market with the ticker, if no perpetual market exists
 * with the given clob pair id, undefined is returned.
 */
export function getClobPairIdFromTicker(ticker: string): string | undefined {
  const perpetualMarket: PerpetualMarketFromDatabase | undefined = getPerpetualMarketFromTicker(
    ticker,
  );

  return perpetualMarket?.clobPairId;
}

/**
 * Gets the perpetual market ticker given a clob pair id.
 * @param clobPairId Clob pair id to get perpetual market for.
 * @returns Ticker for the perpetual market with the clob pair id, if no perpetual market exists
 * with the given clob pair id, undefined is returned.
 */
export function getPerpetualMarketTicker(clobPairId: string): string | undefined {
  const perpetualMarket: PerpetualMarketFromDatabase | undefined = getPerpetualMarketFromClobPairId(
    clobPairId,
  );

  return perpetualMarket?.ticker;
}

/**
 * Gets the perpetual market for a given ticker.
 */
export function getPerpetualMarketFromTicker(
  ticker: string,
): PerpetualMarketFromDatabase | undefined {
  return _.find(
    getPerpetualMarketsList(),
    (perpetualMarket: PerpetualMarketFromDatabase) => {
      return perpetualMarket.ticker === ticker;
    },
  );
}

/**
 * Gets the perpetual market for a given clob pair id.
 */
export function getPerpetualMarketFromClobPairId(
  clobPairId: string,
): PerpetualMarketFromDatabase | undefined {
  return _.find(
    getPerpetualMarketsList(),
    (perpetualMarket: PerpetualMarketFromDatabase) => {
      return perpetualMarket.clobPairId === clobPairId;
    },
  );
}

/**
 * Gets the perpetual market for a given id.
 */
export function getPerpetualMarketFromId(id: string): PerpetualMarketFromDatabase | undefined {
  return idToPerpetualMarket[id];
}

export function getClobPairIdToPerpetualMarket(): Record<string, PerpetualMarketFromDatabase> {
  return _.keyBy(getPerpetualMarketsList(), PerpetualMarketColumns.clobPairId);
}

export function getTickerToPerpetualMarketForTest(): Record<string, PerpetualMarketFromDatabase> {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw Error(
      `getTickerToPerpetualMarketForTest cannot be used in env ${config.NODE_ENV}`);
  }

  return _.keyBy(getPerpetualMarketsList(), PerpetualMarketColumns.ticker);
}

export function getPerpetualMarketsMap(): PerpetualMarketsMap {
  return idToPerpetualMarket;
}

/**
 * Add or updates a perpetual market instance in the in memory cache
 */
export function upsertPerpetualMarket(perpetualMarket: PerpetualMarketFromDatabase): void {
  idToPerpetualMarket[perpetualMarket.id] = perpetualMarket;
}
