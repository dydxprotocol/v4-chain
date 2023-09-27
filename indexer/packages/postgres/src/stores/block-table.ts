import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import BlockModel from '../models/block-model';
import {
  BlockColumns,
  BlockCreateObject,
  BlockFromDatabase,
  BlockQueryConfig,
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
} from '../types';

export async function findAll(
  {
    blockHeight,
    createdOnOrAfter,
    limit,
  }: BlockQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<BlockFromDatabase[]> {
  verifyAllRequiredFields(
    {
      blockHeight,
      createdOnOrAfter,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<BlockModel> = setupBaseQuery<BlockModel>(
    BlockModel,
    options,
  );

  if (blockHeight !== undefined) {
    baseQuery = baseQuery.whereIn(BlockColumns.blockHeight, blockHeight);
  }

  if (createdOnOrAfter !== undefined) {
    baseQuery = baseQuery.where(BlockColumns.time, '>=', createdOnOrAfter);
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
      BlockColumns.blockHeight,
      Ordering.ASC,
    );
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  blockToCreate: BlockCreateObject,
  options: Options = { txId: undefined },
): Promise<BlockFromDatabase> {
  return BlockModel.query(
    Transaction.get(options.txId),
  ).insert(blockToCreate).returning('*');
}

export async function findByBlockHeight(
  blockHeight: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<BlockFromDatabase | undefined> {
  const baseQuery: QueryBuilder<BlockModel> = setupBaseQuery<BlockModel>(
    BlockModel,
    options,
  );
  return baseQuery
    .findById(blockHeight)
    .returning('*');
}

export async function getLatest(
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<BlockFromDatabase | undefined> {
  const baseQuery: QueryBuilder<BlockModel> = setupBaseQuery<BlockModel>(
    BlockModel,
    options,
  );

  const results: BlockModel[] = await baseQuery
    .orderBy(BlockColumns.blockHeight, Ordering.DESC)
    .limit(1)
    .returning('*');

  return results[0];
}
