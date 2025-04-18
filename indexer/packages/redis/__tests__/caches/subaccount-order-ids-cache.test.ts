import {
  redisOrder,
  redisOrderSubaccount3,
  secondRedisOrder,
  subaccountUuid,
  subaccountUuid2,
  subaccountUuid3,
} from './constants';
import {
  redis as client,
} from '../helpers/utils';
import { getOrderIdsForSubaccount, getOrderIdsForSubaccounts } from '../../src/caches/subaccount-order-ids-cache';
import { placeOrder } from '../../src/caches/place-order';
import { deleteAllAsync } from '../../src/helpers/redis';

describe('subaccountOrderIdsCache', () => {
  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('getOrderIdsForSubaccount', () => {
    it('gets empty list for subaccount with no orders', async () => {
      const orderIds: string[] = await getOrderIdsForSubaccount(
        subaccountUuid,
        client,
      );

      expect(orderIds).toHaveLength(0);
    });

    it('gets order ids for subaccount', async () => {
      await Promise.all([
        placeOrder({
          redisOrder,
          client,
        }),
        placeOrder({
          redisOrder: secondRedisOrder,
          client,
        }),
      ]);

      const orderIds: string[] = await getOrderIdsForSubaccount(
        subaccountUuid,
        client,
      );

      expect(orderIds).toHaveLength(2);
      expect(orderIds).toContain(redisOrder.id);
      expect(orderIds).toContain(secondRedisOrder.id);
    });
  });

  describe('getOrderIdsForSubaccounts', () => {
    it('gets empty lists for a list of subaccounts with no orders', async () => {
      const subaccountOrderIds: Record<string, string[]> = await getOrderIdsForSubaccounts(
        [subaccountUuid, subaccountUuid2],
        client,
      );

      expect(subaccountOrderIds).toEqual({
        [subaccountUuid]: [],
        [subaccountUuid2]: [],
      });
    });

    it('gets order ids for a list of subaccounts', async () => {
      await Promise.all([
        // Place two orders for subaccount 0.
        placeOrder({
          redisOrder,
          client,
        }),
        placeOrder({
          redisOrder: secondRedisOrder,
          client,
        }),
        // Place one order for subaccount 3.
        placeOrder({
          redisOrder: redisOrderSubaccount3,
          client,
        }),
      ]);

      const subaccountOrderIds: Record<string, string[]> = await getOrderIdsForSubaccounts(
        [subaccountUuid, subaccountUuid2, subaccountUuid3],
        client,
      );

      expect(subaccountOrderIds).toEqual({
        [subaccountUuid]: [redisOrder.id, secondRedisOrder.id],
        [subaccountUuid2]: [],
        [subaccountUuid3]: [redisOrderSubaccount3.id],
      });
    });

  });
});
