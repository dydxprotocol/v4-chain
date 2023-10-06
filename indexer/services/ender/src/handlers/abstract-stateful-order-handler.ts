import { logger } from '@dydxprotocol-indexer/base';
import {
  OrderFromDatabase, OrderStatus, OrderTable, OrderUpdateObject, OrderCreateObject, SubaccountTable,
  OrderSide, OrderType, protocolTranslations,
  PerpetualMarketFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { IndexerOrderId, IndexerOrder, IndexerOrder_Side } from '@dydxprotocol-indexer/v4-protos';
import { DateTime } from 'luxon';

import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE } from '../constants';
import { getPrice, getSize } from '../lib/helper';
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

  protected async updateOrderStatus(
    orderIdProto: IndexerOrderId,
    status: OrderStatus,
  ): Promise<OrderFromDatabase> {
    const orderId = OrderTable.orderIdToUuid(orderIdProto);
    const orderUpdateObject: OrderUpdateObject = {
      id: orderId,
      status,
      updatedAt: DateTime.fromJSDate(this.block.time!).toISO(),
      updatedAtHeight: this.block.height.toString(),
    };

    const order: OrderFromDatabase | undefined = await OrderTable.update(
      orderUpdateObject,
      { txId: this.txId },
    );
    if (order === undefined) {
      const message: string = `Unable to update order status with orderId: ${orderId}`;
      logger.error({
        at: 'AbstractStatefulOrderHandler#cancelOrder',
        message,
        status,
      });
      throw new Error(message);
    }
    return order;
  }

  /**
   * Upsert order to database, because there may be an existing order with the orderId in the
   * database.
   */
  // eslint-disable-next-line @typescript-eslint/require-await
  protected async upsertOrder(
    perpetualMarket: PerpetualMarketFromDatabase,
    order: IndexerOrder,
    type: OrderType,
    status: OrderStatus,
    triggerPrice?: string,
  ): Promise<OrderFromDatabase> {
    const size: string = getSize(order, perpetualMarket);
    const price: string = getPrice(order, perpetualMarket);

    const orderToCreate: OrderCreateObject = {
      subaccountId: SubaccountTable.subaccountIdToUuid(order.orderId!.subaccountId!),
      clientId: order.orderId!.clientId.toString(),
      clobPairId: order.orderId!.clobPairId.toString(),
      side: order.side === IndexerOrder_Side.SIDE_BUY ? OrderSide.BUY : OrderSide.SELL,
      size,
      totalFilled: '0',
      price,
      type,
      status,
      timeInForce: protocolTranslations.protocolOrderTIFToTIF(order.timeInForce),
      reduceOnly: order.reduceOnly,
      orderFlags: order.orderId!.orderFlags.toString(),
      // On chain orders must have a goodTilBlockTime rather than a goodTilBlock
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(order),
      createdAtHeight: this.block.height.toString(),
      clientMetadata: order.clientMetadata.toString(),
      triggerPrice,
      updatedAt: DateTime.fromJSDate(this.block.time!).toISO(),
      updatedAtHeight: this.block.height.toString(),
    };

    return OrderTable.upsert(orderToCreate, { txId: this.txId });
  }
}
