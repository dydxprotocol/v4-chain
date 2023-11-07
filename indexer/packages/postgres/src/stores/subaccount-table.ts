import { IndexerSubaccountId } from '@dydxprotocol-indexer/v4-protos';
import { PartialModelObject, QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  verifyAllRequiredFields,
  setupBaseQuery,
  rawQuery,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import SubaccountModel from '../models/subaccount-model';
import {
  QueryConfig,
  SubaccountFromDatabase,
  SubaccountQueryConfig,
  SubaccountColumns,
  SubaccountCreateObject,
  Options,
  Ordering,
  QueryableField,
  SubaccountUpdateObject,
} from '../types';

export function uuid(address: string, subaccountNumber: number): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(Buffer.from(`${address}-${subaccountNumber}`, BUFFER_ENCODING_UTF_8));
}

export function subaccountIdToUuid(subaccountId: IndexerSubaccountId): string {
  return uuid(subaccountId.owner, subaccountId.number);
}

export async function findAll(
  {
    id,
    address,
    subaccountNumber,
    updatedBeforeOrAt,
    updatedOnOrAfter,
    limit,
  }: SubaccountQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<SubaccountFromDatabase[]> {
  verifyAllRequiredFields(
    {
      id,
      address,
      subaccountNumber,
      updatedBeforeOrAt,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<SubaccountModel> = setupBaseQuery<SubaccountModel>(
    SubaccountModel,
    options,
  );

  if (id) {
    baseQuery = baseQuery.whereIn(SubaccountColumns.id, id);
  }

  if (address) {
    baseQuery = baseQuery.where(SubaccountColumns.address, address);
  }

  if (subaccountNumber) {
    baseQuery = baseQuery.where(SubaccountColumns.subaccountNumber, subaccountNumber);
  }

  if (updatedBeforeOrAt) {
    baseQuery = baseQuery.where(SubaccountColumns.updatedAt, '<=', updatedBeforeOrAt);
  }

  if (updatedOnOrAfter) {
    baseQuery = baseQuery.where(SubaccountColumns.updatedAt, '>=', updatedOnOrAfter);
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
      SubaccountColumns.address,
      Ordering.ASC,
    ).orderBy(
      SubaccountColumns.subaccountNumber,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function getSubaccountsWithTransfers(
  createdBeforeOrAtHeight: string,
  options: Options = {},
): Promise<SubaccountFromDatabase[]> {
  const queryString: string = `
    SELECT *
    FROM subaccounts
    WHERE id IN (
      SELECT "senderSubaccountId" FROM transfers
      WHERE "createdAtHeight" <= '${createdBeforeOrAtHeight}'
      UNION
      SELECT "recipientSubaccountId" FROM transfers
      WHERE "createdAtHeight" <= '${createdBeforeOrAtHeight}'
    )
  `;

  const result: {
    rows: SubaccountFromDatabase[],
  } = await rawQuery(queryString, options);

  return result.rows;
}

export async function create(
  subaccountToCreate: SubaccountCreateObject,
  options: Options = { txId: undefined },
): Promise<SubaccountFromDatabase> {
  return SubaccountModel.query(
    Transaction.get(options.txId),
  ).insert({
    id: uuid(subaccountToCreate.address, subaccountToCreate.subaccountNumber),
    ...subaccountToCreate,
  }).returning('*');
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<SubaccountFromDatabase | undefined> {
  const baseQuery: QueryBuilder<SubaccountModel> = setupBaseQuery<SubaccountModel>(
    SubaccountModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function upsert(
  subaccountToUpsert: SubaccountCreateObject,
  options: Options = { txId: undefined },
): Promise<SubaccountFromDatabase> {
  const {
    address,
    subaccountNumber,
  } = subaccountToUpsert;
  const createdUuid: string = uuid(
    address,
    subaccountNumber,
  );

  const subaccounts: SubaccountModel[] = await SubaccountModel.query(
    Transaction.get(options.txId),
  ).upsert({
    ...subaccountToUpsert,
    id: createdUuid,
  }).returning('*');
  // should only ever be one subaccount
  return subaccounts[0];
}

export async function update(
  {
    id,
    ...fields
  }: SubaccountUpdateObject,
  options: Options = { txId: undefined },
): Promise<SubaccountFromDatabase | undefined> {
  const subaccount = await SubaccountModel.query(
    Transaction.get(options.txId),
  ).findById(id);
  const updatedSubaccount = await subaccount.$query().patch(fields as PartialModelObject<SubaccountModel>).returning('*');
  // The objection types mistakenly think the query returns an array of Subaccounts.
  return updatedSubaccount as unknown as (SubaccountFromDatabase | undefined);
}
