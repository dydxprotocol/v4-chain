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
import LeaderboardPnlModel from '../models/leaderboard-pnl-model';
import {
  QueryConfig,
  LeaderboardPnlCreateObject,
  LeaderboardPnlFromDatabase,
  LeaderboardPnlColumns,
  LeaderboardPnlQueryConfig,
  Options,
  Ordering,
  QueryableField,
} from '../types';

export async function findAll(
  {
    address,
    timeSpan,
    rank,
    limit,
  }: LeaderboardPnlQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<LeaderboardPnlFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      timeSpan,
      rank,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<LeaderboardPnlModel> = setupBaseQuery<LeaderboardPnlModel>(
    LeaderboardPnlModel,
    options,
  );

  if (address) {
    baseQuery = baseQuery.whereIn(LeaderboardPnlColumns.address, address);
  }

  if (timeSpan) {
    baseQuery = baseQuery.whereIn(LeaderboardPnlColumns.timeSpan, timeSpan);
  }

  if (rank) {
    baseQuery = baseQuery.whereIn(LeaderboardPnlColumns.rank, rank);
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
      LeaderboardPnlColumns.rank,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  leaderboardPnlToCreate: LeaderboardPnlCreateObject,
  options: Options = { txId: undefined },
): Promise<LeaderboardPnlFromDatabase> {
  return LeaderboardPnlModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...leaderboardPnlToCreate,
  }).returning('*');
}

export async function upsert(
  LeaderboardPnlToUpsert: LeaderboardPnlCreateObject,
  options: Options = { txId: undefined },
): Promise<LeaderboardPnlFromDatabase> {
  const leaderboardPnls: LeaderboardPnlModel[] = await LeaderboardPnlModel.query(
    Transaction.get(options.txId),
  ).upsert({
    ...LeaderboardPnlToUpsert,
  }).returning('*');
  return leaderboardPnls[0];
}

export async function bulkUpsert(
  leaderboardPnlObjects: LeaderboardPnlCreateObject[],
  options: Options = { txId: undefined },
): Promise<void> {
  leaderboardPnlObjects.forEach(
    (leaderboardPnlObject: LeaderboardPnlCreateObject) => verifyAllInjectableVariables(
      Object.values(leaderboardPnlObject),
    ),
  );

  const columns: LeaderboardPnlColumns[] = _.keys(
    leaderboardPnlObjects[0]) as LeaderboardPnlColumns[];
  const rows: string[] = setBulkRowsForUpdate<LeaderboardPnlColumns>({
    objectArray: leaderboardPnlObjects,
    columns,
    numericColumns: [
      LeaderboardPnlColumns.rank,
    ],
    stringColumns: [
      LeaderboardPnlColumns.address,
      LeaderboardPnlColumns.timeSpan,
      LeaderboardPnlColumns.currentEquity,
      LeaderboardPnlColumns.pnl,
    ],
  });

  const query: string = generateBulkUpsertString({
    table: LeaderboardPnlModel.tableName,
    objectRows: rows,
    columns,
    uniqueIdentifiers: [LeaderboardPnlColumns.address, LeaderboardPnlColumns.timeSpan],
  });

  const transaction: Knex.Transaction | undefined = Transaction.get(options.txId);
  return transaction
    ? knexPrimary.raw(query).transacting(transaction)
    : knexPrimary.raw(query);
}
