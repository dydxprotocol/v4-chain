import { signerToEcdsaValidator } from '@zerodev/ecdsa-validator';
import { serializePermissionAccount, toPermissionValidator } from '@zerodev/permissions';
import {
  CallPolicyVersion, ParamCondition, toCallPolicy, type CallPolicyParams,
} from '@zerodev/permissions/policies';
import { toECDSASigner } from '@zerodev/permissions/signers';
import { addressToEmptyAccount, createKernelAccount } from '@zerodev/sdk';
import { getEntryPoint, KERNEL_V3_1, KERNEL_V3_3 } from '@zerodev/sdk/constants';
import {
  Abi, Address, Chain, createPublicClient, custom, LocalAccount, PublicClient, toHex,
} from 'viem';
import { avalanche } from 'viem/chains';

import { abi } from '../../src/lib/bridge-abi';

/**
 * A client of a chain on which nothing has happened: no account has code or has sent a
 * transaction, and every contract call reverts. Lets the ZeroDev SDK build and serialize
 * accounts without a network.
 */
export function createOfflineClient(chain: Chain): PublicClient {
  return createPublicClient({
    chain,
    transport: custom({
      // eslint-disable-next-line @typescript-eslint/require-await
      request: async ({ method }: { method: string }) => {
        switch (method) {
          case 'eth_chainId':
            return toHex(chain.id);
          case 'eth_getCode':
            return '0x';
          case 'eth_getTransactionCount':
            return '0x0';
          case 'eth_call':
            throw new Error('execution reverted');
          default:
            throw new Error(`unexpected request: ${method}`);
        }
      },
    // A reverted call is an answer, not a failure to be retried.
    }, { retryCount: 0 }),
  }) as PublicClient;
}

/**
 * Serializes a user's smart account with a session key enabled under `policy`, the way the SDK
 * does it when it holds the user's key: the SDK signs whatever is needed with `owner`.
 */
export async function serializeApprovalWithKey({
  client,
  owner,
  sessionKeyAddress,
  policy,
  accountAddress,
}: {
  client: PublicClient,
  owner: LocalAccount,
  sessionKeyAddress: Address,
  policy: CallPolicyParams<Abi, `0x${string}`>,
  // Avalanche only: the address of the user's smart account.
  accountAddress?: Address,
}): Promise<string> {
  const entryPoint = getEntryPoint('0.7');
  const isAvalanche: boolean = client.chain!.id === avalanche.id;
  const kernelVersion = isAvalanche ? KERNEL_V3_1 : KERNEL_V3_3;
  const regular = await toPermissionValidator(client, {
    entryPoint,
    kernelVersion,
    signer: await toECDSASigner({ signer: addressToEmptyAccount(sessionKeyAddress) }),
    policies: [toCallPolicy(policy)],
  });
  if (isAvalanche) {
    const sudo = await signerToEcdsaValidator(client, { entryPoint, signer: owner, kernelVersion });
    return serializePermissionAccount(await createKernelAccount(client, {
      entryPoint,
      kernelVersion,
      plugins: { sudo, regular },
      address: accountAddress,
    }));
  }
  return serializePermissionAccount(await createKernelAccount(client, {
    entryPoint,
    kernelVersion,
    eip7702Account: owner,
    plugins: { regular },
  }));
}

/**
 * A policy like the one accounts were given for the earlier deposit route: it permits a CCTP v1
 * transfer to a Noble forwarding address.
 */
export function getNobleRoutePolicy(): CallPolicyParams<Abi, `0x${string}`> {
  return {
    policyVersion: CallPolicyVersion.V0_0_5,
    permissions: [
      {
        target: '0x6d4d4A979D766e8b87C985a61B262cfCDf8d6446',
        abi,
        valueLimit: BigInt(0),
        functionName: 'requestCCTPTransferWithCaller',
        args: [
          null,
          null,
          {
            condition: ParamCondition.EQUAL,
            value: '0x000000000000000000000000a8a2d0a34dc2e9a6d2e2c2e2f6a9e7f6e3c1b0a9',
          },
          null,
          null,
          null,
        ],
      },
    ],
  } as CallPolicyParams<Abi, `0x${string}`>;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function decodeApproval(approval: string): any {
  return JSON.parse(Buffer.from(approval, 'base64').toString('utf8'));
}
