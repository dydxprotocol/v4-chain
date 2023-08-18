import { logger } from '@dydxprotocol-indexer/base';
import {
  OrderFromDatabase, OrderStatus, OrderTable, OrderUpdateObject,
} from '@dydxprotocol-indexer/postgres';
import { IndexerOrderId } from '@dydxprotocol-indexer/v4-protos';

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

  protected async cancelOrder(
    orderIdProto: IndexerOrderId,
  ): Promise<OrderFromDatabase> {
    const orderId = OrderTable.orderIdToUuid(orderIdProto);
    const orderUpdateObject: OrderUpdateObject = {
      id: orderId,
      status: OrderStatus.CANCELED,
    };

    const order: OrderFromDatabase | undefined = await OrderTable.update(
      orderUpdateObject,
      { txId: this.txId },
    );
    if (order === undefined) {
      const message: string = `Unable to cancel order with orderId: ${orderId}`;
      logger.error({
        at: 'AbstractStatefulOrderHandler#cancelOrder',
        message,
      });
      throw new Error(message);
    }
    return order;
  }
}
