import {
  deriveLocalizedNotificationMessage,
} from '../src/localization';
import {
  Notification,
  NotificationType,
  LocalizationKey,
  NotificationDynamicFieldKey,
  Deeplink,
  createNotification,
} from '../src/types';

describe('deriveLocalizedNotificationMessage', () => {
  test('should generate a correct message for DepositSuccessNotification', () => {
    const notification = createNotification(NotificationType.DEPOSIT_SUCCESS, {
      [NotificationDynamicFieldKey.AMOUNT]: '1000',
      [NotificationDynamicFieldKey.MARKET]: 'USDT',
    });

    const expected = {
      title: 'Deposit successful',
      body: 'You have successfully deposited 1000 USDT to your dYdX account.',
    };

    const result = deriveLocalizedNotificationMessage(notification);
    expect(result).toEqual(expected);
  });

  test('should generate a correct message for OrderFilledNotification', () => {
    const notification = createNotification(NotificationType.ORDER_FILLED, {
      [NotificationDynamicFieldKey.MARKET]: 'BTC/USD',
      [NotificationDynamicFieldKey.AVERAGE_PRICE]: '45000',
    });

    const expected = {
      title: 'Filled BTC/USD order at 45000.',
      body: 'Order Filled successful',
    };

    const result = deriveLocalizedNotificationMessage(notification);
    expect(result).toEqual(expected);
  });

  test('should throw an error for unknown notification type', () => {
    const unknownNotification = {
      type: 'UNKNOWN_TYPE' as NotificationType,
      titleKey: LocalizationKey.DEPOSIT_SUCCESS_TITLE,
      bodyKey: LocalizationKey.DEPOSIT_SUCCESS_BODY,
      deeplink: Deeplink.DEPOSIT,
      dynamicValues: {
        [NotificationDynamicFieldKey.AMOUNT]: '1000',
        [NotificationDynamicFieldKey.MARKET]: 'USDT',
      },
    } as Notification;

    expect(() => deriveLocalizedNotificationMessage(unknownNotification)).toThrowError('Unknown notification type');
  });
});
