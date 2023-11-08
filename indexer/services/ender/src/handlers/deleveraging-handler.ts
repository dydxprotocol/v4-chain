import { logger } from '@dydxprotocol-indexer/base';
import {
  FillFromDatabase,
  FillModel,
  PerpetualMarketFromDatabase,
  PerpetualMarketModel,
  perpetualMarketRefresher,
  PerpetualPositionFromDatabase,
  PerpetualPositionModel,
  storeHelpers,
  SubaccountMessageContents,
  SubaccountTable,
  TradeMessageContents,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import { DeleveragingEventV1, IndexerSubaccountId } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { DELEVERAGING_EVENT_TYPE, SUBACCOUNT_ORDER_FILL_EVENT_TYPE } from '../constants';
import { generateFillSubaccountMessage, generatePerpetualPositionsContents } from '../helpers/kafka-helper';
import { indexerTendermintEventToTransactionIndex } from '../lib/helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class DeleveragingHandler extends Handler<DeleveragingEventV1> {
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
      // To ensure that SubaccountUpdateEvents and OrderFillEvents for the same subaccount are not
      // processed in parallel
      `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${offsettingSubaccountUuid}`,
      `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${deleveragedSubaccountUuid}`,
      // To ensure that StatefulOrderEvents and OrderFillEvents for the same order are not
      // processed in parallel
      `${DELEVERAGING_EVENT_TYPE}_${offsettingSubaccountUuid}`,
      `${DELEVERAGING_EVENT_TYPE}_${deleveragedSubaccountUuid}`,
    ];
  }

  protected generateConsolidatedKafkaEvent(
    subaccountIdProto: IndexerSubaccountId,
    position: UpdatedPerpetualPositionSubaccountKafkaObject | undefined,
    fill: FillFromDatabase,
    perpetualMarket: PerpetualMarketFromDatabase,
  ): ConsolidatedKafkaEvent {
    const message: SubaccountMessageContents = {
      fills: [
        generateFillSubaccountMessage(fill, perpetualMarket.ticker),
      ],
      perpetualPositions: position === undefined ? undefined : generatePerpetualPositionsContents(
        subaccountIdProto,
        [position],
        perpetualMarketRefresher.getPerpetualMarketsMap(),
      ),
    };
    return this.generateConsolidatedSubaccountKafkaEvent(
      JSON.stringify(message),
      subaccountIdProto,
      undefined,
      true,
      message,
    );
  }

  protected generateTradeKafkaEventFromDeleveraging(
    fill: FillFromDatabase,
  ): ConsolidatedKafkaEvent {
    const tradeContents: TradeMessageContents = {
      trades: [
        {
          id: fill.eventId.toString('hex'),
          size: fill.size,
          price: fill.price,
          side: fill.side.toString(),
          createdAt: fill.createdAt,
          liquidation: false,
          deleveraging: true,
        },
      ],
    };
    return this.generateConsolidatedTradeKafkaEvent(
      JSON.stringify(tradeContents),
      fill.clobPairId,
    );
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const transactionIndex: number = indexerTendermintEventToTransactionIndex(
      this.indexerTendermintEvent,
    );
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_deleveraging_handler(
        ${this.block.height}, 
        '${this.block.time?.toISOString()}', 
        '${JSON.stringify(DeleveragingEventV1.decode(eventDataBinary))}', 
        ${this.indexerTendermintEvent.eventIndex}, 
        ${transactionIndex}, 
        '${this.block.txHashes[transactionIndex]}' 
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'deleveragingHandler#handleViaSqlFunction',
        message: 'Failed to handle DeleveragingEventV1',
        error,
      });
      throw error;
    });
    const liquidatedFill: FillFromDatabase = FillModel.fromJson(
      result.rows[0].result.liquidated_fill) as FillFromDatabase;
    const offsettingFill: FillFromDatabase = FillModel.fromJson(
      result.rows[0].result.offsetting_fill) as FillFromDatabase;
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      result.rows[0].result.perpetual_market) as PerpetualMarketFromDatabase;
    const liquidatedPerpetualPosition:
    PerpetualPositionFromDatabase = PerpetualPositionModel.fromJson(
      result.rows[0].result.liquidated_perpetual_position) as PerpetualPositionFromDatabase;
    const offsettingPerpetualPosition:
    PerpetualPositionFromDatabase = PerpetualPositionModel.fromJson(
      result.rows[0].result.offsetting_perpetual_position) as PerpetualPositionFromDatabase;
    const kafkaEvents: ConsolidatedKafkaEvent[] = [
      this.generateConsolidatedKafkaEvent(
        this.event.liquidated!,
        liquidatedPerpetualPosition,
        liquidatedFill,
        perpetualMarket,
      ),
      this.generateConsolidatedKafkaEvent(
        this.event.offsetting!,
        offsettingPerpetualPosition,
        offsettingFill,
        perpetualMarket,
      ),
      this.generateTradeKafkaEventFromDeleveraging(
        liquidatedFill,
      ),
    ];
    return kafkaEvents;
  }
}
