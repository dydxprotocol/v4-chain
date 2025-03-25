import { testConstants } from '@dydxprotocol-indexer/postgres';
import { addVaultAddress, isVaultAddress } from '../../src/caches/vault-addresses-cache';
import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';

describe('vaultAddressesCache', () => {
  afterEach(async () => {
    await deleteAllAsync(client);
  });

  it('should add a new address and return 1', async () => {
    const address = testConstants.defaultAddress;
    const result = await addVaultAddress(address, client);
    expect(result).toBe(1);

    // Check that the address now exists in the cache.
    const exists = await isVaultAddress(address, client);
    expect(exists).toBe(true);
  });

  it('should add an existing address and return 0', async () => {
    const address = testConstants.defaultAddress;
    // Add the same address twice.
    await addVaultAddress(address, client);
    const result = await addVaultAddress(address, client);
    expect(result).toBe(0);
  });

  it('should add two different addresses', async () => {
    const address1 = testConstants.defaultAddress;
    const address2 = testConstants.defaultAddress3;
    expect(await addVaultAddress(address1, client)).toBe(1);
    expect(await addVaultAddress(address2, client)).toBe(1);
    
    expect(await isVaultAddress(address1, client)).toBe(true);
    expect(await isVaultAddress(address2, client)).toBe(true);
  });
  
  it('should return false if the address is not in cache', async () => {
    const exists = await isVaultAddress('0xABC', client);
    expect(exists).toBe(false);
  });
});
