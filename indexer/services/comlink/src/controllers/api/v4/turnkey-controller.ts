import { randomBytes } from 'crypto';

import { logger, stats } from '@dydxprotocol-indexer/base';
import { TurnkeyUsersTable } from '@dydxprotocol-indexer/postgres';
import { TurnkeyApiClient, TurnkeyApiTypes, Turnkey as TurnkeyServerSDK } from '@turnkey/sdk-server';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Controller, Post, Route, Query,
  Queries,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';

const router: express.Router = express.Router();
const controllerName: string = 'turnkey-controller';

interface TurnkeyAuthResponse {
  organizationId?: string,
  apiKeyId?: string,
  userId?: string,
  session?: string,
  salt: string,
}

interface TurnkeyCreateSuborgResponse {
  subOrgId: string,
  apiKeyId?: string,
  userId?: string,
  salt: string,
}

interface CreateSuborgParams {
  email?: string,
  providerName?: string,
  oidcToken?: string,
  authenticatorName?: string,
  challenge?: string,
  attestation?: TurnkeyApiTypes['v1Attestation'],
}
interface GetSuborgParams {
  email?: string,
  oidcToken?: string,
  credentialId?: string,
}

@Route('turnkey')
export class TurnkeyController extends Controller {
  private parentApiClient: TurnkeyApiClient;
  private bridgeSenderApiClient: TurnkeyApiClient;

  constructor(turnkeyClient?: TurnkeyApiClient, bridgeSenderTurnkeyClient?: TurnkeyApiClient) {
    super();
    logger.info({
      at: 'TurnkeyController#constructor',
      message: 'TurnkeyController constructor',
      params: {
        TURNKEY_API_BASE_URL: config.TURNKEY_API_BASE_URL,
        TURNKEY_API_PRIVATE_KEY: config.TURNKEY_API_PRIVATE_KEY,
        TURNKEY_API_PUBLIC_KEY: config.TURNKEY_API_PUBLIC_KEY,
      },
    });
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

  @Post('/signin')
  async signIn(
    @Query() signinMethod: 'email' | 'social' | 'passkey',
      @Query() userEmail?: string,
      @Query() targetPublicKey?: string,
      @Query() provider?: string,
      @Query() oidcToken?: string,
      @Query() challenge?: string,
      @Queries() attestation?: TurnkeyApiTypes['v1Attestation'],
  ): Promise<TurnkeyAuthResponse> {
    // Determine authentication method
    if (signinMethod === 'email') {
      if (!userEmail || !targetPublicKey) {
        throw new Error('userEmail is required for email signin');
      }
      try {
        logger.info({
          at: 'TurnkeyController#signIn',
          message: 'Email signin',
        });
        const resp = await this.emailSignin(userEmail, targetPublicKey!);
        logger.info({
          at: 'TurnkeyController#signIn',
          message: 'Email auth response',
          params: {
            resp,
          },
        });

        if (resp.userId === undefined || resp.apiKeyId === undefined) {
          throw new Error('Could not send email auth bundle');
        }

        return {
          apiKeyId: resp.apiKeyId,
          userId: resp.userId,
          organizationId: resp.subOrgId,
          salt: resp.salt,
        };

      } catch (error) {
        throw new Error(`Email Signin: ${error}`);
      }
    } else if (signinMethod === 'social') {
      if (!provider || !oidcToken || !targetPublicKey) {
        logger.info({
          at: 'TurnkeyController#signIn',
          message: 'Social signin error',
          params: {
            provider,
            oidcToken,
            targetPublicKey,
          },
        });
        throw new Error('provider, oidcToken, and targetPublicKey are required for social signin');
      }
      try {
        logger.info({
          at: 'TurnkeyController#signIn',
          message: 'Social signin',
        });
        const resp = await this.socialSignin(provider, oidcToken, targetPublicKey);
        logger.info({
          at: 'TurnkeyController#signIn',
          message: 'Social auth response',
          params: {
            resp,
          },
        });
        return {
          session: resp.session,
          salt: resp.salt,
        };
      } catch (error) {
        throw new Error(`Social Signin Error: ${error}`);
      }
    } else if (signinMethod === 'passkey') {
      if (!challenge || !attestation) {
        throw new Error('challenge and attestation are required for passkey signin');
      }
      try {
        const resp = await this.passkeySignin(challenge, 'Passkey', attestation);
        return {
          organizationId: resp.organizationId,
          salt: resp.salt,
        };
      } catch (error) {
        throw new Error(`Passkey Signin Error: ${error}`);
      }
    }
    throw new Error('Invalid signin method. Must be one of: email, social, passkey');
  }

  private getUUID(): string {
    return randomBytes(16).toString('hex');
  }

  // returns the suborgId plus salt if the user exists.
  private async getSuborg(p: GetSuborgParams): Promise<TurnkeyCreateSuborgResponse | undefined> {
    if (p.email) {
      const user = await TurnkeyUsersTable.findByEmail(p.email);
      if (user) {
        // return the suborg id and salt.
        return {
          subOrgId: user.suborgId,
          salt: user.salt,
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
      logger.info({
        at: 'TurnkeyController#getSuborg',
        message: 'Found suborg by credential id',
        params: {
          credentialId: p.credentialId,
          suborgId,
        },
      });
    } else {
      throw new Error('Email is required to create a suborg');
    }

    // find it in our table.
    if (suborgId) {
      const user = await TurnkeyUsersTable.findBySuborgId(suborgId);
      if (user) {
        return {
          subOrgId: suborgId,
          salt: user?.salt || '',
        };
      }
    }
    return undefined;
  }

  // returns the suborgId plus salt and adds the user to the turnkey users table store.
  private async createSuborg(params: CreateSuborgParams): Promise<TurnkeyCreateSuborgResponse> {
    const oauthProviders = [];
    if (params.oidcToken && params.providerName) {
      oauthProviders.push({
        providerName: params.providerName,
        oidcToken: params.oidcToken,
      });
    }

    const authenticators = [];
    if (params.authenticatorName && params.challenge && params.attestation) {

      // serialize the attestation object.
      authenticators.push({
        authenticatorName: params.authenticatorName,
        challenge: params.challenge,
        attestation: params.attestation,
      });

      logger.info({
        at: 'TurnkeyController#createSuborg',
        message: 'Creating suborg with authenticator',
        params: {
          authenticator: params.authenticatorName,
          challenge: params.challenge,
          attestation: params.attestation,
        },
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

    logger.info({
      at: 'TurnkeyController#createSuborg',
      message: 'Created suborg',
      params: {
        subOrg,
      },
    });

    // after creating the orgs, we will need to use the api bridge sender client as
    // parent org api client no longer has permissions to do anything.
    let evmAddress = '';
    let svmAddress = '';
    for (const address of subOrg.wallet?.addresses || []) {
      if (address.startsWith('0x')) {
        // evm always starts with 0x
        evmAddress = address;
      } else {
        // if not evm, then must be svm
        svmAddress = address;
      }
    }

    // generate salt. 256 bit random number
    const salt = this.generateSalt(32);
    // first add to turnkey_users table
    await TurnkeyUsersTable.create({
      suborgId: subOrg.subOrganizationId,
      email: params.email,
      svmAddress,
      evmAddress,
      salt,
      createdAt: new Date().toISOString(),
    });
    // TODO: set the policies on api user

    // Best efforts to check that the subOrg.rootUserIds[0] is the end user
    const user = await this.bridgeSenderApiClient.getUser({
      organizationId: subOrg.subOrganizationId,
      userId: subOrg.rootUserIds?.[0] as string,
    });
    if (
      user.user.authenticators.length > 0 &&
      user.user.authenticators[0].credentialId !== params.attestation?.credentialId
    ) {
      throw new Error('End User not found');
    } else if (user.user.userEmail && user.user.userEmail !== params.email) {
      throw new Error('End User not found');
    }
    // Remove the Delegated Account from the root quorum.
    await this.bridgeSenderApiClient.updateRootQuorum({
      organizationId: subOrg.subOrganizationId,
      threshold: 1,
      userIds: [subOrg.rootUserIds?.[0] as string], // keep end user.
    });
    return {
      subOrgId: subOrg.subOrganizationId,
      salt,
    };
  }

  // email signin creates a suborg if it doesn't already exist.
  private async emailSignin(
    userEmail: string,
    targetPublicKey: string,
  ): Promise<TurnkeyCreateSuborgResponse> {
    // Validate email format
    if (!this.isValidEmail(userEmail)) {
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

    const emailAuthResponse = await this.parentApiClient.emailAuth({
      email: userEmail,
      targetPublicKey,
      emailCustomization: {
        appName: 'dydx',
        logoUrl: 'https://cdn.prod.website-files.com/649ca755d082f1dfc4ed62a4/6870a124cba22652a69c409d_icon%20(1).png',
        magicLinkTemplate: 'https://dydx.trade/login?token=%s',
      },
      invalidateExisting: true,
      organizationId: suborg.subOrgId,
    });
    return {
      subOrgId: suborg.subOrgId,
      apiKeyId: emailAuthResponse.activity.result.emailAuthResult?.apiKeyId,
      userId: emailAuthResponse.activity.result.emailAuthResult?.userId,
      salt: suborg.salt,
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
    };
  }

  // Utility methods
  private isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  private generateChallenge(): string {
    return Buffer.from(crypto.getRandomValues(new Uint8Array(32))).toString('base64url');
  }

  private generateSessionToken(): string {
    return Buffer.from(crypto.getRandomValues(new Uint8Array(32))).toString('base64url');
  }

  private generateSubOrgId(): string {
    return Buffer.from(crypto.getRandomValues(new Uint8Array(16))).toString('hex');
  }

  private generateWalletId(): string {
    return Buffer.from(crypto.getRandomValues(new Uint8Array(16))).toString('hex');
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
    console.log('response', response.organizationIds);

    return response.organizationIds?.[0] || '';
  }

  // this function assumes every user has only one suborg if they have an account with us.
  private async getSuborgByCredentialId(credentialId: string): Promise<string> {
    const response = await this.parentApiClient.getSubOrgIds({
      organizationId: config.TURNKEY_ORGANIZATION_ID,
      filterType: 'CREDENTIAL_ID',
      filterValue: credentialId,
    });
    console.log('response', response.organizationIds);
    return response.organizationIds?.[0] || '';
  }

}

// Validation schemas
const SignInValidationSchema = checkSchema({
  signinMethod: {
    in: ['body'],
    isIn: {
      options: [['social', 'passkey', 'email']],
    },
    errorMessage: 'Must be one of: social, passkey, email',
  },
  userEmail: {
    in: ['body'],
    optional: true,
    isEmail: true,
    errorMessage: 'Must be a valid email address',
  },
  targetPublicKey: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Target public key must be a string',
  },
  // Passkey params
  challenge: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Challenge must be a string',
  },
  attestation: {
    in: ['body'],
    optional: true,
    isObject: true,
    errorMessage: 'Attestation must be a string',
  },
  provider: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Provider must be a string',
  },
  oidcToken: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'OIDC token must be a string',
  },
});

// Express route
router.post(
  '/signin',
  rateLimiterMiddleware(getReqRateLimiter),
  ...SignInValidationSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    try {
      const body = matchedData(req) as {
        signinMethod: 'email',
        userEmail: string,
        targetPublicKey: string,
        provider: string,
        oidcToken: string,
        challenge: string,
        attestation: TurnkeyApiTypes['v1Attestation'],
      };

      const controller: TurnkeyController = new TurnkeyController();
      logger.info({
        at: 'TurnkeyController POST /signin',
        message: 'Signin request',
        params: {
          body,
        },
      });
      const response = await controller.signIn(
        body.signinMethod,
        body.userEmail,
        body.targetPublicKey,
        body.provider,
        body.oidcToken,
        body.challenge,
        body.attestation,
      );

      logger.info({
        at: 'TurnkeyController POST /signin',
        message: 'Signin response',
        params: {
          response,
        },
      });

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

export default router;
