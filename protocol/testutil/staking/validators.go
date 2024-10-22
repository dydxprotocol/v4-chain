package validator_testutils

import (
	"cosmossdk.io/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func BuildTestValidator(name string, bondedTokens math.Int) stakingtypes.ValidatorI {
	switch name {
	case "alice":
		val := stakingtypes.Validator{
			Tokens:          bondedTokens,
			Status:          stakingtypes.Bonded,
			OperatorAddress: string(constants.AliceConsAddress),
		}
		return val
	case "bob":
		val := stakingtypes.Validator{
			Tokens:          bondedTokens,
			Status:          stakingtypes.Bonded,
			OperatorAddress: string(constants.BobConsAddress),
		}
		return val
	case "carl":
		val := stakingtypes.Validator{
			Tokens:          bondedTokens,
			Status:          stakingtypes.Bonded,
			OperatorAddress: string(constants.CarlConsAddress),
		}
		return val
	default:
		return stakingtypes.Validator{}
	}
}

func BuildAndMockTestValidator(
	ctx sdk.Context,
	name string,
	power math.Int,
	mValStore *mocks.ValidatorStore,
) stakingtypes.ValidatorI {
	val := BuildTestValidator(name, power)
	mValStore.On("ValidatorByConsAddr", ctx, GetConsAddressByName(name)).Return(val, nil)
	return val
}

func NewTotalBondedTokensMockReturn(
	ctx sdk.Context,
	names []string,
) *mocks.ValidatorStore {
	mValStore := &mocks.ValidatorStore{}
	for _, name := range names {
		BuildAndMockTestValidator(ctx, name, math.NewInt(500), mValStore)
	}
	mValStore.On("TotalBondedTokens", ctx).Return(math.NewInt(500*int64(len(names))), nil)
	return mValStore
}

func NewTotalBondedTokensValidatorMockReturnWithPowers(
	ctx sdk.Context,
	names []string,
	powers map[string]int64,
) *mocks.ValidatorStore {
	mValStore := &mocks.ValidatorStore{}
	totalPower := math.NewInt(0)
	for _, name := range names {
		BuildAndMockTestValidator(ctx, name, math.NewInt(powers[name]), mValStore)
		totalPower = totalPower.Add(math.NewInt(powers[name]))
	}
	mValStore.On("TotalBondedTokens", ctx).Return(totalPower, nil)
	return mValStore
}

func GetConsAddressByName(name string) sdk.ConsAddress {
	switch name {
	case "alice":
		return constants.AliceConsAddress
	case "bob":
		return constants.BobConsAddress
	case "carl":
		return constants.CarlConsAddress
	default:
		return nil
	}
}
