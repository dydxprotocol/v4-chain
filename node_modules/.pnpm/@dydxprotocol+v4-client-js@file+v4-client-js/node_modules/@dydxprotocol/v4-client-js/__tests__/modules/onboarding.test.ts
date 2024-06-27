import { exportMnemonicAndPrivateKey, deriveHDKeyFromEthereumSignature } from '../../src/lib/onboarding';
import {
  DirectSecp256k1HdWallet,
  DirectSecp256k1Wallet,
  OfflineSigner,
} from '@cosmjs/proto-signing';
import {
  MNEMONIC_FROM_SIGNATURE_RESULT,
  ENTROPY_FROM_SIGNATURE_RESULT,
  PRIVATE_KEY_FROM_SIGNATURE_RESULT,
  PUBLIC_KEY_FROM_SIGNATURE_RESULT,
  SIGNATURE_RESULT,
} from '../helpers/constants';

describe('Onboarding', () => {
  describe('exportMnemonicAndPrivateKey', () => {
    it('Successfully creates HDKey', () => {
      expect(exportMnemonicAndPrivateKey(ENTROPY_FROM_SIGNATURE_RESULT))
        .toEqual({
          mnemonic: MNEMONIC_FROM_SIGNATURE_RESULT,
          privateKey: PRIVATE_KEY_FROM_SIGNATURE_RESULT,
          publicKey: PUBLIC_KEY_FROM_SIGNATURE_RESULT,
        });
    });

    it('Expect mnemonic and private key to generate the same address', async () => {
      const { mnemonic, privateKey } = exportMnemonicAndPrivateKey(ENTROPY_FROM_SIGNATURE_RESULT);

      const wallet: OfflineSigner = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic!, {
        prefix: 'dydx',
      });
      const pkWallet = await DirectSecp256k1Wallet.fromKey(privateKey!, 'dydx');
      const mnemonicAddress = (await wallet.getAccounts())[0].address;
      const pkAddress = (await pkWallet.getAccounts())[0].address;
      expect(mnemonicAddress).toEqual(pkAddress);
    });
  });

  describe('deriveHDKeyFromEthereumSignature', () => {
    it('Throw error if signature not valid', () => {
      expect(() => deriveHDKeyFromEthereumSignature('0x04db21dfsa321b')).toThrow();
    });
    it('Successfully creates HDKey', () => {
      expect(deriveHDKeyFromEthereumSignature(SIGNATURE_RESULT)).toEqual({
        mnemonic: MNEMONIC_FROM_SIGNATURE_RESULT,
        privateKey: PRIVATE_KEY_FROM_SIGNATURE_RESULT,
        publicKey: PUBLIC_KEY_FROM_SIGNATURE_RESULT,
      });
    });

    it('Successfully creates HDKey from signature with different v values', () => {
      const suffixes = ['00', '1b', '01', '1c'];

      const testSignatures = suffixes.map((suffix) => SIGNATURE_RESULT.slice(0, -2) + suffix);
      testSignatures.forEach((sig) => {
        expect(deriveHDKeyFromEthereumSignature(sig)).toEqual({
          mnemonic: MNEMONIC_FROM_SIGNATURE_RESULT,
          privateKey: PRIVATE_KEY_FROM_SIGNATURE_RESULT,
          publicKey: PUBLIC_KEY_FROM_SIGNATURE_RESULT,
        });
      });
    });
  });
});
