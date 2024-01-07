import {
  FillFromDatabase,
  FillModel,
  Liquidity,
  OrderFromDatabase,
  OrderModel,
  PerpetualMarketFromDatabase,
  PerpetualMarketModel,
  PerpetualPositionFromDatabase,
  PerpetualPositionModel,
  OrderStatus,
} from '@dydxprotocol-indexer/postgres';
import { StateFilledQuantumsCache } from '@dydxprotocol-indexer/redis';
import { isStatefulOrder } from '@dydxprotocol-indexer/v4-proto-parser';
import Long from 'long';
import * as pg from 'pg';

import { convertPerpetualPosition } from '../../helpers/kafka-helper';
import { redisClient } from '../../helpers/redis/redis-controller';
import {
  orderFillWithLiquidityToOrderFillEventWithLiquidation,
} from '../../helpers/translation-helper';
import { OrderFillWithLiquidity } from '../../lib/translated-types';
import {
  ConsolidatedKafkaEvent,
  OrderFillEventWithLiquidation,
} from '../../lib/types';
import { AbstractOrderFillHandler } from './abstract-order-fill-handler';

export class LiquidationHandler extends AbstractOrderFillHandler<OrderFillWithLiquidity> {
  eventType: string = 'LiquidationEvent';

  protected getTotalFilled(castedOrderFillEventMessage: OrderFillEventWithLiquidation): Long {
    return this.event.liquidity === Liquidity.TAKER
      ? castedOrderFillEventMessage.totalFilledTaker
      : castedOrderFillEventMessage.totalFilledMaker;
  }

  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const castedLiquidationFillEventMessage:
    OrderFillEventWithLiquidation = orderFillWithLiquidityToOrderFillEventWithLiquidation(
      this.event,
    );
    const field: string = this.event.liquidity === Liquidity.MAKER
      ? 'makerOrder' : 'liquidationOrder';

    const fill: FillFromDatabase = FillModel.fromJson(
      resultRow[field].fill) as FillFromDatabase;
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      resultRow[field].perpetual_market) as PerpetualMarketFromDatabase;
    const position: PerpetualPositionFromDatabase = PerpetualPositionModel.fromJson(
      resultRow[field].perpetual_position) as PerpetualPositionFromDatabase;

    if (this.event.liquidity === Liquidity.MAKER) {
      // Must be done in this order, because fills refer to an order
      // We do not create a taker order for liquidations.
      const makerOrder: OrderFromDatabase = OrderModel.fromJson(
        resultRow[field].order) as OrderFromDatabase;

      // Update the cache tracking the state-filled amount per order for use in vulcan
      await StateFilledQuantumsCache.updateStateFilledQuantums(
        makerOrder!.id,
        this.getTotalFilled(castedLiquidationFillEventMessage).toString(),
        redisClient,
      );

      const kafkaEvents: ConsolidatedKafkaEvent[] = [
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

      // If the order is stateful and fully-filled, send an order removal to vulcan. We only do this
      // for stateful orders as we are guaranteed a stateful order cannot be replaced until the next
      // block.
      if (makerOrder?.status === OrderStatus.FILLED && isStatefulOrder(makerOrder?.orderFlags)) {
        kafkaEvents.push(
          this.getOrderRemoveKafkaEvent(castedLiquidationFillEventMessage.makerOrder!.orderId!),
        );
      }
      return kafkaEvents;
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
