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
  FIREBASE_PRIVATE_KEY_BASE64: parseString({ default: 'BASE64_ENCODED_PRIVATE_KEY' }),

  // APP ID for the Google Firebase Messaging project
  FIREBASE_PROJECT_ID: parseString({ default: 'dydx-v4' }),

  // Client email for the Google Firebase Messaging project
  FIREBASE_CLIENT_EMAIL: parseString({ default: 'firebase-adminsdk-f0joo@dydx-v4.iam.gserviceaccount.com' }),
};

export default parseSchema(notificationsConfigSchema);
