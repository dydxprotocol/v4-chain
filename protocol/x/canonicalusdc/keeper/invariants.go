package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type backingAmounts struct {
	noble            math.Int
	injective        math.Int
	legacyDownstream math.Int
	restricted       math.Int
	pendingInjective math.Int
	pendingNoble     math.Int
	pendingSwap      math.Int
}

func RegisterInvariants(ir sdk.InvariantRegistry, keeper Keeper) {
	ir.RegisterRoute(types.ModuleName, "backing", BackingInvariant(keeper))
}

func BackingInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		controls := keeper.GetControls(ctx)
		if controls.Mode == types.Mode_MODE_DISABLED {
			return "", false
		}
		amounts, err := parseBackingAmounts(keeper.GetLedger(ctx))
		if err != nil {
			return err.Error(), true
		}
		physicalRequired := amounts.injective.Add(amounts.restricted)
		physicalBalance := keeper.bankKeeper.GetBalance(ctx, types.ModuleAddress, controls.InjectiveDenom).Amount
		if physicalBalance.LT(physicalRequired) {
			return fmt.Sprintf(
				"canonical USDC physical reserve %s is below accounted reserve %s",
				physicalBalance,
				physicalRequired,
			), true
		}
		logicalBalance := keeper.bankKeeper.GetBalance(ctx, types.ModuleAddress, controls.LogicalDenom).Amount
		if logicalBalance.LT(amounts.pendingInjective) {
			return fmt.Sprintf(
				"canonical USDC locked logical balance %s is below pending liability %s",
				logicalBalance,
				amounts.pendingInjective,
			), true
		}
		logicalSupply := keeper.bankKeeper.GetSupply(ctx, controls.LogicalDenom).Amount
		accountedBacking := amounts.noble.Add(amounts.injective).Add(amounts.pendingInjective).Add(amounts.pendingSwap)
		if !logicalSupply.Equal(accountedBacking) {
			return fmt.Sprintf(
				"canonical USDC logical supply %s does not equal accounted backing %s",
				logicalSupply,
				accountedBacking,
			), true
		}
		if amounts.legacyDownstream.GT(logicalSupply) {
			return "canonical USDC legacy downstream amount exceeds logical supply", true
		}
		pending := keeper.GetAllPending(ctx)
		if uint32(len(pending)) != keeper.PendingCount(ctx) {
			return "canonical USDC pending count does not match stored settlements", true
		}
		if reason := checkPendingTotals(amounts, pending); reason != "" {
			return reason, true
		}
		if reason := checkParticipantFunding(keeper.GetAllParticipants(ctx), amounts.restricted); reason != "" {
			return reason, true
		}
		return "", false
	}
}

func parseBackingAmounts(ledger types.Ledger) (backingAmounts, error) {
	parsed := backingAmounts{}
	var err error
	parsed.noble, err = types.ParseAmount(ledger.NobleBacking)
	if err != nil {
		return backingAmounts{}, err
	}
	parsed.injective, err = types.ParseAmount(ledger.InjectiveBacking)
	if err != nil {
		return backingAmounts{}, err
	}
	parsed.legacyDownstream, err = types.ParseAmount(ledger.LegacyDownstream)
	if err != nil {
		return backingAmounts{}, err
	}
	parsed.restricted, err = types.ParseAmount(ledger.RestrictedFunding)
	if err != nil {
		return backingAmounts{}, err
	}
	parsed.pendingInjective, err = types.ParseAmount(ledger.PendingInjective)
	if err != nil {
		return backingAmounts{}, err
	}
	parsed.pendingNoble, err = types.ParseAmount(ledger.PendingNoble)
	if err != nil {
		return backingAmounts{}, err
	}
	parsed.pendingSwap, err = types.ParseAmount(ledger.PendingBackingSwap)
	if err != nil {
		return backingAmounts{}, err
	}
	return parsed, nil
}

func checkPendingTotals(amounts backingAmounts, pending []types.PendingSettlement) string {
	totals := map[types.Route]math.Int{
		types.Route_ROUTE_INJECTIVE:    math.ZeroInt(),
		types.Route_ROUTE_NOBLE:        math.ZeroInt(),
		types.Route_ROUTE_BACKING_SWAP: math.ZeroInt(),
	}
	for _, settlement := range pending {
		amount, err := types.ParsePositiveAmount(settlement.Amount)
		if err != nil {
			return "canonical USDC pending settlement has an invalid amount"
		}
		totals[settlement.Route] = totals[settlement.Route].Add(amount)
	}
	if !totals[types.Route_ROUTE_INJECTIVE].Equal(amounts.pendingInjective) ||
		!totals[types.Route_ROUTE_NOBLE].Equal(amounts.pendingNoble) ||
		!totals[types.Route_ROUTE_BACKING_SWAP].Equal(amounts.pendingSwap) {
		return "canonical USDC pending settlement totals do not match the backing ledger"
	}
	return ""
}

func checkParticipantFunding(participants []types.Participant, restricted math.Int) string {
	total := math.ZeroInt()
	for _, participant := range participants {
		funded, err := types.ParseAmount(participant.Funded)
		if err != nil {
			return "canonical USDC participant has an invalid funded amount"
		}
		total = total.Add(funded)
	}
	if !total.Equal(restricted) {
		return "canonical USDC participant funding does not match restricted backing"
	}
	return ""
}
