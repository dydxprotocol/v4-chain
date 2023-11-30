import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import TradingRewardAggregationModel from '../models/trading-reward-aggregation-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  TradingRewardAggregationColumns,
  TradingRewardAggregationCreateObject,
  TradingRewardAggregationFromDatabase,
  TradingRewardAggregationPeriod,
  TradingRewardAggregationQueryConfig,
  TradingRewardAggregationUpdateObject,
} from '../types';

export function uuid(
  address: string,
  period: TradingRewardAggregationPeriod,
  startedAtHeight: string,
): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(Buffer.from(`${address}-${period}-${startedAtHeight}`, BUFFER_ENCODING_UTF_8));
}

export async function findAll(
  {
    address,
    startedAtHeight,
    period,
    limit,
  }: TradingRewardAggregationQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TradingRewardAggregationFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      startedAtHeight,
      period,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery:
  QueryBuilder<TradingRewardAggregationModel> = setupBaseQuery<TradingRewardAggregationModel>(
    TradingRewardAggregationModel,
    options,
  );

  if (address) {
    baseQuery = baseQuery.where(TradingRewardAggregationColumns.address, address);
  }

  if (startedAtHeight) {
    baseQuery = baseQuery.where(TradingRewardAggregationColumns.startedAtHeight, startedAtHeight);
  }

  if (period) {
    baseQuery = baseQuery.where(TradingRewardAggregationColumns.period, period);
  }

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(
        column,
        order,
      );
    }
  } else {
    baseQuery = baseQuery.orderBy(
      TradingRewardAggregationColumns.period,
      Ordering.ASC,
    ).orderBy(
      TradingRewardAggregationColumns.startedAtHeight,
      Ordering.ASC,
    ).orderBy(
      TradingRewardAggregationColumns.address,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  aggregationToCreate: TradingRewardAggregationCreateObject,
  options: Options = { txId: undefined },
): Promise<TradingRewardAggregationFromDatabase> {
  return TradingRewardAggregationModel.query(
    Transaction.get(options.txId),
  ).insert({
    id: uuid(
      aggregationToCreate.address,
      aggregationToCreate.period,
      aggregationToCreate.startedAtHeight,
    ),
    ...aggregationToCreate,
  }).returning('*');
}

export async function findById(
  address: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TradingRewardAggregationFromDatabase | undefined> {
  const baseQuery:
  QueryBuilder<TradingRewardAggregationModel> = setupBaseQuery<TradingRewardAggregationModel>(
    TradingRewardAggregationModel,
    options,
  );
  return baseQuery
    .findById(address)
    .returning('*');
}

export async function update(
  {
    ...fields
  }: TradingRewardAggregationUpdateObject,
  options: Options = { txId: undefined },
): Promise<TradingRewardAggregationFromDatabase | undefined> {
  const aggregation = await TradingRewardAggregationModel.query(
    Transaction.get(options.txId),
    // TODO fix expression typing so we dont have to use any
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ).findById(fields.id).patch(fields as any).returning('*');
  // The objection types mistakenly think the query returns an array of orders.
  return aggregation as unknown as (TradingRewardAggregationFromDatabase | undefined);
}
