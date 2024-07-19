import { deleteAllAsync, setAsync } from '../../src/helpers/redis';
import {
  IndexerOrder,
  RedisOrder,
} from '@dydxprotocol-indexer/v4-protos';
import {
  redis as client,
} from '../helpers/utils';
import { PlaceOrderResult } from '../../src/types';
import { InvalidRedisOrderError } from '../../src/errors';
import { getOrderDataCacheKey } from '../../src/caches/orders-data-cache';
import { placeOrder } from '../../src/caches/place-order';
import {
  order, orderGoodTilBlockTIme, redisOrder, redisOrderGoodTilBlockTime, secondRedisOrder,
} from './constants';
import { expectNewOrder, expectOrderCache, expectOrderExpiry } from './helpers';
import { ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import { getOrderExpiry } from '../../src/caches/helpers';
import Long from 'long';

describe('placeOrder', () => {

  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('placeOrder', () => {
    it('creates new order', async () => {
      const result: PlaceOrderResult = await placeOrder(
        {
          redisOrder,
          client,
        },
      );

      expectNewOrder(result);
      await expectOrderCache(result, redisOrder, 1);
      await expectOrderExpiry(redisOrder);
    });

    it('creates multiple new orders', async () => {
      let result: PlaceOrderResult = await placeOrder(
        {
          redisOrder,
          client,
        },
      );

      expectNewOrder(result);
      await expectOrderCache(result, redisOrder, 1);
      await expectOrderExpiry(redisOrder);

      result = await placeOrder(
        {
          redisOrder: secondRedisOrder,
          client,
        },
      );

      expectNewOrder(result);
      await expectOrderCache(result, secondRedisOrder, 2);
      await expectOrderExpiry(secondRedisOrder);
    });

    it.each([
      [
        'goodTilBlock',
        redisOrder,
        {
          ...order,
          // Change subticks/price to check that data of old order isn't changed
          subticks: Long.fromValue(2_900_000, true),
          goodTilBlock: 1149,
          goodTilBlockTime: undefined,
        } as IndexerOrder,
        '2900.0',
      ],
      [
        'goodTilBlockTime',
        redisOrderGoodTilBlockTime,
        {
          ...orderGoodTilBlockTIme,
          // Change subticks/price to check that data of old order isn't changed
          subticks: Long.fromValue(2_700_000, true),
          goodTilBlock: undefined,
          goodTilBlockTime: 16090,
        } as IndexerOrder,
        '2700.0',
      ],
    ])('does not replace existing order with greater or equal expiry (%s)', async (
      _name: string,
      placedOrder: RedisOrder,
      olderOrder: IndexerOrder,
      olderOrderPrice: string,
    ) => {
      await placeOrder(
        {
          redisOrder: placedOrder,
          client,
        },
      );

      const olderRedisOrder: RedisOrder = {
        ...placedOrder,
        order: olderOrder,
        price: olderOrderPrice,
      };

      const result: PlaceOrderResult = await placeOrder(
        {
          redisOrder: olderRedisOrder,
          client,
        },
      );

      expect(result.placed).toEqual(false);
      expect(result.replaced).toEqual(false);
      expect(result.oldTotalFilledQuantums).toBeUndefined();
      expect(result.restingOnBook).toBeUndefined();
      expect(result.oldOrder).toBeUndefined();

      await expectOrderCache(result, placedOrder, 1);
      await expectOrderExpiry(placedOrder);
    });

    it.each([
      [
        'goodTilBlock',
        redisOrder,
        {
          ...order,
          // Change subticks/price to check that data of old order isn't changed
          subticks: Long.fromValue(3_100_000, true),
          goodTilBlock: 1151,
          goodTilBlockTime: undefined,
        } as IndexerOrder,
        '3100.0',
      ],
      [
        'goodTilBlockTIme',
        redisOrderGoodTilBlockTime,
        {
          ...orderGoodTilBlockTIme,
          // Change subticks/price to check that data of old order isn't changed
          subticks: Long.fromValue(3_200_000, true),
          goodTilBlock: undefined,
          goodTilBlockTime: 17010,
        } as IndexerOrder,
        '3200.0',
      ],
    ])('replaces existing order with lesser expiry (%s), and returns old order', async (
      _name: string,
      placedOrder: RedisOrder,
      newerOrder: IndexerOrder,
      newerOrderPrice: string,
    ) => {
      await placeOrder(
        {
          redisOrder: placedOrder,
          client,
        },
      );

      await setAsync({
        key: getOrderDataCacheKey(placedOrder.order!.orderId!),
        value: `${getOrderExpiry(placedOrder.order!)}_43_true`,
      }, client);

      const newerRedisOrder: RedisOrder = {
        ...placedOrder,
        order: newerOrder,
        price: newerOrderPrice,
      };

      const result: PlaceOrderResult = await placeOrder(
        {
          redisOrder: newerRedisOrder,
          client,
        },
      );

      expect(result.placed).toEqual(false);
      expect(result.replaced).toEqual(true);
      expect(result.oldTotalFilledQuantums).toEqual(43);
      expect(result.restingOnBook).toEqual(true);
      expect(result.oldOrder).toEqual(placedOrder);

      await expectOrderCache(result, newerRedisOrder, 1);
      await expectOrderExpiry(newerRedisOrder);
    });

    it.each([
      ['missing order', { ...redisOrder, order: undefined }],
      ['missing order id', { ...redisOrder, order: { ...order, orderId: undefined } }],
      ['missing subaccount id', {
        ...redisOrder,
        order: {
          ...order,
          orderId: {
            clientId: 1,
            clobPairId: 0,
            orderFlags: ORDER_FLAG_SHORT_TERM,
          },
        },
      }],
    ])('throws error if given order is invalid: %s', async (
      _name: string,
      orderToPlace: RedisOrder,
    ) => {
      await expect(placeOrder({
        redisOrder: orderToPlace,
        client,
      })).rejects.toEqual(expect.any(InvalidRedisOrderError));
    });
  });
});
