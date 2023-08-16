import { logger } from '@dydxprotocol-indexer/base';

import config from '../config';
import { createErrorMessage, createPongMessage } from '../helpers/message';
import { sendMessage } from '../helpers/wss';
import {
  Connection,
  PingMessage,
} from '../types';
import { WS_CLOSE_CODE_POLICY_VIOLATION } from './constants';
import { RateLimiter } from './rate-limit';

export class PingHandler {
  private rateLimiter: RateLimiter;

  constructor() {
    this.rateLimiter = new RateLimiter({
      points: config.RATE_LIMIT_PING_POINTS,
      durationMs: config.RATE_LIMIT_PING_DURATION_MS,
    });
  }

  public handlePing(
    pingMessage: PingMessage,
    connection: Connection,
    connectionId: string,
  ): void {
    const duration: number = this.rateLimiter.rateLimit({
      connectionId,
      key: 'ping',
    });
    if (duration > 0) {
      sendMessage(
        connection.ws,
        connectionId,
        createErrorMessage(
          'Too many ping messages. Please reconnect and try again.',
          connectionId,
          connection.messageId,
        ),
      );

      // Violated rate-limit; disconnect.
      connection.ws.close(
        WS_CLOSE_CODE_POLICY_VIOLATION,
        JSON.stringify({ message: 'Rate limited' }),
      );

      logger.info({
        at: 'ping#handlePing',
        message: 'Connection closed due to violating rate limit',
        connectionId,
      });
      return;
    }

    sendMessage(
      connection.ws,
      connectionId,
      createPongMessage(
        connectionId,
        connection.messageId,
        pingMessage.id,
      ),
    );
  }

  public handleDisconnect(
    connectionId: string,
  ): void {
    this.rateLimiter.removeConnection(connectionId);
  }
}
