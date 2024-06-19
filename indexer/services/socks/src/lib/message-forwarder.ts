import { logger, stats } from '@dydxprotocol-indexer/base';
import _ from 'lodash';

import config from '../config';
import { createChannelBatchDataMessage, createChannelDataMessage } from '../helpers/message';
import { sendMessage } from '../helpers/wss';
import {
  Channel, Connection, MessageToForward, SubscriptionInfo,
} from '../types';
import { Index } from '../websocket/index';
import { MAX_TIMEOUT_INTEGER } from './constants';
import { Subscriptions } from './subscription';

const BATCH_SEND_INTERVAL_MS: number = config.BATCH_SEND_INTERVAL_MS;
const BUFFER_KEY_SEPARATOR: string = ':';

type VersionedContents = {
  contents: string;
  version: string;
  subaccountNumber?: number;
};

export class MessageForwarder {
  private static instance: MessageForwarder;
  private subscriptions: Subscriptions;
  private index: Index;
  private messageBuffer: { [key: string]: VersionedContents[] };
  private batchSending: NodeJS.Timeout;

  private constructor(subscriptions: Subscriptions, index: Index) {
    this.subscriptions = subscriptions;
    this.index = index;
    this.messageBuffer = {};
    this.batchSending = setTimeout(() => {
    }, MAX_TIMEOUT_INTEGER);
  }

  public static getInstance(subscriptions: Subscriptions, index: Index): MessageForwarder {
    if (!MessageForwarder.instance) {
      MessageForwarder.instance = new MessageForwarder(subscriptions, index);
    }
    return MessageForwarder.instance;
  }

  public start(): void {
    this.batchSending = setInterval(
      () => {
        this.forwardBatchedMessages();
      },
      BATCH_SEND_INTERVAL_MS,
    );
  }

  public stop(): void {
    clearInterval(this.batchSending);
  }

  public forwardMessage(message: MessageToForward): void {
    stats.increment(
      `${config.SERVICE_NAME}.message_to_forward`,
      1,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
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
              this.forwardBatchedVersionedMessagesBySubaccountNumber(
                batchedMessages,
                batchedSubscriber,
                channel,
                id,
              );
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
        try {
          this.forwardToClientBatch(
            msgs,
            batchedSubscriber.connectionId,
            channel,
            id,
            version,
            subaccountNumber,
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
        message.subaccountNumber,
      ),
    );
    return 1;
  }
}
