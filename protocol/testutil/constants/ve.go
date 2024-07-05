package constants

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
)

var (
	ValidEmptyExtInfoBytes        []byte = []byte{}
	ValidEmptyCrossChainValidator        = types.CrossChainValidator{}

	AliceCCValidator = types.CrossChainValidator{
		Address:  AliceConsAddress,
		Power:    1,
		Pubkey:   nil,
		OptedOut: false,
	}

	Val1 = sdk.ConsAddress("val1")
	Val2 = sdk.ConsAddress("val2")
	Val3 = sdk.ConsAddress("val3")
)
