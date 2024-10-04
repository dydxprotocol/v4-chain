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
import { knexReadReplica } from '../helpers/knex';

// Assuming block time of 1 second, this should be 4 hours of blocks
const FOUR_HOUR_OF_BLOCKS = Big(3600).times(4);
// Type used for querying for funding index maps for multiple effective heights.
interface FundingIndexUpdatesFromDatabaseWithSearchHeight extends FundingIndexUpdatesFromDatabase {
  // max effective height being queried for
  searchHeight: string,
};

export function uuid(
  blockHeight: string,
  eventId: Buffer,
  perpetualId: string,
): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
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
    ...fundingIndexUpdateToCreate,
    id: uuid(
      fundingIndexUpdateToCreate.effectiveAtHeight,
      fundingIndexUpdateToCreate.eventId,
      fundingIndexUpdateToCreate.perpetualId,
    ),
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
    // Optimization to reduce number of rows needed to scan
    .where(
      FundingIndexUpdatesColumns.effectiveAtHeight,
      '>',
      Big(effectiveBeforeOrAtHeight).minus(FOUR_HOUR_OF_BLOCKS).toFixed(),
    )
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

/**
 * Finds funding index maps for multiple effective before or at heights. Uses a SQL query unnesting
 * an array of effective before or at heights and cross-joining with the funding index updates table
 * to find the closest funding index update per effective before or at height.
 * @param effectiveBeforeOrAtHeights Heights to get funding index maps for.
 * @param options 
 * @returns Object mapping block heights to the respective funding index maps.
 */
export async function findFundingIndexMaps(
  effectiveBeforeOrAtHeights: string[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<{[blockHeight: string]: FundingIndexMap}> {
  const heightNumbers: number[] = effectiveBeforeOrAtHeights
    .map((height: string) => parseInt(height, 10))
    .sort();
  // Get the min height to limit the search to blocks 4 hours or before the min height.
  const minHeight: number = heightNumbers[0];

  const result: {
    rows: FundingIndexUpdatesFromDatabaseWithSearchHeight[]
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT
      DISTINCT ON ("perpetualId", "searchHeight") "perpetualId", "searchHeight",
      "funding_index_updates".*
    FROM
      "funding_index_updates",
      unnest(ARRAY[${heightNumbers.join(',')}]) AS "searchHeight"
    WHERE
      "effectiveAtHeight" > ${Big(minHeight).minus(FOUR_HOUR_OF_BLOCKS).toFixed()} AND
      "effectiveAtHeight" <= "searchHeight"
    ORDER BY
      "perpetualId",
      "searchHeight",
      "effectiveAtHeight" DESC
    `,
  ) as unknown as {
    rows: FundingIndexUpdatesFromDatabaseWithSearchHeight[],
  };

  const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
    {},
    [],
    options,
  );

  const fundingIndexMaps:{[blockHeight: string]: FundingIndexMap} = {};
  for (const height of effectiveBeforeOrAtHeights) {
    fundingIndexMaps[height] = _.reduce(perpetualMarkets,
      (acc: FundingIndexMap, perpetualMarket: PerpetualMarketFromDatabase): FundingIndexMap => {
        acc[perpetualMarket.id] = Big(0);
        return acc;
      },
      {},
    );
  }
  for (const funding of result.rows) {
    fundingIndexMaps[funding.searchHeight][funding.perpetualId] = Big(funding.fundingIndex)
  }

  return fundingIndexMaps;
}
