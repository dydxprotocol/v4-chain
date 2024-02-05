import {
  NodeEnv,
  logger,
  stats,
} from '@dydxprotocol-indexer/base';

import config from '../config';
import * as LiquidityTiersTable from '../stores/liquidity-tiers-table';
import { LiquidityTiersFromDatabase, LiquidityTiersMap, Options } from '../types';
import { startUpdateLoop } from './loopHelper';

let idToLiquidityTier: LiquidityTiersMap = {};

/**
 * Refresh loop to cache the list of all liquidity tiers from the database in-memory.
 */
export async function start(): Promise<void> {
  await startUpdateLoop(
    updateLiquidityTiers,
    config.LIQUIDITY_TIER_REFRESHER_INTERVAL_MS,
    'updateLiquidityTiers',
  );
}

/**
 * Updates in-memory map of liquidity tiers.
 */
export async function updateLiquidityTiers(options?: Options): Promise<void> {
  const startTime: number = Date.now();
  const liquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll(
    {},
    [],
    options || { readReplica: true },
  );

  const tmpIdToLiquidityTier: Record<string, LiquidityTiersFromDatabase> = {};
  liquidityTiers.forEach(
    (liquidityTier: LiquidityTiersFromDatabase) => {
      tmpIdToLiquidityTier[liquidityTier.id] = liquidityTier;
    },
  );

  idToLiquidityTier = tmpIdToLiquidityTier;
  stats.timing(`${config.SERVICE_NAME}.loops.update_liquidity_tiers`, Date.now() - startTime);
}

/**
 * Gets the liquidity tier for a given id.
 */
export function getLiquidityTierFromId(id: number): LiquidityTiersFromDatabase {
  const tier: LiquidityTiersFromDatabase | undefined = idToLiquidityTier[id];
  if (tier === undefined) {
    const message: string = `Unable to find liquidity tier with id: ${id}`;
    logger.error({
      at: 'liquidity-tier-refresher#getLiquidityTierFromId',
      message,
    });
    throw new Error(message);
  }
  return tier;
}

export function getLiquidityTiersMap(): LiquidityTiersMap {
  return idToLiquidityTier;
}

export function upsertLiquidityTier(liquidityTier: LiquidityTiersFromDatabase): void {
  idToLiquidityTier[liquidityTier.id] = liquidityTier;
}

/**
 * Clears the in-memory map of liquidity tier ids to liquidity tiers.
 * Used for testing.
 */
export function clear(): void {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw new Error('clear cannot be used in non-test env');
  }
  idToLiquidityTier = {};
}
