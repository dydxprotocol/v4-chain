import {
  BlockCreateObject,
  BlockTable,
  dbHelpers,
  IsoString,
  SubaccountTable,
  SubaccountUsernamesTable,
  TendermintEventCreateObject,
  TendermintEventTable,
  testConstants,
  testMocks,
  TransferCreateObject,
  TransferTable,
  TransferType,
  WalletTable,
} from '@dydxprotocol-indexer/postgres';
import {
  ParentSubaccountTransferResponseObject,
  RequestMethod,
  TransferBetweenRequest,
  TransferBetweenResponse,
  TransferResponseObject,
} from '../../../../src/types';
import request from 'supertest';
import { getQueryString, sendRequest } from '../../../helpers/helpers';
import { defaultWalletAddress } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import Big from 'big.js';

const defaultWallet = {
  ...testConstants.defaultWallet,
  address: defaultWalletAddress, // defaultWalletAddress != testConstants.defaultWallet.address

};

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
      // use wallet2 to not create duplicate
      await WalletTable.create(testConstants.defaultWallet2);
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

    it('Get /transfers returns transfers/deposits/withdrawals with pagination', async () => {
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
      await WalletTable.create(defaultWallet);
      await Promise.all([
        TransferTable.create(testConstants.defaultTransfer),
        TransferTable.create(transfer2),
        TransferTable.create(testConstants.defaultWithdrawal),
        TransferTable.create(testConstants.defaultDeposit),
      ]);

      const responsePage1: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=1&limit=2`,
      });
      const responsePage2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=2&limit=2`,
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

      expect(responsePage1.body.pageSize).toStrictEqual(2);
      expect(responsePage1.body.offset).toStrictEqual(0);
      expect(responsePage1.body.totalResults).toStrictEqual(4);
      expect(responsePage1.body.transfers).toHaveLength(2);
      expect(responsePage1.body.transfers).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedTransferResponse,
          }),
          expect.objectContaining({
            ...expectedTransfer2Response,
          }),
        ]),
      );

      expect(responsePage2.body.pageSize).toStrictEqual(2);
      expect(responsePage2.body.offset).toStrictEqual(2);
      expect(responsePage2.body.totalResults).toStrictEqual(4);
      expect(responsePage2.body.transfers).toHaveLength(2);
      expect(responsePage2.body.transfers).toEqual(
        expect.arrayContaining([
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
        path: '/v4/transfers?address=invalidaddress&subaccountNumber=100',
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

    it('Get /transfers/parentSubaccountNumber returns more than 10 transfers after filtering', async () => {
      await testMocks.seedData();
      await WalletTable.create(defaultWallet);
      await SubaccountTable.create(testConstants.defaultSubaccount2Num0);

      // Create 50 blocks for all transfers
      const blocks: BlockCreateObject[] = [];
      for (let i = 0; i < 50; i++) {
        blocks.push({
          blockHeight: (3 + i).toString(),
          time: testConstants.createdDateTime.plus({ minutes: i }).toISO(),
        });
      }
      await Promise.all(blocks.map((b) => BlockTable.create(b)));

      // Create 50 TendermintEvents
      const events: TendermintEventCreateObject[] = [];
      for (let i = 0; i < 50; i++) {
        events.push({
          blockHeight: (3 + i).toString(),
          transactionIndex: 0,
          eventIndex: 0,
        });
      }
      await Promise.all(events.map((e) => TendermintEventTable.create(e)));

      // Create 35 same-parent transfers that should be filtered out
      // These are transfers between defaultSubaccount (0) and isolatedSubaccount (128)
      // Both have parent 0 (0 % 128 = 0, 128 % 128 = 0)
      const sameParentTransfers: TransferCreateObject[] = [];
      for (let i = 0; i < 35; i++) {
        const eventId = TendermintEventTable.createEventId(
          events[i].blockHeight,
          events[i].transactionIndex,
          events[i].eventIndex,
        );

        sameParentTransfers.push({
          senderSubaccountId: testConstants.defaultSubaccountId,
          recipientSubaccountId: testConstants.isolatedSubaccountId,
          assetId: testConstants.defaultAsset.id,
          size: `${i + 1}`,
          eventId,
          transactionHash: `same_parent_${i}`,
          createdAt: testConstants.createdDateTime.plus({ minutes: i }).toISO(),
          createdAtHeight: (parseInt(testConstants.createdHeight, 10) + i).toString(),
        });
      }
      await Promise.all(sameParentTransfers.map((t) => TransferTable.create(t)));

      // Create 15 cross-parent transfers (parent 0 -> parent 1)
      // These should ALL be returned
      const crossParentTransfers: TransferCreateObject[] = [];
      for (let i = 0; i < 15; i++) {
        const eventId = TendermintEventTable.createEventId(
          events[35 + i].blockHeight,
          events[35 + i].transactionIndex,
          events[35 + i].eventIndex,
        );

        crossParentTransfers.push({
          senderSubaccountId: testConstants.defaultSubaccountId,
          recipientSubaccountId: testConstants.defaultSubaccountId2,
          assetId: testConstants.defaultAsset.id,
          size: `${i + 100}`,
          eventId,
          transactionHash: `cross_parent_${i}`,
          createdAt: testConstants.createdDateTime.plus({ minutes: 35 + i }).toISO(),
          createdAtHeight: (parseInt(testConstants.createdHeight, 10) + 35 + i).toString(),
        });
      }
      await Promise.all(crossParentTransfers.map((t) => TransferTable.create(t)));

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
      `&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      // Should return all 15 cross-parent transfers
      expect(response.body.transfers.length).toEqual(15);

      // Verify none of the same-parent transfers are included
      response.body.transfers.forEach((transfer: ParentSubaccountTransferResponseObject) => {
        expect(transfer.sender.parentSubaccountNumber)
          .not.toEqual(transfer.recipient.parentSubaccountNumber);
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
        transactionHash: '',
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };
      await WalletTable.create(defaultWallet);
      await SubaccountTable.create(testConstants.defaultSubaccount2Num0);
      await Promise.all([
        TransferTable.create(testConstants.defaultTransfer),
        TransferTable.create(transfer2),
        TransferTable.create(testConstants.defaultWithdrawal),
        TransferTable.create(testConstants.defaultDeposit),
        TransferTable.create(testConstants.defaultTransferWithAlternateAddress),
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

      const expectedTransferWithAlternateAddressResponse: ParentSubaccountTransferResponseObject = {
        id: testConstants.defaultTransferWithAlternateAddressId,
        sender: {
          address: testConstants.defaultAddress2,
          parentSubaccountNumber: testConstants.defaultSubaccount2Num0.subaccountNumber,
        },
        recipient: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        },
        size: testConstants.defaultTransferWithAlternateAddress.size,
        createdAt: testConstants.defaultTransferWithAlternateAddress.createdAt,
        createdAtHeight: testConstants.defaultTransferWithAlternateAddress.createdAtHeight,
        symbol: testConstants.defaultAsset.symbol,
        type: TransferType.TRANSFER_IN,
        transactionHash: testConstants.defaultTransferWithAlternateAddress.transactionHash,
      };

      expect(response.body.transfers).toEqual(
        expect.arrayContaining([
          expect.objectContaining(expectedTransferResponse),
          expect.objectContaining(expectedTransfer2Response),
          expect.objectContaining(expectedWithdrawalResponse),
          expect.objectContaining(expectedDepositResponse),
          expect.objectContaining(expectedTransferWithAlternateAddressResponse),
        ]),
      );
    });

    it('Get /transfers/parentSubaccountNumber returns transfers with pagination', async () => {
      await testMocks.seedData();
      const transfer2: TransferCreateObject = {
        senderSubaccountId: testConstants.defaultSubaccountId2,
        recipientSubaccountId: testConstants.defaultSubaccountId,
        assetId: testConstants.defaultAsset2.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId2,
        transactionHash: '',
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };
      await WalletTable.create(defaultWallet);
      await Promise.all([
        TransferTable.create(testConstants.defaultTransfer),
        TransferTable.create(transfer2),
        TransferTable.create(testConstants.defaultWithdrawal),
        TransferTable.create(testConstants.defaultDeposit),
      ]);

      const responsePage1: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
        `&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=1&limit=2`,
      });

      const responsePage2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
        `&parentSubaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}&page=2&limit=2`,
      });

      expect(responsePage1.body.pageSize).toStrictEqual(2);
      expect(responsePage1.body.offset).toStrictEqual(0);
      expect(responsePage1.body.totalResults).toStrictEqual(4);
      expect(responsePage1.body.transfers).toHaveLength(2);

      expect(responsePage2.body.pageSize).toStrictEqual(2);
      expect(responsePage2.body.offset).toStrictEqual(2);
      expect(responsePage2.body.totalResults).toStrictEqual(4);
      expect(responsePage2.body.transfers).toHaveLength(2);
    });

    it('Get /transfers/parentSubaccountNumber excludes transfers for parent <> child subaccounts', async () => {
      await testMocks.seedData();
      const transfer2: TransferCreateObject = {
        senderSubaccountId: testConstants.defaultSubaccountId,
        recipientSubaccountId: testConstants.isolatedSubaccountId,
        assetId: testConstants.defaultAsset.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId2,
        transactionHash: '',
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };
      const transfer3: TransferCreateObject = {
        senderSubaccountId: testConstants.isolatedSubaccountId2,
        recipientSubaccountId: testConstants.defaultSubaccountId,
        assetId: testConstants.defaultAsset.id,
        size: '5',
        eventId: testConstants.defaultTendermintEventId3,
        transactionHash: '',
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };
      await WalletTable.create(defaultWallet);
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
          expect.objectContaining(expectedTransferResponse),
        ]),
      );
    });

    it('Get /transfers/parentSubaccountNumber includes deposits, withdrawals, and cross-parent transfers', async () => {
      await testMocks.seedData();
      await WalletTable.create(defaultWallet);

      // Create 4 blocks for the transfers
      const blocks: BlockCreateObject[] = [];
      for (let i = 0; i < 4; i++) {
        blocks.push({
          blockHeight: (3 + i).toString(),
          time: testConstants.createdDateTime.plus({ minutes: i }).toISO(),
        });
      }
      await Promise.all(blocks.map((b) => BlockTable.create(b)));

      // Create 4 TendermintEvents
      const events: TendermintEventCreateObject[] = [];
      for (let i = 0; i < 4; i++) {
        events.push({
          blockHeight: (3 + i).toString(),
          transactionIndex: 0,
          eventIndex: 0,
        });
      }
      await Promise.all(events.map((e) => TendermintEventTable.create(e)));

      const eventIds = events.map((e) => TendermintEventTable.createEventId(
        e.blockHeight, e.transactionIndex, e.eventIndex),
      );

      // Transfer from different parent (parent 1) to parent 0 child subaccount
      const transferFromDifferentParent: TransferCreateObject = {
        senderSubaccountId: testConstants.defaultSubaccountId2,
        recipientSubaccountId: testConstants.isolatedSubaccountId,
        assetId: testConstants.defaultAsset.id,
        size: '5',
        eventId: eventIds[0],
        transactionHash: 'transfer_from_different_parent',
        createdAt: testConstants.createdDateTime.toISO(),
        createdAtHeight: testConstants.createdHeight,
      };

      // Transfer from parent 0 child subaccount to different parent (parent 1)
      const transferToDifferentParent: TransferCreateObject = {
        senderSubaccountId: testConstants.isolatedSubaccountId2,
        recipientSubaccountId: testConstants.defaultSubaccountId2,
        assetId: testConstants.defaultAsset.id,
        size: '7',
        eventId: eventIds[1],
        transactionHash: 'transfer_to_different_parent',
        createdAt: testConstants.createdDateTime.plus({ minutes: 1 }).toISO(),
        createdAtHeight: testConstants.createdHeight,
      };

      // Deposit from wallet to parent 0 child subaccount
      const deposit: TransferCreateObject = {
        senderWalletAddress: testConstants.defaultWalletAddress,
        recipientSubaccountId: testConstants.isolatedSubaccountId,
        assetId: testConstants.defaultAsset.id,
        size: '10',
        eventId: eventIds[2],
        transactionHash: 'deposit',
        createdAt: testConstants.createdDateTime.plus({ minutes: 2 }).toISO(),
        createdAtHeight: testConstants.createdHeight,
      };

      // Withdrawal from parent 0 child subaccount to wallet
      const withdrawal: TransferCreateObject = {
        senderSubaccountId: testConstants.isolatedSubaccountId,
        recipientWalletAddress: testConstants.defaultWalletAddress,
        assetId: testConstants.defaultAsset.id,
        size: '12',
        eventId: eventIds[3],
        transactionHash: 'withdrawal',
        createdAt: testConstants.createdDateTime.plus({ minutes: 3 }).toISO(),
        createdAtHeight: testConstants.createdHeight,
      };

      await Promise.all([
        TransferTable.create(transferFromDifferentParent),
        TransferTable.create(transferToDifferentParent),
        TransferTable.create(deposit),
        TransferTable.create(withdrawal),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
      '&parentSubaccountNumber=0',
      });

      expect(response.body.transfers.length).toEqual(4);

      // Verify each transfer type is present
      const transfers = response.body.transfers;

      // Check transfer from different parent
      expect(transfers).toContainEqual(expect.objectContaining({
        sender: expect.objectContaining({
          parentSubaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        }),
        recipient: expect.objectContaining({
          parentSubaccountNumber: 0,
        }),
        size: '5',
        type: TransferType.TRANSFER_IN,
      }));

      // Check transfer to different parent
      expect(transfers).toContainEqual(expect.objectContaining({
        sender: expect.objectContaining({
          parentSubaccountNumber: 0,
        }),
        recipient: expect.objectContaining({
          parentSubaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
        }),
        size: '7',
        type: TransferType.TRANSFER_OUT,
      }));

      // Check deposit
      expect(transfers).toContainEqual(expect.objectContaining({
        sender: {
          address: testConstants.defaultWalletAddress,
        },
        recipient: expect.objectContaining({
          parentSubaccountNumber: 0,
        }),
        size: '10',
        type: TransferType.DEPOSIT,
      }));

      // Check withdrawal
      expect(transfers).toContainEqual(expect.objectContaining({
        sender: expect.objectContaining({
          parentSubaccountNumber: 0,
        }),
        recipient: {
          address: testConstants.defaultWalletAddress,
        },
        size: '12',
        type: TransferType.WITHDRAWAL,
      }));
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

    it('Get /transfers/parentSubaccountNumber with non-existent address returns 404', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/transfers/parentSubaccountNumber?address=invalidaddress&parentSubaccountNumber=0',
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: 'No subaccount found with address invalidaddress and parentSubaccountNumber 0',
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

  describe('GET /transfers/between', () => {
    beforeEach(async () => {
      await testMocks.seedData();
    });

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    const firstTransfer: TransferCreateObject = testConstants.defaultTransfer;
    const secondTransfer: TransferCreateObject = {
      ...testConstants.defaultTransfer,
      size: '5',
      createdAt: testConstants.createdDateTime.plus({ minutes: 1 }).toISO(),
      createdAtHeight: testConstants.createdHeight + 1,
      eventId: testConstants.defaultTendermintEventId2,
    };

    const firstTransferResponse: TransferResponseObject = {
      id: TransferTable.uuid(
        firstTransfer.eventId,
        firstTransfer.assetId,
        firstTransfer.senderSubaccountId,
        firstTransfer.recipientSubaccountId,
      ),
      sender: {
        address: testConstants.defaultSubaccount.address,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      },
      recipient: {
        address: testConstants.defaultSubaccount2.address,
        subaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
      },
      size: firstTransfer.size,
      createdAt: firstTransfer.createdAt,
      createdAtHeight: firstTransfer.createdAtHeight,
      symbol: 'USDC',
      type: TransferType.TRANSFER_OUT,
      transactionHash: firstTransfer.transactionHash,
    };
    const secondTransferResponse: TransferResponseObject = {
      ...firstTransferResponse,
      id: TransferTable.uuid(
        secondTransfer.eventId,
        secondTransfer.assetId,
        secondTransfer.senderSubaccountId,
        secondTransfer.recipientSubaccountId,
      ),
      size: secondTransfer.size,
      createdAt: secondTransfer.createdAt,
      createdAtHeight: secondTransfer.createdAtHeight,
    };

    async function getTransferBetweenResponse(
      createdBeforeOrAtHeight?: number,
      createdBeforeOrAt?: IsoString,
    ): Promise<TransferBetweenResponse> {
      const queryParams: TransferBetweenRequest = {
        sourceAddress: testConstants.defaultSubaccount.address,
        sourceSubaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        recipientAddress: testConstants.defaultSubaccount2.address,
        recipientSubaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
      };

      if (createdBeforeOrAtHeight) {
        queryParams.createdBeforeOrAtHeight = createdBeforeOrAtHeight;
      }

      if (createdBeforeOrAt) {
        queryParams.createdBeforeOrAt = createdBeforeOrAt;
      }

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/transfers/between?${getQueryString(queryParams as any)}`,
      });

      return response.body;
    }

    it('Returns successfully when there are no transfers between wallets', async () => {
      await dbHelpers.clearData();

      const transferBetweenResponse: TransferBetweenResponse = await getTransferBetweenResponse();
      expect(transferBetweenResponse.transfersSubset).toHaveLength(0);
      expect(transferBetweenResponse.totalNetTransfers).toEqual('0');
    });

    it('Returns successfully when source subaccount does not exist', async () => {
      await SubaccountUsernamesTable.deleteBySubaccountId(testConstants.defaultSubaccountId);
      await SubaccountTable.deleteById(testConstants.defaultSubaccountId);

      const transferBetweenResponse: TransferBetweenResponse = await getTransferBetweenResponse();
      expect(transferBetweenResponse.transfersSubset).toHaveLength(0);
      expect(transferBetweenResponse.totalNetTransfers).toEqual('0');
    });

    it('Returns successfully when recipient subaccount does not exist', async () => {
      await SubaccountUsernamesTable.deleteBySubaccountId(testConstants.defaultSubaccountId2);
      await SubaccountTable.deleteById(testConstants.defaultSubaccountId2);

      const transferBetweenResponse: TransferBetweenResponse = await getTransferBetweenResponse();
      expect(transferBetweenResponse.transfersSubset).toHaveLength(0);
      expect(transferBetweenResponse.totalNetTransfers).toEqual('0');

    });

    it('Returns successfully with transfers and net transfers', async () => {
      await Promise.all([
        TransferTable.create(firstTransfer),
        TransferTable.create(secondTransfer),
      ]);

      const transferBetweenResponse: TransferBetweenResponse = await getTransferBetweenResponse();
      expect(transferBetweenResponse.transfersSubset).toHaveLength(2);
      expect(transferBetweenResponse.transfersSubset).toEqual([
        secondTransferResponse,
        firstTransferResponse,
      ]);
      expect(transferBetweenResponse.totalNetTransfers).toEqual(
        Big(firstTransfer.size).plus(secondTransfer.size).toFixed(),
      );
    });

    it('Successfully filters by createdBeforeOrAtHeight', async () => {
      await Promise.all([
        TransferTable.create(firstTransfer),
        TransferTable.create(secondTransfer),
      ]);

      const transferBetweenResponse: TransferBetweenResponse = await getTransferBetweenResponse(
        +firstTransfer.createdAtHeight,
      );
      expect(transferBetweenResponse.transfersSubset).toHaveLength(1);
      expect(transferBetweenResponse.transfersSubset).toEqual([
        firstTransferResponse,
      ]);
      expect(transferBetweenResponse.totalNetTransfers).toEqual(
        Big(firstTransfer.size).plus(secondTransfer.size).toFixed(),
      );
    });

    it('Successfully filters by createdBeforeOrAt', async () => {
      await Promise.all([
        TransferTable.create(firstTransfer),
        TransferTable.create(secondTransfer),
      ]);

      const transferBetweenResponse: TransferBetweenResponse = await getTransferBetweenResponse(
        undefined,
        firstTransfer.createdAt,
      );
      expect(transferBetweenResponse.transfersSubset).toHaveLength(1);
      expect(transferBetweenResponse.transfersSubset).toEqual([
        firstTransferResponse,
      ]);
      expect(transferBetweenResponse.totalNetTransfers).toEqual(
        Big(firstTransfer.size).plus(secondTransfer.size).toFixed(),
      );
    });
  });
});
