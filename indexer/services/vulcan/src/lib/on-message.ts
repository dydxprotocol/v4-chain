import {
  getInstanceId,
  logger,
  stats,
  ParseMessageError,
  STATS_NO_SAMPLING,
} from '@dydxprotocol-indexer/base';
import { KafkaTopics } from '@dydxprotocol-indexer/kafka';
import { OffChainUpdateV1 } from '@dydxprotocol-indexer/v4-protos';
import { IHeaders, KafkaMessage } from 'kafkajs';
import { Handler } from 'src/handlers/handler';

import config from '../config';
import { OrderPlaceHandler } from '../handlers/order-place-handler';
import { OrderRemoveHandler } from '../handlers/order-remove-handler';
import { OrderUpdateHandler } from '../handlers/order-update-handler';
import { DydxRecordHeaderKeys } from './types';

export type HandlerInitializer = new (
  txHash?: string
) => Handler;

function getHandler(update: OffChainUpdateV1): HandlerInitializer | undefined {
  if (update.orderUpdate !== undefined) {
    return OrderUpdateHandler;
  } else if (update.orderPlace !== undefined) {
    return OrderPlaceHandler;
  } else if (update.orderRemove !== undefined) {
    return OrderRemoveHandler;
  }
  return undefined;
}

function getMessageType(update: OffChainUpdateV1): string {
  if (update.orderUpdate !== undefined) {
    return 'orderUpdate';
  } else if (update.orderPlace !== undefined) {
    return 'orderPlace';
  } else if (update.orderRemove !== undefined) {
    return 'orderRemove';
  }
  return 'unknown';
}

export async function onMessage(message: KafkaMessage): Promise<void> {
  stats.increment(
    `${config.SERVICE_NAME}.received_kafka_message`,
    1,
    { instance: getInstanceId() },
  );
  if (!message || !message.value || !message.timestamp) {
    stats.increment(
      `${config.SERVICE_NAME}.empty_kafka_message`,
      1,
      { instance: getInstanceId() },
    );
    logger.error({
      at: 'onMessage#onMessage',
      message: 'Empty message',
    });
    return;
  }

  const start: number = Date.now();
  stats.timing(
    `${config.SERVICE_NAME}.message_time_in_queue`,
    start - Number(message.timestamp),
    STATS_NO_SAMPLING,
    {
      topic: KafkaTopics.TO_VULCAN,
      instance: getInstanceId(),
    },
  );

  const originalMessageTimestamp = message.headers?.message_received_timestamp;
  if (originalMessageTimestamp !== undefined) {
    stats.timing(
      `${config.SERVICE_NAME}.message_time_since_received`,
      start - Number(originalMessageTimestamp),
      STATS_NO_SAMPLING,
      {
        topic: KafkaTopics.TO_VULCAN,
        event_type: String(message.headers?.event_type),
        instance: getInstanceId(),
      },
    );
  }

  const messageValue: Buffer = message.value;
  const offset: string = message.offset;
  let update: OffChainUpdateV1;

  try {
    update = getOffChainUpdate(messageValue, offset);
  } catch (error) {
    logger.crit({
      at: 'onMessage#onMessage',
      message: 'Error: Unable to parse message',
      offset,
      value: message.value,
      error,
    });
    return;
  }

  let success: boolean = false;
  try {
    validateOffChainUpdate(update);

    const handler: Handler = new (getHandler(update))!(
      getTransactionHashFromHeaders(message.headers),
    );

    // If headers don't exist, create them.
    const headers = message.headers ?? {};
    // If the message received timestamp doesn't exist
    // (i.e when a short term order is directly sent to vulcan via full node)
    // set the message_received_timestamp to the message timestamp and the event type
    // to be a short term order event type.
    if (!headers.message_received_timestamp) {
      headers.message_received_timestamp = message.timestamp;
    }
    if (!headers.event_type) {
      if (update.orderPlace) {
        headers.event_type = 'ShortTermOrderPlacement';
      } else if (update.orderRemove) {
        headers.event_type = 'ShortTermOrderRemoval';
      } else if (update.orderUpdate) {
        headers.event_type = 'ShortTermOrderUpdate';
      }
    }

    await handler.handleUpdate(update, headers);

    const postProcessingTime: number = Date.now();
    if (headers.message_received_timestamp !== undefined) {
      stats.timing(
        `${config.SERVICE_NAME}.message_time_since_received_post_processing`,
        postProcessingTime - Number(headers.message_received_timestamp),
        STATS_NO_SAMPLING,
        {
          topic: KafkaTopics.TO_VULCAN,
          event_type: String(headers?.event_type),
          instance: getInstanceId(),
        },
      );
    }

    success = true;
  } catch (error) {
    if (error instanceof ParseMessageError) {
      // Do not re-throw error so message will not be retried
      logger.crit({
        at: 'onMessage#onMessage',
        message: 'Error: Unable to parse message, this must be due to a bug in the V4 node',
        offset,
        update,
        error,
      });
    } else {
      logger.error({
        at: 'onMessage#onMessage',
        message: 'Error: Unable to process message',
        offset,
        update,
        error,
      });
      // Throw error so the message is not acked and will be reprocessed
      throw error;
    }
  } finally {
    stats.timing(
      `${config.SERVICE_NAME}.processed_update.timing`,
      Date.now() - start,
      STATS_NO_SAMPLING,
      {
        success: success.toString(),
        messageType: getMessageType(update),
        instance: getInstanceId(),
      },
    );
  }
}

function getOffChainUpdate(messageValue: Buffer, offset: string): OffChainUpdateV1 {
  const messageValueBinary: Uint8Array = new Uint8Array(messageValue);
  logger.debug({
    at: 'onMessage#getOffChainupdate',
    message: 'Received message',
    offset,
  });

  return OffChainUpdateV1.decode(messageValueBinary);
}

function validateOffChainUpdate(update: OffChainUpdateV1) {
  if (update.orderUpdate === undefined &&
    update.orderPlace === undefined &&
    update.orderRemove === undefined) {
    throw new ParseMessageError('Message does not contain an order update, place, or remove');
  }
}

function getTransactionHashFromHeaders(headers?: IHeaders): string | undefined {
  if (headers === undefined) {
    return undefined;
  }

  const txHashBytes: (
    string |
    Buffer |
    (string | Buffer)[] |
    undefined
  ) = headers[DydxRecordHeaderKeys.TRANSACTION_HASH_KEY];
  if (txHashBytes === undefined) {
    return undefined;
  }

  if (!Buffer.isBuffer(txHashBytes)) {
    return undefined;
  }

  return txHashBytes.toString('hex').toUpperCase();
}
