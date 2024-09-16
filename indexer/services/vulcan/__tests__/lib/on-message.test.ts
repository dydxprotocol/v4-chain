import { logger, stats } from '@dydxprotocol-indexer/base';
import { createKafkaMessage, KafkaTopics } from '@dydxprotocol-indexer/kafka';
import { OffChainUpdateV1 } from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../src/lib/on-message';
import { OrderPlaceHandler } from '../../src/handlers/order-place-handler';
import { OrderRemoveHandler } from '../../src/handlers/order-remove-handler';
import { OrderUpdateHandler } from '../../src/handlers/order-update-handler';
import { redisTestConstants } from '@dydxprotocol-indexer/redis';
import { setTransactionHash } from '../helpers/helpers';

jest.mock('../../src/handlers/order-place-handler');
jest.mock('../../src/handlers/order-remove-handler');
jest.mock('../../src/handlers/order-update-handler');

describe('onMessage', () => {
  const handlerMocks: jest.Mock[] = [
    (OrderPlaceHandler as jest.Mock),
    (OrderRemoveHandler as jest.Mock),
    (OrderUpdateHandler as jest.Mock),
  ];
  const testTxhash: Buffer = Buffer.from('testTxhash');
  let handleUpdateMock: jest.Mock;

  beforeEach(() => {
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(logger, 'crit');
    jest.spyOn(logger, 'error');

    handleUpdateMock = jest.fn();
    handlerMocks.forEach((mockHandler: jest.Mock) => {
      mockHandler.mockReturnValue({
        handleUpdate: handleUpdateMock,
      });
    });
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it.each([
    ['orderPlace', redisTestConstants.orderPlace, (OrderPlaceHandler as jest.Mock)],
    ['orderRemove', redisTestConstants.orderRemove, (OrderRemoveHandler as jest.Mock)],
    ['orderUpdate', redisTestConstants.orderUpdate, (OrderUpdateHandler as jest.Mock)],
  ])('processes updates with the correct handler: %s', async (
    messageType: string,
    updateMessage: any,
    handler: jest.Mock,
  ) => {
    const update: OffChainUpdateV1 = {
      ...updateMessage,
    };
    let message: KafkaMessage = createKafkaMessage(
      Buffer.from(
        Uint8Array.from(
          OffChainUpdateV1.encode(update).finish(),
        ),
      ),
    );
    message = setTransactionHash(message, testTxhash);

    await onMessage(message);

    expect(handler).toHaveBeenCalledTimes(1);
    expect(handleUpdateMock).toHaveBeenCalledWith(update, message.headers ?? {});
    expect(handleUpdateMock).toHaveBeenCalledTimes(1);

    expect(stats.increment).toHaveBeenCalledWith('vulcan.received_kafka_message', 1, { instance: '' });
    expect(stats.timing).toHaveBeenCalledWith(
      'vulcan.message_time_in_queue',
      expect.any(Number),
      1,
      {
        topic: KafkaTopics.TO_VULCAN,
        instance: '',
      },
    );
    expect(stats.timing).toHaveBeenCalledWith(
      'vulcan.processed_update.timing',
      expect.any(Number),
      1,
      {
        success: 'true',
        messageType,
        instance: '',
      },
    );

    handlerMocks.forEach((mockHandler: jest.Mock) => {
      if (mockHandler !== handler) {
        expect(mockHandler).not.toHaveBeenCalled();
      } else {
        expect(mockHandler).toHaveBeenCalledWith(testTxhash.toString('hex').toUpperCase());
      }
    });
  });

  it('logs error and does not process update if message in unparseable', async () => {
    const invalidMessage: KafkaMessage = createKafkaMessage(Buffer.from('abc'));

    await onMessage(invalidMessage);

    handlerMocks.forEach((mockHandler: jest.Mock) => {
      expect(mockHandler).not.toHaveBeenCalled();
    });

    expect(stats.increment).toHaveBeenCalledWith('vulcan.received_kafka_message', 1, { instance: '' });
    expect(stats.timing).toHaveBeenCalledWith(
      'vulcan.message_time_in_queue',
      expect.any(Number),
      1,
      {
        topic: KafkaTopics.TO_VULCAN,
        instance: '',
      },
    );
    expect(logger.crit).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Error: Unable to parse message',
      }),
    );
  });

  it('logs error and does not process update if OffChainUpdate is not valid', async () => {
    const invalidMessage: KafkaMessage = createKafkaMessage(
      Buffer.from(Uint8Array.from(OffChainUpdateV1.encode({}).finish())),
    );

    await onMessage(invalidMessage);

    handlerMocks.forEach((mockHandler: jest.Mock) => {
      expect(mockHandler).not.toHaveBeenCalled();
    });

    expect(stats.increment).toHaveBeenCalledWith('vulcan.received_kafka_message', 1, { instance: '' });
    expect(stats.timing).toHaveBeenCalledWith(
      'vulcan.message_time_in_queue',
      expect.any(Number),
      1,
      {
        topic: KafkaTopics.TO_VULCAN,
        instance: '',
      },
    );
    expect(stats.timing).toHaveBeenCalledWith(
      'vulcan.processed_update.timing',
      expect.any(Number),
      1,
      {
        success: 'false',
        messageType: 'unknown',
        instance: '',
      },
    );
    expect(logger.crit).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Error: Unable to parse message, this must be due to a bug in the V4 node',
      }),
    );
  });

  it('logs error and re-throws if unexpected error occurs while processing update', async () => {
    const unexpectedError: Error = new Error('Unexpected');
    handleUpdateMock.mockImplementation(() => { throw unexpectedError; });
    const message: KafkaMessage = createKafkaMessage(
      Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(redisTestConstants.orderPlace).finish())),
    );

    await expect(onMessage(message)).rejects.toEqual(unexpectedError);
    expect(stats.increment).toHaveBeenCalledWith('vulcan.received_kafka_message', 1, { instance: '' });
    expect(stats.timing).toHaveBeenCalledWith(
      'vulcan.message_time_in_queue',
      expect.any(Number),
      1,
      {
        topic: KafkaTopics.TO_VULCAN,
        instance: '',
      },
    );
    expect(stats.timing).toHaveBeenCalledWith(
      'vulcan.processed_update.timing',
      expect.any(Number),
      1,
      {
        success: 'false',
        messageType: 'orderPlace',
        instance: '',
      },
    );
    expect(logger.error).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Error: Unable to process message',
      }),
    );
  });

  it('logs error and does not process message if message is empty', async () => {
    await onMessage(undefined as any as KafkaMessage);

    expect(stats.increment).toHaveBeenCalledWith('vulcan.received_kafka_message', 1, { instance: '' });
    expect(stats.increment).toHaveBeenCalledWith('vulcan.empty_kafka_message', 1, { instance: '' });
    expect(logger.error).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Empty message',
      }),
    );
    handlerMocks.forEach((mockHandler: jest.Mock) => {
      expect(mockHandler).not.toHaveBeenCalled();
    });
  });
});
