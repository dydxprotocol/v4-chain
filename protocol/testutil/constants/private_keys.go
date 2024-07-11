package constants

import (
	"encoding/base64"
	"fmt"

	"crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmosed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	AlicePrivateKey = privateKeyFromMnenomic(AliceMnenomic)
	BobPrivateKey   = privateKeyFromMnenomic(BobMnenomic)
	CarlPrivateKey  = privateKeyFromMnenomic(CarlMnenomic)
	DavePrivateKey  = privateKeyFromMnenomic(DaveMnenomic)

	AliceEthosPrivateKey = buildPrivKeyFromKeyString("TRJgf7lkTjs/sj43pyweEOanyV7H7fhnVivOi0A4yjW6NjXgCCilX3TshiA8CT/nHxz3brtLh9B/z2fJ4I9N6w==")
	BobEthosPrivateKey   = buildPrivKeyFromKeyString("OFR4w+FC6EMw5fAGTrHVexyPrjzQ7QfqgZOMgVf0izlCUb6Jh7oDJim9jXP1E0koJWUfXhD+pLPgSMZ0YKu7eg==")
	CarlEthosPrivateKey  = buildPrivKeyFromKeyString("3YaBAZLA+sl/E73lLfbFbG0u6DYm33ayr/0UpCt/vFBSLkZ/X6a1ZR0fy7fGWbN0ogP4Xc8rSx9dnvcZnqrqKw==")

	AlicePubKey = AlicePrivateKey.PubKey()
	BobPubKey   = BobPrivateKey.PubKey()
	CarlPubKey  = CarlPrivateKey.PubKey()
	DavePubKey  = DavePrivateKey.PubKey()

	AliceEthosPubKey = AliceEthosPrivateKey.PubKey()
	BobEthosPubKey   = BobEthosPrivateKey.PubKey()
	CarlEthosPubKey  = CarlEthosPrivateKey.PubKey()

	privateKeyMap = map[string]cryptotypes.PrivKey{
		AliceAccAddress.String(): AlicePrivateKey,
		BobAccAddress.String():   BobPrivateKey,
		CarlAccAddress.String():  CarlPrivateKey,
		DaveAccAddress.String():  DavePrivateKey,
	}

	privateConsMap = map[string]cryptotypes.PrivKey{
		AliceConsAddress.String():      AlicePrivateKey,
		BobConsAddress.String():        BobPrivateKey,
		CarlConsAddress.String():       CarlPrivateKey,
		DaveConsAddress.String():       DavePrivateKey,
		AliceEthosConsAddress.String(): AliceEthosPrivateKey,
		BobEthosConsAddress.String():   BobEthosPrivateKey,
		CarlEthosConsAddress.String():  CarlEthosPrivateKey,
	}
)

func buildPrivKeyFromKeyString(privKey string) cryptotypes.PrivKey {
	privKeyBytes, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		panic(fmt.Errorf("failed to decode private key: %w", err))
	}
	key := &cosmosed25519.PrivKey{Key: privKeyBytes[:ed25519.PrivateKeySize]}
	return key
}

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

func GetPublicKeyFromAddress(accAddress string) cryptotypes.PubKey {
	privKey, exists := privateKeyMap[accAddress]
	if !exists {
		panic(fmt.Errorf(
			"unable to look-up private key, acc %s does not match any well known account",
			accAddress))
	}
	return privKey.PubKey()
}

// GetPrivateKeyFromAddress returns the private key for the specified account address.
// Note that this panics if the account address is not one of the well known accounts.
func GetPrivateKeyFromAddress(accAddress string) cryptotypes.PrivKey {
	privKey, exists := privateKeyMap[accAddress]
	if !exists {
		panic(fmt.Errorf(
			"unable to look-up private key, acc %s does not match any well known account",
			accAddress))
	}
	return privKey
}

func GetPrivKeyFromConsAddress(consAddr sdk.ConsAddress) cryptotypes.PrivKey {
	privKey, exists := privateConsMap[consAddr.String()]
	if !exists {
		panic(fmt.Errorf(
			"unable to look-up private key, cons %s does not match any well known account",
			consAddr))
	}
	return privKey
}
