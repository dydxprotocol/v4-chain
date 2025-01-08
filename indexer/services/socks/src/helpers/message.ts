import {
  OutgoingMessageType,
  Channel,
  ErrorMessage,
  ConnectedMessage,
  ChannelDataMessage,
  PongMessage,
  UnsubscribedMessage,
  ChannelBatchDataMessage,
} from '../types';

export function createErrorMessage(
  message: string,
  connectionId: string,
  messageId: number,
  channel?: string,
  id?: string,
): ErrorMessage {
  return {
    type: OutgoingMessageType.ERROR,
    message,
    connection_id: connectionId,
    message_id: messageId,
    channel,
    id,
  };
}

export function createConnectedMessage(connectionId: string): ConnectedMessage {
  return {
    type: OutgoingMessageType.CONNECTED,
    connection_id: connectionId,
    message_id: 0,
  };
}

export function createChannelDataMessage(
  channel: Channel,
  connectionId: string,
  messageId: number,
  id: string,
  version: string,
  // eslint-disable-next-line  @typescript-eslint/no-explicit-any
  contents: any,
  subaccountNumber?: number,
): ChannelDataMessage {
  if (channel === Channel.V4_MARKETS) {
    return {
      type: OutgoingMessageType.CHANNEL_DATA,
      connection_id: connectionId,
      message_id: messageId,
      channel,
      version,
      contents,
    };
  }

  return {
    type: OutgoingMessageType.CHANNEL_DATA,
    connection_id: connectionId,
    message_id: messageId,
    id,
    channel,
    version,
    contents,
    subaccountNumber,
  };
}

export function createChannelBatchDataMessage(
  channel: Channel,
  connectionId: string,
  messageId: number,
  id: string,
  version: string,
  // eslint-disable-next-line  @typescript-eslint/no-explicit-any
  contents: any[],
  subaccountNumber?: number,
): ChannelBatchDataMessage {
  if (channel === Channel.V4_MARKETS) {
    return {
      type: OutgoingMessageType.CHANNEL_BATCH_DATA,
      connection_id: connectionId,
      message_id: messageId,
      channel,
      version,
      contents,
    };
  }

  return {
    type: OutgoingMessageType.CHANNEL_BATCH_DATA,
    connection_id: connectionId,
    message_id: messageId,
    id,
    channel,
    version,
    contents,
    subaccountNumber,
  };
}

export function createPongMessage(
  connectionId: string,
  messageId: number,
  id?: number,
): PongMessage {
  const message: PongMessage = {
    type: OutgoingMessageType.PONG,
    connection_id: connectionId,
    message_id: messageId,
    id,
  };
  if (id) {
    message.id = id;
  }

  return message;
}

export function createUnsubscribedMessage(
  connectionId: string,
  messageId: number,
  channel: Channel,
  id?: string,
): UnsubscribedMessage {
  const message: UnsubscribedMessage = {
    type: OutgoingMessageType.UNSUBSCRIBED,
    connection_id: connectionId,
    message_id: messageId,
    channel,
  };
  if (channel !== Channel.V4_MARKETS) {
    message.id = id;
  }

  return message;
}

// TODO(DEC-248): See if the subscribed string message can be factored out/simplified.
export function createSubscribedMessage(
  channel: Channel,
  id: string,
  contents: string,
  connectionId: string,
  messageId: number,
): string {
  if (channel === Channel.V4_MARKETS) {
    return `{"type":"${OutgoingMessageType.SUBSCRIBED}","connection_id":"${connectionId}",` +
    `"message_id":${messageId},"channel":"${channel}","contents":${contents}}`;
  }

  return `{"type":"${OutgoingMessageType.SUBSCRIBED}","connection_id":"${connectionId}",` +
    `"message_id":${messageId},"channel":"${channel}","id":"${id}","contents":${contents}}`;
}
