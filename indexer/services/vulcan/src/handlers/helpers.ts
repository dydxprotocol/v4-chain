import { StateFilledQuantumsCache } from '@dydxprotocol-indexer/redis';
import { IndexerOrder_Side, RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';

import { redisClient } from '../helpers/redis/redis-controller';
import { OrderbookSide } from '../lib/types';

export function orderSideToOrderbookSide(
  orderSide: IndexerOrder_Side,
): OrderbookSide {
  return orderSide === IndexerOrder_Side.SIDE_BUY ? OrderbookSide.BIDS : OrderbookSide.ASKS;
}

/**
 * Gets the remaining quantums for an order based on the filled amount of the order in state
 * @param order
 * @returns
 */
export async function getStateRemainingQuantums(
  order: RedisOrder,
): Promise<Big> {
  const orderQuantums: Big = Big(order.order!.quantums.toString());
  const stateFilledQuantums: Big = convertToBig(
    await StateFilledQuantumsCache.getStateFilledQuantums(order.id, redisClient),
  );
  return orderQuantums.minus(stateFilledQuantums);
}

function convertToBig(value: string | undefined) {
  if (value === undefined) {
    return Big(0);
  } else {
    return Big(value);
  }
}
