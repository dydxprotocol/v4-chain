import { APITimeInForce, TimeInForce } from '../../src/types';
import {
  getChildSubaccountNums,
  getParentSubaccountNum,
  isOrderTIFPostOnly,
  orderTIFToAPITIF
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

  describe('getChildSubaccountNums', () => {
    it('Gets a list of all possible child subaccount numbers for a parent subaccount number', () => {
      const childSubaccounts = getChildSubaccountNums(0);
      expect(childSubaccounts.length).toEqual(1000);
      expect(childSubaccounts[0]).toEqual(0);
      expect(childSubaccounts[1]).toEqual(128);
      expect(childSubaccounts[999]).toEqual(128 * 999);
    });
  });

  describe('getChildSubaccountNums', () => {
    it('Throws an error if the parent subaccount number is greater than or equal to the maximum parent subaccount number', () => {
      expect(() => getChildSubaccountNums(128)).toThrowError('Parent subaccount number must be less than 128');
    });
  });

  describe('getParentSubaccountNum', () => {
    it('Gets the parent subaccount number from a child subaccount number', () => {
      expect(getParentSubaccountNum(0)).toEqual(0);
      expect(getParentSubaccountNum(128)).toEqual(0);
      expect(getParentSubaccountNum(128 * 999 - 1)).toEqual(127);
    });
  });

  describe('getParentSubaccountNum', () => {
    it('Throws an error if the child subaccount number is greater than the max child subaccount number', () => {
      expect(() => getParentSubaccountNum(128001)).toThrowError('Child subaccount number must be less than 128000');
    });
  });
});
