import { OrderSide } from '@dydxprotocol-indexer/postgres';
import {
  NextFundingCache,
  OrderbookLevelsCache,
  StateFilledQuantumsCache,
} from '@dydxprotocol-indexer/redis';
import Big from 'big.js';

import { redisClient } from '../../src/helpers/redis/redis-controller';

export async function expectNextFundingRate(
  expectedRate: Big | undefined,
  ticker: string,
  defaultFundingRate1H: string = '0',
): Promise<void> {
  const rates: { [ticker: string]: Big | undefined } = await NextFundingCache.getNextFunding(
    redisClient,
    [[ticker, defaultFundingRate1H]],
  );
  expect(rates[ticker]).toEqual(expectedRate);
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

  await OrderbookLevelsCache.updatePriceLevel(
    ticker,
    side,
    price,
    quantums,
    redisClient,
  );
}
