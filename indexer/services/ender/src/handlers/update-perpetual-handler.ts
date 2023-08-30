import assert from 'assert';

import {
  PerpetualMarketFromDatabase, PerpetualMarketTable, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { UpdatePerpetualEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class UpdatePerpetualHandler extends Handler<UpdatePerpetualEventV1> {
  eventType: string = 'UpdatePerpetualEventV1';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const perpetualMarket:
    PerpetualMarketFromDatabase = await this.runFuncWithTimingStatAndErrorLogging(
      this.updatePerpetual(),
      this.generateTimingStatsOptions('update_perpetual'),
    );
    return [
      this.generateConsolidatedMarketKafkaEvent(
        JSON.stringify(generatePerpetualMarketMessage([perpetualMarket])),
      ),
    ];
  }

  private async updatePerpetual(): Promise<PerpetualMarketFromDatabase> {
    const perpetualMarket:
    PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable.update({
      id: this.event.id.toString(),
      ticker: this.event.ticker,
      marketId: this.event.marketId,
      atomicResolution: this.event.atomicResolution,
      liquidityTierId: this.event.liquidityTier,
    }, { txId: this.txId });

    if (perpetualMarket === undefined) {
      this.logAndThrowParseMessageError(
        'Could not find perpetual market with corresponding updatePerpetualEvent.id',
        { event: this.event },
      );
      // This assert should never be hit because a ParseMessageError should be thrown above.
      assert(false);
    }

    await perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);
    return perpetualMarket;
  }
}
