import {
  APIOrderStatus,
  APIOrderStatusEnum,
  apiTranslations,
  IsoString,
  OrderFromDatabase,
  OrderTable,
  PerpetualMarketFromDatabase,
  protocolTranslations,
  SubaccountMessageContents,
  SubaccountTable,
  TimeInForce,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerOrder,
  IndexerOrder_ConditionType,
  OrderPlaceV1_OrderPlacementStatus,
  RedisOrder,
  SubaccountMessage,
} from '@dydxprotocol-indexer/v4-protos';

import { SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION } from './constants';

/**
 * Gets the trigger price for an order, returns undefined if the order has an unspecified condition
 * type
 * @param order
 * @param perpetualMarket
 * @returns
 */
export function getTriggerPrice(
  order: IndexerOrder,
  perpetualMarket: PerpetualMarketFromDatabase,
): string | undefined {
  if (order.conditionType !== IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED) {
    return protocolTranslations.subticksToPrice(
      order.conditionalOrderTriggerSubticks.toString(),
      perpetualMarket,
    );
  }
  return undefined;
}

export function generateSubaccountMessageContents(
  redisOrder: RedisOrder,
  order: OrderFromDatabase | undefined,
  perpetualMarket: PerpetualMarketFromDatabase,
  placementStatus: OrderPlaceV1_OrderPlacementStatus,
  blockHeight: string | undefined,
): SubaccountMessageContents {
  const orderTIF: TimeInForce = protocolTranslations.protocolOrderTIFToTIF(
    redisOrder.order!.timeInForce,
  );
  const status: APIOrderStatus = (
    placementStatus === OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED
      ? APIOrderStatusEnum.OPEN
      : APIOrderStatusEnum.BEST_EFFORT_OPENED
  );
  const createdAtHeight: string | undefined = order?.createdAtHeight;
  const updatedAt: IsoString | undefined = order?.updatedAt;
  const updatedAtHeight: string | undefined = order?.updatedAtHeight;
  const contents: SubaccountMessageContents = {
    orders: [
      {
        id: OrderTable.orderIdToUuid(redisOrder.order!.orderId!),
        subaccountId: SubaccountTable.subaccountIdToUuid(
          redisOrder.order!.orderId!.subaccountId!,
        ),
        clientId: redisOrder.order!.orderId!.clientId.toString(),
        clobPairId: perpetualMarket.clobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(redisOrder.order!.side),
        size: redisOrder.size,
        price: redisOrder.price,
        status,
        type: protocolTranslations.protocolConditionTypeToOrderType(
          redisOrder.order!.conditionType,
        ),
        timeInForce: apiTranslations.orderTIFToAPITIF(orderTIF),
        postOnly: apiTranslations.isOrderTIFPostOnly(orderTIF),
        reduceOnly: redisOrder.order!.reduceOnly,
        orderFlags: redisOrder.order!.orderId!.orderFlags.toString(),
        goodTilBlock: protocolTranslations.getGoodTilBlock(redisOrder.order!)
          ?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(redisOrder.order!),
        ticker: redisOrder.ticker,
        ...(createdAtHeight && { createdAtHeight }),
        ...(updatedAt && { updatedAt }),
        ...(updatedAtHeight && { updatedAtHeight }),
        clientMetadata: redisOrder.order!.clientMetadata.toString(),
        triggerPrice: getTriggerPrice(redisOrder.order!, perpetualMarket),
      },
    ],
    ...(blockHeight && { blockHeight }),
  };
  return contents;
}

export function createSubaccountWebsocketMessage(
  redisOrder: RedisOrder,
  order: OrderFromDatabase | undefined,
  perpetualMarket: PerpetualMarketFromDatabase,
  placementStatus: OrderPlaceV1_OrderPlacementStatus,
  blockHeight: string | undefined,
): Buffer {
  const contents: SubaccountMessageContents = generateSubaccountMessageContents(
    redisOrder,
    order,
    perpetualMarket,
    placementStatus,
    blockHeight,
  );

  const subaccountMessage: SubaccountMessage = SubaccountMessage.fromPartial({
    contents: JSON.stringify(contents),
    subaccountId: redisOrder.order!.orderId!.subaccountId!,
    version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  });

  return Buffer.from(SubaccountMessage.encode(subaccountMessage).finish());
}
