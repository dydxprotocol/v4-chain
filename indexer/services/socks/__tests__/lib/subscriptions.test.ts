import WebSocket from 'ws';
import { Channel, OutgoingMessageType } from '../../src/types';
import { Subscriptions } from '../../src/lib/subscription';
import { sendMessage, sendMessageString } from '../../src/helpers/wss';
import { RateLimiter } from '../../src/lib/rate-limit';
import {
  dbHelpers, testMocks, perpetualMarketRefresher, CandleResolution,
} from '@dydxprotocol-indexer/postgres';
import { btcTicker, invalidChannel, invalidTicker } from '../constants';
import { axiosRequest } from '../../src/lib/axios';
import { AxiosSafeServerError, makeAxiosSafeServerError } from '@dydxprotocol-indexer/base';
import { BlockedError } from '../../src/lib/errors';
import { isRestrictedCountry } from '@dydxprotocol-indexer/compliance';

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
  const mockSubaccountId: string = 'address/0';
  const invalidCandleResolution: string = 'candleResolution';
  const validIds: Record<Channel, string> = {
    [Channel.V4_ACCOUNTS]: mockSubaccountId,
    [Channel.V4_CANDLES]: `${btcTicker}/${CandleResolution.ONE_DAY}`,
    [Channel.V4_MARKETS]: defaultId,
    [Channel.V4_ORDERBOOK]: btcTicker,
    [Channel.V4_TRADES]: btcTicker,
  };
  const invalidIdsMap: Record<Exclude<Channel, Channel.V4_MARKETS>, string[]> = {
    [Channel.V4_ACCOUNTS]: [invalidTicker],
    [Channel.V4_CANDLES]: [
      `${invalidTicker}/${CandleResolution.ONE_DAY}`,
      `${btcTicker}/${invalidCandleResolution}`,
      btcTicker,
    ],
    [Channel.V4_ORDERBOOK]: [invalidTicker],
    [Channel.V4_TRADES]: [invalidTicker],
  };
  const initialResponseUrlPatterns: Record<Channel, string[] | undefined> = {
    [Channel.V4_ACCOUNTS]: [
      '/v4/addresses/.+/subaccountNumber/.+',
      '/v4/orders?.+OPEN,UNTRIGGERED,BEST_EFFORT_OPENED',
    ],
    [Channel.V4_CANDLES]: ['/v4/candles/perpetualMarkets/.+?resolution=.+'],
    [Channel.V4_MARKETS]: ['/v4/perpetualMarkets'],
    [Channel.V4_ORDERBOOK]: ['/v4/orderbooks/perpetualMarket/.+'],
    [Channel.V4_TRADES]: ['/v4/trades/perpetualMarket/.+'],
  };
  const initialMessage: Object = { a: 'b' };
  const restrictedCountry: string = 'US';
  const nonRestrictedCountry: string = 'AR';

  beforeAll(async () => {
    await dbHelpers.migrate();
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterAll(async () => {
    await dbHelpers.clearData();
    await dbHelpers.teardown();
  });

  beforeEach(() => {
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
    (isRestrictedCountry as jest.Mock).mockImplementation((country: string): boolean => {
      return country === restrictedCountry;
    });
  });

  describe('subscribe', () => {
    it.each([
      [Channel.V4_ACCOUNTS, validIds[Channel.V4_ACCOUNTS]],
      [Channel.V4_CANDLES, validIds[Channel.V4_CANDLES]],
      [Channel.V4_MARKETS, validIds[Channel.V4_MARKETS]],
      [Channel.V4_ORDERBOOK, validIds[Channel.V4_ORDERBOOK]],
      [Channel.V4_TRADES, validIds[Channel.V4_TRADES]],
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
        nonRestrictedCountry,
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
          nonRestrictedCountry,
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
            nonRestrictedCountry,
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
        nonRestrictedCountry,
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
        nonRestrictedCountry,
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
        nonRestrictedCountry,
      );

      expect(sendMessageMock).toHaveBeenCalledTimes(1);
      expect(sendMessageMock).toHaveBeenCalledWith(
        mockWs,
        connectionId,
        expect.objectContaining({
          connection_id: connectionId,
          type: 'error',
          message: expectedError.message,
        }));
      expect(subscriptions.subscriptions[Channel.V4_ACCOUNTS]).toBeUndefined();
      expect(subscriptions.subscriptionLists[connectionId]).toBeUndefined();
    });

    it('sends blocked error if subscribing to subaccount from restricted country', async () => {
      const expectedError: BlockedError = new BlockedError();
      await subscriptions.subscribe(
        mockWs,
        Channel.V4_ACCOUNTS,
        connectionId,
        initialMsgId,
        mockSubaccountId,
        false,
        restrictedCountry,
      );

      expect(sendMessageMock).toHaveBeenCalledTimes(1);
      expect(sendMessageMock).toHaveBeenCalledWith(
        mockWs,
        connectionId,
        expect.objectContaining({
          connection_id: connectionId,
          type: 'error',
          message: expectedError.message,
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
        nonRestrictedCountry,
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
      [Channel.V4_ACCOUNTS, validIds[Channel.V4_ACCOUNTS]],
      [Channel.V4_CANDLES, validIds[Channel.V4_CANDLES]],
      [Channel.V4_MARKETS, validIds[Channel.V4_MARKETS]],
      [Channel.V4_ORDERBOOK, validIds[Channel.V4_ORDERBOOK]],
      [Channel.V4_TRADES, validIds[Channel.V4_TRADES]],
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
        nonRestrictedCountry,
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
        nonRestrictedCountry,
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
          validIds[channel],
          false,
          nonRestrictedCountry,
        );
      }));

      for (const channel of Object.values(Channel)) {
        expect(subscriptions.subscriptions[channel][validIds[channel]]).toHaveLength(1);
        expect(subscriptions.subscriptions[channel][validIds[channel]]).toContainEqual(
          expect.objectContaining({ connectionId }),
        );
      }
      expect(subscriptions.subscriptionLists[connectionId]).toHaveLength(
        Object.values(Channel).length,
      );

      subscriptions.remove(connectionId);

      for (const channel of Object.values(Channel)) {
        expect(subscriptions.subscriptions[channel][validIds[channel]]).toHaveLength(0);
      }
      expect(subscriptions.subscriptionLists[connectionId]).toBeUndefined();
    });
  });
});
