import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable,
  WalletTable,
} from '@dydxprotocol-indexer/postgres';
import config from '../config';

const defaultLastUpdateTime: string = '2000-01-01T00:00:00Z';
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

    const lastUpdateTime = persistentCacheEntry 
      ? persistentCacheEntry.value 
      : defaultLastUpdateTime;
    const currentTime = new Date().toISOString();

    // On the first run of this roundtable, we need to calculate the total volume for all historical
    // fills. This is a much more demanding task than regular roundtable runs.
    // At time of commit, the total number of rows in 'fills' table in imperator mainnet is ~250M.
    // This can be processed in ~1min with the introduction of 'createdAt' index in 'fills' table.
    // This is relatively short and significanlty shorter than roundtable task cadence. Hence, 
    // special handling for the first run is not required.
    await WalletTable.updateTotalVolume(lastUpdateTime, currentTime);
    await PersistentCacheTable.upsert({
      key: persistentCacheKey,
      value: currentTime,
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
