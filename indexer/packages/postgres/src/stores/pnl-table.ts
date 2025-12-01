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
  Ordering,
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
    baseQuery = baseQuery.whereIn(PnlColumns.subaccountId, subaccountId);
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

<<<<<<< HEAD
  return handleLimitAndPagination(baseQuery, limit, page);
=======
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
>>>>>>> 3f8d74c1 ([ENG-1369] Fix daily PnL aggregation to exclude child subaccounts created mid-day (#3248))
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
    const count: { count?: string } = (await query
      .clone()
      .clearOrder()
      .count({ count: '*' })
      .first()) as unknown as { count?: string };

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
    baseQuery = baseQuery.whereIn(PnlColumns.subaccountId, subaccountId);
  } else if (parentSubaccount !== undefined) {
    baseQuery = baseQuery.whereIn(
      PnlColumns.subaccountId,
      getSubaccountQueryForParent(parentSubaccount),
    );
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

  const knex = PnlModel.knex();
  // 1. Identify the latest record for each subaccount (with RANK = 1 over entire subaccount)
  // 2. For all other records, rank them within their day (RANK ordered by time ascending)
  // 3. Select the latest record and earliest records for each other day
  const rankQuery = baseQuery.clone()
    .select('*')
    .select(
      knex.raw(`
      RANK() OVER (
        PARTITION BY "${PnlColumns.subaccountId}" 
        ORDER BY "${PnlColumns.createdAtHeight}" DESC
      ) as latest_rank,
      DATE_TRUNC('day', "${PnlColumns.createdAt}" AT TIME ZONE 'UTC') as day_date,
      RANK() OVER (
        PARTITION BY "${PnlColumns.subaccountId}", DATE_TRUNC('day', "${PnlColumns.createdAt}" AT TIME ZONE 'UTC')
        ORDER BY "${PnlColumns.createdAt}" ASC
      ) as earliest_in_day_rank
    `),
    );

  // Now select only records that are either:
  // 1. The very latest for their subaccount (latest_rank = 1), OR
  // 2. The earliest record for their day (day_rank = 1) but NOT the latest day
  const finalQuery = PnlModel.query(Transaction.get(options.txId))
    .with('ranked_pnl', rankQuery)
    .from(
      knex.raw(`
      (
        SELECT DISTINCT ON ("subaccountId", day_date) *
        FROM ranked_pnl
        WHERE 
          -- Either it's the latest record overall
          (latest_rank = 1)
          OR 
          -- Or it's the earliest record of a day
          (earliest_in_day_rank = 1)
        ORDER BY "subaccountId", day_date, latest_rank ASC
      ) AS unique_daily_records
    `),
    )
    .orderBy(PnlColumns.subaccountId, Ordering.ASC)
    .orderBy(PnlColumns.createdAtHeight, Ordering.DESC);

  // Apply pagination if needed
  return handleLimitAndPagination(finalQuery, limit, page);
}
