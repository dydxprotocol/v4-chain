package v_3_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	bridgemoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	perpetualsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	rewardsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vestmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
)

var (
	ICAHostAllowMessages = []string{
		// IBC transfer messages
		sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}),

		// Bank messages
		sdk.MsgTypeURL(&banktypes.MsgSend{}),

		// Staking messages
		sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}),
		sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}),
		sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}),
		sdk.MsgTypeURL(&stakingtypes.MsgCancelUnbondingDelegation{}),

		// Distribution messages
		sdk.MsgTypeURL(&distrtypes.MsgSetWithdrawAddress{}),
		sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}),
		sdk.MsgTypeURL(&distrtypes.MsgFundCommunityPool{}),

		// Gov messages
		sdk.MsgTypeURL(&govtypesv1.MsgVote{}),
	}
	// List of module accounts to check in state.
	// These include all dYdX custom module accounts.
	ModuleAccsToInitialize = []string{
		bridgemoduletypes.ModuleName,
		satypes.ModuleName,
		perpetualsmoduletypes.InsuranceFundName,
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
			macc, isModuleAccount := acc.(sdk.ModuleAccountI)
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
		maccI := (ak.NewAccount(ctx, newModuleAccount)).(sdk.ModuleAccountI) // this set the account number
		ak.SetModuleAccount(ctx, maccI)
		ctx.Logger().Info(fmt.Sprintf(
			"Successfully initialized module account in state: %+v",
			newModuleAccount,
		))
	}
}

func IcaHostKeeperUpgradeHandler(
	ctx sdk.Context,
	vm module.VersionMap,
	mm *module.Manager,
) {
	icaAppModule, ok := mm.Modules[icatypes.ModuleName].(ica.AppModule)
	if !ok {
		panic("Modules[icatypes.ModuleName] is not of type ica.AppModule")
	}
	// Add Interchain Accounts host module
	// set the ICS27 consensus version so InitGenesis is not run
	// initialize ICS27 module
	vm[icatypes.ModuleName] = icaAppModule.ConsensusVersion()

	// controller submodule is not enabled.
	controllerParams := icacontrollertypes.Params{}

	// host submodule params
	hostParams := icahosttypes.Params{
		HostEnabled:   true,
		AllowMessages: ICAHostAllowMessages,
	}

	icaAppModule.InitModule(ctx, controllerParams, hostParams)
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	ak authkeeper.AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info("Running %s Upgrade...", UpgradeName)
		InitializeModuleAccs(sdkCtx, ak)

		// TODO(CORE-824): Initialize ratelimit module params to desired state.

		// TODO(CORE-848): Any unit test after a v3.0.0 upgrade test is added.
		IcaHostKeeperUpgradeHandler(sdkCtx, vm, mm)
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
