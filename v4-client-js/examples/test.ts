import { BECH32_PREFIX } from '../src';
import { CompositeClient } from '../src/clients/composite-client';
import {
  Network, OrderExecution, OrderSide, OrderTimeInForce, OrderType,
} from '../src/clients/constants';
import LocalWallet from '../src/clients/modules/local-wallet';
import { SubaccountInfo } from '../src/clients/subaccount';
import { randomInt } from '../src/lib/utils';
import { DYDX_TEST_MNEMONIC, MAX_CLIENT_ID } from './constants';
import ordersParams from './human_readable_orders.json';

async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function test(): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(DYDX_TEST_MNEMONIC, BECH32_PREFIX);
  console.log(wallet);
  const network = Network.testnet();
  const client = await CompositeClient.connect(network);
  console.log('**Client**');
  console.log(client);
  const subaccount = new SubaccountInfo(wallet, 0);
  for (const orderParams of ordersParams) {
    try {
      const type = OrderType[orderParams.type as keyof typeof OrderType];
      const side = OrderSide[orderParams.side as keyof typeof OrderSide];
      const timeInForceString = orderParams.timeInForce ?? 'GTT';
      const timeInForce = OrderTimeInForce[timeInForceString as keyof typeof OrderTimeInForce];
      const price = orderParams.price ?? 1350;
      const timeInForceSeconds = (timeInForce === OrderTimeInForce.GTT) ? 60 : 0;
      const postOnly = orderParams.postOnly ?? false;
      const tx = await client.placeOrder(
        subaccount,
        'ETH-USD',
        type,
        side,
        price,
        0.01,
        randomInt(MAX_CLIENT_ID),
        timeInForce,
        timeInForceSeconds,
        OrderExecution.DEFAULT,
        postOnly,
        false,
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
