import {
  PnlTickInterval,
  PnlTicksFromDatabase,
} from '../../src/types';
import * as VaultPnlTicksView from '../../src/stores/vault-pnl-ticks-view';
import * as PnlTicksTable from '../../src/stores/pnl-ticks-table';
import * as BlockTable from '../../src/stores/block-table';
import * as VaultTable from '../../src/stores/vault-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import * as WalletTable from '../../src/stores/wallet-table';
import * as SubaccountTable from '../../src/stores/subaccount-table';
import {
  defaultSubaccountId,
  defaultSubaccountIdWithAlternateAddress,
  defaultSubaccountWithAlternateAddress,
  defaultWallet2,
  defaultVault,
  defaultSubaccount,
} from '../helpers/constants';
import { DateTime } from 'luxon';

describe('PnlTicks store', () => {
  beforeEach(async () => {
    await seedData();
    await WalletTable.create(defaultWallet2);
    await SubaccountTable.create(defaultSubaccountWithAlternateAddress);
    await Promise.all([
      VaultTable.create({
        ...defaultVault,
        address: defaultSubaccount.address,
      }),
      VaultTable.create({
        ...defaultVault,
        address: defaultSubaccountWithAlternateAddress.address,
      }),
    ]);
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
    await VaultPnlTicksView.refreshDailyView();
    await VaultPnlTicksView.refreshHourlyView();
    const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
      interval,
      7 * 24 * 60 * 60, // 1 week
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
});
