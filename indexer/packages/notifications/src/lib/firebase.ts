import { logger } from '@dydxprotocol-indexer/base';
import {
  App,
  cert,
  initializeApp,
} from 'firebase-admin/app';
import { getMessaging } from 'firebase-admin/messaging';

import config from '../config';

// Helper function to initialize Firebase App object that is used to send notifications
function initializeFirebaseApp(): App | undefined {
  // Create credentials object from config variables
  const defaultGoogleApplicationCredentials: { [key: string]: string } = {
    project_id: config.FIREBASE_PROJECT_ID,
    private_key: Buffer.from(config.FIREBASE_PRIVATE_KEY_BASE64, 'base64').toString('ascii').replace(/\\n/g, '\n'),
    client_email: config.FIREBASE_CLIENT_EMAIL,
  };

  logger.info({
    at: 'notifications#firebase',
    message: 'Initializing Firebase App',
  });

  let firebaseApp: App;
  try {
    firebaseApp = initializeApp({
      credential: cert(defaultGoogleApplicationCredentials),
    });
  } catch (error) {
    logger.error({
      at: 'notifications#firebase',
      message: 'Failed to initialize Firebase App',
      error,
    });
    return undefined;
  }

  logger.info({
    at: 'notifications#firebase',
    message: 'Firebase App initialized successfully',
  });

  return firebaseApp;
}

const firebaseApp = initializeFirebaseApp();

// Initialize Firebase Messaging if the app was initialized successfully
// This can fail if the credentials passed to the firebaseApp are invalid
let firebaseMessaging = null;
if (firebaseApp) {
  try {
    firebaseMessaging = getMessaging(firebaseApp);
    logger.info({
      at: 'notifications#firebase',
      message: 'Firebase Messaging initialized successfully',
    });
  } catch (error) {
    logger.error({
      at: 'notifications#firebase',
      message: 'Firebase Messaging failed to initialize',
    });
  }
}

export const sendMulticast = firebaseMessaging
  ? firebaseMessaging.sendMulticast.bind(firebaseMessaging)
  : () => {
    logger.error({
      at: 'notifications#firebase',
      message: 'Firebase Messaging is not initialized, sendMulticast is a no-op',
    });
    return Promise.resolve(null);
  };
export { BatchResponse, getMessaging, MulticastMessage } from 'firebase-admin/messaging';
