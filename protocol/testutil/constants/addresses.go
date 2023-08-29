package constants

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	GovAccAddress = sdk.AccAddress(GovPrivateKey.PubKey().Address())

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
)
