import { UpsertVaultEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { Handler } from './handler';
import { ConsolidatedKafkaEvent } from '../lib/types';

export class UpsertVaultHandler extends Handler<UpsertVaultEventV1> {
  eventType: string = 'UpsertVaultEventV1';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(_: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    return [];
  }
}
