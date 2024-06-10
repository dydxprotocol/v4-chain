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
  PerpetualPositionStatus,
} from '@dydxprotocol-indexer/postgres';
import {
  defaultPerpetualPosition,
  createdDateTime,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import { DateTime } from 'luxon';

describe('historical-funding-payment-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await testMocks.seedData();
    await BlockTable.create({
      blockHeight: '5',
      time: '2000-05-25T00:00:00.000Z',
    });

    const fundingIndexUpdate: FundingIndexUpdatesCreateObject = {
      ...testConstants.defaultFundingIndexUpdate,
      oraclePrice: '1000000',
      eventId: testConstants.defaultTendermintEventId2,
      fundingIndex: '1000',
      effectiveAt: DateTime.utc().plus({ days: 1 }).toISO(),
      effectiveAtHeight: '5',
    };

    await Promise.all([
      FundingIndexUpdatesTable.create(testConstants.defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create(fundingIndexUpdate),
    ]);
    await perpetualMarketRefresher.updatePerpetualMarkets();

    const closedPosition1: PerpetualPositionCreateObject = {
      ...testConstants.defaultClosedPerpetualPosition,
      openEventId: testConstants.defaultTendermintEventId2,
      settledFunding: '100',
    };

    const closedPosition2: PerpetualPositionCreateObject = {
      ...testConstants.defaultClosedPerpetualPosition,
      openEventId: testConstants.defaultTendermintEventId3,
      settledFunding: '200',
    };

    await Promise.all([
      PerpetualPositionTable.create(defaultPerpetualPosition),
      PerpetualPositionTable.create(closedPosition1),
      PerpetualPositionTable.create(closedPosition2),
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

  it('Get /historicalFundingPayment returns settled and unsettled funding payments', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFundingPayment/${testConstants.defaultPerpetualMarket.ticker}` +
      `?address=${testConstants.defaultAddress}` +
      `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
    });

    const res = {
      ticker: testConstants.defaultPerpetualMarket.ticker,
      fundingPayments: [
        {
          // Default position size is 10
          // Default funding index is 10050 (last updated)
          // Second funding event index is 1000 (latest)
          // 10 * (10050 - 1000) = 90500
          payment: '90500',
          effectiveAt: createdDateTime.toISO(),
        },
        {
          payment: '200',
          effectiveAt: createdDateTime.toISO(),
        },
        {
          payment: '100',
          effectiveAt: createdDateTime.toISO(),
        },
      ],
    };

    expect(response.body).toEqual(res);
  });

  it('Get /historicalFundingPayment returns funding for liquidated positions', async () => {
    await PerpetualPositionTable.create({
      ...testConstants.defaultClosedPerpetualPosition,
      openEventId: testConstants.defaultTendermintEventId4,
      status: PerpetualPositionStatus.LIQUIDATED,
    });

    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historicalFundingPayment/${testConstants.defaultPerpetualMarket.ticker}` +
       `?address=${testConstants.defaultAddress}` +
       `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
    });

    const res = {
      ticker: testConstants.defaultPerpetualMarket.ticker,
      fundingPayments: [
        {
          // Default position size is 10
          // Default funding index is 10050 (last updated)
          // Second funding event index is 1000 (latest)
          // 10 * (10050 - 1000) = 90500
          payment: '90500',
          effectiveAt: createdDateTime.toISO(),
        },
        {
          payment: '200000',
          effectiveAt: createdDateTime.toISO(),
        },
        {
          payment: '200',
          effectiveAt: createdDateTime.toISO(),
        },
        {
          payment: '100',
          effectiveAt: createdDateTime.toISO(),
        },
      ],
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
