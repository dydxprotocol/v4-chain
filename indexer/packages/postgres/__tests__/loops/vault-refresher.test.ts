import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultAddress, defaultVault } from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';
import { isVault, getVaultAddresses, updateVaults } from '../../src/loops/vault-refresher';

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

  describe('getVaultAddress', () => {
    it('gets all vault addresses', async() => {
      expect(getVaultAddresses()).toEqual(new Set([defaultVault.address]));
    });

    it('checks for vault address', async() => {
      expect(isVault(defaultVault.address)).toBe(true);
    });

    it('checks for non-vault address', async() => {
      expect(isVault(defaultAddress)).toBe(false);
    });
  });
});
