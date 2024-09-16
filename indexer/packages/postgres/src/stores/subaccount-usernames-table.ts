import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexReadReplica } from '../helpers/knex';
import {
  verifyAllRequiredFields,
  setupBaseQuery,
  rawQuery,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import SubaccountUsernamesModel from '../models/subaccount-usernames-model';
import {
  QueryConfig,
  SubaccountUsernamesFromDatabase,
  SubaccountUsernamesQueryConfig,
  SubaccountUsernamesColumns,
  SubaccountUsernamesCreateObject,
  SubaccountsWithoutUsernamesResult,
  Options,
  Ordering,
  QueryableField,
  AddressUsernameFromDatabase,
} from '../types';

export async function findAll(
  {
    username,
    subaccountId,
    limit,
  }: SubaccountUsernamesQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<SubaccountUsernamesFromDatabase[]> {
  verifyAllRequiredFields(
    {
      username,
      subaccountId,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<SubaccountUsernamesModel> = setupBaseQuery<SubaccountUsernamesModel>(
    SubaccountUsernamesModel,
    options,
  );

  if (username) {
    baseQuery = baseQuery.whereIn(SubaccountUsernamesColumns.username, username);
  }

  if (subaccountId) {
    baseQuery = baseQuery.whereIn(SubaccountUsernamesColumns.subaccountId, subaccountId);
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
      SubaccountUsernamesColumns.username,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  subaccountUsernameToCreate: SubaccountUsernamesCreateObject,
  options: Options = { txId: undefined },
): Promise<SubaccountUsernamesFromDatabase> {
  return SubaccountUsernamesModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...subaccountUsernameToCreate,
  }).returning('*');
}

export async function findByUsername(
  username: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<SubaccountUsernamesFromDatabase | undefined> {
  const baseQuery:
  QueryBuilder<SubaccountUsernamesModel> = setupBaseQuery<SubaccountUsernamesModel>(
    SubaccountUsernamesModel,
    options,
  );
  return (await baseQuery).find((subaccountUsername) => subaccountUsername.username === username);
}

export async function getSubaccountsWithoutUsernames(
  options: Options = DEFAULT_POSTGRES_OPTIONS):
  Promise<SubaccountsWithoutUsernamesResult[]> {
  const queryString: string = `
    SELECT id as "subaccountId"
    FROM subaccounts
    WHERE id NOT IN (
      SELECT "subaccountId" FROM subaccount_usernames
    )
    AND subaccounts."subaccountNumber"=0
  `;

  const result: {
    rows: SubaccountsWithoutUsernamesResult[],
  } = await rawQuery(queryString, options);

  return result.rows;
}

export async function findByAddress(
  addresses: string[],
): Promise<AddressUsernameFromDatabase[]> {
  if (addresses.length === 0) {
    return [];
  }

  const result: { rows: AddressUsernameFromDatabase[] } = await knexReadReplica
    .getConnection()
    .raw(
      `
      WITH subaccountIds AS (
        SELECT "id", "address"
        FROM subaccounts
        WHERE "address" = ANY(?)
        AND "subaccountNumber" = 0
      )
      SELECT s."address", u."username"
      FROM subaccountIds s
      LEFT JOIN subaccount_usernames u ON u."subaccountId" = s."id"
      `,
      [addresses],
    );

  return result.rows;
}
