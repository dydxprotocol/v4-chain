import { stats, STATS_NO_SAMPLING } from '@dydxprotocol-indexer/base';
import {
  BatchKafkaProducer,
  KafkaTopics,
  producer,
  ProducerMessage,
  TRADES_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import { FillSubaccountMessageContents, TradeMessageContents } from '@dydxprotocol-indexer/postgres';
import {
  BlockHeightMessage,
  CandleMessage,
  IndexerSubaccountId,
  MarketMessage,
  OffChainUpdateV1,
  SubaccountMessage,
  TradeMessage,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';

import config from '../config';
import { convertToSubaccountMessage } from './helper';
import {
  AnnotatedSubaccountMessage,
  ConsolidatedKafkaEvent,
  SingleTradeMessage,
  VulcanMessage,
} from './types';

type TopicKafkaMessages = {
  topic: KafkaTopics,
  messages: ProducerMessage[],
};

type OrderedMessage = AnnotatedSubaccountMessage | SingleTradeMessage;

type Message = AnnotatedSubaccountMessage | SingleTradeMessage | MarketMessage |
CandleMessage | VulcanMessage | BlockHeightMessage;

export class KafkaPublisher {
  blockHeightMessages: BlockHeightMessage[];
  subaccountMessages: AnnotatedSubaccountMessage[];
  tradeMessages: SingleTradeMessage[];
  marketMessages: MarketMessage[];
  candleMessages: CandleMessage[];
  vulcanMessages: VulcanMessage[];

  constructor() {
    this.blockHeightMessages = [];
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
    this.getMessages(event.topic)!.push(event.message);
  }

  /**
   * Helper function to get messages for a given topic.
   *
   * @param kafkaTopic
   * @private
   */
  private getMessages(kafkaTopic: KafkaTopics): Message[] | undefined {
    switch (kafkaTopic) {
      case KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS:
        return this.subaccountMessages;
      case KafkaTopics.TO_WEBSOCKETS_TRADES:
        return this.tradeMessages;
      case KafkaTopics.TO_WEBSOCKETS_MARKETS:
        return this.marketMessages;
      case KafkaTopics.TO_WEBSOCKETS_CANDLES:
        return this.candleMessages;
      case KafkaTopics.TO_VULCAN:
        return this.vulcanMessages;
      case KafkaTopics.TO_WEBSOCKETS_BLOCK_HEIGHT:
        return this.blockHeightMessages;
      default:
        throw new Error('Invalid Topic');
    }
  }

  /**
   * Sort subaccountMessages that represent fills by block height, transaction index,
   * and event index in ascending order per order id. Only keep the subaccount message if
   * it represents the last fill event per order id, but make sure the subaccount message
   * contents contains all individual fills from that block.
   *
   * Due to separate handlers for order fills, we can be sure that if a message is annotated
   * with a fill, it should only contain data about a single fill / order and not transfers
   * or positions.
   */
  // TODO(IND-453): Generalize this to beyond subaccount messages.
  public aggregateFillEventsForSubaccountMessages() {
    // Create a map to store the last event for fills per order ID
    const lastEventForFills: Record<string, AnnotatedSubaccountMessage> = {};
    // Create a map to store all fill events per order ID
    const allFillEvents: Record<string, FillSubaccountMessageContents[]> = {};
    const nonFillEvents: AnnotatedSubaccountMessage[] = [];

    this.subaccountMessages.forEach((message: AnnotatedSubaccountMessage) => {
      if (message.isFill && message.orderId) {
        const fills:
        FillSubaccountMessageContents[] | undefined = message.subaccountMessageContents?.fills;
        const orderId: string = message.orderId;
        if (fills !== undefined) {
          allFillEvents[orderId] = allFillEvents[orderId]
            ? allFillEvents[orderId].concat(fills)
            : fills;
        }

        // If we haven't seen this order ID before or if the current message
        // was associated with a later event, update the lastFillEvents for this order ID
        if (
          !lastEventForFills[orderId] ||
          this.compareMessages(message, lastEventForFills[orderId]) > 0
        ) {
          lastEventForFills[orderId] = message;
        }
      } else {
        nonFillEvents.push(message);
      }
    });

    // Update the last event for the order ID such that it has all the fills
    // that occurred for the order ID.
    Object.keys(lastEventForFills).forEach((orderId: string) => {
      const lastEvent: AnnotatedSubaccountMessage = lastEventForFills[orderId];
      const fills: FillSubaccountMessageContents[] = allFillEvents[orderId];
      if (fills) {
        lastEvent.subaccountMessageContents!.fills = fills;
        lastEvent.contents = JSON.stringify(lastEvent.subaccountMessageContents);
      }
    });

    this.subaccountMessages = Object.values(lastEventForFills)
      .concat(nonFillEvents)
      .map((annotatedMessage) => convertToSubaccountMessage(annotatedMessage));
    this.sortEvents(this.subaccountMessages);
  }

  /** Helper function to compare two AnnotatedSubaccountMessages based on block height,
   * transaction index, and event index.
   */
  private compareMessages(a: AnnotatedSubaccountMessage, b: AnnotatedSubaccountMessage) {
    if (a.blockHeight === b.blockHeight) {
      if (a.transactionIndex === b.transactionIndex) {
        return a.eventIndex - b.eventIndex;
      }
      return a.transactionIndex - b.transactionIndex;
    }
    return Number(a.blockHeight) - Number(b.blockHeight);
  }

  /**
   * Sort events by block height, transaction index, and event index in ascending order,
   * where the first event should be the earliest event in the block.
   */
  public sortEvents(msgs: OrderedMessage[]) {
    msgs.sort((a: OrderedMessage, b: OrderedMessage) => {
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
    });
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
            config.KAFKA_MAX_BATCH_WEBSOCKET_MESSAGE_SIZE_BYTES,
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
    if (this.blockHeightMessages.length > 0) {
      allTopicKafkaMessages.push({
        topic: KafkaTopics.TO_WEBSOCKETS_BLOCK_HEIGHT,
        messages: _.map(this.blockHeightMessages, (message: BlockHeightMessage) => {
          return {
            value: Buffer.from(BlockHeightMessage.encode(message).finish()),
          };
        }),
      });
    }

    if (this.subaccountMessages.length > 0) {
      this.aggregateFillEventsForSubaccountMessages();

      allTopicKafkaMessages.push({
        topic: KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
        messages: _.map(this.subaccountMessages, (message: SubaccountMessage) => {
          return {
            key: message.subaccountId !== undefined
              ? Buffer.from(
                IndexerSubaccountId.encode(message.subaccountId).finish(),
              ) : undefined,
            value: Buffer.from(SubaccountMessage.encode(message).finish()),
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
            value: Buffer.from(MarketMessage.encode(message).finish()),
          };
        }),
      });
    }

    if (this.candleMessages.length > 0) {
      allTopicKafkaMessages.push({
        topic: KafkaTopics.TO_WEBSOCKETS_CANDLES,
        messages: _.map(this.candleMessages, (message: CandleMessage) => {
          return {
            value: Buffer.from(CandleMessage.encode(message).finish()),
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
            value: Buffer.from(OffChainUpdateV1.encode(message.value).finish()),
            headers: message.headers,
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
          value: Buffer.from(TradeMessage.encode(tradeMessage).finish()),
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
