import {
  stats,
  NodeEnv,
} from '@dydxprotocol-indexer/base';
import _ from 'lodash';

import config from '../config';
import * as VaultTable from '../stores/vault-table';
import { Options, VaultFromDatabase } from '../types';
import { startUpdateLoop } from './loopHelper';

let vaultAddresses: Set<string> = new Set();

/**
 * Refresh loop to cache the list of all vault addresses from the database in-memory.
 */
export async function start(): Promise<void> {
  await startUpdateLoop(
    updateVaults,
    config.VAULT_REFRESHER_INTERVAL_MS,
    'updateVaults',
  );
}

/**
 * Updates in-memory set of vault addresses.
 */
export async function updateVaults(options?: Options): Promise<void> {
  const startTime: number = Date.now();
  const vaults: VaultFromDatabase[] = await VaultTable.findAll(
    {}, [], options || { readReplica: true },
  );
  vaultAddresses = new Set(_.map(vaults, 'address'));
  stats.timing(`${config.SERVICE_NAME}.loops.update_vaults`, Date.now() - startTime);
}

export function getVaultAddresses(): Set<string> {
  return vaultAddresses;
}

export function clear(): void {
  if (config.NODE_ENV !== NodeEnv.TEST) {
    throw new Error('clear cannot be used in non-test env');
  }

  vaultAddresses = new Set();
}

export function isVault(address: string): boolean {
  return vaultAddresses.has(address);
}
