import Knex from 'knex';
import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexPrimary } from '../helpers/knex';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import AffiliateRefereeStatsModel from '../models/affiliate-referee-stats-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  AffiliateRefereeStatsColumns,
  AffiliateRefereeStatsCreateObject,
  AffiliateRefereeStatsFromDatabase,
  Liquidity,
  FillType,
  AffiliateRefereeStatsQueryConfig,
} from '../types';

export async function findAll(
  {
    affiliateAddress,
    limit,
  }: AffiliateRefereeStatsQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateRefereeStatsFromDatabase[]> {
  verifyAllRequiredFields(
    {
      affiliateAddress,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery
  : QueryBuilder<AffiliateRefereeStatsModel> = setupBaseQuery<AffiliateRefereeStatsModel>(
    AffiliateRefereeStatsModel,
    options,
  );

  if (affiliateAddress) {
    baseQuery = baseQuery.where(AffiliateRefereeStatsColumns.affiliateAddress, affiliateAddress);
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
      AffiliateRefereeStatsColumns.affiliateAddress,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  AffiliateRefereeStatsToCreate: AffiliateRefereeStatsCreateObject,
  options: Options = { txId: undefined },
): Promise<AffiliateRefereeStatsFromDatabase> {
  return AffiliateRefereeStatsModel.query(
    Transaction.get(options.txId),
  ).insert(AffiliateRefereeStatsToCreate).returning('*');
}

export async function upsert(
  AffiliateRefereeStatsToUpsert: AffiliateRefereeStatsCreateObject,
  options: Options = { txId: undefined },
): Promise<AffiliateRefereeStatsFromDatabase> {
  const AffiliateRefereeStats:
  AffiliateRefereeStatsModel[] = await AffiliateRefereeStatsModel.query(
    Transaction.get(options.txId),
  ).upsert(AffiliateRefereeStatsToUpsert).returning('*');
  return AffiliateRefereeStats[0];
}

export async function findById(
  refereeAddress: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateRefereeStatsFromDatabase | undefined> {
  const baseQuery:
  QueryBuilder<AffiliateRefereeStatsModel> = setupBaseQuery<AffiliateRefereeStatsModel>(
    AffiliateRefereeStatsModel,
    options,
  );
  return baseQuery
    .findById(refereeAddress)
    .returning('*');
}

/**
 * Updates per-referee stats in the database based on the provided time window.
 *
 * This function aggregates per-referee affiliate stats and fill statistics
 * from various tables. Then it upserts the aggregated data into the `affiliate_referee_stats`
 * table.
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

  const query = `
-- Get metadata for all affiliates
-- Step 1: Get referal height for each affiliate-referee pair
WITH affiliate_metadata_per_referee AS (
  SELECT 
      affiliate_referred_users."affiliateAddress", 
      affiliate_referred_users."refereeAddress", 
--- There should be only one referredAtBlock for each affiliate-referee pair
      MIN("referredAtBlock") AS "referralBlockHeight"
  FROM 
      affiliate_referred_users
  GROUP BY 
      affiliate_referred_users."affiliateAddress", 
      affiliate_referred_users."refereeAddress"
),

-- Calculate per-referee fill related stats for affiliates
-- Step 2a: Inner join affiliate_referred_users with subaccounts
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

-- Step 2b: Filter fills by the given time window
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

-- Step 2c: Inner join filtered_fills with affiliate_referred_subaccounts and filter
affiliate_fills AS (
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

-- Step 2d: Aggregate stats per affiliate-referee tuple
affiliate_stats_per_referee AS (
  SELECT
      affiliate_fills."affiliateAddress",
      affiliate_fills."refereeAddress",
      SUM(affiliate_fills."fee") AS "totalReferredFees",
      SUM(affiliate_fills."affiliateRevShare") AS "affiliateEarnings",
      SUM(CASE WHEN affiliate_fills."liquidity" = '${Liquidity.MAKER}' AND affiliate_fills."fee" > 0 THEN affiliate_fills."fee" ELSE 0 END) AS "referredMakerFees",
      SUM(CASE WHEN affiliate_fills."liquidity" = '${Liquidity.TAKER}' AND affiliate_fills."type" = '${FillType.LIMIT}' THEN affiliate_fills."fee" ELSE 0 END) AS "referredTakerFees",
      SUM(CASE WHEN affiliate_fills."liquidity" = '${Liquidity.MAKER}' AND affiliate_fills."fee" < 0 THEN affiliate_fills."fee" ELSE 0 END) AS "referredMakerRebates",
      SUM(CASE WHEN affiliate_fills."liquidity" = '${Liquidity.TAKER}' AND affiliate_fills."type" = '${FillType.LIQUIDATED}' THEN affiliate_fills."fee" ELSE 0 END) AS "referredLiquidationFees",
      COUNT(CASE WHEN affiliate_fills."liquidity" = '${Liquidity.MAKER}' THEN 1 END) AS "referredMakerTrades",
      COUNT(CASE WHEN affiliate_fills."liquidity" = '${Liquidity.TAKER}' THEN 1 END) AS "referredTakerTrades",
      SUM(affiliate_fills."price" * affiliate_fills."size") AS "referredTotalVolume"
  FROM
      affiliate_fills
  GROUP BY
      affiliate_fills."affiliateAddress",
      affiliate_fills."refereeAddress"
),

-- Step 3a: Prepare data for updating or inserting into the affiliate_referee_stats table
-- Combine metadata with aggregated stats for each affiliate-referee pair
affiliate_referee_stats_update AS (
  SELECT
    affiliate_metadata_per_referee."affiliateAddress",
    affiliate_metadata_per_referee."refereeAddress",
    affiliate_metadata_per_referee."referralBlockHeight",
    COALESCE(affiliate_stats_per_referee."affiliateEarnings", 0) AS "affiliateEarnings",
    COALESCE(affiliate_stats_per_referee."referredMakerTrades", 0) AS "referredMakerTrades",
    COALESCE(affiliate_stats_per_referee."referredTakerTrades", 0) AS "referredTakerTrades",
    COALESCE(affiliate_stats_per_referee."referredMakerFees", 0) AS "referredMakerFees",
    COALESCE(affiliate_stats_per_referee."referredTakerFees", 0) AS "referredTakerFees",
    COALESCE(affiliate_stats_per_referee."referredLiquidationFees", 0) AS "referredLiquidationFees",
    COALESCE(affiliate_stats_per_referee."referredMakerRebates", 0) AS "referredMakerRebates",
    COALESCE(affiliate_stats_per_referee."referredTotalVolume", 0) AS "referredTotalVolume"
  FROM
    affiliate_metadata_per_referee
  LEFT JOIN
    affiliate_stats_per_referee
  ON 
    affiliate_metadata_per_referee."affiliateAddress" = affiliate_stats_per_referee."affiliateAddress"
    AND affiliate_metadata_per_referee."refereeAddress" = affiliate_stats_per_referee."refereeAddress"
)

-- Step 3b: Insert or update the affiliate_referee_stats table
-- Update existing rows with new data or insert new rows if they don't exist
INSERT INTO affiliate_referee_stats (
    "affiliateAddress",
    "refereeAddress",
    "referralBlockHeight",
    "affiliateEarnings", 
    "referredMakerTrades", 
    "referredTakerTrades", 
    "referredMakerFees",
    "referredTakerFees", 
    "referredLiquidationFees",
    "referredMakerRebates",
    "referredTotalVolume"
)
SELECT
    "affiliateAddress",
    "refereeAddress",
    "referralBlockHeight",
    "affiliateEarnings",
    "referredMakerTrades",
    "referredTakerTrades",
    "referredMakerFees",
    "referredTakerFees", 
    "referredLiquidationFees",
    "referredMakerRebates",
    "referredTotalVolume"
FROM 
    affiliate_referee_stats_update
ON CONFLICT ("refereeAddress")
DO UPDATE SET
    "referralBlockHeight" = EXCLUDED."referralBlockHeight",
    "affiliateEarnings" = affiliate_referee_stats."affiliateEarnings" + EXCLUDED."affiliateEarnings",
    "referredMakerTrades" = affiliate_referee_stats."referredMakerTrades" + EXCLUDED."referredMakerTrades",
    "referredTakerTrades" = affiliate_referee_stats."referredTakerTrades" + EXCLUDED."referredTakerTrades",
    "referredMakerFees" = affiliate_referee_stats."referredMakerFees" + EXCLUDED."referredMakerFees",
    "referredTakerFees" = affiliate_referee_stats."referredTakerFees" + EXCLUDED."referredTakerFees",
    "referredLiquidationFees" = affiliate_referee_stats."referredLiquidationFees" + EXCLUDED."referredLiquidationFees",
    "referredMakerRebates" = affiliate_referee_stats."referredMakerRebates" + EXCLUDED."referredMakerRebates",
    "referredTotalVolume" = affiliate_referee_stats."referredTotalVolume" + EXCLUDED."referredTotalVolume";
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
 * @returns {Promise<AffiliateRefereeStatsFromDatabase[]>}
 */
export async function paginatedFindWithAffiliateAddress(
  affiliateAddress: string,
  offset: number,
  limit: number,
  sortByPerRefereeEarning: boolean,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateRefereeStatsFromDatabase[]> {
  let baseQuery:
  QueryBuilder<AffiliateRefereeStatsModel> = setupBaseQuery<AffiliateRefereeStatsModel>(
    AffiliateRefereeStatsModel,
    options,
  );

  baseQuery = baseQuery.where(AffiliateRefereeStatsColumns.affiliateAddress, affiliateAddress);

  if (sortByPerRefereeEarning || offset !== 0) {
    baseQuery = baseQuery.orderBy(AffiliateRefereeStatsColumns.affiliateEarnings, Ordering.DESC)
      .orderBy(AffiliateRefereeStatsColumns.refereeAddress, Ordering.ASC);
  }

  // Apply pagination using offset and limit
  baseQuery = baseQuery.offset(offset).limit(limit);

  return baseQuery.returning('*');
}
