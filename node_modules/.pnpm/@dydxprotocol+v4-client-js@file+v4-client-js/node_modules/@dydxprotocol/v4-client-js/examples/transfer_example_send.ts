import { EncodeObject } from '@cosmjs/proto-signing';
import { Method } from '@cosmjs/tendermint-rpc';
import Long from 'long';

import { TEST_RECIPIENT_ADDRESS } from '../__tests__/helpers/constants';
import { BECH32_PREFIX } from '../src';
import { Network } from '../src/clients/constants';
import LocalWallet from '../src/clients/modules/local-wallet';
import { SubaccountInfo } from '../src/clients/subaccount';
import { ValidatorClient } from '../src/clients/validator-client';
import { DYDX_TEST_MNEMONIC } from './constants';

async function test(): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(
    DYDX_TEST_MNEMONIC,
    BECH32_PREFIX,
  );
  console.log(wallet);

  const client = await ValidatorClient.connect(Network.testnet().validatorConfig);
  console.log('**Client**');
  console.log(client);

  const subaccount = new SubaccountInfo(wallet, 0);

  const amount = new Long(100_000_000);

  const msgs: Promise<EncodeObject[]> = new Promise((resolve) => {
    const msg = client.post.composer.composeMsgSendToken(
      subaccount.address,
      TEST_RECIPIENT_ADDRESS,
      client.config.denoms.CHAINTOKEN_DENOM,
      amount.toString(),
    );

    resolve([msg]);
  });

  const totalFee = await client.post.simulate(
    subaccount.wallet,
    () => msgs,
    undefined,
    undefined,
  );
  console.log('**Total Fee**');
  console.log(totalFee);

  const amountAfterFee = amount.sub(Long.fromString(totalFee.amount[0].amount));
  console.log('**Amount after fee**');
  console.log(amountAfterFee);

  const tx = await client.post.sendToken(
    subaccount,
    TEST_RECIPIENT_ADDRESS,
    client.config.denoms.CHAINTOKEN_DENOM,
    amountAfterFee.toString(),
    false,
    Method.BroadcastTxCommit,
  );
  console.log('**Send**');
  console.log(tx);
}

test()
  .then(() => {})
  .catch((error) => {
    console.log(error.message);
  });
