import {
  BlockTable,
  dbHelpers,
  FundingPaymentsCreateObject,
  FundingPaymentsTable,
  perpetualMarketRefresher,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  FundingPaymentResponseObject,
  RequestMethod,
} from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';

describe('funding-payments-controller#V4', () => {
  const fundingPayment1: FundingPaymentsCreateObject = {
    ...testConstants.defaultFundingPayment,
    createdAt: '2000-05-25T00:00:00.000Z',
  };

  const fundingPayment2: FundingPaymentsCreateObject = {
    ...testConstants.defaultFundingPayment,
    perpetualId: testConstants.defaultPerpetualMarket2.id,
    ticker: testConstants.defaultPerpetualMarket2.ticker,
    createdAt: '2000-05-25T00:00:01.000Z',
  };

  const fundingPayment3: FundingPaymentsCreateObject = {
    ...testConstants.defaultFundingPayment,
    createdAt: '2000-05-25T00:00:02.000Z',
  };

  const expectedFundingPayment1: FundingPaymentResponseObject = {
    createdAt: '2000-05-25T00:00:00.000Z',
    createdAtHeight: testConstants.defaultFundingPayment.createdAtHeight,
    perpetualId: testConstants.defaultFundingPayment.perpetualId,
    ticker: testConstants.defaultFundingPayment.ticker,
    oraclePrice: testConstants.defaultFundingPayment.oraclePrice,
    size: testConstants.defaultFundingPayment.size,
    side: testConstants.defaultFundingPayment.side,
    rate: testConstants.defaultFundingPayment.rate,
    payment: testConstants.defaultFundingPayment.payment,
    subaccountNumber:
      testConstants.defaultSubaccount.subaccountNumber.toString(),
  };

  const expectedFundingPayment2: FundingPaymentResponseObject = {
    ...expectedFundingPayment1,
    perpetualId: testConstants.defaultPerpetualMarket2.id,
    ticker: testConstants.defaultPerpetualMarket2.ticker,
    createdAt: '2000-05-25T00:00:01.000Z',
  };

  const expectedFundingPayment3: FundingPaymentResponseObject = {
    ...expectedFundingPayment1,
    createdAt: '2000-05-25T00:00:02.000Z',
  };

  beforeAll(async () => {
    await dbHelpers.clearData();
    await dbHelpers.migrate();
    await testMocks.seedData();
    await BlockTable.create({
      blockHeight: '7',
      time: '2000-05-25T00:00:00.000Z',
    });

    await Promise.all([
      FundingPaymentsTable.create(fundingPayment1),
      FundingPaymentsTable.create(fundingPayment2),
      FundingPaymentsTable.create(fundingPayment3),
    ]);
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterAll(async () => {
    await dbHelpers.clearData();
    await dbHelpers.teardown();
  });

  it('Get /fundingPayments', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
    });

    expect(response.body.fundingPayments).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedFundingPayment1),
        expect.objectContaining(expectedFundingPayment2),
      ]),
    );
  });

  it('Get /fundingPayments with ticker filter', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&ticker=${testConstants.defaultPerpetualMarket.ticker}`,
    });

    expect(response.body.fundingPayments).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedFundingPayment1),
      ]),
    );
  });

  it('Get /fundingPayments with afterOrAt filter', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&afterOrAt=${testConstants.defaultFundingPayment.createdAt}`,
    });

    expect(response.body.fundingPayments).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedFundingPayment1),
        expect.objectContaining(expectedFundingPayment2),
      ]),
    );
  });

  it('Get /fundingPayments/parentSubaccount', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments/parentSubaccount?address=${testConstants.defaultAddress}&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
    });

    expect(response.body.fundingPayments).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedFundingPayment1),
        expect.objectContaining(expectedFundingPayment2),
        expect.objectContaining(expectedFundingPayment3),
      ]),
    );
  });

  it('Returns 400 with invalid address', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments?address=inv@lid&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      expectedStatus: 400,
    });

    expect(response.body.errors[0]).toEqual(
      expect.objectContaining({
        location: 'query',
        msg: 'address must be a valid dydx address',
        param: 'address',
        value: 'inv@lid',
      }),
    );
  });

  it('Returns 400 with invalid subaccount number', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=invalid`,
      expectedStatus: 400,
    });

    expect(response.body.errors[0]).toEqual(
      expect.objectContaining({
        msg: 'subaccountNumber must be a non-negative integer less than 128001',
        location: 'query',
        param: 'subaccountNumber',
        value: 'invalid',
      }),
    );
  });
});
