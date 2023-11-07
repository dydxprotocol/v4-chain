import { logger } from '@dydxprotocol-indexer/base';
import {
  LiquidityTiersCreateObject,
  LiquidityTiersFromDatabase,
  LiquidityTiersModel,
  LiquidityTiersTable,
  PerpetualMarketFromDatabase,
  liquidityTierRefresher,
  perpetualMarketRefresher,
  protocolTranslations,
  storeHelpers,
} from '@dydxprotocol-indexer/postgres';
import { LiquidityTierUpsertEventV1 } from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';
import * as pg from 'pg';

import config from '../config';
import { QUOTE_CURRENCY_ATOMIC_RESOLUTION } from '../constants';
import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class LiquidityTierHandler extends Handler<LiquidityTierUpsertEventV1> {
  eventType: string = 'LiquidityTierUpsertEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    if (config.USE_LIQUIDITY_TIER_HANDLER_SQL_FUNCTION) {
      return this.handleViaSqlFunction();
    }
    return this.handleViaKnex();
  }

  private async handleViaSqlFunction(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_liquidity_tier_handler(
        '${JSON.stringify(LiquidityTierUpsertEventV1.decode(eventDataBinary))}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'LiquidityTierHandler#handleViaSqlFunction',
        message: 'Failed to handle LiquidityTierUpsertEventV1',
        error,
      });

      throw error;
    });

    const liquidityTier: LiquidityTiersFromDatabase = LiquidityTiersModel.fromJson(
      result.rows[0].result.liquidity_tier) as LiquidityTiersFromDatabase;
    liquidityTierRefresher.upsertLiquidityTier(liquidityTier);
    return this.generateWebsocketEventsForLiquidityTier(liquidityTier);
  }

  private async handleViaKnex(): Promise<ConsolidatedKafkaEvent[]> {
    const liquidityTier:
    LiquidityTiersFromDatabase = await this.runFuncWithTimingStatAndErrorLogging(
      this.upsertLiquidityTier(),
      this.generateTimingStatsOptions('upsert_liquidity_tier'),
    );
    return this.generateWebsocketEventsForLiquidityTier(liquidityTier);
  }

  private async upsertLiquidityTier(): Promise<LiquidityTiersFromDatabase> {
    const liquidityTier: LiquidityTiersFromDatabase = await LiquidityTiersTable.upsert(
      this.getLiquidityTiersCreateObject(this.event),
      { txId: this.txId },
    );
    liquidityTierRefresher.upsertLiquidityTier(liquidityTier);
    return liquidityTier;
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
      ).toFixed(),
    };
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
