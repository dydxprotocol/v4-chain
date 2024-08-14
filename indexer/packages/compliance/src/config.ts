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

  // will be matched with enums in helpers/compliance/compliance-clients and default to
  // ELLIPTIC if unset or an invalid value is set
  COMPLIANCE_DATA_CLIENT: parseString({ default: 'ELLIPTIC' }),

  // Block-list provider environment variables.
  BLOCKED_ADDRESSES: parseString({
    default: '', // comma de-limited
  }),

  // Whitelisted list of dydx addresses
  WHITELISTED_ADDRESSES: parseString({
    default: '', // comma de-limited
  }),

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
