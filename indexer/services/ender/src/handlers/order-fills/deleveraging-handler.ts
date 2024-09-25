import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  FillFromDatabase,
  FillModel,
  MarketFromDatabase,
  MarketModel,
  PerpetualMarketFromDatabase,
  PerpetualMarketModel,
  perpetualMarketRefresher,
  PerpetualPositionFromDatabase,
  PerpetualPositionModel,
  SubaccountTable,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import { DeleveragingEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../../config';
import { SUBACCOUNT_ORDER_FILL_EVENT_TYPE } from '../../constants';
import { annotateWithPnl, convertPerpetualPosition } from '../../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractOrderFillHandler } from './abstract-order-fill-handler';

export class DeleveragingHandler extends AbstractOrderFillHandler<DeleveragingEventV1> {
  eventType: string = 'DeleveragingEvent';

  public getParallelizationIds(): string[] {
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromId(this.event.perpetualId.toString());
    if (perpetualMarket === undefined) {
      logger.error({
        at: 'DeleveragingHandler#internalHandle',
        message: 'Unable to find perpetual market',
        perpetualId: this.event.perpetualId,
        event: this.event,
      });
      throw new Error(`Unable to find perpetual market with perpetualId: ${this.event.perpetualId}`);
    }
    const offsettingSubaccountUuid: string = SubaccountTable
      .uuid(this.event.offsetting!.owner, this.event.offsetting!.number);
    const deleveragedSubaccountUuid: string = SubaccountTable
      .uuid(this.event.liquidated!.owner, this.event.liquidated!.number);
    return [
      `${this.eventType}_${offsettingSubaccountUuid}_${perpetualMarket.clobPairId}`,
      `${this.eventType}_${deleveragedSubaccountUuid}_${perpetualMarket.clobPairId}`,
      // To ensure that SubaccountUpdateEvents, OrderFillEvents, and DeleveragingEvents for the same
      // subaccount are not processed in parallel
      `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${offsettingSubaccountUuid}`,
      `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${deleveragedSubaccountUuid}`,
    ];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const liquidatedFill: FillFromDatabase = FillModel.fromJson(
      resultRow.liquidated_fill) as FillFromDatabase;
    const offsettingFill: FillFromDatabase = FillModel.fromJson(
      resultRow.offsetting_fill) as FillFromDatabase;
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      resultRow.perpetual_market) as PerpetualMarketFromDatabase;
    const market: MarketFromDatabase = MarketModel.fromJson(
      resultRow.market) as MarketFromDatabase;
    const liquidatedPerpetualPosition:
    PerpetualPositionFromDatabase = PerpetualPositionModel.fromJson(
      resultRow.liquidated_perpetual_position) as PerpetualPositionFromDatabase;
    const offsettingPerpetualPosition:
    PerpetualPositionFromDatabase = PerpetualPositionModel.fromJson(
      resultRow.offsetting_perpetual_position) as PerpetualPositionFromDatabase;

    const liquidatedPositionUpdate: UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
      convertPerpetualPosition(liquidatedPerpetualPosition),
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      market,
    );
    const offsettingPositionUpdate: UpdatedPerpetualPositionSubaccountKafkaObject = annotateWithPnl(
      convertPerpetualPosition(offsettingPerpetualPosition),
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      market,
    );
    const kafkaEvents: ConsolidatedKafkaEvent[] = [
      this.generateConsolidatedKafkaEvent(
        this.event.liquidated!,
        undefined,
        liquidatedPositionUpdate,
        liquidatedFill,
        perpetualMarket,
      ),
      this.generateConsolidatedKafkaEvent(
        this.event.offsetting!,
        undefined,
        offsettingPositionUpdate,
        offsettingFill,
        perpetualMarket,
      ),
      this.generateTradeKafkaEventFromTakerOrderFill(
        liquidatedFill,
      ),
    ];
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_deleveraging_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );
    return kafkaEvents;
  }
}
