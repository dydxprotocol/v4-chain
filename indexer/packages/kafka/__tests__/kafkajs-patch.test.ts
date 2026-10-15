// Checks the follower fetching (rackId) fixes in patches/kafkajs+2.2.4.patch are applied.
// eslint-disable-next-line @typescript-eslint/no-var-requires
const ConsumerGroup = require('kafkajs/src/consumer/consumerGroup');

function makeConsumerGroup(
  offlineReplicas: number[],
  preferred: { nodeId: number, expireAt: number },
) {
  const consumerGroup = Object.create(ConsumerGroup.prototype);
  consumerGroup.cluster = {
    findTopicPartitionMetadata: () => [{ partitionId: 0, leader: 1, offlineReplicas }],
    brokerPool: { metadataExpireAt: Date.now() + 300_000 },
  };
  consumerGroup.logger = { info: jest.fn(), debug: jest.fn(), error: jest.fn() };
  consumerGroup.preferredReadReplicasPerTopicPartition = { topic: { 0: preferred } };
  consumerGroup.offsetManager = { setDefaultOffset: jest.fn() };
  return consumerGroup;
}

describe('kafkajs patch: follower fetching', () => {
  const expireAt: number = Date.now() + 60_000;

  it('fetches from the preferred replica while it is online', () => {
    const consumerGroup = makeConsumerGroup([3], { nodeId: 2, expireAt });
    expect(consumerGroup.findReadReplicaForPartitions('topic', [0])).toEqual({ 2: [0] });
    expect(consumerGroup.preferredReadReplicasPerTopicPartition.topic[0]).toBeDefined();
    expect(consumerGroup.cluster.brokerPool.metadataExpireAt).toBeGreaterThan(Date.now() + 1000);
  });

  it('falls back to the leader and refreshes metadata within a second when the replica is offline', () => {
    const consumerGroup = makeConsumerGroup([2], { nodeId: 2, expireAt });
    expect(consumerGroup.findReadReplicaForPartitions('topic', [0])).toEqual({ 1: [0] });
    expect(consumerGroup.cluster.brokerPool.metadataExpireAt)
      .toBeLessThanOrEqual(Date.now() + 1000);
  });

  it('retries the leader instead of resetting the offset when out of range on a follower', async () => {
    const consumerGroup = makeConsumerGroup([], { nodeId: 2, expireAt });
    await consumerGroup.recoverFromOffsetOutOfRange({ topic: 'topic', partition: 0 });
    expect(consumerGroup.preferredReadReplicasPerTopicPartition.topic[0]).toBeUndefined();
    expect(consumerGroup.offsetManager.setDefaultOffset).not.toHaveBeenCalled();
  });
});
