import { getInstanceId, logger, stats } from '@dydxprotocol-indexer/base';
import {
  Batch,
  EachBatchPayload,
  KafkaMessage,
} from 'kafkajs';

import { onMessage } from './on-message';

export async function onBatch(
  payload: EachBatchPayload,
): Promise<void> {
  const batch: Batch = payload.batch;
  const topic: string = batch.topic;
  const partition: string = batch.partition.toString();
  const metricTags: Record<string, string> = { topic, partition, instance: getInstanceId() };
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
  stats.timing(
    'vulcan.batch_time_in_queue',
    batchTimeInQueue,
    metricTags,
  );

  for (let i = 0; i < batch.messages.length; i++) {
    const message: KafkaMessage = batch.messages[i];
    await onMessage(message);
    await payload.heartbeat();
    payload.resolveOffset(message.offset);
  }

  const batchProcessingTime: number = Date.now() - startTime;
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
