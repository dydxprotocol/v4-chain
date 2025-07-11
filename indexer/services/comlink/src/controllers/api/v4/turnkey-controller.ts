import { randomBytes } from 'crypto';

import { stats } from '@dydxprotocol-indexer/base';
import { TurnkeyUsersTable } from '@dydxprotocol-indexer/postgres';
import { TurnkeyApiClient, Turnkey as TurnkeyServerSDK } from '@turnkey/sdk-server';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Controller, Post, Route, Query,
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
  subOrgId: string,
  salt: string,
}

interface TurnkeyCreateSuborgResponse {
  subOrgId: string,
  salt: string,
}

@Route('turnkey')
class TurnkeyController extends Controller {
  private parentApiClient: TurnkeyApiClient;
  private bridgeSenderApiClient: TurnkeyApiClient;

  constructor() {
    super();
    const turnkeyClient = new TurnkeyServerSDK({
      apiBaseUrl: config.TURNKEY_API_BASE_URL,
      apiPrivateKey: config.TURNKEY_API_PRIVATE_KEY,
      apiPublicKey: config.TURNKEY_API_PUBLIC_KEY,
      defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
    });
    this.parentApiClient = turnkeyClient.apiClient();

    const bridgeSenderTurnkeyClient = new TurnkeyServerSDK({
      apiBaseUrl: config.TURNKEY_API_BASE_URL,
      apiPrivateKey: config.TURNKEY_API_SENDER_PRIVATE_KEY,
      apiPublicKey: config.TURNKEY_API_SENDER_PUBLIC_KEY,
      defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
    });
    this.bridgeSenderApiClient = bridgeSenderTurnkeyClient.apiClient();
  }

  @Post('/signin')
  async signIn(
    @Query() signinMethod: 'email' | 'social' | 'passkey',
      @Query() userEmail?: string,
      @Query() targetPublicKey?: string,
  ): Promise<TurnkeyAuthResponse> {
    // Determine authentication method
    if (signinMethod === 'email') {
      if (!userEmail || !targetPublicKey) {
        throw new Error('userEmail is required for email signin');
      }
      try {
        const suborg = await this.emailSignin(userEmail, targetPublicKey!);

        const emailAuthResponse = await this.parentApiClient.emailAuth({
          email: userEmail,
          targetPublicKey: config.TURNKEY_API_PUBLIC_KEY,
          emailCustomization: {
            appName: 'dydx',
            logoUrl: 'https://dydx.trade/currencies/dydx.png',
            magicLinkTemplate: 'https://dydx.trade/login?token=%s',
          },
          invalidateExisting: true,
          organizationId: suborg.subOrgId,
        });

        if (emailAuthResponse.activity.result.emailAuthResult?.userId !== undefined) {
          throw new Error('Could not send email auth bundle');
        }

        return {
          subOrgId: suborg.subOrgId,
          salt: suborg.salt,
        };

      } catch (error) {
        throw new Error('Error signing in');
      }
    }
    throw new Error('Invalid signin method. Must be one of: email, social, passkey');
  }

  private getUUID(): string {
    return randomBytes(16).toString('hex');
  }

  private async createSuborg(email: string): Promise<TurnkeyCreateSuborgResponse> {
    try {

      const subOrg = await this.parentApiClient.createSubOrganization({
        subOrganizationName: this.getUUID(),
        rootUsers: [
          {
            userName: 'End User',
            userEmail: email,
            apiKeys: [],
            authenticators: [],
            oauthProviders: [],
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
        email,
        svmAddress,
        evmAddress,
        salt,
        createdAt: new Date().toISOString(),
      });
      // TODO: set the policies on api user

      // Check that the subOrg.rootUserIds[0] is the end user
      const user = await this.bridgeSenderApiClient.getUser({
        organizationId: subOrg.subOrganizationId,
        userId: subOrg.rootUserIds?.[0] as string,
      });
      if (user.user.userEmail !== email) {
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
    } catch (error) {
      throw new Error('Error signing in');
    }
  }

  // Private helper methods
  private async emailSignin(
    userEmail: string,
    targetPublicKey: string,
  ): Promise<TurnkeyCreateSuborgResponse> {
    // Validate email format
    if (!this.isValidEmail(userEmail)) {
      throw new Error('Invalid email format');
    }
    try {
      let suborg: TurnkeyCreateSuborgResponse;
      // search user by email in table to see if they're already a user.
      const user = await TurnkeyUsersTable.findByEmail(userEmail);
      if (user) {
        // return the suborg id and salt.
        suborg = {
          subOrgId: user.suborgId,
          salt: user.salt,
        };
      } else {
        // if user does not exist, create suborg with email set as root user
        suborg = await this.createSuborg(userEmail);
      }

      this.parentApiClient.emailAuth({
        email: userEmail,
        targetPublicKey,
        invalidateExisting: true,
        organizationId: suborg.subOrgId,
      });
      return suborg;
    } catch (error) {
      throw new Error('Error signing in');
    }
  }

  // private async handleSocialSignIn(
  //   userName: string,
  //   provider: string,
  //   oidc_token: string,
  // ): Promise<TurnkeyAuthResponse> {
  //   // Call Turnkey's v1/submit/oauth_login endpoint
  //   // This is a placeholder - in production, you would call the actual Turnkey API
  //   const subOrgId = this.generateSubOrgId();
  //   const walletId = this.generateWalletId();
  //   const salt = this.generateSalt();

  //   return {
  //     subOrgId,
  //     wallet: {
  //       id: walletId,
  //       name: userName,
  //       accounts: [
  //         {
  //           address: 'ethereum_address_placeholder',
  //           path: "m/44'/60'/0'/0/0",
  //         },
  //         {
  //           address: 'solana_address_placeholder',
  //           path: "m/44'/501'/0'/0'",
  //         },
  //       ],
  //     },
  //     salt,
  //   };
  // }

  // private async handlePasskeySignIn(
  //   userName: string,
  //   challenge: string,
  //   authenticatorName: string,
  //   attestation: {
  //     credentialId: string,
  //     clientDataJson: string,
  //     attestationObject: string,
  //     transports?: string[],
  //   },
  // ): Promise<TurnkeyAuthResponse> {
  //   // Call Turnkey's v1/submit/stamp_login endpoint using passkey as a stamp
  //   // This is a placeholder - in production, you would call the actual Turnkey API
  //   const subOrgId = this.generateSubOrgId();
  //   const walletId = this.generateWalletId();
  //   const salt = this.generateSalt();

  //   return {
  //     subOrgId,
  //     wallet: {
  //       id: walletId,
  //       name: userName,
  //       accounts: [
  //         {
  //           address: 'ethereum_address_placeholder',
  //           path: "m/44'/60'/0'/0/0",
  //         },
  //         {
  //           address: 'solana_address_placeholder',
  //           path: "m/44'/501'/0'/0'",
  //         },
  //       ],
  //     },
  //     salt,
  //   };
  // }

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
  userName: {
    in: ['body'],
    isString: true,
    errorMessage: 'userName must be a string',
  },
  // Passkey params
  // challenge: {
  //   in: ['body'],
  //   optional: true,
  //   isString: true,
  //   errorMessage: 'Challenge must be a string',
  // },
  // authenticatorName: {
  //   in: ['body'],
  //   optional: true,
  //   isString: true,
  //   errorMessage: 'Authenticator name must be a string',
  // },
  // 'attestation.credentialId': {
  //   in: ['body'],
  //   optional: true,
  //   isString: true,
  //   errorMessage: 'Attestation credential ID must be a string',
  // },
  // 'attestation.clientDataJson': {
  //   in: ['body'],
  //   optional: true,
  //   isString: true,
  //   errorMessage: 'Attestation client data JSON must be a string',
  // },
  // 'attestation.attestationObject': {
  //   in: ['body'],
  //   optional: true,
  //   isString: true,
  //   errorMessage: 'Attestation object must be a string',
  // },
  // Email params
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
  // Social params
  // provider: {
  //   in: ['body'],
  //   optional: true,
  //   isString: true,
  //   errorMessage: 'Provider must be a string',
  // },
  // oidc_token: {
  //   in: ['body'],
  //   optional: true,
  //   isString: true,
  //   errorMessage: 'OIDC token must be a string',
  // },
});

// const RegisterValidationSchema = checkSchema({
//   email: {
//     in: ['body'],
//     optional: true,
//     isEmail: true,
//     errorMessage: 'Must be a valid email address',
//   },
//   oauthProvider: {
//     in: ['body'],
//     optional: true,
//     isIn: {
//       options: [['google', 'github', 'apple']],
//     },
//     errorMessage: 'Must be one of: google, github, apple',
//   },
//   oauthToken: {
//     in: ['body'],
//     optional: true,
//     isString: true,
//     errorMessage: 'OAuth token must be a string',
//   },
//   passkeyCredential: {
//     in: ['body'],
//     optional: true,
//     isString: true,
//     errorMessage: 'Passkey credential must be a string',
//   },
//   passkeyChallenge: {
//     in: ['body'],
//     optional: true,
//     isString: true,
//     errorMessage: 'Passkey challenge must be a string',
//   },
//   address: {
//     in: ['body'],
//     isString: true,
//     custom: {
//       options: (value: string) => value.startsWith('dydx1'),
//     },
//     errorMessage: 'Address must be a valid dYdX address starting with dydx1',
//   },
//   timestamp: {
//     in: ['body'],
//     isInt: true,
//     errorMessage: 'Timestamp must be an integer',
//   },
//   message: {
//     in: ['body'],
//     isString: true,
//     errorMessage: 'Message must be a string',
//   },
//   signedMessage: {
//     in: ['body'],
//     optional: true,
//     isString: true,
//     errorMessage: 'Signed message must be a string',
//   },
//   pubKey: {
//     in: ['body'],
//     optional: true,
//     isString: true,
//     errorMessage: 'Public key must be a string',
//   },
//   walletIsKeplr: {
//     in: ['body'],
//     optional: true,
//     isBoolean: true,
//     errorMessage: 'walletIsKeplr must be a boolean',
//   },
// });

// const AddressValidationSchema = checkSchema({
//   address: {
//     in: ['params'],
//     isString: true,
//     custom: {
//       options: (value: string) => value.startsWith('dydx1'),
//     },
//     errorMessage: 'Address must be a valid dYdX address starting with dydx1',
//   },
// });

// const InitiatePasskeyValidationSchema = checkSchema({
//   address: {
//     in: ['body'],
//     isString: true,
//     custom: {
//       options: (value: string) => value.startsWith('dydx1'),
//     },
//     errorMessage: 'Address must be a valid dYdX address starting with dydx1',
//   },
//   action: {
//     in: ['body'],
//     isIn: {
//       options: [['register', 'authenticate']],
//     },
//     errorMessage: 'Action must be either register or authenticate',
//   },
// });

// Express routes
router.post(
  '/turnkeySignin',
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
      };

      const controller: TurnkeyController = new TurnkeyController();
      const response = await controller.signIn(
        body.signinMethod,
        body.userEmail,
        body.targetPublicKey,
      );

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TurnkeyController POST /turnkeySignin',
        'Turnkey sign in error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.post_turnkey_signin.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
