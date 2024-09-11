import { logger } from '@dydxprotocol-indexer/base';
import { sendFirebaseMessage } from '../src/message';
import { sendMulticast } from '../src/lib/firebase';
import { createNotification, NotificationType } from '../src/types';

jest.mock('../src/lib/firebase', () => ({
  sendMulticast: jest.fn(),
}));

describe('sendFirebaseMessage', () => {
  let loggerInfoSpy: jest.SpyInstance;
  let loggerWarnSpy: jest.SpyInstance;
  let loggerErrorSpy: jest.SpyInstance;

  beforeAll(() => {
    loggerInfoSpy = jest.spyOn(logger, 'info').mockImplementation();
    loggerWarnSpy = jest.spyOn(logger, 'warning').mockImplementation();
    loggerErrorSpy = jest.spyOn(logger, 'error').mockImplementation();
  });

  afterAll(() => {
    loggerInfoSpy.mockRestore();
    loggerWarnSpy.mockRestore();
    loggerErrorSpy.mockRestore();
  });

  const defaultToken = {
    token: 'faketoken',
    language: 'en',
  };

  const mockNotification = createNotification(NotificationType.ORDER_FILLED, {
    AMOUNT: '10',
    MARKET: 'BTC-USD',
    AVERAGE_PRICE: '100.50',
  });

  it('should send a Firebase message successfully', async () => {
    await sendFirebaseMessage(
      [{ token: defaultToken.token, language: defaultToken.language }],
      mockNotification,
    );

    expect(sendMulticast).toHaveBeenCalledWith(expect.objectContaining(
      {
        tokens: [defaultToken.token],
        notification: { body: 'Your order for 10 BTC-USD was filled at $100.50', title: 'Order Filled' },
      }));
  });

  it('should log an error if sending the message fails', async () => {
    const mockedSendMulticast = sendMulticast as jest.MockedFunction<typeof sendMulticast>;
    mockedSendMulticast.mockRejectedValueOnce(new Error('Send failed'));

    await sendFirebaseMessage(
      [{ token: defaultToken.token, language: defaultToken.language }],
      mockNotification,
    );

    expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
      message: 'Send failed',
      notificationType: mockNotification.type,
    }));
  });
});
