import {
  OrderTable, protocolTranslations, SubaccountTable, testConstants,
} from '@dydxprotocol-indexer/postgres';
import { ORDER_FLAG_CONDITIONAL, ORDER_FLAG_LONG_TERM, ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerOrderId,
  OrderPlaceV1,
  OrderPlaceV1_OrderPlacementStatus,
  OrderRemoveV1,
  OrderRemoveV1_OrderRemovalStatus,
  OrderUpdateV1,
  IndexerOrder_Side,
  IndexerOrder_TimeInForce,
  RedisOrder,
  RedisOrder_TickerType,
  IndexerSubaccountId,
  IndexerOrder_ConditionType,
  OrderRemovalReason,
} from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';

export type OffChainUpdateOrderPlaceUpdateMessage = {
  orderUpdate: undefined,
  orderRemove: undefined,
  orderPlace: OrderPlaceV1,
};
export type OffChainUpdateOrderRemoveUpdateMessage = {
  orderUpdate: undefined,
  orderRemove: OrderRemoveV1,
  orderPlace: undefined,
};
export type OffChainUpdateOrderUpdateUpdateMessage = {
  orderUpdate: OrderUpdateV1,
  orderRemove: undefined,
  orderPlace: undefined,
};

export const defaultSubaccountId: IndexerSubaccountId = {
  owner: testConstants.defaultSubaccount.address,
  number: testConstants.defaultSubaccount.subaccountNumber,
};
export const defaultOrderId: IndexerOrderId = {
  subaccountId: defaultSubaccountId,
  clientId: 1,
  clobPairId: parseInt(testConstants.defaultPerpetualMarket.clobPairId, 10),
  orderFlags: ORDER_FLAG_SHORT_TERM,
};
export const defaultOrderIdGoodTilBlockTime: IndexerOrderId = {
  subaccountId: defaultSubaccountId,
  clientId: 2,
  clobPairId: parseInt(testConstants.defaultPerpetualMarket.clobPairId, 10),
  orderFlags: ORDER_FLAG_LONG_TERM,
};
export const defaultOrderIdConditional: IndexerOrderId = {
  subaccountId: defaultSubaccountId,
  clientId: 3,
  clobPairId: parseInt(testConstants.defaultPerpetualMarket.clobPairId, 10),
  orderFlags: ORDER_FLAG_CONDITIONAL,
};
export const defaultOrderIdVault: IndexerOrderId = {
  subaccountId: {
    owner: testConstants.defaultVaultAddress,
    number: 0,
  },
  clientId: 4,
  clobPairId: parseInt(testConstants.defaultPerpetualMarket.clobPairId, 10),
  orderFlags: ORDER_FLAG_LONG_TERM,
};

export const defaultSubaccountUuid: string = SubaccountTable.uuid(
  defaultSubaccountId.owner,
  defaultSubaccountId.number,
);

export const defaultLastUpdated: string = '1682000000000';

export const defaultOrder: IndexerOrder = {
  orderId: defaultOrderId,
  side: IndexerOrder_Side.SIDE_BUY,
  quantums: Long.fromValue(1_000_000, true),
  subticks: Long.fromValue(2_000_000, true),
  goodTilBlock: 1150,
  goodTilBlockTime: undefined,
  timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
  reduceOnly: false,
  clientMetadata: 0,
  conditionType: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
  conditionalOrderTriggerSubticks: Long.fromValue(0, true),
  orderRouterAddress: '',
};
export const defaultOrderGoodTilBlockTime: IndexerOrder = {
  ...defaultOrder,
  orderId: defaultOrderIdGoodTilBlockTime,
  goodTilBlockTime: 1_200_000_000,
  goodTilBlock: undefined,
};
export const defaultConditionalOrder: IndexerOrder = {
  ...defaultOrderGoodTilBlockTime,
  orderId: defaultOrderIdConditional,
  conditionType: IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS,
  conditionalOrderTriggerSubticks: Long.fromValue(190_000_000, true),
};
export const defaultOrderFok: IndexerOrder = {
  ...defaultOrder,
  timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
};
export const defaultOrderIoc: IndexerOrder = {
  ...defaultOrder,
  timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
};
export const defaultOrderVault: IndexerOrder = {
  ...defaultOrderGoodTilBlockTime,
  orderId: defaultOrderIdVault,
  timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
};

export const defaultOrderUuid: string = OrderTable.orderIdToUuid(defaultOrderId);
export const defaultOrderUuidGoodTilBlockTime: string = OrderTable.orderIdToUuid(
  defaultOrderIdGoodTilBlockTime,
);
export const defaultOrderUuidConditional: string = OrderTable.orderIdToUuid(
  defaultOrderIdConditional,
);
export const defaultOrderUuidVault: string = OrderTable.orderIdToUuid(
  defaultOrderIdVault,
);

export const defaultPrice = protocolTranslations.subticksToPrice(
  defaultOrder.subticks.toString(),
  testConstants.defaultPerpetualMarket,
);
export const defaultSize = protocolTranslations.quantumsToHumanFixedString(
  defaultOrder.quantums.toString(),
  testConstants.defaultPerpetualMarket.atomicResolution,
);
export const defaultRedisOrder: RedisOrder = {
  id: defaultOrderUuid,
  order: defaultOrder,
  ticker: testConstants.defaultPerpetualMarket.ticker,
  tickerType: RedisOrder_TickerType.TICKER_TYPE_PERPETUAL,
  price: defaultPrice,
  size: defaultSize,
};
export const defaultRedisOrderGoodTilBlockTime: RedisOrder = {
  ...defaultRedisOrder,
  id: defaultOrderUuidGoodTilBlockTime,
  order: defaultOrderGoodTilBlockTime,
};
export const defaultRedisOrderConditional: RedisOrder = {
  ...defaultRedisOrder,
  id: defaultOrderUuidConditional,
  order: defaultConditionalOrder,
};
export const defaultRedisOrderFok: RedisOrder = {
  ...defaultRedisOrder,
  order: defaultOrderFok,
};
export const defaultRedisOrderIoc: RedisOrder = {
  ...defaultRedisOrder,
  order: defaultOrderIoc,
};
export const defaultRedisOrderVault: RedisOrder = {
  ...defaultRedisOrderGoodTilBlockTime,
  id: defaultOrderUuidVault,
  order: defaultOrderVault,
};

export const orderPlace: OffChainUpdateOrderPlaceUpdateMessage = {
  orderUpdate: undefined,
  orderRemove: undefined,
  orderPlace: {
    order: defaultOrder,
    placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
  },
};
export const orderRemove: OffChainUpdateOrderRemoveUpdateMessage = {
  orderPlace: undefined,
  orderUpdate: undefined,
  orderRemove: {
    removedOrderId: defaultOrderId,
    reason: OrderRemovalReason.ORDER_REMOVAL_REASON_INTERNAL_ERROR,
    removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
  },
};
export const orderUpdate: OffChainUpdateOrderUpdateUpdateMessage = {
  orderPlace: undefined,
  orderRemove: undefined,
  orderUpdate: {
    orderId: defaultOrderId,
    totalFilledQuantums: Long.fromValue(250_500, true),
  },
};

export const isolatedSubaccountId: IndexerSubaccountId = {
  owner: testConstants.isolatedSubaccount.address,
  number: testConstants.isolatedSubaccount.subaccountNumber,
};
export const isolatedMarketOrderId: IndexerOrderId = {
  subaccountId: isolatedSubaccountId,
  clientId: 1,
  clobPairId: parseInt(testConstants.isolatedPerpetualMarket.clobPairId, 10),
  orderFlags: ORDER_FLAG_SHORT_TERM,
};
export const isolatedMarketOrder: IndexerOrder = {
  orderId: isolatedMarketOrderId,
  side: IndexerOrder_Side.SIDE_BUY,
  quantums: Long.fromValue(1_000_000, true),
  subticks: Long.fromValue(2_000_000, true),
  goodTilBlock: 1150,
  goodTilBlockTime: undefined,
  timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
  reduceOnly: false,
  clientMetadata: 0,
  conditionType: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
  conditionalOrderTriggerSubticks: Long.fromValue(0, true),
  orderRouterAddress: '',
};

export const isolatedMarketOrderUuid: string = OrderTable.orderIdToUuid(isolatedMarketOrderId);

export const isolatedMarketRedisOrder: RedisOrder = {
  id: isolatedMarketOrderUuid,
  order: isolatedMarketOrder,
  ticker: testConstants.isolatedPerpetualMarket.ticker,
  tickerType: RedisOrder_TickerType.TICKER_TYPE_PERPETUAL,
  price: protocolTranslations.subticksToPrice(
    isolatedMarketOrder.subticks.toString(),
    testConstants.isolatedPerpetualMarket,
  ),
  size: protocolTranslations.quantumsToHumanFixedString(
    isolatedMarketOrder.quantums.toString(),
    testConstants.isolatedPerpetualMarket.atomicResolution,
  ),
};
