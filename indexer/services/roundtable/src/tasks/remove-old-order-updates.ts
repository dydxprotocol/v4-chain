import {
  logger,
  stats,
} from '@dydxprotocol-indexer/base';
import {
  StatefulOrderUpdateInfo,
  StatefulOrderUpdatesCache,
} from '@dydxprotocol-indexer/redis';
import { OrderUpdateV1 } from '@dydxprotocol-indexer/v4-protos';

import config from '../config';
import { redisClient } from '../helpers/redis';

/**
 * This task removes any old cached stateful orer updates from the StatefulOrderUpdates cache
 */
export default async function runTask(): Promise<void> {
  const start: number = Date.now();

  try {
    const oldUpdateCutoff: number = Date.now() - config.OLD_CACHED_ORDER_UPDATES_WINDOW_MS;

    const oldUpdateInfo: StatefulOrderUpdateInfo[] = await StatefulOrderUpdatesCache
      .getOldOrderUpdates(
        oldUpdateCutoff,
        redisClient,
      );
    const removedUpdates: OrderUpdateV1[] = (await Promise.all(
      oldUpdateInfo.map(
        (updateInfo: StatefulOrderUpdateInfo): Promise<OrderUpdateV1 | undefined> => {
          return StatefulOrderUpdatesCache.removeStatefulOrderUpdate(
            updateInfo.orderId,
            updateInfo.timestamp,
            redisClient,
          );
        },
      ),
    )).filter(
      (removedUpdate: OrderUpdateV1 | undefined): removedUpdate is OrderUpdateV1 => {
        if (removedUpdate !== undefined) {
          logger.info({
            at: 'remove-old-order-updates#runTask',
            message: 'Removed old stateful order update',
            removedUpdate,
          });
          return true;
        }
        return false;
      },
    );

    stats.gauge(
      `${config.SERVICE_NAME}.remove_old_order_updates.num_removed`,
      removedUpdates.length,
    );
  } catch (error) {
    logger.error({
      at: 'remove-old-order-updates#runTas',
      message: 'Error occurred in task to remove old stateful order updates',
      error,
    });
  } finally {
    stats.timing(
      `${config.SERVICE_NAME}.remove_old_order_updates`,
      Date.now() - start,
    );
  }
}
