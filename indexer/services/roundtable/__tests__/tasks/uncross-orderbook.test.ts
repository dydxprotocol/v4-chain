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
import {
  RedisClient,
} from 'redis';

jest.mock('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache'),
  deleteStalePriceLevel: jest.fn(),
}));

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

async function updatePriceLevelWithSleep({
  ticker,
  side,
  humanPrice,
  sizeDeltaInQuantums,
  client,
}: {
  ticker: string,
  side: OrderSide,
  humanPrice: string,
  sizeDeltaInQuantums: string,
  client: RedisClient,
}) {
  await OrderbookLevelsCache.updatePriceLevel(
    ticker,
    side,
    humanPrice,
    sizeDeltaInQuantums,
    client,
  );
  await sleep(1000); // sleep for 1 second
}

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
    const market: PerpetualMarketFromDatabase = perpetualMarkets[0];
    await updatePriceLevelWithSleep({
      ticker: market.ticker,
      side: OrderSide.BUY,
      humanPrice: '30100',
      sizeDeltaInQuantums: '2000',
      client: redisClient,
    });
    await updatePriceLevelWithSleep({
      ticker: market.ticker,
      side: OrderSide.SELL,
      humanPrice: '45000',
      sizeDeltaInQuantums: '1000',
      client: redisClient,
    });

    await runTask();
    expect(require('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache').deleteStalePriceLevel).not.toHaveBeenCalled();
  });

  it('removes single crossed bid level', async () => {
    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );
    const market: PerpetualMarketFromDatabase = perpetualMarkets[0];
    await OrderbookLevelsCache.updatePriceLevel(
      market.ticker,
      OrderSide.BUY,
      '45100',
      '2000',
      redisClient,
    );
    await OrderbookLevelsCache.updatePriceLevel(
      market.ticker,
      OrderSide.SELL,
      '45000',
      '1000',
      redisClient,
    );

    await runTask();

    expect(require('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache').deleteStalePriceLevel).toHaveBeenCalledWith(
      expect.any(String),
      OrderSide.BUY,
      '45100',
      expect.any(Number),
      expect.any(RedisClient),
    );
    expect(require('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache').deleteStalePriceLevel).not.toHaveBeenCalledWith(
      expect.any(String),
      OrderSide.SELL,
      '45000',
      expect.any(Number),
      expect.any(RedisClient),
    );
  });

  it('removes first updated when there exist multiple crossed bid levels', async () => {
    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );
    const market: PerpetualMarketFromDatabase = perpetualMarkets[0];
    await updatePriceLevelWithSleep({
      ticker: market.ticker,
      side: OrderSide.BUY,
      humanPrice: '45100',
      sizeDeltaInQuantums: '2000',
      client: redisClient,
    });
    await updatePriceLevelWithSleep({
      ticker: market.ticker,
      side: OrderSide.SELL,
      humanPrice: '45000',
      sizeDeltaInQuantums: '1000',
      client: redisClient,
    });
    await updatePriceLevelWithSleep({
      ticker: market.ticker,
      side: OrderSide.BUY,
      humanPrice: '45200',
      sizeDeltaInQuantums: '2000',
      client: redisClient,
    });

    await runTask();

    // the bids are sorted in descending order
    // the asks are sorted in ascending order
    // the highest bid at price level 45200 was updated after the ask at price level 45000,
    // so the ask at price level 45000 should be removed.
    expect(require('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache').deleteStalePriceLevel).toHaveBeenCalledWith(
      expect.any(String),
      OrderSide.SELL,
      '45000',
      expect.any(Number),
      expect.any(RedisClient),
    );
  });

  it('removes multiple crossed bid levels', async () => {
    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );
    const market: PerpetualMarketFromDatabase = perpetualMarkets[0];
    await updatePriceLevelWithSleep({
      ticker: market.ticker,
      side: OrderSide.BUY,
      humanPrice: '45100',
      sizeDeltaInQuantums: '2000',
      client: redisClient,
    });
    await updatePriceLevelWithSleep({
      ticker: market.ticker,
      side: OrderSide.BUY,
      humanPrice: '45200',
      sizeDeltaInQuantums: '2000',
      client: redisClient,
    });
    await updatePriceLevelWithSleep({
      ticker: market.ticker,
      side: OrderSide.SELL,
      humanPrice: '45000',
      sizeDeltaInQuantums: '1000',
      client: redisClient,
    });

    await runTask();

    expect(require('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache').deleteStalePriceLevel).toHaveBeenCalledTimes(2);
    expect(require('@dydxprotocol-indexer/redis/build/src/caches/orderbook-levels-cache').deleteStalePriceLevel).not.toHaveBeenCalledWith(
      expect.objectContaining({
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
    const market: PerpetualMarketFromDatabase = perpetualMarkets[0];
    await OrderbookLevelsCache.updatePriceLevel(
      market.ticker,
      OrderSide.BUY,
      '45100',
      '2000',
      redisClient,
    );
    await OrderbookLevelsCache.updatePriceLevel(
      market.ticker,
      OrderSide.SELL,
      '45000',
      '1000',
      redisClient,
    );

    await runTask();

    expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
      message: expect.stringContaining('Failed to delete stale bid level for'),
      side: OrderSide.BUY,
      humanPrice: '45100',
    }));
  });
});
