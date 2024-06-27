import { SubaccountCreateObject, SubaccountFromDatabase } from '../../src/types';
import * as SubaccountTable from '../../src/stores/subaccount-table';
import * as TransferTable from '../../src/stores/transfer-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  createdDateTime,
  createdHeight,
  defaultAddress,
  defaultAsset,
  defaultBlock,
  defaultBlock2,
  defaultSubaccount,
  defaultSubaccount2,
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultTendermintEvent,
  defaultTransfer,
} from '../helpers/constants';
import * as AssetTable from '../../src/stores/asset-table';
import * as BlockTable from '../../src/stores/block-table';
import * as TendermintEventTable from '../../src/stores/tendermint-event-table';
import Transaction from '../../src/helpers/transaction';
import _ from 'lodash';
import { DateTime } from 'luxon';

describe('Subaccount store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a Subaccount', async () => {
    await SubaccountTable.create(defaultSubaccount);
  });

  it('Successfully finds all Subaccounts', async () => {
    await Promise.all([
      SubaccountTable.create(defaultSubaccount),
      SubaccountTable.create({
        ...defaultSubaccount,
        subaccountNumber: 1,
      }),
    ]);

    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(subaccounts.length).toEqual(2);
    expect(subaccounts[0]).toEqual(expect.objectContaining(defaultSubaccount));
    expect(subaccounts[1]).toEqual(expect.objectContaining({
      ...defaultSubaccount,
      subaccountNumber: 1,
    }));
  });

  it('Successfully finds Subaccount with address', async () => {
    await Promise.all([
      SubaccountTable.create(defaultSubaccount),
      SubaccountTable.create({
        ...defaultSubaccount,
        address: 'fake_address',
      }),
    ]);

    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        address: defaultSubaccount.address,
      },
      [],
      { readReplica: true },
    );

    expect(subaccounts.length).toEqual(1);
    expect(subaccounts[0]).toEqual(expect.objectContaining(defaultSubaccount));
  });

  it('Successfully finds Subaccount with updatedBeforeOrAt', async () => {
    await Promise.all([
      SubaccountTable.create(defaultSubaccount),
      SubaccountTable.create({
        ...defaultSubaccount,
        address: 'fake_address',
        updatedAt: DateTime.fromISO(defaultSubaccount.updatedAt).plus(1).toISO(),
      }),
    ]);

    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        updatedBeforeOrAt: defaultSubaccount.updatedAt,
      },
      [],
      { readReplica: true },
    );

    expect(subaccounts.length).toEqual(1);
    expect(subaccounts[0]).toEqual(expect.objectContaining(defaultSubaccount));
  });

  it('Successfully finds Subaccount with updatedOnOrAfter', async () => {
    await Promise.all([
      SubaccountTable.create(defaultSubaccount),
      SubaccountTable.create({
        ...defaultSubaccount,
        address: 'fake_address',
        updatedAt: DateTime.fromISO(defaultSubaccount.updatedAt).minus(10).toISO(),
      }),
    ]);

    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        updatedOnOrAfter: defaultSubaccount.updatedAt,
      },
      [],
      { readReplica: true },
    );

    expect(subaccounts.length).toEqual(1);
    expect(subaccounts[0]).toEqual(expect.objectContaining(defaultSubaccount));
  });

  it('Successfully finds a Subaccount', async () => {
    await SubaccountTable.create(defaultSubaccount);

    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
      SubaccountTable.uuid(defaultAddress, 0),
    );

    expect(subaccount).toEqual(expect.objectContaining(defaultSubaccount));
  });

  it('Successfully finds all subaccounts with transfers', async () => {

    const defaultSubaccount3: SubaccountCreateObject = {
      address: defaultAddress,
      subaccountNumber: 3,
      updatedAt: createdDateTime.toISO(),
      updatedAtHeight: createdHeight,
    };
    await Promise.all([
      SubaccountTable.create(defaultSubaccount),
      SubaccountTable.create(defaultSubaccount2),
      SubaccountTable.create(defaultSubaccount3),
      BlockTable.create(defaultBlock),
      AssetTable.create(defaultAsset),
    ]);

    await TendermintEventTable.create(defaultTendermintEvent);
    await TransferTable.create(defaultTransfer);

    const subaccounts: SubaccountFromDatabase[] = await
    SubaccountTable.getSubaccountsWithTransfers(defaultTransfer.createdAtHeight);
    const subaccountIds: string[] = _.map(subaccounts, 'id');
    expect(subaccountIds).toEqual(
      expect.arrayContaining([defaultSubaccountId, defaultSubaccountId2]),
    );
  });

  it('Successfully finds all subaccounts with transfers respects createdAtHeight', async () => {
    const defaultSubaccount3: SubaccountCreateObject = {
      address: defaultAddress,
      subaccountNumber: 3,
      updatedAt: createdDateTime.toISO(),
      updatedAtHeight: createdHeight,
    };
    await Promise.all([
      SubaccountTable.create(defaultSubaccount),
      SubaccountTable.create(defaultSubaccount2),
      SubaccountTable.create(defaultSubaccount3),
      BlockTable.create(defaultBlock),
      AssetTable.create(defaultAsset),
    ]);

    await TendermintEventTable.create(defaultTendermintEvent);
    await TransferTable.create(defaultTransfer);

    const subaccounts: SubaccountFromDatabase[] = await
    SubaccountTable.getSubaccountsWithTransfers('1');
    const subaccountIds: string[] = _.map(subaccounts, 'id');
    expect(subaccountIds.length).toEqual(0);
  });

  it('Successfully finds all subaccounts with transfers respects options', async () => {

    const defaultSubaccount3: SubaccountCreateObject = {
      address: defaultAddress,
      subaccountNumber: 3,
      updatedAt: createdDateTime.toISO(),
      updatedAtHeight: createdHeight,
    };
    await Promise.all([
      SubaccountTable.create(defaultSubaccount),
      SubaccountTable.create(defaultSubaccount2),
      SubaccountTable.create(defaultSubaccount3),
      BlockTable.create(defaultBlock),
      AssetTable.create(defaultAsset),
    ]);

    await TendermintEventTable.create(defaultTendermintEvent);
    await TransferTable.create(defaultTransfer);
    const txId: number = await Transaction.start();
    const subaccounts: SubaccountFromDatabase[] = await
    SubaccountTable.getSubaccountsWithTransfers(
      defaultTransfer.createdAtHeight,
      { txId, readReplica: true },
    );
    const subaccountIds: string[] = _.map(subaccounts, 'id');
    await Transaction.rollback(txId);
    expect(subaccountIds).toEqual(
      expect.arrayContaining([defaultSubaccountId, defaultSubaccountId2]),
    );
  });

  it('Unable finds a Subaccount', async () => {
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
      SubaccountTable.uuid(defaultAddress, 0),
    );
    expect(subaccount).toEqual(undefined);
  });

  it('Successfully creates a Subaccount with updatedAtHeight', async () => {
    await BlockTable.create(defaultBlock);
    await SubaccountTable.create({
      ...defaultSubaccount,
      updatedAtHeight: defaultBlock.blockHeight,
      updatedAt: defaultBlock.time,
    });

    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
      SubaccountTable.uuid(defaultAddress, 0),
    );

    expect(subaccount).toEqual(expect.objectContaining({
      ...defaultSubaccount,
      updatedAtHeight: defaultBlock.blockHeight,
      updatedAt: defaultBlock.time,
    }));
  });

  it('Successfully updates a Subaccount with updatedAtHeight', async () => {
    await Promise.all([
      BlockTable.create(defaultBlock),
      BlockTable.create(defaultBlock2),
    ]);

    await SubaccountTable.create({
      ...defaultSubaccount,
      updatedAtHeight: defaultBlock.blockHeight,
      updatedAt: defaultBlock.time,
    });
    const subaccounts: SubaccountFromDatabase[] = await
    SubaccountTable.findAll({}, [], { readReplica: true });
    expect(subaccounts.length).toEqual(1);

    await SubaccountTable.update({
      id: defaultSubaccountId,
      updatedAtHeight: defaultBlock2.blockHeight,
      updatedAt: defaultBlock2.time,
    });
    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
      SubaccountTable.uuid(defaultAddress, 0),
    );

    expect(subaccount).toEqual(expect.objectContaining({
      ...defaultSubaccount,
      updatedAtHeight: defaultBlock2.blockHeight,
      updatedAt: defaultBlock2.time,
    }));
  });
});
