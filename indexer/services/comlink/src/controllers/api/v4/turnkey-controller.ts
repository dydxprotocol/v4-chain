import { stats } from '@dydxprotocol-indexer/base';
import { TurnkeyUsersTable } from '@dydxprotocol-indexer/postgres';
import { TurnkeyApiClient, TurnkeyApiTypes, Turnkey as TurnkeyServerSDK } from '@turnkey/sdk-server';
import express from 'express';
import { matchedData } from 'express-validator';
import fetch from 'node-fetch';
import {
  Controller, Post, Route, Body,
} from 'tsoa';
import { Address, checksumAddress, recoverMessageAddress } from 'viem';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { PolicyEngine } from '../../../helpers/policy-engine';
import { TurnkeyError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { TurnkeyHelpers } from '../../../lib/turnkey-helpers';
import { CheckSignInSchema, CheckUploadDydxAddressSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  SigninMethod,
  TurnkeyAuthResponse,
} from '../../../types';

// Polyfill fetch globally as it's needed by the turnkey sdk.
/* eslint-disable @typescript-eslint/no-explicit-any */
(global as any).fetch = fetch;

export const router: express.Router = express.Router();
const controllerName: string = 'turnkey-controller';

/**
 * Request interface for user sign-in operations
 */
interface SignInRequest {
  /** The authentication method to use (EMAIL, SOCIAL, or PASSKEY) */
  signinMethod: SigninMethod,
  /** User's email address (required for EMAIL signin method) */
  userEmail?: string,
  /** Target public key for authentication (required for EMAIL and SOCIAL signin methods) */
  targetPublicKey?: string,
  /** OAuth provider name (required for SOCIAL signin method) */
  provider?: string,
  /** OIDC token from OAuth provider (required for SOCIAL signin method) */
  oidcToken?: string,
  /** Challenge string for passkey authentication (required for PASSKEY signin method) */
  challenge?: string,
  /** Attestation object for passkey authentication (required for PASSKEY signin method) */
  attestation?: TurnkeyApiTypes['v1Attestation'],
  /** Optional magic link template URL for email authentication */
  magicLink?: string,
}

@Route('turnkey')
export class TurnkeyController extends Controller {
  /** Main Turnkey API client for user authentication and sub-organization management */
  private turnkeyApiClient: TurnkeyApiClient;
  /** Separate Turnkey API client with sender permissions for initiating bridge transactions */
  private bridgeSenderApiClient: TurnkeyApiClient;
  /** Helper class for Turnkey-specific operations */
  private turnkeyHelpers: TurnkeyHelpers;
  /** Policy engine for configuring strict policies */
  private policyEngine: PolicyEngine;

  constructor(turnkeyClient?: TurnkeyApiClient, bridgeSenderTurnkeyClient?: TurnkeyApiClient) {
    super();
    if (turnkeyClient) {
      this.turnkeyApiClient = turnkeyClient;
    } else {
      this.turnkeyApiClient = new TurnkeyServerSDK({
        apiBaseUrl: config.TURNKEY_API_BASE_URL,
        apiPrivateKey: config.TURNKEY_API_PRIVATE_KEY,
        apiPublicKey: config.TURNKEY_API_PUBLIC_KEY,
        defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
      }).apiClient();
    }
    // Bridge sender client uses different API keys with sender permissions
    // to initiate bridge transactions on behalf of users
    if (bridgeSenderTurnkeyClient) {
      this.bridgeSenderApiClient = bridgeSenderTurnkeyClient;
    } else {
      this.bridgeSenderApiClient = new TurnkeyServerSDK({
        apiBaseUrl: config.TURNKEY_API_BASE_URL,
        apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY,
        apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY,
        defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
      }).apiClient();
    }

    // Initialize the Turnkey helpers with the main API client
    this.turnkeyHelpers = new TurnkeyHelpers(this.turnkeyApiClient, this.bridgeSenderApiClient);
    this.policyEngine = new PolicyEngine(this.bridgeSenderApiClient);
  }

  /**
   * Uploads the dydx address to the turnkey table.
   *
   * Backend won't have this information when we create account for user since you need signature
   * to derive dydx address. Just wait for fe to uplaod to kick off the policy setup.
   */
  @Post('/uploadAddress')
  async uploadAddress(
    @Body() body: { dydxAddress: string, signature: string },
  ): Promise<{ success: boolean }> {
    const { dydxAddress, signature } = body;
    if (!dydxAddress || !signature) {
      throw new TurnkeyError('dydxAddress and signature are required');
    }

    // Recover the signer from the signed dydxAddress message
    let recovered: Address;
    try {
      recovered = await recoverMessageAddress({ message: dydxAddress, signature: signature as `0x${string}` });
    } catch (err) {
      throw TurnkeyHelpers.wrapTurnkeyError(err, 'Failed to recover address from signature');
    }

    // Try to find user by the recovered address, falling back to lowercase variant.
    const evmAddressChecksum = checksumAddress(recovered);
    const user = await TurnkeyUsersTable.findByEvmAddress(evmAddressChecksum);
    if (!user) {
      throw new TurnkeyError('No user found for recovered EVM address');
    }

    await TurnkeyUsersTable.updateDydxAddressByEvmAddress(user.evm_address, dydxAddress);

    // configure the policies now
    await this.policyEngine.configureSolanaPolicy(dydxAddress, user.suborg_id);
    await this.policyEngine.configurePolicy(user.suborg_id, user.evm_address, dydxAddress);
    // this removes self from root quorum.
    // Add back once policy configuration is ready.
    // await this.policyEngine.removeSelfFromRootQuorum(user.suborg_id);

    return { success: true };
  }

  @Post('/signin')
  async signIn(
    @Body() body: SignInRequest,
  ): Promise<TurnkeyAuthResponse> {
    const {
      signinMethod,
      userEmail,
      targetPublicKey,
      provider,
      oidcToken,
      challenge,
      attestation,
      magicLink,
    } = body;
    // Determine authentication method
    if (signinMethod === SigninMethod.EMAIL) {
      if (!userEmail || !targetPublicKey) {
        throw new Error('userEmail and targetPublicKey are required for email signin');
      }
      try {
        const resp = await this.turnkeyHelpers.emailSignin(userEmail, targetPublicKey!, magicLink);
        if (resp.userId === undefined || resp.apiKeyId === undefined) {
          throw new Error('Could not send email auth bundle');
        }

        return {
          apiKeyId: resp.apiKeyId,
          userId: resp.userId,
          organizationId: resp.subOrgId,
          salt: resp.salt,
          dydxAddress: resp.dydxAddress || '',
        };

      } catch (error) {
        throw TurnkeyHelpers.wrapTurnkeyError(error, 'Email signin failed');
      }
    } else if (signinMethod === SigninMethod.SOCIAL) {
      if (!provider || !oidcToken || !targetPublicKey) {
        throw new Error('provider, oidcToken, and targetPublicKey are required for social signin');
      }
      try {
        const resp = await this.turnkeyHelpers.socialSignin(provider, oidcToken, targetPublicKey);
        return {
          session: resp.session,
          salt: resp.salt,
          dydxAddress: resp.dydxAddress || '',
        };
      } catch (error) {
        throw TurnkeyHelpers.wrapTurnkeyError(error, 'Social signin failed');
      }
    } else if (signinMethod === SigninMethod.PASSKEY) {
      if (!challenge || !attestation) {
        throw new Error('Passkey signin requires challenge and attestation.');
      }
      try {
        const resp = await this.turnkeyHelpers.passkeySignin(challenge, 'Passkey', attestation);
        return {
          organizationId: resp.organizationId,
          salt: resp.salt,
          dydxAddress: resp.dydxAddress || '',
        };
      } catch (error) {
        throw TurnkeyHelpers.wrapTurnkeyError(error, 'Passkey signin failed');
      }
    }
    throw new Error(`Invalid signin method. Must be one of: ${SigninMethod.EMAIL}, ${SigninMethod.SOCIAL}, ${SigninMethod.PASSKEY}`);
  }
}

// Express route
router.post(
  '/signin',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckSignInSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    try {
      const body = matchedData(req) as {
        signinMethod: SigninMethod,
        userEmail: string,
        targetPublicKey: string,
        magicLink: string,
        provider: string,
        oidcToken: string,
        challenge: string,
        attestation: TurnkeyApiTypes['v1Attestation'],
      };

      const controller: TurnkeyController = new TurnkeyController();

      const response: TurnkeyAuthResponse = await controller.signIn(body);

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TurnkeyController POST /signin',
        'Turnkey sign in error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.post_signin.timing`,
        Date.now() - start,
      );
    }
  },
);

router.post(
  '/uploadAddress',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckUploadDydxAddressSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const body = matchedData(req) as { dydxAddress: string, signature: string };
      const controller: TurnkeyController = new TurnkeyController();
      const response = await controller.uploadAddress(body);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TurnkeyController POST /uploadAddress',
        'Turnkey uploadAddress error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.post_uploadAddress.timing`,
        Date.now() - start,
      );
    }
  },
);
