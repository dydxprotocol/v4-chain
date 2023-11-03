import { deleteAllAsync, ttl } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import { orderId } from './constants';
import { OrderTable } from '@dydxprotocol-indexer/postgres';
import {
  STATE_FILLED_QUANTUMS_TTL_SECONDS,
  getCacheKey,
  getStateFilledQuantums,
  updateStateFilledQuantums,
} from '../../src/caches/state-filled-quantums-cache';

describe('stateFilledQuantumsCache', () => {
  const orderUuid: string = OrderTable.orderIdToUuid(orderId);

  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('updateStateFilledQuantums', () => {
    it('updates the state filled amount for an order id', async () => {
      const filledQuantums: string = '1000';
      await updateStateFilledQuantums(orderUuid, filledQuantums, client);

      expect(await getStateFilledQuantums(orderUuid, client)).toEqual(filledQuantums);
      expect(await ttl(client, getCacheKey(orderUuid))).toEqual(STATE_FILLED_QUANTUMS_TTL_SECONDS);
    });
  });

  describe('getStateFilledQuantums', () => {
    it('gets the state filled amount for an order id', async () => {
      const filledQuantums: string = '1000';
      await updateStateFilledQuantums(orderUuid, filledQuantums, client);

      expect(await getStateFilledQuantums(orderUuid, client)).toEqual(filledQuantums);
    });

    it('returns undefined if order id does not exist', async () => {
      expect(await getStateFilledQuantums(orderUuid, client)).toEqual(undefined);
    });
  });
});
