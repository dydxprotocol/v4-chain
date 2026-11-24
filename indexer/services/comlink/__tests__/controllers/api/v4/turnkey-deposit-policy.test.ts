import { dbHelpers, TurnkeyUsersTable } from '@dydxprotocol-indexer/postgres';
import request from 'supertest';

import { PolicyEngine } from '../../../../src/helpers/policy-engine';
import { TurnkeyError } from '../../../../src/lib/errors';
import {
  DepositPolicyResponse,
  DepositPolicyStatus,
  RequestMethod,
} from '../../../../src/types';
import {
  DYDX_ADDRESS,
  INJECTIVE_ADDRESS,
  INJECTIVE_EVM_ADDRESS,
} from '../../../helpers/deposit-policy-helpers';
import { sendRequest } from '../../../helpers/helpers';

const PATH: string = '/v4/turnkey/depositPolicy';
const UNKNOWN_DYDX_ADDRESS: string = 'dydx1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5kmz6xt';
const EVM_ADDRESS: string = '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266';
const SIGNATURE: string = `0x${'ab'.repeat(65)}`;
const POLICY: DepositPolicyResponse = {
  dydxAddress: DYDX_ADDRESS,
  injectiveAddress: INJECTIVE_ADDRESS,
  evm: {
    address: EVM_ADDRESS,
    approvals: [
      {
        chainId: '42161',
        status: DepositPolicyStatus.ACTION_REQUIRED,
        typedData: {
          domain: {
            name: 'Kernel', version: '0.3.3', chainId: 42161, verifyingContract: EVM_ADDRESS,
          },
          types: {
            EIP712Domain: [{ name: 'name', type: 'string' }],
            Enable: [{ name: 'validationId', type: 'bytes21' }],
          },
          primaryType: 'Enable',
          message: {
            validationId: '0x02',
            nonce: 1,
            hook: '0x0000000000000000000000000000000000000000',
            validatorData: '0x',
            hookData: '0x',
            selectorData: '0x',
          },
        },
      },
      { chainId: '43114', status: DepositPolicyStatus.UP_TO_DATE },
    ],
  },
  solana: {
    status: DepositPolicyStatus.ACTION_REQUIRED,
    policy: {
      policyName: 'Injective Bridging Policy',
      effect: 'EFFECT_ALLOW',
      condition: 'solana.tx.instructions.count() == 2',
      consensus: "approvers.any(user, user.id == 'api-user-id')",
      notes: 'notes',
    },
    obsoletePolicyIds: ['policy-1'],
  },
};

describe('turnkey-controller deposit policy#V4', () => {
  let getDepositPolicy: jest.SpyInstance;
  let submitDepositApprovals: jest.SpyInstance;

  const post = (body: object, expectedStatus: number): Promise<request.Response> => sendRequest({
    type: RequestMethod.POST,
    path: PATH,
    body,
    expectedStatus,
  });

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    getDepositPolicy = jest.spyOn(PolicyEngine.prototype, 'getDepositPolicy')
      .mockResolvedValue(POLICY);
    submitDepositApprovals = jest.spyOn(PolicyEngine.prototype, 'submitDepositApprovals')
      .mockResolvedValue(['42161']);
    await TurnkeyUsersTable.create({
      suborg_id: 'suborg-1',
      evm_address: EVM_ADDRESS,
      svm_address: 'vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg',
      smart_account_address: '0x9f9D4c574643a293c28753C76fD247Bc3A587CC0',
      salt: 'salt',
      dydx_address: DYDX_ADDRESS,
      created_at: '2026-01-01T00:00:00.000Z',
    });
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.restoreAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET /depositPolicy/:dydxAddress', () => {
    const QUERY: string = `?injectiveAddress=${INJECTIVE_ADDRESS}`;

    it('returns what the user has to approve for the Injective account they name', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `${PATH}/${DYDX_ADDRESS}${QUERY}`,
        expectedStatus: 200,
      });

      expect(response.body).toEqual(POLICY);
      expect(getDepositPolicy).toHaveBeenCalledTimes(1);
      expect(getDepositPolicy).toHaveBeenCalledWith(
        expect.objectContaining({
          suborg_id: 'suborg-1',
          dydx_address: DYDX_ADDRESS,
          evm_address: EVM_ADDRESS,
        }),
        INJECTIVE_ADDRESS,
      );
    });

    it('does not reveal the user\'s sub-organization', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `${PATH}/${DYDX_ADDRESS}${QUERY}`,
        expectedStatus: 200,
      });

      expect(JSON.stringify(response.body)).not.toContain('suborg-1');
    });

    it('answers 404 for a dYdX address with no account', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `${PATH}/${UNKNOWN_DYDX_ADDRESS}${QUERY}`,
        expectedStatus: 404,
      });

      expect(response.body.errors[0].msg).toEqual(
        `No user found with dydx address: ${UNKNOWN_DYDX_ADDRESS}`,
      );
      expect(getDepositPolicy).not.toHaveBeenCalled();
    });

    it.each([
      ['an Injective address', INJECTIVE_ADDRESS],
      ['an EVM address', EVM_ADDRESS],
      ['text', 'not-an-address'],
    ])('answers 400 for %s', async (_name: string, address: string) => {
      await sendRequest({
        type: RequestMethod.GET,
        path: `${PATH}/${address}${QUERY}`,
        expectedStatus: 400,
      });

      expect(getDepositPolicy).not.toHaveBeenCalled();
    });

    it.each([
      ['no Injective address', ''],
      ['an empty Injective address', '?injectiveAddress='],
      ['a dYdX address', `?injectiveAddress=${DYDX_ADDRESS}`],
      ['the EVM form of the Injective address', `?injectiveAddress=${INJECTIVE_EVM_ADDRESS}`],
      [
        'an Injective address with a bad checksum',
        `?injectiveAddress=${INJECTIVE_ADDRESS.slice(0, -1)}q`,
      ],
      [
        'two Injective addresses',
        `?injectiveAddress=${INJECTIVE_ADDRESS}&injectiveAddress=${INJECTIVE_ADDRESS}`,
      ],
    ])('answers 400 for %s as where deposits go', async (_name: string, query: string) => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `${PATH}/${DYDX_ADDRESS}${query}`,
        expectedStatus: 400,
      });

      expect(response.body.errors[0].msg).toEqual(
        'injectiveAddress must be a valid injective address',
      );
      expect(getDepositPolicy).not.toHaveBeenCalled();
    });
  });

  describe('POST /depositPolicy', () => {
    const signatures = [{ chainId: '42161', signature: SIGNATURE }];

    const body = { dydxAddress: DYDX_ADDRESS, injectiveAddress: INJECTIVE_ADDRESS, signatures };

    it('stores the user\'s signatures and names the chains approved', async () => {
      const response: request.Response = await post(body, 200);

      expect(response.body).toEqual({ approvedChainIds: ['42161'] });
      expect(submitDepositApprovals).toHaveBeenCalledTimes(1);
      expect(submitDepositApprovals).toHaveBeenCalledWith(
        expect.objectContaining({ suborg_id: 'suborg-1', dydx_address: DYDX_ADDRESS }),
        INJECTIVE_ADDRESS,
        signatures,
      );
    });

    it('passes on only the chain and signature of each entry', async () => {
      await post({
        ...body,
        signatures: [{ chainId: '42161', signature: SIGNATURE, approval: 'forged' }],
        suborg_id: 'another-suborg',
      }, 200);

      expect(submitDepositApprovals).toHaveBeenCalledWith(
        expect.objectContaining({ suborg_id: 'suborg-1' }),
        INJECTIVE_ADDRESS,
        signatures,
      );
    });

    it('answers 400 with the reason when a signature is refused', async () => {
      submitDepositApprovals.mockRejectedValue(
        new TurnkeyError('Signature for chain 42161 does not match the deposit policy'),
      );

      const response: request.Response = await post(body, 400);

      expect(response.body.errors[0].msg).toEqual(
        'Signature for chain 42161 does not match the deposit policy',
      );
    });

    it('answers 500, revealing nothing, when storing fails unexpectedly', async () => {
      submitDepositApprovals.mockRejectedValue(new Error('connection to 10.0.0.1 refused'));

      const response: request.Response = await post(body, 500);

      expect(response.body).toEqual({ errors: [{ msg: 'Internal Server Error' }] });
    });

    it('answers 404 for a dYdX address with no account', async () => {
      await post({ ...body, dydxAddress: UNKNOWN_DYDX_ADDRESS }, 404);

      expect(submitDepositApprovals).not.toHaveBeenCalled();
    });

    it.each([
      ['no dYdX address', { ...body, dydxAddress: undefined }],
      ['a dYdX address that is not one', { ...body, dydxAddress: EVM_ADDRESS }],
      ['no Injective address', { ...body, injectiveAddress: undefined }],
      ['a dYdX address as the Injective address', { ...body, injectiveAddress: DYDX_ADDRESS }],
      [
        'the EVM form of the Injective address',
        { ...body, injectiveAddress: INJECTIVE_EVM_ADDRESS },
      ],
      ['an Injective address that is not text', { ...body, injectiveAddress: [INJECTIVE_ADDRESS] }],
      ['no signatures', { ...body, signatures: undefined }],
      ['an empty list of signatures', { ...body, signatures: [] }],
      ['signatures that are not a list', { ...body, signatures: SIGNATURE }],
      [
        'more signatures than there are chains',
        {
          ...body,
          signatures: ['1', '10', '8453', '42161', '43114', '137'].map(
            (chainId: string) => ({ chainId, signature: SIGNATURE }),
          ),
        },
      ],
      [
        'a signature of the wrong length',
        { ...body, signatures: [{ chainId: '42161', signature: '0x1234' }] },
      ],
      [
        'a signature that is not hex',
        { ...body, signatures: [{ chainId: '42161', signature: `0x${'zz'.repeat(65)}` }] },
      ],
      ['an entry with no signature', { ...body, signatures: [{ chainId: '42161' }] }],
      [
        'a chain id that is not a number',
        { ...body, signatures: [{ chainId: 'arbitrum', signature: SIGNATURE }] },
      ],
      ['an entry with no chain id', { ...body, signatures: [{ signature: SIGNATURE }] }],
    ])('answers 400 for %s', async (_name: string, invalid: object) => {
      await post(invalid, 400);

      expect(submitDepositApprovals).not.toHaveBeenCalled();
    });
  });
});
