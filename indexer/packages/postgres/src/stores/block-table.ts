import { logger } from '@dydxprotocol-indexer/base';
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

/**
 * Find the first block created on or after the given timestamp.
 * Uses the blocks_time_since_march2025_idx index for efficient querying.
 */
export async function findBlockByCreatedOnOrAfter(
  createdOnOrAfter: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<BlockFromDatabase | undefined> {
  const baseQuery: QueryBuilder<BlockModel> = setupBaseQuery<BlockModel>(
    BlockModel,
    options,
  );

  const blocks = await baseQuery
    .where(BlockColumns.time, '>=', createdOnOrAfter)
    .orderBy(BlockColumns.time, Ordering.ASC)
    .limit(1)
    .returning('*');
  return blocks.length > 0 ? blocks[0] : undefined;
}

// Mark as deprecated to encourage migration to more specific methods
/**
 * @deprecated Use findByBlockHeights or findByCreatedOnOrAfter instead
 */
export async function findAll(
  {
    blockHeight,
    limit,
  }: BlockQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<BlockFromDatabase[]> {
  verifyAllRequiredFields(
    {
      blockHeight,
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
): Promise<BlockFromDatabase> {
  const baseQuery: QueryBuilder<BlockModel> = setupBaseQuery<BlockModel>(
    BlockModel,
    options,
  );

  const results: BlockModel[] = await baseQuery
    .orderBy(BlockColumns.blockHeight, Ordering.DESC)
    .limit(1)
    .returning('*');

  const latestBlock: BlockFromDatabase | undefined = results[0];
  if (latestBlock === undefined) {
    logger.error({
      at: 'block-table#getLatest',
      message: 'Unable to find latest block',
    });
    throw new Error('Unable to find latest block');
  }
  return latestBlock;
}
