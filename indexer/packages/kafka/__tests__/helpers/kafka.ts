import { KafkaMessage } from 'kafkajs';

export function createKafkaMessage(value: Buffer | null = null): KafkaMessage {
  return {
    key: Buffer.from('key'),
    value,
    timestamp: '1687515685000',
    size: 0,
    attributes: 0,
    offset: '0',
  };
}
