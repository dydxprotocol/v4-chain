import {
  dbHelpers,
  OrderFromDatabase,
  MarketFromDatabase,
  OrderStatus,
  OrderType,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  protocolTranslations,
  testConstants,
  testMocks,
  TimeInForce,
  apiTranslations,
  BestEffortOpenedStatus,
  LiquidityTiersFromDatabase,
  helpers,
} from '@dydxprotocol-indexer/postgres';
import { OrderbookLevels, redisTestConstants } from '@dydxprotocol-indexer/redis';
import {
  IndexerOrder_TimeInForce,
  RedisOrder,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import {
  OrderbookLevelsToResponseObject,
  perpetualMarketToResponseObject,
  postgresAndRedisOrderToResponseObject,
  postgresOrderToResponseObject,
  redisOrderToResponseObject,
} from '../../../src/request-helpers/request-transformer';
import { OrderResponseObject } from '../../../src/types';

describe('request-transformer', () => {
  const ticker: string = testConstants.defaultPerpetualMarket.ticker;
  const order: OrderFromDatabase = {
    ...testConstants.defaultOrder,
    id: testConstants.defaultOrderId,
  };

  beforeAll(async () => {
    await dbHelpers.migrate();
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterAll(async () => {
    await dbHelpers.clearData();
    await dbHelpers.teardown();
  });

  describe('perpetualMarketToResponseObject', () => {
    it('successfully converts a perpetual market to a response object', () => {
      const perpetualMarket: PerpetualMarketFromDatabase = testConstants.defaultPerpetualMarket;
      const market: MarketFromDatabase = testConstants.defaultMarket;
      const liquidityTier: LiquidityTiersFromDatabase = testConstants.defaultLiquidityTier;
      expect(perpetualMarketToResponseObject(perpetualMarket, liquidityTier, market)).toEqual(
        {
          clobPairId: perpetualMarket.clobPairId,
          ticker: perpetualMarket.ticker,
          status: perpetualMarket.status,
          oraclePrice: market.oraclePrice,
          priceChange24H: perpetualMarket.priceChange24H,
          volume24H: perpetualMarket.volume24H,
          trades24H: perpetualMarket.trades24H,
          nextFundingRate: perpetualMarket.nextFundingRate,
          initialMarginFraction: helpers.ppmToString(Number(liquidityTier.initialMarginPpm)),
          maintenanceMarginFraction: helpers.ppmToString(
            helpers.getMaintenanceMarginPpm(
              Number(liquidityTier.initialMarginPpm),
              Number(liquidityTier.maintenanceFractionPpm),
            ),
          ),
          openInterest: perpetualMarket.openInterest,
          atomicResolution: perpetualMarket.atomicResolution,
          quantumConversionExponent: perpetualMarket.quantumConversionExponent,
          tickSize: Big(10).pow(-2).toFixed(), // 100 * 1e-8 * 1e-6 / 10e-10 = 1e-2
          stepSize: Big(10).pow(-9).toFixed(), // 10 * 1e-10 = 1e-9
          stepBaseQuantums: perpetualMarket.stepBaseQuantums,
          subticksPerTick: perpetualMarket.subticksPerTick,
          marketType: perpetualMarket.marketType,
          openInterestLowerCap: liquidityTier.openInterestLowerCap,
          openInterestUpperCap: liquidityTier.openInterestUpperCap,
          baseOpenInterest: perpetualMarket.baseOpenInterest,
          defaultFundingRate1H: perpetualMarket.defaultFundingRate1H,
        },
      );
    });
  });

  describe('OrderbookLevelsToResponseObject', () => {
    const perpetualMarket: PerpetualMarketFromDatabase = testConstants.defaultPerpetualMarket;
    const orderbookLevels: OrderbookLevels = {
      bids: [
        { humanPrice: '300.0', quantums: '1000000000', lastUpdated: redisTestConstants.defaultLastUpdated },
        { humanPrice: '250.0', quantums: '350000', lastUpdated: redisTestConstants.defaultLastUpdated },
        { humanPrice: '150.0', quantums: '300', lastUpdated: redisTestConstants.defaultLastUpdated },
      ],
      asks: [
        { humanPrice: '400.0', quantums: '200000000000', lastUpdated: redisTestConstants.defaultLastUpdated },
        { humanPrice: '550.0', quantums: '450000', lastUpdated: redisTestConstants.defaultLastUpdated },
        { humanPrice: '760.0', quantums: '6000', lastUpdated: redisTestConstants.defaultLastUpdated },
      ],
    };
    expect(OrderbookLevelsToResponseObject(orderbookLevels, perpetualMarket)).toEqual({
      bids: [
        { price: '300.0', size: '0.1' }, // 1,000,000,000 * 1e-10
        { price: '250.0', size: '0.000035' }, // 350,000 * 1e-10
        { price: '150.0', size: '0.00000003' }, // 300 * 1e-10
      ],
      asks: [
        { price: '400.0', size: '20' }, // 200,000,000,000 * 1e-10
        { price: '550.0', size: '0.000045' }, // 450,000 * 1e-10 = 4.5e-5
        { price: '760.0', size: '0.0000006' }, // 6,000 * 1e-10
      ],
    });
  });

  describe('postgresAndRedisOrderToResponseObject', () => {
    it('successfully converts a postgres and redis order to a response object', () => {
      const filledOrder: OrderFromDatabase = {
        ...order,
        status: OrderStatus.FILLED,
      };
      const responseObject: OrderResponseObject | undefined = postgresAndRedisOrderToResponseObject(
        filledOrder,
        {
          [testConstants.defaultSubaccountId]:
            testConstants.defaultSubaccount.subaccountNumber,
        },
        redisTestConstants.defaultRedisOrder,
      );
      const expectedRedisOrderTIF: TimeInForce = protocolTranslations.protocolOrderTIFToTIF(
        redisTestConstants.defaultRedisOrder.order!.timeInForce,
      );

      expect(responseObject).not.toBeUndefined();
      expect(responseObject).not.toEqual(
        postgresOrderToResponseObject(
          filledOrder,
          testConstants.defaultSubaccount.subaccountNumber,
        ),
      );
      expect(responseObject).not.toEqual(
        redisOrderToResponseObject(redisTestConstants.defaultRedisOrder),
      );
      expect(responseObject).toEqual({
        ...postgresOrderToResponseObject(
          filledOrder,
          testConstants.defaultSubaccount.subaccountNumber,
        ),
        size: redisTestConstants.defaultRedisOrder.size,
        price: redisTestConstants.defaultRedisOrder.price,
        timeInForce: apiTranslations.orderTIFToAPITIF(expectedRedisOrderTIF),
        postOnly: apiTranslations.isOrderTIFPostOnly(expectedRedisOrderTIF),
        reduceOnly: redisTestConstants.defaultRedisOrder.order!.reduceOnly,
        goodTilBlock: protocolTranslations.getGoodTilBlock(
          redisTestConstants.defaultRedisOrder.order!,
        )?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(
          redisTestConstants.defaultRedisOrder.order!,
        ),
        clientMetadata: redisTestConstants.defaultRedisOrder.order!.clientMetadata.toString(),
      });
    });

    it('successfully converts a postgres order to a response object', () => {
      const responseObject: OrderResponseObject | undefined = postgresAndRedisOrderToResponseObject(
        order,
        {
          [testConstants.defaultSubaccountId]:
            testConstants.defaultSubaccount.subaccountNumber,
        },
      );

      expect(responseObject).not.toBeUndefined();
      expect(responseObject).not.toEqual(
        redisOrderToResponseObject(redisTestConstants.defaultRedisOrder),
      );
      expect(responseObject).toEqual(
        postgresOrderToResponseObject(order, testConstants.defaultSubaccount.subaccountNumber),
      );
    });

    it('successfully converts a redis order to a response object', () => {
      const responseObject: OrderResponseObject | undefined = postgresAndRedisOrderToResponseObject(
        undefined,
        {
          [testConstants.defaultSubaccountId]:
            testConstants.defaultSubaccount.subaccountNumber,
        },
        redisTestConstants.defaultRedisOrder,
      );

      expect(responseObject).not.toBeUndefined();
      expect(responseObject).not.toEqual(
        postgresOrderToResponseObject(order, testConstants.defaultSubaccount.subaccountNumber),
      );
      expect(responseObject).toEqual(
        redisOrderToResponseObject(redisTestConstants.defaultRedisOrder),
      );
    });

    it('successfully converts undefined postgres order and null redis orderto undefined', () => {
      const responseObject: OrderResponseObject | undefined = postgresAndRedisOrderToResponseObject(
        undefined,
        {
          [testConstants.defaultSubaccountId]:
            testConstants.defaultSubaccount.subaccountNumber,
        },
        null,
      );

      expect(responseObject).toBeUndefined();
    });
  });

  describe('postgresOrderToResponseObject', () => {
    it(
      'successfully converts a postgres order with null `goodTilBlockTime` to a response object',
      () => {
        const responseObject: OrderResponseObject = postgresOrderToResponseObject(
          order,
          testConstants.defaultSubaccount.subaccountNumber,
        );

        expect(responseObject).toEqual({
          ...order,
          timeInForce: apiTranslations.orderTIFToAPITIF(order.timeInForce),
          postOnly: apiTranslations.isOrderTIFPostOnly(order.timeInForce),
          ticker,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        });
      },
    );

    it(
      'successfully converts a postgres order with null `goodTilBlock` to a response object',
      () => {
        const orderWithGoodTilBlockTime: OrderFromDatabase = {
          ...testConstants.defaultOrderGoodTilBlockTime,
          id: testConstants.defaultOrderId,
        };
        const responseObject: OrderResponseObject = postgresOrderToResponseObject(
          orderWithGoodTilBlockTime,
          testConstants.defaultSubaccount.subaccountNumber,
        );

        expect(responseObject).toEqual({
          ...orderWithGoodTilBlockTime,
          timeInForce: apiTranslations.orderTIFToAPITIF(order.timeInForce),
          postOnly: apiTranslations.isOrderTIFPostOnly(order.timeInForce),
          ticker,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        });
      },
    );
  });

  describe('redisOrderToResponseObject', () => {
    it.each([
      [
        'default order',
        redisTestConstants.defaultRedisOrder,
      ],
      [
        'FOK TimeinForce, Reduce-only true',
        {
          ...redisTestConstants.defaultRedisOrder,
          order: {
            ...redisTestConstants.defaultOrder,
            timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
            reduceOnly: true,
          },
        },
      ],
      [
        'IOC TimeinForce, Reduce-only false',
        {
          ...redisTestConstants.defaultRedisOrder,
          order: {
            ...redisTestConstants.defaultOrder,
            timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
            reduceOnly: false,
          },
        },
      ],
      [
        'PostOnly TimeinForce, Reduce-only true',
        {
          ...redisTestConstants.defaultRedisOrder,
          order: {
            ...redisTestConstants.defaultOrder,
            timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
            reduceOnly: true,
          },
        },
      ],
      [
        'Unspecified TimeinForce, Reduce-only false',
        {
          ...redisTestConstants.defaultRedisOrder,
          order: {
            ...redisTestConstants.defaultOrder,
            timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
            reduceOnly: false,
          },
        },
      ],
    ])('successfully converts a redis order to a response object: %s', (
      _name: string,
      redisOrder: RedisOrder,
    ) => {
      const responseObject: OrderResponseObject = redisOrderToResponseObject(redisOrder);
      const expectedRedisOrderTIF: TimeInForce = protocolTranslations.protocolOrderTIFToTIF(
        redisOrder.order!.timeInForce,
      );

      expect(responseObject).toEqual({
        id: redisOrder.id,
        subaccountId: redisTestConstants.defaultSubaccountUuid,
        clientId: redisOrder.order!.orderId!.clientId.toString(),
        clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(redisOrder.order!.side),
        size: redisTestConstants.defaultSize,
        totalFilled: '0',
        price: redisTestConstants.defaultPrice,
        type: OrderType.LIMIT,
        status: BestEffortOpenedStatus.BEST_EFFORT_OPENED,
        timeInForce: apiTranslations.orderTIFToAPITIF(expectedRedisOrderTIF),
        postOnly: apiTranslations.isOrderTIFPostOnly(expectedRedisOrderTIF),
        reduceOnly: redisOrder.order!.reduceOnly,
        orderFlags: redisOrder.order!.orderId!.orderFlags.toString(),
        goodTilBlock: protocolTranslations.getGoodTilBlock(redisOrder.order!)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(redisOrder.order!),
        ticker,
        clientMetadata: redisOrder.order!.clientMetadata.toString(),
        subaccountNumber: redisOrder.order!.orderId!.subaccountId!.number,
        orderRouterAddress: redisOrder.order!.orderRouterAddress ?? undefined,
      });
    });
  });
});
