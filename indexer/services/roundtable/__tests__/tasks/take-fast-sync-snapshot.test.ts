import config from '../../src/config';
import { asMock } from '@dydxprotocol-indexer/dev';
import {
  createDBSnapshot,
  getMostRecentDBSnapshotIdentifier,
} from '../../src/helpers/aws';
import takeFastSyncSnapshotTask from '../../src/tasks/take-fast-sync-snapshot';
import { DateTime } from 'luxon';

jest.mock('../../src/helpers/aws');

describe('fast-sync-export-db-snapshot', () => {
  const snapshotIdentifier: string = `${config.FAST_SYNC_SNAPSHOT_IDENTIFIER_PREFIX}-postgres-main-staging-2022-05-03-04-16`;
  beforeAll(() => {
    config.RDS_INSTANCE_NAME = 'postgres-main-staging';
  });

  beforeEach(() => {
    jest.resetAllMocks();
    asMock(getMostRecentDBSnapshotIdentifier).mockImplementation(
      async () => Promise.resolve(snapshotIdentifier),
    );
  });

  afterAll(jest.resetAllMocks);

  it('Last snapshot was taken more than interval ago', async () => {
    await takeFastSyncSnapshotTask();

    expect(createDBSnapshot).toHaveBeenCalled();
  });

  it('Last snapshot was taken less than interval ago', async () => {
    const timestamp: string = DateTime.utc().minus({ minutes: 1 }).toFormat('yyyy-MM-dd-HH-mm');
    asMock(getMostRecentDBSnapshotIdentifier).mockImplementation(
      async () => Promise.resolve(`${config.FAST_SYNC_SNAPSHOT_IDENTIFIER_PREFIX}-postgres-main-staging-${timestamp}`),
    );

    await takeFastSyncSnapshotTask();

    expect(createDBSnapshot).not.toHaveBeenCalled();
  });

  it('No existing snapshot', async () => {
    asMock(getMostRecentDBSnapshotIdentifier).mockImplementation(
      async () => Promise.resolve(undefined),
    );

    await takeFastSyncSnapshotTask();

    expect(createDBSnapshot).toHaveBeenCalled();
  });
});
