package sending_test

import (
	"math/big"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/rand"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sample_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestMsgDepositToSubaccount(t *testing.T) {
	tests := map[string]struct {
		// Account address.
		accountAccAddress sdk.AccAddress

		// Subaccount ID.
		subaccountId satypes.SubaccountId

		// Quantums to transfer.
		quantums *big.Int

		// Asset to transfer.
		asset assetstypes.Asset

		/* Expectations */
		// A string that CheckTx response should contain, if any.
		checkTxResponseContains string

		// Whether CheckTx errors.
		checkTxIsError bool
	}{
		"Deposit from Alice account to Alice subaccount": {
			accountAccAddress: constants.AliceAccAddress,
			subaccountId:      constants.Alice_Num0,
			quantums:          big.NewInt(500_000_000),
			asset:             *constants.Usdc,
		},
		"Deposit from Bob account to Carl subaccount": {
			accountAccAddress: constants.BobAccAddress,
			subaccountId:      constants.Carl_Num0,
			quantums:          big.NewInt(7_000_000),
			asset:             *constants.Usdc,
		},
		// Deposit to a non-existent subaccount will create that subaccount and succeed.
		"Deposit from Bob account to non-existent subaccount": {
			accountAccAddress: constants.BobAccAddress,
			subaccountId: satypes.SubaccountId{
				Owner:  constants.BobAccAddress.String(),
				Number: 104,
			},
			quantums: big.NewInt(7_000_000),
			asset:    *constants.Usdc,
		},
		"Deposit a non-USDC asset": {
			accountAccAddress:       constants.AliceAccAddress,
			subaccountId:            constants.Carl_Num0,
			quantums:                big.NewInt(7_000_000),
			asset:                   *constants.BtcUsd, // non-USDC asset
			checkTxResponseContains: "Non-USDC asset transfer not implemented",
			checkTxIsError:          true,
		},
		"Deposit zero amount": {
			accountAccAddress:       constants.AliceAccAddress,
			subaccountId:            constants.Carl_Num0,
			quantums:                big.NewInt(0), // 0 quantums
			asset:                   *constants.Usdc,
			checkTxResponseContains: "Invalid transfer amount",
			checkTxIsError:          true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Set up tApp.
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp := testapp.NewTestAppBuilder().WithTesting(t).WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts)).Build()
			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			// Clear any messages produced prior to CheckTx calls.
			msgSender.Clear()

			accountBalanceBeforeDeposit := tApp.App.BankKeeper.GetBalance(ctx, tc.accountAccAddress, tc.asset.Denom)
			subaccountQuantumsBeforeDeposit :=
				getSubaccountAssetQuantums(tApp.App.SubaccountsKeeper, ctx, tc.subaccountId, tc.asset)
			_, transferredCoin, _ := tApp.App.AssetsKeeper.ConvertAssetToCoin(ctx, tc.asset.Id, tc.quantums)

			// Construct message.
			msgDepositToSubaccount := sendingtypes.MsgDepositToSubaccount{
				Sender:    tc.accountAccAddress.String(),
				Recipient: tc.subaccountId,
				AssetId:   tc.asset.Id,
				Quantums:  tc.quantums.Uint64(),
			}

			// Invoke CheckTx.
			CheckTx_MsgDepositToSubaccount := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(&msgDepositToSubaccount),
					Gas:                  100_000,
				},
				&msgDepositToSubaccount,
			)
			checkTxResp := tApp.CheckTx(CheckTx_MsgDepositToSubaccount)

			// Check that CheckTx response log contains expected string, if any.
			if tc.checkTxResponseContains != "" {
				require.Contains(t, checkTxResp.Log, tc.checkTxResponseContains)
			}
			// Check that CheckTx succeeds or errors out as expected.
			if tc.checkTxIsError {
				require.Conditionf(t, checkTxResp.IsErr, "Expected CheckTx to error. Response: %+v", checkTxResp)
				return
			}
			require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

			// Check that no indexer events are emitted so far.
			require.Empty(t, msgSender.GetOnchainMessages())
			// Advance to block 3 for transactions to be delivered.
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

			// Check expected account balance.
			accountBalanceAfterDeposit := tApp.App.BankKeeper.GetBalance(ctx, tc.accountAccAddress, tc.asset.Denom)
			require.Equal(t, accountBalanceAfterDeposit, accountBalanceBeforeDeposit.Sub(transferredCoin))
			// Check expected subaccount asset position.
			subaccountQuantumsAfterDeposit :=
				getSubaccountAssetQuantums(tApp.App.SubaccountsKeeper, ctx, tc.subaccountId, tc.asset)
			require.Equal(t,
				subaccountQuantumsAfterDeposit,
				subaccountQuantumsBeforeDeposit.Add(subaccountQuantumsBeforeDeposit, tc.quantums),
			)
			// Check that there are no offchain messages.
			require.Empty(t, msgSender.GetOffchainMessages())
			// Check for expected indexer events.
			expectedOnchainMessages := []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 3,
					Time:   ctx.BlockTime(),
					Events: []*indexer_manager.IndexerTendermintEvent{
						{
							Subtype:             indexerevents.SubtypeSubaccountUpdate,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&tc.subaccountId,
									[]*satypes.PerpetualPosition{},
									[]*satypes.AssetPosition{
										{
											AssetId:  assetstypes.AssetUsdc.Id,
											Quantums: dtypes.NewIntFromBigInt(subaccountQuantumsAfterDeposit),
										},
									},
									nil, // no funding payment should have occurred
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeTransfer,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.TransferEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewDepositEvent(
									tc.accountAccAddress.String(),
									tc.subaccountId,
									tc.asset.Id,
									satypes.BaseQuantums(tc.quantums.Uint64()),
								),
							),
						},
					},
					TxHashes: []string{string(lib.GetTxHash(CheckTx_MsgDepositToSubaccount.GetTx()))},
				},
			)}
			require.ElementsMatch(t, expectedOnchainMessages, msgSender.GetOnchainMessages())
		})
	}
}

func TestMsgDepositToSubaccount_NonExistentAccount(t *testing.T) {
	// Setup tApp.
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	// Generate a random account.
	randomAccount := simtypes.RandomAccounts(rand.NewRand(), 1)[0]

	// Construct message with non-existent account.
	msgDepositToSubaccount := sendingtypes.MsgDepositToSubaccount{
		Sender:    randomAccount.Address.String(),
		Recipient: constants.Alice_Num1,
		AssetId:   assetstypes.AssetUsdc.Id,
		Quantums:  uint64(1_000_000),
	}

	testNonExistentSender(t, &tApp, ctx, &msgDepositToSubaccount, randomAccount.PrivKey)
}

func TestMsgWithdrawFromSubaccount(t *testing.T) {
	tests := map[string]struct {
		// Account address.
		accountAccAddress sdk.AccAddress

		// Subaccount ID.
		subaccountId satypes.SubaccountId

		// Quantums to transfer.
		quantums *big.Int

		// Asset to transfer.
		asset assetstypes.Asset

		/* Expectations */
		// A string that CheckTx response should contain, if any.
		checkTxResponseContains string

		// Whether CheckTx errors.
		checkTxIsError bool
	}{
		"Withdraw from Alice subaccount to Alice account": {
			accountAccAddress: constants.AliceAccAddress,
			subaccountId:      constants.Alice_Num0,
			quantums:          big.NewInt(500_000_000),
			asset:             *constants.Usdc,
		},
		"Withdraw from Bob subaccount to Alice account": {
			accountAccAddress: constants.AliceAccAddress,
			subaccountId:      constants.Bob_Num0,
			quantums:          big.NewInt(7_000_000),
			asset:             *constants.Usdc,
		},
		// Withdrawing to a non-existent account will create that account and succeed.
		"Withdraw from Bob subaccount to non-existent account": {
			accountAccAddress: sdk.MustAccAddressFromBech32(sample_testutil.AccAddress()), // a newly generated account
			subaccountId:      constants.Bob_Num0,
			quantums:          big.NewInt(7_000_000),
			asset:             *constants.Usdc,
		},
		"Withdraw a non-USDC asset": {
			accountAccAddress:       constants.AliceAccAddress,
			subaccountId:            constants.Carl_Num0,
			quantums:                big.NewInt(7_000_000),
			asset:                   *constants.BtcUsd, // non-USDC asset
			checkTxResponseContains: "Non-USDC asset transfer not implemented",
			checkTxIsError:          true,
		},
		"Withdraw zero amount": {
			accountAccAddress:       constants.AliceAccAddress,
			subaccountId:            constants.Carl_Num0,
			quantums:                big.NewInt(0), // 0 quantums
			asset:                   *constants.Usdc,
			checkTxResponseContains: "Invalid transfer amount",
			checkTxIsError:          true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Set up tApp.
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp := testapp.NewTestAppBuilder().WithTesting(t).WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts)).Build()
			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			// Clear any messages produced prior to CheckTx calls.
			msgSender.Clear()

			accountBalanceBeforeWithdraw := tApp.App.BankKeeper.GetBalance(ctx, tc.accountAccAddress, tc.asset.Denom)
			subaccountQuantumsBeforeWithdraw :=
				getSubaccountAssetQuantums(tApp.App.SubaccountsKeeper, ctx, tc.subaccountId, tc.asset)
			_, transferredCoin, _ := tApp.App.AssetsKeeper.ConvertAssetToCoin(ctx, tc.asset.Id, tc.quantums)

			// Construct message.
			msgWithdrawFromSubaccount := sendingtypes.MsgWithdrawFromSubaccount{
				Sender:    tc.subaccountId,
				Recipient: tc.accountAccAddress.String(),
				AssetId:   tc.asset.Id,
				Quantums:  tc.quantums.Uint64(),
			}

			// Invoke CheckTx.
			CheckTx_MsgWithdrawFromSubaccount := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(&msgWithdrawFromSubaccount),
					Gas:                  100_000,
				},
				&msgWithdrawFromSubaccount,
			)
			checkTxResp := tApp.CheckTx(CheckTx_MsgWithdrawFromSubaccount)

			// Check that CheckTx response log contains expected string, if any.
			if tc.checkTxResponseContains != "" {
				require.Contains(t, checkTxResp.Log, tc.checkTxResponseContains)
			}
			// Check that CheckTx succeeds or errors out as expected.
			if tc.checkTxIsError {
				require.Conditionf(t, checkTxResp.IsErr, "Expected CheckTx to error. Response: %+v", checkTxResp)
				return
			}
			require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

			// Check that no indexer events are emitted so far.
			require.Empty(t, msgSender.GetOnchainMessages())
			// Advance to block 3 for transactions to be delivered.
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

			// Check expected account balance.
			accountBalanceAfterWithdraw := tApp.App.BankKeeper.GetBalance(ctx, tc.accountAccAddress, tc.asset.Denom)
			require.Equal(t, accountBalanceAfterWithdraw, accountBalanceBeforeWithdraw.Add(transferredCoin))
			// Check expected subaccount asset position.
			subaccountQuantumsAfterWithdraw :=
				getSubaccountAssetQuantums(tApp.App.SubaccountsKeeper, ctx, tc.subaccountId, tc.asset)
			require.Equal(t,
				subaccountQuantumsAfterWithdraw,
				subaccountQuantumsBeforeWithdraw.Sub(subaccountQuantumsBeforeWithdraw, tc.quantums),
			)
			// Check that there are no offchain messages.
			require.Empty(t, msgSender.GetOffchainMessages())
			// Check for expected indexer events.
			expectedOnchainMessages := []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 3,
					Time:   ctx.BlockTime(),
					Events: []*indexer_manager.IndexerTendermintEvent{
						{
							Subtype:             indexerevents.SubtypeSubaccountUpdate,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&tc.subaccountId,
									[]*satypes.PerpetualPosition{},
									[]*satypes.AssetPosition{
										{
											AssetId:  assetstypes.AssetUsdc.Id,
											Quantums: dtypes.NewIntFromBigInt(subaccountQuantumsAfterWithdraw),
										},
									},
									nil, // no funding payment should have occurred
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeTransfer,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.TransferEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewWithdrawEvent(
									tc.subaccountId,
									tc.accountAccAddress.String(),
									tc.asset.Id,
									satypes.BaseQuantums(tc.quantums.Uint64()),
								),
							),
						},
					},
					TxHashes: []string{string(lib.GetTxHash(CheckTx_MsgWithdrawFromSubaccount.GetTx()))},
				},
			)}
			require.ElementsMatch(t, expectedOnchainMessages, msgSender.GetOnchainMessages())
		})
	}
}

func TestMsgWithdrawFromSubaccount_NonExistentSubaccount(t *testing.T) {
	// Setup tApp.
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	// Generate a random account.
	randomAccount := simtypes.RandomAccounts(rand.NewRand(), 1)[0]

	// Construct message with non-existent subaccount.
	msgWithdrawFromSubaccount := sendingtypes.MsgWithdrawFromSubaccount{
		Sender: satypes.SubaccountId{
			Owner:  randomAccount.Address.String(),
			Number: 0,
		},
		Recipient: constants.AliceAccAddress.String(),
		AssetId:   assetstypes.AssetUsdc.Id,
		Quantums:  uint64(1_000_000),
	}

	testNonExistentSender(t, &tApp, ctx, &msgWithdrawFromSubaccount, randomAccount.PrivKey)
}

// testNonExistentSender is a helper function that tests sending transfer messages with non-existent sender.
func testNonExistentSender(
	t *testing.T,
	tApp *testapp.TestApp,
	ctx sdk.Context,
	message sdk.Msg,
	privKey cryptotypes.PrivKey,
) {
	// Generate signed transaction.
	signedTx, err := sims.GenSignedMockTx(
		rand.NewRand(),
		tApp.App.TxConfig(),
		[]sdk.Msg{message},
		sdk.Coins{},
		100_000, // gas
		ctx.ChainID(),
		[]uint64{0}, // dummy account number
		[]uint64{1}, // dummy sequence number
		privKey,
	)
	require.NoError(t, err)
	// Encode signed transaction.
	bytes, err := tApp.App.TxConfig().TxEncoder()(signedTx)
	require.NoError(t, err)
	// Invoke CheckTx.
	checkTxResp := tApp.CheckTx(
		abcitypes.RequestCheckTx{
			Tx:   bytes,
			Type: abcitypes.CheckTxType_New,
		},
	)

	// Check that CheckTx failed due to unknown address.
	require.Conditionf(t, checkTxResp.IsErr, "Expected CheckTx to error. Response: %+v", checkTxResp)
	require.Contains(t, checkTxResp.Log, "unknown address")
}

// getSubaccountAssetQuantums returns the quantums of an asset that belongs to a subaccount.
func getSubaccountAssetQuantums(
	subaccountsKeeper satypes.SubaccountsKeeper,
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	asset assetstypes.Asset,
) *big.Int {
	subaccount := subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	for _, assetPosition := range subaccount.GetAssetPositions() {
		if assetPosition.AssetId == asset.Id {
			return assetPosition.Quantums.BigInt()
		}
	}
	return big.NewInt(0) // by default, subaccount has 0 of this `asset`.
}
