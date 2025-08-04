import { delay, logger, startBugsnag } from '@dydxprotocol-indexer/base';
import { admin, KafkaTopics, producer } from '@dydxprotocol-indexer/kafka';
import { dbHelpers } from '@dydxprotocol-indexer/postgres';
import { redis } from '@dydxprotocol-indexer/redis';
import { APIGatewayEvent, APIGatewayProxyResult, Context } from 'aws-lambda';
import { ITopicConfig, ITopicMetadata } from 'kafkajs';
import _ from 'lodash';

import config from './config';
import { redisClient } from './redis';
import { sendStatefulOrderMessages } from './vulcan-helpers';

const KAFKA_TOPICS: KafkaTopics[] = [
  KafkaTopics.TO_ENDER,
  KafkaTopics.TO_VULCAN,
  KafkaTopics.TO_WEBSOCKETS_ORDERBOOKS,
  KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
  KafkaTopics.TO_WEBSOCKETS_TRADES,
  KafkaTopics.TO_WEBSOCKETS_MARKETS,
  KafkaTopics.TO_WEBSOCKETS_CANDLES,
  KafkaTopics.TO_WEBSOCKETS_BLOCK_HEIGHT,
];

const DEFAULT_NUM_REPLICAS: number = 3;

const KAFKA_TOPICS_TO_PARTITIONS: { [key in KafkaTopics]: number } = {
  [KafkaTopics.TO_ENDER]: 1,
  [KafkaTopics.TO_VULCAN]: 210,
  [KafkaTopics.TO_WEBSOCKETS_ORDERBOOKS]: 1,
  [KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS]: 3,
  [KafkaTopics.TO_WEBSOCKETS_TRADES]: 1,
  [KafkaTopics.TO_WEBSOCKETS_MARKETS]: 1,
  [KafkaTopics.TO_WEBSOCKETS_CANDLES]: 1,
  [KafkaTopics.TO_WEBSOCKETS_BLOCK_HEIGHT]: 1,
};

export interface BazookaEventJson {
  // Run knex migrations
  migrate: boolean,

  // Rollback the latest batch of knex migrations
  rollback: boolean,

  // Clearing data inside the database, but not deleting the tables and schemas
  clear_db: boolean,

  // Reset the database and all migrations
  reset_db: boolean,

  // Create all kafka topics with replication and parition counts
  create_kafka_topics: boolean,

  // Clearing data inside all topics, not removing the Kafka Topics
  clear_kafka_topics: boolean,

  // Clearing all data in redis
  clear_redis: boolean,

  // Force flag that is required to perform any breaking actions in testnet/mainnet
  // A breaking action is any action in bazooka other that db migration
  force: boolean,

  // Send stateful orders to Vulcan. This is done during Indexer fast sync to
  // uncross the orderbook.
  send_stateful_orders_to_vulcan: boolean,
}

// eslint-disable-next-line @typescript-eslint/require-await
export async function handler(
  event: APIGatewayEvent & BazookaEventJson,
  _context: Context,
): Promise<APIGatewayProxyResult> {
  logger.info({
    at: 'index#handler',
    message: `Event: ${JSON.stringify(event, null, 2)}`,
  });

  startBugsnag();

  if (config.PREVENT_BREAKING_CHANGES_WITHOUT_FORCE && event.force !== true) {
    if (
      event.rollback === true ||
      event.clear_db === true ||
      event.reset_db === true ||
      event.clear_kafka_topics === true ||
      event.clear_redis === true
    ) {
      logger.error({
        at: 'index#handler',
        message: 'Cannot run bazooka without force flag set to "true" because' +
        'PREVENT_BREAKING_CHANGES_WITHOUT_FORCE is enabled',
      });
      throw new Error('Cannot run bazooka without force flag set to "true"');
    }
  }

  // Reset DB and all migrations
  if (event.reset_db) {
    logger.info({
      at: 'index#handler',
      message: 'Reset database',
    });
    await dbHelpers.clearSchema();
    logger.info({
      at: 'index#handler',
      message: 'Successfully reset database',
    });
  }

  if (event.rollback) {
    logger.info({
      at: 'index#handler',
      message: 'Rolling back latest batch of database migrations',
    });
    await dbHelpers.rollback();
    logger.info({
      at: 'index#handler',
      message: 'Successfully rolled back latest batch of database migrations',
    });
  }

  if (event.migrate) {
    logger.info({
      at: 'index#handler',
      message: 'Migrating database',
    });
    await dbHelpers.migrate();
    logger.info({
      at: 'index#handler',
      message: 'Successfully migrated database',
    });
  }

  // Clear DB after migration, because clear_db requires all tables have been created
  // Doesn't run if we've already reset the db.
  if (!event.reset_db && event.clear_db) {
    logger.info({
      at: 'index#handler',
      message: 'Clearing database',
    });
    await dbHelpers.clearData();
    logger.info({
      at: 'index#handler',
      message: 'Successfully cleared database',
    });
  }

  await maybeClearAndCreateKafkaTopics(event);

  if (event.clear_redis) {
    logger.info({
      at: 'index#handler',
      message: 'Clearing redis',
    });
    await redis.deleteAllAsync(redisClient);
    logger.info({
      at: 'index#handler',
      message: 'Successfully cleared redis',
    });
  }

  if (event.send_stateful_orders_to_vulcan) {
    await producer.connect();
    logger.info({
      at: 'index#handler',
      message: 'Sending stateful orders to Vulcan',
    });
    await sendStatefulOrderMessages();
    logger.info({
      at: 'index#handler',
      message: 'Successfully sent stateful orders to Vulcan',
    });
  }

  return {
    statusCode: 200,
    body: JSON.stringify({
      message: 'success',
    }),
  };
}

// eslint-disable-next-line @typescript-eslint/require-await
async function maybeClearAndCreateKafkaTopics(
  event: APIGatewayEvent & BazookaEventJson,
): Promise<void> {
  if (!event.create_kafka_topics && !event.clear_kafka_topics) {
    return;
  }

  logger.info({
    at: 'index#maybeClearAndCreateKafkaTopics',
    message: 'Connecting to Kafka',
  });
  await admin.connect();
  logger.info({
    at: 'index#maybeClearAndCreateKafkaTopics',
    message: 'Successfully connected to Kafka',
  });

  const existingKafkaTopics: string[] = await admin.listTopics();

  if (event.create_kafka_topics) {
    await createKafkaTopics(existingKafkaTopics);
    await partitionKafkaTopics();
  }

  if (event.clear_kafka_topics) {
    await clearKafkaTopics(existingKafkaTopics);
  }
}

async function createKafkaTopics(
  existingKafkaTopics: string[],
): Promise<void> {
  const kafkaTopicsToCreate: KafkaTopics[] = [];

  _.forEach(KAFKA_TOPICS, (kafkaTopic: KafkaTopics) => {
    if (_.includes(existingKafkaTopics, kafkaTopic)) {
      logger.info({
        at: 'index#createKafkaTopics',
        message: `Cannot create kafka topic that does exist: ${kafkaTopic}`,
      });
      return;
    }

    kafkaTopicsToCreate.push(kafkaTopic);
  });

  if (existingKafkaTopics.length === 0) {
    logger.info({
      at: 'index#createKafkaTopics',
      message: 'No kafka topics to create',
    });
    return;
  }

  logger.info({
    at: 'index#createKafkaTopics',
    message: `creating topics: ${kafkaTopicsToCreate}`,
  });
  await admin.createTopics({
    topics: _.map(kafkaTopicsToCreate, (kafkaTopicToCreate: KafkaTopics): ITopicConfig => {
      return {
        topic: kafkaTopicToCreate.toString(),
        numPartitions: KAFKA_TOPICS_TO_PARTITIONS[kafkaTopicToCreate],
        replicationFactor: DEFAULT_NUM_REPLICAS,
      };
    }),
  });
  logger.info({
    at: 'index#createKafkaTopics',
    message: 'Successfully created kafka topics',
  });
}

async function partitionKafkaTopics(): Promise<void> {
  for (const kafkaTopic of KAFKA_TOPICS) {
    const topicMetadata: { topics: Array<ITopicMetadata> } = await admin.fetchTopicMetadata({
      topics: [kafkaTopic],
    });
    if (topicMetadata.topics.length === 1) {
      if (topicMetadata.topics[0].partitions.length !== KAFKA_TOPICS_TO_PARTITIONS[kafkaTopic]) {
        logger.info({
          at: 'index#partitionKafkaTopics',
          message: `Setting topic ${kafkaTopic} to ${KAFKA_TOPICS_TO_PARTITIONS[kafkaTopic]} partitions`,
        });
        await admin.createPartitions({
          validateOnly: false,
          topicPartitions: [{
            topic: kafkaTopic,
            count: KAFKA_TOPICS_TO_PARTITIONS[kafkaTopic],
          }],
        });
        logger.info({
          at: 'index#partitionKafkaTopics',
          message: `Successfully set topic ${kafkaTopic} to ${KAFKA_TOPICS_TO_PARTITIONS[kafkaTopic]} partitions`,
        });
      }
    }
  }
}

async function clearKafkaTopics(
  existingKafkaTopics: string[],
): Promise<void> {
  // Concurrent calls to clear all topics caused the failure:
  // TypeError: Cannot destructure property 'partitions' of 'high.pop(...)' as it is undefined.
  for (const topic of KAFKA_TOPICS) {
    await clearKafkaTopic(
      1,
      config.CLEAR_KAFKA_TOPIC_RETRY_MS,
      config.CLEAR_KAFKA_TOPIC_MAX_RETRIES,
      existingKafkaTopics,
      topic,
    );
  }
}

export async function clearKafkaTopic(
  attempt: number = 1,
  retryMs: number = config.CLEAR_KAFKA_TOPIC_RETRY_MS,
  maxRetries: number = config.CLEAR_KAFKA_TOPIC_MAX_RETRIES,
  existingKafkaTopics: string[],
  kafkaTopic: KafkaTopics,
): Promise<void> {
  const kafkaTopicExists: boolean = _.includes(existingKafkaTopics, kafkaTopic);

  if (!kafkaTopicExists) {
    logger.info({
      at: 'index#clearKafkaTopics',
      message: `Cannot clear kafka topic that does not exist: ${kafkaTopic}`,
    });
    return;
  }

  const topicMetadata: { topics: Array<ITopicMetadata> } = await admin.fetchTopicMetadata({
    topics: [kafkaTopic],
  });

  if (topicMetadata.topics.length !== 1) {
    logger.info({
      at: 'index#clearKafkaTopics',
      message: `Cannot clear kafka topic that does not exist: ${kafkaTopic}`,
    });
    return;
  }

  const numPartitions = topicMetadata.topics[0].partitions.length;

  logger.info({
    at: 'index#clearKafkaTopics',
    message: `Clearing kafka topic: ${kafkaTopic}`,
  });
  try {
    await admin.deleteTopicRecords({
      topic: kafkaTopic,
      partitions: _.times(
        numPartitions,
        (partition: number) => {
          // offset = '-1' to delete all available records on this partition:
          // https://kafka.js.org/docs/admin#a-name-delete-topic-records-a-delete-topic-records
          return { partition, offset: '-1' };
        },
      ),
    });
  } catch (error) {
    logger.error({
      at: 'index#clearKafkaTopics',
      message: 'Failed to delete topic records',
      topic: kafkaTopic,
      error,
      topicMetadata,
    });

    const topicOffsets: Array<{
      high: string,
      low: string,
      partition: number,
      offset: string,
    }> = await admin.fetchTopicOffsets(kafkaTopic);

    logger.error({
      at: 'index#clearKafkaTopics',
      message: 'Failed to delete topic records',
      attempt,
      topic: kafkaTopic,
      topicOffsets,
    });
    if (attempt >= maxRetries) {
      logger.crit({
        at: 'index#clearKafkaTopics',
        message: 'Failed to delete topic records after max retries',
        topic: kafkaTopic,
        topicMetadata,
      });
      throw error;
    }
    await delay(2 ** attempt * retryMs);
    return clearKafkaTopic(
      attempt + 1,
      retryMs,
      maxRetries,
      existingKafkaTopics,
      kafkaTopic,
    );
  }
  logger.info({
    at: 'index#clearKafkaTopics',
    message: `Successfully cleared kafka topic: ${kafkaTopic}`,
  });
}
