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
): Promise<PermissionApprovalFromDatabase[]> {
  const baseQuery: QueryBuilder<PermissionApprovalModel> = setupBaseQuery<PermissionApprovalModel>(
    PermissionApprovalModel,
    options,
  );
  return baseQuery
    .where(PermissionApprovalColumns.suborg_id, suborgId);
}

export async function findBySuborgIdAndChainId(
  suborgId: string,
  chainId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PermissionApprovalFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PermissionApprovalModel> = setupBaseQuery<PermissionApprovalModel>(
    PermissionApprovalModel,
    options,
  );
  return baseQuery
    .where(PermissionApprovalColumns.suborg_id, suborgId)
    .where(PermissionApprovalColumns.chain_id, chainId)
    .first();
}

export async function getApprovalForSuborgAndChain(
  suborgId: string,
  chainId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string | undefined> {
  const approval = await findBySuborgIdAndChainId(suborgId, chainId, options);
  return approval?.approval;
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
    // eslint-disable-next-line @typescript-eslint/naming-convention
    chain_id,
    ...fields
  }: PermissionApprovalCreateObject,
  options: Options = { txId: undefined },
): Promise<PermissionApprovalFromDatabase | undefined> {
  const permissionApproval = await PermissionApprovalModel.query(
    Transaction.get(options.txId),
  ).findById([suborg_id, chain_id]);
  if (!permissionApproval) {
    return undefined;
  }

  const updatedPermissionApproval = await permissionApproval
    .$query()
    .patch(fields)
    .returning('*').first();
  return updatedPermissionApproval;
}

export async function upsert(
  permissionApprovalToUpsert: PermissionApprovalCreateObject,
  options: Options = { txId: undefined },
): Promise<PermissionApprovalFromDatabase> {
  const result = await PermissionApprovalModel.query(
    Transaction.get(options.txId),
  ).upsert(permissionApprovalToUpsert).returning('*').first();
  return result;
}
