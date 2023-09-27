/**
 * Environment variables required by compliance module.
 */

import {
  parseString,
  parseSchema,
  baseConfigSchema,
  parseInteger,
  parseBoolean,
} from '@dydxprotocol-indexer/base';

export const complianceConfigSchema = {
  ...baseConfigSchema,

  // Required environment variables.
  ELLIPTIC_API_KEY: parseString({ default: 'default_elliptic_api_key' }),
  ELLIPTIC_API_SECRET: parseString({ default: '' }),
  ELLIPTIC_MAX_RETRIES: parseInteger({ default: 3 }),
  ELLIPTIC_RISK_SCORE_THRESHOLD: parseInteger({ default: 10 }),

  // Geo-blocking
  RESTRICTED_COUNTRIES: parseString({
    default: '', // comma de-limited
  }),
  INDEXER_LEVEL_GEOBLOCKING_ENABLED: parseBoolean({
    default: true,
  }),
};

export default parseSchema(complianceConfigSchema);
