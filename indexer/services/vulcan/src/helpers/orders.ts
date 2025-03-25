import { VaultAddressesCache } from '@dydxprotocol-indexer/redis';
import { IndexerOrderId } from '@dydxprotocol-indexer/v4-protos';
import { redisClient } from './redis/redis-controller';

/**
 * Check if an order belongs to a vault.
 *
 * @param orderId
 * @returns True if the order is a vault order, false otherwise.
 */
export async function isVaultOrder(orderId: IndexerOrderId): Promise<boolean> {
  return VaultAddressesCache.isVaultAddress(orderId.subaccountId!.owner!, redisClient);
}
