import { logger } from '@dydxprotocol-indexer/base';
import {
  AssetFromDatabase,
  AssetModel,
  AssetTable,
  assetRefresher,
  marketRefresher,
  storeHelpers,
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
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    if (config.USE_ASSET_CREATE_HANDLER_SQL_FUNCTION) {
      return this.handleViaSqlFunction();
    }
    return this.handleViaKnex();
  }

  private async handleViaSqlFunction(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_asset_create_handler(
        '${JSON.stringify(AssetCreateEventV1.decode(eventDataBinary))}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'AssetCreationHandler#handleViaSqlFunction',
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

  private async handleViaKnex(): Promise<ConsolidatedKafkaEvent[]> {
    await this.runFuncWithTimingStatAndErrorLogging(
      this.createAsset(),
      this.generateTimingStatsOptions('create_asset'),
    );
    return [];
  }

  private async createAsset(): Promise<void> {
    if (this.event.hasMarket) {
      marketRefresher.getMarketFromId(
        this.event.marketId,
      );
    }
    const asset: AssetFromDatabase = await AssetTable.create({
      id: this.event.id.toString(),
      symbol: this.event.symbol,
      atomicResolution: this.event.atomicResolution,
      hasMarket: this.event.hasMarket,
      marketId: this.event.marketId,
    }, { txId: this.txId });
    assetRefresher.addAsset(asset);
  }
}
