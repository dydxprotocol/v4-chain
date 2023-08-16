import { Wss } from '../../src/helpers/wss';
import { Subscriptions } from '../../src/lib/subscription';
import config from '../../src/config';
import {
  connect as connectToKafka,
  disconnect as disconnectFromKafka,
} from '../../src/helpers/kafka/kafka-controller';
import { Index } from '../../src/websocket';
import {
  producer,
  WebsocketTopics,
  kafka,
  startConsumer,
  TRADES_WEBSOCKET_MESSAGE_VERSION,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import { MessageForwarder } from '../../src/lib/message-forwarder';
import WebSocket from 'ws';
import {
  Channel,
  ChannelBatchDataMessage,
  ChannelDataMessage,
  IncomingMessageType,
  OutgoingMessage,
  OutgoingMessageType,
  SubscribedMessage,
  WebsocketEvents,
} from '../../src/types';
import { Admin } from 'kafkajs';
import { SubaccountMessage, TradeMessage } from '@dydxprotocol-indexer/v4-protos';
import { dbHelpers, testMocks, perpetualMarketRefresher } from '@dydxprotocol-indexer/postgres';
import {
  btcClobPairId,
  btcTicker,
  defaultSubaccountId,
  ethClobPairId,
  ethTicker,
} from '../constants';
import _ from 'lodash';
import { axiosRequest } from '../../src/lib/axios';

jest.mock('../../src/lib/axios');

describe('message-forwarder', () => {
  let wss: Wss;
  let subscriptions: Subscriptions;
  let index: Index;
  let WS_HOST: string;
  let admin: Admin;

  const baseTradeMessage: TradeMessage = {
    blockHeight: '1',
    contents: '{}',
    clobPairId: btcClobPairId,
    version: TRADES_WEBSOCKET_MESSAGE_VERSION,
  };

  const baseSubaccountMessage: SubaccountMessage = {
    blockHeight: '2',
    transactionIndex: 2,
    eventIndex: 2,
    contents: '{}',
    subaccountId: defaultSubaccountId,
    version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  };

  const btcTradesMessages: TradeMessage[] = [
    {
      ...baseTradeMessage,
      contents: JSON.stringify({ val: 1 }),
    },
    {
      ...baseTradeMessage,
      contents: JSON.stringify({ val: 2 }),
    },
    {
      ...baseTradeMessage,
      contents: JSON.stringify({ val: 3 }),
    },
    {
      ...baseTradeMessage,
      contents: JSON.stringify({ val: 4 }),
    },
    {
      ...baseTradeMessage,
      contents: JSON.stringify({ val: 5 }),
    },
    {
      ...baseTradeMessage,
      contents: JSON.stringify({ val: 6 }),
    },
  ];

  const ethTradesMessages: TradeMessage[] = [
    {
      ...baseTradeMessage,
      clobPairId: ethClobPairId,
      contents: JSON.stringify({ ethVal: 1 }),
    },
  ];

  const btcV2TradesMessages: TradeMessage[] = [
    {
      ...baseTradeMessage,
      version: '2.0.0',
      contents: JSON.stringify({ val: 1 }),
    },
    {
      ...baseTradeMessage,
      version: '2.0.0',
      contents: JSON.stringify({ val: 2 }),
    },
  ];

  const subaccountMessages: SubaccountMessage[] = [
    {
      ...baseSubaccountMessage,
      contents: JSON.stringify({ val: '1' }),
    },
    {
      ...baseSubaccountMessage,
      contents: JSON.stringify({ val: '2' }),
    },
  ];

  const mockAxiosResponse: Object = { a: 'b' };
  const subaccountInitialMessage: Object = {
    ...mockAxiosResponse,
    orders: mockAxiosResponse,
  };

  beforeAll(async () => {
    await dbHelpers.migrate();
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
    admin = kafka.admin();
    await Promise.all([
      connectToKafka(),
      producer.connect(),
      admin.connect(),
    ]);
    await startConsumer();
    await admin.fetchTopicMetadata();
    await admin.deleteTopicRecords({
      topic: WebsocketTopics.TO_WEBSOCKETS_TRADES,
      partitions: [{
        partition: 0,
        offset: '-1',
      }],
    });
  });

  afterAll(async () => {
    await Promise.all([
      disconnectFromKafka(),
      producer.disconnect(),
      admin.disconnect(),
      dbHelpers.clearData(),
    ]);
    await dbHelpers.teardown();
  });

  beforeEach(async () => {
    jest.clearAllMocks();

    config.WS_PORT += 1;
    WS_HOST = `ws://localhost:${config.WS_PORT}`;

    wss = new Wss();
    await wss.start();
    subscriptions = new Subscriptions();
    index = new Index(wss, subscriptions);
    (axiosRequest as jest.Mock).mockImplementation(() => (JSON.stringify(mockAxiosResponse)));
  });

  afterEach(() => {
    jest.clearAllMocks();
    jest.resetAllMocks();
  });

  it('Batch sends messages with different versions', (done: jest.DoneCallback) => {
    const channel: Channel = Channel.V4_TRADES;
    const id: string = btcTicker;

    const messageForwarder: MessageForwarder = new MessageForwarder(subscriptions, index);
    subscriptions.start(messageForwarder.forwardToClient);
    messageForwarder.start();

    const ws = new WebSocket(WS_HOST);
    let connectionId: string;

    ws.on(WebsocketEvents.MESSAGE, async (message) => {
      const msg: OutgoingMessage = JSON.parse(message.toString()) as OutgoingMessage;
      if (msg.message_id === 0) {
        connectionId = msg.connection_id;
      }

      if (msg.message_id === 1) {
        // Check that the initial message is correct.
        checkInitialMessage(
          msg as SubscribedMessage,
          connectionId,
          channel,
          id,
          mockAxiosResponse,
        );

        // Send both BTC and ETH trades messages interleaved
        // await each message to ensure they are sent in order
        for (const tradeMessage of _.concat(
          ethTradesMessages,
          btcTradesMessages,
          ethTradesMessages,
          btcV2TradesMessages,
        )) {
          await producer.send({
            topic: WebsocketTopics.TO_WEBSOCKETS_TRADES,
            messages: [{
              value: Buffer.from(Uint8Array.from(TradeMessage.encode(tradeMessage).finish())),
              partition: 0,
              timestamp: `${Date.now()}`,
            }],
          });
        }
      }

      if (msg.message_id >= 2) {
        const batchMsg: ChannelBatchDataMessage = JSON.parse(
          message.toString(),
        ) as ChannelBatchDataMessage;

        const versionToTradeMessages: _.Dictionary<TradeMessage[]> = _.chain(
          [btcV2TradesMessages, btcTradesMessages],
        )
          .flatten()
          .groupBy((tradeMessage) => tradeMessage.version)
          .value();

        checkVersionedBatchMessage(
          batchMsg,
          connectionId,
          channel,
          id,
          versionToTradeMessages as {string: any[]},
        );
        if (msg.message_id === 3) {
          done();
        }
      }
    });

    ws.on('open', () => {
      ws.send(JSON.stringify({
        type: IncomingMessageType.SUBSCRIBE,
        channel,
        id,
        batched: true,
      }));
    });
  });

  it('Batch sends subaccount messages', (done: jest.DoneCallback) => {
    const channel: Channel = Channel.V4_ACCOUNTS;
    const id: string = `${defaultSubaccountId.owner}/${defaultSubaccountId.number}`;

    const messageForwarder: MessageForwarder = new MessageForwarder(subscriptions, index);
    subscriptions.start(messageForwarder.forwardToClient);
    messageForwarder.start();

    const ws = new WebSocket(WS_HOST);
    let connectionId: string;

    ws.on(WebsocketEvents.MESSAGE, async (message) => {
      const msg: OutgoingMessage = JSON.parse(message.toString()) as OutgoingMessage;
      if (msg.message_id === 0) {
        connectionId = msg.connection_id;
      }

      if (msg.message_id === 1) {
        // Check that the initial message is correct.
        checkInitialMessage(
          msg as SubscribedMessage,
          connectionId,
          channel,
          id,
          subaccountInitialMessage,
        );

        // await each message to ensure they are sent in order
        for (const subaccountMessage of subaccountMessages) {
          await producer.send({
            topic: WebsocketTopics.TO_WEBSOCKETS_SUBACCOUNTS,
            messages: [{
              value: Buffer.from(
                Uint8Array.from(
                  SubaccountMessage.encode(subaccountMessage).finish(),
                ),
              ),
              partition: 0,
              timestamp: `${Date.now()}`,
            }],
          });
        }
      }

      if (msg.message_id >= 2) {
        const batchMsg: ChannelBatchDataMessage = JSON.parse(
          message.toString(),
        ) as ChannelBatchDataMessage;

        checkBatchMessage(
          batchMsg,
          connectionId,
          channel,
          id,
          SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
          subaccountMessages,
        );
        done();
      }
    });

    ws.on('open', () => {
      ws.send(JSON.stringify({
        type: IncomingMessageType.SUBSCRIBE,
        channel,
        id,
        batched: true,
      }));
    });
  });

  it('forwards messages', (done: jest.DoneCallback) => {
    const channel: Channel = Channel.V4_TRADES;
    const id: string = ethTicker;

    const messageForwarder: MessageForwarder = new MessageForwarder(subscriptions, index);
    subscriptions.start(messageForwarder.forwardToClient);
    messageForwarder.start();

    const ws = new WebSocket(WS_HOST);
    let connectionId: string;

    ws.on(WebsocketEvents.MESSAGE, async (message) => {
      const msg: OutgoingMessage = JSON.parse(message.toString()) as OutgoingMessage;
      if (msg.message_id === 0) {
        connectionId = msg.connection_id;
      }

      if (msg.message_id === 1) {
        // Check that the initial message is correct.
        checkInitialMessage(
          msg as SubscribedMessage,
          connectionId,
          channel,
          id,
          mockAxiosResponse,
        );

        // Send both BTC and ETH trades messages
        // await each message to ensure they are sent in order
        for (const tradeMessage of _.concat(
          ethTradesMessages,
          btcTradesMessages,
          ethTradesMessages,
        )) {
          await producer.send({
            topic: WebsocketTopics.TO_WEBSOCKETS_TRADES,
            messages: [{
              value: Buffer.from(Uint8Array.from(TradeMessage.encode(tradeMessage).finish())),
              partition: 0,
              timestamp: `${Date.now()}`,
            }],
          });
        }
      }

      if (msg.message_id >= 2) {
        const forwardedMsg: ChannelDataMessage = JSON.parse(
          message.toString(),
        ) as ChannelDataMessage;

        expect(forwardedMsg.connection_id).toBe(connectionId);
        expect(forwardedMsg.type).toBe(OutgoingMessageType.CHANNEL_DATA);
        expect(forwardedMsg.channel).toBe(channel);
        expect(forwardedMsg.id).toBe(id);
        // Should only receive ETH messages
        expect(forwardedMsg.contents).toEqual(JSON.parse(ethTradesMessages[0].contents));
        expect(forwardedMsg.version).toEqual(TRADES_WEBSOCKET_MESSAGE_VERSION);
        // Only 2 ETH messages should be sent
        if (msg.message_id === 3) {
          done();
        }
      }
    });

    ws.on('open', () => {
      ws.send(JSON.stringify({
        type: IncomingMessageType.SUBSCRIBE,
        channel,
        id,
        batched: false,
      }));
    });
  });
});

function checkInitialMessage(
  subscribedMessage: SubscribedMessage,
  connectionId: string,
  channel: string,
  id: string,
  initialMessage: Object,
): void {
  expect(subscribedMessage.connection_id).toBe(connectionId);
  expect(subscribedMessage.type).toBe(OutgoingMessageType.SUBSCRIBED);
  expect(subscribedMessage.channel).toBe(channel);
  expect(subscribedMessage.id).toBe(id);
  expect(subscribedMessage.contents).toEqual(initialMessage);
}

function checkBatchMessage(
  batchMsg: ChannelBatchDataMessage,
  connectionId: string,
  channel: string,
  id: string,
  version: string,
  expectedMessages: {contents: string}[],
): void {
  expect(batchMsg.connection_id).toBe(connectionId);
  expect(batchMsg.type).toBe(OutgoingMessageType.CHANNEL_BATCH_DATA);
  expect(batchMsg.channel).toBe(channel);
  expect(batchMsg.id).toBe(id);
  expect(batchMsg.contents.length).toBe(expectedMessages.length);
  expect(batchMsg.version).toBe(version);
  batchMsg.contents.forEach(
    (individualMessage: Object, idx: number) => {
      expect(individualMessage).toEqual(JSON.parse(expectedMessages[idx].contents));
    },
  );
}

function checkVersionedBatchMessage(
  batchMsg: ChannelBatchDataMessage,
  connectionId: string,
  channel: string,
  id: string,
  versionToMessages: {string: any[]},
): void {
  expect(batchMsg.connection_id).toBe(connectionId);
  expect(batchMsg.type).toBe(OutgoingMessageType.CHANNEL_BATCH_DATA);
  expect(batchMsg.channel).toBe(channel);
  expect(batchMsg.id).toBe(id);
  _.forEach(versionToMessages, (expectedMessages, version) => {
    if (batchMsg.version === version) {
      expect(batchMsg.contents.length).toBe(expectedMessages.length);
      batchMsg.contents.forEach(
        (individualMessage: Object, idx: number) => {
          expect(individualMessage).toEqual(JSON.parse(expectedMessages[idx].contents));
        },
      );
    }
  });
}
