import { stats } from '@dydxprotocol-indexer/base';
import { OrderTable } from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../../config';
import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

// TODO(IND-334): Rename to LongTermOrderPlacementHandler after deprecating StatefulOrderPlacement
export class StatefulOrderPlacementHandler
  extends AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getOrderId(): string {
    let orderId: string;
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    if (this.event.orderPlace !== undefined) {
      orderId = OrderTable.orderIdToUuid(this.event.orderPlace!.order!.orderId!);
    } else if (this.event.twapOrderPlacement !== undefined) {
      orderId = OrderTable.orderIdToUuid(this.event.twapOrderPlacement!.order!.orderId!);
    } else {
      orderId = OrderTable.orderIdToUuid(this.event.longTermOrderPlacement!.order!.orderId!);
    }
    return orderId;
  }

  public getParallelizationIds(): string[] {
    // Stateful Order Events with the same orderId
    return this.getParallelizationIdsFromOrderId(this.getOrderId());
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    if (!resultRow) {
      return [];
    }

    let order: IndexerOrder;
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    if (this.event.orderPlace !== undefined) {
      order = this.event.orderPlace!.order!;
    } else if (this.event.twapOrderPlacement !== undefined) {
      order = this.event.twapOrderPlacement!.order!;
    } else {
      order = this.event.longTermOrderPlacement!.order!;
    }
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_stateful_order_placement_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );
    return this.createKafkaEvents(order);
  }

  private createKafkaEvents(
    order: IndexerOrder,
  ): ConsolidatedKafkaEvent[] {
    const kafkaEvents: ConsolidatedKafkaEvent[] = [];

    const offChainUpdate: OffChainUpdateV1 = OffChainUpdateV1.fromPartial({
      orderPlace: {
        order,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
      },
    });
    kafkaEvents.push(this.generateConsolidatedVulcanKafkaEvent(
      getOrderIdHash(order.orderId!),
      offChainUpdate,
      {
        message_received_timestamp: this.messageReceivedTimestamp,
        event_type: 'StatefulOrderPlacement',
      },
    ));
    return kafkaEvents;
  }
}
