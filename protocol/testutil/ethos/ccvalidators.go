package ethos_testutils

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
)

func BuildCCValidator(name string, power int64) ccvtypes.CrossChainValidator {
	switch name {
	case "alice":
		val, _ := ccvtypes.NewCCValidator(
			constants.AliceEthosAddressBz,
			power,
			constants.AliceEthosPubKey,
		)
		return val
	case "bob":
		val, _ := ccvtypes.NewCCValidator(
			constants.BobEthosAddressBz,
			power,
			constants.BobEthosPubKey,
		)
		return val
	case "carl":
		val, _ := ccvtypes.NewCCValidator(
			constants.CarlEthosAddressBz,
			power,
			constants.CarlEthosPubKey,
		)
		return val
	default:
		return ccvtypes.CrossChainValidator{}
	}
}

func BuildAndMockCCValidator(
	ctx sdk.Context,
	name string,
	power int64,
	mCCVStore *mocks.CCValidatorStore,
) ccvtypes.CrossChainValidator {
	val := BuildCCValidator(name, power)
	mCCVStore.On("GetCCValidator", ctx, val.Address).Return(val, true)
	return val
}
