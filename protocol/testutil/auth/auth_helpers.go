package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func CreateTestModuleAccount(
	ctx sdk.Context,
	accountKeeper *authkeeper.AccountKeeper,
	moduleName string,
	permissions []string,
) {
	modBaseAcc := authtypes.NewBaseAccount(authtypes.NewModuleAddress(moduleName), nil, 0, 0)
	modAcc := authtypes.NewModuleAccount(modBaseAcc, moduleName, permissions...)
	accountKeeper.SetModuleAccount(ctx, modAcc)
}
