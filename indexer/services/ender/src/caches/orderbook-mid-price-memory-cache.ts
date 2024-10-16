import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  loopHelpers,
} from '@dydxprotocol-indexer/postgres';
import { OrderbookMidPricesCache } from '@dydxprotocol-indexer/redis';

import config from '../config';
import { redisClient } from '../helpers/redis/redis-controller';

interface OrderbookMidPriceCache {
  [ticker: string]: string | undefined,
}

let orderbookMidPriceCache: OrderbookMidPriceCache = {};

/**
 * Refresh loop to cache the list of all perpetual markets from the database in-memory.
 */
export async function start(): Promise<void> {
  await loopHelpers.startUpdateLoop(
    updateOrderbookMidPrices,
    config.ORDERBOOK_MID_PRICE_REFRESH_INTERVAL_MS,
    'updateOrderbookMidPrices',
  );
}

export function getOrderbookMidPrice(ticker: string): string | undefined {
  return orderbookMidPriceCache[ticker];
}

export async function updateOrderbookMidPrices(): Promise<void> {
  const startTime: number = Date.now();
  try {
    const perpetualMarkets: PerpetualMarketFromDatabase[] = Object.values(
      perpetualMarketRefresher.getPerpetualMarketsMap(),
    );

    const tickers: string[] = perpetualMarkets.map((market) => market.ticker);

    orderbookMidPriceCache = await OrderbookMidPricesCache.getMedianPrices(
      redisClient,
      tickers,
    );

    // Log out each median price for each market
    Object.entries(orderbookMidPriceCache).forEach(([ticker, price]) => {
      logger.info({
        at: 'orderbook-mid-price-cache#updateOrderbookMidPrices',
        message: `Median price for market ${ticker}`,
        ticker,
        price,
      });
    });

  } catch (error) {
    logger.error({
      at: 'orderbook-mid-price-cache#updateOrderbookMidPrices',
      message: 'Failed to fetch OrderbookMidPrices',
      error,
    });
  } finally {
    stats.timing(
      `${config.SERVICE_NAME}.update_orderbook_mid_prices_cache.timing`,
      Date.now() - startTime,
    );
  }
}
