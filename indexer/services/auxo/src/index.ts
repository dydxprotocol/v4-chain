import {
  DescribeImagesCommand,
  DescribeImagesCommandOutput,
  ECRClient,
  ImageDetail,
} from '@aws-sdk/client-ecr';
import {
  ContainerDefinition,
  DescribeServicesCommand,
  DescribeServicesCommandOutput,
  DescribeTaskDefinitionCommand,
  DescribeTaskDefinitionCommandOutput,
  ECSClient,
  RegisterTaskDefinitionCommand,
  RegisterTaskDefinitionCommandOutput,
  Service,
  TaskDefinition,
  UpdateServiceCommand,
  UpdateServiceCommandOutput,
} from '@aws-sdk/client-ecs';
import {
  InvokeCommand,
  InvokeCommandOutput,
  LambdaClient,
  UpdateFunctionCodeCommand,
  UpdateFunctionCodeCommandOutput,
} from '@aws-sdk/client-lambda';
import { logger, startBugsnag } from '@dydxprotocol-indexer/base';
import {
  APIGatewayEvent,
  APIGatewayProxyResult,
  Context,
} from 'aws-lambda';
import _ from 'lodash';

import config from './config';
import {
  BAZOOKA_DB_MIGRATION_PAYLOAD,
  BAZOOKA_LAMBDA_FUNCTION_NAME,
  ECS_SERVICE_NAMES,
  SERVICE_NAME_SUFFIX,
} from './constants';
import {
  AuxoEventJson,
  EcsServiceNames,
  TaskDefinitionArnMap,
} from './types';

/**
 * Upgrades all services and run migrations
 * 1. Upgrade Bazooka
 * 2. Run db migration in Bazooka
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

    // 2. Run db migration in Bazooka
    await runDbMigration(lambda);

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

async function upgradeBazooka(
  lambda: LambdaClient,
  ecr: ECRClient,
  event: APIGatewayEvent & AuxoEventJson,
): Promise<void> {
  const imageDetail: ImageDetail = await getImageDetail(ecr, `${event.prefix}-indexer-bazooka`, event);
  const imageUri = `${imageDetail.registryId}.dkr.ecr.${event.region}.amazonaws.com/${imageDetail.repositoryName}@${imageDetail.imageDigest}`;
  logger.info({
    at: 'index#upgradeBazooka',
    message: `Upgrading bazooka to ${imageUri}`,
  });

  // Update Lambda function with the new image
  const response: UpdateFunctionCodeCommandOutput = await lambda.send(
    new UpdateFunctionCodeCommand({
      FunctionName: BAZOOKA_LAMBDA_FUNCTION_NAME,
      ImageUri: imageUri,
    }),
  );
  logger.info({
    at: 'index#upgradeBazooka',
    message: 'Successfully upgraded bazooka',
    response,
  });
}

async function getImageDetail(
  ecr: ECRClient,
  repositoryName: string,
  event: APIGatewayEvent & AuxoEventJson,
): Promise<ImageDetail> {
  logger.info({
    at: 'index#upgradeBazooka',
    message: 'Getting ecr images',
    repositoryName,
    event,
  });
  const images: DescribeImagesCommandOutput = await ecr.send(new DescribeImagesCommand({
    repositoryName,
    imageIds: [
      {
        imageTag: event.upgrade_tag,
      },
    ],
  }));
  logger.info({
    at: 'index#upgradeBazooka',
    message: 'Successfully got ecr images',
    images,
    repositoryName,
    event,
  });

  if (!images.imageDetails || images.imageDetails.length === 0) {
    logger.error({
      at: 'index#upgradeBazooka',
      message: 'Unable to find ecr image',
      imageTag: event.upgrade_tag,
      repositoryName,
      event,
    });
    throw new Error('Unable to find ecr image');
  }
  return images.imageDetails[0];

}

async function runDbMigration(
  lambda: ECRClient,
): Promise<void> {
  logger.info({
    at: 'index#runDbMigration',
    message: 'Running db migration',
  });
  const response: InvokeCommandOutput = await lambda.send(new InvokeCommand({
    FunctionName: BAZOOKA_LAMBDA_FUNCTION_NAME,
    Payload: BAZOOKA_DB_MIGRATION_PAYLOAD,
    // RequestResponse means that the lambda is synchronously invoked
    InvocationType: 'RequestResponse',
  }));
  logger.info({
    at: 'index#runDbMigration',
    message: 'Successfully ran db migration',
    response,
  });
}

async function createNewEcsTaskDefinitions(
  ecs: ECSClient,
  ecr: ECRClient,
  event: APIGatewayEvent & AuxoEventJson,
): Promise<TaskDefinitionArnMap> {
  logger.info({
    at: 'index#createNewEcsTaskDefinitions',
    message: 'Creating new ECS Task Definitions',
  });
  const taskDefinitionArns: string[] = await Promise.all(_.map(
    ECS_SERVICE_NAMES,
    (serviceName: EcsServiceNames) => createNewEcsTaskDefinition(ecs, ecr, event, serviceName),
  ));
  logger.info({
    at: 'index#createNewEcsTaskDefinitions',
    message: 'Created new ECS Task Definition',
  });
  return _.zipObject(ECS_SERVICE_NAMES, taskDefinitionArns);
}

/**
 * @returns The revision number of the new task definition
 */
async function createNewEcsTaskDefinition(
  ecs: ECSClient,
  ecr: ECRClient,
  event: APIGatewayEvent & AuxoEventJson,
  serviceName: EcsServiceNames,
): Promise<string> {
  // Check that the ECR image exists, will throw error here if it does not
  await getImageDetail(ecr, `${event.prefix}-indexer-${serviceName}`, event);

  const taskDefinitionName = `${event.prefix}-indexer-${event.regionAbbrev}-${serviceName}-task`;
  logger.info({
    at: 'index#createNewEcsTaskDefinition',
    message: 'Get existing ECS Task Definition',
    taskDefinitionName,
  });
  const describeResult: DescribeTaskDefinitionCommandOutput = await ecs.send(
    new DescribeTaskDefinitionCommand({
      taskDefinition: taskDefinitionName,
    }),
  );
  logger.info({
    at: 'index#createNewEcsTaskDefinition',
    message: 'Got existing ECS Task Definition',
    taskDefinitionName,
    describeResult,
  });

  if (describeResult.taskDefinition === undefined) {
    logger.error({
      at: 'index#createNewEcsTaskDefinition',
      message: 'Unable to find existing ECS Task Definition',
      taskDefinitionName,
    });
    throw new Error('Unable to find existing ECS Task Definition');
  }

  // All ECS Task Definitions should have two container definitions, the service container
  // , and the datadog agent
  const taskDefinition: TaskDefinition = describeResult.taskDefinition;
  const serviceContainerDefinitionIndex: number = getServiceContainerDefinitionIndex(
    taskDefinition,
  );

  const serviceContainerDefinition:
  ContainerDefinition = taskDefinition.containerDefinitions![serviceContainerDefinitionIndex];
  if (serviceContainerDefinition.image === undefined) {
    logger.error({
      at: 'index#createNewEcsTaskDefinition',
      message: 'No image found in the container definition',
      taskDefinitionName,
    });
    throw new Error('No image found in the container definition');
  }
  const originalImage: string = serviceContainerDefinition.image;
  const updatedContainerDefinitions: ContainerDefinition[] = _.cloneDeep(
    taskDefinition.containerDefinitions!,
  );
  const newImage: string = `${_.split(originalImage, ':')[0]}:${event.upgrade_tag}`;
  updatedContainerDefinitions[serviceContainerDefinitionIndex].image = newImage;

  logger.info({
    at: 'index#createNewEcsTaskDefinition',
    message: 'Registering new task definition',
    taskDefinitionName,
  });
  const registerResult:
  RegisterTaskDefinitionCommandOutput = await ecs.send(new RegisterTaskDefinitionCommand({
    family: taskDefinition.family,
    taskRoleArn: taskDefinition.taskRoleArn,
    executionRoleArn: taskDefinition.executionRoleArn,
    networkMode: taskDefinition.networkMode,
    containerDefinitions: updatedContainerDefinitions,
    volumes: taskDefinition.volumes,
    placementConstraints: taskDefinition.placementConstraints,
    requiresCompatibilities: taskDefinition.requiresCompatibilities,
    cpu: taskDefinition.cpu,
    memory: taskDefinition.memory,
    ipcMode: taskDefinition.ipcMode,
    proxyConfiguration: taskDefinition.proxyConfiguration,
    inferenceAccelerators: taskDefinition.inferenceAccelerators,
    runtimePlatform: taskDefinition.runtimePlatform,
  }));

  if (registerResult.taskDefinition === undefined ||
    registerResult.taskDefinition.taskDefinitionArn === undefined
  ) {
    logger.error({
      at: 'index#createNewEcsTaskDefinition',
      message: 'Failed to register new task definition',
    });
    throw new Error('Failed to register new task definition');
  }

  await waitForTaskDefinitionToRegister(ecs, registerResult);
  return registerResult.taskDefinition.taskDefinitionArn;
}

function getServiceContainerDefinitionIndex(
  taskDefinition: TaskDefinition,
): number {
  const containerDefinitions:
  ContainerDefinition[] | undefined = taskDefinition.containerDefinitions;
  if (containerDefinitions === undefined || containerDefinitions.length === 0) {
    logger.error({
      at: 'index#getServiceTaskDefinition',
      message: 'No container definitions found in the task definition',
      taskDefinition,
    });
    throw new Error('No container definitions found in the task definition');
  }

  const index: number = containerDefinitions.findIndex(
    (containerDefinition: ContainerDefinition) => {
      return _.endsWith(containerDefinition.name, SERVICE_NAME_SUFFIX);
    },
  );
  if (index >= 0) {
    return index;
  }

  logger.error({
    at: 'index#getServiceTaskDefinition',
    message: 'No service container definition found in the task definition',
    containerDefinitions,
  });
  throw new Error('No service container definition found in the task definition');
}

/**
 * Registering a task definition is asynchronous, and this step ensures that the task definition
 * is usable in the ECS service before we attempt to update the ECS service.
 */
async function waitForTaskDefinitionToRegister(
  ecs: ECSClient,
  registerResult: RegisterTaskDefinitionCommandOutput,
): Promise<void> {
  const taskDefinition:
  string = `${registerResult.taskDefinition!.family}:${registerResult.taskDefinition!.revision}`;
  for (let i = 0; i <= config.MAX_TASK_DEFINITION_WAIT_TIME_MS; i += config.SLEEP_TIME_MS) {
    const describeResult: DescribeTaskDefinitionCommandOutput = await ecs.send(
      new DescribeTaskDefinitionCommand({
        taskDefinition,
      }),
    );

    if (describeResult.taskDefinition !== undefined) {
      logger.info({
        at: 'index#waitForTaskDefinitionToRegister',
        message: 'Task definition registered',
        taskDefinition,
        describeResult,
      });
      return;
    }
    logger.info({
      at: 'index#waitForTaskDefinitionToRegister',
      message: `Task definition is undefined, sleeping ${config.SLEEP_TIME_MS}ms`,
    });

    await sleep(config.SLEEP_TIME_MS);
  }
  logger.error({
    at: 'index#waitForTaskDefinitionToRegister',
    message: 'Timed out waiting for task definition to register',
    taskDefinition,
  });
  throw new Error('Timed out waiting for task definition to register');
}

async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function upgradeEcsServices(
  ecs: ECSClient,
  event: APIGatewayEvent & AuxoEventJson,
  taskDefinitionArnMap: TaskDefinitionArnMap,
): Promise<void> {
  logger.info({
    at: 'index#upgradeEcsServices',
    message: 'Describe Services',
  });
  const ecsPrefix: string = `${event.prefix}-indexer-${event.regionAbbrev}`;
  const response: DescribeServicesCommandOutput = await ecs.send(new DescribeServicesCommand({
    cluster: `${ecsPrefix}-cluster`,
    services: _.map(
      ECS_SERVICE_NAMES,
      (serviceName: EcsServiceNames) => {
        return `${ecsPrefix}-${serviceName}`;
      },
    ),
    include: [],
  }));

  logger.info({
    at: 'index#upgradeEcsServices',
    message: 'Described Services',
    response,
  });

  if (response.services === undefined) {
    logger.error({
      at: 'index#upgradeEcsServices',
      message: 'No services found',
    });
    throw new Error('No services found');
  } else if (response.services.length !== ECS_SERVICE_NAMES.length) {
    logger.error({
      at: 'index#upgradeEcsServices',
      message: 'Not all services found',
      numServicesFound: response.services.length,
      services: response.services,
      numServicesExpected: ECS_SERVICE_NAMES.length,
    });
    throw new Error('Not all services found');
  }

  logger.info({
    at: 'index#upgradeEcsServices',
    message: 'Upgrading ECS Services',
  });
  const services: Service[] = response.services;
  await Promise.all(_.map(
    ECS_SERVICE_NAMES,
    (serviceName: string, index: number) => upgradeEcsService(
      ecs,
      services[index],
      taskDefinitionArnMap[serviceName],
    ),
  ));

  logger.info({
    at: 'index#upgradeEcsServices',
    message: 'Upgraded ECS Services',
  });
}

async function upgradeEcsService(
  ecs: ECSClient,
  service: Service,
  taskDefinitionArn: string,
): Promise<void> {
  logger.info({
    at: 'index#upgradeEcsService',
    message: 'Upgrading ECS Service',
    service,
    taskDefinitionArn,
  });
  const response: UpdateServiceCommandOutput = await ecs.send(new UpdateServiceCommand({
    cluster: service.clusterArn,
    service: service.serviceName,
    desiredCount: service.desiredCount,
    taskDefinition: taskDefinitionArn,
    capacityProviderStrategy: service.capacityProviderStrategy,
    deploymentConfiguration: service.deploymentConfiguration,
    networkConfiguration: service.networkConfiguration,
    placementConstraints: service.placementConstraints,
    placementStrategy: service.placementStrategy,
    platformVersion: service.platformVersion,
    healthCheckGracePeriodSeconds: service.healthCheckGracePeriodSeconds,
    enableExecuteCommand: service.enableExecuteCommand,
    enableECSManagedTags: service.enableECSManagedTags,
    loadBalancers: service.loadBalancers,
    propagateTags: service.propagateTags,
    serviceRegistries: service.serviceRegistries,
  }));

  logger.info({
    at: 'index#upgradeEcsService',
    message: 'Upgraded ECS Service',
    serviceName: service.serviceName,
    taskDefinitionArn,
    response,
  });
}
