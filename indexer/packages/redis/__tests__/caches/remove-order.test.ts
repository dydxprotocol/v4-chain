import { deleteAllAsync, setAsync } from '../../src/helpers/redis';
import {
  redis as client,
} from '../helpers/utils';
import { RemoveOrderResult } from '../../src/types';
import { InvalidRedisOrderError } from '../../src/errors';
import { getOrderDataCacheKey } from '../../src/caches/orders-data-cache';
import { placeOrder } from '../../src/caches/place-order';
import {
  order,
  redisOrder,
  secondRedisOrder,
} from './constants';
import {
  expectRemovedOrder,
  expectRemovedOrderResult,
  getOrderByOrderId,
  getOrderData,
  getOrderExpiryFromCache,
  getOrderIdsForSubaccountBySubaccountId,
} from './helpers';
import { removeOrder } from '../../src/caches/remove-order';
import { OrderTable } from '@dydxprotocol-indexer/postgres';
import { ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import { getOrderExpiry } from '../../src/caches/helpers';

describe('removeOrder', () => {
  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('removeOrder', () => {
    it('removes existing order', async () => {
      // Create a new order by placing it
      await placeOrder(
        {
          redisOrder,
          client,
        },
      );

      // Set some additional values for the order data
      await setAsync({
        key: getOrderDataCacheKey(order.orderId!),
        value: '1150_43_true',
      }, client);

      const removeResult: RemoveOrderResult = await removeOrder(
        {
          removedOrderId: redisOrder.order!.orderId!,
          client,
        },
      );

      // Check the result.
      expectRemovedOrderResult({
        result: removeResult,
        removed: true,
        totalFilledQuantums: 43,
        restingOnBook: true,
        removedOrder: redisOrder,
      });

      // Check that the order has been removed from the caches.
      await expectRemovedOrder(redisOrder.order!.orderId!);
    });

    it('removes existing order, leaving another order around', async () => {
      // Create 2 orders by placing them
      await Promise.all([
        await placeOrder(
          {
            redisOrder,
            client,
          },
        ),
        await placeOrder(
          {
            redisOrder: secondRedisOrder,
            client,
          },
        ),
      ]);

      const removeResult: RemoveOrderResult = await removeOrder(
        {
          removedOrderId: secondRedisOrder.order!.orderId!,
          client,
        },
      );

      // Check the result.
      expectRemovedOrderResult({
        result: removeResult,
        removed: true,
        totalFilledQuantums: 0,
        restingOnBook: false,
        removedOrder: secondRedisOrder,
      });

      // Check that the order has been removed from the caches.
      await expectRemovedOrder(secondRedisOrder.order!.orderId!);

      // The first order should still exist in caches.
      expect(await getOrderByOrderId(redisOrder.order!.orderId!, client)).toEqual(redisOrder);
      expect(await getOrderData(redisOrder.order!.orderId!, client)).toEqual(
        `${getOrderExpiry(redisOrder.order!)}_0_false`,
      );
      expect(
        await getOrderIdsForSubaccountBySubaccountId(
          redisOrder.order!.orderId!.subaccountId!,
          client,
        ),
      ).toContain(OrderTable.orderIdToUuid(redisOrder.order!.orderId!));
      expect(
        await getOrderExpiryFromCache(redisOrder.id, client),
      ).toEqual(getOrderExpiry(redisOrder.order!).toString());
    });

    it('returns nothing if order to remove does not exist', async () => {
      const result: RemoveOrderResult = await removeOrder({
        removedOrderId: redisOrder.order!.orderId!,
        client,
      });

      expectRemovedOrderResult({
        result,
        removed: false,
      });
    });

    it('throws error if given order id is missing subaccount id', async () => {
      await expect(removeOrder({
        removedOrderId: { clientId: 1, clobPairId: 0, orderFlags: ORDER_FLAG_SHORT_TERM },
        client,
      })).rejects.toEqual(expect.any(InvalidRedisOrderError));
    });
  });
});
