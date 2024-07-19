import {
  AxiosSafeServerError,
  logger,
  stats,
} from '@dydxprotocol-indexer/base';
import { isRestrictedCountry } from '@dydxprotocol-indexer/compliance';
import { CandleResolution, perpetualMarketRefresher } from '@dydxprotocol-indexer/postgres';
import WebSocket from 'ws';

import config from '../config';
import { createErrorMessage, createSubscribedMessage } from '../helpers/message';
import { sendMessage, sendMessageString } from '../helpers/wss';
import {
  Channel,
  MessageToForward,
  RequestMethod,
  Subscription,
  SubscriptionInfo,
} from '../types';
import { axiosRequest } from './axios';
import { V4_MARKETS_ID, WS_CLOSE_CODE_POLICY_VIOLATION } from './constants';
import { BlockedError, InvalidChannelError } from './errors';
import { RateLimiter } from './rate-limit';

const COMLINK_URL: string = `http://${config.COMLINK_URL}`;
const EMPTY_INITIAL_RESPONSE: string = '{}';

export class Subscriptions {
  // Maps channels and ids to a list of websocket connections subscribed to them
  public subscriptions: { [channel: string]: { [id: string]: SubscriptionInfo[] } };
  // Maps connection ids to the channels and ids the connection is subscribed to
  public subscriptionLists: { [connectionId: string]: Subscription[] };
  // Similar to `subscriptions`, maps channels and ids to websocket connections subscribed to them
  // in batched mode (messages will be sent to the connections in batches from these subscriptions)
  public batchedSubscriptions: { [channel: string]: { [id: string]: SubscriptionInfo[] } };
  // Tracks the # of ids per channel socks will forward messages for
  // Make public to access in tests
  public subscribedIdsPerChannel: { [channel: string]: Set<string> };

  private subscribeRateLimiter: RateLimiter;

  private forwardMessage?: (message: MessageToForward, connectionId: string) => number;

  constructor() {
    this.subscriptionLists = {};
    this.subscriptions = {};
    this.batchedSubscriptions = {};
    this.subscribeRateLimiter = new RateLimiter({
      points: config.RATE_LIMIT_SUBSCRIBE_POINTS,
      durationMs: config.RATE_LIMIT_SUBSCRIBE_DURATION_MS,
    });
    this.subscribedIdsPerChannel = {};
    this.forwardMessage = undefined;
  }

  public start(forwardMessage: (message: MessageToForward, connectionId: string) => number): void {
    this.forwardMessage = forwardMessage;
  }

  /**
   * subscribe handles:
   * - mapping a websocket connection to the channel + id it's subscribing to
   * - fetching and sending the initial response for the channel being subscribed to
   * @param ws websocket connection making the subscription request
   * @param channel channel being subscribed to
   * @param connectionId id of the websocket connection
   * @param messageId message id Starting message id for the connection
   * @param id Specific id the websocket is subscribing to within the channel
   * @param batched Whether messages for the subscription should be sent in batches to the
   * connection
   * @returns
   */
  // eslint-disable-next-line  @typescript-eslint/require-await
  public async subscribe(
    ws: WebSocket,
    channel: Channel,
    connectionId: string,
    messageId: number,
    id?: string,
    batched?: boolean,
    country?: string,
  ): Promise<void> {
    if (this.forwardMessage === undefined) {
      throw new Error('Unexpected error, subscription object is uninitialized.');
    }

    if (!this.validateSubscription(channel, id)) {
      sendMessage(
        ws,
        connectionId,
        createErrorMessage(
          `Invalid subscription id for channel: (${channel}-${id})`,
          connectionId,
          messageId,
        ),
      );
      return;
    }

    const subscriptionId: string = this.normalizeSubscriptionId(id);
    const duration: number = this.subscribeRateLimiter.rateLimit({
      connectionId,
      key: channel + subscriptionId,
    });
    if (duration > 0) {
      sendMessage(
        ws,
        connectionId,
        createErrorMessage(
          `Too many subscribe attempts for channel ${channel}-${subscriptionId}. Please ` +
          ' reconnect and try again.',
          connectionId,
          messageId,
        ),
      );

      // Violated rate-limit; disconnect.
      ws.close(
        WS_CLOSE_CODE_POLICY_VIOLATION,
        JSON.stringify({ message: 'Rate limited' }),
      );

      logger.info({
        at: 'subscription#subscribe',
        message: 'Connection closed due to violating rate limit',
        connectionId,
      });
      return;
    }

    let initialResponse: string;
    const startGetInitialResponse: number = Date.now();
    try {
      initialResponse = await this.getInitialResponsesForChannels(channel, id, country);
    } catch (error) {
      logger.info({
        at: 'Subscription#subscribe',
        message: `Making initial request threw error: ${error.message}`,
        error,
        channel,
        connectionId,
        id,
      });

      // For blocked errors add the erorr message into the error message sent in the websocket
      // connection
      let errorMsg: string = `Internal error, could not fetch data for subscription: ${channel}.`;
      if (error instanceof BlockedError) {
        errorMsg = error.message;
      }

      sendMessage(
        ws,
        connectionId,
        createErrorMessage(
          errorMsg,
          connectionId,
          messageId,
        ),
      );

      stats.increment(
        `${config.SERVICE_NAME}.initial_response_error`,
        1,
        undefined,
        {
          channel,
        },
      );
      return;
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.initial_response_get`,
        Date.now() - startGetInitialResponse,
        undefined,
        {
          channel,
        },
      );
    }

    const subscription: Subscription = {
      channel,
      id: subscriptionId,
      batched,
    };
    if (!this.subscriptionLists[connectionId]) {
      this.subscriptionLists[connectionId] = [];
    }
    if (this.subscriptionLists[connectionId].find(
      (s) => (s.channel === channel && s.id === subscriptionId),
    )) {
      sendMessage(
        ws,
        connectionId,
        createErrorMessage(
          `Invalid subscribe message: already subscribed (${channel}-${subscriptionId})`,
          connectionId,
          messageId,
        ),
      );
      return;
    }

    this.subscriptionLists[connectionId].push(subscription);
    if (!this.subscriptions[channel]) {
      this.subscriptions[channel] = {};
      this.subscribedIdsPerChannel[channel] = new Set<string>();
    }

    if (!this.subscriptions[channel][subscriptionId]) {
      this.subscriptions[channel][subscriptionId] = [];
      this.subscribedIdsPerChannel[channel].add(subscriptionId);
    }

    this.subscriptions[channel][subscriptionId].push({
      connectionId,
      pending: true,
      pendingMessages: [],
      batched,
    });

    const startSend: number = Date.now();
    sendMessageString(
      ws,
      connectionId,
      createSubscribedMessage(
        channel,
        subscriptionId,
        initialResponse,
        connectionId,
        messageId,
      ),
    );

    // Enable forwarding all pending messages to the subscriber.
    for (let i = 0; i < this.subscriptions[channel][subscriptionId].length; i += 1) {
      if (this.subscriptions[channel][subscriptionId][i].connectionId === connectionId) {
        this.subscriptions[channel][subscriptionId][i].pending = false;
        this.subscriptions[channel][subscriptionId][i].pendingMessages.forEach(
          (pendingMessage) => this.forwardMessage!(pendingMessage, connectionId),
        );
        // If the subscription was for batched messages, move the connection inside
        // this.batchedSubscriptions.
        if (batched) {
          this.subscriptions[channel][subscriptionId].splice(i, 1);
          if (!this.batchedSubscriptions[channel]) {
            this.batchedSubscriptions[channel] = {};
          }
          if (!this.batchedSubscriptions[channel][subscriptionId]) {
            this.batchedSubscriptions[channel][subscriptionId] = [];
          }
          this.batchedSubscriptions[channel][subscriptionId].push({
            connectionId,
            pending: false,
            pendingMessages: [],
            batched,
          });
        }
      }
    }

    // Stat every time a subscribe happens so we have up to date stats on datadog.
    stats.gauge(
      `${config.SERVICE_NAME}.subscriptions.channel_size`,
      this.subscribedIdsPerChannel[channel].size,
      {
        channel,
      },
    );
    stats.timing(
      `${config.SERVICE_NAME}.subscribe_send_message`,
      Date.now() - startSend,
    );
  }

  /**
   * unsubscribe handles:
   * - removing a websocket connection from data structures tracking channels + ids the connection
   * is subscribed to for a specific channel + id
   * Note: This is a no-op if the connection is not subscribed to the channel + id
   * @param connectionId Connection id of the websocket unsubscribing
   * @param channel Channel being unsubscribed from
   * @param id Specific id within the channel being unsubscribed from
   */
  public unsubscribe(
    connectionId: string,
    channel: Channel,
    id?: string,
  ): void {
    const subscriptionId: string = this.normalizeSubscriptionId(id);
    if (this.subscriptionLists[connectionId]) {
      this.subscriptionLists[connectionId] = this.subscriptionLists[connectionId].filter(
        (e: Subscription) => (e.channel !== channel || e.id !== subscriptionId),
      );
    }

    let subscribedConnections: number = 0;

    // If there is a list of connections subscribed to the channel and id, remove the connection
    // that is being unsubscribed.
    if (this.subscriptions[channel] && this.subscriptions[channel][subscriptionId]) {
      this.subscriptions[channel][subscriptionId] = this.subscriptions[channel][subscriptionId]
        .filter(
          (e: SubscriptionInfo) => (e.connectionId !== connectionId),
        );
      subscribedConnections += this.subscriptions[channel][subscriptionId].length;
    }

    // If there is a list of connections subscribed to the channel and id for batched messages,
    // remove the connection that is being unsubscribed.
    if (this.batchedSubscriptions[channel] && this.batchedSubscriptions[channel][subscriptionId]) {
      this.batchedSubscriptions[channel][subscriptionId] = this
        .batchedSubscriptions[channel][subscriptionId].filter(
          (e: SubscriptionInfo) => (e.connectionId !== connectionId),
        );
      subscribedConnections += this.batchedSubscriptions[channel][subscriptionId].length;
    }

    // If 0 connections are subscribed to the id for this channel after this unsubscribe, socks will
    // not forward any future messages in the channel with the id.
    // Delete from the set tracking the # of ids for the channel socks will forward messages for.
    if (this.subscribedIdsPerChannel[channel] && subscribedConnections === 0) {
      this.subscribedIdsPerChannel[channel].delete(subscriptionId);
      stats.gauge(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        this.subscribedIdsPerChannel[channel].size,
        {
          channel,
        },
      );
    }
  }

  /**
   * Remove deletes a connection from all data structures tracking subscriptions.
   * @param connectionId Connection id of the connection to be removed.
   */
  public remove(connectionId: string) {
    const subscriptionList = this.subscriptionLists[connectionId];
    if (subscriptionList) {
      subscriptionList.forEach((subscription) => {
        if (subscription.batched) {
          const idx = this.batchedSubscriptions[subscription.channel][subscription.id]
            .findIndex(
              (e: SubscriptionInfo) => e.connectionId === connectionId,
            );
          if (idx >= 0) {
            this.batchedSubscriptions[subscription.channel][subscription.id]
              .splice(idx, 1);
          }
        } else {
          const idx = this.subscriptions[subscription.channel][subscription.id]
            .findIndex(
              (e: SubscriptionInfo) => e.connectionId === connectionId,
            );
          if (idx >= 0) {
            this.subscriptions[subscription.channel][subscription.id]
              .splice(idx, 1);
          }
        }
      });
      delete this.subscriptionLists[connectionId];
    }

    this.subscribeRateLimiter.removeConnection(connectionId);
  }

  /**
   * validateSubscription validates a subscription messages
   * @param channel Channel being subscribed to
   * @param id Specific id within channel being subscribed to
   * @returns
   */
  private validateSubscription(channel: Channel, id?: string): boolean {
    // Only markets channel does not require an id to subscribe to.
    if (channel !== Channel.V4_MARKETS && id === undefined) {
      return false;
    }
    switch (channel) {
      case (Channel.V4_ACCOUNTS): {
        if (id === undefined) {
          return false;
        }
        const parts: string[] = id.split('/');

        // Id for subaccounts channel should be of the format {address}/{subaccountNumber}
        if (parts.length !== 2) {
          return false;
        }

        return true;
      }
      case (Channel.V4_MARKETS): {
        return true;
      }
      case (Channel.V4_ORDERBOOK):
      case (Channel.V4_TRADES):
        if (id === undefined) {
          return false;
        }
        return perpetualMarketRefresher.isValidPerpetualMarketTicker(id);
      case (Channel.V4_CANDLES): {
        if (id === undefined) {
          return false;
        }

        const {
          ticker,
          resolution,
        }: {
          ticker: string,
          resolution?: CandleResolution,
        } = this.parseCandleChannelId(id);
        if (!perpetualMarketRefresher.isValidPerpetualMarketTicker(ticker)) {
          return false;
        }

        return resolution !== undefined;
      }
      default: {
        throw new InvalidChannelError(channel);
      }
    }
  }

  /**
   * Normalizes subscription ids. If the id is undefined, returns the default id for the markets
   * channel, which is the only channel that does not have specific ids to subscribe to.
   * NOTE: Validation of the id and channel will happen in other functions.
   * @param id Subscription id to normalize.
   * @returns Normalized subscription id.
   */
  private normalizeSubscriptionId(id?: string): string {
    return id ?? V4_MARKETS_ID;
  }

  /**
   * Gets the initial response endpoint for a subscription based on the channel and id.
   * @param channel Channel to get the initial response endpoint for.
   * @param id Id of the subscription to get the initial response endpoint for.
   * @returns The endpoint if it exists, or undefined if the channel has no initial response
   * endpoint.
   */
  private getInitialEndpointForSubscription(channel: Channel, id?: string): string | undefined {
    switch (channel) {
      case (Channel.V4_TRADES): {
        if (id === undefined) {
          throw new Error('Invalid undefined channel');
        }

        return `${COMLINK_URL}/v4/trades/perpetualMarket/${id}`;
      }
      case (Channel.V4_MARKETS): {
        return `${COMLINK_URL}/v4/perpetualMarkets`;
      }
      case (Channel.V4_ORDERBOOK): {
        if (id === undefined) {
          throw new Error('Invalid undefined channel');
        }

        return `${COMLINK_URL}/v4/orderbooks/perpetualMarket/${id}`;
      }
      case (Channel.V4_CANDLES): {
        if (id === undefined) {
          throw new Error('Invalid undefined channel');
        }

        const {
          ticker,
          resolution,
        } : {
          ticker: string,
          resolution?: CandleResolution,
        } = this.parseCandleChannelId(id);
        // Resolution is guaranteed to be defined here because it is validated in
        // validateSubscription.
        return `${COMLINK_URL}/v4/candles/perpetualMarkets/${ticker}?resolution=${resolution!}`;
      }
      default: {
        throw new InvalidChannelError(channel);
      }
    }
  }

  private async getInitialResponseForSubaccountSubscription(
    id?: string,
    country?: string,
  ): Promise<string> {
    if (id === undefined) {
      throw new Error('Invalid undefined id');
    }

    // TODO(IND-508): Change this to match technical spec for persistent geo-blocking. This may
    // either have to replicate any blocking logic added on comlink, or re-direct to comlink to
    // determine if subscribing to a specific subaccount is blocked.
    if (country !== undefined && isRestrictedCountry(country)) {
      throw new BlockedError();
    }

    try {
      const {
        address,
        subaccountNumber,
      } : {
        address: string,
        subaccountNumber: string,
      } = this.parseSubaccountChannelId(id);

      const [
        subaccountsResponse,
        ordersResponse,
      ]: [
        string,
        string,
      ] = await Promise.all([
        axiosRequest({
          method: RequestMethod.GET,
          url: `${COMLINK_URL}/v4/addresses/${address}/subaccountNumber/${subaccountNumber}`,
          timeout: config.INITIAL_GET_TIMEOUT_MS,
          transformResponse: (res) => res,
        }),
        // TODO(DEC-1462): Use the /active-orders endpoint once it's added.
        axiosRequest({
          method: RequestMethod.GET,
          url: `${COMLINK_URL}/v4/orders?address=${address}&subaccountNumber=${subaccountNumber}&status=OPEN,UNTRIGGERED,BEST_EFFORT_OPENED`,
          timeout: config.INITIAL_GET_TIMEOUT_MS,
          transformResponse: (res) => res,
        }),
      ]);

      return JSON.stringify({
        ...JSON.parse(subaccountsResponse),
        orders: JSON.parse(ordersResponse),
      });
    } catch (error) {
      // The subaccounts API endpoint returns a 404 for subaccounts that are not indexed, however
      // such subaccounts can be subscribed to and events can be sent when the subaccounts are
      // indexed to an existing subscription.
      if (error instanceof AxiosSafeServerError && (error as AxiosSafeServerError).status === 404) {
        return EMPTY_INITIAL_RESPONSE;
      }
      // 403 indicates a blocked address. Throw a specific error for blocked addresses with a
      // specific error message detailing why the subscription failed due to a blocked address.
      if (error instanceof AxiosSafeServerError && (error as AxiosSafeServerError).status === 403) {
        throw new BlockedError();
      }
      throw error;
    }
  }

  private parseSubaccountChannelId(id: string): {
    address: string,
    subaccountNumber: string,
  } {
    const parts: string[] = id.split('/');
    const address: string = parts[0];
    const subaccountNumber: string = parts[1];
    return { address, subaccountNumber };
  }

  private parseCandleChannelId(id: string): {
    ticker: string,
    resolution?: CandleResolution,
  } {
    // Id for candles channel should be of the format {ticker}/{resolution}
    const parts: string[] = id.split('/');
    const ticker: string = parts[0];
    const resolutionString: string = parts[1];
    const resolution:
    CandleResolution | undefined = Object.values(CandleResolution)
      .find((cr) => cr === resolutionString);
    return { ticker, resolution };
  }

  /**
   * Gets the initial response for a channel.
   * @param channel Channel to get the initial response for.
   * @param id Id fo the subscription to get the initial response for.
   * @returns The initial response for the channel.
   */
  private async getInitialResponsesForChannels(
    channel: Channel,
    id?: string,
    country?: string,
  ): Promise<string> {
    if (channel === Channel.V4_ACCOUNTS) {
      return this.getInitialResponseForSubaccountSubscription(id, country);
    }
    const endpoint: string | undefined = this.getInitialEndpointForSubscription(channel, id);
    // If no endpoint exists, return an empty initial response.
    if (endpoint === undefined) {
      return EMPTY_INITIAL_RESPONSE;
    }

    return axiosRequest({
      method: RequestMethod.GET,
      url: endpoint,
      timeout: config.INITIAL_GET_TIMEOUT_MS,
      transformResponse: (res) => res, // Disables JSON parsing
    });
  }
}
