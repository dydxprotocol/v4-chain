import { EncodeObject } from '@cosmjs/proto-signing';
import Long from 'long';

import { Network } from '../src/clients/constants';
import LocalWallet from '../src/clients/modules/local-wallet';
import { NobleClient } from '../src/clients/noble-client';
import { ValidatorClient } from '../src/clients/validator-client';
import { BECH32_PREFIX, NOBLE_BECH32_PREFIX } from '../src/lib/constants';
import { sleep } from '../src/lib/utils';
import { DYDX_TEST_MNEMONIC } from './constants';

async function test(): Promise<void> {
  const dydxClient = await ValidatorClient.connect(
    Network.testnet().validatorConfig,
  );

  const dydxWallet = await LocalWallet.fromMnemonic(
    DYDX_TEST_MNEMONIC,
    BECH32_PREFIX,
  );
  const nobleWallet = await LocalWallet.fromMnemonic(
    DYDX_TEST_MNEMONIC,
    NOBLE_BECH32_PREFIX,
  );

  const client = new NobleClient('https://rpc.testnet.noble.strange.love');
  await client.connect(nobleWallet);

  if (nobleWallet.address === undefined || dydxWallet.address === undefined) {
    throw new Error('Wallet not found');
  }

  // IBC to noble

  const ibcToNobleMsg: EncodeObject = {
    typeUrl: '/ibc.applications.transfer.v1.MsgTransfer',
    value: {
      sourcePort: 'transfer',
      sourceChannel: 'channel-0',
      token: {
        denom:
          'ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5',
        amount: '1000000',
      },
      sender: dydxWallet.address,
      receiver: nobleWallet.address,
      timeoutTimestamp: Long.fromNumber(
        Math.floor(Date.now() / 1000) * 1e9 + 10 * 60 * 1e9,
      ),
    },
  };

  const msgs = [ibcToNobleMsg];
  const encodeObjects: Promise<EncodeObject[]> = new Promise((resolve) => resolve(msgs),
  );

  await dydxClient.post.send(
    dydxWallet,
    () => {
      return encodeObjects;
    },
    false,
    undefined,
    undefined,
  );

  await sleep(30000);

  try {
    const coins = await client.getAccountBalances();
    console.log('Balances');
    console.log(JSON.stringify(coins));

    // IBC from noble

    const ibcFromNobleMsg: EncodeObject = {
      typeUrl: '/ibc.applications.transfer.v1.MsgTransfer',
      value: {
        sourcePort: 'transfer',
        sourceChannel: 'channel-21',
        token: {
          denom: 'uusdc',
          amount: coins[0].amount,
        },
        sender: nobleWallet.address,
        receiver: dydxWallet.address,
        timeoutTimestamp: Long.fromNumber(
          Math.floor(Date.now() / 1000) * 1e9 + 10 * 60 * 1e9,
        ),
      },
    };
    const fee = await client.simulateTransaction([ibcFromNobleMsg]);

    ibcFromNobleMsg.value.token.amount = (parseInt(ibcFromNobleMsg.value.token.amount, 10) -
      Math.floor(parseInt(fee.amount[0].amount, 10) * 1.4)).toString();

    await client.send([ibcFromNobleMsg]);
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  await sleep(30000);

  try {
    const coin = await client.getAccountBalance('uusdc');
    console.log('Balance');
    console.log(JSON.stringify(coin));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }
}

test()
  .then(() => {})
  .catch((error) => {
    console.log(error.message);
    console.log(error);
  });
