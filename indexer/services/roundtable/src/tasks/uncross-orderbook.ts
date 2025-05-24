import { logger, stats } from '@dydxprotocol-indexer/base';
import { OrderSide, PerpetualMarketFromDatabase, PerpetualMarketTable } from '@dydxprotocol-indexer/postgres';
import { OrderbookLevels, OrderbookLevelsCache } from '@dydxprotocol-indexer/redis';
import Big from 'big.js';

import config from '../config';
import { redisClient } from '../helpers/redis';

/**
 * Task to uncross the orderbook.
 */
export default async function runTask(): Promise<void> {
  const markets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll({}, []);

  for (let i: number = 0; i < markets.length; i++) {
    const market: PerpetualMarketFromDatabase = markets[i];

    const crossedOrderbookLevels: OrderbookLevels = await OrderbookLevelsCache.getOrderBookLevels(
      market.ticker,
      redisClient,
      {
        removeZeros: true,
        sortSides: true,
        uncrossBook: false,
      },
    );

    // Check if the orderbook is crossed
    if (isOrderbookCrossed(crossedOrderbookLevels)) {
      stats.increment(`${config.SERVICE_NAME}.crossed_orderbook`, { ticker: market.ticker });
      await uncrossOrderbook(market, crossedOrderbookLevels);
    }
  }
}

/**
 * Check if the orderbook is crossed.
 */
function isOrderbookCrossed(orderbookLevels: OrderbookLevels): boolean {
  if (orderbookLevels.bids.length > 0 && orderbookLevels.asks.length > 0) {
    const bestBid = Big(orderbookLevels.bids[0].humanPrice);
    const bestAsk = Big(orderbookLevels.asks[0].humanPrice);
    return bestBid.gte(bestAsk);
  }
  return false;
}

/**
 * Uncross the orderbook by removing overlapping bid and ask levels.
 * Remove the price levels that are older.
 */
async function uncrossOrderbook(
  market: PerpetualMarketFromDatabase,
  orderbookLevels: OrderbookLevels,
): Promise<void> {
  const ticker = market.ticker;

  // Remove overlapping levels
  let ai = 0;
  let bi = 0;

  // the bids are sorted in descending order
  // the asks are sorted in ascending order
  while (
    ai < orderbookLevels.asks.length &&
    bi < orderbookLevels.bids.length &&
    Big(orderbookLevels.bids[bi].humanPrice).gte(Big(orderbookLevels.asks[ai].humanPrice))
  ) {
    // Remove the older side. If the ask and bid levels have the same lastUpdated time,
    // remove the bid level.
    if (
      Number(orderbookLevels.bids[bi].lastUpdated) > Number(orderbookLevels.asks[ai].lastUpdated)
    ) {
      ai += 1;
    } else {
      bi += 1;
    }
  }

  // Remove crossed levels from Redis
  const removeBidLevels = orderbookLevels.bids.slice(0, bi);
  const removeAskLevels = orderbookLevels.asks.slice(0, ai);

  logger.info({
    at: 'uncrossOrderbook#uncrossOrderbook',
    message: `Uncrossing orderbook for ${ticker}`,
    removedBids: JSON.stringify(removeBidLevels),
    removedAsks: JSON.stringify(removeAskLevels),
  });

  stats.increment(
    `${config.SERVICE_NAME}.expected_uncross_orderbook_levels`,
    removeBidLevels.length,
    { side: OrderSide.BUY, clobPairId: market.clobPairId, ticker },
  );
  for (const bid of removeBidLevels) {
    const deleted: boolean = await OrderbookLevelsCache.deleteStalePriceLevel(
      ticker,
      OrderSide.BUY,
      bid.humanPrice,
      config.STALE_ORDERBOOK_LEVEL_THRESHOLD_SECONDS,
      redisClient,
    );
    if (!deleted) {
      stats.increment(
        `${config.SERVICE_NAME}.uncross_orderbook_failed`,
        {
          side: OrderSide.BUY,
          clobPairId: market.clobPairId,
          ticker,
        },
      );
      logger.info({
        at: 'uncrossOrderbook#deleteStalePriceLevel',
        message: `Failed to delete stale bid level for ${ticker}`,
        side: OrderSide.BUY,
        humanPrice: bid.humanPrice,
        ticker,
      });
    } else {
      stats.increment(
        `${config.SERVICE_NAME}.uncross_orderbook_succeed`,
        {
          side: OrderSide.BUY,
          clobPairId: market.clobPairId,
          ticker,
        },
      );
    }
  }

  stats.increment(
    `${config.SERVICE_NAME}.expected_uncross_orderbook_levels`,
    removeAskLevels.length,
    { side: OrderSide.SELL, clobPairId: market.clobPairId, ticker },
  );
  for (const ask of removeAskLevels) {
    const deleted: boolean = await OrderbookLevelsCache.deleteStalePriceLevel(
      ticker,
      OrderSide.SELL,
      ask.humanPrice,
      config.STALE_ORDERBOOK_LEVEL_THRESHOLD_SECONDS,
      redisClient,
    );
    if (!deleted) {
      stats.increment(
        `${config.SERVICE_NAME}.uncross_orderbook_failed`,
        {
          side: OrderSide.SELL,
          clobPairId: market.clobPairId,
          ticker,
        },
      );
      logger.info({
        at: 'uncrossOrderbook#deleteStalePriceLevel',
        message: `Failed to delete stale ask level for ${ticker}`,
        side: OrderSide.SELL,
        humanPrice: ask.humanPrice,
        ticker,
      });
    } else {
      stats.increment(
        `${config.SERVICE_NAME}.uncross_orderbook_succeed`,
        {
          side: OrderSide.SELL,
          clobPairId: market.clobPairId,
          ticker,
        },
      );
    }
  }
}
