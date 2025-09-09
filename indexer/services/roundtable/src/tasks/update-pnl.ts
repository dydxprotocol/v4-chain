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

  // Bind the start and end heights to the SQL content and execute within the tx.
  const tx = Transaction.get(txId);
  if (!tx) {
    throw new Error(`Transaction ${txId} not found`);
  }
  await tx.raw(sqlContent, { start, end });

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
 * Updates the PNL (Profit and Loss) table with comprehensive metrics for all relevant subaccounts.
 *
 * The SQL script calculates:
 *
 * 1. Net Transfers: Tracks the cumulative balance of all incoming and outgoing transfers
 *    - Outgoing transfers reduce the balance
 *    - Incoming transfers increase the balance
 *    - Preserves previous net transfer balance from last calculation
 *
 * 2. Total PNL: Tracks the profit and loss from all trading activities
 *    - Includes funding payments received in the period
 *    - Captures changes in position value (mark-to-market)
 *    - Includes cash flows from trading activities (buy/sell)
 *    - Preserves previous PNL from last calculation
 *
 * 3. Equity: Represents the total account value
 *    - Sum of total PNL and net transfers
 *
 * The process:
 * 1. Identifies all subaccounts with either previous PNL records or transfer activity
 * 2. Aggregates transfers (incoming and outgoing) for each subaccount
 * 3. Collects funding data and calculates position values at start and end points
 * 4. Calculates net cash flows from trades (buys and sells)
 * 5. Combines data using the recursive formula: PNL(t+1) = PNL(t) + Funding(t→t+1) +
 *    [Position Value(t+1) - Position Value(t)] + Trade Cash Flow(t→t+1)
 *    where Trade Cash Flow = Sum(Sell Proceeds) - Sum(Buy Costs)
 *
 * This function determines the heights to process based on funding payment events,
 * ensuring PNL is calculated at each height where funding payments occurred,
 * and maintains continuity with previous calculations.
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
      zeroPayments: true,
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
