package aggregator_test

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	valutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/staking"
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

func SetupTest(t *testing.T, vals []string, errorString string, initialSDAIPrice *big.Int) (sdk.Context, veaggregator.VoteAggregator) {
	// ctx, pk, _, _, _, mTimeProvider := keepertest.PricesKeepers(t)
	// mTimeProvider.On("Now").Return(constants.TimeT)

	ctx, _, pk, _, _, _, _, ratelimitKeeper, _, _ := keepertest.SubaccountsKeepers(t, true)

	if initialSDAIPrice != nil {
		ratelimitKeeper.SetSDAIPrice(ctx, initialSDAIPrice)
	}

	mTimeProvider := &mocks.TimeProvider{}
	mTimeProvider.On("Now").Return(constants.TimeT)
	keepertest.CreateTestMarkets(t, ctx, pk)

	mValStore := valutils.NewTotalBondedTokensMockReturn(ctx, vals)

	var pricesAggregatorFn voteweighted.PricesAggregateFn
	var conversionRateAggregatorFn voteweighted.ConversionRateAggregateFn

	if strings.Contains(errorString, "failed to aggregate prices") {
		pricesAggregatorFn = func(ctx sdk.Context, vePrices map[string]map[string]voteweighted.AggregatorPricePair) (map[string]voteweighted.AggregatorPricePair, error) {
			return nil, errors.New(errorString)
		}
	} else {
		pricesAggregatorFn = voteweighted.MedianPrices(
			ctx.Logger(),
			mValStore,
			voteweighted.DefaultPowerThreshold,
		)
	}

	if strings.Contains(errorString, "failed to aggregate sDai conversion rate") {
		conversionRateAggregatorFn = func(ctx sdk.Context, veConversionRates map[string]*big.Int) (*big.Int, error) {
			return nil, errors.New(errorString)
		}
	} else {
		conversionRateAggregatorFn = voteweighted.MedianConversionRate(
			ctx.Logger(),
			mValStore,
			voteweighted.DefaultPowerThreshold,
		)
	}

	handler := veaggregator.NewVeAggregator(
		ctx.Logger(),
		*pk,
		*ratelimitKeeper,
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
		initialSDAIPrice           *big.Int
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
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Single daemon data with conversion rate": {
			validators: []string{"alice"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, "1000000"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: big.NewInt(1000000),
			expectedError:              nil,
		},
		"Success: Multiple price updates, single validator, no conversion rate": {
			validators: []string{"alice"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, ""),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedError: nil,
		},
		"Success: Multiple price updates, single validator with conversion rate": {
			validators: []string{"alice"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		"Success: Single price update, from two validators, without conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Single price update, from two validators with different conversion rates": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		"Success: Single price update, from two validators with same conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000000"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		"Success: Multiple price updates, from two validators with no conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, ""),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates, from two validators with one conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates, from two validators with different conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Multiple price updates, from two validators with same conversion rate": {
			validators: []string{"alice", "bob"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		"Success: Single price update, from multiple validators, without conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Single price update, from multiple validators all conversion rates different": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Single price update, from multiple validators two out of three conversion rates the same": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Single price update, from multiple validators all conversion rates the same": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Multiple price updates, from multiple validators with all empty conversion rates": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPrices, ""),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates, from multiple validators with 2/3 conversion rates empty": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPrices, ""),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates, from multiple validators all conversion rates different": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Single price update from multiple validators with no conversion rate but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Single price update from multiple validators with conversion rate but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates from multiple validators with no conversion rate but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, ""),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates from multiple validators with conversion rate but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from exactly 2/3 validators with no conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, ""),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from exactly2/3 validators with conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
			},
			expectedPrices:             map[string]voteweighted.AggregatorPricePair{},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from multiple validators with no conversion rate but not enough voting power for some prices": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from multiple validators with conversion rate but not enough voting power for some prices and conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidSingleVEPrice, ""),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple prices from multiple validators with conversion rate but not enough voting power for some prices and enough for conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidSingleVEPrice, "1000000000000000000000000001"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: Continues when the validator's prices are malformed with no conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
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
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Continues when the validator's prices are malformed with some conversion rates, but not enough voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPricesWithOneInvalid, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
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
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Continues when the validator's prices are malformed with some conversion rates with exactly 2/3 voting power": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPricesWithOneInvalid, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPricesWithOneInvalid, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPricesWithOneInvalid, ""),
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
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Multiple price updates from >2/3 but not all validators with no conversion rates": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPrices, ""),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: nil,
			expectedError:              nil,
		},
		"Success: Markets don't exist and sDAI price is set": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPricesWithNoMarkets, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPricesWithNoMarkets, ""),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPricesWithNoMarkets, ""),
			},
			initialSDAIPrice: new(big.Int).SetUint64(50000),
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: new(big.Int).SetUint64(50000),
			expectedError:              nil,
		},
		"Success: Default PnL price to Spot price": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPricesOnlySpot, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPricesOnlySpot, "1000000000000000000000000002"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPricesOnlySpot, "1000000000000000000000000002"),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000002"),
			expectedError:              nil,
		},
		"Success: Multiple price updates from >2/3 but not all validators with some different conversion rates": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000002"),
			expectedError:              nil,
		},
		"Success: Multiple price updates from >2/3 but not all validators with all different conversion rates": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
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
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: No prices from multiple validators but all conversion rates valid": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, []vetypes.PricePair{}, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, []vetypes.PricePair{}, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, []vetypes.PricePair{}, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000001"),
			expectedError:              nil,
		},
		"Success: No prices from multiple validators but >2/3 conversion rates valid": {
			validators: []string{"alice", "bob", "carl", "dave"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, []vetypes.PricePair{}, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, []vetypes.PricePair{}, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, []vetypes.PricePair{}, "1000000000000000000000000002"),
			},
			expectedPrices: map[string]voteweighted.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.EthUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
				constants.SolUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.IsoUsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.FiveBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.FiveBillion),
				},
				constants.Iso2UsdPair: {
					SpotPrice: new(big.Int).SetUint64(constants.ThreeBillion),
					PnlPrice:  new(big.Int).SetUint64(constants.ThreeBillion),
				},
			},
			expectedSDaiConversionRate: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedError:              nil,
		},
		// Note: in the below tests, the failure stems from a mock aggregator function returning an error
		"Failure: Correctly returns error for a failure to aggregate prices": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
			},
			expectedPrices:             nil,
			expectedSDaiConversionRate: nil,
			expectedError:              errors.New("failed to aggregate prices"),
		},
		"Failure: Correctly returns error for a failure to aggregate sDai conversion rate": {
			validators: []string{"alice", "bob", "carl"},
			voteInfos: []cometabci.ExtendedVoteInfo{
				mustCreateSignedExtendedVoteInfo(t, constants.AliceConsAddress, constants.ValidVEPrices, "1000000000000000000000000000"),
				mustCreateSignedExtendedVoteInfo(t, constants.BobConsAddress, constants.ValidVEPrices, "1000000000000000000000000001"),
				mustCreateSignedExtendedVoteInfo(t, constants.CarlConsAddress, constants.ValidVEPrices, "1000000000000000000000000002"),
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
				ctx, handler = SetupTest(t, tc.validators, tc.expectedError.Error(), tc.initialSDAIPrice)
			} else {
				ctx, handler = SetupTest(t, tc.validators, "", tc.initialSDAIPrice)
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
