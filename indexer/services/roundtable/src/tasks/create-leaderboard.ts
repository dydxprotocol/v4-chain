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
    await updateLeaderboardPnlTable(timespan);
  };
}

async function updateLeaderboardPnlTable(timespan: LeaderboardPnlTimeSpan) {
  const leaderboardPnlObjects: LeaderboardPnlCreateObject[] = [];

  const lastProcessedTime: string | null = await LeaderboardPnlProcessedCache.getProcessedTime(
    timespan, redisClient);
  const lastProcessedPnlTime: string = (
    await PnlTicksTable.findLatestProcessedBlocktimeAndCount()).maxBlockTime;
  if (lastProcessedTime && Date.parse(lastProcessedTime) >= Date.parse(lastProcessedPnlTime)) {
    logger.info({
      at: 'create-leaderboard#runTask',
      message: 'Skipping run because the current pnl ticks have been processed for the leaderboard',
      pnlTickLatestBlocktime: lastProcessedPnlTime,
      latestBlockTime: lastProcessedTime,
      threshold: config.PNL_TICK_UPDATE_INTERVAL_MS,
    });
    return;
  }

  try {
    leaderboardPnlObjects.push(...await PnlTicksTable.getRankedPnlTicks(timespan));
  } catch (error) {
    logger.error({
      at: 'create-leaderboard#runTask',
      message: `Error when getting ranked pnl ticks for timespan${timespan.toString()}`,
      error,
      timespan,
    });
    return;
  }
  const txId: number = await Transaction.start();
  try {
    await insertLeaderboardPnlObjects(leaderboardPnlObjects, txId);
    await Transaction.commit(txId);
  } catch (error) {
    logger.error({
      at: 'create-leaderboard#runTask',
      message: 'Error when inserting leaderboard pnl objects',
      error,
    });
    await Transaction.rollback(txId);
    return;
  }

  try {
    await LeaderboardPnlProcessedCache.setProcessedTime(timespan, lastProcessedPnlTime, redisClient);
  } catch (error) {
    logger.error({
      at: 'create-leaderboard#runTask',
      message: 'Error when setting processed time',
      error,
    });
  }
}

async function insertLeaderboardPnlObjects(leaderboardPnlObjects: LeaderboardPnlCreateObject[],
  txId: number,
) {
  const chunkedLeaderboardPnlObjects = _.chunk(
    leaderboardPnlObjects, config.LEADERBOARD_PNL_MAX_ROWS_PER_UPSERT);
  for (const leaderboardPnlObjectsToUpsert of chunkedLeaderboardPnlObjects) {
    await LeaderboardPnlTable.bulkUpsert(leaderboardPnlObjectsToUpsert, { txId });
  }
}
