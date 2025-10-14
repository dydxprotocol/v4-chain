import { logger } from '@dydxprotocol-indexer/base';

import { TurnkeyAuthResponse } from '../types';
import { TurnkeyError } from './errors';

/**
 * Helper class for encryption/decryption operations
 */
export class EncryptionHelpers {
  /**
   * Encrypts a TurnkeyAuthResponse using the client's public key
   * For now, we'll use a simple base64 encoding as a placeholder
   * In production, this should use proper public key encryption (e.g., RSA-OAEP)
   *
   * @param payload - The TurnkeyAuthResponse to encrypt
   * @param publicKey - The client's public key (hex string)
   * @returns Encrypted payload as base64 string
   */
  static encryptPayload(
    payload: TurnkeyAuthResponse,
    publicKey: string,
  ): string {
    try {
      // For now, we'll use simple base64 encoding
      // In production, this should be replaced with proper public key encryption
      const payloadString = JSON.stringify(payload);
      const encrypted = Buffer.from(payloadString, 'utf8').toString('base64');

      logger.info({
        at: 'EncryptionHelpers#encryptPayload',
        message: 'Payload encrypted successfully',
        publicKey: `${publicKey.substring(0, 10)}...`,
      });

      return encrypted;
    } catch (error) {
      logger.error({
        at: 'EncryptionHelpers#encryptPayload',
        message: 'Failed to encrypt payload',
        error: error instanceof Error ? error.message : error,
      });
      throw new TurnkeyError(`Failed to encrypt payload: ${error instanceof Error ? error.message : error}`);
    }
  }

  /**
   * Creates a redirect URL for the mobile app
   * @param appScheme - The app scheme (e.g., 'dydxV4')
   * @param encryptedPayload - The encrypted payload
   * @returns The redirect URL
   */
  static createRedirectUrl(appScheme: string, encryptedPayload: string): string {
    const encodedPayload = encodeURIComponent(encryptedPayload);
    return `${appScheme}:///onboard/turnkey?appleLogin=${encodedPayload}`;
  }
}
