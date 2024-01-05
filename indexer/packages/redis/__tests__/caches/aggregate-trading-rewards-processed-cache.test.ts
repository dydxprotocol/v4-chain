import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import {
  getProcessedTime,
  setProcessedTime,
} from '../../src/caches/aggregate-trading-rewards-processed-cache';
import { IsoString, TradingRewardAggregationPeriod } from '@dydxprotocol-indexer/postgres';

describe('aggregateTradingRewardsProcessedCache', () => {
  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  it('successfully sets and gets aggregate trading rewards processed', async () => {
    const initialResult: IsoString | null = await getProcessedTime(
      TradingRewardAggregationPeriod.DAILY,
      client,
    );
    expect(initialResult).toEqual(null);

    const timestamp = '2021-01-01T00:00:00.000Z';
    await setProcessedTime(
      TradingRewardAggregationPeriod.DAILY,
      timestamp,
      client,
    );
    const result: IsoString | null = await getProcessedTime(
      TradingRewardAggregationPeriod.DAILY,
      client,
    );
    expect(result).toEqual(timestamp);
    expect(await getProcessedTime(TradingRewardAggregationPeriod.WEEKLY, client)).toEqual(null);
  });
});
