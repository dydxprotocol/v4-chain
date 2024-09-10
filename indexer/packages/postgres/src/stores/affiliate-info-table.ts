import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexReadReplica } from '../helpers/knex';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import AffiliateInfoModel from '../models/affiliate-info-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  AffiliateInfoColumns,
  AffiliateInfoCreateObject,
  AffiliateInfoFromDatabase,
  AffiliateInfoQueryConfig,
  Liquidity,
  PersistentCacheKeys,
} from '../types';

export async function findAll(
  {
    address,
    limit,
  }: AffiliateInfoQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateInfoFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<AffiliateInfoModel> = setupBaseQuery<AffiliateInfoModel>(
    AffiliateInfoModel,
    options,
  );

  if (address) {
    baseQuery = baseQuery.where(AffiliateInfoColumns.address, address);
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
      AffiliateInfoColumns.address,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  AffiliateInfoToCreate: AffiliateInfoCreateObject,
  options: Options = { txId: undefined },
): Promise<AffiliateInfoFromDatabase> {
  return AffiliateInfoModel.query(
    Transaction.get(options.txId),
  ).insert(AffiliateInfoToCreate).returning('*');
}

export async function upsert(
  AffiliateInfoToUpsert: AffiliateInfoCreateObject,
  options: Options = { txId: undefined },
): Promise<AffiliateInfoFromDatabase> {
  const AffiliateInfos: AffiliateInfoModel[] = await AffiliateInfoModel.query(
    Transaction.get(options.txId),
  ).upsert(AffiliateInfoToUpsert).returning('*');
  // should only ever be one AffiliateInfo
  return AffiliateInfos[0];
}
export async function findById(
  address: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateInfoFromDatabase | undefined> {
  const baseQuery: QueryBuilder<AffiliateInfoModel> = setupBaseQuery<AffiliateInfoModel>(
    AffiliateInfoModel,
    options,
  );
  return baseQuery
    .findById(address)
    .returning('*');
}

export async function updateInfo(
  windowStartTs: string, // exclusive
  windowEndTs: string, // inclusive
) : Promise<void> {

  await knexReadReplica.getConnection().raw(
    `
BEGIN;

-- Get metadata for all affiliates
-- STEP 1: Aggregate affiliate_referred_users
WITH affiliate_metadata AS (
  SELECT 
      "affiliateAddress", 
      COUNT(*) AS "totalReferredUsers",
      MIN("referredAtBlock") AS "firstReferralBlockHeight"
  FROM 
      affiliate_referred_users
  GROUP BY 
      "affiliateAddress"
),

-- Calculate fill related stats for affiliates
-- Step 2a: Inner join affiliate_referred_users with subaccounts to get subaccounts referred by the affiliate
affiliate_referred_subaccounts AS (
  SELECT 
      affiliate_referred_users."affiliateAddress",
      affiliate_referred_users."referredAtBlock",
      subaccounts."id"
  FROM 
      affiliate_referred_users
  INNER JOIN 
      subaccounts
  ON 
      affiliate_referred_users."refereeAddress" = subaccounts."address"
),

-- Step 2b: Filter fills by time window
filtered_fills AS (
  SELECT
      fills."subaccountId",
      fills."liquidity",
      fills."createdAt",
      CAST(fills."fee" AS decimal) AS "fee",
      fills."affiliateRevShare",
      fills."createdAtHeight",
      fills."price",
      fills."size"
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
      affiliate_referred_subaccounts."affiliateAddress",
      affiliate_referred_subaccounts."referredAtBlock"
  FROM 
      filtered_fills
  INNER JOIN
      affiliate_referred_subaccounts
  ON
      filtered_fills."subaccountId" = affiliate_referred_subaccounts."id"
  WHERE 
      filtered_fills."createdAtHeight" >= affiliate_referred_subaccounts."referredAtBlock"
),

-- Step 2d: Groupby to get affiliate level stats
affiliate_stats AS (
  SELECT
      affiliate_fills."affiliateAddress",
      SUM(affiliate_fills."fee") AS "totalReferredFees",
      SUM(affiliate_fills."affiliateRevShare") AS "affiliateEarnings",
      SUM(affiliate_fills."fee") - SUM(affiliate_fills."affiliateRevShare") AS "referredNetProtocolEarnings",
      COUNT(CASE WHEN affiliate_fills."liquidity" = '${Liquidity.MAKER}' THEN 1 END) AS "referredMakerTrades",
      COUNT(CASE WHEN affiliate_fills."liquidity" = '${Liquidity.TAKER}' THEN 1 END) AS "referredTakerTrades",
      SUM(affiliate_fills."price" * affiliate_fills."size") AS "totalReferredVolume"
  FROM
      affiliate_fills
  GROUP BY
      affiliate_fills."affiliateAddress"
),

-- Prepare to update affiliate_info
-- STEP 3a: Left join affiliate_stats onto affiliate_metadata. affiliate_stats only has values for
-- addresses with fills in the time window
affiliate_info_update AS (
  SELECT
    affiliate_metadata."affiliateAddress",
    affiliate_metadata."totalReferredUsers",    
    affiliate_metadata."firstReferralBlockHeight",
    COALESCE(affiliate_stats."totalReferredFees", 0) AS "totalReferredFees",
    COALESCE(affiliate_stats."affiliateEarnings", 0) AS "affiliateEarnings",
    COALESCE(affiliate_stats."referredNetProtocolEarnings", 0) AS "referredNetProtocolEarnings",
    COALESCE(affiliate_stats."referredMakerTrades", 0) AS "referredMakerTrades",
    COALESCE(affiliate_stats."referredTakerTrades", 0) AS "referredTakerTrades",
    COALESCE(affiliate_stats."totalReferredVolume", 0) AS "totalReferredVolume"
  FROM
    affiliate_metadata
  LEFT JOIN
    affiliate_stats
  ON affiliate_metadata."affiliateAddress" = affiliate_stats."affiliateAddress"
)

-- Step 3b: Update/upsert the affiliate info table with the new stats
INSERT INTO affiliate_info (
    "address", 
    "totalReferredUsers",
    "firstReferralBlockHeight",
    "affiliateEarnings", 
    "referredMakerTrades", 
    "referredTakerTrades", 
    "totalReferredFees", 
    "referredNetProtocolEarnings",
    "totalReferredVolume"
)
SELECT
    "affiliateAddress",
    "totalReferredUsers",
    "firstReferralBlockHeight",
    "affiliateEarnings",
    "referredMakerTrades",
    "referredTakerTrades",
    "totalReferredFees",
    "referredNetProtocolEarnings",
    "totalReferredVolume"
FROM 
    affiliate_info_update
ON CONFLICT ("address")
DO UPDATE SET
    "totalReferredUsers" = EXCLUDED."totalReferredUsers",
    "firstReferralBlockHeight" = EXCLUDED."firstReferralBlockHeight",
    "affiliateEarnings" = affiliate_info."affiliateEarnings" + EXCLUDED."affiliateEarnings",
    "referredMakerTrades" = affiliate_info."referredMakerTrades" + EXCLUDED."referredMakerTrades",
    "referredTakerTrades" = affiliate_info."referredTakerTrades" + EXCLUDED."referredTakerTrades",
    "totalReferredFees" = affiliate_info."totalReferredFees" + EXCLUDED."totalReferredFees",
    "referredNetProtocolEarnings" = affiliate_info."referredNetProtocolEarnings" + EXCLUDED."referredNetProtocolEarnings",
    "totalReferredVolume" = affiliate_info."totalReferredVolume" + EXCLUDED."totalReferredVolume";

-- Step 5: Upsert new affiliateInfoLastUpdateTime to persistent_cache table
INSERT INTO persistent_cache (key, value)
VALUES ('${PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME}', '${windowEndTs}')
ON CONFLICT (key) 
DO UPDATE SET value = EXCLUDED.value;

COMMIT;
    `,
  );
}
