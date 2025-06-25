import { readFileSync } from 'fs';
import { join } from 'path';

import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable,
  PersistentCacheKeys,
  PersistentCacheFromDatabase,
  Transaction,
  FundingIndexUpdatesTable,
  Ordering,
} from '@dydxprotocol-indexer/postgres';

import config from '../config';

const defaultLastHeight: string = '0';

/**
 * Process funding payment updates between the specified start and end heights.
 *
 * @param txId Transaction ID for database operations
 * @param start Start block height (inclusive)
 * @param end End block height (inclusive)
 * @returns void
 */
async function processFundingPaymentUpdate(
  txId: number,
  start: string,
  end: string,
  sqlContent: string,
): Promise<void> {
  // Skip processing if no new blocks to process
  if (parseInt(end, 10) <= parseInt(start, 10)) {
    logger.info({
      at: 'update-funding-payments#processFundingPaymentUpdate',
      message: `No new blocks to process. Current: ${end}, Last: ${start}`,
    });
    return;
  }

  // bind the last height and current height to the sql content
  await Transaction.get(txId)?.raw(sqlContent, {
    last_height: start,
    current_height: end,
  });

  // Update the persistent cache with the current height
  await PersistentCacheTable.upsert(
    {
      key: PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
      value: end,
    },
    { txId },
  );

  stats.gauge(`${config.SERVICE_NAME}.update_funding_payments.last_processed_height`, parseInt(end, 10));
}

/**
 * Get from persistent cache table the height where funding payments were last processed.
 * If no last processed height is found, use the default value 0.
 *
 * @returns The last processed height.
 */
async function getLastProcessedHeight(): Promise<string> {
  const lastCache: PersistentCacheFromDatabase | undefined = await PersistentCacheTable.findById(
    PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT,
  );
  if (!lastCache) {
    logger.info({
      at: 'update-funding-payments#getLastProcessedHeight',
      message: `No previous ${PersistentCacheKeys.FUNDING_PAYMENTS_LAST_PROCESSED_HEIGHT} found in persistent cache table. Will use default value: ${defaultLastHeight}`,
    });
    return defaultLastHeight;
  }
  return lastCache.value;
}

/**
 * On execution, it will gather from the funding index updates table all heights that we haven't yet
 * processed. It will then process each height in a loop. It will process by taking the
 * funding_payment table at the last processed height and aggregate the fills at the last height + 1
 * to the current height to create a new perpetual position for the subaccount in order to compute
 * the funding payments.
 *
 * Let [x0, x1, ..., xn] be heights where there was a funding index update and was not previously
 * processed. Then we will process in order [(last_processed, x0), (x0, x1), ..., (xn-1, xn)] such
 * that each funding index update is processed.
 *
 * @returns void
 */
export default async function runTask(): Promise<void> {
  const at: string = 'update-funding-payments#runTask';
  logger.info({ at, message: 'Starting task' });

  // Load funding payments SQL script.
  const sqlPath = join(__dirname, '..', 'scripts', 'update_funding_payments.sql');
  const sqlContent = readFileSync(sqlPath, 'utf8');

  // Get all funding index updates that occurred after last processed height.
  // TODO: Move this logic directly into funding payments SQL script.
  const searchUnprocessedFundingIndexHeightStart: string = (
    parseInt(await getLastProcessedHeight(), 10) + 1
  ).toString();
  const fundingUpdates = await FundingIndexUpdatesTable.findAll(
    {
      effectiveAtOrAfterHeight: searchUnprocessedFundingIndexHeightStart,
      distinctFields: ['effectiveAtHeight'],
    },
    [],
    {
      orderBy: [['effectiveAtHeight', Ordering.ASC]],
    },
  );
  logger.info({
    at,
    message: `Found ${fundingUpdates.length} funding index updates to process.`,
  });

  stats.gauge(
    `${config.SERVICE_NAME}.update_funding_payments.num_funding_index_updates_to_process`,
    fundingUpdates.length,
  );

  // Get unique heights from funding updates.
  const fundingHeights = [...fundingUpdates.map((update) => update.effectiveAtHeight)];
  for (let i = 0; i < fundingHeights.length; i += 1) {
    const txId: number = await Transaction.start();
    try {
      const lastHeight: string = await getLastProcessedHeight();
      const currentHeight: string = fundingHeights[i];
      logger.info({
        at,
        message: 'Computing funding payments',
        positionSnapshotHeight: lastHeight,
        fundingOccurredAtHeight: currentHeight,
      });
      // Compute funding payments where last height is last height that funding payments
      // were process and current height is the height of current funding index update.
      await processFundingPaymentUpdate(txId, lastHeight, currentHeight, sqlContent);
      await Transaction.commit(txId);
      logger.info({
        at,
        message: 'Successfully computed funding payments',
        positionSnapshotHeight: lastHeight,
        fundingOccurredAtHeight: currentHeight,
      });
    } catch (error) {
      await Transaction.rollback(txId);
      logger.error({
        at,
        message: 'Error computing funding payments',
        fundingOccurredAtHeight: fundingHeights[i],
        error,
      });
      throw error;
    }
  }
}
