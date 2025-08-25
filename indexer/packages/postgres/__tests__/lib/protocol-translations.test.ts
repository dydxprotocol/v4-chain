import Big from 'big.js';
import { bigIntToBytes, ORDER_FLAG_LONG_TERM, ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerOrder_ConditionType,
  IndexerOrder_Side,
  IndexerOrder_TimeInForce,
  IndexerOrderId,
} from '@dydxprotocol-indexer/v4-protos';
import {
  funding8HourValuePpmTo1HourRate,
  fundingIndexToHumanFixedString,
  getGoodTilBlock,
  getGoodTilBlockTime,
  getStepSize,
  getTickSize,
  orderTypeToProtocolConditionType,
  priceToSubticks,
  protocolConditionTypeToOrderType,
  protocolOrderTIFToTIF,
  serializedQuantumsToAbsHumanFixedString,
  subticksToPrice,
  tifToProtocolOrderTIF,
} from '../../src/lib/protocol-translations';
import { defaultPerpetualMarket } from '../helpers/constants';
import { OrderType, TimeInForce } from '../../src/types';
import Long from 'long';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';

describe('protocolTranslations', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  const orderIdShortTerm: IndexerOrderId = {
    subaccountId: {
      owner: 'owner',
      number: 0,
    },
    clientId: 1,
    clobPairId: 0,
    orderFlags: ORDER_FLAG_SHORT_TERM,
  };
  const orderIdLongTerm: IndexerOrderId = {
    ...orderIdShortTerm,
    orderFlags: ORDER_FLAG_LONG_TERM,
  };
  const goodTilBlockOrder: IndexerOrder = {
    orderId: orderIdShortTerm,
    side: IndexerOrder_Side.SIDE_BUY,
    subticks: Long.fromValue(1_000_000, true),
    quantums: Long.fromValue(1_000_000, true),
    goodTilBlock: 100,
    timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
    reduceOnly: false,
    clientMetadata: 0,
    conditionType: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
    conditionalOrderTriggerSubticks: Long.fromValue(0, true),
    orderRouterAddress: '',
  };
  const goodTilBlockTimeOrder: IndexerOrder = {
    ...goodTilBlockOrder,
    orderId: orderIdLongTerm,
    goodTilBlock: undefined,
    goodTilBlockTime: 1_500_000_000,
  };
  const expectedGoodTilBlockTimeISO: string = '2017-07-14T02:40:00.000Z';

  describe('getTickSize', () => {
    it('successfully calculates tick size', () => {
      // 100 * 1e-8 * 1e-6 / 10e-10 = 1e-2
      expect(getTickSize(defaultPerpetualMarket)).toEqual(Big(10).pow(-2).toFixed());
    });
  });

  describe('subticksToPrice', () => {
    it('successfully converts subticks to price', () => {
      const subticks = '100';
      // 100 * 1e-8 * 1e-6 / 1e-10 = .01
      expect(subticksToPrice(subticks, defaultPerpetualMarket)).toEqual('0.01');
    });
  });

  describe('priceToSubticks', () => {
    it('successfully converts price to subticks', () => {
      const price = '0.01';
      // .01 * 1e-10 / 1e-6 / 1e-8 = 100
      expect(priceToSubticks(price, defaultPerpetualMarket)).toEqual('100');
    });
  });

  describe('getStepSize', () => {
    it('successfully calculates step size', () => {
      // 10 * 1e-10 = 1e-9
      expect(getStepSize(defaultPerpetualMarket)).toEqual(Big(10).pow(-9).toFixed());
    });
  });

  describe('fundingIndexToHumanFixedString', () => {
    it('successfully gets the human readable form of a funding index value', () => {
      // 1e3 * 1e-6 * 1e-6 / 1e-10 = 1e1
      expect(
        fundingIndexToHumanFixedString(
          '1000',
          defaultPerpetualMarket,
        ),
      ).toEqual(Big(10).pow(1).toFixed());
    });
  });

  describe('funding8HourValuePpmTo1HourRate', () => {
    it('successfully gets the human readable form of a funding rate', () => {
      // 8e6 / 1e-6 / 8 = 1e1
      expect(
        funding8HourValuePpmTo1HourRate(
          8000000,
        ),
      ).toEqual(Big(1).toFixed());
    });
  });

  describe('getGoodTilBlock', () => {
    it('gets goodTilBlock for order', () => {
      expect(getGoodTilBlock(goodTilBlockOrder)).toEqual(100);
    });

    it('returns undefined for order without goodTilBlock', () => {
      expect(getGoodTilBlock(goodTilBlockTimeOrder)).toBeUndefined();
    });
  });

  describe('getGoodTilBlockTime', () => {
    it('gets goodTilBlockTime as ISO string for order', () => {
      expect(getGoodTilBlockTime(goodTilBlockTimeOrder)).toEqual(expectedGoodTilBlockTimeISO);
    });

    it('returns undefined for order without goodTilBlockTime', () => {
      expect(getGoodTilBlockTime(goodTilBlockOrder)).toBeUndefined();
    });
  });

  describe('serializedQuantumsToAbsHumanFixedString', () => {
    it.each([
      [-1_000_000, -5, '10'],
      [-2_000_000_000, -8, '20'],
      [1_000_000_000, -5, '10000'],
      [1_000, -1, '100'],
      [-1_000_000_000, 15, '1000000000000000000000000'],
    ])('successfully converts serialized quantums (%d), atomic resolution (%d) to absolute fixed string',
      (sizeQuantums: number, atomicResolution: number, expectedFixedString: string) => {
        expect(
          serializedQuantumsToAbsHumanFixedString(
            bigIntToBytes(BigInt(sizeQuantums)), atomicResolution),
        ).toEqual(expectedFixedString);
      });
  });

  describe('protocolOrderTIFToTIF', () => {
    it.each([
      ['FOK', IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL, TimeInForce.FOK],
      ['IOC', IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC, TimeInForce.IOC],
      ['UNSPECIFIED', IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED, TimeInForce.GTT],
      ['POST_ONLY', IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY, TimeInForce.POST_ONLY],
    ])('successfully gets TimeInForce given protocol order TIF: %s', (
      _name: string,
      protocolTIF: IndexerOrder_TimeInForce,
      expectedTimeInForce: TimeInForce,
    ) => {
      expect(protocolOrderTIFToTIF(protocolTIF)).toEqual(expectedTimeInForce);
    });

    it('throws error if unrecognized protocolTIF given', () => {
      expect(
        () => {
          protocolOrderTIFToTIF(100 as IndexerOrder_TimeInForce);
        },
      ).toThrow(new Error('Unexpected TimeInForce from protocol: 100'));
    });
  });

  describe('tifToProtocolOrderTIF', () => {
    it.each([
      ['FOK', TimeInForce.FOK, IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL],
      ['IOC', TimeInForce.IOC, IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC],
      ['GTT', TimeInForce.GTT, IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED],
      ['POST_ONLY', TimeInForce.POST_ONLY, IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY],
    ])('successfully gets protocol order TIF given TimeInForce: %s', (
      _name: string,
      timeInForce: TimeInForce,
      expectedProtocolTIF: IndexerOrder_TimeInForce,
    ) => {
      expect(tifToProtocolOrderTIF(timeInForce)).toEqual(expectedProtocolTIF);
    });

    it('throws error if unrecognized TimeInForce given', () => {
      expect(() => {
        tifToProtocolOrderTIF('INVALID' as TimeInForce);
      }).toThrow(new Error('Unexpected TimeInForce: INVALID'));
    });
  });

  describe('protocolConditionTypeToOrderType', () => {
    it.each([
      ['UNSPECIFIED', IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED, OrderType.LIMIT],
      ['UNRECOGNIZED', IndexerOrder_ConditionType.UNRECOGNIZED, OrderType.LIMIT],
      ['STOP_LOSS', IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS, OrderType.STOP_LIMIT],
      ['TAKE_PROFIT', IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT, OrderType.TAKE_PROFIT],
    ])('successfully gets order type given protocol condition type: %s', (
      _name: string,
      conditionType: IndexerOrder_ConditionType,
      orderType: OrderType,
    ) => {
      expect(protocolConditionTypeToOrderType(conditionType)).toEqual(orderType);
    });

    it('throws error if unrecognized ConditionType given', () => {
      expect(() => {
        protocolConditionTypeToOrderType(100 as IndexerOrder_ConditionType);
      }).toThrow(new Error('Unexpected ConditionType: 100'));
    });
  });

  describe('orderTypeToProtocolConditionType', () => {
    it.each([
      ['LIMIT', OrderType.LIMIT, IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED],
      ['MARKET', OrderType.MARKET, IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED],
      [
        'TRAILING_STOP',
        OrderType.TRAILING_STOP,
        IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
      ],
      ['STOP_LIMIT', OrderType.STOP_LIMIT, IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS],
      ['STOP_MARKET', OrderType.STOP_MARKET, IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS],
      ['TAKE_PROFIT', OrderType.TAKE_PROFIT, IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT],
      [
        'TAKE_PROFIT_MARKET',
        OrderType.TAKE_PROFIT_MARKET,
        IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT,
      ],
    ])('successfully gets order type given protocol condition type: %s', (
      _name: string,
      orderType: OrderType,
      conditionType: IndexerOrder_ConditionType,
    ) => {
      expect(orderTypeToProtocolConditionType(orderType)).toEqual(conditionType);
    });

    it('throws error if unrecognized OrderType given', () => {
      expect(() => {
        orderTypeToProtocolConditionType('INVALID' as OrderType);
      }).toThrow(new Error('Unexpected OrderType: INVALID'));
    });
  });
});
