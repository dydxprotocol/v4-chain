import { InfoObject, logger, stats } from '@dydxprotocol-indexer/base';
import RDS from 'aws-sdk/clients/rds';
import { DateTime } from 'luxon';

import config from '../config';
import {
  createDBSnapshot,
  FAST_SYNC_SNAPSHOT_S3_BUCKET_NAME,
  getMostRecentDBSnapshotIdentifier,
  startExportTask,
} from '../helpers/aws';

const statStart: string = `${config.SERVICE_NAME}.fast_sync_export_db_snapshot`;

/**
 * Checks if the difference between two dates is less than a given interval.
 *
 * @param startDate
 * @param endDate
 * @param intervalMs
 */
function isDifferenceLessThanInterval(
  startDate: string,
  endDate: string,
  intervalMs: number,
): boolean {
  const parseDateString = (dateStr: string): Date => {
    const [year, month, day, hour, minute] = dateStr.split('-').map(Number);
    return new Date(year, month, day, hour, minute);
  };

  // Parse the date strings
  const parsedDate1 = parseDateString(startDate);
  const parsedDate2 = parseDateString(endDate);

  // Calculate the difference in milliseconds
  const differenceInMilliseconds = Math.abs(parsedDate1.getTime() - parsedDate2.getTime());

  // Compare with the interval
  return differenceInMilliseconds < intervalMs;
}

export default async function runTask(): Promise<void> {
  const at: string = 'fast-sync-export-db-snapshot#runTask';
  logger.info({ at, message: 'Starting task.' });

  const rds: RDS = new RDS();

  const dateString: string = DateTime.utc().toFormat('yyyy-MM-dd-HH-mm');
  const rdsExportIdentifier: string = `${config.FAST_SYNC_SNAPSHOT_IDENTIFIER_PREFIX}-${config.RDS_INSTANCE_NAME}-${dateString}`;
  // check the time of the last snapshot
  const lastSnapshotIdentifier: string | undefined = await getMostRecentDBSnapshotIdentifier(
    rds,
    config.FAST_SYNC_SNAPSHOT_IDENTIFIER_PREFIX,
  );
  if (lastSnapshotIdentifier !== undefined) {
    const s3Date: string = lastSnapshotIdentifier.split(config.RDS_INSTANCE_NAME)[1].slice(1);
    if (
      isDifferenceLessThanInterval(
        s3Date,
        dateString,
        config.LOOPS_INTERVAL_MS_TAKE_FAST_SYNC_SNAPSHOTS,
      )
    ) {
      logger.info({
        at,
        message: 'Last fast sync db snapshot was taken less than the interval ago',
        interval: config.LOOPS_INTERVAL_MS_TAKE_FAST_SYNC_SNAPSHOTS,
      });
      return;
    }
  }
  // Create the DB snapshot
  await createDBSnapshot(rds, rdsExportIdentifier, config.RDS_INSTANCE_NAME);

  // start S3 Export Job.
  const startExport: number = Date.now();
  try {
    const exportData: RDS.ExportTask = await startExportTask(
      rds,
      rdsExportIdentifier,
      FAST_SYNC_SNAPSHOT_S3_BUCKET_NAME,
      false,
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
