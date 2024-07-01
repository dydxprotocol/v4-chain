import {
  stats,
  logger,
  NodeEnv,
} from '@dydxprotocol-indexer/base';

import config from '../config';
import * as AssetTable from '../stores/asset-table';
import { AssetFromDatabase, AssetsMap, Options } from '../types';
import { startUpdateLoop } from './loopHelper';

let idToAsset: AssetsMap = {};

/**
 * Refresh loop to cache the list of all assets from the database in-memory.
 */
export async function start(): Promise<void> {
  await startUpdateLoop(
    updateAssets,
    config.ASSET_REFRESHER_INTERVAL_MS,
    'updateAssets',
  );
}

/**
 * Updates in-memory map of assets.
 */
export async function updateAssets(options?: Options): Promise<void> {
  const startTime: number = Date.now();
  const assets: AssetFromDatabase[] = await AssetTable.findAll(
    {},
    [],
    options || { readReplica: true },
  );

  const tmpIdToAsset: Record<string, AssetFromDatabase> = {};
  assets.forEach(
    (asset: AssetFromDatabase) => {
      tmpIdToAsset[asset.id] = asset;
    },
  );

  idToAsset = tmpIdToAsset;
  stats.timing(`${config.SERVICE_NAME}.loops.update_assets`, Date.now() - startTime);
}

/**
 * Gets the perpetual market for a given id.
 */
export function getAssetFromId(id: string): AssetFromDatabase {
  const asset: AssetFromDatabase | undefined = idToAsset[id];
  if (asset === undefined) {
    const message: string = `Unable to find asset with assetId: ${id}`;
    logger.error({
      at: 'asset-refresher#getAssetFromId',
      message,
    });
    throw new Error(message);
  }
  return asset;
}

export function getAssetsMap(): AssetsMap {
  return idToAsset;
}

export function clear(): void {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw new Error('clear cannot be used in non-test env');
  }

  idToAsset = {};
}

export function addAsset(asset: AssetFromDatabase): void {
  if (asset.id in idToAsset) {
    throw new Error(`Asset with id ${asset.id} already exists`);
  }

  idToAsset[asset.id] = asset;
}
