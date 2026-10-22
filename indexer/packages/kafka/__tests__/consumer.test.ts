import { execSync } from 'child_process';
import path from 'path';

import { Admin, KafkaMessage } from 'kafkajs';

import { admin } from '../src/admin';
import config from '../src/config';
import {
  consumer,
  getConsecutiveCrashes,
  handleConsumerCrash,
  initConsumer,
  recordCrashAndDecideRestart,
  resetConsecutiveCrashes,
  startConsumer,
  stopConsumer,
  updateOnMessageFunction,
} from '../src/consumer';
import { kafka } from '../src/kafka';
import { producer, disconnect as disconnectProducer } from '../src/producer';

const COMPOSE_FILE: string = path.resolve(process.cwd(), '../../docker-compose.yml');

const sleep = (ms: number): Promise<void> => new Promise((resolve) => {
  setTimeout(resolve, ms);
});

async function waitFor(
  predicate: () => boolean,
  timeoutMs: number,
  intervalMs: number = 500,
): Promise<void> {
  const deadline: number = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    if (predicate()) {
      return;
    }
    // eslint-disable-next-line no-await-in-loop
    await sleep(intervalMs);
  }
  if (!predicate()) {
    throw new Error('waitFor timed out');
  }
}

function mockProcessExit(): jest.SpyInstance {
  return jest.spyOn(process, 'exit').mockImplementation(
    (() => undefined) as unknown as (code?: number) => never,
  );
}

describe('consumer crash handling', () => {
  const originalMaxCrashes: number = config.KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES;
  const originalExitOnCrash: boolean = config.KAFKA_EXIT_PROCESS_ON_CONSUMER_CRASH;

  beforeEach(() => {
    resetConsecutiveCrashes();
  });

  afterEach(() => {
    config.KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES = originalMaxCrashes;
    config.KAFKA_EXIT_PROCESS_ON_CONSUMER_CRASH = originalExitOnCrash;
  });

  it('restarts internally until the consecutive-crash budget is exhausted', () => {
    config.KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES = 3;
    const error: Error = new Error('boom');

    expect(recordCrashAndDecideRestart(error)).toBe(true); // crash 1
    expect(recordCrashAndDecideRestart(error)).toBe(true); // crash 2
    expect(recordCrashAndDecideRestart(error)).toBe(false); // crash 3 -> budget exhausted
    expect(getConsecutiveCrashes()).toBe(3);
  });

  it('resets the crash count after progress, restoring the full budget', () => {
    config.KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES = 3;

    recordCrashAndDecideRestart(new Error('boom'));
    expect(getConsecutiveCrashes()).toBe(1);

    resetConsecutiveCrashes();
    expect(getConsecutiveCrashes()).toBe(0);

    // Budget is full again after a successful message reset the counter.
    expect(recordCrashAndDecideRestart(new Error('boom'))).toBe(true);
  });

  it('force-exits the process on a fatal (non-restartable) crash', () => {
    const exitSpy: jest.SpyInstance = mockProcessExit();

    handleConsumerCrash({ error: new Error('fatal'), restart: false });

    expect(exitSpy).toHaveBeenCalledWith(1);
    exitSpy.mockRestore();
  });

  it('does not exit while the crash is still restartable', () => {
    const exitSpy: jest.SpyInstance = mockProcessExit();

    handleConsumerCrash({ error: new Error('transient'), restart: true });

    expect(exitSpy).not.toHaveBeenCalled();
    exitSpy.mockRestore();
  });

  it('does not exit when process-exit-on-crash is disabled', () => {
    config.KAFKA_EXIT_PROCESS_ON_CONSUMER_CRASH = false;
    const exitSpy: jest.SpyInstance = mockProcessExit();

    handleConsumerCrash({ error: new Error('fatal'), restart: false });

    expect(exitSpy).not.toHaveBeenCalled();
    exitSpy.mockRestore();
  });
});

describe('consumer (integration - requires kafka broker)', () => {
  const originalMaxRetries: number = config.KAFKA_CONSUMER_MAX_RETRIES;
  const originalMaxCrashes: number = config.KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES;
  let startedBroker: boolean = false;
  let topicCounter: number = 0;

  async function isBrokerReachable(): Promise<boolean> {
    const probe: Admin = kafka.admin();
    try {
      await probe.connect();
      await probe.listTopics();
      return true;
    } catch (error) {
      return false;
    } finally {
      await probe.disconnect().catch(() => undefined);
    }
  }

  async function ensureBroker(): Promise<void> {
    if (await isBrokerReachable()) {
      return;
    }
    execSync(`docker compose -f "${COMPOSE_FILE}" up -d kafka`, { stdio: 'ignore' });
    startedBroker = true;
    for (let attempt: number = 0; attempt < 60; attempt += 1) {
      // eslint-disable-next-line no-await-in-loop
      if (await isBrokerReachable()) {
        return;
      }
      // eslint-disable-next-line no-await-in-loop
      await sleep(1000);
    }
    throw new Error('Kafka broker did not become reachable after docker compose up');
  }

  async function createTestTopic(): Promise<string> {
    const topic: string = `test-consumer-${Date.now()}-${topicCounter}`;
    topicCounter += 1;
    await admin.createTopics({ topics: [{ topic, numPartitions: 1 }], waitForLeaders: true });
    return topic;
  }

  beforeAll(async () => {
    await ensureBroker();
    await admin.connect();
    await producer.connect();
  }, 120000);

  afterAll(async () => {
    config.KAFKA_CONSUMER_MAX_RETRIES = originalMaxRetries;
    config.KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES = originalMaxCrashes;
    await stopConsumer().catch(() => undefined);
    await admin.disconnect().catch(() => undefined);
    await disconnectProducer().catch(() => undefined);
    if (startedBroker) {
      execSync(`docker compose -f "${COMPOSE_FILE}" stop kafka`, { stdio: 'ignore' });
    }
  }, 60000);

  it('consumes produced messages and keeps the crash count at zero', async () => {
    const topic: string = await createTestTopic();
    const received: KafkaMessage[] = [];
    updateOnMessageFunction((_topic: string, message: KafkaMessage): Promise<void> => {
      received.push(message);
      return Promise.resolve();
    });

    await initConsumer();
    await consumer!.subscribe({ topic, fromBeginning: true });
    await startConsumer();

    await producer.send({ topic, messages: [{ value: Buffer.from('hello') }] });

    await waitFor(() => received.length > 0, 30000);
    expect(received.length).toBeGreaterThan(0);
    expect(getConsecutiveCrashes()).toBe(0);

    await stopConsumer();
  }, 60000);

  it('force-exits the process after sustained processing failures', async () => {
    // retries: 0 => crash on the first failure with no backoff; budget of 2 => exit on 2nd crash.
    config.KAFKA_CONSUMER_MAX_RETRIES = 0;
    config.KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES = 2;
    const exitSpy: jest.SpyInstance = mockProcessExit();

    const topic: string = await createTestTopic();
    updateOnMessageFunction((): Promise<void> => Promise.reject(
      new Error('simulated postgres outage'),
    ));

    await initConsumer();
    await consumer!.subscribe({ topic, fromBeginning: true });
    await startConsumer();

    // The offset is never committed (processing always throws), so the same message is redelivered
    // after each internal restart, driving repeated crashes until the budget is exhausted.
    await producer.send({ topic, messages: [{ value: Buffer.from('poison') }] });

    await waitFor(() => exitSpy.mock.calls.length > 0, 60000);
    expect(exitSpy).toHaveBeenCalledWith(1);
    expect(getConsecutiveCrashes()).toBeGreaterThanOrEqual(2);

    await stopConsumer().catch(() => undefined);
    exitSpy.mockRestore();
  }, 90000);
});
