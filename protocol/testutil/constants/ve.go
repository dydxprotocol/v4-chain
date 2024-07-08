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

	// AliceConsAddress signing constants.ValidVEPrice prices
	ValidSingleVoteExtInfoBytes []byte = []byte{
		18, 58, 10, 24, 10, 20, 41, 86, 4, 85, 205, 252, 249, 80,
		75, 230, 197, 194, 223, 67, 40, 173, 42, 88, 250, 242, 24,
		1, 26, 28, 10, 7, 8, 2, 18, 3, 2, 27, 95, 10, 7, 8, 1, 18,
		3, 2, 234, 102, 10, 8, 8, 0, 18, 4, 2, 7, 161, 37, 40, 2,
	}
	// AliceConsAddress and BobConsAddress signing constants.ValidVEPrice prices

	ValidMutliVoteExtInfoBytes []byte = []byte{
		18, 58, 10, 24, 10, 20, 41, 86, 4, 85, 205, 252,
		249, 80, 75, 230, 197, 194, 223, 67, 40, 173, 42,
		88, 250, 242, 24, 1, 26, 28, 10, 7, 8, 2, 18, 3, 2,
		27, 95, 10, 7, 8, 1, 18, 3, 2, 234, 102, 10, 8, 8,
		0, 18, 4, 2, 7, 161, 37, 40, 2, 18, 58, 10, 24, 10,
		20, 122, 77, 232, 19, 68, 115, 105, 12, 204, 221, 201,
		90, 226, 45, 39, 145, 179, 104, 173, 72, 24, 1, 26,
		28, 10, 8, 8, 0, 18, 4, 2, 7, 161, 37, 10, 7, 8, 2,
		18, 3, 2, 27, 95, 10, 7, 8, 1, 18, 3, 2, 234, 102, 40, 2,
	}
)
