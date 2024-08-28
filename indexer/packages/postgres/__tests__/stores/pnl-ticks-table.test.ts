import {
  IsoString,
  LeaderboardPnlCreateObject,
  Ordering,
  PnlTicksColumns,
  PnlTicksCreateObject,
} from '../../src/types';
import * as PnlTicksTable from '../../src/stores/pnl-ticks-table';
import * as BlockTable from '../../src/stores/block-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import * as WalletTable from '../../src/stores/wallet-table';
import * as SubaccountTable from '../../src/stores/subaccount-table';
import {
  defaultAddress,
  defaultAddress2,
  defaultBlock,
  defaultBlock2,
  defaultPnlTick,
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultSubaccountIdWithAlternateAddress,
  defaultSubaccountWithAlternateAddress,
  defaultWallet2,
  vaultSubaccount,
  vaultSubaccountId,
  vaultWallet,
} from '../helpers/constants';
import { DateTime } from 'luxon';
import { ZERO_TIME_ISO_8601 } from '../../src/constants';

describe('PnlTicks store', () => {
  beforeEach(async () => {
    await seedData();
    await WalletTable.create(defaultWallet2);
    await SubaccountTable.create(defaultSubaccountWithAlternateAddress);
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

    const { results: pnlTicks } = await PnlTicksTable.findAll({}, [], {
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
    const { results: pnlTicks } = await PnlTicksTable.findAll({}, [], {
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

    const { results: pnlTicks } = await PnlTicksTable.findAll(
      {
        subaccountId: [defaultSubaccountId],
      },
      [],
      { readReplica: true },
    );

    expect(pnlTicks.length).toEqual(2);
  });

  it('Successfully finds PnlTicks using pagination', async () => {
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

    const responsePageOne = await PnlTicksTable.findAll({
      page: 1,
      limit: 1,
    }, [], {
      orderBy: [[PnlTicksColumns.blockHeight, Ordering.DESC]],
    });

    expect(responsePageOne.results.length).toEqual(1);
    expect(responsePageOne.results[0]).toEqual(expect.objectContaining({
      ...defaultPnlTick,
      createdAt: '2020-01-01T00:00:00.000Z',
      blockHeight: '1000',
      blockTime,
    }));
    expect(responsePageOne.offset).toEqual(0);
    expect(responsePageOne.total).toEqual(2);

    const responsePageTwo = await PnlTicksTable.findAll({
      page: 2,
      limit: 1,
    }, [], {
      orderBy: [[PnlTicksColumns.blockHeight, Ordering.DESC]],
    });

    expect(responsePageTwo.results.length).toEqual(1);
    expect(responsePageTwo.results[0]).toEqual(expect.objectContaining(defaultPnlTick));
    expect(responsePageTwo.offset).toEqual(1);
    expect(responsePageTwo.total).toEqual(2);

    const responsePageAllPages = await PnlTicksTable.findAll({
      page: 1,
      limit: 2,
    }, [], {
      orderBy: [[PnlTicksColumns.blockHeight, Ordering.DESC]],
    });

    expect(responsePageAllPages.results.length).toEqual(2);
    expect(responsePageAllPages.results[0]).toEqual(expect.objectContaining({
      ...defaultPnlTick,
      createdAt: '2020-01-01T00:00:00.000Z',
      blockHeight: '1000',
      blockTime,
    }));
    expect(responsePageAllPages.results[1]).toEqual(expect.objectContaining(defaultPnlTick));
    expect(responsePageAllPages.offset).toEqual(0);
    expect(responsePageAllPages.total).toEqual(2);
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

    const {
      maxBlockTime, count,
    }: {
      maxBlockTime: string,
      count: number,
    } = await PnlTicksTable.findLatestProcessedBlocktimeAndCount();

    expect(maxBlockTime).toEqual(blockTime);
    expect(count).toEqual(1);
  });

  it('Successfully finds latest block time without any pnl ticks', async () => {
    const {
      maxBlockTime, count,
    }: {
      maxBlockTime: string,
      count: number,
    } = await PnlTicksTable.findLatestProcessedBlocktimeAndCount();

    expect(maxBlockTime).toEqual(ZERO_TIME_ISO_8601);
    expect(count).toEqual(0);
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

    const leaderboardRankedData: {
      [accountId: string]: PnlTicksCreateObject,
    } = await PnlTicksTable.findMostRecentPnlTickForEachAccount(
      '3',
    );
    expect(leaderboardRankedData[defaultSubaccountId].equity).toEqual('1014');
    expect(leaderboardRankedData[defaultSubaccountId2].equity).toEqual('200');
  });

  const testCases = [
    {
      description: 'Get all time ranked pnl ticks',
      timeSpan: 'ALL_TIME',
      expectedLength: 2,
      expectedResults: [
        {
          address: defaultAddress,
          pnl: '1200',
          currentEquity: '1100',
          timeSpan: 'ALL_TIME',
          rank: '1',
        },
        {
          address: defaultAddress2,
          pnl: '300',
          currentEquity: '200',
          timeSpan: 'ALL_TIME',
          rank: '2',
        },
      ],
    },
    {
      description: 'Get one year ranked pnl ticks with missing pnl for one subaccount',
      timeSpan: 'ONE_YEAR',
      expectedLength: 2,
      expectedResults: [
        {
          address: defaultAddress2,
          pnl: '300',
          currentEquity: '200',
          timeSpan: 'ONE_YEAR',
          rank: '1',
        },
        {
          address: defaultAddress,
          pnl: '40',
          currentEquity: '1100',
          timeSpan: 'ONE_YEAR',
          rank: '2',
        },
      ],
    },
    {
      description: 'Get thirty days ranked pnl ticks',
      timeSpan: 'THIRTY_DAYS',
      expectedLength: 2,
      expectedResults: [
        {
          address: defaultAddress,
          pnl: '30',
          currentEquity: '1100',
          timeSpan: 'THIRTY_DAYS',
          rank: '1',
        },
        {
          address: defaultAddress2,
          pnl: '-30',
          currentEquity: '200',
          timeSpan: 'THIRTY_DAYS',
          rank: '2',
        },
      ],
    },
    {
      description: 'Get seven days ranked pnl ticks',
      timeSpan: 'SEVEN_DAYS',
      expectedLength: 2,
      expectedResults: [
        {
          address: defaultAddress,
          pnl: '20',
          currentEquity: '1100',
          timeSpan: 'SEVEN_DAYS',
          rank: '1',
        },
        {
          address: defaultAddress2,
          pnl: '-20',
          currentEquity: '200',
          timeSpan: 'SEVEN_DAYS',
          rank: '2',
        },
      ],
    },
    {
      description: 'Get one day ranked pnl ticks',
      timeSpan: 'ONE_DAY',
      expectedLength: 2,
      expectedResults: [
        {
          address: defaultAddress,
          pnl: '10',
          currentEquity: '1100',
          timeSpan: 'ONE_DAY',
          rank: '1',
        },
        {
          address: defaultAddress2,
          pnl: '-10',
          currentEquity: '200',
          timeSpan: 'ONE_DAY',
          rank: '2',
        },
      ],
    },
  ];

  it.each(testCases)('$description', async ({ timeSpan, expectedLength, expectedResults }) => {
    await setupRankedPnlTicksData();

    const leaderboardRankedData = await PnlTicksTable.getRankedPnlTicks(timeSpan);

    expect(leaderboardRankedData.length).toEqual(expectedLength);

    expectedResults.forEach((expectedResult, index) => {
      expect(leaderboardRankedData[index]).toEqual(expect.objectContaining(expectedResult));
    });
  });

  it('Ensure that vault addresses are not included in the leaderboard', async () => {
    await setupRankedPnlTicksData();

    await WalletTable.create(vaultWallet);
    await SubaccountTable.create(vaultSubaccount);
    await PnlTicksTable.create({
      subaccountId: vaultSubaccountId,
      equity: '100',
      createdAt: DateTime.utc().toISO(),
      totalPnl: '100',
      netTransfers: '50',
      blockHeight: '9',
      blockTime: defaultBlock.time,
    });

    const leaderboardRankedData: LeaderboardPnlCreateObject[] = await
    PnlTicksTable.getRankedPnlTicks(
      'ALL_TIME',
    );
    expect(leaderboardRankedData.length).toEqual(2);
  });

});

async function setupRankedPnlTicksData() {
  await Promise.all([
    BlockTable.create({
      blockHeight: '3',
      time: defaultBlock.time,
    }),
    BlockTable.create({
      blockHeight: '5',
      time: defaultBlock.time,
    }),
    BlockTable.create({
      blockHeight: '7',
      time: defaultBlock.time,
    }),
    BlockTable.create({
      blockHeight: '9',
      time: defaultBlock.time,
    }),
  ]);
  await PnlTicksTable.createMany([
    {
      subaccountId: defaultSubaccountId,
      equity: '1100',
      createdAt: DateTime.utc().toISO(),
      totalPnl: '1200',
      netTransfers: '50',
      blockHeight: '9',
      blockTime: defaultBlock.time,
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1090',
      createdAt: DateTime.utc().minus({ day: 1 }).toISO(),
      totalPnl: '1190',
      netTransfers: '50',
      blockHeight: '7',
      blockTime: defaultBlock.time,
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1080',
      createdAt: DateTime.utc().minus({ day: 7 }).toISO(),
      totalPnl: '1180',
      netTransfers: '50',
      blockHeight: '5',
      blockTime: defaultBlock.time,
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1070',
      createdAt: DateTime.utc().minus({ day: 30 }).toISO(),
      totalPnl: '1170',
      netTransfers: '50',
      blockHeight: '3',
      blockTime: defaultBlock.time,
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1060',
      createdAt: DateTime.utc().minus({ day: 365 }).toISO(),
      totalPnl: '1160',
      netTransfers: '50',
      blockHeight: '1',
      blockTime: defaultBlock.time,
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '200',
      createdAt: DateTime.utc().toISO(),
      totalPnl: '300',
      netTransfers: '50',
      blockHeight: '9',
      blockTime: defaultBlock.time,
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '210',
      createdAt: DateTime.utc().minus({ day: 1 }).toISO(),
      totalPnl: '310',
      netTransfers: '50',
      blockHeight: '7',
      blockTime: defaultBlock.time,
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '220',
      createdAt: DateTime.utc().minus({ week: 1 }).toISO(),
      totalPnl: '320',
      netTransfers: '50',
      blockHeight: '5',
      blockTime: defaultBlock.time,
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '230',
      createdAt: DateTime.utc().minus({ month: 1 }).toISO(),
      totalPnl: '330',
      netTransfers: '50',
      blockHeight: '3',
      blockTime: defaultBlock.time,
    },
  ]);
}
