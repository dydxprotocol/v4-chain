import {
  getInstanceId, logger, stats, STATS_NO_SAMPLING, wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import { producer } from '@dydxprotocol-indexer/kafka';
import { Message } from 'kafkajs';

import config from '../config';

const queuedMessages: {[topic: string]: Message[]} = {};
const timeouts: {[topic: string]: NodeJS.Timeout } = {};
export const sizeStat: string = `${config.SERVICE_NAME}.flush_websocket.size`;
export const timingStat: string = `${config.SERVICE_NAME}.flush_websocket.timing`;
const sendMessagesTaskname: string = 'produce_websocket_messages';

/**
 * Sends all queued messages to their respective topics
 */
export async function flushAllQueues(): Promise<void> {
  await Promise.all(
    Object.keys(queuedMessages).map(async (topic: string): Promise<void> => {
      if (queuedMessages[topic] !== undefined) {
        return sendMessages(topic);
      }
    }),
  );
}

/**
 * Wrapper to send messages to a Kafka topic. Mesages are batched and sent on an interval
 * or when the batch reaches a configurable maximum size.
 * @param message
 * @param topic
 */
export function sendMessageWrapper(message: Message, topic: string): void {
  if (queuedMessages[topic] === undefined) {
    queuedMessages[topic] = [];
  }
  queuedMessages[topic].push(message);

  if (shouldFlush(topic)) {
    if (timeouts[topic] !== undefined) {
      clearTimeout(timeouts[topic]!);
    }
    wrapBackgroundTask(sendMessages(topic), true, sendMessagesTaskname);
  }

  if (timeouts[topic] === undefined) {
    timeouts[topic] = setTimeout(() => {
      wrapBackgroundTask(sendMessages(topic), true, sendMessagesTaskname);
    },
    config.FLUSH_KAFKA_MESSAGES_INTERVAL_MS,
    );
  }
}

/**
 * Checks if the messages queued for a topic should be sent
 * @param topic
 * @returns
 */
function shouldFlush(topic: string): boolean {
  return (queuedMessages[topic]?.length ?? 0) > config.MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC;
}

/**
 * Sends all queued messages for a topic. If an error occurs during sending, all the messages are
 * re-enqueued. Failures to send are expected to happen rarely and errors are logged / metrics are
 * kept for the number of failures to alert if they start to happen frequently enough to impact
 * performance.
 * If SEND_WEBSOCKET_MESSAGES is false, the messages are not sent and the queue is cleared.
 *
 * @param topic
 * @returns
 */
async function sendMessages(topic: string): Promise<void> {
  delete timeouts[topic];

  const messages: Message[] = queuedMessages[topic];
  if (messages === undefined || messages.length === 0) {
    stats.histogram(sizeStat, 0, STATS_NO_SAMPLING, { topic, success: 'true', instance: getInstanceId() });
    return;
  }
  queuedMessages[topic] = [];
  if (!config.SEND_WEBSOCKET_MESSAGES) return;

  const start: number = Date.now();
  let success: boolean = false;

  try {
    await producer.send({
      topic,
      messages,
    });
    success = true;
  } catch (error) {
    logger.error({
      at: 'send-websocket-helper#sendMessages',
      message: 'Failed to send messages to Kafka',
      topic,
      error,
      numMessages: messages.length,
    });

    // Re-enqueue all messages if they failed to be sent
    messages.forEach((message: Message) => sendMessageWrapper(message, topic));
  } finally {
    const tags: {[name: string]: string} = {
      topic,
      success: success.toString(),
      instance: getInstanceId(),
    };
    stats.histogram(sizeStat, messages.length, STATS_NO_SAMPLING, tags);
    stats.timing(timingStat, Date.now() - start, STATS_NO_SAMPLING, tags);
  }
}
