import { TextEncoder } from 'util';

import { EcsServiceNames } from './types';

export const BAZOOKA_LAMBDA_FUNCTION_NAME: string = 'bazooka_lambda_function';

export const BAZOOKA_DB_MIGRATION_PAYLOAD: Uint8Array = new TextEncoder().encode(
  JSON.stringify({
    migrate: true,
  }),
);

export const BAZOOKA_DB_MIGRATION_AND_CREATE_KAFKA_PAYLOAD: Uint8Array = new TextEncoder().encode(
  JSON.stringify({
    migrate: true,
    create_kafka_topics: true,
  }),
);

export const ECS_SERVICE_NAMES: EcsServiceNames[] = [
  EcsServiceNames.COMLINK,
  EcsServiceNames.ENDER,
  EcsServiceNames.ROUNDTABLE,
  EcsServiceNames.SOCKS,
  EcsServiceNames.VULCAN,
];

export const ECS_DB_WRITER_SERVICE_NAMES: EcsServiceNames[] = [
  EcsServiceNames.ENDER,
  EcsServiceNames.ROUNDTABLE,
];

export const SERVICE_NAME_SUFFIX: string = 'service-container';
