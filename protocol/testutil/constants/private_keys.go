package constants

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	AlicePrivateKey = privateKeyFromMnenomic(AliceMnenomic)
	BobPrivateKey   = privateKeyFromMnenomic(BobMnenomic)
	CarlPrivateKey  = privateKeyFromMnenomic(CarlMnenomic)
	DavePrivateKey  = privateKeyFromMnenomic(DaveMnenomic)

	privateKeyMap = map[string]cryptotypes.PrivKey{
		AliceAccAddress.String(): AlicePrivateKey,
		BobAccAddress.String():   BobPrivateKey,
		CarlAccAddress.String():  CarlPrivateKey,
		DaveAccAddress.String():  DavePrivateKey,
	}
)

func privateKeyFromMnenomic(mnenomic string) cryptotypes.PrivKey {
	kb := keyring.NewInMemory(TestEncodingCfg.Codec)
	_, err := kb.NewAccount("uid", mnenomic, "", sdk.GetConfig().GetFullBIP44Path(), hd.Secp256k1)
	if err != nil {
		panic(err)
	}
	armoredPvKey, err := kb.ExportPrivKeyArmor("uid", "")
	if err != nil {
		panic(err)
	}
	privKey, _, err := crypto.UnarmorDecryptPrivKey(armoredPvKey, "")
	if err != nil {
		panic(err)
	}
	return privKey
}

// GetPrivateKeyFromAddress returns the private key for the specified account address.
// Note that this panics if the account address is not one of the well known accounts.
func GetPrivateKeyFromAddress(accAddress string) cryptotypes.PrivKey {
	privKey, exists := privateKeyMap[accAddress]
	if !exists {
		panic(fmt.Errorf(
			"Unable to look-up private key, acc %s does not match any well known account.",
			accAddress))
	}
	return privKey
}
