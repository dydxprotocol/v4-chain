import { IsoString, OrderTable, SubaccountTable } from '@dydxprotocol-indexer/postgres';
import { ORDER_FLAG_LONG_TERM, ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerOrder_Side,
  IndexerOrder_TimeInForce,
  RedisOrder,
  RedisOrder_TickerType,
  IndexerOrderId,
  IndexerOrder_ConditionType,
} from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';
import { DateTime } from 'luxon';

export const address: string = 'dydxprotocol174e000tqwvszgjxs7yaj844e0m9s6f0m45ws7q';
export const subaccountNumber: number = 0;
export const subaccountNumber2: number = 2;
export const subaccountNumber3: number = 3;
export const clientId: number = 12;
export const subaccountUuid: string = SubaccountTable.uuid(address, subaccountNumber);
export const subaccountUuid2: string = SubaccountTable.uuid(address, subaccountNumber2);
export const subaccountUuid3: string = SubaccountTable.uuid(address, subaccountNumber3);
export const orderId: IndexerOrderId = {
  subaccountId: {
    owner: address,
    number: subaccountNumber,
  },
  clientId,
  clobPairId: 0,
  orderFlags: ORDER_FLAG_SHORT_TERM,
};
export const createdAt: IsoString = DateTime.utc().toISO();
export const createdAtHeight: string = '1';
export const order: IndexerOrder = {
  orderId,
  subticks: Long.fromValue(3_000_000, true),
  quantums: Long.fromValue(5_000_000_000, true),
  side: IndexerOrder_Side.SIDE_BUY,
  goodTilBlock: 1150,
  timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
  reduceOnly: false,
  clientMetadata: 0,
  conditionType: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
  conditionalOrderTriggerSubticks: Long.fromValue(0, true),
  orderRouterAddress: '',
};
export const orderGoodTilBlockTIme: IndexerOrder = {
  ...order,
  orderId: {
    ...orderId,
    orderFlags: ORDER_FLAG_LONG_TERM,
  },
  goodTilBlockTime: 17000,
  goodTilBlock: undefined,
};
export const redisOrder: RedisOrder = {
  id: OrderTable.orderIdToUuid(orderId),
  order,
  price: '3000.0',
  size: '5.0',
  ticker: 'BTC-USD',
  tickerType: RedisOrder_TickerType.TICKER_TYPE_PERPETUAL,
};
export const redisOrderGoodTilBlockTime: RedisOrder = {
  ...redisOrder,
  id: OrderTable.orderIdToUuid(orderGoodTilBlockTIme.orderId!),
  order: orderGoodTilBlockTIme,
};
export const secondRedisOrder: RedisOrder = {
  ...redisOrder,
  id: OrderTable.uuid(
    subaccountUuid,
    (clientId + 1).toString(),
    '0',
    ORDER_FLAG_SHORT_TERM.toString(),
  ),
  order: {
    ...order,
    orderId: {
      ...order.orderId,
      clientId: order.orderId!.clientId + 1,
      clobPairId: 0,
      orderFlags: ORDER_FLAG_SHORT_TERM,
    },
  },
};
export const redisOrderSubaccount3: RedisOrder = {
  ...redisOrder,
  id: OrderTable.uuid(
    subaccountUuid3,
    (clientId + 2).toString(),
    '1',
    ORDER_FLAG_SHORT_TERM.toString(),
  ),
  order: {
    ...order,
    orderId: {
      ...order.orderId,
      subaccountId: {
        owner: address,
        number: subaccountNumber3,
      },
      clientId: clientId + 2,
      clobPairId: 1,
      orderFlags: ORDER_FLAG_SHORT_TERM,
    },
  },
};
