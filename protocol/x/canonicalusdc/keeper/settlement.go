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

// Settle applies the ack or timeout outcome of one pending settlement. It only
// checks the ledger fields and balances it touches.
func (k Keeper) Settle(ctx sdk.Context, pending types.PendingSettlement, success bool, outcome string) error {
	cacheCtx, write := ctx.CacheContext()
	amount, err := types.PositiveAmountOf(pending.Amount)
	if err != nil {
		return err
	}
	ledger := k.GetLedger(cacheCtx)
	switch pending.Route {
	case types.Route_ROUTE_INJECTIVE:
		ledger, err = k.settleInjective(cacheCtx, ledger, pending, amount, success)
	case types.Route_ROUTE_NOBLE:
		ledger, err = settleNoble(ledger, amount, success)
	default:
		return types.ErrSettlementNotFound
	}
	if err != nil {
		return err
	}
	if err := k.SetLedger(cacheCtx, ledger); err != nil {
		return err
	}
	if !k.RemovePending(cacheCtx, pending) {
		return types.ErrSettlementNotFound
	}
	cacheCtx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSettlement,
		sdk.NewAttribute(types.AttributeKeyRoute, pending.Route.String()),
		sdk.NewAttribute(types.AttributeKeySourceChannel, pending.SourceChannel),
		sdk.NewAttribute(types.AttributeKeySequence, fmt.Sprint(pending.Sequence)),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, amount.String()),
		sdk.NewAttribute(types.AttributeKeyOutcome, outcome),
	))
	write()
	return nil
}

// settleInjective burns the locked logical USDC on success. On failure the
// inner transfer has already refunded the physical coins to the module, so the
// logical USDC is unlocked to the sender and the rate-limit charge is undone.
func (k Keeper) settleInjective(
	ctx sdk.Context,
	ledger types.Ledger,
	pending types.PendingSettlement,
	amount math.Int,
	success bool,
) (types.Ledger, error) {
	var err error
	if ledger.PendingInjective, err = subtract(ledger.PendingInjective, amount); err != nil {
		return ledger, err
	}
	logical := sdk.NewCoins(sdk.NewCoin(pending.LogicalDenom, amount))
	if success {
		return ledger, k.bankKeeper.BurnCoins(ctx, types.ModuleName, logical)
	}
	if ledger.InjectiveBacking, err = add(ledger.InjectiveBacking, amount); err != nil {
		return ledger, err
	}
	injective, err := types.AmountOf(ledger.InjectiveBacking)
	if err != nil {
		return ledger, err
	}
	physicalBalance := k.bankKeeper.GetBalance(ctx, types.ModuleAddress, pending.PhysicalDenom).Amount
	if physicalBalance.LT(injective) {
		return ledger, fmt.Errorf("%w: physical refund missing", types.ErrInsufficientBacking)
	}
	sender, err := sdk.AccAddressFromBech32(pending.Sender)
	if err != nil {
		return ledger, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, logical); err != nil {
		return ledger, err
	}
	k.ratelimitKeeper.UndoWithdrawal(ctx, pending.LogicalDenom, amount.BigInt())
	return ledger, nil
}

// settleNoble clears the in-flight amount. On failure the inner transfer has
// already refunded logical USDC to the sender, so it is Noble-backed again.
func settleNoble(ledger types.Ledger, amount math.Int, success bool) (types.Ledger, error) {
	var err error
	if ledger.PendingNoble, err = subtract(ledger.PendingNoble, amount); err != nil {
		return ledger, err
	}
	if !success {
		ledger.NobleBacking, err = add(ledger.NobleBacking, amount)
	}
	return ledger, err
}
