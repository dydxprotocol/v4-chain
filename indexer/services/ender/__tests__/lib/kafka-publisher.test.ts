import {
  KafkaTopics, producer, ProducerMessage, TRADES_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import {
  FillFromDatabase,
  FillTable,
  FillType,
  Liquidity,
  OrderFromDatabase,
  OrderSide,
  OrderStatus,
  SubaccountMessageContents,
  SubaccountTable,
  testConstants,
  TradeContent,
  TradeMessageContents,
  TradeType,
  TransferFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import {
  BlockHeightMessage, IndexerSubaccountId, SubaccountMessage, TradeMessage,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import { AnnotatedSubaccountMessage, ConsolidatedKafkaEvent, SingleTradeMessage } from '../../src/lib/types';

import { KafkaPublisher } from '../../src/lib/kafka-publisher';
import {
  defaultDateTime,
  defaultSubaccountId,
  defaultSubaccountMessage,
  defaultTradeContent,
  defaultTradeKafkaEvent,
  defaultTradeMessage,
  defaultWalletAddress,
} from '../helpers/constants';
import {
  contentToSingleTradeMessage,
  contentToTradeMessage,
  createConsolidatedKafkaEventFromSubaccount,
  createConsolidatedKafkaEventFromTrade,
} from '../helpers/kafka-publisher-helpers';
import {
  generateFillSubaccountMessage,
  generateOrderSubaccountMessage,
  generateTransferContents,
} from '../../src/helpers/kafka-helper';
import { DateTime } from 'luxon';
import { convertToSubaccountMessage } from '../../src/lib/helper';
import { defaultBlock } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

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
        key: Buffer.from(Uint8Array.from(
          IndexerSubaccountId.encode(defaultSubaccountId).finish(),
        )),
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

  it('successfully publishes block height messages', async () => {
    const message: BlockHeightMessage = {
      blockHeight: String(defaultBlock),
      version: '1.0.0',
      time: defaultDateTime.toString(),
    };
    const blockHeightEvent: ConsolidatedKafkaEvent = {
      topic: KafkaTopics.TO_WEBSOCKETS_BLOCK_HEIGHT,
      message,
    };

    const publisher: KafkaPublisher = new KafkaPublisher();
    publisher.addEvents([blockHeightEvent]);

    await publisher.publish();
    expect(producerSendMock).toHaveBeenCalledTimes(1);
    expect(producerSendMock).toHaveBeenCalledWith({
      topic: blockHeightEvent.topic,
      messages: [{
        value: Buffer.from(BlockHeightMessage.encode(blockHeightEvent.message).finish()),
      }],
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

      publisher.sortEvents(publisher.tradeMessages);
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

      publisher.sortEvents(publisher.subaccountMessages);
      expect(publisher.subaccountMessages).toEqual([beforeSubaccount, subaccount, afterSubaccount]);
    });
  });

  describe('aggregateFillEventsForSubaccountMessages', () => {
    const fill: FillFromDatabase = {
      id: FillTable.uuid(testConstants.defaultTendermintEventId, Liquidity.TAKER),
      subaccountId: testConstants.defaultSubaccountId,
      side: OrderSide.BUY,
      liquidity: Liquidity.TAKER,
      type: FillType.LIMIT,
      clobPairId: '1',
      orderId: testConstants.defaultOrderId,
      size: '10',
      price: '20000',
      quoteAmount: '200000',
      eventId: testConstants.defaultTendermintEventId,
      transactionHash: '', // TODO: Add a real transaction Hash
      createdAt: testConstants.createdDateTime.toISO(),
      createdAtHeight: testConstants.createdHeight,
      clientMetadata: '0',
      fee: '1.1',
      affiliateRevShare: '0',
    };
    const order: OrderFromDatabase = {
      ...testConstants.defaultOrderGoodTilBlockTime,
      id: testConstants.defaultOrderId,
    };

    const recipientSubaccountId: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
      owner: 'recipient',
      number: 1,
    });
    const deposit: TransferFromDatabase = {
      id: '',
      senderWalletAddress: defaultWalletAddress,
      recipientSubaccountId: SubaccountTable.uuid(
        recipientSubaccountId.owner,
        recipientSubaccountId.number,
      ),
      assetId: testConstants.defaultAsset.id,
      size: '10',
      eventId: testConstants.defaultTendermintEventId,
      transactionHash: 'hash',
      createdAt: DateTime.utc().toISO(),
      createdAtHeight: '1',
    };
    it('successfully aggregates all fill events per order id and sorts messages', async () => {
      const publisher: KafkaPublisher = new KafkaPublisher();

      // merged with message 3.
      const msg1Contents: SubaccountMessageContents = {
        fills: [
          generateFillSubaccountMessage(fill, 'BTC-USD'),
        ],
        orders: [
          generateOrderSubaccountMessage(order, 'BTC-USD'),
        ],
      };
      const message1: AnnotatedSubaccountMessage = {
        blockHeight: '1',
        transactionIndex: 1,
        eventIndex: 1,
        contents: JSON.stringify(msg1Contents),
        subaccountId: {
          owner: 'owner1',
          number: 0,
        },
        version: '1',
        orderId: 'order1',
        isFill: true,
        subaccountMessageContents: msg1Contents,
      };

      const msg2Contents: SubaccountMessageContents = {
        fills: [
          generateFillSubaccountMessage(fill, 'ETH-USD'),
        ],
      };
      const message2: AnnotatedSubaccountMessage = {
        ...message1,
        transactionIndex: 2,
        contents: JSON.stringify(msg2Contents),
        orderId: 'order2',
        subaccountMessageContents: msg2Contents,
      };

      const msg3Contents: SubaccountMessageContents = {
        fills: [
          generateFillSubaccountMessage({
            ...fill,
            size: '100',
          }, 'BTC-USD'),
        ],
        orders: [
          generateOrderSubaccountMessage({
            ...order,
            status: OrderStatus.FILLED,
          }, 'BTC-USD'),
        ],
      };
      const message3: AnnotatedSubaccountMessage = {
        ...message1,
        transactionIndex: 3,
        contents: JSON.stringify(msg3Contents),
        subaccountMessageContents: msg3Contents,
      };

      // non-fill subaccount message.
      const msg4Contents: SubaccountMessageContents = generateTransferContents(
        deposit,
        testConstants.defaultAsset,
        recipientSubaccountId,
        undefined,
        recipientSubaccountId,
      );
      const message4: AnnotatedSubaccountMessage = {
        ...message1,
        eventIndex: 4,
        orderId: undefined,
        isFill: undefined,
        contents: JSON.stringify(msg4Contents),
      };

      const expectedMergedContents: SubaccountMessageContents = {
        fills: [
          msg1Contents.fills![0],
          msg3Contents.fills![0],
        ],
        orders: [
          msg3Contents.orders![0],
        ],
      };
      const mergedMessage3: AnnotatedSubaccountMessage = {
        ...message3,
        contents: JSON.stringify(expectedMergedContents),
        subaccountMessageContents: expectedMergedContents,
      };

      publisher.addEvents([
        createConsolidatedKafkaEventFromSubaccount(message1),
        createConsolidatedKafkaEventFromSubaccount(message2),
        createConsolidatedKafkaEventFromSubaccount(message3),
        createConsolidatedKafkaEventFromSubaccount(message4),
      ]);

      publisher.aggregateFillEventsForSubaccountMessages();
      const expectedMsgs: SubaccountMessage[] = [
        convertToSubaccountMessage(message4),
        convertToSubaccountMessage(message2),
        convertToSubaccountMessage(mergedMessage3),
      ];
      expect(publisher.subaccountMessages).toEqual(expectedMsgs);

      await publisher.publish();

      expect(producerSendMock).toHaveBeenCalledTimes(1);
      expect(producerSendMock).toHaveBeenCalledWith({
        topic: KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
        messages: _.map(expectedMsgs, (message: SubaccountMessage) => {
          return {
            key: message.subaccountId !== undefined
              ? Buffer.from(Uint8Array.from(
                IndexerSubaccountId.encode(message.subaccountId).finish(),
              )) : undefined,
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
        type: TradeType.LIMIT,
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
        type: TradeType.LIMIT,
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
        type: TradeType.LIMIT,
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
