import { logger } from '@dydxprotocol-indexer/base';
import {
  OrderFromDatabase,
  PerpetualMarketFromDatabase,
  storeHelpers,
  OrderModel,
  PerpetualMarketModel,
  SubaccountFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import SubaccountModel from '@dydxprotocol-indexer/postgres/build/src/models/subaccount-model';
import {
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

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
}
