import _ from 'lodash';
import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS, ZERO_TIME_ISO_8601 } from '../constants';
import { knexReadReplica } from '../helpers/knex';
import { setupBaseQuery, verifyAllInjectableVariables, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import PnlTicksModel from '../models/pnl-ticks-model';
import {
  Options,
  Ordering,
  PnlTicksColumns,
  PnlTicksCreateObject,
  PnlTicksFromDatabase,
  PnlTicksQueryConfig,
  QueryableField,
  QueryConfig,
} from '../types';

export function uuid(
  subaccountId: string,
  createdAt: string,
): string {
  return getUuid(
    Buffer.from(
      `${subaccountId}-${createdAt}`,
      BUFFER_ENCODING_UTF_8),
  );
}

export async function findAll(
  {
    limit,
    id,
    subaccountId,
    createdAt,
    blockHeight,
    blockTime,
    createdBeforeOrAt,
    createdBeforeOrAtBlockHeight,
    createdOnOrAfter,
    createdOnOrAfterBlockHeight,
  }: PnlTicksQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PnlTicksFromDatabase[]> {
  verifyAllRequiredFields(
    {
      limit,
      id,
      subaccountId,
      createdAt,
      blockHeight,
      blockTime,
      createdBeforeOrAt,
      createdBeforeOrAtBlockHeight,
      createdOnOrAfter,
      createdOnOrAfterBlockHeight,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<PnlTicksModel> = setupBaseQuery<PnlTicksModel>(
    PnlTicksModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(PnlTicksColumns.id, id);
  }

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(PnlTicksColumns.subaccountId, subaccountId);
  }

  if (createdAt !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.createdAt, createdAt);
  }

  if (blockHeight !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.blockHeight, blockHeight);
  }

  if (blockTime !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.blockTime, blockTime);
  }

  if (createdBeforeOrAtBlockHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlTicksColumns.blockHeight,
      '<=',
      createdBeforeOrAtBlockHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdOnOrAfterBlockHeight !== undefined) {
    baseQuery = baseQuery.where(
      PnlTicksColumns.blockHeight,
      '>=',
      createdOnOrAfterBlockHeight,
    );
  }

  if (createdOnOrAfter !== undefined) {
    baseQuery = baseQuery.where(PnlTicksColumns.createdAt, '>=', createdOnOrAfter);
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
      PnlTicksColumns.subaccountId,
      Ordering.ASC,
    ).orderBy(
      PnlTicksColumns.blockHeight,
      Ordering.DESC,
    );
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  pnlTicksToCreate: PnlTicksCreateObject,
  options: Options = { txId: undefined },
): Promise<PnlTicksFromDatabase> {
  return PnlTicksModel.query(
    Transaction.get(options.txId),
  ).insert({
    id: uuid(pnlTicksToCreate.subaccountId, pnlTicksToCreate.createdAt),
    ...pnlTicksToCreate,
  }).returning('*');
}

export async function createMany(
  pnlTicks: PnlTicksCreateObject[],
  options: Options = { txId: undefined },
): Promise<PnlTicksFromDatabase[]> {
  const ticks: PnlTicksFromDatabase[] = pnlTicks.map((tick) => ({
    ...tick,
    id: uuid(tick.subaccountId, tick.createdAt),
  }));

  return PnlTicksModel
    .query(Transaction.get(options.txId))
    .insert(ticks)
    .returning('*');
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PnlTicksFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PnlTicksModel> = setupBaseQuery<PnlTicksModel>(
    PnlTicksModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

function convertPnlTicksFromDatabaseToPnlTicksCreateObject(
  pnlTicksFromDatabase: PnlTicksFromDatabase,
): PnlTicksCreateObject {
  return _.omit(pnlTicksFromDatabase, PnlTicksColumns.id);
}

export async function findLatestProcessedBlocktime(): Promise<string> {
  const result: {
    rows: [{ max: string }]
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT MAX("blockTime")
    FROM "pnl_ticks"
    `
    ,
  ) as unknown as { rows: [{ max: string }] };
  return result.rows[0].max || ZERO_TIME_ISO_8601;
}

export async function findMostRecentPnlTickForEachAccount(
  createdOnOrAfterHeight: string,
): Promise<{
  [subaccountId: string]: PnlTicksCreateObject
}> {
  verifyAllInjectableVariables([createdOnOrAfterHeight]);

  const result: {
    rows: PnlTicksFromDatabase[]
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT DISTINCT ON ("subaccountId") *
    FROM "pnl_ticks"
    WHERE "blockHeight" >= '${createdOnOrAfterHeight}'
    ORDER BY "subaccountId" ASC, "blockHeight" DESC, "createdAt" DESC;
    `
    ,
  ) as unknown as { rows: PnlTicksFromDatabase[] };
  return _.keyBy(
    result.rows.map(convertPnlTicksFromDatabaseToPnlTicksCreateObject),
    'subaccountId',
  );
}
