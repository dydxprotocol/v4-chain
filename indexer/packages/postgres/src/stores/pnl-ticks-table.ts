import _ from 'lodash';
import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS, ZERO_TIME_ISO_8601, LEADERBOARD_TIMESPAN } from '../constants';
import { knexReadReplica } from '../helpers/knex';
import { setupBaseQuery, verifyAllInjectableVariables, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import PnlTicksModel from '../models/pnl-ticks-model';
import {
  Options,
  Ordering,
  PnlTicksColumns,
  PnlTicksCreateObject,
  PnlTicksFromDatabase,
  PnlTicksQueryConfig,
  QueryableField,
  QueryConfig,
  PaginationFromDatabase,
  LeaderboardPnlCreateObject,
} from '../types';

export function uuid(
  subaccountId: string,
  createdAt: string,
): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(
    Buffer.from(
      `${subaccountId}-${createdAt}`,
      BUFFER_ENCODING_UTF_8),
  );
}

export async function findAll(
  {
    limit,
    id,
    subaccountId,
    createdAt,
    blockHeight,
    blockTime,
    createdBeforeOrAt,
    createdBeforeOrAtBlockHeight,
    createdOnOrAfter,
    createdOnOrAfterBlockHeight,
    page,
  }: PnlTicksQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<PnlTicksFromDatabase>> {
  verifyAllRequiredFields(
    {
      limit,
      id,
      subaccountId,
      createdAt,
      blockHeight,
      blockTime,
      createdBeforeOrAt,
      createdBeforeOrAtBlockHeight,
      createdOnOrAfter,
      createdOnOrAfterBlockHeight,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<PnlTicksModel> = setupBaseQuery<PnlTicksModel>(
    PnlTicksModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(PnlTicksColumns.id, id);
  }

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(PnlTicksColumns.subaccountId, subaccountId);
  }

  if (createdAt !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.createdAt, createdAt);
  }

  if (blockHeight !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.blockHeight, blockHeight);
  }

  if (blockTime !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.blockTime, blockTime);
  }

  if (createdBeforeOrAtBlockHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlTicksColumns.blockHeight,
      '<=',
      createdBeforeOrAtBlockHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdOnOrAfterBlockHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlTicksColumns.blockHeight,
      '>=',
      createdOnOrAfterBlockHeight,
    );
  }

  if (createdOnOrAfter !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.createdAt, '>=', createdOnOrAfter);
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
      PnlTicksColumns.subaccountId,
      Ordering.ASC,
    ).orderBy(
      PnlTicksColumns.blockHeight,
      Ordering.DESC,
    );
  }

  if (limit !== undefined && page === undefined) {
    baseQuery = baseQuery.limit(limit);
  }

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
    const count: { count?: string } = await baseQuery.clone().clearOrder().count({ count: '*' }).first() as unknown as { count?: string };

    baseQuery = baseQuery.offset(offset).limit(limit);

    return {
      results: await baseQuery.returning('*'),
      limit,
      offset,
      total: parseInt(count.count ?? '0', 10),
    };
  }

  return {
    results: await baseQuery.returning('*'),
  };
}

export async function create(
  pnlTicksToCreate: PnlTicksCreateObject,
  options: Options = { txId: undefined },
): Promise<PnlTicksFromDatabase> {
  return PnlTicksModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...pnlTicksToCreate,
    id: uuid(pnlTicksToCreate.subaccountId, pnlTicksToCreate.createdAt),
  }).returning('*');
}

export async function createMany(
  pnlTicks: PnlTicksCreateObject[],
  options: Options = { txId: undefined },
): Promise<PnlTicksFromDatabase[]> {
  const ticks: PnlTicksFromDatabase[] = pnlTicks.map((tick) => ({
    ...tick,
    id: uuid(tick.subaccountId, tick.createdAt),
  }));

  return PnlTicksModel
    .query(Transaction.get(options.txId))
    .insert(ticks)
    .returning('*');
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PnlTicksFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PnlTicksModel> = setupBaseQuery<PnlTicksModel>(
    PnlTicksModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

function convertPnlTicksFromDatabaseToPnlTicksCreateObject(
  pnlTicksFromDatabase: PnlTicksFromDatabase,
): PnlTicksCreateObject {
  return _.omit(pnlTicksFromDatabase, PnlTicksColumns.id);
}

export async function findLatestProcessedBlocktimeAndCount(): Promise<{
  maxBlockTime: string,
  count: number,
}> {
  const result: {
    rows: [{ max: string, count: number }]
  } = await knexReadReplica.getConnection().raw(
    `
    WITH maxBlockTime AS (
      SELECT MAX("blockTime") as "maxBlockTime"
      FROM "pnl_ticks"
    )
    SELECT
      maxBlockTime."maxBlockTime" as max,
      COUNT(*) as count
    FROM
      "pnl_ticks",
      maxBlockTime
    WHERE
      "pnl_ticks"."blockTime" = maxBlockTime."maxBlockTime"
    GROUP BY 1
    `,
  ) as unknown as { rows: [{ max: string, count: number }] };

  const maxBlockTime = result.rows[0]?.max || ZERO_TIME_ISO_8601;
  const count = Number(result.rows[0]?.count) || 0;

  return {
    maxBlockTime,
    count,
  };
}

function convertTimespanToSQL(timeSpan: string): string {
  const timeSpanEnum: LEADERBOARD_TIMESPAN = LEADERBOARD_TIMESPAN[timeSpan as keyof typeof LEADERBOARD_TIMESPAN];
  switch (timeSpan) {
    case 'ONE_DAY':
      return '1 days';
    case 'SEVEN_DAYS':
      return '7 days';
    case 'THIRTY_DAYS':
      return '30 days';
    case 'ONE_YEAR':
      return '365 days';
    default:
      throw new Error(`Invalid time span: ${timeSpan}`);
  }
}

export async function findMostRecentPnlTickForEachAccount(
  createdOnOrAfterHeight: string,
): Promise<{
  [subaccountId: string]: PnlTicksCreateObject
}> {
  verifyAllInjectableVariables([createdOnOrAfterHeight]);

  const result: {
    rows: PnlTicksFromDatabase[]
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT DISTINCT ON ("subaccountId") *
    FROM "pnl_ticks"
    WHERE "blockHeight" >= '${createdOnOrAfterHeight}'
    ORDER BY "subaccountId" ASC, "blockHeight" DESC, "createdAt" DESC;
    `
    ,
  ) as unknown as { rows: PnlTicksFromDatabase[] };
  return _.keyBy(
    result.rows.map(convertPnlTicksFromDatabaseToPnlTicksCreateObject),
    'subaccountId',
  );
}

export async function getRankedPnlTicks(
  timeSpan: string,
): Promise<LeaderboardPnlCreateObject[]> {
  if (timeSpan === 'ALL_TIME') {
    return getLatestRankedPnlTicks();
  }
  return getRankedPnlTicksForTimeSpan(timeSpan);
}

async function getRankedPnlTicksForTimeSpan(
  timeSpan: string,
): Promise<LeaderboardPnlCreateObject[]> {

  const interval_sql_string: string = convertTimespanToSQL(timeSpan);
  const result: {
    rows: LeaderboardPnlCreateObject[]
  } = await knexReadReplica.getConnection().raw(
    `
    WITH latest_subaccount_pnl_x_days_ago_ranked AS (
      SELECT
          a."subaccountId",
          a."totalPnl",
          b."address",
          ROW_NUMBER() OVER (PARTITION BY a."subaccountId" ORDER BY a."blockHeight" DESC) AS "rn"
      FROM
          pnl_ticks a
      LEFT JOIN
          subaccounts b ON a."subaccountId" = b."id"
      WHERE
          a."createdAt"::date <= (CURRENT_DATE - INTERVAL '${interval_sql_string}')
          AND (b."subaccountNumber" % 128) = 0
    ),
    latest_subaccount_pnl_x_days_ago AS (
      SELECT
          "subaccountId",
          "totalPnl",
          "address"
      FROM 
          latest_subaccount_pnl_x_days_ago_ranked   
      WHERE
          "rn" = 1
    ), latest_pnl_ranked as (
      SELECT
          "subaccountId",
          "totalPnl",
          "netTransfers",
          "equity" as "currentEquity",
          "address",
          ROW_NUMBER() OVER (PARTITION BY "subaccountId" ORDER BY "blockHeight" DESC) AS "rn"
      FROM
          pnl_ticks a left join subaccounts b ON a."subaccountId"=b."id"
      WHERE
          "createdAt"::date = CURRENT_DATE
          AND (b."subaccountNumber" % 128) = 0
    ), latest_pnl as(
      SELECT
          "subaccountId",
          "totalPnl",
          "netTransfers",
          "currentEquity" ,
          "address"
      FROM 
          latest_pnl_ranked
      WHERE
          "rn" = 1
    ), subaccount_pnl_difference as(
      SELECT
        a."address",
        a."totalPnl" - COALESCE(b."totalPnl", 0) as "pnlDifference",
        a."netTransfers" as "netTransfers",
        a."currentEquity" as "currentEquity"
      FROM latest_pnl a left join latest_subaccount_pnl_x_days_ago b
      ON a."subaccountId"=b."subaccountId"
    ), aggregated_results as(
    SELECT
      "address",
      sum(subaccount_pnl_difference."pnlDifference") as "totalPnl",
      sum(subaccount_pnl_difference."netTransfers") as "netTransfers",
      sum(subaccount_pnl_difference."currentEquity") as "currentEquity"
    FROM
      subaccount_pnl_difference
    GROUP BY address
    )
    SELECT
      "address",
      "totalPnl" as "pnl",
      '${timeSpan}' as "timeSpan",
      "currentEquity",
      ROW_NUMBER() over (order by aggregated_results."totalPnl" desc) as rank
    FROM
      aggregated_results;
    `,
  ) as { rows: LeaderboardPnlCreateObject[] };

  return result.rows;
}


async function getLatestRankedPnlTicks(): Promise<LeaderboardPnlCreateObject[]> {
  const result: {
    rows: LeaderboardPnlCreateObject[]
  } = await knexReadReplica.getConnection().raw(
    `
    WITH latest_pnl_ranked as (
      SELECT
          "subaccountId",
          "totalPnl",
          "netTransfers",
          "equity" as "currentEquity",
          "address",
          ROW_NUMBER() OVER (PARTITION BY "subaccountId" ORDER BY "blockHeight" DESC) AS "rn"
      FROM
          pnl_ticks a left join subaccounts b ON a."subaccountId"=b."id"
      WHERE
          "createdAt"::date = CURRENT_DATE
          AND (b."subaccountNumber" % 128) = 0
    ), latest_pnl as(
      SELECT
          "subaccountId",
          "totalPnl",
          "netTransfers",
          "currentEquity" ,
          "address"
      FROM 
          latest_pnl_ranked
      WHERE
          "rn" = 1
    ), aggregated_results as(
    SELECT
      "address",
      sum(latest_pnl."totalPnl") as "totalPnl",
      sum(latest_pnl."netTransfers") as "netTransfers",
      sum(latest_pnl."currentEquity") as "currentEquity"
    FROM
      latest_pnl
    GROUP BY address
    )
    SELECT
      "address",
      "totalPnl" as "pnl",
      'ALL_TIME' as "timeSpan",
      "currentEquity",
      ROW_NUMBER() over (order by aggregated_results."totalPnl" desc) as rank
    FROM
      aggregated_results;
    `,
  ) as { rows: LeaderboardPnlCreateObject[] };

  return result.rows;
}