import _ from 'lodash';
import { DateTime } from 'luxon';
import { QueryBuilder } from 'objection';

import {
  BUFFER_ENCODING_UTF_8,
  DEFAULT_POSTGRES_OPTIONS,
  ZERO_TIME_ISO_8601,
} from '../constants';
import { knexReadReplica } from '../helpers/knex';
import { setupBaseQuery, verifyAllInjectableVariables, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import { VAULTS_CLOB_0_TO_999_STR_CONCAT } from '../lib/helpers';
import { getSubaccountQueryForParent } from '../lib/parent-subaccount-helpers';
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
  LeaderboardPnlTimeSpan,
  PnlTickInterval,
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
    parentSubaccount,
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
  if (subaccountId !== undefined && parentSubaccount !== undefined) {
    throw new Error('Cannot specify both subaccountId and parentSubaccount in pnl ticks query');
  }

  verifyAllRequiredFields(
    {
      limit,
      id,
      subaccountId,
      parentSubaccount,
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
  } else if (parentSubaccount !== undefined) {
    baseQuery = baseQuery.whereIn(
      PnlTicksColumns.subaccountId,
      getSubaccountQueryForParent(parentSubaccount),
    );
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
    rows: [{ max: string, count: number }],
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

export async function findMostRecentPnlTickForEachAccount(
  createdOnOrAfterHeight: string,
): Promise<{
  [subaccountId: string]: PnlTicksCreateObject,
}> {
  verifyAllInjectableVariables([createdOnOrAfterHeight]);

  const result: {
    rows: PnlTicksFromDatabase[],
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
    return getAllTimeRankedPnlTicks();
  }
  return getRankedPnlTicksForTimeSpan(timeSpan);
}

function convertTimespanToSQL(timeSpan: string): string {
  const timeSpanEnum: LeaderboardPnlTimeSpan = LeaderboardPnlTimeSpan[
    timeSpan as keyof typeof LeaderboardPnlTimeSpan];
  switch (timeSpanEnum) {
    case LeaderboardPnlTimeSpan.ONE_DAY:
      return '1 days';
    case LeaderboardPnlTimeSpan.SEVEN_DAYS:
      return '7 days';
    case LeaderboardPnlTimeSpan.THIRTY_DAYS:
      return '30 days';
    case LeaderboardPnlTimeSpan.ONE_YEAR:
      return '365 days';
    default:
      throw new Error(`Invalid time span: ${timeSpan}`);
  }
}

/**
 * Constructs a complex SQL query to calculate the Pnl difference and current equity
 * of subaccounts over a specified time span, ranking them by their PnL.
 *
 *  This has 5 main parts
 * 1. latest_subaccount_pnl_x_days_ago: Identifies the most recent PnL tick for each subaccount
 *    before the specified time span. It filters out subaccounts which are not parent
 *    subaccounts or associated child subaccounts. It also excludes any addresses
 *    that are vault addresses.
 *
 * 2. latest_pnl: Finds the latest PnL tick for each subaccount as of the current date,
 *    applying the same filters as latest_subaccount_pnl_x_days_ago.
 *
 * 3. subaccount_pnl_difference: Calculates the difference in PnL between the
 *    current date and the start of the specified time span for each subaccount.
 *
 * 4. aggregated_results: Aggregates the PnL differences and current equity for
 *    all subaccounts, grouping by address.
 *
 * 5. The final SELECT statement then ranks the addresses based on their total PnL
 *    in descending order, providing a snapshot of subaccount performance over the
 *    specified time span.
 *
*/
async function getRankedPnlTicksForTimeSpan(
  timeSpan: string,
): Promise<LeaderboardPnlCreateObject[]> {
  const vaultAddressesString: string = VAULTS_CLOB_0_TO_999_STR_CONCAT;
  const intervalSqlString: string = convertTimespanToSQL(timeSpan);
  const result: {
    rows: LeaderboardPnlCreateObject[],
  } = await knexReadReplica.getConnection().raw(
    `
    WITH latest_subaccount_pnl_x_days_ago AS (
      SELECT DISTINCT ON (a."subaccountId")
          a."subaccountId",
          a."totalPnl",
          b."address"
      FROM
          pnl_ticks a
      LEFT JOIN
          subaccounts b ON a."subaccountId" = b."id"
      WHERE
          a."createdAt"::date <= (CURRENT_DATE - INTERVAL '${intervalSqlString}')
          AND (b."subaccountNumber" % 128) = 0
          AND b."address" NOT IN (${vaultAddressesString})
      ORDER BY a."subaccountId", a."blockHeight" DESC
    ),
    latest_pnl as (
      SELECT DISTINCT ON (a."subaccountId")
          "subaccountId",
          "totalPnl",
          "equity" as "currentEquity",
          "address"
      FROM
          pnl_ticks a left join subaccounts b ON a."subaccountId"=b."id"
      WHERE
          "createdAt"::date = CURRENT_DATE
          AND (b."subaccountNumber" % 128) = 0
          AND b."address" NOT IN (${vaultAddressesString})
      ORDER BY a."subaccountId", "blockHeight" DESC
    ), 
    subaccount_pnl_difference as(
      SELECT
        a."address",
        a."totalPnl" - COALESCE(b."totalPnl", 0) as "pnlDifference",
        a."currentEquity" as "currentEquity"
      FROM latest_pnl a left join latest_subaccount_pnl_x_days_ago b
      ON a."subaccountId"=b."subaccountId"
    ), aggregated_results as(
    SELECT
      "address",
      sum(subaccount_pnl_difference."pnlDifference") as "totalPnl",
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

/**
 * Constructs a query to calculate and rank the Profit and Loss (PnL) and current equity of
 * subaccounts for the current day. This query is divided into 3 main parts:
 * 1. latest_pnl: This selects the most recent PnL tick for each Parent subaccount
 *    and associated child subaccounts. It filters subaccounts based on the current date.
 *    Additionally, it excludes any addresses that are vault addresses.
 *
 * 2. aggregated_results: This CTE aggregates the results from latest_pnl by address.
 *    It sums up the total PnL and current equity for each address.
 *
 * 3. The final SELECT statement calculates a rank for each address based on the total PnL in
 *    descending order along with associated fields
 */
async function getAllTimeRankedPnlTicks(): Promise<LeaderboardPnlCreateObject[]> {
  const vaultAddressesString: string = VAULTS_CLOB_0_TO_999_STR_CONCAT;
  const result: {
    rows: LeaderboardPnlCreateObject[],
  } = await knexReadReplica.getConnection().raw(
    `
    WITH latest_pnl as (
      SELECT DISTINCT ON (a."subaccountId")
          "subaccountId",
          "totalPnl",
          "equity" as "currentEquity",
          "address"
      FROM
          pnl_ticks a left join subaccounts b ON a."subaccountId"=b."id"
      WHERE
          "createdAt"::date = CURRENT_DATE
          AND (b."subaccountNumber" % 128) = 0
          AND b."address" NOT IN (${vaultAddressesString})
      ORDER BY a."subaccountId", "blockHeight" DESC
    ), aggregated_results as(
    SELECT
      "address",
      sum(latest_pnl."totalPnl") as "totalPnl",
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

/**
 * Constructs a query to get pnl ticks at a specific interval for a set of subaccounts
 * within a time range.
 * Uses a windowing function in the raw query to get the first row of each window of the specific
 * interval time.
 * Currently only supports hourly / daily as the interval.
 * @param interval 'day' or 'hour'.
 * @param timeWindowSeconds Window of time to get pnl ticks for at the specified interval.
 * @param subaccountIds Set of subaccounts to get pnl ticks for.
 * @returns
 */
export async function getPnlTicksAtIntervals(
  interval: PnlTickInterval,
  timeWindowSeconds: number,
  subaccountIds: string[],
  earliestDate: DateTime,
): Promise <PnlTicksFromDatabase[]> {
  if (subaccountIds.length === 0) {
    return [];
  }
  const result: {
    rows: PnlTicksFromDatabase[],
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT
      "id",
      "subaccountId",
      "equity",
      "totalPnl",
      "netTransfers",
      "createdAt",
      "blockHeight",
      "blockTime"
    FROM (
      SELECT
        pnl_ticks.*,
        ROW_NUMBER() OVER (
          PARTITION BY
            "subaccountId", 
            DATE_TRUNC(
              '${interval}',
              "blockTime"
            ) ORDER BY "blockTime"
        ) AS r
      FROM pnl_ticks
      WHERE
        "subaccountId" IN (${subaccountIds.map((id: string) => { return `'${id}'`; }).join(',')}) AND
        "blockTime" >= '${earliestDate.toUTC().toISO()}'::timestamp AND
        "blockTime" > NOW() - INTERVAL '${timeWindowSeconds} second'
    ) AS pnl_intervals
    WHERE
      r = 1
    ORDER BY "subaccountId";
    `,
  ) as unknown as {
    rows: PnlTicksFromDatabase[],
  };

  return result.rows;
}

export async function getLatestPnlTick(
  subaccountIds: string[],
  beforeOrAt: DateTime,
): Promise <PnlTicksFromDatabase[]> {
  if (subaccountIds.length === 0) {
    return [];
  }
  const result: {
    rows: PnlTicksFromDatabase[],
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT
      DISTINCT ON ("subaccountId")
      "id",
      "subaccountId",
      "equity",
      "totalPnl",
      "netTransfers",
      "createdAt",
      "blockHeight",
      "blockTime"
    FROM
      pnl_ticks
    WHERE
      "subaccountId" in (${subaccountIds.map((id: string) => { return `'${id}'`; }).join(',')}) AND
      "blockTime" <= '${beforeOrAt.toUTC().toISO()}'::timestamp AND
      "blockTime" >= '${beforeOrAt.toUTC().minus({ hours: 4 }).toISO()}'::timestamp
    ORDER BY
      "subaccountId",
      "blockTime" DESC
    `,
  ) as unknown as {
    rows: PnlTicksFromDatabase[],
  };

  return result.rows;
}
