import { BECH32_PREFIX, OrderFlags, Order_TimeInForce } from '../src';
import { CompositeClient } from '../src/clients/composite-client';
import {
  Network, OrderSide,
} from '../src/clients/constants';
import LocalWallet from '../src/clients/modules/local-wallet';
import { SubaccountInfo } from '../src/clients/subaccount';
import { randomInt, sleep } from '../src/lib/utils';
import { DYDX_TEST_MNEMONIC, MAX_CLIENT_ID } from './constants';

async function test(): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(DYDX_TEST_MNEMONIC, BECH32_PREFIX);
  console.log(wallet);
  const network = Network.testnet();
  const client = await CompositeClient.connect(network);
  console.log('**Client**');
  console.log(client);
  const subaccount = new SubaccountInfo(wallet, 0);

  const currentBlock = await client.validatorClient.get.latestBlockHeight();
  const nextValidBlockHeight = currentBlock + 1;
  // Note, you can change this to any number between `next_valid_block_height`
  // to `next_valid_block_height + SHORT_BLOCK_WINDOW`
  const goodTilBlock = nextValidBlockHeight + 10;
  const shortTermOrderClientId = randomInt(MAX_CLIENT_ID);
  try {
    // place a short term order
    const tx = await client.placeShortTermOrder(
      subaccount,
      'ETH-USD',
      OrderSide.SELL,
      40000,
      0.01,
      shortTermOrderClientId,
      goodTilBlock,
      Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      false,
    );
    console.log('**Short Term Order Tx**');
    console.log(tx.hash);
  } catch (error) {
    console.log('**Short Term Order Failed**');
    console.log(error.message);
  }

  await sleep(5000);  // wait for placeOrder to complete

  try {
    // cancel the short term order
    const tx = await client.cancelOrder(
      subaccount,
      shortTermOrderClientId,
      OrderFlags.SHORT_TERM,
      'ETH-USD',
      goodTilBlock + 10,
      0,
    );
    console.log('**Cancel Short Term Order Tx**');
    console.log(tx);
  } catch (error) {
    console.log('**Cancel Short Term Order Failed**');
    console.log(error.message);
  }
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
