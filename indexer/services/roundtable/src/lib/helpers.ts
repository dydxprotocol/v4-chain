import { logger } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  SubaccountTable,
  VaultFromDatabase, VaultTable, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
import { DateTime } from 'luxon';

// NEXT: collapse with comlink vault-helpers.

export interface VaultMapping {
  [subaccountId: string]: VaultFromDatabase,
}

export async function getVaultMapping(): Promise<VaultMapping> {
  const vaults: VaultFromDatabase[] = await VaultTable.findAll(
    {},
    [],
    {},
  );
  const vaultMapping: VaultMapping = _.zipObject(
    vaults.map((vault: VaultFromDatabase): string => {
      return SubaccountTable.uuid(vault.address, 0);
    }),
    vaults,
  );
  const validVaultMapping: VaultMapping = {};
  for (const subaccountId of _.keys(vaultMapping)) {
    const perpetual: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(
        vaultMapping[subaccountId].clobPairId,
      );
    if (perpetual === undefined) {
      logger.warning({
        at: 'get-vault-mapping',
        message: `Vault clob pair id ${vaultMapping[subaccountId]} does not correspond to a ` +
          'perpetual market.',
        subaccountId,
      });
      continue;
    }
    validVaultMapping[subaccountId] = vaultMapping[subaccountId];
  }
  return validVaultMapping;
}

export function getVaultPnlStartDate(): DateTime {
  const startDate: DateTime = DateTime.fromISO('2024-01-01T00:00:00Z').toUTC(); // NEXT: use config.
  return startDate;
}
