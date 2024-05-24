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
  IndexerSubaccountId,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  RedisOrder,
  StatefulOrderEventV1,
  SubaccountId,
} from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../../config';
import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

// TODO(IND-334): Rename to LongTermOrderPlacementHandler after deprecating StatefulOrderPlacement
export class StatefulOrderPlacementHandler
  extends AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getOrderId(): string {
    let orderId: string;
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    if (this.event.orderPlace !== undefined) {
      orderId = OrderTable.orderIdToUuid(this.event.orderPlace!.order!.orderId!);
    } else {
      orderId = OrderTable.orderIdToUuid(this.event.longTermOrderPlacement!.order!.orderId!);
    }
    return orderId;
  }

  public getSubaccountId(): IndexerSubaccountId {
    let subaccountId: IndexerSubaccountId;
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    if (this.event.orderPlace !== undefined) {
      subaccountId = this.event.orderPlace!.order!.orderId!.subaccountId!;
    } else {
      subaccountId = this.event.longTermOrderPlacement!.order!.orderId!.subaccountId!;
    }
    return subaccountId;
  }

  public getParallelizationIds(): string[] {
    // Stateful Order Events with the same orderId
    return this.getParallelizationIdsFromOrderId(this.getOrderId());
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    let order: IndexerOrder;
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    if (this.event.orderPlace !== undefined) {
      order = this.event.orderPlace!.order!;
    } else {
      order = this.event.longTermOrderPlacement!.order!;
    }
    return this.createKafkaEvents(order, resultRow);
  }

  private createKafkaEvents(
    order: IndexerOrder,
    resultRow: pg.QueryResultRow,
  ): ConsolidatedKafkaEvent[] {
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
      {
        message_received_timestamp: this.messageReceivedTimestamp,
        event_type: 'StatefulOrderPlacement',
      },
    ));

    if (config.SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_STATEFUL_ORDERS) {
      const perpetualMarket: PerpetualMarketFromDatabase = perpetualMarketRefresher
        .getPerpetualMarketFromClobPairId(order.orderId!.clobPairId.toString())!;
      const dbOrder: OrderFromDatabase = OrderModel.fromJson(resultRow.order) as OrderFromDatabase;
      const redisOrder: RedisOrder = convertToRedisOrder(order, perpetualMarket);
      const subaccountContent: SubaccountMessageContents = generateSubaccountMessageContents(
        redisOrder,
        dbOrder,
        perpetualMarket,
        OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
        this.block.height.toString(),
      );

      const subaccountIdProto: SubaccountId = {
        owner: this.getSubaccountId().owner,
        number: this.getSubaccountId().number,
      };
      kafkaEvents.push(this.generateConsolidatedSubaccountKafkaEvent(
        JSON.stringify(subaccountContent),
        subaccountIdProto,
        this.getOrderId(),
        false,
        subaccountContent,
      ));
    }
    return kafkaEvents;
  }
}
