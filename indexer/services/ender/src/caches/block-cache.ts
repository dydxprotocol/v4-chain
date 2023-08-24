import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  BlockFromDatabase,
  BlockTable,
  IsolationLevel,
  Transaction,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';

import config from '../config';
import { startCandleCache } from './candle-cache';
import { startPriceCache } from './price-cache';

let currentBlockHeight: string = '-1';

export async function refreshBlockCache(txId?: number): Promise<void> {
  const block: BlockFromDatabase | undefined = await BlockTable.getLatest({ txId });
  if (block !== undefined) {
    currentBlockHeight = block.blockHeight;
  }
}

export function getCurrentBlockHeight(): string {
  return currentBlockHeight;
}

export function updateBlockCache(blockHeight: string): void {
  currentBlockHeight = blockHeight;
}

/**
 * If block.height <= currentBlockHeight, then we can skip processing the block.
 * If block.height == currentBlockHeight + 1, then we should process the block.
 * If block.height > currentBlockHeight + 1, then refresh the cache and...
 *  - if block.height <= currentBlockHeight, then we can skip processing the block.
 *  - if block.height >= currentBlockHeight + 1, then we should process the block.
 * @returns true if the block should be skipped, otherwise return false
 */
export async function shouldSkipBlock(
  blockHeight: string,
  canRefreshCache: boolean = true,
): Promise<boolean> {
  if (blockAlreadyProcessed(blockHeight)) {
    stats.increment(`${config.SERVICE_NAME}.block_already_parsed.failure`, 1);
    logger.info({
      at: 'onMessage#onMessage',
      message: `Already processed block with block height: ${blockHeight}, so skipping`,
      blockHeight,
    });
    return true;
  } else if (isNextBlock(blockHeight)) {
    return false;
  } else if (canRefreshCache) {
    const previousBlockHeight: string = getCurrentBlockHeight();
    // Refresh caches including blockCache, in case we are doing an upgrade
    await initializeAllCaches();
    stats.increment(`${config.SERVICE_NAME}.reinitializing_cache`, 1);
    logger.info({
      at: 'onMessage#onMessage',
      message: 'Reinitializing block cache',
      blockHeight,
      previousBlockHeight,
      currentBlockHeight: getCurrentBlockHeight(),
    });
    return shouldSkipBlock(blockHeight, false);
  }

  // For debugging purposes, we want to know if we are skipping blocks
  stats.increment(`${config.SERVICE_NAME}.skipped_block.failure`, 1);
  logger.error({
    at: 'onMessage#onMessage',
    message: 'Indexer full node has skipped a block',
    currentlyProcessingBlockHeight: blockHeight,
    alreadyProcessedBlockHeight: getCurrentBlockHeight(),
  });
  return false;
}

function blockAlreadyProcessed(blockHeight: string): boolean {
  return Big(currentBlockHeight).gte(blockHeight);
}

function isNextBlock(blockHeight: string): boolean {
  return Big(currentBlockHeight).plus(1).eq(blockHeight);
}

/**
 * Initialize all caches, any changes here should be reflected in index.ts.
 * While this isn't directly related to block cache, it has to exist here otherwise there will be a
 * circular dependency between block-cache.ts, index.ts, and the file that holds this function.
 * All caches must be initialized in a Transaction to ensure consistency
 */
export async function initializeAllCaches(): Promise<void> {
  const txId: number = await Transaction.start();
  await Transaction.setIsolationLevel(txId, IsolationLevel.READ_COMMITTED);

  await Promise.all([
    refreshBlockCache(txId),
    // Must be run after perpetualMarketRefresher.start(), because Candle Cache
    // uses the perpetualMarketRefresher cache.
    startCandleCache(txId),
  ]);
  // Must be run after startBlockCache() because it uses the block cache.
  await startPriceCache(getCurrentBlockHeight(), txId);

  await Transaction.rollback(txId);
}
