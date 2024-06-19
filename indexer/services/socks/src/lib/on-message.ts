import { InfoObject, logger, stats } from '@dydxprotocol-indexer/base';
import { updateOnMessageFunction } from '@dydxprotocol-indexer/kafka';
import { KafkaMessage } from 'kafkajs';

import config from '../config';
import { getChannels, getMessagesToForward } from '../helpers/from-kafka-helpers';
import { Channel } from '../types';
import { Index } from '../websocket';
import { MessageForwarder } from './message-forwarder';
import { Subscriptions } from './subscription';

export function start(
  subscriptions: Subscriptions, index: Index,
): void {
  // eslint-disable-next-line @typescript-eslint/require-await
  updateOnMessageFunction(async (topic, message): Promise<void> => {
    return onMessage(topic, message, subscriptions, index);
  });
  MessageForwarder.getInstance(subscriptions, index).start();
}

export function onMessage(
  topic: string,
  message: KafkaMessage,
  subscriptions: Subscriptions,
  index: Index,
): void {
  const startTime: number = Date.now();
  stats.timing(
    `${config.SERVICE_NAME}.message_time_in_queue`,
    startTime - Number(message.timestamp),
    config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
    {
      topic,
    },
  );

  const loggerAt: string = 'MessageForwarder#onMessage';
  const errProps: Partial<InfoObject> = {
    topic,
    offset: message.offset,
  };

  const channels: Channel[] = getChannels(topic);
  if (channels.length === 0) {
    logger.error({
      ...errProps,
      at: loggerAt,
      message: `Unknown kafka topic: ${topic}.`,
    });
    return;
  }
  errProps.channels = channels;

  // Decode the message based on the topic
  const messagesToForward = getMessagesToForward(topic, message);
  for (const messageToForward of messagesToForward) {
    const startForwardMessage: number = Date.now();
    MessageForwarder.getInstance(subscriptions, index).forwardMessage(messageToForward);
    const end: number = Date.now();
    stats.timing(
      `${config.SERVICE_NAME}.forward_message`,
      end - startForwardMessage,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      {
        topic,
        channel: String(messageToForward.channel),
      },
    );

    const originalMessageTimestamp = message.headers?.message_received_timestamp;
    if (originalMessageTimestamp !== undefined) {
      stats.timing(
        `${config.SERVICE_NAME}.message_time_since_received`,
        startForwardMessage - Number(originalMessageTimestamp),
        config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
        {
          topic,
          event_type: String(message.headers?.event_type),
        },
      );
    }
  }
}
