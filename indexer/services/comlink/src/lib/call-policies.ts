import { CallPolicyVersion, ParamCondition, type CallPolicyParams } from '@zerodev/permissions/policies';
import type { Abi } from 'viem';
import {
  arbitrum, avalanche, base, mainnet, optimism,
} from 'viem/chains';

import { abi, cctpV2RelayerAbi } from './bridge-abi';
import { getInjectiveMintRecipient } from './injective-address';
import {
  EVM_BRIDGE_MAX_RELAY_FEE,
  INJECTIVE_CCTP_DOMAIN,
  SKIP_CCTP_V2_RELAYER_ADDRESS,
  usdcAddressByChainId,
} from './smart-contract-constants';

// Neither call sends the chain's native token.
const valueLimit: bigint = BigInt(0);

/**
 * Constructs the policy that lets the bridge session key move a user's USDC on `chainId`, and
 * nothing else. It permits the two calls of a deposit:
 *
 * 1. approving Skip's CCTP v2 relayer to spend the user's USDC, and
 * 2. asking that relayer to burn USDC for minting on Injective to the Injective account the user
 *    named (see `injective-address.ts`), keeping at most `EVM_BRIDGE_MAX_RELAY_FEE` of it as
 *    its fee.
 *
 * Deposits of the native token are not covered: Skip has no route for them to Injective.
 *
 * @param chainId - one of ethereum, arbitrum, base, optimism, avalanche.
 * @returns a function from the user's Injective address to the policy.
 */
function constructPolicy(chainId: string): (injectiveAddress: string) => Promise<CallPolicyParams<Abi, `0x${string}`>> {
  const usdc = usdcAddressByChainId[chainId] as `0x${string}`;
  return (injectiveAddress: string) => {
    return Promise.resolve({
      policyVersion: CallPolicyVersion.V0_0_5,
      permissions: [
        {
          target: usdc,
          abi,
          valueLimit,
          functionName: 'approve',
          args: [
            {
              // The only spender is the relayer.
              condition: ParamCondition.EQUAL,
              value: SKIP_CCTP_V2_RELAYER_ADDRESS,
            },
            null,
          ],
        },
        {
          target: SKIP_CCTP_V2_RELAYER_ADDRESS,
          abi: cctpV2RelayerAbi,
          valueLimit,
          functionName: 'requestCCTPTransferWithCaller',
          args: [
            null,
            {
              // The destination chain is Injective.
              condition: ParamCondition.EQUAL,
              value: INJECTIVE_CCTP_DOMAIN,
            },
            {
              // The mint recipient is the Injective account the user named.
              condition: ParamCondition.EQUAL,
              value: getInjectiveMintRecipient(injectiveAddress),
            },
            {
              // The token burned is USDC.
              condition: ParamCondition.EQUAL,
              value: usdc,
            },
            {
              // The relayer keeps this much of the USDC as its fee, so it is capped.
              condition: ParamCondition.LESS_THAN_OR_EQUAL,
              value: BigInt(EVM_BRIDGE_MAX_RELAY_FEE),
            },
            null,
            null,
            null,
            null,
          ],
        },
      ],
    });
  };
}

export const callPolicyByChainId: Record<string, (injectiveAddress: string) => Promise<CallPolicyParams<Abi, `0x${string}`>>> = {
  [mainnet.id.toString()]: constructPolicy(mainnet.id.toString()),
  [arbitrum.id.toString()]: constructPolicy(arbitrum.id.toString()),
  [avalanche.id.toString()]: constructPolicy(avalanche.id.toString()),
  [optimism.id.toString()]: constructPolicy(optimism.id.toString()),
  [base.id.toString()]: constructPolicy(base.id.toString()),
};

/**
 * The policy a new account starts with: the bridge session key is permitted no calls at all.
 *
 * An account is set up before the user has named the Injective account their deposits go to, so
 * there is nothing yet that the key could be permitted to do. The approval made with this policy
 * still matters: on the EIP-7702 chains it carries the authorization that makes the user's
 * address a smart account, which the approval the user makes later is built on.
 */
export function getNoCallsPolicy(): CallPolicyParams<Abi, `0x${string}`> {
  return {
    policyVersion: CallPolicyVersion.V0_0_5,
    permissions: [],
  };
}

/**
 * Whether a stored approval permits bridging deposits to `injectiveAddress`. It does not if it
 * was made for the earlier route through Noble, for another Injective account, or when the
 * account was set up. A user's approval is replaced when they approve the deposit policy; until
 * then their deposits cannot be bridged.
 *
 * @param approval - a serialized permission account, as stored in the permission_approval table.
 */
export function approvalAllowsInjectiveBridge(approval: string, injectiveAddress: string): boolean {
  try {
    const mintRecipient: string = getInjectiveMintRecipient(injectiveAddress).toLowerCase();
    const params = JSON.parse(Buffer.from(approval, 'base64').toString('utf8'));
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return params.permissionParams.policies.some((policy: any) => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      return policy.policyParams.permissions?.some((permission: any) => {
        return String(permission.target).toLowerCase() ===
            SKIP_CCTP_V2_RELAYER_ADDRESS.toLowerCase() &&
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          permission.rules?.some((rule: any) => {
            return rule.params?.some(
              (param: unknown) => String(param).toLowerCase() === mintRecipient,
            );
          });
      });
    });
  } catch (error) {
    // Not a serialized permission account, not one with a call policy, or not an Injective
    // address.
    return false;
  }
}
