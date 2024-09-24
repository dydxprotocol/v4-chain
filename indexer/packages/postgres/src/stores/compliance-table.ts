import Knex from 'knex';
import _ from 'lodash';
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
import ComplianceDataModel from '../models/compliance-data-model';
import WalletModel from '../models/wallet-model';
import {
  ComplianceDataFromDatabase,
  ComplianceDataQueryConfig,
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
} from '../types';
import {
  ComplianceDataColumns,
  ComplianceDataCreateObject,
  ComplianceDataUpdateObject,
} from '../types/compliance-data-types';

export async function findAll(
  {
    address,
    updatedBeforeOrAt,
    provider,
    blocked,
    limit,
    addressInWalletsTable,
  }: ComplianceDataQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<ComplianceDataFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      updatedBeforeOrAt,
      provider,
      blocked,
      limit,
      addressInWalletsTable,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<ComplianceDataModel> = setupBaseQuery<ComplianceDataModel>(
    ComplianceDataModel,
    options,
  );

  if (address !== undefined) {
    baseQuery = baseQuery.whereIn(ComplianceDataColumns.address, address);
  }

  if (updatedBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(ComplianceDataColumns.updatedAt, '<=', updatedBeforeOrAt);
  }

  if (provider !== undefined) {
    baseQuery = baseQuery.where(ComplianceDataColumns.provider, provider);
  }

  if (blocked !== undefined) {
    baseQuery = baseQuery.where(ComplianceDataColumns.blocked, blocked);
  }

  if (addressInWalletsTable === true) {
    baseQuery = baseQuery.innerJoin(
      WalletModel.tableName,
      `${ComplianceDataModel.tableName}.${ComplianceDataColumns.address}`,
      '=',
      `${WalletModel.tableName}.${WalletModel.idColumn}`);
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
      ComplianceDataColumns.updatedAt,
      Ordering.ASC,
    );
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  complianceDataToCreate: ComplianceDataCreateObject,
  options: Options = { txId: undefined },
): Promise<ComplianceDataFromDatabase> {
  return ComplianceDataModel.query(
    Transaction.get(options.txId),
  ).insert(complianceDataToCreate).returning('*');
}

export async function update(
  {
    address,
    provider,
    ...fields
  }: ComplianceDataUpdateObject,
  options: Options = { txId: undefined },
): Promise<ComplianceDataFromDatabase | undefined> {
  const complianceData = await ComplianceDataModel.query(
    Transaction.get(options.txId),
  ).findById([address, provider]);
  const updatedComplianceData = await complianceData
    .$query()
    .patch(fields as PartialModelObject<ComplianceDataModel>)
    .returning('*');
  // The objection types mistakenly think the query returns an array of ComplianceData.
  return updatedComplianceData as unknown as (ComplianceDataFromDatabase | undefined);
}

export async function upsert(
  complianceDataToUpsert: ComplianceDataCreateObject,
  options: Options = { txId: undefined },
): Promise<ComplianceDataFromDatabase> {
  const updatedComplianceData: ComplianceDataModel[] = await ComplianceDataModel.query(
    Transaction.get(options.txId),
  ).upsert(complianceDataToUpsert).returning('*');

  return updatedComplianceData[0];
}

export async function findByAddressAndProvider(
  address: string,
  provider: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<ComplianceDataFromDatabase | undefined> {
  const baseQuery: QueryBuilder<ComplianceDataModel> = setupBaseQuery<ComplianceDataModel>(
    ComplianceDataModel,
    options,
  );
  return baseQuery
    .findById([address, provider])
    .returning('*');
}

export async function bulkUpsert(
  complianceObjects: ComplianceDataCreateObject[],
  options: Options = { txId: undefined },
): Promise<void> {
  if (complianceObjects.length === 0) {
    return;
  }

  complianceObjects.forEach(
    (complianceObject: ComplianceDataCreateObject) => verifyAllInjectableVariables(
      Object.values(complianceObject),
    ),
  );

  const columns: ComplianceDataColumns[] = _.keys(complianceObjects[0]) as ComplianceDataColumns[];
  const rows: string[] = setBulkRowsForUpdate<ComplianceDataColumns>({
    objectArray: complianceObjects,
    columns,
    booleanColumns: [
      ComplianceDataColumns.blocked,
    ],
    numericColumns: [
      ComplianceDataColumns.riskScore,
    ],
    stringColumns: [
      ComplianceDataColumns.address,
      ComplianceDataColumns.chain,
      ComplianceDataColumns.provider,
    ],
    timestampColumns: [
      ComplianceDataColumns.updatedAt,
    ],
  });

  const query: string = generateBulkUpsertString({
    table: ComplianceDataModel.tableName,
    objectRows: rows,
    columns,
    uniqueIdentifiers: [ComplianceDataColumns.address, ComplianceDataColumns.provider],
  });

  const transaction: Knex.Transaction | undefined = Transaction.get(options.txId);
  return transaction
    ? knexPrimary.raw(query).transacting(transaction)
    : knexPrimary.raw(query);
}
