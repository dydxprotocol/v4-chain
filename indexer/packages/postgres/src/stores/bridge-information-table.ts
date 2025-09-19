import { QueryBuilder } from 'objection';
import { v4 as uuidv4 } from 'uuid';

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
  PaginationFromDatabase,
} from '../types';

export async function findByFromAddressWithTransactionHashFilter(
  fromAddress: string,
  hasTransactionHash: boolean,
  queryOptions: BridgeInformationQueryOptions = {},
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<BridgeInformationFromDatabase>> {
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

  return handleLimitAndPagination(query, queryOptions.limit, queryOptions.page);
}

export async function searchBridgeInformation(
  filters: BridgeInformationQueryFilters = {},
  queryOptions: BridgeInformationQueryOptions = {},
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<BridgeInformationFromDatabase>> {
  const baseQuery: QueryBuilder<BridgeInformationModel> = setupBaseQuery<BridgeInformationModel>(
    BridgeInformationModel,
    options,
  );

  let query = baseQuery;

  // Apply filters
  if (filters.from_addresses) {
    query = query.whereIn(BridgeInformationColumns.from_address, filters.from_addresses);
  }
  if (filters.chain_id) {
    query = query.where(BridgeInformationColumns.chain_id, filters.chain_id);
  }
  if (filters.transaction_hash) {
    query = query.where(BridgeInformationColumns.transaction_hash, filters.transaction_hash);
  }
  if (filters.sinceDate) {
    query = query.where(BridgeInformationColumns.created_at, '>=', filters.sinceDate);
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

  return handleLimitAndPagination(query, queryOptions.limit, queryOptions.page);
}

export async function create(
  bridgeInformationToCreate: BridgeInformationCreateObject,
  options: Options = { txId: undefined },
): Promise<BridgeInformationFromDatabase> {
  // Generate UUID if id is not provided
  const createObject = { ...bridgeInformationToCreate };
  if (!createObject.id) {
    createObject.id = uuidv4();
  }

  return BridgeInformationModel.query(
    Transaction.get(options.txId),
  ).insert(createObject).returning('*');
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

/**
 * Handles pagination and limit logic for bridge information queries
 * @param limit Maximum number of bridge information to return
 * @param page Page number
 * @returns Promise<PaginationFromDatabase<BridgeInformationFromDatabase>>
 */
async function handleLimitAndPagination(
  baseQuery: QueryBuilder<BridgeInformationModel>,
  limit?: number,
  page?: number,
): Promise<PaginationFromDatabase<BridgeInformationFromDatabase>> {
  let query = baseQuery;

  /**
   * If a query is made using a page number, then the limit property is used as 'page limit'
   */
  if (page !== undefined && limit !== undefined) {
    /**
     * We make sure that the page number is always >= 1
     */
    const currentPage: number = Math.max(1, page);
    const offset: number = (currentPage - 1) * limit;

    /**
     * Ensure sorting is applied to maintain consistent pagination results.
     * Also a casting of the ts type is required since the infer of the type
     * obtained from the count is not performed.
     */
    const count: { count?: string } = (await query
      .clone()
      .clearOrder()
      .count({ count: '*' })
      .first()) as unknown as { count?: string };

    query = query.offset(offset).limit(limit);

    const results = (await query.returning('*')) as BridgeInformationFromDatabase[];
    return {
      results,
      limit,
      offset,
      total: parseInt(count.count ?? '0', 10),
    };
  }

  // If no pagination, just apply the limit
  if (limit !== undefined) {
    query = query.limit(limit);
  }

  const results = (await query.returning('*')) as BridgeInformationFromDatabase[];
  return {
    results,
  };
}
