import {
  stats,
  logger,
  NodeEnv,
} from '@dydxprotocol-indexer/base';

import config from '../config';
import * as BlockTable from '../stores/block-table';
import { BlockFromDatabase, Options } from '../types';
import { startUpdateLoop } from './loopHelper';

let latestBlockHeight: string = '';

/**
 * Refresh loop to cache the latest block height from the database in-memory.
 */
export async function start(): Promise<void> {
  await startUpdateLoop(
    updateBlockHeight,
    config.BLOCK_HEIGHT_REFRESHER_INTERVAL_MS,
    'updateBlockHeight',
  );
}

/**
 * Updates in-memory latest block height.
 */
export async function updateBlockHeight(options?: Options): Promise<void> {
  const startTime: number = Date.now();
  try {
    const latestBlock: BlockFromDatabase = await BlockTable.getLatest(
      options || { readReplica: true },
    );
    latestBlockHeight = latestBlock.blockHeight;
    stats.timing(`${config.SERVICE_NAME}.loops.update_block_height`, Date.now() - startTime);
    // eslint-disable-next-line no-empty
  } catch (error) { }
}

/**
 * Gets the latest block height.
 */
export function getLatestBlockHeight(): string {
  if (!latestBlockHeight) {
    const message: string = 'Unable to find latest block height';
    logger.error({
      at: 'block-height-refresher#getLatestBlockHeight',
      message,
    });
    throw new Error(message);
  }
  return latestBlockHeight;
}

export function clear(): void {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw new Error('clear cannot be used in non-test env');
  }

  latestBlockHeight = '';
}
