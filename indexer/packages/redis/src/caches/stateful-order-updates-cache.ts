import { OrderUpdateV1 } from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';
import { Callback, RedisClient } from 'redis';

import { zRangeByScoreAsync } from '../helpers/redis';
import { StatefulOrderUpdateInfo } from '../types';
import { addStatefulOrderUpdateScript, removeStatefulOrderUpdateScript } from './scripts';

// Cache of order ids of the stateful order updates and when the updates were added to teh cache
export const ORDER_UPDATE_IDS_CACHE_KEY: string = 'v4/stateful_order_update_ids';
// Cache of order updates for stateful orders
export const ORDER_UPDATES_CACHE_KEY: string = 'v4/stateful_order_updates';

export async function addStatefulOrderUpdate(
  statefulOrderId: string,
  orderUpdate: OrderUpdateV1,
  updateTimestamp: number,
  client: RedisClient,
): Promise<void> {
  const numKeys: number = 2;
  let evalAsync: (
    orderId: string,
    encodedOrderUpdate: string,
    timestamp: number,
  ) => Promise<void> = (
    orderId,
    encodedOrderUpdate,
    timestamp,
  ) => {
    return new Promise((resolve, reject) => {
      const callback: Callback<void> = (
        err: Error | null,
      ) => {
        if (err) {
          return reject(err);
        }
        return resolve();
      };
      client.evalsha(
        addStatefulOrderUpdateScript.hash,
        numKeys,
        ORDER_UPDATE_IDS_CACHE_KEY,
        ORDER_UPDATES_CACHE_KEY,
        orderId,
        encodedOrderUpdate,
        timestamp,
        callback,
      );
    });
  };

  evalAsync = evalAsync.bind(client);

  return evalAsync(
    statefulOrderId,
    // TODO: use String to directly convert the UInt8Array to a string
    Buffer.from(OrderUpdateV1.encode(orderUpdate).finish()).toString('binary'),
    updateTimestamp,
  );
}

export async function removeStatefulOrderUpdate(
  statefulOrderId: string,
  removeTimestamp: number,
  client: RedisClient,
): Promise<OrderUpdateV1 | undefined> {
  const numKeys: number = 2;
  let evalAsync: (
    orderId: string,
    timestamp: number,
  ) => Promise<OrderUpdateV1 | undefined> = (
    orderId,
    timestamp,
  ) => {
    return new Promise((resolve, reject) => {
      const callback: Callback<string> = (
        err: Error | null,
        results: string,
      ) => {
        if (err) {
          return reject(err);
        }
        if (results === '') {
          return resolve(undefined);
        }

        return resolve(OrderUpdateV1.decode(Buffer.from(results, 'binary')));
      };
      client.evalsha(
        removeStatefulOrderUpdateScript.hash,
        numKeys,
        ORDER_UPDATE_IDS_CACHE_KEY,
        ORDER_UPDATES_CACHE_KEY,
        orderId,
        timestamp,
        callback,
      );
    });
  };

  evalAsync = evalAsync.bind(client);

  return evalAsync(
    statefulOrderId,
    removeTimestamp,
  );
}

export async function getOldOrderUpdates(
  latestTimestamp: number,
  client: RedisClient,
): Promise<StatefulOrderUpdateInfo[]> {
  const rawResults: string[] = await zRangeByScoreAsync({
    key: ORDER_UPDATE_IDS_CACHE_KEY,
    start: -Infinity,
    end: latestTimestamp,
    endIsInclusive: true,
    withScores: true,
  }, client);
  return _.chunk(rawResults, 2).map(
    (keyValuePair: string[]) => {
      return {
        orderId: keyValuePair[0],
        timestamp: Number(keyValuePair[1]),
      };
    },
  );
}
