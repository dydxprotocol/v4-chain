import {
  OrderTable,
  OrderStatus,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  OffChainUpdateV1,
  IndexerOrderId,
  OrderRemoveV1_OrderRemovalStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';

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

  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const orderIdProto: IndexerOrderId = this.event.orderRemoval!.removedOrderId!;
    await this.runFuncWithTimingStatAndErrorLogging(
      this.updateOrderStatus(orderIdProto, OrderStatus.CANCELED),
      this.generateTimingStatsOptions('cancel_order'),
    );

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
