import { logger } from '@dydxprotocol-indexer/base';
import { BuilderCodeParameters, IndexerOrder, IndexerOrder_Side } from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';

import * as OrderTable from '../stores/order-table';
import * as SubaccountTable from '../stores/subaccount-table';
import {
  OrderFromDatabase, OrderSide, PerpetualMarketFromDatabase, SubaccountFromDatabase,
} from '../types';
import { blockTimeFromIsoString } from './helpers';
import {
  humanToQuantums,
  priceToSubticks,
  tifToProtocolOrderTIF,
  orderTypeToProtocolConditionType,
} from './protocol-translations';

/**
 * Converts an order from the database to an IndexerOrder proto.
 * This is used to resend open stateful orders to Vulcan during Indexer fast sync
 * to uncross the orderbook.
 *
 * @param order
 */
export function convertToIndexerOrderWithSubaccount(
  order: OrderFromDatabase,
  perpetualMarket: PerpetualMarketFromDatabase,
  subaccount: SubaccountFromDatabase,
): IndexerOrder {
  if (!OrderTable.isLongTermOrConditionalOrder(order.orderFlags)) {
    logger.error({
      at: 'protocol-translations#convertToIndexerOrder',
      message: 'Order is not a long-term or conditional order',
      order,
    });
    throw new Error(`Order with flags ${order.orderFlags} is not a long-term or conditional order`);
  }
  if (!subaccount === undefined) {
    logger.error({
      at: 'protocol-translations#convertToIndexerOrder',
      message: 'Subaccount for order not found',
      order,
    });
    throw new Error(`Subaccount for order not found: ${order.subaccountId}`);
  }
  const triggerSubticks: Long = (order.triggerPrice === undefined || order.triggerPrice === null)
    ? Long.fromValue(0, true)
    : Long.fromString(priceToSubticks(order.triggerPrice, perpetualMarket), true);
  let builderCodeParameters: BuilderCodeParameters | undefined;

  if (order.builderAddress && order.feePpm) {
    builderCodeParameters = {
      builderAddress: order.builderAddress,
      feePpm: Number(order.feePpm),
    };
  }
  const indexerOrder: IndexerOrder = {
    orderId: {
      subaccountId: {
        owner: subaccount?.address!,
        number: subaccount?.subaccountNumber!,
      },
      clientId: Number(order.clientId),
      clobPairId: Number(order.clobPairId),
      orderFlags: Number(order.orderFlags),
    },
    side: order.side === OrderSide.BUY ? IndexerOrder_Side.SIDE_BUY : IndexerOrder_Side.SIDE_SELL,
    quantums: Long.fromString(humanToQuantums(
      order.size,
      perpetualMarket.atomicResolution,
    ).toFixed(), true),
    subticks: Long.fromString(priceToSubticks(
      order.price,
      perpetualMarket,
    ), true),
    goodTilBlockTime: blockTimeFromIsoString(order.goodTilBlockTime!),
    timeInForce: tifToProtocolOrderTIF(order.timeInForce),
    reduceOnly: order.reduceOnly,
    clientMetadata: Number(order.clientMetadata),
    conditionType: orderTypeToProtocolConditionType(order.type),
    conditionalOrderTriggerSubticks: triggerSubticks,
    builderCodeParams: builderCodeParameters,
    orderRouterAddress: order.orderRouterAddress ?? '',
  };

  return indexerOrder;
}

/**
 * Converts an order from the database to an IndexerOrder proto.
 * This is used to resend open stateful orders to Vulcan during Indexer fast sync
 * to uncross the orderbook.
 *
 * @param order
 */
export async function convertToIndexerOrder(
  order: OrderFromDatabase,
  perpetualMarket: PerpetualMarketFromDatabase,
): Promise<IndexerOrder> {
  const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
    order.subaccountId,
  );

  if (!subaccount === undefined) {
    logger.error({
      at: 'protocol-translations#convertToIndexerOrder',
      message: 'Subaccount for order not found',
      order,
    });
    throw new Error(`Subaccount for order not found: ${order.subaccountId}`);
  }
  return convertToIndexerOrderWithSubaccount(order, perpetualMarket, subaccount!);
}
