import {
  dbHelpers,
  MarketFromDatabase,
  MarketTable,
  perpetualMarketRefresher,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  testConstants,
  testMocks,
  LiquidityTiersFromDatabase,
  LiquidityTiersTable,
  liquidityTierRefresher,
} from '@dydxprotocol-indexer/postgres';
import {
  OrderbookMidPricesCache,
} from '@dydxprotocol-indexer/redis';
import { RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { getQueryString, sendRequest } from '../../../helpers/helpers';
import _ from 'lodash';
import { perpetualMarketToResponseObject } from '../../../../src/request-helpers/request-transformer';

jest.mock('@dydxprotocol-indexer/redis', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/redis'),
  OrderbookMidPricesCache: {
    getMedianPrice: jest.fn(),
  },
}));

describe('perpetual-markets-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await liquidityTierRefresher.updateLiquidityTiers();
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  describe('/', () => {
    const invalidTicker: string = 'UNKNOWN';

    it('Get / gets all tickers', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/perpetualMarkets/',
      });

      // Only two markets
      const perpetualMarkets: PerpetualMarketFromDatabase[] = await
      PerpetualMarketTable.findAll({}, []);
      const markets: MarketFromDatabase[] = await MarketTable.findAll({}, []);
      const liquidityTiers: LiquidityTiersFromDatabase[] = await Promise.all(
        _.map(
          perpetualMarkets,
          async (perpetualMarket) => {
            return await LiquidityTiersTable.findById(
              perpetualMarket.liquidityTierId,
            ) as LiquidityTiersFromDatabase;
          }),
      );

      expectResponseWithMarkets(response, perpetualMarkets, liquidityTiers, markets);
    });

    it('Get / with out of order markets', async () => {
      // Create markets and perpetual markets in different orders.
      await MarketTable.create({
        ...testConstants.defaultMarket,
        id: 99,
        pair: 'XXX-USD',
      });
      await MarketTable.create({
        ...testConstants.defaultMarket,
        id: 100,
        pair: 'YYY-USD',
      });
      await PerpetualMarketTable.create({
        ...testConstants.defaultPerpetualMarket,
        id: '100',
        marketId: 100,
        ticker: 'YYY-USD',
      });
      await PerpetualMarketTable.create({
        ...testConstants.defaultPerpetualMarket,
        id: '99',
        marketId: 99,
        ticker: 'XXX-USD',
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/perpetualMarkets/',
      });

      const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
        {}, []);
      const markets: MarketFromDatabase[] = await MarketTable.findAll({}, []);
      const liquidityTiers: LiquidityTiersFromDatabase[] = await Promise.all(
        _.map(
          perpetualMarkets,
          async (perpetualMarket) => {
            return await LiquidityTiersTable.findById(
              perpetualMarket.liquidityTierId,
            ) as LiquidityTiersFromDatabase;
          }),
      );

      expectResponseWithMarkets(response, perpetualMarkets, liquidityTiers, markets);
    });

    it('Get / gets all markets with limit', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualMarkets?${getQueryString({ limit: 1 })}`,
      });

      // Only one market
      const perpetualMarket:
      PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable.findByTicker(
        testConstants.defaultPerpetualMarket.ticker,
      );
      const market:
      MarketFromDatabase | undefined = await MarketTable.findById(
        testConstants.defaultPerpetualMarket.marketId,
      );
      const liquidityTier:
      LiquidityTiersFromDatabase | undefined = await LiquidityTiersTable.findById(
        testConstants.defaultPerpetualMarket.liquidityTierId,
      );
      expectResponseWithMarkets(
        response,
        [perpetualMarket!],
        [liquidityTier!],
        [market!],
      );
    });

    it('Get / with a ticker in the query gets a market with a matching ticker', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualMarkets?${getQueryString({ ticker: testConstants.defaultPerpetualMarket2.ticker })}`,
      });

      // Only one market
      const perpetualMarket:
      PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable.findByTicker(
        testConstants.defaultPerpetualMarket2.ticker,
      );
      const market:
      MarketFromDatabase | undefined = await MarketTable.findById(
        testConstants.defaultPerpetualMarket2.marketId,
      );
      const liquidityTier:
      LiquidityTiersFromDatabase | undefined = await LiquidityTiersTable.findById(
        testConstants.defaultPerpetualMarket2.liquidityTierId,
      );
      expectResponseWithMarkets(
        response,
        [perpetualMarket!],
        [liquidityTier!],
        [market!],
      );
    });

    it('Returns 404 with unknown ticker', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualMarkets/${invalidTicker}`,
        expectedStatus: 404,
      });

      expect(response.body.error).toContain('Not Found');
    });
  });

  describe('GET /v4/perpetualMarkets/orderbookMidPrices', () => {
    it('returns mid prices for all markets when no tickers are specified', async () => {
      (OrderbookMidPricesCache.getMedianPrice as jest.Mock).mockImplementation((client, ticker) => {
        const prices: {[key: string]: string} = {
          'BTC-USD': '30000.5',
          'ETH-USD': '2000.25',
          'SHIB-USD': '5.75',
        };
        return Promise.resolve(prices[ticker]);
      });

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/perpetualMarkets/orderbookMidPrices',
      });

      expect(response.status).toBe(200);
      expect(response.body).toEqual({
        'BTC-USD': '30000.5',
        'ETH-USD': '2000.25',
        'SHIB-USD': '5.75',
      });
      const numMarkets = (await PerpetualMarketTable.findAll({}, [])).length;
      expect(OrderbookMidPricesCache.getMedianPrice).toHaveBeenCalledTimes(numMarkets);
    });

    it('returns mid prices for multiple specified tickers', async () => {
      (OrderbookMidPricesCache.getMedianPrice as jest.Mock).mockImplementation((client, ticker) => {
        const prices: {[key: string]: string} = {
          'BTC-USD': '30000.5',
          'ETH-USD': '2000.25',
        };
        return Promise.resolve(prices[ticker]);
      });

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/perpetualMarkets/orderbookMidPrices?tickers=BTC-USD&tickers=ETH-USD',
      });

      expect(response.status).toBe(200);
      expect(response.body).toEqual({
        'BTC-USD': '30000.5',
        'ETH-USD': '2000.25',
      });

      expect(OrderbookMidPricesCache.getMedianPrice).toHaveBeenCalledTimes(2);
    });

    it('returns mid prices for one specified ticker', async () => {
      (OrderbookMidPricesCache.getMedianPrice as jest.Mock).mockImplementation((client, ticker) => {
        const prices: {[key: string]: string} = {
          'BTC-USD': '30000.5',
          'ETH-USD': '2000.25',
        };
        return Promise.resolve(prices[ticker]);
      });

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/perpetualMarkets/orderbookMidPrices?tickers=BTC-USD',
      });

      expect(response.status).toBe(200);
      expect(response.body).toEqual({
        'BTC-USD': '30000.5',
      });

      expect(OrderbookMidPricesCache.getMedianPrice).toHaveBeenCalledTimes(1);
    });

    it('omits markets with no mid price', async () => {
      (OrderbookMidPricesCache.getMedianPrice as jest.Mock).mockImplementation((client, ticker) => {
        const prices: {[key: string]: string | null} = {
          'BTC-USD': '30000.5',
          'ETH-USD': null,
          'SHIB-USD': '5.75',
        };
        return Promise.resolve(prices[ticker]);
      });

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/perpetualMarkets/orderbookMidPrices',
      });

      expect(response.status).toBe(200);
      expect(response.body).toEqual({
        'BTC-USD': '30000.5',
        'SHIB-USD': '5.75',
      });
    });

    it('returns an empty object when no markets have mid prices', async () => {
      (OrderbookMidPricesCache.getMedianPrice as jest.Mock).mockResolvedValue(null);

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/perpetualMarkets/orderbookMidPrices',
      });

      expect(response.status).toBe(200);
      expect(response.body).toEqual({});
    });

    it('returns prices only for valid tickers and ignores invalid tickers', async () => {
      (OrderbookMidPricesCache.getMedianPrice as jest.Mock).mockImplementation((client, ticker) => {
        const prices: {[key: string]: string} = {
          'BTC-USD': '30000.5',
        };
        return Promise.resolve(prices[ticker]);
      });

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/perpetualMarkets/orderbookMidPrices?tickers=BTC-USD&tickers=INVALID-TICKER',
      });

      expect(response.status).toBe(200);
      expect(response.body).toEqual({
        'BTC-USD': '30000.5',
      });
      expect(OrderbookMidPricesCache.getMedianPrice).toHaveBeenCalledTimes(1);
    });
  });
});

function expectResponseWithMarkets(
  response: request.Response,
  perpetualMarkets: PerpetualMarketFromDatabase[],
  liquidityTiers: LiquidityTiersFromDatabase[],
  markets: MarketFromDatabase[],
): void {
  expect(_.size(response.body.markets)).toEqual(perpetualMarkets.length);
  expect(_.size(response.body.markets)).toEqual(markets.length);
  expect(_.size(response.body.markets)).toEqual(liquidityTiers.length);

  _.each(_.zip(perpetualMarkets, liquidityTiers, markets), (
    [perpetualMarket, liquidityTier, market]:
    [PerpetualMarketFromDatabase | undefined,
      LiquidityTiersFromDatabase | undefined, MarketFromDatabase | undefined],
  ) => {
    expect(response.body.markets[perpetualMarket!.ticker]).toEqual(
      perpetualMarketToResponseObject(perpetualMarket!, liquidityTier!, market!),
    );
  });
}
