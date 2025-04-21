import {
  dbHelpers,
  testConstants,
  testMocks,
  PnlTicksCreateObject,
  PnlTicksTable,
} from '@dydxprotocol-indexer/postgres';
import { PnlTicksResponseObject, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { pnlTickCreateObjectToResponseObject, sendRequest } from '../../../helpers/helpers';

describe('pnlTicks-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET', () => {
    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('Get /historical-pnl', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const blockHeight: string = '3';
      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historical-pnl?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expectedPnlTickResponse
      : PnlTicksResponseObject = pnlTickCreateObjectToResponseObject({
        ...testConstants.defaultPnlTick,
      });

      const expectedPnlTick2Response
      : PnlTicksResponseObject = pnlTickCreateObjectToResponseObject({
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      });

      expect(response.body.historicalPnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlTick2Response,
          }),
          expect.objectContaining({
            ...expectedPnlTickResponse,
          }),
        ]),
      );
    });

    it('Get /historical-pnl respects pagination', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const blockHeight: string = '1';
      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const responsePage1: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historical-pnl?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=1&limit=1`,
      });

      const responsePage2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historical-pnl?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=2&limit=1`,
      });

      const expectedPnlTickResponse
      : PnlTicksResponseObject = pnlTickCreateObjectToResponseObject({
        ...testConstants.defaultPnlTick,
      });

      const expectedPnlTick2Response
      : PnlTicksResponseObject = pnlTickCreateObjectToResponseObject({
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      });

      expect(responsePage1.body.pageSize).toStrictEqual(1);
      expect(responsePage1.body.offset).toStrictEqual(0);
      expect(responsePage1.body.totalResults).toStrictEqual(2);
      expect(responsePage1.body.historicalPnl).toHaveLength(1);
      expect(responsePage1.body.historicalPnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlTickResponse,
          }),
        ]),
      );

      expect(responsePage2.body.pageSize).toStrictEqual(1);
      expect(responsePage2.body.offset).toStrictEqual(1);
      expect(responsePage2.body.totalResults).toStrictEqual(2);
      expect(responsePage2.body.historicalPnl).toHaveLength(1);
      expect(responsePage2.body.historicalPnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlTick2Response,
          }),
        ]),
      );
    });

    it('Get /historical-pnl respects createdBeforeOrAt and createdBeforeOrAtHeight field', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const blockHeight: string = '1';
      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historical-pnl?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}` +
          `&createdBeforeOrAt=${createdAt}` +
          `&createdBeforeOrAtHeight=${blockHeight}`,
      });

      const expectedPnlTick2Response
      : PnlTicksResponseObject = pnlTickCreateObjectToResponseObject({
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      });

      expect(response.body.historicalPnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlTick2Response,
          }),
        ]),
      );
    });

    it('Get /historical-pnl respects createdOnOrAfter and createdOnOrAfterHeight field', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const blockHeight: string = '1';
      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historical-pnl?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}` +
          `&createdOnOrAfter=${testConstants.defaultPnlTick.createdAt}` +
          `&createdOnOrAfterHeight=${testConstants.defaultPnlTick.blockHeight}`,
      });

      const expectedPnlTickResponse
      : PnlTicksResponseObject = pnlTickCreateObjectToResponseObject({
        ...testConstants.defaultPnlTick,
      });

      expect(response.body.historicalPnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlTickResponse,
          }),
        ]),
      );
    });

    it('Get /historical-pnl respects all created fields', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const blockHeight: string = '1';
      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const createdOnOrAfter: string = '2001-05-25T00:00:00.000Z';
      const createdBeforeOrAt: string = '2005-05-25T00:00:00.000Z';

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historical-pnl?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}` +
          `&createdOnOrAfter=${createdOnOrAfter}` +
          `&createdBeforeOrAt=${createdBeforeOrAt}`,
      });
      expect(response.body.historicalPnl).toHaveLength(0);
    });

    it('Get /historical-pnl returns empty when there are no historical pnl ticks', async () => {
      await testMocks.seedData();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historical-pnl?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      expect(response.body.historicalPnl).toHaveLength(0);
    });

    it('Get /historical-pnl with non-existent address and subaccount number returns 404', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/historical-pnl?address=invalidaddress&subaccountNumber=100',
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: 'No subaccount found with address invalidaddress and subaccountNumber 100',
          },
        ],
      });
    });

    it('Get /historical-pnl/parentSubaccountNumber', async () => {
      await testMocks.seedData();
      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.isolatedSubaccountId,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/historical-pnl/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${parentSubaccountNumber}`,
      });

      const expectedPnlTickResponse: any = {
        // id and subaccountId don't matter
        equity: (parseFloat(testConstants.defaultPnlTick.equity) +
            parseFloat(pnlTick2.equity)).toString(),
        totalPnl: (parseFloat(testConstants.defaultPnlTick.totalPnl) +
            parseFloat(pnlTick2.totalPnl)).toString(),
        netTransfers: (parseFloat(testConstants.defaultPnlTick.netTransfers) +
            parseFloat(pnlTick2.netTransfers)).toString(),
        createdAt: testConstants.defaultPnlTick.createdAt,
        blockHeight: testConstants.defaultPnlTick.blockHeight,
        blockTime: testConstants.defaultPnlTick.blockTime,
      };

      expect(response.body.historicalPnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlTickResponse,
          }),
        ]),
      );
    });
  });

  it('Get /historical-pnl/parentSubaccountNumber with invalid subaccount number returns error', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/historical-pnl/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
          '&parentSubaccountNumber=128',
      expectedStatus: 400,
    });

    expect(response.body).toEqual({
      errors: [
        {
          location: 'query',
          msg: 'parentSubaccountNumber must be a non-negative integer less than 128',
          param: 'parentSubaccountNumber',
          value: '128',
        },
      ],
    });
  });
});
