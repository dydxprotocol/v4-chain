import { ComplianceDataFromDatabase, ComplianceProvider } from '../../src/types';
import * as ComplianceDataTable from '../../src/stores/compliance-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import {
  blockedComplianceData,
  blockedAddress,
  nonBlockedComplianceData,
} from '../helpers/constants';
import { DateTime } from 'luxon';

describe('Compliance data store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates compliance data', async () => {
    await Promise.all([
      ComplianceDataTable.create(blockedComplianceData),
      ComplianceDataTable.create(nonBlockedComplianceData),
    ]);
  });

  it('Successfully finds all compliance data', async () => {
    await Promise.all([
      ComplianceDataTable.create(blockedComplianceData),
      ComplianceDataTable.create(nonBlockedComplianceData),
    ]);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(complianceData.length).toEqual(2);
    expect(complianceData[0]).toEqual(blockedComplianceData);
    expect(complianceData[1]).toEqual(nonBlockedComplianceData);
  });

  it('Successfully finds compliance data with address', async () => {
    await Promise.all([
      ComplianceDataTable.create(blockedComplianceData),
      ComplianceDataTable.create(nonBlockedComplianceData),
    ]);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {
        address: [blockedAddress],
      },
      [],
      { readReplica: true },
    );

    expect(complianceData.length).toEqual(1);
    expect(complianceData[0]).toEqual(expect.objectContaining(blockedComplianceData));
  });

  it('Successfully finds compliance data with updatedBeforeOrAt', async () => {
    await Promise.all([
      ComplianceDataTable.create(blockedComplianceData),
      ComplianceDataTable.create(nonBlockedComplianceData),
    ]);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {
        updatedBeforeOrAt: blockedComplianceData.updatedAt,
      },
      [],
      { readReplica: true },
    );

    expect(complianceData.length).toEqual(1);
    expect(complianceData[0]).toEqual(blockedComplianceData);
  });

  it('Successfully finds compliance data with provider', async () => {
    await Promise.all([
      ComplianceDataTable.create(blockedComplianceData),
      ComplianceDataTable.create(nonBlockedComplianceData),
    ]);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {
        provider: ComplianceProvider.ELLIPTIC,
      },
      [],
      { readReplica: true },
    );

    expect(complianceData.length).toEqual(2);
    expect(complianceData[0]).toEqual(blockedComplianceData);
    expect(complianceData[1]).toEqual(nonBlockedComplianceData);
  });

  it('Successfully finds compliance data with blocked', async () => {
    await Promise.all([
      ComplianceDataTable.create(blockedComplianceData),
      ComplianceDataTable.create(nonBlockedComplianceData),
    ]);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {
        blocked: false,
      },
      [],
      { readReplica: true },
    );

    expect(complianceData.length).toEqual(1);
    expect(complianceData[0]).toEqual(nonBlockedComplianceData);
  });

  it('Successfully finds compliance data by address and provider', async () => {
    await Promise.all([
      ComplianceDataTable.create(blockedComplianceData),
      ComplianceDataTable.create(nonBlockedComplianceData),
    ]);

    const complianceData:
    ComplianceDataFromDatabase | undefined = await ComplianceDataTable.findByAddressAndProvider(
      blockedAddress,
      ComplianceProvider.ELLIPTIC,
      { readReplica: true },
    );

    expect(complianceData).toBeDefined();
    expect(complianceData).toEqual(blockedComplianceData);
  });

  it('Unable finds compliance data', async () => {
    const complianceData:
    ComplianceDataFromDatabase | undefined = await ComplianceDataTable.findByAddressAndProvider(
      blockedAddress,
      ComplianceProvider.ELLIPTIC,
      { readReplica: true },
    );
    expect(complianceData).toEqual(undefined);
  });

  it('Successfully updates compliance data', async () => {
    await ComplianceDataTable.create(nonBlockedComplianceData);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(complianceData.length).toEqual(1);

    const updatedTime: string = DateTime.fromISO(
      nonBlockedComplianceData.updatedAt!,
    ).plus(10).toUTC().toISO();

    await ComplianceDataTable.update({
      address: nonBlockedComplianceData.address,
      provider: nonBlockedComplianceData.provider,
      riskScore: '30.00',
      blocked: true,
      updatedAt: updatedTime,
    });
    const updatedComplianceData:
    ComplianceDataFromDatabase | undefined = await ComplianceDataTable.findByAddressAndProvider(
      nonBlockedComplianceData.address,
      nonBlockedComplianceData.provider,
      { readReplica: true },
    );

    expect(updatedComplianceData).toEqual({
      ...nonBlockedComplianceData,
      riskScore: '30.00',
      blocked: true,
      updatedAt: updatedTime,
    });
  });

  it('Successfully upserts a new compliance data', async () => {
    await ComplianceDataTable.upsert(nonBlockedComplianceData);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(complianceData.length).toEqual(1);
    expect(complianceData[0]).toEqual(nonBlockedComplianceData);
  });

  it('Successfully upserts an existing compliance data', async () => {
    await ComplianceDataTable.upsert(nonBlockedComplianceData);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(complianceData.length).toEqual(1);

    const updatedTime: string = DateTime.fromISO(
      nonBlockedComplianceData.updatedAt!,
    ).plus(10).toUTC().toISO();

    await ComplianceDataTable.upsert({
      address: nonBlockedComplianceData.address,
      provider: nonBlockedComplianceData.provider,
      riskScore: '30.00',
      blocked: true,
      updatedAt: updatedTime,
    });
    const updatedComplianceData:
    ComplianceDataFromDatabase | undefined = await ComplianceDataTable.findByAddressAndProvider(
      nonBlockedComplianceData.address,
      nonBlockedComplianceData.provider,
      { readReplica: true },
    );

    expect(updatedComplianceData).toEqual({
      ...nonBlockedComplianceData,
      riskScore: '30.00',
      blocked: true,
      updatedAt: updatedTime,
    });
  });

  it('Successfully bulk upserts compliance data', async () => {
    await ComplianceDataTable.create(nonBlockedComplianceData);

    let complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(complianceData.length).toEqual(1);

    const updatedTime: string = DateTime.fromISO(
      nonBlockedComplianceData.updatedAt!,
    ).plus(10).toUTC().toISO();

    await ComplianceDataTable.bulkUpsert(
      [
        blockedComplianceData,
        {
          ...nonBlockedComplianceData,
          riskScore: '30.00',
          blocked: true,
          updatedAt: updatedTime,
        },
      ],
    );

    complianceData = await ComplianceDataTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(complianceData.length).toEqual(2);
    expect(complianceData[0]).toEqual(blockedComplianceData);
    expect(complianceData[1]).toEqual({
      ...nonBlockedComplianceData,
      riskScore: '30.00',
      blocked: true,
      updatedAt: updatedTime,
    });
  });
});
