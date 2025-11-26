import { DateTime } from 'luxon';
import { PnlInterval, PnlFromDatabase } from '../../src/types';
import * as VaultPnlView from '../../src/stores/vault-pnl-view';
import * as PnlTable from '../../src/stores/pnl-table';
import * as VaultTable from '../../src/stores/vault-table';
import * as WalletTable from '../../src/stores/wallet-table';
import * as SubaccountTable from '../../src/stores/subaccount-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import {
  defaultSubaccountId,
  defaultSubaccountIdWithAlternateAddress,
  defaultSubaccountWithAlternateAddress,
  defaultWallet2,
  defaultVault,
  defaultSubaccount,
} from '../helpers/constants';

describe('VaultPnl store', () => {
  beforeAll(async () => {
    await migrate();
  });

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

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it.each([
    {
      description: 'Get hourly vault pnl',
      interval: PnlInterval.hour,
    },
    {
      description: 'Get daily vault pnl',
      interval: PnlInterval.day,
    },
  ])('$description', async ({
    interval,
  }: {
    interval: PnlInterval,
  }) => {
    const createdPnl: PnlFromDatabase[] = await setupIntervalPnl();

    await VaultPnlView.refreshDailyView();
    await VaultPnlView.refreshHourlyView();

    const pnlData: PnlFromDatabase[] = await VaultPnlView.getVaultsPnl(
      interval,
      7 * 24 * 60 * 60, // 1 week
      DateTime.fromISO(createdPnl[8].createdAt).plus({ seconds: 1 }),
    );

    // See setup function for created pnl records.
    // Should exclude records that are within the same hour except the first.
    const expectedHourlyPnl: PnlFromDatabase[] = [
      createdPnl[7],
      createdPnl[5],
      createdPnl[2],
      createdPnl[0],
    ];

    // Should exclude records that are within the same day except for the first.
    const expectedDailyPnl: PnlFromDatabase[] = [
      createdPnl[7],
      createdPnl[2],
    ];

    if (interval === PnlInterval.day) {
      expect(pnlData).toEqual(expectedDailyPnl);
    } else if (interval === PnlInterval.hour) {
      expect(pnlData).toEqual(expectedHourlyPnl);
    }
  });

  it('Get latest vault pnl', async () => {
    await setupIntervalPnl();

    await VaultPnlView.refreshHourlyView();

    const latestPnl: PnlFromDatabase[] = await VaultPnlView.getLatestVaultPnl();

    // Should return the most recent PNL for each subaccount
    expect(latestPnl).toHaveLength(2);

    const subaccount1Pnl = latestPnl.find(
      (pnl) => pnl.subaccountId === defaultSubaccountId,
    );
    const subaccount2Pnl = latestPnl.find(
      (pnl) => pnl.subaccountId === defaultSubaccountIdWithAlternateAddress,
    );

    expect(subaccount1Pnl).toBeDefined();
    expect(subaccount1Pnl?.equity).toBe('1100');

    expect(subaccount2Pnl).toBeDefined();
    expect(subaccount2Pnl?.equity).toBe('200');
  });

  async function setupIntervalPnl(): Promise<PnlFromDatabase[]> {
    const currentTime: DateTime = DateTime.utc().startOf('day');
    const tenMinAgo: string = currentTime.minus({ minute: 10 }).toISO();
    const almostTenMinAgo: string = currentTime.minus({ second: 603 }).toISO();
    const twoHoursAgo: string = currentTime.minus({ hour: 2 }).toISO();
    const twoDaysAgo: string = currentTime.minus({ day: 2 }).toISO();
    const monthAgo: string = currentTime.minus({ day: 30 }).toISO();

    const createdPnl: PnlFromDatabase[] = await PnlTable.createMany([
      {
        subaccountId: defaultSubaccountId,
        equity: '1100',
        createdAt: almostTenMinAgo,
        createdAtHeight: '10',
        totalPnl: '1200',
        netTransfers: '50',
      },
      {
        subaccountId: defaultSubaccountId,
        equity: '1090',
        createdAt: tenMinAgo,
        createdAtHeight: '8',
        totalPnl: '1190',
        netTransfers: '50',
      },
      {
        subaccountId: defaultSubaccountId,
        equity: '1080',
        createdAt: twoHoursAgo,
        createdAtHeight: '6',
        totalPnl: '1180',
        netTransfers: '50',
      },
      {
        subaccountId: defaultSubaccountId,
        equity: '1070',
        createdAt: twoDaysAgo,
        createdAtHeight: '4',
        totalPnl: '1170',
        netTransfers: '50',
      },
      {
        subaccountId: defaultSubaccountId,
        equity: '1200',
        createdAt: monthAgo,
        createdAtHeight: '3',
        totalPnl: '1170',
        netTransfers: '50',
      },
      {
        subaccountId: defaultSubaccountIdWithAlternateAddress,
        equity: '200',
        createdAt: almostTenMinAgo,
        createdAtHeight: '10',
        totalPnl: '300',
        netTransfers: '50',
      },
      {
        subaccountId: defaultSubaccountIdWithAlternateAddress,
        equity: '210',
        createdAt: tenMinAgo,
        createdAtHeight: '8',
        totalPnl: '310',
        netTransfers: '50',
      },
      {
        subaccountId: defaultSubaccountIdWithAlternateAddress,
        equity: '220',
        createdAt: twoHoursAgo,
        createdAtHeight: '6',
        totalPnl: '320',
        netTransfers: '50',
      },
      {
        subaccountId: defaultSubaccountIdWithAlternateAddress,
        equity: '230',
        createdAt: twoDaysAgo,
        createdAtHeight: '4',
        totalPnl: '330',
        netTransfers: '50',
      },
    ]);

    return createdPnl;
  }
});
