import {
  dbHelpers,
  PermissionApprovalFromDatabase,
  PermissionApprovalTable,
  TurnkeyUserFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { TurnkeyApiClient } from '@turnkey/sdk-server';
import * as turnkeyViem from '@turnkey/viem';
import { Address, hashTypedData } from 'viem';
import { generatePrivateKey, privateKeyToAccount } from 'viem/accounts';
import {
  arbitrum, avalanche, base, mainnet, optimism,
} from 'viem/chains';

import config from '../../src/config';
import { publicClients } from '../../src/helpers/alchemy-helpers';
import {
  getSolanaBridgePolicy,
  OBSOLETE_SOLANA_BRIDGE_POLICY_NAME,
  PolicyEngine,
  SOLANA_BRIDGE_POLICY_NAME,
} from '../../src/helpers/policy-engine';
import { approvalAllowsInjectiveBridge, callPolicyByChainId } from '../../src/lib/call-policies';
import { TurnkeyError } from '../../src/lib/errors';
import {
  DepositPolicyEvmApproval,
  DepositPolicyResponse,
  DepositPolicySignature,
  DepositPolicyStatus,
  DepositPolicyTypedData,
} from '../../src/types';
import {
  DYDX_ADDRESS,
  evaluateInstructionClauses,
  hashTypedDataV4,
  INJECTIVE_ADDRESS,
  SKIP_SOLANA_BRIDGE_INSTRUCTIONS,
} from './deposit-policy-helpers';
import {
  decodeApproval,
  getNobleRoutePolicy,
  serializeApprovalWithKey,
} from './zerodev-helpers';

jest.mock('@turnkey/viem', () => ({
  createAccount: jest.fn(),
}));

// Chains on which nothing has happened yet, so that no network is needed.
jest.mock('../../src/helpers/alchemy-helpers', () => {
  // eslint-disable-next-line @typescript-eslint/no-shadow, global-require
  const chains = require('viem/chains');
  // eslint-disable-next-line global-require
  const { createOfflineClient } = require('./zerodev-helpers');
  return {
    publicClients: Object.fromEntries(
      [chains.mainnet, chains.arbitrum, chains.base, chains.optimism, chains.avalanche].map(
        (chain) => [chain.id.toString(), createOfflineClient(chain)],
      ),
    ),
  };
});

const EIP_7702_CHAIN_IDS: string[] = [arbitrum, mainnet, optimism, base].map(
  (chain) => chain.id.toString(),
);
const AVALANCHE_CHAIN_ID: string = avalanche.id.toString();
const EVM_CHAIN_IDS: string[] = [...EIP_7702_CHAIN_IDS, AVALANCHE_CHAIN_ID];
const SMART_ACCOUNT_ADDRESS: Address = '0x9f9D4c574643a293c28753C76fD247Bc3A587CC0';
const API_USER_ID: string = 'api-user-id';
const SUBORG_ID: string = 'suborg-existing-user';

describe('policy-engine', () => {
  const owner = privateKeyToAccount(generatePrivateKey());
  const sessionKeyAddress = config.APPROVAL_SIGNER_PUBLIC_ADDRESS as Address;
  const user: TurnkeyUserFromDatabase = {
    suborg_id: SUBORG_ID,
    svm_address: 'vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg',
    evm_address: owner.address,
    smart_account_address: SMART_ACCOUNT_ADDRESS,
    salt: 'salt',
    dydx_address: DYDX_ADDRESS,
    created_at: '2026-01-01T00:00:00.000Z',
  };
  let turnkeyClient: {
    getUsers: jest.Mock,
    getPolicies: jest.Mock,
    createPolicy: jest.Mock,
    updateRootQuorum: jest.Mock,
  };
  let engine: PolicyEngine;

  /** The approval the SDK makes for `chainId` when it holds the user's key. */
  const approvalMadeWithKey = async (
    chainId: string,
    route: 'injective' | 'noble',
  ): Promise<string> => serializeApprovalWithKey({
    client: publicClients[chainId],
    owner,
    sessionKeyAddress,
    policy: route === 'injective'
      ? await callPolicyByChainId[chainId](DYDX_ADDRESS)
      : getNobleRoutePolicy(),
    accountAddress: SMART_ACCOUNT_ADDRESS,
  });

  /** Gives the account the approvals it was set up with, for the earlier route. */
  const seedNobleRouteApprovals = async (chainIds: string[] = EVM_CHAIN_IDS): Promise<void> => {
    for (const chainId of chainIds) {
      await PermissionApprovalTable.create({
        suborg_id: SUBORG_ID,
        chain_id: chainId,
        approval: await approvalMadeWithKey(chainId, 'noble'),
      });
    }
  };

  const getStoredApproval = (chainId: string): Promise<string | undefined> => {
    return PermissionApprovalTable.getApprovalForSuborgAndChain(SUBORG_ID, chainId);
  };

  const getApproval = (
    policy: DepositPolicyResponse,
    chainId: string,
  ): DepositPolicyEvmApproval => {
    return policy.evm.approvals.find((a: DepositPolicyEvmApproval) => a.chainId === chainId)!;
  };

  /** Signs the typed data of each chain that needs it, as the frontend does. */
  const signAll = async (
    policy: DepositPolicyResponse,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    signer: { signTypedData: (typedData: any) => Promise<string> } = owner,
  ): Promise<DepositPolicySignature[]> => Promise.all(
    policy.evm.approvals
      .filter((a: DepositPolicyEvmApproval) => a.status === DepositPolicyStatus.ACTION_REQUIRED)
      .map(async (a: DepositPolicyEvmApproval) => ({
        chainId: a.chainId,
        signature: await signer.signTypedData(a.typedData),
      })),
  );

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(() => {
    turnkeyClient = {
      getUsers: jest.fn().mockResolvedValue({
        users: [
          { userName: 'End User', userId: 'end-user-id' },
          { userName: 'API User', userId: API_USER_ID },
        ],
      }),
      getPolicies: jest.fn().mockResolvedValue({ policies: [] }),
      createPolicy: jest.fn().mockResolvedValue({}),
      updateRootQuorum: jest.fn().mockResolvedValue({}),
    };
    engine = new PolicyEngine(turnkeyClient as unknown as TurnkeyApiClient);
    (turnkeyViem.createAccount as jest.Mock).mockResolvedValue(owner);
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('getSolanaBridgePolicy', () => {
    const policy = getSolanaBridgePolicy(DYDX_ADDRESS, API_USER_ID);
    const [burn, feeTransfer] = SKIP_SOLANA_BRIDGE_INSTRUCTIONS;

    it('allows the bridge API user to sign a CCTP burn to the user\'s Injective account', () => {
      expect(policy).toEqual({
        policyName: 'Injective Bridging Policy',
        effect: 'EFFECT_ALLOW',
        condition: [
          'solana.tx.instructions.count() == 2',
          "solana.tx.instructions[0].program_key == 'CCTPV2vPZJS2u2BBsUoscuikbYjnpFmbFsvVuJdgUMQe'",
          "solana.tx.instructions[0].instruction_data_hex[0..16] == 'd73c3d2e723780b0'",
          "solana.tx.instructions[0].instruction_data_hex[32..40] == '1d000000'",
          'solana.tx.instructions[0].instruction_data_hex[40..104] == ' +
            "'0000000000000000000000000102030405060708090a0b0c0d0e0f1011121314'",
          "solana.tx.instructions[1].program_key == 'TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA'",
          'solana.tx.spl_transfers.count() == 1',
          'solana.tx.spl_transfers[0].amount <= 1000000',
        ].join(' && '),
        consensus: "approvers.any(user, user.id == 'api-user-id')",
        notes: expect.any(String),
      });
    });

    it('matches the instructions of the transaction Skip returns for the route', () => {
      expect(
        evaluateInstructionClauses(policy.condition, SKIP_SOLANA_BRIDGE_INSTRUCTIONS),
      ).toEqual({
        result: true,
        // The two clauses about token transfers, which Turnkey derives from the instructions.
        unevaluated: [
          'solana.tx.spl_transfers.count() == 1',
          'solana.tx.spl_transfers[0].amount <= 1000000',
        ],
      });
      // The one token transfer is the relayer's fee, which is under the cap.
      const fee: bigint = Buffer.from(feeTransfer.dataHex, 'hex').readBigUInt64LE(1);
      expect(fee).toEqual(BigInt(378));
      expect(fee).toBeLessThanOrEqual(BigInt(1_000_000));
    });

    it.each([
      [
        'a burn to another recipient',
        [{ ...burn, dataHex: burn.dataHex.replace('01020304', '11111111') }, feeTransfer],
      ],
      [
        'a burn for another chain',
        [{ ...burn, dataHex: burn.dataHex.replace('1d000000', '04000000') }, feeTransfer],
      ],
      [
        'the CCTP v1 instruction',
        [{ ...burn, dataHex: burn.dataHex.replace('d73c3d2e723780b0', 'a7de137255150e76') }, feeTransfer],
      ],
      [
        'the same data sent to another program',
        [{ ...burn, programKey: 'CCTPiPYPc6AsJuwueEnWgSgucamXDZwBd53dQ11YiKX3' }, feeTransfer],
      ],
      ['a third instruction', [burn, feeTransfer, feeTransfer]],
      ['the burn alone', [burn]],
      ['a second instruction of another program', [burn, { ...feeTransfer, programKey: '11111111111111111111111111111111' }]],
      ['the instructions in the other order', [feeTransfer, burn]],
    ])('does not match %s', (
      _name: string,
      instructions: { programKey: string, dataHex: string }[],
    ) => {
      expect(evaluateInstructionClauses(policy.condition, instructions).result).toEqual(false);
    });
  });

  describe('configurePolicy, for a new account', () => {
    it('approves the Injective route on the EIP-7702 chains', async () => {
      await engine.configurePolicy(SUBORG_ID, owner.address, DYDX_ADDRESS);

      for (const chainId of EIP_7702_CHAIN_IDS) {
        const approval: string | undefined = await getStoredApproval(chainId);
        expect(approval).toEqual(await approvalMadeWithKey(chainId, 'injective'));
        expect(approvalAllowsInjectiveBridge(approval!, DYDX_ADDRESS)).toEqual(true);
      }
    });

    it('creates the Solana policy in the user\'s sub-organization', async () => {
      await engine.configurePolicy(SUBORG_ID, owner.address, DYDX_ADDRESS);

      expect(turnkeyClient.createPolicy).toHaveBeenCalledTimes(1);
      expect(turnkeyClient.createPolicy).toHaveBeenCalledWith({
        organizationId: SUBORG_ID,
        ...getSolanaBridgePolicy(DYDX_ADDRESS, API_USER_ID),
      });
    });
  });

  describe('getDepositPolicy', () => {
    it('names the account deposits arrive at and the address that signs', async () => {
      await seedNobleRouteApprovals();

      const policy: DepositPolicyResponse = await engine.getDepositPolicy(user);

      expect(policy.dydxAddress).toEqual(DYDX_ADDRESS);
      expect(policy.injectiveAddress).toEqual(INJECTIVE_ADDRESS);
      expect(policy.evm.address).toEqual(owner.address);
    });

    it('asks an account approved for the earlier route to sign on every chain', async () => {
      await seedNobleRouteApprovals();

      const policy: DepositPolicyResponse = await engine.getDepositPolicy(user);

      expect(policy.evm.approvals.map((a) => [a.chainId, a.status])).toEqual(
        EVM_CHAIN_IDS.map((chainId: string) => [chainId, DepositPolicyStatus.ACTION_REQUIRED]),
      );
      for (const chainId of EVM_CHAIN_IDS) {
        const typedData: DepositPolicyTypedData = getApproval(policy, chainId).typedData!;
        const isAvalanche: boolean = chainId === AVALANCHE_CHAIN_ID;
        expect(typedData.domain).toEqual({
          name: 'Kernel',
          version: isAvalanche ? '0.3.1' : '0.3.3',
          chainId: Number(chainId),
          // The user's own address, except on Avalanche where they have a smart account.
          verifyingContract: isAvalanche ? SMART_ACCOUNT_ADDRESS : owner.address,
        });
        expect(typedData.primaryType).toEqual('Enable');
        // A complete eth_signTypedData_v4 payload: the type of the domain is spelled out.
        expect(Object.keys(typedData.types)).toEqual(['EIP712Domain', 'Enable']);
        expect(typedData.types.EIP712Domain).toEqual([
          { name: 'name', type: 'string' },
          { name: 'version', type: 'string' },
          { name: 'chainId', type: 'uint256' },
          { name: 'verifyingContract', type: 'address' },
        ]);
        expect(Object.keys(typedData.domain)).toEqual(
          typedData.types.EIP712Domain.map((field) => field.name),
        );
        expect(Object.keys(typedData.message)).toEqual(
          typedData.types.Enable.map((field) => field.name),
        );
      }
    });

    it('gives typed data that survives being sent as JSON', async () => {
      await seedNobleRouteApprovals();

      const policy: DepositPolicyResponse = await engine.getDepositPolicy(user);

      expect(JSON.parse(JSON.stringify(policy))).toEqual(policy);
    });

    it('cannot renew an EIP-7702 chain that has no approval on record', async () => {
      await seedNobleRouteApprovals([arbitrum.id.toString()]);

      const policy: DepositPolicyResponse = await engine.getDepositPolicy(user);

      expect(getApproval(policy, arbitrum.id.toString()).status).toEqual(
        DepositPolicyStatus.ACTION_REQUIRED,
      );
      expect(getApproval(policy, base.id.toString())).toEqual({
        chainId: base.id.toString(),
        status: DepositPolicyStatus.UNAVAILABLE,
      });
      // Avalanche needs nothing from an earlier approval.
      expect(getApproval(policy, AVALANCHE_CHAIN_ID).status).toEqual(
        DepositPolicyStatus.ACTION_REQUIRED,
      );
    });

    it('asks nothing of a chain that is already approved for the Injective route', async () => {
      await seedNobleRouteApprovals();
      await PermissionApprovalTable.upsert({
        suborg_id: SUBORG_ID,
        chain_id: optimism.id.toString(),
        approval: await approvalMadeWithKey(optimism.id.toString(), 'injective'),
      });

      const policy: DepositPolicyResponse = await engine.getDepositPolicy(user);

      expect(getApproval(policy, optimism.id.toString())).toEqual({
        chainId: optimism.id.toString(),
        status: DepositPolicyStatus.UP_TO_DATE,
      });
      expect(getApproval(policy, mainnet.id.toString()).status).toEqual(
        DepositPolicyStatus.ACTION_REQUIRED,
      );
    });

    describe('solana', () => {
      it('asks for the policy to be created, and names the earlier one as obsolete', async () => {
        turnkeyClient.getPolicies.mockResolvedValue({
          policies: [
            { policyId: 'policy-1', policyName: OBSOLETE_SOLANA_BRIDGE_POLICY_NAME },
            { policyId: 'policy-2', policyName: 'Something else' },
          ],
        });

        const { solana } = await engine.getDepositPolicy(user);

        expect(solana).toEqual({
          status: DepositPolicyStatus.ACTION_REQUIRED,
          policy: getSolanaBridgePolicy(DYDX_ADDRESS, API_USER_ID),
          obsoletePolicyIds: ['policy-1'],
        });
        expect(turnkeyClient.getPolicies).toHaveBeenCalledWith({ organizationId: SUBORG_ID });
        // Reading what to do changes nothing.
        expect(turnkeyClient.createPolicy).not.toHaveBeenCalled();
      });

      it('asks nothing once the policy exists', async () => {
        turnkeyClient.getPolicies.mockResolvedValue({
          policies: [
            { policyId: 'policy-1', policyName: OBSOLETE_SOLANA_BRIDGE_POLICY_NAME },
            { policyId: 'policy-3', policyName: SOLANA_BRIDGE_POLICY_NAME },
          ],
        });

        const { solana } = await engine.getDepositPolicy(user);

        expect(solana).toEqual({
          status: DepositPolicyStatus.UP_TO_DATE,
          obsoletePolicyIds: ['policy-1'],
        });
      });

      it('is unavailable when Turnkey cannot be read, without affecting the EVM chains', async () => {
        await seedNobleRouteApprovals();
        turnkeyClient.getPolicies.mockRejectedValue(new Error('turnkey is down'));

        const policy: DepositPolicyResponse = await engine.getDepositPolicy(user);

        expect(policy.solana).toEqual({
          status: DepositPolicyStatus.UNAVAILABLE,
          obsoletePolicyIds: [],
        });
        expect(getApproval(policy, arbitrum.id.toString()).status).toEqual(
          DepositPolicyStatus.ACTION_REQUIRED,
        );
      });
    });
  });

  describe('submitDepositApprovals', () => {
    it('stores, for each chain, the approval the SDK makes when it holds the key', async () => {
      await seedNobleRouteApprovals();
      const signatures: DepositPolicySignature[] = await signAll(
        await engine.getDepositPolicy(user),
      );

      const approved: string[] = await engine.submitDepositApprovals(user, signatures);

      expect(approved).toEqual(EVM_CHAIN_IDS);
      for (const chainId of EVM_CHAIN_IDS) {
        expect(await getStoredApproval(chainId)).toEqual(
          await approvalMadeWithKey(chainId, 'injective'),
        );
      }
    });

    it('accepts signatures made from the payload alone, as eth_signTypedData_v4 defines', async () => {
      await seedNobleRouteApprovals();
      const policy: DepositPolicyResponse = await engine.getDepositPolicy(user);
      // A signer that derives nothing: it hashes the domain with the type the payload gives.
      const signatures: DepositPolicySignature[] = await signAll(policy, {
        signTypedData: (typedData: DepositPolicyTypedData) => {
          return owner.sign({ hash: hashTypedDataV4(typedData) });
        },
      });

      expect(await engine.submitDepositApprovals(user, signatures)).toEqual(EVM_CHAIN_IDS);

      for (const chainId of EVM_CHAIN_IDS) {
        expect(await getStoredApproval(chainId)).toEqual(
          await approvalMadeWithKey(chainId, 'injective'),
        );
      }
    });

    it('gives a payload whose hash does not depend on how the signer finds the domain type', async () => {
      await seedNobleRouteApprovals();
      const policy: DepositPolicyResponse = await engine.getDepositPolicy(user);

      for (const { typedData } of policy.evm.approvals) {
        // With the domain type taken from the payload, and with it derived from the domain.
        const otherTypes = { Enable: typedData!.types.Enable };
        const derived = { ...typedData!, types: otherTypes } as Parameters<typeof hashTypedData>[0];
        expect(hashTypedDataV4(typedData!)).toEqual(hashTypedData(derived));
        // The hash a signer would make if it took a missing domain type to be empty is another
        // one, which is why the payload spells the type out.
        expect(hashTypedDataV4({ ...typedData!, types: { ...otherTypes, EIP712Domain: [] } }))
          .not.toEqual(hashTypedDataV4(typedData!));
      }
    });

    it('keeps the authorization that makes the user\'s address a smart account', async () => {
      await seedNobleRouteApprovals();
      const chainId: string = arbitrum.id.toString();
      const before = decodeApproval((await getStoredApproval(chainId))!);

      await engine.submitDepositApprovals(user, await signAll(await engine.getDepositPolicy(user)));

      const after = decodeApproval((await getStoredApproval(chainId))!);
      expect(before.eip7702Auth).toEqual(
        expect.objectContaining({ chainId: arbitrum.id, nonce: 0 }),
      );
      expect(after.eip7702Auth).toEqual(before.eip7702Auth);
      expect(after.accountParams.accountAddress).toEqual(owner.address);
    });

    it('leaves nothing more to approve', async () => {
      await seedNobleRouteApprovals();
      await engine.submitDepositApprovals(user, await signAll(await engine.getDepositPolicy(user)));

      const policy: DepositPolicyResponse = await engine.getDepositPolicy(user);

      expect(policy.evm.approvals).toEqual(EVM_CHAIN_IDS.map((chainId: string) => ({
        chainId,
        status: DepositPolicyStatus.UP_TO_DATE,
      })));
      for (const chainId of EVM_CHAIN_IDS) {
        expect(
          approvalAllowsInjectiveBridge((await getStoredApproval(chainId))!, DYDX_ADDRESS),
        ).toEqual(true);
      }
    });

    it('stores the approval of one chain, leaving the others as they were', async () => {
      await seedNobleRouteApprovals();
      const chainId: string = base.id.toString();
      const [signature] = (await signAll(await engine.getDepositPolicy(user))).filter(
        (s: DepositPolicySignature) => s.chainId === chainId,
      );

      expect(await engine.submitDepositApprovals(user, [signature])).toEqual([chainId]);

      expect(await getStoredApproval(chainId)).toEqual(
        await approvalMadeWithKey(chainId, 'injective'),
      );
      expect(await getStoredApproval(arbitrum.id.toString())).toEqual(
        await approvalMadeWithKey(arbitrum.id.toString(), 'noble'),
      );
    });

    describe('refuses, storing nothing', () => {
      let before: PermissionApprovalFromDatabase[];
      let signatures: DepositPolicySignature[];

      const expectRefusal = async (
        submitted: DepositPolicySignature[],
        message: string,
      ): Promise<void> => {
        const error: Error = await engine.submitDepositApprovals(user, submitted).then(
          () => { throw new Error('expected a refusal, but the approvals were stored'); },
          (e: Error) => e,
        );
        expect(error).toBeInstanceOf(TurnkeyError);
        expect(error.message).toContain(message);
        expect(await PermissionApprovalTable.findBySuborgId(SUBORG_ID)).toEqual(before);
      };

      beforeEach(async () => {
        await seedNobleRouteApprovals();
        before = await PermissionApprovalTable.findBySuborgId(SUBORG_ID);
        signatures = await signAll(await engine.getDepositPolicy(user));
      });

      it('a signature made with another key', async () => {
        const stranger = privateKeyToAccount(generatePrivateKey());
        const forged: DepositPolicySignature[] = await signAll(
          await engine.getDepositPolicy(user),
          stranger,
        );

        await expectRefusal(forged, 'does not match the deposit policy');
      });

      it('the user\'s signature for another chain', async () => {
        const [first, second] = signatures;

        await expectRefusal(
          [{ chainId: first.chainId, signature: second.signature }],
          `Signature for chain ${first.chainId} does not match the deposit policy`,
        );
      });

      it('the user\'s signature of the policy of another dYdX address', async () => {
        const other: DepositPolicySignature[] = await signAll(await engine.getDepositPolicy({
          ...user,
          dydx_address: 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4',
        }));

        await expectRefusal(other, 'does not match the deposit policy');
      });

      it('every signature when one of them is wrong', async () => {
        const stranger = privateKeyToAccount(generatePrivateKey());
        const last: DepositPolicySignature = signatures[signatures.length - 1];
        const mixed: DepositPolicySignature[] = [
          ...signatures.slice(0, -1),
          {
            chainId: last.chainId,
            signature: await stranger.signMessage({ message: 'not the typed data' }),
          },
        ];

        await expectRefusal(mixed, `Signature for chain ${last.chainId} does not match`);
      });

      it('a signature that is not a signature', async () => {
        await expectRefusal(
          [{ chainId: signatures[0].chainId, signature: '0x1234' }],
          'does not match the deposit policy',
        );
      });

      it('a chain given twice', async () => {
        await expectRefusal([signatures[0], signatures[0]], 'Each chain may appear only once');
      });

      it('a chain that deposits are not bridged from', async () => {
        await expectRefusal(
          [{ chainId: '1776', signature: signatures[0].signature }],
          'Unsupported chain: 1776',
        );
      });
    });
  });
});
