import { stats } from '@dydxprotocol-indexer/base';
import {
  OrderFromDatabase, OrderModel,
  OrderTable,
  PerpetualMarketFromDatabase, PerpetualMarketModel,
  SubaccountMessageContents,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerSubaccountId,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../../config';
import { generateOrderSubaccountMessage } from '../../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

export class ConditionalOrderPlacementHandler extends
  AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getParallelizationIds(): string[] {
    const orderId: string = OrderTable.orderIdToUuid(
      this.event.conditionalOrderPlacement!.order!.orderId!,
    );
    return this.getParallelizationIdsFromOrderId(orderId);
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const order: OrderFromDatabase = OrderModel.fromJson(resultRow.order) as OrderFromDatabase;
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      resultRow.perpetual_market) as PerpetualMarketFromDatabase;

    const subaccountId:
    IndexerSubaccountId = this.event.conditionalOrderPlacement!.order!.orderId!.subaccountId!;
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_conditional_order_placement_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );
    return this.createKafkaEvents(subaccountId, order, perpetualMarket);
  }

  private createKafkaEvents(
    subaccountId: IndexerSubaccountId,
    conditionalOrder: OrderFromDatabase,
    perpetualMarket: PerpetualMarketFromDatabase): ConsolidatedKafkaEvent[] {

    // Since the order isn't placed on the book, no message is sent to vulcan
    // ender needs to send the websocket message indicating the conditional order was placed
    const message: SubaccountMessageContents = {
      orders: [
        generateOrderSubaccountMessage(conditionalOrder, perpetualMarket.ticker),
      ],
      blockHeight: this.block.height.toString(),
    };

    return [
      this.generateConsolidatedSubaccountKafkaEvent(
        JSON.stringify(message),
        subaccountId,
      ),
    ];
  }
}
