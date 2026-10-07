import { IncomingMessage as IncomingMessageHttp } from 'http';
import { Socket } from 'net';

import { logger, stats } from '@dydxprotocol-indexer/base';
import WebSocket from 'ws';

import config from '../../src/config';
import {
  CompressionState,
  flushCompressionStats,
  getPerMessageDeflateOptions,
  prepareUpgradeRequest,
  recordMessageSent,
  trackConnection,
} from '../../src/helpers/ws-compression';
import { sendMessageString, Wss } from '../../src/helpers/wss';

const EXTENSIONS_HEADER: string = 'sec-websocket-extensions';
const OFFER: string = 'permessage-deflate; client_max_window_bits';

describe('ws-compression', () => {
  const defaultRolloutPercent: number = config.WS_COMPRESSION_ROLLOUT_PERCENT;

  function requestWithExtensions(extensions?: string): IncomingMessageHttp {
    const req: IncomingMessageHttp = new IncomingMessageHttp(new Socket());
    if (extensions !== undefined) {
      req.headers[EXTENSIONS_HEADER] = extensions;
    }
    return req;
  }

  function websocketWithExtensions(extensions: string): WebSocket {
    return { extensions } as WebSocket;
  }

  beforeEach(() => {
    jest.spyOn(stats, 'increment');
    // Discard totals left behind by other cases.
    flushCompressionStats();
    (stats.increment as jest.Mock).mockClear();
  });

  afterEach(() => {
    config.WS_COMPRESSION_ROLLOUT_PERCENT = defaultRolloutPercent;
    jest.restoreAllMocks();
  });

  describe('getPerMessageDeflateOptions', () => {
    it('disables the extension when the rollout is at 0', () => {
      config.WS_COMPRESSION_ROLLOUT_PERCENT = 0;
      expect(getPerMessageDeflateOptions()).toEqual(false);
    });

    it('enables the extension when the rollout is above 0', () => {
      config.WS_COMPRESSION_ROLLOUT_PERCENT = 1;
      const options: WebSocket.PerMessageDeflateOptions | false = getPerMessageDeflateOptions();

      expect(options).toEqual(expect.objectContaining({
        threshold: config.WS_COMPRESSION_THRESHOLD_BYTES,
      }));
      // An explicit `false` makes `ws` reject clients that ask for no server context takeover.
      expect(options).not.toHaveProperty('serverNoContextTakeover');
    });
  });

  describe('prepareUpgradeRequest', () => {
    it('removes the offer for a connection outside the rollout', () => {
      config.WS_COMPRESSION_ROLLOUT_PERCENT = 50;
      jest.spyOn(Math, 'random').mockReturnValue(0.5);
      const req: IncomingMessageHttp = requestWithExtensions(OFFER);

      prepareUpgradeRequest(req);

      expect(req.headers[EXTENSIONS_HEADER]).toBeUndefined();
    });

    it('keeps the offer for a connection inside the rollout', () => {
      config.WS_COMPRESSION_ROLLOUT_PERCENT = 50;
      jest.spyOn(Math, 'random').mockReturnValue(0.49);
      const req: IncomingMessageHttp = requestWithExtensions(OFFER);

      prepareUpgradeRequest(req);

      expect(req.headers[EXTENSIONS_HEADER]).toEqual(OFFER);
    });

    it('leaves a request without an offer untouched', () => {
      config.WS_COMPRESSION_ROLLOUT_PERCENT = 100;
      const req: IncomingMessageHttp = requestWithExtensions('some-other-extension');

      prepareUpgradeRequest(req);

      expect(req.headers[EXTENSIONS_HEADER]).toEqual('some-other-extension');
    });
  });

  describe('trackConnection', () => {
    it.each([
      ['a client that made no offer', undefined, '', CompressionState.NOT_OFFERED],
      ['an offer outside the rollout', OFFER, '', CompressionState.OFFERED],
      ['a negotiated extension', OFFER, 'permessage-deflate', CompressionState.COMPRESSED],
    ])('reports %s', (
      _name: string,
      offered: string | undefined,
      negotiated: string,
      expected: CompressionState,
    ) => {
      config.WS_COMPRESSION_ROLLOUT_PERCENT = 0;
      const req: IncomingMessageHttp = requestWithExtensions(offered);
      prepareUpgradeRequest(req);

      expect(trackConnection(websocketWithExtensions(negotiated), req)).toEqual(expected);
      expect(stats.increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.ws_compression.connections`,
        1,
        expect.objectContaining({ state: expected }),
      );
    });
  });

  describe('flushCompressionStats', () => {
    it('emits payload totals by state and threshold and resets them', () => {
      const offeredReq: IncomingMessageHttp = requestWithExtensions(OFFER);
      prepareUpgradeRequest(offeredReq);
      const offeredWs: WebSocket = websocketWithExtensions('');
      trackConnection(offeredWs, offeredReq);
      const untrackedWs: WebSocket = websocketWithExtensions('');

      recordMessageSent(offeredWs, 100);
      recordMessageSent(offeredWs, 50);
      recordMessageSent(offeredWs, config.WS_COMPRESSION_THRESHOLD_BYTES);
      recordMessageSent(untrackedWs, 7);
      (stats.increment as jest.Mock).mockClear();
      flushCompressionStats();

      expect(stats.increment).toHaveBeenCalledTimes(3);
      expect(stats.increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.ws_compression.payload_bytes`,
        150,
        expect.objectContaining({ state: CompressionState.OFFERED, over_threshold: 'false' }),
      );
      expect(stats.increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.ws_compression.payload_bytes`,
        config.WS_COMPRESSION_THRESHOLD_BYTES,
        expect.objectContaining({ state: CompressionState.OFFERED, over_threshold: 'true' }),
      );
      expect(stats.increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.ws_compression.payload_bytes`,
        7,
        expect.objectContaining({ state: CompressionState.NOT_OFFERED, over_threshold: 'false' }),
      );

      (stats.increment as jest.Mock).mockClear();
      flushCompressionStats();
      expect(stats.increment).not.toHaveBeenCalled();
    });
  });

  describe('handshake', () => {
    let wss: Wss;

    async function startServer(rolloutPercent: number): Promise<string> {
      config.WS_COMPRESSION_ROLLOUT_PERCENT = rolloutPercent;
      // Increment port with a large number to ensure it's not used for any other service.
      config.WS_PORT += 1723;
      wss = new Wss();
      await wss.start();
      return `ws://localhost:${config.WS_PORT}`;
    }

    // Resolves with the extensions the client negotiated and the state the server recorded.
    function connect(
      url: string,
      perMessageDeflate: boolean | WebSocket.PerMessageDeflateOptions,
    ): Promise<{ extensions: string, state: CompressionState }> {
      return new Promise((resolve, reject) => {
        const trackedStates: CompressionState[] = [];
        (stats.increment as jest.Mock).mockImplementation(
          (name: string, _value: number, tags: { state: CompressionState }) => {
            if (name === `${config.SERVICE_NAME}.ws_compression.connections`) {
              trackedStates.push(tags.state);
            }
          },
        );
        wss.onConnection(() => {});
        const client: WebSocket = new WebSocket(url, { perMessageDeflate });
        client.on('error', reject);
        client.on('open', () => {
          const extensions: string = client.extensions;
          client.close();
          resolve({ extensions, state: trackedStates[0] });
        });
      });
    }

    afterEach(async () => {
      await wss.close();
    });

    it('does not negotiate compression when the rollout is at 0', async () => {
      const url: string = await startServer(0);
      expect(await connect(url, true)).toEqual({
        extensions: '',
        state: CompressionState.OFFERED,
      });
    });

    it('negotiates compression for an offering client when the rollout is at 100', async () => {
      const url: string = await startServer(100);
      expect(await connect(url, true)).toEqual({
        extensions: 'permessage-deflate',
        state: CompressionState.COMPRESSED,
      });
    });

    it('delivers a message over the threshold intact on a compressed connection', async () => {
      const url: string = await startServer(100);
      const message: string = JSON.stringify({
        contents: Array(config.WS_COMPRESSION_THRESHOLD_BYTES).fill({ price: '1', size: '2' }),
      });
      wss.onConnection((ws: WebSocket) => sendMessageString(ws, 'connectionId', message));

      const received: string = await new Promise((resolve, reject) => {
        const client: WebSocket = new WebSocket(url);
        client.on('error', reject);
        client.on('message', (data: WebSocket.RawData) => {
          client.close();
          resolve(data.toString());
        });
      });

      expect(received).toEqual(message);
    });

    it('compresses only messages that reach the threshold', async () => {
      const url: string = await startServer(100);
      const underThreshold: string = 'a'.repeat(config.WS_COMPRESSION_THRESHOLD_BYTES - 1);
      const overThreshold: string = 'a'.repeat(config.WS_COMPRESSION_THRESHOLD_BYTES);
      wss.onConnection((ws: WebSocket) => {
        // Reply to a message rather than on connection, so that the frames are not delivered in
        // the same packet as the handshake response.
        ws.on('message', () => {
          sendMessageString(ws, 'connectionId', underThreshold);
          sendMessageString(ws, 'connectionId', overThreshold);
        });
      });

      // The raw frames the server sent, alongside the messages the client decoded from them.
      const { frames, messages }: { frames: Buffer, messages: string[] } = await new Promise(
        (resolve, reject) => {
          const chunks: Buffer[] = [];
          const received: string[] = [];
          // Default client offer, so the connection keeps the server compression context.
          const client: WebSocket = new WebSocket(url);
          client.on('error', reject);
          client.on('upgrade', (response: IncomingMessageHttp) => {
            response.socket.on('data', (chunk: Buffer) => chunks.push(chunk));
          });
          client.on('open', () => client.send('ready'));
          client.on('message', (data: WebSocket.RawData) => {
            received.push(data.toString());
            if (received.length === 2) {
              client.close();
              resolve({ frames: Buffer.concat(chunks), messages: received });
            }
          });
        },
      );

      expect(messages).toEqual([underThreshold, overThreshold]);
      // The RSV1 bit of a frame's first byte marks its payload as compressed. The first frame's
      // payload is too long for the 7-bit length field, so it has a 4-byte header.
      const RSV1: number = 0x40;
      const secondFrame: number = 4 + underThreshold.length;
      // eslint-disable-next-line no-bitwise
      expect(frames[0] & RSV1).toEqual(0);
      // eslint-disable-next-line no-bitwise
      expect(frames[secondFrame] & RSV1).toEqual(RSV1);
    });

    it('negotiates compression for a client that asks for no server context takeover', async () => {
      const url: string = await startServer(100);
      // Offers `server_no_context_takeover; client_no_context_takeover`, as Go clients do.
      const offer: WebSocket.PerMessageDeflateOptions = {
        serverNoContextTakeover: true,
        clientNoContextTakeover: true,
      };
      expect(await connect(url, offer)).toEqual({
        extensions: 'permessage-deflate',
        state: CompressionState.COMPRESSED,
      });
    });

    it('does not log sends that were waiting on compression when the socket closed', async () => {
      const url: string = await startServer(100);
      jest.spyOn(logger, 'error');
      const message: string = 'a'.repeat(config.WS_COMPRESSION_THRESHOLD_BYTES);
      const numMessages: number = 5;

      // Resolves once every send has failed.
      await new Promise<void>((resolve, reject) => {
        let numSendErrors: number = 0;
        (stats.increment as jest.Mock).mockImplementation((name: string) => {
          if (name === `${config.SERVICE_NAME}.ws_send.error`) {
            numSendErrors += 1;
            if (numSendErrors === numMessages) {
              resolve();
            }
          }
        });
        wss.onConnection((ws: WebSocket) => {
          for (let i: number = 0; i < numMessages; i++) {
            sendMessageString(ws, 'connectionId', message);
          }
          ws.terminate();
        });
        const client: WebSocket = new WebSocket(url);
        client.on('error', reject);
      });

      expect(logger.error).not.toHaveBeenCalled();
      const closedWhileCompressingCalls: unknown[][] = (stats.increment as jest.Mock).mock.calls
        .filter(([name]: [string]) => {
          return name === `${config.SERVICE_NAME}.ws_send.closed_while_compressing_errors`;
        });
      expect(closedWhileCompressingCalls).toHaveLength(numMessages);
    });

    it('does not negotiate compression for a client that does not offer it', async () => {
      const url: string = await startServer(100);
      expect(await connect(url, false)).toEqual({
        extensions: '',
        state: CompressionState.NOT_OFFERED,
      });
    });
  });
});
