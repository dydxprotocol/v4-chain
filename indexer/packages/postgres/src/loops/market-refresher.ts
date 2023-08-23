import {
  stats, delay, logger, NodeEnv,
} from '@dydxprotocol-indexer/base';

import config from '../config';
import * as MarketTable from '../stores/market-table';
import { MarketFromDatabase, MarketsMap, Options } from '../types';

let idToMarket: MarketsMap = {};

/**
 * Refresh loop to cache the list of all markets from the database in-memory.
 */
export async function start(): Promise<void> {
  for (;;) {
    await updateMarkets();
    await delay(config.MARKET_REFRESHER_INTERVAL_MS);
  }
}

/**
 * Updates in-memory map of markets.
 */
export async function updateMarkets(options?: Options): Promise<void> {
  const startTime: number = Date.now();
  const markets: MarketFromDatabase[] = await MarketTable.findAll(
    {},
    [],
    options || { readReplica: true },
  );

  const tmpIdToMarket: Record<string, MarketFromDatabase> = {};
  markets.forEach(
    (market: MarketFromDatabase) => {
      tmpIdToMarket[market.id] = market;
    },
  );

  idToMarket = tmpIdToMarket;
  stats.timing(`${config.SERVICE_NAME}.loops.update_markets`, Date.now() - startTime);
}

/**
 * Gets the market for a given id.
 */
export function getMarketFromId(id: number): MarketFromDatabase {
  const market: MarketFromDatabase | undefined = idToMarket[id];
  if (market === undefined) {
    const message: string = `Unable to find market with id: ${id}`;
    logger.error({
      at: 'market-refresher#getMarketFromId',
      message,
    });
    throw new Error(message);
  }
  return market;
}

export function getMarketsMap(): MarketsMap {
  return idToMarket;
}

export function clear(): void {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw new Error('clear cannot be used in non-test env');
  }

  idToMarket = {};
}
