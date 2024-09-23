import {
  logger,
} from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
} from '@dydxprotocol-indexer/postgres';
import {
  OrderbookMidPricesCache,
  OrderbookLevelsCache,
} from '@dydxprotocol-indexer/redis';

import { redisClient } from '../helpers/redis';

/**
 * Updates OrderbookMidPricesCache with current orderbook mid price for each market
 */
export default async function runTask(): Promise<void> {
  const markets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll({}, []);

  for (const market of markets) {
    try {
      const price = await OrderbookLevelsCache.getOrderBookMidPrice(market.ticker, redisClient);
      if (price) {
        await OrderbookMidPricesCache.setPrice(redisClient, market.ticker, price);
      } else {
        logger.info({
          at: 'cache-orderbook-mid-prices#runTask',
          message: `undefined price for ${market.ticker}`,
        });
      }
    } catch (error) {
      logger.error({
        at: 'cache-orderbook-mid-prices#runTask',
        message: error.message,
        error,
      });
    }
  }
}
