import Long from 'long';
import {
  BECH32_PREFIX,
  IndexerClient,
  IPlaceOrder,
  LocalWallet,
  Network,
  Order_Side,
  Order_TimeInForce,
  OrderFlags, SocketClient,
  SubaccountInfo,
  ValidatorClient,
} from '@dydxprotocol/v4-client-js';
import {
  DYDX_LOCAL_ADDRESS,
  DYDX_LOCAL_ADDRESS_2,
  DYDX_LOCAL_MNEMONIC,
  DYDX_LOCAL_MNEMONIC_2,
} from './helpers/constants';
import * as utils from './helpers/utils';
import {
  FillTable, FillType, Liquidity, OrderSide, OrderTable, SubaccountTable,
} from '@dydxprotocol-indexer/postgres';

const PERPETUAL_PAIR_BTC_USD: number = 0;
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
    timeInForce: 0,
    orderFlags: 64,
    side: 1,
    clobPairId: PERPETUAL_PAIR_BTC_USD,
    quantums: 10000000,
    subticks: 5000000000,
  },
  {
    mnemonic: DYDX_LOCAL_MNEMONIC_2,
    timeInForce: 0,
    orderFlags: 64,
    side: 2,
    clobPairId: PERPETUAL_PAIR_BTC_USD,
    quantums: 5000000,
    subticks: 5000000000,
  },
];

async function placeOrder(
  mnemonic: string,
  order: IPlaceOrder,
): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(mnemonic, BECH32_PREFIX);
  const client = await ValidatorClient.connect(Network.local().validatorConfig);

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
    const indexerClient = new IndexerClient(Network.local().indexerConfig);

    await utils.sleep(5000);  // wait 5s for orders to be placed & matched
    const [wallet, wallet2] = await Promise.all([
      LocalWallet.fromMnemonic(DYDX_LOCAL_MNEMONIC, BECH32_PREFIX),
      LocalWallet.fromMnemonic(DYDX_LOCAL_MNEMONIC_2, BECH32_PREFIX),
    ]);

    const subaccountId = SubaccountTable.uuid(wallet.address!, 0);
    const subaccountId2 = SubaccountTable.uuid(wallet2.address!, 0);
    const [makerOrders, takerOrders] = await Promise.all([
      OrderTable.findBySubaccountIdAndClobPair(subaccountId, PERPETUAL_PAIR_BTC_USD.toString()),
      OrderTable.findBySubaccountIdAndClobPair(subaccountId2, PERPETUAL_PAIR_BTC_USD.toString()),
    ]);
    expect(makerOrders).toHaveLength(1);
    expect(takerOrders).toHaveLength(1);

    const [makerFills, takerFills] = await Promise.all([
      FillTable.findAll(
        {
          subaccountId: [subaccountId],
        },
        [],
        {},
      ),
      FillTable.findAll(
        {
          subaccountId: [subaccountId2],
        },
        [],
        {},
      ),
    ]);

    expect(makerFills.length).toEqual(1);
    expect(makerFills[0]).toEqual(expect.objectContaining({
      subaccountId,
      side: OrderSide.BUY,
      liquidity: Liquidity.MAKER,
      type: FillType.LIMIT,
      clobPairId: '0',
      orderId: makerOrders[0].id,
      size: '0.0005',
      price: '50000',
      quoteAmount: '25',
      clientMetadata: '0',
      fee: '-0.00275',
    }));

    expect(takerFills.length).toEqual(1);
    expect(takerFills[0]).toEqual(expect.objectContaining({
      subaccountId: subaccountId2,
      side: OrderSide.SELL,
      liquidity: Liquidity.TAKER,
      type: FillType.LIMIT,
      clobPairId: '0',
      orderId: takerOrders[0].id,
      size: '0.0005',
      price: '50000',
      quoteAmount: '25',
      clientMetadata: '0',
      fee: '0.0125',
    }));

    // Check API /v4/perpetualPositions endpoint
    let response = await indexerClient.account.getSubaccountPerpetualPositions(
      DYDX_LOCAL_ADDRESS,
      0,
    );
    expect(response).not.toBeNull();
    let positions = response.positions;
    expect(positions.length).toEqual(1);
    let position: any = positions[0];
    expect(position).toEqual(
      expect.objectContaining({
        market: 'BTC-USD',
        status: 'OPEN',
        side: 'LONG',
        size: '0.0005',
        maxSize: '0.0005',
        entryPrice: '50000',
        exitPrice: null,
        realizedPnl: '0',
        unrealizedPnl: '0',
        closedAt: null,
        sumOpen: '0.0005',
        sumClose: '0',
      }),
    );
    response = await indexerClient.account.getSubaccountPerpetualPositions(
      DYDX_LOCAL_ADDRESS_2,
      0,
    );
    expect(response).not.toBeNull();
    positions = response.positions;
    expect(positions.length).toEqual(1);
    position = positions[0];
    expect(position).toEqual(
      expect.objectContaining({
        market: 'BTC-USD',
        status: 'OPEN',
        side: 'SHORT',
        size: '-0.0005',
        maxSize: '-0.0005',
        entryPrice: '50000',
        exitPrice: null,
        realizedPnl: '0',
        unrealizedPnl: '0',
        closedAt: null,
        sumOpen: '0.0005',
        sumClose: '0',
      }),
    );
  });

  function connectAndValidateSocketClient(): void {
    const mySocket = new SocketClient(
      Network.local().indexerConfig,
      () => {
      },
      () => {
      },
      (message) => {
        if (typeof message.data === 'string') {
          const data = JSON.parse(message.data as string);
          if (data.type === 'connected') {
            mySocket.subscribeToSubaccount(DYDX_LOCAL_ADDRESS, 0);
          } else if (data.type === 'subscribed') {
            expect(data.channel).toEqual('v4_subaccounts');
            expect(data.id).toEqual(`${DYDX_LOCAL_ADDRESS}/0`);
            expect(data.contents.subaccount).toEqual(
              expect.objectContaining({
                address: DYDX_LOCAL_ADDRESS,
                subaccountNumber: 0,
              }),
            );
          } else if (data.type === 'channel_data' && data.contents.transfers) {
            expect(data.contents.transfers).toEqual(
              expect.objectContaining({
                sender: {
                  address: DYDX_LOCAL_ADDRESS,
                },
                recipient: {
                  address: DYDX_LOCAL_ADDRESS,
                  subaccountNumber: 0,
                },
                size: '10',
                symbol: 'USDC',
                type: 'DEPOSIT',
              }),
            );
          }
        }
      },
    );
    mySocket.connect();
  }
});
