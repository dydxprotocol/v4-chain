import { logger } from '@dydxprotocol-indexer/base';
import { sendFirebaseMessage } from '../src/message';
import { sendMulticast } from '../src/lib/firebase';
import { createNotification, NotificationType } from '../src/types';
import { testMocks, dbHelpers } from '@dydxprotocol-indexer/postgres';
import { defaultToken } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

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

  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  const mockNotification = createNotification(NotificationType.ORDER_FILLED, {
    AMOUNT: '10',
    MARKET: 'BTC-USD',
    AVERAGE_PRICE: '100.50',
  });

  it('should send a Firebase message successfully', async () => {
    await sendFirebaseMessage([defaultToken.token], mockNotification);

    expect(sendMulticast).toHaveBeenCalled();
    expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
      message: 'Firebase message sent successfully',
      notificationType: mockNotification.type,
    }));
  });

  it('should log a warning if user has no registration tokens', async () => {
    await sendFirebaseMessage([], mockNotification);

    expect(logger.warning).toHaveBeenCalledWith(expect.objectContaining({
      message: 'Attempted to send Firebase message to user with no registration tokens',
      notificationType: mockNotification.type,
    }));
  });

  it('should log an error if sending the message fails', async () => {
    const mockedSendMulticast = sendMulticast as jest.MockedFunction<typeof sendMulticast>;
    mockedSendMulticast.mockRejectedValueOnce(new Error('Send failed'));

    await sendFirebaseMessage([defaultToken.token], mockNotification);

    expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
      message: 'Failed to send Firebase message',
      notificationType: mockNotification.type,
    }));
  });
});
