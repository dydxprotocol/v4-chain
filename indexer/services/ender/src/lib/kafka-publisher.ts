import { stats, STATS_NO_SAMPLING } from '@dydxprotocol-indexer/base';
import {
  BatchKafkaProducer,
  KafkaTopics,
  producer,
  ProducerMessage,
  TRADES_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import { TradeMessageContents } from '@dydxprotocol-indexer/postgres';
import {
  CandleMessage, MarketMessage, OffChainUpdateV1, SubaccountMessage, TradeMessage,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';

import config from '../config';
import { ConsolidatedKafkaEvent, SingleTradeMessage, VulcanMessage } from './types';

type TopicKafkaMessages = {
  topic: KafkaTopics;
  messages: ProducerMessage[];
};

export class KafkaPublisher {
  subaccountMessages: SubaccountMessage[];
  tradeMessages: SingleTradeMessage[];
  marketMessages: MarketMessage[];
  candleMessages: CandleMessage[];
  vulcanMessages: VulcanMessage[];

  constructor() {
    this.subaccountMessages = [];
    this.tradeMessages = [];
    this.marketMessages = [];
    this.candleMessages = [];
    this.vulcanMessages = [];
  }

  public addEvents(events: ConsolidatedKafkaEvent[]) {
    _.forEach(events, (event: ConsolidatedKafkaEvent) => {
      this.addEvent(event);
    });
  }

  public addEvent(event: ConsolidatedKafkaEvent) {
    switch (event.topic) {
      case KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS:
        this.subaccountMessages.push(event.message);
        break;
      case KafkaTopics.TO_WEBSOCKETS_TRADES:
        this.tradeMessages.push(event.message);
        break;
      case KafkaTopics.TO_WEBSOCKETS_MARKETS:
        this.marketMessages.push(event.message);
        break;
      case KafkaTopics.TO_WEBSOCKETS_CANDLES:
        this.candleMessages.push(event.message);
        break;
      case KafkaTopics.TO_VULCAN:
        this.vulcanMessages.push(event.message);
        break;
      default:
        throw new Error('Invalid Topic');
    }
  }

  /**
   * Sort trade events by block height, transaction index, and event index in ascending order,
   * where the first trade event should be the earliest trade in the block.
   */
  public sortTradeEvents() {
    this.tradeMessages = this.tradeMessages.sort(
      (a: SingleTradeMessage, b: SingleTradeMessage) => {
        if (Big(a.blockHeight).lt(b.blockHeight)) {
          return -1;
        } else if (Big(a.blockHeight).gt(b.blockHeight)) {
          return 1;
        }

        if (a.transactionIndex < b.transactionIndex) {
          return -1;
        } else if (a.transactionIndex > b.transactionIndex) {
          return 1;
        }

        return a.eventIndex < b.eventIndex ? -1 : 1;
      },
    );
  }

  public async publish() {
    const allTopicKafkaMessages:
    TopicKafkaMessages[] = this.generateAllTopicKafkaMessages();

    await Promise.all(
      _.map(
        allTopicKafkaMessages,
        (topicKafkaMessages: TopicKafkaMessages) => {
          const batchProducer: BatchKafkaProducer = new BatchKafkaProducer(
            topicKafkaMessages.topic,
            producer,
          );
          for (const message of topicKafkaMessages.messages) {
            batchProducer.addMessageAndMaybeFlush(message);
          }
          return batchProducer.flush();
        },
      ),
    );
  }

  private generateAllTopicKafkaMessages(): TopicKafkaMessages[] {
    const allTopicKafkaMessages: TopicKafkaMessages[] = [];
    if (this.subaccountMessages.length > 0) {
      allTopicKafkaMessages.push({
        topic: KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
        messages: _.map(this.subaccountMessages, (message: SubaccountMessage) => {
          return {
            value: Buffer.from(Uint8Array.from(SubaccountMessage.encode(message).finish())),
          };
        }),
      });
    }

    if (this.tradeMessages.length > 0) {
      allTopicKafkaMessages.push({
        topic: KafkaTopics.TO_WEBSOCKETS_TRADES,
        messages: this.groupKafkaTradesByClobPairId(),
      });
    }

    if (this.marketMessages.length > 0) {
      allTopicKafkaMessages.push({
        topic: KafkaTopics.TO_WEBSOCKETS_MARKETS,
        messages: _.map(this.marketMessages, (message: MarketMessage) => {
          return {
            value: Buffer.from(Uint8Array.from(MarketMessage.encode(message).finish())),
          };
        }),
      });
    }

    if (this.candleMessages.length > 0) {
      allTopicKafkaMessages.push({
        topic: KafkaTopics.TO_WEBSOCKETS_CANDLES,
        messages: _.map(this.candleMessages, (message: CandleMessage) => {
          return {
            value: Buffer.from(Uint8Array.from(CandleMessage.encode(message).finish())),
          };
        }),
      });
    }

    if (this.vulcanMessages.length > 0) {
      allTopicKafkaMessages.push({
        topic: KafkaTopics.TO_VULCAN,
        messages: _.map(this.vulcanMessages, (message: VulcanMessage) => {
          return {
            key: message.key,
            value: Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(message.value).finish())),
          };
        }),
      });
    }
    return allTopicKafkaMessages;
  }

  /**
   * Groups all trade messages for the trades kafka channel by clobPairId.
   * Expects all trade messages to only contain a single trade, because each OrderFillEvent
   * is individually processed and only contains a single trade each.
   */
  public groupKafkaTradesByClobPairId(): ProducerMessage[] {
    const start: number = Date.now();
    const groupedTradesMessages: _.Dictionary<SingleTradeMessage[]> = _.groupBy(
      this.tradeMessages,
      'clobPairId',
    );

    const groupedMergedTradeMessage: _.Dictionary<TradeMessage> = _.mapValues(
      groupedTradesMessages,
      (clobSpecificTradeMessages: SingleTradeMessage[]) => {
        const tradeContents: TradeMessageContents = _.reduce(
          clobSpecificTradeMessages,
          (result: TradeMessageContents, currentTradeMessage: TradeMessage) => {
            const contents: TradeMessageContents = JSON.parse(currentTradeMessage.contents);
            // content.trades.length == 1 because each OrderFillEvent only has a single fill/trade
            result.trades.push(contents.trades[0]);
            return result;
          },
          { trades: [] },
        );

        return {
          blockHeight: clobSpecificTradeMessages[0].blockHeight,
          contents: JSON.stringify(tradeContents),
          clobPairId: clobSpecificTradeMessages[0].clobPairId,
          version: TRADES_WEBSOCKET_MESSAGE_VERSION,
        };
      },
    );

    const messages: ProducerMessage[] = _.chain(Object.values(groupedMergedTradeMessage))
      .map((tradeMessage: TradeMessage) => {
        return {
          value: Buffer.from(Uint8Array.from(TradeMessage.encode(tradeMessage).finish())),
        };
      })
      .value();
    stats.timing(
      `${config.SERVICE_NAME}.group_kafka_trades.timing`,
      Date.now() - start,
      STATS_NO_SAMPLING,
    );

    return messages;
  }
}
