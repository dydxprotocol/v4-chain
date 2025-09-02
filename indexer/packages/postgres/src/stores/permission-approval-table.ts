import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import PermissionApprovalModel from '../models/permission-approval-model';
import {
  Options,
  PermissionApprovalColumns,
  PermissionApprovalCreateObject,
  PermissionApprovalFromDatabase,
} from '../types';

export async function findBySuborgId(
  suborgId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PermissionApprovalFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PermissionApprovalModel> = setupBaseQuery<PermissionApprovalModel>(
    PermissionApprovalModel,
    options,
  );
  return baseQuery
    .where(PermissionApprovalColumns.suborg_id, suborgId)
    .first();
}

export async function getArbitrumApprovalForSuborg(
  suborgId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string | undefined> {
  const baseQuery: QueryBuilder<PermissionApprovalModel> = setupBaseQuery<PermissionApprovalModel>(
    PermissionApprovalModel,
    options,
  );
  const result = await baseQuery
    .select(PermissionApprovalColumns.arbitrum_approval)
    .where(PermissionApprovalColumns.suborg_id, suborgId)
    .first();

  return result?.arbitrum_approval;
}

export async function getBaseApprovalForSuborg(
  suborgId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string | undefined> {
  const baseQuery: QueryBuilder<PermissionApprovalModel> = setupBaseQuery<PermissionApprovalModel>(
    PermissionApprovalModel,
    options,
  );
  const result = await baseQuery
    .select(PermissionApprovalColumns.base_approval)
    .where(PermissionApprovalColumns.suborg_id, suborgId)
    .first();

  return result?.base_approval;
}

export async function getAvalancheApprovalForSuborg(
  suborgId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string | undefined> {
  const baseQuery: QueryBuilder<PermissionApprovalModel> = setupBaseQuery<PermissionApprovalModel>(
    PermissionApprovalModel,
    options,
  );
  const result = await baseQuery
    .select(PermissionApprovalColumns.avalanche_approval)
    .where(PermissionApprovalColumns.suborg_id, suborgId)
    .first();

  return result?.avalanche_approval;
}

export async function getOptimismApprovalForSuborg(
  suborgId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string | undefined> {
  const baseQuery: QueryBuilder<PermissionApprovalModel> = setupBaseQuery<PermissionApprovalModel>(
    PermissionApprovalModel,
    options,
  );
  const result = await baseQuery
    .select(PermissionApprovalColumns.optimism_approval)
    .where(PermissionApprovalColumns.suborg_id, suborgId)
    .first();

  return result?.optimism_approval;
}

export async function getEthereumApprovalForSuborg(
  suborgId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string | undefined> {
  const baseQuery: QueryBuilder<PermissionApprovalModel> = setupBaseQuery<PermissionApprovalModel>(
    PermissionApprovalModel,
    options,
  );
  const result = await baseQuery
    .select(PermissionApprovalColumns.ethereum_approval)
    .where(PermissionApprovalColumns.suborg_id, suborgId)
    .first();

  return result?.ethereum_approval;
}

export async function create(
  permissionApprovalToCreate: PermissionApprovalCreateObject,
  options: Options = { txId: undefined },
): Promise<PermissionApprovalFromDatabase> {
  return PermissionApprovalModel.query(
    Transaction.get(options.txId),
  ).insert(permissionApprovalToCreate).returning('*');
}

export async function update(
  {
    // eslint-disable-next-line @typescript-eslint/naming-convention
    suborg_id,
    ...fields
  }: PermissionApprovalCreateObject,
  options: Options = { txId: undefined },
): Promise<PermissionApprovalFromDatabase | undefined> {
  const permissionApproval = await PermissionApprovalModel.query(
    Transaction.get(options.txId),
  ).findById(suborg_id);
  if (!permissionApproval) {
    return undefined;
  }

  const updatedPermissionApproval = await permissionApproval
    .$query()
    .patch(fields)
    .returning('*');
  // The objection types mistakenly think the query returns an array of PermissionApproval.
  return updatedPermissionApproval as unknown as (PermissionApprovalFromDatabase | undefined);
}
