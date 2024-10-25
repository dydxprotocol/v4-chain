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

	ValidSingleVoteExtInfoBytes []byte = []byte{
		18, 147, 1, 10, 25, 10, 20, 41, 86, 4, 85, 205,
		252, 249, 80, 75, 230, 197, 194, 223, 67, 40,
		173, 42, 88, 250, 242, 24, 244, 3, 26, 50,
		10, 14, 8, 1, 18, 10, 10, 3, 2, 234, 102, 18,
		3, 2, 234, 102, 10, 16, 8, 0, 18, 12, 10, 4,
		2, 7, 161, 37, 18, 4, 2, 7, 161, 37, 10, 14,
		8, 2, 18, 10, 10, 3, 2, 27, 95, 18, 3, 2, 27,
		95, 34, 64, 106, 234, 170, 223, 62, 29, 79, 10,
		140, 58, 179, 121, 19, 149, 214, 74, 233, 233,
		210, 126, 189, 40, 189, 45, 122, 125, 180, 118,
		224, 167, 49, 132, 9, 72, 145, 58, 191, 97, 204,
		108, 221, 172, 230, 85, 24, 214, 165, 144,
		60, 105, 254, 32, 15, 80, 97, 97, 221, 9, 53,
		91, 223, 253, 93, 82, 40, 2,
	}

	// return value of the following function call
	// (empty prices vote extension from 4 validators with 500000 power and block height 3)
	// validators, err := tApp.App.StakingKeeper.GetBondedValidatorsByPower(tApp.App.NewContextLegacy(true, tApp.header))
	// if err != nil {
	// 	tApp.builder.t.Fatalf("Failed to get bonded validators: %v", err)
	// }
	// localLastCommit := vetestutil.GetEmptyLocalLastCommit(
	// 	validators,
	// 	tApp.App.LastBlockHeight(),
	// 	0,
	// 	"localdydxprotocol",
	// )
	ValidMultiEmptyVoteExtInfoBytes []byte = []byte{
		0x12, 0x60, 0xa, 0x1a, 0xa, 0x14, 0x29, 0x56, 0x4, 0x55, 0xcd,
		0xfc, 0xf9, 0x50, 0x4b, 0xe6, 0xc5, 0xc2, 0xdf, 0x43, 0x28, 0xad,
		0x2a, 0x58, 0xfa, 0xf2, 0x18, 0xa0, 0xc2, 0x1e, 0x22, 0x40, 0xd9,
		0x57, 0xd3, 0x8f, 0xca, 0xaf, 0xe1, 0xe, 0xc0, 0x43, 0x80, 0x84,
		0xee, 0x59, 0x6b, 0x59, 0x89, 0x27, 0x26, 0xa4, 0x6f, 0xe4, 0x21,
		0x78, 0x4a, 0x63, 0x52, 0x6c, 0x2d, 0xe8, 0x80, 0xe8, 0x4f, 0x3e,
		0x60, 0xc8, 0x6b, 0x4a, 0x78, 0x66, 0x21, 0xac, 0x52, 0xf3, 0x58,
		0x80, 0xb5, 0xdf, 0xa1, 0x58, 0xbe, 0x38, 0x63, 0xb2, 0xe6, 0x73,
		0x42, 0xe6, 0x7f, 0x37, 0x3, 0x6e, 0x6e, 0x69, 0x28, 0x2, 0x12,
		0x60, 0xa, 0x1a, 0xa, 0x14, 0x4c, 0x91, 0xa1, 0x7, 0x4c, 0x61,
		0xd6, 0x57, 0x30, 0x95, 0xf8, 0x60, 0xf8, 0x8e, 0x9f, 0xaa, 0x3c,
		0xbb, 0x49, 0xc5, 0x18, 0xa0, 0xc2, 0x1e, 0x22, 0x40, 0x40, 0x2b,
		0xce, 0x1f, 0xc9, 0xfb, 0x54, 0xdc, 0xc1, 0xae, 0x33, 0xdb, 0x89,
		0x5a, 0x62, 0x92, 0xe0, 0xcb, 0x32, 0x6, 0x0, 0x4e, 0x70, 0xd5,
		0xbe, 0xa4, 0x50, 0xb5, 0xa4, 0x96, 0x5e, 0x71, 0x35, 0x35, 0xc,
		0x7b, 0x3e, 0x7d, 0xd5, 0x17, 0xd1, 0xe7, 0xbd, 0xa9, 0x47, 0x9c,
		0x4c, 0x8e, 0xc6, 0x1c, 0xd1, 0x84, 0x81, 0xa3, 0xde, 0x41, 0xa5,
		0x5e, 0xea, 0x2b, 0x76, 0xbe, 0x4a, 0x9e, 0x28, 0x2, 0x12, 0x60,
		0xa, 0x1a, 0xa, 0x14, 0x77, 0x79, 0x4d, 0xcb, 0xbe, 0x97, 0x84,
		0xde, 0xb1, 0x69, 0x1, 0x7f, 0xcf, 0x66, 0x59, 0x4b, 0xe9, 0x7d,
		0xe9, 0x29, 0x18, 0xa0, 0xc2, 0x1e, 0x22, 0x40, 0xdb, 0x1b, 0xa,
		0x2f, 0x48, 0x81, 0x1d, 0xf5, 0x3, 0xb8, 0xdf, 0x10, 0xc3, 0x90,
		0xd8, 0xbd, 0xff, 0xd0, 0x7, 0x58, 0x73, 0x2a, 0xf0, 0x41, 0x38,
		0x5e, 0xed, 0x1b, 0x69, 0x63, 0x9a, 0x63, 0x1b, 0x94, 0x89, 0x2c,
		0x71, 0x1a, 0xac, 0x39, 0x48, 0xc0, 0x1c, 0xa8, 0x38, 0xef, 0x62,
		0x4c, 0x4f, 0x43, 0x65, 0xa7, 0xa5, 0xcc, 0x5d, 0xbe, 0x20, 0xbb,
		0x9d, 0x8d, 0x1, 0x40, 0xc5, 0xbd, 0x28, 0x2, 0x12, 0x60, 0xa,
		0x1a, 0xa, 0x14, 0x7a, 0x4d, 0xe8, 0x13, 0x44, 0x73, 0x69, 0xc,
		0xcc, 0xdd, 0xc9, 0x5a, 0xe2, 0x2d, 0x27, 0x91, 0xb3, 0x68, 0xad,
		0x48, 0x18, 0xa0, 0xc2, 0x1e, 0x22, 0x40, 0xfb, 0x7e, 0x58, 0x4b,
		0xbe, 0x33, 0xa1, 0x18, 0xf7, 0xd2, 0xfb, 0x86, 0xd6, 0xb7, 0x9e,
		0x61, 0xa3, 0x31, 0x98, 0x5b, 0x5f, 0xee, 0x7f, 0xef, 0x72, 0xf1,
		0xcc, 0xc0, 0x5b, 0xf7, 0xc, 0x33, 0x50, 0x2c, 0x58, 0x84, 0x8b,
		0xff, 0xe1, 0x7b, 0xd8, 0x68, 0x33, 0xd7, 0xfd, 0xd1, 0xcc, 0x3e,
		0x77, 0xec, 0xff, 0x5d, 0x4e, 0x93, 0x89, 0xc6, 0xc, 0x15, 0x25,
		0x1f, 0x57, 0x22, 0xf0, 0xed, 0x28, 0x2,
	}
)
