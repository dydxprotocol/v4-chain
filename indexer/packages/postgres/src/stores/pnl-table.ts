import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  setupBaseQuery,
  verifyAllRequiredFields,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import PnlModel from '../models/pnl-model';
import {
  Options,
  PnlFromDatabase,
  PaginationFromDatabase,
  QueryableField,
  QueryConfig,
  PnlQueryConfig,
} from '../types';
import { PnlColumns, PnlCreateObject } from '../types/pnl-types';

export async function findAll(
  {
    limit,
    subaccountId,
    createdAtHeight,
    createdAt,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    createdOnOrAfterHeight,
    createdOnOrAfter,
    page,
  }: PnlQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<PnlFromDatabase>> {
  verifyAllRequiredFields(
    {
      limit,
      subaccountId,
      createdAtHeight,
      createdAt,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdOnOrAfterHeight,
      createdOnOrAfter,
      page,
    } as QueryConfig,
    requiredFields,
  );

  // Check for valid subaccountId
  if (!subaccountId || !Array.isArray(subaccountId) || subaccountId.length === 0) {
    throw new Error('subaccountId array must be provided and non-empty');
  }

  const knex = PnlModel.knex();
  let baseQuery: QueryBuilder<PnlModel> = setupBaseQuery<PnlModel>(
    PnlModel,
    options,
  );

  baseQuery = baseQuery.whereIn(PnlColumns.subaccountId, subaccountId);

  if (createdAt !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, createdAt);
  }

  if (createdAtHeight !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAtHeight, createdAtHeight);
  }

  if (createdBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlColumns.createdAtHeight,
      '<=',
      createdBeforeOrAtHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdOnOrAfterHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlColumns.createdAtHeight,
      '>=',
      createdOnOrAfterHeight,
    );
  }

  if (createdOnOrAfter !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, '>=', createdOnOrAfter);
  }

  // Aggregate by hour across all subaccounts
  const hourlyAggregateQuery = PnlModel.query(Transaction.get(options.txId))
    .select(
      knex.raw('DATE_TRUNC(\'hour\', "createdAt") as "createdAt"'),
      knex.raw('MAX("createdAtHeight") as "createdAtHeight"'),
      knex.raw('SUM(equity::numeric) as equity'),
      knex.raw('SUM("totalPnl"::numeric) as "totalPnl"'),
      knex.raw('SUM("netTransfers"::numeric) as "netTransfers"'),
    )
    .from(baseQuery.as('filtered_pnl'))
    .groupByRaw('DATE_TRUNC(\'hour\', "createdAt")');

  // Apply ordering with same expression
  let finalQuery = hourlyAggregateQuery;
  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      finalQuery = finalQuery.orderBy(column, order);
    }
  } else {
    finalQuery = finalQuery.orderByRaw('DATE_TRUNC(\'hour\', "createdAt") DESC');
  }

  return handleLimitAndPagination(finalQuery, limit, page);
}

export async function create(
  pnlToCreate: PnlCreateObject,
  options: Options = { txId: undefined },
): Promise<PnlFromDatabase> {
  return PnlModel.query(
    Transaction.get(options.txId),
  ).insert(pnlToCreate).returning('*');
}

export async function findById(
  subaccountId: string,
  createdAt: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PnlFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PnlModel> = setupBaseQuery<PnlModel>(
    PnlModel,
    options,
  );
  return baseQuery
    .where(PnlColumns.subaccountId, subaccountId)
    .where(PnlColumns.createdAt, createdAt)
    .first();
}

async function handleLimitAndPagination(
  baseQuery: QueryBuilder<PnlModel>,
  limit?: number,
  page?: number,
): Promise<PaginationFromDatabase<PnlFromDatabase>> {
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
    const countQuery = query
      .modelClass()
      .query()
      .count('* as count')
      .from(query.clone().clearOrder().as('grouped_results'))
      .first();

    const count: { count?: string } = (await countQuery) as unknown as { count?: string };

    query = query.offset(offset).limit(limit);

    const results = (await query) as PnlFromDatabase[];
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

  const results = (await query) as PnlFromDatabase[];
  return {
    results,
  };
}

export async function findAllDailyPnl(
  {
    limit,
    subaccountId,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    createdOnOrAfterHeight,
    createdOnOrAfter,
    page,
  }: PnlQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<PnlFromDatabase>> {
  verifyAllRequiredFields(
    {
      limit,
      subaccountId,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdOnOrAfterHeight,
      createdOnOrAfter,
      page,
    } as QueryConfig,
    requiredFields,
  );

  if (!subaccountId || !Array.isArray(subaccountId) || subaccountId.length === 0) {
    throw new Error('subaccountId array must be provided and non-empty');
  }

  let baseQuery: QueryBuilder<PnlModel> = setupBaseQuery<PnlModel>(
    PnlModel,
    options,
  );

  const knex = PnlModel.knex();
  baseQuery = baseQuery.whereIn(PnlColumns.subaccountId, subaccountId);

  if (createdBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlColumns.createdAtHeight,
      '<=',
      createdBeforeOrAtHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdOnOrAfterHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlColumns.createdAtHeight,
      '>=',
      createdOnOrAfterHeight,
    );
  }

  if (createdOnOrAfter !== undefined) {
    baseQuery = baseQuery.where(PnlColumns.createdAt, '>=', createdOnOrAfter);
  }

  // Step 1: Get first record of each day for each subaccount
  const dailySnapshotsQuery = baseQuery.clone()
    .select(
      knex.raw(`
        DISTINCT ON ("subaccountId", DATE_TRUNC('day', "createdAt"))
        "subaccountId",
        "createdAt",
        "createdAtHeight",
        equity,
        "totalPnl",
        "netTransfers"
      `),
    )
    .orderByRaw('"subaccountId", DATE_TRUNC(\'day\', "createdAt")')
    .orderBy(PnlColumns.createdAt, 'ASC'); // Earliest in day

  // Step 2: Aggregate across subaccounts by day
  const aggregatedQuery = PnlModel.query(Transaction.get(options.txId))
    .with('daily_snapshots', dailySnapshotsQuery)
    .select(
      knex.raw('DATE_TRUNC(\'day\', "createdAt") as "createdAt"'),
      knex.raw('MAX("createdAtHeight") as "createdAtHeight"'),
      knex.raw('SUM(equity::numeric) as equity'),
      knex.raw('SUM("totalPnl"::numeric) as "totalPnl"'),
      knex.raw('SUM("netTransfers"::numeric) as "netTransfers"'),
    )
    .from('daily_snapshots')
    .groupByRaw('DATE_TRUNC(\'day\', "createdAt")')
    .orderByRaw('DATE_TRUNC(\'day\', "createdAt") DESC');

  // Apply pagination if needed
  return handleLimitAndPagination(aggregatedQuery, limit, page);
}
