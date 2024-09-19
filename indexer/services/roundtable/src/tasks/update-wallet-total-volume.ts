import { logger, stats } from '@dydxprotocol-indexer/base';
import { PersistentCacheTable, WalletTable } from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

import config from '../config';

const defaultLastUpdateTime: string = '2020-01-01T00:00:00Z';
const persistentCacheKey: string = 'totalVolumeUpdateTime';

/**
 * Update the total volume for each address in the wallet table.
 */
export default async function runTask(): Promise<void> {
  try {
    const start = Date.now();
    const persistentCacheEntry = await PersistentCacheTable.findById(persistentCacheKey);

    if (!persistentCacheEntry) {
      logger.info({
        at: 'update-address-total-volume#runTask',
        message: `No previous totalVolumeUpdateTime found in persistent cache table. Will use default value: ${defaultLastUpdateTime}`,
      });
    }

    const lastUpdateTime = DateTime.fromISO(persistentCacheEntry
      ? persistentCacheEntry.value
      : defaultLastUpdateTime);
    let currentTime = DateTime.utc();

    // During backfilling, we process one day at a time to reduce roundtable runtime.
    if (currentTime > lastUpdateTime.plus({ days: 1 })) {
      currentTime = lastUpdateTime.plus({ days: 1 });
    }

    await WalletTable.updateTotalVolume(lastUpdateTime.toISO(), currentTime.toISO());
    await PersistentCacheTable.upsert({
      key: persistentCacheKey,
      value: currentTime.toISO(),
    });

    stats.timing(
      `${config.SERVICE_NAME}.update_wallet_total_volume_timing`,
      Date.now() - start,
    );
  } catch (error) {
    logger.error({
      at: 'update-address-total-volume#runTask',
      message: 'Error when updating totalVolume in wallets table',
      error,
    });
  }
}
