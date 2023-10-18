import {
  KafkaTopics, producer, ProducerMessage, TRADES_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import { testConstants, TradeContent, TradeMessageContents } from '@dydxprotocol-indexer/postgres';
import { SubaccountMessage, TradeMessage } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import {
  AnnotatedSubaccountMessage,
  ConsolidatedKafkaEvent,
  convertToSubaccountMessage,
  SingleTradeMessage,
} from '../../src/lib/types';

import { KafkaPublisher } from '../../src/lib/kafka-publisher';
import {
  defaultSubaccountMessage,
  defaultTradeContent,
  defaultTradeKafkaEvent,
  defaultTradeMessage,
} from '../helpers/constants';
import {
  contentToSingleTradeMessage,
  contentToTradeMessage,
  createConsolidatedKafkaEventFromSubaccount,
  createConsolidatedKafkaEventFromTrade,
} from '../helpers/kafka-publisher-helpers';

describe('kafka-publisher', () => {
  let producerSendMock: jest.SpyInstance;
  const subaccountKafkaEvent: ConsolidatedKafkaEvent = {
    topic: KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
    message: defaultSubaccountMessage,
  };

  beforeEach(() => {
    producerSendMock = jest.spyOn(producer, 'send');
    producerSendMock.mockImplementation(() => {});
  });

  it('successfully publishes events', async () => {
    const publisher: KafkaPublisher = new KafkaPublisher();
    publisher.addEvents([subaccountKafkaEvent]);

    await publisher.publish();
    expect(producerSendMock).toHaveBeenCalledTimes(1);
    expect(producerSendMock).toHaveBeenCalledWith({
      topic: subaccountKafkaEvent.topic,
      messages: [{
        value: Buffer.from(SubaccountMessage.encode(subaccountKafkaEvent.message).finish()),
      }],
    });
  });

  it('successfuly publishes and groups trade events', async () => {
    const secondTradeContent: TradeContent = {
      ...defaultTradeContent,
      size: '11',
    };
    const secondTradeContents: TradeMessageContents = {
      trades: [secondTradeContent],
    };
    const secondTradeMessage: SingleTradeMessage = {
      ...defaultTradeMessage,
      contents: JSON.stringify(secondTradeContents),
    };
    const secondTradeKafkaEvent: ConsolidatedKafkaEvent = {
      topic: KafkaTopics.TO_WEBSOCKETS_TRADES,
      message: secondTradeMessage,
    };

    const publisher: KafkaPublisher = new KafkaPublisher();
    publisher.addEvents([defaultTradeKafkaEvent, secondTradeKafkaEvent]);

    await publisher.publish();

    const expectedTradeMessage: TradeMessage = {
      ...defaultTradeMessage,
      contents: JSON.stringify({
        trades: [defaultTradeContent, secondTradeContent],
      }),
    };
    expect(producerSendMock).toHaveBeenCalledTimes(1);
    expect(producerSendMock).toHaveBeenCalledWith({
      topic: defaultTradeKafkaEvent.topic,
      messages: [{ value: Buffer.from(TradeMessage.encode(expectedTradeMessage).finish()) }],
    });
  });

  describe('sortTradeEvents', () => {
    const trade: SingleTradeMessage = contentToSingleTradeMessage(
      {} as TradeContent,
      'BTC-USD',
    );
    const consolidatedTrade:
    ConsolidatedKafkaEvent = createConsolidatedKafkaEventFromTrade(trade);
    it.each([
      [
        'blockHeight',
        {
          ...trade,
          blockHeight: Big(trade.blockHeight).minus(1).toString(),
        },
        {
          ...trade,
          blockHeight: Big(trade.blockHeight).plus(1).toString(),
        },
      ],
      [
        'transactionIndex',
        {
          ...trade,
          transactionIndex: trade.transactionIndex - 1,
        },
        {
          ...trade,
          transactionIndex: trade.transactionIndex + 1,
        },
      ],
      [
        'eventIndex',
        {
          ...trade,
          eventIndex: trade.eventIndex - 1,
        },
        {
          ...trade,
          eventIndex: trade.eventIndex + 1,
        },
      ],
    ])('successfully trades events by %s', (
      _field: string,
      beforeTrade: SingleTradeMessage,
      afterTrade: SingleTradeMessage,
    ) => {
      const publisher: KafkaPublisher = new KafkaPublisher();
      const consolidatedBeforeTrade:
      ConsolidatedKafkaEvent = createConsolidatedKafkaEventFromTrade(
        beforeTrade,
      );
      const consolidatedAfterTrade:
      ConsolidatedKafkaEvent = createConsolidatedKafkaEventFromTrade(
        afterTrade,
      );

      publisher.addEvents([
        consolidatedAfterTrade,
        consolidatedTrade,
        consolidatedBeforeTrade,
      ]);

      publisher.sortEvents(KafkaTopics.TO_WEBSOCKETS_TRADES);
      expect(publisher.tradeMessages).toEqual([beforeTrade, trade, afterTrade]);
    });
  });

  describe('sortSubaccountEvents', () => {
    const subaccount: AnnotatedSubaccountMessage = defaultSubaccountMessage;
    const consolidatedSubaccount:
    ConsolidatedKafkaEvent = createConsolidatedKafkaEventFromSubaccount(subaccount);
    it.each([
      [
        'blockHeight',
        {
          ...subaccount,
          blockHeight: Big(subaccount.blockHeight).minus(1).toString(),
        },
        {
          ...subaccount,
          blockHeight: Big(subaccount.blockHeight).plus(1).toString(),
        },
      ],
      [
        'transactionIndex',
        {
          ...subaccount,
          transactionIndex: subaccount.transactionIndex - 1,
        },
        {
          ...subaccount,
          transactionIndex: subaccount.transactionIndex + 1,
        },
      ],
      [
        'eventIndex',
        {
          ...subaccount,
          eventIndex: subaccount.eventIndex - 1,
        },
        {
          ...subaccount,
          eventIndex: subaccount.eventIndex + 1,
        },
      ],
    ])('successfully subaccounts events by %s', (
      _field: string,
      beforeSubaccount: AnnotatedSubaccountMessage,
      afterSubaccount: AnnotatedSubaccountMessage,
    ) => {
      const publisher: KafkaPublisher = new KafkaPublisher();
      const consolidatedBeforeSubaccount:
      ConsolidatedKafkaEvent = createConsolidatedKafkaEventFromSubaccount(
        beforeSubaccount,
      );
      const consolidatedAfterSubaccount:
      ConsolidatedKafkaEvent = createConsolidatedKafkaEventFromSubaccount(
        afterSubaccount,
      );

      publisher.addEvents([
        consolidatedAfterSubaccount,
        consolidatedSubaccount,
        consolidatedBeforeSubaccount,
      ]);

      publisher.sortEvents(KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS);
      expect(publisher.subaccountMessages).toEqual([beforeSubaccount, subaccount, afterSubaccount]);
    });
  });

  describe('retainLastFillEventsForSubaccountMessages', () => {
    it('successfully retains the last fill event per order id and sorts messages', async () => {
      const publisher: KafkaPublisher = new KafkaPublisher();

      // over-written by message 3.
      const message1: AnnotatedSubaccountMessage = {
        blockHeight: '1',
        transactionIndex: 1,
        eventIndex: 1,
        contents: 'Message 1',
        subaccountId: {
          owner: 'owner1',
          number: 0,
        },
        version: '1',
        orderId: 'order1',
        isFill: true,
      };

      const message2: AnnotatedSubaccountMessage = {
        blockHeight: '1',
        transactionIndex: 2,
        eventIndex: 2,
        contents: 'Message 2',
        subaccountId: {
          owner: 'owner2',
          number: 0,
        },
        version: '1',
        orderId: 'order2',
        isFill: true,
      };

      const message3: AnnotatedSubaccountMessage = {
        blockHeight: '1',
        transactionIndex: 3,
        eventIndex: 3,
        contents: 'Message 3',
        subaccountId: {
          owner: 'owner3',
          number: 0,
        },
        version: '1',
        orderId: 'order1',
        isFill: true,
      };

      // non-fill subaccount message.
      const message4: AnnotatedSubaccountMessage = {
        blockHeight: '1',
        transactionIndex: 3,
        eventIndex: 4,
        contents: 'Message 3',
        subaccountId: {
          owner: 'owner3',
          number: 0,
        },
        version: '4',
      };

      publisher.addEvents([
        createConsolidatedKafkaEventFromSubaccount(message1),
        createConsolidatedKafkaEventFromSubaccount(message2),
        createConsolidatedKafkaEventFromSubaccount(message3),
        createConsolidatedKafkaEventFromSubaccount(message4),
      ]);

      publisher.retainLastFillEventsForSubaccountMessages();
      const expectedMsgs: SubaccountMessage[] = [
        convertToSubaccountMessage(message2),
        convertToSubaccountMessage(message3),
        convertToSubaccountMessage(message4),
      ];
      expect(publisher.subaccountMessages).toEqual(expectedMsgs);

      await publisher.publish();

      expect(producerSendMock).toHaveBeenCalledTimes(1);
      expect(producerSendMock).toHaveBeenCalledWith({
        topic: KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
        messages: _.map(expectedMsgs, (message: SubaccountMessage) => {
          return {
            value: Buffer.from(Uint8Array.from(SubaccountMessage.encode(message).finish())),
          };
        }),
      });
    });
  });

  describe('groupKafkaTradesByClobPairId', () => {
    it('successfully groups kafka trade messages', () => {
      const kafkaPublisher: KafkaPublisher = new KafkaPublisher();
      const btcClobPairId: string = 'BTC-USD';

      const tradeContent1: TradeContent = {
        id: 'trade1',
        size: '10',
        price: '10000',
        side: 'side',
        createdAt: 'today',
        liquidation: false,
      };
      const singleTrade1: SingleTradeMessage = contentToSingleTradeMessage(
        tradeContent1,
        btcClobPairId,
      );

      const tradeContent2: TradeContent = {
        id: 'trade2',
        size: '11',
        price: '12000',
        side: 'side',
        createdAt: 'today',
        liquidation: false,
      };
      const singleTrade2: SingleTradeMessage = contentToSingleTradeMessage(
        tradeContent2,
        btcClobPairId,
      );

      const ethClobPairId: string = 'ETH-USD';
      const tradeContent3: TradeContent = {
        id: 'trade3',
        size: '1',
        price: '1000',
        side: 'side',
        createdAt: 'today',
        liquidation: false,
      };
      const singleTrade3: SingleTradeMessage = contentToSingleTradeMessage(
        tradeContent3,
        ethClobPairId,
      );

      // Add all events
      _.forEach(
        [singleTrade1, singleTrade2, singleTrade3],
        (singleTradeMessage: SingleTradeMessage) => {
          kafkaPublisher.addEvents([
            createConsolidatedKafkaEventFromTrade(singleTradeMessage),
          ]);
        },
      );

      const groupedTrades: ProducerMessage[] = kafkaPublisher.groupKafkaTradesByClobPairId();
      expect(groupedTrades.length).toEqual(2);

      const tradeContents: TradeMessageContents = {
        trades: [tradeContent1, tradeContent2],
      };
      expect(groupedTrades).toContainEqual({
        value: Buffer.from(TradeMessage.encode(
          TradeMessage.fromPartial({
            blockHeight: testConstants.defaultBlock.blockHeight,
            contents: JSON.stringify(tradeContents),
            clobPairId: btcClobPairId,
            version: TRADES_WEBSOCKET_MESSAGE_VERSION,
          }),
        ).finish()),
      });
      const trade3: TradeMessage = contentToTradeMessage(tradeContent3, ethClobPairId);
      expect(groupedTrades).toContainEqual({
        value: Buffer.from(TradeMessage.encode(trade3).finish()),
      });
    });
  });
});
