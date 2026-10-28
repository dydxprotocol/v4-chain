import {
  consumer, startConsumer, updateOnMessageFunction, stopConsumer,
} from '../src/consumer';
import { producer } from '../src/producer';
import { createKafkaMessage } from './helpers/kafka';
import {
  Consumer, ConsumerConfig, Kafka, KafkaMessage,
} from 'kafkajs';
import { TO_ENDER_TOPIC } from '../src';

describe('initConsumer', () => {
  // The config is parsed when its module loads, so load fresh copies of the client and consumer
  // with the given env vars set (or removed, for undefined), and return the consumer's options.
  async function consumerConfigWithEnv(
    env: Record<string, string | undefined>,
  ): Promise<ConsumerConfig> {
    const originalEnv: NodeJS.ProcessEnv = process.env;
    process.env = { ...originalEnv };
    Object.entries(env).forEach(([name, value]: [string, string | undefined]) => {
      if (value === undefined) {
        delete process.env[name];
      } else {
        process.env[name] = value;
      }
    });
    let isolated: { kafka: Kafka, initConsumer: () => Promise<void> } | undefined;
    try {
      jest.isolateModules(() => {
        isolated = {
          // eslint-disable-next-line @typescript-eslint/no-var-requires, global-require
          kafka: require('../src/kafka').kafka,
          // eslint-disable-next-line @typescript-eslint/no-var-requires, global-require
          initConsumer: require('../src/consumer').initConsumer,
        };
      });
    } finally {
      process.env = originalEnv;
    }
    const consumerSpy: jest.SpyInstance = jest.spyOn(isolated!.kafka, 'consumer')
      .mockReturnValue({ on: jest.fn() } as unknown as Consumer);

    await isolated!.initConsumer();

    expect(consumerSpy).toHaveBeenCalledTimes(1);
    return consumerSpy.mock.calls[0][0];
  }

  it('fetches up to 1 MiB of each partition and 4 MiB per request by default', async () => {
    const consumerConfig: ConsumerConfig = await consumerConfigWithEnv({
      KAFKA_FETCH_MAX_BYTES: undefined,
      KAFKA_FETCH_MAX_BYTES_PER_PARTITION: undefined,
    });

    expect(consumerConfig.maxBytesPerPartition).toEqual(1024 * 1024);
    expect(consumerConfig.maxBytes).toEqual(4 * 1024 * 1024);
  });

  it('uses the fetch limits set in the environment', async () => {
    const consumerConfig: ConsumerConfig = await consumerConfigWithEnv({
      KAFKA_FETCH_MAX_BYTES: String(16 * 1024 * 1024),
      KAFKA_FETCH_MAX_BYTES_PER_PARTITION: String(4 * 1024 * 1024),
    });

    expect(consumerConfig.maxBytesPerPartition).toEqual(4 * 1024 * 1024);
    expect(consumerConfig.maxBytes).toEqual(16 * 1024 * 1024);
  });
});

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
