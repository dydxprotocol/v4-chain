import { BECH32_PREFIX } from '../src';
import LocalWallet from '../src/clients/modules/local-wallet';
import { KLYRA_TEST_ADDRESS, KLYRA_TEST_MNEMONIC } from './constants';

async function test(): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(KLYRA_TEST_MNEMONIC, BECH32_PREFIX);
  console.log(wallet);
  const address = wallet.address;
  const addressOk = (address === KLYRA_TEST_ADDRESS);
  console.log(addressOk);
  console.log(address);
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
