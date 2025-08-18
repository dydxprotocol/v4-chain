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

/**
 * Process pnl changes between the specified start and end heights.
 *
 * @param txId Transaction ID for database operations
 * @param start Start block height (exclusive)
 * @param end End block height (inclusive)
 * @returns void
 */
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

  // bind the start height and end height to the sql content
  const result = await Transaction.get(txId)?.raw(sqlContent, {
    start: start,
    end: end,
  });

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
 * Get from persistent cache table the height where pnls were last calculated.
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
 * Updates the PNL (Profit and Loss) table with calculations for all subaccounts with transfer history.
 * 
 * The workflow:
 * 1. Identifies all subaccounts with transfer history up to the specified end height
 * 2. Calculates position effects in two parts:
 *    a) Open position PNL:
 *       - For positions that existed before the start height: Calculates PNL based on price change 
 *         from oracle price at start height to oracle price at end height
 *       - For positions created after start height: Calculates PNL based on price change
 *         from entry price to oracle price at end height
 *    b) Closed position PNL:
 *       - For positions closed between start and end height
 *       - For positions that existed before start height: Uses oracle price at start as reference
 *       - For positions created after start height: Uses entry price as reference
 * 3. Sums up funding payments received in the period between start and end height
 * 4. Calculates total PNL as:
 *    Previous total PNL + Current period funding payments + Current period position effects
 * 
 * The process requires:
 * - Oracle prices at start and end heights for all relevant markets
 * - All open and closed perpetual positions
 * - All funding payments in the period
 * - All previous PNL calculations
 */

export default async function runTask(): Promise<void> {
  const at: string = 'update-pnl#runTask';
  logger.info({ at, message: 'Starting task' });

  // Load SQL script used for data reading and writing.
  const sqlPath = join(__dirname, '..', 'scripts', 'update_pnl.sql');
  const sqlContent = readFileSync(sqlPath, 'utf8');

  // Get all funding payments updates that occurred after last processed height.
  const searchUnprocessedFundingPaymentsHeightStart: string = (
    parseInt(await getLastProcessedHeight(), 10) + 1
  ).toString();

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

  logger.info({
    at,
    message: `Found ${fundingUpdates.results.length} funding periods to process.`,
  });

  stats.gauge(
    `${config.SERVICE_NAME}.update_pnl.num_funding_index_updates_to_process`,
    fundingUpdates.results.length,
  );

  // Get unique heights from funding payments updates.
  const fundingHeights = [...fundingUpdates.results.map((update) => update.createdAtHeight)];

  for (let i = 0; i < fundingHeights.length; i += 1) {
    const txId: number = await Transaction.start();
    try {
      const lastHeight: string = await getLastProcessedHeight();
      const currentHeight: string = fundingHeights[i];

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