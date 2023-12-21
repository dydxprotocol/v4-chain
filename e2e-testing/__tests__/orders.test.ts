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
  HeightResponse,
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
    // cancel the order 60 seconds from now
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
    const indexerClient = new IndexerClient(Network.local().indexerConfig);
    const heightResp: HeightResponse = await indexerClient.utility.getHeight();
    const height: number = heightResp.height;
    connectAndValidateSocketClient();

    // place all orders
    for (const order of orderDetails) {
      const modifiedOrder: IPlaceOrder = defaultOrder;
      modifiedOrder.clientId = Math.floor(Math.random() * 1000000000);
      modifiedOrder.goodTilBlock = 0;
      modifiedOrder.clobPairId = order.clobPairId;
      modifiedOrder.timeInForce = order.timeInForce;
      modifiedOrder.reduceOnly = false;
      modifiedOrder.orderFlags = order.orderFlags;
      modifiedOrder.side = order.side;
      modifiedOrder.quantums = Long.fromNumber(order.quantums);
      modifiedOrder.subticks = Long.fromNumber(order.subticks);

      await placeOrder(order.mnemonic, modifiedOrder);
    }

    await utils.sleep(10000);  // wait 10s for orders to be placed & matched
    const [wallet, wallet2] = await Promise.all([
      LocalWallet.fromMnemonic(DYDX_LOCAL_MNEMONIC, BECH32_PREFIX),
      LocalWallet.fromMnemonic(DYDX_LOCAL_MNEMONIC_2, BECH32_PREFIX),
    ]);

    const subaccountId = SubaccountTable.uuid(wallet.address!, 0);
    const subaccountId2 = SubaccountTable.uuid(wallet2.address!, 0);
    const [makerOrders, takerOrders] = await Promise.all([
      OrderTable.findBySubaccountIdAndClobPairAfterHeight(
        subaccountId,
        PERPETUAL_PAIR_BTC_USD.toString(),
        height,
      ),
      OrderTable.findBySubaccountIdAndClobPairAfterHeight(
        subaccountId2,
        PERPETUAL_PAIR_BTC_USD.toString(),
        height,
      ),
    ]);
    expect(makerOrders).toHaveLength(1);
    expect(takerOrders).toHaveLength(1);

    const [makerFills, takerFills] = await Promise.all([
      FillTable.findAll(
        {
          subaccountId: [subaccountId],
          createdOnOrAfterHeight: height.toString(),
        },
        [],
        {},
      ),
      FillTable.findAll(
        {
          subaccountId: [subaccountId2],
          createdOnOrAfterHeight: height.toString(),
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

    // Check API /v4/orders endpoint
    const [ordersResponse, ordersResponse2] = await Promise.all([
      indexerClient.account.getSubaccountOrders(DYDX_LOCAL_ADDRESS,
        0,
        undefined,
        undefined,
        undefined,
        undefined,
        undefined,
        10,
        undefined,
        undefined,
        true),
      indexerClient.account.getSubaccountOrders(DYDX_LOCAL_ADDRESS_2,
        0,
        undefined,
        undefined,
        undefined,
        undefined,
        undefined,
        10,
        undefined,
        undefined,
        true),
    ]);
    expect(ordersResponse[0]).toEqual(
      expect.objectContaining({
        subaccountId: SubaccountTable.uuid(DYDX_LOCAL_ADDRESS, 0),
        clobPairId: '0',
        side: 'BUY',
        size: '0.001',
        totalFilled: '0.0005',
        price: '50000',
        type: 'LIMIT',
        timeInForce: 'GTT',
        reduceOnly: false,
        orderFlags: '64',
        postOnly: false,
        ticker: 'BTC-USD'
      }),
    );
    expect(ordersResponse2[0]).toEqual(
      expect.objectContaining({
        subaccountId: SubaccountTable.uuid(DYDX_LOCAL_ADDRESS_2, 0),
        clobPairId: '0',
        side: 'SELL',
        size: '0.0005',
        totalFilled: '0.0005',
        price: '50000',
        type: 'LIMIT',
        status: 'FILLED',
        timeInForce: 'GTT',
        reduceOnly: false,
        orderFlags: '64',
        postOnly: false,
        ticker: 'BTC-USD'
      }),
    );

    // Check API /v4/perpetualPositions endpoint
    const [response, response2] = await Promise.all([
      indexerClient.account.getSubaccountPerpetualPositions(DYDX_LOCAL_ADDRESS, 0),
      indexerClient.account.getSubaccountPerpetualPositions(DYDX_LOCAL_ADDRESS_2, 0)
    ]);
    expect(response.positions.length).toEqual(1);
    expect(response.positions[0]).toEqual(
      expect.objectContaining({
        market: 'BTC-USD',
        status: 'OPEN',
        side: 'LONG',
        // size: '0.0005',
        // maxSize: '0.0005',
        entryPrice: '50000',
        exitPrice: null,
        realizedPnl: '0',
        unrealizedPnl: '0',
        closedAt: null,
        // sumOpen: '0.0005',
        sumClose: '0',
      }),
    );
    expect(response2.positions.length).toEqual(1);
    expect(response2.positions[0]).toEqual(
      expect.objectContaining({
        market: 'BTC-USD',
        status: 'OPEN',
        side: 'SHORT',
        // size: '-0.0005',
        // maxSize: '-0.0005',
        entryPrice: '50000',
        exitPrice: null,
        realizedPnl: '0',
        unrealizedPnl: '0',
        closedAt: null,
        // sumOpen: '0.0005',
        sumClose: '0',
      }),
    );


    // Check API /v4/orderbooks endpoint
    const orderbooksResponse = await indexerClient.markets.getPerpetualMarketOrderbook('BTC-USD');
    console.log(`orderbooksResponse: ${JSON.stringify(orderbooksResponse)}`);
    expect(orderbooksResponse).toEqual(
      expect.objectContaining({
        bids:[
          {
            price: '50000',
            size: '0.0005',
          }
        ],
        asks:[]
      }),
    );
  });


  function connectAndValidateSocketClient(): void {
    const mySocket = new SocketClient(
      Network.local().indexerConfig,
      () => {
        console.log('open');
      },
      () => {
        console.log('close');
      },
      (message) => {
        if (typeof message.data === 'string') {
          const data = JSON.parse(message.data as string);
          console.log(`data: ${JSON.stringify(data)}`);
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
          } else if (data.type === 'channel_data' && data.contents.perpetualPositions) {
            console.log(`perpetualPositions data: ${JSON.stringify(data)}`);
            expect(data.contents.perpetualPositions[0]).toEqual(
              expect.objectContaining({
                address: DYDX_LOCAL_ADDRESS,
                subaccountNumber: 0,
                market: 'BTC-USD',
                side: 'LONG',
                status: 'OPEN',
                size: '0.0005',
                maxSize: '0.0005',
                netFunding: '0',
                exitPrice: null,
              }),
            );
          } else if (data.type === 'channel_data' && data.contents.fills) {
            console.log(`fills data: ${JSON.stringify(data)}`);
            expect(data.contents.fills[0]).toEqual(
              expect.objectContaining({
                fee: '-0.00275',
                side: 'BUY',
                size: '0.0005',
                type: 'LIMIT',
                price: '50000',
                liquidity: 'MAKER',
                clobPairId: '0',
                quoteAmount: '25',
                subaccountId: SubaccountTable.uuid(DYDX_LOCAL_ADDRESS, 0),
                clientMetadata: '0',
                ticker: 'BTC-USD'
              }),
            );
          }
        }
      },
    );
    mySocket.connect();
  }
});
