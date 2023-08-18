import { logger } from '@dydxprotocol-indexer/base';
import {
  FillFromDatabase,
  Liquidity,
  OrderFromDatabase,
  OrderTable,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualPositionFromDatabase,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import {
  LiquidationOrderV1, IndexerOrderId,
} from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';

import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE, SUBACCOUNT_ORDER_FILL_EVENT_TYPE } from '../../constants';
import { convertPerpetualPosition } from '../../helpers/kafka-helper';
import {
  ConsolidatedKafkaEvent,
  OrderFillEventWithLiquidation,
  OrderFillEventWithLiquidity,
} from '../../lib/types';
import { AbstractOrderFillHandler, OrderFillEventBase } from './abstract-order-fill-handler';

export class LiquidationHandler extends AbstractOrderFillHandler<OrderFillEventWithLiquidity> {
  eventType: string = 'OrderFillEvent';

  /**
   * @returns the parallelizationIds for the this.event.liquidity order
   */
  public getParallelizationIds(): string[] {
    // OrderFillEvents with the same subaccountId and clobPairId cannot be processed in parallel.
    const liquidatedOrderFill:
    OrderFillEventWithLiquidation = this.event.event as OrderFillEventWithLiquidation;
    if (this.event.liquidity === Liquidity.MAKER) {
      const orderId: IndexerOrderId = liquidatedOrderFill.makerOrder!.orderId!;
      const orderUuid: string = OrderTable.orderIdToUuid(orderId);
      const subaccountUuid: string = SubaccountTable.subaccountIdToUuid(orderId.subaccountId!);
      return [
        `${this.eventType}_${subaccountUuid}_${orderId!.clobPairId}`,
        // To ensure that SubaccountUpdateEvents and OrderFillEvents for the same subaccount are not
        // processed in parallel
        `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${subaccountUuid}`,
        // To ensure that StatefulOrderEvents and OrderFillEvents for the same order are not
        // processed in parallel
        `${STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE}_${orderUuid}`,
      ];
    } else {
      const liquidationOrder: LiquidationOrderV1 = liquidatedOrderFill.liquidationOrder!;
      const subaccountUuid: string = SubaccountTable.subaccountIdToUuid(
        liquidationOrder.liquidated!,
      );
      return [
        `${this.eventType}_${subaccountUuid}_${liquidationOrder.clobPairId}`,
        // To ensure that SubaccountUpdateEvents and OrderFillEvents for the same subaccount are not
        // processed in parallel
        `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${subaccountUuid}`,
        // We do not need to add the StatefulOrderEvent parallelizationId here, because liquidation
        // fills have no order in postgres
      ];
    }
  }

  protected getTotalFilled(castedOrderFillEventMessage: OrderFillEventWithLiquidation): Long {
    return this.event.liquidity === Liquidity.TAKER
      ? castedOrderFillEventMessage.totalFilledTaker
      : castedOrderFillEventMessage.totalFilledMaker;
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const castedLiquidationFillEventMessage:
    OrderFillEventWithLiquidation = this.event.event as OrderFillEventWithLiquidation;
    const clobPairId:
    string = castedLiquidationFillEventMessage.makerOrder.orderId!.clobPairId.toString();
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(clobPairId);
    if (perpetualMarket === undefined) {
      logger.error({
        at: 'liquidationHandler#internalHandle',
        message: 'Unable to find perpetual market',
        clobPairId,
        castedLiquidationFillEventMessage,
      });
      throw new Error(`Unable to find perpetual market with clobPairId: ${clobPairId}`);
    }

    const orderFillBaseEventBase: OrderFillEventBase = this.createEventBaseFromLiquidation(
      castedLiquidationFillEventMessage,
      this.event.liquidity,
    );

    // Must be done in this order, because fills refer to an order
    // We do not create a taker order for liquidations.
    let makerOrder: OrderFromDatabase | undefined;
    if (this.event.liquidity === Liquidity.MAKER) {
      makerOrder = await this.runFuncWithTimingStatAndErrorLogging(
        this.upsertOrderFromEvent(
          perpetualMarket,
          castedLiquidationFillEventMessage.makerOrder,
          this.getTotalFilled(castedLiquidationFillEventMessage),
          false,
        ), this.generateTimingStatsOptions('upsert_maker_order'));
    }

    const fill: FillFromDatabase = await this.runFuncWithTimingStatAndErrorLogging(
      this.createFillFromEvent(perpetualMarket, orderFillBaseEventBase),
      this.generateTimingStatsOptions('create_fill'),
    );

    const position: PerpetualPositionFromDatabase = await this.runFuncWithTimingStatAndErrorLogging(
      this.updatePerpetualPosition(perpetualMarket, orderFillBaseEventBase),
      this.generateTimingStatsOptions('update_perpetual_position'),
    );

    if (this.event.liquidity === Liquidity.MAKER) {
      return [
        this.generateConsolidatedKafkaEvent(
          castedLiquidationFillEventMessage.makerOrder.orderId!.subaccountId!,
          makerOrder,
          convertPerpetualPosition(position),
          fill,
          perpetualMarket,
        ),
        // Update vulcan with the total filled amount of the maker order.
        this.getOrderUpdateKafkaEvent(
          castedLiquidationFillEventMessage.makerOrder!.orderId!,
          castedLiquidationFillEventMessage.totalFilledMaker,
        ),
      ];
    } else {
      return [
        this.generateConsolidatedKafkaEvent(
          castedLiquidationFillEventMessage.liquidationOrder.liquidated!,
          undefined,
          convertPerpetualPosition(position),
          fill,
          perpetualMarket,
        ),
        this.generateTradeKafkaEventFromTakerOrderFill(
          fill,
        ),
      ];
    }
  }
}
