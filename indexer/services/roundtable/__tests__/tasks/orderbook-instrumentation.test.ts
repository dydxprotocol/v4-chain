import { stats } from '@dydxprotocol-indexer/base';
import {
  dbHelpers,
  OrderSide,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  OrderbookLevelsCache,
  redis,
} from '@dydxprotocol-indexer/redis';
import
orderbookInstrumentationTask,
{ priceToSubticks }
  from '../../src/tasks/orderbook-instrumentation';
import { redisClient } from '../../src/helpers/redis';

describe('orderbook-instrumentation', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
    jest.spyOn(stats, 'gauge');
  });

  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await redis.deleteAllAsync(redisClient);
    jest.clearAllMocks();
  });

  it('succeeds with empty orderbook', async () => {
    await orderbookInstrumentationTask();
    expect(stats.gauge).toHaveBeenCalledTimes(0);
  });

  it('succeeds with stats', async () => {
    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
      {},
    );

    await Promise.all(
      perpetualMarkets.map(async (perpetualMarket: PerpetualMarketFromDatabase): Promise<void> => {
        await Promise.all([
          OrderbookLevelsCache.updatePriceLevel(
            perpetualMarket.ticker,
            OrderSide.BUY,
            '45100',
            '2000',
            redisClient,
          ),
          OrderbookLevelsCache.updatePriceLevel(
            perpetualMarket.ticker,
            OrderSide.BUY,
            '45200',
            '1500',
            redisClient,
          ),
          OrderbookLevelsCache.updatePriceLevel(
            perpetualMarket.ticker,
            OrderSide.BUY,
            '45500',
            '500',
            redisClient,
          ),
          OrderbookLevelsCache.updatePriceLevel(
            perpetualMarket.ticker,
            OrderSide.SELL,
            '45000',
            '500',
            redisClient,
          ),
          OrderbookLevelsCache.updatePriceLevel(
            perpetualMarket.ticker,
            OrderSide.SELL,
            '45300',
            '3000',
            redisClient,
          ),
          OrderbookLevelsCache.updatePriceLevel(
            perpetualMarket.ticker,
            OrderSide.SELL,
            '45400',
            '3500',
            redisClient,
          ),
          OrderbookLevelsCache.updatePriceLevel(
            perpetualMarket.ticker,
            OrderSide.SELL,
            '46800',
            '1500',
            redisClient,
          ),
        ]);
      },
      ));

    await orderbookInstrumentationTask();

    perpetualMarkets.forEach((perpetualMarket: PerpetualMarketFromDatabase) => {
      const tags: Object = { ticker: perpetualMarket.ticker };

      // Check for human prices being gauged
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.uncrossed_orderbook.best_bid_human',
        45200,
        tags,
      );
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.crossed_orderbook.best_bid_human',
        45500,
        tags,
      );
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.crossed_orderbook.num_bid_levels',
        3,
        tags,
      );
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.uncrossed_orderbook.best_ask_human',
        45300,
        tags,
      );
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.crossed_orderbook.best_ask_human',
        45000,
        tags,
      );
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.crossed_orderbook.num_ask_levels',
        4,
        tags,
      );

      // Check for subticks being gauged
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.uncrossed_orderbook.best_bid_subticks',
        priceToSubticks('45200', perpetualMarket),
        tags,
      );
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.crossed_orderbook.best_bid_subticks',
        priceToSubticks('45500', perpetualMarket),
        tags,
      );
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.uncrossed_orderbook.best_ask_subticks',
        priceToSubticks('45300', perpetualMarket),
        tags,
      );
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.crossed_orderbook.best_ask_subticks',
        priceToSubticks('45000', perpetualMarket),
        tags,
      );
    });
  });
});
