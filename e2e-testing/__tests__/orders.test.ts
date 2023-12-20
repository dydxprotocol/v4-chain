import Long from 'long';
import {
  BECH32_PREFIX,
  IPlaceOrder,
  LocalWallet,
  Network,
  Order_Side,
  Order_TimeInForce,
  OrderFlags,
  SubaccountInfo,
  ValidatorClient,
} from '@dydxprotocol/v4-client-js';
import { DYDX_LOCAL_MNEMONIC, DYDX_LOCAL_MNEMONIC_2 } from './helpers/constants';

const PERPETUAL_PAIR_BTC_USD: number = 0;
const PERPETUAL_PAIR_ETH_USD: number = 1;
const quantums: Long = new Long(1_000_000_000);
const subticks: Long = new Long(1_000_000_000);

const defaultOrder: IPlaceOrder = {
  clientId: 0,
  orderFlags: OrderFlags.SHORT_TERM,
  clobPairId: PERPETUAL_PAIR_BTC_USD,
  side: Order_Side.SIDE_BUY,
  quantums,
  subticks,
  timeInForce: Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
  reduceOnly: false,
  clientMetadata: 0,
};

type OrderDetails = {
  mnemonic: string;
  timeInForce: number;
  orderFlags: number;
  side: number;
  clobPairId: number;
  quantums: number;
  subticks: number;
};

const orderDetails: OrderDetails[] = [
  {
    mnemonic: DYDX_LOCAL_MNEMONIC,
    timeInForce: 2,
    orderFlags: 64,
    side: 1,
    clobPairId: PERPETUAL_PAIR_BTC_USD,
    quantums: 10000000,
    subticks: 40000000000,
  },
  {
    mnemonic: DYDX_LOCAL_MNEMONIC_2,
    timeInForce: 2,
    orderFlags: 64,
    side: 1,
    clobPairId: PERPETUAL_PAIR_ETH_USD,
    quantums: 10000000,
    subticks: 40000000000,
  },
];

async function placeOrder(
  mnemonic: string,
  order: IPlaceOrder,
): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(mnemonic, BECH32_PREFIX);
  const client = await ValidatorClient.connect(Network.testnet().validatorConfig);

  const subaccount = new SubaccountInfo(wallet, 0);
  const modifiedOrder: IPlaceOrder = order;
  if (order.orderFlags !== 0) {
    modifiedOrder.goodTilBlock = 0;
    const now = new Date();
    const millisecondsPerSecond = 1000;
    const interval = 60 * millisecondsPerSecond;
    const future = new Date(now.valueOf() + interval);
    modifiedOrder.goodTilBlockTime = Math.round(future.getTime() / 1000);
  } else {
    modifiedOrder.goodTilBlockTime = 0;
  }

  await client.post.placeOrderObject(
    subaccount,
    modifiedOrder,
  );
}

describe('orders', () => {
  it('test orders', async () => {
    // place all orders
    for (const order of orderDetails) {
      const modifiedOrder: IPlaceOrder = defaultOrder;
      modifiedOrder.clientId = Math.floor(Math.random() * 1000000000);
      modifiedOrder.goodTilBlock = 0;
      modifiedOrder.clobPairId = order.clobPairId;
      modifiedOrder.timeInForce = order.timeInForce;
      modifiedOrder.reduceOnly = false; // reduceOnly is currently disabled
      modifiedOrder.orderFlags = order.orderFlags;
      modifiedOrder.side = order.side;
      modifiedOrder.quantums = Long.fromNumber(order.quantums);
      modifiedOrder.subticks = Long.fromNumber(order.subticks);

      await placeOrder(order.mnemonic, modifiedOrder);
    }

  });
});
