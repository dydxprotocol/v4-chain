import { logger } from '@dydxprotocol-indexer/base';
import { TokenTable } from '@dydxprotocol-indexer/postgres';

import {
  BatchResponse,
  MulticastMessage,
  sendMulticast,
} from './lib/firebase';
import { deriveLocalizedNotificationMessage } from './localization';
import { Notification } from './types';

export async function sendFirebaseMessage(
  address: string,
  notification: Notification,
): Promise<string> {
  // Re-add once stats are implemented
  // const start = Date.now();

  const tokens = await getUserRegistrationTokens(address);
  if (tokens.length === 0) {
    return 'user has no registration tokens';
  }

  const { title, body } = deriveLocalizedNotificationMessage(notification);
  const link = notification.deeplink;

  const message: MulticastMessage = {
    tokens,
    notification: {
      title,
      body,
    },
    fcmOptions: {
      analyticsLabel: notification.type.toLowerCase(),
    },
    apns: {
      payload: {
        aps: {
          'mutable-content': 1,
        },
        data: {
          firebase: {
            link,
          },
        },
      },
    },
  };

  try {
    if (sendMulticast) {
      const result: BatchResponse = await sendMulticast(message);
      if (result.failureCount && result.failureCount > 0) {
        logger.info({
          at: 'courier-client#firebase',
          message: `Failed to send Firebase message: ${JSON.stringify(message)}`,
          result,
          address,
          notificationType: notification.type,
        });
      }
      return '';
    } else {
      return 'Failed to send FIrebase message, Firebase is not initialized';
    }
  } catch (error) {
    logger.error({
      at: 'courier-client#firebase',
      message: `Failed to send Firebase message: ${JSON.stringify(message)}`,
      error: error as Error,
      address,
      notificationType: notification.type,
    });
    return `failed to send Firebase message due to ${error}`;
  } finally {
    // stats.timing(`${config.SERVICE_NAME}.send_firebase_message.timing`, Date.now() - start);
  }
}

async function getUserRegistrationTokens(address: string): Promise<string[]> {
  const token = await TokenTable.findAll({ address, limit: 10 }, []);
  return token.map((t) => t.token);
}
