import config from '../../src/config';
import refreshVaultPnlTask from '../../src/tasks/refresh-vault-pnl';
import { Settings, DateTime } from 'luxon';
import {
  BlockTable,
  PnlFromDatabase,
  PnlInterval,
  PnlTable,
  PnlTickInterval,
  PnlTicksFromDatabase,
  PnlTicksTable,
  VaultPnlTicksView,
  VaultPnlView,
  VaultTable,
  dbHelpers,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';

jest.mock('../../src/helpers/aws');

describe('refresh-vault-pnl', () => {
  const currentTime: DateTime = DateTime.utc();

  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await Promise.all([
      VaultTable.create({
        ...testConstants.defaultVault,
        address: testConstants.defaultSubaccount.address,
      }),
    ]);
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    // Refresh both old and new views
    await Promise.all([
      VaultPnlTicksView.refreshDailyView(),
      VaultPnlTicksView.refreshHourlyView(),
      VaultPnlView.refreshDailyView(),
      VaultPnlView.refreshHourlyView(),
    ]);
    jest.clearAllMocks();
    Settings.now = () => new Date().valueOf();
  });

  it('refreshes hourly views (old and new) if within time window of an hour', async () => {
    Settings.now = () => currentTime.startOf('hour').plus(
      { milliseconds: config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS - 1 },
    ).valueOf();

    await setupPnlData();

    await refreshVaultPnlTask();

    // Verify old view was refreshed
    const oldPnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
      PnlTickInterval.hour,
      86400,
      currentTime.minus({ day: 1 }),
    );

    // Verify new view was refreshed
    const newPnlData: PnlFromDatabase[] = await VaultPnlView.getVaultsPnl(
      PnlInterval.hour,
      86400,
      currentTime.minus({ day: 1 }),
    );

    // Both views should have the same data
    expect(oldPnlTicks).toHaveLength(1);
    expect(newPnlData).toHaveLength(1);

    // Verify data matches
    expect(oldPnlTicks[0].subaccountId).toEqual(newPnlData[0].subaccountId);
    expect(oldPnlTicks[0].equity).toEqual(newPnlData[0].equity);
    expect(oldPnlTicks[0].totalPnl).toEqual(newPnlData[0].totalPnl);
    expect(oldPnlTicks[0].netTransfers).toEqual(newPnlData[0].netTransfers);
  });

  it('refreshes daily views (old and new) if within time window of a day', async () => {
    Settings.now = () => currentTime.startOf('day').plus(
      { milliseconds: config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS - 1 },
    ).valueOf();

    await setupPnlData();

    await refreshVaultPnlTask();

    // Verify old view was refreshed
    const oldPnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
      PnlTickInterval.day,
      608400,
      currentTime.minus({ day: 7 }),
    );

    // Verify new view was refreshed
    const newPnlData: PnlFromDatabase[] = await VaultPnlView.getVaultsPnl(
      PnlInterval.day,
      608400,
      currentTime.minus({ day: 7 }),
    );

    // Both views should have the same data
    expect(oldPnlTicks).toHaveLength(1);
    expect(newPnlData).toHaveLength(1);

    expect(oldPnlTicks[0].subaccountId).toEqual(newPnlData[0].subaccountId);
    expect(newPnlData[0].subaccountId).toEqual(testConstants.defaultSubaccountId);
  });

  it('does not refresh hourly views if outside of time window of an hour', async () => {
    Settings.now = () => currentTime.startOf('hour').plus(
      { milliseconds: config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS + 1 },
    ).valueOf();

    await setupPnlData();

    await refreshVaultPnlTask();

    // Neither view should have been refreshed
    const oldPnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
      PnlTickInterval.hour,
      86400,
      currentTime.minus({ day: 1 }),
    );

    const newPnlData: PnlFromDatabase[] = await VaultPnlView.getVaultsPnl(
      PnlInterval.hour,
      86400,
      currentTime.minus({ day: 1 }),
    );

    expect(oldPnlTicks).toEqual([]);
    expect(newPnlData).toEqual([]);
  });

  it('does not refresh daily views if outside time window of a day', async () => {
    Settings.now = () => currentTime.startOf('day').plus(
      { milliseconds: config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS + 1 },
    ).valueOf();

    await setupPnlData();

    await refreshVaultPnlTask();

    // Neither view should have been refreshed
    const oldPnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
      PnlTickInterval.day,
      608400,
      currentTime.minus({ day: 7 }),
    );

    const newPnlData: PnlFromDatabase[] = await VaultPnlView.getVaultsPnl(
      PnlInterval.day,
      608400,
      currentTime.minus({ day: 7 }),
    );

    expect(oldPnlTicks).toEqual([]);
    expect(newPnlData).toEqual([]);
  });

  /**
   * Setup PNL data in both old (pnl_ticks) and new (pnl) tables
   * to test that both views refresh correctly during migration period.
   */
  async function setupPnlData(): Promise<void> {
    const twoHoursAgo: string = currentTime.minus({ hour: 2 }).toISO();

    // Create block for old pnl_ticks table (has foreign key to blocks)
    await BlockTable.create({
      blockHeight: '6',
      time: twoHoursAgo,
    });

    // Create data in old pnl_ticks table
    await PnlTicksTable.create({
      subaccountId: testConstants.defaultSubaccountId,
      equity: '1080',
      createdAt: twoHoursAgo,
      totalPnl: '1180',
      netTransfers: '50',
      blockHeight: '6',
      blockTime: twoHoursAgo,
    });

    // Create data in new pnl table
    await PnlTable.create({
      subaccountId: testConstants.defaultSubaccountId,
      equity: '1080',
      createdAt: twoHoursAgo,
      createdAtHeight: '6',
      totalPnl: '1180',
      netTransfers: '50',
    });
  }
});
