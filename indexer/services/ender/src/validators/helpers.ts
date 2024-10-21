import {
  IndexerOrder,
  IndexerOrderId,
  IndexerOrder_Side,
} from '@dydxprotocol-indexer/v4-protos';

export function validateOrderIdAndReturnErrorMessage(orderId: IndexerOrderId): string | undefined {
  if (orderId.subaccountId === undefined) {
    return 'OrderId must contain a subaccountId';
  }

  return undefined;
}

export function validateOrderAndReturnErrorMessage(order: IndexerOrder): string | undefined {
  if (order.orderId === undefined) {
    return 'Order must contain an orderId';
  }

  const errorMessage: string | undefined = validateOrderIdAndReturnErrorMessage(order.orderId);
  if (errorMessage !== undefined) {
    return errorMessage;
  }

  if (order.side === IndexerOrder_Side.SIDE_UNSPECIFIED) {
    return ' Order must specify an order side';
  }

  if (order.goodTilBlock === undefined && order.goodTilBlockTime === undefined) {
    return 'Order must contain a defined goodTilOneof';
  }

  if (order.routerFeePpm < 0 || order.routerFeePpm > 1000000) {
    return 'Router fee ppm must be between 0 and 1000000';
  }

  if (order.routerSubaccountId === undefined && order.routerFeePpm > 0) {
    return 'Router subaccount ID must be set if router fee ppm is greater than 0';
  }

  return undefined;
}
