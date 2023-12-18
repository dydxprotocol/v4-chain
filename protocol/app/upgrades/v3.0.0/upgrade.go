package v_3_0_0

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	bridgemoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	clobmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	rewardsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vestmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
)

var (
	// List of module accounts to check in state.
	// These include all dYdX custom module accounts.
	ModuleAccsToInitialize = []string{
		bridgemoduletypes.ModuleName,
		satypes.ModuleName,
		clobmoduletypes.InsuranceFundName,
		rewardsmoduletypes.TreasuryAccountName,
		rewardsmoduletypes.VesterAccountName,
		vestmoduletypes.CommunityTreasuryAccountName,
		vestmoduletypes.CommunityVesterAccountName,
	}
)

// This module account initialization logic is copied from v2.0.0 upgrade handler.
// Testnet is to be upgraded from v1.0.1 to v3.0.0 directly and needs this logic to fix some accounts.
func InitializeModuleAccs(ctx sdk.Context, ak authkeeper.AccountKeeper) {
	for _, modAccName := range ModuleAccsToInitialize {
		// Get module account and relevant permissions from the accountKeeper.
		addr, perms := ak.GetModuleAddressAndPermissions(modAccName)
		if addr == nil {
			panic(fmt.Sprintf(
				"Did not find %v in `ak.GetModuleAddressAndPermissions`. This is not expected. Skipping.",
				modAccName,
			))
		}

		// Try to get the account in state.
		acc := ak.GetAccount(ctx, addr)
		if acc != nil {
			// Account has been initialized.
			macc, isModuleAccount := acc.(authtypes.ModuleAccountI)
			if isModuleAccount {
				// Module account was correctly initialized. Skipping
				ctx.Logger().Info(fmt.Sprintf(
					"module account %+v was correctly initialized. No-op",
					macc,
				))
				continue
			}
			// Module account has been initialized as a BaseAccount. Change to module account.
			// Note: We need to get the base account to retrieve its account number, and convert it
			// in place into a module account.
			baseAccount, ok := acc.(*authtypes.BaseAccount)
			if !ok {
				panic(fmt.Sprintf(
					"cannot cast %v into a BaseAccount, acc = %+v",
					modAccName,
					acc,
				))
			}
			newModuleAccount := authtypes.NewModuleAccount(
				baseAccount,
				modAccName,
				perms...,
			)
			ak.SetModuleAccount(ctx, newModuleAccount)
			ctx.Logger().Info(fmt.Sprintf(
				"Successfully converted %v to module account in state: %+v",
				modAccName,
				newModuleAccount,
			))
			continue
		}

		// Account has not been initialized at all. Initialize it as module.
		// Implementation taken from
		// https://github.com/dydxprotocol/cosmos-sdk/blob/bdf96fdd/x/auth/keeper/keeper.go#L213
		newModuleAccount := authtypes.NewEmptyModuleAccount(modAccName, perms...)
		maccI := (ak.NewAccount(ctx, newModuleAccount)).(authtypes.ModuleAccountI) // this set the account number
		ak.SetModuleAccount(ctx, maccI)
		ctx.Logger().Info(fmt.Sprintf(
			"Successfully initialized module account in state: %+v",
			newModuleAccount,
		))
	}
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ak authkeeper.AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running %s Upgrade...", UpgradeName)
		InitializeModuleAccs(ctx, ak)

		// TODO(CORE-824): Initialize ratelimit module params to desired state.
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
