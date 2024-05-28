import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { clear, getLatestBlockHeight, updateBlockHeight } from '../../src/loops/block-height-refresher';
import { defaultBlock2 } from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';
import config from '../../src/config';

describe('blockHeightRefresher', () => {
  beforeAll(async () => {
    await migrate();
    await seedData();
    await updateBlockHeight();
  });

  afterAll(async () => {
    await clearData();
    await teardown();
  });

  describe('getLatestBlockHeight', () => {
    it('successfully gets the latest block height after update', async () => {
      await updateBlockHeight();
      expect(getLatestBlockHeight()).toBe(defaultBlock2.blockHeight);
    });
  });

  describe('clear', () => {
    it('throws an error if block height does not exist', () => {
      clear();
      expect(() => getLatestBlockHeight()).toThrowError('Unable to find latest block height');
    });

    it('throws an error when clear is called in non-test environment', () => {
      const originalEnv = config.NODE_ENV;
      config.NODE_ENV = 'production';
      expect(() => clear()).toThrowError('clear cannot be used in non-test env');
      config.NODE_ENV = originalEnv;
    });
  });
});
