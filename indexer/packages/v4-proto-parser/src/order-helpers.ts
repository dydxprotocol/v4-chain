import { createHash } from 'crypto';

import { IndexerOrderId } from '@dydxprotocol-indexer/v4-protos';

/**
 * Gets the hash of an order id. Matches the function in V4 to get the hash of an order.
 * https://github.com/dydxprotocol/v4/blob/311411a3ce92230d4866a7c4abb1422fbc4ef3b9/indexer/off_chain_updates/off_chain_updates.go#L293
 * @param orderId
 * @returns
 */
export function getOrderIdHash(orderId: IndexerOrderId): Buffer {
  const bytes: Buffer = Buffer.from(Uint8Array.from(IndexerOrderId.encode(orderId).finish()));
  return createHash('sha256').update(bytes).digest();
}
