import {
  LocalizationKey,
  Notification,
  NotificationMesage,
  NotificationType,
} from './types';

function replacePlaceholders(template: string, variables: Record<string, string>): string {
  return template.replace(/{(\w+)}/g, (_, key) => variables[key] || `{${key}}`);
}

export function deriveLocalizedNotificationMessage(notification: Notification): NotificationMesage {
  const tempLocalizationFields = {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: 'Deposit successful',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: 'You have successfully deposited {AMOUNT} {MARKET} to your dYdX account.',
    [LocalizationKey.ORDER_FILLED_BODY]: 'Order Filled successful',
    [LocalizationKey.ORDER_FILLED_TITLE]: 'Filled {MARKET} order at {AVERAGE_PRICE}.',
  };

  switch (notification.type) {
    case NotificationType.DEPOSIT_SUCCESS:
      return {
        title: replacePlaceholders(
          tempLocalizationFields[notification.titleKey],
          notification.dynamicValues),
        body: replacePlaceholders(
          tempLocalizationFields[notification.bodyKey],
          notification.dynamicValues,
        ),
      };
    case NotificationType.ORDER_FILLED:
      return {
        title: replacePlaceholders(
          tempLocalizationFields[notification.titleKey],
          notification.dynamicValues),
        body: replacePlaceholders(
          tempLocalizationFields[notification.bodyKey],
          notification.dynamicValues,
        ),
      };
    default:
      throw new Error('Unknown notification type');
  }
}
