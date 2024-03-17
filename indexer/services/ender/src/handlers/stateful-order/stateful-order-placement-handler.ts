import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  OrderPlaceV1_OrderPlacementStatus,
  OffChainUpdateV1,
  IndexerOrder,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../../lib/types';
import { Handler } from '../handler';

// TODO(IND-334): Rename to LongTermOrderPlacementHandler after deprecating StatefulOrderPlacement
export class StatefulOrderPlacementHandler extends Handler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(_: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    let order: IndexerOrder;
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    if (this.event.orderPlace !== undefined) {
      order = this.event.orderPlace!.order!;
    } else {
      order = this.event.longTermOrderPlacement!.order!;
    }
    return this.createKafkaEvents(order);
  }

  private createKafkaEvents(order: IndexerOrder): ConsolidatedKafkaEvent[] {
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
    ));

    return kafkaEvents;
  }
}
