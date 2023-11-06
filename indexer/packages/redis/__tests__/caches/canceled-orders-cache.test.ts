import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import {
  CANCELED_ORDER_WINDOW_SIZE,
  isOrderCanceled,
  addCanceledOrderId,
  removeOrderFromCaches,
  getOrderCanceledStatus,
  addBestEffortCanceledOrderId,
} from '../../src/caches/canceled-orders-cache';
import { CanceledOrderStatus } from '../../src';

describe('cancelledOrdersCache', () => {
  const openOrderId1: string = 'order1';
  const openOrderId2: string = 'order2';
  const openOrderId3: string = 'order3';

  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  it('successfully cancels order', async () => {
    await addCanceledOrderId(openOrderId1, 10, client);
    await addCanceledOrderId(openOrderId2, 20, client);
    await addCanceledOrderId(openOrderId3, 10 + CANCELED_ORDER_WINDOW_SIZE, client);
    const isCanceled1: boolean = await isOrderCanceled(openOrderId1, client);
    const isCanceled2: boolean = await isOrderCanceled(openOrderId2, client);
    const isCanceled3: boolean = await isOrderCanceled(openOrderId3, client);
    expect(isCanceled1).toEqual(true);
    expect(isCanceled2).toEqual(true);
    expect(isCanceled3).toEqual(true);
  });

  it('successfully removes canceled order', async () => {
    await addCanceledOrderId(openOrderId1, 10, client);
    await addCanceledOrderId(openOrderId2, 20, client);
    let isCanceled1: boolean = await isOrderCanceled(openOrderId1, client);
    const isCanceled2: boolean = await isOrderCanceled(openOrderId2, client);
    expect(isCanceled1).toEqual(true);
    expect(isCanceled2).toEqual(true);

    await removeOrderFromCaches(openOrderId1, client);
    isCanceled1 = await isOrderCanceled(openOrderId1, client);
    expect(isCanceled1).toEqual(false);

    await removeOrderFromCaches(openOrderId3, client);
    const isCanceled3: boolean = await isOrderCanceled(openOrderId1, client);
    expect(isCanceled3).toEqual(false);
  });

  it('removes cancelled orders outside of window size', async () => {
    await addCanceledOrderId(openOrderId1, 10, client);
    await addCanceledOrderId(openOrderId2, 20, client);
    await addCanceledOrderId(openOrderId3, 10 + CANCELED_ORDER_WINDOW_SIZE + 1, client);
    const isCanceled1: boolean = await isOrderCanceled(openOrderId1, client);
    const isCanceled2: boolean = await isOrderCanceled(openOrderId2, client);
    const isCanceled3: boolean = await isOrderCanceled(openOrderId3, client);
    expect(isCanceled1).toEqual(false);
    expect(isCanceled2).toEqual(true);
    expect(isCanceled3).toEqual(true);
  });

  describe('getOrderCanceledStatus', () => {
    it('correctly returns CANCELED', async () => {
      await addCanceledOrderId(openOrderId1, 10, client);
      const status: CanceledOrderStatus = await getOrderCanceledStatus(openOrderId1, client);
      expect(status).toEqual(CanceledOrderStatus.CANCELED);
    });

    it('correctly returns BEST_EFFORT_CANCELED', async () => {
      await addBestEffortCanceledOrderId(openOrderId1, 10, client);
      const status: CanceledOrderStatus = await getOrderCanceledStatus(openOrderId1, client);
      expect(status).toEqual(CanceledOrderStatus.BEST_EFFORT_CANCELED);
    });

    it('correctly returns NOT_CANCELED', async () => {
      const status: CanceledOrderStatus = await getOrderCanceledStatus(openOrderId1, client);
      expect(status).toEqual(CanceledOrderStatus.NOT_CANCELED);
    });
  });
});
