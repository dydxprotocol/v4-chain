import { vaultRefresher } from '@dydxprotocol-indexer/postgres';
import { IndexerOrderId } from '@dydxprotocol-indexer/v4-protos';

export function isVaultOrder(orderId: IndexerOrderId): boolean {
  return vaultRefresher.isVault(orderId.subaccountId!.owner!);
}
