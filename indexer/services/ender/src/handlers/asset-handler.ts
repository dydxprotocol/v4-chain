import { stats } from '@dydxprotocol-indexer/base';
import {
  AssetFromDatabase,
  AssetModel,
  assetRefresher,
} from '@dydxprotocol-indexer/postgres';
import { AssetCreateEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../config';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class AssetCreationHandler extends Handler<AssetCreateEventV1> {
  eventType: string = 'AssetCreateEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_asset_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );
    const asset: AssetFromDatabase = AssetModel.fromJson(
      resultRow.asset) as AssetFromDatabase;
    assetRefresher.addAsset(asset);
    return [];
  }
}
