import {
  IndexerTendermintEvent,
  IndexerOrder,
  StatefulOrderEventV1,
  StatefulOrderEventV1_StatefulOrderPlacementV1,
  OrderRemovalReason,
  StatefulOrderEventV1_StatefulOrderRemovalV1,
} from '@dydxprotocol-indexer/v4-protos';

import { Handler, HandlerInitializer } from '../handlers/handler';
import { StatefulOrderPlacementHandler } from '../handlers/stateful-order/stateful-order-placement-handler';
import { StatefulOrderRemovalHandler } from '../handlers/stateful-order/stateful-order-removal-handler';
import { validateOrderAndReturnErrorMessage, validateOrderIdAndReturnErrorMessage } from './helpers';
import { Validator } from './validator';

export class StatefulOrderValidator extends Validator<StatefulOrderEventV1> {
  public validate(): void {
    if (
      this.event.orderPlace === undefined &&
      this.event.orderRemoval === undefined
    ) {
      return this.logAndThrowParseMessageError(
        'One of orderPlace or orderRemoval must be defined in StatefulOrderEvent',
        { event: this.event },
      );
    }
    if (this.event.orderPlace !== undefined) {
      this.validateOrderPlace(this.event.orderPlace);
    } else { // orderRemoval
      this.validateOrderRemoval(this.event.orderRemoval!);
    }
  }

  private validateOrderPlace(
    orderPlace: StatefulOrderEventV1_StatefulOrderPlacementV1,
  ): void {
    const order: IndexerOrder | undefined = orderPlace.order;
    if (order === undefined) {
      return this.logAndThrowParseMessageError(
        'StatefulOrderEvent placement must contain an order',
        { event: this.event },
      );
    }

    const orderErrorMessage: string | undefined = validateOrderAndReturnErrorMessage(order);
    if (orderErrorMessage !== undefined) {
      return this.logAndThrowParseMessageError(
        `StatefulOrderEvent placement: ${orderErrorMessage}`,
        { event: this.event },
      );
    }

    if (order.goodTilBlockTime === undefined) {
      return this.logAndThrowParseMessageError(
        'StatefulOrderEvent placement: order must have goodTilBlockTime',
        { event: this.event },
      );
    }
  }

  private validateOrderRemoval(
    orderRemoval: StatefulOrderEventV1_StatefulOrderRemovalV1,
  ): void {
    if (orderRemoval.removedOrderId === undefined) {
      return this.logAndThrowParseMessageError(
        'StatefulOrderEvent removal must contain an orderId',
        { event: this.event },
      );
    }

    if (orderRemoval.reason === OrderRemovalReason.ORDER_REMOVAL_REASON_UNSPECIFIED) {
      return this.logAndThrowParseMessageError(
        'StatefulOrderEvent removal must contain a valid reason',
        { event: this.event },
      );
    }

    const orderIdErrorMessage: string | undefined = validateOrderIdAndReturnErrorMessage(
      orderRemoval.removedOrderId,
    );
    if (orderIdErrorMessage !== undefined) {
      return this.logAndThrowParseMessageError(
        `StatefulOrderEvent removal ${orderIdErrorMessage}`,
        { event: this.event },
      );
    }
  }

  public getHandlerInitializer() : HandlerInitializer | undefined {
    if (this.event.orderPlace !== undefined) {
      return StatefulOrderPlacementHandler;
    } else if (this.event.orderRemoval !== undefined) {
      return StatefulOrderRemovalHandler;
    }
    return undefined;
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
  ): Handler<StatefulOrderEventV1>[] {
    const Initializer:
    HandlerInitializer | undefined = this.getHandlerInitializer();
    if (Initializer === undefined) {
      this.logAndThrowParseMessageError(
        'Cannot process event',
        { event: this.event },
      );
    }
    // @ts-ignore
    const handler: Handler<StatefulOrderEvent> = new Initializer(
      this.block,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
