import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import AffiliateReferredUsersModel from '../models/affiliate-referred-users-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  AffiliateReferredUsersColumns,
  AffiliateReferredUsersCreateObject,
  AffiliateReferredUserFromDatabase,
  AffiliateReferredUsersQueryConfig,
} from '../types';

export async function findAll(
  {
    affiliateAddress,
    refereeAddress,
    limit,
  }: AffiliateReferredUsersQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateReferredUserFromDatabase[]> {
  verifyAllRequiredFields(
    {
      affiliateAddress,
      refereeAddress,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  // splitting the line after = does not work because it is reformatted to one line by eslint
  // eslint-disable-next-line max-len
  let baseQuery: QueryBuilder<AffiliateReferredUsersModel> = setupBaseQuery<AffiliateReferredUsersModel>(
    AffiliateReferredUsersModel,
    options,
  );

  if (affiliateAddress) {
    baseQuery = baseQuery.where(AffiliateReferredUsersColumns.affiliateAddress, affiliateAddress);
  }

  if (refereeAddress) {
    baseQuery = baseQuery.where(AffiliateReferredUsersColumns.refereeAddress, refereeAddress);
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
      AffiliateReferredUsersColumns.referredAtBlock,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  entryToCreate: AffiliateReferredUsersCreateObject,
  options: Options = { txId: undefined },
): Promise<AffiliateReferredUserFromDatabase> {
  return AffiliateReferredUsersModel.query(
    Transaction.get(options.txId),
  ).insert(entryToCreate).returning('*');
}

export async function findByAffiliateAddress(
  address: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateReferredUserFromDatabase[] | undefined> {
  // splitting the line after = does not work because it is reformatted to one line by eslint
  // eslint-disable-next-line max-len
  const baseQuery: QueryBuilder<AffiliateReferredUsersModel> = setupBaseQuery<AffiliateReferredUsersModel>(
    AffiliateReferredUsersModel,
    options,
  );
  return baseQuery
    .where('affiliateAddress', address)
    .returning('*');
}

export async function findByRefereeAddress(
  address: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateReferredUserFromDatabase | undefined> {
  // splitting the line after = does not work because it is reformatted to one line by eslint
  // eslint-disable-next-line max-len
  const baseQuery: QueryBuilder<AffiliateReferredUsersModel> = setupBaseQuery<AffiliateReferredUsersModel>(
    AffiliateReferredUsersModel,
    options,
  );
  return baseQuery
    .where('refereeAddress', address)
    .returning('*')
    .first(); // should only be one since refereeAddress is primary key
}
