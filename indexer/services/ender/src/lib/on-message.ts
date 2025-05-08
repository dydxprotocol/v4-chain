import {
  logger,
  stats,
  ParseMessageError,
  wrapBackgroundTask,
  STATS_NO_SAMPLING,
} from '@dydxprotocol-indexer/base';
import { KafkaTopics } from '@dydxprotocol-indexer/kafka';
import {
  Transaction,
  IsolationLevel,
  CandleFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerTendermintBlock,
} from '@dydxprotocol-indexer/v4-protos';
import {
  KafkaMessage,
} from 'kafkajs';
import _ from 'lodash';
import { DateTime } from 'luxon';

import {
  shouldSkipBlock,
  updateBlockCache,
} from '../caches/block-cache';
import { updateCandleCacheWithCandle } from '../caches/candle-cache';
import config from '../config';
import { BlockProcessor } from './block-processor';
import { refreshDataCaches } from './cache-manager';
import { CandlesGenerator } from './candles-generator';
import { KafkaPublisher } from './kafka-publisher';

/**
 * @function onMessage
 * @param message the kafka message being processed which should parse
 * into a valid IndexerTendermintBlock
 * @description this method will:
 * - Create a block in `blocks`
 * - Create transactions in `transactions`
 * - Create tendermint events in `tendermint_events`
 * - Create all corresponding objects for each event that occurred in this block
 */
export async function onMessage(message: KafkaMessage): Promise<void> {
  stats.increment(`${config.SERVICE_NAME}.received_kafka_message`, 1);
  const start: number = Date.now();
  const messageTime: number = Number(message.timestamp);

  const indexerTendermintBlock: IndexerTendermintBlock | undefined = getIndexerTendermintBlock(
    message,
  );
  if (indexerTendermintBlock === undefined) {
    return;
  }
  const blockTime: number = indexerTendermintBlock.time!.getTime();
  stats.timing(
    `${config.SERVICE_NAME}.block_time_lag.timing`,
    start - blockTime,
    STATS_NO_SAMPLING,
  );

  const offset = message.offset;
  const blockHeight: string = indexerTendermintBlock.height.toString();
  if (await shouldSkipBlock(blockHeight)) {
    return;
  }
  stats.timing(
    `${config.SERVICE_NAME}.message_time_in_queue`,
    start - messageTime,
    STATS_NO_SAMPLING,
    {
      topic: KafkaTopics.TO_ENDER,
    },
  );

  logger.info({
    at: 'onMessage#onMessage',
    message: 'Processing message',
    offset,
    blockHeight,
    messageTimeInQueue: start - Number(message.timestamp),
    numEvents: indexerTendermintBlock.events.length,
  });

  let success: boolean = false;
  const txId: number = await Transaction.start();
  await Transaction.setIsolationLevel(txId, IsolationLevel.READ_UNCOMMITTED);
  try {
    validateIndexerTendermintBlock(indexerTendermintBlock);

    const blockProcessor: BlockProcessor = new BlockProcessor(
      indexerTendermintBlock,
      txId,
      message.timestamp,
    );
    const kafkaPublisher: KafkaPublisher = await blockProcessor.process();

    const candlesGenerator: CandlesGenerator = new CandlesGenerator(
      kafkaPublisher,
      DateTime.fromJSDate(indexerTendermintBlock.time!),
      txId,
    );
    const candles: CandleFromDatabase[] = await candlesGenerator.updateCandles();
    await Transaction.commit(txId);
    stats.gauge(`${config.SERVICE_NAME}.processing_block_height`, indexerTendermintBlock.height);
    // Update caches after transaction is committed
    updateBlockCache(blockHeight);
    _.forEach(candles, updateCandleCacheWithCandle);

    if (config.SEND_WEBSOCKET_MESSAGES) {
      wrapBackgroundTask(
        kafkaPublisher.publish(),
        false,
        'kafkaPublisher.publish',
      );
    }
    logger.info({
      at: 'onMessage#onMessage',
      message: 'Successfully processed block',
      height: blockHeight,
    });
    success = true;
  } catch (error) {
    await Transaction.rollback(txId);
    await refreshDataCaches();
    stats.increment(`${config.SERVICE_NAME}.update_event_tables.failure`, 1);
    if (error instanceof ParseMessageError) {
      logger.crit({
        at: 'onMessage#onMessage',
        message: 'Error: Unable to parse message, this must be due to a bug in V4 node',
        offset,
        indexerTendermintBlock,
        error,
      });
    } else {
      logger.error({
        at: 'onMessage#onMessage',
        message: 'Error: Unable to process message',
        offset,
        indexerTendermintBlock,
        error,
      });
    }
    // Throw error so the message is not acked and will be reprocessed
    throw error;
  } finally {
    const done: number = Date.now();
    stats.timing(
      `${config.SERVICE_NAME}.processed_block.timing`,
      done - start,
      STATS_NO_SAMPLING,
      { success: success.toString() },
    );
    stats.timing(
      `${config.SERVICE_NAME}.processed_blocktime_lag.timing`,
      done - blockTime,
      STATS_NO_SAMPLING,
      { success: success.toString() },
    );
  }
}

/**
 * Creates an IndexerTendermintBlock from a KafkaMessage. Returns undefined, if there is an issue
 */
function getIndexerTendermintBlock(
  message: KafkaMessage,
): IndexerTendermintBlock | undefined {
  if (!message || !message.value || !message.timestamp) {
    stats.increment(`${config.SERVICE_NAME}.empty_kafka_message`, 1);
    logger.error({
      at: 'onMessage#getIndexerTendermintBlock',
      message: 'Empty message',
    });
    return undefined;
  }
  try {
    const messageValueBinary: Uint8Array = new Uint8Array(message.value);
    logger.info({
      at: 'onMessage#getIndexerTendermintBlock',
      message: 'Received message',
      offset: message.offset,
    });

    const block: IndexerTendermintBlock = IndexerTendermintBlock.decode(
      messageValueBinary,
    );
    logger.info({
      at: 'onMessage#getIndexerTendermintBlock',
      message: 'Parsed message',
      offset: message.offset,
      height: block.height,
      block,
    });

    return block;
  } catch (error) {
    stats.increment(`${config.SERVICE_NAME}.parse_kafka_message.failure`, 1);
    // Does not throw error, because we want to ack this message and skip retry
    logger.crit({
      at: 'onMessage#onMessage',
      message: 'Error: Unable to parse message',
      offset: message.offset,
      value: message.value,
      error,
    });
    return undefined;
  }
}

function validateIndexerTendermintBlock(
  indexerTendermintBlock: IndexerTendermintBlock,
) {
  if (indexerTendermintBlock.time === undefined) {
    logger.error({
      at: 'onMessage#validateIndexerTendermintBlock',
      message: 'Error: IndexerTendermintBlock.time cannot be undefined, this must be due to a bug in V4 node',
      value: indexerTendermintBlock,
    });

    throw new ParseMessageError(
      'IndexerTendermintBlock.time cannot be undefined',
    );
  }
}
