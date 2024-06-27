import { BECH32_PREFIX, OrderFlags } from '../src';
import { CompositeClient } from '../src/clients/composite-client';
import {
  Network, OrderExecution, OrderSide, OrderTimeInForce, OrderType,
} from '../src/clients/constants';
import LocalWallet from '../src/clients/modules/local-wallet';
import { SubaccountInfo } from '../src/clients/subaccount';
import { randomInt } from '../src/lib/utils';
import { DYDX_TEST_MNEMONIC, MAX_CLIENT_ID } from './constants';

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

  /*
  Note this example places a stateful order.
  Programmatic traders should generally not use stateful orders for following reasons:
  - Stateful orders received out of order by validators will fail sequence number validation
    and be dropped.
  - Stateful orders have worse time priority since they are only matched after they are included
    on the block.
  - Stateful order rate limits are more restrictive than Short-Term orders, specifically max 2 per
    block / 20 per 100 blocks.
  - Stateful orders can only be canceled after they’ve been included in a block.
  */
  const longTermOrderClientId = randomInt(MAX_CLIENT_ID);
  try {
    // place a long term order
    const tx = await client.placeOrder(
      subaccount,
      'ETH-USD',
      OrderType.LIMIT,
      OrderSide.SELL,
      40000,
      0.01,
      longTermOrderClientId,
      OrderTimeInForce.GTT,
      60,
      OrderExecution.DEFAULT,
      false,
      false,
    );
    console.log('**Long Term Order Tx**');
    console.log(tx.hash);
  } catch (error) {
    console.log('**Long Term Order Failed**');
    console.log(error.message);
  }

  await sleep(5000);  // wait for placeOrder to complete

  try {
    // cancel the long term order
    const tx = await client.cancelOrder(
      subaccount,
      longTermOrderClientId,
      OrderFlags.LONG_TERM,
      'ETH-USD',
      0,
      120,
    );
    console.log('**Cancel Long Term Order Tx**');
    console.log(tx);
  } catch (error) {
    console.log('**Cancel Long Term Order Failed**');
    console.log(error.message);
  }
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
