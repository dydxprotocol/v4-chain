import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  SubaccountUsernamesTable,
  SubaccountsWithoutUsernamesResult,
} from '@dydxprotocol-indexer/postgres';

import config from '../config';
import { generateUsername } from '../helpers/usernames-helper';

export default async function runTask(): Promise<void> {
  const subaccounts:
  SubaccountsWithoutUsernamesResult[] = await
  SubaccountUsernamesTable.getSubaccountsWithoutUsernames();
  for (const subaccount of subaccounts) {
    const username: string = generateUsername();
    try {
      // if insert fails, try it in the next roundtable cycle
      // There are roughly ~87.5 million possible usernames
      // so the chance of a collision is very low
      await SubaccountUsernamesTable.create({
        username,
        subaccountId: subaccount.subaccountId,
      });
    } catch (e) {
      if (e instanceof Error && e.name === 'UniqueViolationError') {
        stats.increment(
          `${config.SERVICE_NAME}.subaccount-username-generator.collision`, 1);
      } else {
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
}
