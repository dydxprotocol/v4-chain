import yargs from 'yargs';

import { validatePnl, validatePnlForSubaccount } from './helpers/pnl-validation-helpers';
import { runAsyncScript } from './helpers/util';

const args = yargs.options({
  subaccount_id: {
    type: 'string',
    alias: 's',
    description: 'Subaccount to validate pnl for',
  },
  pnl_ids: {
    type: 'string',
    array: true,
    alias: 'p',
    description: 'The list of pnl uuids to validate',
  },
}).argv;

runAsyncScript(async () => {
  if (args.pnl_ids) {
    for (const pnlId of args.pnl_ids) {
      await validatePnl(pnlId);
    }
  }
  if (args.subaccount_id) {
    await validatePnlForSubaccount(args.subaccount_id);
  }
});
