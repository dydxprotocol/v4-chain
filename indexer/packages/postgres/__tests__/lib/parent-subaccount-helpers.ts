import {
  getParentSubaccountNum,
} from '../../src/lib/parent-subaccount-helpers';

describe('getParentSubaccountNum', () => {
  it('Gets the parent subaccount number from a child subaccount number', () => {
    expect(getParentSubaccountNum(0)).toEqual(0);
    expect(getParentSubaccountNum(128)).toEqual(0);
    expect(getParentSubaccountNum(128 * 999 - 1)).toEqual(127);
  });

  it('Throws an error if the child subaccount number is greater than the max child subaccount number', () => {
    expect(() => getParentSubaccountNum(128001)).toThrowError('Child subaccount number must be less than or equal to 128000');
  });
});
