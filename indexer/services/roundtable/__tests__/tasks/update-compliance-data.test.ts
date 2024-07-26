import {
  dbHelpers,
  testConstants,
  testMocks,
  ComplianceTable,
  ComplianceProvider,
  ComplianceDataFromDatabase,
  SubaccountTable,
  SubaccountFromDatabase,
  ComplianceDataColumns,
  Ordering,
  ComplianceDataCreateObject,
} from '@dydxprotocol-indexer/postgres';
import updateComplianceDataTask from '../../src/tasks/update-compliance-data';
import { logger, stats } from '@dydxprotocol-indexer/base';
import _ from 'lodash';
import config from '../../src/config';
import { ClientAndProvider } from '../../src/helpers/compliance-clients';
import { ComplianceClientResponse } from '@dydxprotocol-indexer/compliance';
import { DateTime } from 'luxon';

interface ComplianceClientResponseWithNull extends Omit<ComplianceClientResponse, 'riskScore'> {
  riskScore: string | undefined | null;
}

describe('update-compliance-data', () => {
  let mockProvider: ClientAndProvider;

  const defaultMaxQueries: number = config.MAX_COMPLIANCE_DATA_QUERY_PER_LOOP;

  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    mockProvider = {
      provider: ComplianceProvider.ELLIPTIC,
      client: {
        getComplianceResponse: jest.fn(),
      },
    };
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'gauge');
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.resetAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  it('succeeds with no addresses', async () => {
    // Clear all mock data
    await dbHelpers.clearData();

    await updateComplianceDataTask(mockProvider);

    expectGaugeStats({
      activeAddresses: 0,
      newAddresses: 0,
      oldAddresses: 0,
      addressesScreened: 0,
      upserted: 0,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
  });

  it('succeeds with no active addresses or old addresses', async () => {
    // Seed database with new compliance data and set subaccounts to not active
    await Promise.all([
      ComplianceTable.create({ ...testConstants.nonBlockedComplianceData }),
      setupInitialSubaccounts(config.ACTIVE_ADDRESS_THRESHOLD_SECONDS * 2),
    ]);

    await updateComplianceDataTask(mockProvider);

    expectGaugeStats({
      activeAddresses: 0,
      newAddresses: 0,
      oldAddresses: 0,
      addressesScreened: 0,
      upserted: 0,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
  });

  it('succeeds with blocked active addresses and no old addresses', async () => {
    // Seed database with blocked compliance data older than age threshold and set subaccounts to
    // active
    await setupComplianceData(
      config.MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS * 2,
      {
        ...testConstants.nonBlockedComplianceData,
        blocked: true,
      },
    );

    await updateComplianceDataTask(mockProvider);

    expectGaugeStats({
      // no active addresses as the active address is blocked
      activeAddresses: 0,
      newAddresses: 0,
      oldAddresses: 0,
      addressesScreened: 0,
      upserted: 0,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
  });

  it('succeeds with no active addresses and blocked old addresses', async () => {
    // Seed database with blocked compliance data older than max age threshold and set subaccounts
    // non-active
    await Promise.all([
      setupComplianceData(
        config.MAX_COMPLIANCE_DATA_AGE_SECONDS * 2,
        {
          ...testConstants.nonBlockedComplianceData,
          blocked: true,
        },
      ),
      setupInitialSubaccounts(config.ACTIVE_ADDRESS_THRESHOLD_SECONDS * 2),
    ]);

    await updateComplianceDataTask(mockProvider);

    expectGaugeStats({
      activeAddresses: 0,
      newAddresses: 0,
      // blocked addresses are filtered out when querying for old addresses
      oldAddresses: 0,
      addressesScreened: 0,
      upserted: 0,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
  });

  it('succeeds with a new address to upsert', async () => {
    const riskScore: string = '45.00';
    setupMockProvider(
      mockProvider,
      { [testConstants.defaultAddress]: { blocked: true, riskScore } },
    );

    let complianceData: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
    expect(complianceData).toHaveLength(0);

    await updateComplianceDataTask(mockProvider);

    complianceData = await ComplianceTable.findAll({}, [], {});
    expect(complianceData).toHaveLength(1);
    expect(complianceData[0]).toEqual(expect.objectContaining({
      address: testConstants.defaultAddress,
      provider: mockProvider.provider,
      blocked: true,
      riskScore,
    }));

    expectGaugeStats({
      activeAddresses: 0,
      newAddresses: 1,
      oldAddresses: 0,
      addressesScreened: 1,
      upserted: 1,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
  });

  it('succeeds with an active address to update', async () => {
    // Seed database with compliance data older than the age threshold for active addresses
    await setupComplianceData(config.MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS * 2);

    const riskScore: string = '85.00';
    setupMockProvider(
      mockProvider,
      { [testConstants.defaultAddress]: { blocked: true, riskScore } },
    );

    await updateComplianceDataTask(mockProvider);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
    expect(complianceData).toHaveLength(1);
    expectUpdatedCompliance(
      complianceData[0],
      {
        address: testConstants.defaultAddress,
        blocked: true,
        riskScore,
      },
      mockProvider.provider,
    );

    expectGaugeStats({
      activeAddresses: 1,
      newAddresses: 0,
      oldAddresses: 0,
      addressesScreened: 1,
      upserted: 1,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
  });

  it('succeeds with an active address to update, undefined risk-score', async () => {
    // Seed database with compliance data older than the age threshold for active addresses
    await setupComplianceData(config.MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS * 2);

    setupMockProvider(
      mockProvider,
      { [testConstants.defaultAddress]: { blocked: true, riskScore: undefined } },
    );

    await updateComplianceDataTask(mockProvider);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
    expect(complianceData).toHaveLength(1);
    expectUpdatedCompliance(
      complianceData[0],
      {
        address: testConstants.defaultAddress,
        blocked: true,
        riskScore: null,
      },
      mockProvider.provider,
    );

    expectGaugeStats({
      activeAddresses: 1,
      newAddresses: 0,
      oldAddresses: 0,
      addressesScreened: 1,
      upserted: 1,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
  });

  it('succeeds with an old address to update', async () => {
    // Seed database with old compliance data, and set up subaccounts to not be active
    await Promise.all([
      setupComplianceData(config.MAX_COMPLIANCE_DATA_AGE_SECONDS * 2),
      setupInitialSubaccounts(config.MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS * 2),
    ]);

    const riskScore: string = '75.00';
    setupMockProvider(
      mockProvider,
      { [testConstants.defaultAddress]: { blocked: true, riskScore } },
    );

    await updateComplianceDataTask(mockProvider);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
    expect(complianceData).toHaveLength(1);
    expectUpdatedCompliance(
      complianceData[0],
      {
        address: testConstants.defaultAddress,
        blocked: true,
        riskScore,
      },
      mockProvider.provider,
    );

    expectGaugeStats({
      activeAddresses: 0,
      newAddresses: 0,
      oldAddresses: 1,
      addressesScreened: 1,
      upserted: 1,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
  });

  it('succeeds with updating multiple addresses', async () => {
    // Seed database with old compliance data and set up accounts to not be active
    await Promise.all([
      setupComplianceData(config.MAX_COMPLIANCE_DATA_AGE_SECONDS * 2),
      setupInitialSubaccounts(config.MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS * 2),
    ]);
    // Create a new active subaccount
    await SubaccountTable.create({
      address: testConstants.blockedAddress,
      subaccountNumber: 0,
      updatedAtHeight: '1',
      updatedAt: DateTime.utc().toISO(),
      assetYieldIndex: testConstants.defaultSubaccount.assetYieldIndex,
    });

    const riskScores: string[] = ['75.00', '50.00'];
    setupMockProvider(
      mockProvider,
      {
        [testConstants.defaultAddress]: { blocked: true, riskScore: riskScores[0] },
        [testConstants.blockedAddress]: { blocked: true, riskScore: riskScores[1] },
      },
    );

    await updateComplianceDataTask(mockProvider);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {
      orderBy: [[ComplianceDataColumns.address, Ordering.DESC]],
    });
    expect(complianceData).toHaveLength(2);
    expectUpdatedCompliance(
      complianceData[0],
      {
        address: testConstants.defaultAddress,
        blocked: true,
        riskScore: riskScores[0],
      },
      mockProvider.provider,
    );
    expectUpdatedCompliance(
      complianceData[1],
      {
        address: testConstants.blockedAddress,
        blocked: true,
        riskScore: riskScores[1],
      },
      mockProvider.provider,
    );
    // Both addresses screened
    expectGaugeStats({
      activeAddresses: 0,
      newAddresses: 1,
      oldAddresses: 1,
      addressesScreened: 2,
      upserted: 2,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
  });

  it('succeeds with updating multiple addresses, with failures', async () => {
    // Seed database with old compliance data and set up accounts to not be active
    await Promise.all([
      setupComplianceData(config.MAX_COMPLIANCE_DATA_AGE_SECONDS * 2),
      setupInitialSubaccounts(config.MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS * 2),
    ]);
    // Create a new active subaccount
    await SubaccountTable.create({
      address: testConstants.blockedAddress,
      subaccountNumber: 0,
      updatedAtHeight: '1',
      updatedAt: DateTime.utc().toISO(),
      assetYieldIndex: testConstants.defaultSubaccount.assetYieldIndex,
    });

    const addressWithComplianceError: string = 'dydx1gem4xs643fjhaqvphrvv0adpg4435j7xx9pp4z';
    // Create a new active subaccount that will return an error when queried
    await SubaccountTable.create({
      address: addressWithComplianceError,
      subaccountNumber: 0,
      updatedAtHeight: '1',
      updatedAt: DateTime.utc().toISO(),
      assetYieldIndex: testConstants.defaultSubaccount.assetYieldIndex,
    });

    const riskScores: string[] = ['75.00', '50.00'];
    setupMockProvider(
      mockProvider,
      {
        [testConstants.defaultAddress]: { blocked: true, riskScore: riskScores[0] },
        [testConstants.blockedAddress]: { blocked: true, riskScore: riskScores[1] },
      },
    );

    await updateComplianceDataTask(mockProvider);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {
      orderBy: [[ComplianceDataColumns.address, Ordering.DESC]],
    });
    expect(complianceData).toHaveLength(2);
    expectUpdatedCompliance(
      complianceData[0],
      {
        address: testConstants.defaultAddress,
        blocked: true,
        riskScore: riskScores[0],
      },
      mockProvider.provider,
    );
    expectUpdatedCompliance(
      complianceData[1],
      {
        address: testConstants.blockedAddress,
        blocked: true,
        riskScore: riskScores[1],
      },
      mockProvider.provider,
    );
    // Both addresses screened
    expectGaugeStats({
      activeAddresses: 0,
      newAddresses: 2,
      oldAddresses: 1,
      addressesScreened: 3,
      upserted: 2,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);
    // error log
    expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
      addresses: [addressWithComplianceError],
      errors: [{
        reason: new Error(`Unexpected address ${addressWithComplianceError} passed to provider`),
        status: 'rejected',
      }],
    }));
    // increments stat for failures to get compliance data
    expect(stats.increment).toHaveBeenCalledWith(
      'roundtable.update_compliance_data.get_compliance_data_fail',
      1,
      undefined,
      { provider: mockProvider.provider },
    );
  });

  it('limits number of addresses scanned per run', async () => {
    // Set the limit of addresses to scan to 1
    config.MAX_COMPLIANCE_DATA_QUERY_PER_LOOP = 1;

    // Seed database with compliance data older than the age threshold for active subaccounts
    await Promise.all([
      setupComplianceData(config.MAX_COMPLIANCE_DATA_AGE_SECONDS * 2),
      setupInitialSubaccounts(config.MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS * 2),
    ]);
    // Create a new active subaccount
    await SubaccountTable.create({
      address: testConstants.blockedAddress,
      subaccountNumber: 0,
      updatedAtHeight: '1',
      updatedAt: DateTime.utc().toISO(),
      assetYieldIndex: testConstants.defaultSubaccount.assetYieldIndex,
    });

    const riskScores: string[] = ['75.00', '50.00'];
    setupMockProvider(
      mockProvider,
      {
        [testConstants.defaultAddress]: { blocked: true, riskScore: riskScores[0] },
        [testConstants.blockedAddress]: { blocked: true, riskScore: riskScores[1] },
      },
    );

    // First run should query a new compliance data object for the new active account
    await updateComplianceDataTask(mockProvider);
    let complianceData: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
    expect(complianceData).toHaveLength(2);
    // Should be the second compliance data row as it's newer
    expectUpdatedCompliance(
      complianceData[1],
      {
        address: testConstants.blockedAddress,
        blocked: true,
        riskScore: riskScores[1],
      },
      mockProvider.provider,
    );
    // Only a single address screened
    expectGaugeStats({
      activeAddresses: 0,
      newAddresses: 1,
      oldAddresses: 0,
      addressesScreened: 1,
      upserted: 1,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);

    // Second run should query an updated compliance data object for the address with old
    // compliance data
    await updateComplianceDataTask(mockProvider);
    complianceData = await ComplianceTable.findAll({}, [], {});
    expect(complianceData).toHaveLength(2);
    // Should be the second compliance data row as it's newer
    expectUpdatedCompliance(
      complianceData[1],
      {
        address: testConstants.defaultAddress,
        blocked: true,
        riskScore: riskScores[0],
      },
      mockProvider.provider,
    );

    // Only a single address screened
    expectGaugeStats({
      activeAddresses: 0,
      newAddresses: 0,
      oldAddresses: 1,
      addressesScreened: 1,
      upserted: 1,
    },
    mockProvider.provider,
    );
    expectTimingStats(mockProvider.provider);

    config.MAX_COMPLIANCE_DATA_QUERY_PER_LOOP = defaultMaxQueries;
  });
});

async function setupComplianceData(
  deltaSeconds: number,
  complianceCreate: ComplianceDataCreateObject = testConstants.nonBlockedComplianceData,
): Promise<void> {
  const oldUpdatedAt: string = DateTime.utc().minus(
    { seconds: deltaSeconds },
  ).toUTC().toISO();
  await ComplianceTable.create({
    ...complianceCreate,
    updatedAt: oldUpdatedAt,
  });
  const complianceData: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
  expect(complianceData).toHaveLength(1);
  expect(complianceData[0]).toEqual(expect.objectContaining({
    ...complianceCreate,
    updatedAt: oldUpdatedAt,
  }));
}

async function setupInitialSubaccounts(
  deltaSeconds: number,
): Promise<void> {
  const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll({
    address: testConstants.defaultAddress,
    subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
  }, [], {});
  await setupSubaccounts(deltaSeconds, _.map(subaccounts, 'id'));
}

async function setupSubaccounts(
  deltaSeconds: number,
  subaccountIds: string[],
): Promise<void> {
  const newUpdatedAt: string = DateTime.utc().minus({
    seconds: deltaSeconds,
  }).toUTC().toISO();
  await Promise.all(subaccountIds.map(
    (subaccountId: string) => {
      return SubaccountTable.update(
        { id: subaccountId, 
          updatedAtHeight: '1', 
          updatedAt: newUpdatedAt, 
          assetYieldIndex: testConstants.defaultSubaccount.assetYieldIndex, 
        },
      );
    },
  ));
}

function setupMockProvider(
  clientAndProvider: ClientAndProvider,
  expectedResponses: {[address: string]: {blocked: boolean, riskScore: string | undefined }},
): void {
  // eslint-disable-next-line no-param-reassign
  clientAndProvider.client.getComplianceResponse = jest.fn().mockImplementation(
    // eslint-disable-next-line @typescript-eslint/require-await
    async (address: string): Promise<ComplianceClientResponse> => {
      if (expectedResponses[address] === undefined) {
        throw new Error(`Unexpected address ${address} passed to provider`);
      } else {
        return {
          address,
          ...expectedResponses[address],
        };
      }
    },
  );
}

function expectUpdatedCompliance(
  complianceData: ComplianceDataFromDatabase,
  complianceClientResponse: ComplianceClientResponseWithNull,
  provider: string,
): void {
  expect(complianceData).toEqual(expect.objectContaining({
    ...complianceClientResponse,
    provider,
  }));
  // Updated at should be updated to a time within the last day
  expect(DateTime.fromISO(complianceData.updatedAt).toUnixInteger()).toBeGreaterThan(
    DateTime.utc().minus({ days: 1 }).toUnixInteger());
}

function expectGaugeStats(
  {
    activeAddresses,
    newAddresses,
    oldAddresses,
    addressesScreened,
    upserted,
  }: {
    activeAddresses: number,
    newAddresses: number,
    oldAddresses: number,
    addressesScreened: number,
    upserted: number,
  },
  provider: string,
): void {
  expect(stats.gauge).toHaveBeenCalledWith(
    'roundtable.update_compliance_data.num_active_addresses',
    activeAddresses,
    undefined,
    { provider },
  );
  expect(stats.gauge).toHaveBeenCalledWith(
    'roundtable.update_compliance_data.num_new_addresses',
    newAddresses,
    undefined,
    { provider },
  );
  expect(stats.gauge).toHaveBeenCalledWith(
    'roundtable.update_compliance_data.num_old_addresses',
    oldAddresses,
    undefined,
    { provider },
  );
  expect(stats.gauge).toHaveBeenCalledWith(
    'roundtable.update_compliance_data.num_addresses_to_screen',
    addressesScreened,
    undefined,
    { provider },
  );
  expect(stats.gauge).toHaveBeenCalledWith(
    'roundtable.update_compliance_data.num_upserted',
    upserted,
    undefined,
    { provider },
  );
}

function expectTimingStats(
  provider: string,
): void {
  expect(stats.timing).toHaveBeenCalledWith(
    'roundtable.update_compliance_data.get_active_addresses',
    expect.any(Number),
    undefined,
    { provider },
  );
  expect(stats.timing).toHaveBeenCalledWith(
    'roundtable.update_compliance_data.get_old_addresses',
    expect.any(Number),
    undefined,
    { provider },
  );
  expect(stats.timing).toHaveBeenCalledWith(
    'roundtable.update_compliance_data.query_compliance_data',
    expect.any(Number),
    undefined,
    { provider },
  );
}
