import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultAddress, defaultVault } from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';
import { addVault, getVaultAddresses, isVault, updateVaults } from '../../src/loops/vault-refresher';

describe('vaultRefresher', () => {
  beforeAll(async () => {
    await migrate();
    await seedData();
    await updateVaults();
  });

  afterAll(async () => {
    await clearData();
    await teardown();
  });

  describe('isVault', () => {
    it('checks for vault address', async() => {
      expect(isVault(defaultVault.address)).toBe(true);
    });

    it('checks for non-vault address', async() => {
      expect(isVault(defaultAddress)).toBe(false);
    });
  });

  describe('getVaultAddresses', () => {
    it('gets all vault addresses', async() => {
      expect(getVaultAddresses()).toEqual(new Set([defaultVault.address]));
    });
  });

  describe('addVault', () => {
    it('adds new vault addresses', async() => {
      const newVaultAddr1 = 'dydx1234567';
      const newVaultAddr2 = 'dydx1765432';
      addVault(newVaultAddr1);
      addVault(newVaultAddr2);
      expect(isVault(newVaultAddr1)).toBe(true);
      expect(isVault(newVaultAddr2)).toBe(true);
      expect(getVaultAddresses()).toEqual(new Set([
        defaultVault.address, newVaultAddr1, newVaultAddr2,
      ]));
    });
  });
});
