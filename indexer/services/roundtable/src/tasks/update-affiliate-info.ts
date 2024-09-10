import { logger, stats } from '@dydxprotocol-indexer/base';
import { PersistentCacheTable, AffiliateInfoTable, PersistentCacheKeys } from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

import config from '../config';

const defaultLastUpdateTime: string = '2020-01-01T00:00:00Z';

/**
 * Update the affiliate info for all affiliate addresses.
 */
export default async function runTask(): Promise<void> {
  try {
    const start = Date.now();
    const persistentCacheEntry = await PersistentCacheTable.findById(
      PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
    );

    if (!persistentCacheEntry) {
      logger.info({
        at: 'update-affiliate-info#runTask',
        message: `No previous ${PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME} found in persistent cache table. Will use default value: ${defaultLastUpdateTime}`,
      });
    }

    const lastUpdateTime = DateTime.fromISO(persistentCacheEntry
      ? persistentCacheEntry.value
      : defaultLastUpdateTime);
    let windowEndTime = DateTime.utc();

    // During backfilling, we process one day at a time to reduce roundtable runtime.
    if (windowEndTime > lastUpdateTime.plus({ days: 1 })) {
      windowEndTime = lastUpdateTime.plus({ days: 1 });
    }

    await AffiliateInfoTable.updateInfo(lastUpdateTime.toISO(), windowEndTime.toISO());

    stats.timing(
      `${config.SERVICE_NAME}.update_affiliate_info_timing`,
      Date.now() - start,
    );
  } catch (error) {
    logger.error({
      at: 'update-affiliate-info#runTask',
      message: 'Error when updating affiliate info in affiliate_info table',
      error,
    });
  }
}
