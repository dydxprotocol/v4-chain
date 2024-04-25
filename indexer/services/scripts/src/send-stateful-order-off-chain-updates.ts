import { logger } from '@dydxprotocol-indexer/base';
import {
  BatchKafkaProducer,
  KafkaTopics,
  producer,
  ProducerMessage,
  updateOnMessageFunction,
  consumer,
  startConsumer,
  stopConsumer,
  TO_ENDER_TOPIC,
} from '@dydxprotocol-indexer/kafka';
import {
  OrderFromDatabase,
  OrderSide,
  OrderTable,
  orderTranslations,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerOrderId,
  OffChainUpdateV1,
} from '@dydxprotocol-indexer/v4-protos';
import { Long } from '@dydxprotocol-indexer/v4-protos/build/codegen/helpers';
import { IHeaders } from 'kafkajs';
import _ from 'lodash';
import {
  redis as redisLib,
  OrderbookLevelsCache,
  OrderbookLevels,
  PriceLevel,
} from '@dydxprotocol-indexer/redis'
import {
  RedisClient
} from 'redis';

import config from './config';

import yargs from 'yargs';

import { validatePnl, validatePnlForSubaccount } from './helpers/pnl-validation-helpers';
import { runAsyncScript } from './helpers/util';

interface VulcanMessage {
  key: Buffer,
  value: OffChainUpdateV1,
  headers?: IHeaders,
}

type IndexerOrderIdMap = { [orderUuid: string]: IndexerOrderId };

const res: {
  client: RedisClient,
  connect: () => Promise<void>,
} = redisLib.createRedisClient(config.REDIS_URL, config.REDIS_RECONNECT_TIMEOUT_MS);

const redisClient: RedisClient = res.client;
const connect = res.connect;

/**
 * Sends stateful order messages to Vulcan for all open stateful orders in database.
 *
 */
export async function sendStatefulOrderMessages() {
  try {
    const orders: OrderFromDatabase[] = await
    OrderTable.findOpenLongTermOrConditionalOrders();
    console.log(`Found ${orders.length} open orders.`)
    const books: any = {}
    let missingLevels: number = 0;
    await perpetualMarketRefresher.updatePerpetualMarkets();
    for (const order of orders) {
      const market: PerpetualMarketFromDatabase = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(order.clobPairId)!;
      if (books[order.clobPairId] === undefined) {
        const book: any = await OrderbookLevelsCache.getOrderBookLevels(
          market.ticker,
          redisClient,
        );
        books[order.clobPairId] = book;
        console.log(`${market.ticker} book cached`);
      }
      const book: OrderbookLevels = books[order.clobPairId];
      let side: PriceLevel[] = [];
      if (order.side === OrderSide.BUY) {
        side = book.bids;
      } else {
        side = book.asks;
      }
      for (const level of side) {
        if (order.price === level.humanPrice) {
          continue
        }
      }
      missingLevels += 1;
    }
    console.log(`Missing levels for ${missingLevels} orders`);
    // order uuid -> total filled
    const idToOrderMap: _.Dictionary<OrderFromDatabase> = _.keyBy(orders, 'id');
    const totalFilledQuantumsMap: _.Dictionary<string> = _.mapValues(
      idToOrderMap,
      (order: OrderFromDatabase) => {
        const market: PerpetualMarketFromDatabase = perpetualMarketRefresher
          .getPerpetualMarketFromClobPairId(order.clobPairId)!;
        return protocolTranslations.humanToQuantums(
          order.totalFilled, market.atomicResolution,
        ).toString();
      });

    const indexerOrders: IndexerOrder[] = await Promise.all(
      _.map(
        orders,
        (order: OrderFromDatabase) => {
          const market: PerpetualMarketFromDatabase = perpetualMarketRefresher
            .getPerpetualMarketFromClobPairId(order.clobPairId)!;
          return orderTranslations.convertToIndexerOrder(order, market);
        },
      ),
    );
    const indexerOrderIdMap: IndexerOrderIdMap = _.reduce(
      indexerOrders,
      (orderMap: IndexerOrderIdMap, indexerOrder) => {
        const orderId: IndexerOrderId = indexerOrder.orderId!;
        const uuid: string = OrderTable.orderIdToUuid(orderId);
        // eslint-disable-next-line no-param-reassign
        orderMap[uuid] = orderId;
        return orderMap;
      },
      {},
    );

    const offchainUpdates: OffChainUpdateV1[] = _.map(
      indexerOrderIdMap,
      (orderId, uuid) => {
        if (
          totalFilledQuantumsMap[uuid] === undefined
        ) {
          return undefined;
        }
        return OffChainUpdateV1.fromPartial({
          orderUpdate: {
            orderId,
            totalFilledQuantums: Long.fromValue(totalFilledQuantumsMap[uuid], true),
          },
        });
      }).filter(Boolean) as OffChainUpdateV1[];

    const vulcanMessages: VulcanMessage[] = _.map(
      offchainUpdates,
      (offChainUpdate: OffChainUpdateV1) => {
        if (offChainUpdate.orderUpdate !== undefined) {
          return {
            key: getOrderIdHash(offChainUpdate.orderUpdate!.orderId!),
            value: offChainUpdate,
          };
        }
        throw new Error(`Invalid offchain update: ${offChainUpdate}`);
      });

    const messages: ProducerMessage[] = _.map(vulcanMessages, (message: VulcanMessage) => {
      return {
        key: message.key,
        value: Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(message.value).finish())),
        headers: message.headers,
      };
    });

    const batchProducer: BatchKafkaProducer = new BatchKafkaProducer(
      KafkaTopics.TO_VULCAN,
      producer,
      config.KAFKA_MAX_BATCH_WEBSOCKET_MESSAGE_SIZE_BYTES,
    );
    for (const message of messages) {
      batchProducer.addMessageAndMaybeFlush(message);
    }
    await batchProducer.flush();
    console.log(`Sent ${vulcanMessages.length} off-chain messages.`)
  } catch (error) {
    logger.error({
      at: 'vulcan-helpers#sendStatefulOrderMessages',
      message: 'Error sending stateful order messages to Vulcan',
      error,
    });
  }
}

async function startKafka(): Promise<void> {
  await Promise.all([
    producer.connect(),
  ]);

  logger.info({
    at: 'index#start',
    message: 'Successfully started',
  });
}

async function startRedis(): Promise<void> {
  await Promise.all([
    connect(),
  ]);

  logger.info({
    at: 'index$start',
    message: 'Successfully connected to redis',
  });
}

runAsyncScript(async () => {
  await startKafka()
  await sendStatefulOrderMessages();
});
