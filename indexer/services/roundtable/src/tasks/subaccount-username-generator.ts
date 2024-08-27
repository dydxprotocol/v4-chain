import { randomInt } from 'crypto';

import { logger } from '@dydxprotocol-indexer/base';
import {
  SubaccountUsernamesTable,
  SubaccountsWithoutUsernamesResult,
} from '@dydxprotocol-indexer/postgres';

import config from '../config';
import adjectives from '../lib/adjectives.json';
import nouns from '../lib/nouns.json';

export default async function runTask(): Promise<void> {
  const subaccounts:
  SubaccountsWithoutUsernamesResult[] = await
  SubaccountUsernamesTable.getSubaccountsWithoutUsernames();
  for (const subaccount of subaccounts) {
    const username: string = generateUsername();
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

function generateUsername(): string {
  const randomAdjective: string = adjectives[randomInt(0, adjectives.length)];
  const randomNoun: string = nouns[randomInt(0, nouns.length)];
  const randomNumber: string = randomInt(0, 1000).toString().padStart(
    config.SUBACCOUNT_USERNAME_NUM_RANDOM_DIGITS, '0');

  const capitalizedAdjective: string = randomAdjective.charAt(
    0).toUpperCase() + randomAdjective.slice(1);
  const capitalizedNoun: string = randomNoun.charAt(0).toUpperCase() + randomNoun.slice(1);

  return `${capitalizedAdjective}${capitalizedNoun}${randomNumber}`;
}
