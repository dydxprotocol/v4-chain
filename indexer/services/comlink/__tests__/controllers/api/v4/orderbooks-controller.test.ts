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
      const ticker = testConstants.defaultPerpetualMarket.ticker;
      await Promise.all([
        OrderbookLevelsCache.updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '40000.1',
          '1000000000',
          redisClient,
        ),
        OrderbookLevelsCache.updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '42030.5',
          '2000000000',
          redisClient,
        ),
        OrderbookLevelsCache.updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '47500',
          '3000000000',
          redisClient,
        ),
        // crossing level that will be filtered out, should not be returned in response
        OrderbookLevelsCache.updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '46050',
          '1750000000',
          redisClient,
        ),
        OrderbookLevelsCache.updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '51000.4',
          '1500000000',
          redisClient,
        ),
        OrderbookLevelsCache.updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '50050.2',
          '2500000000',
          redisClient,
        ),
        OrderbookLevelsCache.updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '53200.6',
          '500000000',
          redisClient,
        ),
        OrderbookLevelsCache.updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '60300.8',
          '250000000',
          redisClient,
        ),
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
