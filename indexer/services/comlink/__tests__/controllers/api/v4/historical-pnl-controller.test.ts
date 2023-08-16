import {
  dbHelpers,
  testConstants,
  testMocks,
  PnlTicksCreateObject,
  PnlTicksTable,
} from '@dydxprotocol-indexer/postgres';
import { PnlTicksResponseObject, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';

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

      const expectedPnlTickResponse: PnlTicksResponseObject = {
        ...testConstants.defaultPnlTick,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          testConstants.defaultPnlTick.createdAt,
        ),
      };

      const expectedPnlTick2Response: any = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          createdAt,
        ),
      };

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

      const expectedPnlTick2Response: PnlTicksResponseObject = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          createdAt,
        ),
      };

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

      const expectedPnlTickResponse: PnlTicksResponseObject = {
        ...testConstants.defaultPnlTick,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          testConstants.defaultPnlTick.createdAt,
        ),
      };

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
        path: '/v4/historical-pnl?address=invalid_address&subaccountNumber=100',
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: 'No subaccount found with address invalid_address and subaccountNumber 100',
          },
        ],
      });
    });
  });
});
