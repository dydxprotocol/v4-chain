import { randomBytes } from 'crypto';

import { logger, stats } from '@dydxprotocol-indexer/base';
import { WalletTable, FirebaseNotificationTokenTable } from '@dydxprotocol-indexer/postgres';
import * as dotenv from 'dotenv';
import express, { Request } from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Controller, Get, Post, Route, Body, Path, Query,
} from 'tsoa';
import { TurnkeyApiClient, Turnkey as TurnkeyServerSDK } from "@turnkey/sdk-server";
import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import {
  validateSignature,
  validateSignatureKeplr,
} from '../../../helpers/compliance/compliance-utils';
import { NotFoundError } from '../../../lib/errors';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';

const router: express.Router = express.Router();
const controllerName: string = 'turnkey-controller';

interface TurnkeyAuthResponse {
  subOrgId: string,
  wallet: {
    id: string,
    name: string,
    accounts: Array<{
      address: string,
      path: string,
    }>,
  },
  salt: string,
}

@Route('turnkey')
class TurnkeyController extends Controller {
  private apiClient: TurnkeyApiClient;

  constructor() {
    super();
    const turnkeyClient = new TurnkeyServerSDK({
      apiBaseUrl: config.TURNKEY_API_BASE_URL,
      apiPrivateKey: config.TURNKEY_API_SECRET,
      apiPublicKey: config.TURNKEY_ROOT_API_PUBLIC_KEY,
      defaultOrganizationId: config.TURNKEY_ROOT_ORGANIZATION_ID,
    });
    this.apiClient = turnkeyClient.apiClient();
  }

  @Post('/signin')
  async signIn(
    @Query() signinMethod: 'email' | 'social' | 'passkey',
      @Query() userName?: string,
      @Query() challenge?: string,
      @Query() authenticatorName?: string,
      @Query()
      attestation?: {
        credentialId: string,
        clientDataJson: string,
        attestationObject: string,
        transports?: string[],
      },
      @Query() userEmail?: string,
      @Query() provider?: string,
      @Query() oidcToken?: string,
  ): Promise<TurnkeyAuthResponse> {
    // Determine authentication method
    switch (signinMethod) {
      case 'email':
        if (!userEmail) {
          throw new Error('userEmail is required for email signin');
        }
        return this.handleEmailSignIn(userName, userEmail);

      case 'social':
        if (!provider || !oidcToken) {
          throw new Error('provider and oidcToken are required for social signin');
        }
        return this.handleSocialSignIn(userName, provider, oidcToken);

      case 'passkey':
        if (!challenge || !authenticatorName || !attestation) {
          throw new Error(
            'challenge, authenticatorName, and attestation are required for passkey signin',
          );
        }
        return this.handlePasskeySignIn(userName, challenge, authenticatorName, attestation);

      default:
        throw new Error('Invalid signin method. Must be one of: email, social, passkey');
    }
  }

  @Post('/register')
  async register(@Body() body: TurnkeyRegisterRequest): Promise<TurnkeyAuthResponse> {
    const {
      email,
      oauthProvider,
      oauthToken,
      passkeyCredential,
      passkeyChallenge,
      address,
      timestamp,
      message,
      signedMessage,
      pubKey,
      walletIsKeplr = false,
    } = body;

    // Validate timestamp
    const now = Date.now() / 1000;
    if (Math.abs(now - timestamp) > 300) {
      // 5 minutes
      throw new Error('Timestamp is too old or too far in the future');
    }

    // Validate address format
    if (!address || !address.startsWith('dydx1')) {
      throw new Error('Invalid dYdX address format');
    }

    // Check if user already exists
    const existingWallet = await WalletTable.findById(address);
    if (existingWallet) {
      throw new Error('User already exists with this address');
    }

    // Determine authentication method
    if (email) {
      return this.handleEmailRegister(
        email,
        address,
        timestamp,
        message,
        signedMessage,
        pubKey,
        walletIsKeplr,
      );
    } else if (oauthProvider && oauthToken) {
      return this.handleOAuthRegister(
        oauthProvider,
        oauthToken,
        address,
        timestamp,
        message,
        signedMessage,
        pubKey,
        walletIsKeplr,
      );
    } else if (passkeyCredential && passkeyChallenge) {
      return this.handlePasskeyRegister(
        passkeyCredential,
        passkeyChallenge,
        address,
        timestamp,
        message,
      );
    } else {
      throw new Error(
        'Invalid authentication method. Must provide email, OAuth credentials, or passkey credentials.',
      );
    }
  }

  @Get('/user/:address')
  async getUserInfo(@Path() address: string): Promise<TurnkeyUserResponse> {
    const wallet = await WalletTable.findById(address);
    if (!wallet) {
      throw new NotFoundError(`User not found with address ${address}`);
    }

    // Get Firebase tokens to determine auth methods
    const firebaseTokens = await FirebaseNotificationTokenTable.findAll(
      {
        address,
      },
      [],
    );

    return {
      address: wallet.address,
      email: wallet.email, // Assuming email field exists in wallet
      authMethods: ['email', 'oauth', 'passkey'], // This would be stored in the database
      createdAt: wallet.createdAt || new Date().toISOString(),
      lastSignInAt: wallet.updatedAt,
    };
  }

  @Post('/initiate-passkey')
  async initiatePasskey(
    @Body() body: TurnkeyInitiatePasskeyRequest,
  ): Promise<TurnkeyInitiatePasskeyResponse> {
    const { address, action } = body;

    // Generate challenge
    const challenge = this.generateChallenge();

    if (action === 'register') {
      return {
        challenge,
        credentialCreationOptions: {
          challenge,
          rp: {
            name: 'dYdX',
            id: 'dydx.trade',
          },
          user: {
            id: address,
            name: address,
            displayName: address,
          },
          pubKeyCredParams: [
            {
              type: 'public-key',
              alg: -7, // ES256
            },
          ],
          authenticatorSelection: {
            authenticatorAttachment: 'platform',
            userVerification: 'required',
          },
          timeout: 60000,
          attestation: 'direct',
        },
      };
    } else {
      return {
        challenge,
        credentialRequestOptions: {
          challenge,
          timeout: 60000,
          userVerification: 'required',
          rpId: 'dydx.trade',
        },
      };
    }
  }

  // Private helper methods
  private async handleEmailSignIn(
    userName: string,
    userEmail: string,
  ): Promise<TurnkeyAuthResponse> {
    // Validate email format
    if (!this.isValidEmail(userEmail)) {
      throw new Error('Invalid email format');
    }

    // Call Turnkey's v1/submit/email_auth endpoint
    // This is a placeholder - in production, you would call the actual Turnkey API
    const subOrgId = this.generateSubOrgId();
    const walletId = this.generateWalletId();
    const salt = this.generateSalt();

    return {
      subOrgId,
      wallet: {
        id: walletId,
        name: userName,
        accounts: [
          {
            address: 'ethereum_address_placeholder',
            path: "m/44'/60'/0'/0/0",
          },
          {
            address: 'solana_address_placeholder',
            path: "m/44'/501'/0'/0'",
          },
        ],
      },
      salt,
    };
  }

  private async handleSocialSignIn(
    userName: string,
    provider: string,
    oidc_token: string,
  ): Promise<TurnkeyAuthResponse> {
    // Call Turnkey's v1/submit/oauth_login endpoint
    // This is a placeholder - in production, you would call the actual Turnkey API
    const subOrgId = this.generateSubOrgId();
    const walletId = this.generateWalletId();
    const salt = this.generateSalt();

    return {
      subOrgId,
      wallet: {
        id: walletId,
        name: userName,
        accounts: [
          {
            address: 'ethereum_address_placeholder',
            path: "m/44'/60'/0'/0/0",
          },
          {
            address: 'solana_address_placeholder',
            path: "m/44'/501'/0'/0'",
          },
        ],
      },
      salt,
    };
  }

  private async handlePasskeySignIn(
    userName: string,
    challenge: string,
    authenticatorName: string,
    attestation: {
      credentialId: string,
      clientDataJson: string,
      attestationObject: string,
      transports?: string[],
    },
  ): Promise<TurnkeyAuthResponse> {
    // Call Turnkey's v1/submit/stamp_login endpoint using passkey as a stamp
    // This is a placeholder - in production, you would call the actual Turnkey API
    const subOrgId = this.generateSubOrgId();
    const walletId = this.generateWalletId();
    const salt = this.generateSalt();

    return {
      subOrgId,
      wallet: {
        id: walletId,
        name: userName,
        accounts: [
          {
            address: 'ethereum_address_placeholder',
            path: "m/44'/60'/0'/0/0",
          },
          {
            address: 'solana_address_placeholder',
            path: "m/44'/501'/0'/0'",
          },
        ],
      },
      salt,
    };
  }

  private async handleEmailRegister(
    email: string,
    address: string,
    timestamp: number,
    message: string,
    signedMessage?: string,
    pubKey?: string,
  ): Promise<TurnkeyAuthResponse> {
    // Validate email format
    if (!this.isValidEmail(email)) {
      throw new Error('Invalid email format');
    }

    // Validate signature if provided
    if (signedMessage && pubKey) {
      const isValid = await this.validateDydxSignature(
        address,
        timestamp,
        message,
        signedMessage,
        pubKey,
      );

      if (!isValid) {
        throw new Error('Invalid signature');
      }
    }
    try {
      const subOrg = await apiClient.createSubOrganization({
        organizationId: process.env.TURNKEY_ORGANIZATION_ID!,
        subOrganizationName: 'Sub Org - With Root User',
        rootUsers: [
          {
            userName: 'Root User',
            apiKeys: [
              {
                apiKeyName: 'Sender API Key',
                publicKey: process.env.ROOT_API_PUBLIC_KEY!,
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
    } catch (error) {
      // do something
      logger.error({
        at: 'TurnkeyController#handleEmailRegister',
        message: 'Error creating suborg',
        error,
      });
      throw error;
    }

    const subOrgId = suborg;
    const walletId = this.generateWalletId();
    const salt = this.generateSalt();

    return {
      subOrgId,
      wallet: {
        id: walletId,
        name: address,
        accounts: [
          {
            address: 'ethereum_address_placeholder',
            path: "m/44'/60'/0'/0/0",
          },
          {
            address: 'solana_address_placeholder',
            path: "m/44'/501'/0'/0'",
          },
        ],
      },
      salt,
    };
  }

  private async handleOAuthRegister(
    provider: string,
    token: string,
    address: string,
    timestamp: number,
    message: string,
    signedMessage?: string,
    pubKey?: string,
    walletIsKeplr?: boolean,
  ): Promise<TurnkeyAuthResponse> {
    const userInfo = await this.validateOAuthToken(provider, token);

    if (!userInfo) {
      throw new Error('Invalid OAuth token');
    }

    // Validate signature if provided
    if (signedMessage && pubKey) {
      const isValid = walletIsKeplr
        ? this.validateKeplrSignature(address, message, signedMessage, pubKey)
        : await this.validateDydxSignature(address, timestamp, message, signedMessage, pubKey);

      if (!isValid) {
        throw new Error('Invalid signature');
      }
    }

    // Create wallet entry
    await WalletTable.create({
      address,
      totalTradingRewards: '0',
      totalVolume: '0',
    });

    const subOrgId = this.generateSubOrgId();
    const walletId = this.generateWalletId();
    const salt = this.generateSalt();

    return {
      subOrgId,
      wallet: {
        id: walletId,
        name: address,
        accounts: [
          {
            address: 'ethereum_address_placeholder',
            path: "m/44'/60'/0'/0/0",
          },
          {
            address: 'solana_address_placeholder',
            path: "m/44'/501'/0'/0'",
          },
        ],
      },
      salt,
    };
  }

  private async handlePasskeyRegister(
    credential: string,
    challenge: string,
    address: string,
    timestamp: number,
    message: string,
  ): Promise<TurnkeyAuthResponse> {
    // Validate passkey credential
    const isValid = await this.validatePasskeyCredential(credential, challenge);

    if (!isValid) {
      throw new Error('Invalid passkey credential');
    }

    // Create wallet entry
    await WalletTable.create({
      address,
      totalTradingRewards: '0',
      totalVolume: '0',
    });

    const subOrgId = this.generateSubOrgId();
    const walletId = this.generateWalletId();
    const salt = this.generateSalt();

    return {
      subOrgId,
      wallet: {
        id: walletId,
        name: address,
        accounts: [
          {
            address: 'ethereum_address_placeholder',
            path: "m/44'/60'/0'/0/0",
          },
          {
            address: 'solana_address_placeholder',
            path: "m/44'/501'/0'/0'",
          },
        ],
      },
      salt,
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

  private generateSalt(): string {
    return Buffer.from(crypto.getRandomValues(new Uint8Array(16))).toString('hex');
  }

  private async validateOAuthToken(provider: string, token: string): Promise<any> {
    // Implement OAuth token validation for each provider
    // This is a placeholder - in production, you would validate tokens with the respective OAuth providers
    switch (provider) {
      case 'google':
        return this.validateGoogleToken(token);
      case 'github':
        return this.validateGithubToken(token);
      case 'apple':
        return this.validateAppleToken(token);
      default:
        throw new Error('Unsupported OAuth provider');
    }
  }

  private async validateGoogleToken(token: string): Promise<any> {
    // Placeholder for Google token validation
    return { email: 'user@example.com', address: 'dydx1...' };
  }

  private async validateGithubToken(token: string): Promise<any> {
    // Placeholder for GitHub token validation
    return { email: 'user@example.com', address: 'dydx1...' };
  }

  private async validateAppleToken(token: string): Promise<any> {
    // Placeholder for Apple token validation
    return { email: 'user@example.com', address: 'dydx1...' };
  }

  private async validatePasskeySignature(
    credentialId: string,
    signature: string,
    challenge: string,
  ): Promise<boolean> {
    // Placeholder for passkey signature validation
    // In production, this would verify the WebAuthn signature
    return true;
  }

  private async validatePasskeyCredential(credential: string, challenge: string): Promise<boolean> {
    // Placeholder for passkey credential validation
    // In production, this would verify the WebAuthn credential
    return true;
  }

  private validateKeplrSignature(
    address: string,
    message: string,
    signedMessage: string,
    pubKey: string,
  ): boolean {
    // Use existing validation logic
    try {
      const result = validateSignatureKeplr(
        {} as express.Response,
        address,
        message,
        signedMessage,
        pubKey,
      );
      return result === undefined;
    } catch (error) {
      return false;
    }
  }

  private async validateDydxSignature(
    address: string,
    timestamp: number,
    message: string,
    signedMessage: string,
    pubKey: string,
  ): Promise<boolean> {
    // Use existing validation logic
    try {
      const result = await validateSignature(
        {} as express.Response,
        'CONNECT' as any,
        address,
        timestamp,
        message,
        signedMessage,
        pubKey,
      );
      return result === undefined;
    } catch (error) {
      return false;
    }
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
  challenge: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Challenge must be a string',
  },
  authenticatorName: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Authenticator name must be a string',
  },
  'attestation.credentialId': {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Attestation credential ID must be a string',
  },
  'attestation.clientDataJson': {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Attestation client data JSON must be a string',
  },
  'attestation.attestationObject': {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Attestation object must be a string',
  },
  // Email params
  userEmail: {
    in: ['body'],
    optional: true,
    isEmail: true,
    errorMessage: 'Must be a valid email address',
  },
  // Social params
  provider: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Provider must be a string',
  },
  oidc_token: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'OIDC token must be a string',
  },
});

const RegisterValidationSchema = checkSchema({
  email: {
    in: ['body'],
    optional: true,
    isEmail: true,
    errorMessage: 'Must be a valid email address',
  },
  oauthProvider: {
    in: ['body'],
    optional: true,
    isIn: {
      options: [['google', 'github', 'apple']],
    },
    errorMessage: 'Must be one of: google, github, apple',
  },
  oauthToken: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'OAuth token must be a string',
  },
  passkeyCredential: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Passkey credential must be a string',
  },
  passkeyChallenge: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Passkey challenge must be a string',
  },
  address: {
    in: ['body'],
    isString: true,
    custom: {
      options: (value: string) => value.startsWith('dydx1'),
    },
    errorMessage: 'Address must be a valid dYdX address starting with dydx1',
  },
  timestamp: {
    in: ['body'],
    isInt: true,
    errorMessage: 'Timestamp must be an integer',
  },
  message: {
    in: ['body'],
    isString: true,
    errorMessage: 'Message must be a string',
  },
  signedMessage: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Signed message must be a string',
  },
  pubKey: {
    in: ['body'],
    optional: true,
    isString: true,
    errorMessage: 'Public key must be a string',
  },
  walletIsKeplr: {
    in: ['body'],
    optional: true,
    isBoolean: true,
    errorMessage: 'walletIsKeplr must be a boolean',
  },
});

const AddressValidationSchema = checkSchema({
  address: {
    in: ['params'],
    isString: true,
    custom: {
      options: (value: string) => value.startsWith('dydx1'),
    },
    errorMessage: 'Address must be a valid dYdX address starting with dydx1',
  },
});

const InitiatePasskeyValidationSchema = checkSchema({
  address: {
    in: ['body'],
    isString: true,
    custom: {
      options: (value: string) => value.startsWith('dydx1'),
    },
    errorMessage: 'Address must be a valid dYdX address starting with dydx1',
  },
  action: {
    in: ['body'],
    isIn: {
      options: [['register', 'authenticate']],
    },
    errorMessage: 'Action must be either register or authenticate',
  },
});

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
      const body = matchedData(req) as TurnkeySignInRequest;

      const controller: TurnkeyController = new TurnkeyController();
      const response = await controller.signIn(body);

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

router.post(
  '/register',
  rateLimiterMiddleware(getReqRateLimiter),
  ...RegisterValidationSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    try {
      const body = matchedData(req) as TurnkeyRegisterRequest;

      const controller: TurnkeyController = new TurnkeyController();
      const response = await controller.register(body);

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TurnkeyController POST /register',
        'Turnkey register error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.post_register.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/user/:address',
  rateLimiterMiddleware(getReqRateLimiter),
  ...AddressValidationSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    try {
      const { address } = matchedData(req) as { address: string };

      const controller: TurnkeyController = new TurnkeyController();
      const response = await controller.getUserInfo(address);

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TurnkeyController GET /user/:address',
        'Turnkey get user info error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_user_info.timing`,
        Date.now() - start,
      );
    }
  },
);

router.post(
  '/initiate-passkey',
  rateLimiterMiddleware(getReqRateLimiter),
  ...InitiatePasskeyValidationSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    try {
      const body = matchedData(req) as TurnkeyInitiatePasskeyRequest;

      const controller: TurnkeyController = new TurnkeyController();
      const response = await controller.initiatePasskey(body);

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'TurnkeyController POST /initiate-passkey',
        'Turnkey initiate passkey error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.post_initiate_passkey.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
