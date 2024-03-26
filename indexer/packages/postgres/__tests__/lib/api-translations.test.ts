import { APITimeInForce, TimeInForce } from '../../src/types';
import {
  isOrderTIFPostOnly,
  orderTIFToAPITIF,
} from '../../src/lib/api-translations';

describe('apiTranslations', () => {
  describe('orderTIFToAPITIF', () => {
    it.each([
      ['FOK', TimeInForce.FOK, APITimeInForce.FOK],
      ['IOC', TimeInForce.IOC, APITimeInForce.IOC],
      ['POST_ONLY', TimeInForce.POST_ONLY, APITimeInForce.GTT],
      ['GTT', TimeInForce.GTT, APITimeInForce.GTT],
    ])('Converts order time in force to api time in force: %s', (
      _name: string,
      orderTimeInForce: TimeInForce,
      expectedApiTimeInForce: APITimeInForce,
    ) => {
      expect(orderTIFToAPITIF(orderTimeInForce)).toEqual(expectedApiTimeInForce);
    });
  });

  describe('isOrderTIFPostOnly', () => {
    it.each([
      ['FOK', TimeInForce.FOK, false],
      ['IOC', TimeInForce.IOC, false],
      ['POST_ONLY', TimeInForce.POST_ONLY, true],
      ['GTT', TimeInForce.GTT, false],
    ])('Gets postOnly from order time in force: %s', (
      _name: string,
      orderTimeInForce: TimeInForce,
      expectedPostOnly: boolean,
    ) => {
      expect(isOrderTIFPostOnly(orderTimeInForce)).toEqual(expectedPostOnly);
    });
  });
});
