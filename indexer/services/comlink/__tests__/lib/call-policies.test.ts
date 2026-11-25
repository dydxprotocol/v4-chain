import { toCallPolicy } from '@zerodev/permissions/policies';
import { decodeFunctionData } from 'viem';
import { generatePrivateKey, privateKeyToAccount } from 'viem/accounts';
import {
  arbitrum, avalanche, base, mainnet, optimism,
} from 'viem/chains';

import { cctpV2RelayerAbi } from '../../src/lib/bridge-abi';
import {
  approvalAllowsInjectiveBridge,
  callPolicyByChainId,
  getNoCallsPolicy,
} from '../../src/lib/call-policies';
import { usdcAddressByChainId } from '../../src/lib/smart-contract-constants';
import {
  DYDX_ADDRESS,
  INJECTIVE_ADDRESS,
  INJECTIVE_ADDRESS_WITH_DYDX_BYTES,
  INJECTIVE_MINT_RECIPIENT,
  OTHER_INJECTIVE_ADDRESS,
  OTHER_INJECTIVE_MINT_RECIPIENT,
  SKIP_ARBITRUM_BRIDGE_TX,
} from '../helpers/deposit-policy-helpers';
import {
  createOfflineClient,
  getNobleRoutePolicy,
  serializeApprovalWithKey,
} from '../helpers/zerodev-helpers';

const RELAYER: string = '0x400BB58033a7763A834199190B68F66A2661aE73';
const SESSION_KEY_ADDRESS = '0x3FC11ff27e5373c88EA142d2EdF5492d0839980B';
// The condition codes of the call policy contract that the policy uses.
const EQUAL: number = 0;
const LESS_THAN_OR_EQUAL: number = 4;

interface Rule { offset: number, condition: number, params: string[] }
interface Permission { target: string, selector: string, valueLimit: bigint, rules: Rule[] }

async function getPermissions(chainId: string, injectiveAddress: string): Promise<Permission[]> {
  const policy = toCallPolicy(await callPolicyByChainId[chainId](injectiveAddress));
  return (policy.policyParams as unknown as { permissions: Permission[] }).permissions;
}

/**
 * Whether a call passes a permission, the way the call policy contract checks it: the target and
 * selector match, and each rule's 32-byte word, at its offset into the arguments, equals the
 * rule's value or, for a cap, is at most the rule's value as an unsigned number.
 */
function isPermitted(permission: Permission, to: string, data: string): boolean {
  const args: string = data.slice(8);
  return permission.target.toLowerCase() === to.toLowerCase() &&
    permission.selector.toLowerCase() === `0x${data.slice(0, 8)}` &&
    permission.rules.every((rule: Rule) => {
      const word: string = `0x${args.slice(rule.offset * 2, rule.offset * 2 + 64)}`;
      if (rule.condition === LESS_THAN_OR_EQUAL) {
        expect(rule.params).toHaveLength(1);
        return BigInt(word) <= BigInt(rule.params[0]);
      }
      expect(rule.condition).toEqual(EQUAL);
      return rule.params.some((param: string) => param.toLowerCase() === word);
    });
}

/** Replaces the 32-byte argument at `index` of the bridge call. */
function withArgument(index: number, word: string): string {
  const { data } = SKIP_ARBITRUM_BRIDGE_TX;
  const start: number = 8 + index * 64;
  return `${data.slice(0, start)}${word.padStart(64, '0')}${data.slice(start + 64)}`;
}

describe('call-policies', () => {
  const arbitrumId: string = arbitrum.id.toString();

  describe('callPolicyByChainId', () => {
    it.each([
      mainnet, arbitrum, base, optimism, avalanche,
    ].map((chain) => [chain.name, chain.id.toString()]))(
      'on %s, permits approving the relayer for USDC and bridging to Injective, and no more',
      async (_name: string, chainId: string) => {
        const permissions: Permission[] = await getPermissions(chainId, INJECTIVE_ADDRESS);

        expect(permissions).toHaveLength(2);
        const [approve, bridge] = permissions;
        expect(approve.target).toEqual(usdcAddressByChainId[chainId]);
        // approve(address,uint256)
        expect(approve.selector).toEqual('0x095ea7b3');
        expect(bridge.target).toEqual(RELAYER);
        expect(bridge.selector).toEqual('0xd1cc447d');
        // Neither call may send the native token.
        expect(approve.valueLimit).toEqual(BigInt(0));
        expect(bridge.valueLimit).toEqual(BigInt(0));
      },
    );

    it('reads the transaction Skip returns for the route with its ABI', () => {
      const { functionName, args } = decodeFunctionData({
        abi: cctpV2RelayerAbi,
        data: `0x${SKIP_ARBITRUM_BRIDGE_TX.data}`,
      });

      expect(functionName).toEqual('requestCCTPTransferWithCaller');
      expect(args[0]).toEqual(BigInt(99_999_622));
      expect(args[1]).toEqual(29);
      expect(args[2]).toEqual(INJECTIVE_MINT_RECIPIENT);
      expect(args[3]).toEqual(usdcAddressByChainId[arbitrumId]);
      expect(args[4]).toEqual(BigInt(378));
      expect(args[6]).toEqual(2000);
      expect(args[8]).toEqual('0x');
    });

    it('permits the transaction Skip returns for the route', async () => {
      const [, bridge] = await getPermissions(arbitrumId, INJECTIVE_ADDRESS);

      expect(
        isPermitted(bridge, SKIP_ARBITRUM_BRIDGE_TX.to, SKIP_ARBITRUM_BRIDGE_TX.data),
      ).toEqual(true);
    });

    it('constrains the chain, recipient, token and relayer fee of the bridge call', async () => {
      const [, bridge] = await getPermissions(arbitrumId, INJECTIVE_ADDRESS);

      // Offsets are in bytes into the arguments: the second to fifth of them.
      expect(bridge.rules).toEqual([
        { offset: 32, condition: EQUAL, params: [`0x${'1d'.padStart(64, '0')}`] },
        { offset: 64, condition: EQUAL, params: [INJECTIVE_MINT_RECIPIENT] },
        {
          offset: 96,
          condition: EQUAL,
          params: [`0x${usdcAddressByChainId[arbitrumId].slice(2).padStart(64, '0')}`],
        },
        // 5 USDC.
        { offset: 128, condition: LESS_THAN_OR_EQUAL, params: [`0x${'4c4b40'.padStart(64, '0')}`] },
      ]);
    });

    it('permits any amount', async () => {
      const [, bridge] = await getPermissions(arbitrumId, INJECTIVE_ADDRESS);

      expect(isPermitted(bridge, RELAYER, withArgument(0, 'ffffffffffff'))).toEqual(true);
      expect(isPermitted(bridge, RELAYER, withArgument(0, '1'))).toEqual(true);
    });

    it.each([
      ['no relayer fee', '0', true],
      ['a relayer fee of exactly 5 USDC', (5_000_000).toString(16), true],
      ['a relayer fee one unit over 5 USDC', (5_000_001).toString(16), false],
      ['the whole amount as the relayer fee', (99_999_622).toString(16), false],
      ['the largest possible relayer fee', 'f'.repeat(64), false],
    ])('with %s, permits the bridge call: %s', async (
      _name: string,
      fee: string,
      expected: boolean,
    ) => {
      const [, bridge] = await getPermissions(arbitrumId, INJECTIVE_ADDRESS);

      expect(isPermitted(bridge, RELAYER, withArgument(4, fee))).toEqual(expected);
    });

    it.each([
      ['to another recipient', 2, '0000000000000000000000001111111111111111111111111111111111111111'],
      ['to another chain', 1, '4'],
      ['of another token', 3, '000000000000000000000000ff970a61a04b1ca14834a43f5de4533ebddb5cc8'],
    ])('refuses a transfer %s', async (_name: string, index: number, word: string) => {
      const [, bridge] = await getPermissions(arbitrumId, INJECTIVE_ADDRESS);

      expect(isPermitted(bridge, RELAYER, withArgument(index, word))).toEqual(false);
    });

    it('refuses the same call on another contract', async () => {
      const [, bridge] = await getPermissions(arbitrumId, INJECTIVE_ADDRESS);

      expect(isPermitted(
        bridge,
        '0x6d4d4A979D766e8b87C985a61B262cfCDf8d6446',
        SKIP_ARBITRUM_BRIDGE_TX.data,
      )).toEqual(false);
    });

    it('refuses a transfer to an Injective account other than the one named', async () => {
      const [, bridge] = await getPermissions(arbitrumId, OTHER_INJECTIVE_ADDRESS);

      expect(bridge.rules[1].params).toEqual([OTHER_INJECTIVE_MINT_RECIPIENT]);
      expect(
        isPermitted(bridge, SKIP_ARBITRUM_BRIDGE_TX.to, SKIP_ARBITRUM_BRIDGE_TX.data),
      ).toEqual(false);
    });

    it('refuses to make a policy for something that is not an Injective account', async () => {
      await expect(getPermissions(arbitrumId, DYDX_ADDRESS)).rejects.toThrow(
        'Expected an Injective address',
      );
      await expect(getPermissions(arbitrumId, '')).rejects.toThrow();
    });

    it('permits approving only the relayer to spend USDC', async () => {
      const [approve] = await getPermissions(arbitrumId, INJECTIVE_ADDRESS);
      const usdc: string = usdcAddressByChainId[arbitrumId];
      const amount: string = '5f5e100'.padStart(64, '0');
      const approveCall = (spender: string): string => {
        return `095ea7b3${spender.slice(2).toLowerCase().padStart(64, '0')}${amount}`;
      };

      expect(isPermitted(approve, usdc, approveCall(RELAYER))).toEqual(true);
      expect(isPermitted(
        approve,
        usdc,
        approveCall('0x1111111111111111111111111111111111111111'),
      )).toEqual(false);
      // Not on another token.
      expect(isPermitted(
        approve,
        '0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8',
        approveCall(RELAYER),
      )).toEqual(false);
    });
  });

  describe('getNoCallsPolicy', () => {
    it('permits no calls', () => {
      const policy = toCallPolicy(getNoCallsPolicy());

      expect((policy.policyParams as unknown as { permissions: Permission[] }).permissions)
        .toEqual([]);
    });
  });

  describe('approvalAllowsInjectiveBridge', () => {
    const owner = privateKeyToAccount(generatePrivateKey());
    const client = createOfflineClient(arbitrum);

    it('is true for an approval of the Injective route to the account', async () => {
      const approval: string = await serializeApprovalWithKey({
        client,
        owner,
        sessionKeyAddress: SESSION_KEY_ADDRESS,
        policy: await callPolicyByChainId[arbitrumId](INJECTIVE_ADDRESS),
      });

      expect(approvalAllowsInjectiveBridge(approval, INJECTIVE_ADDRESS)).toEqual(true);
      // The same account, written in upper case.
      expect(
        approvalAllowsInjectiveBridge(approval, INJECTIVE_ADDRESS.toUpperCase()),
      ).toEqual(true);
    });

    it('is false for an approval of the Injective route to another account', async () => {
      const approval: string = await serializeApprovalWithKey({
        client,
        owner,
        sessionKeyAddress: SESSION_KEY_ADDRESS,
        policy: await callPolicyByChainId[arbitrumId](OTHER_INJECTIVE_ADDRESS),
      });

      expect(approvalAllowsInjectiveBridge(approval, INJECTIVE_ADDRESS)).toEqual(false);
      expect(approvalAllowsInjectiveBridge(approval, OTHER_INJECTIVE_ADDRESS)).toEqual(true);
    });

    it('is false for the Injective account with the bytes of the dYdX address', async () => {
      const approval: string = await serializeApprovalWithKey({
        client,
        owner,
        sessionKeyAddress: SESSION_KEY_ADDRESS,
        policy: await callPolicyByChainId[arbitrumId](INJECTIVE_ADDRESS),
      });

      expect(
        approvalAllowsInjectiveBridge(approval, INJECTIVE_ADDRESS_WITH_DYDX_BYTES),
      ).toEqual(false);
    });

    it('is false for the approval a new account is set up with', async () => {
      const approval: string = await serializeApprovalWithKey({
        client,
        owner,
        sessionKeyAddress: SESSION_KEY_ADDRESS,
        policy: getNoCallsPolicy(),
      });

      expect(approvalAllowsInjectiveBridge(approval, INJECTIVE_ADDRESS)).toEqual(false);
    });

    it('is false when asked about something that is not an Injective account', async () => {
      const approval: string = await serializeApprovalWithKey({
        client,
        owner,
        sessionKeyAddress: SESSION_KEY_ADDRESS,
        policy: await callPolicyByChainId[arbitrumId](INJECTIVE_ADDRESS),
      });

      expect(approvalAllowsInjectiveBridge(approval, DYDX_ADDRESS)).toEqual(false);
      expect(approvalAllowsInjectiveBridge(approval, '')).toEqual(false);
    });

    it('is false for an approval of the earlier route through Noble', async () => {
      const approval: string = await serializeApprovalWithKey({
        client,
        owner,
        sessionKeyAddress: SESSION_KEY_ADDRESS,
        policy: getNobleRoutePolicy(),
      });

      expect(approvalAllowsInjectiveBridge(approval, INJECTIVE_ADDRESS)).toEqual(false);
    });

    it.each([
      ['an empty string', ''],
      ['text that is not base64', '0xapproval'],
      ['base64 that is not JSON', Buffer.from('not json').toString('base64')],
      ['JSON of another shape', Buffer.from('{"permissionParams":{}}').toString('base64')],
      ['JSON null', Buffer.from('null').toString('base64')],
    ])('is false for %s', (_name: string, approval: string) => {
      expect(approvalAllowsInjectiveBridge(approval, INJECTIVE_ADDRESS)).toEqual(false);
    });
  });
});
