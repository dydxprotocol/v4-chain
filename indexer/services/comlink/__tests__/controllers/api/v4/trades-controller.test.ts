import {
  DEFAULT_POSTGRES_OPTIONS,
  dbHelpers,
  testMocks,
  testConstants,
  OrderTable,
  FillTable,
  FillFromDatabase,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { RequestMethod, TradeResponseObject } from '../../../../src/types';
import request from 'supertest';
import { createMakerTakerOrderAndFill, sendRequest } from '../../../helpers/helpers';
import { fillToTradeResponseObject } from '../../../../src/request-helpers/request-transformer';

describe('trades-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('/perpetualMarket', () => {
    const invalidTicker: string = 'UNKNOWN';

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('Get /:ticker gets trades for a ticker', async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
      // Order and fill for BTC-USD (maker and taker)
      await createMakerTakerOrderAndFill(
        testConstants.defaultOrder,
        testConstants.defaultFill,
      );

      // Order and fill for ETH-USD (maker and taker)
      const ethSize: string = '600';
      const fills: {
        makerFill: FillFromDatabase,
        takerFill: FillFromDatabase,
      } = await createMakerTakerOrderAndFill(
        {
          ...testConstants.defaultOrder,
          clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
        },
        {
          ...testConstants.defaultFill,
          size: ethSize,
          clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
          eventId: testConstants.defaultTendermintEventId2,
        },
      );

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/trades/perpetualMarket/${testConstants.defaultPerpetualMarket2.ticker}`,
      });

      const expected: TradeResponseObject = fillToTradeResponseObject(fills.takerFill);

      // Only a single trade, BTC fills filtered out, and ETH maker fill filtered out
      expect(response.body.trades).toHaveLength(1);
      expect(response.body.trades).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected,
          }),
        ]),
      );
    });

    it('Get /:ticker gets trades for a ticker in descending order by createdAtHeight', async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
      // Order and fill for BTC-USD (maker and taker)
      const fills1: {
        makerFill: FillFromDatabase,
        takerFill: FillFromDatabase,
      } = await createMakerTakerOrderAndFill(
        testConstants.defaultOrder,
        testConstants.defaultFill,
      );

      const btcSize2: string = '600';
      const fills2: {
        makerFill: FillFromDatabase,
        takerFill: FillFromDatabase,
      } = await createMakerTakerOrderAndFill(
        testConstants.defaultOrder,
        {
          ...testConstants.defaultFill,
          size: btcSize2,
          eventId: testConstants.defaultTendermintEventId2,
          createdAtHeight: '1',
        },
      );

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/trades/perpetualMarket/${testConstants.defaultPerpetualMarket.ticker}`,
      });

      const expected: TradeResponseObject[] = [
        fillToTradeResponseObject(fills1.takerFill),
        fillToTradeResponseObject(fills2.takerFill),
      ];

      // Expect both trades, ordered by createdAtHeight in descending order
      expect(response.body.trades).toHaveLength(2);
      expect(response.body.trades).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected[0],
          }),
          expect.objectContaining({
            ...expected[1],
          }),
        ]),
      );
    });

    it('Get /:ticker gets trades for a ticker in descending order by createdAtHeight and paginated', async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
      // Order and fill for BTC-USD (maker and taker)
      const fills1: {
        makerFill: FillFromDatabase,
        takerFill: FillFromDatabase,
      } = await createMakerTakerOrderAndFill(
        testConstants.defaultOrder,
        testConstants.defaultFill,
      );

      const btcSize2: string = '600';
      const fills2: {
        makerFill: FillFromDatabase,
        takerFill: FillFromDatabase,
      } = await createMakerTakerOrderAndFill(
        testConstants.defaultOrder,
        {
          ...testConstants.defaultFill,
          size: btcSize2,
          eventId: testConstants.defaultTendermintEventId2,
          createdAtHeight: '1',
        },
      );

      const responsePage1: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/trades/perpetualMarket/${testConstants.defaultPerpetualMarket.ticker}?page=1&limit=1`,
      });

      const responsePage2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/trades/perpetualMarket/${testConstants.defaultPerpetualMarket.ticker}?page=2&limit=1`,
      });

      const expected: TradeResponseObject[] = [
        fillToTradeResponseObject(fills1.takerFill),
        fillToTradeResponseObject(fills2.takerFill),
      ];

      // Expect both trades, ordered by createdAtHeight in descending order
      expect(responsePage1.body.pageSize).toStrictEqual(1);
      expect(responsePage1.body.offset).toStrictEqual(0);
      expect(responsePage1.body.totalResults).toStrictEqual(2);
      expect(responsePage1.body.trades).toHaveLength(1);
      expect(responsePage1.body.trades).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected[0],
          }),
        ]),
      );

      expect(responsePage1.body.pageSize).toStrictEqual(1);
      expect(responsePage1.body.offset).toStrictEqual(0);
      expect(responsePage2.body.totalResults).toStrictEqual(2);
      expect(responsePage2.body.trades).toHaveLength(1);
      expect(responsePage2.body.trades).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected[1],
          }),
        ]),
      );
    });

    it('Get /:ticker routes paginated reads through DEFAULT_POSTGRES_OPTIONS (read-replica regression guard)', async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
      await createMakerTakerOrderAndFill(
        testConstants.defaultOrder,
        testConstants.defaultFill,
      );

      // Mutate the live DEFAULT_POSTGRES_OPTIONS object with a marker property. Since the
      // controller spreads this object at call time (`{ ...DEFAULT_POSTGRES_OPTIONS, orderBy }`),
      // the marker will only show up in the findAll call if that spread actually happens -
      // independent of whatever USE_READ_REPLICA currently evaluates to in this test environment.
      const testMarker: unique symbol = Symbol('read-replica-regression-marker');
      (DEFAULT_POSTGRES_OPTIONS as Record<string | symbol, unknown>)[testMarker] = true;
      const findAllSpy: jest.SpyInstance = jest.spyOn(FillTable, 'findAll');

      try {
        await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/trades/perpetualMarket/${testConstants.defaultPerpetualMarket.ticker}?page=1&limit=1`,
        });

        expect(findAllSpy).toHaveBeenCalledWith(
          expect.anything(),
          expect.anything(),
          expect.objectContaining({ [testMarker]: true }),
        );
      } finally {
        delete (DEFAULT_POSTGRES_OPTIONS as Record<string | symbol, unknown>)[testMarker];
        findAllSpy.mockRestore();
      }
    });

    it('Get /:ticker for ticker with no fills', async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
      // Order and fill for BTC-USD
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/trades/perpetualMarket/${testConstants.defaultPerpetualMarket2.ticker}`,
      });

      expect(response.body.trades).toEqual([]);
    });

    it('Get /:ticker for ticker with no fills and paginated', async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
      // Order and fill for BTC-USD
      await OrderTable.create(testConstants.defaultOrder);
      await FillTable.create(testConstants.defaultFill);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/trades/perpetualMarket/${testConstants.defaultPerpetualMarket2.ticker}?page=1&limit=1`,
      });

      expect(response.body.trades).toEqual([]);
    });

    it('Get /:ticker for ticker with price < 1e-6', async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
      // Order and fill for BTC-USD (maker and taker)
      const fills1: {
        makerFill: FillFromDatabase,
        takerFill: FillFromDatabase,
      } = await createMakerTakerOrderAndFill(
        {
          ...testConstants.defaultOrder,
          clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
          price: '0.000000065',
        },
        {
          ...testConstants.defaultFill,
          clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
          price: '0.000000064',
        },
      );

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/trades/perpetualMarket/${testConstants.defaultPerpetualMarket3.ticker}`,
      });

      const expected: TradeResponseObject[] = [
        fillToTradeResponseObject(fills1.takerFill),
      ];

      expect(response.body.trades).toHaveLength(1);
      expect(response.body.trades).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected[0],
          }),
        ]),
      );
    });

    it('Returns 404 with unknown ticker', async () => {
      await testMocks.seedData();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/trades/perpetualMarket/${invalidTicker}`,
        expectedStatus: 400,
      });

      expect(response.body.errors[0]).toEqual(expect.objectContaining({
        msg: 'ticker must be a valid ticker (BTC-USD, etc)',
      }));
    });
  });
});
