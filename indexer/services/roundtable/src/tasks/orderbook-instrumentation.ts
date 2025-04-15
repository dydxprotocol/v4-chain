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
      },
    );
    statOrderbook(uncrossedOrderbookLevels, market, 'uncrossed_orderbook');

    const crossedOrderbookLevels: OrderbookLevels = await OrderbookLevelsCache.getOrderBookLevels(
      market.ticker,
      redisClient,
      {
        sortSides: true,
        uncrossBook: false,
      },
    );
    statOrderbook(crossedOrderbookLevels, market, 'crossed_orderbook');
  }
}

function statOrderbook(
  orderbookLevels: OrderbookLevels,
  perpetualMarket: PerpetualMarketFromDatabase,
  stat: string,
) {
  const ticker: string = perpetualMarket.ticker;
  // When querying the orderbook, the highest/best bid is first and the lowest/best ask is first.
  // Don't stat best bid if there are no bids in the orderbook
  if (orderbookLevels.bids.length > 0) {
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.best_bid_human`,
      Big(orderbookLevels.bids[0].humanPrice).toNumber(),
      { ticker },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.best_bid_subticks`,
      priceToSubticks(orderbookLevels.bids[0].humanPrice, perpetualMarket),
      { ticker },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.num_bid_levels`,
      orderbookLevels.bids.length,
      { ticker },
    );
  }
  // Don't stat best ask if there are no asks in the orderbook
  if (orderbookLevels.asks.length > 0) {
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.best_ask_human`,
      Big(orderbookLevels.asks[0].humanPrice).toNumber(),
      { ticker },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.best_ask_subticks`,
      priceToSubticks(orderbookLevels.asks[0].humanPrice, perpetualMarket),
      { ticker },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${stat}.num_ask_levels`,
      orderbookLevels.asks.length,
      { ticker },
    );
  }
  logger.info({
    at: 'orderbook-instrumentation#statOrderbook',
    message: `Track ${stat} for ${ticker}`,
    bids: JSON.stringify(orderbookLevels.bids.slice(0, 10)),
    asks: JSON.stringify(orderbookLevels.asks.slice(0, 10)),
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
