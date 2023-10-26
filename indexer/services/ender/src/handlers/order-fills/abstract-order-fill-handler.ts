import { logger } from '@dydxprotocol-indexer/base';
import {
  AssetFromDatabase,
  assetRefresher,
  FillCreateObject,
  FillFromDatabase,
  FillTable,
  FillType,
  Liquidity,
  OrderCreateObject,
  OrderFromDatabase,
  OrderSide,
  OrderStatus,
  OrderTable,
  OrderType,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualPositionColumns,
  PerpetualPositionFromDatabase,
  PerpetualPositionTable,
  protocolTranslations,
  SubaccountMessageContents,
  SubaccountTable,
  TendermintEventTable,
  TradeMessageContents,
  UpdatedPerpetualPositionSubaccountKafkaObject,
  USDC_ASSET_ID,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerOrder_Side,
  IndexerSubaccountId,
  LiquidationOrderV1,
  IndexerOrderId, OffChainUpdateV1,
  OrderRemovalReason, OrderRemoveV1_OrderRemovalStatus,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import Long from 'long';
import { DateTime } from 'luxon';

import {
  generateFillSubaccountMessage,
  generateOrderSubaccountMessage,
  generatePerpetualPositionsContents,
  isLiquidation,
} from '../../helpers/kafka-helper';
import {
  getPrice,
  getSize,
  getWeightedAverage,
  indexerTendermintEventToTransactionIndex,
  perpetualPositionAndOrderSideMatching,
} from '../../lib/helper';
import {
  ConsolidatedKafkaEvent,
  OrderFillEventWithLiquidation,
  OrderFillEventWithOrder,
  PriceFields,
  SumFields,
} from '../../lib/types';
import { Handler } from '../handler';

export type OrderFillEventBase = {
  subaccountId: string;
  orderId: string | undefined;
  fillType: FillType;
  clobPairId: string;
  side: OrderSide;
  makerOrder: IndexerOrder,
  fillAmount: Long;
  liquidity: Liquidity;
  clientMetadata?: string;
  fee: Long;
};

export abstract class AbstractOrderFillHandler<T> extends Handler<T> {
  protected liquidityToOrder(
    castedOrderFillEventMessage: OrderFillEventWithOrder, liquidity: Liquidity,
  ): IndexerOrder {
    return liquidity === Liquidity.MAKER
      ? castedOrderFillEventMessage.makerOrder
      : castedOrderFillEventMessage.order;
  }

  protected createEventBase(
    castedOrderFillEventMessage: OrderFillEventWithOrder,
    liquidity: Liquidity,
  ): OrderFillEventBase {
    // event is validated before calling this method, so all fields on the order must exist
    const order: IndexerOrder = this.liquidityToOrder(castedOrderFillEventMessage, liquidity)!;
    return this.createEventBaseFromOrder(
      order,
      castedOrderFillEventMessage.makerOrder,
      castedOrderFillEventMessage.fillAmount,
      liquidity,
      FillType.LIMIT,
      liquidity === Liquidity.MAKER
        ? castedOrderFillEventMessage.makerFee
        : castedOrderFillEventMessage.takerFee,
    );
  }

  protected createEventBaseFromOrder(
    order: IndexerOrder,
    makerOrder: IndexerOrder,
    fillAmount: Long,
    liquidity: Liquidity,
    fillType: FillType,
    fee: Long,
  ): OrderFillEventBase {
    return {
      subaccountId: SubaccountTable.subaccountIdToUuid(order.orderId!.subaccountId!),
      orderId: OrderTable.orderIdToUuid(order.orderId!),
      fillType,
      clobPairId: order.orderId!.clobPairId.toString(),
      side: protocolTranslations.protocolOrderSideToOrderSide(order.side),
      makerOrder,
      fillAmount,
      liquidity,
      clientMetadata: order.clientMetadata.toString(),
      fee,
    };
  }

  protected createEventBaseFromLiquidation(
    castedLiquidationFillEventMessage: OrderFillEventWithLiquidation,
    liquidity: Liquidity,
  ): OrderFillEventBase {
    // event is validated before calling this method, so all fields on the order must exist
    if (liquidity === Liquidity.TAKER) {
      const order: LiquidationOrderV1 = castedLiquidationFillEventMessage.liquidationOrder;
      return {
        subaccountId: SubaccountTable.subaccountIdToUuid(order.liquidated!),
        orderId: undefined,
        fillType: FillType.LIQUIDATED,
        clobPairId: order.clobPairId.toString(),
        side: order.isBuy ? OrderSide.BUY : OrderSide.SELL,
        makerOrder: castedLiquidationFillEventMessage.makerOrder,
        fillAmount: castedLiquidationFillEventMessage.fillAmount,
        liquidity,
        fee: castedLiquidationFillEventMessage.takerFee,
      };
    } else {
      return this.createEventBaseFromOrder(
        castedLiquidationFillEventMessage.makerOrder,
        castedLiquidationFillEventMessage.makerOrder,
        castedLiquidationFillEventMessage.fillAmount,
        liquidity,
        FillType.LIQUIDATION,
        castedLiquidationFillEventMessage.makerFee,
      );
    }
  }

  protected createFillFromEvent(
    perpetualMarket: PerpetualMarketFromDatabase,
    event: OrderFillEventBase,
  ): Promise<FillFromDatabase> {
    // event is validated before calling this method, so all fields on the order must exist
    const eventId: Buffer = TendermintEventTable.createEventId(
      this.block.height.toString(),
      indexerTendermintEventToTransactionIndex(this.indexerTendermintEvent),
      this.indexerTendermintEvent.eventIndex,
    );
    const size: string = protocolTranslations.quantumsToHumanFixedString(
      event.fillAmount.toString(),
      perpetualMarket.atomicResolution,
    );
    const price: string = getPrice(
      event.makerOrder,
      perpetualMarket,
    );
    const transactionIndex: number = indexerTendermintEventToTransactionIndex(
      this.indexerTendermintEvent,
    );
    const asset: AssetFromDatabase = assetRefresher.getAssetFromId(USDC_ASSET_ID);
    const fee: string = protocolTranslations.quantumsToHumanFixedString(
      event.fee.toString(),
      asset.atomicResolution,
    );

    const fillToCreate: FillCreateObject = {
      subaccountId: event.subaccountId,
      side: event.side,
      liquidity: event.liquidity,
      type: event.fillType,
      clobPairId: event.clobPairId,
      orderId: event.orderId,
      size,
      price,
      quoteAmount: Big(size).times(price).toFixed(),
      eventId,
      transactionHash: this.block.txHashes[transactionIndex],
      createdAt: this.timestamp.toISO(),
      createdAtHeight: this.block.height.toString(),
      clientMetadata: event.clientMetadata,
      fee,
    };

    return FillTable.create(fillToCreate, { txId: this.txId });
  }

  protected async getLatestPerpetualPosition(
    perpetualMarket: PerpetualMarketFromDatabase,
    event: OrderFillEventBase,
  ): Promise<PerpetualPositionFromDatabase> {
    const latestPerpetualPositions:
    PerpetualPositionFromDatabase[] = await PerpetualPositionTable.findAll(
      {
        subaccountId: [event.subaccountId],
        perpetualId: [perpetualMarket.id],
        limit: 1,
      },
      [],
      { txId: this.txId },
    );

    if (latestPerpetualPositions.length === 0) {
      logger.error({
        at: 'orderFillHandler#getLatestPerpetualPosition',
        message: 'Unable to find existing perpetual position.',
        blockHeight: this.block.height,
        clobPairId: event.clobPairId,
        subaccountId: event.subaccountId,
        orderId: event.orderId,
      });
      throw new Error(`Unable to find existing perpetual position. blockHeight: ${this.block.height}, clobPairId: ${event.clobPairId}, subaccountId: ${event.subaccountId}, orderId: ${event.orderId}`);
    }

    return latestPerpetualPositions[0];
  }

  protected async updatePerpetualPosition(
    perpetualMarket: PerpetualMarketFromDatabase,
    orderFillEventBase: OrderFillEventBase,
  ): Promise<PerpetualPositionFromDatabase> {
    const latestPerpetualPosition:
    PerpetualPositionFromDatabase = await this.getLatestPerpetualPosition(
      perpetualMarket,
      orderFillEventBase,
    );

    // update (sumOpen and entryPrice) or (sumClose and exitPrice)
    let sumField: SumFields;
    let priceField: PriceFields;
    if (perpetualPositionAndOrderSideMatching(
      latestPerpetualPosition.side, orderFillEventBase.side,
    )) {
      sumField = PerpetualPositionColumns.sumOpen;
      priceField = PerpetualPositionColumns.entryPrice;
    } else {
      sumField = PerpetualPositionColumns.sumClose;
      priceField = PerpetualPositionColumns.exitPrice;
    }

    const size: string = protocolTranslations.quantumsToHumanFixedString(
      orderFillEventBase.fillAmount.toString(),
      perpetualMarket.atomicResolution,
    );
    const price: string = getPrice(
      orderFillEventBase.makerOrder,
      perpetualMarket,
    );

    const updatedPerpetualPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.update(
      {
        id: latestPerpetualPosition.id,
        [sumField]: Big(latestPerpetualPosition[sumField]).plus(size).toFixed(),
        [priceField]: getWeightedAverage(
          latestPerpetualPosition[priceField] ?? '0',
          latestPerpetualPosition[sumField],
          price,
          size,
        ),
      },
      { txId: this.txId },
    );
    if (updatedPerpetualPosition === undefined) {
      logger.error({
        at: 'orderFillHandler#handle',
        message: 'Unable to update perpetual position',
        latestPerpetualPositionId: latestPerpetualPosition.id,
        orderFillEventBase,
      });
      throw new Error(`Unable to update perpetual position with id: ${latestPerpetualPosition.id}`);
    }
    return updatedPerpetualPosition;
  }

  protected upsertOrderFromEvent(
    perpetualMarket: PerpetualMarketFromDatabase,
    order: IndexerOrder,
    totalFilledFromProto: Long,
    isCanceled: boolean,
  ): Promise<OrderFromDatabase> {
    const size: string = getSize(order, perpetualMarket);
    const price: string = getPrice(order, perpetualMarket);
    const totalFilled: string = protocolTranslations.quantumsToHumanFixedString(
      totalFilledFromProto.toString(10),
      perpetualMarket.atomicResolution,
    );

    const orderToCreate: OrderCreateObject = {
      subaccountId: SubaccountTable.subaccountIdToUuid(order.orderId!.subaccountId!),
      clientId: order.orderId!.clientId.toString(),
      clobPairId: order.orderId!.clobPairId.toString(),
      side: order.side === IndexerOrder_Side.SIDE_BUY ? OrderSide.BUY : OrderSide.SELL,
      size,
      totalFilled,
      price,
      type: OrderType.LIMIT, // TODO: Add additional order types once we support
      status: this.getOrderStatus(isCanceled, size, totalFilled),
      timeInForce: protocolTranslations.protocolOrderTIFToTIF(order.timeInForce),
      reduceOnly: order.reduceOnly,
      orderFlags: order.orderId!.orderFlags.toString(),
      goodTilBlock: protocolTranslations.getGoodTilBlock(order)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(order),
      clientMetadata: order.clientMetadata.toString(),
      updatedAt: DateTime.fromJSDate(this.block.time!).toISO(),
      updatedAtHeight: this.block.height.toString(),
    };

    return OrderTable.upsert(orderToCreate, { txId: this.txId });
  }

  protected getOrderStatus(
    isCanceled: boolean,
    size: string,
    totalFilled: string,
  ): OrderStatus {
    if (isCanceled) {
      return OrderStatus.BEST_EFFORT_CANCELED;
    }
    if (Big(size).lte(totalFilled)) {
      return OrderStatus.FILLED;
    }
    return OrderStatus.OPEN;
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
          liquidation: isLiquidation(fill),
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
