import { NodeEnv } from '@dydxprotocol-indexer/base';
import {
  PnlTicksFromDatabase,
  PnlTicksTable,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';

import { getVaultMapping, getVaultPnlStartDate } from '../lib/helpers';
import { VaultMapping } from '../types';

let vaultStartPnl: PnlTicksFromDatabase[] = [];

export async function startVaultStartPnlCache(): Promise<void> {
  const vaultMapping: VaultMapping = await getVaultMapping();
  vaultStartPnl = await PnlTicksTable.getLatestPnlTick(
    _.keys(vaultMapping),
    // Add a buffer of 10 minutes to get the first PnL tick for PnL data as PnL ticks aren't
    // created exactly on the hour.
    getVaultPnlStartDate().plus({ minutes: 10 }),
  );
}

export function getVaultStartPnl(): PnlTicksFromDatabase[] {
  return vaultStartPnl;
}

export function clearVaultStartPnl(): void {
  if (process.env.NODE_ENV !== NodeEnv.TEST) {
    throw Error('cannot clear vault start pnl cache outside of test environment');
  }

  vaultStartPnl = [];
}
