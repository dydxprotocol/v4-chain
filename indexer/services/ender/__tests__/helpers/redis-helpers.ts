import { NodeEnv } from '@dydxprotocol-indexer/base';
import { OrderSide } from '@dydxprotocol-indexer/postgres';
import {
  NextFundingCache,
  OrderbookLevelsCache,
  StateFilledQuantumsCache,
} from '@dydxprotocol-indexer/redis';
import Big from 'big.js';

import { redisClient } from '../../src/helpers/redis/redis-controller';

export async function expectNextFundingRate(
  ticker: string,
  rate: Big | undefined,
): Promise<void> {
  const rates: { [ticker: string]: Big | undefined } = await NextFundingCache.getNextFunding(
    redisClient,
    [ticker],
  );
  expect(rates[ticker]).toEqual(rate);
}

export async function expectStateFilledQuantums(
  orderUuid: string,
  quantums: string,
): Promise<void> {
  const stateFilledQuantums: string | undefined = await StateFilledQuantumsCache
    .getStateFilledQuantums(
      orderUuid,
      redisClient,
    );
  expect(stateFilledQuantums).toBeDefined();
  expect(stateFilledQuantums).toEqual(quantums);
}

export async function updatePriceLevel(
  ticker: string,
  price: string,
  side: OrderSide,
): Promise<void> {
  const quantums: string = '30';

  await OrderbookLevelsCache.updatePriceLevel({
    ticker,
    side,
    humanPrice: price,
    sizeDeltaInQuantums: quantums,
    client: redisClient,
  });
}

export function clearOrderbookLevelsCacheForTests() {
  if (process.env.NODE_ENV !== NodeEnv.TEST) {
    throw Error('cannot clear orderbook levels outside of test environment');
  }

  redisClient.keys('v4/orderbookLevels/*', (err, keys) => {
    if (err) return;

    for (let i = 0; i < keys.length; i++) {
      redisClient.del(keys[i], (_innerErr, _reply) => {
      });
    }
  });
}
