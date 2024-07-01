import { Big } from 'big.js';
import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  verifyAllRequiredFields,
  setupBaseQuery,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import TendermintEventModel from '../models/tendermint-event-model';
import {
  QueryConfig,
  TendermintEventFromDatabase,
  TendermintEventQueryConfig,
  TendermintEventColumns,
  TendermintEventCreateObject,
  Options,
  Ordering,
  QueryableField,
} from '../types';

const THIRTY_TWO_BITS_IN_BYTES: number = 4;

export function createEventId(
  blockHeight: string,
  transactionIndex: number,
  eventIndex: number,
): Buffer {

  const buffer = Buffer.alloc(3 * THIRTY_TWO_BITS_IN_BYTES);
  buffer.writeUInt32BE(Number(blockHeight), 0);
  // transactionIndex is -2 for BEGIN_BLOCK events, and -1 for END_BLOCK events.
  // Increment by 2 to ensure result is >= 0.
  buffer.writeUInt32BE(transactionIndex + 2, THIRTY_TWO_BITS_IN_BYTES);
  buffer.writeUInt32BE(eventIndex, 2 * THIRTY_TWO_BITS_IN_BYTES);

  return buffer;
}

export async function findAll(
  {
    id,
    blockHeight,
    transactionIndex,
    eventIndex,
    limit,
  }: TendermintEventQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TendermintEventFromDatabase[]> {
  verifyAllRequiredFields(
    {
      id,
      blockHeight,
      transactionIndex,
      eventIndex,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<TendermintEventModel> = setupBaseQuery<TendermintEventModel>(
    TendermintEventModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(TendermintEventColumns.id, id);
  }

  if (blockHeight !== undefined) {
    baseQuery = baseQuery.whereIn(TendermintEventColumns.blockHeight, blockHeight);
  }

  if (transactionIndex !== undefined) {
    baseQuery = baseQuery.whereIn(TendermintEventColumns.transactionIndex, transactionIndex);
  }

  if (eventIndex !== undefined) {
    baseQuery = baseQuery.whereIn(TendermintEventColumns.eventIndex, eventIndex);
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
      TendermintEventColumns.id,
      Ordering.ASC,
    );
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  tendermintEventToCreate: TendermintEventCreateObject,
  options: Options = { txId: undefined },
): Promise<TendermintEventFromDatabase> {
  return TendermintEventModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...tendermintEventToCreate,
    id: createEventId(
      tendermintEventToCreate.blockHeight,
      tendermintEventToCreate.transactionIndex,
      tendermintEventToCreate.eventIndex,
    ),
  }).returning('*');
}

export async function findById(
  id: Buffer,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TendermintEventFromDatabase | undefined> {
  const events: TendermintEventFromDatabase[] = await findAll(
    { id: [id] },
    [],
    options,
  );
  if (events.length === 0) {
    return undefined;
  } else {
    return events[0];
  }
}

export function compare(
  eventA: TendermintEventFromDatabase,
  eventB: TendermintEventFromDatabase,
): number {
  if (eventA.blockHeight !== eventB.blockHeight) {
    return Big(eventA.blockHeight).minus(Big(eventB.blockHeight)).toNumber();
  }

  if (eventA.transactionIndex !== eventB.transactionIndex) {
    return eventA.transactionIndex - eventB.transactionIndex;
  }

  if (eventA.eventIndex !== eventB.eventIndex) {
    return eventA.eventIndex - eventB.eventIndex;
  }

  return 0;
}
