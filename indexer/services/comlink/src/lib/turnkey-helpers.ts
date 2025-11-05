import { randomBytes, randomUUID } from 'crypto';

import { logger } from '@dydxprotocol-indexer/base';
import { TurnkeyUsersTable } from '@dydxprotocol-indexer/postgres';
import { TurnkeyApiClient, TurnkeyApiTypes } from '@turnkey/sdk-server';
import { decodeJwt } from 'jose';
import { Address, checksumAddress } from 'viem';

import config from '../config';
import { TURNKEY_EMAIL_CUSTOMIZATION } from '../constants';
import { getSmartAccountAddress } from '../helpers/alchemy-helpers';
import {
  CreateSuborgParams,
  GetSuborgParams,
  TurnkeyCreateSuborgResponse,
  TurnkeyAuthResponse,
} from '../types';
import { TurnkeyError } from './errors';

/**
 * Helper class for Turnkey-specific operations
 */
export class TurnkeyHelpers {
  private turnkeyApiClient: TurnkeyApiClient;

  constructor(turnkeyApiClient: TurnkeyApiClient) {
    this.turnkeyApiClient = turnkeyApiClient;
  }

  /**
   * Generates a random identifier string.
   * Uses crypto.randomUUID() for UUID generation or crypto.randomBytes() for custom-length
   * hex strings.
   *
   * @param type - Type of random string to generate ('uuid' or 'hex')
   * @param bytes - Number of bytes for hex generation (default: 32)
   * @returns A random string
   */
  static generateRandomString(type: 'uuid' | 'hex' = 'uuid', bytes: number = 32): string {
    if (type === 'uuid') {
      return randomUUID();
    }
    return randomBytes(bytes).toString('hex');
  }

  /**
   * Returns the suborgId plus salt if the user exists.
   * Additionally will include the dydxAddress if the user has one uploaded already.
   */
  async getSuborg(p: GetSuborgParams): Promise<TurnkeyCreateSuborgResponse | undefined> {
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

  /**
   * Creates a new Turnkey sub-organization and adds the user to the database.
   * Returns the suborgId plus salt.
   *
   * This sets up the suborg within turnkey for the user. Then upon the user
   * signs in, fe will update the dydx address in the database, then delete the
   * dydx user from the user's wallet. Up until that point, there will not be
   * any funds in the user wallet.
   */
  async createSuborg(params: CreateSuborgParams): Promise<TurnkeyCreateSuborgResponse> {
    // v1OauthProviderParams and v1AuthenticatorParamsV2 are types from @turnkey/sdk-server
    // that define the structure for OAuth providers and authenticator configurations
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

    const subOrg = await this.turnkeyApiClient.createSubOrganization({
      subOrganizationName: TurnkeyHelpers.generateRandomString('uuid'),
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
        walletName: 'User Wallet',
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
    try {
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
    } catch (e) {
      throw new Error(`Failed to derive smart account address: ${e}`);
    }

    // generate salt. 256 bit random number
    const salt = TurnkeyHelpers.generateRandomString('hex', 32);
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

    return {
      subOrgId: subOrg.subOrganizationId,
      salt,
    };
  }

  /**
   * Retrieves the sub-organization ID for a user based on their OIDC token.
   * Makes a remote API call to Turnkey to find the sub-organization.
   * This function assumes every user has only one suborg if they have an account with us.
   *
   * @param oidcToken - The OIDC token from the OAuth provider
   * @returns The sub-organization ID or empty string if not found
   */
  private async getSuborgByOIDCToken(oidcToken: string): Promise<string> {
    const response = await this.turnkeyApiClient.getSubOrgIds({
      organizationId: config.TURNKEY_ORGANIZATION_ID,
      filterType: 'OIDC_TOKEN',
      filterValue: oidcToken,
    });

    return response.organizationIds?.[0] || '';
  }

  /**
   * Retrieves the sub-organization ID for a user based on their credential ID.
   * Makes a remote API call to Turnkey to find the sub-organization.
   * This function assumes every user has only one suborg if they have an account with us.
   *
   * @param credentialId - The credential ID from passkey authentication
   * @returns The sub-organization ID or empty string if not found
   */
  private async getSuborgByCredentialId(credentialId: string): Promise<string> {
    const response = await this.turnkeyApiClient.getSubOrgIds({
      organizationId: config.TURNKEY_ORGANIZATION_ID,
      filterType: 'CREDENTIAL_ID',
      filterValue: credentialId,
    });
    return response.organizationIds?.[0] || '';
  }

  /**
   * Helper method to wrap Turnkey errors with additional context.
   *
   * Remote call failures to Turnkey are handled by wrapping the original error
   * with contextual information about what operation was being performed.
   * This provides better error messages for debugging and user feedback.
   *
   * @param error - The original error from Turnkey API
   * @param contextMessage - Additional context about the operation that failed
   * @returns A wrapped TurnkeyError with enhanced context
   */
  static wrapTurnkeyError(error: unknown, contextMessage: string): TurnkeyError {
    if (error instanceof Error) {
      return new TurnkeyError(
        `${contextMessage}: ${error.message}`,
      );
    }
    return new TurnkeyError(`${contextMessage}: ${String(error)}`);
  }

  /**
   * Performs email authentication with Turnkey.
   * Creates a suborg if it doesn't already exist.
   */
  async emailSignin(
    userEmail: string,
    targetPublicKey: string,
    magicLink?: string,
  ): Promise<TurnkeyCreateSuborgResponse> {
    // lowercase email address.
    const lowerEmail = userEmail.trim().toLowerCase();
    let suborg: TurnkeyCreateSuborgResponse | undefined = await this.getSuborg({
      email: lowerEmail,
    });
    if (!suborg) {
      suborg = await this.createSuborg({
        email: lowerEmail,
      });
    }

    const magicLinkTemplate = config.TURNKEY_MAGIC_LINK_TEMPLATE || magicLink;
    const emailAuthResponse = await this.turnkeyApiClient.emailAuth({
      email: lowerEmail,
      targetPublicKey,
      sendFromEmailAddress: config.TURNKEY_EMAIL_SENDER_ADDRESS,
      sendFromEmailSenderName: config.TURNKEY_EMAIL_SENDER_NAME,
      emailCustomization: {
        appName: TURNKEY_EMAIL_CUSTOMIZATION.APP_NAME,
        logoUrl: TURNKEY_EMAIL_CUSTOMIZATION.LOGO_URL,
        magicLinkTemplate: magicLinkTemplate ? `${magicLinkTemplate}=%s` : undefined,
      },
      invalidateExisting: true, // Invalidates any existing sessions for this user
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

  /**
   * Performs social OAuth authentication with Turnkey.
   * Creates a suborg if one doesn't already exist, then login with the oidc token.
   */
  async socialSignin(
    provider: string,
    oidcToken: string,
    targetPublicKey: string,
  ) {
    // Extract email from Google OIDC token if available
    const extractedEmail = extractEmailFromOidcToken(oidcToken, provider)?.toLowerCase();

    let suborg: TurnkeyCreateSuborgResponse | undefined = await this.getSuborg({
      oidcToken,
    });

    if (!suborg) {
      suborg = await this.getSuborg({
        email: extractedEmail,
      });
      if (suborg) {
        return {
          alreadyExists: true,
        };
      }
    }

    if (!suborg) {
      suborg = await this.createSuborg({
        providerName: provider,
        oidcToken,
        email: extractedEmail, // Include extracted email when creating suborg
      });
    }

    const oauthLoginResponse = await this.turnkeyApiClient.oauthLogin({
      oidcToken,
      publicKey: targetPublicKey,
      invalidateExisting: true, // Invalidates any existing sessions for this user
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
  async passkeySignin(
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
}

/**
 * Extracts email from Google OIDC token payload
 * @param oidcToken - The OIDC token from Google
 * @param providerName - The OAuth provider name
 * @returns The email address if found and provider is Google/Apple, otherwise undefined
 */
export function extractEmailFromOidcToken(
  oidcToken: string,
  providerName: string,
): string | undefined {
  // Only extract email from Google/Apple tokens
  if (providerName.toLowerCase() !== 'google' && providerName.toLowerCase() !== 'apple') {
    return undefined;
  }

  try {
    // Use jose library to decode JWT token without verification
    // We don't verify the signature since we just need to extract the email claim
    const payload = decodeJwt(oidcToken);

    // Extract email from token payload
    const email = payload.email;

    if (!email || typeof email !== 'string') {
      logger.warning({
        at: 'TurnkeyHelpers#extractEmailFromOidcToken',
        message: 'Email not found in OIDC token payload',
        hasEmailField: !!payload.email,
      });
      return undefined;
    }

    return email;
  } catch (error) {
    logger.error({
      at: 'TurnkeyHelpers#extractEmailFromOidcToken',
      message: 'Failed to decode OIDC token',
      error: error instanceof Error ? error.message : error,
    });
    return undefined;
  }
}
