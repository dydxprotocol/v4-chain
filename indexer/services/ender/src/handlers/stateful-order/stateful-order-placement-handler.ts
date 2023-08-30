import { logger } from '@dydxprotocol-indexer/base';
import {
  OrderTable,
  OrderType,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  OrderStatus,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  OrderPlaceV1_OrderPlacementStatus,
  OffChainUpdateV1,
  IndexerOrder,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

// TODO(IND-334): Rename to LongTermOrderPlacementHandler after deprecating StatefulOrderPlacement
export class StatefulOrderPlacementHandler extends
  AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getParallelizationIds(): string[] {
    // Stateful Order Events with the same orderId
    let orderId: string;
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    if (this.event.orderPlace !== undefined) {
      orderId = OrderTable.orderIdToUuid(this.event.orderPlace!.order!.orderId!);
    } else {
      orderId = OrderTable.orderIdToUuid(this.event.longTermOrderPlacement!.order!.orderId!);
    }
    return this.getParallelizationIdsFromOrderId(orderId);
  }

  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    let order: IndexerOrder;
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    if (this.event.orderPlace !== undefined) {
      order = this.event.orderPlace!.order!;
    } else {
      order = this.event.longTermOrderPlacement!.order!;
    }
    const clobPairId: string = order.orderId!.clobPairId.toString();
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(clobPairId);
    if (perpetualMarket === undefined) {
      logger.error({
        at: 'statefulOrderPlacementHandler#internalHandle',
        message: 'Unable to find perpetual market',
        clobPairId,
        order,
      });
      throw new Error(`Unable to find perpetual market with clobPairId: ${clobPairId}`);
    }

    await this.runFuncWithTimingStatAndErrorLogging(
      this.upsertOrder(perpetualMarket!, order, OrderType.LIMIT, OrderStatus.OPEN),
      this.generateTimingStatsOptions('upsert_order'),
    );

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
