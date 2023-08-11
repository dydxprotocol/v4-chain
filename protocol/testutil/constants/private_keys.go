package constants

import (
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
