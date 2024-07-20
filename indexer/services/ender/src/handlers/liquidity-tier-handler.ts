import {
  LiquidityTiersFromDatabase,
  LiquidityTiersModel,
  PerpetualMarketFromDatabase,
  liquidityTierRefresher,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { LiquidityTierUpsertEventV1 } from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';
import * as pg from 'pg';

import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class LiquidityTierHandler extends Handler<LiquidityTierUpsertEventV1> {
  eventType: string = 'LiquidityTierUpsertEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const liquidityTier: LiquidityTiersFromDatabase = LiquidityTiersModel.fromJson(
      resultRow.liquidity_tier) as LiquidityTiersFromDatabase;
    liquidityTierRefresher.upsertLiquidityTier(liquidityTier);
    return this.generateWebsocketEventsForLiquidityTier(liquidityTier);
  }

  private generateWebsocketEventsForLiquidityTier(liquidityTier: LiquidityTiersFromDatabase):
  ConsolidatedKafkaEvent[] {
    const perpetualMarkets: PerpetualMarketFromDatabase[] = _.filter(
      perpetualMarketRefresher.getPerpetualMarketsList(),
      (perpetualMarket: PerpetualMarketFromDatabase) => {
        return perpetualMarket.liquidityTierId === liquidityTier.id;
      },
    );

    if (perpetualMarkets.length === 0) {
      return [];
    }

    return [
      this.generateConsolidatedMarketKafkaEvent(
        JSON.stringify(generatePerpetualMarketMessage(perpetualMarkets)),
      ),
    ];
  }
}
