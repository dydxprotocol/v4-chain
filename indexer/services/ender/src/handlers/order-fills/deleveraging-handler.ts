import {
  FillFromDatabase,
  FillModel,
  PerpetualMarketFromDatabase,
  PerpetualMarketModel,
  PerpetualPositionFromDatabase,
  PerpetualPositionModel,
} from '@dydxprotocol-indexer/postgres';
import { DeleveragingEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractOrderFillHandler } from './abstract-order-fill-handler';

export class DeleveragingHandler extends AbstractOrderFillHandler<DeleveragingEventV1> {
  eventType: string = 'DeleveragingEvent';

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const liquidatedFill: FillFromDatabase = FillModel.fromJson(
      resultRow.liquidated_fill) as FillFromDatabase;
    const offsettingFill: FillFromDatabase = FillModel.fromJson(
      resultRow.offsetting_fill) as FillFromDatabase;
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      resultRow.perpetual_market) as PerpetualMarketFromDatabase;
    const liquidatedPerpetualPosition:
    PerpetualPositionFromDatabase = PerpetualPositionModel.fromJson(
      resultRow.liquidated_perpetual_position) as PerpetualPositionFromDatabase;
    const offsettingPerpetualPosition:
    PerpetualPositionFromDatabase = PerpetualPositionModel.fromJson(
      resultRow.offsetting_perpetual_position) as PerpetualPositionFromDatabase;
    const kafkaEvents: ConsolidatedKafkaEvent[] = [
      this.generateConsolidatedKafkaEvent(
        this.event.liquidated!,
        undefined,
        liquidatedPerpetualPosition,
        liquidatedFill,
        perpetualMarket,
      ),
      this.generateConsolidatedKafkaEvent(
        this.event.offsetting!,
        undefined,
        offsettingPerpetualPosition,
        offsettingFill,
        perpetualMarket,
      ),
      this.generateTradeKafkaEventFromTakerOrderFill(
        liquidatedFill,
      ),
    ];
    return kafkaEvents;
  }
}
