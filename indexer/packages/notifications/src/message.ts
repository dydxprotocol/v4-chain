import { logger, stats } from '@dydxprotocol-indexer/base';

import config from './config';
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
  const start = Date.now();

  // Each set of tokens for a users should have the same language
  const language = tokens[0].language;
  const { title, body } = deriveLocalizedNotificationMessage(
    notification,
    language as LanguageCode,
  );

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
          },
        },
      },
    },
  };

  try {
    const result = await sendMulticast(message);
    if (!result || result?.failureCount > 0) {
      const errorMessages = result?.responses
        .map((response) => response.error?.message)
        .filter(Boolean); // Remove any undefined values

      throw new Error(`Failed to send Firebase message: ${errorMessages?.join(', ') || 'Unknown error'}`);
    }

    logger.info({
      at: 'notifications#firebase',
      message: 'Firebase message sent successfully',
      notificationType: notification.type,
    });
  } catch (error) {
    logger.error({
      at: 'notifications#firebase',
      message: error.message,
      error: error as Error,
      notificationType: notification.type,
    });
  } finally {
    stats.timing(`${config.SERVICE_NAME}.send_firebase_message.timing`, Date.now() - start);
  }
}
