import Big from 'big.js';
import _ from 'lodash';
import { DateTime } from 'luxon';
import { QueryBuilder } from 'objection';

import {
  BUFFER_ENCODING_UTF_8,
  DEFAULT_POSTGRES_OPTIONS,
} from '../constants';
import { knexReadReplica } from '../helpers/knex';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import FillModel from '../models/fill-model';
import {
  CostOfFills,
  FillColumns,
  FillCreateObject,
  FillFromDatabase,
  FillQueryConfig,
  FillUpdateObject,
  Liquidity,
  Market24HourTradeVolumes,
  OpenSizeWithFundingIndex,
  Options,
  OrderedFillsWithFundingIndices,
  Ordering,
  OrderSide,
  PaginationFromDatabase,
  QueryableField,
  QueryConfig,
} from '../types';
import { findIdsForParentSubaccount } from './subaccount-table';

export function uuid(eventId: Buffer, liquidity: Liquidity): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(Buffer.from(`${eventId.toString('hex')}-${liquidity}`, BUFFER_ENCODING_UTF_8));
}

/**
 * Handles pagination and limit logic for fill queries
 * @param baseQuery The base query to apply pagination to
 * @param limit Maximum number of fills to return
 * @param page Page number
 * @returns Promise<PaginationFromDatabase<FillFromDatabase>>
 */
async function handleLimitAndPagination(
  baseQuery: QueryBuilder<FillModel>,
  limit?: number,
  page?: number,
): Promise<PaginationFromDatabase<FillFromDatabase>> {
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
    const count: { count?: string } = await query.clone().clearOrder().count({ count: '*' }).first() as unknown as { count?: string };

    query = query.offset(offset).limit(limit);

    return {
      results: await query.returning('*'),
      limit,
      offset,
      total: parseInt(count.count ?? '0', 10),
    };
  }

  // If no pagination, just apply the limit
  if (limit !== undefined) {
    query = query.limit(limit);
  }

  return {
    results: await query.returning('*'),
  };
}

export async function findAll(
  {
    limit,
    id,
    subaccountId,
    side,
    liquidity,
    type,
    includeTypes,
    excludeTypes,
    clobPairId,
    eventId,
    transactionHash,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    createdOnOrAfterHeight,
    createdOnOrAfter,
    clientMetadata,
    fee,
    page,
    parentSubaccount,
  }: FillQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<FillFromDatabase>> {
  if (subaccountId !== undefined && parentSubaccount !== undefined) {
    throw new Error('Cannot specify both subaccountId and parentSubaccount in order query');
  }

  verifyAllRequiredFields(
    {
      limit,
      id,
      subaccountId,
      side,
      liquidity,
      type,
      clobPairId,
      eventId,
      transactionHash,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdOnOrAfterHeight,
      createdOnOrAfter,
      clientMetadata,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<FillModel> = setupBaseQuery<FillModel>(
    FillModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(FillColumns.id, id);
  }

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(FillColumns.subaccountId, subaccountId);
  } else if (parentSubaccount !== undefined) {
    // PERFORMANCE CRITICAL: Resolve subaccountIds to concrete UUIDs before querying.
    // Using IN (subquery) causes Postgres to misestimate cardinality and scan millions
    // of rows. With explicit UUIDs, Postgres uses optimal index scans per subaccount.
    const subaccountIds = await findIdsForParentSubaccount(parentSubaccount);
    baseQuery = baseQuery.whereIn(FillColumns.subaccountId, subaccountIds);
  }

  if (side !== undefined) {
    baseQuery = baseQuery.where(FillColumns.side, side);
  }

  if (liquidity !== undefined) {
    baseQuery = baseQuery.where(FillColumns.liquidity, liquidity);
  }

  if (type !== undefined) {
    baseQuery = baseQuery.where(FillColumns.type, type);
  }

  if (includeTypes !== undefined && includeTypes.length > 0) {
    baseQuery = baseQuery.whereIn(FillColumns.type, includeTypes);
  }

  if (excludeTypes !== undefined && excludeTypes.length > 0) {
    baseQuery = baseQuery.whereNotIn(FillColumns.type, excludeTypes);
  }

  if (clobPairId !== undefined) {
    baseQuery = baseQuery.where(FillColumns.clobPairId, clobPairId);
  }

  if (eventId !== undefined) {
    baseQuery = baseQuery.where(FillColumns.eventId, eventId);
  }

  if (transactionHash !== undefined) {
    baseQuery = baseQuery.where(FillColumns.transactionHash, transactionHash);
  }

  if (createdBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      FillColumns.createdAtHeight,
      '<=',
      createdBeforeOrAtHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(FillColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdOnOrAfterHeight !== undefined) {
    baseQuery = baseQuery.where(
      FillColumns.createdAtHeight,
      '>=',
      createdOnOrAfterHeight,
    );
  }

  if (createdOnOrAfter !== undefined) {
    baseQuery = baseQuery.where(FillColumns.createdAt, '>=', createdOnOrAfter);
  }

  if (clientMetadata !== undefined) {
    baseQuery = baseQuery.where(FillColumns.clientMetadata, clientMetadata);
  }

  if (fee !== undefined) {
    baseQuery = baseQuery.where(FillColumns.fee, fee);
  }

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(
        column,
        order,
      );
    }
  }

  baseQuery = baseQuery.orderBy(
    FillColumns.createdAtHeight,
    Ordering.DESC,
  );

  return handleLimitAndPagination(baseQuery, limit, page);
}

export async function create(
  fillToCreate: FillCreateObject,
  options: Options = { txId: undefined },
): Promise<FillFromDatabase> {
  return FillModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...fillToCreate,
    id: uuid(fillToCreate.eventId, fillToCreate.liquidity),
  }).returning('*');
}

export async function update(
  {
    ...fields
  }: FillUpdateObject,
  options: Options = { txId: undefined },
): Promise<FillFromDatabase | undefined> {
  const fill = await FillModel.query(
    Transaction.get(options.txId),
    // TODO fix expression typing so we dont have to use any
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ).findById(fields.id).patch(fields as any).returning('*');
  // The objection types mistakenly think the query returns an array of fills.
  return fill as unknown as (FillFromDatabase | undefined);
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<FillFromDatabase | undefined> {
  const baseQuery: QueryBuilder<FillModel> = setupBaseQuery<FillModel>(
    FillModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function get24HourInformation(perpetualMarketIds: string[]): Promise<
  _.Dictionary<Market24HourTradeVolumes>
> {
  // Rename perpetualMarketId to clobPairId
  const oneDayAgo: string = DateTime.local().minus({ days: 1 }).toISO();
  const result: {
    rows: Market24HourTradeVolumes[],
  } = await knexReadReplica.getConnection().raw(
    `SELECT
      "clobPairId",
      COUNT("createdAt") AS "trades24H",
      SUM("quoteAmount") AS "volume24H"
    FROM fills
    WHERE
      "liquidity"='TAKER'
      AND "createdAt">'${oneDayAgo}'
      AND type IN (
        'LIMIT',
        'MARKET',
        'STOP_LIMIT',
        'STOP_MARKET',
        'TRAILING_STOP',
        'TAKE_PROFIT',
        'TAKE_PROFIT_MARKET'
      )
      GROUP BY "clobPairId";
    `,
  ) as unknown as { rows: Market24HourTradeVolumes[] };

  const perpetualMarketStats24hr: _.Dictionary<Market24HourTradeVolumes> = _.keyBy(
    result.rows,
    FillColumns.clobPairId,
  );
  _.forEach(perpetualMarketIds, (clobPairId: string) => {
    if (!perpetualMarketStats24hr[clobPairId]) {
      // no fills exist for this market
      perpetualMarketStats24hr[clobPairId] = {
        clobPairId,
        trades24H: '0',
        volume24H: '0',
      };
    }
  });
  return perpetualMarketStats24hr;
}

/**
 * Returns the cost of all fills for a given subaccount up to a given height.
 *
 * If the subaccount has spent $500 to buy 1 ETH, cost of fills will be -500.
 * If the subaccount shorted 1 ETH at $500, cost of fills will be 500.
 *
 * @param subaccountId
 * @param createdAtHeight
 */
export async function getCostOfFills(
  subaccountId: string,
  createdAtHeight: string,
): Promise<Big> {
  const result: { rows: CostOfFills[] } = await knexReadReplica
    .getConnection()
    .raw(
      `
      SELECT SUM(CASE 
                  WHEN side = 'SELL' THEN price * size 
                  ELSE -1 * price * size 
                END) AS "cost"
      FROM fills 
      WHERE "subaccountId" = ?
        AND "createdAtHeight" <= ?;
      `,
      [subaccountId, createdAtHeight],
    );

  const pnlOfFills = result.rows[0]?.cost || 0;
  return Big(pnlOfFills);
}

/**
 * Returns the total value of all open positions for a given subaccount at a given height.
 *
 * @param subaccountId
 * @param createdAtHeight
 */
export async function getTotalValueOfOpenPositions(
  subaccountId: string,
  createdAtHeight: string,
): Promise<Big> {
  const query = `
    SELECT
      SUM(f.open_size * p.price) AS total_sum
    FROM
      (
        SELECT
          "clobPairId",
          SUM(CASE
            WHEN side = 'SELL' THEN -1 * size
            ELSE size
          END) AS open_size
        FROM
          fills
        WHERE
          "subaccountId" = ?
          AND "createdAtHeight" <= ?
        GROUP BY "clobPairId"
      ) AS f
      JOIN (
        SELECT "clobPairId", "marketId"
        FROM "perpetual_markets"
      ) AS pm
      ON f."clobPairId" = pm."clobPairId"
      JOIN (
        SELECT DISTINCT ON ("marketId") *
        FROM "oracle_prices"
        WHERE "effectiveAtHeight" <= ?
        ORDER BY "marketId", "effectiveAtHeight" DESC
      ) AS p
      ON p."marketId" = pm."marketId";
  `;

  const result = await knexReadReplica
    .getConnection()
    .raw(query, [subaccountId, createdAtHeight, createdAtHeight]);

  const totalSum = result.rows[0]?.total_sum || '0';
  return new Big(totalSum);
}

/**
 * Returns the ordered fills with funding indices for a given clob pair and subaccount
 * prior to a given height.
 *
 * Fills will be paired with the last fill that occurred before it, and contain the latest
 * funding index at the time of the fill.
 *
 * @param clobPairId
 * @param subaccountId
 * @param effectiveBeforeHeight
 */
export async function getOrderedFillsWithFundingIndices(
  clobPairId: string,
  subaccountId: string,
  effectiveBeforeHeight: string,
): Promise<OrderedFillsWithFundingIndices[]> {
  const result: { rows: OrderedFillsWithFundingIndices[] } = await knexReadReplica
    .getConnection()
    .raw(
      `
      WITH input AS (
        SELECT
          f.*,
          fiu."fundingIndex",
          LAG(f."id") OVER (PARTITION BY f."subaccountId", f."clobPairId" ORDER BY f."createdAtHeight") AS last_fill_id
        FROM
          "fills" f
        JOIN
          "perpetual_markets" pm ON f."clobPairId" = pm."clobPairId"
        JOIN
          "funding_index_updates" fiu ON pm."id" = fiu."perpetualId"
        WHERE
          f."subaccountId" = ?
          AND f."clobPairId" = ?
          AND f."createdAtHeight" <= ?
          AND fiu."effectiveAtHeight" = (
            SELECT
              MAX("effectiveAtHeight")
            FROM
              "funding_index_updates"
            WHERE
              "perpetualId" = pm."id"
              AND "effectiveAtHeight" <= f."createdAtHeight"
          )
        ORDER BY
          "createdAtHeight" ASC
      )
      SELECT
        current_fill."id" as "id",
        current_fill."subaccountId" as "subaccountId",
        current_fill."side" as "side",
        current_fill."size" as "size",
        current_fill."createdAtHeight" as "createdAtHeight",
        current_fill."fundingIndex" as "fundingIndex",
        last_fill."id" as "lastFillId",
        last_fill."side" as "lastFillSide",
        last_fill."size" as "lastFillSize",
        last_fill."createdAtHeight" as "lastFillCreatedAtHeight",
        last_fill."fundingIndex" as "lastFillFundingIndex"
      FROM
        input current_fill
      LEFT JOIN input last_fill ON current_fill."last_fill_id" = last_fill."id"
      where current_fill."last_fill_id" is not null
      ORDER BY
        current_fill."createdAtHeight" ASC;
      `,
      [subaccountId, clobPairId, effectiveBeforeHeight],
    );

  return result.rows;
}

/**
 * Returns the paid funding for a given set of ordered fills with funding indices.
 *
 * @param orderedFillsWithFundingIndices
 */
export function getSettledFunding(
  orderedFillsWithFundingIndices: OrderedFillsWithFundingIndices[],
): Big {
  if (orderedFillsWithFundingIndices.length === 0) {
    return Big(0);
  }
  let paidFunding: Big = Big(0);
  let currentSize: Big = Big(orderedFillsWithFundingIndices[0].lastFillSize);
  for (const fill of orderedFillsWithFundingIndices) {
    if (fill.fundingIndex !== fill.lastFillFundingIndex) {
      const currFunding: Big = Big(fill.fundingIndex)
        .minus(fill.lastFillFundingIndex).times(currentSize);
      paidFunding = paidFunding.plus(currFunding);
    }
    currentSize = currentSize.plus(fill.side === OrderSide.BUY ? fill.size : -fill.size);
  }
  return paidFunding;
}

/**
 * Returns the open size with funding index for a given subaccount prior to a given height,
 * for each clob pair.
 *
 * This will be used to compute unsettled funding payments.
 *
 * @param subaccountId
 * @param effectiveBeforeHeight
 */
export async function getOpenSizeWithFundingIndex(
  subaccountId: string,
  effectiveBeforeHeight: string,
): Promise<OpenSizeWithFundingIndex[]> {
  const result: { rows: OpenSizeWithFundingIndex[] } = await knexReadReplica
    .getConnection()
    .raw(
      `
      WITH input AS (
        SELECT 
          f."clobPairId",
          SUM(CASE
            WHEN side = 'SELL' THEN -1 * size
            ELSE size
          END) AS "openSize",
          MAX("createdAtHeight") as "lastFillHeight"
        FROM
          "fills" f
        WHERE
          f."subaccountId" = ?
          AND f."createdAtHeight" <= ?
        GROUP BY
          f."clobPairId"
      )
      SELECT 
        input.*,
        fiu."fundingIndex",
        fiu."effectiveAtHeight" as "fundingIndexHeight"
      FROM
        input
      JOIN
        "perpetual_markets" pm ON input."clobPairId" = pm."clobPairId"
      JOIN
        "funding_index_updates" fiu ON pm."id" = fiu."perpetualId"
      WHERE
        fiu."effectiveAtHeight" = (
          SELECT
            MAX("effectiveAtHeight")
          FROM
            "funding_index_updates"
          WHERE
            "perpetualId" = pm."id"
            AND "effectiveAtHeight" <= input."lastFillHeight"
        )
      ORDER BY input."clobPairId" ASC;
      `,
      [subaccountId, effectiveBeforeHeight],
    );

  return result.rows;
}

/**
 * Returns the unique clob pair ids for a given subaccount's fills prior to a given height.
 *
 * @param subaccountId
 * @param effectiveBeforeOrAtHeight
 * @param options
 */
export async function getClobPairs(
  subaccountId: string,
  effectiveBeforeOrAtHeight: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string[]> {
  const baseQuery: QueryBuilder<FillModel> = setupBaseQuery<
    FillModel>(
      FillModel,
      options,
    );
  const fills: FillFromDatabase[] = await baseQuery
    .distinctOn(FillColumns.clobPairId)
    .where(FillColumns.createdAtHeight, '<=', effectiveBeforeOrAtHeight)
    .where(FillColumns.subaccountId, subaccountId)
    .orderBy(FillColumns.clobPairId, Ordering.ASC)
    .returning('*');
  return _.map(fills, (fill: FillFromDatabase) => fill.clobPairId);
}

/**
 * Returns the fees paid for a given subaccount prior to a given height.
 *
 * @param subaccountId
 * @param effectiveBeforeOrAtHeight
 * @param options
 */
export async function getFeesPaid(
  subaccountId: string,
  effectiveBeforeOrAtHeight: string,
): Promise<Big> {
  const result: { rows: { feesPaid: string }[] } = await knexReadReplica
    .getConnection()
    .raw(
      `
      SELECT
        SUM(CAST(f."fee" AS NUMERIC)) as "feesPaid"
      FROM
        "fills" f
      WHERE
        f."subaccountId" = ?
        AND f."createdAtHeight" <= ?
      `,
      [subaccountId, effectiveBeforeOrAtHeight],
    );

  if (result.rows[0].feesPaid === null) {
    return Big(0);
  }

  return Big(result.rows[0].feesPaid);
}
