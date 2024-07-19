import { stats } from '@dydxprotocol-indexer/base';
import {
  StatefulOrderUpdateInfo,
  StatefulOrderUpdatesCache,
  redis,
  redisTestConstants,
} from '@dydxprotocol-indexer/redis';
import config from '../../src/config';
import removeOldOrderUpdatesTask from '../../src/tasks/remove-old-order-updates';
import { redisClient } from '../../src/helpers/redis';

describe('remove-old-order-updates', () => {
  const fakeTime: Date = new Date(2023, 9, 25, 0, 0, 0, 0);

  beforeAll(() => {
    jest.useFakeTimers().setSystemTime(fakeTime);
  });

  afterAll(() => {
    jest.resetAllMocks();
    jest.useRealTimers();
  });

  beforeEach(() => {
    jest.spyOn(stats, 'gauge');
    jest.clearAllMocks();
  });

  afterEach(async () => {
    await redis.deleteAllAsync(redisClient);
    jest.clearAllMocks();
  });

  it('succeeds with no cached order updates', async () => {
    await removeOldOrderUpdatesTask();
    expect(stats.gauge).toHaveBeenCalledWith(
      `${config.SERVICE_NAME}.remove_old_order_updates.num_removed`,
      0,
    );
  });

  it('succeeds with no old cached order updates', async () => {
    await StatefulOrderUpdatesCache.addStatefulOrderUpdate(
      redisTestConstants.defaultOrderUuidGoodTilBlockTime,
      redisTestConstants.orderUpdate.orderUpdate,
      fakeTime.getTime() - 1,
      redisClient,
    );
    const existingUpdates: StatefulOrderUpdateInfo[] = await StatefulOrderUpdatesCache
      .getOldOrderUpdates(
        fakeTime.getTime(),
        redisClient,
      );
    expect(existingUpdates).toHaveLength(1);

    await removeOldOrderUpdatesTask();
    expect(stats.gauge).toHaveBeenCalledWith(
      `${config.SERVICE_NAME}.remove_old_order_updates.num_removed`,
      0,
    );

    const updatesAfterTask: StatefulOrderUpdateInfo[] = await StatefulOrderUpdatesCache
      .getOldOrderUpdates(
        fakeTime.getTime(),
        redisClient,
      );
    expect(updatesAfterTask).toEqual(existingUpdates);
  });

  it('succeeds with no old cached order updates', async () => {
    await StatefulOrderUpdatesCache.addStatefulOrderUpdate(
      redisTestConstants.defaultOrderUuidGoodTilBlockTime,
      redisTestConstants.orderUpdate.orderUpdate,
      fakeTime.getTime() - 1,
      redisClient,
    );
    const existingUpdates: StatefulOrderUpdateInfo[] = await StatefulOrderUpdatesCache
      .getOldOrderUpdates(
        fakeTime.getTime(),
        redisClient,
      );
    expect(existingUpdates).toHaveLength(1);

    await removeOldOrderUpdatesTask();
    expect(stats.gauge).toHaveBeenCalledWith(
      `${config.SERVICE_NAME}.remove_old_order_updates.num_removed`,
      0,
    );

    const updatesAfterTask: StatefulOrderUpdateInfo[] = await StatefulOrderUpdatesCache
      .getOldOrderUpdates(
        fakeTime.getTime(),
        redisClient,
      );
    expect(updatesAfterTask).toEqual(existingUpdates);
  });

  it('succeeds removing old cached order update', async () => {
    await StatefulOrderUpdatesCache.addStatefulOrderUpdate(
      redisTestConstants.defaultOrderUuidGoodTilBlockTime,
      redisTestConstants.orderUpdate.orderUpdate,
      fakeTime.getTime() - config.OLD_CACHED_ORDER_UPDATES_WINDOW_MS,
      redisClient,
    );
    const existingUpdates: StatefulOrderUpdateInfo[] = await StatefulOrderUpdatesCache
      .getOldOrderUpdates(
        fakeTime.getTime(),
        redisClient,
      );
    expect(existingUpdates).toHaveLength(1);

    await removeOldOrderUpdatesTask();
    expect(stats.gauge).toHaveBeenCalledWith(
      `${config.SERVICE_NAME}.remove_old_order_updates.num_removed`,
      1,
    );

    const updatesAfterTask: StatefulOrderUpdateInfo[] = await StatefulOrderUpdatesCache
      .getOldOrderUpdates(
        fakeTime.getTime(),
        redisClient,
      );
    expect(updatesAfterTask).toHaveLength(0);
  });
});
