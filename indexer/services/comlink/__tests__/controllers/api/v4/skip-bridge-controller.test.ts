import {
  dbHelpers,
  PermissionApprovalTable,
  TurnkeyUserCreateObject,
  TurnkeyUsersTable,
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
import * as skipClient from '@skip-go/client/cjs';
import * as zeroDev from '@zerodev/sdk';
import * as turnkeyViem from '@turnkey/viem';
import * as zerodevPermissions from '@zerodev/permissions';

import { alchemyNetworkToChainIdMap } from '../../../../src/helpers/alchemy-helpers';
import * as skipHelpers from '../../../../src/helpers/skip-helper';
import { RequestMethod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';

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
              '0xaf88d065e77c8cC2239327C5EDb3A432268e5831': { amount: '1000000', valueUsd: '100.0' },
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
