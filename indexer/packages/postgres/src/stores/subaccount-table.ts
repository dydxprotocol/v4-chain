import { IndexerSubaccountId } from '@dydxprotocol-indexer/v4-protos';
import { PartialModelObject, QueryBuilder } from 'objection';

import config from '../config';
import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS, MAX_PARENT_SUBACCOUNTS } from '../constants';
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

  if (subaccountNumber !== undefined) {
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
  options: Options = DEFAULT_POSTGRES_OPTIONS,
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
    ...subaccountToCreate,
    id: uuid(subaccountToCreate.address, subaccountToCreate.subaccountNumber),
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

export async function deleteById(
  id: string,
  options: Options = { txId: undefined },
): Promise<void> {
  if (config.NODE_ENV !== 'test') {
    throw new Error('Subaccount deletion is not allowed in non-test environments');
  }

  await SubaccountModel.query(
    Transaction.get(options.txId),
  ).deleteById(id);
}

/**
 * Retrieves all subaccount IDs associated with a parent subaccount.
 * A subaccount is considered a child of the parent if it has the same address
 * and its subaccount number follows the modulo relationship with the parent.
 *
 * @param parentSubaccount The parent subaccount object with address and subaccountNumber
 * @param options Query options including transaction ID
 * @returns A promise that resolves to an array of subaccount ID strings
 */
export async function findIdsForParentSubaccount(
  parentSubaccount: {
    address: string,
    subaccountNumber: number,
  },
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string[]> {
  // Get all subaccounts for the address
  const subaccounts = await findAll(
    { address: parentSubaccount.address },
    [],
    options,
  );

  // Filter for subaccounts that match the parent relationship
  // (subaccountNumber - parentSubaccountNumber) % MAX_PARENT_SUBACCOUNTS = 0
  return subaccounts
    .filter((subaccount) => (subaccount.subaccountNumber - parentSubaccount.subaccountNumber) %
     MAX_PARENT_SUBACCOUNTS === 0,
    )
    .map((subaccount) => subaccount.id);
}
