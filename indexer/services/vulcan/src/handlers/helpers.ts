import {
  OrderTable,
  PerpetualMarketFromDatabase,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerOrder,
  IndexerOrder_Side,
  RedisOrder,
  RedisOrder_TickerType,
} from '@dydxprotocol-indexer/v4-protos';

import { OrderbookSide } from '../lib/types';

/**
 * Creates a `RedisOrder` given an `Order` and the corresponding `PerpetualMarket` for the `Order`.
 * @param order
 * @param perpetualMarket
 * @returns
 */
export function convertToRedisOrder(
  order: IndexerOrder,
  perpetualMarket: PerpetualMarketFromDatabase,
): RedisOrder {
  return {
    order,
    id: OrderTable.orderIdToUuid(order.orderId!),
    ticker: perpetualMarket.ticker,
    tickerType: RedisOrder_TickerType.TICKER_TYPE_PERPETUAL,
    price: protocolTranslations.subticksToPrice(
      order.subticks.toString(),
      perpetualMarket,
    ),
    size: protocolTranslations.quantumsToHumanFixedString(
      order.quantums.toString(),
      perpetualMarket.atomicResolution,
    ),
  };
}

export function orderSideToOrderbookSide(
  orderSide: IndexerOrder_Side,
): OrderbookSide {
  return orderSide === IndexerOrder_Side.SIDE_BUY ? OrderbookSide.BIDS : OrderbookSide.ASKS;
}
