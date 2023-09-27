package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

/*
main converts a bech32 address to hexadecimal by removing the HRP, separator, and checksum from the address.
This output can be used in an EVM bridge contract to emit the correct event to produce the bech32 address.

Usage:

	go run scripts/bech32_to_hex/bech32_to_hex.go \
		-address <bech32_address>
*/
func main() {
	// ------------ FLAGS ------------

	// Get flags.
	var bech32Address string
	flag.StringVar(&bech32Address, "address", "dydx1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5kmz6xt", "bech32 address")
	flag.Parse()

	// Print the flags used for the user
	fmt.Println("Using the following configuration (modifiable via flags):")
	fmt.Println("address:", bech32Address)
	fmt.Println()

	// ------------ LOGIC ------------

	hrp, decoded, err := bech32.DecodeAndConvert(bech32Address)
	if err != nil {
		log.Fatal(err)
	}

	// ------------ OUTPUT ------------

	fmt.Printf("HRP: %s\n", hrp)
	fmt.Printf(
		"Hexadecimal bytes (length %d): 0x%s\n",
		len(decoded),
		hex.EncodeToString(decoded),
	)
}
