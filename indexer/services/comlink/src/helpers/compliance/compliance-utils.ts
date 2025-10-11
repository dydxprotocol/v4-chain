import {
  ExtendedSecp256k1Signature, Secp256k1, ripemd160, sha256,
} from '@cosmjs/crypto';
import { toBech32 } from '@cosmjs/encoding';
import { GeoOriginHeaders, isRestrictedCountryHeaders } from '@dydxprotocol-indexer/compliance';
import { ComplianceReason } from '@dydxprotocol-indexer/postgres';
import { verifyADR36Amino } from '@keplr-wallet/cosmos';
import express from 'express';
import { DateTime } from 'luxon';

import {
  WRITE_REQUEST_TTL_SECONDS,
  DYDX_ADDRESS_PREFIX,
} from '../../lib/constants';
import { create4xxResponse } from '../../lib/helpers';

export enum ComplianceAction {
  CONNECT = 'CONNECT',
  VALID_SURVEY = 'VALID_SURVEY',
  INVALID_SURVEY = 'INVALID_SURVEY',
}

export enum AccountVerificationRequiredAction {
  REGISTER_TOKEN = 'REGISTER_TOKEN',
  UPDATE_CODE = 'UPDATE_CODE',
}

// TODO: deprecate this pattern.
// Only the origin is pertinent, so store country or region, rather than encoding them.
export function getGeoComplianceReason(
  headers: GeoOriginHeaders,
): ComplianceReason | undefined {
  if (isRestrictedCountryHeaders(headers)) {
    const country: string | undefined = headers['cf-ipcountry'];
    if (country === 'US') {
      return ComplianceReason.US_GEO;
    } else if (country === 'CA') {
      return ComplianceReason.CA_GEO;
    } else if (country === 'GB') {
      return ComplianceReason.GB_GEO;
    } else {
      return ComplianceReason.SANCTIONED_GEO;
    }
  }
  return undefined;
}

function generateAddress(pubkeyArray: Uint8Array): string {
  return toBech32('dydx', ripemd160(sha256(pubkeyArray)));
}

/**
 * Validates a signature by performing various checks including address format,
 * public key correspondence, timestamp validity, and signature verification.
 *
 * @returns {Promise<express.Response | undefined>} Returns undefined if validation
 * is successful. Returns an HTTP response with an error message if validation fails.
 */
export async function validateSignature(
  res: express.Response,
  action: ComplianceAction | AccountVerificationRequiredAction,
  address: string,
  timestamp: number,
  message: string,
  signedMessage: string,
  pubkey: string,
  currentStatus?: string,
): Promise<express.Response| undefined> {
  if (!address.startsWith(DYDX_ADDRESS_PREFIX)) {
    return create4xxResponse(
      res,
      `Address ${address} is not a valid dYdX V4 address`,
    );
  }

  const pubkeyArray: Uint8Array = new Uint8Array(Buffer.from(pubkey, 'base64'));
  if (address !== generateAddress(pubkeyArray)) {
    return create4xxResponse(
      res,
      `Address ${address} does not correspond to the pubkey provided ${pubkey}`,
    );
  }

  // Verify the timestamp is within WRITE_REQUEST_TTL_SECONDS seconds of the current time
  const now = DateTime.now().toSeconds();
  if (Math.abs(now - timestamp) > WRITE_REQUEST_TTL_SECONDS) {
    return create4xxResponse(
      res,
      `Timestamp is not within the valid range of ${WRITE_REQUEST_TTL_SECONDS} seconds`,
    );
  }

  // Prepare the message for verification
  const messageToSign: string = `${message}:${action}"${
    currentStatus || ''
  }:${timestamp}`;
  const messageHash: Uint8Array = sha256(Buffer.from(messageToSign));
  const signedMessageArray: Uint8Array = new Uint8Array(
    Buffer.from(signedMessage, 'base64'),
  );
  const signature: ExtendedSecp256k1Signature = ExtendedSecp256k1Signature.fromFixedLength(
    signedMessageArray,
  );

  // Verify the signature
  const isValidSignature: boolean = await Secp256k1.verifySignature(
    signature,
    messageHash,
    pubkeyArray,
  );
  if (!isValidSignature) {
    return create4xxResponse(res, 'Signature verification failed');
  }

  return undefined;
}

/**
 * Validates a signature using verifyADR36Amino provided by keplr package.
 *
 * @returns {Promise<express.Response | undefined>} Returns undefined if validation
 * is successful. Returns an HTTP response with an error message if validation fails.
 */
export function validateSignatureKeplr(
  res:express.Response,
  address: string,
  message: string,
  signedMessage: string,
  pubkey: string,
): express.Response | undefined {
  const messageToSign: string = message;

  const pubKeyUint = new Uint8Array(Buffer.from(pubkey, 'base64'));
  const signedMessageUint = new Uint8Array(Buffer.from(signedMessage, 'base64'));

  const isVerified = verifyADR36Amino(
    'dydx', address, messageToSign, pubKeyUint, signedMessageUint, 'secp256k1',
  );

  if (!isVerified) {
    return create4xxResponse(
      res,
      'Keplr signature verification failed',
    );
  }

  return undefined;
}
