import { SubaccountMessage } from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';

export interface MessageToForward {
  channel: Channel;
  id: string;
  // eslint-disable-next-line  @typescript-eslint/no-explicit-any
  contents: any;
  version: string;
}

export enum WebsocketTopics {
  TO_WEBSOCKETS_ORDERBOOKS = 'to-websockets-orderbooks',
  TO_WEBSOCKETS_SUBACCOUNTS = 'to-websockets-subaccounts',
  TO_WEBSOCKETS_TRADES = 'to-websockets-trades',
  TO_WEBSOCKETS_MARKETS = 'to-websockets-markets',
  TO_WEBSOCKETS_CANDLES = 'to-websockets-candles',
}

export enum Channel {
  V4_ORDERBOOK = 'v4_orderbook',
  V4_ACCOUNTS = 'v4_subaccounts',
  V4_TRADES = 'v4_trades',
  V4_MARKETS = 'v4_markets',
  V4_CANDLES = 'v4_candles',
}

export const TOPIC_TO_CHANNEL: Record<WebsocketTopics, Channel> = {
  [WebsocketTopics.TO_WEBSOCKETS_CANDLES]: Channel.V4_CANDLES,
  [WebsocketTopics.TO_WEBSOCKETS_MARKETS]: Channel.V4_MARKETS,
  [WebsocketTopics.TO_WEBSOCKETS_ORDERBOOKS]: Channel.V4_ORDERBOOK,
  [WebsocketTopics.TO_WEBSOCKETS_SUBACCOUNTS]: Channel.V4_ACCOUNTS,
  [WebsocketTopics.TO_WEBSOCKETS_TRADES]: Channel.V4_TRADES,
};

export function getChannel(topic: string): Channel | undefined {
  if (!Object.values(WebsocketTopics)
    .some((topicName: string) => {
      return topicName === topic;
    })) {
    throw new Error(`Invalid topic: ${topic}`);
  }

  const topicEnum: WebsocketTopics = <WebsocketTopics>topic;
  return TOPIC_TO_CHANNEL[topicEnum];
}

export function getMessageToForward(
  channel: Channel,
  message: KafkaMessage,
): MessageToForward {
  if (!message || !message.value) {
    throw new Error('Got empty kafka message');
  }

  const messageBinary: Uint8Array = new Uint8Array(message.value);
  switch (channel) {
    case Channel.V4_ACCOUNTS: {
      const subaccountMessage: SubaccountMessage = SubaccountMessage.decode(messageBinary);
      return {
        channel,
        id: getSubaccountMessageId(subaccountMessage),
        contents: JSON.parse(subaccountMessage.contents),
        version: subaccountMessage.version,
      };
    }
    default:
      throw new Error(`Unknown channel: ${channel}`);
  }
}

function getSubaccountMessageId(subaccountMessage: SubaccountMessage): string {
  return `${subaccountMessage.subaccountId!.owner}/${subaccountMessage.subaccountId!.number}`;
}
