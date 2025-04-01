import { PnlTicksFromDatabase } from '@dydxprotocol-indexer/postgres';
import { RedisClient } from 'redis';

import { getAsync } from '../helpers/redis';

const KEY_PREFIX: string = 'v4/megavault-historical-pnl';

export interface CachedMegavaultPnl {
  pnlTicks: PnlTicksFromDatabase[],
  lastUpdated: string,
}

function getKey(resolution: string): string {
  return `${KEY_PREFIX}/${resolution}`;
}

export async function get(
  resolution: string,
  client: RedisClient,
): Promise<CachedMegavaultPnl | null> {
  const value: string | null = await getAsync(getKey(resolution), client);
  if (value === null) {
    return null;
  }
  return JSON.parse(value);
}

export async function set(
  resolution: string,
  pnlTicks: PnlTicksFromDatabase[],
  client: RedisClient,
): Promise<void> {
  const cache: CachedMegavaultPnl = {
    pnlTicks,
    lastUpdated: new Date().toISOString(),
  };
  await client.set(getKey(resolution), JSON.stringify(cache));
}
