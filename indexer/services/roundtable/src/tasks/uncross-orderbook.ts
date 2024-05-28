import { logger } from '@dydxprotocol-indexer/base';
import { OrderSide, PerpetualMarketFromDatabase, PerpetualMarketTable } from '@dydxprotocol-indexer/postgres';
import { OrderbookLevels, OrderbookLevelsCache } from '@dydxprotocol-indexer/redis';
import { deleteStalePriceLevel } from '@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache';
import Big from 'big.js';

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

  while (
    ai < orderbookLevels.asks.length &&
    bi < orderbookLevels.bids.length &&
    Big(orderbookLevels.bids[bi].humanPrice).gte(Big(orderbookLevels.asks[ai].humanPrice))
  ) {
    // Compare the recency and size of the bid and ask to decide which to remove
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

  for (const bid of removeBidLevels) {
    const deleted: boolean = await deleteStalePriceLevel({
      ticker,
      side: OrderSide.BUY,
      humanPrice: bid.humanPrice,
      client: redisClient,
    });
    if (!deleted) {
      logger.info({
        at: 'uncrossOrderbook#deleteStalePriceLevel',
        message: `Failed to delete stale bid level for ${ticker}`,
        side: OrderSide.BUY,
        humanPrice: bid.humanPrice,
      });
    }
  }

  for (const ask of removeAskLevels) {
    const deleted: boolean = await deleteStalePriceLevel({
      ticker,
      side: OrderSide.SELL,
      humanPrice: ask.humanPrice,
      client: redisClient,
    });
    if (!deleted) {
      logger.info({
        at: 'uncrossOrderbook#deleteStalePriceLevel',
        message: `Failed to delete stale ask level for ${ticker}`,
        side: OrderSide.SELL,
        humanPrice: ask.humanPrice,
      });
    }
  }

  logger.info({
    at: 'uncrossOrderbook#uncrossOrderbook',
    message: `Uncrossed orderbook for ${ticker}`,
    removedBids: JSON.stringify(removeBidLevels),
    removedAsks: JSON.stringify(removeAskLevels),
  });
}
