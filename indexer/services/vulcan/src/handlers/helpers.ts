import {
  OrderTable,
  PerpetualMarketFromDatabase,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import { subticksToPrice } from '@dydxprotocol-indexer/postgres/build/src/lib/protocol-translations';
import {
  IndexerOrder,
  IndexerOrder_ConditionType,
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

/**
 * Gets the trigger price for an order, returns undefined if the order has an unspecified condition
 * type
 * @param order
 * @param perpetualMarket
 * @returns
 */
export function getTriggerPrice(
  order: IndexerOrder,
  perpetualMarket: PerpetualMarketFromDatabase,
): string | undefined {
  if (order.conditionType !== IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED) {
    return subticksToPrice(
      order.conditionalOrderTriggerSubticks.toString(),
      perpetualMarket,
    );
  }
  return undefined;
}

export function orderSideToOrderbookSide(
  orderSide: IndexerOrder_Side,
): OrderbookSide {
  return orderSide === IndexerOrder_Side.SIDE_BUY ? OrderbookSide.BIDS : OrderbookSide.ASKS;
}
