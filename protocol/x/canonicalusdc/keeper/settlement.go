package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) MintLogicalAndSend(ctx context.Context, receiver sdk.AccAddress, coins sdk.Coins) error {
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, coins)
}

func (k Keeper) CancelBackingSwap(ctx sdk.Context, controller string, amount math.Int) error {
	if k.GetControls(ctx).Mode != types.Mode_MODE_GRADUAL {
		return types.ErrDisabled
	}
	participant, found := k.GetParticipant(ctx, controller)
	if !found {
		return types.ErrUnauthorized
	}
	funded, _ := types.ParseAmount(participant.Funded)
	if funded.LT(amount) {
		return types.ErrInsufficientBacking
	}
	ledger := k.GetLedger(ctx)
	var err error
	ledger.RestrictedFunding, err = subtract(ledger.RestrictedFunding, amount)
	if err != nil {
		return err
	}
	participant.Funded = funded.Sub(amount).String()
	receiver, err := sdk.AccAddressFromBech32(controller)
	if err != nil {
		return err
	}
	controls := k.GetControls(ctx)
	coins := sdk.NewCoins(sdk.NewCoin(controls.InjectiveDenom, amount))
	cacheCtx, write := ctx.CacheContext()
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(cacheCtx, types.ModuleName, receiver, coins); err != nil {
		return err
	}
	if err := k.SetParticipant(cacheCtx, participant); err != nil {
		return err
	}
	if err := k.SetLedger(cacheCtx, ledger); err != nil {
		return err
	}
	write()
	return nil
}

func (k Keeper) Settle(ctx sdk.Context, pending types.PendingSettlement, success bool, outcome string) error {
	cacheCtx, write := ctx.CacheContext()
	ctx = cacheCtx
	amount, err := types.ParsePositiveAmount(pending.Amount)
	if err != nil {
		return err
	}
	ledger := k.GetLedger(ctx)
	controls := k.GetControls(ctx)
	switch pending.Route {
	case types.Route_ROUTE_INJECTIVE:
		ledger.PendingInjective, err = subtract(ledger.PendingInjective, amount)
		if err != nil {
			return err
		}
		if success {
			err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(pending.LogicalDenom, amount)))
		} else {
			injective, parseErr := types.ParseAmount(ledger.InjectiveBacking)
			if parseErr != nil {
				return parseErr
			}
			restricted, parseErr := types.ParseAmount(ledger.RestrictedFunding)
			if parseErr != nil {
				return parseErr
			}
			requiredAfterRefund := injective.Add(restricted).Add(amount)
			physicalBalance := k.bankKeeper.GetBalance(ctx, types.ModuleAddress, pending.PhysicalDenom).Amount
			if physicalBalance.LT(requiredAfterRefund) {
				return fmt.Errorf("%w: physical refund missing", types.ErrInsufficientBacking)
			}
			var sender sdk.AccAddress
			sender, err = sdk.AccAddressFromBech32(pending.Sender)
			if err == nil {
				err = k.bankKeeper.SendCoinsFromModuleToAccount(
					ctx,
					types.ModuleName,
					sender,
					sdk.NewCoins(sdk.NewCoin(pending.LogicalDenom, amount)),
				)
			}
			ledger.InjectiveBacking = add(ledger.InjectiveBacking, amount)
		}
	case types.Route_ROUTE_NOBLE:
		ledger.PendingNoble, err = subtract(ledger.PendingNoble, amount)
		if !success {
			ledger.NobleBacking = add(ledger.NobleBacking, amount)
		}
	case types.Route_ROUTE_BACKING_SWAP:
		ledger.PendingBackingSwap, err = subtract(ledger.PendingBackingSwap, amount)
		if err != nil {
			return err
		}
		participant, found := k.GetParticipant(ctx, pending.Controller)
		if !found {
			return types.ErrInvalidParticipant
		}
		if success {
			ledger.InjectiveBacking = add(ledger.InjectiveBacking, amount)
			released, _ := types.ParseAmount(participant.Released)
			participant.Released = released.Add(amount).String()
		} else {
			err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(controls.LogicalDenom, amount)))
			ledger.NobleBacking = add(ledger.NobleBacking, amount)
			ledger.RestrictedFunding = add(ledger.RestrictedFunding, amount)
			funded, _ := types.ParseAmount(participant.Funded)
			participant.Funded = funded.Add(amount).String()
		}
		if err == nil {
			err = k.SetParticipant(ctx, participant)
		}
	default:
		return types.ErrSettlementNotFound
	}
	if err != nil {
		return err
	}
	if err := k.SetLedger(ctx, ledger); err != nil {
		return err
	}
	if !k.RemovePending(ctx, pending) {
		return types.ErrSettlementNotFound
	}
	k.markCompleted(ctx, pending.SourceChannel, pending.Sequence)
	if reason, broken := BackingInvariant(k)(ctx); broken {
		return fmt.Errorf("%w: %s", types.ErrInvalidLedger, reason)
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSettlement,
		sdk.NewAttribute(types.AttributeKeyRoute, pending.Route.String()),
		sdk.NewAttribute(types.AttributeKeySourceChannel, pending.SourceChannel),
		sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprint(pending.Sequence)),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, pending.Amount),
		sdk.NewAttribute(types.AttributeKeyOutcome, outcome),
	))
	write()
	return nil
}
