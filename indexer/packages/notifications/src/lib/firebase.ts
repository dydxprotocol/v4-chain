import { logger } from '@dydxprotocol-indexer/base';
import {
  cert,
  initializeApp,
  ServiceAccount,
} from 'firebase-admin/app';
import { getMessaging } from 'firebase-admin/messaging';

import config from '../config';

const initializeFirebaseApp = () => {
  const defaultGoogleApplicationCredentials: { [key: string]: string } = {
    project_id: config.FIREBASE_PROJECT_ID,
    private_key: config.FIREBASE_PRIVATE_KEY,
    client_email: config.FIREBASE_CLIENT_EMAIL,
  };

  logger.info({
    at: 'notifications#firebase',
    message: 'Initializing Firebase App with credentials',
    credentials: defaultGoogleApplicationCredentials,
  });

  const serviceAccount: ServiceAccount = defaultGoogleApplicationCredentials;

  const firebaseApp = initializeApp({
    credential: cert(serviceAccount),
  });

  logger.info({
    at: 'notifications#firebase',
    message: 'Firebase App initialized successfully',
  });

  return firebaseApp;
};

const firebaseApp = initializeFirebaseApp();
const firebaseMessaging = getMessaging(firebaseApp);

export const sendMulticast = firebaseMessaging.sendMulticast.bind(firebaseMessaging);
export { BatchResponse, getMessaging, MulticastMessage } from 'firebase-admin/messaging';
