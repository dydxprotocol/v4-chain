import {
  FillFromDatabase,
  FillModel,
  Liquidity,
  OrderFromDatabase,
  OrderModel,
  OrderTable,
  PerpetualMarketFromDatabase,
  PerpetualMarketModel,
  PerpetualPositionFromDatabase,
  PerpetualPositionModel,
  SubaccountTable,
  OrderStatus,
} from '@dydxprotocol-indexer/postgres';
import { StateFilledQuantumsCache } from '@dydxprotocol-indexer/redis';
import { isStatefulOrder } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrderId, IndexerSubaccountId, IndexerOrder,
} from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';
import * as pg from 'pg';

import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE, SUBACCOUNT_ORDER_FILL_EVENT_TYPE } from '../../constants';
import { convertPerpetualPosition } from '../../helpers/kafka-helper';
import { redisClient } from '../../helpers/redis/redis-controller';
import { orderFillWithLiquidityToOrderFillEventWithOrder } from '../../helpers/translation-helper';
import { OrderFillWithLiquidity } from '../../lib/translated-types';
import { ConsolidatedKafkaEvent, OrderFillEventWithOrder } from '../../lib/types';
import { AbstractOrderFillHandler } from './abstract-order-fill-handler';

export class OrderHandler extends AbstractOrderFillHandler<OrderFillWithLiquidity> {
  eventType: string = 'OrderFillEvent';

  /**
   * @returns the parallelizationIds for the this.event.liquidity order
   */
  public getParallelizationIds(): string[] {
    // OrderFillEvents with the same subaccountId and clobPairId cannot be processed in parallel.
    const castedOrderFillEventMessage:
    OrderFillEventWithOrder = orderFillWithLiquidityToOrderFillEventWithOrder(this.event);
    const orderId: IndexerOrderId = this.event.liquidity === Liquidity.MAKER
      ? castedOrderFillEventMessage.makerOrder!.orderId!
      : castedOrderFillEventMessage.order!.orderId!;
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
  }

  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const kafkaEvents: ConsolidatedKafkaEvent[] = [];

    const castedOrderFillEventMessage:
    OrderFillEventWithOrder = orderFillWithLiquidityToOrderFillEventWithOrder(this.event);
    const field: string = this.event.liquidity === Liquidity.TAKER ? 'order' : 'makerOrder';
    const orderProto: IndexerOrder = this.liquidityToOrder(
      castedOrderFillEventMessage,
      this.event.liquidity,
    );
    const order: OrderFromDatabase = OrderModel.fromJson(
      resultRow[field].order) as OrderFromDatabase;
    const fill: FillFromDatabase = FillModel.fromJson(
      resultRow[field].fill) as FillFromDatabase;
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      resultRow[field].perpetual_market) as PerpetualMarketFromDatabase;
    const position: PerpetualPositionFromDatabase = PerpetualPositionModel.fromJson(
      resultRow[field].perpetual_position) as PerpetualPositionFromDatabase;

    let subaccountId: IndexerSubaccountId;
    if (this.event.liquidity === Liquidity.MAKER) {
      subaccountId = castedOrderFillEventMessage.makerOrder.orderId!.subaccountId!;
    } else {
      subaccountId = castedOrderFillEventMessage.order.orderId!.subaccountId!;
    }
    kafkaEvents.push(
      this.generateConsolidatedKafkaEvent(
        subaccountId,
        order,
        convertPerpetualPosition(position),
        fill,
        perpetualMarket,
      ),
    );

    // Update vulcan with the total filled amount of the order.
    kafkaEvents.push(
      this.getOrderUpdateKafkaEvent(
        orderProto.orderId!,
        this.getTotalFilled(castedOrderFillEventMessage),
      ),
    );

    // Update the cache tracking the state-filled amount per order for use in vulcan
    await StateFilledQuantumsCache.updateStateFilledQuantums(
      order.id,
      this.getTotalFilled(castedOrderFillEventMessage).toString(),
      redisClient,
    );

    // If the order is stateful and fully-filled, send an order removal to vulcan. We only do this
    // for stateful orders as we are guaranteed a stateful order cannot be replaced until the next
    // block.
    if (order.status === OrderStatus.FILLED && isStatefulOrder(order.orderFlags)) {
      kafkaEvents.push(this.getOrderRemoveKafkaEvent(orderProto.orderId!));
    }

    if (this.event.liquidity === Liquidity.TAKER) {
      kafkaEvents.push(this.generateTradeKafkaEventFromTakerOrderFill(fill));
      return kafkaEvents;
    }

    return kafkaEvents;
  }

  protected getTotalFilled(castedOrderFillEventMessage: OrderFillEventWithOrder): Long {
    return this.event.liquidity === Liquidity.TAKER
      ? castedOrderFillEventMessage.totalFilledTaker
      : castedOrderFillEventMessage.totalFilledMaker;
  }
}
