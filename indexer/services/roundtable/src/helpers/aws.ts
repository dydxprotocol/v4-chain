import {
  logger,
} from '@dydxprotocol-indexer/base';
import Athena from 'aws-sdk/clients/athena';
import RDS from 'aws-sdk/clients/rds';
import S3 from 'aws-sdk/clients/s3';

import config from '../config';

const atStart: string = 'aws#';

enum ExportTaskStatus {
  CANCELED = 'canceled',
  CANCELING = 'canceling',
  FAILED = 'failed',
  COMPLETE = 'complete',
}

export const RESEARCH_SNAPSHOT_S3_BUCKET_NAME = config.RESEARCH_SNAPSHOT_S3_BUCKET_ARN.split(':::')[1];
export const RESEARCH_SNAPSHOT_S3_LOCATION_PREFIX = `s3://${RESEARCH_SNAPSHOT_S3_BUCKET_NAME}`;
export const FAST_SYNC_SNAPSHOT_S3_BUCKET_NAME = config.FAST_SYNC_SNAPSHOT_S3_BUCKET_ARN.split(':::')[1];

/**
 * @description Get most recent snapshot identifier for an RDS database.
 * @param rds - RDS client
 * @param snapshotIdentifierPrefixInclude - Only include snapshots with snapshot identifier
 * that starts with prefixInclude
 * @param snapshotIdentifierPrefixExclude - Only include snapshots with snapshot identifier
 * that does not start with prefixExclude
 */
// TODO(CLOB-672): Verify this function returns the most recent DB snapshot.
export async function getMostRecentDBSnapshotIdentifier(
  rds: RDS,
  snapshotIdentifierPrefixInclude?: string,
  snapshotIdentifierPrefixExclude?: string,
): Promise<string | undefined> {
  const awsResponse: RDS.DBSnapshotMessage = await rds.describeDBSnapshots({
    DBInstanceIdentifier: config.RDS_INSTANCE_NAME,
    MaxRecords: 20, // this is the minimum
  }).promise();

  if (awsResponse.DBSnapshots === undefined) {
    throw Error(`No DB snapshots found with identifier: ${config.RDS_INSTANCE_NAME}`);
  }

  let snapshots: RDS.DBSnapshotList = awsResponse.DBSnapshots;
  // Only include snapshots with snapshot identifier that starts with prefixInclude
  if (snapshotIdentifierPrefixInclude !== undefined) {
    snapshots = snapshots
      .filter((snapshot) => snapshot.DBSnapshotIdentifier &&
        snapshot.DBSnapshotIdentifier.startsWith(snapshotIdentifierPrefixInclude),
      );
  }
  if (snapshotIdentifierPrefixExclude !== undefined) {
    snapshots = snapshots
      .filter((snapshot) => snapshot.DBSnapshotIdentifier &&
        !snapshot.DBSnapshotIdentifier.startsWith(snapshotIdentifierPrefixExclude),
      );
  }

  logger.info({
    at: `${atStart}getMostRecentDBSnapshotIdentifier`,
    message: 'Described snapshots for database',
    mostRecentSnapshot: snapshots[snapshots.length - 1],
  });

  return snapshots[snapshots.length - 1]?.DBSnapshotIdentifier;
}

/**
 * @description Create DB snapshot for an RDS database. Only returns when the
 * snapshot is available.
 */
export async function createDBSnapshot(
  rds: RDS,
  snapshotIdentifier: string,
  dbInstanceIdentifier: string,
): Promise<string> {
  const params = {
    DBInstanceIdentifier: dbInstanceIdentifier,
    DBSnapshotIdentifier: snapshotIdentifier,
  };

  try {
    await rds.createDBSnapshot(params).promise();
    // Polling function to check snapshot status. Only return when the snapshot is available.
    const waitForSnapshot = async () => {
      // eslint-disable-next-line no-constant-condition
      while (true) {
        const statusResponse = await rds.describeDBSnapshots(
          { DBSnapshotIdentifier: snapshotIdentifier },
        ).promise();
        const snapshot = statusResponse.DBSnapshots![0];
        if (snapshot.Status === 'available') {
          return snapshot.DBSnapshotIdentifier!;
        } else if (snapshot.Status === 'failed') {
          throw Error(`Snapshot creation failed for identifier: ${snapshotIdentifier}`);
        }

        // Wait for 1 minute before checking again
        await new Promise((resolve) => setTimeout(resolve, 60000));
      }
    };

    return await waitForSnapshot();
  } catch (error) {
    logger.error({
      at: `${atStart}createDBSnapshot`,
      message: 'Failed to create DB snapshot',
      error,
      snapshotIdentifier,
    });
    throw error;
  }
}

/**
 * @description Check if an S3 Object already exists.
 */
export async function checkIfS3ObjectExists(
  s3: S3,
  s3Date: string,
  bucket: string,
): Promise<boolean> {
  const at: string = `${atStart}checkIfS3ObjectExists`;
  const key: string = `${config.RDS_INSTANCE_NAME}-${s3Date}/export_info_${config.RDS_INSTANCE_NAME}-${s3Date}.json`;

  logger.info({
    at,
    message: 'Going to query s3 bucket',
    bucket,
    key,
  });

  try {
    const awsResponse: S3.GetObjectOutput = await s3.getObject({
      Bucket: bucket,
      Key: key,
    }).promise();

    logger.info({
      at,
      message: 'Queried s3 bucket',
      lastModified: awsResponse.LastModified,
    });

    return true;
  } catch (error) {
    logger.info({
      at,
      message: 'Queried s3 bucket and received an error',
      error,
      s3Date,
    });

    if (error.statusCode === 404) {
      return false;
    }

    throw error;
  }
}

/**
 * @description Check if an export job to S3 is currently running.
 */
export async function checkIfExportJobToS3IsOngoing(
  rds: RDS,
  rdsExportIdentifier: string,
): Promise<boolean> {
  const at: string = `${atStart}checkIfExportJobToS3IsOngoing`;

  // get task status
  const awsResponse: RDS.ExportTasksMessage = await rds.describeExportTasks({
    ExportTaskIdentifier: rdsExportIdentifier,
  }).promise();

  logger.info({
    at,
    message: 'Checked if an export task is ongoing',
    data: awsResponse,
    rdsExportIdentifier,
  });

  // check if status was unexpected/invalid for dYdX
  if (
    awsResponse.ExportTasks !== undefined &&
      awsResponse.ExportTasks.length > 0 &&
      awsResponse.ExportTasks[0].Status !== undefined &&
      [
        ExportTaskStatus.CANCELED,
        ExportTaskStatus.CANCELING,
        ExportTaskStatus.FAILED,
      ].includes(awsResponse.ExportTasks[0].Status as ExportTaskStatus)
  ) {
    logger.error({
      at,
      message: 'Unexpected task status',
      exportTask: awsResponse.ExportTasks[0],
      rdsExportIdentifier,
    });

    throw Error('Unexpected task status');
  }

  // return if a task is ongoing:
  // 1. there is a task: not undefined and has a length gt 0
  // 2. the task is ongoing and not complete
  return awsResponse.ExportTasks !== undefined &&
    awsResponse.ExportTasks.length > 0 &&
    awsResponse.ExportTasks[0].Status !== ExportTaskStatus.COMPLETE;
}

/**
 * @description Start an export job from an RDS snapshot to an S3 bucket.
 * Link to API docs: https://docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_StartExportTask.html
 */
export async function startExportTask(
  rds: RDS,
  rdsExportIdentifier: string,
  bucket: string,
  isAutomatedSnapshot: boolean,
): Promise<RDS.ExportTask> {
  // TODO: Add validation
  let sourceArnPrefix: string = `arn:aws:rds:${config.AWS_REGION}:${config.AWS_ACCOUNT_ID}:snapshot:`;
  if (isAutomatedSnapshot) {
    sourceArnPrefix = sourceArnPrefix.concat('rds:');
  }
  const awsResponse: RDS.ExportTask = await rds.startExportTask({
    ExportTaskIdentifier: rdsExportIdentifier,
    S3BucketName: bucket,
    KmsKeyId: config.KMS_KEY_ARN,
    IamRoleArn: config.ECS_TASK_ROLE_ARN,
    SourceArn: `${sourceArnPrefix}${rdsExportIdentifier}`,
  }).promise();

  return awsResponse;
}

/**
 * @description Check if a table exists in Athena.
 */
export async function checkIfTableExistsInAthena(
  athena: Athena,
  table: string,
): Promise<boolean> {
  const at: string = `${atStart}checkIfTableExistsInAthena`;

  try {
    await athena.getTableMetadata({
      CatalogName: config.ATHENA_CATALOG_NAME,
      DatabaseName: config.ATHENA_DATABASE_NAME,
      TableName: table,
    }).promise();

    logger.info({
      at,
      message: 'got table',
      table,
    });

    return true;
  } catch (error) {
    logger.info({
      at,
      message: 'did not get table',
      error,
    });

    if (error.message.includes('EntityNotFoundException')) {
      return false;
    }

    throw error;
  }

}

/**
 * @description Start an Athena query.
 */
export async function startAthenaQuery(
  athena: Athena,
  {
    query,
    timestamp,
  }: {
    query: string,
    timestamp: string,
  },
): Promise<Athena.StartQueryExecutionOutput> {
  return athena.startQueryExecution({
    QueryString: query,
    QueryExecutionContext: {
      Catalog: config.ATHENA_CATALOG_NAME,
      Database: config.ATHENA_DATABASE_NAME,
    },
    ResultConfiguration: {
      OutputLocation: `${RESEARCH_SNAPSHOT_S3_LOCATION_PREFIX}/output/${timestamp}`,
    },
    WorkGroup: config.ATHENA_WORKING_GROUP,
  }).promise();
}
