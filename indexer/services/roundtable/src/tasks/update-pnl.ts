import { readFileSync } from 'fs';
import { join } from 'path';

import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable,
  PersistentCacheKeys,
  PersistentCacheFromDatabase,
  Transaction,
  FundingPaymentsTable,
  Ordering,
} from '@dydxprotocol-indexer/postgres';

import config from '../config';

const defaultLastHeight: string = '0';

async function processPnlUpdate(
  txId: number,
  start: string,
  end: string,
  sqlContent: string,
): Promise<void> {
  // Skip processing if no new blocks to process
  if (parseInt(end, 10) <= parseInt(start, 10)) {
    logger.info({
      at: 'update-pnl#processPnlUpdate',
      message: `No new blocks to process. Current: ${end}, Last: ${start}`,
    });
    return;
  }

  console.log('Processing PNL update from height', start, 'to height', end);

  // Actual logic
  const result = await Transaction.get(txId)?.raw(sqlContent, {
    start: start,
    end: end,
  });

  console.log('SQL execution result:', result);

  // Update the persistent cache with the current height
  await PersistentCacheTable.upsert(
    {
      key: PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT,
      value: end,
    },
    { txId },
  );

  stats.gauge(`${config.SERVICE_NAME}.update_pnl.last_processed_height`, parseInt(end, 10));
}

/**
 * Get from persistent cache table the height where pnl ticks were last processed.
 * If no last processed height is found, use the default value 0.
 *
 * @returns The last processed height.
 */
async function getLastProcessedHeight(): Promise<string> {
  const lastCache: PersistentCacheFromDatabase | undefined = await PersistentCacheTable.findById(
    PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT,
  );
  if (!lastCache) {
    logger.info({
      at: 'update-pnl#getLastProcessedHeight',
      message: `No previous ${PersistentCacheKeys.PNL_LAST_PROCESSED_HEIGHT} found in persistent cache table. Will use default value: ${defaultLastHeight}`,
    });
    return defaultLastHeight;
  }
  return lastCache.value;
}

/**
 * We need to calculate pnl ticks in order [(last_processed, t0), (t0, t1), ..., (tn-1, tn)]
 * where each interval is of 1 hour
 * For each interval, we need to fetch the following information from block height t to
 * block height t+1:
 * 1. all subaccounts with transfer history at block height t
 * 2. all pnl values at block height t-1 for these subaccounts, 0 for new subaccounts
 * 3. all funding payments in this time period
 * 4. all open perpetual positions at block height t
 * 5. all perpetual positions closed from block height t-1 to t
 * 6. all oracle pricing from block height t-1 to t
 */

export default async function runTask(): Promise<void> {
  const at: string = 'update-pnl#runTask';
  logger.info({ at, message: 'Starting task' });

  // Load funding payments SQL script.
  const sqlPath = join(__dirname, '..', 'scripts', 'update_pnl.sql');
  const sqlContent = readFileSync(sqlPath, 'utf8');

  // Get all funding payments updates that occurred after last processed height.
  const searchUnprocessedFundingPaymentsHeightStart: string = (
    parseInt(await getLastProcessedHeight(), 10) + 1
  ).toString();
  console.log('Last height:', searchUnprocessedFundingPaymentsHeightStart);

  const fundingUpdates = await FundingPaymentsTable.findAll(
    {
      createdOnOrAfterHeight: searchUnprocessedFundingPaymentsHeightStart,
      distinctFields: ['createdAtHeight'],
    },
    [],
    {
      orderBy: [['createdAtHeight', Ordering.ASC]],
    },
  );
  console.log('all funding payments:', fundingUpdates.results);                             // TO BE DELETED

  logger.info({
    at,
    message: `Found ${fundingUpdates.results.length} funding periods to process.`,
  });

  stats.gauge(
    `${config.SERVICE_NAME}.update_pnl.num_funding_index_updates_to_process`,
    fundingUpdates.results.length,
  );

  // Get unique heights from funding updates.
  const fundingHeights = [...fundingUpdates.results.map((update) => update.createdAtHeight)];
  console.log('Funding Heights:', fundingHeights);                                          // TO BE DELETED

  for (let i = 0; i < fundingHeights.length; i += 1) {
    const txId: number = await Transaction.start();
    try {
      const lastHeight: string = await getLastProcessedHeight();
      const currentHeight: string = fundingHeights[i];

      console.log(`Processing PNL calculation: ${lastHeight} -> ${currentHeight}`);

      logger.info({
        at,
        message: 'Computing profit and loss',
        positionSnapshotHeight: lastHeight,
        pnlOccurredAtHeight: currentHeight,
      });
      await processPnlUpdate(txId, lastHeight, currentHeight, sqlContent);
      await Transaction.commit(txId);
      logger.info({
        at,
        message: 'Successfully computed profit and loss',
        positionSnapshotHeight: lastHeight,
        pnlOccurredAtHeight: currentHeight,
      });
    } catch (error) {
      await Transaction.rollback(txId);
      logger.error({
        at,
        message: 'Error computing profit and loss',
        pnlOccurredAtHeight: fundingHeights[i],
        error,
      });
      throw error;
    }
  }
}