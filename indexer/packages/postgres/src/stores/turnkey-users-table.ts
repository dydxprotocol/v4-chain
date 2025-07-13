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
    .where(TurnkeyUserColumns.evmAddress, evmAddress)
    .first()
    .returning('*');
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
    .where(TurnkeyUserColumns.svmAddress, svmAddress)
    .first()
    .returning('*');
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
    .first()
    .returning('*');
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
    .where(TurnkeyUserColumns.suborgId, suborgId)
    .first()
    .returning('*');
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
  // should only ever be one turnkey user
  return turnkeyUsers[0];
}
