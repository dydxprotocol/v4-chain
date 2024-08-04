package voteweighted_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/log"
	vemath "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	ethosutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ethos"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
)

type TestVoteWeightedMedianTC struct {
	validators      []string
	validatorPrices map[string]map[string]vemath.AggregatorPricePair
	expectedPrices  map[string]vemath.AggregatorPricePair
	powers          map[string]int64
}

func TestVoteWeightedMedian(t *testing.T) {
	tests := map[string]TestVoteWeightedMedianTC{
		"empty prices": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{},
			expectedPrices:  map[string]vemath.AggregatorPricePair{},
		},
		"single price single validator": {
			validators: []string{"alice"},
			powers: map[string]int64{
				"alice": 500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: constants.Price5Big,
						PnlPrice:  constants.Price5Big,
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
		},
		"single price multiple validators, varied prices": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5001),
						PnlPrice:  big.NewInt(5001),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5002),
						PnlPrice:  big.NewInt(5002),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(5001),
					PnlPrice:  big.NewInt(5001),
				},
			},
		},
		"single price two validators, varied prices": {
			validators: []string{"alice", "bob"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5001),
						PnlPrice:  big.NewInt(5001),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(5000),
					PnlPrice:  big.NewInt(5000),
				},
			},
		},
		"single price two validators, varied prices, one nil": {
			validators: []string{"alice", "bob"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: nil,
						PnlPrice:  nil,
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{},
		},
		"single price multiple validators, varied prices, one nil": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: nil,
						PnlPrice:  nil,
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5002),
						PnlPrice:  big.NewInt(5002),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{},
		},
		"single price multiple validators, varied prices, all nil": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: nil,
						PnlPrice:  nil,
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: nil,
						PnlPrice:  nil,
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: nil,
						PnlPrice:  nil,
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{},
		},
		"single price multiple validators, same price": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: constants.Price5Big,
						PnlPrice:  constants.Price5Big,
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: constants.Price5Big,
						PnlPrice:  constants.Price5Big,
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: constants.Price5Big,
						PnlPrice:  constants.Price5Big,
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
			},
		},
		"multiple prices, single validator": {
			validators: []string{"alice"},
			powers: map[string]int64{
				"alice": 500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: constants.Price5Big,
						PnlPrice:  constants.Price5Big,
					},
					constants.EthUsdPair: {
						SpotPrice: constants.Price6Big,
						PnlPrice:  constants.Price6Big,
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
			},
		},
		"multiple prices, multiple validators, same prices": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: constants.Price5Big,
						PnlPrice:  constants.Price5Big,
					},
					constants.EthUsdPair: {
						SpotPrice: constants.Price6Big,
						PnlPrice:  constants.Price6Big,
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: constants.Price5Big,
						PnlPrice:  constants.Price5Big,
					},
					constants.EthUsdPair: {
						SpotPrice: constants.Price6Big,
						PnlPrice:  constants.Price6Big,
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: constants.Price5Big,
						PnlPrice:  constants.Price5Big,
					},
					constants.EthUsdPair: {
						SpotPrice: constants.Price6Big,
						PnlPrice:  constants.Price6Big,
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: constants.Price5Big,
					PnlPrice:  constants.Price5Big,
				},
				constants.EthUsdPair: {
					SpotPrice: constants.Price6Big,
					PnlPrice:  constants.Price6Big,
				},
			},
		},
		"multiple prices, multiple validators, varied prices": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(6000),
						PnlPrice:  big.NewInt(6000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5001),
						PnlPrice:  big.NewInt(5001),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(6001),
						PnlPrice:  big.NewInt(6001),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5002),
						PnlPrice:  big.NewInt(5002),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(6002),
						PnlPrice:  big.NewInt(6002),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(5001),
					PnlPrice:  big.NewInt(5001),
				},
				constants.EthUsdPair: {
					SpotPrice: big.NewInt(6001),
					PnlPrice:  big.NewInt(6001),
				},
			},
		},
		"multiple prices, multiple validators, varied prices, random differences": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(6000),
						PnlPrice:  big.NewInt(6000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(8500),
						PnlPrice:  big.NewInt(8500),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(7250),
						PnlPrice:  big.NewInt(7250),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4500),
						PnlPrice:  big.NewInt(4500),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5500),
						PnlPrice:  big.NewInt(5500),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(5000),
					PnlPrice:  big.NewInt(5000),
				},
				constants.EthUsdPair: {
					SpotPrice: big.NewInt(6000),
					PnlPrice:  big.NewInt(6000),
				},
			},
		},
		"multiple prices, multiple validators, not enough power": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(6000),
						PnlPrice:  big.NewInt(6000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(6000),
						PnlPrice:  big.NewInt(6000),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{},
		},
		"single price, multiple validators, different stake": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 100,
				"bob":   200,
				"carl":  300,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(6000),
						PnlPrice:  big.NewInt(6000),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(7000),
						PnlPrice:  big.NewInt(7000),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(6000),
					PnlPrice:  big.NewInt(6000),
				},
			},
		},
		"multiple prices, multiple validators, varied prices, varied stake 1": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 300,
				"bob":   1,
				"carl":  1,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4000),
						PnlPrice:  big.NewInt(4000),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(8500),
						PnlPrice:  big.NewInt(8500),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(7250),
						PnlPrice:  big.NewInt(7250),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4500),
						PnlPrice:  big.NewInt(4500),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5500),
						PnlPrice:  big.NewInt(5500),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(4000),
					PnlPrice:  big.NewInt(4000),
				},
				constants.EthUsdPair: {
					SpotPrice: big.NewInt(5000),
					PnlPrice:  big.NewInt(5000),
				},
			},
		},
		"multiple prices, multiple validators, varied prices, varied stake 2": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 20,
				"bob":   70,
				"carl":  80,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4000),
						PnlPrice:  big.NewInt(4000),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5000),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4002),
						PnlPrice:  big.NewInt(4002),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5002),
						PnlPrice:  big.NewInt(5002),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4004),
						PnlPrice:  big.NewInt(4004),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5004),
						PnlPrice:  big.NewInt(5004),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(4002),
					PnlPrice:  big.NewInt(4002),
				},
				constants.EthUsdPair: {
					SpotPrice: big.NewInt(5002),
					PnlPrice:  big.NewInt(5002),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

			mCCVStore := ethosutils.NewGetAllCCValidatorMockReturnWithPowers(ctx, tc.validators, tc.powers)

			medianFn := vemath.Median(
				log.NewNopLogger(),
				mCCVStore,
				vemath.DefaultPowerThreshold,
			)

			prices, err := medianFn(ctx, tc.validatorPrices)

			require.NoError(t, err)
			require.Equal(t, tc.expectedPrices, prices)
		})
	}
}
