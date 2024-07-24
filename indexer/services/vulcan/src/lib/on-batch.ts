import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  Batch,
  EachBatchPayload,
  KafkaMessage,
} from 'kafkajs';

import config from '../config';
import { onMessage } from './on-message';

export async function onBatch(
  payload: EachBatchPayload,
): Promise<void> {
  const batch: Batch = payload.batch;
  const topic: string = batch.topic;
  const partition: string = batch.partition.toString();
  const metricTags: Record<string, string> = { topic, partition };
  if (batch.isEmpty()) {
    logger.error({
      at: 'on-batch#onBatch',
      message: 'Empty batch',
      ...metricTags,
    });
    return;
  }

  const startTime: number = Date.now();
  const firstMessageTimestamp: number = Number(batch.messages[0].timestamp);
  const batchTimeInQueue: number = startTime - firstMessageTimestamp;
  const batchInfo = {
    firstMessageTimestamp: new Date(firstMessageTimestamp).toISOString(),
    batchTimeInQueue,
    messagesInBatch: batch.messages.length,
    firstOffset: batch.firstOffset(),
    lastOffset: batch.lastOffset(),
    ...metricTags,
  };

  logger.info({
    at: 'on-batch#onBatch',
    message: 'Received batch',
    ...batchInfo,
  });
  stats.timing(
    'vulcan.batch_time_in_queue',
    batchTimeInQueue,
    metricTags,
  );

  let lastCommitTime: number = startTime;
  for (let i = 0; i < batch.messages.length; i++) {
    const message: KafkaMessage = batch.messages[i];
    await onMessage(message);

    // Commit every KAFKA_BATCH_PROCESSING_COMMIT_FREQUENCY_MS to reduce number of roundtrips, and
    // also prevent disconnecting from the broker due to inactivity.
    const now: number = Date.now();
    if (now - lastCommitTime > config.KAFKA_BATCH_PROCESSING_COMMIT_FREQUENCY_MS) {
      logger.info({
        at: 'on-batch#onBatch',
        message: 'Committing offsets and sending heart beat',
        ...batchInfo,
      });
      payload.resolveOffset(message.offset);
      await Promise.all([
        payload.heartbeat(),
        // commitOffsetsIfNecessary will respect autoCommitThreshold and will not commit if
        // fewer messages than the threshold have been processed since the last commit.
        payload.commitOffsetsIfNecessary(),
      ]);
      lastCommitTime = now;
    }
  }

  const batchProcessingTime: number = Date.now() - startTime;
  logger.info({
    at: 'on-batch#onBatch',
    message: 'Finished Processing Batch',
    batchProcessingTime,
    ...batchInfo,
  });
  stats.timing(
    'vulcan.batch_processing_time',
    batchProcessingTime,
    metricTags,
  );
  stats.timing(
    'vulcan.batch_size',
    batch.messages.length,
    metricTags,
  );
}
