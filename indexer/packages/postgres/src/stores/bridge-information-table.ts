import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import BridgeInformationModel from '../models/bridge-information-model';
import {
  BridgeInformationColumns,
  BridgeInformationCreateObject,
  BridgeInformationFromDatabase,
  BridgeInformationQueryFilters,
  BridgeInformationQueryOptions,
  Options,
} from '../types';

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<BridgeInformationFromDatabase | undefined> {
  const baseQuery: QueryBuilder<BridgeInformationModel> = setupBaseQuery<BridgeInformationModel>(
    BridgeInformationModel,
    options,
  );
  return baseQuery
    .where(BridgeInformationColumns.id, id)
    .first();
}

export async function findByFromAddress(
  fromAddress: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<BridgeInformationFromDatabase[]> {
  const baseQuery: QueryBuilder<BridgeInformationModel> = setupBaseQuery<BridgeInformationModel>(
    BridgeInformationModel,
    options,
  );
  return baseQuery
    .where(BridgeInformationColumns.from_address, fromAddress)
    .orderBy(BridgeInformationColumns.created_at, 'desc');
}

export async function findByFromAddressWithTransactionHashFilter(
  fromAddress: string,
  hasTransactionHash: boolean,
  queryOptions: BridgeInformationQueryOptions = {},
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<BridgeInformationFromDatabase[]> {
  const baseQuery: QueryBuilder<BridgeInformationModel> = setupBaseQuery<BridgeInformationModel>(
    BridgeInformationModel,
    options,
  );

  let query = baseQuery.where(BridgeInformationColumns.from_address, fromAddress);

  // Filter by transaction_hash null/not null
  if (hasTransactionHash) {
    query = query.whereNotNull(BridgeInformationColumns.transaction_hash);
  } else {
    query = query.whereNull(BridgeInformationColumns.transaction_hash);
  }

  // Apply ordering
  const orderBy = queryOptions.orderBy || 'created_at';
  const orderDirection = queryOptions.orderDirection || 'DESC';
  query = query.orderBy(BridgeInformationColumns[orderBy], orderDirection);

  // Apply pagination
  if (queryOptions.limit) {
    query = query.limit(queryOptions.limit);
  }
  if (queryOptions.offset) {
    query = query.offset(queryOptions.offset);
  }

  return query;
}

export async function searchBridgeInformation(
  filters: BridgeInformationQueryFilters = {},
  queryOptions: BridgeInformationQueryOptions = {},
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<BridgeInformationFromDatabase[]> {
  const baseQuery: QueryBuilder<BridgeInformationModel> = setupBaseQuery<BridgeInformationModel>(
    BridgeInformationModel,
    options,
  );

  let query = baseQuery;

  // Apply filters
  if (filters.from_address) {
    query = query.where(BridgeInformationColumns.from_address, filters.from_address);
  }
  if (filters.chain_id) {
    query = query.where(BridgeInformationColumns.chain_id, filters.chain_id);
  }
  if (filters.transaction_hash) {
    query = query.where(BridgeInformationColumns.transaction_hash, filters.transaction_hash);
  }
  if (filters.has_transaction_hash !== undefined) {
    if (filters.has_transaction_hash) {
      query = query.whereNotNull(BridgeInformationColumns.transaction_hash);
    } else {
      query = query.whereNull(BridgeInformationColumns.transaction_hash);
    }
  }

  // Apply ordering
  const orderBy = queryOptions.orderBy || 'created_at';
  const orderDirection = queryOptions.orderDirection || 'DESC';
  query = query.orderBy(BridgeInformationColumns[orderBy], orderDirection);

  // Apply pagination
  if (queryOptions.limit) {
    query = query.limit(queryOptions.limit);
  }
  if (queryOptions.offset) {
    query = query.offset(queryOptions.offset);
  }

  return query;
}

export async function create(
  bridgeInformationToCreate: BridgeInformationCreateObject,
  options: Options = { txId: undefined },
): Promise<BridgeInformationFromDatabase> {
  return BridgeInformationModel.query(
    Transaction.get(options.txId),
  ).insert(bridgeInformationToCreate).returning('*');
}

export async function upsert(
  bridgeInformationToUpsert: BridgeInformationCreateObject,
  options: Options = { txId: undefined },
): Promise<BridgeInformationFromDatabase> {
  const bridgeInformationRecords: BridgeInformationModel[] = await BridgeInformationModel.query(
    Transaction.get(options.txId),
  ).upsert(bridgeInformationToUpsert).returning('*');
  if (bridgeInformationRecords.length === 0) {
    throw new Error('Upsert failed to return records');
  }
  return bridgeInformationRecords[0];
}

export async function updateTransactionHash(
  id: string,
  transactionHash: string,
  options: Options = { txId: undefined },
): Promise<BridgeInformationFromDatabase | undefined> {
  const existing: BridgeInformationModel | undefined = await BridgeInformationModel.query(
    Transaction.get(options.txId),
  )
    .where(BridgeInformationColumns.id, id)
    .first();

  if (!existing) {
    return undefined;
  }

  const updated = await existing.$query().patch({
    [BridgeInformationColumns.transaction_hash]: transactionHash,
  }).returning('*');

  return updated as unknown as BridgeInformationFromDatabase | undefined;
}
