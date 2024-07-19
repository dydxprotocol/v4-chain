import { connect } from '../../src/helpers/kafka/kafka-controller';
import * as onMessageModule from '../../src/lib/on-message';
import {
  producer,
  startConsumer,
  stopConsumer,
  TO_ENDER_TOPIC,
  createKafkaMessage,
} from '@dydxprotocol-indexer/kafka';
import { KafkaMessage } from 'kafkajs';

// Skipping because timeout could cause tests to be flaky
describe.skip('consumer', () => {
  beforeAll(async () => {
    await connect();
    await startConsumer();
  });

  afterAll(async () => {
    await producer.disconnect();
    await stopConsumer();
  });

  it('is consuming message', async () => {
    jest.spyOn(onMessageModule, 'onMessage');
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
    expect(onMessageModule.onMessage).toHaveBeenCalled();
  });
});
