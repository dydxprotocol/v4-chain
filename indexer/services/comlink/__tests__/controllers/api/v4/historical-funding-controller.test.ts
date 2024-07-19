import {
  BlockTable,
  dbHelpers,
  FundingIndexUpdatesCreateObject,
  FundingIndexUpdatesTable,
  perpetualMarketRefresher,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { HistoricalFundingResponseObject, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';

describe('historical-funding-controller#V4', () => {
  const fundingIndexUpdate2: FundingIndexUpdatesCreateObject = {
    ...testConstants.defaultFundingIndexUpdate,
    oraclePrice: '1000000',
    eventId: testConstants.defaultTendermintEventId2,
    effectiveAtHeight: '5',
  };
  const fundingIndexUpdate3: FundingIndexUpdatesCreateObject = {
    ...testConstants.defaultFundingIndexUpdate,
    perpetualId: testConstants.defaultPerpetualMarket2.id,
    oraclePrice: '100',
  };
  const fundingIndexUpdate4: FundingIndexUpdatesCreateObject = {
    ...testConstants.defaultFundingIndexUpdate,
    perpetualId: testConstants.defaultPerpetualMarket2.id,
    eventId: testConstants.defaultTendermintEventId2,
    oraclePrice: '200',
    effectiveAtHeight: '5',
  };

  const expectedFunding1: HistoricalFundingResponseObject = {
    ticker: testConstants.defaultPerpetualMarket.ticker,
    rate: testConstants.defaultFundingIndexUpdate.rate,
    price: '1000000',
    effectiveAt: testConstants.defaultFundingIndexUpdate.effectiveAt,
    effectiveAtHeight: '5',
  };
  const expectedFunding2: HistoricalFundingResponseObject = {
    ...expectedFunding1,
    price: testConstants.defaultFundingIndexUpdate.oraclePrice,
    effectiveAtHeight: testConstants.defaultFundingIndexUpdate.effectiveAtHeight,
  };
  const expectedFunding3: HistoricalFundingResponseObject = {
    ...expectedFunding1,
    ticker: testConstants.defaultPerpetualMarket2.ticker,
    price: '200',
    effectiveAtHeight: '5',
  };
  const expectedFunding4: HistoricalFundingResponseObject = {
    ...expectedFunding1,
    ticker: testConstants.defaultPerpetualMarket2.ticker,
    price: '100',
    effectiveAtHeight: testConstants.defaultFundingIndexUpdate.effectiveAtHeight,
  };

  beforeAll(async () => {
    await dbHelpers.migrate();
    await testMocks.seedData();
    await BlockTable.create({
      blockHeight: '5',
      time: '2000-05-25T00:00:00.000Z',
    });

    await Promise.all([
      FundingIndexUpdatesTable.create(testConstants.defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create(fundingIndexUpdate2),
      FundingIndexUpdatesTable.create(fundingIndexUpdate3),
      FundingIndexUpdatesTable.create(fundingIndexUpdate4),
    ]);
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterAll(async () => {
    await dbHelpers.clearData();
    await dbHelpers.teardown();
  });

  it('Get /historicalFunding', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFunding/${testConstants.defaultPerpetualMarket.ticker}`,
    });
    const response2: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFunding/${testConstants.defaultPerpetualMarket2.ticker}`,
    });

    expect(response.body.historicalFunding).toEqual(
      [
        expect.objectContaining(expectedFunding1),
        expect.objectContaining(expectedFunding2),
      ],
    );
    expect(response2.body.historicalFunding).toEqual(
      [
        expect.objectContaining(expectedFunding3),
        expect.objectContaining(expectedFunding4),
      ],
    );
  });

  it('Get /historicalFunding respects effectiveBeforeOrAt and effectiveBeforeOrAtHeight field', async () => {
    const blockHeight: string = '3';
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFunding/${testConstants.defaultPerpetualMarket.ticker}?` +
        `effectiveBeforeOrAtHeight=${blockHeight}`,
    });

    expect(response.body.historicalFunding).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedFunding2),
      ]),
    );
  });

  it('Returns 400 with unknown ticker', async () => {
    const unknownTicker = 'DYDX-USD';
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFunding/${unknownTicker}`,
      expectedStatus: 400,
    });

    expect(response.body.errors[0]).toEqual(expect.objectContaining({
      msg: 'ticker must be a valid ticker (BTC-USD, etc)',
    }));
  });
});
