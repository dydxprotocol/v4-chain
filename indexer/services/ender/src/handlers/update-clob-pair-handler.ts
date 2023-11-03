import assert from 'assert';

import { logger } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  PerpetualMarketModel,
  PerpetualMarketTable,
  perpetualMarketRefresher,
  protocolTranslations,
  storeHelpers,
} from '@dydxprotocol-indexer/postgres';
import { UpdateClobPairEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../config';
import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class UpdateClobPairHandler extends Handler<UpdateClobPairEventV1> {
  eventType: string = 'UpdateClobPairEventV1';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    if (config.USE_UPDATE_CLOB_PAIR_HANDLER_SQL_FUNCTION) {
      return this.handleViaSqlFunction();
    }
    return this.handleViaKnex();
  }

  private async handleViaSqlFunction(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_update_clob_pair_handler(
        '${JSON.stringify(UpdateClobPairEventV1.decode(eventDataBinary))}'
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

    return [
      this.generateConsolidatedMarketKafkaEvent(
        JSON.stringify(generatePerpetualMarketMessage([perpetualMarket])),
      ),
    ];
  }

  private async handleViaKnex(): Promise<ConsolidatedKafkaEvent[]> {
    const perpetualMarket:
    PerpetualMarketFromDatabase = await this.runFuncWithTimingStatAndErrorLogging(
      this.updateClobPair(),
      this.generateTimingStatsOptions('update_clob_pair'),
    );
    return [
      this.generateConsolidatedMarketKafkaEvent(
        JSON.stringify(generatePerpetualMarketMessage([perpetualMarket])),
      ),
    ];
  }

  private async updateClobPair(): Promise<PerpetualMarketFromDatabase> {
    // perpetualMarketRefresher.getPerpetualMarketFromClobPairId() cannot be undefined because it
    // is validated by UpdateClobPairValidator.
    const perpetualMarketId: string = perpetualMarketRefresher.getPerpetualMarketFromClobPairId(
      this.event.clobPairId.toString(),
    )!.id;
    const perpetualMarket:
    PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable.update({
      id: perpetualMarketId,
      status: protocolTranslations.clobStatusToMarketStatus(this.event.status),
      quantumConversionExponent: this.event.quantumConversionExponent,
      subticksPerTick: this.event.subticksPerTick,
      stepBaseQuantums: Number(this.event.stepBaseQuantums),
    }, { txId: this.txId });

    if (perpetualMarket === undefined) {
      this.logAndThrowParseMessageError(
        'Could not find perpetual market with corresponding clobPairId',
        { event: this.event },
      );
      // This assert should never be hit because a ParseMessageError should be thrown above.
      assert(false);
    }

    await perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);
    return perpetualMarket;
  }
}
