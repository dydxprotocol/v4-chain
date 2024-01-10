import {
  HistoricalTradingRewardAggregation,
  HistoricalTradingRewardAggregationsResponse,
  RequestMethod,
} from '../../../../src/types';
import { getQueryString, sendRequest } from '../../../helpers/helpers';
import {
  TradingRewardAggregationCreateObject,
  TradingRewardAggregationFromDatabase,
  TradingRewardAggregationPeriod,
  TradingRewardAggregationTable,
  dbHelpers,
  testConstants,
  testConversionHelpers,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { stats } from '@dydxprotocol-indexer/base';
import { DateTime } from 'luxon';
import request from 'supertest';
import { tradingRewardAggregationToResponse } from '../../../../src/request-helpers/request-transformer';

describe('historical-trading-reward-aggregations-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await Promise.all([
      TradingRewardAggregationTable.create(defaultCompletedTradingRewardAggregationCreate),
      TradingRewardAggregationTable.create(defaultIncompleteTradingRewardAggregationCreate),
    ]);
    const aggregations:
    TradingRewardAggregationFromDatabase[] = await TradingRewardAggregationTable.findAll({}, []);

    defaultCompletedTradingRewardAggregation = aggregations[0];
    defaultIncompleteTradingRewardAggregation = aggregations[1];
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  const startedAt: DateTime = testConstants.createdDateTime.startOf('month').toUTC();
  const startedAt2: DateTime = startedAt.plus({ month: 1 });
  const defaultCompletedTradingRewardAggregationCreate: TradingRewardAggregationCreateObject = {
    address: testConstants.defaultAddress,
    startedAt: startedAt.toISO(),
    startedAtHeight: testConstants.defaultBlock.blockHeight,
    endedAt: startedAt2.toISO(),
    endedAtHeight: '10000', // ignored field for the purposes of this test
    period: TradingRewardAggregationPeriod.MONTHLY,
    amount: testConversionHelpers.convertToDenomScale('10'),
  };
  let defaultCompletedTradingRewardAggregation: TradingRewardAggregationFromDatabase;
  const defaultIncompleteTradingRewardAggregationCreate: TradingRewardAggregationCreateObject = {
    address: testConstants.defaultAddress,
    startedAt: startedAt2.toISO(),
    startedAtHeight: testConstants.defaultBlock2.blockHeight,
    period: TradingRewardAggregationPeriod.MONTHLY,
    amount: testConversionHelpers.convertToDenomScale('20'),
  };
  let defaultIncompleteTradingRewardAggregation: TradingRewardAggregationFromDatabase;

  describe('GET', () => {
    it('Get /historicalTradingRewardAggregations/:address returns all valid aggregations', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historicalTradingRewardAggregations/${testConstants.defaultAddress}` +
        `?${getQueryString({ period: TradingRewardAggregationPeriod.MONTHLY })}`,
      });

      const responseBody: HistoricalTradingRewardAggregationsResponse = response.body;
      const rewards: HistoricalTradingRewardAggregation[] = responseBody.rewards;
      expect(rewards.length).toEqual(2);
      expect(rewards[0]).toEqual(tradingRewardAggregationToResponse(
        defaultIncompleteTradingRewardAggregation,
      ));
      expect(rewards[1]).toEqual(tradingRewardAggregationToResponse(
        defaultCompletedTradingRewardAggregation,
      ));
    });

    it('Get /historicalTradingRewardAggregations/:address returns all valid aggregations with limit', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historicalTradingRewardAggregations/${testConstants.defaultAddress}` +
        `?${getQueryString({ period: TradingRewardAggregationPeriod.MONTHLY, limit: 1 })}`,
      });

      const responseBody: HistoricalTradingRewardAggregationsResponse = response.body;
      const rewards: HistoricalTradingRewardAggregation[] = responseBody.rewards;
      expect(rewards.length).toEqual(1);
      expect(rewards[0]).toEqual(tradingRewardAggregationToResponse(
        defaultIncompleteTradingRewardAggregation,
      ));
    });

    it('Get /historicalTradingRewardAggregations/:address returns no aggregations when none exist', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historicalTradingRewardAggregations/${testConstants.defaultAddress}` +
        `?${getQueryString({ period: TradingRewardAggregationPeriod.DAILY })}`,
      });

      const responseBody: HistoricalTradingRewardAggregationsResponse = response.body;
      const rewards: HistoricalTradingRewardAggregation[] = responseBody.rewards;
      expect(rewards.length).toEqual(0);
    });

    it('Get /historicalTradingRewardAggregations/:address returns aggregations with startedAtBeforeOrAt', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historicalTradingRewardAggregations/${testConstants.defaultAddress}` +
        `?${getQueryString({
          period: TradingRewardAggregationPeriod.MONTHLY,
          startingBeforeOrAt: startedAt.toISO(),
        })}`,
      });

      const responseBody: HistoricalTradingRewardAggregationsResponse = response.body;
      const rewards: HistoricalTradingRewardAggregation[] = responseBody.rewards;
      expect(rewards.length).toEqual(1);
      expect(rewards[0]).toEqual(tradingRewardAggregationToResponse(
        defaultCompletedTradingRewardAggregation,
      ));
    });

    it('Get /historicalTradingRewardAggregations/:address returns aggregations with startedAtHeightBeforeOrAt', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historicalTradingRewardAggregations/${testConstants.defaultAddress}` +
        `?${getQueryString({
          period: TradingRewardAggregationPeriod.MONTHLY,
          startingBeforeOrAtHeight: testConstants.defaultBlock.blockHeight,
        })}`,
      });

      const responseBody: HistoricalTradingRewardAggregationsResponse = response.body;
      const rewards: HistoricalTradingRewardAggregation[] = responseBody.rewards;
      expect(rewards.length).toEqual(1);
      expect(rewards[0]).toEqual(tradingRewardAggregationToResponse(
        defaultCompletedTradingRewardAggregation,
      ));
    });
  });
});
