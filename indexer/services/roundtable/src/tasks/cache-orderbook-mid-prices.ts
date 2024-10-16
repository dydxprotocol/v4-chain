import {
  logger,
} from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketTable,
} from '@dydxprotocol-indexer/postgres';
import {
  OrderbookMidPricesCache,
} from '@dydxprotocol-indexer/redis';

import { redisClient } from '../helpers/redis';

/**
 * Updates OrderbookMidPricesCache with current orderbook mid price for each market
 */
export default async function runTask(): Promise<void> {
  const marketTickers: string[] = (await PerpetualMarketTable.findAll(
    {},
    [],
    { readReplica: true },
  )).map((market) => {
    return market.ticker;
  });

  try {
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
