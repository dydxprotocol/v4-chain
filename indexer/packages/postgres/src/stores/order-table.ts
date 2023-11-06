import { IndexerOrderId } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import OrderModel from '../models/order-model';
import {
  Options,
  OrderColumns,
  OrderCreateObject,
  OrderFromDatabase,
  OrderQueryConfig,
  OrderStatus,
  OrderUpdateObject,
  QueryableField,
  QueryConfig,
} from '../types';
import * as SubaccountTable from './subaccount-table';

export function uuid(
  subaccountId: string,
  clientId: string,
  clobPairId: string,
  orderFlags: string,
): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(
    Buffer.from(
      `${subaccountId}-${clientId}-${clobPairId}-${orderFlags}`,
      BUFFER_ENCODING_UTF_8,
    ),
  );
}

/**
 * Expects orderId.subaccountId, and orderId.clientId to exist.
 * @param order
 */
export function orderIdToUuid(orderId: IndexerOrderId): string {
  return uuid(
    SubaccountTable.subaccountIdToUuid(orderId.subaccountId!),
    orderId.clientId.toString(),
    orderId.clobPairId.toString(),
    orderId.orderFlags.toString(),
  );
}

export async function findAll(
  {
    limit,
    id,
    subaccountId,
    clientId,
    clobPairId,
    side,
    size,
    totalFilled,
    price,
    type,
    status,
    reduceOnly,
    orderFlags,
    goodTilBlockBeforeOrAt,
    goodTilBlockTimeBeforeOrAt,
    clientMetadata,
    triggerPrice,
  }: OrderQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<OrderFromDatabase[]> {
  verifyAllRequiredFields(
    {
      limit,
      id,
      subaccountId,
      clientId,
      clobPairId,
      side,
      size,
      totalFilled,
      price,
      type,
      status,
      reduceOnly,
      orderFlags,
      goodTilBlockBeforeOrAt,
      goodTilBlockTimeBeforeOrAt,
      clientMetadata,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<OrderModel> = setupBaseQuery<OrderModel>(
    OrderModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(OrderColumns.id, id);
  }

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(OrderColumns.subaccountId, subaccountId);
  }

  if (clientId !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.clientId, clientId);
  }

  if (clobPairId !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.clobPairId, clobPairId);
  }

  if (side !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.side, side);
  }

  if (size !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.size, size);
  }

  if (totalFilled !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.totalFilled, totalFilled);
  }

  if (price !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.price, price);
  }

  if (type !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.type, type);
  }

  if (status !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.status, status);
  }

  if (reduceOnly !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.reduceOnly, reduceOnly);
  }

  if (orderFlags !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.orderFlags, orderFlags);
  }

  if (clientMetadata !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.clientMetadata, clientMetadata);
  }

  if (triggerPrice !== undefined) {
    baseQuery = baseQuery.where(OrderColumns.triggerPrice, triggerPrice);
  }

  // If filtering by `goodTilBlock`, filter out all rows with NULL `goodTilBlock`
  if (goodTilBlockBeforeOrAt !== undefined) {
    baseQuery = baseQuery.whereNotNull(
      OrderColumns.goodTilBlock,
    ).andWhere(
      OrderColumns.goodTilBlock,
      '<=',
      goodTilBlockBeforeOrAt,
    );
  }

  // If filtering by `goodTilBlockTime`, filter out all rows with NULL `goodTilBlockTime`
  if (goodTilBlockTimeBeforeOrAt !== undefined) {
    baseQuery = baseQuery.whereNotNull(
      OrderColumns.goodTilBlockTime,
    ).andWhere(
      OrderColumns.goodTilBlockTime,
      '<=',
      goodTilBlockTimeBeforeOrAt,
    );
  }

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(
        column,
        order,
      );
    }
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  orderToCreate: OrderCreateObject,
  options: Options = { txId: undefined },
): Promise<OrderFromDatabase> {
  return OrderModel.query(
    Transaction.get(options.txId),
  ).insert({
    id: uuid(
      orderToCreate.subaccountId,
      orderToCreate.clientId,
      orderToCreate.clobPairId,
      orderToCreate.orderFlags,
    ),
    ...orderToCreate,
  }).returning('*');
}

function getOrderStatus(
  orderToUpsert: OrderCreateObject,
  totalFilled: string,
): OrderStatus {
  if (orderToUpsert.status === OrderStatus.BEST_EFFORT_CANCELED) {
    return OrderStatus.BEST_EFFORT_CANCELED;
  }
  if (Big(orderToUpsert.size).lte(totalFilled)) {
    return OrderStatus.FILLED;
  }
  return orderToUpsert.status;
}

export async function update(
  {
    ...fields
  }: OrderUpdateObject,
  options: Options = { txId: undefined },
): Promise<OrderFromDatabase | undefined> {
  const order = await OrderModel.query(
    Transaction.get(options.txId),
    // TODO fix expression typing so we dont have to use any
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ).findById(fields.id).patch(fields as any).returning('*');
  // The objection types mistakenly think the query returns an array of orders.
  return order as unknown as (OrderFromDatabase | undefined);
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<OrderFromDatabase | undefined> {
  const baseQuery: QueryBuilder<OrderModel> = setupBaseQuery<OrderModel>(
    OrderModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function findBySubaccountIdAndClobPair(
  subaccountId: string,
  clobPairId: string,
  options: Options = {},
): Promise<OrderFromDatabase[]> {
  const baseQuery: QueryBuilder<OrderModel> = setupBaseQuery<OrderModel>(
    OrderModel,
    options,
  );

  const orders: OrderFromDatabase[] = await baseQuery
    .where(OrderColumns.subaccountId, subaccountId)
    .where(OrderColumns.clobPairId, clobPairId)
    .returning('*');
  return orders;
}

export async function upsert(
  orderToUpsert: OrderCreateObject,
  options: Options = { txId: undefined },
): Promise<OrderFromDatabase> {
  const orderId: string = uuid(
    orderToUpsert.subaccountId,
    orderToUpsert.clientId,
    orderToUpsert.clobPairId,
    orderToUpsert.orderFlags,
  );

  const order: OrderFromDatabase | undefined = await findById(orderId, options);
  if (order === undefined) {
    return create({
      ...orderToUpsert,
      status: getOrderStatus(orderToUpsert, orderToUpsert.totalFilled),
    }, options);
  }

  const updatedOrder: OrderFromDatabase | undefined = await update({
    ...orderToUpsert,
    status: getOrderStatus(orderToUpsert, orderToUpsert.totalFilled),
    id: orderId,
  }, options);

  if (updatedOrder === undefined) {
    throw Error('order must exist after update');
  }

  return updatedOrder;
}

/**
 * Checks if the order is a long term or conditional order.
 * @param orderFlags
 */
export function isLongTermOrConditionalOrder(orderFlags: string): boolean {
  const flags: number = parseInt(orderFlags, 10);

  const isLongTerm: boolean = Math.floor(flags / 64) % 2 === 1;
  const isConditional: boolean = Math.floor(flags / 32) % 2 === 1;

  return isLongTerm || isConditional;
}

/**
 * Finds all open long term or conditional orders.
 * @param options
 */
export async function findOpenLongTermOrConditionalOrders(
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<OrderFromDatabase[]> {
  const baseQuery: QueryBuilder<OrderModel> = setupBaseQuery<OrderModel>(
    OrderModel,
    options,
  );

  /* eslint-disable */
  return baseQuery
    .where(OrderColumns.status, OrderStatus.OPEN)
    .andWhere(function () {
      this.whereRaw('FLOOR("orderFlags"::integer / 64)::integer % 2 = 1')
        .orWhereRaw('FLOOR("orderFlags"::integer / 32)::integer % 2 = 1');
    })
    .returning('*');
}

export async function updateStaleOrderStatusByIds(
  oldStatus: OrderStatus,
  newStatus: OrderStatus,
  latestBlockHeight: string,
  ids: string[],
  options: Options = { txId: undefined },
): Promise<OrderFromDatabase[]> {
  const updatedOrders: OrderFromDatabase[] = await OrderModel
    .query(Transaction.get(options.txId))
    .where(OrderColumns.status, oldStatus)
    .whereNotNull(OrderColumns.goodTilBlock)
    .where(OrderColumns.goodTilBlock, '<', latestBlockHeight)
    .whereIn(OrderColumns.id, ids)
    .patch({ status: newStatus })
    .returning('*');

  return updatedOrders;
}
