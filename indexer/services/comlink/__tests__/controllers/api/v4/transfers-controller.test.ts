import {
  dbHelpers,
  testConstants,
  testMocks,
  TransferCreateObject,
  TransferTable,
} from '@dydxprotocol-indexer/postgres';
import { RequestMethod, TransferResponseObject } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';

describe('transfers-controller#V4', () => {
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

    it('Get /transfers', async () => {
      await testMocks.seedData();
      const transfer2: TransferCreateObject = {
        senderSubaccountId: testConstants.defaultSubaccountId2,
        recipientSubaccountId: testConstants.defaultSubaccountId,
        assetId: testConstants.defaultAsset2.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId2,
        transactionHash: '', // TODO: Add a real transaction Hash
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };
      await Promise.all([
        TransferTable.create(testConstants.defaultTransfer),
        TransferTable.create(transfer2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expectedTransferResponse: TransferResponseObject = {
        id: testConstants.defaultTransferId,
        senderSubaccountId: testConstants.defaultTransfer.senderSubaccountId,
        recipientSubaccountId: testConstants.defaultTransfer.recipientSubaccountId,
        assetId: testConstants.defaultTransfer.assetId,
        size: testConstants.defaultTransfer.size,
        createdAt: testConstants.defaultTransfer.createdAt,
        createdAtHeight: testConstants.defaultTransfer.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
      };

      const expectedTransfer2Response: TransferResponseObject = {
        id: TransferTable.uuid(
          transfer2.senderSubaccountId,
          transfer2.recipientSubaccountId,
          transfer2.eventId,
          transfer2.assetId,
        ),
        senderSubaccountId: transfer2.senderSubaccountId,
        recipientSubaccountId: transfer2.recipientSubaccountId,
        assetId: transfer2.assetId,
        size: transfer2.size,
        createdAt: transfer2.createdAt,
        createdAtHeight: transfer2.createdAtHeight,
        symbol: testConstants.defaultAsset2.symbol,
      };

      expect(response.body.transfers).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedTransferResponse,
          }),
          expect.objectContaining({
            ...expectedTransfer2Response,
          }),
        ]),
      );
    });

    it('Get /transfers respects createdBeforeOrAt field', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const transfer2: TransferCreateObject = {
        senderSubaccountId: testConstants.defaultSubaccountId2,
        recipientSubaccountId: testConstants.defaultSubaccountId,
        assetId: testConstants.defaultAsset2.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId2,
        transactionHash: '', // TODO: Add a real transaction Hash
        createdAt,
        createdAtHeight: testConstants.createdHeight,
      };
      await Promise.all([
        TransferTable.create(testConstants.defaultTransfer),
        TransferTable.create(transfer2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}` +
          `&createdBeforeOrAt=${createdAt}`,
      });

      const expectedTransfer2Response: TransferResponseObject = {
        id: TransferTable.uuid(
          transfer2.senderSubaccountId,
          transfer2.recipientSubaccountId,
          transfer2.eventId,
          transfer2.assetId,
        ),
        senderSubaccountId: transfer2.senderSubaccountId,
        recipientSubaccountId: transfer2.recipientSubaccountId,
        assetId: transfer2.assetId,
        size: transfer2.size,
        createdAt,
        createdAtHeight: transfer2.createdAtHeight,
        symbol: testConstants.defaultAsset2.symbol,
      };

      expect(response.body.transfers).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedTransfer2Response,
          }),
        ]),
      );
    });

    it('Get /transfers respects createdBeforeOrAtHeight field', async () => {
      await testMocks.seedData();
      const createdAtHeight: string = '5';
      const transfer2: TransferCreateObject = {
        senderSubaccountId: testConstants.defaultSubaccountId2,
        recipientSubaccountId: testConstants.defaultSubaccountId,
        assetId: testConstants.defaultAsset2.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId2,
        transactionHash: '', // TODO: Add a real transaction Hash
        createdAt: testConstants.defaultTransfer.createdAt,
        createdAtHeight,
      };
      await Promise.all([
        TransferTable.create(testConstants.defaultTransfer),
        TransferTable.create(transfer2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}` +
          `&createdBeforeOrAtHeight=${createdAtHeight}`,
      });

      const expectedTransferResponse: TransferResponseObject = {
        id: testConstants.defaultTransferId,
        senderSubaccountId: testConstants.defaultTransfer.senderSubaccountId,
        recipientSubaccountId: testConstants.defaultTransfer.recipientSubaccountId,
        assetId: testConstants.defaultTransfer.assetId,
        size: testConstants.defaultTransfer.size,
        createdAt: testConstants.defaultTransfer.createdAt,
        createdAtHeight: testConstants.defaultTransfer.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
      };

      expect(response.body.transfers).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedTransferResponse,
          }),
        ]),
      );
    });

    it('Get /transfers returns empty when there are no transfers', async () => {
      await testMocks.seedData();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      expect(response.body.transfers).toHaveLength(0);
    });

    it('Get /transfers with non-existent address and subaccount number returns 404', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/transfers?address=invalid_address&subaccountNumber=100',
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
