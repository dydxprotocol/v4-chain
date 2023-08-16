import { getChannel, getMessageToForward } from '../../src/helpers/from-kafka-helpers';
import { InvalidForwardMessageError, InvalidTopicError } from '../../src/lib/errors';
import {
  CandlesChannelMessageToForward,
  Channel,
  MessageToForward,
  OrderbooksChannelMessageToForward,
  SubaccountsChannelMessageToForward,
  TradesChannelMessageToForward,
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
  CandleMessage, CandleMessage_Resolution,
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
  PROTO_TO_CANDLE_RESOLUTION,
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
      const messageToForward: SubaccountsChannelMessageToForward = getMessageToForward(
        Channel.V4_ACCOUNTS,
        message,
      ) as SubaccountsChannelMessageToForward;

      expect(messageToForward.channel).toEqual(Channel.V4_ACCOUNTS);
      expect(messageToForward.id).toEqual(`${defaultOwner}/${defaultAccNumber}`);
      expect(messageToForward.contents).toEqual(defaultContents);
      expect(messageToForward.blockHeight).toEqual(subaccountMessage.blockHeight);
      expect(messageToForward.transactionIndex).toEqual(subaccountMessage.transactionIndex);
      expect(messageToForward.eventIndex).toEqual(subaccountMessage.eventIndex);
    });

    it('gets correct MessageToForward for candles message', () => {
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(CandleMessage.encode(candlesMessage).finish())),
      );
      const messageToForward: CandlesChannelMessageToForward = getMessageToForward(
        Channel.V4_CANDLES,
        message,
      ) as CandlesChannelMessageToForward;

      expect(messageToForward.channel).toEqual(Channel.V4_CANDLES);
      expect(messageToForward.id).toEqual(`${btcTicker}/${CandleResolution.ONE_MINUTE}`);
      expect(messageToForward.contents).toEqual(defaultContents);
      expect(messageToForward.clobPairId).toEqual(candlesMessage.clobPairId);
      if (candlesMessage.resolution !== CandleMessage_Resolution.UNRECOGNIZED) {
        expect(messageToForward.resolution).toEqual(
          PROTO_TO_CANDLE_RESOLUTION[candlesMessage.resolution],
        );
      }
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
      const messageToForward: OrderbooksChannelMessageToForward = getMessageToForward(
        Channel.V4_ORDERBOOK,
        message,
      ) as OrderbooksChannelMessageToForward;

      expect(messageToForward.channel).toEqual(Channel.V4_ORDERBOOK);
      expect(messageToForward.id).toEqual(btcTicker);
      expect(messageToForward.contents).toEqual(defaultContents);
      expect(messageToForward.clobPairId).toEqual(orderbookMessage.clobPairId);
    });

    it('gets correct MessageToForward for trade message', () => {
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(TradeMessage.encode(tradesMessage).finish())),
      );
      const messageToForward: TradesChannelMessageToForward = getMessageToForward(
        Channel.V4_TRADES,
        message,
      ) as TradesChannelMessageToForward;

      expect(messageToForward.channel).toEqual(Channel.V4_TRADES);
      expect(messageToForward.id).toEqual(btcTicker);
      expect(messageToForward.contents).toEqual(defaultContents);
      expect(messageToForward.blockHeight).toEqual(tradesMessage.blockHeight);
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
