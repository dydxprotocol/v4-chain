import { RedisClient } from 'redis';

import { getAsync } from '../helpers/redis';
import {
  CachedMegavaultPnl,
  CachedVaultHistoricalPnl,
  RedisVaultsArray,
  RedisMegavaultPnl,
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
 * Serialize a pnl tick object into a Redi-storage-sufficent format.
 * Rounding is used to reduce the number of decimal places stored.
 * We are fine with losing precision, because the data is only used for display purposes.
 *
 * @param tick - The PNL tick data to serialize
 * @returns [equity, totalPnl, netTransfers, createdAt, blockHeight, blockTime]
 */
function serializePnlTick(tick: CachedPnlTicks): [number, number, number, number, number, number] {
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
function deserializePnlTick(
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
  const serializedVaultsArray = JSON.parse(value);
  return serializedVaultsArray.map((compressed: RedisVaultsArray) => ({
    ticker: compressed[0],
    historicalPnl: compressed[1].map(deserializePnlTick),
  }));
}

export async function setVaultsHistoricalPnl(
  resolution: string,
  vaultsPnl: CachedVaultHistoricalPnl[],
  client: RedisClient,
): Promise<void> {
  const now = new Date().toISOString();

  // Create array of compressed tuples directly without intermediate string conversion
  const serializedVaultsArray = vaultsPnl.map((vault): RedisVaultsArray => [
    vault.ticker,
    vault.historicalPnl.map(serializePnlTick),
  ]);

  await Promise.all([
    client.set(getVaultsHistoricalPnlKey(resolution), JSON.stringify(serializedVaultsArray)),
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
  const compressed: RedisMegavaultPnl = JSON.parse(value);
  return {
    pnlTicks: compressed.map(deserializePnlTick),
  };
}

export async function setMegavaultPnl(
  resolution: string,
  pnlTicks: CachedPnlTicks[],
  client: RedisClient,
): Promise<void> {
  const now = new Date().toISOString();

  // Create compressed array directly
  const compressed: RedisMegavaultPnl = pnlTicks.map(serializePnlTick);

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
