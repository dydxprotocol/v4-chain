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
  ParentSubaccountTransferQueryConfig,
  SubaccountAssetNetTransferMap,
  PaginationFromDatabase,
} from '../types';

export function uuid(
  eventId: Buffer,
  assetId: string,
  senderSubaccountId?: string,
  recipientSubaccountId?: string,
  senderWalletAddress?: string,
  recipientWalletAddress?: string,
): string {
  // TODO(IND-483): Fix all uuid string substitutions to use Array.join.
  return getUuid(
    Buffer.from(
      `${senderSubaccountId}-${recipientSubaccountId}-${senderWalletAddress}-${recipientWalletAddress}-${eventId.toString('hex')}-${assetId}`,
      BUFFER_ENCODING_UTF_8),
  );
}

interface SubaccountAssetNetTransfer {
  subaccountId: string,
  assetId: string,
  totalSize: string,
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
    page,
  }: ToAndFromSubaccountTransferQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<TransferFromDatabase>> {
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

  if (limit !== undefined && page === undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  /**
   * If a query is made using a page number, then the limit property is used as 'page limit'
   * TODO: Improve pagination by adding a required eventId for orderBy clause
   */
  if (page !== undefined && limit !== undefined) {
    /**
     * We make sure that the page number is always >= 1
     */
    const currentPage: number = Math.max(1, page);
    const offset: number = (currentPage - 1) * limit;

    /**
     * Ensure sorting is applied to maintain consistent pagination results.
     * Also a casting of the ts type is required since the infer of the type
     * obtained from the count is not performed.
     */
    const count: { count?: string } = await baseQuery.clone().clearOrder().count({ count: '*' }).first() as unknown as { count?: string };

    baseQuery = baseQuery.offset(offset).limit(limit);

    return {
      results: await baseQuery.returning('*'),
      limit,
      offset,
      total: parseInt(count.count ?? '0', 10),
    };
  }

  return {
    results: await baseQuery.returning('*'),
  };
}

/**
 * Finds all transfers to or from child subaccounts of a parent subaccount,
 * excluding transfers between child subaccounts of the same parent.
 *
 * @param subaccountId - Array of all child subaccount IDs for the parent
 * @param limit - Maximum number of results to return
 * @param createdBeforeOrAtHeight - Filter transfers created at or before this height
 * @param createdBeforeOrAt - Filter transfers created at or before this time
 * @param page - Page number for pagination
 *
 * @returns Paginated list of transfers with same-parent transfers filtered out
 */
export async function findAllToOrFromParentSubaccount(
  {
    subaccountId,
    limit,
    createdBeforeOrAtHeight,
    createdBeforeOrAt,
    page,
  }: ParentSubaccountTransferQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<PaginationFromDatabase<TransferFromDatabase>> {
  verifyAllRequiredFields(
    {
      [QueryableField.LIMIT]: limit,
      [QueryableField.SUBACCOUNT_ID]: subaccountId,
      [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]: createdBeforeOrAtHeight,
      [QueryableField.CREATED_BEFORE_OR_AT]: createdBeforeOrAt,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<TransferModel> = setupBaseQuery<TransferModel>(
    TransferModel,
    options,
  );

  // Join with subaccounts to filter same-parent transfers
  baseQuery = baseQuery
    .leftJoin(
      'subaccounts as sender_sa',
      'transfers.senderSubaccountId',
      'sender_sa.id',
    )
    .leftJoin(
      'subaccounts as recipient_sa',
      'transfers.recipientSubaccountId',
      'recipient_sa.id',
    )
    // Exclude transfers where both sender and recipient are child subaccounts of the same parent
    .whereRaw(`
      NOT (
        transfers."senderSubaccountId" IS NOT NULL 
        AND transfers."recipientSubaccountId" IS NOT NULL
        AND sender_sa.address = recipient_sa.address 
        AND (sender_sa."subaccountNumber" % 128) = (recipient_sa."subaccountNumber" % 128)
      )
    `)
    .select('transfers.*');

  // Filter by child subaccount IDs
  baseQuery = baseQuery.where((queryBuilder) => {
    // eslint-disable-next-line no-void
    void queryBuilder.whereIn(TransferColumns.recipientSubaccountId, subaccountId)
      .orWhereIn(TransferColumns.senderSubaccountId, subaccountId);
  });

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

  if (options.orderBy !== undefined) {
    for (const [column, order] of options.orderBy) {
      baseQuery = baseQuery.orderBy(
        column,
        order,
      );
    }
  }

  if (limit !== undefined && page === undefined) {
    baseQuery = baseQuery.limit(limit);
  }

  // Pagination
  if (page !== undefined && limit !== undefined) {
    const currentPage: number = Math.max(1, page);
    const offset: number = (currentPage - 1) * limit;

    const count: { count?: string } = await baseQuery
      .clone()
      .clearSelect()
      .clearOrder()
      .count('transfers.id as count')
      .first() as unknown as { count?: string };

    baseQuery = baseQuery.offset(offset).limit(limit);

    return {
      results: await baseQuery.returning('*'),
      limit,
      offset,
      total: parseInt(count.count ?? '0', 10),
    };
  }

  return {
    results: await baseQuery.returning('*'),
  };
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
  [assetId: string]: Big,
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
      assetId: string,
      totalSize: string,
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

export async function getNetTransfersBetweenSubaccountIds(
  sourceSubaccountId: string,
  recipientSubaccountId: string,
  assetId: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<string> {
  const queryString: string = `
  SELECT 
    COALESCE(SUM(sub."size"), '0') AS "totalSize"
  FROM (
    SELECT DISTINCT 
      "size" AS "size",
      "id"
    FROM 
      "transfers"
    WHERE "transfers"."assetId" = '${assetId}'
    AND "transfers"."senderSubaccountId" = '${sourceSubaccountId}'
    AND "transfers"."recipientSubaccountId" = '${recipientSubaccountId}'
    UNION 
    SELECT DISTINCT 
      -"size" AS "size",
      "id"
    FROM 
      "transfers"
    WHERE "transfers"."assetId" = '${assetId}'
    AND "transfers"."senderSubaccountId" = '${recipientSubaccountId}'
    AND "transfers"."recipientSubaccountId" = '${sourceSubaccountId}'
  ) AS sub
  `;

  const result: {
    rows: { totalSize: string }[],
  } = await rawQuery(queryString, options);

  // Should only ever return a single row
  return result.rows[0].totalSize;
}

export async function create(
  transferToCreate: TransferCreateObject,
  options: Options = { txId: undefined },
): Promise<TransferFromDatabase> {
  return TransferModel.query(
    Transaction.get(options.txId),
  ).insert({
    ...transferToCreate,
    id: uuid(
      transferToCreate.eventId,
      transferToCreate.assetId,
      transferToCreate.senderSubaccountId,
      transferToCreate.recipientSubaccountId,
      transferToCreate.senderWalletAddress,
      transferToCreate.recipientWalletAddress,
    ),
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

export async function getLastTransferTimeForSubaccounts(
  subaccountIds: string[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<{ [subaccountId: string]: string }> {
  if (!subaccountIds.length) {
    return {};
  }

  let baseQuery: QueryBuilder<TransferModel> = setupBaseQuery<TransferModel>(
    TransferModel,
    options,
  );

  baseQuery = baseQuery
    .select('senderSubaccountId', 'recipientSubaccountId', 'createdAt')
    .where((queryBuilder) => {
      // eslint-disable-next-line no-void
      void queryBuilder.whereIn('senderSubaccountId', subaccountIds)
        .orWhereIn('recipientSubaccountId', subaccountIds);
    })
    .orderBy('createdAt', 'desc');

  const result: TransferFromDatabase[] = await baseQuery;

  const mapping: { [subaccountId: string]: string } = {};

  result.forEach((row) => {
    if (
      row.senderSubaccountId !== undefined &&
      subaccountIds.includes(row.senderSubaccountId)
    ) {
      if (!mapping[row.senderSubaccountId] || row.createdAt > mapping[row.senderSubaccountId]) {
        mapping[row.senderSubaccountId] = row.createdAt;
      }
    }
    if (
      row.recipientSubaccountId !== undefined &&
      subaccountIds.includes(row.recipientSubaccountId)
    ) {
      if (
        !mapping[row.recipientSubaccountId] ||
        row.createdAt > mapping[row.recipientSubaccountId]) {
        mapping[row.recipientSubaccountId] = row.createdAt;
      }
    }
  });

  return mapping;
}
