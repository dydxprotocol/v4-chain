import { keccak256 } from '@cosmjs/crypto';
import { fromBase64, fromHex, toHex } from '@cosmjs/encoding';

import { isValidAccountAddress } from '../../../src/lib/gas-sponsorship/addresses';
import {
  ethereumAddressFromPublicKey,
  injectiveAddressFromPublicKey,
  parseCompressedPublicKey,
  toCanonicalSignature,
  TurnkeyRawSigner,
  TurnkeySponsorSigner,
  verifyDigestSignature,
} from '../../../src/lib/gas-sponsorship/sponsor-signer';
import { LocalKeySigner } from '../../helpers/gas-sponsorship-helpers';

// The public key and address of an account on Injective mainnet.
const MAINNET_PUBLIC_KEY: Uint8Array = fromBase64('Aq8nnehaBv1QXUtBy+VKqYU5m2rL3wqAWqC566gKMq4r');
const MAINNET_ADDRESS: string = 'inj1xdrkxgke8qzey9ze4cx48heav6lzkkszuc224r';

// The well-known key and address of the first account of the Hardhat and Anvil development
// chains.
const DEV_PRIVATE_KEY: string = 'ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80';
const DEV_PUBLIC_KEY: string = '038318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed75';
const DEV_ETHEREUM_ADDRESS: string = '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266';

const SECP256K1_ORDER: bigint = BigInt(
  '0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141',
);
const DIGEST: Uint8Array = keccak256(Buffer.from('a message'));

function hex32(value: bigint): string {
  return value.toString(16).padStart(64, '0');
}

describe('sponsor-signer', () => {
  describe('address derivation', () => {
    it('derives the Injective address of a mainnet account from its public key', () => {
      expect(injectiveAddressFromPublicKey(MAINNET_PUBLIC_KEY, 'inj')).toEqual(MAINNET_ADDRESS);
    });

    it('derives the checksummed Ethereum address of a well-known key', () => {
      expect(ethereumAddressFromPublicKey(fromHex(DEV_PUBLIC_KEY))).toEqual(DEV_ETHEREUM_ADDRESS);
    });
  });

  describe('parseCompressedPublicKey', () => {
    it('parses hex with or without a 0x prefix', () => {
      expect(toHex(parseCompressedPublicKey(DEV_PUBLIC_KEY))).toEqual(DEV_PUBLIC_KEY);
      expect(toHex(parseCompressedPublicKey(`0x${DEV_PUBLIC_KEY}`))).toEqual(DEV_PUBLIC_KEY);
    });

    it.each([
      ['is not hex', 'zz'.repeat(33), 'sponsor public key is not valid hex'],
      [
        'is an uncompressed key',
        `04${'11'.repeat(64)}`,
        'sponsor public key must be a 33-byte compressed secp256k1 key',
      ],
      [
        'is too short',
        DEV_PUBLIC_KEY.slice(2),
        'sponsor public key must be a 33-byte compressed secp256k1 key',
      ],
      [
        'has the wrong prefix byte',
        `05${DEV_PUBLIC_KEY.slice(2)}`,
        'sponsor public key must be a 33-byte compressed secp256k1 key',
      ],
      ['is empty', '', 'sponsor public key must be a 33-byte compressed secp256k1 key'],
    ])('rejects a key that %s', (_name: string, key: string, message: string) => {
      expect(() => parseCompressedPublicKey(key)).toThrow(message);
    });

    it('rejects a key that is not a point on the curve', () => {
      // No point on secp256k1 has this x coordinate.
      expect(() => parseCompressedPublicKey(`02${hex32(BigInt(5))}`)).toThrow();
    });
  });

  describe('isValidAccountAddress', () => {
    it('accepts an account address with the expected prefix', () => {
      expect(isValidAccountAddress(MAINNET_ADDRESS, 'inj')).toEqual(true);
      expect(
        isValidAccountAddress('dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4', 'dydx'),
      ).toEqual(true);
    });

    it.each([
      ['another prefix', MAINNET_ADDRESS, 'dydx'],
      ['a bad checksum', `${MAINNET_ADDRESS.slice(0, -1)}q`, 'inj'],
      ['mixed case', `${MAINNET_ADDRESS.slice(0, -1)}R`, 'inj'],
      ['no data', 'inj', 'inj'],
      ['the empty string', '', 'inj'],
      ['a non-string', 12345, 'inj'],
      ['undefined', undefined, 'inj'],
    ])('rejects %s', (_name: string, address: unknown, prefix: string) => {
      expect(isValidAccountAddress(address, prefix)).toEqual(false);
    });
  });

  describe('toCanonicalSignature', () => {
    it('joins r and s, padded to 32 bytes each', () => {
      expect(toHex(toCanonicalSignature('0x01', '2'))).toEqual(
        `${hex32(BigInt(1))}${hex32(BigInt(2))}`,
      );
    });

    it('leaves a low s as it is', () => {
      const halfOrder: bigint = SECP256K1_ORDER / BigInt(2);

      expect(toHex(toCanonicalSignature('1', hex32(halfOrder))).slice(64)).toEqual(
        hex32(halfOrder),
      );
    });

    it('replaces a high s with its low counterpart', () => {
      const highS: bigint = SECP256K1_ORDER / BigInt(2) + BigInt(1);

      expect(toHex(toCanonicalSignature('1', hex32(highS))).slice(64)).toEqual(
        hex32(SECP256K1_ORDER - highS),
      );
    });

    it.each([
      ['zero r', '0', '1'],
      ['zero s', '1', '0'],
      ['r at the curve order', hex32(SECP256K1_ORDER), '1'],
      ['s at the curve order', '1', hex32(SECP256K1_ORDER)],
    ])('rejects %s', (_name: string, r: string, s: string) => {
      expect(() => toCanonicalSignature(r, s)).toThrow('signature component out of range');
    });

    it.each([
      ['an empty component', '', '1'],
      ['a non-hex component', 'xyz', '1'],
      ['a component longer than 32 bytes', '1', '1'.repeat(65)],
    ])('rejects %s', (_name: string, r: string, s: string) => {
      expect(() => toCanonicalSignature(r, s)).toThrow('is not a 32-byte hex value');
    });
  });

  describe('verifyDigestSignature', () => {
    let signer: LocalKeySigner;
    let signature: Uint8Array;

    beforeAll(async () => {
      signer = await LocalKeySigner.create(DEV_PRIVATE_KEY);
      signature = await signer.signDigest(DIGEST);
    });

    it('accepts a signature of the digest by the key', async () => {
      expect(toHex(signer.publicKey)).toEqual(DEV_PUBLIC_KEY);
      expect(await verifyDigestSignature(signature, DIGEST, signer.publicKey)).toEqual(true);
    });

    it('rejects a signature of another digest', async () => {
      const otherDigest: Uint8Array = keccak256(Buffer.from('another message'));

      expect(await verifyDigestSignature(signature, otherDigest, signer.publicKey)).toEqual(false);
    });

    it('rejects a signature by another key', async () => {
      expect(await verifyDigestSignature(signature, DIGEST, MAINNET_PUBLIC_KEY)).toEqual(false);
    });

    it('rejects a signature of the wrong length', async () => {
      const withRecoveryByte: Uint8Array = Uint8Array.from([...signature, 0]);

      expect(
        await verifyDigestSignature(withRecoveryByte, DIGEST, signer.publicKey),
      ).toEqual(false);
      expect(
        await verifyDigestSignature(signature.slice(0, 63), DIGEST, signer.publicKey),
      ).toEqual(false);
    });
  });

  describe('TurnkeySponsorSigner', () => {
    let localSigner: LocalKeySigner;
    let localSignature: Uint8Array;
    let signRawPayload: jest.Mock;
    let client: TurnkeyRawSigner;

    beforeAll(async () => {
      localSigner = await LocalKeySigner.create(DEV_PRIVATE_KEY);
      localSignature = await localSigner.signDigest(DIGEST);
    });

    beforeEach(() => {
      // Turnkey returns r and s as unprefixed hex.
      signRawPayload = jest.fn().mockResolvedValue({
        r: toHex(localSignature.slice(0, 32)),
        s: toHex(localSignature.slice(32)),
        v: '00',
      });
      client = { signRawPayload };
    });

    it('asks Turnkey to sign the digest as it is, with the key named by its Ethereum address', async () => {
      const signer: TurnkeySponsorSigner = new TurnkeySponsorSigner(client, localSigner.publicKey);

      await signer.signDigest(DIGEST);

      expect(signRawPayload).toHaveBeenCalledTimes(1);
      expect(signRawPayload).toHaveBeenCalledWith({
        signWith: DEV_ETHEREUM_ADDRESS,
        payload: toHex(DIGEST),
        encoding: 'PAYLOAD_ENCODING_HEXADECIMAL',
        hashFunction: 'HASH_FUNCTION_NO_OP',
      });
    });

    it('names the key as configured when given', async () => {
      const signer: TurnkeySponsorSigner = new TurnkeySponsorSigner(
        client,
        localSigner.publicKey,
        'a-private-key-id',
      );

      await signer.signDigest(DIGEST);

      expect(signRawPayload).toHaveBeenCalledWith(
        expect.objectContaining({ signWith: 'a-private-key-id' }),
      );
    });

    it('falls back to the Ethereum address when the configured name is empty', async () => {
      const signer: TurnkeySponsorSigner = new TurnkeySponsorSigner(
        client,
        localSigner.publicKey,
        '',
      );

      await signer.signDigest(DIGEST);

      expect(signRawPayload).toHaveBeenCalledWith(
        expect.objectContaining({ signWith: DEV_ETHEREUM_ADDRESS }),
      );
    });

    it('returns a 64-byte signature that verifies against the public key', async () => {
      const signer: TurnkeySponsorSigner = new TurnkeySponsorSigner(client, localSigner.publicKey);

      const signature: Uint8Array = await signer.signDigest(DIGEST);

      expect(signature).toHaveLength(64);
      expect(await verifyDigestSignature(signature, DIGEST, signer.publicKey)).toEqual(true);
    });

    it('returns a verifying low-s signature when Turnkey returns the high-s form', async () => {
      const s: bigint = BigInt(`0x${toHex(localSignature.slice(32))}`);
      signRawPayload.mockResolvedValue({
        r: toHex(localSignature.slice(0, 32)),
        s: hex32(SECP256K1_ORDER - s),
      });
      const signer: TurnkeySponsorSigner = new TurnkeySponsorSigner(client, localSigner.publicKey);

      const signature: Uint8Array = await signer.signDigest(DIGEST);

      expect(toHex(signature)).toEqual(toHex(localSignature));
      expect(await verifyDigestSignature(signature, DIGEST, signer.publicKey)).toEqual(true);
    });

    it('refuses to sign anything that is not a 32-byte digest', async () => {
      const signer: TurnkeySponsorSigner = new TurnkeySponsorSigner(client, localSigner.publicKey);

      await expect(signer.signDigest(new Uint8Array(31))).rejects.toThrow(
        'expected a 32-byte digest, got 31 bytes',
      );
      expect(signRawPayload).not.toHaveBeenCalled();
    });

    it('passes on a failure from Turnkey', async () => {
      signRawPayload.mockRejectedValue(new Error('policy denied'));
      const signer: TurnkeySponsorSigner = new TurnkeySponsorSigner(client, localSigner.publicKey);

      await expect(signer.signDigest(DIGEST)).rejects.toThrow('policy denied');
    });
  });
});
