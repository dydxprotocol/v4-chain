import {
  IndexerTendermintEvent,
  IndexerOrder,
  StatefulOrderEventV1,
  StatefulOrderEventV1_StatefulOrderPlacementV1,
  OrderRemovalReason,
  StatefulOrderEventV1_StatefulOrderRemovalV1,
  StatefulOrderEventV1_ConditionalOrderPlacementV1,
  StatefulOrderEventV1_ConditionalOrderTriggeredV1,
  StatefulOrderEventV1_LongTermOrderPlacementV1,
} from '@dydxprotocol-indexer/v4-protos';

import { Handler, HandlerInitializer } from '../handlers/handler';
import { ConditionalOrderPlacementHandler } from '../handlers/stateful-order/conditional-order-placement-handler';
import { ConditionalOrderTriggeredHandler } from '../handlers/stateful-order/conditional-order-triggered-handler';
import { StatefulOrderPlacementHandler } from '../handlers/stateful-order/stateful-order-placement-handler';
import { StatefulOrderRemovalHandler } from '../handlers/stateful-order/stateful-order-removal-handler';
import { validateOrderAndReturnErrorMessage, validateOrderIdAndReturnErrorMessage } from './helpers';
import { Validator } from './validator';

export class StatefulOrderValidator extends Validator<StatefulOrderEventV1> {
  public validate(): void {
    if (
      this.event.orderPlace === undefined &&
      this.event.orderRemoval === undefined &&
      this.event.conditionalOrderPlacement === undefined &&
      this.event.conditionalOrderTriggered === undefined &&
      this.event.longTermOrderPlacement === undefined
    ) {
      return this.logAndThrowParseMessageError(
        'One of orderPlace, orderRemoval, conditionalOrderPlacement, conditionalOrderTriggered, ' +
        'longTermOrderPlacement must be defined in StatefulOrderEvent',
        { event: this.event },
      );
    }
    if (this.event.orderPlace !== undefined) {
      this.validateOrderPlace(this.event.orderPlace);
    } else if (this.event.orderRemoval !== undefined) {
      this.validateOrderRemoval(this.event.orderRemoval!);
    } else if (this.event.conditionalOrderPlacement !== undefined) {
      this.validateConditionalOrderPlacement(this.event.conditionalOrderPlacement);
    } else if (this.event.conditionalOrderTriggered !== undefined) {
      this.validateConditionalOrderTriggered(this.event.conditionalOrderTriggered);
    } else { // longTermOrderPlacement
      this.validateLongTermOrderPlacement(this.event.longTermOrderPlacement!);
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

  private validateConditionalOrderPlacement(
    _conditionalOrderPlacement: StatefulOrderEventV1_ConditionalOrderPlacementV1,
  ): void {
    // TODO(IND-334): Implement validation logic
  }

  private validateConditionalOrderTriggered(
    _conditionalOrderTriggered: StatefulOrderEventV1_ConditionalOrderTriggeredV1,
  ): void {
    // TODO(IND-334): Implement validation logic
  }

  private validateLongTermOrderPlacement(
    longTermOrderPlacement: StatefulOrderEventV1_LongTermOrderPlacementV1,
  ): void {
    const order: IndexerOrder | undefined = longTermOrderPlacement.order;
    if (order === undefined) {
      return this.logAndThrowParseMessageError(
        'StatefulOrderEvent long term order placement must contain an order',
        { event: this.event },
      );
    }

    const orderErrorMessage: string | undefined = validateOrderAndReturnErrorMessage(order);
    if (orderErrorMessage !== undefined) {
      return this.logAndThrowParseMessageError(
        `StatefulOrderEvent long term order placement: ${orderErrorMessage}`,
        { event: this.event },
      );
    }

    if (order.goodTilBlockTime === undefined) {
      return this.logAndThrowParseMessageError(
        'StatefulOrderEvent long term order placement: order must have goodTilBlockTime',
        { event: this.event },
      );
    }
  }

  public getHandlerInitializer() : HandlerInitializer | undefined {
    if (this.event.orderPlace !== undefined) {
      return StatefulOrderPlacementHandler;
    } else if (this.event.orderRemoval !== undefined) {
      return StatefulOrderRemovalHandler;
    } else if (this.event.conditionalOrderPlacement !== undefined) {
      return ConditionalOrderPlacementHandler;
    } else if (this.event.conditionalOrderTriggered !== undefined) {
      return ConditionalOrderTriggeredHandler;
    } else if (this.event.longTermOrderPlacement !== undefined) {
      return StatefulOrderPlacementHandler;
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
