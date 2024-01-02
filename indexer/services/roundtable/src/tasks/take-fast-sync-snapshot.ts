import { InfoObject, logger, stats } from '@dydxprotocol-indexer/base';
import RDS from 'aws-sdk/clients/rds';
import S3 from 'aws-sdk/clients/s3';
import { DateTime } from 'luxon';

import config from '../config';
import {
  checkIfExportJobToS3IsOngoing,
  checkIfS3ObjectExists,
  FAST_SYNC_SNAPSHOT_S3_BUCKET_NAME,
  getMostRecentDBSnapshotIdentifier,
  startExportTask,
} from '../helpers/aws';

const statStart: string = `${config.SERVICE_NAME}.fast_sync_export_db_snapshot`;

export default async function runTask(): Promise<void> {
  const at: string = 'fast-sync-export-db-snapshot#runTask';

  const rds: RDS = new RDS();

  // get most recent rds snapshot
  const startDescribe: number = Date.now();
  const dateString: string = DateTime.utc().toFormat('yyyy-MM-dd');
  const mostRecentSnapshot: string = await getMostRecentDBSnapshotIdentifier(rds);
  stats.timing(`${statStart}.describe_rds_snapshots`, Date.now() - startDescribe);

  // dev example: rds:dev-indexer-apne1-db-2023-06-25-18-34
  const s3Date: string = mostRecentSnapshot.split(config.RDS_INSTANCE_NAME)[1].slice(1);
  const s3: S3 = new S3();

  // check if s3 object exists
  const startS3Check: number = Date.now();
  const s3ObjectExists: boolean = await checkIfS3ObjectExists(
    s3,
    s3Date,
    FAST_SYNC_SNAPSHOT_S3_BUCKET_NAME,
  );
  stats.timing(`${statStart}.checkS3Object`, Date.now() - startS3Check);

  const rdsExportIdentifier: string = `${config.RDS_INSTANCE_NAME}-fast-sync-${s3Date}`;

  // If the s3 object exists, return
  if (s3ObjectExists) {
    logger.info({
      at,
      dateString,
      message: 'S3 object exists.',
    });
    return;
  }

  // if we haven't created the object, check if it is being created
  const rdsExportCheck: number = Date.now();
  const exportJobOngoing: boolean = await checkIfExportJobToS3IsOngoing(rds, rdsExportIdentifier);
  stats.timing(`${statStart}.checkRdsExport`, Date.now() - rdsExportCheck);

  if (exportJobOngoing) {
    logger.info({
      at,
      dateString,
      message: 'Will wait for export job to finish',
    });
    return;
  }
  // start Export Job if S3 Object does not exist
  const startExport: number = Date.now();
  try {
    const exportData: RDS.ExportTask = await startExportTask(
      rds,
      rdsExportIdentifier,
      FAST_SYNC_SNAPSHOT_S3_BUCKET_NAME,
    );

    logger.info({
      at,
      message: 'Started an export task',
      exportData,
    });
  } catch (error) { // TODO handle this by finding the most recent snapshot earlier
    const message: InfoObject = {
      at,
      message: 'export to S3 failed',
      error,
    };

    if (error.name === 'DBSnapshotNotFound') {
      stats.increment(`${statStart}.no_s3_snapshot`, 1);

      logger.info(message);
      return;
    }

    logger.error(message);
  } finally {
    stats.timing(`${statStart}.rdsSnapshotExport`, Date.now() - startExport);
  }
}
