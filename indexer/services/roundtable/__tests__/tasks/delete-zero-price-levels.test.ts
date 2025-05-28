import { stats } from '@dydxprotocol-indexer/base';
import {
  dbHelpers, OrderSide, PerpetualMarketFromDatabase, PerpetualMarketTable, testMocks,
} from '@dydxprotocol-indexer/postgres';
import { OrderbookLevels, OrderbookLevelsCache, redis } from '@dydxprotocol-indexer/redis';
import config from '../../src/config';
import deleteZeroPriceLevelsTask from '../../src/tasks/delete-zero-price-levels';
import { redisClient } from '../../src/helpers/redis';

describe('delete-zero-price-levels', () => {
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

  it('succeeds with no levels', async () => {
    await deleteZeroPriceLevelsTask();
    expect(stats.gauge).toHaveBeenCalledWith(
      `${config.SERVICE_NAME}.delete_zero_price_levels.num_levels_deleted`,
      0,
    );
  });

  it('deletes zero price levels for all orderbooks', async () => {
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
            '0',
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
            '0',
            redisClient,
          ),
        ]);
      },
      ));

    await deleteZeroPriceLevelsTask();

    for (const perpetualMarket of perpetualMarkets) {
      const orderbookLevels: OrderbookLevels = await OrderbookLevelsCache.getOrderBookLevels(
        perpetualMarket.ticker,
        redisClient,
        {
          removeZeros: false,
          sortSides: true,
        },
      );

      expect(orderbookLevels.bids).toMatchObject([
        {
          humanPrice: '45100',
          quantums: '2000',
        },
      ]);
      expect(orderbookLevels.asks).toMatchObject([
        {
          humanPrice: '45300',
          quantums: '3000',
        },
      ]);
    }
    expect(stats.gauge).toHaveBeenCalledWith(
      `${config.SERVICE_NAME}.delete_zero_price_levels.num_levels_deleted`,
      2 * perpetualMarkets.length,
    );
  });
});
