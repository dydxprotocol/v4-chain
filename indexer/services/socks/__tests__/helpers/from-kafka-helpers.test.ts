import { getChannel, getMessageToForward } from '../../src/helpers/from-kafka-helpers';
import { InvalidForwardMessageError, InvalidTopicError } from '../../src/lib/errors';
import {
  Channel,
  MessageToForward,
  WebsocketTopics,
} from '../../src/types';
import {
  btcTicker,
  candlesMessage,
  defaultAccNumber,
  defaultContents,
  defaultOwner,
  invalidChannel,
  invalidClobPairId,
  invalidTopic,
  marketsMessage,
  orderbookMessage,
  subaccountMessage,
  tradesMessage,
} from '../constants';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage } from './kafka';
import {
  CandleMessage,
  MarketMessage,
  OrderbookMessage,
  SubaccountMessage,
  TradeMessage,
} from '@dydxprotocol-indexer/v4-protos';
import { V4_MARKETS_ID } from '../../src/lib/constants';
import {
  dbHelpers,
  testMocks,
  perpetualMarketRefresher,
  CandleResolution,
} from '@dydxprotocol-indexer/postgres';

describe('from-kafka-helpers', () => {
  describe('getChannel', () => {
    it.each([
      [WebsocketTopics.TO_WEBSOCKETS_CANDLES, Channel.V4_CANDLES],
      [WebsocketTopics.TO_WEBSOCKETS_MARKETS, Channel.V4_MARKETS],
      [WebsocketTopics.TO_WEBSOCKETS_ORDERBOOKS, Channel.V4_ORDERBOOK],
      [WebsocketTopics.TO_WEBSOCKETS_SUBACCOUNTS, Channel.V4_ACCOUNTS],
      [WebsocketTopics.TO_WEBSOCKETS_TRADES, Channel.V4_TRADES],
    ])('gets correct channel for topic %s', (topic: WebsocketTopics, channel: Channel) => {
      expect(getChannel(topic)).toEqual(channel);
    });

    it('throws InvalidTopicError for invalid topic', () => {
      expect(() => { getChannel(invalidTopic); }).toThrow(new InvalidTopicError(invalidTopic));
    });
  });

  describe('getMessageToForward', () => {
    beforeAll(async () => {
      await dbHelpers.migrate();
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
    });

    afterAll(async () => {
      await dbHelpers.clearData();
      await dbHelpers.teardown();
    });

    it('gets correct MessageToForward for subaccount message', () => {
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(SubaccountMessage.encode(subaccountMessage).finish())),
      );
      const messageToForward: MessageToForward = getMessageToForward(
        Channel.V4_ACCOUNTS,
        message,
      );

      expect(messageToForward.channel).toEqual(Channel.V4_ACCOUNTS);
      expect(messageToForward.id).toEqual(`${defaultOwner}/${defaultAccNumber}`);
      expect(messageToForward.contents).toEqual(defaultContents);
    });

    it('gets correct MessageToForward for candles message', () => {
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(CandleMessage.encode(candlesMessage).finish())),
      );
      const messageToForward: MessageToForward = getMessageToForward(
        Channel.V4_CANDLES,
        message,
      );

      expect(messageToForward.channel).toEqual(Channel.V4_CANDLES);
      expect(messageToForward.id).toEqual(`${btcTicker}/${CandleResolution.ONE_MINUTE}`);
      expect(messageToForward.contents).toEqual(defaultContents);
    });

    it('gets correct MessageToForward for market message', () => {
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(MarketMessage.encode(marketsMessage).finish())),
      );
      const messageToForward: MessageToForward = getMessageToForward(Channel.V4_MARKETS, message);

      expect(messageToForward.channel).toEqual(Channel.V4_MARKETS);
      expect(messageToForward.id).toEqual(V4_MARKETS_ID);
      expect(messageToForward.contents).toEqual(defaultContents);
    });

    it('gets correct MessageToForward for orderbook message', () => {
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(OrderbookMessage.encode(orderbookMessage).finish())),
      );
      const messageToForward: MessageToForward = getMessageToForward(
        Channel.V4_ORDERBOOK,
        message,
      );

      expect(messageToForward.channel).toEqual(Channel.V4_ORDERBOOK);
      expect(messageToForward.id).toEqual(btcTicker);
      expect(messageToForward.contents).toEqual(defaultContents);
    });

    it('gets correct MessageToForward for trade message', () => {
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(TradeMessage.encode(tradesMessage).finish())),
      );
      const messageToForward: MessageToForward = getMessageToForward(
        Channel.V4_TRADES,
        message,
      );

      expect(messageToForward.channel).toEqual(Channel.V4_TRADES);
      expect(messageToForward.id).toEqual(btcTicker);
      expect(messageToForward.contents).toEqual(defaultContents);
    });

    it('throws InvalidForwardMessageError for empty message', () => {
      const message: KafkaMessage = createKafkaMessage(null);

      expect(() => { getMessageToForward(Channel.V4_ACCOUNTS, message); }).toThrow(
        new InvalidForwardMessageError('Got empty kafka message'),
      );
    });

    it('throws InvalidForwardMessageError for invalid channel', () => {
      const message: KafkaMessage = createKafkaMessage(Buffer.from(''));

      expect(() => { getMessageToForward((invalidChannel as Channel), message); }).toThrow(
        new InvalidForwardMessageError(`Unknown channel: ${invalidChannel}`),
      );
    });

    it('throw InvalidForwardMessageError for invalid clobPairId on candles message', () => {
      const invalidCandlesMessage: CandleMessage = {
        ...candlesMessage,
        clobPairId: invalidClobPairId,
      };
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(CandleMessage.encode(invalidCandlesMessage).finish())),
      );

      expect(() => { getMessageToForward(Channel.V4_CANDLES, message); }).toThrow(
        new InvalidForwardMessageError(`Invalid clob pair id: ${invalidClobPairId}`),
      );
    });

    it('throw InvalidForwardMessageError for invalid clobPairId on orderbook message', () => {
      const invalidOrderbookMessage: OrderbookMessage = {
        ...orderbookMessage,
        clobPairId: invalidClobPairId,
      };
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(OrderbookMessage.encode(invalidOrderbookMessage).finish())),
      );

      expect(() => { getMessageToForward(Channel.V4_ORDERBOOK, message); }).toThrow(
        new InvalidForwardMessageError(`Invalid clob pair id: ${invalidClobPairId}`),
      );
    });

    it('throw InvalidForwardMessageError for invalid clobPairId on trade message', () => {
      const invalidTradeMessage: TradeMessage = {
        ...tradesMessage,
        clobPairId: invalidClobPairId,
      };
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(TradeMessage.encode(invalidTradeMessage).finish())),
      );

      expect(() => { getMessageToForward(Channel.V4_TRADES, message); }).toThrow(
        new InvalidForwardMessageError(`Invalid clob pair id: ${invalidClobPairId}`),
      );
    });
  });
});
