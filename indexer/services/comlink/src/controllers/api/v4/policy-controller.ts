import { logger } from '@dydxprotocol-indexer/base';
import { TurnkeyApiClient, Turnkey as TurnkeyServerSDK } from '@turnkey/sdk-server';
import { createAccount } from '@turnkey/viem';
import { signerToEcdsaValidator } from '@zerodev/ecdsa-validator';
import { serializePermissionAccount, toPermissionValidator } from '@zerodev/permissions';
import { toCallPolicy, CallPolicyVersion } from '@zerodev/permissions/policies';
import { toECDSASigner } from '@zerodev/permissions/signers';
import { addressToEmptyAccount, createKernelAccount } from '@zerodev/sdk';
import { KERNEL_V3_1, KERNEL_V3_3 } from '@zerodev/sdk/constants';
import { Controller } from 'tsoa';
import { avalanche } from 'viem/chains';

import config from '../../../config';
import { abi } from '../../../helpers/abi';
import { publicClients } from '../../../helpers/alchemy-helpers';
import { suborgToApproval } from '../../../helpers/skip-helper';
import { entryPoint } from '../../../lib/smart-contract-constants';

export class PolicyController extends Controller {
  private turnkeySenderClient: TurnkeyApiClient;

  constructor(turnkeySenderClient?: TurnkeyApiClient) {
    super();
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
    chainId: string,
    fromAddress: string,
  ) {

    const turnkeyAccount = await createAccount({
      // @ts-ignore
      client: this.turnkeySenderClient,
      organizationId: suborgId,
      signWith: fromAddress,
    });

    let kernelVersion = KERNEL_V3_3;
    if (chainId === avalanche.id.toString()) {
      kernelVersion = KERNEL_V3_1;
      // use this for avalanche to create the ecdsa validator.
      await signerToEcdsaValidator(publicClients[chainId], {
        entryPoint,
        kernelVersion,
        signer: turnkeyAccount,
      });
    }

    // Create an "empty account" as the signer -- you only need the public
    // key (address) to do this.
    const emptyAccount = addressToEmptyAccount(config.MASTER_SIGNER_PUBLIC as `0x${string}`);
    const emptySessionKeySigner = await toECDSASigner({ signer: emptyAccount });

    try {
      const permissionPlugin = await toPermissionValidator(publicClients[chainId], {
        entryPoint,
        kernelVersion,
        signer: emptySessionKeySigner,
        policies: [
          toCallPolicy({
            policyVersion: CallPolicyVersion.V0_0_4,
            permissions: [
              {
                target: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831',
                abi,
                valueLimit: BigInt(1000000000000000000),
                functionName: 'approve',
              },
              {
                target: '0x4c58aE019E54D10594F1Aa26ABF385B6fb17A52d',
                abi,
                valueLimit: BigInt(1000000000000000000),
                functionName: 'submitOrder',
              },
              {
                target: '0x4c58aE019E54D10594F1Aa26ABF385B6fb17A52d',
                abi,
                valueLimit: BigInt(1000000000000000000),
                functionName: 'swapAndSubmitOrder',
              },
            ],
          }),
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
      const approval = await serializePermissionAccount(sessionAccount);
      logger.info({
        at: 'policy-controller#configurePolicy',
        message: 'Approval obtained',
        approval,
        suborgId,
      });
      suborgToApproval.set(suborgId, approval);
    } catch (error) {
      logger.error({ at: 'policy-controller#configurePolicy', message: 'Error configuring policy', error });
      throw error;
    }
  }
}
