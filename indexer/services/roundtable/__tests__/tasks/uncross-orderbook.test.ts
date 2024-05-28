import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  dbHelpers,
  OrderSide,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { OrderbookLevelsCache, redis } from '@dydxprotocol-indexer/redis';
import runTask from '../../src/tasks/uncross-orderbook';
import { redisClient } from '../../src/helpers/redis';

jest.mock('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache'),
  deleteStalePriceLevel: jest.fn(),
}));

describe('uncross-orderbook', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
    jest.spyOn(logger, 'info');
    jest.spyOn(stats, 'gauge');
  });

  beforeEach(async () => {
    await testMocks.seedData();
    jest.clearAllMocks();
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

  it('succeeds without any crossed levels', async () => {
    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );
    const market = perpetualMarkets[0];
    await OrderbookLevelsCache.updatePriceLevel({
      ticker: market.ticker,
      side: OrderSide.BUY,
      humanPrice: '30100',
      sizeDeltaInQuantums: '2000',
      client: redisClient,
    });
    await OrderbookLevelsCache.updatePriceLevel({
      ticker: market.ticker,
      side: OrderSide.SELL,
      humanPrice: '45000',
      sizeDeltaInQuantums: '1000',
      client: redisClient,
    });

    await runTask();
    expect(logger.info).not.toHaveBeenCalled();
  });

  it('removes crossed bid and ask levels', async () => {
    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );
    const market = perpetualMarkets[0];
    await OrderbookLevelsCache.updatePriceLevel({
      ticker: market.ticker,
      side: OrderSide.BUY,
      humanPrice: '45100',
      sizeDeltaInQuantums: '2000',
      client: redisClient,
    });
    await OrderbookLevelsCache.updatePriceLevel({
      ticker: market.ticker,
      side: OrderSide.SELL,
      humanPrice: '45000',
      sizeDeltaInQuantums: '1000',
      client: redisClient,
    });

    await runTask();

    expect(require('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache').deleteStalePriceLevel).toHaveBeenCalledWith(expect.objectContaining({
      side: OrderSide.BUY,
      humanPrice: '45100',
    }));
    expect(require('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache').deleteStalePriceLevel).not.toHaveBeenCalledWith(expect.objectContaining({
      side: OrderSide.SELL,
      humanPrice: '45000',
    }));
  });

  it('logs a failure to delete stale bid or ask level', async () => {
    const { deleteStalePriceLevel } = require('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache');
    deleteStalePriceLevel.mockImplementationOnce(() => false);

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );
    const market = perpetualMarkets[0];
    await OrderbookLevelsCache.updatePriceLevel({
      ticker: market.ticker,
      side: OrderSide.BUY,
      humanPrice: '45100',
      sizeDeltaInQuantums: '2000',
      client: redisClient,
    });
    await OrderbookLevelsCache.updatePriceLevel({
      ticker: market.ticker,
      side: OrderSide.SELL,
      humanPrice: '45000',
      sizeDeltaInQuantums: '1000',
      client: redisClient,
    });

    await runTask();

    expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
      message: expect.stringContaining('Failed to delete stale bid level for'),
      side: OrderSide.BUY,
      humanPrice: '45100',
    }));
  });
});
