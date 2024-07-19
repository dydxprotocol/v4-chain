import _ from 'lodash';
import { Callback } from 'redis';
import { deleteAllAsync, hGetAsync } from '../../src/helpers/redis';
import {
  redis as client,

} from '../helpers/utils';
import {
  updatePriceLevel,
  getOrderBookLevels,
  getKey,
  deleteZeroPriceLevel,
} from '../../src/caches/orderbook-levels-cache';
import { OrderSide } from '@dydxprotocol-indexer/postgres';
import { OrderbookLevels, PriceLevel } from '../../src/types';
import { InvalidOptionsError, InvalidPriceLevelUpdateError } from '../../src/errors';
import { logger } from '@dydxprotocol-indexer/base';

describe('orderbookLevelsCache', () => {
  const ticker: string = 'BTC-USD';

  beforeEach(async () => {
    await deleteAllAsync(client);
    jest.spyOn(logger, 'crit');
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('updatePriceLevel', () => {
    it.each([
      ['BTC-USD bid, $45,100, +1000 quantums', 'BTC-USD', OrderSide.BUY, '45100', '1000'],
      ['ETH-USD bid, $3,350, +300 quantums', 'ETH-USD', OrderSide.BUY, '3350', '300'],
      ['BTC-USD ask, $51,100, +1500 quantums', 'BTC-USD', OrderSide.BUY, '51100', '1500'],
      ['ETH-USD ask, $3,950, +320 quantums', 'ETH-USD', OrderSide.BUY, '3950', '320'],
    ])('updates the total size in quantums for a price/side/ticker, new price level: %s', async (
      _name: string,
      tickerForTest: string,
      side: OrderSide,
      humanPrice: string,
      sizeDeltaInQuantums: string,
    ) => {

      const updatedQuantums: number = await updatePriceLevel({
        ticker: tickerForTest,
        side,
        humanPrice,
        sizeDeltaInQuantums,
        client,
      });
      const orderbook: OrderbookLevels = await getOrderBookLevels(tickerForTest, client);

      expect(updatedQuantums.toString()).toEqual(sizeDeltaInQuantums);
      let levels: PriceLevel[];
      if (side === OrderSide.BUY) {
        levels = orderbook.bids;
      } else {
        levels = orderbook.asks;
      }
      expect(levels).toMatchObject([{
        humanPrice,
        quantums: sizeDeltaInQuantums,
      }]);
    });

    it.each([
      [
        'BTC-USD bid, $45,100, +1000 quantums, existing 1500 quantums',
        'BTC-USD',
        OrderSide.BUY,
        '45100',
        '1000',
        '1500',
        '2500',
      ],
      [
        'BTC-USD bid, $45,100, -1000 quantums, existing 1500 quantums',
        'BTC-USD',
        OrderSide.BUY,
        '45100',
        '-1000',
        '1500',
        '500',
      ],
      [
        'BTC-USD ask, $62,100, +2000 quantums, existing 3000 quantums',
        'BTC-USD',
        OrderSide.BUY,
        '62100',
        '2000',
        '3000',
        '5000',
      ],
      [
        'BTC-USD bid, $62,100, -2000 quantums, existing 3000 quantums',
        'BTC-USD',
        OrderSide.BUY,
        '62100',
        '-2000',
        '3000',
        '1000',
      ],
    ])('updates the total size in quantums for a price/side/ticker, existing price level: %s', async (
      _name: string,
      tickerForTest: string,
      side: OrderSide,
      humanPrice: string,
      sizeDeltaInQuantums: string,
      existingQuantums: string,
      expectedQuantums: string,
    ) => {
      // Set existing quantums for the level
      await updatePriceLevel({
        ticker: tickerForTest,
        side,
        humanPrice,
        sizeDeltaInQuantums: existingQuantums,
        client,
      });

      const updatedQuantums: number = await updatePriceLevel({
        ticker: tickerForTest,
        side,
        humanPrice,
        sizeDeltaInQuantums,
        client,
      });
      const orderbook: OrderbookLevels = await getOrderBookLevels(tickerForTest, client);

      expect(updatedQuantums.toString()).toEqual(expectedQuantums);
      let levels: PriceLevel[];
      if (side === OrderSide.BUY) {
        levels = orderbook.bids;
      } else {
        levels = orderbook.asks;
      }
      expect(levels).toMatchObject([{
        humanPrice,
        quantums: expectedQuantums,
      }]);
    });

    it('can update price level to zero quantums', async () => {
      const humanPrice: string = '50000';
      const quantums: string = '3000';
      const updateQuantums: string = '-3000';
      // set up initial quantums
      await updatePriceLevel({
        ticker,
        side: OrderSide.BUY,
        humanPrice,
        sizeDeltaInQuantums: quantums,
        client,
      });

      let orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client);
      expect(orderbookLevels.bids).toMatchObject([{
        humanPrice,
        quantums,
      }]);

      // Update with delta to set quantums to 0
      const result: number = await updatePriceLevel({
        ticker,
        side: OrderSide.BUY,
        humanPrice,
        sizeDeltaInQuantums: updateQuantums,
        client,
      });
      orderbookLevels = await getOrderBookLevels(ticker, client);

      expect(result).toEqual(0);
      expect(orderbookLevels.bids).toEqual([]);
    });

    it('throws error if update will cause quantums to be negative', async () => {
      const humanPrice: string = '50000';
      const quantums: string = '1000';
      const invalidDelta: string = '-2000';
      const resultingQuantums: string = '-1000';
      // Set existing quantums for the level
      await updatePriceLevel({
        ticker,
        side: OrderSide.BUY,
        humanPrice,
        sizeDeltaInQuantums: quantums,
        client,
      });

      // Test that an error is thrown if the update results in a negative quantums for the price
      // level
      await expect(updatePriceLevel({
        ticker,
        side: OrderSide.BUY,
        humanPrice,
        sizeDeltaInQuantums: invalidDelta,
        client,
      })).rejects.toBeInstanceOf(InvalidPriceLevelUpdateError);
      expect(logger.crit).toHaveBeenCalledTimes(1);
      await expect(updatePriceLevel({
        ticker,
        side: OrderSide.BUY,
        humanPrice,
        sizeDeltaInQuantums: invalidDelta,
        client,
      })).rejects.toEqual(expect.objectContaining({
        message: expect.stringContaining(resultingQuantums),
      }));

      // Expect that the value in the orderbook is unchanged
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client);
      expect(orderbookLevels.bids).toMatchObject([{
        humanPrice,
        quantums,
      }]);
    });
  });

  describe('getOrderbookLevels', () => {
    beforeEach(async () => {
      // Setup some price levels
      // Bids
      await Promise.all([
        updatePriceLevel({
          ticker,
          side: OrderSide.BUY,
          humanPrice: '45100',
          sizeDeltaInQuantums: '2000',
          client,
        }),
        updatePriceLevel({
          ticker,
          side: OrderSide.BUY,
          humanPrice: '45200',
          sizeDeltaInQuantums: '5000',
          client,
        }),
        updatePriceLevel({
          ticker,
          side: OrderSide.BUY,
          humanPrice: '49300',
          sizeDeltaInQuantums: '6000',
          client,
        }),
        // Asks
        updatePriceLevel({
          ticker,
          side: OrderSide.SELL,
          humanPrice: '52200',
          sizeDeltaInQuantums: '8000',
          client,
        }),
        updatePriceLevel({
          ticker,
          side: OrderSide.SELL,
          humanPrice: '54200',
          sizeDeltaInQuantums: '4000',
          client,
        }),
        updatePriceLevel({
          ticker,
          side: OrderSide.SELL,
          humanPrice: '59300',
          sizeDeltaInQuantums: '1000',
          client,
        }),
        // Zero quantums bid/ask
        updatePriceLevel({
          ticker,
          side: OrderSide.BUY,
          humanPrice: '40000',
          sizeDeltaInQuantums: '0',
          client,
        }),
        updatePriceLevel({
          ticker,
          side: OrderSide.SELL,
          humanPrice: '65000',
          sizeDeltaInQuantums: '0',
          client,
        }),
        // Crossing price level, that will be removed when uncrossing
        updatePriceLevel({
          ticker,
          side: OrderSide.BUY,
          humanPrice: '52300',
          sizeDeltaInQuantums: '7000',
          client,
        }),
      ]);
    });

    afterEach(() => {
      jest.restoreAllMocks();
    });

    it('gets price levels, no options, removes zeros', async () => {
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client);

      expect(orderbookLevels.bids).toMatchObject([
        { humanPrice: '45100', quantums: '2000' },
        { humanPrice: '45200', quantums: '5000' },
        { humanPrice: '49300', quantums: '6000' },
        { humanPrice: '52300', quantums: '7000' },
      ]);
      expect(orderbookLevels.asks).toMatchObject([
        { humanPrice: '52200', quantums: '8000' },
        { humanPrice: '54200', quantums: '4000' },
        { humanPrice: '59300', quantums: '1000' },
      ]);
    });

    it('gets price levels, does not remove zeros', async () => {
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(
        ticker,
        client,
        {
          removeZeros: false,
        },
      );

      expect(orderbookLevels.bids).toMatchObject([
        { humanPrice: '40000', quantums: '0' },
        { humanPrice: '45100', quantums: '2000' },
        { humanPrice: '45200', quantums: '5000' },
        { humanPrice: '49300', quantums: '6000' },
        { humanPrice: '52300', quantums: '7000' },
      ]);
      expect(orderbookLevels.asks).toMatchObject([
        { humanPrice: '52200', quantums: '8000' },
        { humanPrice: '54200', quantums: '4000' },
        { humanPrice: '59300', quantums: '1000' },
        { humanPrice: '65000', quantums: '0' },
      ]);
    });

    it('gets price levels, sorts levels', async () => {
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(
        ticker,
        client,
        {
          sortSides: true,
        },
      );

      expect(orderbookLevels.bids).toMatchObject([
        { humanPrice: '52300', quantums: '7000' },
        { humanPrice: '49300', quantums: '6000' },
        { humanPrice: '45200', quantums: '5000' },
        { humanPrice: '45100', quantums: '2000' },
      ]);
      expect(orderbookLevels.asks).toMatchObject([
        { humanPrice: '52200', quantums: '8000' },
        { humanPrice: '54200', quantums: '4000' },
        { humanPrice: '59300', quantums: '1000' },
      ]);
    });

    it('gets price levels, sorts levels, and limits', async () => {
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(
        ticker,
        client,
        {
          sortSides: true,
          limitPerSide: 1,
        },
      );

      expect(orderbookLevels.bids).toMatchObject([
        { humanPrice: '52300', quantums: '7000' },
      ]);
      expect(orderbookLevels.asks).toMatchObject([
        { humanPrice: '52200', quantums: '8000' },
      ]);
    });

    it('gets price levels, sorts levels and uncrosses', async () => {
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(
        ticker,
        client,
        {
          sortSides: true,
          uncrossBook: true,
        },
      );

      expect(orderbookLevels.bids).toMatchObject([
        { humanPrice: '49300', quantums: '6000' },
        { humanPrice: '45200', quantums: '5000' },
        { humanPrice: '45100', quantums: '2000' },
      ]);
      expect(orderbookLevels.asks).toMatchObject([
        { humanPrice: '52200', quantums: '8000' },
        { humanPrice: '54200', quantums: '4000' },
        { humanPrice: '59300', quantums: '1000' },
      ]);
    });

    it('gets price levels, sorts levels and uncrosses (with respect to lastUpdated)', async () => {
      const expectedBids: PriceLevel[] = [
        { humanPrice: '49300', quantums: '6000', lastUpdated: '10' },
        { humanPrice: '45200', quantums: '5000', lastUpdated: '10' },
        { humanPrice: '45100', quantums: '2000', lastUpdated: '10' },
      ];
      const expectedAsks: PriceLevel[] = [
        { humanPrice: '52200', quantums: '8000', lastUpdated: '10' },
        { humanPrice: '54200', quantums: '4000', lastUpdated: '10' },
        { humanPrice: '59300', quantums: '1000', lastUpdated: '10' },
      ];

      const discardedBids: PriceLevel[] = [
        { humanPrice: '60000', quantums: '10000', lastUpdated: '5' }, // out-of-date
        { humanPrice: '59000', quantums: '1000', lastUpdated: '10' }, // quantums < SELL quantums
        { humanPrice: '58000', quantums: '8000', lastUpdated: '10' }, // matches SELL, default to SELL
      ];
      const discardedAsks: PriceLevel[] = [
        { humanPrice: '39000', quantums: '1000', lastUpdated: '10' }, // quantums < BUY quantums
        { humanPrice: '40000', quantums: '10000', lastUpdated: '5' }, // out-of-date
      ];

      function convertToRedisResults(priceLevels: PriceLevel[][]): string[][] {
        const redisQuantumsResult: string[] = [];
        const redisLastUpdatedResult: string[] = [];
        for (const lst of priceLevels) {
          redisQuantumsResult.unshift(..._.flatten(_.map(lst, (priceLevel: PriceLevel) => {
            return [priceLevel.humanPrice, priceLevel.quantums];
          })));
          redisLastUpdatedResult.unshift(..._.flatten(_.map(lst, (priceLevel: PriceLevel) => {
            return [priceLevel.humanPrice, priceLevel.lastUpdated];
          })));
        }
        return [redisQuantumsResult, redisLastUpdatedResult];
      }

      const mockRedisResultBids: string[][] = convertToRedisResults([expectedBids, discardedBids]);
      const mockRedisResultAsks: string[][] = convertToRedisResults([expectedAsks, discardedAsks]);

      const mockEvalBids: jest.Func = jest.fn((a, b, c, d, cb) => {
        (cb as Callback<string[][]>)(null, mockRedisResultBids);
      });
      const mockEvalAsks: jest.Func = jest.fn((a, b, c, d, cb) => {
        (cb as Callback<string[][]>)(null, mockRedisResultAsks);
      });

      jest.spyOn(client, 'evalsha')
        .mockImplementationOnce(mockEvalBids)
        .mockImplementationOnce(mockEvalAsks);

      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(
        ticker,
        client,
        {
          sortSides: true,
          uncrossBook: true,
        },
      );

      expect(orderbookLevels.bids).toEqual(expectedBids);
      expect(orderbookLevels.asks).toEqual(expectedAsks);
    });

    it('gets price levels, sorts levels, uncrosses, and limits', async () => {
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(
        ticker,
        client,
        {
          sortSides: true,
          uncrossBook: true,
          limitPerSide: 2,
        },
      );

      expect(orderbookLevels.bids).toMatchObject([
        { humanPrice: '49300', quantums: '6000' },
        { humanPrice: '45200', quantums: '5000' },
      ]);
      expect(orderbookLevels.asks).toMatchObject([
        { humanPrice: '52200', quantums: '8000' },
        { humanPrice: '54200', quantums: '4000' },
      ]);
    });

    it('gets price levels, sorts levels, uncrosses, and limits (no effect)', async () => {
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(
        ticker,
        client,
        {
          sortSides: true,
          uncrossBook: true,
          limitPerSide: 20,
        },
      );

      expect(orderbookLevels.bids).toMatchObject([
        { humanPrice: '49300', quantums: '6000' },
        { humanPrice: '45200', quantums: '5000' },
        { humanPrice: '45100', quantums: '2000' },
      ]);
      expect(orderbookLevels.asks).toMatchObject([
        { humanPrice: '52200', quantums: '8000' },
        { humanPrice: '54200', quantums: '4000' },
        { humanPrice: '59300', quantums: '1000' },
      ]);
    });

    it('throws error if sortSides is not true but uncrossBook is true', async () => {
      await expect(
        getOrderBookLevels(
          ticker,
          client,
          {
            uncrossBook: true,
          },
        ),
      ).rejects.toBeInstanceOf(InvalidOptionsError);
    });

    it.each([
      [-1],
      [0],
      [1],
    ])(
      'throws error if sortSides is not true but limitPerSide is %d',
      async (
        limitPerSide: number,
      ) => {
        await expect(
          getOrderBookLevels(
            ticker,
            client,
            { limitPerSide },
          ),
        ).rejects.toBeInstanceOf(InvalidOptionsError);
      },
    );
  });

  describe('deleteZeroPriceLevel', () => {
    const humanPrice: string = '45100';

    it('deletes zero price level', async () => {
      await updatePriceLevel({
        ticker,
        side: OrderSide.BUY,
        humanPrice: '45100',
        sizeDeltaInQuantums: '0',
        client,
      });

      let size: string | null = await hGetAsync(
        {
          hash: getKey(ticker, OrderSide.BUY),
          key: humanPrice,
        },
        client,
      );

      expect(size).toEqual('0');

      const deleted: boolean = await deleteZeroPriceLevel({
        ticker,
        side: OrderSide.BUY,
        humanPrice,
        client,
      });

      size = await hGetAsync(
        {
          hash: getKey(ticker, OrderSide.BUY),
          key: humanPrice,
        },
        client,
      );

      expect(deleted).toEqual(true);
      expect(size).toBeNull();
    });

    it('does not delete non-zero price level', async () => {
      await updatePriceLevel({
        ticker,
        side: OrderSide.BUY,
        humanPrice: '45100',
        sizeDeltaInQuantums: '10',
        client,
      });

      let size: string | null = await hGetAsync(
        {
          hash: getKey(ticker, OrderSide.BUY),
          key: humanPrice,
        },
        client,
      );

      expect(size).toEqual('10');

      const deleted: boolean = await deleteZeroPriceLevel({
        ticker,
        side: OrderSide.BUY,
        humanPrice,
        client,
      });

      size = await hGetAsync(
        {
          hash: getKey(ticker, OrderSide.BUY),
          key: humanPrice,
        },
        client,
      );

      expect(deleted).toEqual(false);
      expect(size).toEqual('10');
    });
  });
});
