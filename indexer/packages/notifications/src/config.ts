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
  FIREBASE_PROJECT_ID: parseString({ default: '' }),

  // Client email for the Google Firebase Messaging project
  FIREBASE_CLIENT_EMAIL: parseString({ default: '' }),
};

export default parseSchema(notificationsConfigSchema);
