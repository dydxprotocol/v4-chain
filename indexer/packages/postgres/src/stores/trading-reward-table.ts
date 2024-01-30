import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import TradingRewardModel from '../models/trading-reward-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  TradingRewardColumns,
  TradingRewardCreateObject,
  TradingRewardFromDatabase,
  TradingRewardQueryConfig,
} from '../types';

export function uuid(address: string, blockHeight: string): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(Buffer.from(`${address}-${blockHeight}`, BUFFER_ENCODING_UTF_8));
}

export async function findAll(
  {
    address,
    blockHeight,
    blockTimeBeforeOrAt,
    blockTimeAfterOrAt,
    blockTimeBefore,
    blockHeightBeforeOrAt,
    limit,
  }: TradingRewardQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TradingRewardFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      blockHeight,
      blockTimeBeforeOrAt,
      blockTimeAfterOrAt,
      blockTimeBefore,
      blockHeightBeforeOrAt,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<TradingRewardModel> = setupBaseQuery<TradingRewardModel>(
    TradingRewardModel,
    options,
  );

  if (address) {
    baseQuery = baseQuery.where(TradingRewardColumns.address, address);
  }

  if (blockHeight) {
    baseQuery = baseQuery.where(TradingRewardColumns.blockHeight, blockHeight);
  }

  if (blockTimeBeforeOrAt) {
    baseQuery = baseQuery.where(TradingRewardColumns.blockTime, '<=', blockTimeBeforeOrAt);
  }

  if (blockTimeAfterOrAt) {
    baseQuery = baseQuery.where(TradingRewardColumns.blockTime, '>=', blockTimeAfterOrAt);
  }

  if (blockTimeBefore) {
    baseQuery = baseQuery.where(TradingRewardColumns.blockTime, '<', blockTimeBefore);
  }

  if (blockHeightBeforeOrAt) {
    baseQuery = baseQuery.where(TradingRewardColumns.blockHeight, '<=', blockHeightBeforeOrAt);
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
      TradingRewardColumns.blockHeight,
      Ordering.DESC,
    ).orderBy(
      TradingRewardColumns.address,
      Ordering.DESC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  tradingRewardToCreate: TradingRewardCreateObject,
  options: Options = { txId: undefined },
): Promise<TradingRewardFromDatabase> {
  return TradingRewardModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...tradingRewardToCreate,
    id: uuid(tradingRewardToCreate.address, tradingRewardToCreate.blockHeight),
  }).returning('*');
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TradingRewardFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TradingRewardModel> = setupBaseQuery<TradingRewardModel>(
    TradingRewardModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}
