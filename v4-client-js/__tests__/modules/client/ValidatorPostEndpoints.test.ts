import { Order_Side } from '@dydxprotocol/v4-proto/src/codegen/dydxprotocol/clob/order';

import { Network } from '../../../src/clients/constants';
import LocalWallet from '../../../src/clients/modules/local-wallet';
import { SubaccountInfo } from '../../../src/clients/subaccount';
import { IPlaceOrder } from '../../../src/clients/types';
import { ValidatorClient } from '../../../src/clients/validator-client';
import { randomInt } from '../../../src/lib/utils';
import { DYDX_TEST_MNEMONIC, defaultOrder } from '../../../examples/constants';
import { BECH32_PREFIX } from '../../../src';

function dummyOrder(height: number): IPlaceOrder {
  const placeOrder = defaultOrder;
  placeOrder.clientId = randomInt(1000000000);
  placeOrder.goodTilBlock = height + 3;
  // placeOrder.goodTilBlockTime = height + 3;
  const random = randomInt(1000);
  if ((random % 2) === 0) {
    placeOrder.side = Order_Side.SIDE_BUY;
  } else {
    placeOrder.side = Order_Side.SIDE_SELL;
  }
  return placeOrder;
}

describe('Validator Client', () => {
  let wallet: LocalWallet;
  let client: ValidatorClient;

  describe('Post', () => {
    beforeEach(async () => {
      wallet = await LocalWallet.fromMnemonic(DYDX_TEST_MNEMONIC, BECH32_PREFIX);
      client = await ValidatorClient.connect(Network.testnet().validatorConfig);
    });

    it('PlaceOrder', async () => {
      console.log('**Client**');
      console.log(client);
      const address = wallet.address!;
      const account = await client.get.getAccount(address);
      console.log('**Account**');
      console.log(account);
      const height = await client.get.latestBlockHeight();
      const subaccount = new SubaccountInfo(wallet, 0);
      const placeOrder = dummyOrder(height);
      placeOrder.clientId = randomInt(1_000_000_000);
      const tx = await client.post.placeOrderObject(
        subaccount,
        placeOrder,
      );
      console.log('**Order Tx**');
      console.log(tx);
    });
  });
});
