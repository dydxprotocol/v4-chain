import { TradingRewardAggregationFromDatabase, TradingRewardAggregationPeriod } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  createdDateTime,
  defaultSubaccountId,
  defaultTradingRewardAggregation,
  defaultTradingRewardAggregationId,
} from '../helpers/constants';
import * as TradingRewardAggregationTable from '../../src/stores/trading-reward-aggregation-table';
import { BlockTable } from '../../src';
import { seedData } from '../helpers/mock-generators';
import { denomToHumanReadableConversion } from '../helpers/conversion-helpers';

describe('TradingRewardAggregation store', () => {
  beforeAll(async () => {
    await migrate();
  });

  beforeEach(async () => {
    await seedData();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a TradingRewardAggregation', async () => {
    await TradingRewardAggregationTable.create(defaultTradingRewardAggregation);
  });

  it('Successfully finds all TradingRewardAggregations', async () => {
    await Promise.all([
      TradingRewardAggregationTable.create(defaultTradingRewardAggregation),
      TradingRewardAggregationTable.create({
        ...defaultTradingRewardAggregation,
        startedAtHeight: '1',
      }),
    ]);

    const tradingRewardAggregations:
    TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(tradingRewardAggregations.length).toEqual(2);
    expect(tradingRewardAggregations[0]).toEqual(expect.objectContaining({
      ...defaultTradingRewardAggregation,
      startedAtHeight: '1',
    }));
    expect(tradingRewardAggregations[1]).toEqual(
      expect.objectContaining(defaultTradingRewardAggregation),
    );
  });

  it('Successfully finds latest monthly TradingRewardAggregation', async () => {
    await Promise.all([
      BlockTable.create({
        blockHeight: '100',
        time: createdDateTime.toISO(),
      }),
    ]);

    await Promise.all([
      TradingRewardAggregationTable.create({
        ...defaultTradingRewardAggregation,
        period: TradingRewardAggregationPeriod.MONTHLY,
      }),
      TradingRewardAggregationTable.create({
        ...defaultTradingRewardAggregation,
        startedAtHeight: '100',
        period: TradingRewardAggregationPeriod.MONTHLY,
      }),
    ]);

    const tradingRewardAggregation:
    TradingRewardAggregationFromDatabase | undefined = await TradingRewardAggregationTable
      .getLatestAggregatedTradeReward(TradingRewardAggregationPeriod.MONTHLY);

    expect(tradingRewardAggregation).toEqual(
      expect.objectContaining({
        ...defaultTradingRewardAggregation,
        startedAtHeight: '100',
        period: TradingRewardAggregationPeriod.MONTHLY,
      }),
    );
  });

  it('Successfully finds a TradingRewardAggregation', async () => {
    await TradingRewardAggregationTable.create(defaultTradingRewardAggregation);

    const tradingRewardAggregation:
    TradingRewardAggregationFromDatabase | undefined = await TradingRewardAggregationTable.findById(
      defaultTradingRewardAggregationId,
    );

    expect(tradingRewardAggregation).toEqual(
      expect.objectContaining(defaultTradingRewardAggregation),
    );
  });

  it('Successfully returns undefined when updating a nonexistent TradingRewardAggregation', async () => {
    const fakeUpdate:
    TradingRewardAggregationFromDatabase | undefined = await TradingRewardAggregationTable.update({
      id: defaultSubaccountId,
    });
    expect(fakeUpdate).toBeUndefined();
  });

  it('Successfully updates an existing TradingRewardAggregation', async () => {
    await TradingRewardAggregationTable.create(defaultTradingRewardAggregation);

    const amount: string = denomToHumanReadableConversion(100000);
    const endedAt: string = '2021-01-01T00:00:00.000Z';
    const endedAtHeight: string = '1000';
    const update:
    TradingRewardAggregationFromDatabase | undefined = await TradingRewardAggregationTable.update({
      id: defaultTradingRewardAggregationId,
      endedAt,
      endedAtHeight,
      amount,
    });
    expect(update).toEqual({
      ...defaultTradingRewardAggregation,
      id: defaultTradingRewardAggregationId,
      endedAt,
      endedAtHeight,
      amount,
    });
  });

  it('Successfully deleted trading reward aggregations after a certain height', async () => {
    await Promise.all([
      BlockTable.create({
        blockHeight: '100',
        time: createdDateTime.toISO(),
      }),
      BlockTable.create({
        blockHeight: '101',
        time: createdDateTime.toISO(),
      }),
    ]);

    await Promise.all([
      TradingRewardAggregationTable.create(defaultTradingRewardAggregation),
      TradingRewardAggregationTable.create({
        ...defaultTradingRewardAggregation,
        startedAtHeight: '100',
      }),
      TradingRewardAggregationTable.create({
        ...defaultTradingRewardAggregation,
        startedAtHeight: '101',
      }),
    ]);

    await TradingRewardAggregationTable.deleteAll({
      startedAtHeightOrAfter: '100',
    });

    const tradingRewardAggregations:
    TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll(
      {},
      [],
    );

    expect(tradingRewardAggregations.length).toEqual(1);
    expect(tradingRewardAggregations[0]).toEqual(
      expect.objectContaining(defaultTradingRewardAggregation),
    );
  });
});
