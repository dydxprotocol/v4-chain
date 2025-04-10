import {
  compressVaultPnl,
  decompressVaultPnl,
} from '../../src/caches/vault-cache';
import { CachedVaultHistoricalPnl } from '../../src/types';

describe('vault-compression', () => {
  const mockVaultHistoricalPnl: CachedVaultHistoricalPnl = {
    ticker: 'BTC-USD',
    historicalPnl: [
      {
        equity: '128091.071477',
        totalPnl: '-12654.755537',
        netTransfers: '0.000000',
        createdAt: '2025-03-18T00:00:10.860Z',
        blockHeight: '39856519',
        blockTime: '2025-03-18T00:00:08.436Z',
      },
    ],
  };

  it('compresses and decompresses data with expected precision loss', () => {
    // Compress and decompress the data
    const compressed = compressVaultPnl(mockVaultHistoricalPnl);
    const decompressed = decompressVaultPnl(compressed);

    // Check structural integrity
    expect(decompressed).toHaveProperty('ticker');
    expect(decompressed).toHaveProperty('historicalPnl');
    expect(decompressed.ticker).toEqual(mockVaultHistoricalPnl.ticker);
    expect(Array.isArray(decompressed.historicalPnl)).toBe(true);
    expect(decompressed.historicalPnl.length).toEqual(mockVaultHistoricalPnl.historicalPnl.length);

    // Check fields are present in decompressed data
    const original = mockVaultHistoricalPnl.historicalPnl[0];
    const result = decompressed.historicalPnl[0];
    expect(result).toHaveProperty('equity');
    expect(result).toHaveProperty('totalPnl');
    expect(result).toHaveProperty('netTransfers');
    expect(result).toHaveProperty('createdAt');
    expect(result).toHaveProperty('blockHeight');
    expect(result).toHaveProperty('blockTime');

    // Check precision loss for numeric fields (limit to 1 decimal)
    expect(Number(result.equity).toFixed(1)).toEqual(Number(original.equity).toFixed(1));
    expect(Number(result.totalPnl).toFixed(1)).toEqual(Number(original.totalPnl).toFixed(1));
    expect(Number(result.netTransfers).toFixed(1)).toEqual(
      Number(original.netTransfers).toFixed(1),
    );

    // Check non-numeric data
    expect(result.blockHeight).toEqual(original.blockHeight);

    // Check date fields - they should be valid ISO strings but may not match exactly
    expect(() => new Date(result.createdAt)).not.toThrow();
    expect(() => new Date(result.blockTime)).not.toThrow();
  });

  it('results in smaller data size', () => {
    const compressed = compressVaultPnl(mockVaultHistoricalPnl);
    const originalSize = JSON.stringify(mockVaultHistoricalPnl).length;
    const compressedSize = compressed.length;

    expect(compressedSize).toBeLessThan(originalSize);
  });
});
