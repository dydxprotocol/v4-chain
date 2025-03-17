import {
  ONE_DAY_IN_MILLISECONDS,
  ONE_MINUTE_IN_MILLISECONDS,
  floorDate,
  logger,
  runFuncWithTimingStat,
} from '@dydxprotocol-indexer/base';
import {
  BlockFromDatabase,
  BlockTable,
  IsoString,
  IsolationLevel,
  Ordering,
  TradingRewardAggregationColumns,
  TradingRewardAggregationCreateObject,
  TradingRewardAggregationFromDatabase,
  TradingRewardAggregationPeriod,
  TradingRewardAggregationTable,
  TradingRewardAggregationUpdateObject,
  TradingRewardColumns,
  TradingRewardFromDatabase,
  TradingRewardTable,
  Transaction,
} from '@dydxprotocol-indexer/postgres';
import { AggregateTradingRewardsProcessedCache } from '@dydxprotocol-indexer/redis';
import Big from 'big.js';
import _ from 'lodash';
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
type IntervalTradingRewardsByAddress = _.Dictionary<string>;
interface AggregationUpdateAndCreateObjects {
  updateObjects: TradingRewardAggregationUpdateObject[],
  createObjects: TradingRewardAggregationCreateObject[],
}

enum DateTimeUnit {
  DAY = 'day',
  WEEK = 'week',
  MONTH = 'month',
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
      period: this.period,
      start: interval.start.toISO(),
      end: interval.end.toISO(),
    });

    const intervalTradingRewardsByAddress:
    IntervalTradingRewardsByAddress = await runFuncWithTimingStat(
      this.getIntervalTradingRewardsByAddress(
        interval,
      ),
      this.generateTimingStatsOptions('getIntervalTradingRewardsByAddress'),
    );
    await this.updateTradingRewardsAggregation(interval, intervalTradingRewardsByAddress);
    await this.setProcessedTime(
      interval.end.toISO(),
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
      await runFuncWithTimingStat(
        this.deleteIncompleteAggregatedTradingReward(latestAggregation),
        this.generateTimingStatsOptions('deleteIncompleteAggregatedTradingReward'),
      );
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
      period: this.period,
    });
    await TradingRewardAggregationTable.deleteAll({
      period: this.period,
      startedAtHeightOrAfter: latestAggregation.startedAtHeight,
    });
    logger.info({
      at: 'aggregate-trading-rewards#deleteIncompleteAggregatedTradingReward',
      message: `Deleted the last ${this.period} aggregated trading rewards`,
      period: this.period,
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
        message: 'AggregateTradingRewardsProcessedCache is empty',
        period: this.period,
      });
      const nextStartTime: DateTime = await this.getNextIntervalStartWhenCacheEmpty();
      await this.setProcessedTime(
        nextStartTime.toISO(),
      );

      return this.generateInterval(nextStartTime, latestBlock);
    }

    const startTime: DateTime = DateTime.fromISO(processedTime, UTC_OPTIONS);
    return this.generateInterval(startTime, latestBlock);
  }

  /**
   * Returns the start time of the next interval to process if the
   * AggregateTradingRewardProcessedCache is empty.
   * - If there is a most recent complete aggregation for this period,
   * returns the end time of the most recent aggregation.
   * - If there is a trading reward in the database,
   * returns the block time of the trading reward.
   * - Otherwise returns the start time of the first block in the database.
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

    const firstTradingReward: TradingRewardFromDatabase[] = await TradingRewardTable.findAll({
      limit: 1,
    }, [], { orderBy: [[TradingRewardColumns.blockTime, Ordering.ASC]] });
    if (firstTradingReward.length === 0) {
      logger.error({
        at: 'aggregate-trading-rewards#getNextIntervalStartWhenCacheEmpty',
        message: 'No trading rewards in database',
        period: this.period,
      });
      throw new Error('No trading rewards in database');
    }

    return DateTime.fromISO(firstTradingReward[0].blockTime, UTC_OPTIONS);
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

  private async getIntervalTradingRewardsByAddress(
    interval: Interval,
  ): Promise<IntervalTradingRewardsByAddress> {
    const tradingRewards: TradingRewardFromDatabase[] = await TradingRewardTable.findAll({
      blockTimeAfterOrAt: interval.start.toISO(),
      blockTimeBefore: interval.end.toISO(),
    }, []);

    const tradingRewardsByAddress: _.Dictionary<TradingRewardFromDatabase[]> = _.groupBy(
      tradingRewards,
      TradingRewardColumns.address,
    );

    return _.mapValues(
      tradingRewardsByAddress,
      (tradingRewardsForAddress: TradingRewardFromDatabase[]) => {
        return _.reduce(
          tradingRewardsForAddress,
          (sum: string, tradingReward: TradingRewardFromDatabase) => {
            return Big(sum).plus(tradingReward.amount).toFixed();
          },
          '0',
        );
      },
    );
  }

  private async updateTradingRewardsAggregation(
    interval: Interval,
    intervalTradingRewardsByAddress: IntervalTradingRewardsByAddress,
  ): Promise<void> {
    let aggregationUpdateAndCreateObjects:
    AggregationUpdateAndCreateObjects = await this.getAggregationUpdateAndCreateObjectsFromInterval(
      interval,
      intervalTradingRewardsByAddress,
    );

    // If interval.end is the end of this.period, then we need to set the endedAt and endedAtHeight
    // for all the aggregation objects.
    if (this.isEndofPeriod(interval.end)) {
      aggregationUpdateAndCreateObjects = await this.addEndedAtAndEndedAtHeightUpdates(
        aggregationUpdateAndCreateObjects,
        interval,
      );
    }

    const txId: number = await Transaction.start();
    await Transaction.setIsolationLevel(txId, IsolationLevel.READ_UNCOMMITTED);
    try {
      await this.createAndUpdateAggregations(aggregationUpdateAndCreateObjects, txId);
      await Transaction.commit(txId);
      logger.info({
        at: 'aggregate-trading-rewards#updateTradingRewardsAggregation',
        message: 'Updated trading rewards aggregation',
        period: this.period,
        start: interval.start.toISO(),
        end: interval.end.toISO(),
      });
    } catch (error) {
      await Transaction.rollback(txId);
      logger.info({
        at: 'aggregate-trading-rewards#updateTradingRewardsAggregation',
        message: 'Failed to update trading rewards aggregation',
        period: this.period,
        error: error.message,
        start: interval.start.toISO(),
        end: interval.end.toISO(),
      });
      throw error;
    }
  }

  /**
   * Generate all update and create objects for the interval, by aggregating all trading rewards
   * for the interval.
   */
  private async getAggregationUpdateAndCreateObjectsFromInterval(
    interval: Interval,
    intervalTradingRewardsByAddress: IntervalTradingRewardsByAddress,
  ): Promise<AggregationUpdateAndCreateObjects> {
    const tradingRewardAddresses: string[] = Object.keys(intervalTradingRewardsByAddress);

    const startedAt: string = this.getStartedAt(interval);
    const startedAtHeight: string = await this.getNextBlock(startedAt);
    const existingAggregateTradingRewards:
    TradingRewardAggregationFromDatabase[] = await runFuncWithTimingStat(
      TradingRewardAggregationTable.findAll({
        addresses: tradingRewardAddresses,
        period: this.period,
        startedAtHeight,
      }, []),
      this.generateTimingStatsOptions('findAllExistingAggregations'),
    );
    const existingAggregateTradingRewardsMap:
    { [address: string]: TradingRewardAggregationFromDatabase } = _.keyBy(
      existingAggregateTradingRewards,
      TradingRewardAggregationColumns.address,
    );

    const aggregateTradingRewardsToUpdate: TradingRewardAggregationUpdateObject[] = _.intersection(
      tradingRewardAddresses,
      Object.keys(existingAggregateTradingRewardsMap),
    ).map((address: string) => {
      return {
        id: TradingRewardAggregationTable.uuid(address, this.period, startedAtHeight),
        amount: Big(intervalTradingRewardsByAddress[address])
          .plus(existingAggregateTradingRewardsMap[address].amount).toFixed(),
      };
    });

    const aggregateTradingRewardsToCreate: TradingRewardAggregationCreateObject[] = _.difference(
      tradingRewardAddresses,
      Object.keys(existingAggregateTradingRewardsMap),
    ).map((address: string) => {
      return {
        address,
        startedAt,
        startedAtHeight,
        period: this.period,
        amount: intervalTradingRewardsByAddress[address],
      };
    });

    return {
      updateObjects: aggregateTradingRewardsToUpdate,
      createObjects: aggregateTradingRewardsToCreate,
    };
  }

  private getStartedAt(interval: Interval): IsoString {
    return interval.start.startOf(this.getDateTimeUnit()).toISO();
  }

  private async getNextBlock(time: IsoString): Promise<string> {
    const block: BlockFromDatabase | undefined = await BlockTable.findBlockByCreatedOnOrAfter(
      time,
      { readReplica: true },
    );

    if (block === undefined) {
      logger.error({
        at: 'aggregate-trading-rewards#getStartedAtHeight',
        message: 'No blocks found after time, this should never happen',
        period: this.period,
        time,
      });
      throw new Error(`No blocks found after ${time}`);
    }
    return block.blockHeight;
  }

  private isEndofPeriod(endTime: DateTime): boolean {
    return endTime.startOf(this.getDateTimeUnit()).toISO() === endTime.toISO();
  }

  private async addEndedAtAndEndedAtHeightUpdates(
    aggregationUpdateAndCreateObjects: AggregationUpdateAndCreateObjects,
    interval: Interval,
  ): Promise<AggregationUpdateAndCreateObjects> {
    const startedAt: string = this.getStartedAt(interval);
    const startedAtHeight: string = await this.getNextBlock(startedAt);
    const endedAt: IsoString = interval.end.toISO();
    // endedAtHeight is the first block created before endedAt
    const endedAtHeight: string = Big(await this.getNextBlock(endedAt)).minus(1).toFixed();

    const allIncompleteAggregation: TradingRewardAggregationFromDatabase[] = await
    TradingRewardAggregationTable.findAll({
      period: this.period,
      startedAtHeight,
    }, []);
    const allIncompleteAggregationAddresses: string[] = _.map(
      allIncompleteAggregation,
      TradingRewardAggregationColumns.address,
    );

    const otherAggregationAddresses: string[] = _.difference(
      allIncompleteAggregationAddresses,
      _.map(
        aggregationUpdateAndCreateObjects.updateObjects,
        TradingRewardAggregationColumns.address,
      ),
    );
    const otherAggregationUpdateObjects: TradingRewardAggregationUpdateObject[] = _.map(
      otherAggregationAddresses,
      (address: string) => {
        return {
          id: TradingRewardAggregationTable.uuid(address, this.period, startedAtHeight),
          endedAt,
          endedAtHeight,
        };
      },
    );
    return {
      createObjects: _.map(
        aggregationUpdateAndCreateObjects.createObjects,
        (createObject: TradingRewardAggregationCreateObject) => {
          return {
            ...createObject,
            endedAt,
            endedAtHeight,
          };
        },
      ),
      updateObjects: _.map(
        aggregationUpdateAndCreateObjects.updateObjects,
        (updateObject: TradingRewardAggregationUpdateObject) => {
          return {
            ...updateObject,
            endedAt,
            endedAtHeight,
          } as TradingRewardAggregationUpdateObject;
        },
      ).concat(otherAggregationUpdateObjects),
    };
  }

  private getDateTimeUnit(): DateTimeUnit {
    switch (this.period) {
      case TradingRewardAggregationPeriod.DAILY:
        return DateTimeUnit.DAY;
      case TradingRewardAggregationPeriod.WEEKLY:
        return DateTimeUnit.WEEK;
      case TradingRewardAggregationPeriod.MONTHLY:
        return DateTimeUnit.MONTH;
      default:
        throw new Error(`Invalid period ${this.period}`);
    }
  }

  private async createAndUpdateAggregations(
    aggregationUpdateAndCreateObjects: AggregationUpdateAndCreateObjects,
    txId: number,
  ): Promise<void> {
    const createObjectsChunks: TradingRewardAggregationCreateObject[][] = _.chunk(
      aggregationUpdateAndCreateObjects.createObjects,
      config.AGGREGATE_TRADING_REWARDS_CHUNK_SIZE,
    );
    for (const createObjectsChunk of createObjectsChunks) {
      logger.info({
        at: 'aggregate-trading-rewards#setAggregationUpdateAndCreateObjects',
        message: 'Creating trading reward aggregations',
        period: this.period,
        count: createObjectsChunk.length,
        createObjectsChunk: JSON.stringify(createObjectsChunk),
      });
      await runFuncWithTimingStat(
        Promise.all(_.map(
          createObjectsChunk,
          (createObject: TradingRewardAggregationCreateObject) => {
            return TradingRewardAggregationTable.create(createObject, { txId });
          },
        )),
        this.generateTimingStatsOptions('createChunk'),
      );
      logger.info({
        at: 'aggregate-trading-rewards#setAggregationUpdateAndCreateObjects',
        message: 'Created trading reward aggregations',
        period: this.period,
        count: createObjectsChunk.length,
      });
    }

    const updateObjectsChunks: TradingRewardAggregationUpdateObject[][] = _.chunk(
      aggregationUpdateAndCreateObjects.updateObjects,
      config.AGGREGATE_TRADING_REWARDS_CHUNK_SIZE,
    );
    for (const updateObjectsChunk of updateObjectsChunks) {
      logger.info({
        at: 'aggregate-trading-rewards#setAggregationUpdateAndCreateObjects',
        message: 'Updating trading reward aggregations',
        period: this.period,
        count: updateObjectsChunk.length,
        updateObjectsChunk: JSON.stringify(updateObjectsChunk),
      });
      await runFuncWithTimingStat(
        Promise.all(_.map(
          updateObjectsChunk,
          (updateObject: TradingRewardAggregationUpdateObject) => {
            return TradingRewardAggregationTable.update(updateObject, { txId });
          },
        )),
        this.generateTimingStatsOptions('updateChunk'),
      );
      logger.info({
        at: 'aggregate-trading-rewards#setAggregationUpdateAndCreateObjects',
        message: 'Updated trading reward aggregations',
        period: this.period,
        count: updateObjectsChunk.length,
      });
    }
  }

  private async setProcessedTime(processedTime: IsoString): Promise<void> {
    logger.info({
      at: 'aggregate-trading-rewards#setProcessedTime',
      message: 'Setting processed time',
      period: this.period,
      processedTime,
    });
    await AggregateTradingRewardsProcessedCache.setProcessedTime(
      this.period,
      processedTime,
      redisClient,
    );
    logger.info({
      at: 'aggregate-trading-rewards#setProcessedTime',
      message: 'Set processed time',
      period: this.period,
      processedTime,
    });
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  protected generateTimingStatsOptions(fnName: string): any {
    return {
      taskName: 'aggregate-trading-rewards',
      fnName,
    };
  }
}
