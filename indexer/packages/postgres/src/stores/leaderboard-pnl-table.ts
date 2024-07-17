import Knex from 'knex';
import _ from 'lodash';
import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexPrimary } from '../helpers/knex';
import {
  verifyAllRequiredFields,
  setupBaseQuery,
  verifyAllInjectableVariables,
  setBulkRowsForUpdate,
  generateBulkUpsertString,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import LeaderboardPNLModel from '../models/leaderboard-pnl-model';
import {
  QueryConfig,
  LeaderboardPNLCreateObject,
  LeaderboardPNLFromDatabase,
  LeaderboardPNLColumns,
  LeaderboardPNLQueryConfig,
  Options,
  Ordering,
  QueryableField,
} from '../types';

export async function findAll(
  {
    subaccountId,
    timeSpan,
    rank,
    limit,
  }: LeaderboardPNLQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<LeaderboardPNLFromDatabase[]> {
  verifyAllRequiredFields(
    {
      subaccountId,
      timeSpan,
      rank,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<LeaderboardPNLModel> = setupBaseQuery<LeaderboardPNLModel>(
    LeaderboardPNLModel,
    options,
  );

  if (subaccountId) {
    baseQuery = baseQuery.whereIn(LeaderboardPNLColumns.subaccountId, subaccountId);
  }

  if (timeSpan) {
    baseQuery = baseQuery.whereIn(LeaderboardPNLColumns.timeSpan, timeSpan);
  }

  if (rank) {
    baseQuery = baseQuery.whereIn(LeaderboardPNLColumns.rank, rank);
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
      LeaderboardPNLColumns.rank,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  leaderboardPNLToCreate: LeaderboardPNLCreateObject,
  options: Options = { txId: undefined },
): Promise<LeaderboardPNLFromDatabase> {
  return LeaderboardPNLModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...leaderboardPNLToCreate,
  }).returning('*');
}

export async function upsert(
  LeaderboardPNLToUpsert: LeaderboardPNLCreateObject,
  options: Options = { txId: undefined },
): Promise<LeaderboardPNLFromDatabase> {
  const leaderboardPNLs: LeaderboardPNLModel[] = await LeaderboardPNLModel.query(
    Transaction.get(options.txId),
  ).upsert({
    ...LeaderboardPNLToUpsert,
  }).returning('*');
  return leaderboardPNLs[0];
}

export async function bulkUpsert(
  leaderboardPNLObjects: LeaderboardPNLCreateObject[],
  options: Options = { txId: undefined },
): Promise<void> {

  leaderboardPNLObjects.forEach(
    (leaderboardPNLObject: LeaderboardPNLCreateObject) => verifyAllInjectableVariables(
      Object.values(leaderboardPNLObject),
    ),
  );

  const columns: LeaderboardPNLColumns[] = _.keys(
    leaderboardPNLObjects[0]) as LeaderboardPNLColumns[];
  const rows: string[] = setBulkRowsForUpdate<LeaderboardPNLColumns>({
    objectArray: leaderboardPNLObjects,
    columns,
    numericColumns: [
      LeaderboardPNLColumns.rank,
    ],
    stringColumns: [
      LeaderboardPNLColumns.subaccountId,
      LeaderboardPNLColumns.timeSpan,
      LeaderboardPNLColumns.currentEquity,
      LeaderboardPNLColumns.pnl,
    ],
  });

  const query: string = generateBulkUpsertString({
    table: LeaderboardPNLModel.tableName,
    objectRows: rows,
    columns,
    uniqueIdentifiers: [LeaderboardPNLColumns.subaccountId, LeaderboardPNLColumns.timeSpan],
  });

  const transaction: Knex.Transaction | undefined = Transaction.get(options.txId);
  return transaction
    ? knexPrimary.raw(query).transacting(transaction)
    : knexPrimary.raw(query);
}
