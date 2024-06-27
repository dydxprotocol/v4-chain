import _ from 'lodash';
import { PartialModelObject, QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import MarketModel from '../models/market-model';
import {
  MarketColumns,
  MarketCreateObject,
  MarketFromDatabase,
  MarketQueryConfig,
  MarketsMap,
  MarketUpdateObject,
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
} from '../types';

export async function findAll(
  {
    limit,
    id,
    pair,
  }: MarketQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<MarketFromDatabase[]> {
  verifyAllRequiredFields(
    {
      id,
      pair,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<MarketModel> = setupBaseQuery<MarketModel>(
    MarketModel,
    options,
  );

  if (id) {
    baseQuery = baseQuery.whereIn(MarketColumns.id, id);
  }

  if (pair) {
    baseQuery = baseQuery.whereIn(MarketColumns.pair, pair);
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
      MarketColumns.id,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  MarketToCreate: MarketCreateObject,
  options: Options = { txId: undefined },
): Promise<MarketFromDatabase> {
  return MarketModel.query(
    Transaction.get(options.txId),
  ).insert(MarketToCreate).returning('*');
}

export async function update(
  {
    id,
    ...fields
  }: MarketUpdateObject,
  options: Options = { txId: undefined },
): Promise<MarketFromDatabase | undefined> {
  const market = await MarketModel.query(
    Transaction.get(options.txId),
    // TODO fix expression typing so we dont have to use any
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ).findById(id);
  const updatedMarket = await market
    .$query(Transaction.get(options.txId))
    .patch(fields as PartialModelObject<MarketModel>).returning('*');
  // The objection types mistakenly think the query returns an array of markets.
  return updatedMarket as unknown as (MarketFromDatabase | undefined);
}

export async function findById(
  id: number,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<MarketFromDatabase | undefined> {
  const baseQuery: QueryBuilder<MarketModel> = setupBaseQuery<MarketModel>(
    MarketModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function findByPair(
  pair: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<MarketFromDatabase | undefined> {
  const baseQuery: QueryBuilder<MarketModel> = setupBaseQuery<MarketModel>(
    MarketModel,
    options,
  );

  const markets: MarketFromDatabase[] = await baseQuery
    .where(MarketColumns.pair, pair)
    .returning('*');

  if (markets.length === 0) {
    return undefined;
  }
  return markets[0];
}

export async function getMarketsMap(
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<MarketsMap> {
  const markets: MarketFromDatabase[] = await findAll(
    {},
    [],
    options,
  );
  return _.keyBy(markets, MarketColumns.id);
}
