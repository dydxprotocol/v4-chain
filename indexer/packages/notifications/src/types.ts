export enum NotificationType {
  DEPOSIT_SUCCESS = 'DEPOSIT_SUCCESS',
  ORDER_FILLED = 'ORDER_FILLED',
  ORDER_TRIGGERED = 'ORDER_TRIGGERED',
  LIQUIDATION = 'LIQUIDATION',
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
export enum LocalizationKey {
  DEPOSIT_SUCCESS_TITLE = 'DEPOSIT_SUCCESS_TITLE',
  DEPOSIT_SUCCESS_BODY = 'DEPOSIT_SUCCESS_BODY',
  ORDER_FILLED_TITLE = 'ORDER_FILLED_TITLE',
  ORDER_FILLED_BODY = 'ORDER_FILLED_BODY',
  ORDER_TRIGGERED_TITLE = 'ORDER_TRIGGERED_TITLE',
  ORDER_TRIGGERED_BODY = 'ORDER_TRIGGERED_BODY',
}

// Deeplinks for each notification
export enum Deeplink {
  DEPOSIT = '/profile',
  ORDER_FILLED = '/profile',
  ORDER_TRIGGERED = '/profile',
}

export enum Topic {
  TRADING = 'trading',
  PRICE_ALERTS = 'price_alerts',
}

interface BaseNotification <T extends Record<string, string>> {
  type: NotificationType,
  titleKey: LocalizationKey;
  bodyKey: LocalizationKey;
  topic: Topic;
  deeplink: Deeplink;
  dynamicValues: T,
}

interface DepositSuccessNotification extends BaseNotification<{
  [NotificationDynamicFieldKey.AMOUNT]: string;
  [NotificationDynamicFieldKey.MARKET]: string;
}> {
  type: NotificationType.DEPOSIT_SUCCESS;
  titleKey: LocalizationKey.DEPOSIT_SUCCESS_TITLE;
  bodyKey: LocalizationKey.DEPOSIT_SUCCESS_BODY;
  topic: Topic.TRADING;
  dynamicValues: {
    [NotificationDynamicFieldKey.AMOUNT]: string;
    [NotificationDynamicFieldKey.MARKET]: string;
  }
}

interface OrderFilledNotification extends BaseNotification <{
  [NotificationDynamicFieldKey.MARKET]: string;
  [NotificationDynamicFieldKey.AVERAGE_PRICE]: string;
}>{
  type: NotificationType.ORDER_FILLED;
  titleKey: LocalizationKey.ORDER_FILLED_TITLE;
  bodyKey: LocalizationKey.ORDER_FILLED_BODY;
  topic: Topic.TRADING;
  dynamicValues: {
    [NotificationDynamicFieldKey.MARKET]: string;
    [NotificationDynamicFieldKey.AMOUNT]: string;
    [NotificationDynamicFieldKey.AVERAGE_PRICE]: string;
  };
}

interface OrderTriggeredNotification extends BaseNotification <{
  [NotificationDynamicFieldKey.MARKET]: string;
  [NotificationDynamicFieldKey.PRICE]: string;
}>{
  type: NotificationType.ORDER_TRIGGERED;
  titleKey: LocalizationKey.ORDER_TRIGGERED_TITLE;
  bodyKey: LocalizationKey.ORDER_TRIGGERED_BODY;
  topic: Topic.TRADING;
  dynamicValues: {
    [NotificationDynamicFieldKey.MARKET]: string;
    [NotificationDynamicFieldKey.AMOUNT]: string;
    [NotificationDynamicFieldKey.PRICE]: string;
  };
}

export type NotificationMesage = {
  title: string;
  body: string;
};

export type Notification =
DepositSuccessNotification |
OrderFilledNotification |
OrderTriggeredNotification;

export function createNotification<T extends NotificationType>(
  type: T,
  dynamicValues: T extends NotificationType.DEPOSIT_SUCCESS
    ? DepositSuccessNotification['dynamicValues']
    : T extends NotificationType.ORDER_FILLED
      ? OrderFilledNotification['dynamicValues']
      : T extends NotificationType.ORDER_TRIGGERED
        ? OrderTriggeredNotification['dynamicValues'] : never,
): Notification {
  switch (type) {
    case NotificationType.DEPOSIT_SUCCESS:
      return {
        type: NotificationType.DEPOSIT_SUCCESS,
        titleKey: LocalizationKey.DEPOSIT_SUCCESS_TITLE,
        bodyKey: LocalizationKey.DEPOSIT_SUCCESS_BODY,
        topic: Topic.TRADING,
        deeplink: Deeplink.DEPOSIT,
        dynamicValues: dynamicValues as DepositSuccessNotification['dynamicValues'],
      };
    case NotificationType.ORDER_FILLED:
      return {
        type: NotificationType.ORDER_FILLED,
        titleKey: LocalizationKey.ORDER_FILLED_TITLE,
        bodyKey: LocalizationKey.ORDER_FILLED_BODY,
        topic: Topic.TRADING,
        deeplink: Deeplink.ORDER_FILLED,
        dynamicValues: dynamicValues as OrderFilledNotification['dynamicValues'],
      };
    case NotificationType.ORDER_TRIGGERED:
      return {
        type: NotificationType.ORDER_TRIGGERED,
        titleKey: LocalizationKey.ORDER_TRIGGERED_TITLE,
        bodyKey: LocalizationKey.ORDER_TRIGGERED_BODY,
        topic: Topic.TRADING,
        deeplink: Deeplink.ORDER_TRIGGERED,
        dynamicValues: dynamicValues as OrderTriggeredNotification['dynamicValues'],
      };

      // Add other cases for new notification types here
    default:
      throw new Error('Unknown notification type');
  }
}
