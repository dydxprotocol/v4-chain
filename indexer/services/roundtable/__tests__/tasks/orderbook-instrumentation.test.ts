import { stats } from '@dydxprotocol-indexer/base';
import {
  dbHelpers,
  OrderSide,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  OpenOrdersCache,
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
    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
      {},
    );

    await orderbookInstrumentationTask();

    perpetualMarkets.forEach((perpetualMarket: PerpetualMarketFromDatabase) => {
      const tags: Object = { clob_pair_id: perpetualMarket.clobPairId };
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.open_orders_count',
        0,
        tags,
      );
    });

    // Only open orders stat should have been sent
    expect(stats.gauge).toHaveBeenCalledTimes(perpetualMarkets.length);
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
          OrderbookLevelsCache.updatePriceLevel({
            ticker: perpetualMarket.ticker,
            side: OrderSide.BUY,
            humanPrice: '45100',
            sizeDeltaInQuantums: '2000',
            client: redisClient,
          }),
          OrderbookLevelsCache.updatePriceLevel({
            ticker: perpetualMarket.ticker,
            side: OrderSide.BUY,
            humanPrice: '45200',
            sizeDeltaInQuantums: '1500',
            client: redisClient,
          }),
          OrderbookLevelsCache.updatePriceLevel({
            ticker: perpetualMarket.ticker,
            side: OrderSide.BUY,
            humanPrice: '45500',
            sizeDeltaInQuantums: '500',
            client: redisClient,
          }),
          OrderbookLevelsCache.updatePriceLevel({
            ticker: perpetualMarket.ticker,
            side: OrderSide.SELL,
            humanPrice: '45000',
            sizeDeltaInQuantums: '500',
            client: redisClient,
          }),
          OrderbookLevelsCache.updatePriceLevel({
            ticker: perpetualMarket.ticker,
            side: OrderSide.SELL,
            humanPrice: '45300',
            sizeDeltaInQuantums: '3000',
            client: redisClient,
          }),
          OrderbookLevelsCache.updatePriceLevel({
            ticker: perpetualMarket.ticker,
            side: OrderSide.SELL,
            humanPrice: '45400',
            sizeDeltaInQuantums: '3500',
            client: redisClient,
          }),
          OpenOrdersCache.addOpenOrder('orderUuid', perpetualMarket.clobPairId, redisClient),
        ]);
      },
      ));

    await orderbookInstrumentationTask();

    perpetualMarkets.forEach((perpetualMarket: PerpetualMarketFromDatabase) => {
      const tags: Object = { clob_pair_id: perpetualMarket.clobPairId };

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
        'roundtable.uncrossed_orderbook.best_ask_human',
        45300,
        tags,
      );
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.crossed_orderbook.best_ask_human',
        45000,
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

      // Check for open order count being gauged
      expect(stats.gauge).toHaveBeenCalledWith(
        'roundtable.open_orders_count',
        1,
        tags,
      );
    });
  });
});
