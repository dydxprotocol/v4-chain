import {
  AssetFromDatabase, AssetTable, assetRefresher, marketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { AssetCreateEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class AssetCreationHandler extends Handler<AssetCreateEventV1> {
  eventType: string = 'AssetCreateEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
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
