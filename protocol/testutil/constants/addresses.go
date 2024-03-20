package constants

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var (
	AliceAccAddress  = sdk.AccAddress(AlicePrivateKey.PubKey().Address())
	BobAccAddress    = sdk.AccAddress(BobPrivateKey.PubKey().Address())
	CarlAccAddress   = sdk.AccAddress(CarlPrivateKey.PubKey().Address())
	DaveAccAddress   = sdk.AccAddress(DavePrivateKey.PubKey().Address())
	AliceValAddress  = sdk.ValAddress(AlicePrivateKey.PubKey().Address())
	BobValAddress    = sdk.ValAddress(BobPrivateKey.PubKey().Address())
	CarlValAddress   = sdk.ValAddress(CarlPrivateKey.PubKey().Address())
	DaveValAddress   = sdk.ValAddress(DavePrivateKey.PubKey().Address())
	AliceConsAddress = sdk.ConsAddress(AlicePrivateKey.PubKey().Address())
	BobConsAddress   = sdk.ConsAddress(BobPrivateKey.PubKey().Address())
	CarlConsAddress  = sdk.ConsAddress(CarlPrivateKey.PubKey().Address())
	DaveConsAddress  = sdk.ConsAddress(DavePrivateKey.PubKey().Address())

	// Collateral pool addresses for isolated perpetuals.
	IsoCollateralPoolAddress  = authtypes.NewModuleAddress(satypes.ModuleName + ":3")
	Iso2CollateralPoolAddress = authtypes.NewModuleAddress(satypes.ModuleName + ":4")
)
