import Knex from 'knex';
import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  setupBaseQuery,
  verifyAllRequiredFields,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getSubaccountQueryForParent } from '../lib/parent-subaccount-helpers';
import PnlModel from '../models/pnl-model';
import {
  Options,
  PnlFromDatabase,
  PaginationFromDatabase,
  QueryableField,
  QueryConfig,
  PnlQueryConfig,
  Ordering,
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
    parentSubaccount,
  }: PnlQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<PnlFromDatabase>> {
  if (parentSubaccount !== undefined && subaccountId !== undefined) {
    throw new Error('Cannot specify both parentSubaccount and subaccountId');
  }

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
      parentSubaccount,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<PnlModel> = setupBaseQuery<PnlModel>(
    PnlModel,
    options,
  );

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(PnlColumns.subaccountId, Array.isArray(subaccountId)
      ? subaccountId : [subaccountId]);
  } else if (parentSubaccount !== undefined) {
    baseQuery = baseQuery.whereIn(
      PnlColumns.subaccountId,
      getSubaccountQueryForParent(parentSubaccount),
    );
  }

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

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(column, order);
    }
  } else {
    baseQuery = baseQuery.orderBy(
      PnlColumns.subaccountId,
      Ordering.ASC,
    ).orderBy(
      PnlColumns.createdAtHeight,
      Ordering.DESC,
    );
  }

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
    const count: { count?: string } = (await baseQuery
      .clone()
      .clearOrder()
      .count({ count: '*' })
      .first()) as unknown as { count?: string };

    baseQuery = baseQuery.offset(offset).limit(limit);

    const results = (await baseQuery.returning('*')) as PnlFromDatabase[];
    return {
      results,
      limit,
      offset,
      total: parseInt(count.count ?? '0', 10),
    };
  }

  // If no pagination, just apply the limit
  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  const results = (await baseQuery.returning('*')) as PnlFromDatabase[];
  return {
    results,
  };
}

export async function findAllHourlyAggregate(
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
  const aggregateBase = setupBaseQuery<PnlModel>(PnlModel, options);
  const knex = (aggregateBase as unknown as { knex?: () => Knex }).knex?.() ?? PnlModel.knex();
  const hourlyAggregateQuery = aggregateBase
    .clearSelect()
    .select(
      knex.raw('DATE_TRUNC(\'hour\', "createdAt") as "createdAt"'),
      knex.raw('MAX("createdAtHeight"::bigint)::text as "createdAtHeight"'),
      knex.raw('SUM(equity::numeric) as equity'),
      knex.raw('SUM("totalPnl"::numeric) as "totalPnl"'),
      knex.raw('SUM("netTransfers"::numeric) as "netTransfers"'),
    )
    .from(baseQuery.as('filtered_pnl'))
    .groupByRaw('DATE_TRUNC(\'hour\', "createdAt")');

  // Apply ordering
  let finalQuery = hourlyAggregateQuery;

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      if (column === 'createdAtHeight') {
        finalQuery = finalQuery.orderByRaw(`MAX("createdAtHeight"::bigint) ${order}`);
      } else if (column === PnlColumns.createdAt) {
        finalQuery = finalQuery.orderByRaw(`DATE_TRUNC('hour', "${column}") ${order}`);
      } else {
        finalQuery = finalQuery.orderBy(column as string, order);
      }
    }
  } else {
    finalQuery = finalQuery.orderByRaw('DATE_TRUNC(\'hour\', "createdAt") DESC');
  }

  return handleLimitAndPagination(finalQuery, limit, page, options);
}

export async function findAllDailyAggregate(
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

  const dailyBase = setupBaseQuery<PnlModel>(PnlModel, options);
  const knex = (dailyBase as unknown as { knex?: () => Knex }).knex?.() ?? PnlModel.knex();
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

  // Step 1: Find the earliest timestamp for each day across ALL subaccounts
  const earliestTimestampPerDay = baseQuery.clone()
    .clearSelect()
    .select(
      knex.raw('DATE_TRUNC(\'day\', "createdAt") as day_date'),
      knex.raw('MIN("createdAt") as earliest_timestamp'),
    )
    .groupByRaw('DATE_TRUNC(\'day\', "createdAt")');

  // Step 2: Get all records at those earliest timestamps
  const dailySnapshotsQuery = baseQuery.clone()
    .select(
      PnlColumns.createdAt,
      PnlColumns.createdAtHeight,
      PnlColumns.equity,
      PnlColumns.totalPnl,
      PnlColumns.netTransfers,
    )
    .innerJoin(
      earliestTimestampPerDay.as('earliest'),
      function joinCondition() {
        this.on(knex.raw('DATE_TRUNC(\'day\', "pnl"."createdAt") = "earliest"."day_date"'))
          .andOn(knex.raw('"pnl"."createdAt" = "earliest"."earliest_timestamp"'));
      },
    );

  // Step 3: Aggregate across subaccounts by day
  const aggregatedQuery = dailyBase
    .clearSelect()
    .with('daily_snapshots', dailySnapshotsQuery)
    .select(
      knex.raw('DATE_TRUNC(\'day\', "createdAt") as "createdAt"'),
      knex.raw('MAX("createdAtHeight"::bigint)::text as "createdAtHeight"'),
      knex.raw('SUM(equity::numeric) as equity'),
      knex.raw('SUM("totalPnl"::numeric) as "totalPnl"'),
      knex.raw('SUM("netTransfers"::numeric) as "netTransfers"'),
    )
    .from('daily_snapshots')
    .groupByRaw('DATE_TRUNC(\'day\', "createdAt")')
    .orderByRaw('DATE_TRUNC(\'day\', "createdAt") DESC');

  return handleLimitAndPagination(aggregatedQuery, limit, page, options);
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
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<PnlFromDatabase>> {
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
    const ModelClass = baseQuery.modelClass();
    const aggregatedQueryKnex = baseQuery.clone().toKnexQuery();

    // Build count query using Objection to preserve transaction
    const countQueryBuilder = ModelClass
      .query(Transaction.get(options.txId))
      .count('* as count')
      .from(aggregatedQueryKnex.as('subquery'))
      .first();

    const countResult = await countQueryBuilder as unknown as { count?: string | number } |
    undefined;

    let total = 0;
    if (countResult?.count) {
      total = typeof countResult.count === 'string'
        ? parseInt(countResult.count, 10)
        : countResult.count as number;
    }

    // Apply pagination
    const paginatedQuery = baseQuery.offset(offset).limit(limit);
    const results = (await paginatedQuery) as PnlFromDatabase[];

    return {
      results,
      limit,
      offset,
      total,
    };
  }

  // If no pagination, just apply the limit
  let query = baseQuery;
  if (limit !== undefined) {
    query = query.limit(limit);
  }

  const results = (await query) as PnlFromDatabase[];

  return {
    results,
  };
}
