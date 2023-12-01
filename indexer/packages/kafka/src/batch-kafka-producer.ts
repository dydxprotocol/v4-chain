import { logger } from '@dydxprotocol-indexer/base';
import { Producer, RecordMetadata } from 'kafkajs';
import _ from 'lodash';

import { KafkaTopics } from './types';

/**
 * Single message sent to the producer.
 */
export type ProducerMessage = {
  key?: Buffer,
  value: Buffer,
};

/**
 * Groups messages for a single kafka topic into batches to send fewer ProducerRecords.
 */
export class BatchKafkaProducer {
  maxBatchSizeBytes: number;
  producer: Producer;
  topic: KafkaTopics;

  producerMessages: ProducerMessage[];
  producerPromises: Promise<RecordMetadata[]>[];
  currentSize: number;

  constructor(
    topic: KafkaTopics,
    producer: Producer,
    // Note that default parameters are bound during module load time making it difficult
    // to modify the parameter during a test so we explicitly require callers to pass in
    // config.KAFKA_MAX_BATCH_WEBSOCKET_MESSAGE_SIZE_BYTES.
    maxBatchSizeBytes: number,
  ) {
    this.maxBatchSizeBytes = maxBatchSizeBytes;
    this.producer = producer;
    this.topic = topic;

    this.producerMessages = [];
    this.producerPromises = [];
    this.currentSize = 0;
  }

  /**
   * Add a message to the current batch. If the message size would push the current batch size over
   * the maxBatchSizeBytes, the current batch (without this message) is flushed first, then the
   * message is added to a new batch.
   */
  public addMessageAndMaybeFlush(message: ProducerMessage): void {
    const keyByteLength: number = message.key === undefined ? 0 : message.key.byteLength;
    const msgBuffer: Buffer = message.value;
    if (this.currentSize + msgBuffer.byteLength + keyByteLength > this.maxBatchSizeBytes) {
      this.sendBatch();
    }
    this.producerMessages.push({ key: message.key, value: msgBuffer });
    this.currentSize += msgBuffer.byteLength;
    this.currentSize += keyByteLength;
  }

  public async flush(): Promise<RecordMetadata[][]> {
    this.sendBatch();
    // TODO(IND-198): Log an error when kafka producer fails
    return Promise.all(this.producerPromises);
  }

  private sendBatch(): void {
    if (!_.isEmpty(this.producerMessages)) {
      this.producerPromises.push(
        this.producer.send({ topic: this.topic, messages: this.producerMessages }),
      );
    }
    logger.info({
      at: 'BatchMessenger#sendBatch',
      message: 'Produced kafka batch',
      currentSize: this.currentSize,
      producerMessages: JSON.stringify(this.producerMessages),
      recalculatedCurrentSize: this.producerMessages.reduce(
        (acc: number, msg: ProducerMessage) => acc + msg.value.byteLength,
        0,
      ),
      topic: this.topic,
    });
    this.producerMessages = [];
    this.currentSize = 0;
  }
}
