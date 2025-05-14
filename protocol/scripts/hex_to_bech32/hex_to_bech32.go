package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MinAddrLen = 20
	MaxAddrLen = 32
)

/*
main converts a hexadecimal address to bech32

Usage:

	go run scripts/hex_to_bech32/hex_to_bech32.go \
		-hex <bech32_address>
*/
func main() {
	// ------------ FLAGS ------------
	var hexAddress string
	flag.StringVar(&hexAddress, "hex", "0xda49e72c3577cc08bbe5b64d4a89e8e808e22f28", "hexadecimal address")
	flag.Parse()

	fmt.Println("Using the following configuration (modifiable via flags):")
	fmt.Println("hex address:", hexAddress)
	fmt.Println()

	// ------------ LOGIC ------------
	hexAddress = strings.TrimPrefix(hexAddress, "0x")

	bytes, err := hex.DecodeString(hexAddress)
	if err != nil {
		panic(err)
	}
	address := PadOrTruncateAddress(bytes)
	bech32Address := sdk.MustBech32ifyAddressBytes("dydx", address)

	// ------------ OUTPUT ------------
	fmt.Println(bech32Address)
}

func PadOrTruncateAddress(address []byte) []byte {
	if len(address) > MaxAddrLen {
		return address[:MaxAddrLen]
	} else if len(address) < MinAddrLen {
		return append(address, make([]byte, MinAddrLen-len(address))...)
	}
	return address
}
