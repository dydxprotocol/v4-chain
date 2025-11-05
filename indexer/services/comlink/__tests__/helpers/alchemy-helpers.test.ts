import { addAddressesToAlchemyWebhook, registerAddressWithAlchemyWebhook } from '../../src/helpers/alchemy-helpers';
import { dbHelpers, TurnkeyUserCreateObject, TurnkeyUsersTable } from '@dydxprotocol-indexer/postgres';
import config from '../../src/config';

// Mock fetch globally
global.fetch = jest.fn();

// Mock the ZeroDev SDK functions
jest.mock('@zerodev/ecdsa-validator', () => ({
  getKernelAddressFromECDSA: jest.fn(),
}));

jest.mock('@zerodev/sdk/constants', () => ({
  getEntryPoint: jest.fn(),
  KERNEL_V3_1: '0.3.1',
}));

// Mock viem
jest.mock('viem', () => ({
  createPublicClient: jest.fn(),
  http: jest.fn(),
  Address: jest.fn(),
}));

// Mock config
jest.mock('../../src/config', () => ({
  ALCHEMY_AUTH_TOKEN: 'test-auth-token',
  ALCHEMY_WEBHOOK_ID: 'test-webhook-id',
  ALCHEMY_WEBHOOK_UPDATE_URL: 'https://dashboard.alchemy.com/api/update-webhook-addresses',
  ETHEREUM_WEBHOOK_ID: 'wh_ctbkt6y9hez91xr2',
  ARBITRUM_WEBHOOK_ID: 'wh_ltwqwcsrx1b8lgry',
  AVALANCHE_WEBHOOK_ID: 'wh_52wz9dbxywxov2dm',
  BASE_WEBHOOK_ID: 'wh_lpjn5gnwj0ll0gap',
  OPTIMISM_WEBHOOK_ID: 'wh_7eo900bsg8rkvo6z',
  SOLANA_WEBHOOK_ID: 'wh_eqxyotjv478gscpo',
}));

describe('alchemy-helpers', () => {
  const mockFetch = fetch as jest.MockedFunction<typeof fetch>;
  // const mockLogger = logger as jest.Mocked<typeof logger>;
  beforeAll(async () => {
    // Mock the database function
    await dbHelpers.clearData();
    await dbHelpers.migrate();
    const mockUser: TurnkeyUserCreateObject = {
      suborg_id: 'test-org',
      email: 'test@example.com',
      salt: 'test-salt',
      created_at: new Date().toISOString(),
      evm_address: '0x1234567890123456789012345678901234567890',
      svm_address: 'ABC123DEF456GHI789JKL012MNO345PQR678STU901VWX234YZA567',
    };
    await TurnkeyUsersTable.create(mockUser);
  });
  afterAll(async () => {
    // Mock the database function
    await dbHelpers.clearData();
    await dbHelpers.teardown();
  });

  beforeEach(() => {
    jest.clearAllMocks();
    mockFetch.mockClear();
  });

  describe('registerAddressWithAlchemyWebhook', () => {
    it('should successfully register an address with Alchemy webhook', async () => {
      const address = '0x1234567890123456789012345678901234567890';
      const webhookId = 'wh_test123';

      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      await registerAddressWithAlchemyWebhook(address, webhookId);

      expect(mockFetch).toHaveBeenCalledWith(
        config.ALCHEMY_WEBHOOK_UPDATE_URL,
        {
          method: 'PATCH',
          headers: {
            'X-Alchemy-Token': 'test-auth-token',
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            webhook_id: webhookId,
            addresses_to_add: [address],
            addresses_to_remove: [],
          }),
        },
      );
    });

    it('should throw error when webhook registration fails', async () => {
      const address = '0x1234567890123456789012345678901234567890';
      const webhookId = 'wh_test123';
      const errorResponse = '{"error": "Webhook not found"}';

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        statusText: 'Not Found',
        text: () => Promise.resolve(errorResponse),
      } as Response);

      await expect(registerAddressWithAlchemyWebhook(address, webhookId)).rejects.toThrow(
        'Failed to register address with Alchemy webhook: Not Found - {"error": "Webhook not found"}',
      );

      expect(mockFetch).toHaveBeenCalledTimes(1);
    });

    it('should handle network errors', async () => {
      const address = '0x1234567890123456789012345678901234567890';
      const webhookId = 'wh_test123';

      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(registerAddressWithAlchemyWebhook(address, webhookId)).rejects.toThrow('Network error');
    });
  });

  describe('addAddressesToAlchemyWebhook', () => {
    it('should successfully add EVM and SVM addresses to all webhooks', async () => {
      const evmAddress = '0x1234567890123456789012345678901234567890';
      const svmAddress = 'ABC123DEF456GHI789JKL012MNO345PQR678STU901VWX234YZA567';
      const smartAccountAddress = '0x9876543210987654321098765432109876543210';
      // Mock successful responses for all webhook registrations
      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      // Mock the ZeroDev SDK response
      const { getKernelAddressFromECDSA } = require('@zerodev/ecdsa-validator');
      (getKernelAddressFromECDSA as jest.Mock).mockResolvedValue(smartAccountAddress);

      // Mock the viem functions
      const { createPublicClient, http } = require('viem');
      (createPublicClient as jest.Mock).mockReturnValue({});
      (http as jest.Mock).mockReturnValue({});

      await addAddressesToAlchemyWebhook(evmAddress, svmAddress);

      // Should be called for each EVM chain + Solana
      const expectedWebhookIds = [
        'wh_ctbkt6y9hez91xr2', // mainnet
        'wh_ltwqwcsrx1b8lgry', // arbitrum
        'wh_52wz9dbxywxov2dm', // avalanche
        'wh_lpjn5gnwj0ll0gap', // base
        'wh_7eo900bsg8rkvo6z', // optimism
        'wh_eqxyotjv478gscpo', // solana
      ];

      expect(mockFetch).toHaveBeenCalledTimes(expectedWebhookIds.length);
    });

    it('should handle partial failures and continue with other webhooks', async () => {
      const evmAddress = '0x1234567890123456789012345678901234567890';
      const svmAddress = 'ABC123DEF456GHI789JKL012MNO345PQR678STU901VWX234YZA567';

      // Mock some successful and some failed responses
      mockFetch
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response) // First call succeeds
        .mockResolvedValueOnce({
          ok: false,
          status: 404,
          statusText: 'Not Found',
          text: () => Promise.resolve('{"error": "Webhook not found"}'),
        } as Response) // second network call fails
        .mockResolvedValueOnce({
          ok: false,
          status: 404,
          statusText: 'Not Found',
          text: () => Promise.resolve('{"error": "Webhook not found"}'),
        } as Response) // third network call fails
        .mockResolvedValueOnce({
          ok: false,
          status: 404,
          statusText: 'Not Found',
          text: () => Promise.resolve('{"error": "Webhook not found"}'),
        } as Response) // fourth network call fails
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response) // fifth call succeeds
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response); // sixth call (Solana) succeeds

      await addAddressesToAlchemyWebhook(evmAddress, svmAddress);

      // Should still be called for all webhooks despite some failures
      // Failed webhooks will retry 3 times each, so: 1 + 3 + 3 + 3 + 1 + 1 = 12 calls
      // Plus 1 additional call from getKernelAddressFromECDSA for Avalanche chain = 13 calls
      expect(mockFetch).toHaveBeenCalledTimes(13);
    });

    it('should handle missing EVM address', async () => {
      const svmAddress = 'ABC123DEF456GHI789JKL012MNO345PQR678STU901VWX234YZA567';

      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      await addAddressesToAlchemyWebhook(undefined, svmAddress);

      // Should only be called once for Solana
      expect(mockFetch).toHaveBeenCalledTimes(1);
      expect(mockFetch).toHaveBeenCalledWith(
        config.ALCHEMY_WEBHOOK_UPDATE_URL,
        expect.objectContaining({
          body: JSON.stringify({
            webhook_id: 'wh_eqxyotjv478gscpo',
            addresses_to_add: [svmAddress],
            addresses_to_remove: [],
          }),
        }),
      );
    });

    it('should handle missing SVM address', async () => {
      const evmAddress = '0x1234567890123456789012345678901234567890';

      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      await addAddressesToAlchemyWebhook(evmAddress, '');

      // Should only be called for EVM chains (5 times)
      expect(mockFetch).toHaveBeenCalledTimes(5);

      // Verify no Solana webhook call
      const calls = mockFetch.mock.calls;
      calls.forEach((call) => {
        const body = JSON.parse(call[1]!.body as string);
        expect(body.webhook_id).not.toBe('wh_eqxyotjv478gscpo');
      });
    });

    it('should handle retry logic correctly', async () => {
      const evmAddress = '0x1234567890123456789012345678901234567890';

      // Mock first two calls to fail, third to succeed for the first chain
      // Then all subsequent chains succeed immediately
      mockFetch
        .mockResolvedValueOnce({
          ok: false,
          status: 500,
          statusText: 'Internal Server Error',
          text: () => Promise.resolve('{"error": "Server error"}'),
        } as Response) // First attempt for chain 1 fails
        .mockResolvedValueOnce({
          ok: false,
          status: 500,
          statusText: 'Internal Server Error',
          text: () => Promise.resolve('{"error": "Server error"}'),
        } as Response) // Second attempt for chain 1 fails
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response) // Third attempt for chain 1 succeeds
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response) // Chain 2 succeeds
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response) // Chain 3 succeeds
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response) // Chain 4 succeeds
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response); // Chain 5 succeeds

      // Test the retry logic by calling the registration function directly
      // Note: This tests the retry logic through the main function
      await addAddressesToAlchemyWebhook(evmAddress, '');

      // Should be called multiple times due to retries for the first chain, then continue
      expect(mockFetch).toHaveBeenCalledTimes(7); // 3 attempts for first chain + 4 more chains
    });
  });

  describe('edge cases', () => {
    it('should handle empty addresses gracefully', async () => {
      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      await addAddressesToAlchemyWebhook('', '');

      // Should not make any webhook calls
      expect(mockFetch).not.toHaveBeenCalled();
    });

    it('should handle addresses that do not exist in database', async () => {
      const malformedAddress = 'invalid-address';

      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      await expect(addAddressesToAlchemyWebhook(malformedAddress, '')).rejects.toThrow(
        'EVM address does not exist in the database: invalid-address',
      );

      // Should not make any webhook calls
      expect(mockFetch).not.toHaveBeenCalled();
    });
  });
});
