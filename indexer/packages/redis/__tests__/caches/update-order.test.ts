import { deleteAllAsync, setAsync } from '../../src/helpers/redis';
import {
  redis as client,
} from '../helpers/utils';
import { UpdateOrderResult } from '../../src/types';
import { InvalidRedisOrderError, InvalidTotalFilledQuantumsError } from '../../src/errors';
import { getOrderDataCacheKey } from '../../src/caches/orders-data-cache';
import { placeOrder } from '../../src/caches/place-order';
import {
  order,
  redisOrder,
} from './constants';
import {
  expectUpdateOrderResult,
  getOrderData,
} from './helpers';
import { updateOrder } from '../../src/caches/update-order';
import { ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import { getOrderExpiry } from '../../src/caches/helpers';

describe('updateOrder', () => {
  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('updateOrder', () => {
    it('updates existing order', async () => {
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

      const newTotalFilledQuantums: number = 50;
      const updateResult: UpdateOrderResult = await updateOrder(
        {
          updatedOrderId: redisOrder.order!.orderId!,
          newTotalFilledQuantums,
          client,
        },
      );

      // Check the result.
      expectUpdateOrderResult({
        result: updateResult,
        updated: true,
        oldTotalFilledQuantums: 43,
        oldRestingOnBook: true,
        order: redisOrder,
      });

      // Check that the order data has been updated
      const updatedOrderData: string | null = await getOrderData(
        redisOrder.order!.orderId!,
        client,
      );
      expect(updatedOrderData).toEqual(
        `${getOrderExpiry(redisOrder.order!)}_${newTotalFilledQuantums}_true`,
      );
    });

    it('updating new placed order sets resting on book to true', async () => {
      // Create a new order by placing it
      await placeOrder(
        {
          redisOrder,
          client,
        },
      );

      const newTotalFilledQuantums: number = 50;
      const updateResult: UpdateOrderResult = await updateOrder(
        {
          updatedOrderId: redisOrder.order!.orderId!,
          newTotalFilledQuantums,
          client,
        },
      );

      // Check the result.
      expectUpdateOrderResult({
        result: updateResult,
        updated: true,
        oldTotalFilledQuantums: 0,
        oldRestingOnBook: false,
        order: redisOrder,
      });

      // Check that the order data has been updated
      const updatedOrderData: string | null = await getOrderData(
        redisOrder.order!.orderId!,
        client,
      );
      expect(updatedOrderData).toEqual(
        `${getOrderExpiry(redisOrder.order!)}_${newTotalFilledQuantums}_true`,
      );
    });

    it('does nothing if an order is not updated', async () => {
      const result: UpdateOrderResult = await updateOrder({
        updatedOrderId: redisOrder.order!.orderId!,
        newTotalFilledQuantums: 1,
        client,
      });

      // Check the result
      expectUpdateOrderResult({
        result,
        updated: false,
      });

      // No order data should exist if an order was not updated
      const orderData: string | null = await getOrderData(redisOrder.order!.orderId!, client);
      expect(orderData).toEqual(null);
    });

    it('throws error if given order id is missing subaccount id', async () => {
      await expect(updateOrder({
        updatedOrderId: { clientId: 1, clobPairId: 0, orderFlags: ORDER_FLAG_SHORT_TERM },
        newTotalFilledQuantums: 1,
        client,
      })).rejects.toEqual(expect.any(InvalidRedisOrderError));
    });

    it('throws error if given new total filled quantums is negative', async () => {
      await expect(updateOrder({
        updatedOrderId: redisOrder.order!.orderId!,
        newTotalFilledQuantums: -1,
        client,
      })).rejects.toEqual(expect.any(InvalidTotalFilledQuantumsError));
    });
  });
});
