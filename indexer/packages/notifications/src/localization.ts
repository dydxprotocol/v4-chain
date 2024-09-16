import { LOCALIZED_MESSAGES } from './localized-messages';
import {
  Notification,
  LanguageCode,
} from './types';

function replacePlaceholders(template: string, variables: Record<string, string>): string {
  // The template string contains placeholders in the format "{KEY}".
  // For example: 'Your order for {AMOUNT} {MARKET} was filled at ${AVERAGE_PRICE}'.
  // This function replaces these placeholders with corresponding values
  // from the "variables" object.
  // If the key inside "{}" exists in the "variables" object, it is replaced
  // with the matching value.
  // If the key does not exist, the placeholder remains unchanged in the resulting string.
  // The .replace method uses a regular expression to find all words inside "{}"
  // and replaces them as described.
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
