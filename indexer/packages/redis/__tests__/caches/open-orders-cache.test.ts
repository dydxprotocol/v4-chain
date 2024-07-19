import { deleteAllAsync } from '../../src/helpers/redis';
import {
  redis as client,
} from '../helpers/utils';
import {
  addOpenOrder,
  removeOpenOrder,
  getOpenOrderIds,
} from '../../src/caches/open-orders-cache';

describe('ordersCache', () => {
  const fakeClobPairId: string = '1';
  const fakeClobPairId2: string = '2';
  const openOrderId1: string = 'order1';
  const openOrderId2: string = 'order2';
  const openOrderId3: string = 'order3';

  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  it('successfully adds open order', async () => {
    await addOpenOrder(openOrderId1, fakeClobPairId, client);
    await addOpenOrder(openOrderId2, fakeClobPairId, client);

    const openOrderIds: string[] = await getOpenOrderIds(fakeClobPairId, client);
    expect(openOrderIds).toHaveLength(2);
    expect(openOrderIds).toContain(openOrderId1);
    expect(openOrderIds).toContain(openOrderId2);
  });

  it('ignores adding open order that was previously added', async () => {
    await addOpenOrder(openOrderId1, fakeClobPairId, client);
    await addOpenOrder(openOrderId1, fakeClobPairId, client);

    const openOrderIds: string[] = await getOpenOrderIds(fakeClobPairId, client);
    expect(openOrderIds).toHaveLength(1);
    expect(openOrderIds).toContain(openOrderId1);
  });

  it('successfully removes open order', async () => {
    await addOpenOrder(openOrderId1, fakeClobPairId, client);
    await addOpenOrder(openOrderId2, fakeClobPairId, client);

    let openOrderIds: string[] = await getOpenOrderIds(fakeClobPairId, client);
    expect(openOrderIds).toHaveLength(2);
    expect(openOrderIds).toContain(openOrderId1);
    expect(openOrderIds).toContain(openOrderId2);

    await removeOpenOrder(openOrderId1, fakeClobPairId, client);
    openOrderIds = await getOpenOrderIds(fakeClobPairId, client);
    expect(openOrderIds).toHaveLength(1);
    expect(openOrderIds).toContain(openOrderId2);
  });

  it('ignores removing open order that doesn\'t exist', async () => {
    await addOpenOrder(openOrderId1, fakeClobPairId, client);
    await addOpenOrder(openOrderId2, fakeClobPairId, client);

    let openOrderIds: string[] = await getOpenOrderIds(fakeClobPairId, client);
    expect(openOrderIds).toHaveLength(2);
    expect(openOrderIds).toContain(openOrderId1);
    expect(openOrderIds).toContain(openOrderId2);

    await removeOpenOrder(openOrderId3, fakeClobPairId, client);
    openOrderIds = await getOpenOrderIds(fakeClobPairId, client);
    expect(openOrderIds).toHaveLength(2);
    expect(openOrderIds).toContain(openOrderId1);
    expect(openOrderIds).toContain(openOrderId2);
  });

  it('tracks separate clob pair ids in separate caches', async () => {
    await addOpenOrder(openOrderId1, fakeClobPairId, client);
    await addOpenOrder(openOrderId2, fakeClobPairId2, client);

    const openOrderIdsClob1: string[] = await getOpenOrderIds(fakeClobPairId, client);
    const openOrderIdsClob2: string[] = await getOpenOrderIds(fakeClobPairId2, client);
    expect(openOrderIdsClob1).toEqual([openOrderId1]);
    expect(openOrderIdsClob2).toEqual([openOrderId2]);
  });
});
