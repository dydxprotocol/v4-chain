import {
  deriveLocalizedNotificationMessage,
} from '../src/localization';
import {
  NotificationType,
  NotificationDynamicFieldKey,
  createNotification,
} from '../src/types';

describe('deriveLocalizedNotificationMessage', () => {
  test('should generate a correct message for DepositSuccessNotification', () => {
    const notification = createNotification(NotificationType.DEPOSIT_SUCCESS, {
      [NotificationDynamicFieldKey.AMOUNT]: '1000',
      [NotificationDynamicFieldKey.MARKET]: 'USDT',
    });

    const expected = {
      title: 'Deposit Successful',
      body: 'You have successfully deposited 1000 USDT to your dYdX account.',
    };

    const result = deriveLocalizedNotificationMessage(notification);
    expect(result).toEqual(expected);
  });

  test('should generate a correct message for OrderFilledNotification', () => {
    const notification = createNotification(NotificationType.ORDER_FILLED, {
      [NotificationDynamicFieldKey.MARKET]: 'BTC/USD',
      [NotificationDynamicFieldKey.AVERAGE_PRICE]: '45000',
      [NotificationDynamicFieldKey.AMOUNT]: '1000',
    });

    const expected = {
      title: 'Your order for 1000 BTC/USD was filled at $45000',
      body: 'Order Filled',
    };

    const result = deriveLocalizedNotificationMessage(notification);
    expect(result).toEqual(expected);
  });
});
