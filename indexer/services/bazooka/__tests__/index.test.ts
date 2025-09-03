import { admin, KafkaTopics } from '@dydxprotocol-indexer/kafka';
import { BazookaEventJson, clearKafkaTopic, handler } from '../src';
import { APIGatewayEvent, Context } from 'aws-lambda';
import config from '../src/config';

describe('index', () => {
  describe('clearKafkaTopic', () => {
    const adminDeleteSpy: jest.SpyInstance = jest.spyOn(admin, 'deleteTopicRecords');
    const adminFetchSpy: jest.SpyInstance = jest.spyOn(admin, 'fetchTopicOffsets');
    const fetchTopicMetadataSpy: jest.SpyInstance = jest.spyOn(admin, 'fetchTopicMetadata');
    adminFetchSpy.mockResolvedValue([]);
    fetchTopicMetadataSpy.mockResolvedValue({});

    beforeEach(() => {
      adminDeleteSpy.mockReset();
      adminFetchSpy.mockReset();
      fetchTopicMetadataSpy.mockReset();
    });

    it('successfully clears with retry', async () => {
      adminDeleteSpy
        .mockRejectedValueOnce(new Error('test'))
        .mockResolvedValueOnce(Promise<void>);
      fetchTopicMetadataSpy
        .mockResolvedValue({
          topics: [
            {
              topic: KafkaTopics.TO_ENDER,
              partitions: [
                {},
                {},
              ],
            },
          ],
        });

      await clearKafkaTopic(1,
        5,
        3,
        [KafkaTopics.TO_ENDER],
        KafkaTopics.TO_ENDER);
      expect(adminDeleteSpy).toHaveBeenCalledTimes(2);
    });

    it('throws error after max retry', async () => {
      adminDeleteSpy
        .mockRejectedValueOnce(new Error('test'))
        .mockRejectedValueOnce(new Error('test'))
        .mockRejectedValueOnce(new Error('test'))
        .mockRejectedValueOnce(new Error('test'));
      fetchTopicMetadataSpy
        .mockResolvedValue({
          topics: [
            {
              topic: KafkaTopics.TO_ENDER,
              partitions: [
                {},
                {},
              ],
            },
          ],
        });

      await expect(async () => {
        await clearKafkaTopic(1,
          5,
          3,
          [KafkaTopics.TO_ENDER],
          KafkaTopics.TO_ENDER);
      }).rejects.toThrowError('test');
      expect(adminDeleteSpy).toHaveBeenCalledTimes(3);
    });
  });

  describe('handler', () => {
    afterEach(() => {
      config.PREVENT_BREAKING_CHANGES_WITHOUT_FORCE = false;
    });

    it('smoke test', async () => {
      await handler({
        migrate: false,
        rollback: false,
        clear_db: false,
        reset_db: false,
        create_kafka_topics: false,
        clear_kafka_topics: false,
        clear_redis: false,
        force: false,
      } as APIGatewayEvent & BazookaEventJson, {} as Context);
    });

    it.each([
      [{
        migrate: false,
        rollback: true,
        clear_db: false,
        reset_db: false,
        create_kafka_topics: false,
        clear_kafka_topics: false,
        clear_redis: false,
        force: false,
      } as APIGatewayEvent & BazookaEventJson],
      [{
        migrate: false,
        rollback: false,
        clear_db: true,
        reset_db: false,
        create_kafka_topics: false,
        clear_kafka_topics: false,
        clear_redis: false,
        force: false,
      } as APIGatewayEvent & BazookaEventJson],
      [{
        migrate: false,
        rollback: false,
        clear_db: false,
        reset_db: true,
        create_kafka_topics: false,
        clear_kafka_topics: false,
        clear_redis: false,
        force: false,
      } as APIGatewayEvent & BazookaEventJson],
      [{
        migrate: false,
        rollback: false,
        clear_db: false,
        reset_db: false,
        create_kafka_topics: false,
        clear_kafka_topics: true,
        clear_redis: false,
        force: false,
      } as APIGatewayEvent & BazookaEventJson],
      [{
        migrate: false,
        rollback: false,
        clear_db: false,
        reset_db: false,
        create_kafka_topics: false,
        clear_kafka_topics: false,
        clear_redis: true,
        force: false,
      } as APIGatewayEvent & BazookaEventJson],
    ])('Throws error if attempting to force non migrate', async (event: APIGatewayEvent & BazookaEventJson) => {
      config.PREVENT_BREAKING_CHANGES_WITHOUT_FORCE = true;
      await expect(
        handler(event, {} as Context),
      ).rejects.toThrowError('Cannot run bazooka without force flag set to "true"');
    });
  });
});
