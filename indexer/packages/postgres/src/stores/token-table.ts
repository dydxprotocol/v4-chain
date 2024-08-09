import { DateTime } from 'luxon';
import { PartialModelObject, QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import TokenModel from '../models/token-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  TokenColumns,
  TokenCreateObject,
  TokenFromDatabase,
  TokenQueryConfig,
  TokenUpdateObject,
} from '../types';

export async function findAll(
  {
    address,
    limit,
  }: TokenQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TokenFromDatabase[]> {
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
    baseQuery = baseQuery.where(TokenColumns.address, address);
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
      TokenColumns.updatedAt,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  tokenToCreate: TokenCreateObject,
  options: Options = { txId: undefined },
): Promise<TokenFromDatabase> {
  return TokenModel.query(
    Transaction.get(options.txId),
  ).insert(tokenToCreate).returning('*');
}

export async function update(
  {
    token,
    ...fields
  }: TokenUpdateObject,
  options: Options = { txId: undefined },
): Promise<TokenFromDatabase> {
  const existingToken = await TokenModel.query(
    Transaction.get(options.txId),
  ).findOne({ token });
  const updatedToken = await existingToken.$query().patch(fields as PartialModelObject<TokenModel>).returning('*');
  return updatedToken as unknown as TokenFromDatabase;
}

export async function upsert(
  tokenToUpsert: TokenCreateObject,
  options: Options = { txId: undefined },
): Promise<TokenFromDatabase> {
  const existingToken = await TokenModel.query(
    Transaction.get(options.txId),
  ).findOne({ token: tokenToUpsert.token });

  if (existingToken) {
    return update(tokenToUpsert, options);
  } else {
    return create(tokenToUpsert, options);
  }
}

export async function findByToken(
  token: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TokenFromDatabase | undefined> {
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
  options: Options = { txId: undefined },
): Promise<TokenFromDatabase> {
  return upsert(
    {
      token,
      address,
      updatedAt: DateTime.now().toISO(),
    },
    options,
  );
}
