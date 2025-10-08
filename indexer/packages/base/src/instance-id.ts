import { v4 as uuidv4 } from 'uuid';

import { axiosRequest } from './axios';
import config from './config';
import logger from './logger';

let INSTANCE_ID: string = '';

export function getInstanceId(): string {
  return INSTANCE_ID;
}

export async function setInstanceId(instanceId? : string): Promise<void> {
  if (instanceId !== undefined) {
    INSTANCE_ID = instanceId;
    return;
  }

  if (INSTANCE_ID !== '') {
    return;
  }
  if (config.ECS_CONTAINER_METADATA_URI_V4 !== '' &&
      (
        config.isProduction() || config.isStaging()
      )
  ) {
    // https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html
    const taskUrl = `${config.ECS_CONTAINER_METADATA_URI_V4}/task`;
    try {
      const response = await axiosRequest({
        method: 'GET',
        url: taskUrl,
      }) as { TaskARN: string };
      INSTANCE_ID = response.TaskARN;
    } catch (error) {
      logger.error({
        at: 'instance-id#setInstanceId',
        message: 'Failed to retrieve task arn from metadata endpoint. Falling back to uuid.',
        error,
        taskUrl,
      });
      INSTANCE_ID = uuidv4();
    }
  } else {
    INSTANCE_ID = uuidv4();

  }
}

// Exported for tests
export function resetForTests(): void {
  if (!config.isTest()) {
    throw new Error(`resetForTests() cannot be called for env: ${config.NODE_ENV}`);
  }
  INSTANCE_ID = '';
}
