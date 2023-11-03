import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import {
  addStatefulOrderUpdate,
  removeStatefulOrderUpdate,
  getOldOrderUpdates,
} from '../../src/caches/stateful-order-updates-cache';
import { IndexerOrderId, OrderUpdateV1 } from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';
import { orderId } from './constants';
import { OrderTable } from '@dydxprotocol-indexer/postgres';
import { StatefulOrderUpdateInfo } from 'packages/redis/src';

describe('statefulOrderUpdatesCache', () => {
  const orderUpdate: OrderUpdateV1 = {
    orderId,
    totalFilledQuantums: Long.fromNumber(100, true),
  };
  const orderUuid: string = OrderTable.orderIdToUuid(orderId);
  const initialTimestamp: number = Date.now();
  const olderTimestamp: number = initialTimestamp - 10;
  const newerTimestamp: number = initialTimestamp + 10;

  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('addStatefulOrderUpdate', () => {
    it('adds stateful order update to cache', async () => {
      await addStatefulOrderUpdate(
        orderUuid,
        orderUpdate,
        initialTimestamp,
        client,
      );

      const statefulOrderUpdateInfo: StatefulOrderUpdateInfo[] = await getOldOrderUpdates(
        newerTimestamp, client,
      );
      const removedUpdate: OrderUpdateV1 | undefined = await removeStatefulOrderUpdate(
        orderUuid, initialTimestamp, client,
      );

      expect(statefulOrderUpdateInfo).toHaveLength(1);
      expect(statefulOrderUpdateInfo[0]).toEqual({
        orderId: orderUuid,
        timestamp: initialTimestamp,
      });
      expect(removedUpdate).toBeDefined();
      expect(removedUpdate).toEqual(orderUpdate);
    });
  });

  describe('removeStatefulorderUpdate', () => {
    it('removes and returns existing stateful order update from cache', async () => {
      await addStatefulOrderUpdate(
        orderUuid,
        orderUpdate,
        initialTimestamp,
        client,
      );

      const removedUpdate: OrderUpdateV1 | undefined = await removeStatefulOrderUpdate(
        orderUuid, initialTimestamp, client,
      );
      const statefulOrderUpdateInfo: StatefulOrderUpdateInfo[] = await getOldOrderUpdates(
        newerTimestamp, client,
      );

      expect(removedUpdate).toBeDefined();
      expect(removedUpdate).toEqual(orderUpdate);
      expect(statefulOrderUpdateInfo).toHaveLength(0);
    });

    it('does not remove existing stateful order update if timestamp is lower', async () => {
      await addStatefulOrderUpdate(
        orderUuid,
        orderUpdate,
        initialTimestamp,
        client,
      );

      const removedUpdate: OrderUpdateV1 | undefined = await removeStatefulOrderUpdate(
        orderUuid, olderTimestamp, client,
      );
      const statefulOrderUpdateInfo: StatefulOrderUpdateInfo[] = await getOldOrderUpdates(
        newerTimestamp, client,
      );

      expect(removedUpdate).toBeUndefined();
      expect(statefulOrderUpdateInfo).toHaveLength(1);
      expect(statefulOrderUpdateInfo[0]).toEqual({
        orderId: orderUuid,
        timestamp: initialTimestamp,
      });
    });

    it('removes non-existing order and returns undefined', async () => {
      const removedUpdate: OrderUpdateV1 | undefined = await removeStatefulOrderUpdate(
        orderUuid, initialTimestamp, client,
      );

      expect(removedUpdate).toBeUndefined();
    });
  });

  describe('getOldOrderUpdates', () => {
    const orderId2: IndexerOrderId = {
      ...orderId,
      clientId: 45,
    };
    const orderUuid2: string = OrderTable.orderIdToUuid(orderId2);
    const orderUpdate2: OrderUpdateV1 = {
      ...orderUpdate,
      orderId: orderId2,
    };

    beforeEach(async () => {
      await Promise.all([
        addStatefulOrderUpdate(
          orderUuid,
          orderUpdate,
          initialTimestamp,
          client,
        ),
        addStatefulOrderUpdate(
          orderUuid2,
          orderUpdate2,
          olderTimestamp,
          client,
        )],
      );
    });

    it('returns stateful order update info older than the threshold', async () => {
      const statefulOrderUpdateInfo: StatefulOrderUpdateInfo[] = await getOldOrderUpdates(
        olderTimestamp, client,
      );

      expect(statefulOrderUpdateInfo).toEqual([{
        orderId: orderUuid2,
        timestamp: olderTimestamp,
      }]);
    });

    it('returns multiple stateful order update info older than the threshold', async () => {
      const statefulOrderUpdateInfo: StatefulOrderUpdateInfo[] = await getOldOrderUpdates(
        initialTimestamp, client,
      );

      expect(statefulOrderUpdateInfo).toEqual([{
        orderId: orderUuid2,
        timestamp: olderTimestamp,
      }, {
        orderId: orderUuid,
        timestamp: initialTimestamp,
      }]);
    });
  });
});
