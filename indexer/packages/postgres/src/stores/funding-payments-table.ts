import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import FundingPaymentsModel from '../models/funding-payments-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  FundingPaymentsColumns,
  FundingPaymentsCreateObject,
  FundingPaymentsFromDatabase,
  FundingPaymentsQueryConfig,
} from '../types';

export async function findAll(
  {
    limit,
    subaccountId,
    perpetualId,
    ticker,
    createdAtHeight,
    createdAt,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    createdOnOrAfterHeight,
    createdOnOrAfter,
  }: FundingPaymentsQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<FundingPaymentsFromDatabase[]> {
  verifyAllRequiredFields(
    {
      limit,
      subaccountId,
      perpetualId,
      ticker,
      createdAtHeight,
      createdAt,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdOnOrAfterHeight,
      createdOnOrAfter,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<FundingPaymentsModel> = setupBaseQuery<FundingPaymentsModel>(
    FundingPaymentsModel,
    options,
  );

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(FundingPaymentsColumns.subaccountId, subaccountId);
  }

  if (perpetualId !== undefined) {
    baseQuery = baseQuery.whereIn(FundingPaymentsColumns.perpetualId, perpetualId);
  }

  if (ticker !== undefined) {
    baseQuery = baseQuery.where(FundingPaymentsColumns.ticker, ticker);
  }

  if (createdAtHeight !== undefined) {
    baseQuery = baseQuery.where(FundingPaymentsColumns.createdAtHeight, createdAtHeight);
  }

  if (createdAt !== undefined) {
    baseQuery = baseQuery.where(FundingPaymentsColumns.createdAt, createdAt);
  }

  if (createdBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      FundingPaymentsColumns.createdAtHeight,
      '<=',
      createdBeforeOrAtHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(FundingPaymentsColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdOnOrAfterHeight !== undefined) {
    baseQuery = baseQuery.where(
      FundingPaymentsColumns.createdAtHeight,
      '>=',
      createdOnOrAfterHeight,
    );
  }

  if (createdOnOrAfter !== undefined) {
    baseQuery = baseQuery.where(FundingPaymentsColumns.createdAt, '>=', createdOnOrAfter);
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
      FundingPaymentsColumns.createdAtHeight,
      Ordering.DESC,
    );
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  fundingPaymentToCreate: FundingPaymentsCreateObject,
  options: Options = { txId: undefined },
): Promise<FundingPaymentsFromDatabase> {
  return FundingPaymentsModel.query(
    Transaction.get(options.txId),
  ).insert(fundingPaymentToCreate).returning('*');
}

export async function findById(
  subaccountId: string,
  createdAt: string,
  ticker: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<FundingPaymentsFromDatabase | undefined> {
  const baseQuery: QueryBuilder<FundingPaymentsModel> = setupBaseQuery<FundingPaymentsModel>(
    FundingPaymentsModel,
    options,
  );
  return baseQuery
    .where(FundingPaymentsColumns.subaccountId, subaccountId)
    .where(FundingPaymentsColumns.createdAt, createdAt)
    .where(FundingPaymentsColumns.ticker, ticker)
    .first()
    .returning('*');
}
