import {
  stats,
  logger,
  InfoObject,
  safeJsonStringify,
} from '@dydxprotocol-indexer/base';
import { updateOnMessageFunction, updateOnBatchFunction } from '@dydxprotocol-indexer/kafka';
import { Batch, EachBatchPayload, KafkaMessage } from 'kafkajs';
import _ from 'lodash';

import config from '../config';
import {
  getChannel,
  getMessageToForward,
} from '../helpers/from-kafka-helpers';
import {
  createChannelDataMessage,
  createChannelBatchDataMessage,
} from '../helpers/message';
import { sendMessage } from '../helpers/wss';
import {
  MessageToForward,
  Channel,
  SubscriptionInfo,
  Connection,
} from '../types';
import { Index } from '../websocket/index';
import { MAX_TIMEOUT_INTEGER } from './constants';
import { Subscriptions } from './subscription';

const BATCH_SEND_INTERVAL_MS: number = config.BATCH_SEND_INTERVAL_MS;
const BUFFER_KEY_SEPARATOR: string = ':';

type VersionedContents = {
  contents: string;
  version: string;
};

export class MessageForwarder {
  private subscriptions: Subscriptions;
  private index: Index;
  private started: boolean;
  private stopped: boolean;
  private messageBuffer: { [key: string]: VersionedContents[] };
  private batchSending: NodeJS.Timeout;

  constructor(
    subscriptions: Subscriptions,
    index: Index,
  ) {
    this.subscriptions = subscriptions;
    this.index = index;
    this.started = false;
    this.stopped = false;
    this.messageBuffer = {};
    this.batchSending = setTimeout(() => {}, MAX_TIMEOUT_INTEGER);
  }

  public start(): void {
    if (this.started) {
      throw new Error('MessageForwarder already started');
    }

    // Kafkajs requires the function passed into `eachMessage` be an async function.
    if (config.BATCH_PROCESSING_ENABLED) {
      logger.info({
        at: 'consumers#connect',
        message: 'Batch processing enabled',
      });
      // eslint-disable-next-line @typescript-eslint/require-await
      updateOnBatchFunction(async (payload: EachBatchPayload) => this.onBatch(payload));
    } else {
      logger.info({
        at: 'consumers#connect',
        message: 'Batch processing disabled. Processing each message individually',
      });
      updateOnMessageFunction(async (
        topic: string,
        message: KafkaMessage,
      // eslint-disable-next-line @typescript-eslint/require-await
      ): Promise<void> => this.onMessage(topic, message));
    }

    this.started = true;
    this.batchSending = setInterval(
      () => { this.forwardBatchedMessages(); },
      BATCH_SEND_INTERVAL_MS,
    );
  }

  public stop(): void {
    if (this.stopped) {
      throw new Error('MessageForwarder already stopped');
    }
    if (!this.started) {
      throw new Error('MessageForwarder not started');
    }
    clearInterval(this.batchSending);
  }

  public onBatch(
    payload: EachBatchPayload,
  ): void {
    const batch: Batch = payload.batch;
    const topic: string = batch.topic;
    const partition: string = batch.partition.toString();
    const metricTags: Record<string, string> = { topic, partition };
    if (batch.isEmpty()) {
      logger.error({
        at: 'on-batch#onBatch',
        message: 'Empty batch',
        ...metricTags,
      });
      return;
    }

    const startTime: number = Date.now();
    const firstMessageTimestamp: number = Number(batch.messages[0].timestamp);
    const batchTimeInQueue: number = startTime - firstMessageTimestamp;
    const batchInfo = {
      firstMessageTimestamp: new Date(firstMessageTimestamp).toISOString(),
      batchTimeInQueue,
      messagesInBatch: batch.messages.length,
      firstOffset: batch.firstOffset(),
      lastOffset: batch.lastOffset(),
      ...metricTags,
    };

    logger.info({
      at: 'on-batch#onBatch',
      message: 'Received batch',
      ...batchInfo,
    });
    stats.timing(
      'socks.batch_time_in_queue',
      batchTimeInQueue,
      metricTags,
    );

    for (let i = 0; i < batch.messages.length; i++) {
      const message: KafkaMessage = batch.messages[i];
      this.onMessage(topic, message);
    }

    const batchProcessingTime: number = Date.now() - startTime;
    logger.info({
      at: 'on-batch#onBatch',
      message: 'Finished Processing Batch',
      batchProcessingTime,
      ...batchInfo,
    });
    stats.timing(
      'socks.batch_processing_time',
      batchProcessingTime,
      metricTags,
    );
    stats.timing(
      'socks.batch_size',
      batch.messages.length,
      metricTags,
    );
  }

  public onMessage(topic: string, message: KafkaMessage): void {
    stats.timing(
      `${config.SERVICE_NAME}.message_time_in_queue`,
      Date.now() - Number(message.timestamp),
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

    const channel: Channel | undefined = getChannel(topic);
    if (channel === undefined) {
      logger.error({
        ...errProps,
        at: loggerAt,
        message: `Unknown kafka topic: ${topic}.`,
      });
      return;
    }
    errProps.channel = channel;

    let messageToForward: MessageToForward;
    try {
      messageToForward = getMessageToForward(channel, message);
    } catch (error) {
      logger.error({
        ...errProps,
        at: loggerAt,
        message: 'Failed to get message to forward from kafka message',
        kafkaMessage: safeJsonStringify(message),
        error,
      });
      return;
    }

    const start: number = Date.now();
    this.forwardMessage(messageToForward);
    const end: number = Date.now();
    stats.timing(
      `${config.SERVICE_NAME}.forward_message`,
      end - start,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      {
        topic,
        channel: String(channel),
      },
    );
  }

  public forwardMessage(message: MessageToForward): void {
    stats.increment(
      `${config.SERVICE_NAME}.message_to_forward`,
      1,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
    );

    if (!this.subscriptions.subscriptions[message.channel] &&
      !this.subscriptions.batchedSubscriptions[message.channel]) {
      logger.info({
        at: 'message-forwarder#forwardMessage',
        message: 'No clients to forward to',
        messageId: message.id,
        messageChannel: message.channel,
        contents: message.contents,
      });
      return;
    }

    const id: string = message.id;
    let subscriptions: SubscriptionInfo[] = [];
    if (this.subscriptions.subscriptions[message.channel]) {
      subscriptions = this.subscriptions.subscriptions[message.channel][id] || [];
    }
    let forwardedToSubscribers: boolean = false;

    if (subscriptions.length > 0) {
      if (message.channel !== Channel.V4_ORDERBOOK ||
          (
            // Don't log orderbook messages unless enabled
            message.channel === Channel.V4_ORDERBOOK && config.ENABLE_ORDERBOOK_LOGS
          )
      ) {
        logger.info({
          at: 'message-forwarder#forwardMessage',
          message: 'Forwarding message to clients..',
          messageContents: message,
          connectionIds: subscriptions.map((s: SubscriptionInfo) => s.connectionId),
        });
      }
    }

    // Buffer messages if the subscription is for batched messages
    if (this.subscriptions.batchedSubscriptions[message.channel] &&
       this.subscriptions.batchedSubscriptions[message.channel][message.id]) {
      const bufferKey: string = this.getMessageBufferKey(
        message.channel,
        message.id,
      );
      if (!this.messageBuffer[bufferKey]) {
        this.messageBuffer[bufferKey] = [];
      }
      this.messageBuffer[bufferKey].push({
        contents: message.contents,
        version: message.version,
      } as VersionedContents);
      forwardedToSubscribers = true;
    }

    // Send message to client if the subscription is not batched
    if (subscriptions.length > 0) {
      let numClientsForwarded: number = 0;
      subscriptions.forEach(
        (subscription: SubscriptionInfo) => {
          if (subscription.pending) {
            subscription.pendingMessages.push(message);
            return;
          }
          numClientsForwarded += this.forwardToClient(message, subscription.connectionId);
        },
      );
      stats.increment(
        `${config.SERVICE_NAME}.forward_to_client_success`,
        numClientsForwarded,
        config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      );
      forwardedToSubscribers = true;
    }

    // Don't double count a message that has both batched subscribers and non-batched subscribers
    if (forwardedToSubscribers) {
      stats.increment(
        `${config.SERVICE_NAME}.forward_message_with_subscribers`,
        1,
      );
    }
  }

  public forwardBatchedMessages(): void {
    const bufferKeys: string[] = Object.keys(this.messageBuffer);
    bufferKeys.forEach(
      (bufferKey: string) => {
        const batchedMessages: VersionedContents[] = this.messageBuffer[bufferKey];
        if (batchedMessages.length > 0) {
          const {
            channel,
            channelString,
            id,
          } = this.parseMessageBufferKey(bufferKey);
          if (!this.subscriptions.batchedSubscriptions[channelString]) {
            return;
          }
          const batchedSubscribers: SubscriptionInfo[] = this
            .subscriptions
            .batchedSubscriptions[channelString][id];
          batchedSubscribers.forEach(
            (batchedSubscriber: SubscriptionInfo) => {
              const batchedVersionedMessages: _.Dictionary<VersionedContents[]> = _.groupBy(
                batchedMessages,
                (c) => c.version,
              );
              _.forEach(batchedVersionedMessages, (msgs, version) => {
                try {
                  this.forwardToClientBatch(
                    msgs,
                    batchedSubscriber.connectionId,
                    channel,
                    id,
                    version,
                  );
                } catch (error) {
                  logger.error({
                    at: 'message-forwarder#forwardBatchedMessages',
                    message: error.message,
                    connectionId: batchedSubscriber.connectionId,
                    error,
                  });
                }
              });
            },
          );
        }
      },
    );
    this.messageBuffer = {};
  }

  public forwardToClientBatch(
    batchedMessages: VersionedContents[],
    connectionId: string,
    channel: Channel,
    id: string,
    version: string,
  ): void {
    const connection: Connection = this.index.connections[connectionId];
    if (!connection) {
      logger.info({
        at: 'message-forwarder#forwardToClientBatch',
        message: 'Attempted to forward batched messages, but connection did not exist',
        connectionId,
      });
      stats.increment(`${config.SERVICE_NAME}.forward_to_client_batch_error`, 1);
      this.subscriptions.unsubscribe(connectionId, channel, id);
      return;
    }

    this.index.connections[connectionId].messageId += 1;
    stats.increment(`${config.SERVICE_NAME}.forward_to_client_batch_success`, 1);
    sendMessage(
      connection.ws,
      connectionId,
      createChannelBatchDataMessage(
        channel,
        connectionId,
        this.index.connections[connectionId].messageId,
        id,
        version,
        batchedMessages.map((c) => c.contents),
      ),
    );
  }

  private getMessageBufferKey(channel: Channel, id: string): string {
    return `${channel}${BUFFER_KEY_SEPARATOR}${id}`;
  }

  private parseMessageBufferKey(
    bufferKey: string,
  ): {
      channel: Channel,
      channelString: string,
      id: string
    } {
    const [channelString, id]: string[] = bufferKey.split(BUFFER_KEY_SEPARATOR);
    const channel: Channel = channelString as Channel;
    return {
      channel,
      channelString,
      id,
    };
  }

  public forwardToClient(message: MessageToForward, connectionId: string): number {
    const connection: Connection = this.index.connections[connectionId];
    if (!connection) {
      logger.info({
        at: 'message-forwarder#forwardToClient',
        message: 'Attempted to forward message, but connection did not exist',
        connectionId,
      });
      stats.increment(`${config.SERVICE_NAME}.forward_to_client_error`, 1);
      this.subscriptions.unsubscribe(connectionId, message.channel, message.id);
      return 0;
    }

    this.index.connections[connectionId].messageId += 1;

    sendMessage(
      connection.ws,
      connectionId,
      createChannelDataMessage(
        message.channel,
        connectionId,
        this.index.connections[connectionId].messageId,
        message.id,
        message.version,
        message.contents,
      ),
    );
    return 1;
  }
}
