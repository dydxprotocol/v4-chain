package keeper

import (
	"math/big"
	"time"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	gometrics "github.com/hashicorp/go-metrics"
)

// ProcessTransfer transfers quote balance between two subaccounts.
func (k Keeper) ProcessTransfer(
	ctx sdk.Context,
	pendingTransfer *types.Transfer,
) (err error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.ProcessTransfer, metrics.Latency)

	err = k.subaccountsKeeper.TransferFundsFromSubaccountToSubaccount(
		ctx,
		pendingTransfer.Sender,
		pendingTransfer.Recipient,
		pendingTransfer.AssetId,
		pendingTransfer.GetBigQuantums(),
	)
	if err != nil {
		return err
	}

	recipientAddr := pendingTransfer.Recipient.MustGetAccAddress()

	// Create an account for the recipient address if one does not exist.
	// This is copied from https://sourcegraph.com/github.com/cosmos/cosmos-sdk/-/blob/x/bank/keeper/send.go?L199-203
	if exists := k.accountKeeper.HasAccount(ctx, recipientAddr); !exists {
		defer telemetry.IncrCounter(1, types.ModuleName, metrics.New, metrics.Account)
		k.accountKeeper.SetAccount(
			ctx,
			k.accountKeeper.NewAccountWithAddress(ctx, recipientAddr),
		)
	}

	// Add transfer event to Indexer block message.
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeTransfer,
		indexerevents.TransferEventVersion,
		indexer_manager.GetBytes(
			k.GenerateTransferEvent(pendingTransfer),
		),
	)

	return nil
}

// GenerateTransferEvent takes in a transfer and returns a transfer event.
func (k Keeper) GenerateTransferEvent(transfer *types.Transfer) *indexerevents.TransferEventV1 {
	return indexerevents.NewTransferEvent(
		satypes.SubaccountId{
			Owner:  transfer.Sender.Owner,
			Number: transfer.Sender.Number,
		},
		satypes.SubaccountId{
			Owner:  transfer.Recipient.Owner,
			Number: transfer.Recipient.Number,
		},
		transfer.AssetId,
		satypes.BaseQuantums(transfer.Amount),
	)
}

// ProcessDepositToSubaccount transfers quote balance from an account to a subaccount.
func (k Keeper) ProcessDepositToSubaccount(
	ctx sdk.Context,
	msgDepositToSubaccount *types.MsgDepositToSubaccount,
) (err error) {
	// Emit metric on latency.
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.ProcessDepositToSubaccount,
		metrics.Latency)

	// Convert sender address string to an sdk.AccAddress.
	senderAccAddress, err := sdk.AccAddressFromBech32(msgDepositToSubaccount.Sender)
	if err != nil {
		return err
	}

	// Invoke account-to-subaccount transfer keeper method in subaccounts.
	err = k.subaccountsKeeper.DepositFundsFromAccountToSubaccount(
		ctx,
		senderAccAddress,
		msgDepositToSubaccount.Recipient,
		msgDepositToSubaccount.AssetId,
		new(big.Int).SetUint64(msgDepositToSubaccount.Quantums),
	)

	// Emit gauge metric with labels if deposit to subaccount succeeds.
	if err == nil {
		metrics.EmitTelemetryWithLabelsForExecMode(
			ctx,
			// sdk.ExecModeFinalize is used here to ensure metrics are only emitted in the Finalize ExecMode.
			[]sdk.ExecMode{sdk.ExecModeFinalize},
			metrics.SetGaugeWithLabels,
			metrics.SendingProcessDepositToSubaccount,
			float32(msgDepositToSubaccount.Quantums),
			metrics.GetLabelForIntValue(metrics.AssetId, int(msgDepositToSubaccount.AssetId)),
		)

		// Add deposit event to Indexer block message.
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeTransfer,
			indexerevents.TransferEventVersion,
			indexer_manager.GetBytes(
				k.GenerateDepositEvent(msgDepositToSubaccount),
			),
		)
	}

	return err
}

// GenerateDepositEvent takes in a deposit and returns a deposit event.
func (k Keeper) GenerateDepositEvent(deposit *types.MsgDepositToSubaccount) *indexerevents.TransferEventV1 {
	return indexerevents.NewDepositEvent(
		deposit.Sender,
		satypes.SubaccountId{
			Owner:  deposit.Recipient.Owner,
			Number: deposit.Recipient.Number,
		},
		deposit.AssetId,
		satypes.BaseQuantums(deposit.Quantums),
	)
}

// ProcessWithdrawFromSubaccount transfers quote balance from a subaccount to an account.
func (k Keeper) ProcessWithdrawFromSubaccount(
	ctx sdk.Context,
	msgWithdrawFromSubaccount *types.MsgWithdrawFromSubaccount,
) (err error) {
	// Emit metric on latency.
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.ProcessWithdrawFromSubaccount,
		metrics.Latency)

	// Convert recipient address string to an sdk.AccAddress.
	recipientAccAddress, err := sdk.AccAddressFromBech32(msgWithdrawFromSubaccount.Recipient)
	if err != nil {
		return err
	}

	// Invoke subaccount-to-account transfer keeper method in subaccounts.
	err = k.subaccountsKeeper.WithdrawFundsFromSubaccountToAccount(
		ctx,
		msgWithdrawFromSubaccount.Sender,
		recipientAccAddress,
		msgWithdrawFromSubaccount.AssetId,
		new(big.Int).SetUint64(msgWithdrawFromSubaccount.Quantums),
	)

	// Emit gauge metric with labels if withdrawal from subaccount succeeds.
	if err == nil {
		telemetry.SetGaugeWithLabels(
			[]string{
				types.ModuleName,
				metrics.ProcessWithdrawFromSubaccount,
			},
			float32(msgWithdrawFromSubaccount.Quantums),
			[]gometrics.Label{
				metrics.GetLabelForIntValue(metrics.AssetId, int(msgWithdrawFromSubaccount.AssetId)),
			},
		)

		// Add withdraw event to Indexer block message.
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeTransfer,
			indexerevents.TransferEventVersion,
			indexer_manager.GetBytes(
				k.GenerateWithdrawEvent(msgWithdrawFromSubaccount),
			),
		)
	}

	return err
}

// GenerateWithdrawEvent takes in a withdrawal and returns a withdraw event.
func (k Keeper) GenerateWithdrawEvent(withdraw *types.MsgWithdrawFromSubaccount) *indexerevents.TransferEventV1 {
	return indexerevents.NewWithdrawEvent(
		satypes.SubaccountId{
			Owner:  withdraw.Sender.Owner,
			Number: withdraw.Sender.Number,
		},
		withdraw.Recipient,
		withdraw.AssetId,
		satypes.BaseQuantums(withdraw.Quantums),
	)
}

// SendFromModuleToAccount sends coins from a module account to an `x/bank` recipient.
func (k Keeper) SendFromModuleToAccount(
	ctx sdk.Context,
	msg *types.MsgSendFromModuleToAccount,
) (err error) {
	if err = msg.ValidateBasic(); err != nil {
		return err
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		msg.GetSenderModuleName(),
		sdk.MustAccAddressFromBech32(msg.GetRecipient()),
		sdk.NewCoins(msg.Coin),
	)
}

// SendFromAccountToAccount sends coins from one `x/bank` account to another `x/bank` account.
func (k Keeper) SendFromAccountToAccount(
	ctx sdk.Context,
	msg *types.MsgSendFromAccountToAccount,
) (err error) {
	senderAddr := sdk.MustAccAddressFromBech32(msg.GetSender())
	recipientAddr := sdk.MustAccAddressFromBech32(msg.GetRecipient())

	return k.bankKeeper.SendCoins(
		ctx,
		senderAddr,
		recipientAddr,
		sdk.NewCoins(msg.Coin),
	)
}
