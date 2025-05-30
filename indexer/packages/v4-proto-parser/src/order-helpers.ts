import { createHash } from 'crypto';

import { IndexerOrderId, IndexerOrder_TimeInForce } from '@dydxprotocol-indexer/v4-protos';

import { ORDER_FLAG_CONDITIONAL, ORDER_FLAG_LONG_TERM } from './constants';

/**
 * Gets the hash of an order id. Matches the function in V4 to get the hash of an order.
 * https://github.com/dydxprotocol/v4/blob/311411a3ce92230d4866a7c4abb1422fbc4ef3b9/indexer/off_chain_updates/off_chain_updates.go#L293
 * @param orderId
 * @returns
 */
export function getOrderIdHash(orderId: IndexerOrderId): Buffer {
  const bytes: Buffer = Buffer.from(IndexerOrderId.encode(orderId).finish());
  return createHash('sha256').update(bytes).digest();
}

export function isStatefulOrder(orderFlag: number | String): boolean {
  const numberOrderFlag: number = Number(orderFlag);
  // A string that is not a number will be converted to NaN, and should return false.
  if (Number.isNaN(numberOrderFlag)) {
    return false;
  }
  return numberOrderFlag === ORDER_FLAG_CONDITIONAL || numberOrderFlag === ORDER_FLAG_LONG_TERM;
}

export function isLongTermOrder(orderFlag: number | String): boolean {
  const numberOrderFlag: number = Number(orderFlag);
  // A string that is not a number will be converted to NaN, and should return false.
  if (Number.isNaN(numberOrderFlag)) {
    return false;
  }
  return numberOrderFlag === ORDER_FLAG_LONG_TERM;
}

export function requiresImmediateExecution(tif: IndexerOrder_TimeInForce): boolean {
  return (
    tif === IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL ||
    tif === IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC
  );
}
