import { ComplianceDataFromDatabase } from '../../src/types';
import * as ComplianceDataTable from '../../src/stores/compliance-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import {
  sanctionedComplianceData,
  sanctionedAddress,
  nonSanctionedComplianceData,
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
      ComplianceDataTable.create(sanctionedComplianceData),
      ComplianceDataTable.create(nonSanctionedComplianceData),
    ]);
  });

  it('Successfully finds all compliance data', async () => {
    await Promise.all([
      ComplianceDataTable.create(sanctionedComplianceData),
      ComplianceDataTable.create(nonSanctionedComplianceData),
    ]);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(complianceData.length).toEqual(2);
    expect(complianceData[0]).toEqual(expect.objectContaining(sanctionedComplianceData));
    expect(complianceData[1]).toEqual(expect.objectContaining(nonSanctionedComplianceData));
  });

  it('Successfully finds compliance data with updatedBeforeOrAt', async () => {
    await Promise.all([
      ComplianceDataTable.create(sanctionedComplianceData),
      ComplianceDataTable.create(nonSanctionedComplianceData),
    ]);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {
        updatedBeforeOrAt: sanctionedComplianceData.updatedAt,
      },
      [],
      { readReplica: true },
    );

    expect(complianceData.length).toEqual(1);
    expect(complianceData[0]).toEqual(expect.objectContaining(sanctionedComplianceData));
  });

  it('Successfully finds compliance data by address', async () => {
    await Promise.all([
      ComplianceDataTable.create(sanctionedComplianceData),
      ComplianceDataTable.create(nonSanctionedComplianceData),
    ]);

    const complianceData:
    ComplianceDataFromDatabase | undefined = await ComplianceDataTable.findByAddress(
      sanctionedAddress,
      { readReplica: true },
    );

    expect(complianceData).toBeDefined();
    expect(complianceData).toEqual(expect.objectContaining(sanctionedComplianceData));
  });

  it('Unable finds compliance data', async () => {
    const complianceData:
    ComplianceDataFromDatabase | undefined = await ComplianceDataTable.findByAddress(
      sanctionedAddress,
      { readReplica: true },
    );
    expect(complianceData).toEqual(undefined);
  });

  it('Successfully updates compliance data', async () => {
    await ComplianceDataTable.create(nonSanctionedComplianceData);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(complianceData.length).toEqual(1);

    const updatedTime: string = DateTime.fromISO(
      nonSanctionedComplianceData.updatedAt,
    ).plus(10).toUTC().toISO();

    await ComplianceDataTable.update({
      address: nonSanctionedComplianceData.address,
      riskScore: '30.00',
      sanctioned: true,
      updatedAt: updatedTime,
    });
    const updatedComplianceData:
    ComplianceDataFromDatabase | undefined = await ComplianceDataTable.findByAddress(
      nonSanctionedComplianceData.address,
      { readReplica: true },
    );

    expect(updatedComplianceData).toEqual(expect.objectContaining({
      ...nonSanctionedComplianceData,
      riskScore: '30.00',
      sanctioned: true,
      updatedAt: updatedTime,
    }));
  });
});
