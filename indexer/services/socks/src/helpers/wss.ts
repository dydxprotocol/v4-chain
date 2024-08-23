import { stats, getInstanceId, logger } from '@dydxprotocol-indexer/base';
import WebSocket from 'ws';

import config from '../config';
import {
  WS_CLOSE_CODE_ABNORMAL_CLOSURE,
  ERR_WRITE_STREAM_DESTROYED,
  WEBSOCKET_NOT_OPEN,
} from '../lib/constants';
import { IncomingMessage, OutgoingMessage, WebsocketEvents } from '../types';

export class Wss {
  private wss: WebSocket.Server;
  private started: boolean;
  private closed: boolean;

  constructor() {
    this.started = false;
    this.closed = false;
    this.wss = new WebSocket.Server({
      port: config.WS_PORT,
    });
  }

  public async start(): Promise<void> {
    if (this.started) {
      throw new Error('Wss already started');
    }

    this.started = true;

    this.wss.on(WebsocketEvents.ERROR, (error: Error) => {
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
      this.wss.on(WebsocketEvents.LISTENING, resolve);
    });
  }

  public onConnection(callback: (ws: WebSocket, req: IncomingMessage) => void): void {
    this.wss.on(WebsocketEvents.CONNECTION, callback);
  }

  public async close(): Promise<void> {
    if (this.closed) {
      throw new Error('Wss already closed');
    }
    if (!this.started) {
      throw new Error('Wss not started');
    }

    this.wss.close();
    this.closed = true;

    await new Promise((resolve) => {
      this.wss.on(WebsocketEvents.CLOSE, resolve);
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
      at: 'wss#sendMessage',
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

  ws.send(message, (error) => {
    if (error) {
      stats.increment(
        `${config.SERVICE_NAME}.ws_send.error`,
        1,
        config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
        {
          instance: getInstanceId(),
          code: (error as WssError)?.code,
        },
      );
      const errorLog = { // type is InfoObject in node-service-base
        at: 'wss#sendMessage',
        message: `Failed to send message: ${error.message}`,
        error,
        connectionId,
        code: (error as WssError)?.code,
      };
      logger.error(errorLog);
      try {
        ws.removeAllListeners();
        ws.close(
          WS_CLOSE_CODE_ABNORMAL_CLOSURE,
          `client returned ${error?.message} error`,
        );
      } catch (closeError) {
        const closeErrorLog = {
          at: 'wss#sendMessage',
          message: `Failed to close connection: ${closeError.message}`,
          connectionId,
          closeError,
        };
        if (closeError?.message.includes?.(ERR_WRITE_STREAM_DESTROYED)) {
          // This error means the underlying Socket was destroyed
          // Don't log an error as this can be expected when clients disconnect abruptly and
          // can happen to multiple closes while the close handshake is going on
          stats.increment(
            `${config.SERVICE_NAME}.ws_send.stream_destroyed_errors`,
            1,
            {
              action: 'close',
              instance: getInstanceId(),
            },
          );
        } else {
          logger.error(closeErrorLog);
        }
      }
    }
  });
}
