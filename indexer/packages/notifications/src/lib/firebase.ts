import {
  cert,
  initializeApp,
  ServiceAccount,
} from 'firebase-admin/app';
import { getMessaging } from 'firebase-admin/messaging';

import config from '../config';

const defaultGoogleApplicationCredentials: { [key: string]: string } = {
  project_id: config.FIREBASE_PROJECT_ID,
  private_key: Buffer.from(config.FIREBASE_PRIVATE_KEY, 'base64').toString('utf8'),
  client_email: config.FIREBASE_CLIENT_EMAIL,
};
const serviceAccount: ServiceAccount = JSON.parse(
  JSON.stringify(defaultGoogleApplicationCredentials),
);

const firebaseApp = initializeApp({
  credential: cert(serviceAccount),
});
const firebaseMessaging = getMessaging(firebaseApp);

export const sendMulticast = firebaseMessaging.sendMulticast.bind(firebaseMessaging);
export { BatchResponse, getMessaging, MulticastMessage } from 'firebase-admin/messaging';
