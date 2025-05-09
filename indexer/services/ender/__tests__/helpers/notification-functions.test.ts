import {
  sendOrderFilledNotification,
  sendOrderTriggeredNotification,
} from '../../src/helpers/notifications/notifications-functions';
import {
  dbHelpers, FillFromDatabase,
  OrderFromDatabase,
  PerpetualMarketFromDatabase,
  PerpetualMarketStatus,
  PerpetualMarketType,
  SubaccountTable,
  testMocks,
} from '@dydxprotocol-indexer/postgres';

import {
  createNotification,
  sendFirebaseMessage,
  NotificationType,
} from '@dydxprotocol-indexer/notifications';
import {
  defaultSubaccountId,
  defaultMarket,
  defaultFirebaseNotificationToken,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

// Mock only the sendFirebaseMessage function
jest.mock('@dydxprotocol-indexer/notifications', () => {
  const actualModule = jest.requireActual('@dydxprotocol-indexer/notifications');
  return {
    ...actualModule, // keep all other exports intact
    sendFirebaseMessage: jest.fn(),
    createNotification: jest.fn(),
  };
});

const mockMarket: PerpetualMarketFromDatabase = {
  id: '1',
  clobPairId: '1',
  ticker: 'BTC-USD',
  marketId: 1,
  status: PerpetualMarketStatus.ACTIVE,
  priceChange24H: '0',
  volume24H: '0',
  trades24H: 0,
  nextFundingRate: '0',
  openInterest: '0',
  quantumConversionExponent: 1,
  atomicResolution: 1,
  subticksPerTick: 1,
  stepBaseQuantums: 1,
  liquidityTierId: 1,
  marketType: PerpetualMarketType.ISOLATED,
  baseOpenInterest: '0',
  defaultFundingRate1H: '0',
};

describe('notification functions', () => {
  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });
  describe('sendOrderFilledNotification', () => {
    it('should create and send an order filled notification', async () => {
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
      const mockFill: FillFromDatabase = {
        id: '1',
        subaccountId: defaultSubaccountId,
        side: 'BUY',
        liquidity: 'TAKER',
        type: 'LIMIT',
        clobPairId: String(defaultMarket.id),
        size: '5',
        price: '100.25',
        quoteAmount: '501.25',
        eventId: Buffer.from('1'),
        transactionHash: '0x1234567890abcdef',
        createdAt: new Date().toISOString(),
        createdAtHeight: '900001',
      } as FillFromDatabase;

      await sendOrderFilledNotification(mockOrder, mockMarket, mockFill);

      // Assert that createNotification was called with correct arguments
      expect(createNotification).toHaveBeenCalledWith(
        NotificationType.ORDER_FILLED,
        {
          AMOUNT: '5',
          MARKET: 'BTC-USD',
          AVERAGE_PRICE: '100.25',
        },
      );

      expect(sendFirebaseMessage).toHaveBeenCalledWith(
        [
          expect.objectContaining({
            token: defaultFirebaseNotificationToken.token,
            language: defaultFirebaseNotificationToken.language,
          }),
        ],
        undefined,
      );
    });

    describe('sendOrderTriggeredNotification', () => {
      it('should create and send an order triggered notification', async () => {
        const subaccount = await SubaccountTable.findById(defaultSubaccountId);
        const mockOrder: OrderFromDatabase = {
          id: '1',
          subaccountId: subaccount!.id,
          clientId: '1',
          clobPairId: '1',
          side: 'BUY',
          size: '10',
          price: '100.50',
          type: 'LIMIT',
          status: 'OPEN',
          timeInForce: 'GTT',
          reduceOnly: false,
          orderFlags: '0',
          goodTilBlock: '1000000',
          createdAtHeight: '900000',
          clientMetadata: '0',
          triggerPrice: '99.00',
          updatedAt: new Date().toISOString(),
          updatedAtHeight: '900001',
        } as OrderFromDatabase;

        await sendOrderTriggeredNotification(mockOrder, mockMarket, subaccount!);

        expect(createNotification).toHaveBeenCalledWith(
          NotificationType.ORDER_TRIGGERED,
          {
            MARKET: 'BTC-USD',
            PRICE: '100.50',
            AMOUNT: '10',
          },
        );

        expect(sendFirebaseMessage).toHaveBeenCalledWith([expect.objectContaining(
          {
            token: defaultFirebaseNotificationToken.token,
            language: defaultFirebaseNotificationToken.language,
          },
        )], undefined);
      });
    });
  });
});
