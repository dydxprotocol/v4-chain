package constants

import (
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	AliceAccAddress       = sdk.AccAddress(AlicePrivateKey.PubKey().Address())
	BobAccAddress         = sdk.AccAddress(BobPrivateKey.PubKey().Address())
	CarlAccAddress        = sdk.AccAddress(CarlPrivateKey.PubKey().Address())
	DaveAccAddress        = sdk.AccAddress(DavePrivateKey.PubKey().Address())
	AliceValAddress       = sdk.ValAddress(AlicePrivateKey.PubKey().Address())
	BobValAddress         = sdk.ValAddress(BobPrivateKey.PubKey().Address())
	CarlValAddress        = sdk.ValAddress(CarlPrivateKey.PubKey().Address())
	DaveValAddress        = sdk.ValAddress(DavePrivateKey.PubKey().Address())
	AliceConsAddress, _   = sdk.ConsAddressFromBech32("klyravalcons1zf9csp5ygq95cqyxh48w3qkuckmpealrhxq0ye")
	BobConsAddress, _     = sdk.ConsAddressFromBech32("klyravalcons1s7wykslt83kayxuaktep9fw8qxe5n73up9h3xp")
	CarlConsAddress, _    = sdk.ConsAddressFromBech32("klyravalcons1vy0nrh7l4rtezrsakaadz4mngwlpdmhyretgwy")
	DaveConsAddress, _    = sdk.ConsAddressFromBech32("klyravalcons1stjspktkshgcsv8sneqk2vs2ws0nw2wrnjkt6l")
	AliceValidatorAddress = sdk.ValAddress(AlicePrivateKey.PubKey().Address())
	BobValidatorAddress   = sdk.ValAddress(BobPrivateKey.PubKey().Address())
	CarlValidatorAddress  = sdk.ValAddress(CarlPrivateKey.PubKey().Address())
	DaveValidatorAddress  = sdk.ValAddress(DavePrivateKey.PubKey().Address())
	AliceAddressBz        = AlicePrivateKey.PubKey().Address().Bytes()
	BobAddressBz          = BobPrivateKey.PubKey().Address().Bytes()
	CarlAddressBz         = CarlPrivateKey.PubKey().Address().Bytes()
	DaveAddressBz         = DavePrivateKey.PubKey().Address().Bytes()

	// Collateral pool addresses for isolated perpetuals.
	IsoCollateralPoolAddress  = authtypes.NewModuleAddress(satypes.ModuleName + ":3")
	Iso2CollateralPoolAddress = authtypes.NewModuleAddress(satypes.ModuleName + ":4")
)
