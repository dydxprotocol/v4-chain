import { logger, stats } from '@dydxprotocol-indexer/base';
import { FirebaseNotificationTokenTable } from '@dydxprotocol-indexer/postgres';

import config from '../config';

const statStart: string = `${config.SERVICE_NAME}.delete_old_firebase_notification_tokens`;

export default async function runTask(): Promise<void> {
  const at: string = 'delete-old-firebase-notification-tokens#runTask';
  const startDeleteOldFirebase: number = Date.now();
  // Delete old snapshots.
  stats.timing(statStart, Date.now() - startDeleteOldFirebase);

  try {
    // Delete tokens older than a month
    const oneMonthAgo = new Date();
    oneMonthAgo.setMonth(oneMonthAgo.getMonth() - 1);

    const tokensToDelete = await FirebaseNotificationTokenTable.findAll(
      {
        updatedBeforeOrAt: oneMonthAgo.toISOString(),
      },
      [],
    );

    if (tokensToDelete.length > 0) {
      await FirebaseNotificationTokenTable.deleteMany(
        tokensToDelete.map((tokenRecord) => tokenRecord.token),
      ).then((count) => {
        stats.increment(`${config.SERVICE_NAME}.firebase_notification_tokens_deleted`, count);
      });
    }
  } catch (error) {
    logger.info({ at, error, message: 'Failed to delete old Firebase notification tokens' });
  } finally {
    stats.timing(statStart, Date.now() - startDeleteOldFirebase);
  }
}
