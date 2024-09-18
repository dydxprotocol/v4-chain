import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable, WalletTable, PersistentCacheKeys, PersistentCacheFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

import config from '../config';

const defaultLastUpdateTime: string = '2020-01-01T00:00:00Z';

/**
 * Update the total volume for each addresses in the wallet table who filled recently.
 */
export default async function runTask(): Promise<void> {
  try {
    const start = Date.now();
    const persistentCacheEntry: PersistentCacheFromDatabase | undefined = await PersistentCacheTable
      .findById(PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME);

    if (!persistentCacheEntry) {
      logger.info({
        at: 'update-wallet-total-volume#runTask',
        message: `No previous ${PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME} found in persistent cache table. Will use default value: ${defaultLastUpdateTime}`,
      });
    }

    const lastUpdateTime: DateTime = DateTime.fromISO(persistentCacheEntry
      ? persistentCacheEntry.value
      : defaultLastUpdateTime);
    let windowEndTime = DateTime.utc();

    // During backfilling, we process one day at a time to reduce roundtable runtime.
    if (windowEndTime > lastUpdateTime.plus({ days: 1 })) {
      windowEndTime = lastUpdateTime.plus({ days: 1 });
    }

    await WalletTable.updateTotalVolume(lastUpdateTime.toISO(), windowEndTime.toISO());

    stats.timing(
      `${config.SERVICE_NAME}.update_wallet_total_volume_timing`,
      Date.now() - start,
    );
  } catch (error) {
    logger.error({
      at: 'update-wallet-total-volume#runTask',
      message: 'Error when updating totalVolume in wallets table',
      error,
    });
  }
}
