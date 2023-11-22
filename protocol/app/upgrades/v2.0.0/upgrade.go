package v_2_0_0

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	bridgemoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	rewardsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	vestmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"

	// Modules
	clobmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var (
	// // Module account permissions. Subset of `maccPerms` from `app.go`. `
	// MaccPerms = map[string][]string{
	// 	// -------- dYdX custom module accounts --------
	// 	// bridge module account mints tokens for bridged funds.
	// 	bridgemoduletypes.ModuleName: {authtypes.Minter},
	// 	// subaccounts module account holds tokens for all subaccounts.
	// 	satypes.ModuleName: nil,
	// 	// clob insurance fund account manages insurance fund for liquidations.
	// 	clobmoduletypes.InsuranceFundName: nil,
	// 	// rewards treasury account distribute funds trading accounts.
	// 	rewardsmoduletypes.TreasuryAccountName: nil,
	// 	// rewards vester account vest rewards tokens into the rewards treasury.
	// 	rewardsmoduletypes.VesterAccountName: nil,
	// 	// community treasury account holds funds for community use.
	// 	vestmoduletypes.CommunityTreasuryAccountName: nil,
	// 	// community vester account vests funds into the community treasury.
	// 	vestmoduletypes.CommunityVesterAccountName: nil,
	// }
	// List of module accounts to initialize in state.
	// These include all dYdX module accounts.
	// Used to deterministically iterate over MaccPerms above.
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

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ak authkeeper.AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running v2.0.0 Upgrade...")

		for _, modAccName := range ModuleAccsToInitialize {
			addr, perms := ak.GetModuleAddressAndPermissions(modAccName)
			if addr == nil {
				panic(fmt.Sprintf(
					"Did not find %v in `ak.GetModuleAddressAndPermissions`. This is not expected. Skipping.",
					modAccName,
				))
			}

			acc := ak.GetAccount(ctx, addr)
			if acc != nil {
				// Account has been initalized.
				macc, isModuleAccount := acc.(types.ModuleAccountI)
				if isModuleAccount {
					// Module account was correctly initialized. Skipping
					ctx.Logger().Info(fmt.Sprintf(
						"module account %+v was correctly initialized. No-op",
						macc,
					))
				} else {
					// Module account has been initalized as a BaseAccount. Change to module account.
					// Note: We need to get the base account to retrieve its account number, and convert it
					// in place into a module account.
					baseAccount, ok := acc.(*types.BaseAccount)
					if !ok {
						panic(fmt.Sprintf(
							"cannot cast %v into a BaseAccount",
							modAccName,
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
				}
			}

			// Account has not been initialized at all. Initialize it as module.
			// Implementation taken from
			// https://github.com/dydxprotocol/cosmos-sdk/blob/bdf96fdd/x/auth/keeper/keeper.go#L213
			newModuleAccount := authtypes.NewEmptyModuleAccount(modAccName, perms...)
			maccI := (ak.NewAccount(ctx, newModuleAccount)).(types.ModuleAccountI) // this set the account number
			ak.SetModuleAccount(ctx, maccI)
			ctx.Logger().Info(fmt.Sprintf(
				"Successfully initialized module account in state: %+v",
				newModuleAccount,
			))
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
