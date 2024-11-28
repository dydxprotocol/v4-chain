import Knex from 'knex';
import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexPrimary } from '../helpers/knex';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import AffiliatePerRefereeStatsModel from '../models/affiliate-per-referee-stats-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  AffiliatePerRefereeStatsColumns,
  AffiliatePerRefereeStatsCreateObject,
  AffiliatePerRefereeStatsFromDatabase,
  Liquidity,
  FillType,
  AffiliatePerRefereeStatsQueryConfig,
} from '../types';

export async function findAll(
  {
    affiliateAddress,
    limit,
  }: AffiliatePerRefereeStatsQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliatePerRefereeStatsFromDatabase[]> {
  verifyAllRequiredFields(
    {
      affiliateAddress,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery
  : QueryBuilder<AffiliatePerRefereeStatsModel> = setupBaseQuery<AffiliatePerRefereeStatsModel>(
    AffiliatePerRefereeStatsModel,
    options,
  );

  if (affiliateAddress) {
    baseQuery = baseQuery.where(AffiliatePerRefereeStatsColumns.affiliateAddress, affiliateAddress);
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
      AffiliatePerRefereeStatsColumns.affiliateAddress,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  affiliatePerRefereeStatsToCreate: AffiliatePerRefereeStatsCreateObject,
  options: Options = { txId: undefined },
): Promise<AffiliatePerRefereeStatsFromDatabase> {
  return AffiliatePerRefereeStatsModel.query(
    Transaction.get(options.txId),
  ).insert(affiliatePerRefereeStatsToCreate).returning('*');
}

export async function upsert(
  affiliatePerRefereeStatsToUpsert: AffiliatePerRefereeStatsCreateObject,
  options: Options = { txId: undefined },
): Promise<AffiliatePerRefereeStatsFromDatabase> {
  const affiliatePerRefereeStats:
  AffiliatePerRefereeStatsModel[] = await AffiliatePerRefereeStatsModel.query(
    Transaction.get(options.txId),
  ).upsert(affiliatePerRefereeStatsToUpsert).returning('*');
  return affiliatePerRefereeStats[0];
}

export async function findById(
  affiliateAddress: string,
  refereeAddress: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliatePerRefereeStatsFromDatabase | undefined> {
  const baseQuery:
  QueryBuilder<AffiliatePerRefereeStatsModel> = setupBaseQuery<AffiliatePerRefereeStatsModel>(
    AffiliatePerRefereeStatsModel,
    options,
  );
  return baseQuery
    .findById([affiliateAddress, refereeAddress])
    .returning('*');
}

/**
 * Updates per-referee stats in the database based on the provided time window.
 *
 * This function aggregates per-referee affiliate stats and fill statistics
 * from various tables. Then it upserts the aggregated data into the `affiliate_info` table.
 *
 * @async
 * @function updateStats
 * @param {string} windowStartTs - The exclusive start timestamp for filtering fills.
 * @param {string} windowEndTs - The inclusive end timestamp for filtering fill.
 * @param {number} [txId] - Optional transaction ID.
 * @returns {Promise<void>}
 */
export async function updateStats(
  windowStartTs: string, // exclusive
  windowEndTs: string, // inclusive
  txId: number | undefined = undefined,
) : Promise<void> {
  const transaction: Knex.Transaction | undefined = Transaction.get(txId);

  // NEXT TODO
  const query = `
-- Aggregates stats per affiliate-referee tuple
WITH affiliate_metadata_per_referee AS (
  SELECT 
      affiliate_referred_users."affiliateAddress", 
      affiliate_referred_users."refereeAddress", 
      MIN("referredAtBlock") AS "firstReferralBlockHeight"
  FROM 
      affiliate_referred_users
  GROUP BY 
      affiliate_referred_users."affiliateAddress", 
      affiliate_referred_users."refereeAddress"
),

affiliate_referred_subaccounts AS (
  SELECT 
      affiliate_referred_users."affiliateAddress",
      affiliate_referred_users."refereeAddress",
      affiliate_referred_users."referredAtBlock",
      subaccounts."id" AS "subaccountId"
  FROM 
      affiliate_referred_users
  INNER JOIN 
      subaccounts
  ON 
      affiliate_referred_users."refereeAddress" = subaccounts."address"
),

filtered_fills AS (
  SELECT
      fills."subaccountId",
      fills."liquidity",
      fills."createdAt",
      CAST(fills."fee" AS decimal) AS "fee",
      fills."affiliateRevShare",
      fills."createdAtHeight",
      fills."price",
      fills."size",
      fills."type"
  FROM 
      fills
  WHERE 
      fills."createdAt" > '${windowStartTs}'
      AND fills."createdAt" <= '${windowEndTs}'
),

affiliate_fills_per_referee AS (
  SELECT
      filtered_fills."subaccountId",
      filtered_fills."liquidity",
      filtered_fills."createdAt",
      filtered_fills."fee",
      filtered_fills."affiliateRevShare",
      filtered_fills."price",
      filtered_fills."size",
      filtered_fills."type",
      affiliate_referred_subaccounts."affiliateAddress",
      affiliate_referred_subaccounts."refereeAddress",
      affiliate_referred_subaccounts."referredAtBlock"
  FROM 
      filtered_fills
  INNER JOIN
      affiliate_referred_subaccounts
  ON
      filtered_fills."subaccountId" = affiliate_referred_subaccounts."subaccountId"
  WHERE 
      filtered_fills."createdAtHeight" >= affiliate_referred_subaccounts."referredAtBlock"
),

affiliate_stats_per_referee AS (
  SELECT
      affiliate_fills_per_referee."affiliateAddress",
      affiliate_fills_per_referee."refereeAddress",
      SUM(affiliate_fills_per_referee."affiliateRevShare") AS "affiliateEarnings",
      SUM(CASE WHEN affiliate_fills_per_referee."liquidity" = '${Liquidity.MAKER}' AND affiliate_fills_per_referee."fee" > 0 THEN affiliate_fills_per_referee."fee" ELSE 0 END) AS "totalReferredMakerFees",
      SUM(CASE WHEN affiliate_fills_per_referee."liquidity" = '${Liquidity.TAKER}' AND affiliate_fills_per_referee."type" = '${FillType.LIMIT}' THEN affiliate_fills_per_referee."fee" ELSE 0 END) AS "totalReferredTakerFees",
      SUM(CASE WHEN affiliate_fills_per_referee."liquidity" = '${Liquidity.MAKER}' AND affiliate_fills_per_referee."fee" < 0 THEN affiliate_fills_per_referee."fee" ELSE 0 END) AS "totalReferredMakerRebates",
      SUM(CASE WHEN affiliate_fills_per_referee."liquidity" = '${Liquidity.TAKER}' AND affiliate_fills_per_referee."type" = '${FillType.LIQUIDATED}' THEN affiliate_fills_per_referee."fee" ELSE 0 END) AS "totalReferredLiquidationFees",
      COUNT(CASE WHEN affiliate_fills_per_referee."liquidity" = '${Liquidity.MAKER}' THEN 1 END) AS "referredMakerTrades",
      COUNT(CASE WHEN affiliate_fills_per_referee."liquidity" = '${Liquidity.TAKER}' THEN 1 END) AS "referredTakerTrades",
      SUM(affiliate_fills_per_referee."price" * affiliate_fills_per_referee."size") AS "referredTotalVolume"
  FROM
      affiliate_fills_per_referee
  GROUP BY
      affiliate_fills_per_referee."affiliateAddress",
      affiliate_fills_per_referee."refereeAddress"
),

affiliate_info_per_referee_update AS (
  SELECT
    affiliate_metadata_per_referee."affiliateAddress",
    affiliate_metadata_per_referee."refereeAddress",
    affiliate_metadata_per_referee."firstReferralBlockHeight",
    COALESCE(affiliate_stats_per_referee."affiliateEarnings", 0) AS "affiliateEarnings",
    COALESCE(affiliate_stats_per_referee."referredMakerTrades", 0) AS "referredMakerTrades",
    COALESCE(affiliate_stats_per_referee."referredTakerTrades", 0) AS "referredTakerTrades",
    COALESCE(affiliate_stats_per_referee."totalReferredMakerFees", 0) AS "totalReferredMakerFees",
    COALESCE(affiliate_stats_per_referee."totalReferredTakerFees", 0) AS "totalReferredTakerFees",
    COALESCE(affiliate_stats_per_referee."totalReferredMakerRebates", 0) AS "totalReferredMakerRebates",
    COALESCE(affiliate_stats_per_referee."referredTotalVolume", 0) AS "referredTotalVolume"
  FROM
    affiliate_metadata_per_referee
  LEFT JOIN
    affiliate_stats_per_referee
  ON 
    affiliate_metadata_per_referee."affiliateAddress" = affiliate_stats_per_referee."affiliateAddress"
    AND affiliate_metadata_per_referee."refereeAddress" = affiliate_stats_per_referee."refereeAddress"
)

INSERT INTO affiliate_per_referee_stats (
    "affiliateAddress",
    "refereeAddress",
    "affiliateEarnings",
    "referredMakerTrades",
    "referredTakerTrades",
    "referredTotalVolume",
    "firstReferralBlockHeight",
    "totalReferredTakerFees",
    "totalReferredMakerFees",
    "totalReferredMakerRebates"
)
SELECT
    "affiliateAddress",
    "refereeAddress",
    "affiliateEarnings",
    "referredMakerTrades",
    "referredTakerTrades",
    "referredTotalVolume",
    "firstReferralBlockHeight",
    "totalReferredTakerFees",
    "totalReferredMakerFees",
    "totalReferredMakerRebates"
FROM 
    affiliate_info_per_referee_update
ON CONFLICT ("affiliateAddress", "refereeAddress")
DO UPDATE SET
    "affiliateEarnings" = EXCLUDED."affiliateEarnings",
    "referredMakerTrades" = EXCLUDED."referredMakerTrades",
    "referredTakerTrades" = EXCLUDED."referredTakerTrades",
    "referredTotalVolume" = EXCLUDED."referredTotalVolume",
    "firstReferralBlockHeight" = EXCLUDED."firstReferralBlockHeight",
    "totalReferredTakerFees" = EXCLUDED."totalReferredTakerFees",
    "totalReferredMakerFees" = EXCLUDED."totalReferredMakerFees",
    "totalReferredMakerRebates" = EXCLUDED."totalReferredMakerRebates";

    `;

  return transaction
    ? knexPrimary.raw(query).transacting(transaction)
    : knexPrimary.raw(query);
}

/**
 * Given affiliate address, finds per-referee stats from the database with optional sorting,
 * and offset based pagination.
 *
 * @async
 * @function paginatedFindWithAffiliateAddress
 * @param {string} affiliateAddress - the affiliate address to query for.
 * @param {number} offset - The offset for pagination.
 * @param {number} limit - The maximum number of records to return.
 * @param {boolean} sortByPerRefereeEarning - Sort the results by per-referee-earnings.
 * @param {Options} [options=DEFAULT_POSTGRES_OPTIONS] - Optional config for database interaction.
 * @returns {Promise<AffiliatePerRefereeStatsFromDatabase[]>}
 */
export async function paginatedFindWithAffiliateAddress(
  affiliateAddress: string,
  offset: number,
  limit: number,
  sortByPerRefereeEarning: boolean,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliatePerRefereeStatsFromDatabase[]> {
  let baseQuery:
  QueryBuilder<AffiliatePerRefereeStatsModel> = setupBaseQuery<AffiliatePerRefereeStatsModel>(
    AffiliatePerRefereeStatsModel,
    options,
  );

  baseQuery = baseQuery.where(AffiliatePerRefereeStatsColumns.affiliateAddress, affiliateAddress);

  if (sortByPerRefereeEarning || offset !== 0) {
    baseQuery = baseQuery.orderBy(AffiliatePerRefereeStatsColumns.affiliateEarnings, Ordering.DESC)
      .orderBy(AffiliatePerRefereeStatsColumns.refereeAddress, Ordering.ASC);
  }

  // Apply pagination using offset and limit
  baseQuery = baseQuery.offset(offset).limit(limit);

  return baseQuery.returning('*');
}
