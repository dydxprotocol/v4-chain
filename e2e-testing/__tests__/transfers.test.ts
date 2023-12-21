import Long from 'long';
import {
  BECH32_PREFIX,
  IndexerClient,
  LocalWallet,
  Network,
  SocketClient,
  SubaccountInfo,
  ValidatorClient,
} from '@dydxprotocol/v4-client-js';
import {
  Ordering,
  SubaccountTable,
  TransferColumns,
  TransferFromDatabase,
  TransferTable,
} from '@dydxprotocol-indexer/postgres';
import * as utils from './helpers/utils';
import { DYDX_LOCAL_ADDRESS, DYDX_LOCAL_MNEMONIC } from './helpers/constants';

describe('transfers', () => {
  it('test deposit', async () => {
    connectAndValidateSocketClient();
    const wallet = await LocalWallet.fromMnemonic(DYDX_LOCAL_MNEMONIC, BECH32_PREFIX);

    const validatorClient = await ValidatorClient.connect(Network.local().validatorConfig);
    const indexerClient = new IndexerClient(Network.local().indexerConfig);

    const subaccount = new SubaccountInfo(wallet, 0);

    // Check USDC asset position before
    let assetPosResp: any = await indexerClient.account.getSubaccountAssetPositions(
      DYDX_LOCAL_ADDRESS,
      0,
    );
    expect(assetPosResp).not.toBeNull();
    const positions = assetPosResp.positions;
    const usdcPositionSizeBefore = positions.length !== undefined && positions.length > 0 ? positions[0].size : '0';

    // Deposit
    await validatorClient.post.deposit(
      subaccount,
      0,
      new Long(10_000_000),
    );

    // TODO(IND-547): investigate deterministically advancing network height
    await utils.sleep(5000);  // wait 5s for deposit to complete
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

    // Check API /v4/transfers endpoint
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

    // Check API /v4/assetPositions endpoint
    assetPosResp = await indexerClient.account.getSubaccountAssetPositions(DYDX_LOCAL_ADDRESS, 0);
    expect(assetPosResp).not.toBeNull();
    const usdcPositionSizeAfter = assetPosResp.positions[0].size;
    // expect usdcPositionSizeAfter to be usdcPositionSizeBefore + 10
    expect(usdcPositionSizeAfter).toEqual((parseInt(usdcPositionSizeBefore, 10) + 10).toString());
  });

  function connectAndValidateSocketClient(): void {
    const mySocket = new SocketClient(
      Network.local().indexerConfig,
      () => {
      },
      () => {
      },
      (message) => {
        if (typeof message.data === 'string') {
          const data = JSON.parse(message.data as string);
          if (data.type === 'connected') {
            mySocket.subscribeToSubaccount(DYDX_LOCAL_ADDRESS, 0);
          } else if (data.type === 'subscribed') {
            expect(data.channel).toEqual('v4_subaccounts');
            expect(data.id).toEqual(`${DYDX_LOCAL_ADDRESS}/0`);
            expect(data.contents.subaccount).toEqual(
              expect.objectContaining({
                address: DYDX_LOCAL_ADDRESS,
                subaccountNumber: 0,
              }),
            );
          } else if (data.type === 'channel_data' && data.contents.transfers) {
            expect(data.contents.transfers).toEqual(
              expect.objectContaining({
                sender: {
                  address: DYDX_LOCAL_ADDRESS,
                },
                recipient: {
                  address: DYDX_LOCAL_ADDRESS,
                  subaccountNumber: 0,
                },
                size: '10',
                symbol: 'USDC',
                type: 'DEPOSIT',
              }),
            );
          }
        }
      },
    );
    mySocket.connect();
  }
});
