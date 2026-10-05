import { logger } from '@dydxprotocol-indexer/base';
import {
  parentSubaccountHelpers,
  perpetualMarketRefresher,
  PROTO_TO_CANDLE_RESOLUTION,
  SubaccountMessageContents,
} from '@dydxprotocol-indexer/postgres';
import { getParentSubaccountNum } from '@dydxprotocol-indexer/postgres/build/src/lib/parent-subaccount-helpers';
import {
  BlockHeightMessage,
  CandleMessage,
  CandleMessage_Resolution,
  createProtoReader,
  MarketMessage,
  OrderbookMessage,
  SubaccountMessage,
  TradeMessage,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';

import { TOPIC_TO_CHANNEL, V4_BLOCK_HEIGHT_ID, V4_MARKETS_ID } from '../lib/constants';
import { InvalidForwardMessageError, InvalidTopicError } from '../lib/errors';
import { Channel, MessageToForward, WebsocketTopic } from '../types';

export function getChannels(topic: WebsocketTopic): Channel[] {
  if (!Object.values(WebsocketTopic)
    .some((topicName: string) => { return topicName === topic; })) {
    throw new InvalidTopicError(topic);
  }

  const topicEnum: WebsocketTopic = <WebsocketTopic> topic;
  return TOPIC_TO_CHANNEL[topicEnum];
}

export type HasSubscribers = (channel: Channel, id: string) => boolean;

/**
 * Every socks task consumes every message but most have no subscriber on a given task, so the id
 * is resolved first and `contents` is only parsed for channels `hasSubscribers` accepts.
 */
export function getMessagesToForward(
  topic: string,
  message: KafkaMessage,
  hasSubscribers: HasSubscribers = () => true,
): MessageToForward[] {
  if (!message || !message.value) {
    throw new InvalidForwardMessageError('Got empty kafka message');
  }

  switch (topic) {
    case WebsocketTopic.TO_WEBSOCKETS_CANDLES: {
      const candleMessage: CandleMessage = CandleMessage.decode(createProtoReader(message.value));
      const id: string = getCandleMessageId(candleMessage);
      if (!hasSubscribers(Channel.V4_CANDLES, id)) {
        return [];
      }
      return [{
        channel: Channel.V4_CANDLES,
        id,
        contents: JSON.parse(candleMessage.contents),
        version: candleMessage.version,
      }];
    }
    case WebsocketTopic.TO_WEBSOCKETS_MARKETS: {
      const marketMessage: MarketMessage = MarketMessage.decode(createProtoReader(message.value));
      if (!hasSubscribers(Channel.V4_MARKETS, V4_MARKETS_ID)) {
        return [];
      }
      return [{
        channel: Channel.V4_MARKETS,
        id: V4_MARKETS_ID,
        contents: JSON.parse(marketMessage.contents),
        version: marketMessage.version,
      }];
    }
    case WebsocketTopic.TO_WEBSOCKETS_ORDERBOOKS: {
      const orderbookMessage: OrderbookMessage = OrderbookMessage.decode(
        createProtoReader(message.value),
      );
      const id: string = getTickerOrThrow(orderbookMessage.clobPairId);
      if (!hasSubscribers(Channel.V4_ORDERBOOK, id)) {
        return [];
      }
      return [{
        channel: Channel.V4_ORDERBOOK,
        id,
        contents: JSON.parse(orderbookMessage.contents),
        version: orderbookMessage.version,
      }];
    }
    case WebsocketTopic.TO_WEBSOCKETS_TRADES: {
      const tradeMessage: TradeMessage = TradeMessage.decode(createProtoReader(message.value));
      const id: string = getTickerOrThrow(tradeMessage.clobPairId);
      if (!hasSubscribers(Channel.V4_TRADES, id)) {
        return [];
      }
      return [{
        channel: Channel.V4_TRADES,
        id,
        contents: JSON.parse(tradeMessage.contents),
        version: tradeMessage.version,
      }];
    }
    case WebsocketTopic.TO_WEBSOCKETS_SUBACCOUNTS: {
      const subaccountMessage: SubaccountMessage = SubaccountMessage.decode(
        createProtoReader(message.value),
      );
      const messages: MessageToForward[] = [];
      const accountId: string = getSubaccountMessageId(subaccountMessage);
      if (hasSubscribers(Channel.V4_ACCOUNTS, accountId)) {
        messages.push({
          channel: Channel.V4_ACCOUNTS,
          id: accountId,
          contents: JSON.parse(subaccountMessage.contents),
          version: subaccountMessage.version,
        });
      }
      const parentAccountId: string = getParentSubaccountMessageId(subaccountMessage);
      if (hasSubscribers(Channel.V4_PARENT_ACCOUNTS, parentAccountId)) {
        messages.push({
          channel: Channel.V4_PARENT_ACCOUNTS,
          id: parentAccountId,
          subaccountNumber: subaccountMessage.subaccountId!.number,
          contents: getParentSubaccountContents(subaccountMessage),
          version: subaccountMessage.version,
        });
      }
      return messages;
    }
    case WebsocketTopic.TO_WEBSOCKETS_BLOCK_HEIGHT: {
      const blockHeightMessage: BlockHeightMessage = BlockHeightMessage.decode(
        createProtoReader(message.value),
      );
      if (!hasSubscribers(Channel.V4_BLOCK_HEIGHT, V4_BLOCK_HEIGHT_ID)) {
        return [];
      }
      return [{
        channel: Channel.V4_BLOCK_HEIGHT,
        id: V4_BLOCK_HEIGHT_ID,
        version: blockHeightMessage.version,
        contents: {
          blockHeight: blockHeightMessage.blockHeight,
          time: blockHeightMessage.time,
        },
      }];
    }
    default:
      throw new InvalidForwardMessageError(`Unknown topic: ${topic}`);
  }
}

function getTickerOrThrow(clobPairId: string): string {
  const ticker: string | undefined = perpetualMarketRefresher.getPerpetualMarketTicker(clobPairId);
  if (ticker === undefined) {
    throw new InvalidForwardMessageError(`Invalid clob pair id: ${clobPairId}`);
  }

  return ticker;
}

function getSubaccountMessageId(subaccountMessage: SubaccountMessage): string {
  return `${subaccountMessage.subaccountId!.owner}/${subaccountMessage.subaccountId!.number}`;
}

function getParentSubaccountMessageId(subaccountMessage: SubaccountMessage): string {
  const parentSubaccountNumber: number = parentSubaccountHelpers.getParentSubaccountNum(
    subaccountMessage.subaccountId!.number,
  );
  return `${subaccountMessage.subaccountId!.owner}/${parentSubaccountNumber}`;
}

function getCandleMessageId(candleMessage: CandleMessage): string {
  const ticker: string = getTickerOrThrow(candleMessage.clobPairId);
  if (candleMessage.resolution === CandleMessage_Resolution.UNRECOGNIZED) {
    // This should never happen, but in the off chance that it does, log an error and this message
    // should never be published
    logger.error({
      at: 'from-kafka-helpers#getCandleMessageId',
      message: 'Unrecognized candle resolution',
    });
    return `${ticker}/`;
  }
  return `${ticker}/${PROTO_TO_CANDLE_RESOLUTION[candleMessage.resolution]}`;
}

function getParentSubaccountContents(msg: SubaccountMessage): SubaccountMessageContents {
  // Filter out transfers between child subaccounts of the same parent subaccount.
  const contents: SubaccountMessageContents = JSON.parse(msg.contents) as SubaccountMessageContents;
  if (contents.transfers === undefined) {
    return contents;
  }

  if (contents.transfers.sender.address !== contents.transfers.recipient.address) {
    return contents;
  }

  const senderSubaccount: number | undefined = contents.transfers.sender.subaccountNumber;
  const recipientSubaccount: number | undefined = contents.transfers.recipient.subaccountNumber;

  if (senderSubaccount === undefined || recipientSubaccount === undefined) {
    return contents;
  }

  const senderParentSubaccountNumber: number = getParentSubaccountNum(senderSubaccount);
  const recipientParentSubaccountNumber: number = getParentSubaccountNum(recipientSubaccount);

  if (senderParentSubaccountNumber !== recipientParentSubaccountNumber) {
    return contents;
  }

  delete contents.transfers;
  return contents;
}
