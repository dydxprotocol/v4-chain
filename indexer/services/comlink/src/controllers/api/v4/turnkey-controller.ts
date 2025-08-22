import { randomBytes } from 'crypto';

import { logger, stats } from '@dydxprotocol-indexer/base';
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
import { addAddressesToAlchemyWebhook, getSmartAccountAddress } from '../../../helpers/alchemy-helpers';
import { isValidEmail } from '../../../helpers/utility/validation';
import { TurnkeyError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckSignInSchema, CheckUploadDydxAddressSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import {
  SigninMethod,
  TurnkeyAuthResponse,
  TurnkeyCreateSuborgResponse,
  CreateSuborgParams,
  GetSuborgParams,
} from '../../../types';

// Polyfill fetch globally as it's needed by the turnkey sdk.
/* eslint-disable @typescript-eslint/no-explicit-any */
(global as any).fetch = fetch;

export const router: express.Router = express.Router();
const controllerName: string = 'turnkey-controller';

interface SignInRequest {
  signinMethod: SigninMethod,
  userEmail?: string,
  targetPublicKey?: string,
  provider?: string,
  oidcToken?: string,
  challenge?: string,
  attestation?: TurnkeyApiTypes['v1Attestation'],
  magicLink?: string,
}

@Route('turnkey')
export class TurnkeyController extends Controller {
  private parentApiClient: TurnkeyApiClient;
  private bridgeSenderApiClient: TurnkeyApiClient;

  constructor(turnkeyClient?: TurnkeyApiClient, bridgeSenderTurnkeyClient?: TurnkeyApiClient) {
    super();
    if (turnkeyClient) {
      this.parentApiClient = turnkeyClient;
    } else {
      this.parentApiClient = new TurnkeyServerSDK({
        apiBaseUrl: config.TURNKEY_API_BASE_URL,
        apiPrivateKey: config.TURNKEY_API_PRIVATE_KEY,
        apiPublicKey: config.TURNKEY_API_PUBLIC_KEY,
        defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
      }).apiClient();
    }
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
  }

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
      throw this.wrapTurnkeyError(err, 'Failed to recover address from signature');
    }

    // Try to find user by the recovered address, falling back to lowercase variant
    const evmAddressChecksum = checksumAddress(recovered);
    const user = await TurnkeyUsersTable.findByEvmAddress(evmAddressChecksum);
    if (!user) {
      throw new TurnkeyError('No user found for recovered EVM address');
    }

    await TurnkeyUsersTable.updateDydxAddressByEvmAddress(user.evm_address, dydxAddress);

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
        throw new Error('userEmail is required for email signin');
      }
      try {
        const resp = await this.emailSignin(userEmail, targetPublicKey!, magicLink);
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
        throw this.wrapTurnkeyError(error, 'Email signin failed');
      }
    } else if (signinMethod === SigninMethod.SOCIAL) {
      if (!provider || !oidcToken || !targetPublicKey) {
        throw new Error('provider, oidcToken, and targetPublicKey are required for social signin');
      }
      try {
        const resp = await this.socialSignin(provider, oidcToken, targetPublicKey);
        return {
          session: resp.session,
          salt: resp.salt,
          dydxAddress: resp.dydxAddress || '',
        };
      } catch (error) {
        throw this.wrapTurnkeyError(error, 'Social signin failed');
      }
    } else if (signinMethod === SigninMethod.PASSKEY) {
      if (!challenge || !attestation) {
        throw new Error('challenge and attestation are required for passkey signin');
      }
      try {
        const resp = await this.passkeySignin(challenge, 'Passkey', attestation);
        return {
          organizationId: resp.organizationId,
          salt: resp.salt,
          dydxAddress: resp.dydxAddress || '',
        };
      } catch (error) {
        throw this.wrapTurnkeyError(error, 'Passkey signin failed');
      }
    }
    throw new Error(`Invalid signin method. Must be one of: ${SigninMethod.EMAIL}, ${SigninMethod.SOCIAL}, ${SigninMethod.PASSKEY}`);
  }

  private getUUID(): string {
    return randomBytes(16).toString('hex');
  }

  /*
   * Returns the suborgId plus salt if the user exists.
   * Additionally will include the dydxAddress if the user has one uploaded already.
   *
   */
  private async getSuborg(p: GetSuborgParams): Promise<TurnkeyCreateSuborgResponse | undefined> {
    if (p.email) {
      const user = await TurnkeyUsersTable.findByEmail(p.email);
      if (user) {
        // return the suborg id and salt.
        return {
          subOrgId: user.suborg_id,
          salt: user.salt,
          dydxAddress: user.dydx_address || '',
        };
      }
      return undefined;
    }

    // if we don't have an email, we need to find the suborg id by oidc token or credential id.
    let suborgId: string;
    if (p.oidcToken) {
      suborgId = await this.getSuborgByOIDCToken(p.oidcToken);
    } else if (p.credentialId) {
      suborgId = await this.getSuborgByCredentialId(p.credentialId);
    } else if (p.email) {
      suborgId = await this.getSuborgByEmail(p.email);
    } else {
      throw new Error('One of email, oidcToken, or credentialId is required');
    }

    // find it in our table.
    if (suborgId) {
      const user = await TurnkeyUsersTable.findBySuborgId(suborgId);
      if (user) {
        return {
          subOrgId: user?.suborg_id || '',
          salt: user?.salt || '',
          dydxAddress: user?.dydx_address || '',
        };
      }
    }
    return undefined;
  }

  // returns the suborgId plus salt and adds the user to the turnkey users table store.
  private async createSuborg(params: CreateSuborgParams): Promise<TurnkeyCreateSuborgResponse> {
    const oauthProviders: TurnkeyApiTypes['v1OauthProviderParams'][] = [];
    if (params.oidcToken && params.providerName) {
      oauthProviders.push({
        providerName: params.providerName,
        oidcToken: params.oidcToken,
      });
    }

    const authenticators: TurnkeyApiTypes['v1AuthenticatorParamsV2'][] = [];
    if (params.authenticatorName && params.challenge && params.attestation) {

      // serialize the attestation object.
      authenticators.push({
        authenticatorName: params.authenticatorName,
        challenge: params.challenge,
        attestation: params.attestation,
      });

    }
    const subOrg = await this.parentApiClient.createSubOrganization({
      subOrganizationName: this.getUUID(),
      rootUsers: [
        {
          userName: 'End User',
          userEmail: params.email,
          apiKeys: [],
          authenticators,
          oauthProviders,
        },
        {
          userName: 'API User',
          apiKeys: [
            {
              apiKeyName: 'Bridge API Key',
              publicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY,
              curveType: 'API_KEY_CURVE_P256',
            },
          ],
          authenticators: [],
          oauthProviders: [],
        },
      ],
      rootQuorumThreshold: 1,
      wallet: {
        walletName: 'Default ETH Wallet',
        accounts: [
          {
            curve: 'CURVE_SECP256K1',
            pathFormat: 'PATH_FORMAT_BIP32',
            path: "m/44'/60'/0'/0/0",
            addressFormat: 'ADDRESS_FORMAT_ETHEREUM',
          },
          {
            curve: 'CURVE_ED25519',
            pathFormat: 'PATH_FORMAT_BIP32',
            path: "m/44'/501'/0'/0'", // Standard Solana derivation path
            addressFormat: 'ADDRESS_FORMAT_SOLANA',
          },
        ],
      },
    });

    // after creating the orgs, we will need to use the api bridge sender client as
    // parent org api client no longer has permissions to do anything.
    let evmAddress = '';
    let svmAddress = '';
    // smart account address can be derived offchain before we send any user ops.
    // smart account address is needed by the frontend to display the correct
    // deposit address for the avalanche chain since it does not support eip7702.
    let smartAccountAddress = '';
    for (const address of subOrg.wallet?.addresses || []) {
      if (address.startsWith('0x')) {
        // evm always starts with 0x
        evmAddress = address;
        smartAccountAddress = await getSmartAccountAddress(evmAddress);
        smartAccountAddress = checksumAddress(smartAccountAddress as Address);
      } else {
        // if not evm, then must be svm
        svmAddress = address;
      }
    }

    // generate salt. 256 bit random number
    const salt = this.generateSalt(32);
    // first add to turnkey_users table
    await TurnkeyUsersTable.create({
      suborg_id: subOrg.subOrganizationId,
      email: params.email,
      svm_address: svmAddress,
      evm_address: evmAddress,
      smart_account_address: smartAccountAddress,
      salt,
      created_at: new Date().toISOString(),
    });

    // need to also add the svm and evm addresses to the alchemy hook
    if (evmAddress && svmAddress) {
      // We don't need to wait for it since
      // frontend doesn't really neeed the results???
      addAddressesToAlchemyWebhook(evmAddress, svmAddress).catch((error) => {
        logger.error({
          message: 'Failed to add addresses to alchemy webhook',
          error,
          at: new Date().toISOString(),
          evmAddress,
          svmAddress,
        });
      });
    }
    return {
      subOrgId: subOrg.subOrganizationId,
      salt,
    };
  }

  // email signin creates a suborg if it doesn't already exist.
  private async emailSignin(
    userEmail: string,
    targetPublicKey: string,
    magicLink?: string,
  ): Promise<TurnkeyCreateSuborgResponse> {
    // Validate email format
    if (!isValidEmail(userEmail)) {
      throw new Error('Invalid email format');
    }
    let suborg: TurnkeyCreateSuborgResponse | undefined = await this.getSuborg({
      email: userEmail,
    });
    if (!suborg) {
      suborg = await this.createSuborg({
        email: userEmail,
      });
    }

    // Validate magic link template if provided
    const magicLinkTemplate = config.TURNKEY_MAGIC_LINK_TEMPLATE || magicLink;
    if (magicLinkTemplate) {
      try {
        // eslint-disable-next-line no-new
        new URL(magicLinkTemplate.replace('%s', 'test'));
      } catch {
        throw new Error('Invalid magic link template URL');
      }
    }
    const emailAuthResponse = await this.parentApiClient.emailAuth({
      email: userEmail,
      targetPublicKey,
      emailCustomization: {
        appName: 'dydx',
        logoUrl: 'https://cdn.prod.website-files.com/649ca755d082f1dfc4ed62a4/6870a124cba22652a69c409d_icon%20(1).png',
        magicLinkTemplate: magicLinkTemplate ? `${magicLinkTemplate}=%s` : undefined,
      },
      invalidateExisting: true,
      organizationId: suborg.subOrgId,
    });

    return {
      subOrgId: suborg.subOrgId,
      apiKeyId: emailAuthResponse.activity.result.emailAuthResult?.apiKeyId,
      userId: emailAuthResponse.activity.result.emailAuthResult?.userId,
      salt: suborg.salt,
      dydxAddress: suborg.dydxAddress || '',
    };
  }

  // creates a suborg if one doesn't already exist, then login with the oidc token.
  private async socialSignin(
    provider: string,
    oidcToken: string,
    targetPublicKey: string,
  ): Promise<TurnkeyAuthResponse> {
    let suborg: TurnkeyCreateSuborgResponse | undefined = await this.getSuborg({
      oidcToken,
    });

    if (!suborg) {
      suborg = await this.createSuborg({
        providerName: provider,
        oidcToken,
      });
    }

    const oauthLoginResponse = await this.parentApiClient.oauthLogin({
      oidcToken,
      publicKey: targetPublicKey,
      invalidateExisting: true,
      organizationId: suborg.subOrgId,
    });
    return {
      session: oauthLoginResponse.activity.result.oauthLoginResult?.session,
      salt: suborg.salt,
      dydxAddress: suborg.dydxAddress || '',
    };
  }

  // does not return a session as front end can just call the stampLogin endpoint.
  // front end should just call the stampLogin endpoint and use this signin method
  // as a way to get the salt.
  private async passkeySignin(
    challenge: string,
    authenticatorName: string,
    attestation: TurnkeyApiTypes['v1Attestation'],
  ): Promise<TurnkeyAuthResponse> {
    let suborg: TurnkeyCreateSuborgResponse | undefined = await this.getSuborg({
      credentialId: attestation.credentialId,
    });
    if (!suborg) {
      suborg = await this.createSuborg({
        authenticatorName,
        challenge,
        attestation,
      });
    }

    return {
      organizationId: suborg.subOrgId,
      salt: suborg.salt,
      dydxAddress: suborg.dydxAddress || '',
    };
  }

  // default 32 bytes.
  private generateSalt(bytes: number = 32): string {
    return randomBytes(bytes).toString('hex');
  }

  // this function assumes every user has only one suborg if they have an account with us.
  private async getSuborgByOIDCToken(oidcToken: string): Promise<string> {
    const response = await this.parentApiClient.getSubOrgIds({
      organizationId: config.TURNKEY_ORGANIZATION_ID,
      filterType: 'OIDC_TOKEN',
      filterValue: oidcToken,
    });

    return response.organizationIds?.[0] || '';
  }

  // this function assumes every user has only one suborg if they have an account with us.
  private async getSuborgByCredentialId(credentialId: string): Promise<string> {
    const response = await this.parentApiClient.getSubOrgIds({
      organizationId: config.TURNKEY_ORGANIZATION_ID,
      filterType: 'CREDENTIAL_ID',
      filterValue: credentialId,
    });
    return response.organizationIds?.[0] || '';
  }

  // this function assumes every user has only one suborg if they have an account with us.
  private async getSuborgByEmail(email: string): Promise<string> {
    const response = await this.parentApiClient.getSubOrgIds({
      organizationId: config.TURNKEY_ORGANIZATION_ID,
      filterType: 'EMAIL',
      filterValue: email,
    });
    return response.organizationIds?.[0] || '';
  }

  // Helper method to wrap Turnkey errors with additional context
  private wrapTurnkeyError(error: unknown, contextMessage: string): TurnkeyError {
    if (error instanceof Error) {
      return new TurnkeyError(
        `${contextMessage}: ${error.message}`,
      );
    }
    return new TurnkeyError(`${contextMessage}: ${String(error)}`);
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
