import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE } from '../constants';
import { Handler } from './handler';

export abstract class AbstractStatefulOrderHandler<T> extends Handler<T> {
  public getParallelizationIdsFromOrderId(orderId: string): string[] {
    return [
      `${this.eventType}_${orderId}`,
      // To ensure that StatefulOrderEvents and OrderFillEvents for the same order are not
      // processed in parallel
      `${STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE}_${orderId}`,
    ];
  }
}
