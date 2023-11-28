import { logger } from '@dydxprotocol-indexer/base';
import {
  AssetFromDatabase,
  AssetModel,
  assetRefresher,
  storeHelpers,
} from '@dydxprotocol-indexer/postgres';
import { AssetCreateEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class AssetCreationHandler extends Handler<AssetCreateEventV1> {
  eventType: string = 'AssetCreateEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_asset_create_handler(
        '${JSON.stringify(AssetCreateEventV1.decode(eventDataBinary))}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'AssetCreationHandler#internalHandle',
        message: 'Failed to handle AssetCreateEventV1',
        error,
      });

      throw error;
    });

    const asset: AssetFromDatabase = AssetModel.fromJson(
      result.rows[0].result.asset) as AssetFromDatabase;
    assetRefresher.addAsset(asset);
    return [];
  }
}
