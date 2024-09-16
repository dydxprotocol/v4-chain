import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import AffiliateInfoModel from '../models/affiliate-info-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  AffiliateInfoColumns,
  AffiliateInfoCreateObject,
  AffiliateInfoFromDatabase,
  AffiliateInfoQueryConfig,
} from '../types';

export async function findAll(
  {
    address,
    limit,
  }: AffiliateInfoQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateInfoFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<AffiliateInfoModel> = setupBaseQuery<AffiliateInfoModel>(
    AffiliateInfoModel,
    options,
  );

  if (address) {
    baseQuery = baseQuery.where(AffiliateInfoColumns.address, address);
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
      AffiliateInfoColumns.address,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  AffiliateInfoToCreate: AffiliateInfoCreateObject,
  options: Options = { txId: undefined },
): Promise<AffiliateInfoFromDatabase> {
  return AffiliateInfoModel.query(
    Transaction.get(options.txId),
  ).insert(AffiliateInfoToCreate).returning('*');
}

export async function upsert(
  AffiliateInfoToUpsert: AffiliateInfoCreateObject,
  options: Options = { txId: undefined },
): Promise<AffiliateInfoFromDatabase> {
  const AffiliateInfos: AffiliateInfoModel[] = await AffiliateInfoModel.query(
    Transaction.get(options.txId),
  ).upsert(AffiliateInfoToUpsert).returning('*');
  // should only ever be one AffiliateInfo
  return AffiliateInfos[0];
}

export async function findById(
  address: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateInfoFromDatabase | undefined> {
  const baseQuery: QueryBuilder<AffiliateInfoModel> = setupBaseQuery<AffiliateInfoModel>(
    AffiliateInfoModel,
    options,
  );
  return baseQuery
    .findById(address)
    .returning('*');
}

export async function paginatedFindWithAddressFilter(
  addressFilter: string[],
  offset: number,
  limit: number,
  sortByAffiliateEarning: boolean,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AffiliateInfoFromDatabase[] | undefined> {
  let baseQuery: QueryBuilder<AffiliateInfoModel> = setupBaseQuery<AffiliateInfoModel>(
    AffiliateInfoModel,
    options,
  );

  // Apply address filter if provided
  if (addressFilter.length > 0) {
    baseQuery = baseQuery.whereIn(AffiliateInfoColumns.address, addressFilter);
  }

  // Sorting by affiliate earnings or default sorting by address
  if (sortByAffiliateEarning) {
    baseQuery = baseQuery.orderBy(AffiliateInfoColumns.affiliateEarnings, Ordering.DESC);
  }

  // Apply pagination using offset and limit
  baseQuery = baseQuery.offset(offset).limit(limit);

  // Returning all fields
  return baseQuery.returning('*');
}
