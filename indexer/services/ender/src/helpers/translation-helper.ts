import { SubaccountUpdateEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { OrderFillWithLiquidity, SubaccountUpdate } from '../lib/translated-types';
import { OrderFillEventWithLiquidation, OrderFillEventWithLiquidity, OrderFillEventWithOrder } from '../lib/types';

export function subaccountUpdateEventV1ToSubaccountUpdate(
  event: SubaccountUpdateEventV1,
): SubaccountUpdate {
  return {
    subaccountId: event.subaccountId,
    updatedPerpetualPositions: event.updatedPerpetualPositions,
    updatedAssetPositions: event.updatedAssetPositions,
  };
}

export function orderFillEventV1ToOrderFill(
  event: OrderFillEventWithLiquidity,
): OrderFillWithLiquidity {
  return {
    makerOrder: event.event.makerOrder,
    order: event.event.order,
    liquidationOrder: event.event.liquidationOrder,
    fillAmount: event.event.fillAmount,
    makerFee: event.event.makerFee,
    takerFee: event.event.takerFee,
    totalFilledMaker: event.event.totalFilledMaker,
    totalFilledTaker: event.event.totalFilledTaker,
    liquidity: event.liquidity,
    affiliateRevShare: event.event.affiliateRevShare,
    makerBuilderFee: event.event.makerBuilderFee,
    takerBuilderFee: event.event.takerBuilderFee,
    makerBuilderAddress: event.event.makerBuilderAddress,
    takerBuilderAddress: event.event.takerBuilderAddress,
  };
}

export function orderFillWithLiquidityToOrderFillEventWithOrder(
  orderFillWithLiquidity: OrderFillWithLiquidity,
): OrderFillEventWithOrder {
  return {
    makerOrder: orderFillWithLiquidity.makerOrder!,
    order: orderFillWithLiquidity.order!,
    fillAmount: orderFillWithLiquidity.fillAmount,
    totalFilledMaker: orderFillWithLiquidity.totalFilledMaker,
    totalFilledTaker: orderFillWithLiquidity.totalFilledTaker,
    makerFee: orderFillWithLiquidity.makerFee,
    takerFee: orderFillWithLiquidity.takerFee,
    affiliateRevShare: orderFillWithLiquidity.affiliateRevShare,
    makerBuilderFee: orderFillWithLiquidity.makerBuilderFee,
    takerBuilderFee: orderFillWithLiquidity.takerBuilderFee,
    makerBuilderAddress: orderFillWithLiquidity.makerBuilderAddress,
    takerBuilderAddress: orderFillWithLiquidity.takerBuilderAddress,
  };
}

export function orderFillWithLiquidityToOrderFillEventWithLiquidation(
  orderFillWithLiquidity: OrderFillWithLiquidity,
): OrderFillEventWithLiquidation {
  return {
    makerOrder: orderFillWithLiquidity.makerOrder!,
    liquidationOrder: orderFillWithLiquidity.liquidationOrder!,
    fillAmount: orderFillWithLiquidity.fillAmount,
    totalFilledMaker: orderFillWithLiquidity.totalFilledMaker,
    totalFilledTaker: orderFillWithLiquidity.totalFilledTaker,
    makerFee: orderFillWithLiquidity.makerFee,
    takerFee: orderFillWithLiquidity.takerFee,
    affiliateRevShare: orderFillWithLiquidity.affiliateRevShare,
    makerBuilderFee: orderFillWithLiquidity.makerBuilderFee,
    takerBuilderFee: orderFillWithLiquidity.takerBuilderFee,
    makerBuilderAddress: orderFillWithLiquidity.makerBuilderAddress,
    takerBuilderAddress: orderFillWithLiquidity.takerBuilderAddress,
  };
}
