import {
  BlockFromDatabase,
  BlockTable,
  TradingRewardFromDatabase,
} from "@dydxprotocol-indexer/postgres";
import { DateTime } from "luxon";

/**
 * Task: Aggregate Trading Rewards
 * Description: This task aggregates trading rewards for a specific period of time.
 * It retrieves trading data from the database, calculates the rewards, and stores the aggregated
 * results.
 */
interface Interval {
  start: DateTime;
  end: DateTime;
}

interface SortedTradingRewardData {
  [address: string]: TradingRewardFromDatabase[];
}

export default async function runTask(): Promise<void> {
  // TODO(IND-499): Add resetting aggregation data when cache is empty
  const interval: Interval | undefined = await getTradingRewardDataToProcessInterval();

  const tradingRewardData: TradingRewardFromDatabase[] = await getTradingRewardDataToProcess(interval);
  const sortedTradingRewardData: SortedTradingRewardData = sortTradingRewardData(tradingRewardData);
  await updateTradingRewardsAggregation(sortedTradingRewardData);
  // TODO(IND-499): Update AggregateTradingRewardsProcessedCache
}

async function getTradingRewardDataToProcessInterval(): Promise<Interval> {
  const latestBlock: BlockFromDatabase = await BlockTable.getLatest();

  // TODO(IND-499): Setup AggregateTradingRewardsProcessedCache for start time and add end time
  return {
    start: DateTime.fromISO(latestBlock.time),
    end: DateTime.fromISO(latestBlock.time),
  };
}

async function getTradingRewardDataToProcess(interval: Interval): Promise<TradingRewardFromDatabase[]> {
  // TODO: Implement
  return [];
}

function sortTradingRewardData(tradingRewardData: TradingRewardFromDatabase[]): SortedTradingRewardData {
  // TODO: Implement
  return {};
}

async function updateTradingRewardsAggregation(sortedTradingRewardData: SortedTradingRewardData): Promise<void> {
  // TODO: Implement
}
