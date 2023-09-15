package keeper_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/common"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
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
		bytes := indexer_manager.GetBytesFromEventData(event.Data)
		unmarshaler := common.UnmarshalerImpl{}
		var transfer indexerevents.TransferEventV1
		err := unmarshaler.Unmarshal(bytes, &transfer)
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
		bytes := indexer_manager.GetBytesFromEventData(event.Data)
		unmarshaler := common.UnmarshalerImpl{}
		var deposit indexerevents.TransferEventV1
		err := unmarshaler.Unmarshal(bytes, &deposit)
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
		bytes := indexer_manager.GetBytesFromEventData(event.Data)
		unmarshaler := common.UnmarshalerImpl{}
		var withdraw indexerevents.TransferEventV1
		err := unmarshaler.Unmarshal(bytes, &withdraw)
		if err != nil {
			panic(err)
		}
		withdraws = append(withdraws, &withdraw)
	}
	require.Contains(t, withdraws, expectedEvent)
}

func runProcessTransferTest(t *testing.T, tc TransferTestCase) {
	ctx, keeper, accountKeeper, pricesKeeper, perpKeeper, aKeeper, saKeeper, _ := keepertest.SendingKeepers(t)
	ctx = ctx.WithBlockHeight(5)
	keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ctx, perpKeeper)

	perpetuals := []perptypes.Perpetual{
		constants.BtcUsd_100PercentMarginRequirement,
	}
	require.NoError(t, keepertest.CreateUsdcAsset(ctx, aKeeper))

	for _, p := range perpetuals {
		_, err := perpKeeper.CreatePerpetual(
			ctx,
			p.Params.Id,
			p.Params.Ticker,
			p.Params.MarketId,
			p.Params.AtomicResolution,
			p.Params.DefaultFundingPpm,
			p.Params.LiquidityTier,
		)
		require.NoError(t, err)
	}

	for _, s := range tc.subaccounts {
		saKeeper.SetSubaccount(
			ctx,
			s,
		)
		accountKeeper.SetAccount(
			ctx,
			accountKeeper.NewAccountWithAddress(ctx, s.GetId().MustGetAccAddress()),
		)
	}

	err := keeper.ProcessTransfer(ctx, tc.transfer)
	for subaccountId, expectedQuoteBalance := range tc.expectedSubaccountBalance {
		subaccount := saKeeper.GetSubaccount(ctx, subaccountId)
		require.Equal(t, expectedQuoteBalance, subaccount.GetUsdcPosition())
	}
	if tc.expectedErr != "" {
		require.ErrorContains(t, err, tc.expectedErr)
	} else {
		require.NoError(t, err)
		assertTransferEventInIndexerBlock(
			t,
			keeper,
			ctx,
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
	ctx, keeper, accountKeeper, pricesKeeper, perpKeeper, aKeeper, saKeeper, _ := keepertest.SendingKeepers(t)
	ctx = ctx.WithBlockHeight(5)
	keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ctx, perpKeeper)

	perpetuals := []perptypes.Perpetual{
		constants.BtcUsd_100PercentMarginRequirement,
	}
	require.NoError(t, keepertest.CreateUsdcAsset(ctx, aKeeper))

	for _, p := range perpetuals {
		_, err := perpKeeper.CreatePerpetual(
			ctx,
			p.Params.Id,
			p.Params.Ticker,
			p.Params.MarketId,
			p.Params.AtomicResolution,
			p.Params.DefaultFundingPpm,
			p.Params.LiquidityTier,
		)
		require.NoError(t, err)
	}
	saKeeper.SetSubaccount(ctx, constants.Carl_Num0_599USD)
	accountKeeper.SetAccount(
		ctx,
		accountKeeper.NewAccountWithAddress(ctx, constants.Carl_Num0.MustGetAccAddress()),
	)

	// Create a sample recipient address.
	recipient := sample.AccAddress()
	recipientAddr, err := sdk.AccAddressFromBech32(recipient)
	require.NoError(t, err)

	// Verify that the recipient account does not exist.
	require.False(t, accountKeeper.HasAccount(ctx, recipientAddr))

	// Process the transfer.
	transfer := types.Transfer{
		Sender: constants.Carl_Num0,
		Recipient: satypes.SubaccountId{
			Owner:  recipient,
			Number: uint32(0),
		},
		AssetId: lib.UsdcAssetId,
		Amount:  500_000_000, // $500
	}
	err = keeper.ProcessTransfer(ctx, &transfer)
	require.NoError(t, err)

	// The account should've been created for the recipient address.
	require.True(t, accountKeeper.HasAccount(ctx, recipientAddr))
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
				AssetId:   lib.UsdcAssetId,
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
			ctx, keeper, _, _, _, _, _, _ :=
				keepertest.SendingKeepersWithSubaccountsKeeper(t, mockSubaccountsKeeper)
			// Set up mock calls.
			if tc.setUpMocks != nil {
				mockCall := mockSubaccountsKeeper.On(
					"DepositFundsFromAccountToSubaccount",
					ctx,
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
					keeper.ProcessDepositToSubaccount(ctx, &msg)
				})
			} else {
				err := keeper.ProcessDepositToSubaccount(ctx, &msg)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else if tc.expectedErrContains != "" {
					require.Contains(t, err.Error(), tc.expectedErrContains)
				} else {
					require.NoError(t, err)

					// Verify that corresponding indexer deposit event was emitted.
					assertDepositEventInIndexerBlock(t, ctx, keeper, &msg)
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
				AssetId:   lib.UsdcAssetId,
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
			ctx, keeper, _, _, _, _, _, _ :=
				keepertest.SendingKeepersWithSubaccountsKeeper(t, mockSubaccountsKeeper)
			// Set up mock calls.
			if tc.setUpMocks != nil {
				mockCall := mockSubaccountsKeeper.On(
					"WithdrawFundsFromSubaccountToAccount",
					ctx,
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
					keeper.ProcessWithdrawFromSubaccount(ctx, &msg)
				})
			} else {
				err := keeper.ProcessWithdrawFromSubaccount(ctx, &msg)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else if tc.expectedErrContains != "" {
					require.Contains(t, err.Error(), tc.expectedErrContains)
				} else {
					require.NoError(t, err)

					// Verify that corresponding indexer withdraw event was emitted.
					assertWithdrawEventInIndexerBlock(t, ctx, keeper, &msg)
				}
			}
		})
	}
}
