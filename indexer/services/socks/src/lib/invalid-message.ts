import { logger } from '@dydxprotocol-indexer/base';

import config from '../config';
import { createErrorMessage } from '../helpers/message';
import { sendMessage } from '../helpers/wss';
import { Connection } from '../types';
import { WS_CLOSE_CODE_POLICY_VIOLATION } from './constants';
import { RateLimiter } from './rate-limit';

export const RATE_LIMITED: string = JSON.stringify({ message: 'Rate limited' });

export class InvalidMessageHandler {
  private rateLimiter: RateLimiter;

  constructor() {
    this.rateLimiter = new RateLimiter({
      points: config.RATE_LIMIT_INVALID_MESSAGE_POINTS,
      durationMs: config.RATE_LIMIT_INVALID_MESSAGE_DURATION_MS,
    });
  }

  public handleInvalidMessage(
    responseMessage: string,
    connection: Connection,
    connectionId: string,
  ): void {
    const duration: number = this.rateLimiter.rateLimit({
      connectionId,
      key: 'invalidMessage',
    });
    if (duration > 0) {
      sendMessage(
        connection.ws,
        connectionId,
        // TODO we should not send their connection id in the error message
        createErrorMessage(
          'Too many invalid messages. Please reconnect and try again.',
          connectionId,
          connection.messageId,
        ),
      );

      connection.ws.close(
        WS_CLOSE_CODE_POLICY_VIOLATION,
        RATE_LIMITED,
      );

      logger.info({
        at: 'invalid-message#handleInvalidMessage',
        message: 'Connection closed due to violating rate limit',
        connectionId,
      });
      return;
    }

    sendMessage(
      connection.ws,
      connectionId,
      createErrorMessage(
        responseMessage,
        connectionId,
        connection.messageId,
      ),
    );
  }

  public handleDisconnect(
    connectionId: string,
  ): void {
    this.rateLimiter.removeConnection(connectionId);
  }
}
