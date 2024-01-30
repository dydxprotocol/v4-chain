import {
  BlockTable,
  IsoString,
  TradingRewardAggregationCreateObject,
  TradingRewardAggregationFromDatabase,
  TradingRewardAggregationPeriod,
  TradingRewardAggregationTable,
  TradingRewardCreateObject,
  TradingRewardTable,
  dbHelpers,
  testConstants,
  testConversionHelpers,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import generateTaskFromPeriod, { AggregateTradingReward } from '../../src/tasks/aggregate-trading-rewards';
import { logger } from '@dydxprotocol-indexer/base';
import { DateTime, Interval } from 'luxon';
import { UTC_OPTIONS } from '../../src/lib/constants';
import { redisClient } from '../../src/helpers/redis';
import { AggregateTradingRewardsProcessedCache, redis } from '@dydxprotocol-indexer/redis';
import config from '../../src/config';

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
    await redis.deleteAllAsync(redisClient);
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
    amount: testConversionHelpers.convertToDenomScale('10'),
  };
  const defaultMonthlyTradingRewardAggregation2: TradingRewardAggregationCreateObject = {
    address: testConstants.defaultAddress,
    startedAt: startedAt2.toISO(),
    startedAtHeight: testConstants.defaultBlock2.blockHeight,
    period: TradingRewardAggregationPeriod.MONTHLY,
    amount: testConversionHelpers.convertToDenomScale('10'),
  };

  describe('maybeDeleteIncompleteAggregatedTradingReward', () => {
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
        await aggregateTradingReward.maybeDeleteIncompleteAggregatedTradingReward();
        await validateNumberOfAggregations(1);
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
        await aggregateTradingReward.maybeDeleteIncompleteAggregatedTradingReward();

        await validateNumberOfAggregations(2);
      },
    );
  });

  describe('getTradingRewardDataToProcessInterval', () => {
    it.each([
      TradingRewardAggregationPeriod.DAILY,
      TradingRewardAggregationPeriod.WEEKLY,
      TradingRewardAggregationPeriod.MONTHLY,
    ])('Throws error if there are no blocks in the database', async (
      period: TradingRewardAggregationPeriod,
    ) => {
      await dbHelpers.clearData();
      const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(period);
      await expect(aggregateTradingReward.getTradingRewardDataToProcessInterval())
        .rejects.toEqual(new Error('Unable to find latest block'));
    });

    it.each([
      TradingRewardAggregationPeriod.DAILY,
      TradingRewardAggregationPeriod.WEEKLY,
      TradingRewardAggregationPeriod.MONTHLY,
    ])('Throws error if cache is empty no aggregations, and no trading rewards', async (
      period: TradingRewardAggregationPeriod,
    ) => {
      const firstBlockTime: DateTime = DateTime.fromISO(
        testConstants.defaultBlock.time,
        UTC_OPTIONS,
      ).toUTC();
      await createBlockWithTime(firstBlockTime.plus({ hours: 1 }));
      const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(period);
      await expect(aggregateTradingReward.getTradingRewardDataToProcessInterval())
        .rejects.toEqual(new Error('No trading rewards in database'));
    });

    it(
      'Throws error interval when cache is empty and no aggregation for the period and trading reward',
      async () => {
        await TradingRewardAggregationTable.create({
          ...defaultMonthlyTradingRewardAggregation,
          period: TradingRewardAggregationPeriod.WEEKLY,
        });
        const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
          TradingRewardAggregationPeriod.MONTHLY,
        );
        await expect(aggregateTradingReward.getTradingRewardDataToProcessInterval())
          .rejects.toEqual(new Error('No trading rewards in database'));
      },
    );

    it(
      'Successfully returns interval when cache is empty and no aggregation and trading reward exists',
      async () => {
        await Promise.all([
          TradingRewardTable.create({
            address: testConstants.defaultAddress,
            blockTime: startedAt.plus({ minutes: 2, seconds: 20 }).toISO(), // random constants
            blockHeight: '1',
            amount: '10',
          }),
          createBlockWithTime(startedAt.plus({ hours: 10 })),
        ]);
        const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
          TradingRewardAggregationPeriod.MONTHLY,
        );
        const interval:
        Interval = await aggregateTradingReward.getTradingRewardDataToProcessInterval();

        expect(interval).toEqual(Interval.fromDateTimes(
          startedAt.plus({ minutes: 2, seconds: 20 }),
          startedAt.plus({
            minutes: 2,
            seconds: 20,
            milliseconds: config.AGGREGATE_TRADING_REWARDS_MAX_INTERVAL_SIZE_MS,
          }),
        ));
      },
    );

    it(
      'Successfully returns interval when cache is empty and a complete aggregations exist',
      async () => {
        await Promise.all([
          TradingRewardAggregationTable.create(defaultMonthlyTradingRewardAggregation),
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
        Interval = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).not.toBeUndefined();
        expect(interval).toEqual(Interval.fromDateTimes(
          startedAt2,
          startedAt2.plus({ milliseconds: config.AGGREGATE_TRADING_REWARDS_MAX_INTERVAL_SIZE_MS }),
        ));

        await validateNumberOfAggregations(2);
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
        Interval = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).toEqual(Interval.fromDateTimes(endedAt2, endedAt2));

        await validateNumberOfAggregations(2);
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
        Interval = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).toEqual(Interval.fromDateTimes(endedAt2, endedAt2.plus({ minutes: 1 })));

        await validateNumberOfAggregations(2);
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
        Interval = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).toEqual(Interval.fromDateTimes(
          endedAt2,
          endedAt2.plus({ milliseconds: config.AGGREGATE_TRADING_REWARDS_MAX_INTERVAL_SIZE_MS }),
        ));

        await validateNumberOfAggregations(2);
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
        Interval = await aggregateTradingReward.getTradingRewardDataToProcessInterval();
        expect(interval).toEqual(Interval.fromDateTimes(
          endedAt2.plus({ hours: 23, minutes: 55 }),
          endedAt2.plus({ days: 1 })),
        );

        await validateNumberOfAggregations(2);
      });
  });

  describe('runTask', () => {
    beforeEach(async () => {
      // In the tests below we have block at the following time:
      // - 3 at 1/1/24 00:00:30
      // - 4 at 1/1/24 23:00:30
      // - 3 at 1/2/24 00:00:30
      await Promise.all([
        BlockTable.create({
          blockHeight: '3',
          time: thirdBlockTime,
        }),
        BlockTable.create({
          blockHeight: '4',
          time: fourthBlockTime,
        }),
        BlockTable.create({
          blockHeight: '5',
          time: fifthBlockTime,
        }),
      ]);
    });

    const thirdBlockTime: IsoString = startedAt.plus({ seconds: 30 }).toISO();
    const fourthBlockTime: IsoString = startedAt.plus({ hours: 23, seconds: 30 }).toISO();
    const fifthBlockTime: IsoString = startedAt.plus({ day: 1, seconds: 30 }).toISO();
    const defaultTradingReward: TradingRewardCreateObject = {
      address: testConstants.defaultAddress,
      blockTime: thirdBlockTime,
      blockHeight: '3',
      amount: testConversionHelpers.convertToDenomScale('10'),
    };

    const intervalToBeProcessed: Interval = Interval.fromDateTimes(
      startedAt,
      startedAt.plus({ hours: 1 }),
    );
    const defaultCreatedTradingRewardAggregation: TradingRewardAggregationFromDatabase = {
      id: TradingRewardAggregationTable.uuid(
        testConstants.defaultAddress,
        TradingRewardAggregationPeriod.DAILY,
        '3',
      ),
      address: testConstants.defaultAddress,
      startedAt: startedAt.toISO(),
      startedAtHeight: '3',
      period: TradingRewardAggregationPeriod.DAILY,
      amount: testConversionHelpers.convertToDenomScale('10'),
    };

    it('Successfully logs and exits if there are no blocks in the database', async () => {
      await dbHelpers.clearData();
      await expect(generateTaskFromPeriod(TradingRewardAggregationPeriod.DAILY)())
        .rejects.toEqual(new Error('Unable to find latest block'));
    });

    it('Successfully creates new aggregations', async () => {
      await Promise.all([
        TradingRewardTable.create(defaultTradingReward),
        AggregateTradingRewardsProcessedCache.setProcessedTime(
          TradingRewardAggregationPeriod.DAILY,
          intervalToBeProcessed.start.toISO(),
          redisClient,
        ),
        TradingRewardAggregationTable.create({
          ...defaultCreatedTradingRewardAggregation,
          startedAt: intervalToBeProcessed.start.minus({ days: 1 }).toISO(),
          startedAtHeight: '1',
        }),
      ]);

      const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
        TradingRewardAggregationPeriod.DAILY,
      );
      await aggregateTradingReward.runTask();

      await expectAggregateTradingRewardsProcessedCache(
        TradingRewardAggregationPeriod.DAILY,
        intervalToBeProcessed.end.toISO(),
      );
      await validateNumberOfAggregations(2);
      await validateAggregationWithExpectedValue(defaultCreatedTradingRewardAggregation);
    });

    it('Successfully updates aggregation amounts', async () => {
      await Promise.all([
        TradingRewardTable.create(defaultTradingReward),
        TradingRewardAggregationTable.create(defaultCreatedTradingRewardAggregation),
        AggregateTradingRewardsProcessedCache.setProcessedTime(
          TradingRewardAggregationPeriod.DAILY,
          intervalToBeProcessed.start.toISO(),
          redisClient,
        ),
      ]);

      const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
        TradingRewardAggregationPeriod.DAILY,
      );
      await aggregateTradingReward.runTask();

      await expectAggregateTradingRewardsProcessedCache(
        TradingRewardAggregationPeriod.DAILY,
        intervalToBeProcessed.end.toISO(),
      );
      await validateNumberOfAggregations(1);
      await validateAggregationWithExpectedValue({
        ...defaultCreatedTradingRewardAggregation,
        amount: testConversionHelpers.convertToDenomScale('20'),
      });
    });

    it('Successfully creates new aggregations and sets endAt and endAtHeight', async () => {
      await Promise.all([
        TradingRewardTable.create({
          ...defaultTradingReward,
          blockTime: fourthBlockTime,
          blockHeight: '4',
        }),
        AggregateTradingRewardsProcessedCache.setProcessedTime(
          TradingRewardAggregationPeriod.DAILY,
          startedAt.plus({ hours: 23 }).toISO(),
          redisClient,
        ),
      ]);

      const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
        TradingRewardAggregationPeriod.DAILY,
      );
      await aggregateTradingReward.runTask();

      await expectAggregateTradingRewardsProcessedCache(
        TradingRewardAggregationPeriod.DAILY,
        startedAt.plus({ days: 1 }).toISO(),
      );
      await validateNumberOfAggregations(1);
      await validateAggregationWithExpectedValue({
        ...defaultCreatedTradingRewardAggregation,
        endedAt: startedAt.plus({ days: 1 }).toISO(),
        endedAtHeight: '4',
      });
    });

    it('Successfully updates aggregation amount and sets endAt and endAtHeight', async () => {
      await Promise.all([
        TradingRewardTable.create({
          ...defaultTradingReward,
          blockTime: fourthBlockTime,
          blockHeight: '4',
        }),
        TradingRewardAggregationTable.create(defaultCreatedTradingRewardAggregation),
        AggregateTradingRewardsProcessedCache.setProcessedTime(
          TradingRewardAggregationPeriod.DAILY,
          startedAt.plus({ hours: 23 }).toISO(),
          redisClient,
        ),
      ]);

      const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
        TradingRewardAggregationPeriod.DAILY,
      );
      await aggregateTradingReward.runTask();

      await expectAggregateTradingRewardsProcessedCache(
        TradingRewardAggregationPeriod.DAILY,
        startedAt.plus({ days: 1 }).toISO(),
      );
      await validateNumberOfAggregations(1);
      await validateAggregationWithExpectedValue({
        ...defaultCreatedTradingRewardAggregation,
        endedAt: startedAt.plus({ days: 1 }).toISO(),
        endedAtHeight: '4',
        amount: testConversionHelpers.convertToDenomScale('20'),
      });
    });

    it('Successfully updates aggregation with no amount update and sets endAt and endAtHeight', async () => {
      await Promise.all([
        TradingRewardAggregationTable.create(defaultCreatedTradingRewardAggregation),
        AggregateTradingRewardsProcessedCache.setProcessedTime(
          TradingRewardAggregationPeriod.DAILY,
          startedAt.plus({ hours: 23 }).toISO(),
          redisClient,
        ),
      ]);

      const aggregateTradingReward: AggregateTradingReward = new AggregateTradingReward(
        TradingRewardAggregationPeriod.DAILY,
      );
      await aggregateTradingReward.runTask();

      await expectAggregateTradingRewardsProcessedCache(
        TradingRewardAggregationPeriod.DAILY,
        startedAt.plus({ days: 1 }).toISO(),
      );
      await validateNumberOfAggregations(1);
      await validateAggregationWithExpectedValue({
        ...defaultCreatedTradingRewardAggregation,
        endedAt: startedAt.plus({ days: 1 }).toISO(),
        endedAtHeight: '4',
        amount: testConversionHelpers.convertToDenomScale('10'),
      });
    });
  });
});

async function createBlockWithTime(time: DateTime): Promise<void> {
  await BlockTable.create({
    blockHeight: '3',
    time: time.toISO(),
  });
}

async function validateNumberOfAggregations(
  expectedNumberOfAggregations: number,
): Promise<void> {
  const aggregations:
  TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll({}, []);
  expect(aggregations.length).toEqual(expectedNumberOfAggregations);
}

async function validateAggregationWithExpectedValue(
  expectedAggregation: TradingRewardAggregationFromDatabase,
): Promise<void> {
  const aggregation:
  TradingRewardAggregationFromDatabase | undefined = await TradingRewardAggregationTable.findById(
    expectedAggregation.id,
  );

  expect(aggregation).toEqual(expect.objectContaining(expectedAggregation));
}

async function expectAggregateTradingRewardsProcessedCache(
  period: TradingRewardAggregationPeriod,
  processedTime: string,
): Promise<void> {
  const cacheValue: string | null = await AggregateTradingRewardsProcessedCache.getProcessedTime(
    period,
    redisClient,
  );
  expect(cacheValue).toEqual(processedTime);
}
