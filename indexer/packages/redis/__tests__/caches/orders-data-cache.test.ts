import { IndexerOrderId } from '@dydxprotocol-indexer/v4-protos';
import { getOrderData, getOrderDataCacheKey } from '../../src/caches/orders-data-cache';
import { deleteAllAsync, setAsync } from '../../src/helpers/redis';
import { OrderData } from '../../src/types';
import {
  redis as client,
} from '../helpers/utils';
import { redisOrder } from './constants';

describe('ordersDataCache', () => {
  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  const defaultOrderId: IndexerOrderId = redisOrder.order!.orderId!;
  const dataCacheKey: string = getOrderDataCacheKey(defaultOrderId);

  describe('getOrderData', () => {
    it('gets an existing order', async () => {
      await setAsync(
        {
          key: dataCacheKey,
          value: '1150_43_true',
        },
        client,
      );
      const orderData: OrderData | null = await getOrderData(
        defaultOrderId,
        client,
      );
      expect(orderData).not.toBeNull();
      expect(orderData).toEqual({
        goodTilBlock: '1150',
        totalFilledQuantums: '43',
        restingOnBook: true,
      });
    });

    it('returns null for an non-existent order', async () => {
      const orderData: OrderData | null = await getOrderData(
        defaultOrderId,
        client,
      );
      expect(orderData).toBeNull();
    });
  });
});
