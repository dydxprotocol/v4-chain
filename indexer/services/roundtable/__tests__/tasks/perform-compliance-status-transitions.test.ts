import {
  dbHelpers,
  testConstants,
  ComplianceStatusTable,
  ComplianceStatus,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import performComplianceStatusTransitionsTask from '../../src/tasks/perform-compliance-status-transitions';
import { logger, stats } from '@dydxprotocol-indexer/base';
import config from '../../src/config';
import { DateTime } from 'luxon';

describe('update-close-only-status', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'gauge');
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  it('succeeds with no CLOSE_ONLY addresses', async () => {
    await performComplianceStatusTransitionsTask();

    // Assert no addresses were updated
    expect(stats.gauge).toHaveBeenCalledWith(
      `${config.SERVICE_NAME}.num_stale_close_only_updated.count`,
      0,
    );
  });

  it('updates CLOSE_ONLY addresses older than 7 days to BLOCKED', async () => {
    config.CLOSE_ONLY_TO_BLOCKED_DAYS = 7;
    // Seed database with CLOSE_ONLY compliance status older than 7 days
    const oldUpdatedAt = DateTime.utc().minus({ days: 8 }).toISO();
    const newTs = DateTime.utc().toISO();
    await Promise.all([
      ComplianceStatusTable.create({
        address: testConstants.blockedAddress,
        status: ComplianceStatus.CLOSE_ONLY,
        createdAt: oldUpdatedAt,
        updatedAt: oldUpdatedAt,
      }),
      ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.CLOSE_ONLY,
        createdAt: newTs,
        updatedAt: newTs,
      }),
    ]);

    await performComplianceStatusTransitionsTask();

    // Assert the status was updated to BLOCKED
    const updatedStatus = await ComplianceStatusTable.findAll(
      { address: [testConstants.blockedAddress] },
      [],
      {},
    );
    expect(updatedStatus[0].status).toEqual(ComplianceStatus.BLOCKED);
    expect(updatedStatus[0].updatedAt).not.toEqual(oldUpdatedAt);
    const nonUpdatedStatus = await ComplianceStatusTable.findAll(
      { address: [testConstants.defaultAddress] },
      [],
      {},
    );
    expect(nonUpdatedStatus[0]).toEqual(expect.objectContaining({
      address: testConstants.defaultAddress,
      status: ComplianceStatus.CLOSE_ONLY,
      createdAt: newTs,
      updatedAt: newTs,
    }));

    // Assert the stats were correctly recorded
    expect(stats.gauge).toHaveBeenCalledWith(
      `${config.SERVICE_NAME}.num_stale_close_only_updated.count`,
      1,
    );
  });
});
