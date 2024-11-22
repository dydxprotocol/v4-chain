import config from '../../src/config';
import refreshVaulPnlTask from '../../src/tasks/refresh-vault-pnl';
import { Settings, DateTime } from 'luxon';
import {
  BlockTable,
  PnlTickInterval,
  PnlTicksFromDatabase,
  PnlTicksTable,
  VaultPnlTicksView,
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
    await VaultPnlTicksView.refreshDailyView();
    await VaultPnlTicksView.refreshHourlyView();
    jest.clearAllMocks();
    Settings.now = () => new Date().valueOf();
  });

  it('refreshes hourly view if within time window of an hour', async () => {
    Settings.now = () => currentTime.startOf('hour').plus(
      { milliseconds: config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS - 1 },
    ).valueOf();
    const pnlTick: PnlTicksFromDatabase = await setupPnlTick();
    await refreshVaulPnlTask();

    const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
      PnlTickInterval.hour,
      86400,
      currentTime.minus({ day: 1 }),
    );
    expect(pnlTicks).toEqual([pnlTick]);
  });

  it('refreshes daily view if within time window of a day', async () => {
    Settings.now = () => currentTime.startOf('day').plus(
      { milliseconds: config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS - 1 },
    ).valueOf();
    const pnlTick: PnlTicksFromDatabase = await setupPnlTick();
    await refreshVaulPnlTask();

    const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
      PnlTickInterval.day,
      608400,
      currentTime.minus({ day: 7 }),
    );
    expect(pnlTicks).toEqual([pnlTick]);
  });

  it('does not refresh hourly view if outside of time window of an hour', async () => {
    Settings.now = () => currentTime.startOf('hour').plus(
      { milliseconds: config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS + 1 },
    ).valueOf();
    await setupPnlTick();
    await refreshVaulPnlTask();

    const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
      PnlTickInterval.hour,
      86400,
      currentTime.minus({ day: 1 }),
    );
    expect(pnlTicks).toEqual([]);
  });

  it('does not refresh daily view if outside time window of a day', async () => {
    Settings.now = () => currentTime.startOf('day').plus(
      { milliseconds: config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS + 1 },
    ).valueOf();
    await setupPnlTick();
    await refreshVaulPnlTask();

    const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
      PnlTickInterval.day,
      608400,
      currentTime.minus({ day: 7 }),
    );
    expect(pnlTicks).toEqual([]);
  });

  async function setupPnlTick(): Promise<PnlTicksFromDatabase> {
    const twoHoursAgo: string = currentTime.minus({ hour: 2 }).toISO();
    await Promise.all([
      BlockTable.create({
        blockHeight: '6',
        time: twoHoursAgo,
      }),
    ]);
    const createdTick: PnlTicksFromDatabase = await PnlTicksTable.create(
      {
        subaccountId: testConstants.defaultSubaccountId,
        equity: '1080',
        createdAt: twoHoursAgo,
        totalPnl: '1180',
        netTransfers: '50',
        blockHeight: '6',
        blockTime: twoHoursAgo,
      },
    );
    return createdTick;
  }
});
