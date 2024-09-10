import { PartialModelObject, QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
import { knexReadReplica } from '../helpers/knex';
import { setupBaseQuery, verifyAllRequiredFields } from '../helpers/stores-helpers';
import Transaction from '../helpers/transaction';
import WalletModel from '../models/wallet-model';
import {
  Options,
  Ordering,
  QueryableField,
  QueryConfig,
  WalletColumns,
  WalletCreateObject,
  WalletFromDatabase,
  WalletQueryConfig,
  WalletUpdateObject,
  PersistentCacheKeys,
} from '../types';

export async function findAll(
  {
    address,
    isWhitelistAffiliate,
    limit,
  }: WalletQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<WalletFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
      isWhitelistAffiliate,
      limit,
    } as QueryConfig,
    requiredFields,
  );

  let baseQuery: QueryBuilder<WalletModel> = setupBaseQuery<WalletModel>(
    WalletModel,
    options,
  );

  if (address) {
    baseQuery = baseQuery.where(WalletColumns.address, address);
  }

  if (isWhitelistAffiliate !== undefined) {
    baseQuery = baseQuery.where(WalletColumns.isWhitelistAffiliate, isWhitelistAffiliate);
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
      WalletColumns.address,
      Ordering.ASC,
    );
  }

  if (limit) {
    baseQuery = baseQuery.limit(limit);
  }

  return baseQuery.returning('*');
}

export async function create(
  walletToCreate: WalletCreateObject,
  options: Options = { txId: undefined },
): Promise<WalletFromDatabase> {
  return WalletModel.query(
    Transaction.get(options.txId),
  ).insert(walletToCreate).returning('*');
}

export async function update(
  {
    address,
    ...fields
  }: WalletUpdateObject,
  options: Options = { txId: undefined },
): Promise<WalletFromDatabase | undefined> {
  const wallet = await WalletModel.query(
    Transaction.get(options.txId),
  ).findById(address);
  const updatedWallet = await wallet.$query().patch(fields as PartialModelObject<WalletModel>).returning('*');
  // The objection types mistakenly think the query returns an array of Wallets.
  return updatedWallet as unknown as (WalletFromDatabase | undefined);
}

export async function upsert(
  walletToUpsert: WalletCreateObject,
  options: Options = { txId: undefined },
): Promise<WalletFromDatabase> {
  const wallets: WalletModel[] = await WalletModel.query(
    Transaction.get(options.txId),
  ).upsert(walletToUpsert).returning('*');
  // should only ever be one wallet
  return wallets[0];
}

export async function findById(
  address: string,
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<WalletFromDatabase | undefined> {
  const baseQuery: QueryBuilder<WalletModel> = setupBaseQuery<WalletModel>(
    WalletModel,
    options,
  );
  return baseQuery
    .findById(address)
    .returning('*');
}

/**
 * Calculates the total volume in a given time window for each address and adds the values to the
 * existing totalVolume values.
 *
 * @param windowStartTs - The start timestamp of the time window (exclusive).
 * @param windowEndTs - The end timestamp of the time window (inclusive).
 */
export async function updateTotalVolume(
  windowStartTs: string,
  windowEndTs: string,
) : Promise<void> {

  await knexReadReplica.getConnection().raw(
    `
    BEGIN;

    WITH fills_total AS (
      -- Step 1: Calculate total volume for each subaccountId
      SELECT "subaccountId", SUM("price" * "size") AS "totalVolume"
      FROM fills
      WHERE "createdAt" > '${windowStartTs}' AND "createdAt" <= '${windowEndTs}'
      GROUP BY "subaccountId"
    ),
    subaccount_volume AS (
      -- Step 2: Merge with subaccounts table to get the address
      SELECT s."address", f."totalVolume"
      FROM fills_total f
      JOIN subaccounts s
      ON f."subaccountId" = s."id"
    ),
    address_volume AS (
      -- Step 3: Group by address and sum the totalVolume
      SELECT "address", SUM("totalVolume") AS "totalVolume"
      FROM subaccount_volume
      GROUP BY "address"
    )
    -- Step 4: Left join the result with the wallets table and update the total volume
    UPDATE wallets
    SET "totalVolume" = COALESCE(wallets."totalVolume", 0) + av."totalVolume"
    FROM address_volume av
    WHERE wallets."address" = av."address";

    -- Step 5: Upsert new totalVolumeUpdateTime to persistent_cache table
    INSERT INTO persistent_cache (key, value)
    VALUES ('${PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME}', '${windowEndTs}')
    ON CONFLICT (key) 
    DO UPDATE SET value = EXCLUDED.value;

    COMMIT;
    `,
  );
}
