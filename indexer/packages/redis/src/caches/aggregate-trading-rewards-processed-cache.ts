import { IsoString, TradingRewardAggregationPeriod } from '@dydxprotocol-indexer/postgres';
import { RedisClient } from 'redis';

import { getAsync } from '../helpers/redis';

export const AGGREGATE_TRADING_REWARDS_PROCESSED_CACHE_KEY: string = 'v4/aggregate_trading_rewards_processed/';

function getKey(period: TradingRewardAggregationPeriod): string {
  return `${AGGREGATE_TRADING_REWARDS_PROCESSED_CACHE_KEY}${period}`;
}

export async function getProcessedTime(
  period: TradingRewardAggregationPeriod,
  client: RedisClient,
): Promise<IsoString | null> {
  return getAsync(getKey(period), client);
}

export async function setProcessedTime(
  period: TradingRewardAggregationPeriod,
  timestamp: IsoString,
  client: RedisClient,
): Promise<void> {
  await client.set(getKey(period), timestamp);
}
