import config from '../../src/config';
import { asMock } from '@dydxprotocol-indexer/dev';
import {
  checkIfExportJobToS3IsOngoing,
  checkIfS3ObjectExists,
  checkIfTableExistsInAthena,
  getMostRecentDBSnapshotIdentifier,
  startAthenaQuery,
  startExportTask,
} from '../../src/helpers/aws';
import updateResearchEnvironmentTask from '../../src/tasks/update-research-environment';

jest.mock('../../src/helpers/aws');

describe('update-research-environment', () => {
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

    await updateResearchEnvironmentTask();

    expect(checkIfExportJobToS3IsOngoing).not.toHaveBeenCalled();
    expect(startExportTask).not.toHaveBeenCalled();
  });

  it('export job in progress', async () => {
    asMock(checkIfExportJobToS3IsOngoing).mockImplementation(
      async () => Promise.resolve(true));

    await updateResearchEnvironmentTask();

    expect(startExportTask).not.toHaveBeenCalled();
  });

  it('start export job', async () => {
    await updateResearchEnvironmentTask();

    expect(startExportTask).toHaveBeenCalled();
  });

  it('Athena tables exist', async () => {
    asMock(checkIfS3ObjectExists).mockImplementation(async () => Promise.resolve(true));
    asMock(checkIfTableExistsInAthena).mockImplementation(async () => Promise.resolve(true));

    await updateResearchEnvironmentTask();

    expect(startAthenaQuery).not.toHaveBeenCalled();
  });

  it('Athena tables do not exist', async () => {
    asMock(checkIfS3ObjectExists).mockImplementation(async () => Promise.resolve(true));

    await updateResearchEnvironmentTask();

    expect(startAthenaQuery).toHaveBeenCalled();
  });
});
