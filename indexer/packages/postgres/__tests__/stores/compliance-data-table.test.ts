import { ComplianceDataFromDatabase, ComplianceProvider } from '../../src/types';
import * as ComplianceDataTable from '../../src/stores/compliance-table';
import * as WalletTable from '../../src/stores/wallet-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import {
  blockedComplianceData,
  blockedAddress,
  nonBlockedComplianceData,
  defaultWallet,
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

  it('Successfully filters by onlyDydxAddressWithDeposit', async () => {
    // Create two compliance entries, one with a corresponding wallet entry and another without
    await Promise.all([
      WalletTable.create(defaultWallet),
      ComplianceDataTable.create(nonBlockedComplianceData),
      ComplianceDataTable.create({
        ...nonBlockedComplianceData,
        address: 'not_dydx_address',
      }),
    ]);

    const complianceData: ComplianceDataFromDatabase[] = await ComplianceDataTable.findAll(
      {
        addressInWalletsTable: true,
      },
      [],
      { readReplica: true },
    );

    expect(complianceData.length).toEqual(1);
    expect(complianceData[0]).toEqual(nonBlockedComplianceData);
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

    const updatedTime1: string = DateTime.fromISO(
      nonBlockedComplianceData.updatedAt!,
    ).plus(10).toUTC().toISO();
    const updatedTime2: string = DateTime.fromISO(
      nonBlockedComplianceData.updatedAt!,
    ).plus(20).toUTC().toISO();
    const otherAddress: string = 'dydx1scu097p2sstqzupe6t687kpc2w4sv665fedctf';

    await ComplianceDataTable.bulkUpsert(
      [
        blockedComplianceData,
        {
          ...nonBlockedComplianceData,
          riskScore: '30.00',
          blocked: true,
          updatedAt: updatedTime1,
        },
        {
          ...nonBlockedComplianceData,
          address: otherAddress,
          riskScore: undefined,
          blocked: false,
          updatedAt: updatedTime2,
        },
      ],
    );

    complianceData = await ComplianceDataTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(complianceData.length).toEqual(3);
    expect(complianceData[0]).toEqual(blockedComplianceData);
    expect(complianceData[1]).toEqual({
      ...nonBlockedComplianceData,
      riskScore: '30.00',
      blocked: true,
      updatedAt: updatedTime1,
    });
    expect(complianceData[2]).toEqual({
      ...nonBlockedComplianceData,
      address: otherAddress,
      riskScore: null,
      blocked: false,
      updatedAt: updatedTime2,
    });
  });
});
