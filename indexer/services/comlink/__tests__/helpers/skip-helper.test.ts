import {
  limitAmount,
  nobleToHex,
  encodeToHexAndPad,
  nobleToSolana,
  getSlippageTolerancePercent,
} from '../../src/helpers/skip-helper';
import { getETHPrice } from '../../src/helpers/alchemy-helpers';
import { ethDenomByChainId, ETH_WEI_QUANTUM, ETH_USDC_QUANTUM } from '../../src/lib/smart-contract-constants';
import { logger } from '@dydxprotocol-indexer/base';

// Mock dependencies
jest.mock('../../src/helpers/alchemy-helpers', () => ({
  getETHPrice: jest.fn(),
}));

jest.mock('../../src/config', () => ({
  MAXIMUM_BRIDGE_AMOUNT_USDC: 99900,
  TURNKEY_API_BASE_URL: 'https://api.turnkey.com',
  TURNKEY_API_SENDER_PUBLIC_KEY: 'test-public-key',
  TURNKEY_API_SENDER_PRIVATE_KEY: 'test-private-key',
  SKIP_SLIPPAGE_TOLERANCE_USDC: 100,
  SKIP_SLIPPAGE_TOLERANCE_PERCENTAGE: '0.1',
}));

jest.mock('@dydxprotocol-indexer/base', () => ({
  logger: {
    error: jest.fn(),
  },
}));

jest.mock('@dydxprotocol-indexer/postgres', () => ({
  PermissionApprovalTable: {
    findBySuborgIdAndChainId: jest.fn(),
  },
}));

// Mock external dependencies
jest.mock('@skip-go/client/cjs', () => ({
  route: jest.fn(),
  messages: jest.fn(),
}));

jest.mock('@turnkey/sdk-server', () => ({
  Turnkey: jest.fn().mockImplementation(() => ({
    apiClient: jest.fn(() => ({})),
  })),
}));

jest.mock('@turnkey/solana', () => ({
  TurnkeySigner: jest.fn().mockImplementation(() => ({
    signTransaction: jest.fn(),
    signAllTransactions: jest.fn(),
  })),
}));

jest.mock('@turnkey/viem', () => ({
  createAccount: jest.fn(),
}));

jest.mock('@zerodev/permissions', () => ({
  deserializePermissionAccount: jest.fn(),
}));

jest.mock('@zerodev/permissions/signers', () => ({
  toECDSASigner: jest.fn(),
}));

jest.mock('@zerodev/sdk', () => ({
  CreateKernelAccountReturnType: jest.fn(),
}));

jest.mock('@zerodev/sdk/constants', () => ({
  getEntryPoint: jest.fn(() => '0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789'),
  KERNEL_V3_1: '0.3.1',
  KERNEL_V3_3: '0.3.3',
}));

jest.mock('bech32', () => ({
  decode: jest.fn(),
  fromWords: jest.fn(),
}));

jest.mock('bs58', () => ({
  encode: jest.fn(),
}));

jest.mock('viem', () => ({
  encodeFunctionData: jest.fn(),
}));

jest.mock('viem/chains', () => ({
  mainnet: { id: 1 },
  arbitrum: { id: 42161 },
  avalanche: { id: 43114 },
  base: { id: 8453 },
  optimism: { id: 10 },
}));

describe('skip-helper', () => {
  const mockGetETHPrice = getETHPrice as jest.MockedFunction<typeof getETHPrice>;
  const mockLogger = logger as jest.Mocked<typeof logger>;

  beforeEach(() => {
    jest.clearAllMocks();

    // Set up bech32 mocks
    const { decode, fromWords } = require('bech32');
    decode.mockImplementation((addr: string) => {
      if (addr.startsWith('noble')) {
        if (addr === 'noble1invalid') {
          return { prefix: 'noble', words: [] }; // Empty words to trigger length error
        }
        return { prefix: 'noble', words: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20] };
      } else if (addr.startsWith('dydx')) {
        if (addr === 'dydx1invalid') {
          return { prefix: 'dydx', words: [] }; // Empty words to trigger length error
        }
        return { prefix: 'dydx', words: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20] };
      } else if (addr.startsWith('cosmos')) {
        return { prefix: 'cosmos', words: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20] };
      } else {
        return { prefix: 'unknown', words: [] };
      }
    });

    fromWords.mockImplementation((words: number[]) => {
      // Return different buffer lengths based on the words array to simulate
      // different address lengths
      if (words.length === 0) {
        return Buffer.alloc(19, 0x01); // Wrong length to trigger error
      }
      return Buffer.alloc(20, 0x01); // Correct length
    });

    // Set up bs58 mock
    const { encode: bs58Encode } = require('bs58');
    bs58Encode.mockImplementation((_buffer: Buffer) => '1A2B3C4D5E6F7G8H9I0J1K2L3M4N5O6P7Q8R9S0T1U2V'); // 44 char base58 string
  });

  describe('limitAmount', () => {
    describe('ETH branch', () => {
      const testChainId = '1'; // mainnet
      const ethDenom = ethDenomByChainId[testChainId];

      it('should limit ETH amount when calculated max is lower than input amount', async () => {
        const inputAmount = '200000000000000000000'; // 200 ETH in wei (much larger than max)
        const mockEthPrice = 2000; // $2000 per ETH
        const expectedMaxWei = Math.floor((99900 / mockEthPrice) * ETH_WEI_QUANTUM);

        mockGetETHPrice.mockResolvedValue(mockEthPrice);

        const result = await limitAmount(testChainId, inputAmount, ethDenom);

        expect(mockGetETHPrice).toHaveBeenCalledTimes(1);
        expect(result).toBe(expectedMaxWei.toString());
      });

      it('should return original amount when input is lower than calculated max', async () => {
        const inputAmount = '1000000000000000032'; // 1 ETH in wei
        const mockEthPrice = 2000; // $2000 per ETH

        mockGetETHPrice.mockResolvedValue(mockEthPrice);

        const result = await limitAmount(testChainId, inputAmount, ethDenom);

        expect(mockGetETHPrice).toHaveBeenCalledTimes(1);
        expect(result).toBe(inputAmount);
      });

      it('should handle exact match between input and calculated max', async () => {
        const mockEthPrice = 2000; // $2000 per ETH
        const expectedMaxWei = Math.floor((99900 / mockEthPrice) * ETH_WEI_QUANTUM);
        const inputAmount = expectedMaxWei.toString();

        mockGetETHPrice.mockResolvedValue(mockEthPrice);

        const result = await limitAmount(testChainId, inputAmount, ethDenom);

        expect(mockGetETHPrice).toHaveBeenCalledTimes(1);
        expect(result).toBe(inputAmount);
      });

      it('should handle very high ETH price', async () => {
        const inputAmount = '1000000000000000003'; // 1 ETH in wei
        const mockEthPrice = 100000; // $100,000 per ETH
        const expectedMaxWei = Math.floor((99900 / mockEthPrice) * ETH_WEI_QUANTUM);

        mockGetETHPrice.mockResolvedValue(mockEthPrice);

        const result = await limitAmount(testChainId, inputAmount, ethDenom);

        expect(result).toBe(expectedMaxWei.toString());
      });

      it('should handle very low ETH price', async () => {
        const inputAmount = '1000000000000000000'; // 1 ETH in wei
        const mockEthPrice = 1; // $1 per ETH

        mockGetETHPrice.mockResolvedValue(mockEthPrice);

        const result = await limitAmount(testChainId, inputAmount, ethDenom);

        expect(result).toBe(inputAmount); // Should be limited by input amount
      });

      it('should throw error and log when getETHPrice fails', async () => {
        const inputAmount = '1000000000000000000';
        const error = new Error('API failure');

        mockGetETHPrice.mockRejectedValue(error);

        await expect(limitAmount(testChainId, inputAmount, ethDenom)).rejects.toThrow('API failure');

        expect(mockLogger.error).toHaveBeenCalledWith({
          at: 'skip-helper#limitAmount',
          message: 'Failed to get ETH price',
          error,
        });
      });

      it('should handle zero ETH price', async () => {
        const inputAmount = '1000000000000000000';
        const mockEthPrice = 0;

        mockGetETHPrice.mockResolvedValue(mockEthPrice);

        // When price is 0, maxDepositInWei becomes Infinity, which can't be converted to BigInt
        await expect(limitAmount(testChainId, inputAmount, ethDenom)).rejects.toThrow();
      });

      it('should handle negative ETH price', async () => {
        const inputAmount = '1000000000000000000';
        const mockEthPrice = -1000;

        mockGetETHPrice.mockResolvedValue(mockEthPrice);

        const result = await limitAmount(testChainId, inputAmount, ethDenom);

        // When price is negative, maxDepositInWei becomes negative,
        // so min should return the negative value
        const expectedMaxWei = Math.floor((99900 / mockEthPrice) * ETH_WEI_QUANTUM);
        expect(result).toBe(expectedMaxWei.toString());
      });
    });

    describe('USDC branch', () => {
      const testChainId = '1';
      const usdcDenom = 'usdc-address';

      it('should limit USDC amount when input exceeds maximum', async () => {
        const inputAmount = '200000000000'; // 200,000 USDC (exceeds 99,900 limit)
        const expectedMax = 99900 * ETH_USDC_QUANTUM;

        const result = await limitAmount(testChainId, inputAmount, usdcDenom);

        expect(result).toBe(expectedMax.toString());
      });

      it('should return original amount when input is within limit', async () => {
        const inputAmount = '50000000032'; // 50,000 USDC (within 99,900 limit)

        const result = await limitAmount(testChainId, inputAmount, usdcDenom);

        expect(result).toBe(inputAmount);
      });

      it('should handle exact maximum amount', async () => {
        const inputAmount = (99900 * ETH_USDC_QUANTUM).toString();

        const result = await limitAmount(testChainId, inputAmount, usdcDenom);

        expect(result).toBe(inputAmount);
      });

      it('should handle zero amount', async () => {
        const inputAmount = '0';

        const result = await limitAmount(testChainId, inputAmount, usdcDenom);

        expect(result).toBe('0');
      });

      it('should handle very large amount', async () => {
        const inputAmount = '999999999999999999999999999'; // Very large number
        const expectedMax = 99900 * ETH_USDC_QUANTUM;

        const result = await limitAmount(testChainId, inputAmount, usdcDenom);

        expect(result).toBe(expectedMax.toString());
      });

      it('should handle string amount that parses to negative number', async () => {
        const inputAmount = '-10000000000'; // Negative amount

        const result = await limitAmount(testChainId, inputAmount, usdcDenom);

        // min() should return the negative value since it's smaller than the positive max
        expect(result).toBe(inputAmount);
      });
    });

    describe('edge cases', () => {
      it('should handle empty string amount', async () => {
        // BigInt('') returns 0n, so it should work
        const result = await limitAmount('1', '', 'usdc-address');
        expect(result).toBe('0');
      });

      it('should handle non-numeric string amount', async () => {
        // BigInt('not-a-number') will throw, so we expect it to throw
        await expect(limitAmount('1', 'not-a-number', 'usdc-address')).rejects.toThrow();
      });

      it('should handle different chain IDs for ETH', async () => {
        const arbitrumChainId = '42161';
        const arbitrumEthDenom = ethDenomByChainId[arbitrumChainId];
        const inputAmount = '1000000000000000000';
        const mockEthPrice = 2000;

        mockGetETHPrice.mockResolvedValue(mockEthPrice);

        const result = await limitAmount(arbitrumChainId, inputAmount, arbitrumEthDenom);

        expect(mockGetETHPrice).toHaveBeenCalledTimes(1);
        expect(result).toBe(inputAmount);
      });
    });
  });

  describe('nobleToHex', () => {
    it('should convert noble address to hex with proper padding', () => {
      // This is a mock noble address - in real usage it would be a valid bech32 address
      const nobleAddress = 'noble1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu';

      const result = nobleToHex(nobleAddress);

      expect(result).toMatch(/^0x[0-9a-f]{64}$/); // 32 bytes = 64 hex chars
      expect(result.startsWith('0x000000000000000000000000')).toBe(true); // Should be padded with 12 zero bytes
    });

    it('should convert dydx address to hex with proper padding', () => {
      const dydxAddress = 'dydx1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu';

      const result = nobleToHex(dydxAddress);

      expect(result).toMatch(/^0x[0-9a-f]{64}$/);
      expect(result.startsWith('0x000000000000000000000000')).toBe(true);
    });

    it('should throw error for invalid HRP', () => {
      const invalidAddress = 'cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu';

      expect(() => nobleToHex(invalidAddress)).toThrow('Invalid HRP: expected "noble", got "cosmos"');
    });

    it('should throw error for invalid address length', () => {
      // Use an address that will result in empty words array
      const invalidAddress = 'noble1invalid'; // This will trigger empty words

      expect(() => nobleToHex(invalidAddress)).toThrow();
    });
  });

  describe('encodeToHexAndPad', () => {
    it('should encode string to hex and pad with 277 zero bytes', () => {
      const input = 'test string';
      const result = encodeToHexAndPad(input);

      expect(result).toMatch(/^0x[0-9a-f]+$/);
      expect(result.length).toBe(2 + 277 * 2 + input.length * 2); // 0x + 277*2 + input hex length
    });

    it('should handle empty string', () => {
      const input = '';
      const result = encodeToHexAndPad(input);

      expect(result).toMatch(/^0x[0-9a-f]+$/);
      expect(result.length).toBe(2 + 277 * 2); // 0x + 277*2 + 0
    });

    it('should handle special characters', () => {
      const input = '!@#$%^&*()';
      const result = encodeToHexAndPad(input);

      expect(result).toMatch(/^0x[0-9a-f]+$/);
      expect(result.length).toBe(2 + 277 * 2 + input.length * 2);
    });

    it('should handle unicode characters', () => {
      const input = '测试字符串';
      const result = encodeToHexAndPad(input);

      expect(result).toMatch(/^0x[0-9a-f]+$/);
      expect(result.length).toBe(2 + 277 * 2 + Buffer.byteLength(input) * 2);
    });

    it('should handle very long string', () => {
      const input = 'a'.repeat(1000);
      const result = encodeToHexAndPad(input);

      expect(result).toMatch(/^0x[0-9a-f]+$/);
      expect(result.length).toBe(2 + 277 * 2 + input.length * 2);
    });
  });

  describe('nobleToSolana', () => {
    it('should convert noble address to Solana base58 with proper padding', () => {
      const nobleAddress = 'noble1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu';

      const result = nobleToSolana(nobleAddress);

      expect(result).toMatch(/^[1-9A-Za-z0]{44}$/); // Base58 format, 32 bytes = ~44 chars
      expect(typeof result).toBe('string');
    });

    it('should throw error for invalid HRP', () => {
      const invalidAddress = 'dydx1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu';

      expect(() => nobleToSolana(invalidAddress)).toThrow('Unexpected HRP "dydx". Expected "noble".');
    });

    it('should throw error for invalid payload length', () => {
      const invalidAddress = 'noble1invalid'; // This will trigger empty words

      expect(() => nobleToSolana(invalidAddress)).toThrow();
    });

    it('should handle case insensitive input', () => {
      const nobleAddress = 'NOBLE1QYPQXPQ9QCRSSZG2PVXQ6RS0ZQG3YYC5LZV7XU';

      const result = nobleToSolana(nobleAddress);

      expect(result).toMatch(/^[1-9A-Za-z0]{44}$/);
    });
  });

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
