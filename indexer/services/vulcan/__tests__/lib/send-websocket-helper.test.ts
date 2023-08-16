import {
  delay, logger, stats, STATS_NO_SAMPLING, wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import { synchronizeWrapBackgroundTask } from '@dydxprotocol-indexer/dev';
import { producer, WebsocketTopics } from '@dydxprotocol-indexer/kafka';
import {
  flushAllQueues, sendWebsocketWrapper, sizeStat, timingStat,
} from '../../src/lib/send-websocket-helper';
import config from '../../src/config';
import { Message, ProducerRecord } from 'kafkajs';

jest.mock('@dydxprotocol-indexer/base', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/base'),
  wrapBackgroundTask: jest.fn(),
}));

describe('send-websocket-helper', () => {
  let producerSendSpy: jest.SpyInstance;
  let logErrorSpy: jest.SpyInstance;
  let statsTimingSpy: jest.SpyInstance;
  let statsHistogramSpy: jest.SpyInstance;

  beforeAll(() => {
    jest.useFakeTimers();
  });

  beforeEach(() => {
    synchronizeWrapBackgroundTask(wrapBackgroundTask);
    producerSendSpy = jest.spyOn(producer, 'send').mockReturnThis();
    logErrorSpy = jest.spyOn(logger, 'error').mockReturnThis();
    statsTimingSpy = jest.spyOn(stats, 'timing');
    statsHistogramSpy = jest.spyOn(stats, 'histogram');
  });

  afterEach(() => {
    jest.clearAllTimers();
  });

  afterAll(() => {
    jest.useRealTimers();
  });

  describe('flushAllQueues', () => {
    it('sends messages for all message queues', async () => {
      const expectedMessagesSent: {[topic: string]: ProducerRecord} = {};
      Object.values(WebsocketTopics).forEach((topic: string) => {
        const messages: Buffer[] = [];
        for (let i: number = 0; i < config.MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC; i++) {
          const message: Buffer = Buffer.from(i.toString());
          sendWebsocketWrapper(message, topic);
          messages.push(message);
        }
        expectedMessagesSent[topic] = {
          topic,
          messages: messages.map((message: Buffer): Message => { return { value: message }; }),
        };
      });

      await flushAllQueues();

      expect(producerSendSpy).toBeCalledTimes(Object.keys(expectedMessagesSent).length);
      expect(statsTimingSpy).toBeCalledTimes(Object.keys(expectedMessagesSent).length);
      expect(statsHistogramSpy).toBeCalledTimes(Object.keys(expectedMessagesSent).length);
      Object.keys(expectedMessagesSent).forEach((topic: string) => {
        expect(producerSendSpy).toHaveBeenCalledWith(expectedMessagesSent[topic]);
        expectStats(
          statsTimingSpy,
          statsHistogramSpy,
          topic,
          config.MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC,
          true,
        );
      });
    });
  });

  describe('sendWebsocketWrapper', () => {
    it('sends messages for a topic on an interval in batches', async () => {
      const messageVal1: Buffer = Buffer.from('some message');
      const messageVal2: Buffer = Buffer.from('another message');
      const topic: string = 'some-topic';
      const expectedMessage: ProducerRecord = {
        topic,
        messages: [{ value: messageVal1 }, { value: messageVal2 }],
      };

      sendWebsocketWrapper(messageVal1, topic);

      // No messages should be sent if no timers have been run
      expect(producerSendSpy).not.toHaveBeenCalled();

      sendWebsocketWrapper(messageVal2, topic);

      // No messages should be sent if no timers have been run
      expect(producerSendSpy).not.toHaveBeenCalled();

      jest.runOnlyPendingTimers();
      // Both messages should be sent in one batch
      expect(producerSendSpy).toHaveBeenCalledTimes(1);
      expect(producerSendSpy).toHaveBeenCalledWith(expectedMessage);

      // Wait for mock producer.send function to complete
      await delay(1);
      expect(statsTimingSpy).toBeCalledTimes(1);
      expect(statsHistogramSpy).toBeCalledTimes(1);
      expectStats(statsTimingSpy, statsHistogramSpy, topic, 2, true);
    });

    it(
      'sends messages for a topic when number of messages cross configure max threshold',
      async () => {
        const topic: string = 'some-topic';
        const expectedMessage: ProducerRecord = sendMessagesForTest(
          config.MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC + 1,
          topic,
        );

        // All messages should be sent in one batch, no timers need to be run
        expect(producerSendSpy).toHaveBeenCalledTimes(1);
        expect(producerSendSpy).toHaveBeenCalledWith(expectedMessage);

        // Wait for mock producer.send function to complete
        await delay(1);
        expect(statsTimingSpy).toBeCalledTimes(1);
        expect(statsHistogramSpy).toBeCalledTimes(1);
        expectStats(
          statsTimingSpy,
          statsHistogramSpy,
          topic,
          config.MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC + 1,
          true,
        );

        // Run any remaining timers
        jest.runOnlyPendingTimers();
      },
    );

    it('logs errors and re-enqueues messages if sending failed', async () => {
      producerSendSpy
        .mockImplementationOnce(() => { throw new Error(); })
        .mockImplementationOnce(() => undefined);
      const topic: string = 'some-topic';
      const expectedMessage: ProducerRecord = sendMessagesForTest(
        config.MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC,
        topic,
      );

      // First attempt to send messages to the producer errors, and should log an error
      jest.runOnlyPendingTimers();
      expect(producerSendSpy).toHaveBeenCalledTimes(1);

      // Wait for mock producer.send function to complete
      await delay(1);
      expect(statsTimingSpy).toBeCalledTimes(1);
      expect(statsHistogramSpy).toBeCalledTimes(1);
      expectStats(
        statsTimingSpy,
        statsHistogramSpy,
        topic,
        config.MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC,
        false,
      );
      expect(logErrorSpy).toHaveBeenCalledTimes(1);
      expect(logErrorSpy).toHaveBeenCalledWith(expect.objectContaining({
        message: 'Failed to send messages to Kafka',
        topic,
      }));

      // Second attempt should succeed and re-send the previous messages
      jest.runOnlyPendingTimers();
      expect(producerSendSpy).toHaveBeenCalledTimes(2);
      expect(producerSendSpy).toHaveBeenNthCalledWith(2, expectedMessage);

      // Wait for mock producer.send function to complete
      await delay(1);
      expect(statsTimingSpy).toBeCalledTimes(2);
      expect(statsHistogramSpy).toBeCalledTimes(2);
      expectStats(
        statsTimingSpy,
        statsHistogramSpy,
        topic,
        config.MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC,
        true,
      );
    });

    it('respects SEND_WEBSOCKET_MESSAGES flag', () => {
      config.SEND_WEBSOCKET_MESSAGES = false;
      const messageVal1: Buffer = Buffer.from('some message');
      const topic: string = 'some-topic';
      sendWebsocketWrapper(messageVal1, topic);

      jest.runOnlyPendingTimers();
      // Both messages should be sent in one batch
      expect(producerSendSpy).toHaveBeenCalledTimes(0);
    });
  });
});

function sendMessagesForTest(numMessages: number, topic: string): ProducerRecord {
  const expectedMessage: ProducerRecord = { topic, messages: [] };
  for (let i: number = 0; i < numMessages; i++) {
    const messageVal: Buffer = Buffer.from(i.toString());
    sendWebsocketWrapper(messageVal, topic);
    expectedMessage.messages.push({ value: messageVal });
  }
  return expectedMessage;
}

function expectStats(
  timingSpy: jest.SpyInstance,
  histogramSpy: jest.SpyInstance,
  topic: string,
  size: number,
  success: boolean,
): void {
  const tags: {[name: string]: string} = {
    topic,
    success: success.toString(),
  };
  expect(timingSpy).toHaveBeenCalledWith(timingStat, expect.any(Number), STATS_NO_SAMPLING, tags);
  expect(histogramSpy).toHaveBeenCalledWith(sizeStat, size, STATS_NO_SAMPLING, tags);
}
