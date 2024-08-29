import {
  FillFromDatabase,
  FillType,
  fillTypeToTradeType,
  Liquidity,
  OrderFromDatabase,
  OrderSide,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  SubaccountMessageContents,
  TradeMessageContents,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerSubaccountId,
  IndexerOrderId, OffChainUpdateV1,
  OrderRemovalReason, OrderRemoveV1_OrderRemovalStatus,
} from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';

import {
  generateFillSubaccountMessage,
  generateOrderSubaccountMessage,
  generatePerpetualPositionsContents,
} from '../../helpers/kafka-helper';
import {
  ConsolidatedKafkaEvent,
  OrderFillEventWithOrder,
} from '../../lib/types';
import { Handler } from '../handler';

export type OrderFillEventBase = {
  subaccountId: string,
  orderId: string | undefined,
  fillType: FillType,
  clobPairId: string,
  side: OrderSide,
  makerOrder: IndexerOrder,
  fillAmount: Long,
  liquidity: Liquidity,
  clientMetadata?: string,
  fee: Long,
};

export abstract class AbstractOrderFillHandler<T> extends Handler<T> {
  protected liquidityToOrder(
    castedOrderFillEventMessage: OrderFillEventWithOrder, liquidity: Liquidity,
  ): IndexerOrder {
    return liquidity === Liquidity.MAKER
      ? castedOrderFillEventMessage.makerOrder
      : castedOrderFillEventMessage.order;
  }

  /**
   * @param order - order may be undefined if the fill is a liquidation and this is the TAKER
   */
  protected generateConsolidatedKafkaEvent(
    subaccountIdProto: IndexerSubaccountId,
    order: OrderFromDatabase | undefined,
    position: UpdatedPerpetualPositionSubaccountKafkaObject | undefined,
    fill: FillFromDatabase,
    perpetualMarket: PerpetualMarketFromDatabase,
  ): ConsolidatedKafkaEvent {
    const message: SubaccountMessageContents = {
      fills: [
        generateFillSubaccountMessage(fill, perpetualMarket.ticker),
      ],
      perpetualPositions: position === undefined ? undefined : generatePerpetualPositionsContents(
        subaccountIdProto,
        [position],
        perpetualMarketRefresher.getPerpetualMarketsMap(),
      ),
      blockHeight: this.block.height.toString(),
    };
    if (order !== undefined) {
      message.orders = [
        generateOrderSubaccountMessage(order, perpetualMarket.ticker),
      ];
    }
    return this.generateConsolidatedSubaccountKafkaEvent(
      JSON.stringify(message),
      subaccountIdProto,
      order?.id,
      true,
      message,
    );
  }

  protected generateTradeKafkaEventFromTakerOrderFill(
    fill: FillFromDatabase,
  ): ConsolidatedKafkaEvent {
    const tradeContents: TradeMessageContents = {
      trades: [
        {
          id: fill.eventId.toString('hex'),
          size: fill.size,
          price: fill.price,
          side: fill.side.toString(),
          createdAt: fill.createdAt,
          type: fillTypeToTradeType(fill.type),
        },
      ],
    };
    return this.generateConsolidatedTradeKafkaEvent(
      JSON.stringify(tradeContents),
      fill.clobPairId,
    );
  }

  /**
   * Get a ConsolidatedKafkaEvent containing an order update to be sent to vulcan to update the
   * total filled amount of the order.
   * @param orderId
   * @param totalFilledQuantums
   * @returns
   */
  protected getOrderUpdateKafkaEvent(
    orderId: IndexerOrderId,
    totalFilledQuantums: Long,
  ): ConsolidatedKafkaEvent {
    const offChainUpdate: OffChainUpdateV1 = OffChainUpdateV1.fromPartial({
      orderUpdate: {
        orderId,
        totalFilledQuantums,
      },
    });
    return this.generateConsolidatedVulcanKafkaEvent(
      getOrderIdHash(orderId),
      offChainUpdate,
    );
  }

  /**
   * Get a ConsolidatedKafkaEvent containing an order remove to be sent to vulcan to remove a fully
   * filled order.
   * @param orderId
   * @returns
   */
  protected getOrderRemoveKafkaEvent(
    orderId: IndexerOrderId,
  ): ConsolidatedKafkaEvent {
    const offchainUpdate: OffChainUpdateV1 = OffChainUpdateV1.fromPartial({
      orderRemove: {
        removedOrderId: orderId,
        reason: OrderRemovalReason.ORDER_REMOVAL_REASON_FULLY_FILLED,
        removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_FILLED,
      },
    });
    return this.generateConsolidatedVulcanKafkaEvent(
      getOrderIdHash(orderId),
      offchainUpdate,
    );
  }
}
