import { logger, stats } from '@dydxprotocol-indexer/base';
import RDS from 'aws-sdk/clients/rds';

import config from '../config';
import { deleteOldFastSyncSnapshots } from '../helpers/aws';

const statStart: string = `${config.SERVICE_NAME}.delete_old_fast_sync_snapshots`;

export default async function runTask(): Promise<void> {
  const at: string = 'delete-old-fast-sync-snapshots#runTask';
  logger.info({ at, message: 'Starting task.' });

  const rds: RDS = new RDS();

  const startDeleteOldSnapshot: number = Date.now();
  // Delete old snapshots.
  await deleteOldFastSyncSnapshots(rds);
  stats.timing(`${statStart}.deleteOldSnapshots`, Date.now() - startDeleteOldSnapshot);
}
