import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable,
  AffiliateInfoTable,
  PersistentCacheKeys,
  PersistentCacheFromDatabase,
  BlockFromDatabase,
  BlockTable,
  Transaction,
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

import config from '../config';

const defaultLastUpdateTime: string = '2024-09-16T00:00:00Z';

/**
 * Update the affiliate info for all affiliate addresses.
 */
export default async function runTask(): Promise<void> {
  const taskStart: number = Date.now();
  // Wrap getting cache, updating info, and setting cache in one transaction so that persistent
  // cache and affilitate info table are in sync.
  const txId: number = await Transaction.start();
  try {
    const latestBlock: BlockFromDatabase = await BlockTable.getLatest();
    if (latestBlock.time === null) {
      throw Error('Failed to get latest block time');
    }
    const persistentCacheEntry: PersistentCacheFromDatabase | undefined = await PersistentCacheTable
      .findById(PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME, { txId });
    if (!persistentCacheEntry) {
      logger.info({
        at: 'update-affiliate-info#runTask',
        message: `No previous ${PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME} found in persistent cache table. Will use default value: ${defaultLastUpdateTime}`,
      });
    }
    const windowStartTime: DateTime = DateTime.fromISO(persistentCacheEntry
      ? persistentCacheEntry.value
      : defaultLastUpdateTime);

    // Track how long ago the last update time (windowStartTime) in persistent cache was
    stats.gauge(
      `${config.SERVICE_NAME}.persistent_cache_${PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME}_lag_seconds`,
      DateTime.utc().diff(windowStartTime).as('seconds'),
      { cache: PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME },
    );

    let windowEndTime = DateTime.fromISO(latestBlock.time);
    // During backfilling, we process one day at a time to reduce roundtable runtime.
    if (windowEndTime > windowStartTime.plus({ days: 1 })) {
      windowEndTime = windowStartTime.plus({ days: 1 });
    }

    logger.info({
      at: 'update-affiliate-info#runTask',
      message: `Updating affiliate info from ${windowStartTime.toISO()} to ${windowEndTime.toISO()}`,
    });
    const updateAffiliateInfoStartTime: number = Date.now();
    await AffiliateInfoTable.updateInfo(windowStartTime.toISO(), windowEndTime.toISO(), txId);
    await PersistentCacheTable.upsert({
      key: PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
      value: windowEndTime.toISO(),
    }, { txId });

    await Transaction.commit(txId);
    stats.timing(
      `${config.SERVICE_NAME}.update-affiliate-info.update-txn.timing`,
      Date.now() - updateAffiliateInfoStartTime,
    );
  } catch (error) {
    await Transaction.rollback(txId);
    logger.error({
      at: 'update-affiliate-info#runTask',
      message: 'Error when updating affiliate info in affiliate_info table',
      error,
    });
  }

  stats.timing(`${config.SERVICE_NAME}.update-affiliate-info.total.timing`, Date.now() - taskStart);
}
