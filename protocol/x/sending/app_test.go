package sending_test

import (
	"bytes"
	"math/big"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
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

func TestMsgCreateTransfer(t *testing.T) {
	tests := map[string]struct {
		/* Setup */
		// Initial balance of sender subaccount.
		senderInitialBalance uint64

		// Whether recipient subaccount exists.
		recipientDoesNotExist bool

		// Sender subaccount ID.
		senderSubaccountId satypes.SubaccountId

		// Recipient subaccount ID.
		recipientSubaccountId satypes.SubaccountId

		// Asset to transfer.
		asset assetstypes.Asset

		// Amount to transfer.
		amount uint64

		/* Expectations */
		// A string that CheckTx response should contain, if any.
		checkTxResponseContains string

		// Whether CheckTx fails.
		checkTxFails bool

		// Whether DeliverTx fails.
		deliverTxFails bool
	}{
		"Success: transfer from Alice subaccount 1 to Alice subaccount 2": {
			senderInitialBalance:  600_000_000,
			senderSubaccountId:    constants.Alice_Num0,
			recipientSubaccountId: constants.Alice_Num1,
			asset:                 *constants.Usdc,
			amount:                500_000_000,
		},
		"Success: transfer from Bob subaccount to Carl subaccount": {
			senderInitialBalance:  10_000_000,
			senderSubaccountId:    constants.Bob_Num0,
			recipientSubaccountId: constants.Carl_Num0,
			asset:                 *constants.Usdc,
			amount:                7_654_321,
		},
		// Transfer to a non-existent subaccount will create that subaccount and succeed.
		"Success: transfer from Alice subaccount to non-existent subaccount": {
			senderInitialBalance:  10_000_000,
			recipientDoesNotExist: true,
			senderSubaccountId:    constants.Alice_Num0,
			recipientSubaccountId: satypes.SubaccountId{
				Owner:  constants.BobAccAddress.String(),
				Number: 104,
			},
			asset:  *constants.Usdc,
			amount: 3_000_000,
		},
		"Failure: transfer more than balance": {
			senderInitialBalance:  600_000_000,
			senderSubaccountId:    constants.Alice_Num0,
			recipientSubaccountId: constants.Alice_Num1,
			asset:                 *constants.Usdc,
			amount:                600_000_001,
			deliverTxFails:        true,
		},
		"Failure: transfer a non-USDC asset": {
			senderSubaccountId:      constants.Alice_Num0,
			recipientSubaccountId:   constants.Alice_Num1,
			asset:                   *constants.BtcUsd, // non-USDC asset
			amount:                  7_000_000,
			checkTxResponseContains: "Non-USDC asset transfer not implemented",
			checkTxFails:            true,
		},
		"Failure: transfer zero amount": {
			senderSubaccountId:      constants.Alice_Num0,
			recipientSubaccountId:   constants.Alice_Num1,
			asset:                   *constants.Usdc,
			amount:                  0,
			checkTxResponseContains: "Invalid transfer amount",
			checkTxFails:            true,
		},
		"Failure: transfer from a subaccount to itself": {
			senderSubaccountId:      constants.Bob_Num0,
			recipientSubaccountId:   constants.Bob_Num0,
			asset:                   *constants.Usdc,
			amount:                  123_456,
			checkTxResponseContains: "Sender is the same as recipient",
			checkTxFails:            true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Set up tApp with indexer and sender subaccount balance of USDC.
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &(tc.senderSubaccountId),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId: constants.Usdc.Id,
										Index:   0,
										Quantums: dtypes.NewIntFromUint64(
											tc.senderInitialBalance,
										),
									},
								},
							},
						}
						if !tc.recipientDoesNotExist {
							genesisState.Subaccounts = append(
								genesisState.Subaccounts,
								satypes.Subaccount{
									Id: &(tc.recipientSubaccountId),
									AssetPositions: []*satypes.AssetPosition{
										{
											AssetId: constants.Usdc.Id,
											Index:   0,
											Quantums: dtypes.NewIntFromUint64(
												rand.NewRand().Uint64(),
											),
										},
									},
								},
							)
						}
					},
				)
				return genesis
			}).WithAppOptions(appOpts).Build()
			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Clear any messages produced prior to CheckTx calls.
			msgSender.Clear()

			senderQuantumsBeforeTransfer :=
				getSubaccountAssetQuantums(tApp.App.SubaccountsKeeper, ctx, tc.senderSubaccountId, tc.asset)
			recipientQuantumsBeforeTransfer :=
				getSubaccountAssetQuantums(tApp.App.SubaccountsKeeper, ctx, tc.recipientSubaccountId, tc.asset)

			// Construct message.
			msgCreateTransfer := sendingtypes.MsgCreateTransfer{
				Transfer: &sendingtypes.Transfer{
					Sender:    tc.senderSubaccountId,
					Recipient: tc.recipientSubaccountId,
					AssetId:   tc.asset.Id,
					Amount:    tc.amount,
				},
			}

			// Invoke CheckTx.
			CheckTx_MsgCreateTransfer := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(&msgCreateTransfer),
					Gas:                  100_000,
					FeeAmt:               constants.TestFeeCoins_5Cents,
				},
				&msgCreateTransfer,
			)
			checkTxResp := tApp.CheckTx(CheckTx_MsgCreateTransfer)

			// Check that CheckTx response log contains expected string, if any.
			if tc.checkTxResponseContains != "" {
				require.Contains(t, checkTxResp.Log, tc.checkTxResponseContains)
			}
			// Check that CheckTx succeeds or fails as expected.
			if tc.checkTxFails {
				require.Conditionf(t, checkTxResp.IsErr, "Expected CheckTx to error. Response: %+v", checkTxResp)
				return
			}
			require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

			// Check that no indexer events are emitted so far.
			require.Empty(t, msgSender.GetOnchainMessages())

			if tc.deliverTxFails {
				// Check that DeliverTx fails on MsgCreateTransfer.
				tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
					ValidateDeliverTxs: func(
						context sdk.Context,
						request abcitypes.RequestDeliverTx,
						response abcitypes.ResponseDeliverTx,
						_ int,
					) (haltChain bool) {
						if bytes.Equal(request.Tx, CheckTx_MsgCreateTransfer.Tx) {
							require.True(t, response.IsErr())
						} else {
							require.True(t, response.IsOK())
						}
						return false
					},
				})
				return
			} else {
				// Advance to block 3.
				ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			}

			// Verify expected sender subaccount balance.
			senderQuantumsAfterTransfer :=
				getSubaccountAssetQuantums(tApp.App.SubaccountsKeeper, ctx, tc.senderSubaccountId, tc.asset)
			require.Equal(
				t,
				new(big.Int).Sub(
					senderQuantumsBeforeTransfer,
					new(big.Int).SetUint64(tc.amount),
				),
				senderQuantumsAfterTransfer,
			)
			// Verify expected recipient subaccount balance.
			recipientQuantumsAfterTransfer :=
				getSubaccountAssetQuantums(tApp.App.SubaccountsKeeper, ctx, tc.recipientSubaccountId, tc.asset)
			require.Equal(
				t,
				new(big.Int).Add(
					recipientQuantumsBeforeTransfer,
					new(big.Int).SetUint64(tc.amount),
				),
				recipientQuantumsAfterTransfer,
			)
			// Verify that there are no offchain messages.
			require.Empty(t, msgSender.GetOffchainMessages())
			// Verify expected indexer events.
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
									&tc.senderSubaccountId,
									[]*satypes.PerpetualPosition{},
									[]*satypes.AssetPosition{
										{
											AssetId:  assetstypes.AssetUsdc.Id,
											Quantums: dtypes.NewIntFromBigInt(senderQuantumsAfterTransfer),
										},
									},
									nil, // no funding payment should have occurred
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeSubaccountUpdate,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&tc.recipientSubaccountId,
									[]*satypes.PerpetualPosition{},
									[]*satypes.AssetPosition{
										{
											AssetId:  assetstypes.AssetUsdc.Id,
											Quantums: dtypes.NewIntFromBigInt(recipientQuantumsAfterTransfer),
										},
									},
									nil, // no funding payment should have occurred
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeTransfer,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
							Version:             indexerevents.TransferEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewTransferEvent(
									tc.senderSubaccountId,
									tc.recipientSubaccountId,
									tc.asset.Id,
									satypes.BaseQuantums(tc.amount),
								),
							),
						},
					},
					TxHashes: []string{string(lib.GetTxHash(CheckTx_MsgCreateTransfer.GetTx()))},
				},
			)}
			require.ElementsMatch(t, expectedOnchainMessages, msgSender.GetOnchainMessages())
		})
	}
}

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
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).Build()
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
					FeeAmt:               constants.TestFeeCoins_5Cents,
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
			require.Equal(
				t,
				accountBalanceAfterDeposit,
				accountBalanceBeforeDeposit.Sub(transferredCoin).Sub(constants.TestFeeCoins_5Cents[0]),
			)
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
	tApp := testapp.NewTestAppBuilder(t).Build()
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

	testNonExistentSender(t, tApp, ctx, &msgDepositToSubaccount, randomAccount.PrivKey)
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
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).Build()
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
					Gas:                  constants.TestGasLimit,
					FeeAmt:               constants.TestFeeCoins_5Cents,
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
			if tc.subaccountId.Owner == tc.accountAccAddress.String() {
				accountBalanceAfterWithdraw = accountBalanceAfterWithdraw.Add(constants.TestFeeCoins_5Cents[0])
			}
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
	tApp := testapp.NewTestAppBuilder(t).Build()
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

	testNonExistentSender(t, tApp, ctx, &msgWithdrawFromSubaccount, randomAccount.PrivKey)
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
		constants.TestFeeCoins_5Cents,
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
