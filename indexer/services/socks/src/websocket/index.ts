import {
  stats, logger, safeJsonStringify, InfoObject,
} from '@dydxprotocol-indexer/base';
import { v4 as uuidv4 } from 'uuid';
import WebSocket from 'ws';

import config from '../config';
import {
  createErrorMessage,
  createConnectedMessage,
  createUnsubscribedMessage,
} from '../helpers/message';
import { Wss, sendMessage } from '../helpers/wss';
import { ERR_INVALID_WEBSOCKET_FRAME, WS_CLOSE_CODE_SERVICE_RESTART } from '../lib/constants';
import { InvalidMessageHandler } from '../lib/invalid-message';
import { PingHandler } from '../lib/ping';
import { Subscriptions } from '../lib/subscription';
import {
  IncomingMessageType,
  Channel,
  IncomingMessage,
  SubscribeMessage,
  UnsubscribeMessage,
  Connection,
  PingMessage,
  ALL_CHANNELS,
  WebsocketEvents,
} from '../types';
import { CountryRestrictor } from './restrict-countries';

const HEARTBEAT_INTERVAL_MS: number = config.WS_HEARTBEAT_INTERVAL_MS;
const HEARTBEAT_TIMEOUT_MS: number = config.WS_HEARTBEAT_TIMEOUT_MS;

export class Index {
  // Map of connection ids (UUID/V4) to the websocket connections.
  public connections: { [connectionId: string]: Connection };

  // Websocket server.
  private wss: Wss;
  // Subscriptions tracking object (see lib/subscriptions.ts).
  private subscriptions: Subscriptions;
  // Handlers for pings and invalid messages.
  private pingHandler: PingHandler;
  private invalidMessageHandler: InvalidMessageHandler;
  private countryRestrictor: CountryRestrictor;

  constructor(wss: Wss, subscriptions: Subscriptions) {
    this.wss = wss;
    this.connections = {};
    this.subscriptions = subscriptions;
    this.pingHandler = new PingHandler();
    this.invalidMessageHandler = new InvalidMessageHandler();
    this.countryRestrictor = new CountryRestrictor();

    // Attach the new connection handler to the websocket server.
    this.wss.onConnection((ws: WebSocket, req: IncomingMessage) => this.onConnection(ws, req));
  }

  public async close(): Promise<void> {
    logger.info({
      at: 'index#close',
      message: 'Closing websocket server...',
    });

    Object.keys(this.connections).forEach((connectionId: string) => {
      try {
        this.connections[connectionId].ws.close(
          WS_CLOSE_CODE_SERVICE_RESTART,
          JSON.stringify({ message: 'Service restarting' }),
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
    // Terminate the connection if the connection requestion originated from a restricted country
    if (this.countryRestrictor.isRestrictedCountry(req)) {
      return ws.terminate();
    }

    const connectionId: string = uuidv4();

    this.connections[connectionId] = {
      ws,
      messageId: 0,
      countryCode: this.countryRestrictor.getCountry(req),
    };

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
    stats.increment(`${config.SERVICE_NAME}.num_connections`, 1);
    stats.gauge(
      `${config.SERVICE_NAME}.num_concurrent_connections`,
      numConcurrentConnections,
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

    // Attach message handler to connection.
    ws.on(WebsocketEvents.MESSAGE, (message: WebSocket.RawData) => {
      try {
        this.onMessage(connectionId, message);
      } catch (error) {
        logger.error({
          at: 'index#onMessage',
          message: `Failed to handle message: ${error.message}`,
          error,
          connectionId,
          messageContents: safeJsonStringify(message),
        });
        sendMessage(
          this.connections[connectionId].ws,
          connectionId,
          createErrorMessage(
            'Internal error',
            connectionId,
            this.connections[connectionId].messageId,
          ),
        );
      }
    });

    // Start sending periodic pings (heartbeat) to connection.
    this.connections[connectionId].heartbeat = setInterval(
      () => this.heartbeat(connectionId),
      HEARTBEAT_INTERVAL_MS,
    );

    // Attach handler for pongs (response to heartbeat pings) from connection.
    this.connections[connectionId].ws.on(WebsocketEvents.PONG, () => {
      logger.info({
        at: 'index#onPong',
        message: 'Received pong',
        connectionId,
      });

      // Clear the delayed disconnect set by the heartbeat handler when a pong is received.
      if (this.connections[connectionId].disconnect) {
        clearTimeout(this.connections[connectionId].disconnect);
        delete this.connections[connectionId].disconnect;
      }
    });

    // Attach handler for close events from the connection.
    this.connections[connectionId].ws.on(WebsocketEvents.CLOSE, (code: number, reason: Buffer) => {
      logger.info({
        at: 'index#onClose',
        message: 'Connection closed',
        connectionId,
        code,
        reason: safeJsonStringify(reason), // `reason` could be a Buffer, which cannot be logged
      });

      stats.increment(
        `${config.SERVICE_NAME}.num_disconnects`,
        1,
        { code: String(code), reason: String(reason) },
      );

      this.disconnect(connectionId);
    });

    // Attach error handler to connection.
    this.connections[connectionId].ws.on(WebsocketEvents.ERROR, (error: Error) => {
      const errorLog: InfoObject = {
        at: 'index#onError',
        message: `Connection threw error: ${error}`,
        connectionId,
        error,
      };
      if (error?.message.includes?.(ERR_INVALID_WEBSOCKET_FRAME)) {
        // Clients can send invalid frames that cause an error event, don't log an error as this can
        // happen and there's no way to mitigate the error being emitted by the ws library
        logger.info(errorLog);
      } else {
        logger.error(errorLog);
      }
    });
  }

  /**
   * Handler for messages received from a websocket connection.
   * Parses then message, then forwards subscribe/unsubscribe requests to a Subscriptions object,
   * and handles pings/invalid messages.
   * @param connectionId Id of the websocket connection.
   * @param message Message received from the websocket connection.
   * @returns
   */
  private onMessage(connectionId: string, message: WebSocket.Data): void {
    stats.increment(`${config.SERVICE_NAME}.on_message`, 1);
    if (!this.connections[connectionId]) {
      logger.info({
        at: 'index#onMessage',
        message: 'Received message for closed connection',
        connectionId,
        messageContents: safeJsonStringify(message),
      });
      return;
    }

    this.connections[connectionId].messageId += 1;

    const messageStr = message.toString();

    let parsed: IncomingMessage;
    try {
      parsed = JSON.parse(messageStr);
    } catch (error) {
      this.invalidMessageHandler.handleInvalidMessage(
        'Invalid message: could not parse',
        this.connections[connectionId],
        connectionId,
      );
      return;
    }

    if (!parsed.type) {
      this.invalidMessageHandler.handleInvalidMessage(
        'Invalid message: type is required',
        this.connections[connectionId],
        connectionId,
      );
      return;
    }

    switch (parsed.type) {
      case IncomingMessageType.SUBSCRIBE: {
        const subscribeMessage = parsed as SubscribeMessage;
        if (!this.validateSubscriptionMessage(connectionId, subscribeMessage)) {
          return;
        }

        logger.info({
          at: 'index#onSubscribe',
          message: 'Received websocket subscribe',
          connectionId,
          messageContents: safeJsonStringify(message),
        });

        // Do not wait for this.
        this.subscriptions.subscribe(
          this.connections[connectionId].ws,
          subscribeMessage.channel,
          connectionId,
          this.connections[connectionId].messageId,
          subscribeMessage.id,
          subscribeMessage.batched,
          this.connections[connectionId].countryCode,
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
        if (!this.validateSubscriptionMessage(connectionId, unsubscribeMessage)) {
          return;
        }

        this.subscriptions.unsubscribe(
          connectionId,
          unsubscribeMessage.channel,
          unsubscribeMessage.id,
        );

        sendMessage(
          this.connections[connectionId].ws,
          connectionId,
          createUnsubscribedMessage(
            connectionId,
            this.connections[connectionId].messageId,
            unsubscribeMessage.channel,
            unsubscribeMessage.id,
          ),
        );
        break;
      }
      case IncomingMessageType.PING: {
        this.pingHandler.handlePing(
          parsed as PingMessage,
          this.connections[connectionId],
          connectionId,
        );
        break;
      }
      default: {
        this.invalidMessageHandler.handleInvalidMessage(
          `Invalid message type: ${parsed.type}`,
          this.connections[connectionId],
          connectionId,
        );
        return;
      }
    }
    stats.increment(`${config.SERVICE_NAME}.message_received_${parsed.type}`, 1);
  }

  /**
   * Heartbeat function to periodically send pings to a websocket connection, and disconnect if a
   * pong isn't received from the websocket connection within `HEARTBEAT_TIMEOUT_MS` from the ping.
   * @param connectionId Id of the websocket connection heartbeat pings are sent to.
   * @returns
   */
  private heartbeat(connectionId: string): void {
    if (
      !this.connections[connectionId] ||
      this.connections[connectionId].disconnect
    ) {
      return;
    }
    try {
      logger.info({
        at: 'index#heartbeat',
        message: 'Sending heartbeat ping',
        connectionId,
      });

      // Disconnect the websocket connection after `HEARTBEAT_TIMEOUT_MS` from the ping. This
      // timeout is cleared when a pong is received.
      this.connections[connectionId].disconnect = setTimeout(
        () => this.disconnect(connectionId),
        HEARTBEAT_TIMEOUT_MS,
      );
      this.connections[connectionId].ws.ping();
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
   * Disconnect a websocket connection.
   * @param connectionId Connection id of the websocket connection to disconnect.
   * @returns
   */
  private disconnect(connectionId: string): void {
    if (!this.connections[connectionId]) {
      return;
    }
    try {
      logger.info({
        at: 'index#disconnect',
        message: 'Disconnecting websocket connection',
        connectionId,
      });

      // Remove periodic job sending heartbeat pings to websocket connection.
      if (this.connections[connectionId].heartbeat) {
        clearInterval(this.connections[connectionId].heartbeat);
      }
      this.connections[connectionId].ws.terminate();

      // Delete subscription data.
      this.subscriptions.remove(connectionId);
      this.pingHandler.handleDisconnect(connectionId);
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
   * @param connectionId Id of connection that sent the subscription message.
   * @param message Subscription message.
   * @returns true if valid, otherwise false.
   */
  private validateSubscriptionMessage(
    connectionId: string,
    message: SubscribeMessage | UnsubscribeMessage,
  ): boolean {
    if (!message.channel) {
      sendMessage(
        this.connections[connectionId].ws,
        connectionId,
        createErrorMessage(
          'Invalid subscribe message: channel is required',
          connectionId,
          this.connections[connectionId].messageId,
        ),
      );
      return false;
    }
    if (!ALL_CHANNELS.includes(message.channel)) {
      sendMessage(
        this.connections[connectionId].ws,
        connectionId,
        createErrorMessage(
          `Invalid channel: ${message.channel}`,
          connectionId,
          this.connections[connectionId].messageId,
        ),
      );
      return false;
    }
    if (!this.validateSubscriptionForChannel(message)) {
      sendMessage(
        this.connections[connectionId].ws,
        connectionId,
        createErrorMessage(
          `Invalid id: ${message.id}`,
          connectionId,
          this.connections[connectionId].messageId,
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
    if (message.channel === Channel.V4_MARKETS) {
      return true;
    }
    return message.id !== undefined && typeof message.id === 'string';
  }
}
