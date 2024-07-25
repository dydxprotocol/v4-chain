import {
  cert,
  initializeApp,
  ServiceAccount,
} from 'firebase-admin/app';
import { getMessaging } from 'firebase-admin/messaging';

// TODO: Move out to a config file or env vars
const GOOGLE_APPLICATION_CREDENTIALS = '';
const FIREBASE_DATABASE_URL = '';

const serviceAccount: ServiceAccount = JSON.parse(
  Buffer.from(GOOGLE_APPLICATION_CREDENTIALS, 'base64').toString('ascii'),
);
const firebaseApp = initializeApp({
  credential: cert(serviceAccount),
  databaseURL: FIREBASE_DATABASE_URL,
});
const firebaseMessaging = getMessaging(firebaseApp);

export const sendMulticast = firebaseMessaging.sendMulticast.bind(firebaseMessaging);
export { BatchResponse, getMessaging, MulticastMessage } from 'firebase-admin/messaging';
