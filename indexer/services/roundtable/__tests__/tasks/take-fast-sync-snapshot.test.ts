import config from '../../src/config';
import { asMock } from '@dydxprotocol-indexer/dev';
import {
  checkIfExportJobToS3IsOngoing,
  checkIfS3ObjectExists,
  getMostRecentDBSnapshotIdentifier,
  startExportTask,
} from '../../src/helpers/aws';
import takeFastSyncSnapshotTask from '../../src/tasks/take-fast-sync-snapshot';

jest.mock('../../src/helpers/aws');

describe('fast-sync-export-db-snapshot', () => {
  beforeAll(() => {
    config.RDS_INSTANCE_NAME = 'postgres-main-staging';
  });

  beforeEach(() => {
    jest.resetAllMocks();
    asMock(getMostRecentDBSnapshotIdentifier).mockImplementation(async () => Promise.resolve('postgres-main-staging-2022-05-03-04-16'));
  });

  afterAll(jest.resetAllMocks);

  it('s3Object exists', async () => {
    asMock(checkIfS3ObjectExists).mockImplementation(async () => Promise.resolve(true));

    await takeFastSyncSnapshotTask();

    expect(checkIfExportJobToS3IsOngoing).not.toHaveBeenCalled();
    expect(startExportTask).not.toHaveBeenCalled();
  });

  it('export job in progress', async () => {
    asMock(checkIfExportJobToS3IsOngoing).mockImplementation(
      async () => Promise.resolve(true));

    await takeFastSyncSnapshotTask();

    expect(startExportTask).not.toHaveBeenCalled();
  });

  it('start export job', async () => {
    await takeFastSyncSnapshotTask();

    expect(startExportTask).toHaveBeenCalled();
  });
});
