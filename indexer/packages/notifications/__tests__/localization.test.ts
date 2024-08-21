import {
  deriveLocalizedNotificationMessage,
} from '../src/localization';
import {
  NotificationType,
  NotificationDynamicFieldKey,
  createNotification,
  isValidLanguageCode,
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
      title: 'Order Filled',
      body: 'Your order for 1000 BTC/USD was filled at $45000',
    };

    const result = deriveLocalizedNotificationMessage(notification);
    expect(result).toEqual(expected);
  });

  describe('isValidLanguageCode', () => {
    test('should return true for valid language codes', () => {
      const validCodes = ['en', 'es', 'fr', 'de', 'it', 'ja', 'ko', 'zh'];
      validCodes.forEach((code) => {
        expect(isValidLanguageCode(code)).toBe(true);
      });
    });

    test('should return false for invalid language codes', () => {
      const invalidCodes = ['', 'EN', 'eng', 'esp', 'fra', 'deu', 'ita', 'jpn', 'kor', 'zho', 'xx'];
      invalidCodes.forEach((code) => {
        expect(isValidLanguageCode(code)).toBe(false);
      });
    });

    test('should return false for non-string inputs', () => {
      const nonStringInputs = [null, undefined, 123, {}, []];
      nonStringInputs.forEach((input) => {
        expect(isValidLanguageCode(input as any)).toBe(false);
      });
    });
  });
});
