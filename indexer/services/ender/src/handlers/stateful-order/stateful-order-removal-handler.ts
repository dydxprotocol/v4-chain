import { stats } from '@dydxprotocol-indexer/base';
import {
  OrderTable,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  OffChainUpdateV1,
  IndexerOrderId,
  OrderRemoveV1_OrderRemovalStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../../config';
import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

export class StatefulOrderRemovalHandler extends
  AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getParallelizationIds(): string[] {
    // Stateful Order Events with the same orderId
    const orderId: string = OrderTable.orderIdToUuid(this.event.orderRemoval!.removedOrderId!);
    return this.getParallelizationIdsFromOrderId(orderId);
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const orderIdProto: IndexerOrderId = this.event.orderRemoval!.removedOrderId!;
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_stateful_order_removal_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );
    return this.createKafkaEvents(orderIdProto);
  }

  private createKafkaEvents(orderIdProto: IndexerOrderId): ConsolidatedKafkaEvent[] {
    const offChainUpdate: OffChainUpdateV1 = OffChainUpdateV1.fromPartial({
      orderRemove: {
        removedOrderId: orderIdProto,
        reason: this.event.orderRemoval!.reason,
        removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED,
      },
    });

    return [
      this.generateConsolidatedVulcanKafkaEvent(
        getOrderIdHash(orderIdProto),
        offChainUpdate,
        {
          message_received_timestamp: this.messageReceivedTimestamp,
          event_type: 'StatefulOrderRemoval',
        },
      ),
    ];
  }
}
