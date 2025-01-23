import { stats } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase, PerpetualMarketModel,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import {
  PerpetualMarketCreateEventV1, PerpetualMarketCreateEventV2,
  PerpetualMarketCreateEventV3,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../config';
import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class PerpetualMarketCreationHandler extends Handler<
  PerpetualMarketCreateEventV1 | PerpetualMarketCreateEventV2
  | PerpetualMarketCreateEventV3
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
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_perpetual_market_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );
    return [
      this.generateConsolidatedMarketKafkaEvent(
        JSON.stringify(generatePerpetualMarketMessage([perpetualMarket])),
      ),
    ];
  }
}
