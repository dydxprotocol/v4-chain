import {
  producer,
  WebsocketTopics,
  kafka,
  startConsumer,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import { dbHelpers, testMocks, perpetualMarketRefresher } from '@dydxprotocol-indexer/postgres';
import { SubaccountMessage } from '@dydxprotocol-indexer/v4-protos';
import { Admin } from 'kafkajs';
import WebSocket from 'ws';

import config from '../../src/config';
import {
  connect as connectToKafka,
  disconnect as disconnectFromKafka,
} from '../../src/helpers/kafka/kafka-controller';
import { Wss } from '../../src/helpers/wss';
import { axiosRequest } from '../../src/lib/axios';
import { MessageForwarder } from '../../src/lib/message-forwarder';
import { Subscriptions } from '../../src/lib/subscription';
import {
  Channel,
  ChannelBatchDataMessage,
  IncomingMessageType,
  OutgoingMessage,
  OutgoingMessageType,
  SubscribedMessage,
  WebsocketEvents,
} from '../../src/types';
import { Index } from '../../src/websocket';
import {
  defaultChildAccNumber,
  defaultChildSubaccountId,
  defaultSubaccountId,
} from '../constants';

jest.mock('../../src/lib/axios');

describe('message-forwarder', () => {
  let wss: Wss;
  let subscriptions: Subscriptions;
  let index: Index;
  let WS_HOST: string;
  let admin: Admin;

  const baseSubaccountMessage: SubaccountMessage = {
    blockHeight: '2',
    transactionIndex: 2,
    eventIndex: 2,
    contents: '{}',
    subaccountId: defaultSubaccountId,
    version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  };

  const childSubaccountMessage: SubaccountMessage = {
    blockHeight: '2',
    transactionIndex: 2,
    eventIndex: 2,
    contents: '{}',
    subaccountId: defaultChildSubaccountId,
    version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  };

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

  const childSubaccountMessages: SubaccountMessage[] = [
    {
      ...childSubaccountMessage,
      contents: JSON.stringify({ val: '1' }),
    },
    {
      ...childSubaccountMessage,
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

  it('Batch sends subaccount messages to parent subaccount channel', (done: jest.DoneCallback) => {
    const channel: Channel = Channel.V4_PARENT_ACCOUNTS;
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
        for (const subaccountMessage of childSubaccountMessages) {
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
          172,
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
  subaccountNumber?: number,
): void {
  expect(batchMsg.connection_id).toBe(connectionId);
  expect(batchMsg.type).toBe(OutgoingMessageType.CHANNEL_BATCH_DATA);
  expect(batchMsg.channel).toBe(channel);
  expect(batchMsg.id).toBe(id);
  expect(batchMsg.contents.length).toBe(expectedMessages.length);
  expect(batchMsg.version).toBe(version);
  expect(batchMsg.subaccountNumber).toBe(subaccountNumber);
  batchMsg.contents.forEach(
    (individualMessage: Object, idx: number) => {
      expect(individualMessage).toEqual(JSON.parse(expectedMessages[idx].contents));
    },
  );
}
