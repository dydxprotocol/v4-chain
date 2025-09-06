import { logger } from '@dydxprotocol-indexer/base';
import { PermissionApprovalTable } from '@dydxprotocol-indexer/postgres';
import { TurnkeyApiClient, Turnkey as TurnkeyServerSDK } from '@turnkey/sdk-server';
import { createAccount } from '@turnkey/viem';
import { signerToEcdsaValidator } from '@zerodev/ecdsa-validator';
import { serializePermissionAccount, toPermissionValidator } from '@zerodev/permissions';
import { toCallPolicy } from '@zerodev/permissions/policies';
import { toECDSASigner } from '@zerodev/permissions/signers';
import { addressToEmptyAccount, createKernelAccount } from '@zerodev/sdk';
import { KERNEL_V3_1, KERNEL_V3_3 } from '@zerodev/sdk/constants';
import bs58 from 'bs58';
import { LocalAccount } from 'viem';
import {
  avalanche, optimism, base, arbitrum,
} from 'viem/chains';

import config from '../config';
import { avalancheCallPolicy, chainIdToCallPolicy } from '../lib/call-policies';
import { entryPoint } from '../lib/smart-contract-constants';
import { publicClients } from './alchemy-helpers';
import { nobleToSolana } from './skip-helper';

export class PolicyEngine {
  private turnkeySenderClient: TurnkeyApiClient;

  constructor(turnkeySenderClient?: TurnkeyApiClient) {
    if (!turnkeySenderClient) {
      this.turnkeySenderClient = new TurnkeyServerSDK({
        apiBaseUrl: config.TURNKEY_API_BASE_URL,
        apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY,
        apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY,
        defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
      }).apiClient();
    } else {
      this.turnkeySenderClient = turnkeySenderClient;
    }
  }

  async configurePolicy(
    suborgId: string, // sender api must have access to this.
    fromAddress: string,
  ) {

    const turnkeyAccount = await createAccount({
      // @ts-ignore
      client: this.turnkeySenderClient,
      organizationId: suborgId,
      signWith: fromAddress,
    });

    const chains = [arbitrum.id.toString(), base.id.toString(), optimism.id.toString()];
    for (const chain of chains) {
      try {
        const approval = await getApprovalFor7702Evm(turnkeyAccount, chain);
        logger.info({
          at: 'policy-controller#configurePolicy',
          message: `Approval obtained for chain ${chain} and suborg ${suborgId}`,
        });
        await PermissionApprovalTable.create({
          suborg_id: suborgId,
          chain_id: chain,
          approval,
        });
      } catch (error) {
        logger.error({ at: 'policy-controller#configurePolicy', message: 'Error configuring policy', error });
        throw error;
      }
    }

    const avalancheApproval = await getApprovalForAvalanche(turnkeyAccount);
    await PermissionApprovalTable.create({
      suborg_id: suborgId,
      chain_id: avalanche.id.toString(),
      approval: avalancheApproval,
    });
  }

  async getAPIUserId(suborgId: string): Promise<string> {
    // query users from turnkey for   the suborg id and find out what the api user's
    // userId is to configure a policy on it.
    const users = await this.turnkeySenderClient.getUsers({ organizationId: suborgId });
    const apiUser = users.users.filter((user) => user.userName === 'API User');
    if (apiUser.length !== 1) {
      throw new Error(`Expected 1 API user, got ${apiUser.length}`);
    }
    const userId = apiUser[0].userId;

    return userId;
  }

  async configureSolanaPolicy(dydxAddress: string, suborgId: string) {
    const userId = await this.getAPIUserId(suborgId);
    // get noble forwarding address for the dydx user.
    const nobleForwardingAddress = await getNobleForwardingAddress(dydxAddress);
    const solanaAddress = nobleToSolana(nobleForwardingAddress);

    const depositForBurnWithCallerHex = 'a7de137255150e76';
    // left pads it the hex version of the solana address so that total is 32 bytes.
    const hexData = solanaAddressToPaddedHex(solanaAddress);
    await this.turnkeySenderClient.createPolicy({
      organizationId: suborgId,
      policyName: 'Solana Bridging Policy',
      condition: `solana.tx.instructions[0].instruction_data_hex[0..16] == '${depositForBurnWithCallerHex}' && solana.tx.instructions[0].instruction_data_hex[40..104] == '${hexData}'`,
      consensus: `approvers.any(user, user.id == '${userId}')`,
      effect: 'EFFECT_ALLOW',
      notes: 'Solana bridge policy',
    });

    // TODO remove self from root quorum.
  }

  async removeSelfFromRootQuorum(suborgId: string) {
    const users = await this.turnkeySenderClient.getUsers({ organizationId: suborgId });
    const userIds = users.users.filter((user) => user.userName !== 'API User').map((x) => x.userId);

    await this.turnkeySenderClient.updateRootQuorum({
      organizationId: suborgId,
      threshold: 1,
      userIds,
    });
  }
}

async function getNobleForwardingAddress(nobleAddress: string): Promise<string> {
  const dydxNobleChannel = 33;
  const endpoint = `https://api.noble.xyz/noble/forwarding/v1/address/channel-${dydxNobleChannel}/${nobleAddress}/`;

  const ac = new AbortController();
  const timeout = setTimeout(() => ac.abort(), 10_000);
  try {
    const response = await fetch(endpoint, {
      signal: ac.signal,
    });
    if (!response.ok) {
      throw new Error(`failed to get a forwarding address: ${response.statusText}`);
    }
    const data = await response.json();
    if (!data || (data && !data.address)) {
      throw new Error('failed to get a forwarding address');
    }
    return data.address;
  } catch (e) {
    throw new Error(`failed to get a forwarding address: ${e}`);
  } finally {
    clearTimeout(timeout);
  }
}

function solanaAddressToPaddedHex(solanaAddress: string): string {
  // Remove '0x' if present and ensure lowercase
  let hex = solanaAddress.startsWith('0x')
    ? solanaAddress.slice(2)
    : solanaAddress;

  // If it's a base58 Solana address, decode to bytes and then to hex
  // Solana addresses are usually base58, not hex
  // We'll use bs58 to decode
  // If bs58 is not available, throw an error
  let bytes: Uint8Array;
  try {
    // @ts-ignore
    bytes = bs58.decode(solanaAddress);
    hex = Buffer.from(bytes).toString('hex');
  } catch (e) {
    // If bs58 is not available or decoding fails, assume it's already hex
    if (!/^[0-9a-fA-F]+$/.test(hex)) {
      throw new Error('Invalid Solana address: not base58 or hex');
    }
  }

  // Pad to 32 bytes (64 hex chars)
  if (hex.length > 64) {
    throw new Error('Hex Solana address is longer than 32 bytes');
  }
  return hex.padStart(64, '0');
}

async function getApprovalForAvalanche(turnkeyAccount: LocalAccount) {
  const ecdsaValidator = await signerToEcdsaValidator(publicClients[avalanche.id.toString()], {
    entryPoint,
    signer: turnkeyAccount,
    kernelVersion: KERNEL_V3_1,
  });

  // Create an "empty account" as the signer -- you only need the public
  // key (address) to do this.
  const emptyAccount = addressToEmptyAccount(config.MASTER_SIGNER_PUBLIC as `0x${string}`);
  const emptySessionKeySigner = await toECDSASigner({ signer: emptyAccount });

  const permissionPlugin = await toPermissionValidator(publicClients[avalanche.id.toString()], {
    entryPoint,
    signer: emptySessionKeySigner,
    policies: [
      toCallPolicy(avalancheCallPolicy),
    ],
    kernelVersion: KERNEL_V3_1,
  });

  const sessionKeyAccount = await createKernelAccount(publicClients[avalanche.id.toString()], {
    entryPoint,
    plugins: {
      sudo: ecdsaValidator,
      regular: permissionPlugin,
    },
    kernelVersion: KERNEL_V3_1,
  });

  return serializePermissionAccount(sessionKeyAccount);
}

async function getApprovalFor7702Evm(turnkeyAccount: LocalAccount, chainId: string) {

  const kernelVersion = KERNEL_V3_3;
  const emptyAccount = addressToEmptyAccount(config.MASTER_SIGNER_PUBLIC as `0x${string}`);
  const emptySessionKeySigner = await toECDSASigner({ signer: emptyAccount });
  const permissionPlugin = await toPermissionValidator(publicClients[chainId], {
    entryPoint,
    kernelVersion,
    signer: emptySessionKeySigner,
    policies: [
      toCallPolicy(chainIdToCallPolicy[chainId]),
    ],
  });

  const sessionAccount = await createKernelAccount(publicClients[chainId], {
    entryPoint,
    kernelVersion,
    eip7702Account: turnkeyAccount,
    plugins: {
      regular: permissionPlugin,
    },
  });
  return serializePermissionAccount(sessionAccount);
}
