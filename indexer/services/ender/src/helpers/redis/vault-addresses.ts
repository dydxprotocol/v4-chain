import { stats } from "@dydxprotocol-indexer/base";
import { Options, VaultFromDatabase, VaultTable } from "@dydxprotocol-indexer/postgres";
import { VaultAddressesCache } from "@dydxprotocol-indexer/redis";

import config from '../../config';
import { redisClient } from "./redis-controller";

// /**
//  * Initialize vault addresses cache by adding all vaults in database.
//  */
// export async function initializeVaultAddressesCache(options?: Options): Promise<void> {
//   const start: number = Date.now();
//   const vaults: VaultFromDatabase[] = await VaultTable.findAll(
//       {}, [], options || { readReplica: true },
//   );
//   await Promise.all(
//       vaults.map(vault => VaultAddressesCache.addVaultAddress(vault.address, redisClient))
//   );
//   stats.timing(`${config.SERVICE_NAME}.initialize_vault_addresses_cache`, Date.now() - start);
// }

/**
 * Check if the given address is a vault address.
 */
export async function isVaultAddress(address: string): Promise<boolean> {
  return VaultAddressesCache.isVaultAddress(address, redisClient);
}
