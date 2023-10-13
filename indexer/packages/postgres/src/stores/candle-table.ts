import _ from 'lodash';
import { PartialModelObject, QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
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
    id: uuid(candle.startedAt, candle.ticker, candle.resolution),
    ...candle,
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
  resolutions: CandleResolution[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<CandlesMap> {
  const candlesMap: CandlesMap = {};

  await Promise.all(
    _.map(
      tickers,
      async (ticker: string) => {
        candlesMap[ticker] = {};
        const findLatestCandles: Promise<CandleFromDatabase | undefined>[] = resolutions.map(
          (resolution: CandleResolution) => findLatest(
            ticker,
            resolution,
            options,
          ),
        );

        // Map each resolution to its respective candle
        const allLatestCandles: (CandleFromDatabase | undefined)[] = await Promise.all(
          findLatestCandles,
        );
        _.forEach(allLatestCandles, (candle: CandleFromDatabase | undefined) => {
          if (candle !== undefined) {
            candlesMap[ticker][candle.resolution] = candle;
          }
        });
      },
    ),
  );

  return candlesMap;
}
