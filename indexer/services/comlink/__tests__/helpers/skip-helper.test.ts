import { getSlippageTolerancePercent } from '../../src/helpers/skip-helper';

// Mock config
jest.mock('../../src/config', () => ({
  SKIP_SLIPPAGE_TOLERANCE_USDC: 100,
  SKIP_SLIPPAGE_TOLERANCE_PERCENTAGE: '0.1', // 10%
}));

describe('skip-helper', () => {
  describe('getSlippageTolerancePercent', () => {
    it('should return the percentage-based tolerance when it is smaller than USDC-based tolerance', () => {
      // When estimatedAmountOut is large, USDC-based tolerance becomes very small
      const estimatedAmountOut = '1000000000'; // 1000 USDC (assuming 6 decimals)
      // USDC-based: 100 * 100000000 / 1000000000 = 0.1 (10%)
      const result = getSlippageTolerancePercent(estimatedAmountOut);

      // Should return the smaller value: 0.1% default.
      expect(result).toBe('0.1');
    });

    it('should return the USDC-based tolerance when it is smaller than percentage-based tolerance', () => {
      // When estimatedAmountOut is small, USDC-based tolerance becomes large
      const estimatedAmountOut = '1000000000000'; // 1,000,000 USDC (1 million with 6 decimals)
      const result = getSlippageTolerancePercent(estimatedAmountOut);

      // USDC-based: 10_000_000_000 / 1_000_000_000_000 = 0.01 (1%)
      // Percentage-based: 0.1 (10%)
      // Should return the smaller value: 0.01 (1%)
      expect(result).toBe('0.01');
    });
  });
});
