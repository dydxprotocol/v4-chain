import {
  IsoString,
  Ordering,
  PnlTicksColumns,
  PnlTicksCreateObject,
  PnlTicksFromDatabase,
} from '../../src/types';
import * as PnlTicksTable from '../../src/stores/pnl-ticks-table';
import * as BlockTable from '../../src/stores/block-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import {
  defaultBlock, defaultBlock2,
  defaultPnlTick,
  defaultSubaccountId,
  defaultSubaccountId2,
} from '../helpers/constants';
import { DateTime } from 'luxon';
import { ZERO_TIME_ISO_8601 } from '../../src/constants';

describe('PnlTicks store', () => {
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

  it('Successfully creates a Pnl Tick', async () => {
    await PnlTicksTable.create(defaultPnlTick);
  });

  it('Successfully creates multiple Pnl Ticks', async () => {
    await BlockTable.create({
      ...defaultBlock,
      blockHeight: '5',
    });
    const pnlTick2: PnlTicksCreateObject = {
      subaccountId: defaultSubaccountId,
      equity: '5',
      totalPnl: '5',
      netTransfers: '5',
      createdAt: '2020-01-01T00:00:00.000Z',
      blockHeight: '5',
      blockTime: defaultBlock.time,
    };
    await Promise.all([
      PnlTicksTable.create(defaultPnlTick),
      PnlTicksTable.create(pnlTick2),
    ]);

    const pnlTicks: PnlTicksFromDatabase[] = await PnlTicksTable.findAll({}, [], {
      orderBy: [[PnlTicksColumns.blockHeight, Ordering.ASC]],
    });

    expect(pnlTicks.length).toEqual(2);
    expect(pnlTicks[0]).toEqual(expect.objectContaining(defaultPnlTick));
    expect(pnlTicks[1]).toEqual(expect.objectContaining(pnlTick2));
  });

  it('createMany Pnl Ticks', async () => {
    const pnlTick2: PnlTicksCreateObject = {
      subaccountId: defaultSubaccountId,
      equity: '5',
      totalPnl: '5',
      netTransfers: '5',
      createdAt: '2020-01-01T00:00:00.000Z',
      blockHeight: '5',
      blockTime: defaultBlock.time,
    };
    await PnlTicksTable.createMany([defaultPnlTick, pnlTick2]);
    const pnlTicks: PnlTicksFromDatabase[] = await PnlTicksTable.findAll({}, [], {
      orderBy: [[PnlTicksColumns.blockHeight, Ordering.ASC]],
    });

    expect(pnlTicks.length).toEqual(2);
    expect(pnlTicks[0]).toEqual(expect.objectContaining(defaultPnlTick));
    expect(pnlTicks[1]).toEqual(expect.objectContaining(pnlTick2));
  });

  it('Successfully finds PnlTicks with subaccountId', async () => {
    await Promise.all([
      PnlTicksTable.create(defaultPnlTick),
      PnlTicksTable.create({
        ...defaultPnlTick,
        createdAt: '2020-01-01T00:00:00.000Z',
      }),
      PnlTicksTable.create({
        ...defaultPnlTick,
        subaccountId: defaultSubaccountId2,
        createdAt: '2020-01-01T00:00:00.000Z',
      }),
    ]);

    const pnlTicks: PnlTicksFromDatabase[] = await PnlTicksTable.findAll(
      {
        subaccountId: [defaultSubaccountId],
      },
      [],
      { readReplica: true },
    );

    expect(pnlTicks.length).toEqual(2);
  });

  it('Successfully finds latest block time', async () => {
    const blockTime: IsoString = '2023-01-01T00:00:00.000Z';
    await Promise.all([
      PnlTicksTable.create(defaultPnlTick),
      PnlTicksTable.create({
        ...defaultPnlTick,
        createdAt: '2020-01-01T00:00:00.000Z',
        blockHeight: '1000',
        blockTime,
      }),
    ]);

    const latestBlocktime: string = await PnlTicksTable.findLatestProcessedBlocktime();

    expect(latestBlocktime).toEqual(blockTime);
  });

  it('Successfully finds latest block time without any pnl ticks', async () => {
    const latestBlocktime: string = await PnlTicksTable.findLatestProcessedBlocktime();
    expect(latestBlocktime).toEqual(ZERO_TIME_ISO_8601);
  });

  it('createMany PnlTicks, find most recent pnl ticks for each account', async () => {
    await Promise.all([
      BlockTable.create({
        blockHeight: '3',
        time: defaultBlock.time,
      }),
      BlockTable.create({
        blockHeight: '5',
        time: defaultBlock.time,
      }),
    ]);
    await PnlTicksTable.createMany([
      {
        subaccountId: defaultSubaccountId,
        equity: '1092',
        createdAt: DateTime.utc().minus({ hours: 1 }).toISO(),
        totalPnl: '1000',
        netTransfers: '50',
        blockHeight: defaultBlock.blockHeight,
        blockTime: defaultBlock.time,
      },
      {
        subaccountId: defaultSubaccountId,
        equity: '1097',
        createdAt: DateTime.utc().minus({ hours: 3 }).toISO(),
        totalPnl: '1000',
        netTransfers: '50',
        blockHeight: '3',
        blockTime: defaultBlock.time,
      },
      {
        subaccountId: defaultSubaccountId,
        equity: '1011',
        createdAt: DateTime.utc().minus({ hours: 11 }).toISO(),
        totalPnl: '1000',
        netTransfers: '50',
        blockHeight: '5',
        blockTime: defaultBlock.time,
      },
      {
        subaccountId: defaultSubaccountId,
        equity: '1014',
        createdAt: DateTime.utc().minus({ hours: 9 }).toISO(),
        totalPnl: '1000',
        netTransfers: '50',
        blockHeight: '5',
        blockTime: defaultBlock.time,
      },
      {
        subaccountId: defaultSubaccountId2,
        equity: '100',
        createdAt: new Date().toISOString(),
        totalPnl: '1000',
        netTransfers: '50',
        blockHeight: '2',
        blockTime: defaultBlock2.time,
      },
      {
        subaccountId: defaultSubaccountId2,
        equity: '200',
        createdAt: DateTime.utc().minus({ hours: 9 }).toISO(),
        totalPnl: '1000',
        netTransfers: '50',
        blockHeight: '5',
        blockTime: defaultBlock.time,
      },
    ]);

    const mostRecent: {
      [accountId: string]: PnlTicksCreateObject
    } = await PnlTicksTable.findMostRecentPnlTickForEachAccount(
      '3',
    );
    expect(mostRecent[defaultSubaccountId].equity).toEqual('1014');
    expect(mostRecent[defaultSubaccountId2].equity).toEqual('200');
  });
});
