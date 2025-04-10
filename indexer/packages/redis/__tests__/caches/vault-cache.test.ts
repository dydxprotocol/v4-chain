import {
  getMegavaultPnl,
  setMegavaultPnl,
  getMegavaultPnlCacheTimestamp,
  getVaultsHistoricalPnl,
  setVaultsHistoricalPnl,
  getVaultsHistoricalPnlCacheTimestamp,
  compressVaultPnl,
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
    equity: '1000.1234',
    totalPnl: '100.1234',
    netTransfers: '0.1234',
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

      // Check structure and ticker match
      expect(result![0]).toHaveProperty('ticker', 'BTC-USD');
      expect(result![0]).toHaveProperty('historicalPnl');
      expect(result![0].historicalPnl.length).toEqual(1);

      // Check field existence
      const pnlData = result![0].historicalPnl[0];
      expect(pnlData).toHaveProperty('equity');
      expect(pnlData).toHaveProperty('totalPnl');
      expect(pnlData).toHaveProperty('netTransfers');
      expect(pnlData).toHaveProperty('createdAt');
      expect(pnlData).toHaveProperty('blockHeight');
      expect(pnlData).toHaveProperty('blockTime');

      // Check blockHeight (should match exactly)
      expect(pnlData.blockHeight).toEqual(mockVaultHistoricalPnl[0].historicalPnl[0].blockHeight);

      // Check numeric fields with expected precision loss (1 decimal place)
      expect(Number(pnlData.equity).toFixed(1)).toEqual(
        Number(mockVaultHistoricalPnl[0].historicalPnl[0].equity).toFixed(1),
      );
      expect(Number(pnlData.totalPnl).toFixed(1)).toEqual(
        Number(mockVaultHistoricalPnl[0].historicalPnl[0].totalPnl).toFixed(1),
      );
      expect(Number(pnlData.netTransfers).toFixed(1)).toEqual(
        Number(mockVaultHistoricalPnl[0].historicalPnl[0].netTransfers).toFixed(1),
      );

      const timestamp = await getVaultsHistoricalPnlCacheTimestamp(resolution, client);
      expect(timestamp).not.toBeNull();
      expect(timestamp!.toISOString()).toBe(testTimestamp);
    });

    it('creates a compressed string with the expected format', () => {
      const singleVault = mockVaultHistoricalPnl[0];
      const compressedString = compressVaultPnl(singleVault);
      
      // Check that it's a valid JSON string
      expect(() => JSON.parse(compressedString)).not.toThrow();
      
      // Parse the compressed data
      const parsed = JSON.parse(compressedString);
      
      // Check that it has the expected structure
      expect(Array.isArray(parsed)).toBe(true);
      expect(parsed.length).toBe(2);
      expect(parsed[0]).toBe('BTC-USD'); // ticker
      expect(Array.isArray(parsed[1])).toBe(true); // array of historical entries
      expect(parsed[1].length).toBe(1); // single historical entry
      
      // Check the format of the historical data entry
      const entry = parsed[1][0];
      expect(Array.isArray(entry)).toBe(true);
      expect(entry.length).toBe(6);
      
      // Check the types of each element in the entry
      expect(typeof entry[0]).toBe('string'); // equity as string with 1 decimal
      expect(typeof entry[1]).toBe('string'); // totalPnl as string with 1 decimal
      expect(typeof entry[2]).toBe('string'); // netTransfers as string with 1 decimal
      expect(typeof entry[3]).toBe('number'); // createdAt timestamp as number
      expect(typeof entry[4]).toBe('number'); // blockHeight as number
      expect(typeof entry[5]).toBe('number'); // blockTime timestamp as number
      
      // Check specific values
      expect(entry[0]).toBe('1000.1'); // equity with 1 decimal place
      expect(entry[1]).toBe('100.1'); // totalPnl with 1 decimal place
      expect(entry[2]).toBe('0.1'); // netTransfers with 1 decimal place
      expect(entry[4]).toBe(10); // blockHeight as number
      
      // Check timestamps
      const createdAtDate = new Date(entry[3] * 1000);
      const blockTimeDate = new Date(entry[5] * 1000);
      expect(createdAtDate.toISOString()).toBe(testTimestamp);
      expect(blockTimeDate.toISOString()).toBe(testTimestamp);
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
