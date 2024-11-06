import {
  AssetFromDatabase,
  AssetModel,
  assetRefresher,
} from '@klyraprotocol-indexer/postgres';
import { AssetCreateEventV1 } from '@klyraprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { Handler } from './handler';
import { ConsolidatedKafkaEvent } from '../lib/types';

export class AssetCreationHandler extends Handler<AssetCreateEventV1> {
  eventType: string = 'AssetCreateEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const asset: AssetFromDatabase = AssetModel.fromJson(
      resultRow.asset) as AssetFromDatabase;
    assetRefresher.addAsset(asset);
    return [];
  }
}
