import {
  dbHelpers,
  TurnkeyUserCreateObject,
  TurnkeyUsersTable,
} from '@dydxprotocol-indexer/postgres';
import { stats } from '@dydxprotocol-indexer/base';
import request from 'supertest';
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
});
