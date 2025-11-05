import { GeoOriginStatus } from '@dydxprotocol-indexer/compliance';
import {
  dbHelpers,
  PermissionApprovalTable,
  TurnkeyUserCreateObject,
  TurnkeyUsersTable,
  BridgeInformationTable,
  BridgeInformationCreateObject,
} from '@dydxprotocol-indexer/postgres';
import { stats } from '@dydxprotocol-indexer/base';
import request from 'supertest';

jest.mock('@skip-go/client/cjs', () => ({
  __esModule: true,
  balances: jest.fn(),
  route: jest.fn(),
  executeRoute: jest.fn(),
  setClientOptions: jest.fn(),
  TransferStatus: {},
}));
jest.mock('@zerodev/sdk', () => ({
  __esModule: true,
  createKernelAccount: jest.fn(),
  createZeroDevPaymasterClient: jest.fn(),
  getUserOperationGasPrice: jest.fn(),
  createKernelAccountClient: jest.fn(),
  KERNEL_V3_3: 'KERNEL_V3_3',
  KERNEL_V3_1: 'KERNEL_V3_1',
}));
jest.mock('@zerodev/permissions/signers', () => ({
  __esModule: true,
  toECDSASigner: jest.fn(),
}));
jest.mock('@zerodev/permissions', () => ({
  __esModule: true,
  deserializePermissionAccount: jest.fn(),
}));
jest.mock('@turnkey/viem', () => ({
  __esModule: true,
  createAccount: jest.fn(),
}));
jest.mock('viem/accounts', () => ({
  __esModule: true,
  privateKeyToAccount: jest.fn(),
}));
import * as skipClient from '@skip-go/client/cjs';
import * as zeroDev from '@zerodev/sdk';
import * as turnkeyViem from '@turnkey/viem';
import * as zerodevPermissions from '@zerodev/permissions';
import * as zerodevPermissionsSigners from '@zerodev/permissions/signers';
import { privateKeyToAccount } from 'viem/accounts';

import { alchemyNetworkToChainIdMap } from '../../../../src/helpers/alchemy-helpers';
import * as skipHelpers from '../../../../src/helpers/skip-helper';
import { RequestMethod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';

const geoOriginHeaders = {
  'geo-origin-country': 'AR', // Argentina
  'geo-origin-region': 'AR-V', // Tierra del Fuego
  'geo-origin-status': GeoOriginStatus.OK,
};

describe('skip-bridge-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  describe('GET /getDepositAddress/:dydxAddress', () => {
    const testDydxAddress2 = 'dydx1234567890123456789012345678901234567890';
    const invalidDydxAddress = 'invalid-address';

    const mockTurnkeyUser2: TurnkeyUserCreateObject = {
      suborg_id: 'test-suborg-id-2',
      svm_address: 'svm1234567890123456789012345678901234567890',
      evm_address: '0x1234567890123456789012345678901234567890',
      smart_account_address: '0x9876543210987654321098765432109876543210',
      salt: 'test-salt-2',
      dydx_address: testDydxAddress2,
      created_at: new Date().toISOString(),
    };

    describe('Success cases', () => {
      it('should return deposit addresses for existing user', async () => {
        // Create a user in the database
        await TurnkeyUsersTable.create(mockTurnkeyUser2);

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDepositAddress/${testDydxAddress2}`,
          expectedStatus: 200,
          headers: geoOriginHeaders,
        });

        expect(response.body).toEqual({
          evmAddress: mockTurnkeyUser2.evm_address,
          avalancheAddress: mockTurnkeyUser2.smart_account_address,
          svmAddress: mockTurnkeyUser2.svm_address,
        });

        // Verify stats were called
        expect(stats.timing).toHaveBeenCalledWith(
          expect.stringContaining('bridging-controller.get_deposit_address.timing'),
          expect.any(Number),
        );
      });
    });

    describe('Error cases', () => {
      it('should return 404 for non-existent user (not test address)', async () => {
        const nonExistentAddress = 'dydx1nonexistentaddress123456789012345678';

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDepositAddress/${nonExistentAddress}`,
          expectedStatus: 404,
        });

        expect(response.body).toEqual({
          error: 'User not found',
          message: `No user found with dydx address: ${nonExistentAddress}`,
        });
      });

      it('should return 400 for invalid dydx address format', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDepositAddress/${invalidDydxAddress}`,
          expectedStatus: 400,
        });

        // The validation middleware should catch this
        expect(response.body.errors).toBeDefined();
        expect(response.body.errors[0].msg).toContain('address must be a valid dydx address');
      });

      it('should return 400 for empty dydx address parameter', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: '/v4/bridging/getDepositAddress/',
          expectedStatus: 404, // Express router will return 404 for missing param
        });

        expect(response.body.error).toBe('Not Found');
      });
    });

    describe('Response format validation', () => {
      it('should return response with correct field names and types', async () => {
        await TurnkeyUsersTable.create(mockTurnkeyUser2);

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDepositAddress/${testDydxAddress2}`,
          expectedStatus: 200,
        });

        // Verify response structure
        expect(response.body).toHaveProperty('evmAddress');
        expect(response.body).toHaveProperty('avalancheAddress');
        expect(response.body).toHaveProperty('svmAddress');

        // Verify field types
        expect(typeof response.body.evmAddress).toBe('string');
        expect(typeof response.body.avalancheAddress).toBe('string');
        expect(typeof response.body.svmAddress).toBe('string');

        // Verify no extra fields
        const expectedKeys = ['evmAddress', 'avalancheAddress', 'svmAddress'];
        const actualKeys = Object.keys(response.body);
        expect(actualKeys.sort()).toEqual(expectedKeys.sort());
      });

      it('should handle null/undefined values in database gracefully', async () => {
        // Create user with some null values
        const userWithNulls = {
          ...mockTurnkeyUser2,
          smart_account_address: undefined,
        };
        await TurnkeyUsersTable.create(userWithNulls);

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDepositAddress/${testDydxAddress2}`,
          expectedStatus: 200,
        });

        expect(response.body.evmAddress).toBe(userWithNulls.evm_address);
        expect(response.body.avalancheAddress).toBeNull();
        expect(response.body.svmAddress).toBe(userWithNulls.svm_address);
      });
    });
  });

  describe('GET /getDeposits/:dydxAddress', () => {
    const testDydxAddress = 'dydx13s025jax0dzw4kuan9hc7e2qkkmcuaz6y93hky';
    const invalidDydxAddress = 'invalid-address';

    const mockTurnkeyUser: TurnkeyUserCreateObject = {
      suborg_id: 'test-suborg-deposits',
      svm_address: 'svm1234567890123456789012345678901234567890',
      evm_address: '0x1234567890123456789012345678901234567890',
      smart_account_address: '0x9876543210987654321098765432109876543210',
      salt: 'test-salt-deposits',
      dydx_address: testDydxAddress,
      created_at: new Date().toISOString(),
    };

    const mockBridgeInfo1: BridgeInformationCreateObject = {
      from_address: mockTurnkeyUser.evm_address,
      chain_id: 'ethereum',
      amount: '1000000',
      transaction_hash: '0xabc123',
      created_at: '2025-09-18T10:00:00.000Z',
    };

    const mockBridgeInfo2: BridgeInformationCreateObject = {
      from_address: mockTurnkeyUser.smart_account_address!,
      chain_id: 'polygon',
      amount: '2000000',
      transaction_hash: '0xdef456',
      created_at: '2025-09-18T12:00:00.000Z',
    };

    const mockBridgeInfo3: BridgeInformationCreateObject = {
      from_address: mockTurnkeyUser.svm_address,
      chain_id: 'solana',
      amount: '3000000',
      created_at: '2025-09-17T08:00:00.000Z', // Before today
    };

    describe('Success cases', () => {
      beforeEach(async () => {
        await TurnkeyUsersTable.create(mockTurnkeyUser);
        await BridgeInformationTable.create(mockBridgeInfo1);
        await BridgeInformationTable.create(mockBridgeInfo2);
        await BridgeInformationTable.create(mockBridgeInfo3);
      });

      it('should return all deposits for existing user without date filter', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}`,
          expectedStatus: 200,
        });

        expect(response.body.deposits).toBeDefined();
        expect(response.body.deposits.results).toHaveLength(3);
        expect(response.body.total).toBe(3);

        // Verify all user's addresses are included
        const fromAddresses = response.body.deposits.results.map((d: any) => d.from_address);
        expect(fromAddresses).toContain(mockTurnkeyUser.evm_address);
        expect(fromAddresses).toContain(mockTurnkeyUser.smart_account_address);
        expect(fromAddresses).toContain(mockTurnkeyUser.svm_address);
      });

      it('should filter deposits by createdOnOrAfter date', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?createdOnOrAfter=2025-09-18T00:00:00.000Z`,
          expectedStatus: 200,
        });

        expect(response.body.deposits.results).toHaveLength(2);
        expect(response.body.since).toBe('2025-09-18T00:00:00.000Z');
        expect(response.body.total).toBe(2);

        // Should only include deposits from today (mockBridgeInfo1 and mockBridgeInfo2)
        const amounts = response.body.deposits.results.map((d: any) => d.amount);
        expect(amounts).toContain('1000000');
        expect(amounts).toContain('2000000');
        expect(amounts).not.toContain('3000000'); // This one is from yesterday
      });

      it('should support pagination with limit and page', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?limit=2&page=1`,
          expectedStatus: 200,
        });

        expect(response.body.deposits.results).toHaveLength(2);
        expect(response.body.deposits.limit).toBe(2);
        expect(response.body.deposits.offset).toBe(0);
        expect(response.body.deposits.total).toBe(3);

        // Test page 2
        const page2Response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?limit=2&page=2`,
          expectedStatus: 200,
        });

        expect(page2Response.body.deposits.results).toHaveLength(1);
        expect(page2Response.body.deposits.offset).toBe(2);
      });

      it('should support limit without pagination', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?limit=2`,
          expectedStatus: 200,
        });

        expect(response.body.deposits.results).toHaveLength(2);
        expect(response.body.deposits.limit).toBeUndefined();
        expect(response.body.deposits.offset).toBeUndefined();
        expect(response.body.deposits.total).toBeUndefined();
      });

      it('should combine date filter and pagination', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?createdOnOrAfter=2025-09-18T00:00:00.000Z&limit=1&page=1`,
          expectedStatus: 200,
        });

        expect(response.body.deposits.results).toHaveLength(1);
        expect(response.body.since).toBe('2025-09-18T00:00:00.000Z');
        expect(response.body.deposits.total).toBe(2); // Only 2 records match the date filter
      });
    });

    describe('Error cases', () => {
      it('should return 404 for non-existent user', async () => {
        const nonExistentAddress = 'dydx1nonexistentaddress123456789012345678';

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${nonExistentAddress}`,
          expectedStatus: 404,
        });

        expect(response.body).toEqual({
          error: 'User not found',
          message: `No user found with dydx address: ${nonExistentAddress}`,
        });
      });

      it('should return 404 for user without smart_account_address', async () => {
        const userWithoutSmartAccount = {
          ...mockTurnkeyUser,
          dydx_address: 'dydx1testwithoutsmartaccount123456789012345',
          smart_account_address: undefined,
        };
        await TurnkeyUsersTable.create(userWithoutSmartAccount);

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${userWithoutSmartAccount.dydx_address}`,
          expectedStatus: 404,
        });

        expect(response.body).toEqual({
          error: 'User not found',
          message: `No user found with dydx address: ${userWithoutSmartAccount.dydx_address}`,
        });
      });

      it('should return 400 for invalid dydx address format', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${invalidDydxAddress}`,
          expectedStatus: 400,
        });

        expect(response.body.errors).toBeDefined();
        expect(response.body.errors[0].msg).toContain('address must be a valid dydx address');
      });

      it('should return 400 for invalid date format', async () => {
        await TurnkeyUsersTable.create(mockTurnkeyUser);

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?createdOnOrAfter=invalid-date`,
          expectedStatus: 400,
        });

        expect(response.body.errors).toBeDefined();
      });

      it('should return 400 for invalid pagination parameters', async () => {
        await TurnkeyUsersTable.create(mockTurnkeyUser);

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?limit=invalid&page=invalid`,
          expectedStatus: 400,
        });

        expect(response.body.errors).toBeDefined();
      });
    });

    describe('Response format validation', () => {
      beforeEach(async () => {
        await TurnkeyUsersTable.create(mockTurnkeyUser);
        await BridgeInformationTable.create(mockBridgeInfo1);
      });

      it('should return response with correct structure', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}`,
          expectedStatus: 200,
        });

        // Verify response structure
        expect(response.body).toHaveProperty('deposits');
        expect(response.body).toHaveProperty('total');
        expect(response.body.deposits).toHaveProperty('results');
        expect(Array.isArray(response.body.deposits.results)).toBe(true);

        // Verify deposit record structure
        const deposit = response.body.deposits.results[0];
        expect(deposit).toHaveProperty('id');
        expect(deposit).toHaveProperty('from_address');
        expect(deposit).toHaveProperty('chain_id');
        expect(deposit).toHaveProperty('amount');
        expect(deposit).toHaveProperty('transaction_hash');
        expect(deposit).toHaveProperty('created_at');

        // Verify field types
        expect(typeof deposit.id).toBe('string');
        expect(typeof deposit.from_address).toBe('string');
        expect(typeof deposit.chain_id).toBe('string');
        expect(typeof deposit.amount).toBe('string');
        expect(typeof deposit.created_at).toBe('string');
      });

      it('should include since field when date filter is provided', async () => {
        const sinceDate = '2025-09-18T00:00:00.000Z';
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?createdOnOrAfter=${sinceDate}`,
          expectedStatus: 200,
        });

        expect(response.body.since).toBe(sinceDate);
      });

      it('should not include since field when no date filter is provided', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}`,
          expectedStatus: 200,
        });

        expect(response.body.since).toBeUndefined();
      });

      it('should return deposits in descending order by created_at', async () => {
        // Create additional record with later timestamp
        await BridgeInformationTable.create({
          ...mockBridgeInfo2,
          created_at: '2025-09-18T15:00:00.000Z',
        });

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}`,
          expectedStatus: 200,
        });

        const results = response.body.deposits.results;
        expect(results).toHaveLength(2);

        // Should be in descending order (newest first)
        expect(new Date(results[0].created_at).getTime()).toBeGreaterThan(
          new Date(results[1].created_at).getTime(),
        );
      });
    });

    describe('Edge cases', () => {
      beforeEach(async () => {
        await TurnkeyUsersTable.create(mockTurnkeyUser);
      });

      it('should return empty results when user has no deposits', async () => {
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}`,
          expectedStatus: 200,
        });

        expect(response.body.deposits.results).toHaveLength(0);
        expect(response.body.total).toBe(0);
      });

      it('should return empty results when date filter excludes all deposits', async () => {
        await BridgeInformationTable.create(mockBridgeInfo3); // From 2025-09-17

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?createdOnOrAfter=2025-09-19T00:00:00.000Z`,
          expectedStatus: 200,
        });

        expect(response.body.deposits.results).toHaveLength(0);
        expect(response.body.total).toBe(0);
        expect(response.body.since).toBe('2025-09-19T00:00:00.000Z');
      });

      it('should handle deposits with null transaction_hash', async () => {
        const bridgeInfoWithoutTx = {
          ...mockBridgeInfo1,
          transaction_hash: undefined,
        };
        await BridgeInformationTable.create(bridgeInfoWithoutTx);

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}`,
          expectedStatus: 200,
        });

        expect(response.body.deposits.results).toHaveLength(1);
        expect(response.body.deposits.results[0].transaction_hash).toBeNull();
      });

      it('should handle very old date filter', async () => {
        await BridgeInformationTable.create(mockBridgeInfo1);

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?createdOnOrAfter=2020-01-01T00:00:00.000Z`,
          expectedStatus: 200,
        });

        expect(response.body.deposits.results).toHaveLength(1);
        expect(response.body.since).toBe('2020-01-01T00:00:00.000Z');
      });

      it('should handle future date filter', async () => {
        await BridgeInformationTable.create(mockBridgeInfo1);

        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}?createdOnOrAfter=2030-01-01T00:00:00.000Z`,
          expectedStatus: 200,
        });

        expect(response.body.deposits.results).toHaveLength(0);
        expect(response.body.total).toBe(0);
      });
    });

    describe('Performance and stats', () => {
      beforeEach(async () => {
        await TurnkeyUsersTable.create(mockTurnkeyUser);
      });

      it('should record timing stats', async () => {
        await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/bridging/getDeposits/${testDydxAddress}`,
          expectedStatus: 200,
        });

        expect(stats.timing).toHaveBeenCalledWith(
          expect.stringContaining('bridging-controller.get_deposit_address.timing'),
          expect.any(Number),
        );
      });
    });
  });

  describe('POST /startBridge', () => {
    const arbNetwork = 'ARB_MAINNET';
    const arbChainId = alchemyNetworkToChainIdMap[arbNetwork];
    const evmFromAddress = '0x1111111111111111111111111111111111111111';
    const dydxAddress = 'dydx1234567890123456789012345678901234567899';

    beforeEach(async () => {
      // Create a user record needed by sweep/startEvmBridge
      await TurnkeyUsersTable.create({
        suborg_id: 'suborg-arb',
        evm_address: evmFromAddress,
        svm_address: 'svm1234567890123456789012345678901234567890',
        dydx_address: dydxAddress,
        salt: 'salt',
        created_at: new Date().toISOString(),
      } as unknown as TurnkeyUserCreateObject);

      // Mock external dependencies used during EVM bridge flow
      (skipClient.balances as unknown as jest.Mock).mockResolvedValue({
        chains: {
          [arbChainId]: {
            denoms: {
              // Provide large USD value to pass threshold check
              '0xaf88d065e77c8cC2239327C5EDb3A432268e5831': { amount: '100000000', valueUsd: '100.0' },
              // ETH denom can be present or omitted
            },
          },
        },
      });
      jest.spyOn(skipHelpers, 'getSkipCallData').mockResolvedValue([] as any);

      // Kernel account and client mocks
      (turnkeyViem.createAccount as jest.Mock).mockResolvedValue({} as any);
      (zeroDev.createKernelAccount as jest.Mock).mockResolvedValue({} as any);
      (zeroDev.createZeroDevPaymasterClient as jest.Mock).mockReturnValue({} as any);
      (zeroDev.getUserOperationGasPrice as jest.Mock).mockResolvedValue({
        maxFeePerGas: 1n,
        maxPriorityFeePerGas: 1n,
      } as any);

      const mockKernelClient = {
        account: {
          encodeCalls: jest.fn().mockResolvedValue('0xdeadbeef'),
        },
        sendUserOperation: jest.fn().mockResolvedValue('0xuserophash'),
        waitForUserOperationReceipt: jest.fn().mockResolvedValue({
          receipt: { transactionHash: '0xth' },
        }),
      };
      (zeroDev.createKernelAccountClient as jest.Mock).mockReturnValue(mockKernelClient as any);

      // Mock deserializePermissionAccount to return a mock account
      (zerodevPermissions.deserializePermissionAccount as jest.Mock).mockResolvedValue({
        address: '0xmockaccount',
        source: 'mock',
      } as any);
      // Mock toECDSASigner to return a mock signer
      (zerodevPermissionsSigners.toECDSASigner as jest.Mock).mockResolvedValue({
        sign: jest.fn().mockResolvedValue('0xsigned'),
        getAddress: jest.fn().mockReturnValue('0xsigner'),
      } as any);

      // Mock privateKeyToAccount to return a mock account
      (privateKeyToAccount as jest.Mock).mockReturnValue({
        address: '0xprivatekeyaccount',
        privateKey: new Uint8Array(32),
        publicKey: new Uint8Array(64),
      } as any);
    });

    it('should process EVM bridge request and send a user operation', async () => {
      await PermissionApprovalTable.create({
        suborg_id: 'suborg-arb',
        chain_id: arbChainId,
        approval: '0xapproval',
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/bridging/startBridge',
        expectedStatus: 200,
        body: {
          id: 'id',
          type: 'BLOCK_ACTIVITY',
          webhookId: 'wh_xxx',
          event: {
            network: arbNetwork,
            activity: [
              {
                fromAddress: '0x0000000000000000000000000000000000000000',
                toAddress: evmFromAddress,
                asset: 'USDC',
                value: '100',
              },
            ],
          },
        },
      });

      expect(response.status).toBe(200);
      // Called to build call data for bridging
      expect(skipHelpers.getSkipCallData).toHaveBeenCalled();
      // User operation should be sent once
      expect(zeroDev.createKernelAccountClient).toHaveBeenCalledTimes(1);
      const client = (
        zeroDev.createKernelAccountClient as unknown as jest.Mock
      ).mock.results[0].value;
      expect(client.sendUserOperation).toHaveBeenCalledTimes(1);
    });

    it('should return 500 for unsupported network', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/bridging/startBridge',
        expectedStatus: 500,
        body: {
          id: 'id',
          type: 'BLOCK_ACTIVITY',
          webhookId: 'wh_xxx',
          event: {
            network: 'UNSUPPORTED_NETWORK',
            activity: [
              {
                toAddress: evmFromAddress,
              },
            ],
          },
        },
      });

      expect(response.body.errors[0].msg).toBe('Internal Server Error');
    });
  });
});
