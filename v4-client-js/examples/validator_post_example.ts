import Long from 'long';
import protobuf from 'protobufjs';

import { BECH32_PREFIX } from '../src';
import { Network } from '../src/clients/constants';
import LocalWallet from '../src/clients/modules/local-wallet';
import { SubaccountInfo } from '../src/clients/subaccount';
import { IPlaceOrder } from '../src/clients/types';
import { ValidatorClient } from '../src/clients/validator-client';
import { randomInt } from '../src/lib/utils';
import { DYDX_TEST_MNEMONIC, defaultOrder } from './constants';
import ordersParams from './raw_orders.json';

// Required for encoding and decoding queries that are of type Long.
// Must be done once but since the individal modules should be usable
// - must be set in each module that encounters encoding/decoding Longs.
// Reference: https://github.com/protobufjs/protobuf.js/issues/921
protobuf.util.Long = Long;
protobuf.configure();

function dummyOrder(height: number): IPlaceOrder {
  const placeOrder = defaultOrder;
  placeOrder.clientId = randomInt(1000000000);
  placeOrder.goodTilBlock = height + 3;
  return placeOrder;
}

async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function test(): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(DYDX_TEST_MNEMONIC, BECH32_PREFIX);
  console.log(wallet);

  const client = await ValidatorClient.connect(Network.testnet().validatorConfig);
  console.log('**Client**');
  console.log(client);

  const value1 = Long.fromNumber(400000000000);
  console.log(value1.toString());

  const subaccount = new SubaccountInfo(wallet, 0);
  for (const orderParams of ordersParams) {
    const height = await client.get.latestBlockHeight();
    const placeOrder = dummyOrder(height);

    placeOrder.timeInForce = orderParams.timeInForce;
    placeOrder.reduceOnly = false; // reduceOnly is currently disabled
    placeOrder.orderFlags = orderParams.orderFlags;
    placeOrder.side = orderParams.side;
    placeOrder.quantums = Long.fromNumber(orderParams.quantums);
    placeOrder.subticks = Long.fromNumber(orderParams.subticks);
    try {
      if (placeOrder.orderFlags !== 0) {
        placeOrder.goodTilBlock = 0;
        const now = new Date();
        const millisecondsPerSecond = 1000;
        const interval = 60 * millisecondsPerSecond;
        const future = new Date(now.valueOf() + interval);
        placeOrder.goodTilBlockTime = Math.round(future.getTime() / 1000);
      } else {
        placeOrder.goodTilBlockTime = 0;
      }

      const tx = await client.post.placeOrderObject(
        subaccount,
        placeOrder,
      );
      console.log('**Order Tx**');
      console.log(tx);
    } catch (error) {
      console.log(error.message);
    }

    await sleep(5000);  // wait for placeOrder to complete
  }
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
