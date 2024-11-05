import {
  PerpetualMarketFromDatabase, PerpetualMarketModel,
  perpetualMarketRefresher,
} from '@klyraprotocol-indexer/postgres';
import { PerpetualMarketCreateEventV1, PerpetualMarketCreateEventV2 } from '@klyraprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { Handler } from './handler';
import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';

export class PerpetualMarketCreationHandler extends Handler<
  PerpetualMarketCreateEventV1 | PerpetualMarketCreateEventV2
> {
  eventType: string = 'PerpetualMarketCreateEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      resultRow.perpetual_market) as PerpetualMarketFromDatabase;

    perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);
    return [
      this.generateConsolidatedMarketKafkaEvent(
        JSON.stringify(generatePerpetualMarketMessage([perpetualMarket])),
      ),
    ];
  }
}
