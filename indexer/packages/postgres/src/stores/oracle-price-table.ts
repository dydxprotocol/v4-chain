import Knex from 'knex';
import _ from 'lodash';
import { DateTime } from 'luxon';
import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexReadReplica } from '../helpers/knex';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import OraclePriceModel from '../models/oracle-price-model';
import {
  Options,
  OraclePriceColumns,
  OraclePriceCreateObject,
  OraclePriceFromDatabase,
  OraclePriceQueryConfig,
  Ordering,
  PriceMap,
  QueryableField,
  QueryConfig,
} from '../types';

export function uuid(
  marketId: number, height: string,
): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(Buffer.from(`${marketId.toString()}-${height}`, BUFFER_ENCODING_UTF_8));
}

export async function findAll(
  {
    limit,
    id,
    marketId,
    price,
    effectiveAt,
    effectiveAtHeight,
    effectiveBeforeOrAt,
    effectiveBeforeOrAtHeight,
  }: OraclePriceQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<OraclePriceFromDatabase[]> {
  verifyAllRequiredFields(
    {
      limit,
      id,
      marketId,
      price,
      effectiveAt,
      effectiveAtHeight,
      effectiveBeforeOrAt,
      effectiveBeforeOrAtHeight,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<OraclePriceModel> = setupBaseQuery<OraclePriceModel>(
    OraclePriceModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(OraclePriceColumns.id, id);
  }
  if (marketId !== undefined) {
    baseQuery = baseQuery.whereIn(OraclePriceColumns.marketId, marketId);
  }
  if (price !== undefined) {
    baseQuery = baseQuery.whereIn(OraclePriceColumns.price, price);
  }
  if (effectiveAt !== undefined) {
    baseQuery = baseQuery.where(OraclePriceColumns.effectiveAt, effectiveAt);
  }
  if (effectiveAtHeight !== undefined) {
    baseQuery = baseQuery.where(OraclePriceColumns.effectiveAtHeight, effectiveAtHeight);
  }

  if (effectiveBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(OraclePriceColumns.effectiveAt, '<=', effectiveBeforeOrAt);
  }

  if (effectiveBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      OraclePriceColumns.effectiveAtHeight,
      '<=',
      effectiveBeforeOrAtHeight,
    );
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
      OraclePriceColumns.effectiveAtHeight,
      Ordering.DESC,
    );
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  oraclePriceToCreate: OraclePriceCreateObject,
  options: Options = { txId: undefined },
): Promise<OraclePriceFromDatabase> {
  return OraclePriceModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...oraclePriceToCreate,
    id: uuid(oraclePriceToCreate.marketId, oraclePriceToCreate.effectiveAtHeight),
  }).returning('*');
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<OraclePriceFromDatabase | undefined> {
  const baseQuery: QueryBuilder<OraclePriceModel> = setupBaseQuery<OraclePriceModel>(
    OraclePriceModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function findOraclePricesInReverseChronologicalOrder(
  marketId: number,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<OraclePriceFromDatabase[] | undefined> {
  const baseQuery: QueryBuilder<OraclePriceModel> = setupBaseQuery<OraclePriceModel>(
    OraclePriceModel,
    options,
  );

  return baseQuery
    .where(OraclePriceColumns.marketId, marketId)
    .orderBy(OraclePriceColumns.effectiveAtHeight, Ordering.DESC)
    .returning('*');
}

export async function findMostRecentMarketOraclePrice(
  marketId: number,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<OraclePriceFromDatabase | undefined> {
  const baseQuery: QueryBuilder<OraclePriceModel> = setupBaseQuery<OraclePriceModel>(
    OraclePriceModel,
    options,
  );

  const oraclePrices: OraclePriceFromDatabase[] = await baseQuery
    .where(OraclePriceColumns.marketId, marketId)
    .orderBy(OraclePriceColumns.effectiveAtHeight, Ordering.DESC)
    .limit(1)
    .returning('*');

  if (oraclePrices.length === 0) {
    return undefined;
  }
  return oraclePrices[0];
}

function constructPriceMap(oraclePrices: OraclePriceFromDatabase[]): PriceMap {
  return _.reduce(oraclePrices, (acc: PriceMap, oraclePrice: OraclePriceFromDatabase) => {
    acc[oraclePrice.marketId] = oraclePrice.price;
    return acc;
  }, {});
}

export async function findLatestPricesByDateTime(
  latestDateTimeISO: string,
): Promise<PriceMap> {
  // Use raw query with LEFT JOIN LATERAL for better performance.
  // This query enables Postgres to utilize the index on the effectiveAt column
  // for individual markets.
  const query = `
    SELECT m.id AS "marketId",
           p."price",
           p."effectiveAt",
           p."effectiveAtHeight",
           p."id"
    FROM "markets" m
    LEFT JOIN LATERAL (
      SELECT "id", "price", "effectiveAt", "effectiveAtHeight"
      FROM "oracle_prices"
      WHERE "marketId" = m.id
      AND "effectiveAt" <= ?
      ORDER BY "effectiveAt" DESC
      LIMIT 1
    ) p ON TRUE
    WHERE p."price" IS NOT NULL
  `;

  const result = await knexReadReplica.getConnection().raw(
    query,
    [latestDateTimeISO],
  ) as unknown as { rows: OraclePriceFromDatabase[] };

  return constructPriceMap(result.rows);
}

export async function getPricesFrom24hAgo(
): Promise<PriceMap> {
  const oneDayAgo: string = DateTime.utc().minus({ days: 1 }).toISO();
  return findLatestPricesByDateTime(oneDayAgo);
}

export async function getLatestPrices(): Promise<PriceMap> {
  const now: string = DateTime.utc().toISO();
  return findLatestPricesByDateTime(now);
}

export async function findLatestPricesBeforeOrAtHeight(
  effectiveBeforeOrAtHeight: string,
  transaction?: Knex.Transaction,
): Promise<PriceMap> {
  const query: string = `
    SELECT "marketId", "price"
    FROM "oracle_prices"
    WHERE ("marketId", "effectiveAtHeight") IN (
      SELECT "marketId", MAX("effectiveAtHeight")
      FROM "oracle_prices"
      WHERE "effectiveAtHeight" <= ?
      GROUP BY "marketId");
  `;
  let result: { rows: OraclePriceFromDatabase[] };
  if (transaction === undefined) {
    result = await knexReadReplica.getConnection().raw(
      query,
      [effectiveBeforeOrAtHeight],
    ) as unknown as { rows: OraclePriceFromDatabase[] };
  } else {
    result = await knexReadReplica.getConnection().raw(
      query,
      [effectiveBeforeOrAtHeight],
    ).transacting(transaction) as unknown as { rows: OraclePriceFromDatabase[] };
  }

  return constructPriceMap(result.rows);
}
