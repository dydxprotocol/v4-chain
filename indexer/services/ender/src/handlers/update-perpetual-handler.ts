import { stats } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualMarketModel,
} from '@dydxprotocol-indexer/postgres';
import { UpdatePerpetualEventV1, UpdatePerpetualEventV2 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../config';
import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class UpdatePerpetualHandler extends Handler<
  UpdatePerpetualEventV1 | UpdatePerpetualEventV2> {
  eventType: string = 'UpdatePerpetualEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      resultRow.perpetual_market) as PerpetualMarketFromDatabase;

    await perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_update_perpetual_event.sql_latency`,
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
