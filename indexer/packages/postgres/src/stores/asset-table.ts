import { PartialModelObject, QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import AssetModel from '../models/asset-model';
import {
  AssetColumns,
  AssetCreateObject,
  AssetFromDatabase,
  AssetQueryConfig,
  AssetUpdateObject,
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
} from '../types';

export async function findAll(
  {
    limit,
    id,
    symbol,
    atomicResolution,
    hasMarket,
    marketId,
  }: AssetQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AssetFromDatabase[]> {
  verifyAllRequiredFields(
    {
      id,
      symbol,
      atomicResolution,
      hasMarket,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<AssetModel> = setupBaseQuery<AssetModel>(
    AssetModel,
    options,
  );

  if (id) {
    baseQuery = baseQuery.whereIn(AssetColumns.id, id);
  }

  if (symbol) {
    baseQuery = baseQuery.where(AssetColumns.symbol, symbol);
  }

  if (atomicResolution) {
    baseQuery = baseQuery.where(AssetColumns.atomicResolution, atomicResolution);
  }

  if (hasMarket) {
    baseQuery = baseQuery.where(AssetColumns.hasMarket, hasMarket);
  }

  if (marketId) {
    baseQuery = baseQuery.where(AssetColumns.marketId, marketId);
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
      AssetColumns.symbol,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  assetToCreate: AssetCreateObject,
  options: Options = { txId: undefined },
): Promise<AssetFromDatabase> {
  return AssetModel.query(
    Transaction.get(options.txId),
  ).insert(assetToCreate).returning('*');
}

export async function update(
  {
    id,
    ...fields
  }: AssetUpdateObject,
  options: Options = { txId: undefined },
): Promise<AssetFromDatabase | undefined> {
  const asset = await AssetModel.query(
    Transaction.get(options.txId),
    // TODO fix expression typing so we dont have to use any
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ).findById(id);
  const updatedAsset = await asset.$query().patch(fields as PartialModelObject<AssetModel>).returning('*');
  // The objection types mistakenly think the query returns an array of assets.
  return updatedAsset as unknown as (AssetFromDatabase | undefined);
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AssetFromDatabase | undefined> {
  const baseQuery: QueryBuilder<AssetModel> = setupBaseQuery<AssetModel>(
    AssetModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}
