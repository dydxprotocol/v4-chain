import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  OffChainUpdateV1,
  IndexerOrderId,
  OrderRemoveV1_OrderRemovalStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../../lib/types';
import { Handler } from '../handler';

export class StatefulOrderRemovalHandler extends Handler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

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

    return [
      this.generateConsolidatedVulcanKafkaEvent(
        getOrderIdHash(orderIdProto),
        offChainUpdate,
      ),
    ];
  }
}
