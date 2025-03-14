import { PartialModelObject, QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexReadReplica } from '../helpers/knex';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import CandleModel from '../models/candle-model';
import {
  CandleColumns,
  CandleCreateObject,
  CandleFromDatabase,
  CandleQueryConfig,
  CandleResolution,
  CandlesMap,
  CandleUpdateObject,
  IsoString,
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
} from '../types';

export function uuid(startedAt: IsoString, ticker: string, resolution: CandleResolution): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(Buffer.from(`${startedAt}-${ticker}-${resolution}`, BUFFER_ENCODING_UTF_8));
}

/**
 * Find all candles that match the given query config. fromIso is inclusive and toISO is exclusive.
 */
export async function findAll(
  {
    limit,
    id,
    ticker,
    resolution,
    fromISO,
    toISO,
  }: CandleQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<CandleFromDatabase[]> {
  verifyAllRequiredFields(
    {
      id,
      ticker,
      resolution,
      fromISO,
      toISO,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<CandleModel> = setupBaseQuery<CandleModel>(
    CandleModel,
    options,
  );

  if (id) {
    baseQuery = baseQuery.whereIn(CandleColumns.id, id);
  }

  if (ticker) {
    baseQuery = baseQuery.whereIn(CandleColumns.ticker, ticker);
  }

  if (resolution) {
    baseQuery = baseQuery.where(CandleColumns.resolution, resolution);
  }

  if (fromISO) {
    baseQuery = baseQuery.where(CandleColumns.startedAt, '>=', fromISO);
  }

  if (toISO) {
    baseQuery = baseQuery.where(CandleColumns.startedAt, '<', toISO);
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
      CandleColumns.ticker,
      Ordering.DESC,
    ).orderBy(
      CandleColumns.resolution,
      Ordering.DESC,
    ).orderBy(
      CandleColumns.startedAt,
      Ordering.DESC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  candle: CandleCreateObject,
  options: Options = { txId: undefined },
): Promise<CandleFromDatabase> {
  return CandleModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...candle,
    id: uuid(candle.startedAt, candle.ticker, candle.resolution),
  }).returning('*');
}

export async function update(
  {
    id,
    ...fields
  }: CandleUpdateObject,
  options: Options = { txId: undefined },
): Promise<CandleFromDatabase | undefined> {
  const candle = await CandleModel.query(
    Transaction.get(options.txId),
    // TODO fix expression typing so we dont have to use any
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ).findById(id);
  const updatedCandle = await candle.$query().patch(fields as PartialModelObject<CandleModel>).returning('*');
  // The objection types mistakenly think the query returns an array of candles.
  return updatedCandle as unknown as (CandleFromDatabase | undefined);
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<CandleFromDatabase | undefined> {
  const baseQuery: QueryBuilder<CandleModel> = setupBaseQuery<CandleModel>(
    CandleModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function findLatest(
  ticker: string,
  resolution: CandleResolution,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<CandleFromDatabase | undefined> {
  const candles: CandleFromDatabase[] = await findAll(
    {
      ticker: [ticker],
      resolution,
      limit: 1,
    },
    [QueryableField.TICKER, QueryableField.RESOLUTION],
    {
      ...options,
      // Ordered by 'startedAt' descending, so the most recent candle will be only candle returned
      orderBy: [[CandleColumns.startedAt, Ordering.DESC]],
    },
  );

  if (candles.length === 1) {
    return candles[0];
  }

  return undefined;
}

export async function findCandlesMap(
  tickers: string[],
  latestBlockTime: IsoString,
): Promise<CandlesMap> {
  if (tickers.length === 0) {
    return {};
  }

  const candlesMap: CandlesMap = {};
  for (const ticker of tickers) {
    candlesMap[ticker] = {};
  }

  const minuteCandlesResult: {
    rows: CandleFromDatabase[],
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT DISTINCT ON (
      ticker,
      resolution
    ) candles.* FROM 
      candles
    WHERE
      "ticker" IN (${tickers.map((ticker) => { return `'${ticker}'`; }).join(',')}) AND
      "startedAt" > ?::timestamptz - INTERVAL '3 hours' AND
      resolution IN ('1MIN', '5MINS', '15MINS', '30MINS', '1HOUR')
    ORDER BY
      ticker,
      resolution,
      "startedAt" DESC;
    `,
    [latestBlockTime],
  ) as unknown as {
    rows: CandleFromDatabase[],
  };

  const hourDayCandlesResult: {
    rows: CandleFromDatabase[],
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT DISTINCT ON (
      ticker,
      resolution
    ) candles.* FROM 
      candles
    WHERE
      "ticker" IN (${tickers.map((ticker) => { return `'${ticker}'`; }).join(',')}) AND
      "startedAt" > ?::timestamptz - INTERVAL '2 days' AND
      resolution IN ('4HOURS', '1DAY')
    ORDER BY
      ticker,
      resolution,
      "startedAt" DESC;
    `,
    [latestBlockTime],
  ) as unknown as {
    rows: CandleFromDatabase[],
  };

  const latestCandles: CandleFromDatabase[] = minuteCandlesResult.rows
    .concat(hourDayCandlesResult.rows);
  for (const candle of latestCandles) {
    if (candlesMap[candle.ticker] === undefined) {
      candlesMap[candle.ticker] = {};
    }
    candlesMap[candle.ticker][candle.resolution] = candle;
  }

  return candlesMap;
}

/**
 * Find all candles for a given resolution within a lookback period.
 * Uses Objection.js query builder for type safety and query construction.
 * @param resolution - The candle resolution to query for
 * @param lookbackMs - Number of milliseconds to look back from now
 * @param options - Query options
 */
export async function findByResAndLookbackPeriod(
  resolution: CandleResolution,
  lookbackMs: number,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<CandleFromDatabase[]> {
  const baseQuery: QueryBuilder<CandleModel> = setupBaseQuery<CandleModel>(
    CandleModel,
    options,
  );

  // Calculate lookback time in JS instead of using NOW() - interval in SQL.
  // This allows the query planner to prove the startedAt condition will match
  // candles_resolution_started_at_1_4_hour_feb2025_idx's requirements
  const lookbackTime: IsoString = new Date(Date.now() - lookbackMs).toISOString();

  return baseQuery
    .where(CandleColumns.resolution, resolution)
    .where(CandleColumns.startedAt, '>=', lookbackTime)
    .orderBy(CandleColumns.startedAt, Ordering.DESC)
    .returning('*');
}
