import { NodeEnv, logger, stats } from '@dydxprotocol-indexer/base';
import {
  BlockFromDatabase,
  BlockTable,
  IsolationLevel,
  Transaction,
  assetRefresher,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';

import config from '../config';
import { startCandleCache } from './candle-cache';

const INITIAL_BLOCK_HEIGHT: string = '-1';

let currentBlockHeight: string = INITIAL_BLOCK_HEIGHT;

export async function refreshBlockCache(txId?: number): Promise<void> {
  try {
    const block: BlockFromDatabase = await BlockTable.getLatest({ txId });
    currentBlockHeight = block.blockHeight;
  } catch (error) { // Unable to find latest block
    logger.info({
      at: 'block-cache#refreshBlockCache',
      message: 'Unable to refresh block cache most likely due to unable to find latest block',
      error,
    });

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
    stats.increment(`${config.SERVICE_NAME}.block_already_parsed`, 1);
    logger.info({
      at: 'onMessage#onMessage',
      message: `Already processed block with block height: ${blockHeight}, so skipping`,
      blockHeight,
    });
    return true;
  } else if (isNextBlock(blockHeight)) {
    logger.info({
      at: 'block-cache#shouldSkipBlock',
      message: 'Block will be processed',
      blockHeight,
      currentBlockHeight,
    });
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
  return true;
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
  const start: number = Date.now();
  const txId: number = await Transaction.start();
  await Transaction.setIsolationLevel(txId, IsolationLevel.READ_COMMITTED);

  await Promise.all([
    refreshBlockCache(txId),
    // Must be run after perpetualMarketRefresher.start(), because Candle Cache
    // uses the perpetualMarketRefresher cache.
    startCandleCache(txId),
    perpetualMarketRefresher.updatePerpetualMarkets({ txId }),
    assetRefresher.updateAssets({ txId }),
  ]);

  await Transaction.rollback(txId);
  stats.timing(
    `${config.SERVICE_NAME}.initialize_caches`,
    Date.now() - start,
  );
}

export function resetBlockCache(): void {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw new Error('resetBlockCache cannot be used in non-test env');
  }

  currentBlockHeight = INITIAL_BLOCK_HEIGHT;
}
