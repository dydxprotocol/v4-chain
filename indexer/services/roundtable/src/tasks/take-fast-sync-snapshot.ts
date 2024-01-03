import { InfoObject, logger, stats } from '@dydxprotocol-indexer/base';
import RDS from 'aws-sdk/clients/rds';
import { DateTime } from 'luxon';

import config from '../config';
import {
  checkIfExportJobToS3IsOngoing,
  FAST_SYNC_SNAPSHOT_S3_BUCKET_NAME,
  startExportTask,
} from '../helpers/aws';

const statStart: string = `${config.SERVICE_NAME}.fast_sync_export_db_snapshot`;

export default async function runTask(): Promise<void> {
  const at: string = 'fast-sync-export-db-snapshot#runTask';
  logger.info({ at, message: 'Starting task.' });

  const rds: RDS = new RDS();

  const dateString: string = DateTime.utc().toFormat('yyyy-MM-dd-HH-mm');
  const rdsExportIdentifier: string = `${config.RDS_INSTANCE_NAME}-fast-sync-${dateString}`;

  // check if it is being created
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

  // start Export Job if not already started
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
