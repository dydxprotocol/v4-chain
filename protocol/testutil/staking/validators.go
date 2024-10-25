package validator_testutils

import (
	"cosmossdk.io/math"
	lib "github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func BuildTestValidator(name string, bondedTokens math.Int) stakingtypes.ValidatorI {
	alicePubKey, _ := codectypes.NewAnyWithValue(constants.AlicePubKey)
	bobPubKey, _ := codectypes.NewAnyWithValue(constants.BobPubKey)
	carlPubKey, _ := codectypes.NewAnyWithValue(constants.CarlPubKey)
	switch name {
	case "alice":
		val := stakingtypes.Validator{
			Tokens:          bondedTokens,
			Status:          stakingtypes.Bonded,
			ConsensusPubkey: alicePubKey,
		}
		return val
	case "bob":
		val := stakingtypes.Validator{
			Tokens:          bondedTokens,
			Status:          stakingtypes.Bonded,
			ConsensusPubkey: bobPubKey,
		}
		return val
	case "carl":
		val := stakingtypes.Validator{
			Tokens:          bondedTokens,
			Status:          stakingtypes.Bonded,
			ConsensusPubkey: carlPubKey,
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
	mValStore.On("TotalBondedTokens", ctx).Return(ConvertPowerToTokens(500*int64(len(names))), nil)
	return mValStore
}

func NewTotalBondedTokensValidatorMockReturnWithPowers(
	ctx sdk.Context,
	names []string,
	powers map[string]int64,
) *mocks.ValidatorStore {
	mValStore := &mocks.ValidatorStore{}
	totalPower := int64(0)
	for _, name := range names {
		BuildAndMockTestValidator(ctx, name, math.NewInt(powers[name]), mValStore)
		totalPower += powers[name]
	}
	mValStore.On("TotalBondedTokens", ctx).Return(ConvertPowerToTokens(totalPower), nil)
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
	case "dave":
		return constants.DaveConsAddress
	default:
		return nil
	}
}

func ConvertPowerToTokens(power int64) math.Int {
	powerReduction := lib.PowerReduction
	return math.NewInt(power).Mul(powerReduction)
}
