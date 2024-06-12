import { generateSubaccountMessageContents } from '@dydxprotocol-indexer/kafka';
import {
  OrderFromDatabase,
  OrderModel,
  OrderTable,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  SubaccountMessageContents,
} from '@dydxprotocol-indexer/postgres';
import { convertToRedisOrder } from '@dydxprotocol-indexer/redis';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerOrderId,
  IndexerSubaccountId,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  RedisOrder,
  StatefulOrderEventV1,
  SubaccountId,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

export class StatefulOrderReplacementHandler
  extends AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  private getOrderId(): string {
    const orderId = OrderTable.orderIdToUuid(this.event.orderReplacement!.order!.orderId!);
    return orderId;
  }

  public getParallelizationIds(): string[] {
    // Stateful Order Events with the same orderId
    return this.getParallelizationIdsFromOrderId(this.getOrderId());
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(_: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const oldOrderId = this.event.orderReplacement!.oldOrderId!;
    const order = this.event.orderReplacement!.order!;
    return this.createKafkaEvents(oldOrderId, order);
  }

  private createKafkaEvents(
    oldOrderId: IndexerOrderId,
    order: IndexerOrder,
  ): ConsolidatedKafkaEvent[] {
    const kafkaEvents: ConsolidatedKafkaEvent[] = [];

    const offChainUpdate: OffChainUpdateV1 = OffChainUpdateV1.fromPartial({
      orderReplace: {
        oldOrderId,
        order,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
      },
    });
    kafkaEvents.push(this.generateConsolidatedVulcanKafkaEvent(
      getOrderIdHash(order.orderId!),
      offChainUpdate,
      {
        message_received_timestamp: this.messageReceivedTimestamp,
        event_type: 'StatefulOrderReplacement',
      },
    ));

    return kafkaEvents;
  }
}
