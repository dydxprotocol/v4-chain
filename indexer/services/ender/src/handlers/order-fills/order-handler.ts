import { logger } from '@dydxprotocol-indexer/base';
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
  storeHelpers,
  SubaccountTable,
  USDC_ASSET_ID,
  OrderStatus,
} from '@dydxprotocol-indexer/postgres';
import { CanceledOrderStatus, CanceledOrdersCache, StateFilledQuantumsCache } from '@dydxprotocol-indexer/redis';
import { isStatefulOrder } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  OrderFillEventV1, IndexerOrderId, IndexerSubaccountId, IndexerOrder,
} from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';
import * as pg from 'pg';

import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE, SUBACCOUNT_ORDER_FILL_EVENT_TYPE } from '../../constants';
import { convertPerpetualPosition } from '../../helpers/kafka-helper';
import { redisClient } from '../../helpers/redis/redis-controller';
import { orderFillWithLiquidityToOrderFillEventWithOrder } from '../../helpers/translation-helper';
import { indexerTendermintEventToTransactionIndex } from '../../lib/helper';
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

  protected getTotalFilled(castedOrderFillEventMessage: OrderFillEventWithOrder): Long {
    return this.event.liquidity === Liquidity.TAKER
      ? castedOrderFillEventMessage.totalFilledTaker
      : castedOrderFillEventMessage.totalFilledMaker;
  }

  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const transactionIndex: number = indexerTendermintEventToTransactionIndex(
      this.indexerTendermintEvent,
    );
    const kafkaEvents: ConsolidatedKafkaEvent[] = [];

    const castedOrderFillEventMessage:
    OrderFillEventWithOrder = orderFillWithLiquidityToOrderFillEventWithOrder(this.event);
    const field: string = this.event.liquidity === Liquidity.TAKER ? 'order' : 'makerOrder';
    const orderProto: IndexerOrder = this.liquidityToOrder(
      castedOrderFillEventMessage,
      this.event.liquidity,
    );
    const orderUuid = OrderTable.orderIdToUuid(orderProto.orderId!);
    const canceledOrderStatus:
    CanceledOrderStatus = await CanceledOrdersCache.getOrderCanceledStatus(
      orderUuid,
      redisClient,
    );

    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_order_fill_handler_per_order(
        '${field}', 
        ${this.block.height}, 
        '${this.block.time?.toISOString()}', 
        '${JSON.stringify(OrderFillEventV1.decode(eventDataBinary))}', 
        ${this.indexerTendermintEvent.eventIndex}, 
        ${transactionIndex}, 
        '${this.block.txHashes[transactionIndex]}', 
        '${this.event.liquidity}', 
        'LIMIT',
        '${USDC_ASSET_ID}',
        '${canceledOrderStatus}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'orderHandler#internalHandle',
        message: 'Failed to handle OrderFillEventV1',
        error,
      });
      throw error;
    });
    const order: OrderFromDatabase = OrderModel.fromJson(
      result.rows[0].result.order) as OrderFromDatabase;
    const fill: FillFromDatabase = FillModel.fromJson(
      result.rows[0].result.fill) as FillFromDatabase;
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      result.rows[0].result.perpetual_market) as PerpetualMarketFromDatabase;
    const position: PerpetualPositionFromDatabase = PerpetualPositionModel.fromJson(
      result.rows[0].result.perpetual_position) as PerpetualPositionFromDatabase;

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
}
