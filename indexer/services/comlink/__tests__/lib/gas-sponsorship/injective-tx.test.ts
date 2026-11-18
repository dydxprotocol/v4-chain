import { fromBase64, toBase64, toHex } from '@cosmjs/encoding';

import {
  assembleSignedTx,
  buildGrantTx,
  encodeEthSecp256k1PubKey,
  getTxHash,
  GrantTxParams,
  timestampFromDate,
  UnsignedTx,
} from '../../../src/lib/gas-sponsorship/injective-tx';
import { verifyDigestSignature } from '../../../src/lib/gas-sponsorship/sponsor-signer';

/**
 * Two fee-grant transactions taken from Injective mainnet, with the parameters that produce
 * them. Rebuilding each and getting the same bytes, and a sign digest the on-chain signature
 * verifies against, shows that what this module encodes and signs is what the chain accepts.
 */
const MAINNET_GRANTER: string = 'inj1xdrkxgke8qzey9ze4cx48heav6lzkkszuc224r';
const MAINNET_GRANTER_PUBLIC_KEY: Uint8Array = fromBase64(
  'Aq8nnehaBv1QXUtBy+VKqYU5m2rL3wqAWqC566gKMq4r',
);
const MAINNET_GRANTER_ACCOUNT_NUMBER: bigint = BigInt(1682332);
const MAINNET_ALLOWED_MESSAGES: string[] = [
  '/cosmos.authz.v1beta1.MsgExec',
  '/injective.wasmx.v1.MsgExecuteContractCompat',
];

interface MainnetFixture {
  name: string,
  txHash: string,
  params: GrantTxParams,
  signature: string,
  txBytes: string,
}

const GRANT_ONLY: MainnetFixture = {
  name: 'a grant on its own (block 192295530)',
  txHash: '504DB8C20E7367E1EAC27FD2B84E39784B61FEB77FAC064C0C3709AF1EE8D147',
  params: {
    chainId: 'injective-1',
    granter: MAINNET_GRANTER,
    grantee: 'inj1edsansmcyfr5ucw3vxwv0ts8jg02ugzhdsa6d8',
    granterPublicKey: MAINNET_GRANTER_PUBLIC_KEY,
    accountNumber: MAINNET_GRANTER_ACCOUNT_NUMBER,
    sequence: BigInt(55553),
    spendLimit: { denom: 'inj', amount: BigInt('35200000000000000') },
    // 2026-11-19T17:16:31.486304377Z
    expiration: { seconds: BigInt(Date.UTC(2026, 10, 19, 17, 16, 31) / 1000), nanos: 486304377 },
    allowedMessages: MAINNET_ALLOWED_MESSAGES,
    fee: { denom: 'inj', amount: BigInt('48000000000000') },
    gasLimit: BigInt(300000),
    revokeExisting: false,
  },
  signature: 'aznuadr0mIGOvp5F/7whosoCaG/1Oykq8JkYTOYCQARSe03vZ+9B6sWfs0fb/Q5I6PzGIqGm+PSXwxSQXhn3hgA=',
  txBytes: 'CuACCt0CCiovY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuTXNnR3JhbnRBbGxvd2FuY2USrgIKKmluajF4ZHJr' +
    'eGdrZThxemV5OXplNGN4NDhoZWF2Nmx6a2tzenVjMjI0chIqaW5qMWVkc2Fuc21jeWZyNXVjdzN2eHd2MHRzOGpn' +
    'MDJ1Z3poZHNhNmQ4GtMBCiwvY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuQWxsb3dlZE1zZ0FsbG93YW5jZRKiAQpT' +
    'CicvY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuQmFzaWNBbGxvd2FuY2USKAoYCgNpbmoSETM1MjAwMDAwMDAwMDAw' +
    'MDAwEgwI79381wYQ+dTx5wESHS9jb3Ntb3MuYXV0aHoudjFiZXRhMS5Nc2dFeGVjEiwvaW5qZWN0aXZlLndhc214' +
    'LnYxLk1zZ0V4ZWN1dGVDb250cmFjdENvbXBhdBJ/CmAKVAotL2luamVjdGl2ZS5jcnlwdG8udjFiZXRhMS5ldGhz' +
    'ZWNwMjU2azEuUHViS2V5EiMKIQKvJ53oWgb9UF1LQcvlSqmFOZtqy98KgFqgueuoCjKuKxIECgIIARiBsgMSGwoV' +
    'CgNpbmoSDjQ4MDAwMDAwMDAwMDAwEOCnEhpBaznuadr0mIGOvp5F/7whosoCaG/1Oykq8JkYTOYCQARSe03vZ+9B' +
    '6sWfs0fb/Q5I6PzGIqGm+PSXwxSQXhn3hgA=',
};

const REVOKE_THEN_GRANT: MainnetFixture = {
  name: 'a revoke followed by a grant (block 192295827)',
  txHash: 'C554B3ACB4F98161CCFE7938399A73C4D452FBD9EF92F88C4F79E68896905BE4',
  params: {
    chainId: 'injective-1',
    granter: MAINNET_GRANTER,
    grantee: 'inj1frwsqymyg3wnfhy4mqrq9cghwwzst0tgsmfal7',
    granterPublicKey: MAINNET_GRANTER_PUBLIC_KEY,
    accountNumber: MAINNET_GRANTER_ACCOUNT_NUMBER,
    sequence: BigInt(55554),
    spendLimit: { denom: 'inj', amount: BigInt('35200000000000000') },
    // 2026-11-19T17:19:29.652016699Z
    expiration: { seconds: BigInt(Date.UTC(2026, 10, 19, 17, 19, 29) / 1000), nanos: 652016699 },
    allowedMessages: MAINNET_ALLOWED_MESSAGES,
    fee: { denom: 'inj', amount: BigInt('96000000000000') },
    gasLimit: BigInt(600000),
    revokeExisting: true,
  },
  signature: 'ilgYbrEYeAk2LvFoK7l5BrZkgAHItMeM3oQCCkWGpZMHz7ExpX2GcBYZyMb1KYvSRlpL/xw4KHAX0EhlLOCMOgA=',
  txBytes: 'CuoDCocBCisvY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuTXNnUmV2b2tlQWxsb3dhbmNlElgKKmluajF4ZHJr' +
    'eGdrZThxemV5OXplNGN4NDhoZWF2Nmx6a2tzenVjMjI0chIqaW5qMWZyd3NxeW15ZzN3bmZoeTRtcXJxOWNnaHd3' +
    'enN0MHRnc21mYWw3Ct0CCiovY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuTXNnR3JhbnRBbGxvd2FuY2USrgIKKmlu' +
    'ajF4ZHJreGdrZThxemV5OXplNGN4NDhoZWF2Nmx6a2tzenVjMjI0chIqaW5qMWZyd3NxeW15ZzN3bmZoeTRtcXJx' +
    'OWNnaHd3enN0MHRnc21mYWw3GtMBCiwvY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuQWxsb3dlZE1zZ0FsbG93YW5j' +
    'ZRKiAQpTCicvY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuQmFzaWNBbGxvd2FuY2USKAoYCgNpbmoSETM1MjAwMDAw' +
    'MDAwMDAwMDAwEgwIod/81wYQu/jztgISHS9jb3Ntb3MuYXV0aHoudjFiZXRhMS5Nc2dFeGVjEiwvaW5qZWN0aXZl' +
    'Lndhc214LnYxLk1zZ0V4ZWN1dGVDb250cmFjdENvbXBhdBJ/CmAKVAotL2luamVjdGl2ZS5jcnlwdG8udjFiZXRh' +
    'MS5ldGhzZWNwMjU2azEuUHViS2V5EiMKIQKvJ53oWgb9UF1LQcvlSqmFOZtqy98KgFqgueuoCjKuKxIECgIIARiC' +
    'sgMSGwoVCgNpbmoSDjk2MDAwMDAwMDAwMDAwEMDPJBpBilgYbrEYeAk2LvFoK7l5BrZkgAHItMeM3oQCCkWGpZMH' +
    'z7ExpX2GcBYZyMb1KYvSRlpL/xw4KHAX0EhlLOCMOgA=',
};

describe('injective-tx', () => {
  describe.each([GRANT_ONLY, REVOKE_THEN_GRANT])('rebuilding $name', (fixture: MainnetFixture) => {
    const unsignedTx: UnsignedTx = buildGrantTx(fixture.params);
    // Signatures on chain carry a trailing recovery byte, which verification ignores.
    const signature: Uint8Array = fromBase64(fixture.signature);

    it('produces the same transaction bytes', () => {
      expect(toBase64(assembleSignedTx(unsignedTx, signature))).toEqual(fixture.txBytes);
    });

    it('produces a sign digest that the on-chain signature is valid for', async () => {
      expect(signature).toHaveLength(65);
      expect(
        await verifyDigestSignature(
          signature.slice(0, 64),
          unsignedTx.signDigest,
          fixture.params.granterPublicKey,
        ),
      ).toEqual(true);
    });

    it('produces a different sign digest if anything signed changes', async () => {
      const changes: Partial<GrantTxParams>[] = [
        { chainId: 'injective-888' },
        { accountNumber: fixture.params.accountNumber + BigInt(1) },
        { sequence: fixture.params.sequence + BigInt(1) },
        { grantee: 'inj1v4d4hf84ea7cf74pcgzghhf5ek8q6empm7gejp' },
        { spendLimit: { denom: 'inj', amount: BigInt('35200000000000001') } },
        {
          expiration: {
            seconds: fixture.params.expiration.seconds,
            nanos: fixture.params.expiration.nanos + 1,
          },
        },
        { allowedMessages: ['/cosmos.authz.v1beta1.MsgExec'] },
        { fee: { denom: 'inj', amount: BigInt(1) } },
        { gasLimit: fixture.params.gasLimit + BigInt(1) },
        { revokeExisting: !fixture.params.revokeExisting },
      ];

      for (const change of changes) {
        const changed: UnsignedTx = buildGrantTx({ ...fixture.params, ...change });
        expect(toHex(changed.signDigest)).not.toEqual(toHex(unsignedTx.signDigest));
        expect(
          await verifyDigestSignature(
            signature.slice(0, 64),
            changed.signDigest,
            fixture.params.granterPublicKey,
          ),
        ).toEqual(false);
      }
    });

    it('computes the hash the transaction has on chain', () => {
      expect(getTxHash(fromBase64(fixture.txBytes))).toEqual(fixture.txHash);
    });
  });

  describe('timestampFromDate', () => {
    it('splits a date into whole seconds and nanoseconds', () => {
      expect(timestampFromDate(new Date('2026-11-19T17:16:31.486Z'))).toEqual({
        seconds: BigInt(1795108591),
        nanos: 486000000,
      });
    });

    it('has no nanoseconds for a whole second', () => {
      expect(timestampFromDate(new Date('2026-11-19T17:16:31.000Z'))).toEqual({
        seconds: BigInt(1795108591),
        nanos: 0,
      });
    });

    it('rejects an invalid date', () => {
      expect(() => timestampFromDate(new Date('not a date'))).toThrow('invalid date');
    });
  });

  describe('encodeEthSecp256k1PubKey', () => {
    it('rejects a key that is not a 33-byte compressed key', () => {
      expect(() => encodeEthSecp256k1PubKey(new Uint8Array(65))).toThrow(
        'expected a 33-byte compressed public key, got 65 bytes',
      );
    });
  });
});
