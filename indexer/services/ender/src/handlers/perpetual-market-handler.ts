import {
  helpers,
  marketRefresher,
  perpetualMarketRefresher,
  PerpetualMarketTable,
  liquidityTierRefresher,
} from '@dydxprotocol-indexer/postgres';
import { PerpetualMarketCreateEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class PerpetualMarketCreationHandler extends Handler<PerpetualMarketCreateEventV1> {
  eventType: string = 'PerpetualMarketCreateEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    await this.runFuncWithTimingStatAndErrorLogging(
      this.createPerpetualMarket(),
      this.generateTimingStatsOptions('create_perpetual_market'),
    );
    return [];
  }

  private async createPerpetualMarket(): Promise<void> {
    marketRefresher.getMarketFromId(
      this.event.marketId,
    );
    liquidityTierRefresher.getLiquidityTierFromId(
      this.event.liquidityTier,
    );
    await PerpetualMarketTable.create(
      helpers.getPerpetualMarketCreateObject(this.event),
      { txId: this.txId },
    );
    await Promise.all([
      perpetualMarketRefresher.updatePerpetualMarkets({ txId: this.txId }),
      liquidityTierRefresher.updateLiquidityTiers({ txId: this.txId }),
    ]);
  }
}
