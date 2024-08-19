import { sendOrderFilledNotification } from '../../src/helpers/notifications/notifications-functions';
import { OrderFromDatabase } from '@dydxprotocol-indexer/postgres';

import { createNotification, sendFirebaseMessage, NotificationType } from '@dydxprotocol-indexer/notifications';
import { defaultSubaccountId, defaultMarket } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

// Mock only the sendFirebaseMessage function
jest.mock('@dydxprotocol-indexer/notifications', () => {
  const actualModule = jest.requireActual('@dydxprotocol-indexer/notifications');
  return {
    ...actualModule, // keep all other exports intact
    sendFirebaseMessage: jest.fn(),
    createNotification: jest.fn(),
  };
});

describe('sendOrderFilledNotification', () => {
  it('should create and send a notification', async () => {
    const mockOrder: OrderFromDatabase = {
      id: '1',
      subaccountId: defaultSubaccountId,
      clientId: '1',
      clobPairId: String(defaultMarket.id),
      side: 'BUY',
      size: '10',
      totalFilled: '0',
      price: '100.50',
      type: 'LIMIT',
      status: 'OPEN',
      timeInForce: 'GTT',
      reduceOnly: false,
      orderFlags: '0',
      goodTilBlock: '1000000',
      createdAtHeight: '900000',
      clientMetadata: '0',
      triggerPrice: undefined,
      updatedAt: new Date().toISOString(),
      updatedAtHeight: '900001',
    } as OrderFromDatabase;

    await sendOrderFilledNotification(mockOrder);

    // Assert that createNotification was called with correct arguments
    expect(createNotification).toHaveBeenCalledWith(
      NotificationType.ORDER_FILLED,
      {
        AMOUNT: '10',
        MARKET: 'BTC-USD',
        AVERAGE_PRICE: '100.50',
      },
    );

    // Assert that sendFirebaseMessage was called with correct arguments
    expect(sendFirebaseMessage).toHaveBeenCalledWith(defaultSubaccountId, undefined);
  });

  it('should throw an error if market is not found', async () => {
    const mockOrder: OrderFromDatabase = {
      id: '1',
      subaccountId: 'subaccount123',
      clientId: '1',
      clobPairId: '1',
      side: 'BUY',
      size: '10',
      totalFilled: '0',
      price: '100.50',
      type: 'LIMIT',
      status: 'OPEN',
      timeInForce: 'GTT',
      reduceOnly: false,
      orderFlags: '0',
      goodTilBlock: '1000000',
      createdAtHeight: '900000',
      clientMetadata: '0',
      triggerPrice: undefined,
      updatedAt: new Date().toISOString(),
      updatedAtHeight: '900001',
    } as OrderFromDatabase;

    await expect(sendOrderFilledNotification(mockOrder)).rejects.toThrow('sendOrderFilledNotification # Market not found');
  });
});
