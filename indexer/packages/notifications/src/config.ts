/**
 * Environment variables required for Notifications module.
 */

import {
  parseString,
  parseSchema,
  baseConfigSchema,
} from '@dydxprotocol-indexer/base';

export const notificationsConfigSchema = {
  ...baseConfigSchema,

  // Private Key for the Google Firebase Messaging project
  FIREBASE_PRIVATE_KEY: parseString({ default: '' }),

  // APP ID for the Google Firebase Messaging project
  FIREBASE_APP_ID: parseString({ default: '' }),
};

export default parseSchema(notificationsConfigSchema);
