import Long from 'long';

import {
  Network,
  LocalWallet,
  ValidatorClient,
  BECH32_PREFIX,
  SubaccountInfo,
  IndexerClient,
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

function sleep(milliseconds: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, milliseconds));
}

describe('transfers', () => {
  it('test deposit', async () => {
    const wallet = await LocalWallet.fromMnemonic(DYDX_LOCAL_MNEMONIC, BECH32_PREFIX);

    const validatorClient = await ValidatorClient.connect(Network.local().validatorConfig);
    const indexerClient = new IndexerClient(Network.local().indexerConfig);

    const subaccount = new SubaccountInfo(wallet, 0);

    // Check USDC asset position before
    let assetPosResp: any = await indexerClient.account.getSubaccountAssetPositions(DYDX_LOCAL_ADDRESS, 0);
    expect(assetPosResp).not.toBeNull();
    const positions = assetPosResp.positions;
    const usdcPositionSizeBefore = positions.length !== undefined && positions.length > 0 ? positions[0].size : '0';

    const tx = await validatorClient.post.deposit(
      subaccount,
      0,
      new Long(10_000_000),
    );

    await sleep(15000);  // wait 15s for deposit to complete
    const defaultSubaccountId: string = SubaccountTable.uuid(wallet.address!, 0);

    // Check DB
    const transfers: TransferFromDatabase[] = await TransferTable.findAllToOrFromSubaccountId(
      { subaccountId: [defaultSubaccountId] },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(transfers.length).toEqual(1);
    expect(transfers[0]).toEqual(
      expect.objectContaining({
        recipientSubaccountId: defaultSubaccountId,
        senderWalletAddress: wallet.address!,
        size: '10',
      }),
    );

    // Check API
    const response = await indexerClient.account.getSubaccountTransfers(DYDX_LOCAL_ADDRESS, 0);
    expect(response).not.toBeNull();
    const transfersFromApi = response.transfers;
    expect(transfersFromApi).not.toBeNull();
    const transfer: any = transfersFromApi[0];
    expect(transfer).toEqual(
      expect.objectContaining({
        sender: {
          address: wallet.address!,
        },
        recipient: {
          address: wallet.address!,
          subaccountNumber: 0,
        },
        size: '10',
        symbol: 'USDC',
        type: 'DEPOSIT',
      }),
    );
    // Check asset position size is incremented too
    assetPosResp = await indexerClient.account.getSubaccountAssetPositions(DYDX_LOCAL_ADDRESS, 0);
    expect(assetPosResp).not.toBeNull();
    const usdcPositionSizeAfter = assetPosResp.positions[0].size;
    // expect usdcPositionSizeAfter to be usdcPositionSizeBefore + 10
    expect(usdcPositionSizeAfter).toEqual((parseInt(usdcPositionSizeBefore) + 10).toString());
  });

});
