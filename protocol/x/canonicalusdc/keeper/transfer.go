package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type transferAuthorizationKey struct{}

type TransferDecorator struct {
	keeper *Keeper
}

func NewTransferDecorator(keeper *Keeper) transfertypes.MsgServer {
	return TransferDecorator{keeper: keeper}
}

var _ transfertypes.MsgServer = TransferDecorator{}

func (d TransferDecorator) Transfer(
	goCtx context.Context,
	msg *transfertypes.MsgTransfer,
) (*transfertypes.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	controls := d.keeper.GetControls(ctx)
	if msg == nil || msg.Token.Denom != controls.LogicalDenom || controls.Mode == types.Mode_MODE_DISABLED {
		return d.keeper.TransferMsgServer().Transfer(goCtx, msg)
	}
	if controls.Mode == types.Mode_MODE_PAUSED {
		return nil, types.ErrPaused
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if msg.SourcePort != transfertypes.PortID {
		return nil, fmt.Errorf("%w: canonical transfers require the transfer port", types.ErrUnsupportedChannel)
	}
	amount := msg.Token.Amount
	if err := validateTransferAmount(controls, amount); err != nil {
		return nil, err
	}
	switch msg.SourceChannel {
	case controls.InjectiveChannel:
		return d.transferInjective(ctx, controls, msg, amount)
	case controls.NobleChannel:
		return d.transferNoble(ctx, controls, msg, amount)
	default:
		return nil, fmt.Errorf("%w: %s", types.ErrUnsupportedChannel, msg.SourceChannel)
	}
}

func (d TransferDecorator) UpdateParams(
	ctx context.Context,
	msg *transfertypes.MsgUpdateParams,
) (*transfertypes.MsgUpdateParamsResponse, error) {
	return d.keeper.TransferMsgServer().UpdateParams(ctx, msg)
}

func validateTransferAmount(controls types.Controls, amount math.Int) error {
	if !amount.IsPositive() {
		return types.ErrInvalidAmount
	}
	maxAmount, err := types.PositiveAmountOf(controls.MaxTransferAmount)
	if err != nil || amount.GT(maxAmount) {
		return types.ErrInvalidAmount
	}
	return nil
}

func (d TransferDecorator) transferInjective(
	ctx sdk.Context,
	controls types.Controls,
	msg *transfertypes.MsgTransfer,
	amount math.Int,
) (*transfertypes.MsgTransferResponse, error) {
	ledger := d.keeper.GetLedger(ctx)
	var err error
	if ledger.InjectiveBacking, err = subtract(ledger.InjectiveBacking, amount); err != nil {
		return nil, err
	}
	if ledger.PendingInjective, err = add(ledger.PendingInjective, amount); err != nil {
		return nil, err
	}
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	cacheCtx, write := ctx.CacheContext()
	// The physical packet carries the Injective denom, so x/ratelimit does not
	// charge it against logical USDC capacity.
	if err := d.keeper.ratelimitKeeper.ProcessWithdrawal(cacheCtx, controls.LogicalDenom, amount.BigInt()); err != nil {
		return nil, err
	}
	coins := sdk.NewCoins(sdk.NewCoin(controls.LogicalDenom, amount))
	if err := d.keeper.bankKeeper.SendCoinsFromAccountToModule(cacheCtx, sender, types.ModuleName, coins); err != nil {
		return nil, err
	}
	physical := *msg
	physical.Token = sdk.NewCoin(controls.InjectiveDenom, amount)
	physical.Sender = types.ModuleAddress.String()
	authorizedCtx := cacheCtx.WithValue(transferAuthorizationKey{}, true)
	response, err := d.keeper.TransferMsgServer().Transfer(authorizedCtx, &physical)
	if err != nil {
		return nil, err
	}
	pending := newPending(types.Route_ROUTE_INJECTIVE, controls, msg, response.Sequence, controls.InjectiveDenom, amount)
	if err := d.keeper.SetPending(cacheCtx, pending); err != nil {
		return nil, err
	}
	if err := d.keeper.SetLedger(cacheCtx, ledger); err != nil {
		return nil, err
	}
	emitLogicalTransfer(cacheCtx, msg)
	emitPending(cacheCtx, pending)
	write()
	return response, nil
}

func (d TransferDecorator) transferNoble(
	ctx sdk.Context,
	controls types.Controls,
	msg *transfertypes.MsgTransfer,
	amount math.Int,
) (*transfertypes.MsgTransferResponse, error) {
	if !controls.NobleWithdrawalsEnabled {
		return nil, types.ErrDisabled
	}
	ledger := d.keeper.GetLedger(ctx)
	var err error
	if ledger.NobleBacking, err = subtract(ledger.NobleBacking, amount); err != nil {
		return nil, err
	}
	if ledger.PendingNoble, err = add(ledger.PendingNoble, amount); err != nil {
		return nil, err
	}
	cacheCtx, write := ctx.CacheContext()
	authorizedCtx := cacheCtx.WithValue(transferAuthorizationKey{}, true)
	response, err := d.keeper.TransferMsgServer().Transfer(authorizedCtx, msg)
	if err != nil {
		return nil, err
	}
	pending := newPending(types.Route_ROUTE_NOBLE, controls, msg, response.Sequence, controls.LogicalDenom, amount)
	if err := d.keeper.SetPending(cacheCtx, pending); err != nil {
		return nil, err
	}
	if err := d.keeper.SetLedger(cacheCtx, ledger); err != nil {
		return nil, err
	}
	emitPending(cacheCtx, pending)
	write()
	return response, nil
}

func newPending(
	route types.Route,
	controls types.Controls,
	msg *transfertypes.MsgTransfer,
	sequence uint64,
	physicalDenom string,
	amount math.Int,
) types.PendingSettlement {
	return types.PendingSettlement{
		Route:            route,
		SourceChannel:    msg.SourceChannel,
		Sequence:         sequence,
		Sender:           msg.Sender,
		LogicalReceiver:  msg.Receiver,
		PhysicalReceiver: msg.Receiver,
		LogicalDenom:     controls.LogicalDenom,
		PhysicalDenom:    physicalDenom,
		Amount:           types.NewAmount(amount),
	}
}

func emitLogicalTransfer(ctx sdk.Context, msg *transfertypes.MsgTransfer) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		transfertypes.EventTypeTransfer,
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(transfertypes.AttributeKeyReceiver, msg.Receiver),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, msg.Token.Amount.String()),
		sdk.NewAttribute(transfertypes.AttributeKeyDenom, msg.Token.Denom),
		sdk.NewAttribute(transfertypes.AttributeKeyMemo, msg.Memo),
	))
}

func emitPending(ctx sdk.Context, pending types.PendingSettlement) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCanonicalTransfer,
		sdk.NewAttribute(types.AttributeKeyRoute, pending.Route.String()),
		sdk.NewAttribute(types.AttributeKeyLogicalDenom, pending.LogicalDenom),
		sdk.NewAttribute(types.AttributeKeyPhysicalDenom, pending.PhysicalDenom),
		sdk.NewAttribute(types.AttributeKeySourceChannel, pending.SourceChannel),
		sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprint(pending.Sequence)),
		sdk.NewAttribute(types.AttributeKeyPending, "true"),
	))
}
