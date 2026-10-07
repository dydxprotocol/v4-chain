import { IncomingMessage as IncomingMessageHttp, OutgoingHttpHeaders } from 'http';
import { Duplex } from 'stream';

import { stats, getInstanceId, logger } from '@dydxprotocol-indexer/base';
import WebSocket from 'ws';

import config from '../config';
import {
  WS_CLOSE_CODE_ABNORMAL_CLOSURE,
  ERR_SOCKET_CLOSED_WHILE_COMPRESSING,
  ERR_WRITE_STREAM_DESTROYED,
  RATE_LIMIT_REASON_HEADER,
  RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE,
  WEBSOCKET_NOT_OPEN,
} from '../lib/constants';
import { reconnectPenalty } from '../lib/reconnect-penalty';
import { IncomingMessage, OutgoingMessage, WebsocketEvent } from '../types';
import { getClientIp } from './header-utils';
import {
  flushCompressionStats,
  getPerMessageDeflateOptions,
  prepareUpgradeRequest,
  recordMessageSent,
  shouldCompress,
  trackConnection,
} from './ws-compression';

// Status returned for a websocket upgrade refused because the client is in the reconnect
// penalty box. It is an ordinary HTTP response, which is what makes it visible to Cloudflare.
const HTTP_TOO_MANY_REQUESTS: number = 429;

function incrementSendErrorStats(instanceId: string, error: WssError): void {
  stats.increment(
    `${config.SERVICE_NAME}.ws_send.error`,
    1,
    config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
    {
      instance: instanceId,
      code: error?.code,
    },
  );
}

function incrementStreamDestroyedErrorStats(instanceId: string): void {
  stats.increment(
    `${config.SERVICE_NAME}.ws_send.stream_destroyed_errors`,
    1,
    {
      action: 'close',
      instance: instanceId,
    },
  );
}

function incrementWriteEpipeErrorStats(instanceId: string): void {
  stats.increment(
    `${config.SERVICE_NAME}.ws_send.write_epipe_errors`,
    1,
    {
      action: 'close',
      instance: instanceId,
    },
  );
}

function incrementClosedWhileCompressingErrorStats(instanceId: string): void {
  stats.increment(
    `${config.SERVICE_NAME}.ws_send.closed_while_compressing_errors`,
    1,
    {
      action: 'close',
      instance: instanceId,
    },
  );
}

/**
 * Refuses a websocket upgrade from a client that is serving a reconnect penalty.
 *
 * Rejecting at the handshake, rather than accepting and immediately closing, is what makes the
 * penalty legible to the edge: Cloudflare cannot inspect websocket frames, but it can count HTTP
 * 429 responses and match on the reason header, and escalate to a longer block of its own.
 * `Retry-After` tells a well-behaved client how long to wait without needing any of that.
 */
export function verifyClientNotPenalized(
  info: { origin: string, secure: boolean, req: IncomingMessageHttp },
  callback: (
    result: boolean,
    code?: number,
    statusMessage?: string,
    headers?: OutgoingHttpHeaders,
  ) => void,
): void {
  const clientIp: string | undefined = getClientIp(info.req);
  const remainingPenaltyMs: number = reconnectPenalty.getRemainingPenaltyMs(clientIp);

  if (remainingPenaltyMs <= 0) {
    callback(true);
    return;
  }

  const retryAfterSeconds: number = Math.ceil(remainingPenaltyMs / 1000);

  stats.increment(
    `${config.SERVICE_NAME}.connection_rejected_penalty`,
    1,
    undefined,
    {
      instance: getInstanceId(),
    },
  );

  logger.info({
    at: 'wss#verifyClientNotPenalized',
    message: 'Rejected websocket upgrade for client serving a reconnect penalty',
    retryAfterSeconds,
  });

  callback(
    false,
    HTTP_TOO_MANY_REQUESTS,
    'Too Many Requests',
    {
      'Retry-After': String(retryAfterSeconds),
      [RATE_LIMIT_REASON_HEADER]: RATE_LIMIT_REASON_SUBSCRIPTION_LIMIT_ABUSE,
    },
  );
}

/**
 * `ws` negotiates extensions before it calls `verifyClient`, so the upgrade request has to be
 * prepared here for the compression rollout to decide whether an offer is honoured.
 */
class WebSocketServer extends WebSocket.Server {
  public handleUpgrade(
    request: IncomingMessageHttp,
    socket: Duplex,
    upgradeHead: Buffer,
    callback: (client: WebSocket, request: IncomingMessageHttp) => void,
  ): void {
    prepareUpgradeRequest(request);
    super.handleUpgrade(request, socket, upgradeHead, callback);
  }
}

export class Wss {
  private wss: WebSocket.Server;
  private started: boolean;
  private closed: boolean;
  private compressionStatsInterval?: NodeJS.Timeout;

  constructor() {
    this.started = false;
    this.closed = false;

    const serverOptions: WebSocket.ServerOptions = {
      port: config.WS_PORT,
      allowSynchronousEvents: true,
      autoPong: true,
      verifyClient: verifyClientNotPenalized,
      perMessageDeflate: getPerMessageDeflateOptions(),
    };
    this.wss = new WebSocketServer(serverOptions);
  }

  public async start(): Promise<void> {
    if (this.started) {
      throw new Error('Wss already started');
    }

    this.started = true;

    this.wss.on(WebsocketEvent.ERROR, (error: Error) => {
      logger.error({
        at: 'wss#onError',
        message: `WebSocket server threw error: ${error.message}`,
        error,
      });
    });

    await new Promise((resolve) => {
      logger.info({
        at: 'wss#onListening',
        message: 'Listening for websocket connections',
      });
      this.wss.on(WebsocketEvent.LISTENING, resolve);
    });

    this.compressionStatsInterval = setInterval(
      flushCompressionStats,
      config.WS_COMPRESSION_METRIC_INTERVAL_MS,
    );
    // Metrics alone should not keep the process alive.
    this.compressionStatsInterval.unref();
  }

  public onConnection(callback: (ws: WebSocket, req: IncomingMessage) => void): void {
    this.wss.on(WebsocketEvent.CONNECTION, (ws: WebSocket, req: IncomingMessage) => {
      trackConnection(ws, req);
      callback(ws, req);
    });
  }

  public async close(): Promise<void> {
    if (this.closed) {
      throw new Error('Wss already closed');
    }
    if (!this.started) {
      throw new Error('Wss not started');
    }

    if (this.compressionStatsInterval) {
      clearInterval(this.compressionStatsInterval);
    }
    this.wss.close();
    this.closed = true;

    await new Promise((resolve) => {
      this.wss.on(WebsocketEvent.CLOSE, resolve);
    });
  }
}

export class WssError extends Error {
  public code: string;

  constructor(message: string, code: string) {
    super(message);
    this.code = code;
  }
}

export function sendMessage(
  ws: WebSocket,
  connectionId: string,
  message: OutgoingMessage,
): void {
  sendMessageString(ws, connectionId, JSON.stringify(message));
}

export function sendMessageString(
  ws: WebSocket,
  connectionId: string,
  message: string,
): void {
  if (ws.readyState !== WebSocket.OPEN) {
    logger.info({
      at: 'wss#sendMessageString',
      message: 'Not sending message because websocket is not open',
      connectionId,
      readyState: ws.readyState,
    });
    stats.increment(
      `${config.SERVICE_NAME}.ws_message_not_sent`,
      1,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      {
        instance: getInstanceId(),
        reason: WEBSOCKET_NOT_OPEN,
        readyState: ws.readyState.toString(),
      },
    );
    return;
  }

  const payloadLength: number = Buffer.byteLength(message);
  recordMessageSent(ws, payloadLength);
  ws.send(message, { compress: shouldCompress(payloadLength) }, (error) => {
    if (error) {
      const instanceId = getInstanceId();
      incrementSendErrorStats(instanceId, error as WssError);
      // Don't log to avoid bursts when clients disconnect abruptly
      if (error?.message.includes?.(ERR_WRITE_STREAM_DESTROYED)) {
        incrementStreamDestroyedErrorStats(instanceId);
      } else if (error?.message.includes?.('EPIPE')) {
        incrementWriteEpipeErrorStats(instanceId);
      } else if (error?.message.includes?.(ERR_SOCKET_CLOSED_WHILE_COMPRESSING)) {
        incrementClosedWhileCompressingErrorStats(instanceId);
      } else {
        const errorLog = { // type is InfoObject in node-service-base
          at: 'wss#sendMessageString',
          message: `Failed to send message: ${error.message}`,
          error,
          connectionId,
          code: (error as WssError)?.code,
        };
        logger.error(errorLog);
      }
      try {
        ws.close(
          WS_CLOSE_CODE_ABNORMAL_CLOSURE,
          error?.message,
        );
      } catch (closeError) {
        // These errors indicate the underlying Socket was destroyed
        // Don't log an error as this can be expected when clients disconnect abruptly and
        // can happen to multiple closes while the close handshake is going on
        if (closeError?.message.includes?.(ERR_WRITE_STREAM_DESTROYED)) {
          incrementStreamDestroyedErrorStats(instanceId);
        } else if (closeError?.message.includes?.('EPIPE')) {
          incrementWriteEpipeErrorStats(instanceId);
        } else {
          const closeErrorLog = {
            at: 'wss#sendMessageString',
            message: `Failed to close connection: ${closeError.message}`,
            connectionId,
            closeError,
          };
          logger.error(closeErrorLog);
        }
      }
    }
  });
}
