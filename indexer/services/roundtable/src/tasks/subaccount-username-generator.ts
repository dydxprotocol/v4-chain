import { logger } from '@dydxprotocol-indexer/base';
import {
  SubaccountUsernamesTable,
  SubaccountsWithoutUsernamesResult,
} from '@dydxprotocol-indexer/postgres';
import { generateUsername } from 'unique-username-generator';

import config from '../config';

export default async function runTask(): Promise<void> {
  const subaccounts:
  SubaccountsWithoutUsernamesResult[] = await
  SubaccountUsernamesTable.getSubaccountsWithoutUsernames();
  for (const subaccount of subaccounts) {
    const username: string = generateUsername('', config.SUBACCOUNT_USERNAME_NUM_RANDOM_DIGITS);
    try {
      // if insert fails, try it in the next roundtable cycle
      // There are roughly 50 Billion possible usernames with 3 random digits
      // so the chance of a collision is very low
      await SubaccountUsernamesTable.create({
        username,
        subaccountId: subaccount.subaccountId,
      });
    } catch (e) {
      logger.error({
        at: 'subaccount-username-generator#runTask',
        message: 'Failed to insert username for subaccount',
        subaccountId: subaccount.subaccountId,
        username,
        error: e,
      });
    }
  }
}
