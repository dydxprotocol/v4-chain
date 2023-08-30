import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import {
  CANCELED_ORDER_WINDOW_SIZE,
  isOrderCanceled,
  addCanceledOrderId,
  removeOrderFromCache,
} from '../../src/caches/canceled-orders-cache';

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

    let numRemoved: number = await removeOrderFromCache(openOrderId1, client);
    expect(numRemoved).toEqual(1);
    isCanceled1 = await isOrderCanceled(openOrderId1, client);
    expect(isCanceled1).toEqual(false);

    numRemoved = await removeOrderFromCache(openOrderId3, client);
    expect(numRemoved).toEqual(0);
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

});
