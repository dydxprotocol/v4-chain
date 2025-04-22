import { logger } from '@dydxprotocol-indexer/base';
import Knex from 'knex';
import _ from 'lodash';
import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexPrimary } from '../helpers/knex';
import {
  generateBulkUpdateString,
  setBulkRowsForUpdate,
  setupBaseQuery,
  verifyAllInjectableVariables,
  verifyAllRequiredFields,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import PerpetualMarketModel from '../models/perpetual-market-model';
import {
  Options,
  Ordering,
  PerpetualMarketColumns,
  PerpetualMarketCreateObject,
  PerpetualMarketFromDatabase,
  PerpetualMarketQueryConfig,
  PerpetualMarketUpdateObject,
  PerpetualMarketWithMarket,
  QueryableField,
  QueryConfig,
} from '../types';

export async function findAll(
  {
    id,
    marketId,
    liquidityTierId,
    limit,
    joinWithMarkets = false,
  }: PerpetualMarketQueryConfig & { joinWithMarkets?: boolean },
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PerpetualMarketFromDatabase[] | PerpetualMarketWithMarket[]> {
  verifyAllRequiredFields(
    {
      id,
      marketId,
      liquidityTierId,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<PerpetualMarketModel> = setupBaseQuery<PerpetualMarketModel>(
    PerpetualMarketModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(PerpetualMarketColumns.id, id);
  }

  if (marketId !== undefined) {
    baseQuery = baseQuery.whereIn(PerpetualMarketColumns.marketId, marketId);
  }

  if (liquidityTierId !== undefined) {
    baseQuery = baseQuery.whereIn(PerpetualMarketColumns.liquidityTierId, liquidityTierId);
  }

  if (joinWithMarkets) {
    baseQuery = baseQuery
      .joinRelated('market')
      .select([
        `${PerpetualMarketModel.tableName}.*`,
        'market.pair',
        'market.exponent',
        'market.minPriceChangePpm',
        'market.oraclePrice',
      ]);
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
      PerpetualMarketColumns.id,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  perpetualMarketToCreate: PerpetualMarketCreateObject,
  options: Options = { txId: undefined },
): Promise<PerpetualMarketFromDatabase> {
  return PerpetualMarketModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...perpetualMarketToCreate,
  }).returning('*');
}

export async function update(
  {
    ...fields
  }: PerpetualMarketUpdateObject,
  options: Options = { txId: undefined },
): Promise<PerpetualMarketFromDatabase | undefined> {
  const perpetualMarket = await PerpetualMarketModel.query(
    Transaction.get(options.txId),
  // TODO fix expression typing so we dont have to use any
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ).findById(fields.id!).patch(fields as any).returning('*');
  // The objection types mistakenly think the query returns an array of perpetualMarkets.
  return perpetualMarket as unknown as (PerpetualMarketFromDatabase | undefined);
}

export async function updateByMarketId(
  {
    ...fields
  }: PerpetualMarketUpdateObject,
  options: Options = { txId: undefined },
): Promise<PerpetualMarketFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PerpetualMarketModel> = setupBaseQuery<PerpetualMarketModel>(
    PerpetualMarketModel,
    options,
  );

  const perpetualMarkets: PerpetualMarketFromDatabase[] = await baseQuery
    .where(PerpetualMarketColumns.marketId, fields.marketId!)
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    .patch(fields as any)
    .returning('*');

  if (perpetualMarkets.length === 0) {
    return undefined;
  }
  return perpetualMarkets[0];
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PerpetualMarketFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PerpetualMarketModel> = setupBaseQuery<PerpetualMarketModel>(
    PerpetualMarketModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function findByClobPairId(
  clobPairId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PerpetualMarketFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PerpetualMarketModel> = setupBaseQuery<PerpetualMarketModel>(
    PerpetualMarketModel,
    options,
  );

  const perpetualMarkets: PerpetualMarketFromDatabase[] = await baseQuery
    .where(PerpetualMarketColumns.clobPairId, clobPairId)
    .returning('*');

  if (perpetualMarkets.length === 0) {
    return undefined;
  }
  return perpetualMarkets[0];
}

export async function findByMarketId(
  marketId: number,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PerpetualMarketFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PerpetualMarketModel> = setupBaseQuery<PerpetualMarketModel>(
    PerpetualMarketModel,
    options,
  );

  const perpetualMarkets: PerpetualMarketFromDatabase[] = await baseQuery
    .where(PerpetualMarketColumns.marketId, marketId)
    .returning('*');

  if (perpetualMarkets.length === 0) {
    return undefined;
  }
  return perpetualMarkets[0];
}

export async function findByTicker(
  ticker: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PerpetualMarketFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PerpetualMarketModel> = setupBaseQuery<PerpetualMarketModel>(
    PerpetualMarketModel,
    options,
  );

  const perpetualMarkets: PerpetualMarketFromDatabase[] = await baseQuery
    .where(PerpetualMarketColumns.ticker, ticker)
    .returning('*');

  if (perpetualMarkets.length === 0) {
    return undefined;
  }

  if (perpetualMarkets.length > 1) {
    logger.error({
      at: 'perpetualMarketTable#findByTicker',
      message: 'More than one market with ticker',
      ticker,
    });
  }

  return perpetualMarkets[0];
}

export async function updateMarketCheckerFields(
  markets: {
    id: string,
    volume24H: string,
    trades24H: number,
    priceChange24H: string,
    openInterest: string,
    nextFundingRate: string,
  }[],
  transaction?: Knex.Transaction,
): Promise<void> {
  if (markets.length === 0) {
    return;
  }

  markets.forEach((market) => verifyAllInjectableVariables(Object.values(market)));

  const columns = _.keys(markets[0]) as PerpetualMarketColumns[];
  const marketRows: string[] = setBulkRowsForUpdate<PerpetualMarketColumns>({
    objectArray: markets,
    columns,
    stringColumns: [
      PerpetualMarketColumns.status,
      PerpetualMarketColumns.ticker,
    ],
    numericColumns: [
      PerpetualMarketColumns.id,
      PerpetualMarketColumns.volume24H,
      PerpetualMarketColumns.trades24H,
      PerpetualMarketColumns.priceChange24H,
      PerpetualMarketColumns.openInterest,
      PerpetualMarketColumns.nextFundingRate,
    ],
  });

  const query: string = generateBulkUpdateString({
    table: 'perpetual_markets',
    objectRows: marketRows,
    columns,
    isUuid: false,
    uniqueIdentifier: PerpetualMarketColumns.id,
  });

  if (transaction) {
    await knexPrimary.raw(query).transacting(transaction);
  } else {
    await knexPrimary.raw(query);
  }
}
