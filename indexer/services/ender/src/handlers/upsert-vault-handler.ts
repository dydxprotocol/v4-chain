import { VaultFromDatabase, VaultModel } from '@dydxprotocol-indexer/postgres';
import { VaultAddressesCache } from '@dydxprotocol-indexer/redis';
import { UpsertVaultEventV1 } from '@dydxprotocol-indexer/v4-protos';
import { redisClient } from '../helpers/redis/redis-controller';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class UpsertVaultHandler extends Handler<UpsertVaultEventV1> {
  eventType: string = 'UpsertVaultEventV1';

  public getParallelizationIds(): string[] {
    return [];
  }

  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const vault: VaultFromDatabase = VaultModel.fromJson(resultRow.vault) as VaultFromDatabase;
    console.log('ender upsert vault handler', vault.address);
    await VaultAddressesCache.addVaultAddress(vault.address, redisClient);

    return [];
  }
}
