import {
  OrderFromDatabase,
  OrderTable,
  PerpetualMarketFromDatabase,
  orderTranslations,
  SubaccountFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';

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
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const result:
    [OrderFromDatabase,
      PerpetualMarketFromDatabase,
      SubaccountFromDatabase | undefined] = await this.handleEventViaSqlFunction();

    const order: IndexerOrder = orderTranslations.convertToIndexerOrderWithSubaccount(
      result[0], result[1], result[2]!);
    return this.createKafkaEvents(order);
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
      ),
    ];
  }
}
