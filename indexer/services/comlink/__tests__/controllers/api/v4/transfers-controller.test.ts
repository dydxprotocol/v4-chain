import {
  dbHelpers,
  testConstants,
  testMocks,
  TransferCreateObject,
  TransferTable,
  TransferType,
} from '@dydxprotocol-indexer/postgres';
import { ParentSubaccountTransferResponseObject, RequestMethod, TransferResponseObject } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';
import {
  createdDateTime, createdHeight,
  defaultAsset,
  defaultTendermintEventId4,
  defaultWalletAddress,
  isolatedSubaccountId,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

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

    it('Get /transfers returns transfers/deposits/withdrawals', async () => {
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
        TransferTable.create(testConstants.defaultWithdrawal),
        TransferTable.create(testConstants.defaultDeposit),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expectedTransferResponse: TransferResponseObject = {
        id: testConstants.defaultTransferId,
        sender: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        },
        size: testConstants.defaultTransfer.size,
        createdAt: testConstants.defaultTransfer.createdAt,
        createdAtHeight: testConstants.defaultTransfer.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.TRANSFER_OUT,
        transactionHash: testConstants.defaultTransfer.transactionHash,
      };

      const expectedTransfer2Response: TransferResponseObject = {
        id: TransferTable.uuid(
          transfer2.eventId,
          transfer2.assetId,
          transfer2.senderSubaccountId,
          transfer2.recipientSubaccountId,
          transfer2.senderWalletAddress,
          transfer2.recipientWalletAddress,
        ),
        sender: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        size: transfer2.size,
        createdAt: transfer2.createdAt,
        createdAtHeight: transfer2.createdAtHeight,
        symbol: testConstants.defaultAsset2.symbol,
        type: TransferType.TRANSFER_IN,
        transactionHash: transfer2.transactionHash,
      };

      const expectedDepositResponse: TransferResponseObject = {
        id: testConstants.defaultDepositId,
        sender: {
          address: testConstants.defaultWalletAddress,
        },
        recipient: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        size: testConstants.defaultDeposit.size,
        createdAt: testConstants.defaultDeposit.createdAt,
        createdAtHeight: testConstants.defaultDeposit.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.DEPOSIT,
        transactionHash: testConstants.defaultDeposit.transactionHash,
      };

      const expectedWithdrawalResponse: TransferResponseObject = {
        id: testConstants.defaultWithdrawalId,
        sender: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultWalletAddress,
        },
        size: testConstants.defaultWithdrawal.size,
        createdAt: testConstants.defaultWithdrawal.createdAt,
        createdAtHeight: testConstants.defaultWithdrawal.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.WITHDRAWAL,
        transactionHash: testConstants.defaultWithdrawal.transactionHash,
      };

      expect(response.body.transfers).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedTransferResponse,
          }),
          expect.objectContaining({
            ...expectedTransfer2Response,
          }),
          expect.objectContaining({
            ...expectedWithdrawalResponse,
          }),
          expect.objectContaining({
            ...expectedDepositResponse,
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
          transfer2.eventId,
          transfer2.assetId,
          transfer2.senderSubaccountId,
          transfer2.recipientSubaccountId,
          transfer2.senderWalletAddress,
          transfer2.recipientWalletAddress,
        ),
        sender: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        size: transfer2.size,
        createdAt,
        createdAtHeight: transfer2.createdAtHeight,
        symbol: testConstants.defaultAsset2.symbol,
        type: TransferType.TRANSFER_IN,
        transactionHash: transfer2.transactionHash,
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
        sender: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        },
        size: testConstants.defaultTransfer.size,
        createdAt: testConstants.defaultTransfer.createdAt,
        createdAtHeight: testConstants.defaultTransfer.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.TRANSFER_OUT,
        transactionHash: testConstants.defaultTransfer.transactionHash,
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
    it('Get /transfers/parentSubaccountNumber returns transfers/deposits/withdrawals', async () => {
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
        TransferTable.create(testConstants.defaultWithdrawal),
        TransferTable.create(testConstants.defaultDeposit),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expectedTransferResponse: ParentSubaccountTransferResponseObject = {
        id: testConstants.defaultTransferId,
        sender: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        },
        size: testConstants.defaultTransfer.size,
        createdAt: testConstants.defaultTransfer.createdAt,
        createdAtHeight: testConstants.defaultTransfer.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.TRANSFER_OUT,
        transactionHash: testConstants.defaultTransfer.transactionHash,
      };

      const expectedTransfer2Response: ParentSubaccountTransferResponseObject = {
        id: TransferTable.uuid(
          transfer2.eventId,
          transfer2.assetId,
          transfer2.senderSubaccountId,
          transfer2.recipientSubaccountId,
          transfer2.senderWalletAddress,
          transfer2.recipientWalletAddress,
        ),
        sender: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        size: transfer2.size,
        createdAt: transfer2.createdAt,
        createdAtHeight: transfer2.createdAtHeight,
        symbol: testConstants.defaultAsset2.symbol,
        type: TransferType.TRANSFER_IN,
        transactionHash: transfer2.transactionHash,
      };

      const expectedDepositResponse: ParentSubaccountTransferResponseObject = {
        id: testConstants.defaultDepositId,
        sender: {
          address: testConstants.defaultWalletAddress,
        },
        recipient: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        size: testConstants.defaultDeposit.size,
        createdAt: testConstants.defaultDeposit.createdAt,
        createdAtHeight: testConstants.defaultDeposit.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.DEPOSIT,
        transactionHash: testConstants.defaultDeposit.transactionHash,
      };

      const expectedWithdrawalResponse: ParentSubaccountTransferResponseObject = {
        id: testConstants.defaultWithdrawalId,
        sender: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultWalletAddress,
        },
        size: testConstants.defaultWithdrawal.size,
        createdAt: testConstants.defaultWithdrawal.createdAt,
        createdAtHeight: testConstants.defaultWithdrawal.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.WITHDRAWAL,
        transactionHash: testConstants.defaultWithdrawal.transactionHash,
      };

      expect(response.body.transfers).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedTransferResponse,
          }),
          expect.objectContaining({
            ...expectedTransfer2Response,
          }),
          expect.objectContaining({
            ...expectedWithdrawalResponse,
          }),
          expect.objectContaining({
            ...expectedDepositResponse,
          }),
        ]),
      );
    });

    it('Get /transfers/parentSubaccountNumber excludes transfers for parent <> child subaccounts', async () => {
      await testMocks.seedData();
      const transfer2: TransferCreateObject = {
        senderSubaccountId: testConstants.defaultSubaccountId,
        recipientSubaccountId: testConstants.isolatedSubaccountId,
        assetId: testConstants.defaultAsset.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId2,
        transactionHash: '', // TODO: Add a real transaction Hash
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };
      const transfer3: TransferCreateObject = {
        senderSubaccountId: testConstants.isolatedSubaccountId2,
        recipientSubaccountId: testConstants.defaultSubaccountId,
        assetId: testConstants.defaultAsset.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId3,
        transactionHash: '', // TODO: Add a real transaction Hash
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };
      await Promise.all([
        TransferTable.create(testConstants.defaultTransfer),
        TransferTable.create(transfer2),
        TransferTable.create(transfer3),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expectedTransferResponse: ParentSubaccountTransferResponseObject = {
        id: testConstants.defaultTransferId,
        sender: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        },
        size: testConstants.defaultTransfer.size,
        createdAt: testConstants.defaultTransfer.createdAt,
        createdAtHeight: testConstants.defaultTransfer.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.TRANSFER_OUT,
        transactionHash: testConstants.defaultTransfer.transactionHash,
      };

      expect(response.body.transfers.length).toEqual(1);
      expect(response.body.transfers).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedTransferResponse,
          }),
        ]),
      );
    });

    it('Get /transfers/parentSubaccountNumber includes transfers for wallets/subaccounts(non parent) <> child subaccounts', async () => {
      await testMocks.seedData();
      const transferFromNonParent: TransferCreateObject = {
        senderSubaccountId: testConstants.defaultSubaccountId2,
        recipientSubaccountId: testConstants.isolatedSubaccountId,
        assetId: testConstants.defaultAsset.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId2,
        transactionHash: '', // TODO: Add a real transaction Hash
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };
      const transferToNonParent: TransferCreateObject = {
        senderSubaccountId: testConstants.isolatedSubaccountId2,
        recipientSubaccountId: testConstants.defaultSubaccountId2,
        assetId: testConstants.defaultAsset.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId3,
        transactionHash: '', // TODO: Add a real transaction Hash
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };
      const depositToChildSA: TransferCreateObject = {
        senderWalletAddress: defaultWalletAddress,
        recipientSubaccountId: isolatedSubaccountId,
        assetId: defaultAsset.id,
        size: '10',
        eventId: defaultTendermintEventId4,
        transactionHash: '', // TODO: Add a real transaction Hash
        createdAt: createdDateTime.toISO(),
        createdAtHeight: createdHeight,
      };
      const withdrawFromChildSA: TransferCreateObject = {
        senderSubaccountId: isolatedSubaccountId,
        recipientWalletAddress: defaultWalletAddress,
        assetId: defaultAsset.id,
        size: '10',
        eventId: defaultTendermintEventId4,
        transactionHash: '', // TODO: Add a real transaction Hash
        createdAt: createdDateTime.toISO(),
        createdAtHeight: createdHeight,
      };
      await Promise.all([
        TransferTable.create(transferFromNonParent),
        TransferTable.create(transferToNonParent),
        TransferTable.create(depositToChildSA),
        TransferTable.create(withdrawFromChildSA),
      ]);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${parentSubaccountNumber}`,
      });

      const expectedTransferResponse1: ParentSubaccountTransferResponseObject = {
        id: TransferTable.uuid(
          transferFromNonParent.eventId,
          transferFromNonParent.assetId,
          transferFromNonParent.senderSubaccountId,
          transferFromNonParent.recipientSubaccountId,
          transferFromNonParent.senderWalletAddress,
          transferFromNonParent.recipientWalletAddress,
        ),
        sender: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
        },
        size: transferFromNonParent.size,
        createdAt: transferFromNonParent.createdAt,
        createdAtHeight: transferFromNonParent.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.TRANSFER_IN,
        transactionHash: transferFromNonParent.transactionHash,
      };
      const expectedTransferResponse2: ParentSubaccountTransferResponseObject = {
        id: TransferTable.uuid(
          transferToNonParent.eventId,
          transferToNonParent.assetId,
          transferToNonParent.senderSubaccountId,
          transferToNonParent.recipientSubaccountId,
          transferToNonParent.senderWalletAddress,
          transferToNonParent.recipientWalletAddress,
        ),
        sender: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
        },
        recipient: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        },
        size: transferToNonParent.size,
        createdAt: transferToNonParent.createdAt,
        createdAtHeight: transferToNonParent.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.TRANSFER_OUT,
        transactionHash: transferToNonParent.transactionHash,
      };
      const expectedDepositResponse: ParentSubaccountTransferResponseObject = {
        id: TransferTable.uuid(
          depositToChildSA.eventId,
          depositToChildSA.assetId,
          depositToChildSA.senderSubaccountId,
          depositToChildSA.recipientSubaccountId,
          depositToChildSA.senderWalletAddress,
          depositToChildSA.recipientWalletAddress,
        ),
        sender: {
          address: testConstants.defaultWalletAddress,
        },
        recipient: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
        },
        size: depositToChildSA.size,
        createdAt: depositToChildSA.createdAt,
        createdAtHeight: depositToChildSA.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.DEPOSIT,
        transactionHash: depositToChildSA.transactionHash,
      };
      const expectedWithdrawalResponse: ParentSubaccountTransferResponseObject = {
        id: TransferTable.uuid(
          withdrawFromChildSA.eventId,
          withdrawFromChildSA.assetId,
          withdrawFromChildSA.senderSubaccountId,
          withdrawFromChildSA.recipientSubaccountId,
          withdrawFromChildSA.senderWalletAddress,
          withdrawFromChildSA.recipientWalletAddress,
        ),
        sender: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
        },
        recipient: {
          address: testConstants.defaultWalletAddress,
        },
        size: withdrawFromChildSA.size,
        createdAt: withdrawFromChildSA.createdAt,
        createdAtHeight: withdrawFromChildSA.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.WITHDRAWAL,
        transactionHash: withdrawFromChildSA.transactionHash,
      };

      expect(response.body.transfers.length).toEqual(4);
      expect(response.body.transfers).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedTransferResponse1,
          }),
          expect.objectContaining({
            ...expectedTransferResponse2,
          }),
          expect.objectContaining({
            ...expectedDepositResponse,
          }),
          expect.objectContaining({
            ...expectedWithdrawalResponse,
          }),
        ]),
      );
    });

    it('Get /transfers/parentSubaccountNumber returns empty when there are no transfers', async () => {
      await testMocks.seedData();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      expect(response.body.transfers).toHaveLength(0);
    });

    it('Get /transfers/parentSubaccountNumber with non-existent address and subaccount number returns 404', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/transfers/parentSubaccountNumber?address=invalid_address&parentSubaccountNumber=100',
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: 'No subaccount found with address invalid_address and parentSubaccountNumber 100',
          },
        ],
      });
    });

    it('Get /transfers/parentSubaccountNumber with invalid parentSubaccountNumber', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
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
});
