import {
  ONE_DAY_IN_MILLISECONDS,
  ONE_MINUTE_IN_MILLISECONDS,
  floorDate,
  logger,
} from '@dydxprotocol-indexer/base';
import {
  BlockFromDatabase,
  BlockTable,
  IsoString,
  TradingRewardAggregationFromDatabase,
  TradingRewardAggregationPeriod,
  TradingRewardAggregationTable,
  TradingRewardFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { AggregateTradingRewardsProcessedCache } from '@dydxprotocol-indexer/redis';
import { DateTime, Interval } from 'luxon';

import config from '../config';
import { redisClient } from '../helpers/redis';
import { UTC_OPTIONS } from '../lib/constants';

/**
 * Task: Aggregate Trading Rewards
 * Description: This task aggregates trading rewards for a specific period of time.
 * It retrieves trading data from the database, calculates the rewards, and stores the aggregated
 * results.
 */
interface SortedTradingRewardData {
  [address: string]: TradingRewardFromDatabase[];
}

export default function generateTaskFromPeriod(
  period: TradingRewardAggregationPeriod,
): () => Promise<void> {
  return async () => {
    const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(period);
    await aggregateTradingReward.runTask();
  };
}

export class AggregateTradingReward {
  period: TradingRewardAggregationPeriod;

  constructor(period: TradingRewardAggregationPeriod) {
    this.period = period;
  }

  async runTask(): Promise<void> {
    await this.maybeDeleteIncompleteAggregatedTradingReward();
    const interval: Interval = await this.getTradingRewardDataToProcessInterval();
    logger.info({
      at: 'aggregate-trading-rewards#runTask',
      message: 'Generated interval to aggregate trading rewards',
      start: interval.start.toISO(),
      end: interval.end.toISO(),
    });

    const tradingRewardData:
    TradingRewardFromDatabase[] = await this.getTradingRewardDataToProcess(interval);
    const sortedTradingRewardData: SortedTradingRewardData = this.sortTradingRewardData(
      tradingRewardData,
    );
    await this.updateTradingRewardsAggregation(sortedTradingRewardData);
    await AggregateTradingRewardsProcessedCache.setProcessedTime(
      this.period,
      interval.end.toISO(),
      redisClient,
    );
  }

  /**
   * If the latest processed time is null (should only happen during a fast sync),
   * and the latest period of aggregated trading rewards is incomplete.
   */
  async maybeDeleteIncompleteAggregatedTradingReward(): Promise<void> {
    const processedTime:
    IsoString | null = await AggregateTradingRewardsProcessedCache.getProcessedTime(
      this.period,
      redisClient,
    );
    if (processedTime !== null) {
      return;
    }
    const latestAggregation:
    TradingRewardAggregationFromDatabase | undefined = await
    TradingRewardAggregationTable.getLatestAggregatedTradeReward(this.period);

    // endedAt is only set when the entire interval has been processed for an aggregation
    if (latestAggregation !== undefined && latestAggregation.endedAt === null) {
      await this.deleteIncompleteAggregatedTradingReward(latestAggregation);
    }
  }

  /**
   * Deletes the latest this.period of aggregated trading rewards if it is incomplete. This is
   * called when the processedTime is null, and the latest aggregated trading rewards is incomplete.
   * We delete the latest this.period of aggregated trading rewards data because we don't know how
   * much data was processed within the interval, and we don't want to double count rewards.
   */
  private async deleteIncompleteAggregatedTradingReward(
    latestAggregation: TradingRewardAggregationFromDatabase,
  ): Promise<void> {
    logger.info({
      at: 'aggregate-trading-rewards#deleteIncompleteAggregatedTradingReward',
      message: `Deleting the latest ${this.period} aggregated trading rewards.`,
    });
    await TradingRewardAggregationTable.deleteAll({
      period: this.period,
      startedAtHeightOrAfter: latestAggregation.startedAtHeight,
    });
    logger.info({
      at: 'aggregate-trading-rewards#deleteIncompleteAggregatedTradingReward',
      message: `Deleted the last ${this.period} aggregated trading rewards`,
      height: latestAggregation.startedAtHeight,
      time: latestAggregation.startedAt,
    });
  }

  /**
   * Gets the interval of time to aggregate trading rewards for.
   * If there are no blocks in the database, then throw an error.
   * If There is no processedTime in the cache, then delete the latest month of data,
   * and reprocess that data or start from block 1.
   * If the processedTime is not null, and blocks exist in the database, then process up to the
   * next config.AGGREGATE_TRADING_REWARDS_MAX_INTERVAL_SIZE_MS of data.
   */
  async getTradingRewardDataToProcessInterval(): Promise<Interval> {
    const processedTime:
    IsoString | null = await AggregateTradingRewardsProcessedCache.getProcessedTime(
      this.period,
      redisClient,
    );
    const latestBlock: BlockFromDatabase = await BlockTable.getLatest();
    if (processedTime === null) {
      logger.info({
        at: 'aggregate-trading-rewards#getTradingRewardDataToProcessInterval',
        message: 'Resetting AggregateTradingRewardsProcessedCache',
      });
      const nextStartTime: DateTime = await this.getNextIntervalStartWhenCacheEmpty();
      await AggregateTradingRewardsProcessedCache.setProcessedTime(
        this.period,
        nextStartTime.toISO(),
        redisClient,
      );

      return this.generateInterval(nextStartTime, latestBlock);
    }

    const startTime: DateTime = DateTime.fromISO(processedTime, UTC_OPTIONS);
    return this.generateInterval(startTime, latestBlock);
  }

  /**
   * Returns the start time of the next interval to process if the
   * AggregateTradingRewardProcessedCache is empty. If there is a most recent complete aggregation
   * for this period, returns the end time of the most recent aggregation, otherwise returns the
   * start time of the first block in the database.
   */
  private async getNextIntervalStartWhenCacheEmpty(): Promise<DateTime> {
    const latestAggregation:
    TradingRewardAggregationFromDatabase | undefined = await
    TradingRewardAggregationTable.getLatestAggregatedTradeReward(this.period);
    // Since we've deleted the incomplete aggregations, we can assume that the latestAggregation
    // is complete.
    if (latestAggregation !== undefined) {
      return DateTime.fromISO(latestAggregation.endedAt!, UTC_OPTIONS);
    }

    // Since we were able to find the latest block, we assume we can find the first block
    const firstBlock: BlockFromDatabase[] = await BlockTable.findAll({
      blockHeight: ['1'],
      limit: 1,
    }, []);
    return DateTime.fromISO(firstBlock[0].time, UTC_OPTIONS);
  }

  /**
   * Generate the interval that will be processed. The end time of the interval is calculated from
   * a start time and the latest block. This will be the earliest of the following:
   * 1. The next day
   * 2. Start time plus config.AGGREGATE_TRADING_REWARDS_MAX_INTERVAL_SIZE_MS
   * 3. The start of the minute of the latest block
   * @param startTime - startTime of the interval
   * @param latestBlock
   * @returns
   */
  private generateInterval(
    startTime: DateTime,
    latestBlock: BlockFromDatabase,
  ): Interval {
    const latestBlockTime: Date = DateTime.fromISO(latestBlock.time, UTC_OPTIONS).toJSDate();

    // The most recent start of a minute. i.e 12:02:33 will be rounded to 12:02:00
    const normalizedLatestBlockTime: Date = floorDate(
      latestBlockTime,
      ONE_MINUTE_IN_MILLISECONDS,
    );

    const nextDay: Date = startTime.plus({ days: 1 }).toJSDate();
    const normalizedNextDay: Date = floorDate(nextDay, ONE_DAY_IN_MILLISECONDS);

    const startDate: Date = startTime.toJSDate();
    const startTimePlusMaxIntervalSize: Date = DateTime.fromJSDate(startDate).plus(
      { milliseconds: config.AGGREGATE_TRADING_REWARDS_MAX_INTERVAL_SIZE_MS },
    ).toJSDate();
    const endTime: Date = new Date(Math.min(
      normalizedLatestBlockTime.getTime(),
      normalizedNextDay.getTime(),
      startTimePlusMaxIntervalSize.getTime(),
    ));
    const endDateTime: DateTime = DateTime.fromJSDate(endTime).toUTC();
    return Interval.fromDateTimes(startTime, endDateTime);
  }

  private async getTradingRewardDataToProcess(
    _interval: Interval,
  ): Promise<TradingRewardFromDatabase[]> {
    // TODO: Implement
    return Promise.resolve([]);
  }

  private sortTradingRewardData(
    _tradingRewardData: TradingRewardFromDatabase[],
  ): SortedTradingRewardData {
    // TODO: Implement
    return {};
  }

  private async updateTradingRewardsAggregation(
    _sortedTradingRewardData: SortedTradingRewardData,
  ): Promise<void> {
    // TODO: Implement
  }
}
