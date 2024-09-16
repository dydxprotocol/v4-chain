// Types of notifications that can be sent
export enum NotificationType {
  DEPOSIT_SUCCESS = 'DEPOSIT_SUCCESS',
  ORDER_FILLED = 'ORDER_FILLED',
  ORDER_TRIGGERED = 'ORDER_TRIGGERED',
}

// Keys for the dynamic values that are used in the notification messages
// Each key corresponds to a placeholder in the localizable strings for each notification
export enum NotificationDynamicFieldKey {
  AMOUNT = 'AMOUNT',
  AVERAGE_PRICE = 'AVERAGE_PRICE',
  PRICE = 'PRICE',
  FILLED_AMOUNT = 'FILLED_AMOUNT',
  MARKET = 'MARKET',
  SIDE = 'SIDE',
}

// Keys for the strings that are contained in the localzation file
// for each notification body and title
export enum LocalizationBodyKey {
  DEPOSIT_SUCCESS_BODY = 'DEPOSIT_SUCCESS_BODY',
  ORDER_FILLED_BODY = 'ORDER_FILLED_BODY',
  ORDER_TRIGGERED_BODY = 'ORDER_TRIGGERED_BODY',
}

export enum LocalizationTitleKey {
  DEPOSIT_SUCCESS_TITLE = 'DEPOSIT_SUCCESS_TITLE',
  ORDER_FILLED_TITLE = 'ORDER_FILLED_TITLE',
  ORDER_TRIGGERED_TITLE = 'ORDER_TRIGGERED_TITLE',
}

export type LocalizationKey = LocalizationBodyKey | LocalizationTitleKey;

// Topics for each notification
// Topics are used to send notifications to specific topics in Firebase
export enum Topic {
  TRADING = 'trading',
  PRICE_ALERTS = 'price_alerts',
}

export type LanguageCode = 'en' | 'es' | 'fr' | 'de' | 'it' | 'ja' | 'ko' | 'zh';
export function isValidLanguageCode(code: string): code is LanguageCode {
  return ['en', 'es', 'fr', 'de', 'it', 'ja', 'ko', 'zh'].includes(code);
}

interface BaseNotification <T extends Record<string, string>> {
  type: NotificationType,
  titleKey: LocalizationTitleKey,
  bodyKey: LocalizationBodyKey,
  topic: Topic,
  dynamicValues: T,
}

interface DepositSuccessNotification extends BaseNotification<{
  [NotificationDynamicFieldKey.AMOUNT]: string,
  [NotificationDynamicFieldKey.MARKET]: string,
}> {
  type: NotificationType.DEPOSIT_SUCCESS,
  titleKey: LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE,
  bodyKey: LocalizationBodyKey.DEPOSIT_SUCCESS_BODY,
  topic: Topic.TRADING,
  dynamicValues: {
    [NotificationDynamicFieldKey.AMOUNT]: string,
    [NotificationDynamicFieldKey.MARKET]: string,
  },
}

interface OrderFilledNotification extends BaseNotification <{
  [NotificationDynamicFieldKey.MARKET]: string,
  [NotificationDynamicFieldKey.AVERAGE_PRICE]: string,
}>{
  type: NotificationType.ORDER_FILLED,
  titleKey: LocalizationTitleKey.ORDER_FILLED_TITLE,
  bodyKey: LocalizationBodyKey.ORDER_FILLED_BODY,
  topic: Topic.TRADING,
  dynamicValues: {
    [NotificationDynamicFieldKey.MARKET]: string,
    [NotificationDynamicFieldKey.AMOUNT]: string,
    [NotificationDynamicFieldKey.AVERAGE_PRICE]: string,
  },
}

interface OrderTriggeredNotification extends BaseNotification <{
  [NotificationDynamicFieldKey.MARKET]: string,
  [NotificationDynamicFieldKey.PRICE]: string,
}>{
  type: NotificationType.ORDER_TRIGGERED,
  titleKey: LocalizationTitleKey.ORDER_TRIGGERED_TITLE,
  bodyKey: LocalizationBodyKey.ORDER_TRIGGERED_BODY,
  topic: Topic.TRADING,
  dynamicValues: {
    [NotificationDynamicFieldKey.MARKET]: string,
    [NotificationDynamicFieldKey.AMOUNT]: string,
    [NotificationDynamicFieldKey.PRICE]: string,
  },
}

export type Notification =
DepositSuccessNotification |
OrderFilledNotification |
OrderTriggeredNotification;

// Factory function to create notifications.
//
// dynamicValues is a conditional type that changes based on the notification type:
// Below can be read as, if notificationType is DEPOSIT_SUCCESS then dynamicValues must
// match the type of DepositSuccessNotification['dynamicValues']
export function createNotification<T extends NotificationType>(
  notificationType: T,
  dynamicValues: T extends NotificationType.DEPOSIT_SUCCESS
    ? DepositSuccessNotification['dynamicValues']
    : T extends NotificationType.ORDER_FILLED
      ? OrderFilledNotification['dynamicValues']
      : T extends NotificationType.ORDER_TRIGGERED
        ? OrderTriggeredNotification['dynamicValues'] : never,
): Notification {
  switch (notificationType) {
    case NotificationType.DEPOSIT_SUCCESS:
      return {
        type: NotificationType.DEPOSIT_SUCCESS,
        titleKey: LocalizationTitleKey.DEPOSIT_SUCCESS_TITLE,
        bodyKey: LocalizationBodyKey.DEPOSIT_SUCCESS_BODY,
        topic: Topic.TRADING,
        dynamicValues: dynamicValues as DepositSuccessNotification['dynamicValues'],
      } as DepositSuccessNotification;
    case NotificationType.ORDER_FILLED:
      return {
        type: NotificationType.ORDER_FILLED,
        titleKey: LocalizationTitleKey.ORDER_FILLED_TITLE,
        bodyKey: LocalizationBodyKey.ORDER_FILLED_BODY,
        topic: Topic.TRADING,
        dynamicValues: dynamicValues as OrderFilledNotification['dynamicValues'],
      } as OrderFilledNotification;
    case NotificationType.ORDER_TRIGGERED:
      return {
        type: NotificationType.ORDER_TRIGGERED,
        titleKey: LocalizationTitleKey.ORDER_TRIGGERED_TITLE,
        bodyKey: LocalizationBodyKey.ORDER_TRIGGERED_BODY,
        topic: Topic.TRADING,
        dynamicValues: dynamicValues as OrderTriggeredNotification['dynamicValues'],
      } as OrderTriggeredNotification;
      // Add other cases for new notification types here
    default:
      throw new Error('Unknown notification type');
  }
}
