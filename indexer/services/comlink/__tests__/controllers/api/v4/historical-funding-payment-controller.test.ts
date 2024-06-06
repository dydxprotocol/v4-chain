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
  PerpetualPositionTable,
  PerpetualPositionCreateObject,
} from '@dydxprotocol-indexer/postgres';
import { defaultPerpetualPosition } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('historical-funding-payment-controller#V4', () => {
  // TODO:(ADAM) - Clean up this, see which updates you actually need
  // I think you just need to call update perpetual markets

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

  // create PerpetualPositions
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

    const perpetualPosition1: PerpetualPositionCreateObject = {
      ...testConstants.defaultPerpetualPosition,
      openEventId: testConstants.defaultTendermintEventId2,
      settledFunding: '100',
    };

    const perpetualPosition2: PerpetualPositionCreateObject = {
      ...testConstants.defaultPerpetualPosition,
      openEventId: testConstants.defaultTendermintEventId3,
      settledFunding: '200',
    };

    await Promise.all([
      PerpetualPositionTable.create(defaultPerpetualPosition),
      PerpetualPositionTable.create(perpetualPosition1),
      PerpetualPositionTable.create(perpetualPosition2),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
      }),
    ]);
  });

  afterAll(async () => {
    await dbHelpers.clearData();
    await dbHelpers.teardown();
  });

  it('Get /historicalFundingPayment', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFundingPayment/${testConstants.defaultPerpetualMarket.ticker}` +
      `?address=${testConstants.defaultAddress}` +
      `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
    });

    const res = {
      historicalFundingPayments: [{
        ticker: testConstants.defaultPerpetualMarket.ticker,
        payment: '200',
        effectiveAt: '2021-01-01T00:00:00.000Z',
      }],
    };

    expect(response.body).toEqual(res);
  });

  it('Get /historicalFundingPayment returns 400 with an invalid ticker', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFundingPayment/${testConstants.invalidTicker}?` +
      `address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      expectedStatus: 400,
    });

    expect(response.body).toEqual({
      errors: [
        {
          location: 'params',
          msg: 'ticker must be a valid ticker (BTC-USD, etc)',
          param: 'ticker',
          value: testConstants.invalidTicker,
        },
      ],
    });
  });

  it('Get /historicalFundingPayment returns 400 with an invalid subaccount', async () => {
    const invalidSubaccount = 12900000000;
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFundingPayment/${testConstants.defaultPerpetualMarket.ticker}?` +
      `&address=${testConstants.defaultAddress}&subaccountNumber=${invalidSubaccount}`,
      expectedStatus: 400,
    });

    expect(response.body).toEqual({
      errors: [
        {
          location: 'query',
          msg: 'subaccountNumber must be a non-negative integer less than 128001',
          param: 'subaccountNumber',
          value: invalidSubaccount.toString(),
        },
      ],
    });
  });
});
