import { logger } from '@dydxprotocol-indexer/base';
import {
  OrderCreateObject,
  OrderFromDatabase,
  OrderSide,
  OrderStatus,
  OrderTable,
  OrderType,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  protocolTranslations,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  OrderPlaceV1_OrderPlacementStatus,
  OffChainUpdateV1,
  IndexerOrder,
  IndexerOrder_Side,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';

import { getPrice, getSize } from '../../lib/helper';
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
      this.upsertOrder(perpetualMarket!, order),
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

  /**
   * Upsert order to database, because there may be an existing order with the orderId in the
   * database.
   */
  // eslint-disable-next-line @typescript-eslint/require-await
  protected async upsertOrder(
    perpetualMarket: PerpetualMarketFromDatabase,
    order: IndexerOrder,
  ): Promise<OrderFromDatabase> {
    const size: string = getSize(order, perpetualMarket);
    const price: string = getPrice(order, perpetualMarket);

    const orderToCreate: OrderCreateObject = {
      subaccountId: SubaccountTable.subaccountIdToUuid(order.orderId!.subaccountId!),
      clientId: order.orderId!.clientId.toString(),
      clobPairId: order.orderId!.clobPairId.toString(),
      side: order.side === IndexerOrder_Side.SIDE_BUY ? OrderSide.BUY : OrderSide.SELL,
      size,
      totalFilled: '0',
      price,
      type: OrderType.LIMIT, // TODO: Add additional order types once we support
      status: OrderStatus.OPEN,
      timeInForce: protocolTranslations.protocolOrderTIFToTIF(order.timeInForce),
      reduceOnly: order.reduceOnly,
      orderFlags: order.orderId!.orderFlags.toString(),
      // On chain orders must have a goodTilBlockTime rather than a goodTilBlock
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(order),
      createdAtHeight: this.block.height.toString(),
      clientMetadata: order.clientMetadata.toString(),
    };

    return OrderTable.upsert(orderToCreate, { txId: this.txId });
  }
}
