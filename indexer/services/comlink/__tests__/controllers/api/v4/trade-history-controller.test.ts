import {
  dbHelpers,
  testMocks,
  testConstants,
  OrderTable,
  FillTable,
  OrderFromDatabase,
  FillFromDatabase,
  perpetualMarketRefresher,
  OrderSide,
  OrderType,
  FillType,
  Liquidity,
} from '@dydxprotocol-indexer/postgres';

import { MarketType, RequestMethod, TradeHistoryType } from '../../../../src/types';
import { getQueryString, sendRequest } from '../../../helpers/helpers';

describe('trade-history-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET /', () => {
    const defaultSubaccountNumber: number = testConstants.defaultSubaccount.subaccountNumber;
    const defaultAddress: string = testConstants.defaultSubaccount.address;
    const defaultMarket: string = testConstants.defaultPerpetualMarket.ticker;

    beforeEach(async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
    });

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('returns trade history for a single OPEN fill', async () => {
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory?address=${defaultAddress}` +
          `&subaccountNumber=${defaultSubaccountNumber}`,
      });

      expect(response.body.tradeHistory).toHaveLength(1);
      expect(response.body.tradeHistory[0]).toEqual(
        expect.objectContaining({
          action: TradeHistoryType.OPEN,
          side: OrderSide.BUY,
          marketId: defaultMarket,
          prevSize: '0',
          orderId: testConstants.defaultOrderId,
          orderType: OrderType.LIMIT,
        }),
      );
      expect(response.body.totalResults).toBe(1);
      expect(response.body.pageSize).toBeDefined();
      expect(response.body.offset).toBe(0);
    });

    it('returns OPEN then CLOSE for buy-then-sell of same size', async () => {
      // Create buy order + fill (OPEN)
      const buyOrder: OrderFromDatabase = await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create({
        ...testConstants.defaultFill,
        orderId: buyOrder.id,
        side: OrderSide.BUY,
        size: '5',
        price: '20000',
        quoteAmount: '100000',
        createdAt: '2023-01-22T00:00:00.000Z',
        createdAtHeight: '1',
      });

      // Create sell order + fill (CLOSE)
      const sellOrder: OrderFromDatabase = await OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '2',
        side: OrderSide.SELL,
      });
      const sellFill: FillFromDatabase = await FillTable.create({
        ...testConstants.defaultFill,
        orderId: sellOrder.id,
        side: OrderSide.SELL,
        size: '5',
        price: '21000',
        quoteAmount: '105000',
        eventId: testConstants.defaultTendermintEventId2,
        createdAt: '2023-01-22T00:01:00.000Z',
        createdAtHeight: '2',
      });

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory?address=${defaultAddress}` +
          `&subaccountNumber=${defaultSubaccountNumber}`,
      });

      // Should be sorted most recent first: CLOSE, then OPEN
      expect(response.body.tradeHistory).toHaveLength(2);
      expect(response.body.tradeHistory[0].action).toBe(TradeHistoryType.CLOSE);
      expect(response.body.tradeHistory[0].side).toBe(OrderSide.SELL);
      expect(response.body.tradeHistory[0].orderId).toBe(sellOrder.id);
      // PnL: (21000 - 20000) * 5 = 5000
      expect(response.body.tradeHistory[0].netRealizedPnl).toBe('5000');

      expect(response.body.tradeHistory[1].action).toBe(TradeHistoryType.OPEN);
      expect(response.body.tradeHistory[1].side).toBe(OrderSide.BUY);
      expect(response.body.tradeHistory[1].netRealizedPnl).toBe('0');
    });

    it('filters by market when market and marketType are provided', async () => {
      // BTC-USD order + fill
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      // ETH-USD order + fill
      const ethOrder: OrderFromDatabase = await OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '3',
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
      });
      await FillTable.create({
        ...testConstants.defaultFill,
        orderId: ethOrder.id,
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
        eventId: testConstants.defaultTendermintEventId2,
      });

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory?address=${defaultAddress}` +
          `&subaccountNumber=${defaultSubaccountNumber}` +
          `&market=${testConstants.defaultPerpetualMarket2.ticker}` +
          `&marketType=${MarketType.PERPETUAL}`,
      });

      expect(response.body.tradeHistory).toHaveLength(1);
      expect(response.body.tradeHistory[0].marketId).toBe(
        testConstants.defaultPerpetualMarket2.ticker,
      );
    });

    it('returns empty trade history when no fills exist', async () => {
      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory?address=${defaultAddress}` +
          `&subaccountNumber=${defaultSubaccountNumber}`,
      });

      expect(response.body.tradeHistory).toEqual([]);
      expect(response.body.totalResults).toBe(0);
    });

    it('paginates results correctly', async () => {
      // Create buy order + fill (OPEN)
      const buyOrder: OrderFromDatabase = await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create({
        ...testConstants.defaultFill,
        orderId: buyOrder.id,
        side: OrderSide.BUY,
        size: '5',
        price: '20000',
        quoteAmount: '100000',
        createdAt: '2023-01-22T00:00:00.000Z',
        createdAtHeight: '1',
      });

      // Create sell order + fill (CLOSE)
      const sellOrder: OrderFromDatabase = await OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '2',
        side: OrderSide.SELL,
      });
      await FillTable.create({
        ...testConstants.defaultFill,
        orderId: sellOrder.id,
        side: OrderSide.SELL,
        size: '5',
        price: '21000',
        quoteAmount: '105000',
        eventId: testConstants.defaultTendermintEventId2,
        createdAt: '2023-01-22T00:01:00.000Z',
        createdAtHeight: '2',
      });

      // Page 1 with limit 1
      const page1 = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory?address=${defaultAddress}` +
          `&subaccountNumber=${defaultSubaccountNumber}&page=1&limit=1`,
      });

      expect(page1.body.tradeHistory).toHaveLength(1);
      expect(page1.body.pageSize).toBe(1);
      expect(page1.body.totalResults).toBe(2);
      expect(page1.body.offset).toBe(0);

      // Page 2 with limit 1
      const page2 = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory?address=${defaultAddress}` +
          `&subaccountNumber=${defaultSubaccountNumber}&page=2&limit=1`,
      });

      expect(page2.body.tradeHistory).toHaveLength(1);
      expect(page2.body.pageSize).toBe(1);
      expect(page2.body.totalResults).toBe(2);
      expect(page2.body.offset).toBe(1);
    });

    it('returns 404 for unknown market', async () => {
      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory?address=${defaultAddress}` +
          `&subaccountNumber=${defaultSubaccountNumber}` +
          `&market=UNKNOWN&marketType=${MarketType.PERPETUAL}`,
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: `UNKNOWN not found in markets of type ${MarketType.PERPETUAL}`,
          },
        ],
      });
    });

    it.each([
      [
        'market without marketType',
        {
          address: testConstants.defaultAddress,
          subaccountNumber: 0,
          market: 'BTC-USD',
        },
        'marketType',
        'marketType must be provided if market is provided',
      ],
      [
        'marketType without market',
        {
          address: testConstants.defaultAddress,
          subaccountNumber: 0,
          marketType: MarketType.PERPETUAL,
        },
        'market',
        'market must be provided if marketType is provided',
      ],
      [
        'invalid marketType',
        {
          address: testConstants.defaultAddress,
          subaccountNumber: 0,
          marketType: 'INVALID',
          market: 'BTC-USD',
        },
        'marketType',
        'marketType must be a valid market type (PERPETUAL/SPOT)',
      ],
    ])('returns 400 when validation fails: %s', async (
      _reason: string,
      queryParams: Record<string, string | number>,
      fieldWithError: string,
      expectedErrorMsg: string,
    ) => {
      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory?${getQueryString(queryParams)}`,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: fieldWithError,
            msg: expectedErrorMsg,
          }),
        ]),
      }));
    });

    it('returns liquidation trade with null orderType', async () => {
      // Create a fill without an orderId (liquidation)
      await FillTable.create({
        ...testConstants.defaultFill,
        orderId: undefined,
        side: OrderSide.SELL,
        size: '5',
        price: '20000',
        quoteAmount: '100000',
        type: FillType.LIQUIDATION,
        liquidity: Liquidity.TAKER,
      });

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory?address=${defaultAddress}` +
          `&subaccountNumber=${defaultSubaccountNumber}`,
      });

      expect(response.body.tradeHistory).toHaveLength(1);
      expect(response.body.tradeHistory[0].orderType).toBeNull();
      expect(response.body.tradeHistory[0].orderId).toBeNull();
    });
  });

  describe('GET /parentSubaccountNumber', () => {
    beforeEach(async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
    });

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('returns trade history for parent subaccount including child subaccounts', async () => {
      // Fill for default subaccount (subaccountNumber=0)
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      // Fill for isolated subaccount (subaccountNumber=128, child of parent 0)
      await OrderTable.create(testConstants.isolatedMarketOrder);
      await FillTable.create(testConstants.isolatedMarketFill);

      const parentSubaccountNumber = 0;
      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
          `&parentSubaccountNumber=${parentSubaccountNumber}`,
      });

      // Both fills should be returned (from subaccount 0 and subaccount 128)
      expect(response.body.tradeHistory.length).toBeGreaterThanOrEqual(2);
      expect(response.body.totalResults).toBeGreaterThanOrEqual(2);
    });

    it('paginates parent subaccount results', async () => {
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      await OrderTable.create(testConstants.isolatedMarketOrder);
      await FillTable.create(testConstants.isolatedMarketFill);

      const parentSubaccountNumber = 0;

      const page1 = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
          `&parentSubaccountNumber=${parentSubaccountNumber}&page=1&limit=1`,
      });

      expect(page1.body.tradeHistory).toHaveLength(1);
      expect(page1.body.pageSize).toBe(1);
      expect(page1.body.totalResults).toBeGreaterThanOrEqual(2);
    });

    it.each([
      [
        'market without marketType',
        {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          market: 'BTC-USD',
        },
        'marketType',
        'marketType must be provided if market is provided',
      ],
      [
        'marketType without market',
        {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          marketType: MarketType.PERPETUAL,
        },
        'market',
        'market must be provided if marketType is provided',
      ],
    ])('returns 400 when validation fails for parentSubaccount: %s', async (
      _reason: string,
      queryParams: Record<string, string | number>,
      fieldWithError: string,
      expectedErrorMsg: string,
    ) => {
      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/tradeHistory/parentSubaccountNumber?${getQueryString(queryParams)}`,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: fieldWithError,
            msg: expectedErrorMsg,
          }),
        ]),
      }));
    });
  });
});
