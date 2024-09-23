import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable,
  WalletTable,
  PersistentCacheKeys,
  PersistentCacheFromDatabase,
  Transaction,
  BlockFromDatabase,
  BlockTable,
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

import config from '../config';

const defaultLastUpdateTime: string = '2023-10-26T00:00:00Z';

/**
 * Update the total volume for each addresses in the wallet table who filled recently.
 */
export default async function runTask(): Promise<void> {
  // Wrap getting cache, updating info, and setting cache in one transaction so that persistent
  // cache and affilitate info table are in sync.
  const txId: number = await Transaction.start();
  try {
    const latestBlock: BlockFromDatabase = await BlockTable.getLatest();
    if (latestBlock.time === null) {
      throw Error('Failed to get latest block time');
    }

    const persistentCacheEntry: PersistentCacheFromDatabase | undefined = await PersistentCacheTable
      .findById(PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME, { txId });

    if (!persistentCacheEntry) {
      logger.info({
        at: 'update-wallet-total-volume#runTask',
        message: `No previous ${PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME} found in persistent cache table. Will use default value: ${defaultLastUpdateTime}`,
      });
    }

    const windowStartTime: DateTime = DateTime.fromISO(persistentCacheEntry
      ? persistentCacheEntry.value
      : defaultLastUpdateTime);

    // Track how long ago the last update time (windowStartTime) in persistent cache was
    stats.gauge(
      `${config.SERVICE_NAME}.persistent_cache_${PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME}_lag_seconds`,
      DateTime.utc().diff(windowStartTime).as('seconds'),
      { cache: PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME },
    );

    let windowEndTime = DateTime.fromISO(latestBlock.time);
    // During backfilling, we process one day at a time to reduce roundtable runtime.
    if (windowEndTime > windowStartTime.plus({ days: 1 })) {
      windowEndTime = windowStartTime.plus({ days: 1 });
    }

    logger.info({
      at: 'update-wallet-total-volume#runTask',
      message: `Updating wallet total volume from ${windowStartTime.toISO()} to ${windowEndTime.toISO()}`,
    });
    await WalletTable.updateTotalVolume(windowStartTime.toISO(), windowEndTime.toISO(), txId);
    await PersistentCacheTable.upsert({
      key: PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME,
      value: windowEndTime.toISO(),
    }, { txId });

    await Transaction.commit(txId);
  } catch (error) {
    await Transaction.rollback(txId);
    logger.error({
      at: 'update-wallet-total-volume#runTask',
      message: 'Error when updating totalVolume in wallets table',
      error,
    });
  }
}
