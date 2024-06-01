import Long from 'long';

import { BECH32_PREFIX } from '../src';
import { Network } from '../src/clients/constants';
import LocalWallet from '../src/clients/modules/local-wallet';
import { SubaccountInfo } from '../src/clients/subaccount';
import { ValidatorClient } from '../src/clients/validator-client';
import { DYDX_TEST_MNEMONIC } from './constants';

async function test(): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(DYDX_TEST_MNEMONIC, BECH32_PREFIX);
  console.log(wallet);

  const client = await ValidatorClient.connect(Network.testnet().validatorConfig);
  console.log('**Client**');
  console.log(client);

  const subaccount = new SubaccountInfo(wallet, 0);

  const tx = await client.post.transfer(
    subaccount,
    subaccount.address,
    1,
    0,
    new Long(10_000_000),
  );
  console.log('**Transfer Tx**');
  console.log(tx);
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
