import { RedisClient } from 'redis';

import { getAsync } from '../helpers/redis';
import {
  CachedMegavaultPnl,
  CachedVaultHistoricalPnl,
  CompressedVaultPnl,
  CompressedMegavaultPnl,
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
 * Common function to compress a CachedPnlTicks object into a more storage-efficient format.
 *
 * @param tick - The PNL tick data to compress
 * @returns [equity, totalPnl, netTransfers, createdAt, blockHeight, blockTime]
 */
function compressPnlTick(tick: CachedPnlTicks): [number, number, number, number, number, number] {
  return [
    Math.round(Number(tick.equity)),
    Math.round(Number(tick.totalPnl)),
    Math.round(Number(tick.netTransfers)),
    Math.floor(new Date(tick.createdAt).getTime() / 1000),
    Number(tick.blockHeight),
    Math.floor(new Date(tick.blockTime).getTime() / 1000),
  ];
}

/**
 * Common function to decompress a PNL tick array back into a CachedPnlTicks object.
 *
 * @param data - [equity, totalPnl, netTransfers, createdAt, blockHeight, blockTime]
 * @returns The decompressed CachedPnlTicks object
 */
function decompressPnlTick(
  [e, p, n, c, h, t]: [number, number, number, number, number, number],
): CachedPnlTicks {
  return {
    equity: e.toString(),
    totalPnl: p.toString(),
    netTransfers: n.toString(),
    createdAt: new Date(c * 1000).toISOString(),
    blockHeight: h.toString(),
    blockTime: new Date(t * 1000).toISOString(),
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

  // Parse the compressed vault PNL data array directly
  const compressedVaultsArray = JSON.parse(value);
  return compressedVaultsArray.map((compressed: CompressedVaultPnl) => ({
    ticker: compressed[0],
    historicalPnl: compressed[1].map(decompressPnlTick),
  }));
}

export async function setVaultsHistoricalPnl(
  resolution: string,
  vaultsPnl: CachedVaultHistoricalPnl[],
  client: RedisClient,
): Promise<void> {
  const now = new Date().toISOString();

  // Create array of compressed tuples directly without intermediate string conversion
  const compressedVaultsArray = vaultsPnl.map((vault): CompressedVaultPnl => [
    vault.ticker,
    vault.historicalPnl.map(compressPnlTick),
  ]);

  await Promise.all([
    client.set(getVaultsHistoricalPnlKey(resolution), JSON.stringify(compressedVaultsArray)),
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

  // Parse the compressed megavault PNL data directly
  const compressed: CompressedMegavaultPnl = JSON.parse(value);
  return {
    pnlTicks: compressed.map(decompressPnlTick),
  };
}

export async function setMegavaultPnl(
  resolution: string,
  pnlTicks: CachedPnlTicks[],
  client: RedisClient,
): Promise<void> {
  const now = new Date().toISOString();

  // Create compressed array directly
  const compressed: CompressedMegavaultPnl = pnlTicks.map(compressPnlTick);

  await Promise.all([
    client.set(getMegavaultHistoricalPnlKey(resolution), JSON.stringify(compressed)),
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
