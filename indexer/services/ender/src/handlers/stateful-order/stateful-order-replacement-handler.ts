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

import config from '../../config';
import { ConsolidatedKafkaEvent } from '../../lib/types';
import { AbstractStatefulOrderHandler } from '../abstract-stateful-order-handler';

export class StatefulOrderReplacementHandler
  extends AbstractStatefulOrderHandler<StatefulOrderEventV1> {
  eventType: string = 'StatefulOrderEvent';

  public getOrderId(): string {
    const orderId = OrderTable.orderIdToUuid(this.event.orderReplacement!.order!.orderId!);
    return orderId;
  }

  public getSubaccountId(): IndexerSubaccountId {
    const subaccountId = this.event.orderReplacement!.order!.orderId!.subaccountId!;
    return subaccountId;
  }

  public getParallelizationIds(): string[] {
    // Stateful Order Events with the same orderId
    return this.getParallelizationIdsFromOrderId(this.getOrderId());
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const oldOrderId = this.event.orderReplacement!.oldOrderId!;
    const order = this.event.orderReplacement!.order!;
    return this.createKafkaEvents(oldOrderId, order, resultRow);
  }

  private createKafkaEvents(
    oldOrderId: IndexerOrderId,
    order: IndexerOrder,
    resultRow: pg.QueryResultRow,
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
