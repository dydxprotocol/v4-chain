import Big from 'big.js';
import _ from 'lodash';
import { QueryBuilder } from 'objection';

import { BUFFER_ENCODING_UTF_8, DEFAULT_POSTGRES_OPTIONS } from '../constants';
import {
  verifyAllRequiredFields,
  setupBaseQuery,
  rawQuery,
} from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import { getUuid } from '../helpers/uuid';
import TransferModel from '../models/transfer-model';
import {
  QueryConfig,
  TransferFromDatabase,
  TransferQueryConfig,
  TransferColumns,
  TransferCreateObject,
  Options,
  QueryableField,
  ToAndFromSubaccountTransferQueryConfig,
  SubaccountAssetNetTransferMap,
} from '../types';

export function uuid(
  eventId: Buffer,
  assetId: string,
  senderSubaccountId?: string,
  recipientSubaccountId?: string,
  senderWalletAddress?: string,
  recipientWalletAddress?: string,
): string {
  return getUuid(
    Buffer.from(
      `${senderSubaccountId}-${recipientSubaccountId}-${senderWalletAddress}-${recipientWalletAddress}-${eventId.toString('hex')}-${assetId}`,
      BUFFER_ENCODING_UTF_8),
  );
}

interface SubaccountAssetNetTransfer {
  subaccountId: string;
  assetId: string;
  totalSize: string;
}

export async function findAll(
  {
    limit,
    id,
    senderSubaccountId,
    recipientSubaccountId,
    senderWalletAddress,
    recipientWalletAddress,
    assetId,
    size,
    eventId,
    transactionHash,
    createdAt,
    createdAtHeight,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    createdAfter,
    createdAfterHeight,
  }: TransferQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TransferFromDatabase[]> {
  verifyAllRequiredFields(
    {
      limit,
      id,
      senderSubaccountId,
      recipientSubaccountId,
      senderWalletAddress,
      recipientWalletAddress,
      assetId,
      size,
      eventId,
      transactionHash,
      createdAt,
      createdAtHeight,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdAfter,
      createdAfterHeight,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<TransferModel> = setupBaseQuery<TransferModel>(
    TransferModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.id, id);
  }

  if (senderSubaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.senderSubaccountId, senderSubaccountId);
  }

  if (recipientSubaccountId !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.recipientSubaccountId, recipientSubaccountId);
  }

  if (senderWalletAddress !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.senderWalletAddress, senderWalletAddress);
  }

  if (recipientWalletAddress !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.recipientWalletAddress, recipientWalletAddress);
  }

  if (assetId !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.assetId, assetId);
  }

  if (size !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.size, size);
  }

  if (eventId !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.eventId, eventId);
  }

  if (transactionHash !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.transactionHash, transactionHash);
  }

  if (createdAt !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.createdAt, createdAt);
  }

  if (createdAtHeight !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.createdAtHeight, createdAtHeight);
  }

  if (createdBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      TransferColumns.createdAtHeight,
      '<=',
      createdBeforeOrAtHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(TransferColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdAfterHeight !== undefined) {
    baseQuery = baseQuery.where(
      TransferColumns.createdAtHeight,
      '>',
      createdAfterHeight,
    );
  }

  if (createdAfter !== undefined) {
    baseQuery = baseQuery.where(TransferColumns.createdAt, '>', createdAfter);
  }

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(
        column,
        order,
      );
    }
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

// finds all transfers to or from the subaccount id
export async function findAllToOrFromSubaccountId(
  {
    limit,
    id,
    subaccountId,
    assetId,
    size,
    eventId,
    transactionHash,
    createdAt,
    createdAtHeight,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    createdAfterHeight,
    createdAfter,
  }: ToAndFromSubaccountTransferQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TransferFromDatabase[]> {
  verifyAllRequiredFields(
    {
      limit,
      id,
      subaccountId,
      assetId,
      size,
      eventId,
      transactionHash,
      createdAt,
      createdAtHeight,
      createdBeforeOrAtHeight,
      createdBeforeOrAt,
      createdAfterHeight,
      createdAfter,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<TransferModel> = setupBaseQuery<TransferModel>(
    TransferModel,
    options,
  );

  if (id !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.id, id);
  }

  if (subaccountId !== undefined) {
    baseQuery = baseQuery.where((queryBuilder) => {
      // eslint-disable-next-line no-void
      void queryBuilder.whereIn(TransferColumns.recipientSubaccountId, subaccountId)
        .orWhereIn(TransferColumns.senderSubaccountId, subaccountId);
    });
  }

  if (assetId !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.assetId, assetId);
  }

  if (size !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.size, size);
  }

  if (eventId !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.eventId, eventId);
  }

  if (transactionHash !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.transactionHash, transactionHash);
  }

  if (createdAt !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.createdAt, createdAt);
  }

  if (createdAtHeight !== undefined) {
    baseQuery = baseQuery.whereIn(TransferColumns.createdAtHeight, createdAtHeight);
  }

  if (createdBeforeOrAtHeight !== undefined) {
    baseQuery = baseQuery.where(
      TransferColumns.createdAtHeight,
      '<=',
      createdBeforeOrAtHeight,
    );
  }

  if (createdBeforeOrAt !== undefined) {
    baseQuery = baseQuery.where(TransferColumns.createdAt, '<=', createdBeforeOrAt);
  }

  if (createdAfterHeight !== undefined) {
    baseQuery = baseQuery.where(
      TransferColumns.createdAtHeight,
      '>',
      createdAfterHeight,
    );
  }

  if (createdAfter !== undefined) {
    baseQuery = baseQuery.where(TransferColumns.createdAt, '>', createdAfter);
  }

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(
        column,
        order,
      );
    }
  }

  if (limit !== undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

function convertToSubaccountAssetMap(
  transfers: SubaccountAssetNetTransfer[],
): SubaccountAssetNetTransferMap {
  const assetGroups: _.Dictionary<SubaccountAssetNetTransfer[]> = _.groupBy(
    transfers,
    'subaccountId');

  return _.mapValues(assetGroups, (group: SubaccountAssetNetTransfer[]) => {
    return _.reduce(group, (result: {}, asset: SubaccountAssetNetTransfer) => {
      return { ...result, [asset.assetId]: asset.totalSize };
    }, {});
  });
}

export interface AssetTransferMap {
  [assetId: string]: Big;
}

export async function getNetTransfersBetweenBlockHeightsForSubaccount(
  subaccountId: string,
  createdAfterHeight: string,
  createdBeforeOrAtHeight: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<AssetTransferMap> {
  const queryString: string = `
    SELECT 
      "assetId",
      SUM(
        CASE 
          WHEN "senderSubaccountId" = '${subaccountId}' THEN -"size"
          ELSE "size"
        END
      ) AS "totalSize"
    FROM 
      "transfers"
    WHERE 
      (
        "senderSubaccountId" = '${subaccountId}' 
        OR "recipientSubaccountId" = '${subaccountId}'
      )
      AND "createdAtHeight" > ${createdAfterHeight} 
      AND "createdAtHeight" <= ${createdBeforeOrAtHeight}
    GROUP BY 
      "assetId";
  `;

  const result: {
    rows: {
      assetId: string;
      totalSize: string;
    }[],
  } = await rawQuery(queryString, options);
  return _.mapValues(_.keyBy(result.rows, 'assetId'), (row: { assetId: string, totalSize: string }) => {
    return Big(row.totalSize);
  });
}

export async function getNetTransfersPerSubaccount(
  createdBeforeOrAtHeight: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<SubaccountAssetNetTransferMap> {
  // Get the net value of transfers since beginning of time up until createdBeforeOrAtHeight
  // for all subaccounts. If a subaccount is sending an asset, the value will be negative.
  // If a subaccount is receiving an asset, the value will be positive.
  const queryString: string = `
  SELECT 
    sub."subaccountId",
    sub."assetId",
    SUM(sub."size") AS "totalSize"
  FROM (
    SELECT DISTINCT 
      "senderSubaccountId" AS "subaccountId",
      "assetId",
      -"size" AS "size",
      "id"
    FROM 
      "transfers"
    WHERE "transfers"."createdAtHeight" <= ${createdBeforeOrAtHeight}
    UNION 
    SELECT DISTINCT 
      "recipientSubaccountId" AS "subaccountId",
      "assetId",
      "size" AS "size",
      "id"
    FROM 
      "transfers"
    WHERE "transfers"."createdAtHeight" <= ${createdBeforeOrAtHeight}
  ) AS sub
  GROUP BY 
    sub."subaccountId",
    sub."assetId";
  `;

  const result: {
    rows: SubaccountAssetNetTransfer[],
  } = await rawQuery(queryString, options);

  const assetsPerSubaccount: SubaccountAssetNetTransfer[] = result.rows;
  return convertToSubaccountAssetMap(assetsPerSubaccount);
}

export async function create(
  transferToCreate: TransferCreateObject,
  options: Options = { txId: undefined },
): Promise<TransferFromDatabase> {
  return TransferModel.query(
    Transaction.get(options.txId),
  ).insert({
    id: uuid(
      transferToCreate.eventId,
      transferToCreate.assetId,
      transferToCreate.senderSubaccountId,
      transferToCreate.recipientSubaccountId,
      transferToCreate.senderWalletAddress,
      transferToCreate.recipientWalletAddress,
    ),
    ...transferToCreate,
  }).returning('*');
}

export async function findById(
  id: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<TransferFromDatabase | undefined> {
  const baseQuery: QueryBuilder<TransferModel> = setupBaseQuery<TransferModel>(
    TransferModel,
    options,
  );
  return baseQuery
    .findById(id)
    .returning('*');
}
