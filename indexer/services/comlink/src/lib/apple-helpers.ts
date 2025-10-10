import { logger } from '@dydxprotocol-indexer/base';
import { SignJWT, JWTPayload } from 'jose';
import fetch from 'node-fetch';

import { AppleJWTClaims, AppleTokenResponse } from '../types';
import { TurnkeyError } from './errors';

/**
 * Helper class for Apple Sign-In operations
 */
export class AppleHelpers {
  /**
   * Generates a JWT client secret for Apple Sign-In authentication
   * @param teamId - Apple Developer Team ID
   * @param serviceId - Apple Service ID (client_id)
   * @param keyId - Apple Key ID
   * @param privateKey - Apple private key in PEM format
   * @returns JWT client secret string
   */
  static async generateClientSecret(
    teamId: string,
    serviceId: string,
    keyId: string,
    privateKey: string,
  ): Promise<string> {
    try {
      const now = Math.floor(Date.now() / 1000);
      const exp = now + (60 * 60 * 24 * 180); // 6 months max

      const claims: AppleJWTClaims & JWTPayload = {
        iss: teamId,
        iat: now,
        exp,
        aud: 'https://appleid.apple.com',
        sub: serviceId,
      };

      // Parse the private key
      const key = await this.parsePrivateKey(privateKey);

      // Create and sign the JWT
      const jwt = await new SignJWT(claims)
        .setProtectedHeader({
          alg: 'ES256',
          kid: keyId,
        })
        .sign(key);

      return jwt;
    } catch (error) {
      logger.error({
        at: 'AppleHelpers#generateClientSecret',
        message: 'Failed to generate Apple client secret',
        error: error instanceof Error ? error.message : error,
      });
      throw new TurnkeyError(`Failed to generate Apple client secret: ${error instanceof Error ? error.message : String(error)}`);
    }
  }

  /**
   * Exchanges Apple authorization code for ID token
   * @param code - Authorization code from Apple
   * @param teamId - Apple Developer Team ID
   * @param serviceId - Apple Service ID
   * @param keyId - Apple Key ID
   * @param privateKey - Apple private key
   * @returns Apple token response with ID token
   */
  static async fetchTokenFromCode(
    code: string,
    teamId: string,
    serviceId: string,
    keyId: string,
    privateKey: string,
  ): Promise<AppleTokenResponse> {
    try {
      const clientSecret = await this.generateClientSecret(teamId, serviceId, keyId, privateKey);

      const bodyParams = new URLSearchParams({
        client_id: serviceId,
        client_secret: clientSecret,
        code,
        grant_type: 'authorization_code',
      });

      const response = await fetch('https://appleid.apple.com/auth/token', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: bodyParams.toString(),
      });

      if (!response.ok) {
        const errorText = await response.text();
        logger.error({
          at: 'AppleHelpers#fetchTokenFromCode',
          message: 'Apple token exchange failed',
          status: response.status,
          errorText,
        });
        throw new TurnkeyError(`Apple token exchange failed: ${response.status} ${errorText}`);
      }

      const tokenResponse = await response.json() as AppleTokenResponse;

      if (!tokenResponse.id_token) {
        throw new TurnkeyError('No ID token received from Apple');
      }

      return tokenResponse;
    } catch (error) {
      logger.error({
        at: 'AppleHelpers#fetchTokenFromCode',
        message: 'Failed to fetch Apple token',
        error: error instanceof Error ? error.message : error,
      });
      throw new TurnkeyError(`Failed to fetch Apple token: ${error instanceof Error ? error.message : String(error)}`);
    }
  }

  /**
   * Parses Apple private key from PEM format
   * @param privateKey - Private key in PEM format
   * @returns CryptoKey for signing
   */
  private static async parsePrivateKey(privateKey: string): Promise<CryptoKey> {
    try {
      // Clean up the private key string
      const cleanKey = privateKey
        .replace(/-----BEGIN PRIVATE KEY-----/, '')
        .replace(/-----END PRIVATE KEY-----/, '')
        .replace(/\s/g, '');

      const keyBuffer = Buffer.from(cleanKey, 'base64');

      return await crypto.subtle.importKey(
        'pkcs8',
        keyBuffer,
        {
          name: 'ECDSA',
          namedCurve: 'P-256',
        },
        false,
        ['sign'],
      );
    } catch (error) {
      logger.error({
        at: 'AppleHelpers#parsePrivateKey',
        message: 'Failed to parse Apple private key',
        error: error instanceof Error ? error.message : error,
      });
      throw new TurnkeyError(`Failed to parse Apple private key: ${error instanceof Error ? error.message : String(error)}`);
    }
  }

  /**
   * Extracts email from Apple ID token
   * @param idToken - Apple ID token
   * @returns Email address if found
   */
  static extractEmailFromIdToken(idToken: string): string | undefined {
    try {
      // Decode the JWT payload (without verification since we just need to extract email)
      const parts = idToken.split('.');
      if (parts.length !== 3) {
        throw new Error('Invalid JWT format');
      }

      const payload = JSON.parse(Buffer.from(parts[1], 'base64').toString());
      return payload.email;
    } catch (error) {
      logger.warning({
        at: 'AppleHelpers#extractEmailFromIdToken',
        message: 'Failed to extract email from Apple ID token',
        error: error instanceof Error ? error.message : error,
      });
      return undefined;
    }
  }
}
