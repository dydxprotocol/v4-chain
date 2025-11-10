import {
  BlockCreateObject,
  Ordering,
  SubaccountAssetNetTransferMap,
  TendermintEventCreateObject,
  TransferColumns,
  TransferCreateObject,
  TransferFromDatabase,
} from '../../src/types';
import * as TransferTable from '../../src/stores/transfer-table';
import { AssetTransferMap } from '../../src/stores/transfer-table';
import * as WalletTable from '../../src/stores/wallet-table';
import * as SubaccountTable from '../../src/stores/subaccount-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import {
  createdDateTime,
  createdHeight,
  defaultAsset,
  defaultAsset2,
  defaultDeposit,
  defaultSubaccount3,
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultSubaccountId3,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
  defaultTransfer,
  defaultTransfer2,
  defaultTransfer3,
  defaultWalletAddress,
  defaultWithdrawal,
  isolatedSubaccountId,
} from '../helpers/constants';
import Big from 'big.js';
import { CheckViolationError } from 'objection';
import { DateTime } from 'luxon';
import { BlockTable, TendermintEventTable, USDC_ASSET_ID } from '../../src';

describe('Transfer store', () => {
  beforeEach(async () => {
    await seedData();
  });

  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a Transfer', async () => {
    await TransferTable.create(defaultTransfer);
  });

  it('Successfully creates multiple transfers', async () => {
    const transfer2: TransferCreateObject = {
      senderSubaccountId: defaultSubaccountId2,
      recipientSubaccountId: defaultSubaccountId,
      assetId: defaultAsset2.id,
      size: '5',
      eventId: defaultTendermintEventId2,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: createdDateTime.toISO(),
      createdAtHeight: createdHeight,
    };
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create(transfer2),
    ]);

    const transfers: TransferFromDatabase[] = await TransferTable.findAll({}, [], {
      orderBy: [[TransferColumns.id, Ordering.ASC]],
    });

    expect(transfers.length).toEqual(2);
    expect(transfers[0]).toEqual(expect.objectContaining(defaultTransfer));
    expect(transfers[1]).toEqual(expect.objectContaining(transfer2));
  });

  it('Successfully finds all Transfers', async () => {
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create({
        ...defaultTransfer,
        eventId: defaultTendermintEventId2,
      }),
    ]);

    const transfers: TransferFromDatabase[] = await TransferTable.findAll({}, [], {
      orderBy: [[TransferColumns.id, Ordering.ASC]],
    });

    expect(transfers.length).toEqual(2);
    expect(transfers[0]).toEqual(expect.objectContaining(defaultTransfer));
    expect(transfers[1]).toEqual(expect.objectContaining({
      ...defaultTransfer,
      eventId: defaultTendermintEventId2,
    }));
  });

  it('Successfully finds all transfers to and from subaccount', async () => {
    const transfer2: TransferCreateObject = {
      senderSubaccountId: defaultSubaccountId2,
      recipientSubaccountId: defaultSubaccountId,
      assetId: defaultAsset2.id,
      size: '5',
      eventId: defaultTendermintEventId2,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: createdDateTime.toISO(),
      createdAtHeight: createdHeight,
    };
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create(transfer2),
    ]);

    const { results: transfers } = await TransferTable.findAllToOrFromSubaccountId(
      { subaccountId: [defaultSubaccountId] },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(transfers.length).toEqual(2);
    expect(transfers[0]).toEqual(expect.objectContaining(defaultTransfer));
    expect(transfers[1]).toEqual(expect.objectContaining(transfer2));
  });

  it('Successfully finds all transfers to and from subaccount w/ event id', async () => {
    const transfer2: TransferCreateObject = {
      senderSubaccountId: defaultSubaccountId2,
      recipientSubaccountId: defaultSubaccountId,
      assetId: defaultAsset2.id,
      size: '5',
      eventId: defaultTendermintEventId2,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: createdDateTime.toISO(),
      createdAtHeight: createdHeight,
    };
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create(transfer2),
    ]);

    const { results: transfers } = await TransferTable.findAllToOrFromSubaccountId(
      {
        subaccountId: [defaultSubaccountId],
        eventId: [defaultTendermintEventId],
      },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(transfers.length).toEqual(1);
    expect(transfers[0]).toEqual(expect.objectContaining(defaultTransfer));
  });

  it('Successfully finds all transfers to and from subaccount using pagination', async () => {
    const transfer2: TransferCreateObject = {
      senderSubaccountId: defaultSubaccountId2,
      recipientSubaccountId: defaultSubaccountId,
      assetId: defaultAsset2.id,
      size: '5',
      eventId: defaultTendermintEventId2,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: createdDateTime.toISO(),
      createdAtHeight: createdHeight,
    };
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create(transfer2),
    ]);

    const responsePageOne = await TransferTable.findAllToOrFromSubaccountId(
      { subaccountId: [defaultSubaccountId], page: 1, limit: 1 },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(responsePageOne.results.length).toEqual(1);
    expect(responsePageOne.results[0]).toEqual(expect.objectContaining(defaultTransfer));
    expect(responsePageOne.offset).toEqual(0);
    expect(responsePageOne.total).toEqual(2);

    const responsePageTwo = await TransferTable.findAllToOrFromSubaccountId(
      { subaccountId: [defaultSubaccountId], page: 2, limit: 1 },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(responsePageTwo.results.length).toEqual(1);
    expect(responsePageTwo.results[0]).toEqual(expect.objectContaining(transfer2));
    expect(responsePageTwo.offset).toEqual(1);
    expect(responsePageTwo.total).toEqual(2);

    const responsePageAllPages = await TransferTable.findAllToOrFromSubaccountId(
      { subaccountId: [defaultSubaccountId], page: 1, limit: 2 },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(responsePageAllPages.results.length).toEqual(2);
    expect(responsePageAllPages.results[0]).toEqual(expect.objectContaining(defaultTransfer));
    expect(responsePageAllPages.results[1]).toEqual(expect.objectContaining(transfer2));
    expect(responsePageAllPages.offset).toEqual(0);
    expect(responsePageAllPages.total).toEqual(2);
  });

  it('Successfully finds Transfer with eventId', async () => {
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create({
        ...defaultTransfer,
        eventId: defaultTendermintEventId2,
      }),
    ]);

    const transfers: TransferFromDatabase[] = await TransferTable.findAll(
      {
        eventId: [defaultTendermintEventId2],
      },
      [],
      { readReplica: true },
    );

    expect(transfers.length).toEqual(1);
    expect(transfers[0]).toEqual(expect.objectContaining({
      ...defaultTransfer,
      eventId: defaultTendermintEventId2,
    }));
  });

  it('Successfully finds all Transfers before or at the height', async () => {
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create({
        ...defaultTransfer,
        eventId: defaultTendermintEventId2,
        createdAtHeight: '5',
      }),
    ]);

    const transfers: TransferFromDatabase[] = await TransferTable.findAll(
      { createdBeforeOrAtHeight: defaultTransfer.createdAtHeight }, [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(transfers.length).toEqual(1);
    expect(transfers[0]).toEqual(expect.objectContaining(defaultTransfer));
  });

  it('Successfully finds all Transfers created after height', async () => {
    const transfer2: TransferCreateObject = {
      ...defaultTransfer,
      eventId: defaultTendermintEventId2,
      createdAtHeight: '5',
    };
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create(transfer2),
    ]);

    const transfers: TransferFromDatabase[] = await TransferTable.findAll(
      { createdAfterHeight: '4' }, [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(transfers.length).toEqual(1);
    expect(transfers[0]).toEqual(expect.objectContaining(transfer2));
  });

  it('Successfully finds all transfers to and from subaccount created before time', async () => {
    const transfer2: TransferCreateObject = {
      senderSubaccountId: defaultSubaccountId2,
      recipientSubaccountId: defaultSubaccountId,
      assetId: defaultAsset2.id,
      size: '5',
      eventId: defaultTendermintEventId2,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: '1982-05-25T00:00:00.000Z',
      createdAtHeight: '100',
    };
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create(transfer2),
    ]);

    const { results: transfers } = await TransferTable.findAllToOrFromSubaccountId(
      {
        subaccountId: [defaultSubaccountId],
        createdBeforeOrAt: '2000-05-25T00:00:00.000Z',
      },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(transfers.length).toEqual(1);
    expect(transfers[0]).toEqual(expect.objectContaining(transfer2));
  });

  it('Successfully finds a Transfer', async () => {
    await TransferTable.create(defaultTransfer);

    const transfer: TransferFromDatabase | undefined = await TransferTable.findById(
      TransferTable.uuid(
        defaultTransfer.eventId,
        defaultTransfer.assetId,
        defaultTransfer.senderSubaccountId,
        defaultTransfer.recipientSubaccountId,
        defaultTransfer.senderWalletAddress,
        defaultTransfer.recipientWalletAddress,
      ),
    );

    expect(transfer).toEqual(expect.objectContaining(defaultTransfer));
  });

  it('Recipient/sender must exist', async () => {
    await WalletTable.create({
      address: defaultWalletAddress,
      totalTradingRewards: '0',
      totalVolume: '0',
    });
    const invalidDeposit: TransferCreateObject = {
      ...defaultDeposit,
      recipientWalletAddress: defaultWalletAddress,
    };

    await expect(TransferTable.create(invalidDeposit)).rejects.toBeInstanceOf(CheckViolationError);
    const invalidWithdrawal: TransferCreateObject = {
      ...defaultWithdrawal,
      senderWalletAddress: defaultWalletAddress,
    };

    await expect(TransferTable.create(invalidWithdrawal))
      .rejects.toBeInstanceOf(CheckViolationError);

    const invalidTfer: TransferCreateObject = {
      recipientWalletAddress: defaultWalletAddress,
      assetId: defaultAsset.id,
      size: '10',
      eventId: defaultTendermintEventId,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: createdDateTime.toISO(),
      createdAtHeight: createdHeight,
    };
    await expect(TransferTable.create(invalidTfer)).rejects.toBeInstanceOf(CheckViolationError);
  });

  it('Successfully creates/finds a transfer/deposit/withdrawal', async () => {
    await WalletTable.create({
      address: defaultWalletAddress,
      totalTradingRewards: '0',
      totalVolume: '0',
    });
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create(defaultDeposit),
      TransferTable.create(defaultWithdrawal),
    ]);

    const [
      transfer,
      deposit,
      withdrawal,
    ]: [
      TransferFromDatabase | undefined,
      TransferFromDatabase | undefined,
      TransferFromDatabase | undefined,
    ] = await Promise.all([
      TransferTable.findById(
        TransferTable.uuid(
          defaultTransfer.eventId,
          defaultTransfer.assetId,
          defaultTransfer.senderSubaccountId,
          defaultTransfer.recipientSubaccountId,
          defaultTransfer.senderWalletAddress,
          defaultTransfer.recipientWalletAddress,
        ),
      ),
      TransferTable.findById(
        TransferTable.uuid(
          defaultDeposit.eventId,
          defaultDeposit.assetId,
          defaultDeposit.senderSubaccountId,
          defaultDeposit.recipientSubaccountId,
          defaultDeposit.senderWalletAddress,
          defaultDeposit.recipientWalletAddress,
        ),
      ),
      TransferTable.findById(
        TransferTable.uuid(
          defaultWithdrawal.eventId,
          defaultWithdrawal.assetId,
          defaultWithdrawal.senderSubaccountId,
          defaultWithdrawal.recipientSubaccountId,
          defaultWithdrawal.senderWalletAddress,
          defaultWithdrawal.recipientWalletAddress,
        ),
      ),
    ]);

    expect(transfer).toEqual(expect.objectContaining(defaultTransfer));
    expect(deposit).toEqual(expect.objectContaining(defaultDeposit));
    expect(withdrawal).toEqual(expect.objectContaining(defaultWithdrawal));
  });

  it('Successfully gets total transfers per subaccount', async () => {
    await SubaccountTable.create(defaultSubaccount3);
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create(defaultTransfer2),
      TransferTable.create(defaultTransfer3),
    ]);

    const transferMap: SubaccountAssetNetTransferMap = await
    TransferTable.getNetTransfersPerSubaccount(createdHeight);

    expect(transferMap).toEqual(expect.objectContaining({
      [defaultSubaccountId]: {
        [defaultAsset.id]: '-10',
      },
      [defaultSubaccountId2]: {
        [defaultAsset.id]: '15',
        [defaultAsset2.id]: '5',
      },
      [defaultSubaccountId3]: {
        [defaultAsset.id]: '-5',
        [defaultAsset2.id]: '-5',
      },
    }));
  });

  it('Successfully gets total transfers per subaccount with duplicate transfer amounts', async () => {
    await SubaccountTable.create(defaultSubaccount3);
    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create({
        ...defaultTransfer,
        eventId: defaultTendermintEventId2,
      }),
      TransferTable.create(defaultTransfer2),
      TransferTable.create(defaultTransfer3),
    ]);

    const transferMap: SubaccountAssetNetTransferMap = await
    TransferTable.getNetTransfersPerSubaccount(createdHeight);

    expect(transferMap).toEqual(expect.objectContaining({
      [defaultSubaccountId]: {
        [defaultAsset.id]: '-20',
      },
      [defaultSubaccountId2]: {
        [defaultAsset.id]: '25',
        [defaultAsset2.id]: '5',
      },
      [defaultSubaccountId3]: {
        [defaultAsset.id]: '-5',
        [defaultAsset2.id]: '-5',
      },
    }));
  });

  it('Successfully gets total decimal value transfers per subaccount and respects createdBeforeOrAtHeight', async () => {
    await SubaccountTable.create(defaultSubaccount3);
    await Promise.all([
      TransferTable.create({
        ...defaultTransfer,
        size: '10.5',
      }),
      TransferTable.create({
        ...defaultTransfer2,
        size: '5.2',
      }),
      // this transfer is ignored because createdAtHeight > createdBeforeOrAtHeight
      TransferTable.create({
        ...defaultTransfer3,
        size: '5.3',
        createdAtHeight: '5',
      }),
    ]);

    const transferMap: SubaccountAssetNetTransferMap = await
    TransferTable.getNetTransfersPerSubaccount(createdHeight);

    expect(transferMap).toEqual(expect.objectContaining({
      [defaultSubaccountId]: {
        [defaultAsset.id]: '-10.5',
      },
      [defaultSubaccountId2]: {
        [defaultAsset.id]: '15.7',
      },
      [defaultSubaccountId3]: {
        [defaultAsset.id]: '-5.2',
      },
    }));
  });

  it('Successfully gets net transfers between block heights for a subaccount', async () => {
    await SubaccountTable.create(defaultSubaccount3);
    await Promise.all([
      TransferTable.create({
        ...defaultTransfer,
        createdAtHeight: '1',
      }),
      TransferTable.create({
        ...defaultTransfer,
        size: '10.5',
        createdAtHeight: '2',
        eventId: defaultTendermintEventId2,
      }),
      TransferTable.create({
        ...defaultTransfer2,
        size: '5.2',
        createdAtHeight: '3',
      }),
      TransferTable.create({
        ...defaultTransfer3,
        size: '5.3',
        createdAtHeight: '5',
      }),
    ]);

    const [
      transferMap,
      transferMap2,
      transferMap21,
      transferMap3,
    ]: [
      AssetTransferMap,
      AssetTransferMap,
      AssetTransferMap,
      AssetTransferMap,
    ] = await Promise.all([
      TransferTable.getNetTransfersBetweenBlockHeightsForSubaccount(
        defaultSubaccountId,
        '1',
        '3',
      ),
      TransferTable.getNetTransfersBetweenBlockHeightsForSubaccount(
        defaultSubaccountId2,
        '1',
        '3',
      ),
      TransferTable.getNetTransfersBetweenBlockHeightsForSubaccount(
        defaultSubaccountId2,
        '1',
        '5',
      ),
      TransferTable.getNetTransfersBetweenBlockHeightsForSubaccount(
        defaultSubaccountId3,
        '2',
        '5',
      ),
    ]);
    expect(transferMap).toEqual({
      [defaultAsset.id]: Big('-10.5'),
    });
    expect(transferMap2).toEqual({
      [defaultAsset.id]: Big('15.7'),
    });
    expect(transferMap21).toEqual({
      [defaultAsset.id]: Big('15.7'),
      [defaultAsset2.id]: Big('5.3'),
    });
    expect(transferMap3).toEqual({
      [defaultAsset.id]: Big('-5.2'),
      [defaultAsset2.id]: Big('-5.3'),
    });
  });

  it('Successfully gets the latest createdAt for subaccounts', async () => {
    const now = DateTime.utc();

    const transfer2 = {
      senderSubaccountId: defaultSubaccountId2,
      recipientSubaccountId: defaultSubaccountId,
      assetId: defaultAsset2.id,
      size: '5',
      eventId: defaultTendermintEventId3,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: now.minus({ hours: 2 }).toISO(),
      createdAtHeight: createdHeight,
    };

    const transfer3 = {
      senderSubaccountId: defaultSubaccountId2,
      recipientSubaccountId: defaultSubaccountId,
      assetId: defaultAsset2.id,
      size: '5',
      eventId: defaultTendermintEventId2,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: now.minus({ hours: 1 }).toISO(),
      createdAtHeight: createdHeight,
    };

    await Promise.all([
      TransferTable.create(defaultTransfer),
      TransferTable.create(transfer2),
      TransferTable.create(transfer3),
    ]);

    const transferTimes: { [subaccountId: string]: string } = await
    TransferTable.getLastTransferTimeForSubaccounts(
      [defaultSubaccountId, defaultSubaccountId2],
    );

    expect(transferTimes[defaultSubaccountId]).toEqual(defaultTransfer.createdAt);
    expect(transferTimes[defaultSubaccountId2]).toEqual(defaultTransfer.createdAt);
  });

  describe('getNetTransfersBetweenSubaccountIds', () => {
    it('Successfully gets total net Transfers between two subaccounts', async () => {
      await Promise.all([
        TransferTable.create({
          ...defaultTransfer,
          size: '20',
        }),
        TransferTable.create({
          ...defaultTransfer,
          size: '30',
          eventId: defaultTendermintEventId2,
        }),
        TransferTable.create({
          ...defaultTransfer,
          senderSubaccountId: defaultSubaccountId2,
          recipientSubaccountId: defaultSubaccountId,
          size: '10',
          eventId: defaultTendermintEventId3,
        }),
      ]);

      const netTransfers: string = await TransferTable.getNetTransfersBetweenSubaccountIds(
        defaultSubaccountId,
        defaultSubaccountId2,
        USDC_ASSET_ID,
      );

      expect(netTransfers).toEqual('40'); // 20 + 30 - 10

      // Test the other way around
      const negativeNetTransfers: string = await TransferTable.getNetTransfersBetweenSubaccountIds(
        defaultSubaccountId2,
        defaultSubaccountId,
        USDC_ASSET_ID,
      );

      expect(negativeNetTransfers).toEqual('-40'); // 10 - 20 - 30
    });

    it('Successfully gets total net Transfers between two subaccounts with no transfers', async () => {
      const netTransfers: string = await TransferTable.getNetTransfersBetweenSubaccountIds(
        defaultSubaccountId,
        defaultSubaccountId2,
        USDC_ASSET_ID,
      );

      expect(netTransfers).toEqual('0');
    });
  });

  describe('findAllToOrFromParentSubaccount', () => {
    it('Successfully excludes transfers between child subaccounts of the same parent', async () => {
    // defaultSubaccount (subaccount 0) -> isolatedSubaccount (subaccount 128)
    // Both have parent subaccount 0 (0 % 128 = 0, 128 % 128 = 0)
      const sameParentTransfer: TransferCreateObject = {
        senderSubaccountId: defaultSubaccountId,
        recipientSubaccountId: isolatedSubaccountId,
        assetId: defaultAsset.id,
        size: '100',
        eventId: defaultTendermintEventId,
        transactionHash: '',
        createdAt: createdDateTime.toISO(),
        createdAtHeight: createdHeight,
      };

      // defaultSubaccount (parent 0) -> defaultSubaccount2 (parent 1)
      const differentParentTransfer: TransferCreateObject = {
        senderSubaccountId: defaultSubaccountId,
        recipientSubaccountId: defaultSubaccountId2,
        assetId: defaultAsset.id,
        size: '200',
        eventId: defaultTendermintEventId2,
        transactionHash: '',
        createdAt: createdDateTime.toISO(),
        createdAtHeight: createdHeight,
      };

      await Promise.all([
        TransferTable.create(sameParentTransfer),
        TransferTable.create(differentParentTransfer),
      ]);

      const { results: transfers } = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [defaultSubaccountId, isolatedSubaccountId],
          limit: 100,
        },
        [],
      );

      expect(transfers.length).toEqual(1);
      expect(transfers[0]).toEqual(expect.objectContaining(differentParentTransfer));
    });

    it('Successfully includes transfers from different addresses', async () => {
      const crossAddressTransfer: TransferCreateObject = {
        senderSubaccountId: defaultSubaccountId2,
        recipientSubaccountId: defaultSubaccountId,
        assetId: defaultAsset.id,
        size: '150',
        eventId: defaultTendermintEventId,
        transactionHash: '',
        createdAt: createdDateTime.toISO(),
        createdAtHeight: createdHeight,
      };

      await TransferTable.create(crossAddressTransfer);

      const { results: transfers } = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [defaultSubaccountId],
          limit: 100,
        },
        [],
      );

      expect(transfers.length).toEqual(1);
      expect(transfers[0]).toEqual(expect.objectContaining(crossAddressTransfer));
    });

    it('Successfully includes deposits to child subaccounts', async () => {
      await WalletTable.create({
        address: defaultWalletAddress,
        totalTradingRewards: '0',
        totalVolume: '0',
      });

      const deposit: TransferCreateObject = {
        senderWalletAddress: defaultWalletAddress,
        recipientSubaccountId: isolatedSubaccountId,
        assetId: defaultAsset.id,
        size: '500',
        eventId: defaultTendermintEventId,
        transactionHash: '',
        createdAt: createdDateTime.toISO(),
        createdAtHeight: createdHeight,
      };

      await TransferTable.create(deposit);

      const { results: transfers } = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [isolatedSubaccountId],
          limit: 100,
        },
        [],
      );

      expect(transfers.length).toEqual(1);
      expect(transfers[0]).toEqual(expect.objectContaining(deposit));
    });

    it('Successfully includes withdrawals from child subaccounts', async () => {
      await WalletTable.create({
        address: defaultWalletAddress,
        totalTradingRewards: '0',
        totalVolume: '0',
      });

      const withdrawal: TransferCreateObject = {
        senderSubaccountId: isolatedSubaccountId,
        recipientWalletAddress: defaultWalletAddress,
        assetId: defaultAsset.id,
        size: '300',
        eventId: defaultTendermintEventId,
        transactionHash: '',
        createdAt: createdDateTime.toISO(),
        createdAtHeight: createdHeight,
      };

      await TransferTable.create(withdrawal);

      const { results: transfers } = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [isolatedSubaccountId],
          limit: 100,
        },
        [],
      );

      expect(transfers.length).toEqual(1);
      expect(transfers[0]).toEqual(expect.objectContaining(withdrawal));
    });

    it('Successfully respects limit parameter', async () => {
    // Create 5 blocks first
      const blocks: BlockCreateObject[] = [];
      for (let i = 0; i < 5; i++) {
        blocks.push({
          blockHeight: (3 + i).toString(),
          time: createdDateTime.plus({ minutes: i }).toISO(),
        });
      }
      await Promise.all(blocks.map((b) => BlockTable.create(b)));

      // Create 5 TendermintEvents
      const events: TendermintEventCreateObject[] = [];
      for (let i = 0; i < 5; i++) {
        events.push({
          blockHeight: (3 + i).toString(),
          transactionIndex: 0,
          eventIndex: 0,
        });
      }
      await Promise.all(events.map((e) => TendermintEventTable.create(e)));

      // Create 5 transfers
      const transfers: TransferCreateObject[] = [];
      for (let i = 0; i < 5; i++) {
        const eventId = TendermintEventTable.createEventId(
          events[i].blockHeight,
          events[i].transactionIndex,
          events[i].eventIndex,
        );

        transfers.push({
          senderSubaccountId: defaultSubaccountId,
          recipientSubaccountId: defaultSubaccountId2,
          assetId: defaultAsset.id,
          size: `${i + 1}`,
          eventId,
          transactionHash: `hash${i}`,
          createdAt: createdDateTime.plus({ minutes: i }).toISO(),
          createdAtHeight: (parseInt(createdHeight, 10) + i).toString(),
        });
      }

      await Promise.all(transfers.map((t) => TransferTable.create(t)));

      const { results: resultTransfers } = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [defaultSubaccountId],
          limit: 3,
        },
        [],
      );

      expect(resultTransfers.length).toEqual(3);
    });

    it('Successfully finds all transfers to and from parent subaccount using pagination', async () => {
    // Create 4 blocks first
      const blocks: BlockCreateObject[] = [];
      for (let i = 0; i < 4; i++) {
        blocks.push({
          blockHeight: (3 + i).toString(),
          time: createdDateTime.plus({ minutes: i }).toISO(),
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

      // Create 4 transfers
      const transfers: TransferCreateObject[] = [];
      for (let i = 0; i < 4; i++) {
        const eventId = TendermintEventTable.createEventId(
          events[i].blockHeight,
          events[i].transactionIndex,
          events[i].eventIndex,
        );

        transfers.push({
          senderSubaccountId: defaultSubaccountId,
          recipientSubaccountId: defaultSubaccountId2,
          assetId: defaultAsset.id,
          size: `${i + 1}`,
          eventId,
          transactionHash: `hash${i}`,
          createdAt: createdDateTime.plus({ minutes: i }).toISO(),
          createdAtHeight: (parseInt(createdHeight, 10) + i).toString(),
        });
      }

      await Promise.all(transfers.map((t) => TransferTable.create(t)));

      const responsePageOne = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [defaultSubaccountId],
          limit: 2,
          page: 1,
        },
        [],
        { orderBy: [[TransferColumns.id, Ordering.ASC]] },
      );

      expect(responsePageOne.results.length).toEqual(2);
      expect(responsePageOne.offset).toEqual(0);
      expect(responsePageOne.total).toEqual(4);
      expect(responsePageOne.limit).toEqual(2);

      const responsePageTwo = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [defaultSubaccountId],
          limit: 2,
          page: 2,
        },
        [],
        { orderBy: [[TransferColumns.id, Ordering.ASC]] },
      );

      expect(responsePageTwo.results.length).toEqual(2);
      expect(responsePageTwo.offset).toEqual(2);
      expect(responsePageTwo.total).toEqual(4);
      expect(responsePageTwo.limit).toEqual(2);
    });

    it('Successfully finds all transfers before or at the height', async () => {
      const transfer1: TransferCreateObject = {
        senderSubaccountId: defaultSubaccountId,
        recipientSubaccountId: defaultSubaccountId2,
        assetId: defaultAsset.id,
        size: '100',
        eventId: defaultTendermintEventId,
        transactionHash: '',
        createdAt: createdDateTime.toISO(),
        createdAtHeight: '10',
      };

      const transfer2: TransferCreateObject = {
        senderSubaccountId: defaultSubaccountId,
        recipientSubaccountId: defaultSubaccountId2,
        assetId: defaultAsset.id,
        size: '200',
        eventId: defaultTendermintEventId2,
        transactionHash: '',
        createdAt: createdDateTime.plus({ minutes: 1 }).toISO(),
        createdAtHeight: '20',
      };

      await Promise.all([
        TransferTable.create(transfer1),
        TransferTable.create(transfer2),
      ]);

      const { results: transfers } = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [defaultSubaccountId],
          limit: 100,
          createdBeforeOrAtHeight: '10',
        },
        [],
      );

      expect(transfers.length).toEqual(1);
      expect(transfers[0]).toEqual(expect.objectContaining(transfer1));
    });

    it('Successfully finds all transfers before or at the time', async () => {
      const createdAt1 = '2000-01-01T00:00:00.000Z';
      const createdAt2 = '2000-01-02T00:00:00.000Z';

      const transfer1: TransferCreateObject = {
        senderSubaccountId: defaultSubaccountId,
        recipientSubaccountId: defaultSubaccountId2,
        assetId: defaultAsset.id,
        size: '100',
        eventId: defaultTendermintEventId,
        transactionHash: '',
        createdAt: createdAt1,
        createdAtHeight: createdHeight,
      };

      const transfer2: TransferCreateObject = {
        senderSubaccountId: defaultSubaccountId,
        recipientSubaccountId: defaultSubaccountId2,
        assetId: defaultAsset.id,
        size: '200',
        eventId: defaultTendermintEventId2,
        transactionHash: '',
        createdAt: createdAt2,
        createdAtHeight: (parseInt(createdHeight, 10) + 1).toString(),
      };

      await Promise.all([
        TransferTable.create(transfer1),
        TransferTable.create(transfer2),
      ]);

      const { results: transfers } = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [defaultSubaccountId],
          limit: 100,
          createdBeforeOrAt: createdAt1,
        },
        [],
      );

      expect(transfers.length).toEqual(1);
      expect(transfers[0]).toEqual(expect.objectContaining(transfer1));
    });

    it('Successfully returns empty results when no transfers match criteria', async () => {
      const { results: transfers } = await TransferTable.findAllToOrFromParentSubaccount(
        {
          subaccountId: [defaultSubaccountId],
          limit: 100,
        },
        [],
      );

      expect(transfers.length).toEqual(0);
    });
  });
});
