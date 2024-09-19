package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// SweepMainVaultBankBalances deposits any usdc balance from the Megavault main vault bank balance
// into the Megavault main vault subaccount balance.
func (k Keeper) SweepMainVaultBankBalance(
	ctx sdk.Context,
) {
	mainVaultBalance := k.bankKeeper.GetBalance(
		ctx,
		types.MegavaultMainAddress,
		assettypes.AssetUsdc.Denom,
	)
	// No funds to sweep
	if mainVaultBalance.Amount.BigInt().Sign() <= 0 {
		return
	}

	err := k.subaccountsKeeper.DepositFundsFromAccountToSubaccount(
		ctx,
		types.MegavaultMainAddress,
		types.MegavaultMainSubaccount,
		assettypes.AssetUsdc.Id,
		mainVaultBalance.Amount.BigInt(),
	)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			"SweepMainVaultBankBalance: Failed to sweep funds from main vault bank balance to subaccount",
			err,
		)
		return
	}
}
