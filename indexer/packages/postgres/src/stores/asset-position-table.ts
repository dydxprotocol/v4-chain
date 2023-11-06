import Big from 'big.js';
import _ from 'lodash';
import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS, USDC_ASSET_ID } from '../constants';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import AssetPositionModel from '../models/asset-position-model';
import {
  AssetPositionColumns,
  AssetPositionCreateObject,
  AssetPositionFromDatabase,
  AssetPositionQueryConfig,
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  SubaccountUsdcMap,
} from '../types';

export function uuid(subaccountId: string, assetId: string): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(Buffer.from(`${subaccountId}-${assetId}`, BUFFER_ENCODING_UTF_8));
}

export async function findAll(
  {
    limit,
    id,
    subaccountId,
    assetId,
    size,
    isLong,
  }: AssetPositionQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AssetPositionFromDatabase[]> {
  verifyAllRequiredFields(
    {
      limit,
      id,
      subaccountId,
      assetId,
      size,
      isLong,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<AssetPositionModel> = setupBaseQuery<AssetPositionModel>(
    AssetPositionModel,
    options,
  );

  if (id) {
    baseQuery = baseQuery.whereIn(AssetPositionColumns.id, id);
  }
  if (subaccountId) {
    baseQuery = baseQuery.whereIn(AssetPositionColumns.subaccountId, subaccountId);
  }
  if (assetId) {
    baseQuery = baseQuery.whereIn(AssetPositionColumns.assetId, assetId);
  }
  if (size) {
    baseQuery = baseQuery.where(AssetPositionColumns.size, size);
  }
  if (isLong) {
    baseQuery = baseQuery.where(AssetPositionColumns.isLong, isLong);
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
      AssetPositionColumns.assetId,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function upsert(
  assetPositionToUpsert: AssetPositionCreateObject,
  options: Options = { txId: undefined },
): Promise<AssetPositionFromDatabase> {
  const {
    subaccountId,
    assetId,
  } = assetPositionToUpsert;
  const createdUuid: string = uuid(
    subaccountId,
    assetId,
  );

  const assetPositions: AssetPositionModel[] = await AssetPositionModel.query(
    Transaction.get(options.txId),
  ).upsert({
    ...assetPositionToUpsert,
    id: createdUuid,
  }).returning('*');
  // should only ever be one asset position
  return assetPositions[0];
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AssetPositionFromDatabase | undefined> {
  const baseQuery: QueryBuilder<AssetPositionModel> = setupBaseQuery<AssetPositionModel>(
    AssetPositionModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}

function convertToSubaccountUsdcMap(
  usdcPositions: AssetPositionFromDatabase[],
): SubaccountUsdcMap {
  return _.mapValues(_.keyBy(usdcPositions, 'subaccountId'), (asset) => {
    return Big(asset.isLong ? asset.size : -asset.size);
  });
}

export async function findUsdcPositionForSubaccounts(
  subaccountIds: string[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<{ [subaccountId: string]: Big }> {
  const positions: AssetPositionFromDatabase[] = await findAll(
    {
      subaccountId: subaccountIds,
      assetId: [USDC_ASSET_ID],
    },
    [],
    options,
  );
  if (positions.length === 0) {
    return {};
  }
  return convertToSubaccountUsdcMap(positions);
}
