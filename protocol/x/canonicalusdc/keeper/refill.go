package keeper

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	ratelimitutil "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Injective withdrawals are charged to the logical denom's x/ratelimit
// capacity, but an Injective deposit arrives as the physical voucher, which
// x/ratelimit does not credit back, so a withdraw-and-redeposit loop would
// drain the shared capacity for everyone. Refilling the whole deposit instead
// would let counterfeit Injective USDC buy exit capacity. A deposit therefore
// refills only capacity that Injective withdrawals spent and that is still
// missing from the bucket, capped at one period's baseline. Noble withdrawals
// never accrue: x/ratelimit credits Noble deposits itself.

// refillCap is the largest x/ratelimit baseline for the denom; zero when the
// denom has no limiters.
func (k Keeper) refillCap(ctx sdk.Context, denom string) math.Int {
	supply := k.bankKeeper.GetSupply(ctx, denom).Amount.BigInt()
	cap := new(big.Int)
	for _, limiter := range k.ratelimitKeeper.GetLimitParams(ctx, denom).Limiters {
		if baseline := ratelimitutil.GetBaseline(supply, limiter); baseline.Cmp(cap) > 0 {
			cap = baseline
		}
	}
	return math.NewIntFromBigInt(cap)
}

// missingCapacity is the smallest shortfall below its baseline across the
// denom's limiters, so a refill never lifts any of them above its baseline.
func (k Keeper) missingCapacity(ctx sdk.Context, denom string) math.Int {
	supply := k.bankKeeper.GetSupply(ctx, denom).Amount.BigInt()
	limiters := k.ratelimitKeeper.GetLimitParams(ctx, denom).Limiters
	capacities := k.ratelimitKeeper.GetDenomCapacity(ctx, denom).CapacityList
	if len(limiters) == 0 || len(capacities) != len(limiters) {
		return math.ZeroInt()
	}
	missing := math.ZeroInt()
	for i, limiter := range limiters {
		shortfall := math.NewIntFromBigInt(ratelimitutil.GetBaseline(supply, limiter)).
			Sub(math.NewIntFromBigInt(capacities[i].BigInt()))
		if shortfall.IsNegative() {
			shortfall = math.ZeroInt()
		}
		if i == 0 || shortfall.LT(missing) {
			missing = shortfall
		}
	}
	return missing
}

// addRefillable records an Injective withdrawal's charge on the ledger value,
// capped. The caller persists the ledger.
func (k Keeper) addRefillable(
	ctx sdk.Context,
	ledger types.Ledger,
	denom string,
	amount math.Int,
) (types.Ledger, error) {
	current, err := types.AmountOf(ledger.RefillableCapacity)
	if err != nil {
		return ledger, err
	}
	next := current.Add(amount)
	if cap := k.refillCap(ctx, denom); next.GT(cap) {
		next = cap
	}
	ledger.RefillableCapacity = types.NewAmount(next)
	return ledger, nil
}

// subtractRefillable lowers the refillable capacity for a charge that
// x/ratelimit undid, flooring at zero. The caller persists the ledger.
func subtractRefillable(ledger types.Ledger, amount math.Int) (types.Ledger, error) {
	current, err := types.AmountOf(ledger.RefillableCapacity)
	if err != nil {
		return ledger, err
	}
	next := current.Sub(amount)
	if next.IsNegative() {
		next = math.ZeroInt()
	}
	ledger.RefillableCapacity = types.NewAmount(next)
	return ledger, nil
}

// RefillCapacity restores logical withdrawal capacity for an Injective
// deposit: at most the deposit, at most what Injective withdrawals spent
// under the current cap, and at most what is still missing from the bucket.
// The deposit consumes its share of the refillable amount even when the
// bucket has already recovered on its own.
func (k Keeper) RefillCapacity(ctx sdk.Context, denom string, amount math.Int) error {
	ledger := k.GetLedger(ctx)
	refillable, err := types.AmountOf(ledger.RefillableCapacity)
	if err != nil {
		return err
	}
	refillable = math.MinInt(refillable, k.refillCap(ctx, denom))
	used := math.MinInt(amount, refillable)
	if !used.IsPositive() {
		return nil
	}
	if credit := math.MinInt(used, k.missingCapacity(ctx, denom)); credit.IsPositive() {
		k.ratelimitKeeper.IncrementCapacitiesForDenom(ctx, denom, credit.BigInt())
	}
	ledger.RefillableCapacity = types.NewAmount(refillable.Sub(used))
	return k.SetLedger(ctx, ledger)
}
