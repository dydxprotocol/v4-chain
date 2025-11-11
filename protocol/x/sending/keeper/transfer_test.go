package keeper_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TransferTestCase struct {
	// Setup.
	subaccounts []satypes.Subaccount
	transfer    *types.Transfer
	// Expectations.
	expectedSubaccountBalance map[satypes.SubaccountId]*big.Int
	expectedErr               string
}

// assertTransferEventInIndexerBlock verifies that the transfer has a corresponding transfer
// event in the Indexer Kafka message.
func assertTransferEventInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	transfer *types.Transfer,
) {
	block := k.GetIndexerEventManager().ProduceBlock(ctx)
	expectedEvent := k.GenerateTransferEvent(transfer)
	var transfers []*indexerevents.TransferEventV1
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeTransfer {
			continue
		}
		var transfer indexerevents.TransferEventV1
		err := proto.Unmarshal(event.DataBytes, &transfer)
		if err != nil {
			panic(err)
		}
		transfers = append(transfers, &transfer)
	}
	require.Contains(t, transfers, expectedEvent)
}

// assertDepositEventInIndexerBlock verifies that the deposit has a corresponding deposit
// event in the Indexer Kafka message.
func assertDepositEventInIndexerBlock(
	t *testing.T,
	ctx sdk.Context,
	k *keeper.Keeper,
	deposit *types.MsgDepositToSubaccount,
) {
	block := k.GetIndexerEventManager().ProduceBlock(ctx)
	expectedEvent := k.GenerateDepositEvent(deposit)
	var deposits []*indexerevents.TransferEventV1
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeTransfer {
			continue
		}
		var deposit indexerevents.TransferEventV1
		err := proto.Unmarshal(event.DataBytes, &deposit)
		if err != nil {
			panic(err)
		}
		deposits = append(deposits, &deposit)
	}
	require.Contains(t, deposits, expectedEvent)
}

// assertWithdrawEventInIndexerBlock verifies that the withdraw has a corresponding withdraw
// event in the Indexer Kafka message.
func assertWithdrawEventInIndexerBlock(
	t *testing.T,
	ctx sdk.Context,
	k *keeper.Keeper,
	withdraw *types.MsgWithdrawFromSubaccount,
) {
	block := k.GetIndexerEventManager().ProduceBlock(ctx)
	expectedEvent := k.GenerateWithdrawEvent(withdraw)
	var withdraws []*indexerevents.TransferEventV1
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeTransfer {
			continue
		}
		var withdraw indexerevents.TransferEventV1
		err := proto.Unmarshal(event.DataBytes, &withdraw)
		if err != nil {
			panic(err)
		}
		withdraws = append(withdraws, &withdraw)
	}
	require.Contains(t, withdraws, expectedEvent)
}

func runProcessTransferTest(t *testing.T, tc TransferTestCase) {
	ks := keepertest.SendingKeepers(t)
	ks.Ctx = ks.Ctx.WithBlockHeight(5)
	keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

	perpetuals := []perptypes.Perpetual{
		constants.BtcUsd_100PercentMarginRequirement,
	}
	require.NoError(t, keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper))

	for _, p := range perpetuals {
		_, err := ks.PerpetualsKeeper.CreatePerpetual(
			ks.Ctx,
			p.Params.Id,
			p.Params.Ticker,
			p.Params.MarketId,
			p.Params.AtomicResolution,
			p.Params.DefaultFundingPpm,
			p.Params.LiquidityTier,
			p.Params.MarketType,
		)
		require.NoError(t, err)
	}

	for _, s := range tc.subaccounts {
		ks.SubaccountsKeeper.SetSubaccount(
			ks.Ctx,
			s,
		)
		ks.AccountKeeper.SetAccount(
			ks.Ctx,
			ks.AccountKeeper.NewAccountWithAddress(ks.Ctx, s.GetId().MustGetAccAddress()),
		)
	}

	err := ks.SendingKeeper.ProcessTransfer(ks.Ctx, tc.transfer)
	for subaccountId, expectedQuoteBalance := range tc.expectedSubaccountBalance {
		subaccount := ks.SubaccountsKeeper.GetSubaccount(ks.Ctx, subaccountId)
		require.Equal(t, expectedQuoteBalance, subaccount.GetUsdcPosition())
	}
	if tc.expectedErr != "" {
		require.ErrorContains(t, err, tc.expectedErr)
	} else {
		require.NoError(t, err)
		assertTransferEventInIndexerBlock(
			t,
			ks.SendingKeeper,
			ks.Ctx,
			tc.transfer,
		)
	}
}

func TestProcessTransfer(t *testing.T) {
	tests := map[string]TransferTestCase{
		"Transfer succeeds": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
				constants.Dave_Num0_599USD,
			},
			transfer: &constants.Transfer_Carl_Num0_Dave_Num0_Quote_500,
			expectedSubaccountBalance: map[satypes.SubaccountId]*big.Int{
				constants.Carl_Num0: big.NewInt(99_000_000),
				constants.Dave_Num0: big.NewInt(1_099_000_000),
			},
		},
		"Transfer succeeds - recipient does not exist": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			transfer: &constants.Transfer_Carl_Num0_Dave_Num0_Quote_500,
			expectedSubaccountBalance: map[satypes.SubaccountId]*big.Int{
				constants.Carl_Num0: big.NewInt(99_000_000),
				constants.Dave_Num0: big.NewInt(500_000_000),
			},
		},
		"Sender does not have sufficient balance": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
				constants.Dave_Num0_599USD,
			},
			transfer: &constants.Transfer_Carl_Num0_Dave_Num0_Quote_600,
			expectedSubaccountBalance: map[satypes.SubaccountId]*big.Int{
				constants.Carl_Num0: big.NewInt(599_000_000), // balance unchanged
				constants.Dave_Num0: big.NewInt(599_000_000), // balance unchanged
			},
			expectedErr: fmt.Sprintf(
				"Subaccount with id %v failed with UpdateResult: NewlyUndercollateralized: failed to apply subaccount updates",
				constants.Carl_Num0,
			),
		},
		"Sender is under collateralized": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_599USD,
			},
			transfer: &constants.Transfer_Carl_Num0_Dave_Num0_Quote_600,
			expectedSubaccountBalance: map[satypes.SubaccountId]*big.Int{
				constants.Carl_Num0: big.NewInt(100_000_000_000), // balance unchanged
				constants.Dave_Num0: big.NewInt(599_000_000),     // balance unchanged
			},
			expectedErr: fmt.Sprintf(
				"Subaccount with id %v failed with UpdateResult: NewlyUndercollateralized: failed to apply subaccount updates",
				constants.Carl_Num0,
			),
		},
		"Receiver is under collateralized (transfer still succeeds)": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			transfer: &constants.Transfer_Carl_Num0_Dave_Num0_Quote_500,
			expectedSubaccountBalance: map[satypes.SubaccountId]*big.Int{
				constants.Carl_Num0: big.NewInt(99_000_000),
				constants.Dave_Num0: big.NewInt(50_500_000_000),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessTransferTest(t, tc)
		})
	}
}

func TestProcessTransfer_CreateRecipientAccount(t *testing.T) {
	ks := keepertest.SendingKeepers(t)
	ks.Ctx = ks.Ctx.WithBlockHeight(5)
	keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

	perpetuals := []perptypes.Perpetual{
		constants.BtcUsd_100PercentMarginRequirement,
	}
	require.NoError(t, keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper))

	for _, p := range perpetuals {
		_, err := ks.PerpetualsKeeper.CreatePerpetual(
			ks.Ctx,
			p.Params.Id,
			p.Params.Ticker,
			p.Params.MarketId,
			p.Params.AtomicResolution,
			p.Params.DefaultFundingPpm,
			p.Params.LiquidityTier,
			p.Params.MarketType,
		)
		require.NoError(t, err)
	}
	ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, constants.Carl_Num0_599USD)
	ks.AccountKeeper.SetAccount(
		ks.Ctx,
		ks.AccountKeeper.NewAccountWithAddress(ks.Ctx, constants.Carl_Num0.MustGetAccAddress()),
	)

	// Create a sample recipient address.
	recipient := sample.AccAddress()
	recipientAddr, err := sdk.AccAddressFromBech32(recipient)
	require.NoError(t, err)

	// Verify that the recipient account does not exist.
	require.False(t, ks.AccountKeeper.HasAccount(ks.Ctx, recipientAddr))

	// Process the transfer.
	transfer := types.Transfer{
		Sender: constants.Carl_Num0,
		Recipient: satypes.SubaccountId{
			Owner:  recipient,
			Number: uint32(0),
		},
		AssetId: assettypes.AssetUsdc.Id,
		Amount:  500_000_000, // $500
	}
	err = ks.SendingKeeper.ProcessTransfer(ks.Ctx, &transfer)
	require.NoError(t, err)

	// The account should've been created for the recipient address.
	require.True(t, ks.AccountKeeper.HasAccount(ks.Ctx, recipientAddr))
}

func TestProcessDepositToSubaccount(t *testing.T) {
	testError := errors.New("error")

	tests := map[string]struct {
		msg                 types.MsgDepositToSubaccount
		expectedErr         error
		expectedErrContains string
		shouldPanic         bool
		setUpMocks          func(mckCall *mock.Call)
	}{
		"Success": {
			msg:         constants.MsgDepositToSubaccount_Alice_To_Carl_Num0_750,
			expectedErr: nil,
			setUpMocks: func(mckCall *mock.Call) {
				mckCall.Return(nil)
			},
		},
		"Propagate error": {
			msg:         constants.MsgDepositToSubaccount_Alice_To_Carl_Num0_750,
			expectedErr: testError,
			setUpMocks: func(mckCall *mock.Call) {
				mckCall.Return(testError)
			},
		},
		"Propagate panic": {
			msg:         constants.MsgDepositToSubaccount_Alice_To_Carl_Num0_750,
			expectedErr: testError,
			shouldPanic: true,
			setUpMocks: func(mckCall *mock.Call) {
				mckCall.Panic(testError.Error())
			},
		},
		"Bad sender address string": {
			msg: types.MsgDepositToSubaccount{
				Sender:    "1234567", // bad address string
				Recipient: constants.Alice_Num0,
				AssetId:   assettypes.AssetUsdc.Id,
				Quantums:  750_000_000,
			},
			expectedErrContains: "decoding bech32 failed",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := tc.msg
			// Create mock subaccounts keeper.
			mockSubaccountsKeeper := &mocks.SubaccountsKeeper{}
			// Create sending keeper with mock subaccounts keeper.
			ks := keepertest.SendingKeepersWithSubaccountsKeeper(t, mockSubaccountsKeeper)
			// Set up mock calls.
			if tc.setUpMocks != nil {
				mockCall := mockSubaccountsKeeper.On(
					"DepositFundsFromAccountToSubaccount",
					ks.Ctx,
					sdk.MustAccAddressFromBech32(msg.Sender),
					msg.Recipient,
					msg.AssetId,
					new(big.Int).SetUint64(msg.Quantums),
				)
				tc.setUpMocks(mockCall)
			}

			if tc.shouldPanic {
				require.PanicsWithValue(t, tc.expectedErr.Error(), func() {
					//nolint:errcheck
					ks.SendingKeeper.ProcessDepositToSubaccount(ks.Ctx, &msg)
				})
			} else {
				err := ks.SendingKeeper.ProcessDepositToSubaccount(ks.Ctx, &msg)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else if tc.expectedErrContains != "" {
					require.Contains(t, err.Error(), tc.expectedErrContains)
				} else {
					require.NoError(t, err)

					// Verify that corresponding indexer deposit event was emitted.
					assertDepositEventInIndexerBlock(t, ks.Ctx, ks.SendingKeeper, &msg)
				}
			}
		})
	}
}

func TestProcessWithdrawFromSubaccount(t *testing.T) {
	testError := errors.New("error")

	tests := map[string]struct {
		msg                 types.MsgWithdrawFromSubaccount
		expectedErr         error
		expectedErrContains string
		shouldPanic         bool
		setUpMocks          func(mckCall *mock.Call)
	}{
		"Success": {
			msg:         constants.MsgWithdrawFromSubaccount_Carl_Num0_To_Alice_750,
			expectedErr: nil,
			setUpMocks: func(mckCall *mock.Call) {
				mckCall.Return(nil)
			},
		},
		"Propagate error": {
			msg:         constants.MsgWithdrawFromSubaccount_Carl_Num0_To_Alice_750,
			expectedErr: testError,
			setUpMocks: func(mckCall *mock.Call) {
				mckCall.Return(testError)
			},
		},
		"Propagate panic": {
			msg:         constants.MsgWithdrawFromSubaccount_Carl_Num0_To_Alice_750,
			expectedErr: testError,
			shouldPanic: true,
			setUpMocks: func(mckCall *mock.Call) {
				mckCall.Panic(testError.Error())
			},
		},
		"Bad recipient address string": {
			msg: types.MsgWithdrawFromSubaccount{
				Sender:    constants.Alice_Num0,
				Recipient: "1234567", // bad address string
				AssetId:   assettypes.AssetUsdc.Id,
				Quantums:  750_000_000,
			},
			expectedErrContains: "decoding bech32 failed",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := tc.msg
			// Create mock subaccounts keeper.
			mockSubaccountsKeeper := &mocks.SubaccountsKeeper{}
			// Create sending keeper with mock subaccounts keeper.
			ks := keepertest.SendingKeepersWithSubaccountsKeeper(t, mockSubaccountsKeeper)
			// Set up mock calls.
			if tc.setUpMocks != nil {
				mockCall := mockSubaccountsKeeper.On(
					"WithdrawFundsFromSubaccountToAccount",
					ks.Ctx,
					msg.Sender,
					sdk.MustAccAddressFromBech32(msg.Recipient),
					msg.AssetId,
					new(big.Int).SetUint64(msg.Quantums),
				)
				tc.setUpMocks(mockCall)
			}

			if tc.shouldPanic {
				require.PanicsWithValue(t, tc.expectedErr.Error(), func() {
					//nolint:errcheck
					ks.SendingKeeper.ProcessWithdrawFromSubaccount(ks.Ctx, &msg)
				})
			} else {
				err := ks.SendingKeeper.ProcessWithdrawFromSubaccount(ks.Ctx, &msg)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else if tc.expectedErrContains != "" {
					require.Contains(t, err.Error(), tc.expectedErrContains)
				} else {
					require.NoError(t, err)

					// Verify that corresponding indexer withdraw event was emitted.
					assertWithdrawEventInIndexerBlock(t, ks.Ctx, ks.SendingKeeper, &msg)
				}
			}
		})
	}
}

func TestSendFromModuleToAccount(t *testing.T) {
	testDenom := "TestSendFromModuleToAccount/Coin"
	testModuleName := "bridge"

	tests := map[string]struct {
		// Setup.
		initialModuleBalance int64
		balanceToSend        int64
		recipientAddress     string

		// Expectations
		expectedErrContains string
	}{
		"Success - send to user account": {
			initialModuleBalance: 1000,
			balanceToSend:        100,
			recipientAddress:     constants.AliceAccAddress.String(),
		},
		"Success - send to module account": {
			initialModuleBalance: 1000,
			balanceToSend:        100,
			recipientAddress:     authtypes.NewModuleAddress("community_treasury").String(),
		},
		"Success - send to self": {
			initialModuleBalance: 1000,
			balanceToSend:        100,
			recipientAddress:     authtypes.NewModuleAddress(testModuleName).String(),
		},
		"Success - send 0 amount": {
			initialModuleBalance: 700,
			balanceToSend:        0,
			recipientAddress:     authtypes.NewModuleAddress(testModuleName).String(),
		},
		"Error - insufficient fund": {
			initialModuleBalance: 100,
			balanceToSend:        101,
			recipientAddress:     constants.AliceAccAddress.String(),
			expectedErrContains:  "insufficient funds",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initiate keepers and fund module with initial balance.
			ks := keepertest.SendingKeepers(t)
			ctx, sendingKeeper, bankKeeper, accountKeeper := ks.Ctx, ks.SendingKeeper, ks.BankKeeper, ks.AccountKeeper
			err := bankKeeper.MintCoins(
				ctx,
				testModuleName,
				sdk.NewCoins(sdk.NewCoin(testDenom, sdkmath.NewInt(int64(tc.initialModuleBalance)))),
			)
			require.NoError(t, err)
			startingModuleBalance := bankKeeper.GetBalance(
				ctx,
				accountKeeper.GetModuleAddress(testModuleName),
				testDenom,
			)
			startingRecipientBalance := bankKeeper.GetBalance(
				ctx,
				sdk.MustAccAddressFromBech32(tc.recipientAddress),
				testDenom,
			)

			// Send coins from module to account.
			err = sendingKeeper.SendFromModuleToAccount(
				ctx,
				&types.MsgSendFromModuleToAccount{
					Authority:        lib.GovModuleAddress.String(),
					SenderModuleName: testModuleName,
					Recipient:        tc.recipientAddress,
					Coin:             sdk.NewCoin(testDenom, sdkmath.NewInt(int64(tc.balanceToSend))),
				},
			)

			// Verify ending balances and error.
			endingModuleBalance := bankKeeper.GetBalance(
				ctx,
				accountKeeper.GetModuleAddress(testModuleName),
				testDenom,
			)
			endingRecipientBalance := bankKeeper.GetBalance(
				ctx,
				sdk.MustAccAddressFromBech32(tc.recipientAddress),
				testDenom,
			)
			if tc.expectedErrContains != "" { // if error should occur.
				// Verify that error is as expected.
				require.ErrorContains(t, err, tc.expectedErrContains)
				// Verify that module balance is unchanged.
				require.Equal(
					t,
					startingModuleBalance.Amount.Int64(),
					endingModuleBalance.Amount.Int64(),
				)
				// Verify that recipient balance is unchanged.
				require.Equal(
					t,
					startingRecipientBalance.Amount.Int64(),
					endingRecipientBalance.Amount.Int64(),
				)
			} else { // if send should succeed.
				// Verify that no error occurred.
				require.NoError(t, err)
				if tc.recipientAddress == authtypes.NewModuleAddress(testModuleName).String() {
					// If module sent to itself, verify that module balance is unchanged.
					require.Equal(t, startingModuleBalance, endingModuleBalance)
				} else {
					// Otherwise, verify that module balance is reduced by amount sent.
					require.Equal(
						t,
						startingModuleBalance.Amount.Int64()-tc.balanceToSend,
						endingModuleBalance.Amount.Int64(),
					)
					// Verify that recipient balance is increased by amount sent.
					require.Equal(
						t,
						startingRecipientBalance.Amount.Int64()+tc.balanceToSend,
						endingRecipientBalance.Amount.Int64(),
					)
				}
			}
		})
	}
}

func TestSendFromModuleToAccount_InvalidMsg(t *testing.T) {
	msgEmptySender := &types.MsgSendFromModuleToAccount{
		Authority:        lib.GovModuleAddress.String(),
		SenderModuleName: "",
		Recipient:        constants.AliceAccAddress.String(),
		Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(100)),
	}

	ks := keepertest.SendingKeepers(t)
	err := ks.SendingKeeper.SendFromModuleToAccount(ks.Ctx, msgEmptySender)
	require.ErrorContains(t, err, "Module name is empty")
}

func TestSendFromModuleToAccount_NonExistentSenderModule(t *testing.T) {
	msgNonExistentSender := &types.MsgSendFromModuleToAccount{
		Authority:        lib.GovModuleAddress.String(),
		SenderModuleName: "nonexistent",
		Recipient:        constants.AliceAccAddress.String(),
		Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(100)),
	}

	// Calling SendFromModuleToAccount with a non-existent sender module will panic.
	defer func() {
		if r := recover(); r != nil {
			require.ErrorContains(t, r.(error), "module account nonexistent does not exist")
		}
	}()
	ks := keepertest.SendingKeepers(t)
	err := ks.SendingKeeper.SendFromModuleToAccount(ks.Ctx, msgNonExistentSender)
	require.NoError(t, err) // this line is never reached, just here for lint check.
}

func TestSendFromModuleToAccount_InvalidRecipient(t *testing.T) {
	ks := keepertest.SendingKeepers(t)
	err := ks.SendingKeeper.SendFromModuleToAccount(
		ks.Ctx,
		&types.MsgSendFromModuleToAccount{
			Authority:        lib.GovModuleAddress.String(),
			SenderModuleName: "bridge",
			Recipient:        "dydx1abc", // invalid recipient address
			Coin:             sdk.NewCoin("dv4tnt", sdkmath.NewInt(1)),
		},
	)
	require.ErrorContains(t, err, "Account address is invalid")
}

func TestSendFromAccountToAccount(t *testing.T) {
	testDenom := "TestSendFromAccountToAccount/Coin"
	testModuleName := "bridge"
	tests := map[string]struct {
		// Setup.
		initialSenderBalance int64
		balanceToSend        int64
		senderAddress        string
		recipientAddress     string
		// Expectations
		expectedErrContains string
	}{
		"Success - send between user accounts": {
			initialSenderBalance: 1000,
			balanceToSend:        100,
			senderAddress:        constants.AliceAccAddress.String(),
			recipientAddress:     constants.BobAccAddress.String(),
		},
		"Success - send to module account": {
			initialSenderBalance: 1000,
			balanceToSend:        100,
			senderAddress:        constants.AliceAccAddress.String(),
			recipientAddress:     authtypes.NewModuleAddress("community_treasury").String(),
		},
		"Success - send to self": {
			initialSenderBalance: 1000,
			balanceToSend:        100,
			senderAddress:        constants.AliceAccAddress.String(),
			recipientAddress:     constants.AliceAccAddress.String(),
		},
		"Success - send 0 amount": {
			initialSenderBalance: 700,
			balanceToSend:        0,
			senderAddress:        constants.AliceAccAddress.String(),
			recipientAddress:     constants.BobAccAddress.String(),
		},
		"Error - insufficient funds": {
			initialSenderBalance: 100,
			balanceToSend:        101,
			senderAddress:        constants.AliceAccAddress.String(),
			recipientAddress:     constants.BobAccAddress.String(),
			expectedErrContains:  "insufficient funds",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initiate keepers and fund sender with initial balance.
			ks := keepertest.SendingKeepers(t)
			ctx, sendingKeeper, bankKeeper := ks.Ctx, ks.SendingKeeper, ks.BankKeeper

			senderAddr := sdk.MustAccAddressFromBech32(tc.senderAddress)

			// Mint coins to bridging module and transfer to sender account
			err := bankKeeper.MintCoins(
				ctx,
				testModuleName,
				sdk.NewCoins(sdk.NewCoin(testDenom, sdkmath.NewInt(tc.initialSenderBalance))),
			)
			require.NoError(t, err)
			err = bankKeeper.SendCoinsFromModuleToAccount(
				ctx,
				testModuleName,
				senderAddr,
				sdk.NewCoins(sdk.NewCoin(testDenom, sdkmath.NewInt(tc.initialSenderBalance))),
			)
			require.NoError(t, err)

			startingSenderBalance := bankKeeper.GetBalance(
				ctx,
				senderAddr,
				testDenom,
			)
			startingRecipientBalance := bankKeeper.GetBalance(
				ctx,
				sdk.MustAccAddressFromBech32(tc.recipientAddress),
				testDenom,
			)

			// Send coins from account to account.
			err = sendingKeeper.SendFromAccountToAccount(
				ctx,
				&types.MsgSendFromAccountToAccount{
					Authority: lib.GovModuleAddress.String(),
					Sender:    tc.senderAddress,
					Recipient: tc.recipientAddress,
					Coin:      sdk.NewCoin(testDenom, sdkmath.NewInt(tc.balanceToSend)),
				},
			)

			// Verify ending balances and error.
			endingSenderBalance := bankKeeper.GetBalance(
				ctx,
				senderAddr,
				testDenom,
			)
			endingRecipientBalance := bankKeeper.GetBalance(
				ctx,
				sdk.MustAccAddressFromBech32(tc.recipientAddress),
				testDenom,
			)

			if tc.expectedErrContains != "" { // if error should occur.
				// Verify that error is as expected.
				require.ErrorContains(t, err, tc.expectedErrContains)
				// Verify that sender balance is unchanged.
				require.Equal(
					t,
					startingSenderBalance.Amount.Int64(),
					endingSenderBalance.Amount.Int64(),
				)
				// Verify that recipient balance is unchanged.
				require.Equal(
					t,
					startingRecipientBalance.Amount.Int64(),
					endingRecipientBalance.Amount.Int64(),
				)
			} else { // if send should succeed.
				// Verify that no error occurred.
				require.NoError(t, err)
				if tc.senderAddress == tc.recipientAddress {
					// If account sent to itself, verify that balance is unchanged.
					require.Equal(t, startingSenderBalance, endingSenderBalance)
				} else {
					// Otherwise, verify that sender balance is reduced by amount sent.
					require.Equal(
						t,
						startingSenderBalance.Amount.Int64()-tc.balanceToSend,
						endingSenderBalance.Amount.Int64(),
					)
					// Verify that recipient balance is increased by amount sent.
					require.Equal(
						t,
						startingRecipientBalance.Amount.Int64()+tc.balanceToSend,
						endingRecipientBalance.Amount.Int64(),
					)
				}
			}
		})
	}
}
