import { QueryBuilder, raw } from 'objection';
import { v4 as uuidv4 } from 'uuid';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import GasSponsorshipModel from '../models/gas-sponsorship-model';
import {
  GasSponsorshipColumns,
  GasSponsorshipCreateObject,
  GasSponsorshipFromDatabase,
  GasSponsorshipQueryFilters,
  GasSponsorshipStatus,
  GasSponsorshipUpdateObject,
  Options,
} from '../types';

function applyFilters(
  baseQuery: QueryBuilder<GasSponsorshipModel>,
  filters: GasSponsorshipQueryFilters,
): QueryBuilder<GasSponsorshipModel> {
  let query: QueryBuilder<GasSponsorshipModel> = baseQuery;
  if (filters.chain_id !== undefined) {
    query = query.where(GasSponsorshipColumns.chain_id, filters.chain_id);
  }
  if (filters.grantee_address !== undefined) {
    query = query.where(GasSponsorshipColumns.grantee_address, filters.grantee_address);
  }
  if (filters.dydx_address !== undefined) {
    query = query.where(GasSponsorshipColumns.dydx_address, filters.dydx_address);
  }
  if (filters.statuses !== undefined) {
    query = query.whereIn(GasSponsorshipColumns.status, filters.statuses);
  }
  if (filters.createdOnOrAfter !== undefined) {
    query = query.where(GasSponsorshipColumns.created_at, '>=', filters.createdOnOrAfter);
  }
  return query;
}

/**
 * Finds sponsorships matching the filters, newest first.
 */
export async function findAll(
  filters: GasSponsorshipQueryFilters = {},
  limit?: number,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<GasSponsorshipFromDatabase[]> {
  let query: QueryBuilder<GasSponsorshipModel> = applyFilters(
    setupBaseQuery<GasSponsorshipModel>(GasSponsorshipModel, options),
    filters,
  ).orderBy(GasSponsorshipColumns.created_at, 'DESC');

  if (limit !== undefined) {
    query = query.limit(limit);
  }
  return (await query.returning('*')) as GasSponsorshipFromDatabase[];
}

/**
 * Finds the most recently created sponsorship matching the filters.
 */
export async function findLatest(
  filters: GasSponsorshipQueryFilters = {},
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<GasSponsorshipFromDatabase | undefined> {
  const results: GasSponsorshipFromDatabase[] = await findAll(filters, 1, options);
  return results[0];
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<GasSponsorshipFromDatabase | undefined> {
  const baseQuery: QueryBuilder<GasSponsorshipModel> = setupBaseQuery<GasSponsorshipModel>(
    GasSponsorshipModel,
    options,
  );
  return baseQuery.findById(id).returning('*') as unknown as Promise<
    GasSponsorshipFromDatabase | undefined
  >;
}

/**
 * Returns how much has been committed to sponsorships created on or after `createdOnOrAfter`, as
 * an integer string in the smallest unit of the chain's fee denom.
 *
 * A sponsorship commits its full spend limit plus the fee paid to grant it. Failed sponsorships
 * gave no allowance and are not counted. Unused allowance is still counted, so this is an upper
 * bound on what was actually spent.
 */
export async function getTotalCommittedSince(
  chainId: string,
  createdOnOrAfter: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string> {
  const baseQuery: QueryBuilder<GasSponsorshipModel> = setupBaseQuery<GasSponsorshipModel>(
    GasSponsorshipModel,
    options,
  );
  const result: { total?: string } | undefined = (await baseQuery
    .where(GasSponsorshipColumns.chain_id, chainId)
    .where(GasSponsorshipColumns.created_at, '>=', createdOnOrAfter)
    .whereNot(GasSponsorshipColumns.status, GasSponsorshipStatus.FAILED)
    .select(raw('COALESCE(SUM("spend_limit" + "fee_amount"), 0)::text AS "total"'))
    .first()) as unknown as { total?: string } | undefined;

  return result?.total ?? '0';
}

export async function create(
  gasSponsorshipToCreate: GasSponsorshipCreateObject,
  options: Options = { txId: undefined },
): Promise<GasSponsorshipFromDatabase> {
  const createObject: GasSponsorshipCreateObject = { ...gasSponsorshipToCreate };
  if (!createObject.id) {
    createObject.id = uuidv4();
  }

  return GasSponsorshipModel.query(
    Transaction.get(options.txId),
  ).insert(createObject).returning('*');
}

export async function update(
  id: string,
  fields: GasSponsorshipUpdateObject,
  options: Options = { txId: undefined },
): Promise<GasSponsorshipFromDatabase | undefined> {
  const existing: GasSponsorshipModel | undefined = await GasSponsorshipModel.query(
    Transaction.get(options.txId),
  ).findById(id);

  if (!existing) {
    return undefined;
  }

  const updated = await existing.$query(
    Transaction.get(options.txId),
  ).patch(fields).returning('*');
  // The objection types don't model returning('*') on an instance patch, which gives the row.
  return updated as unknown as GasSponsorshipFromDatabase | undefined;
}
