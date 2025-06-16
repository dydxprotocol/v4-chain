import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getSubaccountQueryForParent } from '../lib/parent-subaccount-helpers';
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
  PaginationFromDatabase,
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
    page,
    parentSubaccount,
  }: FundingPaymentsQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<FundingPaymentsFromDatabase>> {
  if (parentSubaccount !== undefined && subaccountId !== undefined) {
    throw new Error('Cannot specify both parentSubaccount and subaccountId');
  }

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
  } else if (parentSubaccount !== undefined) {
    baseQuery = baseQuery.whereIn(
      FundingPaymentsColumns.subaccountId,
      getSubaccountQueryForParent(parentSubaccount),
    );
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
      baseQuery = baseQuery.orderBy(column, order);
    }
  } else {
    baseQuery = baseQuery.orderBy(FundingPaymentsColumns.createdAtHeight, Ordering.DESC);
  }

  return handleLimitAndPagination(baseQuery, limit, page);
}

export async function create(
  fundingPaymentToCreate: FundingPaymentsCreateObject,
  options: Options = { txId: undefined },
): Promise<FundingPaymentsFromDatabase> {
  return FundingPaymentsModel.query(Transaction.get(options.txId))
    .insert(fundingPaymentToCreate)
    .returning('*');
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

/**
 * Handles pagination and limit logic for funding payment queries
 * @param baseQuery The base query to apply pagination to
 * @param limit Maximum number of funding payments to return
 * @param page Page number
 * @returns Promise<PaginationFromDatabase<FundingPaymentsFromDatabase>>
 */
async function handleLimitAndPagination(
  baseQuery: QueryBuilder<FundingPaymentsModel>,
  limit?: number,
  page?: number,
): Promise<PaginationFromDatabase<FundingPaymentsFromDatabase>> {
  let query = baseQuery;

  /**
   * If a query is made using a page number, then the limit property is used as 'page limit'
   */
  if (page !== undefined && limit !== undefined) {
    /**
     * We make sure that the page number is always >= 1
     */
    const currentPage: number = Math.max(1, page);
    const offset: number = (currentPage - 1) * limit;

    /**
     * Ensure sorting is applied to maintain consistent pagination results.
     * Also a casting of the ts type is required since the infer of the type
     * obtained from the count is not performed.
     */
    const count: { count?: string } = (await query
      .clone()
      .clearOrder()
      .count({ count: '*' })
      .first()) as unknown as { count?: string };

    query = query.offset(offset).limit(limit);

    const results = (await query.returning('*')) as FundingPaymentsFromDatabase[];
    return {
      results,
      limit,
      offset,
      total: parseInt(count.count ?? '0', 10),
    };
  }

  // If no pagination, just apply the limit
  if (limit !== undefined) {
    query = query.limit(limit);
  }

  const results = (await query.returning('*')) as FundingPaymentsFromDatabase[];
  return {
    results,
  };
}
