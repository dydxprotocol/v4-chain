import {
  logger,
} from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import {
  OrderbookMidPricesCache,
} from '@dydxprotocol-indexer/redis';

import { redisClient } from '../helpers/redis';

/**
 * Updates OrderbookMidPricesCache with current orderbook mid price for each market
 */
export default async function runTask(): Promise<void> {
  try {
    const perpetualMarkets = Object.values(perpetualMarketRefresher.getPerpetualMarketsMap());
    const marketTickers = perpetualMarkets.map(
      (market: PerpetualMarketFromDatabase) => market.ticker,
    );

    if (marketTickers.length === 0) {
      throw new Error('perpetualMarketRefresher is empty');
    }

    logger.info({
      at: 'cache-orderbook-mid-prices#runTask',
      message: 'Caching orderbook mid prices for markets',
      markets: marketTickers.join(', '),
    });
    await OrderbookMidPricesCache.fetchAndCacheOrderbookMidPrices(redisClient, marketTickers);
  } catch (error) {
    logger.error({
      at: 'cache-orderbook-mid-prices#runTask',
      message: error.message,
      error,
    });
  }
}
