import {
  redis as client,
} from '../helpers/utils';
import { SubaccountTable } from '@dydxprotocol-indexer/postgres';
import { getOrderIdsForSubaccount } from '../../src/caches/subaccount-order-ids-cache';
import { placeOrder } from '../../src/caches/place-order';
import {
  address, redisOrder, secondRedisOrder, subaccountNumber,
} from './constants';

describe('subaccountOrderIdsCache', () => {
  describe('getOrderIdsForSubaccount', () => {
    it('gets empty list for subaccount with no orders', async () => {
      const orderIds: string[] = await getOrderIdsForSubaccount(
        SubaccountTable.uuid(address, subaccountNumber),
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
        SubaccountTable.uuid(address, subaccountNumber),
        client,
      );

      expect(orderIds).toHaveLength(2);
      expect(orderIds).toContain(redisOrder.id);
      expect(orderIds).toContain(secondRedisOrder.id);
    });
  });
});
