import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  verifyAllRequiredFields,
  setupBaseQuery,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import TransactionModel from '../models/transaction-model';
import {
  QueryConfig,
  TransactionFromDatabase,
  TransactionQueryConfig,
  TransactionColumns,
  TransactionCreateObject,
  Options,
  Ordering,
  QueryableField,
} from '../types';

export function uuid(blockHeight: string, transactionIndex: number): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(Buffer.from(`${blockHeight}-${transactionIndex}`, BUFFER_ENCODING_UTF_8));
}

export async function findAll(
  {
    id,
    blockHeight,
    transactionIndex,
    transactionHash,
    limit,
  }: TransactionQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TransactionFromDatabase[]> {
  verifyAllRequiredFields(
    {
      id,
      blockHeight,
      transactionIndex,
      transactionHash,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<TransactionModel> = setupBaseQuery<TransactionModel>(
    TransactionModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(TransactionColumns.id, id);
  }

  if (blockHeight !== undefined) {
    baseQuery = baseQuery.whereIn(TransactionColumns.blockHeight, blockHeight);
  }

  if (transactionIndex !== undefined) {
    baseQuery = baseQuery.whereIn(TransactionColumns.transactionIndex, transactionIndex);
  }

  if (transactionHash !== undefined) {
    baseQuery = baseQuery.whereIn(TransactionColumns.transactionHash, transactionHash);
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
      TransactionColumns.blockHeight,
      Ordering.ASC,
    ).orderBy(
      TransactionColumns.transactionIndex,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  transactionToCreate: TransactionCreateObject,
  options: Options = { txId: undefined },
): Promise<TransactionFromDatabase> {
  return TransactionModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...transactionToCreate,
    id: uuid(transactionToCreate.blockHeight, transactionToCreate.transactionIndex),
  }).returning('*');
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TransactionFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TransactionModel> = setupBaseQuery<TransactionModel>(
    TransactionModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}
