import { promisify } from 'util';

import { SubaccountTable } from '@dydxprotocol-indexer/postgres';
import { IndexerSubaccountId } from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';
import { RedisClient } from 'redis';

import { hGetAllAsync } from '../helpers/redis';

// Cache of subaccount uuid to list of order uuids for the subaccount
export const SUBACCOUNT_ORDERS_KEY_PREFIX: string = 'v4/subaccountOrderIds/';

/**
 * Get order ids for a subaccount.
 * @param subaccountUuid Indexer assigned UUID of the subaccount to get order ids for.
 * @param client Redis client.
 * @returns List of indexer assigned order uuids that belong to the subaccount.
 */
export async function getOrderIdsForSubaccount(
  subaccountUuid: string,
  client: RedisClient,
): Promise<string[]> {
  return _.keys(await hGetAllAsync(getSubaccountOrderIdsCacheKeyWithUUID(subaccountUuid), client));
}

/**
 * Get order ids for a list of subaccounts.
 * @param subaccountUuids List of indexer subaccount ids to get order ids for.
 * @param client Redis client.
 * @returns Map of indexer subaccount id to list of indexer order ids that belong to the subaccount.
 */
export async function getOrderIdsForSubaccounts(
  subaccountUuids: string[],
  client: RedisClient,
): Promise<Record<string, string[]>> {
  // Pipeline all hgetalls.
  const multi = client.multi();
  for (const uuid of subaccountUuids) {
    const key = getSubaccountOrderIdsCacheKeyWithUUID(uuid);
    multi.hgetall(key);
  }
  const execAsync = promisify(multi.exec).bind(multi);
  const allOrderIds = await execAsync();

  const subaccountUuidToOrderIds: Record<string, string[]> = {};
  subaccountUuids.forEach((uuid, i) => {
    subaccountUuidToOrderIds[uuid] = allOrderIds[i] ? Object.keys(allOrderIds[i]) : [];
  });

  return subaccountUuidToOrderIds;
}

export function getSubaccountOrderIdsCacheKey(subaccountId: IndexerSubaccountId): string {
  return getSubaccountOrderIdsCacheKeyWithUUID(SubaccountTable.subaccountIdToUuid(subaccountId));
}

export function getSubaccountOrderIdsCacheKeyWithUUID(subaccountUuid: string): string {
  return `${SUBACCOUNT_ORDERS_KEY_PREFIX}${subaccountUuid}`;
}
