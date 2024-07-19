import {
  IndexerOrder,
  IndexerOrderId,
  RedisOrder,
} from '@dydxprotocol-indexer/v4-protos';

import { InvalidRedisOrderError } from '../errors';

export function validateRedisOrder(redisOrder: RedisOrder): void {
  if (redisOrder.order === undefined) {
    throw new InvalidRedisOrderError('Order proto cannot be undefined.');
  }

  if (redisOrder.order.orderId === undefined) {
    throw new InvalidRedisOrderError('Order id in Order proto cannot be undefined.');
  }

  validateOrderId(redisOrder.order.orderId);
}

export function validateOrderId(orderId: IndexerOrderId): void {
  if (orderId.subaccountId === undefined) {
    throw new InvalidRedisOrderError('Subaccount id in Order proto cannot be undefined.');
  }
}

export function getOrderExpiry(order: IndexerOrder): number {
  // Protocol guarantees that an order can only have `goodTilBlock` or `goodTilBlockTime` and will
  // never be replaced with another order that has a different `goodTilOneof`.
  if (order.goodTilBlock !== undefined) {
    return order.goodTilBlock;
  } else if (order.goodTilBlockTime !== undefined) {
    return order.goodTilBlockTime;
  } else {
    throw new InvalidRedisOrderError('Order proto has inavlid goodTilOneOf field');
  }
}
