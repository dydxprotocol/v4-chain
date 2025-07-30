/**
 * Environment variables required by all services using base.
 */

import {
  parseInteger,
  parseString,
  parseSchema,
  parseBoolean,
} from './config-util';
import { NodeEnv } from './types';

// Use a string that passes validation when creating the Bugsnag object.
const DEFAULT_BUGSNAG_KEY = '00000000000000000000000000000000';

export const baseConfigSecrets: (keyof typeof baseConfigSchema)[] = [
  'BUGSNAG_KEY',
];

export const baseConfigSchema = {
  // Required environment variables.
  BUGSNAG_KEY: parseString({
    default: DEFAULT_BUGSNAG_KEY,
    requireInEnv: [NodeEnv.PRODUCTION, NodeEnv.STAGING],
  }),
  BUGSNAG_RELEASE_STAGE: parseString({
    default: null,
    requireInEnv: [NodeEnv.PRODUCTION, NodeEnv.STAGING],
  }),
  SEND_BUGSNAG_ERRORS: parseBoolean({
    default: true,
  }),
  SERVICE_NAME: parseString({ default: '' }),

  TURNKEY_API_BASE_URL: parseString({ default: 'https://api.turnkey.com' }),
  // API keys for root user on parent org to use to create suborgs.
  TURNKEY_API_PRIVATE_KEY: parseString({ default: '' }),
  TURNKEY_API_PUBLIC_KEY: parseString({ default: '' }),
  // API keys for senders to use to start bridging.
  TURNKEY_API_SENDER_PRIVATE_KEY: parseString({ default: '' }),
  TURNKEY_API_SENDER_PUBLIC_KEY: parseString({ default: '' }),
  TURNKEY_ORGANIZATION_ID: parseString({ default: '' }),
  // Alchemy API key for skip bridge.
  ALCHEMY_KEY: parseString({ default: '' }),
  // ZeroDev RPC for skip bridge.
  ZERODEV_RPC: parseString({ default: '' }),

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
