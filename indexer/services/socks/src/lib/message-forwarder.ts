import {
  stats,
  getInstanceId,
  logger,
  InfoObject,
} from '@dydxprotocol-indexer/base';
import { updateOnBatchFunction, updateOnMessageFunction } from '@dydxprotocol-indexer/kafka';
import {
  Batch,
  EachBatchPayload,
  KafkaMessage,
} from 'kafkajs';
import _ from 'lodash';

import config from '../config';
import {
  getChannels,
  getMessagesToForward,
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
  WebsocketTopic,
} from '../types';
import { Index } from '../websocket/index';
import { MAX_TIMEOUT_INTEGER } from './constants';
import { Subscriptions } from './subscription';

const BATCH_SEND_INTERVAL_MS: number = config.BATCH_SEND_INTERVAL_MS;
const BUFFER_KEY_SEPARATOR: string = ':';

type VersionedContents = {
  contents: string,
  version: string,
  subaccountNumber?: number,
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
    this.batchSending = setTimeout(() => { }, MAX_TIMEOUT_INTEGER);
  }

  public start(): void {
    if (this.started) {
      throw new Error('MessageForwarder already started');
    }

    if (config.BATCH_PROCESSING_ENABLED) {
      logger.info({
        at: 'consumers#connect',
        message: 'Batch processing enabled',
      });
      updateOnBatchFunction(async (payload: EachBatchPayload): Promise<void> => {
        return this.onBatch(payload);
      });
    } else {
      logger.info({
        at: 'consumers#connect',
        message: 'Batch processing disabled. Processing each message individually',
      });
      // Kafkajs requires the function passed into `eachMessage` be an async function.
      // eslint-disable-next-line @typescript-eslint/require-await
      updateOnMessageFunction(async (topic, message): Promise<void> => {
        return this.onMessage(topic, message);
      });
    }

    this.started = true;
    this.batchSending = setInterval(
      () => { this.forwardBatchedMessages(); },
      BATCH_SEND_INTERVAL_MS,
    );
  }

  public async onBatch(
    payload: EachBatchPayload,
  ): Promise<void> {
    const batch: Batch = payload.batch;
    const topic: string = batch.topic;
    const partition: string = batch.partition.toString();
    const metricTags: Record<string, string> = {
      topic,
      partition,
      instance: getInstanceId(),
    };
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
    stats.timing(
      'socks.batch_time_in_queue',
      batchTimeInQueue,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      metricTags,
    );

    let lastCommitTime: number = startTime;
    for (let i = 0; i < batch.messages.length; i++) {
      const message: KafkaMessage = batch.messages[i];
      await this.onMessage(batch.topic, message);

      // Commit every KAFKA_BATCH_PROCESSING_COMMIT_FREQUENCY_MS to reduce number of roundtrips, and
      // also prevent disconnecting from the broker due to inactivity.
      const now: number = Date.now();
      if (now - lastCommitTime > config.KAFKA_BATCH_PROCESSING_COMMIT_FREQUENCY_MS) {
        logger.info({
          at: 'on-batch#onBatch',
          message: 'Committing offsets and sending heart beat',
          ...batchInfo,
        });
        payload.resolveOffset(message.offset);
        await Promise.all([
          payload.heartbeat(),
          // commitOffsetsIfNecessary will respect autoCommitThreshold and will not commit if
          // fewer messages than the threshold have been processed since the last commit.
          payload.commitOffsetsIfNecessary(),
        ]);
        lastCommitTime = now;
      }
    }

    const batchProcessingTime: number = Date.now() - startTime;
    stats.timing(
      'socks.batch_processing_time',
      batchProcessingTime,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      metricTags,
    );
    stats.timing(
      'socks.batch_size',
      batch.messages.length,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      metricTags,
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

  public onMessage(topic: string, message: KafkaMessage): void {
    const start: number = Date.now();
    stats.timing(
      `${config.SERVICE_NAME}.message_time_in_queue`,
      start - Number(message.timestamp),
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      {
        instance: getInstanceId(),
        topic,
      },
    );

    const loggerAt: string = 'MessageForwarder#onMessage';
    const errProps: Partial<InfoObject> = {
      topic,
      offset: message.offset,
    };

    const channels: Channel[] = getChannels(topic as WebsocketTopic);
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
    for (const messageToForward of getMessagesToForward(topic, message)) {
      const startForwardMessage: number = Date.now();
      this.forwardMessage(messageToForward);
      const end: number = Date.now();
      stats.timing(
        `${config.SERVICE_NAME}.forward_message`,
        end - startForwardMessage,
        config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
        {
          instance: getInstanceId(),
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
            instance: getInstanceId(),
            topic,
            event_type: String(message.headers?.event_type),
          },
        );
      }
    }
  }

  public forwardMessage(message: MessageToForward): void {
    stats.increment(
      `${config.SERVICE_NAME}.message_to_forward`,
      1,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      {
        instance: getInstanceId(),
      },
    );

    if (!this.subscriptions.subscriptions[message.channel] &&
      !this.subscriptions.batchedSubscriptions[message.channel]) {
      return;
    }

    const id: string = message.id;
    let subscriptions: SubscriptionInfo[] = [];
    if (this.subscriptions.subscriptions[message.channel]) {
      subscriptions = this.subscriptions.subscriptions[message.channel][id] || [];
    }
    let forwardedToSubscribers: boolean = false;

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
        subaccountNumber: message.subaccountNumber,
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
        {
          instance: getInstanceId(),
        },
      );
      forwardedToSubscribers = true;
    }

    // Don't double count a message that has both batched subscribers and non-batched subscribers
    if (forwardedToSubscribers) {
      stats.increment(
        `${config.SERVICE_NAME}.forward_message_with_subscribers`,
        1,
        config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
        {
          instance: getInstanceId(),
        },
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
              try {
                this.forwardBatchedVersionedMessagesBySubaccountNumber(
                  batchedMessages,
                  batchedSubscriber,
                  channel,
                  id,
                );
              } catch (error) {
                // catch error outside of loop to stop forwarding messages
                logger.error({
                  at: 'message-forwarder#forwardBatchedMessages',
                  message: error.message,
                  connectionId: batchedSubscriber.connectionId,
                  error,
                });
                throw error;
              }
            },
          );
        }
      },
    );
    this.messageBuffer = {};
  }

  private forwardBatchedVersionedMessagesBySubaccountNumber(
    batchedMessages: VersionedContents[],
    batchedSubscriber: SubscriptionInfo,
    channel: Channel,
    id: string,
  ): void {
    const batchedVersionedMessages: _.Dictionary<VersionedContents[]> = _.groupBy(
      batchedMessages,
      (c) => c.version,
    );
    _.forEach(batchedVersionedMessages, (versionedMsgs, version) => {
      const batchedMessagesBySubaccountNumber: _.Dictionary<VersionedContents[]> = _.groupBy(
        versionedMsgs,
        (c) => c.subaccountNumber,
      );
      _.forEach(batchedMessagesBySubaccountNumber, (msgs, subaccountNumberKey) => {
        const subaccountNumber: number | undefined = Number.isNaN(Number(subaccountNumberKey))
          ? undefined
          : Number(subaccountNumberKey);
        this.forwardToClientBatch(
          msgs,
          batchedSubscriber.connectionId,
          channel,
          id,
          version,
          subaccountNumber,
        );
      });
    });
  }

  public forwardToClientBatch(
    batchedMessages: VersionedContents[],
    connectionId: string,
    channel: Channel,
    id: string,
    version: string,
    subaccountNumber?: number,
  ): void {
    const connection: Connection = this.index.connections[connectionId];
    if (!connection) {
      logger.info({
        at: 'message-forwarder#forwardToClientBatch',
        message: 'Attempted to forward batched messages, but connection did not exist',
        connectionId,
      });
      stats.increment(
        `${config.SERVICE_NAME}.forward_to_client_batch_error`,
        1,
        {
          instance: getInstanceId(),
        },
      );
      this.subscriptions.unsubscribe(connectionId, channel, id);
      return;
    }

    this.index.connections[connectionId].messageId += 1;
    stats.increment(
      `${config.SERVICE_NAME}.forward_to_client_batch_success`,
      1,
      {
        instance: getInstanceId(),
      },
    );
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
        subaccountNumber,
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
      id: string,
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
      stats.increment(
        `${config.SERVICE_NAME}.forward_to_client_error`,
        1,
        {
          instance: getInstanceId(),
        },
      );
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
        message.subaccountNumber,
      ),
    );
    return 1;
  }
}
