import { RedisClient } from 'redis';

import { Options, VaultFromDatabase, VaultTable } from "@dydxprotocol-indexer/postgres";

// Cache key for vault addresses.
export const VAULT_ADDRESSES_CACHE_KEY: string = 'v4/vault_addresses';

/**
 * isVaultAddress returns true if a given address is a vault address, false otherwise.
 *
 * @param address
 * @param client
 * @returns A promise that resolves to true if the given address is a vault address.
 */
export async function isVaultAddress(
  address: string,
  client: RedisClient,
): Promise<boolean> {
  return new Promise((resolve, reject) => {
    client.sismember(VAULT_ADDRESSES_CACHE_KEY, address, (err, result) => {
      if (err) {
        return reject(err);
      }

      // `result` is 1 if `address` exists in the set, 0 otherwise.
      resolve(result === 1);
    });
  });
}

/**
 * Adds an address to vault address cache.
 * 
 * Ender calls this function when a new vault is created.
 *
 * @param address
 * @param client
 * @returns A promise that resolves to 1 if a new address is added, 0 if address already exists.
 */
export async function addVaultAddress(
  address: string,
  client: RedisClient,
): Promise<number> {
  return new Promise((resolve, reject) => {
    client.sadd(VAULT_ADDRESSES_CACHE_KEY, address, (err, result) => {
      if (err) {
        return reject(err);
      }
      resolve(result);
    });
  });
}

/**
 * Initialize vault addresses cache by adding all vaults in database.
 * 
 * Ender and Vulcan both call this function upon start-up.
 *
 * @param client
 * @param options
 */
export async function initialize(
  client: RedisClient,
  options?: Options,
): Promise<void> {
  const vaults: VaultFromDatabase[] = await VaultTable.findAll(
      {}, [], options || { readReplica: true },
  );
  await Promise.all(
      vaults.map(vault => addVaultAddress(vault.address, client))
  );
}
