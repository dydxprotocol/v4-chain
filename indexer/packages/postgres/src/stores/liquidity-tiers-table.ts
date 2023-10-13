import { PartialModelObject, QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import LiquidityTiersModel from '../models/liquidity-tiers-model';
import {
  LiquidityTiersColumns,
  LiquidityTiersCreateObject,
  LiquidityTiersFromDatabase,
  LiquidityTiersQueryConfig,
  LiquidityTiersUpdateObject,
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
} from '../types';

export async function findAll(
  {
    limit,
    id,
  }: LiquidityTiersQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<LiquidityTiersFromDatabase[]> {
  verifyAllRequiredFields(
    {
      limit,
      id,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<LiquidityTiersModel> = setupBaseQuery<LiquidityTiersModel>(
    LiquidityTiersModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(LiquidityTiersColumns.id, id);
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
      LiquidityTiersColumns.id,
      Ordering.ASC,
    );
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  liquidityTierToCreate: LiquidityTiersCreateObject,
  options: Options = { txId: undefined },
): Promise<LiquidityTiersFromDatabase> {
  return LiquidityTiersModel.query(
    Transaction.get(options.txId),
  ).insert(liquidityTierToCreate).returning('*');
}

export async function findById(
  id: number,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<LiquidityTiersFromDatabase | undefined> {
  const baseQuery: QueryBuilder<LiquidityTiersModel> = setupBaseQuery<LiquidityTiersModel>(
    LiquidityTiersModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function update(
  {
    id,
    ...fields
  }: LiquidityTiersUpdateObject,
  options: Options = { txId: undefined },
): Promise<LiquidityTiersFromDatabase | undefined> {
  const liquidityTier = await LiquidityTiersModel.query(
    Transaction.get(options.txId),
    // TODO fix expression typing so we dont have to use any
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ).findById(id);
  const updatedLiquidityTiers = await liquidityTier.$query().patch(fields as PartialModelObject<LiquidityTiersModel>).returning('*');
  // The objection types mistakenly think the query returns an array of liquidityTiers.
  return updatedLiquidityTiers as unknown as (LiquidityTiersFromDatabase | undefined);
}

export async function upsert(
  tierToUpsert: LiquidityTiersCreateObject,
  options: Options = { txId: undefined },
): Promise<LiquidityTiersFromDatabase> {
  const tiers: LiquidityTiersModel[] = await LiquidityTiersModel.query(
    Transaction.get(options.txId),
  ).upsert(tierToUpsert).returning('*');
  // should only ever be one tier
  return tiers[0];
}
