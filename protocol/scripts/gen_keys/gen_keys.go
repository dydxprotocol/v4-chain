package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	// Set the desired Bech32 prefix for accounts
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("ethosvaloper", "ethosvaloperpub")

	// Initialize a codec that is required for Keyring
	registry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(registry)

	// Create an in-memory keyring with a codec
	kr, err := keyring.New("myKeyring", keyring.BackendMemory, "", os.Stdin, marshaler)
	if err != nil {
		log.Fatalf("Failed to create keyring: %v", err)
	}

	// The mnemonic
	mnemonic := "merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small"

	// HD path for the first account under BIP44: m/44'/118'/0'/0/0
	hdPath := hd.NewFundraiserParams(0, 118, 0).String()

	// Info holds the address and keys
	info, err := kr.NewAccount("myAccount", mnemonic, "", hdPath, hd.Secp256k1)
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}

	// Get the address in Bech32 format
	address, err := info.GetAddress()
	if err != nil {
		log.Fatalf("Failed to get address: %v", err)
	}

	fmt.Println("New Address:", address.String())
}
