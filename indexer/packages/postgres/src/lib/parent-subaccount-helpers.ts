import {
  CHILD_SUBACCOUNT_MULTIPLIER,
  MAX_PARENT_SUBACCOUNTS,
} from '../constants';

export function getParentSubaccountNum(childSubaccountNum: number): number {
  if (childSubaccountNum > MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER) {
    throw new Error(`Child subaccount number must be less than or equal to ${MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER}`);
  }
  return childSubaccountNum % MAX_PARENT_SUBACCOUNTS;
}
