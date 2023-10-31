import { logger } from '@dydxprotocol-indexer/base';
import {
  FillCreateObject,
  FillFromDatabase,
  FillTable,
  FillType,
  Liquidity,
  OrderSide,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualPositionColumns,
  PerpetualPositionFromDatabase,
  PerpetualPositionTable,
  protocolTranslations,
  SubaccountMessageContents,
  SubaccountTable,
  TendermintEventTable,
  TradeMessageContents,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import { DeleveragingEventV1, IndexerSubaccountId } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';

import { generateFillSubaccountMessage, generatePerpetualPositionsContents } from '../helpers/kafka-helper';
import {
  getWeightedAverage,
  indexerTendermintEventToTransactionIndex,
  perpetualPositionAndOrderSideMatching,
} from '../lib/helper';
import { ConsolidatedKafkaEvent, PriceFields, SumFields } from '../lib/types';
import { Handler } from './handler';

export class DeleveragingHandler extends Handler<DeleveragingEventV1> {
  eventType: string = 'DeleveragingEvent';

  public getParallelizationIds(): string[] {
    return [];
  }

  protected createFillsFromEvent(
    perpetualMarket: PerpetualMarketFromDatabase,
    event: DeleveragingEventV1,
  ): Promise<FillFromDatabase>[] {
    const eventId: Buffer = TendermintEventTable.createEventId(
      this.block.height.toString(),
      indexerTendermintEventToTransactionIndex(this.indexerTendermintEvent),
      this.indexerTendermintEvent.eventIndex,
    );
    const size: string = protocolTranslations.quantumsToHumanFixedString(
      event.fillAmount.toString(),
      perpetualMarket.atomicResolution,
    );
    const price: string = protocolTranslations.subticksToPrice(
      event.subticks.toString(10),
      perpetualMarket,
    );
    const transactionIndex: number = indexerTendermintEventToTransactionIndex(
      this.indexerTendermintEvent,
    );

    const liquidatedSubaccountFill: FillCreateObject = {
      subaccountId: SubaccountTable.uuid(event.liquidated!.owner, event.liquidated!.number),
      side: event.isBuy ? OrderSide.BUY : OrderSide.SELL,
      liquidity: Liquidity.TAKER,
      type: FillType.DELEVERAGED,
      clobPairId: event.clobPairId.toString(),
      size,
      price,
      quoteAmount: Big(size).times(price).toFixed(),
      eventId,
      transactionHash: this.block.txHashes[transactionIndex],
      createdAt: this.timestamp.toISO(),
      createdAtHeight: this.block.height.toString(),
      fee: '0',
    };
    const offsettingSubaccountFill: FillCreateObject = {
      ...liquidatedSubaccountFill,
      side: event.isBuy ? OrderSide.SELL : OrderSide.BUY,
      liquidity: Liquidity.MAKER,
      type: FillType.OFFSETTING,
    };

    return [
      FillTable.create(liquidatedSubaccountFill, { txId: this.txId }),
      FillTable.create(offsettingSubaccountFill, { txId: this.txId }),
    ];
  }

  protected async getLatestPerpetualPosition(
    perpetualMarket: PerpetualMarketFromDatabase,
    event: DeleveragingEventV1,
    deleveraged: boolean,
  ): Promise<PerpetualPositionFromDatabase> {
    const latestPerpetualPositions:
    PerpetualPositionFromDatabase[] = await PerpetualPositionTable.findAll(
      {
        subaccountId: deleveraged ? [
          SubaccountTable.uuid(event.liquidated!.owner, event.liquidated!.number),
        ] : [
          SubaccountTable.uuid(event.offsetting!.owner, event.offsetting!.number),
        ],
        perpetualId: [perpetualMarket.id],
        limit: 1,
      },
      [],
      { txId: this.txId },
    );

    if (latestPerpetualPositions.length === 0) {
      logger.error({
        at: 'deleveragingHandler#getLatestPerpetualPosition',
        message: 'Unable to find existing perpetual position.',
        blockHeight: this.block.height,
        clobPairId: event.clobPairId,
        subaccountId: deleveraged
          ? SubaccountTable.uuid(event.liquidated!.owner, event.liquidated!.number)
          : SubaccountTable.uuid(event.offsetting!.owner, event.offsetting!.number),
      });
      throw new Error('Unable to find existing perpetual position');
    }

    return latestPerpetualPositions[0];
  }

  protected getOrderSide(
    event: DeleveragingEventV1,
    deleveraged: boolean,
  ): OrderSide {
    if (deleveraged) {
      return event.isBuy ? OrderSide.BUY : OrderSide.SELL;
    }
    return event.isBuy ? OrderSide.SELL : OrderSide.BUY;
  }

  protected async updatePerpetualPosition(
    perpetualMarket: PerpetualMarketFromDatabase,
    event: DeleveragingEventV1,
    deleveraged: boolean,
  ): Promise<PerpetualPositionFromDatabase> {
    const latestPerpetualPosition:
    PerpetualPositionFromDatabase = await this.getLatestPerpetualPosition(
      perpetualMarket,
      event,
      deleveraged,
    );

    // update (sumOpen and entryPrice) or (sumClose and exitPrice)
    let sumField: SumFields;
    let priceField: PriceFields;
    if (perpetualPositionAndOrderSideMatching(
      latestPerpetualPosition.side,
      this.getOrderSide(event, deleveraged),
    )) {
      sumField = PerpetualPositionColumns.sumOpen;
      priceField = PerpetualPositionColumns.entryPrice;
    } else {
      sumField = PerpetualPositionColumns.sumClose;
      priceField = PerpetualPositionColumns.exitPrice;
    }

    const size: string = protocolTranslations.quantumsToHumanFixedString(
      event.fillAmount.toString(),
      perpetualMarket.atomicResolution,
    );
    const price: string = protocolTranslations.subticksToPrice(
      event.subticks.toString(10),
      perpetualMarket,
    );

    const updatedPerpetualPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.update(
      {
        id: latestPerpetualPosition.id,
        [sumField]: Big(latestPerpetualPosition[sumField]).plus(size).toFixed(),
        [priceField]: getWeightedAverage(
          latestPerpetualPosition[priceField] ?? '0',
          latestPerpetualPosition[sumField],
          price,
          size,
        ),
      },
      { txId: this.txId },
    );
    if (updatedPerpetualPosition === undefined) {
      logger.error({
        at: 'deleveragingHandler#handle',
        message: 'Unable to update perpetual position',
        latestPerpetualPositionId: latestPerpetualPosition.id,
        event,
      });
      throw new Error(`Unable to update perpetual position with id: ${latestPerpetualPosition.id}`);
    }
    return updatedPerpetualPosition;
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
    const clobPairId:
    string = this.event.clobPairId.toString();
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(clobPairId);
    if (perpetualMarket === undefined) {
      logger.error({
        at: 'DeleveragingHandler#internalHandle',
        message: 'Unable to find perpetual market',
        clobPairId,
        event: this.event,
      });
      throw new Error(`Unable to find perpetual market with clobPairId: ${clobPairId}`);
    }
    const fills: FillFromDatabase[] = await this.runFuncWithTimingStatAndErrorLogging(
      Promise.all(
        this.createFillsFromEvent(perpetualMarket, this.event),
      ),
      this.generateTimingStatsOptions('create_fill'),
    );

    const positions: PerpetualPositionFromDatabase[] = await
    this.runFuncWithTimingStatAndErrorLogging(
      Promise.all([
        this.updatePerpetualPosition(perpetualMarket, this.event, true),
        this.updatePerpetualPosition(perpetualMarket, this.event, false),
      ]),
      this.generateTimingStatsOptions('update_perpetual_position'),
    );
    const kafkaEvents: ConsolidatedKafkaEvent[] = [
      this.generateConsolidatedKafkaEvent(
        this.event.liquidated!,
        positions[0],
        fills[0],
        perpetualMarket,
      ),
      this.generateConsolidatedKafkaEvent(
        this.event.offsetting!,
        positions[1],
        fills[1],
        perpetualMarket,
      ),
      this.generateTradeKafkaEventFromDeleveraging(
        fills[0],
      ),
    ];
    return kafkaEvents;
  }
}
