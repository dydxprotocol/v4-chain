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
  const testPnlTick = {
    equity: '1000.1234',
    totalPnl: '100.1234',
    netTransfers: '0.1234',
    createdAt: testTimestamp,
    blockHeight: '10',
    blockTime: testTimestamp,
  };
  const testPnlTick2 = {
    equity: '12345.678',
    totalPnl: '6789.0123',
    netTransfers: '0.1234',
    createdAt: testTimestamp,
    blockHeight: '987654321',
    blockTime: testTimestamp,
  };

  const testMegavaultPnl: CachedMegavaultPnl = {
    pnlTicks: [testPnlTick, testPnlTick2],
  };

  const testVaultHistoricalPnl: CachedVaultHistoricalPnl[] = [{
    ticker: 'BTC-USD',
    historicalPnl: [testPnlTick, testPnlTick2],
  }, {
    ticker: 'ETH-USD',
    historicalPnl: [testPnlTick],
  }];

  describe('megavault cache', () => {
    it('returns null when no data is cached', async () => {
      const result = await getMegavaultPnl(resolution, client);
      expect(result).toBeNull();

      const timestamp = await getMegavaultPnlCacheTimestamp(resolution, client);
      expect(timestamp).toBeNull();
    });

    it('successfully sets and gets megavault PNL data', async () => {
      await setMegavaultPnl(resolution, testMegavaultPnl.pnlTicks, client);

      const result = await getMegavaultPnl(resolution, client);
      expect(result).not.toBeNull();

      // Verify array length
      expect(result!.pnlTicks).toHaveLength(testMegavaultPnl.pnlTicks.length);

      // Check each PNL tick
      testMegavaultPnl.pnlTicks.forEach((expectedTick, index) => {
        const actualTick = result!.pnlTicks[index];

        // Check properties exist
        expect(actualTick).toHaveProperty('equity');
        expect(actualTick).toHaveProperty('totalPnl');
        expect(actualTick).toHaveProperty('netTransfers');
        expect(actualTick).toHaveProperty('blockHeight');
        expect(actualTick).toHaveProperty('blockTime');
        expect(actualTick).toHaveProperty('createdAt');

        // Check numeric fields with rounding
        expect(Math.round(Number(actualTick.equity))).toEqual(
          Math.round(Number(expectedTick.equity)),
        );
        expect(Math.round(Number(actualTick.totalPnl))).toEqual(
          Math.round(Number(expectedTick.totalPnl)),
        );
        expect(Math.round(Number(actualTick.netTransfers))).toEqual(
          Math.round(Number(expectedTick.netTransfers)),
        );

        // Check non-numeric fields
        expect(actualTick.blockHeight).toEqual(expectedTick.blockHeight);
        expect(actualTick.blockTime).toEqual(expectedTick.blockTime);
        expect(actualTick.createdAt).toEqual(expectedTick.createdAt);
      });

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
      await setVaultsHistoricalPnl(resolution, testVaultHistoricalPnl, client);

      const result = await getVaultsHistoricalPnl(resolution, client);
      expect(result).not.toBeNull();
      expect(result!.length).toEqual(2);

      // Check both objects have correct structure
      testVaultHistoricalPnl.forEach((expected, index) => {
        const actual = result![index];

        // Check ticker
        expect(actual.ticker).toEqual(expected.ticker);

        // Check array length matches expected (BTC has 2 ticks, ETH has 1)
        expect(actual.historicalPnl.length).toEqual(expected.historicalPnl.length);

        // Check each historical PNL tick
        actual.historicalPnl.forEach((actualPnlTick, tickIndex) => {
          const expectedPnlTick = expected.historicalPnl[tickIndex];

          // Check entire structure with transformed numbers to account for rounding
          expect({
            ...actualPnlTick,
            equity: Math.round(Number(actualPnlTick.equity)).toString(),
            totalPnl: Math.round(Number(actualPnlTick.totalPnl)).toString(),
            netTransfers: Math.round(Number(actualPnlTick.netTransfers)).toString(),
          }).toEqual({
            ...expectedPnlTick,
            equity: Math.round(Number(expectedPnlTick.equity)).toString(),
            totalPnl: Math.round(Number(expectedPnlTick.totalPnl)).toString(),
            netTransfers: Math.round(Number(expectedPnlTick.netTransfers)).toString(),
          });
        });
      });

      const timestamp = await getVaultsHistoricalPnlCacheTimestamp(resolution, client);
      expect(timestamp).not.toBeNull();
      expect(timestamp!.toISOString()).toBe(testTimestamp);
    });
  });

  describe('cache timestamps', () => {
    it('returns correct timestamp for megavault', async () => {
      await setMegavaultPnl(resolution, testMegavaultPnl.pnlTicks, client);
      const timestamp = await getMegavaultPnlCacheTimestamp(resolution, client);

      expect(timestamp).not.toBeNull();
      expect(timestamp instanceof Date).toBe(true);
      expect(timestamp!.toISOString()).toBe(testTimestamp);
    });

    it('returns correct timestamp for vaults', async () => {
      await setVaultsHistoricalPnl(resolution, testVaultHistoricalPnl, client);
      const timestamp = await getVaultsHistoricalPnlCacheTimestamp(resolution, client);

      expect(timestamp).not.toBeNull();
      expect(timestamp instanceof Date).toBe(true);
      expect(timestamp!.toISOString()).toBe(testTimestamp);
    });
  });
});
