import {
  logger,
  stats,
} from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  QUOTE_CURRENCY_ATOMIC_RESOLUTION,
} from '@dydxprotocol-indexer/postgres';
import {
  OpenOrdersCache,
  OrderbookLevels,
  OrderbookLevelsCache,
} from '@dydxprotocol-indexer/redis';
import Big from 'big.js';

import config from '../config';
import { redisClient } from '../helpers/redis';

/**
 * Instrument data on the orderbook to be used for analytics.
 */
export default async function runTask(): Promise<void> {
  const markets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll({}, []);

  for (let i: number = 0; i < markets.length; i++) {
    // Track the best bid and ask in each market
    const market: PerpetualMarketFromDatabase = markets[i];
    const uncrossedOrderbookLevels: OrderbookLevels = await OrderbookLevelsCache.getOrderBookLevels(
      market.ticker,
      redisClient,
      {
        sortSides: true,
        uncrossBook: true,
        limitPerSide: 10,
      },
    );
    statOrderbook(uncrossedOrderbookLevels, market, 'uncrossed_orderbook');

    const crossedOrderbookLevels: OrderbookLevels = await OrderbookLevelsCache.getOrderBookLevels(
      market.ticker,
      redisClient,
      {
        sortSides: true,
        uncrossBook: false,
        limitPerSide: 10,
      },
    );
    statOrderbook(crossedOrderbookLevels, market, 'crossed_orderbook');
    const openOrders: string[] = await OpenOrdersCache.getOpenOrderIds(
      market.clobPairId,
      redisClient,
    );
    stats.gauge(
      `${config.SERVICE_NAME}.open_orders_count`,
      openOrders.length,
      { clob_pair_id: market.clobPairId },
    );
  }
}

function statOrderbook(
  orderbookLevels: OrderbookLevels,
  perpetualMarket: PerpetualMarketFromDatabase,
  stat: string,
) {
  const clobPairId: string = perpetualMarket.clobPairId;
  // When querying the orderbook, the highest/best bid is first and the lowest/best ask is first.
  // Don't stat best bid if there are no bids in the orderbook
  if (orderbookLevels.bids.length > 0) {
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.best_bid_human`,
      Big(orderbookLevels.bids[0].humanPrice).toNumber(),
      { clob_pair_id: clobPairId },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.best_bid_subticks`,
      priceToSubticks(orderbookLevels.bids[0].humanPrice, perpetualMarket),
      { clob_pair_id: clobPairId },
    );
  }
  // Don't stat best ask if there are no asks in the orderbook
  if (orderbookLevels.asks.length > 0) {
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.best_ask_human`,
      Big(orderbookLevels.asks[0].humanPrice).toNumber(),
      { clob_pair_id: clobPairId },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.best_ask_subticks`,
      priceToSubticks(orderbookLevels.asks[0].humanPrice, perpetualMarket),
      { clob_pair_id: clobPairId },
    );
  }
  logger.info({
    at: 'orderbook-instrumentation#statOrderbook',
    message: `Track ${stat} for ${clobPairId}`,
    bids: JSON.stringify(orderbookLevels.bids),
    asks: JSON.stringify(orderbookLevels.asks),
  });
}

export function priceToSubticks(
  humanPrice: string,
  perpetualMarket: PerpetualMarketFromDatabase,
): number {
  return Big(humanPrice)
    .div(Big(10).pow(perpetualMarket.quantumConversionExponent))
    .div(Big(10).pow(QUOTE_CURRENCY_ATOMIC_RESOLUTION))
    .times(Big(10).pow(perpetualMarket.atomicResolution))
    .toNumber();
}
