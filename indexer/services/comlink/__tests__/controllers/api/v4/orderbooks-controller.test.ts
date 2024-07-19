import {
  dbHelpers,
  testConstants,
  OrderSide,
  testMocks,
  perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';
import { OrderbookLevelsCache, redis } from '@dydxprotocol-indexer/redis';
import { redisClient } from '../../../../src/helpers/redis/redis-controller';

describe('orderbooks-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('/perpetualMarkets', () => {
    beforeEach(async () => {
      await redis.deleteAllAsync(redisClient);
    });

    afterEach(async () => {
      await redis.deleteAllAsync(redisClient);
      await dbHelpers.clearData();
    });

    it('Get /:ticker gets orderbook for perpetual market with the matching ticker', async () => {
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
      const commonPriceLevel: any = {
        ticker: testConstants.defaultPerpetualMarket.ticker,
        client: redisClient,
      };
      await Promise.all([
        OrderbookLevelsCache.updatePriceLevel({
          ...commonPriceLevel,
          humanPrice: '40000.1',
          side: OrderSide.BUY,
          sizeDeltaInQuantums: '1000000000',
        }),
        OrderbookLevelsCache.updatePriceLevel({
          ...commonPriceLevel,
          humanPrice: '42030.5',
          side: OrderSide.BUY,
          sizeDeltaInQuantums: '2000000000',
        }),
        OrderbookLevelsCache.updatePriceLevel({
          ...commonPriceLevel,
          humanPrice: '47500',
          side: OrderSide.BUY,
          sizeDeltaInQuantums: '3000000000',
        }),
        // crossing level that will be filtered out, should not be returned in response
        OrderbookLevelsCache.updatePriceLevel({
          ...commonPriceLevel,
          humanPrice: '46050',
          side: OrderSide.SELL,
          sizeDeltaInQuantums: '1750000000',
        }),
        OrderbookLevelsCache.updatePriceLevel({
          ...commonPriceLevel,
          humanPrice: '51000.4',
          side: OrderSide.SELL,
          sizeDeltaInQuantums: '1500000000',
        }),
        OrderbookLevelsCache.updatePriceLevel({
          ...commonPriceLevel,
          humanPrice: '50050.2',
          side: OrderSide.SELL,
          sizeDeltaInQuantums: '2500000000',
        }),
        OrderbookLevelsCache.updatePriceLevel({
          ...commonPriceLevel,
          humanPrice: '53200.6',
          side: OrderSide.SELL,
          sizeDeltaInQuantums: '500000000',
        }),
        OrderbookLevelsCache.updatePriceLevel({
          ...commonPriceLevel,
          humanPrice: '60300.8',
          side: OrderSide.SELL,
          sizeDeltaInQuantums: '250000000',
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orderbooks/perpetualMarket/${testConstants.defaultPerpetualMarket.ticker}`,
      });

      expect(response.body.bids).toHaveLength(3);
      expect(response.body.asks).toHaveLength(4);
      // Check that the bids are sorted by price in descending order
      expect(response.body.bids).toEqual([
        { price: '47500', size: '0.3' }, // 3,000,000,000 * 1e-10
        { price: '42030.5', size: '0.2' }, // 2,000,000,000 * 1e-10
        { price: '40000.1', size: '0.1' }, // 1,000,000,000 * 1e-10
      ]);
      // Check that the asks are sorted by price in ascending order
      expect(response.body.asks).toEqual([
        { price: '50050.2', size: '0.25' }, // 2,500,000,000 * 1e-10
        { price: '51000.4', size: '0.15' }, // 1,500,000,000 * 1e-10
        { price: '53200.6', size: '0.05' }, // 500,000,000 * 1e-10
        { price: '60300.8', size: '0.025' }, // 250,000,000 * 1e-10
      ]);
    });

    it('Get /:ticker returns 404 if ticker is not found', async () => {
      const invalidTicker: string = 'invalid-invalid';

      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orderbooks/perpetualMarket/${invalidTicker}`,
        expectedStatus: 404,
      });
    });
  });
});
