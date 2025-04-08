import { PnlTicksFromDatabase } from '@dydxprotocol-indexer/postgres';
import { RedisClient } from 'redis';

import { getAsync } from '../helpers/redis';
import { CachedMegavaultPnl, CachedVaultHistoricalPnl } from '../types';

const KEY_PREFIX: string = 'v4/vault';

function getMegavaultHistoricalPnlKey(resolution: string): string {
  return `${KEY_PREFIX}/mv_pnl/${resolution}`;
}

function getMegavaultHistoricalTimestampKey(resolution: string): string {
  return `${KEY_PREFIX}/mv_pnl/timestamp/${resolution}`;
}

function getVaultsHistoricalPnlKey(resolution: string): string {
  return `${KEY_PREFIX}/vaults_pnl/${resolution}`;
}

function getVaultsHistoricalPnlTimestampKey(resolution: string): string {
  return `${KEY_PREFIX}/vaults_pnl/timestamp/${resolution}`;
}

/**
* Cache for /vaults/historicalPnl endpoint
**/

export async function getVaultsHistoricalPnl(
  resolution: string,
  client: RedisClient,
): Promise<CachedVaultHistoricalPnl[] | null> {
  const value: string | null = await getAsync(getVaultsHistoricalPnlKey(resolution), client);
  if (value === null) {
    return null;
  }
  return JSON.parse(value);
}

export async function setVaultsHistoricalPnl(
  resolution: string,
  vaultsPnl: CachedVaultHistoricalPnl[],
  client: RedisClient,
): Promise<void> {
  const now = new Date().toISOString();
  await Promise.all([
    client.set(getVaultsHistoricalPnlKey(resolution), JSON.stringify(vaultsPnl)),
    client.set(getVaultsHistoricalPnlTimestampKey(resolution), now),
  ]);
}

export async function getVaultsHistoricalPnlCacheTimestamp(
  resolution: string,
  client: RedisClient,
): Promise<Date | null> {
  const timestamp = await getAsync(getVaultsHistoricalPnlTimestampKey(resolution), client);
  return timestamp ? new Date(timestamp) : null;
}

/**
* Cache for /megavault/historicalPnl endpoint
**/

export async function getMegavaultPnl(
  resolution: string,
  client: RedisClient,
): Promise<CachedMegavaultPnl | null> {
  const value: string | null = await getAsync(getMegavaultHistoricalPnlKey(resolution), client);
  if (value === null) {
    return null;
  }
  return JSON.parse(value);
}

export async function setMegavaultPnl(
  resolution: string,
  pnlTicks: PnlTicksFromDatabase[],
  client: RedisClient,
): Promise<void> {
  const now = new Date().toISOString();
  const cache: CachedMegavaultPnl = {
    pnlTicks,
    lastUpdated: now,
  };
  await Promise.all([
    client.set(getMegavaultHistoricalPnlKey(resolution), JSON.stringify(cache)),
    client.set(getMegavaultHistoricalTimestampKey(resolution), now),
  ]);
}

export async function getMegavaultPnlCacheTimestamp(
  resolution: string,
  client: RedisClient,
): Promise<Date | null> {
  const timestamp = await getAsync(getMegavaultHistoricalTimestampKey(resolution), client);
  return timestamp ? new Date(timestamp) : null;
}
