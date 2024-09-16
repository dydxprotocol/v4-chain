import { DateTime } from 'luxon';
import { PartialModelObject, QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import TokenModel from '../models/firebase-notification-token-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  FirebaseNotificationTokenColumns,
  FirebaseNotificationTokenCreateObject,
  FirebaseNotificationTokenFromDatabase,
  FirebaseNotificationTokenQueryConfig,
  FirebaseNotificationTokenUpdateObject,
} from '../types';

export async function findAll(
  {
    address,
    limit,
    updatedBeforeOrAt,
  }: FirebaseNotificationTokenQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<FirebaseNotificationTokenFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<TokenModel> = setupBaseQuery<TokenModel>(
    TokenModel,
    options,
  );

  if (address) {
    baseQuery = baseQuery.where(FirebaseNotificationTokenColumns.address, address);
  }

  if (updatedBeforeOrAt) {
    baseQuery = baseQuery.where(FirebaseNotificationTokenColumns.updatedAt, '<=', updatedBeforeOrAt);
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
      FirebaseNotificationTokenColumns.updatedAt,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  tokenToCreate: FirebaseNotificationTokenCreateObject,
  options: Options = { txId: undefined },
): Promise<FirebaseNotificationTokenFromDatabase> {
  return TokenModel.query(
    Transaction.get(options.txId),
  ).insert(tokenToCreate).returning('*');
}

export async function update(
  {
    token,
    ...fields
  }: FirebaseNotificationTokenUpdateObject,
  options: Options = { txId: undefined },
): Promise<FirebaseNotificationTokenFromDatabase> {
  const existingToken = await TokenModel.query(
    Transaction.get(options.txId),
  ).findOne({ token });
  const updatedToken = await existingToken.$query().patch(fields as PartialModelObject<TokenModel>).returning('*');
  return updatedToken as unknown as FirebaseNotificationTokenFromDatabase;
}

export async function upsert(
  tokenToUpsert: FirebaseNotificationTokenCreateObject,
  options: Options = { txId: undefined },
): Promise<FirebaseNotificationTokenFromDatabase> {
  const existingToken = await TokenModel.query(
    Transaction.get(options.txId),
  ).findOne({ token: tokenToUpsert.token });

  if (existingToken) {
    return update(tokenToUpsert, options);
  } else {
    return create(tokenToUpsert, options);
  }
}

export async function deleteMany(
  tokens: string[],
  options: Options = { txId: undefined },
): Promise<number> {
  const baseQuery: QueryBuilder<TokenModel> = setupBaseQuery<TokenModel>(
    TokenModel,
    options,
  );

  const result = await baseQuery
    .delete()
    .whereIn('token', tokens);
  return result;
}

export async function findByToken(
  token: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<FirebaseNotificationTokenFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TokenModel> = setupBaseQuery<TokenModel>(
    TokenModel,
    options,
  );
  return baseQuery
    .findOne({ token })
    .returning('*');
}

export async function registerToken(
  token: string,
  address: string,
  language: string,
  options: Options = { txId: undefined },
): Promise<FirebaseNotificationTokenFromDatabase> {
  return upsert(
    {
      token,
      address,
      updatedAt: DateTime.now().toISO(),
      language,
    },
    options,
  );
}
