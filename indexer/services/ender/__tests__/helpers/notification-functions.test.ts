import {
  sendOrderFilledNotification,
  sendOrderTriggeredNotification,
} from '../../src/helpers/notifications/notifications-functions';
import {
  dbHelpers,
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
  defaultToken,
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

      await sendOrderFilledNotification(mockOrder, mockMarket);

      // Assert that createNotification was called with correct arguments
      expect(createNotification).toHaveBeenCalledWith(
        NotificationType.ORDER_FILLED,
        {
          AMOUNT: '10',
          MARKET: 'BTC-USD',
          AVERAGE_PRICE: '100.50',
        },
      );

      // Assert that sendFirebaseMessage was called with correct arguments, default wallet
      // is expected because mockOrder uses defaultSubaccountId
      expect(sendFirebaseMessage).toHaveBeenCalledWith([defaultToken.token], undefined);
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

        expect(sendFirebaseMessage).toHaveBeenCalledWith([defaultToken.token], undefined);
      });
    });
  });
});
