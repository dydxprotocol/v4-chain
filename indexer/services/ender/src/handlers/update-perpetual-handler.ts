import { logger } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  storeHelpers,
  PerpetualMarketModel,
} from '@dydxprotocol-indexer/postgres';
import { UpdatePerpetualEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

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
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_update_perpetual_handler(
        '${JSON.stringify(UpdatePerpetualEventV1.decode(eventDataBinary))}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'UpdatePerpetualHandler#handleViaSqlFunction',
        message: 'Failed to handle UpdatePerpetualEventV1',
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
}
