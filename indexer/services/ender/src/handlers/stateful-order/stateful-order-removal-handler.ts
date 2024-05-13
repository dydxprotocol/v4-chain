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

import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';
import { logger, stats } from '@dydxprotocol-indexer/base';

export class StatefulOrderRemovalHandler extends
  AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getParallelizationIds(): string[] {
    // Stateful Order Events with the same orderId
    const orderId: string = OrderTable.orderIdToUuid(this.event.orderRemoval!.removedOrderId!);
    return this.getParallelizationIdsFromOrderId(orderId);
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(_: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const orderIdProto: IndexerOrderId = this.event.orderRemoval!.removedOrderId!;
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

    const messageKey: Buffer = Buffer.from(orderIdProto.clobPairId.toString());
    logger.info({
      at: 'handlers#stateful-order-removal',
      message: `Clob pair ID ${orderIdProto.clobPairId}`,
    });

    return [
      this.generateConsolidatedVulcanKafkaEvent(
        messageKey,
        offChainUpdate,
        {
          message_received_timestamp: this.messageReceivedTimestamp,
          event_type: 'StatefulOrderRemoval',
        },
      ),
    ];
  }
}
