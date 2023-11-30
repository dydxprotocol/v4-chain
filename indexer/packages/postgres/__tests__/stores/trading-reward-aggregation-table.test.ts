import { TradingRewardAggregationFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  defaultSubaccountId,
  defaultTradingRewardAggregation,
  defaultTradingRewardAggregationId,
  defaultWallet,
} from '../helpers/constants';
import * as TradingRewardAggregationTable from '../../src/stores/trading-reward-aggregation-table';
import { WalletTable } from '../../src';
import { seedData } from '../helpers/mock-generators';

describe('TradingRewardAggregation store', () => {
  beforeAll(async () => {
    await migrate();
  });

  beforeEach(async () => {
    await seedData();
    await WalletTable.create(defaultWallet);
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

    const amount: string = '100000.00';
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
});
