import {
  consumer, startConsumer, updateOnMessageFunction, stopConsumer,
} from '../src/consumer';
import { producer } from '../src/producer';
import { createKafkaMessage } from './helpers/kafka';
import { KafkaMessage } from 'kafkajs';
import { TO_ENDER_TOPIC } from '../src';

// Skipping because timeout could cause tests to be flaky
describe.skip('consumer', () => {
  beforeAll(async () => {
    await Promise.all([
      consumer!.connect(),
      producer.connect(),
    ]);
    await consumer!.subscribe({ topic: TO_ENDER_TOPIC });
    await startConsumer();
  });

  afterAll(async () => {
    await Promise.all([
      stopConsumer(),
      producer.disconnect(),
    ]);
  });

  it('is consuming message', async () => {
    const onMessageFn: (topic: string, message: KafkaMessage) => Promise<void> = jest.fn();
    updateOnMessageFunction(onMessageFn);
    const kafkaMessage: KafkaMessage = createKafkaMessage(null);

    await producer.send({
      topic: TO_ENDER_TOPIC,
      messages: [{
        value: kafkaMessage.value,
        timestamp: `${Date.now()}`,
      }],
    });
    await new Promise((resolve) => {
      setTimeout(resolve, 2000);
    });
    expect(onMessageFn).toHaveBeenCalled();
  });
});
