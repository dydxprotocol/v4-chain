import { logger } from '@dydxprotocol-indexer/base';
import {
  LeaderboardPnlTimeSpan,
  LeaderboardPnlCreateObject,
  LeaderboardPnlTable,
  PnlTicksTable,
  Transaction,
} from '@dydxprotocol-indexer/postgres';
import { LeaderboardPnlProcessedCache } from '@dydxprotocol-indexer/redis';
import _ from 'lodash';

import config from '../config';
import { redisClient } from '../helpers/redis';

export default function generateLeaderboardTaskFromTimespan(
  timespan: LeaderboardPnlTimeSpan,
): () => Promise<void> {
  return async () => {
    const leaderboardPnlProcessor: LeaderboardPnlProcessor = new LeaderboardPnlProcessor(timespan);
    await leaderboardPnlProcessor.runTask();

  };
}

class LeaderboardPnlProcessor {

  constructor(
    private timespan: LeaderboardPnlTimeSpan,
  ) {
    this.timespan = timespan;
  }

  async runTask(): Promise<void> {
    await this.updateLeaderboardPnlTable(this.timespan);
  }

  /**
   * Updates the leaderboard PnL table for a specified time span.
   * This function performs several operations to ensure the leaderboard PnL table is up-to-date:
   * 1. Check for updates: It first retrieves the last processed time for the given
   *    time span from a cache. It then fetches the latest processed block time from
   *    the PnL ticks table. If the last processed time is greater than or equal to
   *    the latest block time, it logs a message indicating that the current PnL ticks
   *    have already been processed for the leaderboard and exits the function.
   *
   * 2. Retrieve PnL Objects: If updates are needed, it retrieves
   *    leaderboard PnL objects for the specified time span by calling
   *    `getLeaderboardPnlObjects`.
   *
   * 3. Insert PnL Objects: It then inserts the leaderboard PnL objects into the
   *    leaderboard PnL table by calling `insertLeaderboardPnlObjects`.
   *
   * 4. Update Cache: After successful insertion, it updates the cache with the
   *    latest processed block time for the given time span.
   *
   */
  async updateLeaderboardPnlTable(timespan: LeaderboardPnlTimeSpan) {
    const lastProcessedTime: string | null = await LeaderboardPnlProcessedCache.getProcessedTime(
      timespan, redisClient);
    const lastProcessedPnlTime: string = (
      await PnlTicksTable.findLatestProcessedBlocktimeAndCount()).maxBlockTime;

    // Check if the last processed time is greater than or equal to the latest block time.
    // In cases where indexer is from previous state, the last processed time may be null.
    // In that case update the leaderboard.
    if (lastProcessedTime && Date.parse(lastProcessedTime) >= Date.parse(lastProcessedPnlTime)) {
      logger.info({
        at: 'create-leaderboard#runTask',
        message: 'Skipping run because the current pnl ticks have been processed for the leaderboard',
        pnlTickLatestBlocktime: lastProcessedPnlTime,
        latestBlockTime: lastProcessedTime,
        threshold: config.PNL_TICK_UPDATE_INTERVAL_MS,
        timespan,
      });
      return;
    }
    try {
      const leaderboardPnlObjects: LeaderboardPnlCreateObject[] = await
      this.getLeaderboardPnlObjects(timespan);
      await this.insertLeaderboardPnlObjects(leaderboardPnlObjects);
      await this.updateLeaderboardPnlProcessedCache(timespan, lastProcessedPnlTime);
    } catch (error) {
      logger.error({
        at: 'create-leaderboard#runTask',
        message: 'Error when updating leaderboard pnl table',
        error,
        timespan,
      });
    }
  }

  async getLeaderboardPnlObjects(timespan: LeaderboardPnlTimeSpan) {
    const leaderboardPnlObjects: LeaderboardPnlCreateObject[] = [];
    try {
      leaderboardPnlObjects.push(...await PnlTicksTable.getRankedPnlTicks(timespan));
    } catch (error) {
      logger.error({
        at: 'create-leaderboard#runTask',
        message: `Error when getting ranked pnl ticks for timespan${timespan.toString()}`,
        error,
        timespan,
      });
      throw error;
    }
    return leaderboardPnlObjects;
  }

  async insertLeaderboardPnlObjects(leaderboardPnlObjects: LeaderboardPnlCreateObject[],
  ) {
    const txId: number = await Transaction.start();
    try {
      const chunkedLeaderboardPnlObjects = _.chunk(
        leaderboardPnlObjects,
        config.LEADERBOARD_PNL_MAX_ROWS_PER_UPSERT);
      for (const leaderboardPnlObjectsToUpsert of chunkedLeaderboardPnlObjects) {
        await LeaderboardPnlTable.bulkUpsert(leaderboardPnlObjectsToUpsert, { txId });
      }
      await Transaction.commit(txId);
    } catch (error) {
      logger.error({
        at: 'create-leaderboard#runTask',
        message: 'Error when inserting leaderboard pnl objects',
        error,
      });
      await Transaction.rollback(txId);
      throw error;
    }
  }

  async updateLeaderboardPnlProcessedCache(timespan: LeaderboardPnlTimeSpan,
    lastProcessedPnlTime: string) {
    try {
      await LeaderboardPnlProcessedCache.setProcessedTime(timespan, lastProcessedPnlTime,
        redisClient);
    } catch (error) {
      logger.error({
        at: 'create-leaderboard#runTask',
        message: 'Error when setting processed time',
        error,
      });
    }
  }

}
