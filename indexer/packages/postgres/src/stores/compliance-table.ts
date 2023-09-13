import { PartialModelObject, QueryBuilder } from 'objection';

import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import ComplianceDataModel from '../models/compliance-data-model';
import {
  ComplianceDataFromDatabase,
  ComplianceDataQueryConfig,
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
} from '../types';
import { ComplianceDataColumns, ComplianceDataCreateObject, ComplianceDataUpdateObject } from '../types/compliance-data-types';

export async function findAll(
  {
    updatedBeforeOrAt,
    provider,
    limit,
  }: ComplianceDataQueryConfig,
  requiredFields: QueryableField[],
  options: Options = {},
): Promise<ComplianceDataFromDatabase[]> {
  verifyAllRequiredFields(
    {
      updatedBeforeOrAt,
      provider,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<ComplianceDataModel> = setupBaseQuery<ComplianceDataModel>(
    ComplianceDataModel,
    options,
  );

  if (updatedBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(ComplianceDataColumns.updatedAt, '<=', updatedBeforeOrAt);
  }

  if (provider !== undefined) {
    baseQuery = baseQuery.where(ComplianceDataColumns.provider, provider);
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

export async function findByAddressAndProvider(
  address: string,
  provider: string,
  options: Options = {},
): Promise<ComplianceDataFromDatabase | undefined> {
  const baseQuery: QueryBuilder<ComplianceDataModel> = setupBaseQuery<ComplianceDataModel>(
    ComplianceDataModel,
    options,
  );
  return baseQuery
    .findById([address, provider])
    .returning('*');
}
