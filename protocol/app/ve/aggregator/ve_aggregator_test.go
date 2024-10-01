package aggregator_test

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vemath "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	ethosutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ethos"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TODO: Create error test cases for all the functions separately.

var (
	voteCodec = vecodec.NewDefaultVoteExtensionCodec()
	extCodec  = vecodec.NewDefaultExtendedCommitCodec()
)

func SetupTest(t *testing.T, vals []string, errorString string) (sdk.Context, veaggregator.VoteAggregator) {
	ctx, pk, _, _, _, mTimeProvider := keepertest.PricesKeepers(t)
	mTimeProvider.On("Now").Return(constants.TimeT)

	keepertest.CreateTestMarkets(t, ctx, pk)

	mCCVStore := ethosutils.NewGetAllCCValidatorMockReturn(ctx, vals)

	var pricesAggregatorFn vemath.PricesAggregateFn
	var conversionRateAggregatorFn vemath.ConversionRateAggregateFn

	if strings.Contains(errorString, "failed to aggregate prices") {
		pricesAggregatorFn = func(ctx sdk.Context, vePrices map[string]map[string]vemath.AggregatorPricePair) (map[string]vemath.AggregatorPricePair, error) {
			return nil, fmt.Errorf(errorString)
		}
	} else {
		pricesAggregatorFn = voteweighted.MedianPrices(
			ctx.Logger(),
			mCCVStore,
			voteweighted.DefaultPowerThreshold,
		)
	}

	if strings.Contains(errorString, "failed to aggregate sDai conversion rate") {
		conversionRateAggregatorFn = func(ctx sdk.Context, veConversionRates map[string]*big.Int) (*big.Int, error) {
			return nil, fmt.Errorf(errorString)
		}
	} else {
		conversionRateAggregatorFn = voteweighted.MedianConversionRate(
			ctx.Logger(),
			mCCVStore,
			voteweighted.DefaultPowerThreshold,
		)
	}

	handler := veaggregator.NewVeAggregator(
		ctx.Logger(),
		*pk,
		pricesAggregatorFn,
		conversionRateAggregatorFn,
	)
	return ctx, handler
}

func mustCreateSignedExtendedVoteInfo(
	t *testing.T,
	consAddress sdk.ConsAddress,
	prices []vetypes.PricePair,
	sdaiConversionRate string,
) cometabci.ExtendedVoteInfo {
	voteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
		vetesting.NewDefaultSignedVeInfo(
			consAddress,
			prices,
			sdaiConversionRate,
		),
	)
	require.NoError(t, err)
	return voteInfo
}

func TestAggregateDaemonVEIntoFinalPricesAndConversionRate(t *testing.T) {
	tests := map[string]struct {
		validators                 []string
		voteInfos                  []cometabci.ExtendedVoteInfo
		expectedPrices             map[string]voteweighted.AggregatorPricePair
		expectedSDaiConversionRate *big.Int
		expectedError              error
	}{
		"Success: no daemon data": {
			validators:                 []string{"alice"},
			voteInfos:                  []cometabci.ExtendedVoteInfo{},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Single daemon data, empty conversion rate": {
			validators: []string{"alice"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Single daemon data with conversion rate": {
			validators: []string{"alice"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, "1000000"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: big.NewInt(1000000),
			expectedError:              nil,
		},
		"Success: Multiple price updates, single validator, no conversion rate": {
			validators: []string{"alice"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedError: nil,
		},
		"Success: Multiple price updates, single validator with conversion rate": {
			validators: []string{"alice"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		"Success: Single price update, from two validators, without conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Single price update, from two validators with different conversion rates": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		"Success: Single price update, from two validators with same conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000000"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		"Success: Multiple price updates, from two validators with no conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates, from two validators with one conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates, from two validators with different conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Multiple price updates, from two validators with same conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		"Success: Single price update, from multiple validators, without conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Single price update, from multiple validators all conversion rates different": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Single price update, from multiple validators two out of three conversion rates the same": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Single price update, from multiple validators all conversion rates the same": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Multiple price updates, from multiple validators with all empty conversion rates": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPrices, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates, from multiple validators with 2/3 conversion rates empty": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPrices, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates, from multiple validators all conversion rates different": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Single price update from multiple validators with no conversion rate but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Single price update from multiple validators with conversion rate but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates from multiple validators with no conversion rate but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, ""),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates from multiple validators with conversion rate but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from exactly 2/3 validators with no conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, ""),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from exactly2/3 validators with conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from multiple validators with no conversion rate but not enough voting power for some prices": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from multiple validators with conversion rate but not enough voting power for some prices and conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from multiple validators with conversion rate but not enough voting power for some prices and enough for conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Continues when the validator's prices are malformed with no conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Continues when the validator's prices are malformed with some conversion rates, but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPricesWithOneInvalid, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Continues when the validator's prices are malformed with some conversion rates with exactly 2/3 voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPricesWithOneInvalid, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPricesWithOneInvalid, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates from >2/3 but not all validators with no conversion rates": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPrices, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates from >2/3 but not all validators with some different conversion rates": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000002"),
			expectedError:              nil,
		},
		"Success: Multiple price updates from >2/3 but not all validators with all different conversion rates": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
				constants.SolUsdPair: {
					SpotPrice: constants.Price7Big,
					PnlPrice:  constants.Price7Big,
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: No prices from multiple validators but all conversion rates valid": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, []vetypes.PricePair{}, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, []vetypes.PricePair{}, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, []vetypes.PricePair{}, "1000000000000000000000000002"),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: No prices from multiple validators but >2/3 conversion rates valid": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, []vetypes.PricePair{}, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, []vetypes.PricePair{}, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, []vetypes.PricePair{}, "1000000000000000000000000002"),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		// Note: in the below tests, the failure stems from a mock aggregator function returning an error
		"Failure: Correctly returns error for a failure to aggregate prices": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
			},
			expectedPrices:             nil,
			expectedSDaiConversionRate: nil,
			expectedError:              errors.New("failed to aggregate prices"),
		},
		"Failure: Correctly returns error for a failure to aggregate sDai conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlEthosConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
			},
			expectedPrices:             nil,
			expectedSDaiConversionRate: nil,
			expectedError:              errors.New("failed to aggregate sDai conversion rate"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var ctx sdk.Context
			var handler veaggregator.VoteAggregator
			if tc.expectedError != nil {
				ctx, handler = SetupTest(t, tc.validators, tc.expectedError.Error())
			} else {
				ctx, handler = SetupTest(t, tc.validators, "")
			}

			_, commitBz, err := vetesting.CreateExtendedCommitInfo(tc.voteInfos)
			require.NoError(t, err)

			proposal := [][]byte{commitBz}
			votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
			require.NoError(t, err)

			prices, sDaiConversionRate, err := handler.AggregateDaemonVEIntoFinalPricesAndConversionRate(ctx, votes)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedPrices, prices)
			require.Equal(t, tc.expectedSDaiConversionRate, sDaiConversionRate)
		})
	}
}
