import { LOCALIZED_MESSAGES } from './localized-messages';
import {
  Notification,
  LanguageCode,
} from './types';

function replacePlaceholders(template: string, variables: Record<string, string>): string {
  return template.replace(/{(\w+)}/g, (_, key) => variables[key] || `{${key}}`);
}

type NotificationMessage = {
  title: string,
  body: string,
};

export function deriveLocalizedNotificationMessage(
  notification: Notification,
  languageCode: LanguageCode = 'en',
): NotificationMessage {
  const localizationFields = LOCALIZED_MESSAGES[languageCode] || LOCALIZED_MESSAGES.en;

  return {
    title: replacePlaceholders(
      localizationFields[notification.titleKey],
      notification.dynamicValues,
    ),
    body: replacePlaceholders(
      localizationFields[notification.bodyKey],
      notification.dynamicValues,
    ),
  };
}
