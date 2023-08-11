package constants

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	AliceAccAddress  = sdk.AccAddress(AlicePrivateKey.PubKey().Address())
	BobAccAddress    = sdk.AccAddress(BobPrivateKey.PubKey().Address())
	CarlAccAddress   = sdk.AccAddress(CarlPrivateKey.PubKey().Address())
	DaveAccAddress   = sdk.AccAddress(DavePrivateKey.PubKey().Address())
	AliceValAddress  = sdk.ValAddress(AlicePrivateKey.PubKey().Address())
	BobValAddress    = sdk.ValAddress(BobPrivateKey.PubKey().Address())
	CarlValAddress   = sdk.ValAddress(BobPrivateKey.PubKey().Address())
	DaveValAddress   = sdk.ValAddress(BobPrivateKey.PubKey().Address())
	AliceConsAddress = sdk.ConsAddress(AlicePrivateKey.PubKey().Address())
	BobConsAddress   = sdk.ConsAddress(BobPrivateKey.PubKey().Address())
	CarlConsAddress  = sdk.ConsAddress(CarlPrivateKey.PubKey().Address())
	DaveConsAddress  = sdk.ConsAddress(DavePrivateKey.PubKey().Address())
)
