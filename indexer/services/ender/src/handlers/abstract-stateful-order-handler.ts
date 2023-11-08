import { logger } from '@dydxprotocol-indexer/base';
import {
  OrderFromDatabase,
  OrderStatus,
  OrderTable,
  OrderUpdateObject,
  OrderCreateObject,
  SubaccountTable,
  OrderSide,
  OrderType,
  protocolTranslations,
  PerpetualMarketFromDatabase,
  storeHelpers,
  OrderModel,
  PerpetualMarketModel,
  SubaccountFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import SubaccountModel from '@dydxprotocol-indexer/postgres/build/src/models/subaccount-model';
import {
  IndexerOrderId,
  IndexerOrder,
  IndexerOrder_Side,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { DateTime } from 'luxon';
import * as pg from 'pg';

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

  protected async handleEventViaSqlFunction():
  Promise<[OrderFromDatabase,
    PerpetualMarketFromDatabase,
    SubaccountFromDatabase | undefined]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_stateful_order_handler(
        ${this.block.height},
        '${this.block.time?.toISOString()}',
        '${JSON.stringify(StatefulOrderEventV1.decode(eventDataBinary))}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'AbstractStatefulOrderHandler#handleEventViaSqlFunction',
        message: 'Failed to handle StatefulOrderEventV1',
        error,
      });
      throw error;
    });

    return [
      OrderModel.fromJson(result.rows[0].result.order) as OrderFromDatabase,
      PerpetualMarketModel.fromJson(
        result.rows[0].result.perpetual_market) as PerpetualMarketFromDatabase,
      result.rows[0].result.subaccount
        ? SubaccountModel.fromJson(result.rows[0].result.subaccount) as SubaccountFromDatabase
        : undefined,
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
