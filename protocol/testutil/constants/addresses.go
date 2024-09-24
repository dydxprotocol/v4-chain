package constants

import (
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
	AliceAddressBz   = AlicePrivateKey.PubKey().Address().Bytes()
	BobAddressBz     = BobPrivateKey.PubKey().Address().Bytes()
	CarlAddressBz    = CarlPrivateKey.PubKey().Address().Bytes()
	DaveAddressBz    = DavePrivateKey.PubKey().Address().Bytes()

	// Collateral pool addresses for isolated perpetuals.
	IsoCollateralPoolAddress  = authtypes.NewModuleAddress(satypes.ModuleName + ":3")
	Iso2CollateralPoolAddress = authtypes.NewModuleAddress(satypes.ModuleName + ":4")

	AliceEthosConsAddress = sdk.ConsAddress(AliceEthosPrivateKey.PubKey().Address())
	BobEthosConsAddress   = sdk.ConsAddress(BobEthosPrivateKey.PubKey().Address())
	CarlEthosConsAddress  = sdk.ConsAddress(CarlEthosPrivateKey.PubKey().Address())

	AliceEthosAddressBz = AliceEthosPrivateKey.PubKey().Address().Bytes()
	BobEthosAddressBz   = BobEthosPrivateKey.PubKey().Address().Bytes()
	CarlEthosAddressBz  = CarlEthosPrivateKey.PubKey().Address().Bytes()
)
