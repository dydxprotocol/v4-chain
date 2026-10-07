import { IncomingMessage as IncomingMessageHttp } from 'http';

import { stats, getInstanceId } from '@dydxprotocol-indexer/base';
import WebSocket from 'ws';

import config from '../config';

const EXTENSIONS_HEADER: string = 'sec-websocket-extensions';
const PERMESSAGE_DEFLATE: string = 'permessage-deflate';
// Favours speed over ratio. Compression runs on the libuv threadpool, but its cost still scales
// with outbound volume, and most of the saving on repetitive JSON is available at low levels.
const DEFLATE_LEVEL: number = 3;

export enum CompressionState {
  // The client did not offer permessage-deflate, so it can never receive compressed messages.
  NOT_OFFERED = 'not_offered',
  // The client offered permessage-deflate but the connection is outside the rollout.
  OFFERED = 'offered',
  // permessage-deflate was negotiated; outbound messages over the threshold are compressed.
  COMPRESSED = 'compressed',
}

// Upgrade requests whose client offered permessage-deflate, recorded before the offer is
// possibly stripped so that the connection can still be counted as one that would accept it.
const requestsOfferingCompression: WeakSet<IncomingMessageHttp> = new WeakSet();
const stateByWebsocket: WeakMap<WebSocket, CompressionState> = new WeakMap();
const payloadLengthByState: Record<CompressionState, number> = {
  [CompressionState.NOT_OFFERED]: 0,
  [CompressionState.OFFERED]: 0,
  [CompressionState.COMPRESSED]: 0,
};

/**
 * Options for the `ws` permessage-deflate extension, or `false` to leave it disabled, which is
 * the behaviour of the server before compression was configurable.
 */
export function getPerMessageDeflateOptions(): WebSocket.PerMessageDeflateOptions | false {
  if (config.WS_COMPRESSION_ROLLOUT_PERCENT <= 0) {
    return false;
  }
  return {
    // Keep the compression context between messages. Channel messages are small and repetitive,
    // so most of the saving comes from referring back to earlier messages on the connection.
    serverNoContextTakeover: false,
    // Clients send very little, so do not hold a decompression context per connection for it.
    clientNoContextTakeover: true,
    zlibDeflateOptions: { level: DEFLATE_LEVEL },
    threshold: config.WS_COMPRESSION_THRESHOLD_BYTES,
  };
}

/**
 * Must run on an upgrade request before `ws` negotiates extensions. Records whether the client
 * offered permessage-deflate, then removes the offer for connections outside the rollout so that
 * their handshake is identical to one made with compression disabled.
 */
export function prepareUpgradeRequest(req: IncomingMessageHttp): void {
  const offered: string | string[] | undefined = req.headers[EXTENSIONS_HEADER];
  if (typeof offered !== 'string' || !offered.includes(PERMESSAGE_DEFLATE)) {
    return;
  }
  requestsOfferingCompression.add(req);
  if (Math.random() * 100 >= config.WS_COMPRESSION_ROLLOUT_PERCENT) {
    // eslint-disable-next-line no-param-reassign
    delete req.headers[EXTENSIONS_HEADER];
  }
}

/**
 * Records the compression state of a newly established connection.
 */
export function trackConnection(ws: WebSocket, req: IncomingMessageHttp): CompressionState {
  let state: CompressionState = CompressionState.NOT_OFFERED;
  if (ws.extensions?.includes(PERMESSAGE_DEFLATE)) {
    state = CompressionState.COMPRESSED;
  } else if (requestsOfferingCompression.has(req)) {
    state = CompressionState.OFFERED;
  }
  stateByWebsocket.set(ws, state);
  stats.increment(
    `${config.SERVICE_NAME}.ws_compression.connections`,
    1,
    {
      instance: getInstanceId(),
      state,
    },
  );
  return state;
}

/**
 * Adds an outbound message to the running total for its connection's compression state. The
 * totals are emitted by `flushCompressionStats` rather than here, as this runs on every send.
 * @param payloadLength Length of the message string. Messages are ASCII JSON, so this is the
 * payload size in bytes before compression.
 */
export function recordMessageSent(ws: WebSocket, payloadLength: number): void {
  payloadLengthByState[stateByWebsocket.get(ws) ?? CompressionState.NOT_OFFERED] += payloadLength;
}

/**
 * Emits and resets the outbound payload totals. The share sent to connections that offered or
 * negotiated compression is the share of websocket egress that compression can reduce.
 */
export function flushCompressionStats(): void {
  const instance: string = getInstanceId();
  Object.values(CompressionState).forEach((state: CompressionState) => {
    const payloadLength: number = payloadLengthByState[state];
    if (payloadLength === 0) {
      return;
    }
    payloadLengthByState[state] = 0;
    stats.increment(
      `${config.SERVICE_NAME}.ws_compression.payload_bytes`,
      payloadLength,
      {
        instance,
        state,
      },
    );
  });
}
