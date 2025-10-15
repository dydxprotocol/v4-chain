import {
  AxiosSafeServerError,
  getInstanceId,
  logger,
  stats,
} from '@dydxprotocol-indexer/base';
import { GeoOriginHeaders } from '@dydxprotocol-indexer/compliance';
import {
  APIOrderStatus,
  BestEffortOpenedStatus,
  blockHeightRefresher,
  CHILD_SUBACCOUNT_MULTIPLIER,
  CandleResolution,
  MAX_PARENT_SUBACCOUNTS,
  OrderStatus,
  perpetualMarketRefresher,
  OrderFromDatabase,
} from '@dydxprotocol-indexer/postgres';
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
import { V4_BLOCK_HEIGHT_ID, V4_MARKETS_ID, WS_CLOSE_CODE_POLICY_VIOLATION } from './constants';
import { BlockedError, InvalidChannelError } from './errors';
import { RateLimiter } from './rate-limit';

const COMLINK_URL: string = `http://${config.COMLINK_URL}`;
const EMPTY_INITIAL_RESPONSE: string = '{}';
const VALID_ORDER_STATUS_FOR_INITIAL_SUBACCOUNT_RESPONSE: APIOrderStatus[] = [
  OrderStatus.OPEN,
  OrderStatus.UNTRIGGERED,
  BestEffortOpenedStatus.BEST_EFFORT_OPENED,
];

const VALID_ORDER_STATUS: string = VALID_ORDER_STATUS_FOR_INITIAL_SUBACCOUNT_RESPONSE.join(',');

const CHANNEL_CONNECTION_LIMITS: { [channel: string]: number } = {
  [Channel.V4_ACCOUNTS]: config.V4_ACCOUNTS_CHANNEL_LIMIT,
  [Channel.V4_BLOCK_HEIGHT]: 1,
  [Channel.V4_CANDLES]: config.V4_CANDLES_CHANNEL_LIMIT,
  [Channel.V4_MARKETS]: config.V4_MARKETS_CHANNEL_LIMIT,
  [Channel.V4_ORDERBOOK]: config.V4_ORDERBOOK_CHANNEL_LIMIT,
  [Channel.V4_PARENT_ACCOUNTS]: config.V4_PARENT_ACCOUNTS_CHANNEL_LIMIT,
  [Channel.V4_TRADES]: config.V4_TRADES_CHANNEL_LIMIT,
};

export class Subscriptions {
  private forwardMessage?: (message: MessageToForward, connectionId: string) => number;
  private subscriptionMetricsInterval?: NodeJS.Timeout;
  private subscribeRateLimiter: RateLimiter;
  public batchedSubscriptions: { [channel: string]: { [id: string]: SubscriptionInfo[] } };
  public subsByChannelByConnectionId: { [channel: string]: { [connectionId: string]: number } };
  public subscriptionLists: { [connectionId: string]: Subscription[] };
  public subscriptions: { [channel: string]: { [id: string]: SubscriptionInfo[] } };

  constructor() {
    this.subscriptionLists = {};
    this.subscriptions = {};
    this.batchedSubscriptions = {};
    this.subscribeRateLimiter = new RateLimiter({
      points: config.RATE_LIMIT_SUBSCRIBE_POINTS,
      durationMs: config.RATE_LIMIT_SUBSCRIBE_DURATION_MS,
    });
    this.subsByChannelByConnectionId = {};
    this.forwardMessage = undefined;
  }

  public start(forwardMessage: (message: MessageToForward, connectionId: string) => number): void {
    this.forwardMessage = forwardMessage;

    this.subscriptionMetricsInterval = setInterval(
      () => this.emitSubscriptionMetrics(),
      config.SUBSCRIPTION_METRIC_INTERVAL_MS,
    );
  }

  public incrementSubscriptions(channel: Channel, connectionId: string): number {
    this.subsByChannelByConnectionId[channel] ??= {};
    this.subsByChannelByConnectionId[channel][connectionId] ??= 0;
    this.subsByChannelByConnectionId[channel][connectionId] += 1;
    return this.subsByChannelByConnectionId[channel][connectionId];
  }

  public decrementSubscriptions(channel: Channel, connectionId: string): number {
    this.subsByChannelByConnectionId[channel] ??= {};
    this.subsByChannelByConnectionId[channel][connectionId] ??= 0;
    this.subsByChannelByConnectionId[channel][connectionId] -= 1;
    if (this.subsByChannelByConnectionId[channel][connectionId] < 0) {
      this.subsByChannelByConnectionId[channel][connectionId] = 0;
    }
    return this.subsByChannelByConnectionId[channel][connectionId];
  }

  public removeSubscriptions(connectionId: string): void {
    Object.keys(this.subsByChannelByConnectionId).forEach((channel) => {
      if (this.subsByChannelByConnectionId[channel][connectionId] !== undefined) {
        delete this.subsByChannelByConnectionId[channel][connectionId];
      }
    });
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
    geoOriginHeaders?: GeoOriginHeaders,
  ): Promise<void> {
    const activeSubscriptions = this.incrementSubscriptions(channel, connectionId);

    if (this.forwardMessage === undefined) {
      this.decrementSubscriptions(channel, connectionId);
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
      this.decrementSubscriptions(channel, connectionId);
      return;
    }

    const channelSubscriptionsLimit = CHANNEL_CONNECTION_LIMITS[channel];

    if (activeSubscriptions > channelSubscriptionsLimit) {
      stats.increment(
        `${config.SERVICE_NAME}.subscriptions_limit_reached`,
        1,
        undefined,
        { channel, instance: getInstanceId() },
      );

      sendMessage(
        ws,
        connectionId,
        createErrorMessage(
          `Per-connection subscription limit reached for ${channel} (limit=${channelSubscriptionsLimit}).`,
          connectionId,
          messageId,
        ),
      );
      this.decrementSubscriptions(channel, connectionId);
      return;
    }

    const subscriptionId: string = this.normalizeSubscriptionId(channel, id);
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
      this.decrementSubscriptions(channel, connectionId);
      return;
    }

    let initialResponse: string;
    const startGetInitialResponse: number = Date.now();
    try {
      initialResponse = await this.getInitialResponsesForChannels(channel, id, geoOriginHeaders);
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
          channel,
          id,
        ),
      );

      stats.increment(
        `${config.SERVICE_NAME}.initial_response_error`,
        1,
        undefined,
        {
          instance: getInstanceId(),
          channel,
        },
      );
      this.decrementSubscriptions(channel, connectionId);
      return;
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.initial_response_get`,
        Date.now() - startGetInitialResponse,
        undefined,
        { channel, instance: getInstanceId() },
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
          channel,
          id,
        ),
      );
      this.decrementSubscriptions(channel, connectionId);
      return;
    }

    this.subscriptionLists[connectionId].push(subscription);
    this.subscriptions[channel] ??= {};
    this.subscriptions[channel][subscriptionId] ??= [];

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
    stats.timing(
      `${config.SERVICE_NAME}.subscribe_send_message`,
      Date.now() - startSend,
      { instance: getInstanceId() },
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

    let removed = false;

    // remove subscription from subscription list
    const subscriptionId: string = this.normalizeSubscriptionId(channel, id);
    if (this.subscriptionLists[connectionId]) {
      const idx = this.subscriptionLists[connectionId]
        .findIndex((e: Subscription) => e.channel === channel && e.id === subscriptionId);
      if (idx >= 0) {
        this.subscriptionLists[connectionId].splice(idx, 1);
        removed = true;
      }
    }

    // remove subscription from batched and non-batched subscriptions
    for (const subscriptions of [this.subscriptions[channel], this.batchedSubscriptions[channel]]) {
      if (subscriptions && subscriptions[subscriptionId]) {
        const idx = subscriptions[subscriptionId]
          .findIndex((e: SubscriptionInfo) => e.connectionId === connectionId);
        if (idx >= 0) {
          subscriptions[subscriptionId].splice(idx, 1);
          removed = true;
        }
      }
    }

    if (removed) {
      this.decrementSubscriptions(channel, connectionId);
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
        for (const subscriptions of [
          this.batchedSubscriptions[subscription.channel],
          this.subscriptions[subscription.channel],
        ]) {
          if (subscriptions && subscriptions[subscription.id]) {
            const idx = subscriptions[subscription.id]
              .findIndex((e: SubscriptionInfo) => e.connectionId === connectionId);
            if (idx >= 0) {
              subscriptions[subscription.id].splice(idx, 1);
            }
          }
        }
      });
      delete this.subscriptionLists[connectionId];
    }
    this.removeSubscriptions(connectionId);
    this.subscribeRateLimiter.removeConnection(connectionId);
  }

  /**
   * validateSubscription validates a subscription messages
   * @param channel Channel being subscribed to
   * @param id Specific id within channel being subscribed to
   * @returns
   */
  private validateSubscription(channel: Channel, id?: string): boolean {
    // Only markets & block height channels do not require an id to subscribe to.
    switch (channel) {
      case (Channel.V4_BLOCK_HEIGHT):
      case (Channel.V4_MARKETS):
        return true;
      default:
        if (id === undefined) {
          return false;
        }
        break;
    }
    switch (channel) {
      case (Channel.V4_ACCOUNTS): {
        return this.validateSubaccountChannelId(
          id,
          MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER,
        );
      }
      case (Channel.V4_CANDLES): {
        const { ticker, resolution } = this.parseCandleChannelId(id);
        if (!perpetualMarketRefresher.isValidPerpetualMarketTicker(ticker)) {
          return false;
        }
        return resolution !== undefined;
      }
      case (Channel.V4_PARENT_ACCOUNTS): {
        return this.validateSubaccountChannelId(id, MAX_PARENT_SUBACCOUNTS);
      }
      case (Channel.V4_ORDERBOOK):
      case (Channel.V4_TRADES):
        return perpetualMarketRefresher.isValidPerpetualMarketTicker(id);
      default: {
        throw new InvalidChannelError(channel);
      }
    }
  }

  /**
   * Normalizes subscription ids. If the id is undefined, returns the default id for the markets
   * channel or block height channel which are the only channels that don't
   * have specific ids to subscribe to.
   * NOTE: Validation of the id and channel will happen in other functions.
   * @param id Subscription id to normalize.
   * @returns Normalized subscription id.
   */
  private normalizeSubscriptionId(channel: Channel, id?: string): string {
    if (channel === Channel.V4_BLOCK_HEIGHT) {
      return id ?? V4_BLOCK_HEIGHT_ID;
    }
    return id ?? V4_MARKETS_ID;
  }

  private validateSubaccountChannelId(id: string, maxSubaccountNumber: number): boolean {
    // Id for subaccounts channel should be of the format {address}/{subaccountNumber}
    const parts: string[] = id.split('/');
    if (parts.length !== 2) {
      return false;
    }

    if (Number.isNaN(Number(parts[1]))) {
      return false;
    }

    return Number(parts[1]) < maxSubaccountNumber;
  }

  /**
   * Gets the initial response endpoint for a subscription based on the channel and id.
   * @param channel Channel to get the initial response endpoint for.
   * @param id Id of the subscription to get the initial response endpoint for.
   * @returns The endpoint if it exists, or undefined if the channel has no initial response
   * endpoint.
   */
  private getInitialEndpointForSubscription(
    channel: Exclude<Channel, Channel.V4_ACCOUNTS | Channel.V4_PARENT_ACCOUNTS>,
    id?: string,
  ): string | undefined {
    switch (channel) {
      case (Channel.V4_BLOCK_HEIGHT): {
        return `${COMLINK_URL}/v4/height`;
      }
      case (Channel.V4_MARKETS): {
        return `${COMLINK_URL}/v4/perpetualMarkets`;
      }
      case (Channel.V4_TRADES): {
        if (id === undefined) { throw new Error('Invalid undefined channel'); }
        return `${COMLINK_URL}/v4/trades/perpetualMarket/${id}`;
      }
      case (Channel.V4_ORDERBOOK): {
        if (id === undefined) { throw new Error('Invalid undefined channel'); }
        return `${COMLINK_URL}/v4/orderbooks/perpetualMarket/${id}`;
      }
      case (Channel.V4_CANDLES): {
        if (id === undefined) { throw new Error('Invalid undefined channel'); }
        const {
          ticker,
          resolution,
        }: {
          ticker: string,
          resolution?: CandleResolution,
        } = this.parseCandleChannelId(id);
        // Resolution is guaranteed to be defined here because it is validated in
        // validateSubscription.
        return `${COMLINK_URL}/v4/candles/perpetualMarkets/${ticker}?resolution=${resolution!}`;
      }
      default: throw new InvalidChannelError(channel);
    }
  }

  // TODO deduplicate with getInitialResponseForParentSubaccountSubscription
  private async getInitialResponseForSubaccountSubscription(
    id: string,
    geoOriginHeaders?: GeoOriginHeaders,
  ): Promise<string> {

    try {
      const { address, subaccountNumber } = this.parseSubaccountChannelId(id);
      const blockHeightString: string = await blockHeightRefresher.getLatestBlockHeight();
      const blockHeight: number = parseInt(blockHeightString, 10);

      const [
        subaccountsResponse,
        ordersResponse,
        currentBestEffortCanceledOrdersResponse,
      ]: string[] = await Promise.all([
        axiosRequest({
          method: RequestMethod.GET,
          url: `${COMLINK_URL}/v4/addresses/${address}/subaccountNumber/${subaccountNumber}`,
          timeout: config.INITIAL_GET_TIMEOUT_MS,
          headers: geoOriginHeaders || {},
          transformResponse: (res) => res,
        }),
        // TODO(DEC-1462): Use the /active-orders endpoint once it's added.
        axiosRequest({
          method: RequestMethod.GET,
          url: `${COMLINK_URL}/v4/orders?address=${address}&subaccountNumber=${subaccountNumber}&status=${VALID_ORDER_STATUS}`,
          timeout: config.INITIAL_GET_TIMEOUT_MS,
          headers: geoOriginHeaders || {},
          transformResponse: (res) => res,
        }),
        axiosRequest({
          method: RequestMethod.GET,
          url: `${COMLINK_URL}/v4/orders?address=${address}&subaccountNumber=${subaccountNumber}&status=BEST_EFFORT_CANCELED&goodTilBlockAfter=${Math.max(blockHeight - 20, 1)}`,
          timeout: config.INITIAL_GET_TIMEOUT_MS,
          headers: geoOriginHeaders || {},
          transformResponse: (res) => res,
        }),
      ]);

      const orders: OrderFromDatabase[] = JSON.parse(ordersResponse);
      const currentBestEffortCanceledOrders: OrderFromDatabase[] = JSON.parse(
        currentBestEffortCanceledOrdersResponse,
      );
      const allOrders: OrderFromDatabase[] = orders.concat(currentBestEffortCanceledOrders);

      return JSON.stringify({
        ...JSON.parse(subaccountsResponse),
        orders: allOrders,
        blockHeight: blockHeightString,
      });
    } catch (error) {
      logger.error({
        at: 'getInitialResponseForSubaccountSubscription',
        message: 'Error on getting initial response for subaccount subscription',
        id,
        error,
      });
      // The subaccount may initially be invalid but become valid later
      if (error instanceof AxiosSafeServerError && (error as AxiosSafeServerError).status === 404) {
        return EMPTY_INITIAL_RESPONSE;
      }
      if (error instanceof AxiosSafeServerError && (error as AxiosSafeServerError).status === 403) {
        throw new BlockedError();
      }
      throw error;
    }
  }

  private async getInitialResponseForParentSubaccountSubscription(
    id: string,
    geoOriginHeaders?: GeoOriginHeaders,
  ): Promise<string> {

    try {
      const {
        address,
        subaccountNumber,
      }: {
        address: string,
        subaccountNumber: string,
      } = this.parseSubaccountChannelId(id);

      const blockHeight: string = await blockHeightRefresher.getLatestBlockHeight();
      const numBlockHeight: number = parseInt(blockHeight, 10);

      const [
        subaccountsResponse,
        ordersResponse,
        currentBestEffortCanceledOrdersResponse,
      ]: string[] = await Promise.all([
        axiosRequest({
          method: RequestMethod.GET,
          url: `${COMLINK_URL}/v4/addresses/${address}/parentSubaccountNumber/${subaccountNumber}`,
          timeout: config.INITIAL_GET_TIMEOUT_MS,
          headers: geoOriginHeaders || {},
          transformResponse: (res) => res,
        }),
        axiosRequest({
          method: RequestMethod.GET,
          url: `${COMLINK_URL}/v4/orders/parentSubaccountNumber?address=${address}&parentSubaccountNumber=${subaccountNumber}&status=${VALID_ORDER_STATUS}`,
          timeout: config.INITIAL_GET_TIMEOUT_MS,
          headers: geoOriginHeaders || {},
          transformResponse: (res) => res,
        }),
        axiosRequest({
          method: RequestMethod.GET,
          url: `${COMLINK_URL}/v4/orders/parentSubaccountNumber?address=${address}&parentSubaccountNumber=${subaccountNumber}&status=BEST_EFFORT_CANCELED&goodTilBlockAfter=${Math.max(numBlockHeight - 20, 1)}`,
          timeout: config.INITIAL_GET_TIMEOUT_MS,
          headers: geoOriginHeaders || {},
          transformResponse: (res) => res,
        }),
      ]);

      const orders: OrderFromDatabase[] = JSON.parse(ordersResponse);
      const currentBestEffortCanceledOrders: OrderFromDatabase[] = JSON.parse(
        currentBestEffortCanceledOrdersResponse,
      );
      const allOrders: OrderFromDatabase[] = orders.concat(currentBestEffortCanceledOrders);

      return JSON.stringify({
        ...JSON.parse(subaccountsResponse),
        orders: allOrders,
        blockHeight,
      });
    } catch (error) {
      logger.error({
        at: 'getInitialResponseForParentSubaccountSubscription',
        message: 'Error on getting initial response for subaccount subscription',
        id,
        error,
      });
      // The subaccounts may initially be invalid but become valid later
      if (error instanceof AxiosSafeServerError && (error as AxiosSafeServerError).status === 404) {
        return EMPTY_INITIAL_RESPONSE;
      }
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
    const resolution: CandleResolution | undefined = Object.values(CandleResolution)
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
    geoOriginHeaders?: GeoOriginHeaders,
  ): Promise<string> {
    let endpoint: string | undefined;
    switch (channel) {
      case (Channel.V4_ACCOUNTS):
        if (id === undefined) { throw new Error('Invalid undefined id'); }
        return this.getInitialResponseForSubaccountSubscription(id, geoOriginHeaders);
      case (Channel.V4_PARENT_ACCOUNTS):
        if (id === undefined) { throw new Error('Invalid undefined id'); }
        return this.getInitialResponseForParentSubaccountSubscription(id, geoOriginHeaders);
      default:
        endpoint = this.getInitialEndpointForSubscription(channel, id);
        if (endpoint === undefined) { return EMPTY_INITIAL_RESPONSE; }
        return axiosRequest({
          method: RequestMethod.GET,
          url: endpoint,
          timeout: config.INITIAL_GET_TIMEOUT_MS,
          headers: geoOriginHeaders || {},
          transformResponse: (res) => res, // Disables JSON parsing
        });
    }
  }

  private emitSubscriptionMetrics(): void {
    const maxSubscriptionsByChannel: { [channel: string]: number } = {};
    const maxIdByChannel: { [channel: string]: string } = {};
    const subscriptionsByChannel: { [channel: string]: number } = {};

    Object.entries(this.subsByChannelByConnectionId).forEach(
      ([channel, subscribedIdsByConnectionId]) => {
        let maxId: string = '';
        let maxSubscriptions: number = 0;
        subscriptionsByChannel[channel] = 0;
        Object.entries(subscribedIdsByConnectionId).forEach(([connectionId, subscriptions]) => {
          subscriptionsByChannel[channel] += subscriptions;
          if (subscriptions > (maxSubscriptions || 0)) {
            maxSubscriptions = subscriptions;
            maxId = connectionId;
          }
        });
        maxIdByChannel[channel] = maxId;
        maxSubscriptionsByChannel[channel] = maxSubscriptions;
      });

    const instanceId = getInstanceId();

    Object.entries(maxSubscriptionsByChannel).forEach(([channel, count]) => {
      stats.gauge(
        `${config.SERVICE_NAME}.largest_subscriber`,
        count,
        {
          channel,
          instance: instanceId,
        },
      );
    });

    Object.entries(subscriptionsByChannel).forEach(([channel, count]) => {
      stats.gauge(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        count,
        {
          channel,
          instance: instanceId,
        },
      );
    });

    if (Object.keys(maxSubscriptionsByChannel).length > 0) {
      logger.info({
        at: 'Subscriptions#emitSubscriptionMetrics',
        message: 'Max subscriptions by channel',
        maxSubscriptionsByChannel,
      });
    }

    if (Object.keys(maxIdByChannel).length > 0) {
      logger.info({
        at: 'Subscriptions#emitSubscriptionMetrics',
        message: 'Max id by channel',
        maxIdByChannel,
      });
    }
  }

  public stop(): void {
    if (this.subscriptionMetricsInterval !== undefined) {
      clearInterval(this.subscriptionMetricsInterval);
    }
  }
}
