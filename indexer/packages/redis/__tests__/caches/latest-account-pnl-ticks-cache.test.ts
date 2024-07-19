import { getAll, set } from '../../src/caches/latest-account-pnl-ticks-cache';
import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import { subaccountUuid } from './constants';
import { PnlTickForSubaccounts } from '../../src';
import { testConstants } from '@dydxprotocol-indexer/postgres';

describe('latestAccountPnlTicksCache', () => {
  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('getAll', () => {
    it('get all latest subaccounts', async () => {
      const ticksForSubaccounts: PnlTickForSubaccounts = {
        [subaccountUuid]: testConstants.defaultPnlTick,
      };
      await set(
        ticksForSubaccounts,
        client,
      );

      expect(await getAll(client)).toEqual(ticksForSubaccounts);
    });

    it('returns empty object for an non-existent key', async () => {
      expect(await getAll(client)).toEqual({});
    });
  });
});
