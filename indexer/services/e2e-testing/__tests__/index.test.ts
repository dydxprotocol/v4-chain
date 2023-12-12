import Long from 'long';

import {
  Network,
  LocalWallet,
  ValidatorClient,
  BECH32_PREFIX,
  SubaccountInfo,
} from '@dydxprotocol/v4-client-js';
import {
  Ordering,
  TransferColumns,
  TransferFromDatabase,
  TransferTable,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';

export const DYDX_LOCAL_ADDRESS = 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4';
export const DYDX_LOCAL_MNEMONIC = 'merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small';

describe('index', () => {
  it('test transfers', async () => {
    const wallet = await LocalWallet.fromMnemonic(DYDX_LOCAL_MNEMONIC, BECH32_PREFIX);
    // console.log(wallet);

    const client = await ValidatorClient.connect(Network.local().validatorConfig);
    // console.log('**Client**');
    // console.log(client);

    const subaccount = new SubaccountInfo(wallet, 0);
    const tx = await client.post.deposit(
      subaccount,
      0,
      new Long(10_000_000),
    );
    // console.log('**Deposit Tx**');
    // console.log(tx);

    const defaultSubaccountId: string = SubaccountTable.uuid(wallet.address!, 0);
    const transfers: TransferFromDatabase[] = await TransferTable.findAllToOrFromSubaccountId(
      { subaccountId: [defaultSubaccountId] },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(transfers.length).toBeGreaterThan(0);
  });
});
