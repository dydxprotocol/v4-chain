import { logger } from '@dydxprotocol-indexer/base';
import {
  OrderFromDatabase,
  OrderStatus,
  OrderTable,
  PerpetualMarketFromDatabase,
  orderTranslations,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerOrderId,
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
    const orderIdProto: IndexerOrderId = this.event.conditionalOrderTriggered!.triggeredOrderId!;
    const orderFromDatabase: OrderFromDatabase = await this.runFuncWithTimingStatAndErrorLogging(
      this.updateOrderStatus(orderIdProto, OrderStatus.OPEN),
      this.generateTimingStatsOptions('trigger_order'),
    );

    const clobPairId: string = orderIdProto.clobPairId.toString();
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(clobPairId);
    if (perpetualMarket === undefined) {
      logger.error({
        at: 'statefulOrderPlacementHandler#internalHandle',
        message: 'Unable to find perpetual market',
        clobPairId,
        orderIdProto,
      });
      throw new Error(`Unable to find perpetual market with clobPairId: ${clobPairId}`);
    }

    // The conditional order was triggered, so send a message to vulcan to place the order
    const order: IndexerOrder = await orderTranslations.convertToIndexerOrder(
      orderFromDatabase,
      perpetualMarket,
    );
    const offChainUpdate: OffChainUpdateV1 = OffChainUpdateV1.fromPartial({
      orderPlace: {
        order,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
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
