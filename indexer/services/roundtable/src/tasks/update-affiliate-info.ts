import { logger } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable, AffiliateInfoTable, PersistentCacheKeys, PersistentCacheFromDatabase,
  BlockFromDatabase,
  BlockTable,
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

const defaultLastUpdateTime: string = '2024-09-16T00:00:00Z';

/**
 * Update the affiliate info for all affiliate addresses.
 */
export default async function runTask(): Promise<void> {
  try {
    const persistentCacheEntry: PersistentCacheFromDatabase | undefined = await PersistentCacheTable
      .findById(PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME);

    if (!persistentCacheEntry) {
      logger.info({
        at: 'update-affiliate-info#runTask',
        message: `No previous ${PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME} found in persistent cache table. Will use default value: ${defaultLastUpdateTime}`,
      });
    }

    const lastUpdateTime: DateTime = DateTime.fromISO(persistentCacheEntry
      ? persistentCacheEntry.value
      : defaultLastUpdateTime);

    const latestBlock: BlockFromDatabase = await BlockTable.getLatest();
    if (latestBlock.time === null) {
      throw Error('Failed to get latest block time');
    }
    let windowEndTime = DateTime.fromISO(latestBlock.time);

    // During backfilling, we process one day at a time to reduce roundtable runtime.
    if (windowEndTime > lastUpdateTime.plus({ days: 1 })) {
      windowEndTime = lastUpdateTime.plus({ days: 1 });
    }

    await AffiliateInfoTable.updateInfo(lastUpdateTime.toISO(), windowEndTime.toISO());

  } catch (error) {
    logger.error({
      at: 'update-affiliate-info#runTask',
      message: 'Error when updating affiliate info in affiliate_info table',
      error,
    });
  }
}
