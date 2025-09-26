import { GeoOriginHeaders, GeoOriginStatus } from '@dydxprotocol-indexer/compliance';
import crypto from 'node:crypto';
import { Index } from '../../src/websocket/index';
import WebSocket from 'ws';
import { Wss, sendMessage } from '../../src/helpers/wss';
import { WS_CLOSE_HEARTBEAT_TIMEOUT } from '../../src/lib/constants';
import { Subscriptions } from '../../src/lib/subscription';
import { IncomingMessage } from 'http';
import { Socket } from 'net';
import {
  IncomingMessageType,
  OutgoingMessageType,
  Channel,
  ALL_CHANNELS,
  WebsocketEvent,
  Connection,
} from '../../src/types';
import { InvalidMessageHandler } from '../../src/lib/invalid-message';

jest.mock('node:crypto');
jest.mock('../../src/helpers/wss');
jest.mock('../../src/lib/subscription');
jest.mock('../../src/lib/invalid-message');

describe('Index', () => {
  let index: Index;
  let mockWss: Wss;
  let websocket: WebSocket;
  let mockSub: Subscriptions;
  let mockConnect: (ws: WebSocket, req: IncomingMessage) => void;
  let wsCloseSpy: jest.SpyInstance;
  let wsOnSpy: jest.SpyInstance;
  let wsPingSpy: jest.SpyInstance;
  let wsTerminateSpy: jest.SpyInstance;
  let disconnectSpy: jest.SpyInstance;
  let invalidMsgHandlerSpy: jest.SpyInstance;

  const connectionId: string = 'conId';
  const geoOriginHeaders: GeoOriginHeaders = {
    'geo-origin-country': 'AR', // Argentina
    'geo-origin-region': 'AR-V', // Tierra del Fuego
    'geo-origin-status': GeoOriginStatus.OK,
  };

  beforeAll(() => {
    jest.useFakeTimers();
  });

  afterAll(() => {
    jest.resetAllMocks();
    jest.useRealTimers();
  });

  beforeEach(() => {
    jest.clearAllTimers();
    (Wss as unknown as jest.Mock).mockClear();
    (Subscriptions as unknown as jest.Mock).mockClear();
    (sendMessage as unknown as jest.Mock).mockClear();
    mockWss = new Wss();
    websocket = new WebSocket(null as any as string, [], { autoPong: true } as any);
    wsCloseSpy = jest.spyOn(websocket, 'close');
    wsOnSpy = jest.spyOn(websocket, 'on');
    wsPingSpy = jest.spyOn(websocket, 'ping').mockImplementation(jest.fn());
    wsTerminateSpy = jest.spyOn(websocket, 'terminate');
    mockWss.onConnection = jest.fn().mockImplementation(
      (cb: (ws: WebSocket, req: IncomingMessage) => void) => {
        mockConnect = cb;
      },
    );
    mockSub = new Subscriptions();
    invalidMsgHandlerSpy = jest.spyOn(InvalidMessageHandler.prototype, 'handleInvalidMessage');
    index = new Index(mockWss, mockSub);
    disconnectSpy = jest.spyOn(index, 'disconnect');
  });

  describe('connection', () => {
    it('adds connection to index, sends connection message, and attaches handlers', () => {
      (crypto.randomUUID as unknown as jest.Mock).mockReturnValueOnce(connectionId);
      mockConnect(websocket, new IncomingMessage(new Socket()));

      // Test that the connection is tracked.
      expect(index.connections[connectionId]).not.toBeUndefined();
      expect(index.connections[connectionId].ws).toEqual(websocket);
      expect(index.connections[connectionId].messageId).toEqual(0);

      // Test that handlers are attached.
      expect(wsOnSpy).toHaveBeenCalledTimes(4);
      expect(wsOnSpy).toHaveBeenCalledWith(WebsocketEvent.MESSAGE, expect.anything());
      expect(wsOnSpy).toHaveBeenCalledWith(WebsocketEvent.CLOSE, expect.anything());
      expect(wsOnSpy).toHaveBeenCalledWith(WebsocketEvent.ERROR, expect.anything());
      expect(wsOnSpy).toHaveBeenCalledWith(WebsocketEvent.PONG, expect.anything());

      // Test that a connection messages is sent.
      expect(sendMessage).toHaveBeenCalledTimes(1);
      expect(sendMessage).toHaveBeenCalledWith(
        websocket,
        connectionId,
        expect.objectContaining({
          type: OutgoingMessageType.CONNECTED,
        }),
      );
    });
  });

  describe('handlers', () => {
    beforeEach(() => {
      // Connect to the index before starting each test.
      (crypto.randomUUID as unknown as jest.Mock).mockReturnValueOnce(connectionId);
      const incomingMessage: IncomingMessage = new IncomingMessage(new Socket());
      incomingMessage.headers['geo-origin-country'] = geoOriginHeaders['geo-origin-country'];
      incomingMessage.headers['geo-origin-region'] = geoOriginHeaders['geo-origin-region'];
      incomingMessage.headers['geo-origin-status'] = geoOriginHeaders['geo-origin-status'];
      mockConnect(websocket, incomingMessage);
    });

    describe('message', () => {
      const subId: string = 'subId';
      const unparseable: string = '{';
      const invalidMessageEmpty: IncomingMessage = createIncomingMessage({});
      const invalidMessageInvalidType: IncomingMessage = createIncomingMessage({ type: 'i' });

      it.each([
        ['unparseable', unparseable, 'Invalid message: could not parse'],
        ['missing type', JSON.stringify(invalidMessageEmpty), 'Invalid message: type is required'],
        ['invalid type', JSON.stringify(invalidMessageInvalidType), 'Invalid message type: i'],
      ])('handles invalid message: %s', (_name: string, message: string, err: string) => {
        websocket.emit('message', message);

        expect(invalidMsgHandlerSpy).toHaveBeenCalledTimes(1);
        expect(invalidMsgHandlerSpy).toHaveBeenCalledWith(
          err,
          expect.objectContaining({
            ws: websocket,
            messageId: index.connections[connectionId].messageId,
          }),
          connectionId,
        );
      });

      // Nested parameterized test of invalid subscribe and unsubscribe message handling.
      for (const type of [IncomingMessageType.SUBSCRIBE, IncomingMessageType.UNSUBSCRIBE]) {
        it.each([
          [
            'missing channel and id',
            JSON.stringify(createIncomingMessage({ type })),
            'Invalid subscribe message: channel is required',
          ],
          [
            'invalid channel',
            JSON.stringify(createIncomingMessage(
              {
                type,
                channel: 'invalid',
              },
            )),
            'Invalid channel: invalid',
          ],
          [
            'missing id',
            JSON.stringify(createIncomingMessage(
              {
                type,
                channel: Channel.V4_ACCOUNTS,
              },
            )),
            'Invalid id: undefined',
          ],
        // eslint-disable-next-line  no-loop-func
        ])(`handles invalid ${type} message: %s`, (_name: string, message: string, err: string) => {
          websocket.emit(WebsocketEvent.MESSAGE, message);

          // Should be the second call, first call is to send the connected message.
          expect(sendMessage).toHaveBeenNthCalledWith(
            2,
            websocket,
            connectionId,
            {
              type: 'error',
              message: err,
              connection_id: connectionId,
              message_id: index.connections[connectionId].messageId,
            },
          );
        });
      }

      it.each(
        ALL_CHANNELS.map((channel: Channel) => { return [channel]; }),
      )('handles valid subscription message for channel: %s', (channel: Channel) => {
        // Test that markets work with a missing id.
        const id: string | undefined = (
          channel === Channel.V4_MARKETS || channel === Channel.V4_BLOCK_HEIGHT
        ) ? undefined : subId;
        const isBatched: boolean = false;
        const subMessage: IncomingMessage = createIncomingMessage({
          type: IncomingMessageType.SUBSCRIBE,
          channel,
          id,
          batched: isBatched,
        });
        websocket.emit(WebsocketEvent.MESSAGE, JSON.stringify(subMessage));

        expect(mockSub.subscribe).toHaveBeenCalledTimes(1);
        expect(mockSub.subscribe).toHaveBeenCalledWith(
          websocket,
          channel,
          connectionId,
          index.connections[connectionId].messageId,
          id,
          isBatched,
          geoOriginHeaders,
        );
      });

      it.each(
        ALL_CHANNELS.map((channel: Channel) => { return [channel]; }),
      )('handles valid unsubscribe message for channel: %s', (channel: Channel) => {
        // Test that markets work with a missing id.
        const id: string | undefined = (
          channel === Channel.V4_MARKETS || channel === Channel.V4_BLOCK_HEIGHT
        ) ? undefined : subId;
        const unSubMessage: IncomingMessage = createIncomingMessage({
          type: IncomingMessageType.UNSUBSCRIBE,
          channel,
          id,
        });
        websocket.emit(WebsocketEvent.MESSAGE, JSON.stringify(unSubMessage));

        expect(mockSub.unsubscribe).toHaveBeenCalledTimes(1);
        expect(mockSub.unsubscribe).toHaveBeenCalledWith(
          connectionId,
          channel,
          id,
        );
        expect(sendMessage).toHaveBeenNthCalledWith(
          2,
          websocket,
          connectionId,
          {
            channel,
            connection_id: connectionId,
            id,
            message_id: index.connections[connectionId].messageId,
            type: OutgoingMessageType.UNSUBSCRIBED,
          },
        );
      });
    });

    describe('close', () => {
      it('disconnects connection on close', () => {
        const connection: Connection = index.connections[connectionId];
        expect(disconnectSpy).not.toHaveBeenCalled();
        websocket.emit(WebsocketEvent.CLOSE);
        jest.runAllTimers();
        expect(disconnectSpy).toHaveBeenCalledTimes(1);
        expect(disconnectSpy).toHaveBeenCalledWith(connection);
      });

      it('handles reason as a Buffer object', () => {
        const dummyCode: number = 21;
        const bufferReason: Buffer = Buffer.from('bufferReason');
        jest.spyOn(websocket, 'terminate').mockImplementation(jest.fn());
        websocket.emit(WebsocketEvent.CLOSE, dummyCode, bufferReason);

        expect(wsTerminateSpy).toHaveBeenCalledTimes(1);
        expect(mockSub.remove).toHaveBeenCalledWith(connectionId);
        expect(index.connections[connectionId]).toBeUndefined();
      });
    });

    describe('pong', () => {
      it('removes delayed disconnect on pong', () => {
        // Run pending timers to start heartbeat to attach delayed disconnect.
        jest.runOnlyPendingTimers();
        jest.spyOn(websocket, 'terminate').mockImplementation(jest.fn());
        websocket.emit(WebsocketEvent.PONG);

        expect(index.connections[connectionId].disconnect).toBeUndefined();

        // Run pending timers to check connection wasn't disconnected on a timer.
        jest.runOnlyPendingTimers();
        expect(wsTerminateSpy).not.toHaveBeenCalled();
      });
    });

    describe('heartbeat', () => {
      it('sends heartbeat ping', () => {
        // Run pending timers to start heartbeat.
        jest.runOnlyPendingTimers();

        // Test that a heartbeat was sent.
        expect(wsPingSpy).toHaveBeenCalledTimes(1);
      });

      it('disconnects if pong isn\'t received', () => {
        // Run pending timers to start heartbeat to attach delayed disconnect.
        jest.runOnlyPendingTimers();

        expect(index.connections[connectionId].disconnect).not.toBeUndefined();

        // Run pending timers to check connection was disconnected on a timer.
        jest.runOnlyPendingTimers();
        expect(wsTerminateSpy).toHaveBeenCalledTimes(1);
        expect(wsCloseSpy).toHaveBeenCalledTimes(1);
        expect(wsCloseSpy).toHaveBeenCalledWith(WS_CLOSE_HEARTBEAT_TIMEOUT, 'Heartbeat timeout');
      });
    });
  });

  describe('close', () => {
    it('closes all connections, then closes server', async () => {
      jest.spyOn(websocket, 'close');
      mockConnect(websocket, new IncomingMessage(new Socket()));
      await index.close();

      expect(wsCloseSpy).toHaveBeenCalledTimes(1);
      expect(mockWss.close).toHaveBeenCalledTimes(1);
    });
  });
});

function createIncomingMessage(properties: any): IncomingMessage {
  const message: IncomingMessage = new IncomingMessage(new Socket());
  return {
    ...message,
    ...properties,
  };
}
