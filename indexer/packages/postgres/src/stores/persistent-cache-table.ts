import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import PersistentCacheModel from '../models/persistent-cache-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  PersistentCacheColumns,
  PersistentCacheCreateObject,
  PersistentCacheFromDatabase,
  PersistentCacheQueryConfig,
} from '../types';

export async function findAll(
  {
    key,
    limit,
  }: PersistentCacheQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PersistentCacheFromDatabase[]> {
  verifyAllRequiredFields(
    {
      key,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<PersistentCacheModel> = setupBaseQuery<PersistentCacheModel>(
    PersistentCacheModel,
    options,
  );

  if (key) {
    baseQuery = baseQuery.where(PersistentCacheColumns.key, key);
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
      PersistentCacheColumns.key,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  kvToCreate: PersistentCacheCreateObject,
  options: Options = { txId: undefined },
): Promise<PersistentCacheFromDatabase> {
  return PersistentCacheModel.query(
    Transaction.get(options.txId),
  ).insert(kvToCreate).returning('*');
}

export async function upsert(
  kvToUpsert: PersistentCacheCreateObject,
  options: Options = { txId: undefined },
): Promise<PersistentCacheFromDatabase> {
  const kvs: PersistentCacheModel[] = await PersistentCacheModel.query(
    Transaction.get(options.txId),
  ).upsert(kvToUpsert).returning('*');
  // should only ever be one key value pair
  return kvs[0];
}

export async function findById(
  kv: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PersistentCacheFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PersistentCacheModel> = setupBaseQuery<PersistentCacheModel>(
    PersistentCacheModel,
    options,
  );
  return baseQuery
    .findById(kv)
    .returning('*');
}
