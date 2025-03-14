import { vaultRefresher, VaultFromDatabase, VaultModel } from '@dydxprotocol-indexer/postgres';
import { UpsertVaultEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class UpsertVaultHandler extends Handler<UpsertVaultEventV1> {
  eventType: string = 'UpsertVaultEventV1';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const vault: VaultFromDatabase = VaultModel.fromJson(resultRow.vault) as VaultFromDatabase;
    vaultRefresher.addVault(vault.address);

    return [];
  }
}
