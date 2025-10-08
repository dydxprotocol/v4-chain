import WebSocket from 'ws';
import { Channel, OutgoingMessageType } from '../../src/types';
import { Subscriptions } from '../../src/lib/subscription';
import { sendMessage, sendMessageString } from '../../src/helpers/wss';
import { RateLimiter } from '../../src/lib/rate-limit';
import {
  dbHelpers,
  testMocks,
  perpetualMarketRefresher,
  CandleResolution,
  MAX_PARENT_SUBACCOUNTS,
  blockHeightRefresher,
} from '@dydxprotocol-indexer/postgres';
import {
  btcTicker, ethTicker, invalidChannel, invalidTicker,
} from '../constants';
import { axiosRequest } from '../../src/lib/axios';
import {
  AxiosSafeServerError, makeAxiosSafeServerError, stats, setInstanceId,
} from '@dydxprotocol-indexer/base';
import { BlockedError } from '../../src/lib/errors';
import config from '../../src/config';

jest.mock('ws');
jest.mock('../../src/helpers/wss');
jest.mock('../../src/lib/axios');
jest.mock('@dydxprotocol-indexer/compliance');

describe('Subscriptions', () => {
  let subscriptions: Subscriptions;
  let mockWs: WebSocket;
  let sendMessageMock: jest.Mock;
  let sendMessageStringMock: jest.Mock;
  let rateLimiterSpy: jest.SpyInstance;
  let axiosRequestMock: jest.Mock;

  const connectionId: string = 'connectionId';
  const initialMsgId: number = 1;
  const defaultId: string = 'id';
  const defaultId1: string = 'id1';
  const mockSubaccountId: string = 'address/0';
  const mockSubaccountId1: string = 'address/1';
  const invalidCandleResolution: string = 'candleResolution';
  const singleIds: Record<Channel, string> = {
    [Channel.V4_ACCOUNTS]: mockSubaccountId,
    [Channel.V4_CANDLES]: `${btcTicker}/${CandleResolution.ONE_DAY}`,
    [Channel.V4_MARKETS]: defaultId,
    [Channel.V4_ORDERBOOK]: btcTicker,
    [Channel.V4_TRADES]: btcTicker,
    [Channel.V4_PARENT_ACCOUNTS]: mockSubaccountId,
    [Channel.V4_BLOCK_HEIGHT]: defaultId,
  };
  const multipleIds: Record<Channel, string[]> = {
    [Channel.V4_ACCOUNTS]: [mockSubaccountId, mockSubaccountId1],
    [Channel.V4_CANDLES]: [`${btcTicker}/${CandleResolution.ONE_DAY}`, `${btcTicker}/${CandleResolution.ONE_DAY}`],
    [Channel.V4_MARKETS]: [defaultId, defaultId1],
    [Channel.V4_ORDERBOOK]: [btcTicker, ethTicker],
    [Channel.V4_TRADES]: [btcTicker, ethTicker],
    [Channel.V4_PARENT_ACCOUNTS]: [mockSubaccountId, mockSubaccountId1],
    [Channel.V4_BLOCK_HEIGHT]: [defaultId, defaultId1],
  };
  const invalidIdsMap: Record<Channel, string[]> = {
    [Channel.V4_ACCOUNTS]: [invalidTicker],
    [Channel.V4_CANDLES]: [
      `${invalidTicker}/${CandleResolution.ONE_DAY}`,
      `${btcTicker}/${invalidCandleResolution}`,
      btcTicker,
    ],
    [Channel.V4_ORDERBOOK]: [invalidTicker],
    [Channel.V4_TRADES]: [invalidTicker],
    [Channel.V4_PARENT_ACCOUNTS]: [`address/${MAX_PARENT_SUBACCOUNTS}`],
    [Channel.V4_BLOCK_HEIGHT]: ['unused'],
    [Channel.V4_MARKETS]: ['unused'],
  };
  const initialResponseUrlPatterns: Record<Channel, string[] | undefined> = {
    [Channel.V4_ACCOUNTS]: [
      '/v4/addresses/.+/subaccountNumber/.+',
      '/v4/orders?.+subaccountNumber.+OPEN,UNTRIGGERED,BEST_EFFORT_OPENED',
      '/v4/orders?.+subaccountNumber.+BEST_EFFORT_CANCELED.+goodTilBlockAfter=[0-9]+',
    ],
    [Channel.V4_CANDLES]: ['/v4/candles/perpetualMarkets/.+?resolution=.+'],
    [Channel.V4_MARKETS]: ['/v4/perpetualMarkets'],
    [Channel.V4_ORDERBOOK]: ['/v4/orderbooks/perpetualMarket/.+'],
    [Channel.V4_TRADES]: ['/v4/trades/perpetualMarket/.+'],
    [Channel.V4_PARENT_ACCOUNTS]: [
      '/v4/addresses/.+/parentSubaccountNumber/.+',
      '/v4/orders/parentSubaccountNumber?.+parentSubaccountNumber.+OPEN,UNTRIGGERED,BEST_EFFORT_OPENED',
      '/v4/orders/parentSubaccountNumber?.+parentSubaccountNumber.+BEST_EFFORT_CANCELED.+goodTilBlockAfter=[0-9]+',
    ],
    [Channel.V4_BLOCK_HEIGHT]: ['v4/height'],
  };
  const initialMessage: Object = ['a', 'b'];
  const country: string = 'AR';

  beforeAll(async () => {
    await dbHelpers.migrate();
    await testMocks.seedData();
    await Promise.all([
      perpetualMarketRefresher.updatePerpetualMarkets(),
      blockHeightRefresher.updateBlockHeight(),
    ]);
    config.SERVICE_NAME = 'socks-test';
    await setInstanceId('test-instance-id');
  });

  afterAll(async () => {
    await dbHelpers.clearData();
    await dbHelpers.teardown();
  });

  beforeEach(() => {
    jest.useFakeTimers();
    (WebSocket as unknown as jest.Mock).mockClear();
    subscriptions = new Subscriptions();
    subscriptions.start(jest.fn());
    mockWs = new WebSocket(null);
    sendMessageMock = (sendMessage as jest.Mock);
    sendMessageStringMock = (sendMessageString as jest.Mock);
    rateLimiterSpy = jest.spyOn(RateLimiter.prototype, 'rateLimit');
    axiosRequestMock = (axiosRequest as jest.Mock);
    axiosRequestMock.mockClear();
    axiosRequestMock.mockImplementation(() => (JSON.stringify(initialMessage)));
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  describe('subscribe', () => {
    it.each([
      [Channel.V4_ACCOUNTS, singleIds[Channel.V4_ACCOUNTS]],
      [Channel.V4_CANDLES, singleIds[Channel.V4_CANDLES]],
      [Channel.V4_MARKETS, singleIds[Channel.V4_MARKETS]],
      [Channel.V4_ORDERBOOK, singleIds[Channel.V4_ORDERBOOK]],
      [Channel.V4_TRADES, singleIds[Channel.V4_TRADES]],
      [Channel.V4_PARENT_ACCOUNTS, singleIds[Channel.V4_PARENT_ACCOUNTS]],
      [Channel.V4_BLOCK_HEIGHT, singleIds[Channel.V4_BLOCK_HEIGHT]],
    ])('handles valid subscription request to channel %s', async (
      channel: Channel,
      id: string,
    ) => {
      await subscriptions.subscribe(
        mockWs,
        channel,
        connectionId,
        initialMsgId,
        id,
        false,
        country,
      );

      expect(sendMessageStringMock).toHaveBeenCalledTimes(1);
      expect(sendMessageStringMock).toHaveBeenCalledWith(
        mockWs,
        connectionId,
        expect.stringContaining(OutgoingMessageType.SUBSCRIBED),
      );
      expect(subscriptions.subscriptions[channel][id])
        .toContainEqual(expect.objectContaining({ connectionId }));
      expect(subscriptions.subscriptionLists[connectionId]).toHaveLength(1);
      expect(subscriptions.subscriptionLists[connectionId])
        .toContainEqual(expect.objectContaining({ channel, id }));

      const urlPatterns: string[] | undefined = initialResponseUrlPatterns[channel];
      if (urlPatterns !== undefined) {
        for (const urlPattern of urlPatterns) {
          expect(axiosRequestMock).toHaveBeenCalledWith(expect.objectContaining({
            url: expect.stringMatching(RegExp(urlPattern)),
            headers: {
              'cf-ipcountry': country,
            },
          }));
        }
      } else {
        expect(axiosRequestMock).not.toHaveBeenCalled();
      }
    });

    it.each([
      [Channel.V4_ACCOUNTS, invalidIdsMap[Channel.V4_ACCOUNTS]],
      [Channel.V4_CANDLES, invalidIdsMap[Channel.V4_CANDLES]],
      [Channel.V4_ORDERBOOK, invalidIdsMap[Channel.V4_ORDERBOOK]],
      [Channel.V4_TRADES, invalidIdsMap[Channel.V4_TRADES]],
      [Channel.V4_PARENT_ACCOUNTS, invalidIdsMap[Channel.V4_PARENT_ACCOUNTS]],
    ])('sends error message if invalid subscription request to channel %s', async (
      channel: Channel,
      invalidIds: string[],
    ) => {
      for (const id of invalidIds) {
        await subscriptions.subscribe(
          mockWs,
          channel,
          connectionId,
          initialMsgId,
          id,
          false,
        );

        expect(sendMessageMock).toHaveBeenCalledTimes(1);
        expect(sendMessageMock).toHaveBeenCalledWith(
          mockWs,
          connectionId,
          expect.objectContaining({
            connection_id: connectionId,
            type: 'error',
            message: `Invalid subscription id for channel: (${channel}-${id})`,
          }),
        );
        expect(subscriptions.subscriptions[channel]).toBeUndefined();

        sendMessageMock.mockClear();
      }
    });

    it('throws error if channel is invalid', async () => {
      await expect(
        async () => {
          await subscriptions.subscribe(
            mockWs,
            (invalidChannel as Channel),
            connectionId,
            initialMsgId,
            defaultId,
            false,
          );
        },
      ).rejects.toEqual(new Error(`Invalid channel: ${invalidChannel}`));
    });

    it('sends error message if rate limit exceeded', async () => {
      rateLimiterSpy.mockImplementation(() => 1);
      await subscriptions.subscribe(
        mockWs,
        Channel.V4_ACCOUNTS,
        connectionId,
        initialMsgId,
        mockSubaccountId,
        false,
      );

      expect(sendMessageMock).toHaveBeenCalledTimes(1);
      expect(sendMessageMock).toHaveBeenCalledWith(
        mockWs,
        connectionId,
        expect.objectContaining({
          message: expect.stringContaining('Too many subscribe attempts'),
        }));
      expect(subscriptions.subscriptions[Channel.V4_ACCOUNTS]).toBeUndefined();
      expect(subscriptions.subscriptionLists[connectionId]).toBeUndefined();
    });

    it('sends error message if initial message request fails', async () => {
      axiosRequestMock.mockImplementation(() => { throw Error(); });
      await subscriptions.subscribe(
        mockWs,
        Channel.V4_ACCOUNTS,
        connectionId,
        initialMsgId,
        mockSubaccountId,
        false,
      );

      expect(sendMessageMock).toHaveBeenCalledTimes(1);
      expect(sendMessageMock).toHaveBeenCalledWith(
        mockWs,
        connectionId,
        expect.objectContaining({
          connection_id: connectionId,
          type: 'error',
          message: expect.stringContaining(
            `Internal error, could not fetch data for subscription: ${Channel.V4_ACCOUNTS}`,
          ),
          channel: Channel.V4_ACCOUNTS,
          id: mockSubaccountId,
        }));
      expect(subscriptions.subscriptions[Channel.V4_ACCOUNTS]).toBeUndefined();
      expect(subscriptions.subscriptionLists[connectionId]).toBeUndefined();
    });

    it('sends blocked error message if initial message request fails with 403', async () => {
      const expectedError: BlockedError = new BlockedError();
      axiosRequestMock.mockImplementation(
        () => {
          throw new AxiosSafeServerError({
            data: {},
            status: 403,
            statusText: '',
          }, {});
        },
      );
      await subscriptions.subscribe(
        mockWs,
        Channel.V4_ACCOUNTS,
        connectionId,
        initialMsgId,
        mockSubaccountId,
        false,
        country,
      );

      expect(sendMessageMock).toHaveBeenCalledTimes(1);
      expect(sendMessageMock).toHaveBeenCalledWith(
        mockWs,
        connectionId,
        expect.objectContaining({
          connection_id: connectionId,
          type: 'error',
          message: expectedError.message,
          channel: Channel.V4_ACCOUNTS,
          id: mockSubaccountId,
        }));
      expect(subscriptions.subscriptions[Channel.V4_ACCOUNTS]).toBeUndefined();
      expect(subscriptions.subscriptionLists[connectionId]).toBeUndefined();
    });

    it('sends empty contents if initial message request fails with 404 for accounts', async () => {
      axiosRequestMock.mockImplementation(() => {
        return Promise.reject(makeAxiosSafeServerError(404, '', ''));
      });
      await subscriptions.subscribe(
        mockWs,
        Channel.V4_ACCOUNTS,
        connectionId,
        initialMsgId,
        mockSubaccountId,
        false,
        country,
      );

      expect(sendMessageStringMock).toHaveBeenCalledTimes(1);
      expect(sendMessageStringMock).toHaveBeenCalledWith(
        mockWs,
        connectionId,
        expect.stringContaining(OutgoingMessageType.SUBSCRIBED),
      );
      expect(subscriptions.subscriptions[Channel.V4_ACCOUNTS][mockSubaccountId])
        .toContainEqual(expect.objectContaining({ connectionId }));
      expect(subscriptions.subscriptionLists[connectionId]).toHaveLength(1);
      expect(subscriptions.subscriptionLists[connectionId])
        .toContainEqual(
          expect.objectContaining({ channel: Channel.V4_ACCOUNTS, id: mockSubaccountId }),
        );
    });
  });

  describe('unsubscribe', () => {
    it.each([
      [Channel.V4_ACCOUNTS, singleIds[Channel.V4_ACCOUNTS]],
      [Channel.V4_CANDLES, singleIds[Channel.V4_CANDLES]],
      [Channel.V4_MARKETS, singleIds[Channel.V4_MARKETS]],
      [Channel.V4_ORDERBOOK, singleIds[Channel.V4_ORDERBOOK]],
      [Channel.V4_TRADES, singleIds[Channel.V4_TRADES]],
      [Channel.V4_BLOCK_HEIGHT, singleIds[Channel.V4_BLOCK_HEIGHT]],
    ])('handles valid unsubscription request to channel %s', async (
      channel: Channel,
      id: string,
    ) => {
      await subscriptions.subscribe(
        mockWs,
        channel,
        connectionId,
        initialMsgId,
        id,
        false,
        country,
      );
      subscriptions.unsubscribe(
        connectionId,
        channel,
        id,
      );

      expect(subscriptions.subscriptions[channel][id]).toHaveLength(0);
      expect(subscriptions.subscriptionLists[connectionId]).toHaveLength(0);
    });

    it('is no-op if connection is not subscribed to channel and id', async () => {
      await subscriptions.subscribe(
        mockWs,
        Channel.V4_ACCOUNTS,
        connectionId,
        initialMsgId,
        mockSubaccountId,
        false,
        country,
      );
      subscriptions.unsubscribe(
        connectionId,
        Channel.V4_CANDLES,
        defaultId,
      );

      expect(subscriptions.subscriptions[Channel.V4_ACCOUNTS][mockSubaccountId]).toHaveLength(1);
      expect(subscriptions.subscriptions[Channel.V4_CANDLES]).toBeUndefined();
      expect(subscriptions.subscriptionLists[connectionId]).toHaveLength(1);
    });
  });

  describe('remove', () => {
    it('removes connection id from all subscriptions', async () => {
      await Promise.all(Object.values(Channel).map((channel: Channel): Promise<void> => {
        return subscriptions.subscribe(
          mockWs,
          channel,
          connectionId,
          initialMsgId,
          singleIds[channel],
          false,
        );
      }));

      for (const channel of Object.values(Channel)) {
        expect(subscriptions.subscriptions[channel][singleIds[channel]]).toHaveLength(1);
        expect(subscriptions.subscriptions[channel][singleIds[channel]]).toContainEqual(
          expect.objectContaining({ connectionId }),
        );
      }
      expect(subscriptions.subscriptionLists[connectionId]).toHaveLength(
        Object.values(Channel).length,
      );

      subscriptions.remove(connectionId);

      for (const channel of Object.values(Channel)) {
        expect(subscriptions.subscriptions[channel][singleIds[channel]]).toHaveLength(0);
      }
      expect(subscriptions.subscriptionLists[connectionId]).toBeUndefined();
    });
  });

  describe('emitLargestSubscriberMetric', () => {

    it('emits metrics for largest subscriber per channel', async () => {
      const statsSpy = jest.spyOn(stats, 'gauge');

      // Subscribe connection 1 to multiple channels
      await subscriptions.subscribe(
        mockWs,
        Channel.V4_ACCOUNTS,
        'connection1',
        initialMsgId,
        singleIds[Channel.V4_ACCOUNTS],
        false,
      );
      await subscriptions.subscribe(
        mockWs,
        Channel.V4_TRADES,
        'connection1',
        initialMsgId,
        singleIds[Channel.V4_TRADES],
        false,
      );
      await subscriptions.subscribe(
        mockWs,
        Channel.V4_ORDERBOOK,
        'connection1',
        initialMsgId,
        singleIds[Channel.V4_ORDERBOOK],
        false,
      );

      // verify metrics emitted for each channel with subscriptions
      // accounts: 1
      expect(statsSpy).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        1,
        {
          channel: Channel.V4_ACCOUNTS,
          instance: 'test-instance-id',
        },
      );
      // trades: 1
      expect(statsSpy).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        1,
        {
          channel: Channel.V4_TRADES,
          instance: 'test-instance-id',
        },
      );
      // orders: 1
      expect(statsSpy).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        1,
        {
          channel: Channel.V4_ORDERBOOK,
          instance: 'test-instance-id',
        },
      );
      jest.advanceTimersByTime(config.LARGEST_SUBSCRIBER_METRIC_INTERVAL_MS);
      // accounts: 1
      expect(statsSpy).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        1,
        {
          channel: Channel.V4_ACCOUNTS,
          instance: 'test-instance-id',
        },
      );

      // trades: 1
      expect(statsSpy).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        1,
        {
          channel: Channel.V4_TRADES,
          instance: 'test-instance-id',
        },
      );

      // books: 1
      expect(statsSpy).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        1,
        {
          channel: Channel.V4_ORDERBOOK,
          instance: 'test-instance-id',
        },
      );

      // Subscribe connection 2 to fewer channels
      // for each id in multipleIds[Channel.V4_ACCOUNTS], subscribe to the channel
      for (const id of multipleIds[Channel.V4_ACCOUNTS]) {
        await subscriptions.subscribe(
          mockWs,
          Channel.V4_ACCOUNTS,
          'connection2',
          initialMsgId,
          id,
          false,
        );
      }
      // for each id in multipleIds[Channel.V4_TRADES], subscribe to the channel
      for (const id of multipleIds[Channel.V4_TRADES]) {
        await subscriptions.subscribe(
          mockWs,
          Channel.V4_TRADES,
          'connection2',
          initialMsgId,
          id,
          false,
        );
      }

      // Advance timers to trigger metric emission
      jest.advanceTimersByTime(config.LARGEST_SUBSCRIBER_METRIC_INTERVAL_MS);

      // accounts: 2
      expect(statsSpy).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        2,
        {
          channel: Channel.V4_ACCOUNTS,
          instance: 'test-instance-id',
        },
      );

      // orders: 2
      expect(statsSpy).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.subscriptions.channel_size`,
        2,
        {
          channel: Channel.V4_TRADES,
          instance: 'test-instance-id',
        },
      );
    });

    it('does not emit metrics when no subscriptions exist', () => {
      const statsSpy = jest.spyOn(stats, 'gauge');

      // Advance timers to trigger metric emission
      jest.advanceTimersByTime(config.LARGEST_SUBSCRIBER_METRIC_INTERVAL_MS);

      // Should not emit any largest_subscriber metrics
      expect(statsSpy).not.toHaveBeenCalledWith(
        expect.stringContaining('largest_subscriber'),
        expect.anything(),
        expect.anything(),
      );
    });
  });
});
