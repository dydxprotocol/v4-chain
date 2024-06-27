import { BECH32_PREFIX } from '../src';
import LocalWallet from '../src/clients/modules/local-wallet';
import { DYDX_TEST_ADDRESS, DYDX_TEST_MNEMONIC } from './constants';

async function test(): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(DYDX_TEST_MNEMONIC, BECH32_PREFIX);
  console.log(wallet);
  const address = wallet.address;
  const addressOk = (address === DYDX_TEST_ADDRESS);
  console.log(addressOk);
  console.log(address);
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
