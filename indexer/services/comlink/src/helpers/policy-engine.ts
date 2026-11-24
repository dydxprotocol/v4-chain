import { logger } from '@dydxprotocol-indexer/base';
import {
  PermissionApprovalTable,
  Transaction,
  TurnkeyUserFromDatabase,
  TurnkeyUsersTable,
} from '@dydxprotocol-indexer/postgres';
import { TurnkeyApiClient, Turnkey as TurnkeyServerSDK } from '@turnkey/sdk-server';
import { createAccount } from '@turnkey/viem';
import { signerToEcdsaValidator } from '@zerodev/ecdsa-validator';
import { serializePermissionAccount, toPermissionValidator } from '@zerodev/permissions';
import { toCallPolicy, type CallPolicyParams } from '@zerodev/permissions/policies';
import { toECDSASigner } from '@zerodev/permissions/signers';
import {
  addressToEmptyAccount, createKernelAccount, CreateKernelAccountReturnType,
} from '@zerodev/sdk';
import { KERNEL_V3_1, KERNEL_V3_3 } from '@zerodev/sdk/constants';
import {
  Abi, Address, getTypesForEIP712Domain, LocalAccount, PublicClient, verifyTypedData,
} from 'viem';
import {
  avalanche, optimism, base, arbitrum, mainnet,
} from 'viem/chains';

import config from '../config';
import {
  approvalAllowsInjectiveBridge,
  callPolicyByChainId,
  getNoCallsPolicy,
} from '../lib/call-policies';
import { TurnkeyError } from '../lib/errors';
import { getInjectiveMintRecipient, normalizeInjectiveAddress } from '../lib/injective-address';
import {
  entryPoint,
  INJECTIVE_CCTP_DOMAIN,
  SOLANA_BRIDGE_MAX_RELAY_FEE,
  SOLANA_CCTP_V2_DEPOSIT_FOR_BURN_DISCRIMINATOR,
  SOLANA_CCTP_V2_PROGRAM,
  SOLANA_TOKEN_PROGRAM,
} from '../lib/smart-contract-constants';
import {
  DepositPolicyEvmApproval,
  DepositPolicyResponse,
  DepositPolicySignature,
  DepositPolicyStatus,
  DepositPolicyTurnkeyPolicy,
  DepositPolicyTypedData,
} from '../types';
import { publicClients } from './alchemy-helpers';

// EVM chains where the user's own address becomes the smart account, through EIP-7702.
const EIP_7702_CHAIN_IDS: string[] = [
  arbitrum.id.toString(), mainnet.id.toString(), optimism.id.toString(), base.id.toString(),
];
// Avalanche does not support EIP-7702, so the user has a separate smart account there.
const AVALANCHE_CHAIN_ID: string = avalanche.id.toString();
const EVM_CHAIN_IDS: string[] = [...EIP_7702_CHAIN_IDS, AVALANCHE_CHAIN_ID];

export const SOLANA_BRIDGE_POLICY_NAME: string = 'Injective Bridging Policy';
// The policy for the earlier deposit route, through Noble.
export const OBSOLETE_SOLANA_BRIDGE_POLICY_NAME: string = 'Solana Bridging Policy';

/**
 * The Turnkey policy that lets the bridge API user sign, with the user's Solana key, the
 * transaction of a deposit and nothing else. That transaction has two instructions:
 *
 * 1. CCTP v2 `deposit_for_burn`, whose data is an 8-byte discriminator, the amount (u64), the
 *    destination domain (u32, little-endian) and the 32-byte mint recipient. The domain must be
 *    Injective's and the recipient the Injective account the user named.
 * 2. A token transfer that pays Skip's relayer its fee, capped so that it cannot be used to
 *    move the user's balance.
 */
export function getSolanaBridgePolicy(
  injectiveAddress: string,
  apiUserId: string,
): DepositPolicyTurnkeyPolicy {
  const domain: Buffer = Buffer.alloc(4);
  domain.writeUInt32LE(INJECTIVE_CCTP_DOMAIN);
  const mintRecipient: string = getInjectiveMintRecipient(injectiveAddress).slice(2).toLowerCase();
  const burn: string = 'solana.tx.instructions[0]';
  return {
    policyName: SOLANA_BRIDGE_POLICY_NAME,
    effect: 'EFFECT_ALLOW',
    condition: [
      'solana.tx.instructions.count() == 2',
      `${burn}.program_key == '${SOLANA_CCTP_V2_PROGRAM}'`,
      `${burn}.instruction_data_hex[0..16] == '${SOLANA_CCTP_V2_DEPOSIT_FOR_BURN_DISCRIMINATOR}'`,
      `${burn}.instruction_data_hex[32..40] == '${domain.toString('hex')}'`,
      `${burn}.instruction_data_hex[40..104] == '${mintRecipient}'`,
      `solana.tx.instructions[1].program_key == '${SOLANA_TOKEN_PROGRAM}'`,
      'solana.tx.spl_transfers.count() == 1',
      `solana.tx.spl_transfers[0].amount <= ${SOLANA_BRIDGE_MAX_RELAY_FEE}`,
    ].join(' && '),
    consensus: `approvers.any(user, user.id == '${apiUserId}')`,
    notes: 'Lets dYdX bridge USDC deposits to this account on Injective',
  };
}

/**
 * Builds the user's smart account on `chainId` with the bridge session key enabled on it,
 * limited to the calls `policy` permits.
 *
 * @param owner - the user's EVM account. It only needs to be able to sign if the account is
 * serialized without an enable signature.
 * @param accountAddress - the address of the smart account, if known. Only used on Avalanche,
 * where it differs from the owner's address.
 */
async function buildSessionAccount(
  chainId: string,
  owner: LocalAccount,
  policy: CallPolicyParams<Abi, `0x${string}`>,
  accountAddress?: Address,
): Promise<CreateKernelAccountReturnType<'0.7'>> {
  const client: PublicClient = publicClients[chainId];
  const kernelVersion = chainId === AVALANCHE_CHAIN_ID ? KERNEL_V3_1 : KERNEL_V3_3;
  // Create an "empty account" as the signer -- you only need the public
  // key (address) to do this.
  const sessionKeySigner = await toECDSASigner({
    signer: addressToEmptyAccount(config.APPROVAL_SIGNER_PUBLIC_ADDRESS as `0x${string}`),
  });
  const permissionPlugin = await toPermissionValidator(client, {
    entryPoint,
    kernelVersion,
    signer: sessionKeySigner,
    policies: [
      toCallPolicy(policy),
    ],
  });

  if (chainId === AVALANCHE_CHAIN_ID) {
    const ecdsaValidator = await signerToEcdsaValidator(client, {
      entryPoint,
      signer: owner,
      kernelVersion,
    });
    return createKernelAccount(client, {
      entryPoint,
      kernelVersion,
      plugins: {
        sudo: ecdsaValidator,
        regular: permissionPlugin,
      },
      address: accountAddress,
    });
  }
  return createKernelAccount(client, {
    entryPoint,
    kernelVersion,
    eip7702Account: owner,
    plugins: {
      regular: permissionPlugin,
    },
  });
}

/**
 * The EIP-7702 authorization kept in a stored approval, if it has one. An approval made before
 * the account was first used carries the authorization that turns the user's address into a
 * smart account, and a new approval for the same account needs to carry it too.
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function getStoredEip7702Authorization(approval: string | undefined): any {
  if (approval === undefined) {
    return undefined;
  }
  try {
    return JSON.parse(Buffer.from(approval, 'base64').toString('utf8')).eip7702Auth ?? undefined;
  } catch (error) {
    return undefined;
  }
}

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

  /**
   * Sets up a new user's smart accounts. Runs while the API user is still in the root quorum of
   * the user's sub-organization, so it can sign for the user.
   *
   * This does not yet let deposits be bridged. They go to an Injective account that only the
   * user can name, so the bridge session key is enabled with no calls permitted (see
   * `getNoCallsPolicy`), and no Solana policy is created. The user grants the rest through
   * `getDepositPolicy`, as the user of an account set up for an earlier route does.
   */
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

    // Process all chains in parallel
    await Promise.all([
      // Process EVM chains in parallel
      ...EIP_7702_CHAIN_IDS.map(async (chain) => {
        try {
          const exists = await PermissionApprovalTable.findBySuborgIdAndChainId(suborgId, chain);
          if (exists) {
            logger.info({
              at: 'policy-controller#configurePolicy',
              message: `Policy already exists for chain ${chain} and suborg ${suborgId}, skipping`,
            });
            return;
          }
          const approval = await serializePermissionAccount(
            await buildSessionAccount(chain, turnkeyAccount, getNoCallsPolicy()),
          );
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
            AVALANCHE_CHAIN_ID,
          );
          if (!exists) {
            const avalancheApproval = await serializePermissionAccount(
              await buildSessionAccount(AVALANCHE_CHAIN_ID, turnkeyAccount, getNoCallsPolicy()),
            );
            await PermissionApprovalTable.create({
              suborg_id: suborgId,
              chain_id: AVALANCHE_CHAIN_ID,
              approval: avalancheApproval,
            });
          } else {
            logger.info({
              at: 'policy-controller#configurePolicy',
              message: `Policy already exists for chain ${AVALANCHE_CHAIN_ID} and suborg ${suborgId}, skipping`,
            });
          }
        } catch (error) {
          logger.error({ at: 'policy-controller#configurePolicy', message: 'Error configuring avalanche policy', error });
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

  async removeSelfFromRootQuorum(suborgId: string) {
    const users = await this.turnkeySenderClient.getUsers({ organizationId: suborgId });
    const userIds = users.users.filter((user) => user.userName !== 'API User').map((x) => x.userId);

    await this.turnkeySenderClient.updateRootQuorum({
      organizationId: suborgId,
      threshold: 1,
      userIds,
    });
  }

  /**
   * Works out what a user has to do before their deposits can be bridged to `injectiveAddress`,
   * the Injective account they name for them.
   *
   * Once an account is set up, the API user is no longer in the root quorum of the user's
   * sub-organization and cannot change what it is permitted to do there. The user has to grant
   * it, in two parts:
   *
   * - EVM: for each chain, sign typed data that enables the bridge session key on their smart
   *   account, and hand the signatures to `submitDepositApprovals`.
   * - Solana: create a Turnkey policy in their sub-organization.
   *
   * Both name `injectiveAddress` as the only recipient, so the address is part of what the user
   * approves. Asking changes nothing: the address is only stored by `submitDepositApprovals`.
   *
   * @throws if `injectiveAddress` is not the address of an Injective account.
   */
  async getDepositPolicy(
    user: TurnkeyUserFromDatabase,
    injectiveAddress: string,
  ): Promise<DepositPolicyResponse> {
    const dydxAddress: string = user.dydx_address!;
    const recipient: string = normalizeInjectiveAddress(injectiveAddress);
    const [approvals, solana]: [
      DepositPolicyEvmApproval[],
      DepositPolicyResponse['solana'],
    ] = await Promise.all([
      Promise.all(
        EVM_CHAIN_IDS.map((chainId: string) => this.getEvmApproval(user, recipient, chainId)),
      ),
      this.getSolanaPolicy(user, recipient),
    ]);
    return {
      dydxAddress,
      injectiveAddress: recipient,
      evm: {
        address: user.evm_address,
        approvals,
      },
      solana,
    };
  }

  /**
   * Stores the approvals made by signing the typed data from `getDepositPolicy`, and
   * `injectiveAddress` as the account the user's deposits are bridged to. Each signature is
   * checked against the typed data for its chain and that address as it stands now, so only the
   * user can change their approval or the address, and only to the policy this service asks for.
   *
   * The approvals and the address are stored together or not at all. A chain that is not among
   * `signatures` keeps its approval: if that one is for another address, deposits on that chain
   * are not bridged until the user approves it as well.
   *
   * @returns the ids of the chains whose approval was stored.
   * @throws TurnkeyError if a signature is not the user's signature of the typed data.
   */
  async submitDepositApprovals(
    user: TurnkeyUserFromDatabase,
    injectiveAddress: string,
    signatures: DepositPolicySignature[],
  ): Promise<string[]> {
    const recipient: string = normalizeInjectiveAddress(injectiveAddress);
    const chainIds: string[] = signatures.map((s: DepositPolicySignature) => s.chainId);
    if (new Set(chainIds).size !== chainIds.length) {
      throw new TurnkeyError('Each chain may appear only once');
    }
    const unsupported: string | undefined = chainIds.find(
      (chainId: string) => !EVM_CHAIN_IDS.includes(chainId),
    );
    if (unsupported !== undefined) {
      throw new TurnkeyError(`Unsupported chain: ${unsupported}`);
    }

    // Check every signature before storing anything.
    const approvals: { chainId: string, approval: string }[] = await Promise.all(
      signatures.map(async ({ chainId, signature }: DepositPolicySignature) => {
        return {
          chainId,
          approval: await this.buildApproval(user, recipient, chainId, signature),
        };
      }),
    );
    const txId: number = await Transaction.start();
    try {
      for (const { chainId, approval } of approvals) {
        await PermissionApprovalTable.upsert({
          suborg_id: user.suborg_id,
          chain_id: chainId,
          approval,
        }, { txId });
      }
      await TurnkeyUsersTable.updateInjectiveAddressBySuborgId(
        user.suborg_id,
        recipient,
        { txId },
      );
      await Transaction.commit(txId);
    } catch (error) {
      await Transaction.rollback(txId);
      throw error;
    }
    return chainIds;
  }

  private async getEvmApproval(
    user: TurnkeyUserFromDatabase,
    injectiveAddress: string,
    chainId: string,
  ): Promise<DepositPolicyEvmApproval> {
    try {
      const stored: string | undefined = await PermissionApprovalTable.getApprovalForSuborgAndChain(
        user.suborg_id,
        chainId,
      );
      // Deposits are bridged to the address on record, within what the approval permits. Both
      // have to be for this address before there is nothing left to do.
      if (
        stored !== undefined &&
        approvalAllowsInjectiveBridge(stored, injectiveAddress) &&
        user.injective_address === injectiveAddress
      ) {
        return { chainId, status: DepositPolicyStatus.UP_TO_DATE };
      }
      // The user's address became a smart account through the approval made when the account
      // was set up. Without that there is nothing to build on.
      if (stored === undefined && EIP_7702_CHAIN_IDS.includes(chainId)) {
        logger.error({
          at: 'policy-engine#getEvmApproval',
          message: 'No approval on record to renew',
          suborgId: user.suborg_id,
          chainId,
        });
        return { chainId, status: DepositPolicyStatus.UNAVAILABLE };
      }
      const { typedData } = await this.getEnableTypedData(user, injectiveAddress, chainId);
      return { chainId, status: DepositPolicyStatus.ACTION_REQUIRED, typedData };
    } catch (error) {
      logger.error({
        at: 'policy-engine#getEvmApproval',
        message: 'Failed to prepare the deposit approval',
        suborgId: user.suborg_id,
        chainId,
        error,
      });
      return { chainId, status: DepositPolicyStatus.UNAVAILABLE };
    }
  }

  private async getSolanaPolicy(
    user: TurnkeyUserFromDatabase,
    injectiveAddress: string,
  ): Promise<DepositPolicyResponse['solana']> {
    try {
      const [apiUserId, { policies }] = await Promise.all([
        this.getAPIUserId(user.suborg_id),
        this.turnkeySenderClient.getPolicies({ organizationId: user.suborg_id }),
      ]);
      const bridgePolicies = policies.filter((existing) => {
        return existing.policyName === SOLANA_BRIDGE_POLICY_NAME ||
          existing.policyName === OBSOLETE_SOLANA_BRIDGE_POLICY_NAME;
      });
      // The policy is in place if a bridge policy for the Injective route names this account as
      // the recipient. It is recognised by that rather than by its whole condition, so as not
      // to depend on Turnkey returning a condition character for character as it was created.
      // Any other bridge policy is for the earlier route or for another Injective account.
      const recipient: string = `'${getInjectiveMintRecipient(injectiveAddress).slice(2)}'`;
      const isCurrent = (existing: typeof bridgePolicies[number]): boolean => {
        return existing.policyName === SOLANA_BRIDGE_POLICY_NAME &&
          (existing.condition ?? '').toLowerCase().includes(recipient.toLowerCase());
      };
      const obsoletePolicyIds: string[] = bridgePolicies
        .filter((existing) => !isCurrent(existing))
        .map((existing) => existing.policyId);
      if (bridgePolicies.some(isCurrent)) {
        return { status: DepositPolicyStatus.UP_TO_DATE, obsoletePolicyIds };
      }
      return {
        status: DepositPolicyStatus.ACTION_REQUIRED,
        policy: getSolanaBridgePolicy(injectiveAddress, apiUserId),
        obsoletePolicyIds,
      };
    } catch (error) {
      logger.error({
        at: 'policy-engine#getSolanaPolicy',
        message: 'Failed to prepare the Solana deposit policy',
        suborgId: user.suborg_id,
        error,
      });
      return { status: DepositPolicyStatus.UNAVAILABLE, obsoletePolicyIds: [] };
    }
  }

  /**
   * Builds the user's smart account on `chainId` without their key, and the typed data whose
   * signature by the user enables the bridge session key on it.
   */
  private async getEnableTypedData(
    user: TurnkeyUserFromDatabase,
    injectiveAddress: string,
    chainId: string,
  ): Promise<{ account: CreateKernelAccountReturnType<'0.7'>, typedData: DepositPolicyTypedData }> {
    const account: CreateKernelAccountReturnType<'0.7'> = await buildSessionAccount(
      chainId,
      // Stands in for the user's key, which signs elsewhere.
      addressToEmptyAccount(user.evm_address as Address),
      await callPolicyByChainId[chainId](injectiveAddress),
      user.smart_account_address as Address | undefined,
    );
    const typedData = await account.kernelPluginManager.getPluginsEnableTypedData(account.address);
    // The SDK leaves the type of the domain out, to be derived from the domain's fields. Not
    // every signer derives it: one that takes a missing domain type to be empty would sign a
    // different hash. Spelling it out makes this the complete `eth_signTypedData_v4` payload,
    // which hashes the same either way.
    const complete = {
      ...typedData,
      types: {
        EIP712Domain: getTypesForEIP712Domain({ domain: typedData.domain }),
        ...typedData.types,
      },
    };
    // Through JSON, so that it is exactly what the user is given to sign.
    return { account, typedData: JSON.parse(JSON.stringify(complete)) };
  }

  private async buildApproval(
    user: TurnkeyUserFromDatabase,
    injectiveAddress: string,
    chainId: string,
    signature: string,
  ): Promise<string> {
    const { account, typedData } = await this.getEnableTypedData(user, injectiveAddress, chainId);
    let isValid: boolean = false;
    try {
      isValid = await verifyTypedData({
        ...typedData,
        primaryType: 'Enable',
        address: user.evm_address as Address,
        signature: signature as `0x${string}`,
      } as Parameters<typeof verifyTypedData>[0]);
    } catch (error) {
      // A malformed signature.
    }
    if (!isValid) {
      throw new TurnkeyError(
        `Signature for chain ${chainId} does not match the deposit policy; request it again`,
      );
    }

    const stored: string | undefined = await PermissionApprovalTable.getApprovalForSuborgAndChain(
      user.suborg_id,
      chainId,
    );
    try {
      return await serializePermissionAccount(
        account,
        undefined,
        signature as `0x${string}`,
        getStoredEip7702Authorization(stored),
      );
    } catch (error) {
      logger.error({
        at: 'policy-engine#buildApproval',
        message: 'Failed to assemble the deposit approval',
        suborgId: user.suborg_id,
        chainId,
        error,
      });
      throw new TurnkeyError(`Could not store the approval for chain ${chainId}`);
    }
  }
}
