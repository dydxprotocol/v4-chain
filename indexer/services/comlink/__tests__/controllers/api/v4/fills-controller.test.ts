import {
  dbHelpers,
  testMocks,
  testConstants,
  OrderTable,
  FillTable,
  OrderFromDatabase,
  perpetualMarketRefresher,
  FillFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { FillResponseObject, MarketType, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import {
  getQueryString,
  sendRequest,
  fillResponseObjectFromFillCreateObject,
} from '../../../helpers/helpers';

describe('fills-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET', () => {
    const defaultSubaccountNumber: number = testConstants.defaultSubaccount.subaccountNumber;
    const defaultAddress: string = testConstants.defaultSubaccount.address;
    const defaultMarket: string = testConstants.defaultPerpetualMarket.ticker;
    const invalidMarket: string = 'UNKNOWN';

    beforeEach(async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
    });

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('Get /fills gets fills', async () => {
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expected: Partial<FillResponseObject> = {
        side: testConstants.defaultFill.side,
        liquidity: testConstants.defaultFill.liquidity,
        market: testConstants.defaultPerpetualMarket.ticker,
        marketType: MarketType.PERPETUAL,
        price: testConstants.defaultFill.price,
        size: testConstants.defaultFill.size,
        fee: testConstants.defaultFill.fee,
        affiliateRevShare: testConstants.defaultFill.affiliateRevShare,
        type: testConstants.defaultFill.type,
        orderId: testConstants.defaultFill.orderId,
        createdAt: testConstants.defaultFill.createdAt,
        createdAtHeight: testConstants.defaultFill.createdAtHeight,
      };

      expect(response.body.fills).toHaveLength(1);
      expect(response.body.fills).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected,
          }),
        ]),
      );
    });

    it('Get /fills with market gets fills for market', async () => {
      // Order and fill for BTC-USD
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      // Order and fill for ETH-USD
      const ethOrder: OrderFromDatabase = await OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '3',
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
      });
      const ethFill: FillFromDatabase = await FillTable.create({
        ...testConstants.defaultFill,
        orderId: ethOrder.id,
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
        eventId: testConstants.defaultTendermintEventId2,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}` +
          `&market=${testConstants.defaultPerpetualMarket2.ticker}&marketType=${MarketType.PERPETUAL}`,
      });

      const expected: Partial<FillResponseObject> = {
        side: ethFill.side,
        liquidity: ethFill.liquidity,
        market: testConstants.defaultPerpetualMarket2.ticker,
        marketType: MarketType.PERPETUAL,
        price: ethFill.price,
        size: ethFill.size,
        fee: ethFill.fee,
        affiliateRevShare: ethFill.affiliateRevShare,
        type: ethFill.type,
        orderId: ethOrder.id,
        createdAt: ethFill.createdAt,
        createdAtHeight: ethFill.createdAtHeight,
        subaccountNumber: defaultSubaccountNumber,
      };

      // Only the ETH-USD order should be returned
      expect(response.body.fills).toHaveLength(1);
      expect(response.body.fills).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected,
          }),
        ]),
      );
    });

    it('Get /fills with market gets correctly ordered fills', async () => {
      // Order and fill for BTC-USD
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create({
        ...testConstants.defaultFill,
        eventId: testConstants.defaultTendermintEventId2,
      });

      // Order and fill for ETH-USD
      const ethOrder: OrderFromDatabase = await OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '3',
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
      });
      const ethFill: FillFromDatabase = await FillTable.create({
        ...testConstants.defaultFill,
        orderId: ethOrder.id,
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
        createdAtHeight: '1',
      });

      let response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expected: Partial<FillResponseObject>[] = [
        {
          side: testConstants.defaultFill.side,
          liquidity: testConstants.defaultFill.liquidity,
          market: testConstants.defaultPerpetualMarket.ticker,
          marketType: MarketType.PERPETUAL,
          price: testConstants.defaultFill.price,
          size: testConstants.defaultFill.size,
          fee: testConstants.defaultFill.fee,
          affiliateRevShare: testConstants.defaultFill.affiliateRevShare,
          type: testConstants.defaultFill.type,
          orderId: testConstants.defaultFill.orderId,
          createdAt: testConstants.defaultFill.createdAt,
          createdAtHeight: testConstants.defaultFill.createdAtHeight,
        },
        {
          side: ethFill.side,
          liquidity: ethFill.liquidity,
          market: testConstants.defaultPerpetualMarket2.ticker,
          marketType: MarketType.PERPETUAL,
          price: ethFill.price,
          size: ethFill.size,
          fee: ethFill.fee,
          type: ethFill.type,
          orderId: ethOrder.id,
          createdAt: ethFill.createdAt,
          createdAtHeight: ethFill.createdAtHeight,
        },
      ];

      // Page is not specified, so fills should be returned sorted by createdAtHeight
      // in descending order.
      expect(response.body.fills).toHaveLength(2);
      expect(response.body.fills).toEqual(
        [
          expect.objectContaining({
            ...expected[0],
          }),
          expect.objectContaining({
            ...expected[1],
          }),
        ],
      );

      response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=1&limit=2`,
      });
      // Page is specified, so fills should be sorted by eventId in ascending order.
      expect(response.body.fills).toHaveLength(2);
      expect(response.body.fills).toEqual(
        [
          expect.objectContaining({
            ...expected[1],
          }),
          expect.objectContaining({
            ...expected[0],
          }),
        ],
      );
    });

    it('Get /fills with market gets fills ordered by createdAtHeight descending and paginated', async () => {
      // Order and fill for BTC-USD
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      // Order and fill for ETH-USD
      const ethOrder: OrderFromDatabase = await OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '3',
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
      });
      const ethFill: FillFromDatabase = await FillTable.create({
        ...testConstants.defaultFill,
        orderId: ethOrder.id,
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
        eventId: testConstants.defaultTendermintEventId2,
        createdAtHeight: '1',
      });

      const responsePage1: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=1&limit=1`,
      });

      const responsePage2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=2&limit=1`,
      });

      const expected: Partial<FillResponseObject>[] = [
        {
          side: testConstants.defaultFill.side,
          liquidity: testConstants.defaultFill.liquidity,
          market: testConstants.defaultPerpetualMarket.ticker,
          marketType: MarketType.PERPETUAL,
          price: testConstants.defaultFill.price,
          size: testConstants.defaultFill.size,
          fee: testConstants.defaultFill.fee,
          affiliateRevShare: testConstants.defaultFill.affiliateRevShare,
          type: testConstants.defaultFill.type,
          orderId: testConstants.defaultFill.orderId,
          createdAt: testConstants.defaultFill.createdAt,
          createdAtHeight: testConstants.defaultFill.createdAtHeight,
        },
        {
          side: ethFill.side,
          liquidity: ethFill.liquidity,
          market: testConstants.defaultPerpetualMarket2.ticker,
          marketType: MarketType.PERPETUAL,
          price: ethFill.price,
          size: ethFill.size,
          fee: ethFill.fee,
          affiliateRevShare: ethFill.affiliateRevShare,
          type: ethFill.type,
          orderId: ethOrder.id,
          createdAt: ethFill.createdAt,
          createdAtHeight: ethFill.createdAtHeight,
        },
      ];

      expect(responsePage1.body.pageSize).toStrictEqual(1);
      expect(responsePage1.body.offset).toStrictEqual(0);
      expect(responsePage1.body.totalResults).toStrictEqual(2);
      expect(responsePage1.body.fills).toHaveLength(1);
      expect(responsePage1.body.fills).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected[0],
          }),
        ]),
      );

      expect(responsePage2.body.pageSize).toStrictEqual(1);
      expect(responsePage2.body.offset).toStrictEqual(1);
      expect(responsePage2.body.totalResults).toStrictEqual(2);
      expect(responsePage2.body.fills).toHaveLength(1);
      expect(responsePage2.body.fills).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected[1],
          }),
        ]),
      );
    });

    it('Get /fills with market with no fills', async () => {
      // Order and fill for BTC-USD
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}` +
          `&market=${testConstants.defaultPerpetualMarket2.ticker}&marketType=${MarketType.PERPETUAL}`,
      });

      expect(response.body.fills).toEqual([]);
    });

    it.each([
      [
        'market passed in without marketType',
        {
          address: defaultAddress,
          subaccountNumber: defaultSubaccountNumber,
          market: defaultMarket,
        },
        'marketType',
        'marketType must be provided if market is provided',
      ],
      [
        'marketType passed in without market',
        {
          address: defaultAddress,
          subaccountNumber: defaultSubaccountNumber,
          marketType: MarketType.PERPETUAL,
        },
        'market',
        'market must be provided if marketType is provided',
      ],
      [
        'invalid marketType',
        {
          address: defaultAddress,
          subaccountNumber: defaultSubaccountNumber,
          marketType: 'INVALID',
          market: defaultMarket,
        },
        'marketType',
        'marketType must be a valid market type (PERPETUAL/SPOT)',
      ],
    ])('Returns 400 when validation fails: %s', async (
      _reason: string,
      queryParams: {
        address?: string,
        subaccountNumber?: number,
        market?: string,
        marketType?: string,
        createdBeforeOrAt?: string,
        createdBeforeOrAtHeight?: number,
      },
      fieldWithError: string,
      expectedErrorMsg: string,
    ) => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills?${getQueryString(queryParams)}`,
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

    it('Returns 404 with unknown market and type', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}` +
          `&market=${invalidMarket}&marketType=${MarketType.PERPETUAL}`,
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: `${invalidMarket} not found in markets of type ${MarketType.PERPETUAL}`,
          },
        ],
      });
    });

    it('Get /fills/parentSubaccountNumber gets fills', async () => {
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);
      await OrderTable.create(testConstants.isolatedMarketOrder);
      await FillTable.create(testConstants.isolatedMarketFill);
      await FillTable.create(testConstants.isolatedMarketFill2);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${parentSubaccountNumber}`,
      });

      // Use fillResponseObjectFromFillCreateObject to create expectedFills
      const expectedFills: Partial<FillResponseObject>[] = [
        fillResponseObjectFromFillCreateObject(testConstants.defaultFill, defaultSubaccountNumber),
        fillResponseObjectFromFillCreateObject(testConstants.isolatedMarketFill,
          testConstants.isolatedSubaccount.subaccountNumber),
        fillResponseObjectFromFillCreateObject(testConstants.isolatedMarketFill2,
          testConstants.isolatedSubaccount2.subaccountNumber),
      ];

      expect(response.body.fills).toHaveLength(3);
      expect(response.body.fills).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedFills[0],
          }),
          expect.objectContaining({
            ...expectedFills[1],
          }),
          expect.objectContaining({
            ...expectedFills[2],
          }),
        ]),
      );
    });

    it('Get /fills/parentSubaccountNumber gets fills ordered by createdAtHeight descending and paginated', async () => {
      // Order and fill for BTC-USD
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      // Order and fill for ETH-USD
      const ethOrder: OrderFromDatabase = await OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '3',
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
      });
      const ethFill: FillFromDatabase = await FillTable.create({
        ...testConstants.defaultFill,
        orderId: ethOrder.id,
        clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
        eventId: testConstants.defaultTendermintEventId2,
        createdAtHeight: '1',
      });

      const responsePage1: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
          `&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=1&limit=1`,
      });

      const responsePage2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
          `&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=2&limit=1`,
      });

      const expected: Partial<FillResponseObject>[] = [
        {
          side: testConstants.defaultFill.side,
          liquidity: testConstants.defaultFill.liquidity,
          market: testConstants.defaultPerpetualMarket.ticker,
          marketType: MarketType.PERPETUAL,
          price: testConstants.defaultFill.price,
          size: testConstants.defaultFill.size,
          fee: testConstants.defaultFill.fee,
          affiliateRevShare: testConstants.defaultFill.affiliateRevShare,
          type: testConstants.defaultFill.type,
          orderId: testConstants.defaultFill.orderId,
          createdAt: testConstants.defaultFill.createdAt,
          createdAtHeight: testConstants.defaultFill.createdAtHeight,
        },
        {
          side: ethFill.side,
          liquidity: ethFill.liquidity,
          market: testConstants.defaultPerpetualMarket2.ticker,
          marketType: MarketType.PERPETUAL,
          price: ethFill.price,
          size: ethFill.size,
          fee: ethFill.fee,
          affiliateRevShare: ethFill.affiliateRevShare,
          type: ethFill.type,
          orderId: ethOrder.id,
          createdAt: ethFill.createdAt,
          createdAtHeight: ethFill.createdAtHeight,
        },
      ];

      expect(responsePage1.body.pageSize).toStrictEqual(1);
      expect(responsePage1.body.offset).toStrictEqual(0);
      expect(responsePage1.body.totalResults).toStrictEqual(2);
      expect(responsePage1.body.fills).toHaveLength(1);
      expect(responsePage1.body.fills).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected[0],
          }),
        ]),
      );

      expect(responsePage2.body.pageSize).toStrictEqual(1);
      expect(responsePage2.body.offset).toStrictEqual(1);
      expect(responsePage2.body.totalResults).toStrictEqual(2);
      expect(responsePage2.body.fills).toHaveLength(1);
      expect(responsePage2.body.fills).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected[1],
          }),
        ]),
      );
    });

    it.each([
      [
        'market passed in without marketType',
        {
          address: defaultAddress,
          subaccountNumber: defaultSubaccountNumber,
          market: defaultMarket,
        },
        'marketType',
        'marketType must be provided if market is provided',
      ],
      [
        'marketType passed in without market',
        {
          address: defaultAddress,
          subaccountNumber: defaultSubaccountNumber,
          marketType: MarketType.PERPETUAL,
        },
        'market',
        'market must be provided if marketType is provided',
      ],
      [
        'invalid marketType',
        {
          address: defaultAddress,
          subaccountNumber: defaultSubaccountNumber,
          marketType: 'INVALID',
          market: defaultMarket,
        },
        'marketType',
        'marketType must be a valid market type (PERPETUAL/SPOT)',
      ],
    ])('Returns 400 when validation fails for parentSubaccount endpoint: %s', async (
      _reason: string,
      queryParams: {
        address?: string,
        parentSubaccountNumber?: number,
        market?: string,
        marketType?: string,
        createdBeforeOrAt?: string,
        createdBeforeOrAtHeight?: number,
      },
      fieldWithError: string,
      expectedErrorMsg: string,
    ) => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fills/parentSubaccountNumber?${getQueryString(queryParams)}`,
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
