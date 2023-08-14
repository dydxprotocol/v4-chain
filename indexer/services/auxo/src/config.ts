/**
 * Environment variables required by Auxo.
 */

import {
  parseSchema,
  baseConfigSchema,
  parseNumber,
} from '@dydxprotocol-indexer/base';

export const configSchema = {
  ...baseConfigSchema,
  // Max amount of time we want to wait for a task definition to be created
  MAX_TASK_DEFINITION_WAIT_TIME_MS: parseNumber({
    default: 60_000, // 60s
  }),
  SLEEP_TIME_MS: parseNumber({
    default: 5_000, // 5s
  }),
};

export default parseSchema(configSchema);
