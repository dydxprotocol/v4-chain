import { PnlTicksCreateObject } from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
import { RedisClient } from 'redis';

import { hGetAllAsync, hSetAsync } from '../helpers/redis';
import { PnlTickForSubaccounts } from '../types';

const KEY: string = 'v4/latest-accounts-pnl-tick';

export async function getAll(client: RedisClient): Promise<PnlTickForSubaccounts> {
  const stringMap: { [subaccountId: string]: string } = await hGetAllAsync(KEY, client);
  return _.mapValues(stringMap, (val) => JSON.parse(val));
}

export async function set(
  ticksForAccount: PnlTickForSubaccounts,
  client: RedisClient,
): Promise<number> {
  return hSetAsync({
    hash: KEY,
    pairs: _.mapValues(ticksForAccount, (val: PnlTicksCreateObject) => JSON.stringify(val)),
  }, client);
}
