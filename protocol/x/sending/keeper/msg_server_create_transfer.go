package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
)

func (k msgServer) CreateTransfer(
	goCtx context.Context,
	msg *types.MsgCreateTransfer,
) (*types.MsgCreateTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Process the transfer by applying subaccount updates.
	err := k.Keeper.ProcessTransfer(ctx, msg.Transfer)
	if err != nil {
		telemetry.IncrCounter(1, types.ModuleName, metrics.Transfer, metrics.Error)
		return nil, err
	}

	telemetry.IncrCounter(1, types.ModuleName, metrics.Transfer, metrics.Success)

	// emit create_transfer event
	ctx.EventManager().EmitEvent(
		types.NewCreateTransferEvent(
			msg.Transfer.Sender,
			msg.Transfer.Recipient,
			msg.Transfer.AssetId,
			msg.Transfer.Amount,
		),
	)

	return &types.MsgCreateTransferResponse{}, nil
}

// DepositToSubaccount initiates a transfer from sender (an `x/banks` account)
// to a recipient (an `x/subaccounts` subaccount).
func (k msgServer) DepositToSubaccount(
	goCtx context.Context,
	msg *types.MsgDepositToSubaccount,
) (*types.MsgDepositToSubaccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Process deposit from account to subaccount.
	err := k.Keeper.ProcessDepositToSubaccount(ctx, msg)
	if err != nil {
		telemetry.IncrCounter(1, types.ModuleName, metrics.ProcessDepositToSubaccount, metrics.Error)
		return nil, err
	}
	telemetry.IncrCounter(1, types.ModuleName, metrics.ProcessDepositToSubaccount, metrics.Success)

	// emit deposit_to_subaccount event
	ctx.EventManager().EmitEvent(
		types.NewDepositToSubaccountEvent(
			sdk.MustAccAddressFromBech32(msg.Sender),
			msg.Recipient,
			msg.AssetId,
			msg.Quantums,
		),
	)

	return &types.MsgDepositToSubaccountResponse{}, nil
}

// WithdrawFromSubaccount initiates a transfer from sender (an `x/subaccounts` subaccount)
// to a recipient (an `x/banks` account).
func (k msgServer) WithdrawFromSubaccount(
	goCtx context.Context,
	msg *types.MsgWithdrawFromSubaccount,
) (*types.MsgWithdrawFromSubaccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Process withdrawal from subaccount to account.
	err := k.Keeper.ProcessWithdrawFromSubaccount(ctx, msg)
	if err != nil {
		telemetry.IncrCounter(1, types.ModuleName, metrics.ProcessWithdrawFromSubaccount, metrics.Error)
		return nil, err
	}
	telemetry.IncrCounter(1, types.ModuleName, metrics.ProcessWithdrawFromSubaccount, metrics.Success)

	// emit withdraw_from_subaccount event
	ctx.EventManager().EmitEvent(
		types.NewWithdrawFromSubaccountEvent(
			msg.Sender,
			sdk.MustAccAddressFromBech32(msg.Recipient),
			msg.AssetId,
			msg.Quantums,
		),
	)

	return &types.MsgWithdrawFromSubaccountResponse{}, nil
}
