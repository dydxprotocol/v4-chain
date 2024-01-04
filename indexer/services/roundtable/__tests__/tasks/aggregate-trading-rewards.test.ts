import {
  BlockTable,
  TradingRewardAggregationCreateObject,
  TradingRewardAggregationFromDatabase,
  TradingRewardAggregationPeriod,
  TradingRewardAggregationTable,
  dbHelpers, testConstants, testMocks,
} from '@dydxprotocol-indexer/postgres';
import generateTaskFromPeriod, { AggregateTradingReward } from '../../src/tasks/aggregate-trading-rewards';
import { logger } from '@dydxprotocol-indexer/base';
import { DateTime, Interval } from 'luxon';
import { UTC_OPTIONS } from '../../src/lib/constants';
import { deleteAllAsync } from '@dydxprotocol-indexer/redis/build/src/helpers/redis';
import { redisClient } from '../../src/helpers/redis';
import { AggregateTradingRewardsProcessedCache } from '@dydxprotocol-indexer/redis';

describe('aggregate-trading-rewards', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    jest.spyOn(logger, 'error');
    jest.spyOn(logger, 'info');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await deleteAllAsync(redisClient);
    jest.resetAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  const startedAt: DateTime = testConstants.createdDateTime.startOf('month').toUTC();
  const startedAt2: DateTime = startedAt.plus({ month: 1 });
  const endedAt2: DateTime = startedAt2.plus({ month: 1 });
  const defaultMonthlyTradingRewardAggregation: TradingRewardAggregationCreateObject = {
    address: testConstants.defaultAddress,
    startedAt: startedAt.toISO(),
    startedAtHeight: testConstants.defaultBlock.blockHeight,
    endedAt: startedAt2.toISO(),
    endedAtHeight: '10000', // ignored field for the purposes of this test
    period: TradingRewardAggregationPeriod.MONTHLY,
    amount: '10',
  };
  const defaultMonthlyTradingRewardAggregation2: TradingRewardAggregationCreateObject = {
    address: testConstants.defaultAddress,
    startedAt: startedAt2.toISO(),
    startedAtHeight: testConstants.defaultBlock2.blockHeight,
    period: TradingRewardAggregationPeriod.MONTHLY,
    amount: '10',
  };

  describe('getTradingRewardDataToProcessInterval', () => {
    it.each([
      TradingRewardAggregationPeriod.DAILY,
      TradingRewardAggregationPeriod.WEEKLY,
      TradingRewardAggregationPeriod.MONTHLY,
    ])('Successfully returns undefined if there are no blocks in the database', async (
      period: TradingRewardAggregationPeriod,
    ) => {
      await dbHelpers.clearData();
      const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(period);
      const interval:
      Interval | undefined = await aggregateTradingReward.getTradingRewardDataToProcessInterval();

      expect(interval).toBeUndefined();
      expect(logger.info).toHaveBeenCalledWith(
        expect.objectContaining({
          message:
            'Unable to aggregate trading rewards because there are no blocks in the database.',
        }),
      );
    });

    it.each([
      TradingRewardAggregationPeriod.DAILY,
      TradingRewardAggregationPeriod.WEEKLY,
      TradingRewardAggregationPeriod.MONTHLY,
    ])('Successfully returns first block time if cache is empty and no aggregations', async (
      period: TradingRewardAggregationPeriod,
    ) => {
      const firstBlockTime: DateTime = DateTime.fromISO(
        testConstants.defaultBlock.time,
        UTC_OPTIONS,
      ).toUTC();
      await createBlockWithTime(firstBlockTime.plus({ hours: 1 }));
      const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(period);
      const interval:
      Interval | undefined = await aggregateTradingReward.getTradingRewardDataToProcessInterval();

      expect(interval).not.toBeUndefined();
      expect(interval).toEqual(Interval.fromDateTimes(
        firstBlockTime,
        firstBlockTime.plus({ hours: 1 })),
      );
    });

    it(
      'Deletes incomplete aggregations when cache is empty and only one incomplete aggregations exist',
      async () => {
        await Promise.all([
          TradingRewardAggregationTable.create(defaultMonthlyTradingRewardAggregation2),
          TradingRewardAggregationTable.create({
            ...defaultMonthlyTradingRewardAggregation2,
            period: TradingRewardAggregationPeriod.WEEKLY,
          }),
        ]);
        const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
          TradingRewardAggregationPeriod.MONTHLY,
        );
        const interval:
        Interval | undefined = await aggregateTradingReward.getTradingRewardDataToProcessInterval();

        const firstBlockTime: DateTime = DateTime.fromISO(
          testConstants.defaultBlock.time,
          UTC_OPTIONS,
        ).toUTC();
        expect(interval).toEqual(Interval.fromDateTimes(
          firstBlockTime,
          firstBlockTime.plus({ hours: 1 })),
        );

        const aggregations:
        TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll(
          {},
          [],
        );
        expect(aggregations.length).toEqual(1);
      },
    );

    it(
      'Deletes incomplete aggregations when cache is empty and multiple aggregations exist',
      async () => {
        await Promise.all([
          TradingRewardAggregationTable.create(defaultMonthlyTradingRewardAggregation),
          TradingRewardAggregationTable.create(defaultMonthlyTradingRewardAggregation2),
          TradingRewardAggregationTable.create({
            ...defaultMonthlyTradingRewardAggregation2,
            period: TradingRewardAggregationPeriod.WEEKLY,
          }),
          createBlockWithTime(startedAt2.plus({ hours: 1 })),
        ]);
        const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
          TradingRewardAggregationPeriod.MONTHLY,
        );
        const interval:
        Interval | undefined = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).not.toBeUndefined();
        expect(interval).toEqual(Interval.fromDateTimes(startedAt2, startedAt2.plus({ hours: 1 })));

        const aggregations:
        TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll(
          {},
          [],
        );
        expect(aggregations.length).toEqual(2);
      });

    it(
      'Successfully returns interval when cache is populated and not enough blocks have been processed',
      async () => {
        await Promise.all([
          TradingRewardAggregationTable.create(defaultMonthlyTradingRewardAggregation),
          TradingRewardAggregationTable.create({
            ...defaultMonthlyTradingRewardAggregation2,
            endedAt: endedAt2.toISO(),
            endedAtHeight: '10000', // ignored field for the purposes of this test
          }),
          createBlockWithTime(endedAt2.plus({ seconds: 59 })),
          AggregateTradingRewardsProcessedCache.setProcessedTime(
            TradingRewardAggregationPeriod.MONTHLY,
            endedAt2.toISO(),
            redisClient,
          ),
        ]);

        const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
          TradingRewardAggregationPeriod.MONTHLY,
        );
        const interval:
        Interval | undefined = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).toEqual(Interval.fromDateTimes(endedAt2, endedAt2));

        const aggregations:
        TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll(
          {},
          [],
        );
        expect(aggregations.length).toEqual(2);
      });

    it(
      'Successfully returns interval when cache is populated and >= minutes of blocks are unprocessed',
      async () => {
        await Promise.all([
          TradingRewardAggregationTable.create(defaultMonthlyTradingRewardAggregation),
          TradingRewardAggregationTable.create({
            ...defaultMonthlyTradingRewardAggregation2,
            endedAt: endedAt2.toISO(),
            endedAtHeight: '10000', // ignored field for the purposes of this test
          }),
          createBlockWithTime(endedAt2.plus({ seconds: 61 })),
          AggregateTradingRewardsProcessedCache.setProcessedTime(
            TradingRewardAggregationPeriod.MONTHLY,
            endedAt2.toISO(),
            redisClient,
          ),
        ]);

        const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
          TradingRewardAggregationPeriod.MONTHLY,
        );
        const interval:
        Interval | undefined = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).toEqual(Interval.fromDateTimes(endedAt2, endedAt2.plus({ minutes: 1 })));

        const aggregations:
        TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll(
          {},
          [],
        );
        expect(aggregations.length).toEqual(2);
      });

    it(
      'Successfully returns interval when cache is populated and >= 1hr of blocks are unprocessed',
      async () => {
        await Promise.all([
          TradingRewardAggregationTable.create(defaultMonthlyTradingRewardAggregation),
          TradingRewardAggregationTable.create({
            ...defaultMonthlyTradingRewardAggregation2,
            endedAt: endedAt2.toISO(),
            endedAtHeight: '10000', // ignored field for the purposes of this test
          }),
          createBlockWithTime(endedAt2.plus({ minutes: 61 })),
          AggregateTradingRewardsProcessedCache.setProcessedTime(
            TradingRewardAggregationPeriod.MONTHLY,
            endedAt2.toISO(),
            redisClient,
          ),
        ]);

        const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
          TradingRewardAggregationPeriod.MONTHLY,
        );
        const interval:
        Interval | undefined = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).toEqual(Interval.fromDateTimes(endedAt2, endedAt2.plus({ hour: 1 })));

        const aggregations:
        TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll(
          {},
          [],
        );
        expect(aggregations.length).toEqual(2);
      });

    it(
      'Successfully returns interval when cache is populated close to EOD',
      async () => {
        await Promise.all([
          TradingRewardAggregationTable.create(defaultMonthlyTradingRewardAggregation),
          TradingRewardAggregationTable.create({
            ...defaultMonthlyTradingRewardAggregation2,
            endedAt: endedAt2.toISO(),
            endedAtHeight: '10000', // ignored field for the purposes of this test
          }),
          createBlockWithTime(endedAt2.plus({ hours: 25 })),
          AggregateTradingRewardsProcessedCache.setProcessedTime(
            TradingRewardAggregationPeriod.MONTHLY,
            endedAt2.plus({ hours: 23, minutes: 55 }).toISO(),
            redisClient,
          ),
        ]);

        const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
          TradingRewardAggregationPeriod.MONTHLY,
        );
        const interval:
        Interval | undefined = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).toEqual(Interval.fromDateTimes(
          endedAt2.plus({ hours: 23, minutes: 55 }),
          endedAt2.plus({ days: 1 })),
        );

        const aggregations:
        TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll(
          {},
          [],
        );
        expect(aggregations.length).toEqual(2);
      });
  });

  describe('runTask', () => {
    it('Successfully logs and exits if there are no blocks in the database', async () => {
      await dbHelpers.clearData();
      await generateTaskFromPeriod(TradingRewardAggregationPeriod.MONTHLY)();

      expect(logger.info).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'No interval to aggregate trading rewards',
        }),
      );
    });
  });
});

async function createBlockWithTime(time: DateTime): Promise<void> {
  await BlockTable.create({
    blockHeight: '3',
    time: time.toISO(),
  });
}
