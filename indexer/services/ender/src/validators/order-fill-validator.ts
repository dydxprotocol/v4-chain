import { Liquidity } from '@dydxprotocol-indexer/postgres';
import {
  IndexerTendermintEvent,
  LiquidationOrderV1,
  IndexerOrder, OrderFillEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';

import { Handler, HandlerInitializer } from '../handlers/handler';
import { LiquidationHandler } from '../handlers/order-fills/liquidation-handler';
import { OrderHandler } from '../handlers/order-fills/order-handler';
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
          indexerTendermintEvent,
          txId,
          orderFillEventV1ToOrderFill(orderFillEventWithLiquidity),
        );
      },
    );
  }
}
