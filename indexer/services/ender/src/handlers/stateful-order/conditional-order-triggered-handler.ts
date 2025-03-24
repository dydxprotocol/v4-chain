import { stats } from '@dydxprotocol-indexer/base';
import {
  OrderFromDatabase,
  OrderTable,
  PerpetualMarketFromDatabase,
  orderTranslations,
  SubaccountFromDatabase, OrderModel, PerpetualMarketModel,
} from '@dydxprotocol-indexer/postgres';
import SubaccountModel from '@dydxprotocol-indexer/postgres/build/src/models/subaccount-model';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../../config';
import { sendOrderTriggeredNotification } from '../../helpers/notifications/notifications-functions';
import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

export class ConditionalOrderTriggeredHandler extends
  AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getParallelizationIds(): string[] {
    const orderId: string = OrderTable.orderIdToUuid(
      this.event.conditionalOrderTriggered!.triggeredOrderId!,
    );
    return this.getParallelizationIdsFromOrderId(orderId);
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const order: OrderFromDatabase = OrderModel.fromJson(resultRow.order) as OrderFromDatabase;
    const perpetualMarket: PerpetualMarketFromDatabase = PerpetualMarketModel.fromJson(
      resultRow.perpetual_market) as PerpetualMarketFromDatabase;
    const subaccount: SubaccountFromDatabase = SubaccountModel.fromJson(
      resultRow.subaccount) as SubaccountFromDatabase;

    const indexerOrder: IndexerOrder = orderTranslations.convertToIndexerOrderWithSubaccount(
      order, perpetualMarket, subaccount);
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_conditional_order_triggered_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );
    await sendOrderTriggeredNotification(order, perpetualMarket, subaccount);
    return this.createKafkaEvents(indexerOrder);
  }

  private createKafkaEvents(order: IndexerOrder): ConsolidatedKafkaEvent[] {
    const offChainUpdate: OffChainUpdateV1 = OffChainUpdateV1.fromPartial({
      orderPlace: {
        order,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
      },
    });

    return [
      this.generateConsolidatedVulcanKafkaEvent(
        getOrderIdHash(order.orderId!),
        offChainUpdate,
        {
          message_received_timestamp: this.messageReceivedTimestamp,
          event_type: 'ConditionalOrderTriggered',
        },
      ),
    ];
  }
}
