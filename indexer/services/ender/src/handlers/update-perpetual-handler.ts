import {
  PerpetualMarketFromDatabase, PerpetualMarketTable, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { UpdatePerpetualEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class UpdatePerpetualHandler extends Handler<UpdatePerpetualEventV1> {
  eventType: string = 'UpdatePerpetualEventV1';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    await this.runFuncWithTimingStatAndErrorLogging(
      this.updatePerpetual(),
      this.generateTimingStatsOptions('update_perpetual'),
    );
    return [];
  }

  private async updatePerpetual(): Promise<void> {
    const perpetualMarket:
    PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable.update({
      id: this.event.id.toString(),
      ticker: this.event.ticker,
      marketId: this.event.marketId,
      atomicResolution: this.event.atomicResolution,
      liquidityTierId: this.event.liquidityTier,
    }, { txId: this.txId });

    if (perpetualMarket === undefined) {
      return this.logAndThrowParseMessageError(
        'Could not find perpetual market with corresponding updatePerpetualEvent.id',
        { event: this.event },
      );
    }
    await perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);
  }
}
