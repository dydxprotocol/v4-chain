import Knex from 'knex';
import { PartialModelObject, QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexPrimary } from '../helpers/knex';
import {
  generateBulkUpsertString,
  setBulkRowsForUpdate,
  setupBaseQuery,
  verifyAllInjectableVariables,
  verifyAllRequiredFields,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import ComplianceStatusModel from '../models/compliance-status-model';
import {
  ComplianceStatusColumns,
  ComplianceStatusCreateObject,
  ComplianceStatusFromDatabase,
  ComplianceStatusQueryConfig,
  ComplianceStatusUpdateObject,
  ComplianceStatusUpsertObject,
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
} from '../types';

export async function findAll(
  {
    address,
    status,
    reason,
    createdBeforeOrAt,
    updatedBeforeOrAt,
    limit,
  }: ComplianceStatusQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<ComplianceStatusFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      status,
      reason,
      createdBeforeOrAt,
      updatedBeforeOrAt,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<ComplianceStatusModel> = setupBaseQuery<ComplianceStatusModel>(
    ComplianceStatusModel,
    options,
  );

  if (address !== undefined) {
    baseQuery = baseQuery.whereIn(ComplianceStatusColumns.address, address);
  }

  if (status !== undefined) {
    baseQuery = baseQuery.whereIn(ComplianceStatusColumns.status, status);
  }

  if (reason !== undefined) {
    baseQuery = baseQuery.where(ComplianceStatusColumns.reason, reason);
  }

  if (updatedBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(ComplianceStatusColumns.updatedAt, '<=', updatedBeforeOrAt);
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(ComplianceStatusColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(column, order);
    }
  } else {
    baseQuery = baseQuery.orderBy(ComplianceStatusColumns.updatedAt, Ordering.ASC);
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  complianceStatusToCreate: ComplianceStatusCreateObject,
  options: Options = { txId: undefined },
): Promise<ComplianceStatusFromDatabase> {
  return ComplianceStatusModel.query(
    Transaction.get(options.txId),
  ).insert(complianceStatusToCreate).returning('*');
}

export async function update(
  {
    address,
    ...fields
  }: ComplianceStatusUpdateObject,
  options: Options = { txId: undefined },
): Promise<ComplianceStatusFromDatabase | undefined> {
  const complianceStatus = await ComplianceStatusModel.query(
    Transaction.get(options.txId),
  ).findById(address);
  const updatedComplianceStatus = await complianceStatus
    .$query()
    .patch(fields as PartialModelObject<ComplianceStatusModel>)
    .returning('*');
  // The objection types mistakenly think the query returns an array of ComplianceStatus.
  return updatedComplianceStatus as unknown as (ComplianceStatusFromDatabase | undefined);
}

export async function upsert(
  complianceStatusToUpsert: ComplianceStatusUpsertObject,
  options: Options = { txId: undefined },
): Promise<ComplianceStatusFromDatabase> {
  const complianceStatus: ComplianceStatusFromDatabase | undefined = await
  ComplianceStatusModel.query(
    Transaction.get(options.txId),
  ).findById(complianceStatusToUpsert.address);

  if (complianceStatus === undefined) {
    return create(complianceStatusToUpsert, options);
  }

  const updatedComplianceStatus: ComplianceStatusFromDatabase | undefined = await update({
    ...complianceStatusToUpsert,
  }, options);

  if (updatedComplianceStatus === undefined) {
    throw Error('Compliance status must exist after update');
  }

  return updatedComplianceStatus;
}

export async function bulkUpsert(
  complianceStatusObjects: ComplianceStatusUpsertObject[],
  options: Options = { txId: undefined },
): Promise<void> {
  if (complianceStatusObjects.length === 0) {
    return;
  }

  complianceStatusObjects.forEach(
    (complianceStatusObject: ComplianceStatusUpsertObject) => verifyAllInjectableVariables(
      Object.values(complianceStatusObject),
    ),
  );

  const columns: ComplianceStatusColumns[] = [
    ComplianceStatusColumns.address,
    ComplianceStatusColumns.status,
    ComplianceStatusColumns.reason,
    ComplianceStatusColumns.updatedAt,
  ];
  const rows: string[] = setBulkRowsForUpdate<ComplianceStatusColumns>({
    objectArray: complianceStatusObjects,
    columns,
    stringColumns: [
      ComplianceStatusColumns.address,
    ],
    timestampColumns: [
      ComplianceStatusColumns.updatedAt,
    ],
    enumColumns: [
      ComplianceStatusColumns.status,
      ComplianceStatusColumns.reason,
    ],
  });

  const query: string = generateBulkUpsertString({
    table: ComplianceStatusModel.tableName,
    objectRows: rows,
    columns,
    uniqueIdentifiers: [ComplianceStatusColumns.address],
  });

  const transaction: Knex.Transaction | undefined = Transaction.get(options.txId);
  return transaction
    ? knexPrimary.raw(query).transacting(transaction)
    : knexPrimary.raw(query);
}
