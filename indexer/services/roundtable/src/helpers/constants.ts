import { testConstants } from '@klyraprotocol-indexer/postgres';
import { PnlTickForSubaccounts } from '@klyraprotocol-indexer/redis';

export const defaultPnlTickForSubaccounts: PnlTickForSubaccounts = {
  [testConstants.defaultSubaccountId]: testConstants.defaultPnlTick,
  [testConstants.defaultSubaccountId2]: {
    ...testConstants.defaultPnlTick,
    subaccountId: testConstants.defaultSubaccountId2,
    equity: '9000',
  },
};

export const defaultZeroPerpYieldIndex: string = '0/1';
