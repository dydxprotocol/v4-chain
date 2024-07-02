import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  BlockTable,
  SubaccountTable,
  TransferTable,
  testMocks,
  dbHelpers,
  testConstants,
} from '@dydxprotocol-indexer/postgres';
import runTask from '../../src/tasks/pnl-instrumentation';
import { getMostRecentPnlTicksForEachAccount } from '../../src/helpers/pnl-ticks-helper';
import { DateTime } from 'luxon';
import config from '../../src/config';
import { asMock } from '@dydxprotocol-indexer/dev';

jest.mock('../../src/helpers/pnl-ticks-helper');

describe('pnl-instrumentation', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
    jest.spyOn(stats, 'gauge');
    jest.spyOn(logger, 'error');
  });

  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  it('succeeds with no stale PNL subaccounts', async () => {
    jest.spyOn(BlockTable, 'getLatest').mockResolvedValue({
      blockHeight: '12345',
    } as any);

    jest.spyOn(SubaccountTable, 'getSubaccountsWithTransfers').mockResolvedValue([
      { id: 'subaccount1' },
      { id: 'subaccount2' },
    ] as any);

    asMock(getMostRecentPnlTicksForEachAccount).mockImplementation(
      async () => Promise.resolve({
        subaccount1: {
          ...testConstants.defaultPnlTick,
          blockTime: DateTime.utc().minus({ hours: 1 }).toISO(),
        },
        subaccount2: {
          ...testConstants.defaultPnlTick,
          blockTime: DateTime.utc().minus({ hours: 1 }).toISO(),
        },
      }),
    );

    await runTask();

    expect(stats.gauge).toHaveBeenCalledWith(`${config.SERVICE_NAME}.pnl_stale_subaccounts`, 0);
    expect(stats.gauge).toHaveBeenCalledWith(`${config.SERVICE_NAME}.pnl_stale_subaccounts_with_prior_pnl`, 0);
    expect(stats.gauge).toHaveBeenCalledWith(`${config.SERVICE_NAME}.pnl_stale_subaccounts_without_prior_pnl`, 0);
    expect(logger.error).toHaveBeenCalledTimes(0);
  });

  it('succeeds with stale PNL subaccounts', async () => {
    jest.spyOn(BlockTable, 'getLatest').mockResolvedValue({
      blockHeight: '12345',
    } as any);

    jest.spyOn(SubaccountTable, 'getSubaccountsWithTransfers').mockResolvedValue([
      { id: 'subaccount1' },
      { id: 'subaccount2' },
    ] as any);

    asMock(getMostRecentPnlTicksForEachAccount).mockImplementation(
      async () => Promise.resolve({
        subaccount1: {
          ...testConstants.defaultPnlTick,
          blockTime: DateTime.utc().minus({ hours: 3 }).toISO(),
        },
        subaccount2: {
          ...testConstants.defaultPnlTick,
          blockTime: DateTime.utc().minus({ hours: 3 }).toISO(),
        },
      }),
    );

    await runTask();

    expect(stats.gauge).toHaveBeenCalledWith(`${config.SERVICE_NAME}.pnl_stale_subaccounts`, 2);
    expect(stats.gauge).toHaveBeenCalledWith(`${config.SERVICE_NAME}.pnl_stale_subaccounts_with_prior_pnl`, 2);
    expect(stats.gauge).toHaveBeenCalledWith(`${config.SERVICE_NAME}.pnl_stale_subaccounts_without_prior_pnl`, 0);
    expect(logger.error).toHaveBeenCalledWith({
      at: 'pnl-instrumentation#statPnl',
      message: 'Subaccount ids with stale PNL data',
      stalePnlSubaccounts: ['subaccount1', 'subaccount2'],
      staleTransferSubaccounts: [],
    });
  });

  it('succeeds with stale transfer subaccounts', async () => {
    jest.spyOn(BlockTable, 'getLatest').mockResolvedValue({
      blockHeight: '12345',
    } as any);

    jest.spyOn(SubaccountTable, 'getSubaccountsWithTransfers').mockResolvedValue([
      { id: 'subaccount1' },
      { id: 'subaccount2' },
    ] as any);

    asMock(getMostRecentPnlTicksForEachAccount).mockImplementation(
      async () => Promise.resolve({
        subaccount1: {
          ...testConstants.defaultPnlTick,
          blockTime: DateTime.utc().minus({ hours: 1 }).toISO(),
        },
      }),
    );

    jest.spyOn(TransferTable, 'getLastTransferTimeForSubaccounts').mockResolvedValue({
      subaccount2: DateTime.utc().minus({ hours: 3 }).toISO(),
    });

    await runTask();

    expect(stats.gauge).toHaveBeenCalledWith(`${config.SERVICE_NAME}.pnl_stale_subaccounts`, 1);
    expect(stats.gauge).toHaveBeenCalledWith(`${config.SERVICE_NAME}.pnl_stale_subaccounts_with_prior_pnl`, 0);
    expect(stats.gauge).toHaveBeenCalledWith(`${config.SERVICE_NAME}.pnl_stale_subaccounts_without_prior_pnl`, 1);
    expect(logger.error).toHaveBeenCalledWith({
      at: 'pnl-instrumentation#statPnl',
      message: 'Subaccount ids with stale PNL data',
      stalePnlSubaccounts: [],
      staleTransferSubaccounts: ['subaccount2'],
    });
  });
});
