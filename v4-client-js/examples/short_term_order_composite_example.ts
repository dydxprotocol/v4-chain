import { Order_TimeInForce } from '@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order';

import { BECH32_PREFIX } from '../src';
import { CompositeClient } from '../src/clients/composite-client';
import {
  Network, OrderExecution, OrderSide,
} from '../src/clients/constants';
import LocalWallet from '../src/clients/modules/local-wallet';
import { SubaccountInfo } from '../src/clients/subaccount';
import { randomInt } from '../src/lib/utils';
import { DYDX_TEST_MNEMONIC } from './constants';
import ordersParams from './human_readable_short_term_orders.json';

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
      const side = OrderSide[orderParams.side as keyof typeof OrderSide];
      const price = orderParams.price ?? 1350;

      const currentBlock = await client.validatorClient.get.latestBlockHeight();
      const nextValidBlockHeight = currentBlock + 1;
      // Note, you can change this to any number between `next_valid_block_height`
      // to `next_valid_block_height + SHORT_BLOCK_WINDOW`
      const goodTilBlock = nextValidBlockHeight + 10;

      const timeInForce = orderExecutionToTimeInForce(orderParams.timeInForce);

      // uint32
      const clientId = randomInt(2 ** 32 - 1);

      const tx = await client.placeShortTermOrder(
        subaccount,
        'ETH-USD',
        side,
        price,
        0.01,
        clientId,
        goodTilBlock,
        timeInForce,
        false,
      );
      console.log('**Order Tx**');
      console.log(tx.hash.toString());
    } catch (error) {
      console.log(error.message);
    }

    await sleep(5000);  // wait for placeOrder to complete
  }
}

function orderExecutionToTimeInForce(orderExecution: string): Order_TimeInForce {
  switch (orderExecution) {
    case OrderExecution.DEFAULT:
      return Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED;
    case OrderExecution.FOK:
      return Order_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL;
    case OrderExecution.IOC:
      return Order_TimeInForce.TIME_IN_FORCE_IOC;
    case OrderExecution.POST_ONLY:
      return Order_TimeInForce.TIME_IN_FORCE_POST_ONLY;
    default:
      throw new Error('Unrecognized order execution');
  }
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
