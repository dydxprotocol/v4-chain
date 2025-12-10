import { QueryBuilder } from 'knex';

import {
  CHILD_SUBACCOUNT_MULTIPLIER,
  MAX_PARENT_SUBACCOUNTS,
} from '../constants';
import { knexReadReplica } from '../helpers/knex';

export function getParentSubaccountNum(childSubaccountNum: number): number {
  if (childSubaccountNum > MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER) {
    throw new Error(`Child subaccount number must be less than or equal to ${MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER}`);
  }
  return childSubaccountNum % MAX_PARENT_SUBACCOUNTS;
}

/**
 * Creates a subquery to find all subaccounts associated with a parent subaccount.
 * This is a common query used across tables when filtering by parent subaccount.
 *
 * @param parentSubaccount The parent subaccount object with address and subaccountNumber
 * @returns A knex query that selects subaccount IDs
 */
export function getSubaccountQueryForParent(parentSubaccount: {
  address: string,
  subaccountNumber: number,
}): QueryBuilder {
  return knexReadReplica.getConnection()
    .select('id as subaccountId')
    .from('subaccounts')
    .where('address', parentSubaccount.address)
    .andWhereRaw(
      '("subaccountNumber" - ?) % ? = 0',
      [parentSubaccount.subaccountNumber, MAX_PARENT_SUBACCOUNTS],
    );
}
