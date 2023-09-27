import { QueryBuilder } from 'objection';

import { DEFAULT_POSTGRES_OPTIONS } from '../constants';
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
} from '../types';

export async function findAll(
  {
    address,
    limit,
  }: WalletQueryConfig,
  requiredFields: QueryableField[],
  options: Options = DEFAULT_POSTGRES_OPTIONS,
): Promise<WalletFromDatabase[]> {
  verifyAllRequiredFields(
    {
      address,
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
