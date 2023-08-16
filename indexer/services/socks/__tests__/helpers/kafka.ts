import { KafkaMessage } from 'kafkajs';

export function createKafkaMessage(value: Buffer | null): KafkaMessage {
  return {
    key: Buffer.from('key'),
    value,
    timestamp: 'timestamp',
    size: 0,
    attributes: 0,
    offset: '0',
  };
}
