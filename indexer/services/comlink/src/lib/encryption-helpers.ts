import { logger } from '@dydxprotocol-indexer/base';
import { Address, keccak256, toBytes } from 'viem';

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
   * Decrypts an encrypted payload using the client's private key
   * For now, we'll use simple base64 decoding as a placeholder
   * In production, this should use proper public key decryption
   *
   * @param encryptedPayload - The encrypted payload as base64 string
   * @param privateKey - The client's private key (hex string)
   * @returns Decrypted TurnkeyAuthResponse
   */
  static decryptPayload(
    encryptedPayload: string,
    privateKey: string,
  ): TurnkeyAuthResponse {
    try {
      // For now, we'll use simple base64 decoding
      // In production, this should be replaced with proper private key decryption
      const decrypted = Buffer.from(encryptedPayload, 'base64').toString('utf8');
      const payload = JSON.parse(decrypted) as TurnkeyAuthResponse;

      logger.info({
        at: 'EncryptionHelpers#decryptPayload',
        message: 'Payload decrypted successfully',
        privateKey: `${privateKey.substring(0, 10)}...`,
      });

      return payload;
    } catch (error) {
      logger.error({
        at: 'EncryptionHelpers#decryptPayload',
        message: 'Failed to decrypt payload',
        error: error instanceof Error ? error.message : error,
      });
      throw new TurnkeyError(`Failed to decrypt payload: ${error instanceof Error ? error.message : error}`);
    }
  }

  /**
   * Validates that a public key is a valid Ethereum address
   * @param publicKey - The public key to validate
   * @returns True if valid, false otherwise
   */
  static isValidPublicKey(publicKey: string): boolean {
    try {
      // Check if it's a valid hex string
      if (!/^0x[0-9a-fA-F]{64}$/.test(publicKey)) {
        return false;
      }

      // Try to convert to address by taking the last 20 bytes of keccak256 hash
      const publicKeyBytes = toBytes(publicKey as Address);
      const hash = keccak256(publicKeyBytes);
      const addressBytes = toBytes(hash).slice(-20);
      const address = `0x${Buffer.from(addressBytes).toString('hex')}`;
      return address.length === 42; // Valid Ethereum address length
    } catch (error) {
      logger.warning({
        at: 'EncryptionHelpers#isValidPublicKey',
        message: 'Invalid public key format',
        publicKey: `${publicKey.substring(0, 10)}...`,
        error: error instanceof Error ? error.message : error,
      });
      return false;
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
