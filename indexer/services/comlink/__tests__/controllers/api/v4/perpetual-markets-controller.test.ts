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
import { RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { getQueryString, sendRequest } from '../../../helpers/helpers';
import _ from 'lodash';
import { perpetualMarketToResponseObject } from '../../../../src/request-helpers/request-transformer';

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
  });

  afterEach(async () => {
    await dbHelpers.clearData();
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

    it('Returns 400 with unknown ticker or market', async () => {
      let response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualMarkets?${getQueryString({ ticker: invalidTicker })}`,
        expectedStatus: 400,
      });
      expect(response.body.errors[0].msg).toContain('ticker must be a valid ticker');

      response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualMarkets?${getQueryString({ market: invalidTicker })}`,
        expectedStatus: 400,
      });
      expect(response.body.errors[0].msg).toContain('ticker must be a valid ticker');
    });

    it('Returns 400 when both ticker and market are provided', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualMarkets?${getQueryString({
          ticker: testConstants.defaultPerpetualMarket.ticker,
          market: testConstants.defaultPerpetualMarket.ticker,
        })}`,
        expectedStatus: 400,
      });

      expect(response.body.errors[0].msg).toContain('Only one of ticker or market may be provided');
    });

    it('Market parameter functions the same as ticker parameter', async () => {
      // Get response using ticker parameter
      const tickerResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualMarkets?${getQueryString({
          ticker: testConstants.defaultPerpetualMarket.ticker,
        })}`,
      });

      // Get response using market parameter
      const marketResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualMarkets?${getQueryString({
          market: testConstants.defaultPerpetualMarket.ticker,
        })}`,
      });

      // Only one market should be returned in both cases
      const perpetualMarket: PerpetualMarketFromDatabase | undefined = await
      PerpetualMarketTable.findByTicker(testConstants.defaultPerpetualMarket.ticker);
      const market: MarketFromDatabase | undefined = await
      MarketTable.findById(testConstants.defaultPerpetualMarket.marketId);
      const liquidityTier: LiquidityTiersFromDatabase | undefined = await
      LiquidityTiersTable.findById(testConstants.defaultPerpetualMarket.liquidityTierId);

      // Verify both responses contain the same data
      expectResponseWithMarkets(
        tickerResponse,
        [perpetualMarket!],
        [liquidityTier!],
        [market!],
      );

      expectResponseWithMarkets(
        marketResponse,
        [perpetualMarket!],
        [liquidityTier!],
        [market!],
      );

      // The response bodies should be identical
      expect(tickerResponse.body).toEqual(marketResponse.body);
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
