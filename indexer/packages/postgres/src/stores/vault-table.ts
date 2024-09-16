import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  verifyAllRequiredFields,
  setupBaseQuery,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import VaultModel from '../models/vault-model';
import {
  QueryConfig,
  VaultQueryConfig,
  VaultColumns,
  Options,
  Ordering,
  QueryableField,
  VaultFromDatabase,
  VaultCreateObject,
} from '../types';

export async function findAll(
  {
    address,
    clobPairId,
    status,
    limit,
  }: VaultQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<VaultFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      clobPairId,
      status,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<VaultModel> = setupBaseQuery<VaultModel>(
    VaultModel,
    options,
  );

  if (address) {
    baseQuery = baseQuery.whereIn(VaultColumns.address, address);
  }

  if (clobPairId) {
    baseQuery = baseQuery.whereIn(VaultColumns.clobPairId, clobPairId);
  }

  if (status) {
    baseQuery = baseQuery.whereIn(VaultColumns.status, status);
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
      VaultColumns.clobPairId,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  vaultToCreate: VaultCreateObject,
  options: Options = { txId: undefined },
): Promise<VaultFromDatabase> {
  return VaultModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...vaultToCreate,
  });
}

export async function upsert(
  vaultToUpsert: VaultCreateObject,
  options: Options = { txId: undefined },
): Promise<VaultFromDatabase> {
  const vaults: VaultModel[] = await VaultModel.query(
    Transaction.get(options.txId),
  ).upsert({
    ...vaultToUpsert,
  }).returning('*');
  return vaults[0];
}
