import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';
import { RequestMethod } from '../../../../src/types';
import {
  BlockTable,
  dbHelpers,
  testConstants,
  FundingIndexUpdatesCreateObject,
  FundingIndexUpdatesTable,
  perpetualMarketRefresher,
  testMocks,
} from '@dydxprotocol-indexer/postgres';

describe('historical-funding-payment-controller#V4', () => {
  // TODO(Adam): Clean up this code setup, its the same in historical funding controller tests

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

  it('Get /historicalFundingPayment', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFundingPayment/${testConstants.defaultPerpetualMarket.ticker}`,
    });

    const res = {
      historicalFundingPayments: [{ ticker: testConstants.defaultPerpetualMarket.ticker }],
    };

    expect(response.body).toEqual(res);
  });
});
