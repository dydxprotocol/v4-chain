import { logger } from '@dydxprotocol-indexer/base';
import { addAddressesToAlchemyWebhook, registerAddressWithAlchemyWebhook } from '../../src/helpers/alchemy-helpers';

// Mock fetch globally
global.fetch = jest.fn();

// Mock the logger
jest.mock('@dydxprotocol-indexer/base', () => ({
  logger: {
    info: jest.fn(),
    error: jest.fn(),
    warning: jest.fn(),
  },
}));

// Mock config
jest.mock('../../src/config', () => ({
  ALCHEMY_AUTH_TOKEN: 'test-auth-token',
  ALCHEMY_WEBHOOK_ID: 'test-webhook-id',
}));

describe('alchemy-helpers', () => {
  const mockFetch = fetch as jest.MockedFunction<typeof fetch>;
  const mockLogger = logger as jest.Mocked<typeof logger>;

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
        'https://dashboard.alchemy.com/api/update-webhook-addresses',
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

      expect(mockLogger.info).toHaveBeenCalledWith({
        at: 'TurnkeyController#registerAddressWithAlchemyWebhook',
        message: `Address ${address} successfully added to Alchemy webhook`,
        address,
        webhookId: 'test-webhook-id',
      });
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

      // Mock successful responses for all webhook registrations
      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      await addAddressesToAlchemyWebhook(evmAddress, svmAddress);

      // Should be called for each EVM chain + Solana
      const expectedWebhookIds = [
        'wh_ys5e0lhw2iaq0wge', // mainnet
        'wh_arbitrum',          // arbitrum
        'wh_avalanche',         // avalanche
        'wh_base',              // base
        'wh_optimism',          // optimism
        'wh_solana',            // solana
      ];

      expect(mockFetch).toHaveBeenCalledTimes(expectedWebhookIds.length);

      expect(mockLogger.info).toHaveBeenCalledWith({
        at: 'TurnkeyController#addAddressesToAlchemyWebhook',
        message: 'Successfully added addresses to Alchemy webhook',
        evmAddress,
        svmAddress,
      });
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
        .mockResolvedValue({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response); // Third call succeeds

      await addAddressesToAlchemyWebhook(evmAddress, svmAddress);

      // Should still be called for all webhooks despite some failures
      expect(mockFetch).toHaveBeenCalledTimes(8);

      // Should log errors for failed registrations
      expect(mockLogger.error).toHaveBeenCalledWith({
        at: 'TurnkeyController#addAddressesToAlchemyWebhook',
        message: expect.stringContaining('Failed to register EVM address with webhook for chain'),
        error: expect.any(Error),
        address: evmAddress,
        chainId: expect.any(String),
        webhookId: expect.any(String),
      });

      // Should still log overall success
      expect(mockLogger.info).toHaveBeenCalledWith({
        at: 'TurnkeyController#addAddressesToAlchemyWebhook',
        message: 'Successfully added addresses to Alchemy webhook',
        evmAddress,
        svmAddress,
      });
    });

    it('should handle missing EVM address', async () => {
      const svmAddress = 'ABC123DEF456GHI789JKL012MNO345PQR678STU901VWX234YZA567';

      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      await addAddressesToAlchemyWebhook('', svmAddress);

      // Should only be called once for Solana
      expect(mockFetch).toHaveBeenCalledTimes(1);
      expect(mockFetch).toHaveBeenCalledWith(
        'https://dashboard.alchemy.com/api/update-webhook-addresses',
        expect.objectContaining({
          body: JSON.stringify({
            webhook_id: 'wh_solana',
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
        expect(body.webhook_id).not.toBe('wh_solana');
      });
    });

    it('should handle retry logic correctly', async () => {
      const evmAddress = '0x1234567890123456789012345678901234567890';

      // Mock first two calls to fail, third to succeed
      mockFetch
        .mockResolvedValueOnce({
          ok: false,
          status: 500,
          statusText: 'Internal Server Error',
          text: () => Promise.resolve('{"error": "Server error"}'),
        } as Response)
        .mockResolvedValueOnce({
          ok: false,
          status: 500,
          statusText: 'Internal Server Error',
          text: () => Promise.resolve('{"error": "Server error"}'),
        } as Response)
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          statusText: 'OK',
        } as Response);

      // Test the retry logic by calling the registration function directly
      // Note: This tests the retry logic through the main function
      await addAddressesToAlchemyWebhook(evmAddress, '');

      // Should be called multiple times due to retries
      expect(mockFetch).toHaveBeenCalledTimes(15); // 5 chains * 3 attempts each
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

    it('should handle malformed addresses', async () => {
      const malformedAddress = 'invalid-address';

      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      await addAddressesToAlchemyWebhook(malformedAddress, '');

      // Should still attempt to register the malformed address
      expect(mockFetch).toHaveBeenCalledTimes(5);
    });

    it('should handle very long addresses', async () => {
      const longAddress = `0x${'1'.repeat(100)}`;

      mockFetch.mockResolvedValue({
        ok: true,
        status: 200,
        statusText: 'OK',
      } as Response);

      await addAddressesToAlchemyWebhook(longAddress, '');

      expect(mockFetch).toHaveBeenCalledTimes(5);

      // Verify the long address is included in the request
      const calls = mockFetch.mock.calls;
      calls.forEach((call) => {
        const body = JSON.parse(call[1]!.body as string);
        expect(body.addresses_to_add).toContain(longAddress);
      });
    });
  });
});
