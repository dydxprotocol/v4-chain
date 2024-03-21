import { KafkaTopics } from '../src';
import { BatchKafkaProducer, ProducerMessage } from '../src/batch-kafka-producer';
import { producer } from '../src/producer';
import { IHeaders } from 'kafkajs';
import _ from 'lodash';

interface TestMessage {
  key?: string,
  value: string,
  headers?: IHeaders,
}

function testMessage2ProducerMessage(data: TestMessage): ProducerMessage {
  const key: Buffer | undefined = data.key === undefined ? undefined : Buffer.from(data.key);
  return { key, value: Buffer.from(data.value), headers: data.headers };
}

function testMessage2ProducerMessages(data: TestMessage[]): ProducerMessage[] {
  return _.map(data, (d) => testMessage2ProducerMessage(d));
}

describe('batch-kafka-producer', () => {
  let producerSendMock: jest.SpyInstance;
  beforeAll(() => {
    producerSendMock = jest.spyOn(producer, 'send');
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  afterAll(() => {
    jest.resetAllMocks();
  });

  it.each([
    [
      'will send key if key is not undefined',
      5,
      [{ key: '1', value: 'a' }, { key: '2', value: 'b' }, { key: '3', value: 'c', headers: { timestamp: 'value' } }],
      [[{ key: '1', value: 'a' }, { key: '2', value: 'b' }]],
      [{ key: '3', value: 'c', headers: { timestamp: 'value' } }],
    ],
    [
      'will not send message until the batch size is reached',
      5,
      [{ value: 'a' }, { value: 'b' }, { value: 'c' }, { value: 'd' }],
      [],
      [{ value: 'a' }, { value: 'b' }, { value: 'c' }, { value: 'd' }],
    ],
    [
      'will send message when new message would surpass buffer size',
      5,
      [{ value: 'a' }, { value: 'b' }, { value: 'c' }, { value: 'd' }, { value: 'e' }, { value: 'f' }],
      [[{ value: 'a' }, { value: 'b' }, { value: 'c' }, { value: 'd' }, { value: 'e' }]],
      [{ value: 'f' }],
    ],
    [
      'maxBatchSize uses bytelength of input message to determine current batch fill size',
      5,
      [{ value: 'hiya' }, { value: 'there' }, { value: 'how' }, { value: 'are' }, { value: 'you' }],
      [
        [{ value: 'hiya' }],
        [{ value: 'there' }],
        [{ value: 'how' }],
        [{ value: 'are' }],
      ],
      [{ value: 'you' }],
    ],
    [
      'will batch messages that fit within maxBatchSize',
      6,
      [
        { value: 'hiya' },
        { value: 'a' },
        { value: 'b' },
        { value: 'there' },
        { value: 'c' },
        { value: 'd' },
        { value: 'how' },
        { value: 'e' },
        { value: 'f' },
        { value: 'are' },
        { value: 'g' },
        { value: 'h' },
        { value: 'you' },
        { value: 'i' },
      ],
      [
        [{ value: 'hiya' }, { value: 'a' }, { value: 'b' }],
        [{ value: 'there' }, { value: 'c' }],
        [{ value: 'd' }, { value: 'how' }, { value: 'e' }, { value: 'f' }],
        [{ value: 'are' }, { value: 'g' }, { value: 'h' }],
      ],
      [{ value: 'you' }, { value: 'i' }],
    ],
  ])('%s', async (
    _name: string,
    batchSize: number,
    messages: TestMessage[],
    expectedMessagesPerCall: TestMessage[][],
    expectedMessagesOnFlush: TestMessage[],
  ) => {
    const topic: KafkaTopics = KafkaTopics.TO_VULCAN;
    const batchProducer: BatchKafkaProducer = new BatchKafkaProducer(topic, producer, batchSize);

    for (const msg of messages) {
      const key: Buffer | undefined = msg.key === undefined ? undefined : Buffer.from(msg.key);
      batchProducer.addMessageAndMaybeFlush(
        { value: Buffer.from(msg.value), key, headers: msg.headers },
      );
    }

    expect(producerSendMock.mock.calls).toHaveLength(expectedMessagesPerCall.length);
    for (const [index, expectedMessages] of expectedMessagesPerCall.entries()) {
      expect(producerSendMock.mock.calls[index]).toEqual([
        { topic, messages: testMessage2ProducerMessages(expectedMessages) },
      ]);
    }

    await batchProducer.flush();
    expect(producerSendMock.mock.lastCall).toEqual([
      { topic, messages: testMessage2ProducerMessages(expectedMessagesOnFlush) },
    ]);
  });
});
