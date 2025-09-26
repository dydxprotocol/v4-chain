import { randomUUID } from 'node:crypto';

import {
  InfoObject, getInstanceId, logger, safeJsonStringify, stats,
} from '@dydxprotocol-indexer/base';
import WebSocket from 'ws';

import config from '../config';
import { getGeoOriginHeaders } from '../helpers/header-utils';
import { createConnectedMessage, createErrorMessage, createUnsubscribedMessage } from '../helpers/message';
import { sendMessage, Wss } from '../helpers/wss';
import { ERR_INVALID_WEBSOCKET_FRAME, WS_CLOSE_CODE_SERVICE_RESTART, WS_CLOSE_HEARTBEAT_TIMEOUT } from '../lib/constants';
import { InvalidMessageHandler } from '../lib/invalid-message';
import { Subscriptions } from '../lib/subscription';
import {
  ALL_CHANNELS,
  Channel,
  Connection,
  IncomingMessage,
  IncomingMessageType,
  SubscribeMessage,
  UnsubscribeMessage,
  WebsocketEvent,
} from '../types';

const HEARTBEAT_INTERVAL_MS: number = config.WS_HEARTBEAT_INTERVAL_MS;
const HEARTBEAT_TIMEOUT_MS: number = config.WS_HEARTBEAT_TIMEOUT_MS;

export class Index {
  public connections: { [connectionId: string]: Connection };

  private wss: Wss;
  private subscriptions: Subscriptions;
  private invalidMessageHandler: InvalidMessageHandler;

  constructor(wss: Wss, subscriptions: Subscriptions) {
    this.wss = wss;
    this.connections = {};
    this.subscriptions = subscriptions;
    this.invalidMessageHandler = new InvalidMessageHandler();
    this.wss.onConnection((ws: WebSocket, req: IncomingMessage) => this.onConnection(ws, req));
  }

  public async close(): Promise<void> {
    logger.info({
      at: 'index#close',
      message: 'Closing websocket server...',
    });

    const restartMessage: string = JSON.stringify({ message: 'Service restarting' });
    Object.keys(this.connections).forEach((connectionId: string) => {
      try {
        this.connections[connectionId].ws.close(
          WS_CLOSE_CODE_SERVICE_RESTART,
          restartMessage,
        );
        logger.info({
          at: 'index#close',
          message: 'Connection successfully closed',
          connectionId,
        });
      } catch (error) {
        logger.error({
          at: 'index#close',
          message: `Failed to close connection: ${error.message}`,
          connectionId,
          error,
        });
      }
    });
    await this.wss.close();

    logger.info({
      at: 'index#close',
      message: 'Finished closing websocket server',
    });
  }

  /**
   * Handler for new websocket connections.
   * Stores the connection and assigns a UUID to the connection to keep track of it, along with
   * a message id counter to track the messages sent to the connection.
   * Attaches message, heartbeat, error, and close handlers to the connection.
   * @param ws New websocket connection.
   * @param req HTTP request accompanying new connection request.
   */
  private onConnection(ws: WebSocket, req: IncomingMessage): void {

    const instanceId: string = getInstanceId();

    const connectionId: string = randomUUID();
    this.connections[connectionId] = {
      ws,
      messageId: 0,
      id: connectionId,
      geoOriginHeaders: getGeoOriginHeaders(req),
    } as Connection;
    const connection: Connection = this.connections[connectionId];

    const numConcurrentConnections: number = Object.keys(this.connections).length;
    logger.info({
      at: 'index#onConnection',
      message: 'Received websocket connection',
      url: ws.url,
      protocol: ws.protocol,
      connectionId,
      headers: req.headers,
      numConcurrentConnections,
    });
    stats.increment(
      `${config.SERVICE_NAME}.num_connections`,
      1,
      {
        instance: instanceId,
      },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.num_concurrent_connections`,
      numConcurrentConnections,
      {
        instance: instanceId,
      },
    );

    try {
      sendMessage(ws, connectionId, createConnectedMessage(connectionId));
    } catch (error) {
      logger.error({
        at: 'index#onConnection',
        message: `Failed to send connected message: ${error.message}`,
        error,
        connectionId,
      });
    }

    ws.on(WebsocketEvent.MESSAGE, (message: WebSocket.RawData) => {
      try {
        this.onMessage(connection, message);
      } catch (error) {
        logger.error({
          at: 'index#onMessage',
          message: `Failed to handle message: ${error.message}`,
          error,
          connectionId,
          messageContents: safeJsonStringify(message),
        });
        sendMessage(
          connection.ws,
          connectionId,
          createErrorMessage(
            'Internal error',
            connectionId,
            connection.messageId,
          ),
        );
      }
    });

    ws.on(WebsocketEvent.CLOSE, (code: number, reason: Buffer) => {
      // TODO remove use as developer indicates
      const reasonStr = safeJsonStringify(reason);

      logger.info({
        at: 'index#onClose',
        message: 'Connection closed',
        connectionId,
        code,
        // TODO remove comment
        reason: reasonStr, // `reason` could be a Buffer, which cannot be logged
      });

      stats.increment(
        `${config.SERVICE_NAME}.num_disconnects`,
        1,
        {
          code: String(code),
          reason: reasonStr,
          instance: instanceId,
        },
      );

      this.disconnect(connection);
    });

    ws.on(WebsocketEvent.ERROR, (error: Error) => {
      const errorLog: InfoObject = {
        at: 'index#onError',
        message: `Connection threw error: ${error}`,
        connectionId,
        error,
      };
      if (error?.message.includes?.(ERR_INVALID_WEBSOCKET_FRAME)) {
        // Clients sending invalid frames is not considered an error
        logger.info(errorLog);
      } else {
        logger.error(errorLog);
      }
    });

    ws.on(WebsocketEvent.PONG, () => {
      // Clear the delayed disconnect set by the heartbeat handler when a pong is received.
      if (connection.disconnect) {
        clearTimeout(connection.disconnect);
        delete connection.disconnect;
      }
    });

    connection.heartbeat = setInterval(
      () => this.heartbeat(connection),
      HEARTBEAT_INTERVAL_MS,
    );
  }

  /**
   * Handler for messages received from a websocket connection.
   * Parses then message, then forwards subscribe/unsubscribe requests to a Subscriptions object,
   * and handles pings/invalid messages.
   * @param connectionId Id of the websocket connection.
   * @param message Message received from the websocket connection
   * @returns
   */
  private onMessage(connection: Connection, message: WebSocket.Data): void {
    const connectionId: string = connection.id;
    const instanceId: string = getInstanceId();

    stats.increment(
      `${config.SERVICE_NAME}.on_message`,
      1,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      {
        instance: instanceId,
      },
    );
    if (!connection) {
      logger.info({
        at: 'index#onMessage',
        message: 'Received message for closed connection',
        connectionId,
        messageContents: safeJsonStringify(message),
      });
      return;
    }

    let parsed: IncomingMessage;
    try {
      parsed = JSON.parse(message.toString());
    } catch (error) {
      this.invalidMessageHandler.handleInvalidMessage(
        'Invalid message: could not parse',
        connection,
        connectionId,
      );
      return;
    }

    if (!parsed.type) {
      this.invalidMessageHandler.handleInvalidMessage(
        'Invalid message: type is required',
        connection,
        connectionId,
      );
      return;
    }

    switch (parsed.type) {
      case IncomingMessageType.SUBSCRIBE: {
        const subscribeMessage = parsed as SubscribeMessage;
        if (!this.validateSubscriptionMessage(connection, subscribeMessage)) {
          return;
        }

        logger.info({
          at: 'index#onSubscribe',
          message: 'Received websocket subscribe',
          connectionId,
          messageContents: safeJsonStringify(message),
        });

        // eslint-disable-next-line  no-param-reassign
        connection.messageId += 1;

        // Do not wait for this.
        this.subscriptions.subscribe(
          connection.ws,
          subscribeMessage.channel,
          connectionId,
          connection.messageId,
          subscribeMessage.id,
          subscribeMessage.batched,
          connection.geoOriginHeaders,
        ).catch((error: Error) => logger.error({
          at: 'Subscription#subscribe',
          message: `Subscribing threw error: ${error.message}`,
          error,
          connectionId,
        }));
        break;
      }
      case IncomingMessageType.UNSUBSCRIBE: {
        const unsubscribeMessage = parsed as UnsubscribeMessage;
        if (!this.validateSubscriptionMessage(connection, unsubscribeMessage)) {
          return;
        }

        this.subscriptions.unsubscribe(
          connectionId,
          unsubscribeMessage.channel,
          unsubscribeMessage.id,
        );

        // eslint-disable-next-line  no-param-reassign
        connection.messageId += 1;

        sendMessage(
          connection.ws,
          connectionId,
          createUnsubscribedMessage(
            connectionId,
            connection.messageId,
            unsubscribeMessage.channel,
            unsubscribeMessage.id,
          ),
        );
        break;
      }
      // Handle pings by doing nothing. The ws library automatically responds to pings with pongs.
      case IncomingMessageType.PING: {
        break;
      }
      default: {
        this.invalidMessageHandler.handleInvalidMessage(
          `Invalid message type: ${parsed.type}`,
          connection,
          connectionId,
        );
        return;
      }
    }
    stats.increment(
      `${config.SERVICE_NAME}.message_received_${parsed.type}`,
      1,
      config.MESSAGE_FORWARDER_STATSD_SAMPLE_RATE,
      {
        instance: instanceId,
      },
    );
  }

  /**
   * Heartbeat function to periodically send pings to a websocket connection, and disconnect if a
   * pong isn't received from the websocket connection within `HEARTBEAT_TIMEOUT_MS` from the ping.
   * @param connection websocket connection heartbeat pings are sent to
   * @returns
   */
  private heartbeat(connection: Connection): void {
    const connectionId: string = connection.id;
    if (!connection || connection.disconnect) {
      return;
    }

    try {
      logger.info({
        at: 'index#heartbeat',
        message: 'Sending heartbeat ping',
        connectionId,
      });

      // eslint-disable-next-line  no-param-reassign
      connection.disconnect = setTimeout(
        () => this.heartbeatDisconnect(connection),
        HEARTBEAT_TIMEOUT_MS,
      );
      // Fence: ping after starting timer to disconnect unhealthy outbound
      connection.ws.ping();
    } catch (error) {
      logger.error({
        at: 'index#heartbeat',
        message: `Heartbeat threw error: ${error.message}`,
        error,
        connectionId,
      });
    }
  }

  /**
   * Attempt to close connection with heartbeat timeout error
   * If an error occurs during closing, disconnect and log error
   * @param connection websocket connection to disconnect
   */
  private heartbeatDisconnect(connection: Connection): void {
    const connectionId: string = connection.id;
    const ws: WebSocket = connection.ws;
    try {
      ws.close(
        WS_CLOSE_HEARTBEAT_TIMEOUT,
        'Heartbeat timeout',
      );
    } catch (error) {
      this.disconnect(connection);
      logger.error({
        at: 'index#dc',
        message: 'Error closing websocket after heartbeat timeout',
        connectionId,
        error,
      });
    }
  }

  /**
   * Disconnect a websocket connection:
   * - Remove all event listeners for the connection
   * - Terminate the connection
   * - Remove from subscriptions
   * - Remove from invalid message handler rate limiter
   * - Remove from connections map
   * @param connection websocket connection to disconnect
   * @returns
   */
  public disconnect(connection: Connection): void {
    const connectionId: string = connection.id;
    if (!connection) {
      return;
    }
    try {
      logger.info({
        at: 'index#disconnect',
        message: 'Disconnecting websocket connection',
        connectionId,
      });

      if (connection.heartbeat) {
        clearInterval(connection.heartbeat);
      }
      const ws: WebSocket = connection.ws;
      ws.removeAllListeners();
      ws.terminate();

      this.subscriptions.remove(connectionId);
      this.invalidMessageHandler.handleDisconnect(connectionId);
      delete this.connections[connectionId];
    } catch (error) {
      logger.error({
        at: 'index#disconnect',
        message: `Disconnecting threw error: ${error.message}`,
        error,
        connectionId,
      });
    }
  }

  /**
   * Performs basic validation of a subscription message. Checks that a valid channel exists as a
   * property on the message, and that the id is valid for the channel.
   * @param connection connection that sent the subscription message
   * @param message Subscripe or Unsubscribe message
   * @returns true if valid, otherwise false.
   */
  private validateSubscriptionMessage(
    connection: Connection,
    message: SubscribeMessage | UnsubscribeMessage,
  ): boolean {
    const connectionId: string = connection.id;
    const ws: WebSocket = connection.ws;
    if (!message.channel) {
      sendMessage(
        ws,
        connectionId,
        createErrorMessage(
          'Invalid subscribe message: channel is required',
          connectionId,
          connection.messageId,
        ),
      );
      return false;
    }
    if (!ALL_CHANNELS.includes(message.channel)) {
      sendMessage(
        ws,
        connectionId,
        createErrorMessage(
          `Invalid channel: ${message.channel}`,
          connectionId,
          connection.messageId,
        ),
      );
      return false;
    }
    if (!this.validateSubscriptionForChannel(message)) {
      sendMessage(
        ws,
        connectionId,
        createErrorMessage(
          `Invalid id: ${message.id}`,
          connectionId,
          connection.messageId,
        ),
      );
      return false;
    }

    return true;
  }

  /**
   * Performs basic validation for subscribing to a channel. Checks that an id exists (except for
   * the markets channel).
   * @param message Subscription message to validate.
   * @returns true if valid, false otherwise.
   */
  private validateSubscriptionForChannel(
    message: SubscribeMessage | UnsubscribeMessage,
  ): boolean {
    if (message.channel === Channel.V4_MARKETS || message.channel === Channel.V4_BLOCK_HEIGHT) {
      return true;
    }
    return message.id !== undefined && typeof message.id === 'string';
  }
}
