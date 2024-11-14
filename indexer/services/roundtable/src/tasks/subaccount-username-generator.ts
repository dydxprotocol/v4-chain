import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  SubaccountUsernamesTable,
  SubaccountsWithoutUsernamesResult,
} from '@dydxprotocol-indexer/postgres';

import config from '../config';
import { generateUsernameForSubaccount } from '../helpers/usernames-helper';

export default async function runTask(): Promise<void> {
  const subaccountZerosWithoutUsername:
  SubaccountsWithoutUsernamesResult[] = await
  SubaccountUsernamesTable.getSubaccountZerosWithoutUsernames(
    config.SUBACCOUNT_USERNAME_BATCH_SIZE,
  );
  for (const subaccount of subaccountZerosWithoutUsername) {
    for (let i = 0; i < config.ATTEMPT_PER_SUBACCOUNT; i++) {
      const username: string = generateUsernameForSubaccount(
        subaccount.subaccountId,
        // Always use subaccountNum 0 for generation. Effectively we are
        // generating one username per address. The fact that we are storing
        // in the `subaccount_usernames` table is a tech debt.
        0,
        // generation nonce
        i,
      );
      try {
        await SubaccountUsernamesTable.create({
          username,
          subaccountId: subaccount.subaccountId,
        });
        // If success, break from loop and move to next subaccount.
        break;
      } catch (e) {
        // There are roughly ~225 million possible usernames
        // so the chance of a collision is very lo
        if (e instanceof Error && e.name === 'UniqueViolationError') {
          stats.increment(
            `${config.SERVICE_NAME}.subaccount-username-generator.collision`, 1);
          logger.info({
            at: 'subaccount-username-generator#runTask',
            message: 'username collision',
            address: subaccount.address,
            subaccountId: subaccount.subaccountId,
            username,
            error: e,
          });
        } else {
          logger.error({
            at: 'subaccount-username-generator#runTask',
            message: 'Failed to insert username for subaccount',
            address: subaccount.address,
            subaccountId: subaccount.subaccountId,
            username,
            error: e,
          });
        }
      }
    }
  }
}
