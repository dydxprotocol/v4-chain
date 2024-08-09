import { LeaderboardPnlTimeSpan } from '@dydxprotocol-indexer/postgres';
import { RedisClient } from 'redis';

import { getAsync } from '../helpers/redis';

export const LEADERBOARD_PNL_TIMESPAN_PROCESSED_CACHE_KEY: string = 'v4/leaderboard_pnl_processed/';

function getKey(period: LeaderboardPnlTimeSpan): string {
  return `${LEADERBOARD_PNL_TIMESPAN_PROCESSED_CACHE_KEY}${period}`;
}

export async function getProcessedTime(
  timespan: LeaderboardPnlTimeSpan,
  client: RedisClient,
): Promise<string | null> {
  return getAsync(getKey(timespan), client);
}

export async function setProcessedTime(
  period: LeaderboardPnlTimeSpan,
  timestamp: string,
  client: RedisClient,
): Promise<void> {
  await client.set(getKey(period), timestamp);
}
