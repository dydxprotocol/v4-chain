import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  SubaccountUsernamesTable,
  SubaccountsWithoutUsernamesResult,
  Transaction,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';

import config from '../config';
import { generateUsernameForSubaccount } from '../helpers/usernames-helper';

export default async function runTask(): Promise<void> {
  const taskStart: number = Date.now();

  const subaccountZerosWithoutUsername:
  SubaccountsWithoutUsernamesResult[] = await
  SubaccountUsernamesTable.getSubaccountZerosWithoutUsernames(
    config.SUBACCOUNT_USERNAME_BATCH_SIZE,
  );
  stats.timing(
    `${config.SERVICE_NAME}.get_subaccount_zeros_without_usernames.timing`,
    Date.now() - taskStart,
  );

  const txId: number = await Transaction.start();
  const txnStart: number = Date.now();
  try {
    let successCount: number = 0;
    for (const subaccount of subaccountZerosWithoutUsername) {
      for (let i = 0; i < config.ATTEMPT_PER_SUBACCOUNT; i++) {
        const username: string = generateUsernameForSubaccount(
          subaccount.address,
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
          }, { txId });
          // If success, break from loop and move to next subaccount.
          successCount += 1;
          break;
        } catch (e) {
          // There are roughly ~225 million possible usernames
          // so the chance of collision is very low.
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
    await Transaction.commit(txId);
    const subaccountAddresses = _.map(
      subaccountZerosWithoutUsername,
      (subaccount) => subaccount.address,
    );
    stats.timing(
      `${config.SERVICE_NAME}.subaccount_username_generator.txn.timing`,
      Date.now() - txnStart,
    );
    logger.info({
      at: 'subaccount-username-generator#runTask',
      message: 'Generated usernames',
      batchSize: subaccountZerosWithoutUsername.length,
      successCount,
      addressSample: subaccountAddresses.slice(0, 10),
      duration: Date.now() - taskStart,
    });
  } catch (error) {
    await Transaction.rollback(txId);
    logger.error({
      at: 'subaccount-username-generator#runTask',
      message: 'Error when updating totalVolume in wallets table',
      error,
    });
  }

  stats.timing(
    `${config.SERVICE_NAME}.subaccount_username_generator.total.timing`,
    Date.now() - taskStart,
  );
}
