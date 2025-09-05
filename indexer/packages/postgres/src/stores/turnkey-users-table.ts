import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import TurnkeyUserModel from '../models/turnkey-user-model';
import {
  Options,
  TurnkeyUserColumns,
  TurnkeyUserCreateObject,
  TurnkeyUserFromDatabase,
} from '../types';

export async function findByEvmAddress(
  evmAddress: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TurnkeyUserFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TurnkeyUserModel> = setupBaseQuery<TurnkeyUserModel>(
    TurnkeyUserModel,
    options,
  );
  return baseQuery
    .where(TurnkeyUserColumns.evm_address, evmAddress)
    .first();
}

export async function findBySvmAddress(
  svmAddress: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TurnkeyUserFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TurnkeyUserModel> = setupBaseQuery<TurnkeyUserModel>(
    TurnkeyUserModel,
    options,
  );
  return baseQuery
    .where(TurnkeyUserColumns.svm_address, svmAddress)
    .first();
}

export async function findBySmartAccountAddress(
  smartAccountAddress: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TurnkeyUserFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TurnkeyUserModel> = setupBaseQuery<TurnkeyUserModel>(
    TurnkeyUserModel,
    options,
  );
  return baseQuery
    .where(TurnkeyUserColumns.smart_account_address, smartAccountAddress)
    .first();
}

export async function findByDydxAddress(
  dydxAddress: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TurnkeyUserFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TurnkeyUserModel> = setupBaseQuery<TurnkeyUserModel>(
    TurnkeyUserModel,
    options,
  );
  return baseQuery
    .where(TurnkeyUserColumns.dydx_address, dydxAddress)
    .first();
}

export async function findByEmail(
  email: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TurnkeyUserFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TurnkeyUserModel> = setupBaseQuery<TurnkeyUserModel>(
    TurnkeyUserModel,
    options,
  );
  return baseQuery
    .where(TurnkeyUserColumns.email, email)
    .first();
}

export async function findBySuborgId(
  suborgId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TurnkeyUserFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TurnkeyUserModel> = setupBaseQuery<TurnkeyUserModel>(
    TurnkeyUserModel,
    options,
  );
  return baseQuery
    .where(TurnkeyUserColumns.suborg_id, suborgId)
    .first();
}

export async function create(
  turnkeyUserToCreate: TurnkeyUserCreateObject,
  options: Options = { txId: undefined },
): Promise<TurnkeyUserFromDatabase> {
  return TurnkeyUserModel.query(
    Transaction.get(options.txId),
  ).insert(turnkeyUserToCreate).returning('*');
}

export async function upsert(
  turnkeyUserToUpsert: TurnkeyUserCreateObject,
  options: Options = { txId: undefined },
): Promise<TurnkeyUserFromDatabase> {
  const turnkeyUsers: TurnkeyUserModel[] = await TurnkeyUserModel.query(
    Transaction.get(options.txId),
  ).upsert(turnkeyUserToUpsert).returning('*');
  if (turnkeyUsers.length === 0) {
    throw new Error('Upsert failed to return records');
  }
  return turnkeyUsers[0];
}

export async function updateDydxAddressByEvmAddress(
  evmAddress: string,
  dydxAddress: string,
  options: Options = { txId: undefined },
): Promise<TurnkeyUserFromDatabase | undefined> {
  const existing: TurnkeyUserModel | undefined = await TurnkeyUserModel.query(
    Transaction.get(options.txId),
  )
    .where(TurnkeyUserColumns.evm_address, evmAddress)
    .first();

  if (!existing) {
    return undefined;
  }

  const updated = await existing.$query().patch({
    [TurnkeyUserColumns.dydx_address]: dydxAddress,
  }).returning('*');

  return updated as unknown as TurnkeyUserFromDatabase | undefined;
}
