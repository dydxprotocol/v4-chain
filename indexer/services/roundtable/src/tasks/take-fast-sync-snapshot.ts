import { logger, stats } from '@dydxprotocol-indexer/base';
import RDS from 'aws-sdk/clients/rds';
import { DateTime } from 'luxon';

import config from '../config';
import {
  createDBSnapshot,
  getMostRecentDBSnapshotIdentifier,
} from '../helpers/aws';

const statStart: string = `${config.SERVICE_NAME}.fast_sync_export_db_snapshot`;

export default async function runTask(): Promise<void> {
  const at: string = 'fast-sync-export-db-snapshot#runTask';
  logger.info({ at, message: 'Starting task.' });

  const rds: RDS = new RDS();

  const dateString: string = DateTime.utc().toFormat('yyyy-MM-dd-HH-mm');
  const snapshotIdentifier: string = `${config.FAST_SYNC_SNAPSHOT_IDENTIFIER_PREFIX}-${config.RDS_INSTANCE_NAME}-${dateString}`;
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
      stats.increment(`${statStart}.existingDbSnapshot`, 1);
      logger.info({
        at,
        message: 'Last fast sync db snapshot was taken less than the interval ago',
        interval: config.LOOPS_INTERVAL_MS_TAKE_FAST_SYNC_SNAPSHOTS,
        currentDate: dateString,
        lastSnapshotDate: s3Date,
      });
      return;
    }
  }
  // Create the DB snapshot
  const startSnapshot: number = Date.now();
  const createdSnapshotIdentifier: string = await
  createDBSnapshot(rds, snapshotIdentifier, config.RDS_INSTANCE_NAME);
  logger.info({ at, message: 'Created DB snapshot.', snapshotIdentifier: createdSnapshotIdentifier });
  stats.timing(`${statStart}.createDbSnapshot`, Date.now() - startSnapshot);
}

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
