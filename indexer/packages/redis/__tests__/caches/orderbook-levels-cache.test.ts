import _ from 'lodash';
import { Callback, RedisClient } from 'redis';
import { deleteAllAsync, hGetAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import {
  deleteStalePriceLevel,
  deleteZeroPriceLevel,
  getKey,
  getLastUpdatedKey,
  getOrderBookLevels,
  getOrderBookMidPrice,
  updatePriceLevel,
} from '../../src/caches/orderbook-levels-cache';
import { OrderSide } from '@dydxprotocol-indexer/postgres';
import { OrderbookLevels, PriceLevel } from '../../src/types';
import { InvalidOptionsError } from '../../src/errors';
import { logger } from '@dydxprotocol-indexer/base';

interface MockMulti {
  hgetall(key: string): MockMulti,
  exec(cb: Callback<{ [field: string]: string }[]>): boolean,
}

function getMockMulti(results: { [field: string]: string }[]): MockMulti {
  return {
    hgetall: jest.fn().mockReturnThis(),
    exec: jest.fn((cb: Callback<{ [field: string]: string }[]>) => {
      cb(null, results);
      return true;
    }),
  };
}

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
    ])(
      'updates the total size in quantums for a price/side/ticker, new price level: %s',
      async (
        _name: string,
        tickerForTest: string,
        side: OrderSide,
        humanPrice: string,
        sizeDeltaInQuantums: string,
      ) => {
        const updatedQuantums: number = await updatePriceLevel(
          tickerForTest,
          side,
          humanPrice,
          sizeDeltaInQuantums,
          client,
        );
        const orderbook: OrderbookLevels = await getOrderBookLevels(tickerForTest, client);

        expect(updatedQuantums.toString()).toEqual(sizeDeltaInQuantums);
        let levels: PriceLevel[];
        if (side === OrderSide.BUY) {
          levels = orderbook.bids;
        } else {
          levels = orderbook.asks;
        }
        expect(levels).toMatchObject([
          {
            humanPrice,
            quantums: sizeDeltaInQuantums,
          },
        ]);
      },
    );

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
    ])(
      'updates the total size in quantums for a price/side/ticker, existing price level: %s',
      async (
        _name: string,
        tickerForTest: string,
        side: OrderSide,
        humanPrice: string,
        sizeDeltaInQuantums: string,
        existingQuantums: string,
        expectedQuantums: string,
      ) => {
        // Set existing quantums for the level
        await updatePriceLevel(
          tickerForTest,
          side,
          humanPrice,
          existingQuantums,
          client,
        );

        const updatedQuantums: number = await updatePriceLevel(
          tickerForTest,
          side,
          humanPrice,
          sizeDeltaInQuantums,
          client,
        );
        const orderbook: OrderbookLevels = await getOrderBookLevels(tickerForTest, client);

        expect(updatedQuantums.toString()).toEqual(expectedQuantums);
        let levels: PriceLevel[];
        if (side === OrderSide.BUY) {
          levels = orderbook.bids;
        } else {
          levels = orderbook.asks;
        }
        expect(levels).toMatchObject([
          {
            humanPrice,
            quantums: expectedQuantums,
          },
        ]);
      },
    );

    it('can update price level to zero quantums', async () => {
      const humanPrice: string = '50000';
      const quantums: string = '3000';
      const updateQuantums: string = '-3000';
      // set up initial quantums
      await updatePriceLevel(
        ticker,
        OrderSide.BUY,
        humanPrice,
        quantums,
        client,
      );

      let orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client);
      expect(orderbookLevels.bids).toMatchObject([
        {
          humanPrice,
          quantums,
        },
      ]);

      // Update with delta to set quantums to 0
      const result: number = await updatePriceLevel(
        ticker,
        OrderSide.BUY,
        humanPrice,
        updateQuantums,
        client,
      );
      orderbookLevels = await getOrderBookLevels(ticker, client);

      expect(result).toEqual(0);
      expect(orderbookLevels.bids).toEqual([]);
    });

    it('sets price level to 0 if update will cause quantums to be negative', async () => {
      const humanPrice: string = '50000';
      const quantums: string = '1000';
      const invalidDelta: string = '-2000';
      // Set existing quantums for the level
      await updatePriceLevel(
        ticker,
        OrderSide.BUY,
        humanPrice,
        quantums,
        client,
      );

      await updatePriceLevel(
        ticker,
        OrderSide.BUY,
        humanPrice,
        invalidDelta,
        client,
      );
      expect(logger.crit).toHaveBeenCalledTimes(1);

      // Expect that the value in the orderbook is set to 0
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client, {
        removeZeros: false,
      });
      expect(orderbookLevels.bids).toMatchObject([
        {
          humanPrice,
          quantums: '0',
        },
      ]);
    });
  });

  describe('getOrderbookLevels', () => {
    beforeEach(async () => {
      // Setup some price levels
      // Bids
      await Promise.all([
        updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '45100',
          '2000',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '45200',
          '5000',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '49300',
          '6000',
          client,
        ),
        // Asks
        updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '52200',
          '8000',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '54200',
          '4000',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '59300',
          '1000',
          client,
        ),
        // Zero quantums bid/ask
        updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '40000',
          '0',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '65000',
          '0',
          client,
        ),
        // Crossing price level, that will be removed when uncrossing
        updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '52300',
          '7000',
          client,
        ),
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
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client, {
        removeZeros: false,
      });

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
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client, {
        sortSides: true,
      });

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
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client, {
        sortSides: true,
        limitPerSide: 1,
      });

      expect(orderbookLevels.bids).toMatchObject([{ humanPrice: '52300', quantums: '7000' }]);
      expect(orderbookLevels.asks).toMatchObject([{ humanPrice: '52200', quantums: '8000' }]);
    });

    it('gets price levels, sorts levels and uncrosses', async () => {
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client, {
        sortSides: true,
        uncrossBook: true,
      });

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
      const discardedBids: PriceLevel[] = [
        { humanPrice: '60000', quantums: '10000', lastUpdated: '5' }, // out-of-date
        { humanPrice: '59000', quantums: '1000', lastUpdated: '10' }, // quantums < SELL quantums
        { humanPrice: '58000', quantums: '8000', lastUpdated: '10' }, // matches SELL, default SELL
      ];
      const bidsQuantums: { [field: string]: string } = _.fromPairs(
        _.map(_.flatten([expectedBids, discardedBids]), (priceLevel: PriceLevel) => [
          priceLevel.humanPrice,
          priceLevel.quantums,
        ]),
      );
      const bidsLastUpdated: { [field: string]: string } = _.fromPairs(
        _.map(_.flatten([expectedBids, discardedBids]), (priceLevel: PriceLevel) => [
          priceLevel.humanPrice,
          priceLevel.lastUpdated,
        ]),
      );

      const expectedAsks: PriceLevel[] = [
        { humanPrice: '52200', quantums: '8000', lastUpdated: '10' },
        { humanPrice: '54200', quantums: '4000', lastUpdated: '10' },
        { humanPrice: '59300', quantums: '1000', lastUpdated: '10' },
      ];
      const discardedAsks: PriceLevel[] = [
        { humanPrice: '39000', quantums: '1000', lastUpdated: '10' }, // quantums < BUY quantums
        { humanPrice: '40000', quantums: '10000', lastUpdated: '5' }, // out-of-date
      ];
      const asksQuantums: { [field: string]: string } = _.fromPairs(
        _.map(_.flatten([expectedAsks, discardedAsks]), (priceLevel: PriceLevel) => [
          priceLevel.humanPrice,
          priceLevel.quantums,
        ]),
      );
      const asksLastUpdated: { [field: string]: string } = _.fromPairs(
        _.map(_.flatten([expectedAsks, discardedAsks]), (priceLevel: PriceLevel) => [
          priceLevel.humanPrice,
          priceLevel.lastUpdated,
        ]),
      );

      jest
        .spyOn(RedisClient.prototype, 'multi')
        .mockReturnValueOnce(getMockMulti([bidsQuantums, bidsLastUpdated]))
        .mockReturnValueOnce(getMockMulti([asksQuantums, asksLastUpdated]));

      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client, {
        sortSides: true,
        uncrossBook: true,
      });

      expect(orderbookLevels.bids).toEqual(expectedBids);
      expect(orderbookLevels.asks).toEqual(expectedAsks);
    });

    it('gets price levels, sorts levels, uncrosses, and limits', async () => {
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client, {
        sortSides: true,
        uncrossBook: true,
        limitPerSide: 2,
      });

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
      const orderbookLevels: OrderbookLevels = await getOrderBookLevels(ticker, client, {
        sortSides: true,
        uncrossBook: true,
        limitPerSide: 20,
      });

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
        getOrderBookLevels(ticker, client, {
          uncrossBook: true,
        }),
      ).rejects.toBeInstanceOf(InvalidOptionsError);
    });

    it.each([[-1], [0], [1]])(
      'throws error if sortSides is not true but limitPerSide is %d',
      async (limitPerSide: number) => {
        await expect(getOrderBookLevels(ticker, client, { limitPerSide })).rejects.toBeInstanceOf(
          InvalidOptionsError,
        );
      },
    );
  });

  describe('deleteZeroPriceLevel', () => {
    const humanPrice: string = '45100';

    it('deletes zero price level', async () => {
      await updatePriceLevel(
        ticker,
        OrderSide.BUY,
        '45100',
        '0',
        client,
      );

      let size: string | null = await hGetAsync(
        {
          hash: getKey(ticker, OrderSide.BUY),
          key: humanPrice,
        },
        client,
      );

      expect(size).toEqual('0');

      const deleted: boolean = await deleteZeroPriceLevel(
        ticker,
        OrderSide.BUY,
        humanPrice,
        client,
      );

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
      await updatePriceLevel(
        ticker,
        OrderSide.BUY,
        '45100',
        '10',
        client,
      );

      let size: string | null = await hGetAsync(
        {
          hash: getKey(ticker, OrderSide.BUY),
          key: '45100',
        },
        client,
      );

      expect(size).toEqual('10');

      const deleted: boolean = await deleteZeroPriceLevel(
        ticker,
        OrderSide.BUY,
        '45100',
        client,
      );

      size = await hGetAsync(
        {
          hash: getKey(ticker, OrderSide.BUY),
          key: '45100',
        },
        client,
      );

      expect(deleted).toEqual(false);
      expect(size).toEqual('10');
    });
  });

  describe('deleteStalePriceLevel', () => {
    const humanPrice: string = '45100';

    it('deletes stale price level', async () => {
      await updatePriceLevel(
        ticker,
        OrderSide.BUY,
        humanPrice,
        '100',
        client,
      );

      client.hset(getLastUpdatedKey(ticker, OrderSide.BUY), humanPrice, Date.now() / 1000 - 20);

      let size: string | null = await hGetAsync(
        {
          hash: getKey(ticker, OrderSide.BUY),
          key: humanPrice,
        },
        client,
      );

      expect(size).toEqual('100');

      const deleted: boolean = await deleteStalePriceLevel(
        ticker,
        OrderSide.BUY,
        humanPrice,
        10,
        client,
      );

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

    it('does not delete recent price level', async () => {
      await updatePriceLevel(
        ticker,
        OrderSide.BUY,
        humanPrice,
        '10',
        client,
      );

      const deleted: boolean = await deleteStalePriceLevel(
        ticker,
        OrderSide.BUY,
        humanPrice,
        10,
        client,
      );

      const size: string | null = await hGetAsync(
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

  describe('getMidPrice', () => {
    beforeEach(() => {
      jest.restoreAllMocks();
      jest.restoreAllMocks();
    });
    afterEach(() => {
      jest.restoreAllMocks();
      jest.restoreAllMocks();
    });

    it('returns the correct mid price', async () => {
      await Promise.all([
        updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '45200',
          '2000',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '45100',
          '2000',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.BUY,
          '45300',
          '2000',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '45500',
          '2000',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '45400',
          '2000',
          client,
        ),
        updatePriceLevel(
          ticker,
          OrderSide.SELL,
          '45600',
          '2000',
          client,
        ),
      ]);

      const midPrice = await getOrderBookMidPrice(ticker, client);
      expect(midPrice).toEqual('45350');
    });
  });

  it('returns the correct mid price for very small numbers', async () => {
    await Promise.all([
      updatePriceLevel(
        ticker,
        OrderSide.SELL,
        '0.000000002346',
        '2000',
        client,
      ),
      updatePriceLevel(
        ticker,
        OrderSide.BUY,
        '0.000000002344',
        '2000',
        client,
      ),
    ]);

    const midPrice = await getOrderBookMidPrice(ticker, client);
    expect(midPrice).toEqual('0.000000002345');
  });

  it('returns the approprite amount of decimal precision', async () => {
    await Promise.all([
      updatePriceLevel(
        ticker,
        OrderSide.SELL,
        '1.02',
        '2000',
        client,
      ),
      updatePriceLevel(
        ticker,
        OrderSide.BUY,
        '1.01',
        '2000',
        client,
      ),
    ]);

    const midPrice = await getOrderBookMidPrice(ticker, client);
    expect(midPrice).toEqual('1.015');
  });

  it('returns undefined if there are no bids or asks', async () => {
    await updatePriceLevel(
      ticker,
      OrderSide.SELL,
      '45400',
      '2000',
      client,
    );

    const midPrice = await getOrderBookMidPrice(ticker, client);
    expect(midPrice).toBeUndefined();
  });

  it('returns undefined if humanPrice is NaN', async () => {
    await updatePriceLevel(
      ticker,
      OrderSide.SELL,
      'nan',
      '2000',
      client,
    );

    const midPrice = await getOrderBookMidPrice(ticker, client);

    expect(midPrice).toBeUndefined();
  });
});
