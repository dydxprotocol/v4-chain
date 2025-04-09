import {
  getMegavaultPnl,
  setMegavaultPnl,
  getMegavaultPnlCacheTimestamp,
  getVaultsHistoricalPnl,
  setVaultsHistoricalPnl,
  getVaultsHistoricalPnlCacheTimestamp,
} from '../../src/caches/vault-cache';
import { deleteAllAsync } from '../../src/helpers/redis';
import { CachedMegavaultPnl, CachedVaultHistoricalPnl } from '../../src/types';
import { redis as client } from '../helpers/utils';

const testTimestamp = '2025-04-01T00:00:00.000Z';

describe('vault-cache', () => {
  beforeEach(async () => {
    await deleteAllAsync(client);
    jest.useFakeTimers();
    // Set a fixed timestamp for tests
    jest.setSystemTime(new Date(testTimestamp));
  });

  afterEach(async () => {
    await deleteAllAsync(client);
    jest.useRealTimers();
  });

  const resolution = '1H';
  const mockPnlTick = {
    id: 'id1',
    subaccountId: 'subaccount1',
    equity: '1000',
    totalPnl: '100',
    netTransfers: '0',
    createdAt: testTimestamp,
    blockHeight: '10',
    blockTime: testTimestamp,
  };

  const mockMegavaultPnl: CachedMegavaultPnl = {
    pnlTicks: [mockPnlTick],
    lastUpdated: testTimestamp,
  };

  const mockVaultHistoricalPnl: CachedVaultHistoricalPnl[] = [{
    ticker: 'BTC-USD',
    historicalPnl: [{
      equity: mockPnlTick.equity,
      totalPnl: mockPnlTick.totalPnl,
      netTransfers: mockPnlTick.netTransfers,
      createdAt: mockPnlTick.createdAt,
      blockHeight: mockPnlTick.blockHeight,
      blockTime: mockPnlTick.blockTime,
    }],
  }];

  describe('megavault cache', () => {
    it('returns null when no data is cached', async () => {
      const result = await getMegavaultPnl(resolution, client);
      expect(result).toBeNull();

      const timestamp = await getMegavaultPnlCacheTimestamp(resolution, client);
      expect(timestamp).toBeNull();
    });

    it('successfully sets and gets megavault PNL data', async () => {
      await setMegavaultPnl(resolution, mockMegavaultPnl.pnlTicks, client);

      const result = await getMegavaultPnl(resolution, client);
      expect(result).not.toBeNull();
      expect(result!.pnlTicks).toEqual(mockMegavaultPnl.pnlTicks);

      const timestamp = await getMegavaultPnlCacheTimestamp(resolution, client);
      expect(timestamp).not.toBeNull();
      expect(timestamp!.toISOString()).toBe(testTimestamp);
    });
  });

  describe('vaults historical PNL cache', () => {
    it('returns null when no data is cached', async () => {
      const result = await getVaultsHistoricalPnl(resolution, client);
      expect(result).toBeNull();

      const timestamp = await getVaultsHistoricalPnlCacheTimestamp(resolution, client);
      expect(timestamp).toBeNull();
    });

    it('successfully sets and gets vaults historical PNL data', async () => {
      await setVaultsHistoricalPnl(resolution, mockVaultHistoricalPnl, client);

      const result = await getVaultsHistoricalPnl(resolution, client);
      expect(result).not.toBeNull();
      expect(result).toEqual(mockVaultHistoricalPnl);

      const timestamp = await getVaultsHistoricalPnlCacheTimestamp(resolution, client);
      expect(timestamp).not.toBeNull();
      expect(timestamp!.toISOString()).toBe(testTimestamp);
    });
  });

  describe('cache timestamps', () => {
    it('returns correct timestamp for megavault', async () => {
      await setMegavaultPnl(resolution, mockMegavaultPnl.pnlTicks, client);
      const timestamp = await getMegavaultPnlCacheTimestamp(resolution, client);

      expect(timestamp).not.toBeNull();
      expect(timestamp instanceof Date).toBe(true);
      expect(timestamp!.toISOString()).toBe(testTimestamp);
    });

    it('returns correct timestamp for vaults', async () => {
      await setVaultsHistoricalPnl(resolution, mockVaultHistoricalPnl, client);
      const timestamp = await getVaultsHistoricalPnlCacheTimestamp(resolution, client);

      expect(timestamp).not.toBeNull();
      expect(timestamp instanceof Date).toBe(true);
      expect(timestamp!.toISOString()).toBe(testTimestamp);
    });
  });
});
