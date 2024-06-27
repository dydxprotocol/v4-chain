import { Network } from '../../../src/clients/constants';
import LocalWallet from '../../../src/clients/modules/local-wallet';
import { SubaccountInfo } from '../../../src/clients/subaccount';
import { ValidatorClient } from '../../../src/clients/validator-client';
import { DYDX_TEST_MNEMONIC } from '../../../examples/constants';
import Long from 'long';
import { BECH32_PREFIX } from '../../../src';

async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

describe('Validator Client', () => {
  let wallet: LocalWallet;
  let subaccount: SubaccountInfo;
  let client: ValidatorClient;

  describe('Transfers', () => {
    beforeEach(async () => {
      wallet = await LocalWallet.fromMnemonic(DYDX_TEST_MNEMONIC, BECH32_PREFIX);
      subaccount = new SubaccountInfo(wallet, 0);
      client = await ValidatorClient.connect(Network.testnet().validatorConfig);
      await sleep(5000);  // wait for withdraw to complete
    });

    it('Withdraw', async () => {
      const tx = await client.post.withdraw(
        subaccount,
        0,
        new Long(1_00_000_000),
        undefined,
      );
      console.log('**Withdraw Tx**');
      console.log(tx);
    });

    it('Deposit', async () => {
      const tx = await client.post.deposit(
        subaccount,
        0,
        new Long(1_000_000),
      );
      console.log('**Deposit Tx**');
      console.log(tx);
    });

    it('Transfer', async () => {
      const tx = await client.post.transfer(
        subaccount,
        subaccount.address,
        1,
        0,
        new Long(1_000),
      );
      console.log('**Transfer Tx**');
      console.log(tx);
    });
  });
});
