import {
  stats,
} from '@dydxprotocol-indexer/base';
import {
  OrderSide,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
} from '@dydxprotocol-indexer/postgres';
import { OrderbookLevels, OrderbookLevelsCache, PriceLevel } from '@dydxprotocol-indexer/redis';
import _ from 'lodash';

import config from '../config';
import { redisClient } from '../helpers/redis';

export default async function runTask(): Promise<void> {
  let numDeletedZeroPriceLevels: number = 0;

  const perpMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll({}, []);
  for (const perpetualMarket of perpMarkets) {
    const priceLevels: OrderbookLevels = await OrderbookLevelsCache.getOrderBookLevels(
      perpetualMarket.ticker,
      redisClient,
      {
        removeZeros: false,
      },
    );

    const zeroBidLevels: PriceLevel[] = priceLevels.bids.filter(
      (level: PriceLevel): boolean => level.quantums === '0',
    );
    const zeroAskLevels: PriceLevel[] = priceLevels.asks.filter(
      (level: PriceLevel): boolean => level.quantums === '0',
    );

    const removedLevels: boolean [] = await Promise.all(_.flatten([
      _.map(zeroBidLevels, (zeroBidLevel: PriceLevel): Promise<boolean> => {
        return OrderbookLevelsCache.deleteZeroPriceLevel(
          perpetualMarket.ticker,
          OrderSide.BUY,
          zeroBidLevel.humanPrice,
          redisClient,
        );
      }),
      _.map(zeroAskLevels, (zeroAskLevel: PriceLevel): Promise<boolean> => {
        return OrderbookLevelsCache.deleteZeroPriceLevel(
          perpetualMarket.ticker,
          OrderSide.SELL,
          zeroAskLevel.humanPrice,
          redisClient,
        );
      }),
    ]));

    numDeletedZeroPriceLevels += _.filter(
      removedLevels,
      (removed: boolean): boolean => removed === true,
    ).length;
  }

  stats.gauge(
    `${config.SERVICE_NAME}.delete_zero_price_levels.num_levels_deleted`,
    numDeletedZeroPriceLevels,
  );
}
