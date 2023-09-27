import Big from 'big.js';
import _ from 'lodash';
import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import FundingIndexUpdatesModel from '../models/funding-index-updates-model';
import {
  Options,
  FundingIndexUpdatesColumns,
  FundingIndexUpdatesCreateObject,
  FundingIndexUpdatesFromDatabase,
  FundingIndexUpdatesQueryConfig,
  Ordering,
  QueryableField,
  QueryConfig,
  FundingIndexMap,
  PerpetualMarketFromDatabase,
} from '../types';
import * as PerpetualMarketTable from './perpetual-market-table';

export function uuid(
  blockHeight: string,
  eventId: Buffer,
  perpetualId: string,
): string {
  return getUuid(Buffer.from(`${blockHeight}-${eventId.toString('hex')}-${perpetualId}`, BUFFER_ENCODING_UTF_8));
}

export async function findAll(
  {
    limit,
    id,
    perpetualId,
    eventId,
    effectiveAt,
    effectiveAtHeight,
    effectiveBeforeOrAt,
    effectiveBeforeOrAtHeight,
  }: FundingIndexUpdatesQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<FundingIndexUpdatesFromDatabase[]> {
  verifyAllRequiredFields(
    {
      limit,
      id,
      perpetualId,
      eventId,
      effectiveAt,
      effectiveAtHeight,
      effectiveBeforeOrAt,
      effectiveBeforeOrAtHeight,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<FundingIndexUpdatesModel> = setupBaseQuery<FundingIndexUpdatesModel>(
    FundingIndexUpdatesModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(FundingIndexUpdatesColumns.id, id);
  }
  if (perpetualId !== undefined) {
    baseQuery = baseQuery.whereIn(FundingIndexUpdatesColumns.perpetualId, perpetualId);
  }
  if (eventId !== undefined) {
    baseQuery = baseQuery.where(FundingIndexUpdatesColumns.eventId, eventId);
  }
  if (effectiveAt !== undefined) {
    baseQuery = baseQuery.where(FundingIndexUpdatesColumns.effectiveAt, effectiveAt);
  }
  if (effectiveAtHeight !== undefined) {
    baseQuery = baseQuery.where(FundingIndexUpdatesColumns.effectiveAtHeight, effectiveAtHeight);
  }

  if (effectiveBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(FundingIndexUpdatesColumns.effectiveAt, '<=', effectiveBeforeOrAt);
  }

  if (effectiveBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      FundingIndexUpdatesColumns.effectiveAtHeight,
      '<=',
      effectiveBeforeOrAtHeight,
    );
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
      FundingIndexUpdatesColumns.effectiveAtHeight,
      Ordering.DESC,
    ).orderBy(
      FundingIndexUpdatesColumns.eventId,
      Ordering.DESC,
    );
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  fundingIndexUpdateToCreate: FundingIndexUpdatesCreateObject,
  options: Options = { txId: undefined },
): Promise<FundingIndexUpdatesFromDatabase> {
  return FundingIndexUpdatesModel.query(
    Transaction.get(options.txId),
  ).insert({
    id: uuid(
      fundingIndexUpdateToCreate.effectiveAtHeight,
      fundingIndexUpdateToCreate.eventId,
      fundingIndexUpdateToCreate.perpetualId,
    ),
    ...fundingIndexUpdateToCreate,
  }).returning('*');
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<FundingIndexUpdatesFromDatabase | undefined> {
  const baseQuery: QueryBuilder<FundingIndexUpdatesModel> = setupBaseQuery<
    FundingIndexUpdatesModel>(
      FundingIndexUpdatesModel,
      options,
    );
  return baseQuery
    .findById(id)
    .returning('*');
}

export async function findMostRecentMarketFundingIndexUpdate(
  perpetualId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<FundingIndexUpdatesFromDatabase | undefined> {
  const baseQuery: QueryBuilder<FundingIndexUpdatesModel> = setupBaseQuery<
    FundingIndexUpdatesModel>(
      FundingIndexUpdatesModel,
      options,
    );

  const fundingIndexUpdates: FundingIndexUpdatesFromDatabase[] = await baseQuery
    .where(FundingIndexUpdatesColumns.perpetualId, perpetualId)
    .orderBy(FundingIndexUpdatesColumns.effectiveAtHeight, Ordering.DESC)
    .limit(1)
    .returning('*');

  if (fundingIndexUpdates.length === 0) {
    return undefined;
  }
  return fundingIndexUpdates[0];
}

export async function findFundingIndexMap(
  effectiveBeforeOrAtHeight: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<FundingIndexMap> {
  // TODO(IND-39): Remove this default when the database is seeded using events emitted from
  // protocol during genesis.
  // Default funding index per perpetual market is 0.
  const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
    {},
    [],
    options,
  );
  const initialFundingIndexMap: FundingIndexMap = _.reduce(perpetualMarkets,
    (acc: FundingIndexMap, perpetualMarket: PerpetualMarketFromDatabase): FundingIndexMap => {
      acc[perpetualMarket.id] = Big(0);
      return acc;
    },
    {},
  );

  const baseQuery: QueryBuilder<FundingIndexUpdatesModel> = setupBaseQuery<
    FundingIndexUpdatesModel>(
      FundingIndexUpdatesModel,
      options,
    );

  const fundingIndexUpdates: FundingIndexUpdatesFromDatabase[] = await baseQuery
    .distinctOn(FundingIndexUpdatesColumns.perpetualId)
    .where(FundingIndexUpdatesColumns.effectiveAtHeight, '<=', effectiveBeforeOrAtHeight)
    .orderBy(FundingIndexUpdatesColumns.perpetualId)
    .orderBy(FundingIndexUpdatesColumns.effectiveAtHeight, Ordering.DESC)
    .returning('*');

  return _.reduce(fundingIndexUpdates,
    (acc: FundingIndexMap, fundingIndexUpdate: FundingIndexUpdatesFromDatabase) => {
      acc[fundingIndexUpdate.perpetualId] = Big(fundingIndexUpdate.fundingIndex);
      return acc;
    },
    initialFundingIndexMap,
  );
}
