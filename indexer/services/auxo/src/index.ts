import { ECRClient } from '@aws-sdk/client-ecr';
import { ECSClient } from '@aws-sdk/client-ecs';
import { LambdaClient } from '@aws-sdk/client-lambda';
import { logger, startBugsnag } from '@dydxprotocol-indexer/base';
import { APIGatewayEvent, APIGatewayProxyResult, Context } from 'aws-lambda';

import {
  upgradeBazooka, runDbAndKafkaMigration, createNewEcsTaskDefinitions, upgradeEcsServices,
} from './helpers';
import { AuxoEventJson, TaskDefinitionArnMap } from './types';

/**
 * Upgrades all services and run migrations
 * 1. Upgrade Bazooka
 * 2. Run db migration in Bazooka, and update kafka topics
 * 3. Create new ECS Task Definition for ECS Services with new image
 * 4. Upgrade all ECS Services (comlink, ender, roundtable, socks, vulcan)
 */
// eslint-disable-next-line @typescript-eslint/require-await
export async function handler(
  event: APIGatewayEvent & AuxoEventJson,
  _context: Context,
): Promise<APIGatewayProxyResult> {
  logger.info({
    at: 'index#handler',
    message: `Event: ${JSON.stringify(event, null, 2)}`,
  });

  startBugsnag();

  const region = event.region;

  try {
    // Initialize clients
    const ecs: ECSClient = new ECSClient({ region });
    const lambda: LambdaClient = new LambdaClient({ region });
    const ecr: ECRClient = new ECRClient({ region });
    // 1. Upgrade Bazooka
    await upgradeBazooka(lambda, ecr, event);

    // 2. Run db migration in Bazooka,
    // boolean flag used to determine if new kafka topics should be created
    await runDbAndKafkaMigration(event.addNewKafkaTopics, lambda);

    if (event.onlyRunDbMigrationAndCreateKafkaTopics) {
      return {
        statusCode: 200,
        body: JSON.stringify({
          message: 'success',
        }),
      };
    }

    // 3. Create new ECS Task Definition for ECS Services with new image
    const taskDefinitionArnMap: TaskDefinitionArnMap = await createNewEcsTaskDefinitions(
      ecs,
      ecr,
      event,
    );

    // 4. Upgrade all ECS Services (comlink, ender, roundtable, socks, vulcan)
    await upgradeEcsServices(ecs, event, taskDefinitionArnMap);
  } catch (error) {
    logger.error({
      at: 'index#handler',
      message: 'Error upgrading services',
      error,
    });
    throw error;
  }

  return {
    statusCode: 200,
    body: JSON.stringify({
      message: 'success',
    }),
  };
}
