import {
  LiquidityTiersCreateObject,
  LiquidityTiersFromDatabase,
  LiquidityTiersTable,
  liquidityTierRefresher,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import { LiquidityTierUpsertEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { QUOTE_CURRENCY_ATOMIC_RESOLUTION } from '../constants';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class LiquidityTierHandler extends Handler<LiquidityTierUpsertEventV1> {
  eventType: string = 'LiquidityTierUpsertEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    await this.runFuncWithTimingStatAndErrorLogging(
      this.upsertLiquidityTier(),
      this.generateTimingStatsOptions('upsert_liquidity_tier'),
    );
    return [];
  }

  private async upsertLiquidityTier(): Promise<void> {
    const liquidityTier: LiquidityTiersFromDatabase = await LiquidityTiersTable.upsert(
      this.getLiquidityTiersCreateObject(this.event),
      { txId: this.txId },
    );
    liquidityTierRefresher.upsertLiquidityTier(liquidityTier);
  }

  private getLiquidityTiersCreateObject(liquidityTier: LiquidityTierUpsertEventV1):
  LiquidityTiersCreateObject {
    return {
      id: liquidityTier.id,
      name: liquidityTier.name,
      initialMarginPpm: liquidityTier.initialMarginPpm.toString(),
      maintenanceFractionPpm: liquidityTier.maintenanceFractionPpm.toString(),
      basePositionNotional: protocolTranslations.quantumsToHuman(
        liquidityTier.basePositionNotional.toString(),
        QUOTE_CURRENCY_ATOMIC_RESOLUTION,
      ).toFixed(6),
    };
  }
}
