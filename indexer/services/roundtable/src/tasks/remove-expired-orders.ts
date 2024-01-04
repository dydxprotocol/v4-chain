import {
  logger,
  stats,
} from '@dydxprotocol-indexer/base';
import { BatchKafkaProducer, KafkaTopics, producer } from '@dydxprotocol-indexer/kafka';
import { BlockTable, BlockFromDatabase } from '@dydxprotocol-indexer/postgres';
import {
  OrderData,
  OrderExpiryCache,
  OrdersCache,
  OrdersDataCache,
} from '@dydxprotocol-indexer/redis';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import { IndexerOrder, RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';

import config from '../config';
import { redisClient } from '../helpers/redis';
import { getExpiredOffChainUpdateMessage } from '../helpers/websocket';

/**
 * This task uses the latest block height (minus a buffer) to pull expired order UUIDs from the
 * expiry cache and send messages to Vulcan to expire them.
 *
 * NOTE: Every message sent will use the reason ORDER_REMOVAL_REASON_INDEXER_EXPIRED.
 */
export default async function runTask(): Promise<void> {
  const start: number = Date.now();
  const block: BlockFromDatabase = await BlockTable.getLatest({ readReplica: true });

  try {
    // Only need to expire short-term orders because long-term OrderRemoves will be on-chain.
    // Short-term orders exclusively use blockHeight.
    const expiryCutoff: number = (
      +block.blockHeight - config.BLOCKS_TO_DELAY_EXPIRY_BEFORE_SENDING_REMOVES
    );
    const orderUuids: string[] = await OrderExpiryCache.getOrderExpiries({
      latestExpiry: expiryCutoff,
      latestExpiryIsInclusive: true,
    }, redisClient) as string[];

    const batchKafkaProducer: BatchKafkaProducer = new BatchKafkaProducer(
      KafkaTopics.TO_VULCAN,
      producer,
      config.KAFKA_MAX_BATCH_WEBSOCKET_MESSAGE_SIZE_BYTES,
    );

    const allPromises: Promise<void>[] = _.map(orderUuids, (orderUuid: string) => {
      return addExpireOrderMessageOrFail(orderUuid, expiryCutoff, batchKafkaProducer);
    });

    // The promises listed make calls to `batchKafkaProducer.addMessageAndMaybeFlush`, so we need to
    // wait for their completion before we're certain all the messages have been added to the
    // `batchKafkaProducer`. Only after can we call `.flush()`.
    // Use `allSettled` because if a single order uuid fails, we still want to expire the rest.
    await Promise.allSettled(allPromises);
    await batchKafkaProducer.flush();
  } catch (error) {
    logger.error({
      at: 'remove-expired-orders#runTask',
      message: 'Error occurred in task to remove expired orders',
      error,
    });
  } finally {
    stats.timing(
      `${config.SERVICE_NAME}.remove_expired_orders.timing`,
      Date.now() - start,
    );
  }
}

async function addExpireOrderMessageOrFail(
  orderUuid: string,
  expiryCutoff: number,
  messenger: BatchKafkaProducer,
): Promise<void> {
  try {
    const [redisOrder, orderData]: [RedisOrder | null, OrderData | null] = await Promise.all([
      OrdersCache.getOrder(orderUuid, redisClient),
      OrdersDataCache.getOrderDataWithUUID(orderUuid, redisClient),
    ]);

    if (redisOrder == null || orderData == null) {
      stats.increment(`${config.SERVICE_NAME}.expired_order_data_not_found`, 1);
      logger.info({
        at: 'remove-expired-orders#runTask',
        message: 'Unable to retrieve data from one of { orders-cache, orders-data-cache }',
        orderUuid,
        redisOrder,
        orderData,
      });
      return;
    }
    const order: IndexerOrder = redisOrder.order!;
    if (order.goodTilBlock! > expiryCutoff) {
      stats.increment(`${config.SERVICE_NAME}.indexer_expired_order_has_newer_expiry`, 1);
      logger.info({
        at: 'remove-expired-orders#runTask',
        message: 'Order expiry lower than order (goodTilBlock - buffer)',
        expectedExpiry: `${expiryCutoff} or lower`,
        actualExpiry: order.goodTilBlock!,
        orderUuid,
        redisOrder,
        orderData,
      });
      return;
    }
    if (order.quantums.lte(orderData.totalFilledQuantums)) {
      stats.increment(`${config.SERVICE_NAME}.fully_filled_orders_expired_by_roundtable`, 1);
    }

    messenger.addMessageAndMaybeFlush({
      key: getOrderIdHash(order.orderId!),
      value: getExpiredOffChainUpdateMessage(order.orderId!),
    });

    stats.increment(`${config.SERVICE_NAME}.expiry_message_sent`, 1);
  } catch (error) {
    logger.error({
      at: 'remove-expired-orders#runTask',
      message: 'Encountered error expiring order',
      orderUuid,
      expiryCutoff,
      error,
    });
  }
}
