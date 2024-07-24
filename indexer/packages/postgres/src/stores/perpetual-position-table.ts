import Knex from 'knex';
import _ from 'lodash';
import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { logAndThrowValidationError } from '../helpers/error-helpers';
import { knexPrimary, knexReadReplica } from '../helpers/knex';
import {
  verifyAllRequiredFields,
  setupBaseQuery,
  verifyAllInjectableVariables,
  setBulkRowsForUpdate,
  generateBulkUpdateString,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import PerpetualPositionModel from '../models/perpetual-position-model';
import {
  QueryConfig,
  PerpetualPositionFromDatabase,
  PerpetualPositionQueryConfig,
  PerpetualPositionColumns,
  PerpetualPositionCreateObject,
  PerpetualPositionStatus,
  PerpetualPositionUpdateObject,
  PerpetualPositionCloseObject,
  Options,
  Ordering,
  QueryableField,
  MarketOpenInterest,
  PerpetualPositionSubaccountUpdateObject,
  SubaccountToPerpetualPositionsMap,
} from '../types';

const DEFAULT_CREATE_FIELDS = {
  sumOpen: '0',
  sumClose: '0',
  entryPrice: '0',
};

const DEFAULT_SUBACCOUNT_UPDATE_DEFAULT_POSITION_FIELDS = {
  closedAt: null,
  closedAtHeight: null,
  closeEventId: null,
};

export function uuid(subaccountId: string, openEventId: Buffer): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(Buffer.from(`${subaccountId}-${openEventId.toString('hex')}`, BUFFER_ENCODING_UTF_8));
}

export async function findAll(
  {
    id,
    subaccountId,
    perpetualId,
    side,
    status,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    limit,
  }: PerpetualPositionQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PerpetualPositionFromDatabase[]> {
  verifyAllRequiredFields(
    {
      id,
      subaccountId,
      perpetualId,
      side,
      status,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<PerpetualPositionModel> = setupBaseQuery<PerpetualPositionModel>(
    PerpetualPositionModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(PerpetualPositionColumns.id, id);
  }

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(PerpetualPositionColumns.subaccountId, subaccountId);
  }

  if (perpetualId !== undefined) {
    baseQuery = baseQuery.whereIn(PerpetualPositionColumns.perpetualId, perpetualId);
  }

  if (side !== undefined) {
    baseQuery = baseQuery.where(PerpetualPositionColumns.side, side);
  }

  if (status !== undefined) {
    baseQuery = baseQuery.whereIn(PerpetualPositionColumns.status, status);
  }

  if (createdBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      PerpetualPositionColumns.createdAtHeight,
      '<=',
      createdBeforeOrAtHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(PerpetualPositionColumns.createdAt, '<=', createdBeforeOrAt);
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
      PerpetualPositionColumns.subaccountId,
      Ordering.ASC,
    ).orderBy(
      PerpetualPositionColumns.openEventId,
      Ordering.DESC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  perpetualPosition: PerpetualPositionCreateObject,
  options: Options = { txId: undefined },
): Promise<PerpetualPositionFromDatabase> {
  const perpetualPositionToCreate = {
    ...DEFAULT_CREATE_FIELDS,
    ...perpetualPosition,
  };

  return PerpetualPositionModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...perpetualPositionToCreate,
    id: uuid(perpetualPositionToCreate.subaccountId, perpetualPositionToCreate.openEventId),
  }).returning('*');
}

export async function update(
  {
    id,
    ...fields
  }: PerpetualPositionUpdateObject,
  options: Options = { txId: undefined },
): Promise<PerpetualPositionFromDatabase | undefined> {
  const perpetualPosition = await PerpetualPositionModel.query(
    Transaction.get(options.txId),
  // TODO fix expression typing so we dont have to use any
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ).findById(id).patch(fields as any).returning('*');
  // The objection types mistakenly think the query returns an array of perpetual positions.
  return perpetualPosition as unknown as (PerpetualPositionFromDatabase | undefined);
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PerpetualPositionFromDatabase | undefined> {
  const baseQuery: QueryBuilder<PerpetualPositionModel> = setupBaseQuery<PerpetualPositionModel>(
    PerpetualPositionModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function findOpenPositionForSubaccountPerpetual(
  subaccountId: string,
  perpetualId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PerpetualPositionFromDatabase | undefined> {
  const positions: PerpetualPositionFromDatabase[] = await findAll(
    {
      subaccountId: [subaccountId],
      perpetualId: [perpetualId],
      status: [PerpetualPositionStatus.OPEN],
    },
    [],
    options,
  );
  if (positions.length === 0) {
    return undefined;
  }

  return positions[0];
}

export async function findOpenPositionsForSubaccounts(
  subaccountIds: string[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<SubaccountToPerpetualPositionsMap> {
  const positions: PerpetualPositionFromDatabase[] = await findAll(
    {
      subaccountId: subaccountIds,
      status: [PerpetualPositionStatus.OPEN],
    },
    [],
    options,
  );
  if (positions.length === 0) {
    return {};
  }

  return _.reduce(positions,
    (acc: SubaccountToPerpetualPositionsMap, position: PerpetualPositionFromDatabase) => {
      const { subaccountId, perpetualId } = position;
      acc[subaccountId] = acc[subaccountId] || {};
      acc[subaccountId][perpetualId] = position;
      return acc;
    },
    {});
}

export async function closePosition(
  existingPosition: PerpetualPositionFromDatabase,
  perpetualPositionCloseObject: PerpetualPositionCloseObject,
  options: Options = { txId: undefined },
): Promise<PerpetualPositionFromDatabase | undefined> {
  const updateObject: PerpetualPositionUpdateObject = closePositionUpdateObject(
    existingPosition,
    perpetualPositionCloseObject,
  );
  return update(updateObject, options);
}

/**
 * Validates close position and returns the update object to update the position.
 */
export function closePositionUpdateObject(
  existingPosition: PerpetualPositionFromDatabase,
  perpetualPositionCloseObject: PerpetualPositionCloseObject,
): PerpetualPositionSubaccountUpdateObject {
  validateClosePosition(existingPosition);
  return {
    id: perpetualPositionCloseObject.id,
    closedAt: perpetualPositionCloseObject.closedAt,
    closedAtHeight: perpetualPositionCloseObject.closedAtHeight,
    closeEventId: perpetualPositionCloseObject.closeEventId,
    lastEventId: perpetualPositionCloseObject.closeEventId,
    settledFunding: perpetualPositionCloseObject.settledFunding,
    status: PerpetualPositionStatus.CLOSED,
    size: '0',
  };
}

/**
 * Throws an error if the position to close is already closed.
 * @param existingPosition
 */
function validateClosePosition(
  existingPosition: PerpetualPositionFromDatabase,
) {
  if (existingPosition.status === PerpetualPositionStatus.CLOSED) {
    logAndThrowValidationError('Unable to close because position is closed');
  }
}

// TODO(DEC-1821): Fix getOpenInterestLong to only return data for the ids requested
export async function getOpenInterestLong(perpetualMarketIds: string[]): Promise<
  _.Dictionary<MarketOpenInterest>
> {
  if (perpetualMarketIds.length === 0) {
    return {};
  }
  const perpetualMarketIdsSqlArray = `(${perpetualMarketIds.join(',')})`;
  const result: {
    rows: MarketOpenInterest[],
  } = await knexReadReplica.getConnection().raw(
    `SELECT
      "perpetualId" AS "perpetualMarketId",
      sum(size) AS "openInterest"
    FROM perpetual_positions
    WHERE "side"='LONG'
      AND "status"='OPEN'
      AND "perpetualId" IN ${perpetualMarketIdsSqlArray}
    GROUP BY "perpetualId";
    `,
  ) as unknown as {
    rows: MarketOpenInterest[],
  };

  const openInterestStats: {
    [perpetualMarketId: string]: MarketOpenInterest,
  } = _.keyBy(result.rows, 'perpetualMarketId');
  Object.values(perpetualMarketIds).forEach((perpetualMarketId) => {
    if (!openInterestStats[perpetualMarketId]) {
      // no positions exist for this market, set to 0
      openInterestStats[perpetualMarketId] = {
        perpetualMarketId,
        openInterest: '0',
      };
    }
  });

  return openInterestStats;
}

export async function bulkCreate(
  positions: PerpetualPositionCreateObject[],
  options: Options = { txId: undefined },
): Promise<PerpetualPositionFromDatabase[]> {
  const perpetualPositionsToCreate:
  PerpetualPositionCreateObject[] = _.map(positions, (position: PerpetualPositionCreateObject) => {
    return {
      ...DEFAULT_CREATE_FIELDS,
      ...position,
    };
  });

  return PerpetualPositionModel.query(
    Transaction.get(options.txId),
  ).insert(
    perpetualPositionsToCreate.map((position) => ({
      id: uuid(position.subaccountId, position.openEventId),
      ...position,
    })),
  ).returning('*');
}

/**
 * Bulk update for processing SubaccountUpdateEvents. Updates the following fields:
 * - closedAt
 * - closedAtHeight
 * - closeEventId
 * - lastEventId
 * - settledFunding
 * - status
 * - size
 * - maxSize
 * maxSize is calculated algorithmically based on the previous maxSize and the new size.
 */
export async function bulkUpdateSubaccountFields(
  positions: PerpetualPositionSubaccountUpdateObject[],
  options: Options = { txId: undefined },
): Promise<void> {
  if (positions.length === 0) {
    return;
  }
  const positionUpdatesWithDefaultValues: PerpetualPositionSubaccountUpdateObject[] = _.map(
    positions,
    (position: PerpetualPositionSubaccountUpdateObject) => {
      return {
        ...DEFAULT_SUBACCOUNT_UPDATE_DEFAULT_POSITION_FIELDS,
        ...position,
        [PerpetualPositionColumns.perpetualId]: undefined,
      };
    },
  );
  positionUpdatesWithDefaultValues.forEach(
    (position) => verifyAllInjectableVariables(Object.values(position)),
  );

  const columns: PerpetualPositionColumns[] = [
    PerpetualPositionColumns.id,
    PerpetualPositionColumns.closedAt,
    PerpetualPositionColumns.closedAtHeight,
    PerpetualPositionColumns.closeEventId,
    PerpetualPositionColumns.lastEventId,
    PerpetualPositionColumns.settledFunding,
    PerpetualPositionColumns.status,
    PerpetualPositionColumns.size,
  ];
  const positionRows: string[] = setBulkRowsForUpdate<PerpetualPositionColumns>({
    objectArray: positionUpdatesWithDefaultValues,
    columns,
    stringColumns: [
      PerpetualPositionColumns.id,
      PerpetualPositionColumns.status,
    ],
    numericColumns: [
      PerpetualPositionColumns.settledFunding,
      PerpetualPositionColumns.size,
    ],
    bigintColumns: [
      PerpetualPositionColumns.closedAtHeight,
    ],
    timestampColumns: [
      PerpetualPositionColumns.closedAt,
    ],
    binaryColumns: [
      PerpetualPositionColumns.closeEventId,
      PerpetualPositionColumns.lastEventId,
    ],
  });

  const query: string = generateBulkUpdateString({
    table: PerpetualPositionModel.tableName,
    objectRows: positionRows,
    columns,
    isUuid: true,
    uniqueIdentifier: PerpetualPositionColumns.id,
    setFieldsToAppend: [
      `"${PerpetualPositionColumns.maxSize}" = GREATEST("${PerpetualPositionColumns.maxSize}", c."${PerpetualPositionColumns.size}")`,
    ],
  });

  const transaction: Knex.Transaction | undefined = Transaction.get(options.txId);
  return transaction
    ? knexPrimary.raw(query).transacting(transaction)
    : knexPrimary.raw(query);
}
