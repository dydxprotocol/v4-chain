import { readFileSync } from 'fs';
import { join } from 'path';

import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable,
  PersistentCacheKeys,
  PersistentCacheFromDatabase,
  Transaction,
  BlockFromDatabase,
  BlockTable,
} from '@dydxprotocol-indexer/postgres';

import config from '../config';

const defaultLastHeight: string = '0';
const statStart: string = `${config.SERVICE_NAME}.aggregate_data`;

/**
 * Execute the update_funding_payments.sql file to perform data aggregation.
 */
export default async function runTask(): Promise<void> {
  const at: string = 'aggregate-data#runTask';
  logger.info({ at, message: 'Starting aggregate data task.' });

  const taskStart: number = Date.now();
  // Wrap getting cache, updating info, and setting cache in one transaction so that persistent
  // cache and funding payments are in sync.
  const txId: number = await Transaction.start();
  try {
    const latestBlock: BlockFromDatabase = await BlockTable.getLatest();
    if (latestBlock.blockHeight === null) {
      throw Error('Failed to get latest block height');
    }

    const persistentCacheEntry: PersistentCacheFromDatabase | undefined = await PersistentCacheTable
      .findById(PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT, { txId });

    if (!persistentCacheEntry) {
      logger.info({
        at: 'update-funding-payments#runTask',
        message: `No previous ${PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT} found in persistent cache table. Will use default value: ${defaultLastHeight}`,
      });
    }

    const lastHeight: string = persistentCacheEntry?.value ?? defaultLastHeight;
    const currentHeight: string = latestBlock.blockHeight;
    // Load and execute the update_funding_payments.sql file
    const sqlPath = join(__dirname, '..', 'scripts', 'update_funding_payments.sql');
    const sqlContent = readFileSync(sqlPath, 'utf8');

    // bind the last height and current height to the sql content
    await Transaction.get(txId)?.raw(sqlContent, {
      last_height: lastHeight,
      current_height: currentHeight
    });

    // Update the persistent cache with the current height
    await PersistentCacheTable.upsert({
      key: PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
      value: currentHeight,
    }, { txId });

    stats.timing(`${statStart}.executeAggregate`, Date.now() - taskStart);
    logger.info({ at, message: 'Successfully executed aggregate task.' });

    await Transaction.commit(txId);
  } catch (error) {
    await Transaction.rollback(txId);
    logger.error({
      at,
      message: 'Error executing aggregate task',
      error,
    });
    throw error;
  }

  stats.timing(`${config.SERVICE_NAME}.update-funding-payments.total.timing`, Date.now() - taskStart);
}
