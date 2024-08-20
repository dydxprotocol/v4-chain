import {
  LocalizationKey,
  Notification,
  NotificationMesage,
} from './types';

function replacePlaceholders(template: string, variables: Record<string, string>): string {
  return template.replace(/{(\w+)}/g, (_, key) => variables[key] || `{${key}}`);
}

export function deriveLocalizedNotificationMessage(notification: Notification): NotificationMesage {
  const tempLocalizationFields = {
    [LocalizationKey.DEPOSIT_SUCCESS_TITLE]: 'Deposit Successful',
    [LocalizationKey.DEPOSIT_SUCCESS_BODY]: 'You have successfully deposited {AMOUNT} {MARKET} to your dYdX account.',
    [LocalizationKey.ORDER_FILLED_TITLE]: 'Order Filled',
    // eslint-disable-next-line no-template-curly-in-string
    [LocalizationKey.ORDER_FILLED_BODY]: 'Your order for {AMOUNT} {MARKET} was filled at ${AVERAGE_PRICE}',
    // eslint-disable-next-line no-template-curly-in-string
    [LocalizationKey.ORDER_TRIGGERED_BODY]: 'Your order for {AMOUNT} {MARKET} was triggered at ${PRICE}',
    [LocalizationKey.ORDER_TRIGGERED_TITLE]: '{MARKET} Order Triggered',
  };

  return {
    title: replacePlaceholders(
      tempLocalizationFields[notification.titleKey],
      notification.dynamicValues),
    body: replacePlaceholders(
      tempLocalizationFields[notification.bodyKey],
      notification.dynamicValues,
    ),
  };
}
