import { TIME_IN_FORCE_TO_API_TIME_IN_FORCE, CHILD_SUBACCOUNT_MULTIPLIER, MAX_PARENT_SUBACCOUNTS } from '../constants';
import { APITimeInForce, TimeInForce } from '../types';

/**
 * Gets post-only boolean from Indexer TimeInForce value.
 * Only returns true if the passed in TimeInForce value is POST_ONLY
 * @param timeInForce
 * @returns
 */
export function isOrderTIFPostOnly(timeInForce: TimeInForce): boolean {
  return timeInForce === TimeInForce.POST_ONLY;
}

/**
 * Converts Indexer TimeInForce to APITimeInForce value.
 * Special cases POST_ONLY as GTT.
 * @param timeInForce
 * @returns
 */
export function orderTIFToAPITIF(timeInForce: TimeInForce): APITimeInForce {
  return TIME_IN_FORCE_TO_API_TIME_IN_FORCE[timeInForce];
}

/**
 * Gets a list of all possible child subaccount numbers for a parent subaccount number
 * Child subaccounts = [128*0+parentSubaccount, 128*1+parentSubaccount ... 128*999+parentSubaccount]
 * @param parentSubaccount
 * @returns
 */
export function getChildSubaccountNums(parentSubaccountNum: number): number[] {
  if (parentSubaccountNum >= MAX_PARENT_SUBACCOUNTS) {
    throw new Error(`Parent subaccount number must be less than ${MAX_PARENT_SUBACCOUNTS}`);
  }
  return Array.from({ length: CHILD_SUBACCOUNT_MULTIPLIER },
    (_, i) => MAX_PARENT_SUBACCOUNTS * i + parentSubaccountNum);
}

/**
 * Gets the parent subaccount number from a child subaccount number
 * Parent subaccount = childSubaccount % 128
 * @param childSubaccountNum
 * @returns
 */
export function getParentSubaccountNum(childSubaccountNum: number): number {
  if (childSubaccountNum > MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER) {
    throw new Error(`Child subaccount number must be less than ${MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER}`);
  }
  return childSubaccountNum % MAX_PARENT_SUBACCOUNTS;
}
