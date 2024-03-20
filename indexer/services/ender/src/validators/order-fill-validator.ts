import {
  Liquidity,
  OrderTable,
} from '@dydxprotocol-indexer/postgres';
import { CanceledOrdersCache } from '@dydxprotocol-indexer/redis';
import {
  IndexerTendermintEvent,
  LiquidationOrderV1,
  IndexerOrder, OrderFillEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';

import { Handler, HandlerInitializer } from '../handlers/handler';
import { LiquidationHandler } from '../handlers/order-fills/liquidation-handler';
import { OrderHandler } from '../handlers/order-fills/order-handler';
import { redisClient } from '../helpers/redis/redis-controller';
import { orderFillEventV1ToOrderFill } from '../helpers/translation-helper';
import { OrderFillWithLiquidity } from '../lib/translated-types';
import { OrderFillEventWithLiquidity } from '../lib/types';
import { validateOrderAndReturnErrorMessage } from './helpers';
import { Validator } from './validator';

export class OrderFillValidator extends Validator<OrderFillEventV1> {
  public validate(): void {
    if (this.event.makerOrder === undefined) {
      return this.logAndThrowParseMessageError(
        'OrderFillEvent must contain a maker order',
        { event: this.event },
      );
    }

    this.validateOrder(this.event.makerOrder, Liquidity.MAKER);
    if (this.event.order) {
      this.validateOrder(this.event.order, Liquidity.TAKER);
    } else {
      this.validateLiquidationOrder(this.event.liquidationOrder!);
    }
  }

  public async getEventForBlockProcessor(): Promise<OrderFillEventV1> {
    // If event.order is populated then this means it is not a liquidation
    // order, and therefore we need to know the canceled order status stored
    // in redis to correctly update the database.
    if (this.event.order) {
      return Promise.all([
        CanceledOrdersCache.getOrderCanceledStatus(
          OrderTable.orderIdToUuid(this.event.makerOrder!.orderId!),
          redisClient,
        ),
        CanceledOrdersCache.getOrderCanceledStatus(
          OrderTable.orderIdToUuid(this.event.order.orderId!),
          redisClient,
        ),
      ],
      ).then((canceledOrderStatuses) => {
        return {
          makerCanceledOrderStatus: canceledOrderStatuses[0],
          takerCanceledOrderstatus: canceledOrderStatuses[1],
          ...this.event,
        };
      });
    }

    return this.event;
  }

  private validateOrder(
    order: IndexerOrder,
    liquidity: Liquidity,
  ): void {
    const orderName: string = liquidity === Liquidity.MAKER ? 'makerOrder' : 'takerOrder';

    const errorMessage: string | undefined = validateOrderAndReturnErrorMessage(order);
    if (errorMessage !== undefined) {
      return this.logAndThrowParseMessageError(
        `OrderFillEvent must contain a ${orderName}: ${errorMessage}`,
        { event: this.event },
      );
    }
  }

  private validateLiquidationOrder(
    liquidationOrder: LiquidationOrderV1,
  ): void {
    if (liquidationOrder.liquidated === undefined) {
      return this.logAndThrowParseMessageError(
        'LiquidationOrder must contain a liquidated subaccountId',
        { event: this.event },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    __: string,
  ): Handler<OrderFillWithLiquidity>[] {
    const orderFillEventsWithLiquidity: OrderFillEventWithLiquidity[] = [
      {
        event: this.event,
        liquidity: Liquidity.MAKER,
      },
      {
        event: this.event,
        liquidity: Liquidity.TAKER,
      },
    ];

    const Initializer:
    HandlerInitializer | undefined = this.event.order === undefined
      ? LiquidationHandler : OrderHandler;

    return _.map(
      orderFillEventsWithLiquidity,
      (orderFillEventWithLiquidity: OrderFillEventWithLiquidity) => {
        return new Initializer(
          this.block,
          this.blockEventIndex,
          indexerTendermintEvent,
          txId,
          orderFillEventV1ToOrderFill(orderFillEventWithLiquidity),
        );
      },
    );
  }
}
