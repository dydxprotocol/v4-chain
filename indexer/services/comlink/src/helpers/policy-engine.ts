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
  avalanche, optimism, base, arbitrum, mainnet,
} from 'viem/chains';

import config from '../config';
import { callPolicyByChainId } from '../lib/call-policies';
import { entryPoint } from '../lib/smart-contract-constants';
import { publicClients } from './alchemy-helpers';
import { getNobleForwardingAddress, nobleToSolana } from './skip-helper';

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
    dydxAddress: string,
  ) {

    const turnkeyAccount = await createAccount({
      // @ts-ignore
      client: this.turnkeySenderClient,
      organizationId: suborgId,
      signWith: fromAddress,
    });

    const chains = [
      arbitrum.id.toString(), mainnet.id.toString(), optimism.id.toString(), base.id.toString(),
    ];

    // Process all chains in parallel
    await Promise.all([
      // Process EVM chains in parallel
      ...chains.map(async (chain) => {
        try {
          const exists = await PermissionApprovalTable.findBySuborgIdAndChainId(suborgId, chain);
          if (exists) {
            logger.info({
              at: 'policy-controller#configurePolicy',
              message: `Policy already exists for chain ${chain} and suborg ${suborgId}, skipping`,
            });
            return;
          }
          const approval = await getApprovalFor7702Evm(turnkeyAccount, chain, dydxAddress);
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
          logger.error({ at: 'policy-controller#configurePolicy', message: `Error configuring policy for chain ${chain}`, error });
          throw error;
        }
      }),
      // Process Avalanche chain
      (async () => {
        try {
          const exists = await PermissionApprovalTable.findBySuborgIdAndChainId(
            suborgId,
            avalanche.id.toString(),
          );
          if (!exists) {
            const avalancheApproval = await getApprovalForAvalanche(turnkeyAccount, dydxAddress);
            await PermissionApprovalTable.create({
              suborg_id: suborgId,
              chain_id: avalanche.id.toString(),
              approval: avalancheApproval,
            });
          } else {
            logger.info({
              at: 'policy-controller#configurePolicy',
              message: `Policy already exists for chain ${avalanche.id.toString()} and suborg ${suborgId}, skipping`,
            });
          }
        } catch (error) {
          logger.error({ at: 'policy-controller#configurePolicy', message: 'Error configuring avalanche policy', error });
        }
      })(),
      // Process Solana chain
      (async () => {
        try {
          await this.configureSolanaPolicy(dydxAddress, suborgId);
        } catch (error) {
          logger.error({ at: 'policy-controller#configurePolicy', message: 'Error configuring solana policy', error });
        }
      })(),
    ]);
  }

  async getAPIUserId(suborgId: string): Promise<string> {
    // query users from turnkey for the suborg id and find out what the api user's
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
    // Get noble forwarding address for the dydx user.
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

async function getApprovalForAvalanche(turnkeyAccount: LocalAccount, dydxAddress: string) {
  const ecdsaValidator = await signerToEcdsaValidator(publicClients[avalanche.id.toString()], {
    entryPoint,
    signer: turnkeyAccount,
    kernelVersion: KERNEL_V3_1,
  });

  // Create an "empty account" as the signer -- you only need the public
  // key (address) to do this.
  const emptyAccount = addressToEmptyAccount(config.APPROVAL_SIGNER_PUBLIC_ADDRESS as `0x${string}`);
  const emptySessionKeySigner = await toECDSASigner({ signer: emptyAccount });

  const permissionPlugin = await toPermissionValidator(publicClients[avalanche.id.toString()], {
    entryPoint,
    signer: emptySessionKeySigner,
    policies: [
      toCallPolicy(await callPolicyByChainId[avalanche.id.toString()](dydxAddress)),
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

async function getApprovalFor7702Evm(
  turnkeyAccount: LocalAccount,
  chainId: string,
  dydxAddress: string,
) {
  const callPolicy = await callPolicyByChainId[chainId](dydxAddress);
  const kernelVersion = KERNEL_V3_3;
  const emptyAccount = addressToEmptyAccount(config.APPROVAL_SIGNER_PUBLIC_ADDRESS as `0x${string}`);
  const emptySessionKeySigner = await toECDSASigner({ signer: emptyAccount });
  const permissionPlugin = await toPermissionValidator(publicClients[chainId], {
    entryPoint,
    kernelVersion,
    signer: emptySessionKeySigner,
    policies: [
      toCallPolicy(
        callPolicy,
      ),
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
