import { logger } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable,
  AffiliateInfoTable,
  PersistentCacheKeys,
  PersistentCacheFromDatabase,
  BlockFromDatabase,
  BlockTable,
  Transaction,
  IsolationLevel,
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

const defaultLastUpdateTime: string = '2024-09-16T00:00:00Z';

/**
 * Update the affiliate info for all affiliate addresses.
 */
export default async function runTask(): Promise<void> {
  const latestBlock: BlockFromDatabase = await BlockTable.getLatest();
  if (latestBlock.time === null) {
    throw Error('Failed to get latest block time');
  }

  // Wrap getting cache, updating info, and setting cache in one transaction with row locking to
  // prevent race condition on persistent cache rows between read and write.
  const txId: number = await Transaction.start();
  await Transaction.setIsolationLevel(txId, IsolationLevel.REPEATABLE_READ);
  try {
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

    let windowEndTime = DateTime.fromISO(latestBlock.time);
    // During backfilling, we process one day at a time to reduce roundtable runtime.
    if (windowEndTime > windowStartTime.plus({ days: 1 })) {
      windowEndTime = windowStartTime.plus({ days: 1 });
    }

    await AffiliateInfoTable.updateInfo(windowStartTime.toISO(), windowEndTime.toISO(), { txId });
    await PersistentCacheTable.upsert({
      key: PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
      value: windowEndTime.toISO(),
    }, { txId });

    await Transaction.commit(txId);
  } catch (error) {
    await Transaction.rollback(txId);
    logger.error({
      at: 'update-affiliate-info#runTask',
      message: 'Error when updating affiliate info in affiliate_info table',
      error,
    });
  }
}
