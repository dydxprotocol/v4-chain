import {
  HistoricalBlockTradingReward,
  HistoricalBlockTradingRewardsResponse,
  HistoricalTradingRewardAggregation,
  HistoricalTradingRewardAggregationsResponse,
  RequestMethod,
} from '../../../../src/types';
import { getQueryString, sendRequest } from '../../../helpers/helpers';
import {
  TradingRewardCreateObject,
  TradingRewardFromDatabase,
  TradingRewardTable,
  dbHelpers,
  testConstants,
  testConversionHelpers,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { stats } from '@dydxprotocol-indexer/base';
import request from 'supertest';
import { tradingRewardToResponse } from '../../../../src/request-helpers/request-transformer';

describe('historical-block-trading-reward-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await Promise.all([
      TradingRewardTable.create(defaultTradingRewardCreate),
      TradingRewardTable.create(defaultTradingRewardCreate2),
    ]);

    const rewards: TradingRewardFromDatabase[] = await TradingRewardTable.findAll({}, []);

    defaultTradingReward = rewards[1];
    defaultTradingReward2 = rewards[0];
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  const defaultTradingRewardCreate: TradingRewardCreateObject = {
    address: testConstants.defaultAddress,
    blockTime: testConstants.defaultBlock.time,
    blockHeight: testConstants.defaultBlock.blockHeight,
    amount: testConversionHelpers.convertToDenomScale('10'),
  };
  let defaultTradingReward: TradingRewardFromDatabase;
  const defaultTradingRewardCreate2: TradingRewardCreateObject = {
    address: testConstants.defaultAddress,
    blockTime: testConstants.defaultBlock2.time,
    blockHeight: testConstants.defaultBlock2.blockHeight,
    amount: testConversionHelpers.convertToDenomScale('5'),
  };
  let defaultTradingReward2: TradingRewardFromDatabase;

  describe('GET', () => {
    it('Get /historicalBlockTradingReward/:address returns all valid rewards', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historicalBlockTradingRewards/${testConstants.defaultAddress}`,
      });

      const responseBody: HistoricalTradingRewardAggregationsResponse = response.body;
      const rewards: HistoricalTradingRewardAggregation[] = responseBody.rewards;
      expect(rewards.length).toEqual(2);
      expect(rewards[0]).toEqual(tradingRewardToResponse(
        defaultTradingReward2,
      ));
      expect(rewards[1]).toEqual(tradingRewardToResponse(
        defaultTradingReward,
      ));
    });

    it('Get /historicalBlockTradingRewards/:address returns all valid rewards with limit', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historicalBlockTradingRewards/${testConstants.defaultAddress}` +
        `?${getQueryString({ limit: 1 })}`,
      });

      const responseBody: HistoricalBlockTradingRewardsResponse = response.body;
      const rewards: HistoricalBlockTradingReward[] = responseBody.rewards;
      expect(rewards.length).toEqual(1);
      expect(rewards[0]).toEqual(tradingRewardToResponse(
        defaultTradingReward2,
      ));
    });

    it('Get /historicalBlockTradingRewards/:address returns no rewards when none exist', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/historicalBlockTradingRewards/fakeAddress',
      });

      const responseBody: HistoricalBlockTradingRewardsResponse = response.body;
      const rewards: HistoricalBlockTradingReward[] = responseBody.rewards;
      expect(rewards.length).toEqual(0);
    });

    it('Get /historicalBlockTradingRewards/:address returns rewards with blockTimeBeforeOrAt', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historicalBlockTradingRewards/${testConstants.defaultAddress}` +
        `?${getQueryString({ startingBeforeOrAt: testConstants.defaultBlock.time })}`,
      });

      const responseBody: HistoricalBlockTradingRewardsResponse = response.body;
      const rewards: HistoricalBlockTradingReward[] = responseBody.rewards;
      expect(rewards.length).toEqual(1);
      expect(rewards[0]).toEqual(tradingRewardToResponse(
        defaultTradingReward,
      ));
    });

    it('Get /historicalBlockTradingRewards/:address returns rewards with blockHeightBeforeOrAt', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historicalBlockTradingRewards/${testConstants.defaultAddress}` +
        `?${getQueryString({ startingBeforeOrAtHeight: testConstants.defaultBlock.blockHeight })}`,
      });

      const responseBody: HistoricalBlockTradingRewardsResponse = response.body;
      const rewards: HistoricalBlockTradingReward[] = responseBody.rewards;
      expect(rewards.length).toEqual(1);
      expect(rewards[0]).toEqual(tradingRewardToResponse(
        defaultTradingReward,
      ));
    });
  });
});
