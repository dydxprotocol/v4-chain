import { logger } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketCreateObject,
  PerpetualMarketFromDatabase, PerpetualMarketModel,
  perpetualMarketRefresher,
  PerpetualMarketTable,
  protocolTranslations, storeHelpers,
} from '@dydxprotocol-indexer/postgres';
import { PerpetualMarketCreateEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../config';
import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class PerpetualMarketCreationHandler extends Handler<PerpetualMarketCreateEventV1> {
  eventType: string = 'PerpetualMarketCreateEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    if (config.USE_PERPETUAL_MARKET_HANDLER_SQL_FUNCTION) {
      return this.handleViaSqlFunction();
    }
    return this.handleViaKnex();
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  private async handleViaSqlFunction(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_perpetual_market_handler(
        '${JSON.stringify(PerpetualMarketCreateEventV1.decode(eventDataBinary))}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'PerpetualMarketCreationHandler#handleViaSqlFunction',
        message: 'Failed to handle PerpetualMarketCreateEventV1',
        error,
      });

      throw error;
    });

    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      result.rows[0].result.perpetual_market) as PerpetualMarketFromDatabase;

    perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);
    return [
      this.generateConsolidatedMarketKafkaEvent(
        JSON.stringify(generatePerpetualMarketMessage([perpetualMarket])),
      ),
    ];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  private async handleViaKnex(): Promise<ConsolidatedKafkaEvent[]> {
    const perpetualMarket:
    PerpetualMarketFromDatabase = await this.runFuncWithTimingStatAndErrorLogging(
      this.createPerpetualMarket(),
      this.generateTimingStatsOptions('create_perpetual_market'),
    );
    return [
      this.generateConsolidatedMarketKafkaEvent(
        JSON.stringify(generatePerpetualMarketMessage([perpetualMarket])),
      ),
    ];
  }

  private async createPerpetualMarket(): Promise<PerpetualMarketFromDatabase> {
    const perpetualMarket: PerpetualMarketFromDatabase = await PerpetualMarketTable.create(
      this.getPerpetualMarketCreateObject(this.event),
      { txId: this.txId },
    );
    perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);
    return perpetualMarket;
  }

  /**
   * @description Given a PerpetualMarketCreateEventV1 event, generate the `PerpetualMarket`
   * to create.
   */
  private getPerpetualMarketCreateObject(
    perpetualMarketCreateEventV1: PerpetualMarketCreateEventV1,
  ): PerpetualMarketCreateObject {
    return {
      id: perpetualMarketCreateEventV1.id.toString(),
      clobPairId: perpetualMarketCreateEventV1.clobPairId.toString(),
      ticker: perpetualMarketCreateEventV1.ticker,
      marketId: perpetualMarketCreateEventV1.marketId,
      status: protocolTranslations.clobStatusToMarketStatus(perpetualMarketCreateEventV1.status),
      lastPrice: '0',
      priceChange24H: '0',
      trades24H: 0,
      volume24H: '0',
      nextFundingRate: '0',
      openInterest: '0',
      quantumConversionExponent: perpetualMarketCreateEventV1.quantumConversionExponent,
      atomicResolution: perpetualMarketCreateEventV1.atomicResolution,
      subticksPerTick: perpetualMarketCreateEventV1.subticksPerTick,
      stepBaseQuantums: Number(perpetualMarketCreateEventV1.stepBaseQuantums),
      liquidityTierId: perpetualMarketCreateEventV1.liquidityTier,
    };
  }
}
