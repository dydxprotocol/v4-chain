/**
 * Environment variables required by all services using base.
 */

import {
  parseInteger,
  parseString,
  parseSchema,
  parseBoolean,
} from './config-util';

export const baseConfigSchema = {
  // Required environment variables.
  SERVICE_NAME: parseString({ default: '' }),

  // Optional environment variables.
  NODE_ENV: parseString({ default: null }),
  ENABLE_LOGS_IN_TEST: parseBoolean({ default: false }),
  STATSD_HOST: parseString({ default: 'localhost' }),
  STATSD_PORT: parseInteger({ default: 8125 }),
  LOG_LEVEL: parseString({ default: 'debug' }),
  ECS_CONTAINER_METADATA_URI_V4: parseString({ default: '' }),
  AWS_REGION: parseString({ default: '' }),
};

export default parseSchema(baseConfigSchema);
