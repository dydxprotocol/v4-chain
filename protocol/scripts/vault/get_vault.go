package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	appconfig "github.com/dydxprotocol/v4-chain/protocol/app/config"
	vaultcli "github.com/dydxprotocol/v4-chain/protocol/x/vault/client/cli"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

/*
main derives a vault (in essence a subaccount) given its type and number.
This output can be used in a transaction to deposit into a vault.

Usage:

	go run scripts/vault/get_vault.go -type <vault_type> -number <vault_number>
*/
func main() {
	// ------------ FLAGS ------------
	// Get flags.
	var vaultTypeStr string
	var vaultNumberStr string
	flag.StringVar(&vaultTypeStr, "type", "clob", "vault type")
	flag.StringVar(&vaultNumberStr, "number", "0", "vault number")
	flag.Parse()

	// Convert vault type string to VaultType.
	vaultType, err := vaultcli.GetVaultTypeFromString(vaultTypeStr)
	if err != nil {
		log.Fatal(err)
	}

	// Convert vault number string to uint32.
	vaultNumberParsed, err := strconv.ParseUint(vaultNumberStr, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	vaultNumber := uint32(vaultNumberParsed)

	// Print the flags used.
	fmt.Println("Using the following configuration (modifiable via flags):")
	fmt.Println("type:", vaultType)
	fmt.Println("number:", vaultNumber)
	fmt.Println()

	// ------------ LOGIC ------------
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(appconfig.Bech32PrefixAccAddr, appconfig.Bech32PrefixAccPub)

	vaultId := vaulttypes.VaultId{
		Type:   vaultType,
		Number: vaultNumber,
	}
	if err != nil {
		log.Fatal(err)
	}

	vault := vaultId.ToSubaccountId()

	// ------------ OUTPUT ------------
	fmt.Printf("Vault:\n  Owner: %s\n  Number: %d\n", vault.Owner, vault.Number)
}
