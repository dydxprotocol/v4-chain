package client

import (
	simappparams "cosmossdk.io/simapp/params"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/go-bip39"
)

func NewTestKeyring(
	encodingCfg simappparams.EncodingConfig,
	keyringDir string,
	accountName string,
) (keyring.Keyring, *keyring.Record) {
	kr, _ := keyring.New("chain", "test", keyringDir, nil, encodingCfg.Codec)
	algos, _ := kr.SupportedAlgorithms()
	algo, _ := keyring.NewSigningAlgoFromString("secp256k1", algos)
	entropySeed, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropySeed)
	info, _ := kr.NewAccount(accountName, mnemonic, "", "", algo)
	return kr, info
}

func PopulateFromFields(ctx client.Context, bech32Address string) client.Context {
	fromAddress, fromName, _, _ := client.GetFromFields(
		ctx,
		ctx.Keyring,
		bech32Address,
	)
	return ctx.
		WithFromAddress(fromAddress).
		WithFromName(fromName).
		WithFrom(fromName)
}
