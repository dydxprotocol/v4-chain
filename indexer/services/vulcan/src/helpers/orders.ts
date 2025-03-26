import { libHelpers } from '@dydxprotocol-indexer/postgres';
import { IndexerOrderId } from '@dydxprotocol-indexer/v4-protos';

export function isVaultOrder(orderId: IndexerOrderId): boolean {
  return libHelpers.VAULTS_CLOB_0_TO_999.has(orderId.subaccountId!.owner!);
}
