import { logger } from '@dydxprotocol-indexer/base';
import {
  perpetualMarketRefresher,
  PROTO_TO_CANDLE_RESOLUTION,
  parentSubaccountHelpers,
  SubaccountMessageContents,
} from '@dydxprotocol-indexer/postgres';
import { getParentSubaccountNum } from '@dydxprotocol-indexer/postgres/build/src/lib/parent-subaccount-helpers';
import {
  CandleMessage,
  CandleMessage_Resolution,
  MarketMessage,
  OrderbookMessage,
  TradeMessage,
  SubaccountMessage,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';

import { TOPIC_TO_CHANNEL, V4_MARKETS_ID } from '../lib/constants';
import { InvalidForwardMessageError, InvalidTopicError } from '../lib/errors';
import { Channel, MessageToForward, WebsocketTopics } from '../types';

export function getChannels(topic: string): Channel[] {
  if (!Object.values(WebsocketTopics)
    .some((topicName: string) => { return topicName === topic; })) {
    throw new InvalidTopicError(topic);
  }

  const topicEnum: WebsocketTopics = <WebsocketTopics> topic;
  return TOPIC_TO_CHANNEL[topicEnum];
}

// TODO: remove this function and fix all tests to use getMessagesToForward instead
export function getMessageToForward(
  channel: Channel,
  message: KafkaMessage,
): MessageToForward {
  if (!message || !message.value) {
    throw new InvalidForwardMessageError('Got empty kafka message');
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
    case Channel.V4_CANDLES: {
      const candleMessage: CandleMessage = CandleMessage.decode(messageBinary);
      if (candleMessage.resolution === CandleMessage_Resolution.UNRECOGNIZED) {
        throw new InvalidForwardMessageError(`Unrecognized candle resolution: ${candleMessage.resolution}`);
      }
      return {
        channel,
        id: getCandleMessageId(candleMessage),
        contents: JSON.parse(candleMessage.contents),
        version: candleMessage.version,
      };
    }
    case Channel.V4_MARKETS: {
      const marketMessage: MarketMessage = MarketMessage.decode(messageBinary);
      return {
        channel,
        id: V4_MARKETS_ID,
        contents: JSON.parse(marketMessage.contents),
        version: marketMessage.version,
      };
    }
    case Channel.V4_ORDERBOOK: {
      const orderbookMessage: OrderbookMessage = OrderbookMessage.decode(messageBinary);
      return {
        channel,
        id: getTickerOrThrow(orderbookMessage.clobPairId),
        contents: JSON.parse(orderbookMessage.contents),
        version: orderbookMessage.version,
      };
    }
    case Channel.V4_TRADES: {
      const tradeMessage: TradeMessage = TradeMessage.decode(messageBinary);
      return {
        channel,
        id: getTickerOrThrow(tradeMessage.clobPairId),
        contents: JSON.parse(tradeMessage.contents),
        version: tradeMessage.version,
      };
    }
    case Channel.V4_PARENT_ACCOUNTS: {
      const subaccountMessage: SubaccountMessage = SubaccountMessage.decode(messageBinary);
      return {
        channel,
        id: getParentSubaccountMessageId(subaccountMessage),
        subaccountNumber: subaccountMessage.subaccountId!.number,
        contents: getParentSubaccountContents(subaccountMessage),
        version: subaccountMessage.version,
      };
    }
    default:
      throw new InvalidForwardMessageError(`Unknown channel: ${channel}`);
  }
}

export function getMessagesToForward(topic: string, message: KafkaMessage): MessageToForward[] {
  if (!message || !message.value) {
    throw new InvalidForwardMessageError('Got empty kafka message');
  }
  const messageBinary: Uint8Array = new Uint8Array(message.value);

  switch (topic) {
    case WebsocketTopics.TO_WEBSOCKETS_CANDLES: {
      const candleMessage: CandleMessage = CandleMessage.decode(messageBinary);
      return [{
        channel: Channel.V4_CANDLES,
        id: getCandleMessageId(candleMessage),
        contents: JSON.parse(candleMessage.contents),
        version: candleMessage.version,
      }];
    }
    case WebsocketTopics.TO_WEBSOCKETS_MARKETS: {
      const marketMessage: MarketMessage = MarketMessage.decode(messageBinary);
      return [{
        channel: Channel.V4_MARKETS,
        id: V4_MARKETS_ID,
        contents: JSON.parse(marketMessage.contents),
        version: marketMessage.version,
      }];
    }
    case WebsocketTopics.TO_WEBSOCKETS_ORDERBOOKS: {
      const orderbookMessage: OrderbookMessage = OrderbookMessage.decode(messageBinary);
      return [{
        channel: Channel.V4_ORDERBOOK,
        id: getTickerOrThrow(orderbookMessage.clobPairId),
        contents: JSON.parse(orderbookMessage.contents),
        version: orderbookMessage.version,
      }];
    }
    case WebsocketTopics.TO_WEBSOCKETS_TRADES: {
      const tradeMessage: TradeMessage = TradeMessage.decode(messageBinary);
      return [{
        channel: Channel.V4_TRADES,
        id: getTickerOrThrow(tradeMessage.clobPairId),
        contents: JSON.parse(tradeMessage.contents),
        version: tradeMessage.version,
      }];
    }
    case WebsocketTopics.TO_WEBSOCKETS_SUBACCOUNTS: {
      const subaccountMessage: SubaccountMessage = SubaccountMessage.decode(messageBinary);
      return [{
        channel: Channel.V4_ACCOUNTS,
        id: getSubaccountMessageId(subaccountMessage),
        contents: JSON.parse(subaccountMessage.contents),
        version: subaccountMessage.version,
      },
      {
        channel: Channel.V4_PARENT_ACCOUNTS,
        id: getParentSubaccountMessageId(subaccountMessage),
        subaccountNumber: subaccountMessage.subaccountId!.number,
        contents: getParentSubaccountContents(subaccountMessage),
        version: subaccountMessage.version,
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

  const senderAddress: string = contents.transfers.sender.address;
  const recipientAddress: string = contents.transfers.recipient.address;

  if (senderAddress !== recipientAddress) {
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
