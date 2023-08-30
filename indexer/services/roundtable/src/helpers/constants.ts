import { testConstants } from '@dydxprotocol-indexer/postgres';
import { PnlTickForSubaccounts } from '@dydxprotocol-indexer/redis';

export const defaultPnlTickForSubaccounts: PnlTickForSubaccounts = {
  [testConstants.defaultSubaccountId]: testConstants.defaultPnlTick,
  [testConstants.defaultSubaccountId2]: {
    ...testConstants.defaultPnlTick,
    subaccountId: testConstants.defaultSubaccountId2,
    equity: '9000',
  },
};
