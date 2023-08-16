import { OrderTable, SubaccountTable } from '@dydxprotocol-indexer/postgres';
import { ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrderId,
  RedisOrder,
  IndexerSubaccountId,
} from '@dydxprotocol-indexer/v4-protos';
import { RedisClient } from 'redis';

import { getOrderExpiry } from '../../src/caches/helpers';
import { ORDER_EXPIRY_CACHE_KEY } from '../../src/caches/order-expiry-cache';
import { getOrder, ORDERS_CACHE_KEY_PREFIX } from '../../src/caches/orders-cache';
import { getOrderDataCacheKey } from '../../src/caches/orders-data-cache';
import { getOrderIdsForSubaccount } from '../../src/caches/subaccount-order-ids-cache';
import { getAsync, zScoreAsync } from '../../src/helpers/redis';
import { PlaceOrderResult, RemoveOrderResult, UpdateOrderResult } from '../../src/types';
import { redis as client } from '../helpers/utils';

export async function expectOrderCache(
  result: PlaceOrderResult,
  redisOrder: RedisOrder,
  numOrders: number,
): Promise<void> {
  const actualOrder: RedisOrder | null = await getOrderByOrderId(
    redisOrder.order!.orderId!,
    client,
  );
  expect(actualOrder).toEqual(redisOrder);

  // If no order was replaced, total filled should be 0 for the order
  const totalFilled: number = result.replaced ? result.oldTotalFilledQuantums! : 0;
  const actualOrderData: string | null = await getOrderData(redisOrder.order!.orderId!, client);
  expect(actualOrderData).toEqual(
    // when placed, an order is never resting on the book
    `${getOrderExpiry(redisOrder.order!)}_${totalFilled}_false`,
  );

  const subaccountOrders: string[] = await getOrderIdsForSubaccountBySubaccountId(
    redisOrder.order!.orderId!.subaccountId!,
    client,
  );
  expect(subaccountOrders).toHaveLength(numOrders);
  expect(subaccountOrders).toContain(redisOrder.id);
}

export function expectNewOrder(
  result: PlaceOrderResult,
): void {
  expect(result.placed).toEqual(true);
  expect(result.replaced).toEqual(false);
  expect(result.oldTotalFilledQuantums).toBeUndefined();
  expect(result.restingOnBook).toBeUndefined();
  expect(result.oldOrder).toBeUndefined();
}

export async function expectOrderExpiry(order: RedisOrder): Promise<void> {
  const orderUuid: string = order.id;
  const isShortTermOrder: boolean = order.order!.orderId!.orderFlags === ORDER_FLAG_SHORT_TERM;
  const expectedExpiry: string = getOrderExpiry(order.order!).toString();

  const expiry: string | null = await getOrderExpiryFromCache(orderUuid, client);
  if (isShortTermOrder) {
    expect(expiry).toEqual(expectedExpiry);
  } else {
    expect(expiry).toBeNull();
  }
}

export function expectRemovedOrderResult({
  result,
  removed,
  totalFilledQuantums,
  restingOnBook,
  removedOrder,
}:{
  result: RemoveOrderResult,
  removed: boolean,
  totalFilledQuantums?: number,
  restingOnBook?: boolean,
  removedOrder?: RedisOrder,
},
): void {
  expect(result.removed).toEqual(removed);
  expect(result.totalFilledQuantums).toEqual(totalFilledQuantums);
  expect(result.restingOnBook).toEqual(restingOnBook);
  expect(result.removedOrder).toEqual(removedOrder);
}

export async function expectRemovedOrder(
  removedOrderId: IndexerOrderId,
): Promise<void> {
  const orderUuid: string = OrderTable.orderIdToUuid(removedOrderId);
  const nonexistentOrder: RedisOrder | null = await getOrderByOrderId(removedOrderId, client);
  const nonexistentOrderData: string | null = await getOrderData(removedOrderId, client);
  const subaccountOrderIds: string[] = await getOrderIdsForSubaccount(
    SubaccountTable.subaccountIdToUuid(removedOrderId.subaccountId!),
    client,
  );
  const nonexistentExpiry: string|null = await zScoreAsync(
    { hash: ORDER_EXPIRY_CACHE_KEY, key: orderUuid },
    client,
  );
  expect(nonexistentOrder).toEqual(null);
  expect(nonexistentOrderData).toEqual(null);
  expect(subaccountOrderIds).not.toContain(orderUuid);
  expect(nonexistentExpiry).toBeNull();
}

export function expectUpdateOrderResult({
  result,
  updated,
  oldTotalFilledQuantums,
  oldRestingOnBook,
  order,
}: {
  result: UpdateOrderResult,
  updated: boolean,
  oldTotalFilledQuantums?: number,
  oldRestingOnBook?: boolean,
  order?: RedisOrder,
}): void {
  expect(result.updated).toEqual(updated);
  expect(result.oldTotalFilledQuantums).toEqual(oldTotalFilledQuantums);
  expect(result.oldRestingOnBook).toEqual(oldRestingOnBook);
  expect(result.order).toEqual(order);
}

export async function getOrderByOrderId(
  orderId: IndexerOrderId,
  redisClient: RedisClient,
): Promise<RedisOrder | null> {
  return getOrder(OrderTable.orderIdToUuid(orderId), redisClient);
}

export async function getOrderData(
  orderId: IndexerOrderId,
  redisClient: RedisClient,
): Promise<string | null> {
  return getAsync(getOrderDataCacheKey(orderId), redisClient);
}

export async function getOrderExpiryFromCache(
  orderId: string,
  redisClient: RedisClient,
): Promise<string | null> {
  return zScoreAsync({ hash: ORDER_EXPIRY_CACHE_KEY, key: orderId }, redisClient);
}

export async function getOrderIdsForSubaccountBySubaccountId(
  subaccountId: IndexerSubaccountId,
  redisClient: RedisClient,
): Promise<string[]> {
  return getOrderIdsForSubaccount(SubaccountTable.subaccountIdToUuid(subaccountId), redisClient);
}

export function getOrderCacheKey(orderId: IndexerOrderId): string {
  return getOrderCacheKeyWithUUID(OrderTable.orderIdToUuid(orderId));
}

export function getOrderCacheKeyWithUUID(orderUuid: string): string {
  return `${ORDERS_CACHE_KEY_PREFIX}${orderUuid}`;
}
