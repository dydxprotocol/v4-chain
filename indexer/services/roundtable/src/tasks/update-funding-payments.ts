import { readFileSync } from 'fs';
import { join } from 'path';

import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PersistentCacheTable,
  PersistentCacheKeys,
  PersistentCacheFromDatabase,
  Transaction,
  FundingIndexUpdatesTable,
} from '@dydxprotocol-indexer/postgres';

import config from '../config';

const defaultLastHeight: string = '0';
const statStart: string = `${config.SERVICE_NAME}.aggregate_data`;

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
}

/**
 * Get the last processed height from the persistent cache table.
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
 * Execute the update_funding_payments.sql file to perform data aggregation.
 *
 * On execution, it will gather from the funding index updates table all heights that
 * we haven't yet processed. It will then process each height in a loop, retrying up to 3 times.
 * It will process by taking the funding_payment table at the last processed height and aggregate
 * the fills at the last height + 1 to the current height to create a new perpetual position for the
 * subaccount in order to compute the funding payments.
 *
 * @returns void
 */
export default async function runTask(): Promise<void> {
  const at: string = 'aggregate-data#runTask';
  logger.info({ at, message: 'Starting aggregate data task.' });
  const taskStart: number = Date.now();

  // Load and execute the update_funding_payments.sql file
  const sqlPath = join(__dirname, '..', 'scripts', 'update_funding_payments.sql');
  const sqlContent = readFileSync(sqlPath, 'utf8');

  // Get all unique effectiveAtHeights from funding index updates since the last processed height.
  const lastProcessedHeight: string = await getLastProcessedHeight();
  const fundingUpdates = await FundingIndexUpdatesTable.findAllHeightsStartingAt(
    lastProcessedHeight,
  );
  for (let i = 0; i < fundingUpdates.length; i += 1) {
    // retry up to 3 times.
    for (let retries = 0; retries < 3; retries += 1) {
      const txId: number = await Transaction.start();
      try {
        // start transaction with last processed height.
        const lastHeight: string = await getLastProcessedHeight();
        // get the current height from the funding index updates.
        const currentHeight: string = fundingUpdates[i];
        // compute the funding payments.
        await processFundingPaymentUpdate(txId, lastHeight, currentHeight, sqlContent);
        logger.info({
          at,
          message: 'Successfully processed funding payment update for heights ',
          start: lastHeight,
          end: currentHeight,
        });
        stats.timing(`${statStart}.executeAggregate`, Date.now() - taskStart);
      } catch (error) {
        await Transaction.rollback(txId);
        logger.error({
          at,
          message: 'Error processing funding payment update',
          end: fundingUpdates[i],
          error,
        });
        if (retries === 2) {
          throw error;
        }
      } finally {
        await Transaction.commit(txId);
        // break out of retry loop.
        retries = 3;
      }
    }
  }
  stats.timing(
    `${config.SERVICE_NAME}.update-funding-payments.total.timing`,
    Date.now() - taskStart,
  );
}
