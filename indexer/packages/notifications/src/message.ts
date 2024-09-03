import { logger } from '@dydxprotocol-indexer/base';

import {
  MulticastMessage,
  sendMulticast,
} from './lib/firebase';
import { deriveLocalizedNotificationMessage } from './localization';
import { LanguageCode, Notification } from './types';

export async function sendFirebaseMessage(
  tokens: {token: string, language: string}[],
  notification: Notification,
): Promise<void> {
  // Re-add once stats are implemented
  // const start = Date.now();

  if (tokens.length === 0) {
    logger.warning({
      at: 'notifications#firebase',
      message: 'Attempted to send Firebase message to user with no registration tokens',
      tokens,
      notificationType: notification.type,
    });
    return;
  }

  // Each set of tokens for a users should have the same language
  const language = tokens[0].language;
  const { title, body } = deriveLocalizedNotificationMessage(
    notification,
    language as LanguageCode,
  );
  const link = notification.deeplink;

  const message: MulticastMessage = {
    tokens: tokens.map((token) => token.token),
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
    const result = await sendMulticast(message);
    if (result?.failureCount && result?.failureCount > 0) {
      throw new Error('Failed to send Firebase message');
    }
  } catch (error) {
    logger.error({
      at: 'notifications#firebase',
      message: 'Failed to send Firebase message',
      error: error as Error,
      notificationType: notification.type,
    });
  } finally {
    // stats.timing(`${config.SERVICE_NAME}.send_firebase_message.timing`, Date.now() - start);
    logger.info({
      at: 'notifications#firebase',
      message: 'Firebase message sent successfully',
      notificationType: notification.type,
    });
  }
}
