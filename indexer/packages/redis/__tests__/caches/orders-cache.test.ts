import { getOrder } from '../../src/caches/orders-cache';
import { deleteAllAsync } from '../../src/helpers/redis';
import {
  redis as client,
} from '../helpers/utils';
import { placeOrder } from '../../src/caches/place-order';
import { redisOrder } from './constants';

describe('ordersCache', () => {
  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('getOrder', () => {
    it('gets an existing order', async () => {
      await placeOrder({
        redisOrder,
        client,
      });

      expect(await getOrder(redisOrder.id, client)).toEqual(redisOrder);
    });

    it('returns null for an non-existent order', async () => {
      expect(await getOrder(redisOrder.id, client)).toEqual(null);
    });
  });
});
