import {
  ComplianceStatus,
  ComplianceStatusFromDatabase,
  ComplianceDataColumns,
  Ordering,
  ComplianceStatusUpsertObject,
} from '../../src/types';
import * as ComplianceStatusTable from '../../src/stores/compliance-status-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  compliantStatusData,
  defaultAddress,
  noncompliantStatusData,
  noncompliantStatusUpsertData,
} from '../helpers/constants';
import { DateTime } from 'luxon';

describe('Compliance status store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates compliance status', async () => {
    await Promise.all([
      ComplianceStatusTable.create(compliantStatusData),
      ComplianceStatusTable.create(noncompliantStatusData),
    ]);
  });

  it('Successfully creates compliance status without createdAt/updatedAt', async () => {
    await ComplianceStatusTable.create({
      address: defaultAddress,
      status: ComplianceStatus.COMPLIANT,
    });
    const complianceStatus: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(complianceStatus.length).toEqual(1);
    expect(complianceStatus[0]).toEqual(expect.objectContaining({
      address: defaultAddress,
      status: ComplianceStatus.COMPLIANT,
    }));
  });

  it('Successfully finds all compliance status', async () => {
    await Promise.all([
      ComplianceStatusTable.create(compliantStatusData),
      ComplianceStatusTable.create(noncompliantStatusData),
    ]);

    const complianceStatus: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(complianceStatus.length).toEqual(2);
    expect(complianceStatus[0]).toEqual({
      ...compliantStatusData,
      reason: null,
    });
    expect(complianceStatus[1]).toEqual(noncompliantStatusData);
  });

  it('Successfully finds compliance status with address', async () => {
    await Promise.all([
      ComplianceStatusTable.create(compliantStatusData),
      ComplianceStatusTable.create(noncompliantStatusData),
    ]);

    const complianceStatus: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll(
      {
        address: [noncompliantStatusData.address],
      },
      [],
    );

    expect(complianceStatus.length).toEqual(1);
    expect(complianceStatus[0]).toEqual(expect.objectContaining(noncompliantStatusData));
  });

  it('Successfully finds compliance status with status', async () => {
    await Promise.all([
      ComplianceStatusTable.create(compliantStatusData),
      ComplianceStatusTable.create(noncompliantStatusData),
    ]);

    const complianceStatus: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll(
      {
        status: [ComplianceStatus.COMPLIANT],
      },
      [],
      { readReplica: true },
    );

    expect(complianceStatus.length).toEqual(1);
    expect(complianceStatus[0]).toEqual(expect.objectContaining({
      ...compliantStatusData,
      reason: null,
    }));
  });

  it('Successfully updates compliance status', async () => {
    await ComplianceStatusTable.create(noncompliantStatusData);

    const updatedTime: string = DateTime.fromISO(
      noncompliantStatusData.createdAt!,
    ).plus({ minutes: 10 }).toUTC().toISO();

    await ComplianceStatusTable.update({
      address: noncompliantStatusData.address,
      status: ComplianceStatus.CLOSE_ONLY,
      updatedAt: updatedTime,
    });

    const updatedComplianceStatus: ComplianceStatusFromDatabase[] = await
    ComplianceStatusTable.findAll(
      {
        address: [noncompliantStatusData.address],
      },
      [],
    );

    expect(updatedComplianceStatus.length).toEqual(1);
    expect(updatedComplianceStatus[0]).toEqual({
      ...noncompliantStatusData,
      status: ComplianceStatus.CLOSE_ONLY,
      updatedAt: updatedTime,
    });
  });

  it('Successfully upserts a new compliance status', async () => {
    await ComplianceStatusTable.upsert(noncompliantStatusUpsertData);

    const complianceStatus: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(complianceStatus.length).toEqual(1);
    expect(complianceStatus[0]).toEqual(expect.objectContaining(noncompliantStatusUpsertData));
  });

  it('Successfully upserts an existing compliance status', async () => {
    await ComplianceStatusTable.upsert(noncompliantStatusUpsertData);

    await ComplianceStatusTable.upsert({
      ...noncompliantStatusUpsertData,
      status: ComplianceStatus.CLOSE_ONLY,
    });

    const updatedComplianceStatus: ComplianceStatusFromDatabase[] = await
    ComplianceStatusTable.findAll(
      {
        address: [noncompliantStatusData.address],
      },
      [],
      { readReplica: true },
    );
    expect(updatedComplianceStatus.length).toEqual(1);
    expect(updatedComplianceStatus[0]).toEqual(
      expect.objectContaining({
        ...noncompliantStatusUpsertData,
        status: ComplianceStatus.CLOSE_ONLY,
      }),
    );
  });

  it('Successfully bulk upserts compliance status', async () => {
    const compliantUpsertStatusData: ComplianceStatusUpsertObject = {
      address: defaultAddress,
      status: ComplianceStatus.COMPLIANT,
      updatedAt: DateTime.utc().toISO(),
    };
    await ComplianceStatusTable.create(noncompliantStatusUpsertData);
    const otherAddress: string = '0x123456789abcdef';

    await ComplianceStatusTable.bulkUpsert(
      [
        compliantUpsertStatusData,
        {
          ...noncompliantStatusUpsertData,
          status: ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY,
        },
        {
          ...noncompliantStatusUpsertData,
          address: otherAddress,
          status: ComplianceStatus.COMPLIANT,
        },
      ],
    );

    const complianceStatus = await ComplianceStatusTable.findAll(
      {},
      [],
      {
        orderBy: [[ComplianceDataColumns.address, Ordering.DESC]],
      },
    );

    expect(complianceStatus.length).toEqual(3);
    expect(complianceStatus).toEqual(expect.arrayContaining([
      expect.objectContaining({
        ...compliantUpsertStatusData,
        reason: null,
      }),
      expect.objectContaining({
        ...noncompliantStatusUpsertData,
        status: ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY,
      }),
      expect.objectContaining({
        ...noncompliantStatusUpsertData,
        address: otherAddress,
        status: ComplianceStatus.COMPLIANT,
      }),
    ]));
  });
});
