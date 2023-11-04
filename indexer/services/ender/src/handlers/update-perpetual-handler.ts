import assert from 'assert';

import { logger } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  perpetualMarketRefresher,
  storeHelpers,
  PerpetualMarketModel,
} from '@dydxprotocol-indexer/postgres';
import { UpdatePerpetualEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../config';
import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class UpdatePerpetualHandler extends Handler<UpdatePerpetualEventV1> {
  eventType: string = 'UpdatePerpetualEventV1';

  public getParallelizationIds(): string[] {
    return [];
  }

  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    if (config.USE_UPDATE_PERPETUAL_HANDLER_SQL_FUNCTION) {
      return this.handleViaSqlFunction();
    }
    return this.handleViaKnex();
  }

  private async handleViaSqlFunction(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_update_perpetual_handler(
        '${JSON.stringify(UpdatePerpetualEventV1.decode(eventDataBinary))}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'UpdateClobPairHandler#handleViaSqlFunction',
        message: 'Failed to handle UpdateClobPairEventV1',
        error,
      });

      throw error;
    });

    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      result.rows[0].result.perpetual_market) as PerpetualMarketFromDatabase;

    await perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);

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
