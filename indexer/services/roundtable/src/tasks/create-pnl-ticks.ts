import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  BlockFromDatabase,
  BlockTable,
  IsolationLevel,
  PnlTicksCreateObject,
  PnlTicksTable,
  Transaction,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { LatestAccountPnlTicksCache } from '@dydxprotocol-indexer/redis';
import _ from 'lodash';

import config from '../config';
import { getPnlTicksCreateObjects, normalizeStartTime } from '../helpers/pnl-ticks-helper';
import { redisClient } from '../helpers/redis';

export default async function runTask(): Promise<void> {
  const startGetNewTicks: number = Date.now();
  const [
    block,
    {
      maxBlockTime,
      count,
    },
  ]: [
    BlockFromDatabase,
    {
      maxBlockTime: string,
      count: number,
    },
  ] = await Promise.all([
    BlockTable.getLatest({ readReplica: true }),
    PnlTicksTable.findLatestProcessedBlocktimeAndCount(),
  ]);
  const latestBlockTime: string = block.time;
  const latestBlockHeight: string = block.blockHeight;
  // Check that the latest block time is within PNL_TICK_UPDATE_INTERVAL_MS of the last computed
  // PNL tick block time.
  if (
    Date.parse(latestBlockTime) - normalizeStartTime(new Date(maxBlockTime)).getTime() <
    config.PNL_TICK_UPDATE_INTERVAL_MS && count < config.PNL_TICK_MAX_ACCOUNTS_PER_RUN
  ) {
    logger.info({
      at: 'create-pnl-ticks#runTask',
      message: 'Skipping run because update interval has not been reached and all subaccounts have been processed',
      pnlTickLatestBlocktime: maxBlockTime,
      latestBlockTime,
      threshold: config.PNL_TICK_UPDATE_INTERVAL_MS,
    });
    return;
  }

  // Start a transaction to ensure different table reads are consistent. Use a repeatable read
  // to ensure all reads within the transaction are consistent.
  const txId: number = await Transaction.start();
  await Transaction.setIsolationLevel(txId, IsolationLevel.REPEATABLE_READ);
  let newTicksToCreate: PnlTicksCreateObject[] = [];
  try {
    await perpetualMarketRefresher.updatePerpetualMarkets();
    newTicksToCreate = await getPnlTicksCreateObjects(txId);
  } catch (error) {
    logger.error({
      at: 'create-pnl-ticks#runTask',
      message: 'Error when getting pnl ticks',
      error,
      latestBlockHeight,
      latestBlockTime,
      txId,
    });
    return;
  } finally {
    // Make sure to always roll-back the transaction so there are no hanging DB connections.
    // Transaction is read-only, so roll back.
    await Transaction.rollback(txId);
  }

  stats.timing(
    `${config.SERVICE_NAME}_get_ticks_timing`,
    new Date().getTime() - startGetNewTicks,
  );

  const startNewTicksCreation: number = new Date().getTime();
  const newTicks: PnlTicksCreateObject[] = await batchCreateNewTicks(
    newTicksToCreate,
  );
  const newestTicksPerSubaccount: { [subaccountId: string]: PnlTicksCreateObject } = _.keyBy(
    newTicks,
    'subaccountId',
  );
  await LatestAccountPnlTicksCache.set(newestTicksPerSubaccount, redisClient);

  stats.timing(
    `${config.SERVICE_NAME}_generate_ticks_timing`,
    new Date().getTime() - startNewTicksCreation,
  );
}

async function batchCreateNewTicks(
  ticks: PnlTicksCreateObject[],
): Promise<PnlTicksCreateObject[]> {
  const newTicks: PnlTicksCreateObject[] = [];
  // Break messages into chunks of length (at-most) PNL_TICK_MAX_ROWS_PER_UPSERT.
  const chunkedTicks = _.chunk(ticks, config.PNL_TICK_MAX_ROWS_PER_UPSERT);
  for (const ticksToCreate of chunkedTicks) {
    const txId: number = await Transaction.start();
    try {
      await PnlTicksTable.createMany(ticksToCreate, { txId });
      await Transaction.commit(txId);
    } catch (e) {
      await Transaction.rollback(txId);
      throw e;
    }
    newTicks.push(...ticksToCreate);
  }

  stats.gauge(
    `${config.SERVICE_NAME}_pnl_tick_create_count`,
    newTicks.length,
  );

  return newTicks;
}
