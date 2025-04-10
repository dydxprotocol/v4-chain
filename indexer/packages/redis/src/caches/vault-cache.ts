import { RedisClient } from 'redis';

import { getAsync } from '../helpers/redis';
import {
  CachedMegavaultPnl,
  CachedVaultHistoricalPnl,
  CompressedVaultPnl,
  CachedPnlTicks,
} from '../types';

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
 * Compresses a CachedVaultHistoricalPnl object into a more storage-efficient format.
 * Reduces size by:
 * 1. Using arrays instead of objects with named fields
 * 2. Converting ISO dates to Unix timestamps
 * 3. Limiting decimal precision to 1 place
 *
 * @param data - The vault historical PNL data to compress
 * @returns Compressed JSON string
 */
export function compressVaultPnl(data: CachedVaultHistoricalPnl): string {
  const compressed: CompressedVaultPnl = [
    data.ticker,
    data.historicalPnl.map((tick) => [
      Number(tick.equity).toFixed(1),
      Number(tick.totalPnl).toFixed(1),
      Number(tick.netTransfers).toFixed(1),
      Math.floor(new Date(tick.createdAt).getTime() / 1000),
      Number(tick.blockHeight),
      Math.floor(new Date(tick.blockTime).getTime() / 1000),
    ]),
  ];
  return JSON.stringify(compressed);
}

/**
 * Decompresses a string created by compressVaultPnl back into a CachedVaultHistoricalPnl object.
 *
 * @param compressedData - The compressed JSON string
 * @returns The decompressed vault historical PNL data
 */
export function decompressVaultPnl(compressedData: string): CachedVaultHistoricalPnl {
  const [ticker, historicalData]: CompressedVaultPnl = JSON.parse(compressedData);
  return {
    ticker,
    historicalPnl: historicalData.map(([e, p, n, c, h, t]) => ({
      equity: e,
      totalPnl: p,
      netTransfers: n,
      createdAt: new Date(c * 1000).toISOString(),
      blockHeight: h.toString(),
      blockTime: new Date(t * 1000).toISOString(),
    })),
  };
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

  // Parse as an array of compressed vault PNL data
  const compressedVaults = JSON.parse(value);
  return compressedVaults.map((compressed: string) => decompressVaultPnl(compressed));
}

export async function setVaultsHistoricalPnl(
  resolution: string,
  vaultsPnl: CachedVaultHistoricalPnl[],
  client: RedisClient,
): Promise<void> {
  const now = new Date().toISOString();

  // Compress each vault's PNL data
  const compressedVaults = vaultsPnl.map((vault) => compressVaultPnl(vault));

  await Promise.all([
    client.set(getVaultsHistoricalPnlKey(resolution), JSON.stringify(compressedVaults)),
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
  pnlTicks: CachedPnlTicks[],
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
