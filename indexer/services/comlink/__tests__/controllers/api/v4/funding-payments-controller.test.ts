import {
  BlockTable,
  DEFAULT_POSTGRES_OPTIONS,
  dbHelpers,
  FundingPaymentsCreateObject,
  FundingPaymentsTable,
  perpetualMarketRefresher,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { FundingPaymentResponseObject, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import config from '../../../../src/config';
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
    createdAtHeight: '2',
  };

  const fundingPayment3: FundingPaymentsCreateObject = {
    ...testConstants.defaultFundingPayment,
    createdAt: '2000-05-25T00:00:02.000Z',
    createdAtHeight: '3',
  };

  const fundingPayment4: FundingPaymentsCreateObject = {
    ...testConstants.defaultFundingPayment,
    createdAt: '2000-05-25T00:00:03.000Z',
    payment: '0',
    createdAtHeight: '4',
  };

  const fundingPayment5: FundingPaymentsCreateObject = {
    ...testConstants.defaultFundingPayment,
    createdAt: '2000-05-25T00:00:04.000Z',
    payment: '-1',
    createdAtHeight: '5',
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
    subaccountNumber: testConstants.defaultSubaccount.subaccountNumber.toString(),
    fundingIndex: testConstants.defaultFundingPayment.fundingIndex,
  };

  const expectedFundingPayment2: FundingPaymentResponseObject = {
    ...expectedFundingPayment1,
    perpetualId: testConstants.defaultPerpetualMarket2.id,
    ticker: testConstants.defaultPerpetualMarket2.ticker,
    createdAt: '2000-05-25T00:00:01.000Z',
    createdAtHeight: '2',
  };

  const expectedFundingPayment3: FundingPaymentResponseObject = {
    ...expectedFundingPayment1,
    createdAt: '2000-05-25T00:00:02.000Z',
    createdAtHeight: '3',
  };

  const expectedFundingPayment4: FundingPaymentResponseObject = {
    ...expectedFundingPayment1,
    createdAt: '2000-05-25T00:00:03.000Z',
    payment: '0',
    createdAtHeight: '4',
  };

  const expectedFundingPayment5: FundingPaymentResponseObject = {
    ...expectedFundingPayment1,
    createdAt: '2000-05-25T00:00:04.000Z',
    payment: '-1',
    createdAtHeight: '5',
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
      FundingPaymentsTable.create(fundingPayment4),
      FundingPaymentsTable.create(fundingPayment5),
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
      path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&showZeroPayments=true`,
    });

    expect(response.body.fundingPayments).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedFundingPayment1),
        expect.objectContaining(expectedFundingPayment2),
        expect.objectContaining(expectedFundingPayment3),
        expect.objectContaining(expectedFundingPayment4),
        expect.objectContaining(expectedFundingPayment5),
      ]),
    );
  });

  it('Get /fundingPayments with ticker filter', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&ticker=${testConstants.defaultPerpetualMarket.ticker}&showZeroPayments=true`,
    });

    expect(response.body.fundingPayments).toEqual(
      expect.arrayContaining([expect.objectContaining(expectedFundingPayment1)]),
    );
  });

  it('Get /fundingPayments with createdOnOrAfter filter', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&createdOnOrAfter=${fundingPayment2.createdAt}&showZeroPayments=true`,
    });

    expect(response.body.fundingPayments).toHaveLength(4);
    expect(response.body.fundingPayments).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedFundingPayment2),
        expect.objectContaining(expectedFundingPayment3),
        expect.objectContaining(expectedFundingPayment4),
        expect.objectContaining(expectedFundingPayment5),
      ]),
    );
  });

  it('Get /fundingPayments/parentSubaccount', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments/parentSubaccount?address=${testConstants.defaultAddress}&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&showZeroPayments=true`,
    });

    expect(response.body.fundingPayments).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedFundingPayment1),
        expect.objectContaining(expectedFundingPayment2),
        expect.objectContaining(expectedFundingPayment3),
        expect.objectContaining(expectedFundingPayment4),
        expect.objectContaining(expectedFundingPayment5),
      ]),
    );
  });

  it('Gets /fundingPayments/parentSubaccount with pagination', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments/parentSubaccount?address=${testConstants.defaultAddress}&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&showZeroPayments=true&limit=1&page=1`,
    });

    expect(response.body.fundingPayments.length).toEqual(1);
    const response2: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments/parentSubaccount?address=${testConstants.defaultAddress}&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&showZeroPayments=true&limit=1&page=2`,
    });

    expect(response2.body.fundingPayments.length).toEqual(1);

    // expects page 1 and page 2 to be different
    expect(response.body.fundingPayments[0]).not.toEqual(
      response2.body.fundingPayments[0],
    );
  });

  it('Get /fundingPayments routes paginated reads through DEFAULT_POSTGRES_OPTIONS (read-replica regression guard)', async () => {
    // Mutate the live DEFAULT_POSTGRES_OPTIONS object with a marker property. Since the
    // controller spreads this object at call time (`{ ...DEFAULT_POSTGRES_OPTIONS, orderBy }`),
    // the marker will only show up in the findAll call if that spread actually happens -
    // independent of whatever USE_READ_REPLICA currently evaluates to in this test environment.
    const testMarker: unique symbol = Symbol('read-replica-regression-marker');
    (DEFAULT_POSTGRES_OPTIONS as Record<string | symbol, unknown>)[testMarker] = true;
    const findAllSpy: jest.SpyInstance = jest.spyOn(FundingPaymentsTable, 'findAll');

    try {
      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&showZeroPayments=true&limit=1&page=1`,
      });

      expect(findAllSpy).toHaveBeenCalledWith(
        expect.anything(),
        expect.anything(),
        expect.objectContaining({ [testMarker]: true }),
      );
    } finally {
      delete (DEFAULT_POSTGRES_OPTIONS as Record<string | symbol, unknown>)[testMarker];
      findAllSpy.mockRestore();
    }
  });

  it('Get /fundingPayments/parentSubaccount routes paginated reads through DEFAULT_POSTGRES_OPTIONS (read-replica regression guard)', async () => {
    const testMarker: unique symbol = Symbol('read-replica-regression-marker');
    (DEFAULT_POSTGRES_OPTIONS as Record<string | symbol, unknown>)[testMarker] = true;
    const findAllSpy: jest.SpyInstance = jest.spyOn(FundingPaymentsTable, 'findAll');

    try {
      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/fundingPayments/parentSubaccount?address=${testConstants.defaultAddress}&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&showZeroPayments=true&limit=1&page=1`,
      });

      expect(findAllSpy).toHaveBeenCalledWith(
        expect.anything(),
        expect.anything(),
        expect.objectContaining({ [testMarker]: true }),
      );
    } finally {
      delete (DEFAULT_POSTGRES_OPTIONS as Record<string | symbol, unknown>)[testMarker];
      findAllSpy.mockRestore();
    }
  });

  it('Get /fundingPayments rejects a deep page/limit combination end-to-end through the real app (offset ceiling)', async () => {
    // Not the exact incident values (page=97-111 & limit=1000, offset ~96k-110k) - just deep
    // enough to clear MAX_PAGINATION_OFFSET with the same shape of request (a real HTTP call
    // through the full route: rate-limiter middleware, CheckPaginationSchema, and the controller
    // handler) rather than exercising the validator or the rate limiter in isolation.
    const limit: number = 1000;
    const page: number = Math.floor(config.MAX_PAGINATION_OFFSET / limit) + 2;

    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&showZeroPayments=true&limit=${limit}&page=${page}`,
      expectedStatus: 400,
    });

    expect(response.body.errors).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          location: 'query',
          param: 'page',
          msg: expect.stringContaining('exceeds MAX_PAGINATION_OFFSET'),
        }),
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

  it('Get /fundingPayments with paymentGreaterThanOrEqual filter omits 0 payment', async () => {
    // default no show zero payments
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/fundingPayments?address=${testConstants.defaultAddress}&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
    });

    expect(response.body.fundingPayments).toHaveLength(4);
    expect(response.body.fundingPayments).toEqual(
      expect.arrayContaining([
        expect.objectContaining(expectedFundingPayment1),
        expect.objectContaining(expectedFundingPayment2),
        expect.objectContaining(expectedFundingPayment3),
        expect.objectContaining(expectedFundingPayment5),
      ]),
    );
  });
});
