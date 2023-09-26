/**
 * Environment variables required by compliance module.
 */

import {
  parseString,
  parseSchema,
  baseConfigSchema,
  parseInteger,
} from '@dydxprotocol-indexer/base';

export const complianceConfigSchema = {
  ...baseConfigSchema,

  // Required environment variables.
  ELLIPTIC_API_KEY: parseString({ default: 'default_elliptic_api_key' }),
  ELLIPTIC_API_SECRET: parseString({ default: '' }),
  ELLIPTIC_MAX_RETRIES: parseInteger({ default: 3 }),
  ELLIPTIC_RISK_SCORE_THRESHOLD: parseInteger({ default: 10 }),
};

export default parseSchema(complianceConfigSchema);
