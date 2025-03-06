import {
  IsoString,
  LeaderboardPnlCreateObject,
  Ordering,
  PnlTickInterval,
  PnlTicksColumns,
  PnlTicksCreateObject,
  PnlTicksFromDatabase,
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

  it.each([
    {
      description: 'Get hourly pnl ticks',
      interval: PnlTickInterval.hour,
    },
    {
      description: 'Get daily pnl ticks',
      interval: PnlTickInterval.day,
    },
  ])('$description', async ({
    interval,
  }: {
    interval: PnlTickInterval,
  }) => {
    const createdTicks: PnlTicksFromDatabase[] = await setupIntervalPnlTicks();
    const pnlTicks: PnlTicksFromDatabase[] = await PnlTicksTable.getPnlTicksAtIntervals(
      interval,
      7 * 24 * 60 * 60, // 1 week
      [defaultSubaccountId, defaultSubaccountIdWithAlternateAddress],
      DateTime.fromISO(createdTicks[8].blockTime).plus({ seconds: 1 }),
    );
    // See setup function for created ticks.
    // Should exclude tick that is within the same hour except the first.
    const expectedHourlyTicks: PnlTicksFromDatabase[] = [
      createdTicks[7],
      createdTicks[5],
      createdTicks[2],
      createdTicks[0],
    ];
    // Should exclude ticks that is within the same day except for the first.
    const expectedDailyTicks: PnlTicksFromDatabase[] = [
      createdTicks[7],
      createdTicks[2],
    ];

    if (interval === PnlTickInterval.day) {
      expect(pnlTicks).toEqual(expectedDailyTicks);
    } else if (interval === PnlTickInterval.hour) {
      expect(pnlTicks).toEqual(expectedHourlyTicks);
    }
  });

  it('Gets latest pnl ticks for subaccounts before or at given date', async () => {
    const createdTicks: PnlTicksFromDatabase[] = await setupIntervalPnlTicks();
    const latestTicks: PnlTicksFromDatabase[] = await PnlTicksTable.getLatestPnlTick(
      [defaultSubaccountId, defaultSubaccountIdWithAlternateAddress],
      DateTime.fromISO(createdTicks[8].blockTime).plus({ seconds: 1 }),
    );
    expect(latestTicks).toEqual([createdTicks[8], createdTicks[3]]);
  });

  it('Gets empty pnl ticks for subaccounts before or at date earlier than all pnl data', async () => {
    const createdTicks: PnlTicksFromDatabase[] = await setupIntervalPnlTicks();
    const latestTicks: PnlTicksFromDatabase[] = await PnlTicksTable.getLatestPnlTick(
      [defaultSubaccountId, defaultSubaccountIdWithAlternateAddress],
      DateTime.fromISO(createdTicks[0].blockTime).minus({ years: 1 }),
    );
    expect(latestTicks).toEqual([]);
  });

  it('Gets empty pnl ticks for subaccounts before or at date if no subaccounts given', async () => {
    const createdTicks: PnlTicksFromDatabase[] = await setupIntervalPnlTicks();
    const latestTicks: PnlTicksFromDatabase[] = await PnlTicksTable.getLatestPnlTick(
      [],
      DateTime.fromISO(createdTicks[0].blockTime).plus({ years: 1 }),
    );
    expect(latestTicks).toEqual([]);
  });

});

async function setupRankedPnlTicksData() {
  const now = DateTime.utc().startOf('day');
  const thirtyDaysAgo = now.minus({ days: 30 });
  const sevenDaysAgo = now.minus({ days: 7 });
  const oneDayAgo = now.minus({ days: 1 });
  const oneYearAgo = now.minus({ years: 1 });

  await PnlTicksTable.createMany([
    {
      subaccountId: defaultSubaccountId,
      equity: '1100',
      totalPnl: '1200',
      netTransfers: '50',
      createdAt: now.toISO(),
      blockHeight: '9',
      blockTime: now.toISO(),
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1090',
      totalPnl: '1190',
      netTransfers: '50',
      createdAt: oneDayAgo.toISO(),
      blockHeight: '7',
      blockTime: oneDayAgo.toISO(),
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1080',
      totalPnl: '1180',
      netTransfers: '50',
      createdAt: sevenDaysAgo.toISO(),
      blockHeight: '5',
      blockTime: sevenDaysAgo.toISO(),
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1070',
      totalPnl: '1170',
      netTransfers: '50',
      createdAt: thirtyDaysAgo.toISO(),
      blockHeight: '3',
      blockTime: thirtyDaysAgo.toISO(),
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1060',
      totalPnl: '1160',
      netTransfers: '50',
      createdAt: oneYearAgo.toISO(),
      blockHeight: '1',
      blockTime: oneYearAgo.toISO(),
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '200',
      createdAt: now.toISO(),
      totalPnl: '300',
      netTransfers: '50',
      blockHeight: '9',
      blockTime: now.toISO(),
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '210',
      createdAt: oneDayAgo.toISO(),
      totalPnl: '310',
      netTransfers: '50',
      blockHeight: '7',
      blockTime: oneDayAgo.toISO(),
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '220',
      totalPnl: '320',
      netTransfers: '50',
      createdAt: sevenDaysAgo.toISO(),
      blockHeight: '5',
      blockTime: sevenDaysAgo.toISO(),
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '230',
      totalPnl: '330',
      netTransfers: '50',
      createdAt: thirtyDaysAgo.toISO(),
      blockHeight: '3',
      blockTime: thirtyDaysAgo.toISO(),
    },
  ]);
}

async function setupIntervalPnlTicks(): Promise<PnlTicksFromDatabase[]> {
  const currentTime: DateTime = DateTime.utc().startOf('day');
  const tenMinAgo: string = currentTime.minus({ minute: 10 }).toISO();
  const almostTenMinAgo: string = currentTime.minus({ second: 603 }).toISO();
  const twoHoursAgo: string = currentTime.minus({ hour: 2 }).toISO();
  const twoDaysAgo: string = currentTime.minus({ day: 2 }).toISO();
  const monthAgo: string = currentTime.minus({ day: 30 }).toISO();
  await Promise.all([
    BlockTable.create({
      blockHeight: '3',
      time: monthAgo,
    }),
    BlockTable.create({
      blockHeight: '4',
      time: twoDaysAgo,
    }),
    BlockTable.create({
      blockHeight: '6',
      time: twoHoursAgo,
    }),
    BlockTable.create({
      blockHeight: '8',
      time: almostTenMinAgo,
    }),
    BlockTable.create({
      blockHeight: '10',
      time: tenMinAgo,
    }),
  ]);
  const createdTicks: PnlTicksFromDatabase[] = await PnlTicksTable.createMany([
    {
      subaccountId: defaultSubaccountId,
      equity: '1100',
      createdAt: almostTenMinAgo,
      totalPnl: '1200',
      netTransfers: '50',
      blockHeight: '10',
      blockTime: almostTenMinAgo,
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1090',
      createdAt: tenMinAgo,
      totalPnl: '1190',
      netTransfers: '50',
      blockHeight: '8',
      blockTime: tenMinAgo,
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1080',
      createdAt: twoHoursAgo,
      totalPnl: '1180',
      netTransfers: '50',
      blockHeight: '6',
      blockTime: twoHoursAgo,
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1070',
      createdAt: twoDaysAgo,
      totalPnl: '1170',
      netTransfers: '50',
      blockHeight: '4',
      blockTime: twoDaysAgo,
    },
    {
      subaccountId: defaultSubaccountId,
      equity: '1200',
      createdAt: monthAgo,
      totalPnl: '1170',
      netTransfers: '50',
      blockHeight: '3',
      blockTime: monthAgo,
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '200',
      createdAt: almostTenMinAgo,
      totalPnl: '300',
      netTransfers: '50',
      blockHeight: '10',
      blockTime: almostTenMinAgo,
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '210',
      createdAt: tenMinAgo,
      totalPnl: '310',
      netTransfers: '50',
      blockHeight: '8',
      blockTime: tenMinAgo,
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '220',
      createdAt: twoHoursAgo,
      totalPnl: '320',
      netTransfers: '50',
      blockHeight: '6',
      blockTime: twoHoursAgo,
    },
    {
      subaccountId: defaultSubaccountIdWithAlternateAddress,
      equity: '230',
      createdAt: twoDaysAgo,
      totalPnl: '330',
      netTransfers: '50',
      blockHeight: '4',
      blockTime: twoDaysAgo,
    },
  ]);
  return createdTicks;
}
