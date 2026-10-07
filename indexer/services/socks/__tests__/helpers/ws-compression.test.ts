import { IncomingMessage as IncomingMessageHttp } from 'http';
import { Socket } from 'net';

import { stats } from '@dydxprotocol-indexer/base';
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

    it('keeps the server compression context when the rollout is above 0', () => {
      config.WS_COMPRESSION_ROLLOUT_PERCENT = 1;
      expect(getPerMessageDeflateOptions()).toEqual(expect.objectContaining({
        serverNoContextTakeover: false,
        threshold: config.WS_COMPRESSION_THRESHOLD_BYTES,
      }));
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
    it('emits payload totals by state and resets them', () => {
      const offeredReq: IncomingMessageHttp = requestWithExtensions(OFFER);
      prepareUpgradeRequest(offeredReq);
      const offeredWs: WebSocket = websocketWithExtensions('');
      trackConnection(offeredWs, offeredReq);
      const untrackedWs: WebSocket = websocketWithExtensions('');

      recordMessageSent(offeredWs, 100);
      recordMessageSent(offeredWs, 50);
      recordMessageSent(untrackedWs, 7);
      (stats.increment as jest.Mock).mockClear();
      flushCompressionStats();

      expect(stats.increment).toHaveBeenCalledTimes(2);
      expect(stats.increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.ws_compression.payload_bytes`,
        150,
        expect.objectContaining({ state: CompressionState.OFFERED }),
      );
      expect(stats.increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.ws_compression.payload_bytes`,
        7,
        expect.objectContaining({ state: CompressionState.NOT_OFFERED }),
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
      perMessageDeflate: boolean,
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

    it('does not negotiate compression for a client that does not offer it', async () => {
      const url: string = await startServer(100);
      expect(await connect(url, false)).toEqual({
        extensions: '',
        state: CompressionState.NOT_OFFERED,
      });
    });
  });
});
