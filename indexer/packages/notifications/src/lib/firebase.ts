import { logger } from '@dydxprotocol-indexer/base';
import {
  App,
  cert,
  initializeApp,
  ServiceAccount,
} from 'firebase-admin/app';
import { getMessaging } from 'firebase-admin/messaging';

import config from '../config';

const initializeFirebaseApp = () => {
  const defaultGoogleApplicationCredentials: { [key: string]: string } = {
    project_id: config.FIREBASE_PROJECT_ID,
    private_key: Buffer.from(config.FIREBASE_PRIVATE_KEY_BASE64, 'base64').toString('ascii').replace(/\\n/g, '\n'),
    client_email: config.FIREBASE_CLIENT_EMAIL,
  };

  logger.info({
    at: 'notifications#firebase',
    message: 'Initializing Firebase App',
  });

  const serviceAccount: ServiceAccount = defaultGoogleApplicationCredentials;

  let firebaseApp: App;
  try {
    firebaseApp = initializeApp({
      credential: cert(serviceAccount),
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
};

const firebaseApp = initializeFirebaseApp();
// Initialize Firebase Messaging if the app was initialized successfully
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
