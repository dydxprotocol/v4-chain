import { logger } from '@dydxprotocol-indexer/base';
import { PartialModelObject, QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  verifyAllRequiredFields,
  setupBaseQuery,
  rawQuery,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import {
  QueryConfig,
  YieldParamsFromDatabase,
  YieldParamsQueryConfig,
  YieldParamsColumns,
  YieldParamsCreateObject,
  Options,
  Ordering,
  QueryableField,
} from '../types';
import YieldParamsModel from '../models/yield-params-model';

export function uuid(createdAtHeight: string): string {
    // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
    return getUuid(Buffer.from(`${createdAtHeight}`, BUFFER_ENCODING_UTF_8));
  }

export async function findAll(
  {
    id,
    createdAtHeight,
    createdBeforeOrAtHeight,
    createdAfterHeight,
    createdAt,
    createdBeforeOrAt,
    createdAfter,
    assetYieldIndex,
    sDAIPrice,
    limit,
  }: YieldParamsQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<YieldParamsFromDatabase[]> {
  verifyAllRequiredFields(
    {
      id,
      createdAtHeight,
      createdBeforeOrAtHeight,
      createdAfterHeight,
      createdAt,
      createdBeforeOrAt,
      createdAfter,
      assetYieldIndex,
      sDAIPrice,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<YieldParamsModel> = setupBaseQuery<YieldParamsModel>(
    YieldParamsModel,
    options,
  );

  if (id) {
    baseQuery = baseQuery.whereIn(YieldParamsColumns.id, id);
  }

  if (assetYieldIndex) {
    baseQuery = baseQuery.where(YieldParamsColumns.assetYieldIndex, assetYieldIndex)
  }

  if (sDAIPrice) {
    baseQuery = baseQuery.where(YieldParamsColumns.sDAIPrice, sDAIPrice)
  }

  if (createdAt) {
    baseQuery = baseQuery.where(YieldParamsColumns.createdAt, createdAt);
  }

  if (createdAtHeight) {
    baseQuery = baseQuery.whereIn(YieldParamsColumns.createdAtHeight, createdAtHeight);
  }

  if (createdBeforeOrAt) {
    baseQuery = baseQuery.where(YieldParamsColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdBeforeOrAtHeight) {
    baseQuery = baseQuery.where(YieldParamsColumns.createdAtHeight, '<=', createdBeforeOrAtHeight);
  }

  if (createdAfter) {
    baseQuery = baseQuery.where(YieldParamsColumns.createdAt, '>', createdAfter);
  }

  if (createdAfterHeight) {
    baseQuery = baseQuery.where(YieldParamsColumns.createdAtHeight, '>', createdAfterHeight);
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
      YieldParamsColumns.assetYieldIndex,
      Ordering.DESC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  yieldParamsToCreate: YieldParamsCreateObject,
  options: Options = { txId: undefined },
): Promise<YieldParamsFromDatabase> {
  return YieldParamsModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...yieldParamsToCreate,
    id: uuid(yieldParamsToCreate.createdAtHeight),
  }).returning('*');
}

export async function findById(
    id: string,
    options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<YieldParamsFromDatabase | undefined> {
    const baseQuery: QueryBuilder<YieldParamsModel> = setupBaseQuery<YieldParamsModel>(
      YieldParamsModel,
      options,
    );
    return baseQuery
      .findById(id)
      .returning('*');
}

export async function getLatest(
    options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<YieldParamsFromDatabase> {
    const baseQuery: QueryBuilder<YieldParamsModel> = setupBaseQuery<YieldParamsModel>(
        YieldParamsModel,
        options,
    );

    const results: YieldParamsFromDatabase[] = await baseQuery
        .orderBy(YieldParamsColumns.createdAtHeight, Ordering.DESC)
        .limit(1)
        .returning('*');

    const latestYieldParams: YieldParamsFromDatabase | undefined = results[0];
    if (latestYieldParams === undefined) {
        logger.error({
        at: 'yield-params-table#getLatest',
        message: 'Unable to find latest yield params',
        });
        throw new Error('Unable to find latest yield params');
    }
    return latestYieldParams;
}
  