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
    id: uuid(oraclePriceToCreate.marketId, oraclePriceToCreate.effectiveAtHeight),
    ...oraclePriceToCreate,
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

async function findLatestPricesByDateTime(
  latestDateTimeISO: string,
): Promise<PriceMap> {
  const baseQuery: QueryBuilder<OraclePriceModel> = setupBaseQuery<OraclePriceModel>(
    OraclePriceModel,
    { readReplica: true },
  );

  const innerQuery: QueryBuilder<OraclePriceModel> = setupBaseQuery<OraclePriceModel>(
    OraclePriceModel,
    { readReplica: true },
  );

  const subQuery = innerQuery
    .select('marketId')
    .max('effectiveAt as maxEffectiveAt')
    .where('effectiveAt', '<=', latestDateTimeISO)
    .groupBy('marketId');

  const oraclePrices: OraclePriceFromDatabase[] = await baseQuery
    .innerJoin(subQuery.as('sub'), function () {
      this
        .on('oracle_prices.marketId', '=', 'sub.marketId')
        .andOn('oracle_prices.effectiveAt', '=', 'sub.maxEffectiveAt');
    })
    .returning('*');

  return constructPriceMap(oraclePrices);
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

export async function findLatestPrices(
  effectiveBeforeOrAtHeight: string,
  transaction?: Knex.Transaction,
): Promise<PriceMap> {
  const query: string = `
    SELECT "marketId", "price"
    FROM "oracle_prices"
    WHERE ("marketId", "effectiveAtHeight") IN (
      SELECT "marketId", MAX("effectiveAtHeight")
      FROM "oracle_prices"
      WHERE "effectiveAtHeight" <= '${effectiveBeforeOrAtHeight}'
      GROUP BY "marketId");
  `;
  let result: { rows: OraclePriceFromDatabase[] };
  if (transaction === undefined) {
    result = await knexReadReplica.getConnection().raw(
      query,
    ) as unknown as { rows: OraclePriceFromDatabase[] };
  } else {
    result = await knexReadReplica.getConnection().raw(
      query,
    ).transacting(transaction) as unknown as { rows: OraclePriceFromDatabase[] };
  }

  return constructPriceMap(result.rows);
}
