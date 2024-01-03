import {
  BECH32_PREFIX,
  HeightResponse,
  IndexerClient,
  IPlaceOrder,
  LocalWallet,
  Network,
  SocketClient,
  SubaccountInfo,
  ValidatorClient,
} from '@dydxprotocol/v4-client-js';
import {
  DYDX_LOCAL_ADDRESS,
  DYDX_LOCAL_ADDRESS_2,
  DYDX_LOCAL_MNEMONIC,
  DYDX_LOCAL_MNEMONIC_2,
  orderDetails,
  PERPETUAL_PAIR_BTC_USD,
} from './helpers/constants';
import { DateTime } from 'luxon';
import * as utils from './helpers/utils';
import { connectAndValidateSocketClient, createModifiedOrder } from './helpers/utils';
import {
  CandleResolution,
  FillTable,
  FillType,
  helpers,
  IsoString,
  Liquidity,
  OrderSide,
  OrderTable,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';

async function placeOrder(
  mnemonic: string,
  order: IPlaceOrder,
): Promise<void> {
  const wallet = await LocalWallet.fromMnemonic(mnemonic, BECH32_PREFIX);
  const client = await ValidatorClient.connect(Network.local().validatorConfig);

  const subaccount = new SubaccountInfo(wallet, 0);
  const modifiedOrder: IPlaceOrder = order;
  if (order.orderFlags !== 0) {
    // cancel the order 30 seconds from now
    modifiedOrder.goodTilBlock = 0;
    const now = new Date();
    const millisecondsPerSecond = 1000;
    const interval = 30 * millisecondsPerSecond;
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
    connectAndValidateSocketClient(validateOrders);

    // place all orders
    for (const order of orderDetails) {
      const modifiedOrder: IPlaceOrder = createModifiedOrder(order);

      await placeOrder(order.mnemonic, modifiedOrder);
    }
    const candleStart: IsoString = helpers.calculateNormalizedCandleStartTime(
      DateTime.utc(),
      CandleResolution.ONE_MINUTE,
    ).toISO();

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
        ticker: 'BTC-USD',
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
        ticker: 'BTC-USD',
      }),
    );

    // Check API /v4/perpetualPositions endpoint
    const [response, response2] = await Promise.all([
      indexerClient.account.getSubaccountPerpetualPositions(DYDX_LOCAL_ADDRESS, 0),
      indexerClient.account.getSubaccountPerpetualPositions(DYDX_LOCAL_ADDRESS_2, 0),
    ]);
    expect(response.positions.length).toEqual(1);
    expect(response.positions[0]).toEqual(
      expect.objectContaining({
        market: 'BTC-USD',
        status: 'OPEN',
        side: 'LONG',
        entryPrice: '50000',
      }),
    );
    expect(response2.positions.length).toEqual(1);
    expect(response2.positions[0]).toEqual(
      expect.objectContaining({
        market: 'BTC-USD',
        status: 'OPEN',
        side: 'SHORT',
        entryPrice: '50000',
      }),
    );

    // Check API /v4/orderbooks endpoint
    const orderbooksResponse = await indexerClient.markets.getPerpetualMarketOrderbook('BTC-USD');
    expect(orderbooksResponse).toEqual(
      expect.objectContaining({
        bids: [
          {
            price: '50000',
            size: '0.0005',
          },
        ],
        asks: [],
      }),
    );

    // Check API /v4/candles endpoint
    const candlesResponse = await indexerClient.markets.getPerpetualMarketCandles(
      'BTC-USD',
      CandleResolution.ONE_MINUTE,
      undefined,
      undefined,
      1,
    );
    expect(candlesResponse.candles[0]).toEqual(
      expect.objectContaining({
        startedAt: candleStart,
        ticker: 'BTC-USD',
        resolution: '1MIN',
        low: '50000',
        high: '50000',
        open: '50000',
        close: '50000',
        baseTokenVolume: '0.0005',
        usdVolume: '25',
        trades: 1,
      }),
    );
  });

  function validateOrders(data: any, socketClient: SocketClient): void {
    if (data.type === 'connected') {
      socketClient.subscribeToSubaccount(DYDX_LOCAL_ADDRESS, 0);
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
      expect(data.contents.perpetualPositions[0]).toEqual(
        expect.objectContaining({
          address: DYDX_LOCAL_ADDRESS,
          subaccountNumber: 0,
          market: 'BTC-USD',
          side: 'LONG',
          status: 'OPEN',
          netFunding: '0',
          exitPrice: null,
        }),
      );
    } else if (data.type === 'channel_data' && data.contents.fills) {
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
          ticker: 'BTC-USD',
        }),
      );
    }
  }
});
