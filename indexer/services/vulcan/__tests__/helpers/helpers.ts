import { createKafkaMessage } from '@dydxprotocol-indexer/kafka';
import { OrderSide } from '@dydxprotocol-indexer/postgres';
import {
  redisTestConstants,
  OrderbookLevelsCache,
  CanceledOrdersCache,
  CanceledOrderStatus,
} from '@dydxprotocol-indexer/redis';
import { OffChainUpdateV1 } from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';

import { redisClient } from '../../src/helpers/redis/redis-controller';
import { onMessage } from '../../src/lib/on-message';
import { DydxRecordHeaderKeys } from '../../src/lib/types';
import { defaultKafkaHeaders } from './constants';

export async function handleInitialOrderPlace(
  orderPlace: redisTestConstants.OffChainUpdateOrderPlaceUpdateMessage,
): Promise<void> {
  const update: OffChainUpdateV1 = {
    ...orderPlace,
  };
  const message: KafkaMessage = createKafkaMessage(
    Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(update).finish())),
  );
  message.headers = defaultKafkaHeaders;

  await onMessage(message);
}

export async function handleOrderUpdate(
  orderUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage,
): Promise<void> {
  const update: OffChainUpdateV1 = {
    ...orderUpdate,
  };
  const message: KafkaMessage = createKafkaMessage(
    Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(update).finish())),
  );
  message.headers = defaultKafkaHeaders;

  await onMessage(message);
}

export async function expectOrderbookLevelCache(
  ticker: string,
  orderSide: OrderSide,
  humanPrice: string,
  size: string,
): Promise<void> {
  const level: string = await OrderbookLevelsCache.getOrderbookLevel(
    ticker,
    orderSide,
    humanPrice,
    redisClient,
  );
  expect(level).toEqual(size);
}

export function setTransactionHash(
  kafkaMessage: KafkaMessage,
  txHash: Buffer,
): KafkaMessage {
  const messageWithTxhash: KafkaMessage = {
    ...kafkaMessage,
  };
  if (kafkaMessage.headers === undefined) {
    messageWithTxhash.headers = {};
  }

  messageWithTxhash.headers![DydxRecordHeaderKeys.TRANSACTION_HASH_KEY] = txHash;
  return messageWithTxhash;
}

export async function expectCanceledOrderStatus(
  orderId: string,
  canceledOrderStatus: CanceledOrderStatus,
) {
  expect(await CanceledOrdersCache.getOrderCanceledStatus(orderId, redisClient)).toEqual(
    canceledOrderStatus,
  );
}
