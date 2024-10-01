package voteweighted_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
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

func TestMedianPrices(t *testing.T) {
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
					SpotPrice: big.NewInt(6500),
					PnlPrice:  big.NewInt(6500),
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
		"single price, multiple validators, same stake, different spot vs pnl": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 100,
				"bob":   100,
				"carl":  100,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5001),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5001),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5001),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(5000),
					PnlPrice:  big.NewInt(5001),
				},
			},
		},
		"single price, multiple validators, not enough stake, different spot vs pnl": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 100,
				"bob":   100,
				"carl":  100,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5001),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5001),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{},
		},
		"single price, multiple validators, same, stake, weird prices": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 100,
				"bob":   100,
				"carl":  100,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4000),
						PnlPrice:  big.NewInt(800),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4000),
						PnlPrice:  big.NewInt(8000),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4000),
						PnlPrice:  big.NewInt(8000),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(4000),
					PnlPrice:  big.NewInt(8000),
				},
			},
		},
		"multiple prices, multiple validators, same stake, different spot vs pnl": {
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
						PnlPrice:  big.NewInt(5001),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(6000),
						PnlPrice:  big.NewInt(6001),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(8500),
						PnlPrice:  big.NewInt(8501),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(7250),
						PnlPrice:  big.NewInt(7251),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4500),
						PnlPrice:  big.NewInt(4501),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5500),
						PnlPrice:  big.NewInt(5501),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(5000),
					PnlPrice:  big.NewInt(5001),
				},
				constants.EthUsdPair: {
					SpotPrice: big.NewInt(6000),
					PnlPrice:  big.NewInt(6001),
				},
			},
		},
		"multiple prices, multiple validators, insufficient stake, different spot vs pnl": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 100,
				"bob":   100,
				"carl":  100,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5001),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(6000),
						PnlPrice:  big.NewInt(6001),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5001),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(6000),
						PnlPrice:  big.NewInt(6001),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{},
		},
		"multiple prices, multiple validators, different spot vs pnl, varied stake 1": {
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
						PnlPrice:  big.NewInt(4005),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5005),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(8500),
						PnlPrice:  big.NewInt(8505),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(7250),
						PnlPrice:  big.NewInt(7255),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4500),
						PnlPrice:  big.NewInt(4505),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5500),
						PnlPrice:  big.NewInt(5505),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(4000),
					PnlPrice:  big.NewInt(4005),
				},
				constants.EthUsdPair: {
					SpotPrice: big.NewInt(5000),
					PnlPrice:  big.NewInt(5005),
				},
			},
		},
		"multiple prices, multiple validators, different spot vs pnl, varied stake 2": {
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
						PnlPrice:  big.NewInt(4005),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5005),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4002),
						PnlPrice:  big.NewInt(4007),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5002),
						PnlPrice:  big.NewInt(5007),
					},
				},
				constants.CarlEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4004),
						PnlPrice:  big.NewInt(4009),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5004),
						PnlPrice:  big.NewInt(5009),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(4002),
					PnlPrice:  big.NewInt(4007),
				},
				constants.EthUsdPair: {
					SpotPrice: big.NewInt(5002),
					PnlPrice:  big.NewInt(5007),
				},
			},
		},
		"voters barely have > 2/3": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 5000001,
				"bob":   5000000,
				"carl":  5000000,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4000),
						PnlPrice:  big.NewInt(4005),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5005),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4002),
						PnlPrice:  big.NewInt(4007),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5002),
						PnlPrice:  big.NewInt(5007),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(4000),
					PnlPrice:  big.NewInt(4005),
				},
				constants.EthUsdPair: {
					SpotPrice: big.NewInt(5000),
					PnlPrice:  big.NewInt(5005),
				},
			},
		},
		"voters barely have > 2/3 with large powers": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500000000000001,
				"bob":   500000000000000,
				"carl":  500000000000000,
			},
			validatorPrices: map[string]map[string]vemath.AggregatorPricePair{
				constants.AliceEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4000),
						PnlPrice:  big.NewInt(4005),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5000),
						PnlPrice:  big.NewInt(5005),
					},
				},
				constants.BobEthosConsAddress.String(): {
					constants.BtcUsdPair: {
						SpotPrice: big.NewInt(4002),
						PnlPrice:  big.NewInt(4007),
					},
					constants.EthUsdPair: {
						SpotPrice: big.NewInt(5002),
						PnlPrice:  big.NewInt(5007),
					},
				},
			},
			expectedPrices: map[string]vemath.AggregatorPricePair{
				constants.BtcUsdPair: {
					SpotPrice: big.NewInt(4000),
					PnlPrice:  big.NewInt(4005),
				},
				constants.EthUsdPair: {
					SpotPrice: big.NewInt(5000),
					PnlPrice:  big.NewInt(5005),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

			mCCVStore := ethosutils.NewGetAllCCValidatorMockReturnWithPowers(ctx, tc.validators, tc.powers)

			medianFn := vemath.MedianPrices(
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

func TestMedianConversionRate(t *testing.T) {
	tests := map[string]struct {
		validators               []string
		powers                   map[string]int64
		validatorConversionRates map[string]*big.Int
		expectedConversionRate   *big.Int
	}{
		"empty conversion rates": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorConversionRates: map[string]*big.Int{},
			expectedConversionRate:   nil,
		},
		"single validator": {
			validators: []string{"alice"},
			powers: map[string]int64{
				"alice": 500,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(5000),
			},
			expectedConversionRate: big.NewInt(5000),
		},
		"two validators, varied conversion rates": {
			validators: []string{"alice", "bob"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(5000),
				constants.BobEthosConsAddress.String():   big.NewInt(5001),
			},
			expectedConversionRate: big.NewInt(5000),
		},
		"multiple validators, varied conversion rates": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(5000),
				constants.BobEthosConsAddress.String():   big.NewInt(5001),
				constants.CarlEthosConsAddress.String():  big.NewInt(5002),
			},
			expectedConversionRate: big.NewInt(5001),
		},
		"two validators, varied conversion rates, one nil": {
			validators: []string{"alice", "bob"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(5000),
				constants.BobEthosConsAddress.String():   nil,
			},
			expectedConversionRate: nil,
		},
		"multiple validators, varied rates, one nil, returns nil": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(5000),
				constants.BobEthosConsAddress.String():   nil,
				constants.CarlEthosConsAddress.String():  big.NewInt(5002),
			},
			expectedConversionRate: nil,
		},
		"multiple validators with different power, varied conversion rates, one nil, returns non-nil": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 5000001,
				"bob":   5000000,
				"carl":  5000000,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(5000),
				constants.BobEthosConsAddress.String():   nil,
				constants.CarlEthosConsAddress.String():  big.NewInt(5002),
			},
			expectedConversionRate: big.NewInt(5000),
		},
		"multiple validators with different large powers, varied conversion rates, one nil, returns non-nil": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500000000000001,
				"bob":   500000000000000,
				"carl":  500000000000000,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(5000),
				constants.BobEthosConsAddress.String():   nil,
				constants.CarlEthosConsAddress.String():  big.NewInt(5002),
			},
			expectedConversionRate: big.NewInt(5000),
		},
		"multiple validators, varied rates, all nil": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): nil,
				constants.BobEthosConsAddress.String():   nil,
				constants.CarlEthosConsAddress.String():  nil,
			},
			expectedConversionRate: nil,
		},
		"multiple validators, same rate": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): constants.Price5Big,
				constants.BobEthosConsAddress.String():   constants.Price5Big,
				constants.CarlEthosConsAddress.String():  constants.Price5Big,
			},
			expectedConversionRate: constants.Price5Big,
		},
		"multiple validators, not enough power to return non-nil": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 500,
				"bob":   500,
				"carl":  500,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(5000),
				constants.BobEthosConsAddress.String():   big.NewInt(5000),
			},
			expectedConversionRate: nil,
		},
		"multiple validators, different stake": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 100,
				"bob":   200,
				"carl":  300,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(5000),
				constants.BobEthosConsAddress.String():   big.NewInt(6000),
				constants.CarlEthosConsAddress.String():  big.NewInt(7000),
			},
			expectedConversionRate: big.NewInt(6500),
		},
		"multiple validators, varied rates, varied stake 1": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 300,
				"bob":   1,
				"carl":  1,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(4000),
				constants.BobEthosConsAddress.String():   big.NewInt(8500),
				constants.CarlEthosConsAddress.String():  big.NewInt(4500),
			},
			expectedConversionRate: big.NewInt(4000),
		},
		"multiple validators, varied rates, varied stake 2": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 20,
				"bob":   70,
				"carl":  80,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(4000),
				constants.BobEthosConsAddress.String():   big.NewInt(4002),
				constants.CarlEthosConsAddress.String():  big.NewInt(4004),
			},
			expectedConversionRate: big.NewInt(4002),
		},
		"submitted conversion rate with no stake": {
			validators: []string{"alice", "bob", "carl"},
			powers: map[string]int64{
				"alice": 60,
				"bob":   50,
			},
			validatorConversionRates: map[string]*big.Int{
				constants.AliceEthosConsAddress.String(): big.NewInt(4000),
				constants.BobEthosConsAddress.String():   big.NewInt(4002),
				constants.CarlEthosConsAddress.String():  big.NewInt(4004),
			},
			expectedConversionRate: big.NewInt(4000),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

			mCCVStore := ethosutils.NewGetAllCCValidatorMockReturnWithPowers(ctx, tc.validators, tc.powers)

			medianFn := vemath.MedianConversionRate(
				log.NewNopLogger(),
				mCCVStore,
				vemath.DefaultPowerThreshold,
			)

			conversionRate, err := medianFn(ctx, tc.validatorConversionRates)

			require.NoError(t, err)
			require.Equal(t, tc.expectedConversionRate, conversionRate)
		})
	}
}

func TestComputeMedian(t *testing.T) {
	tests := map[string]struct {
		prices      []vemath.PricePerValidator
		totalWeight math.Int
		expected    *big.Int
	}{
		"single price": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 100, Price: big.NewInt(5000)},
			},
			totalWeight: math.NewInt(100),
			expected:    big.NewInt(5000),
		},
		"two prices, equal weight": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 50, Price: big.NewInt(5000)},
				{VoteWeight: 50, Price: big.NewInt(6000)},
			},
			totalWeight: math.NewInt(100),
			expected:    big.NewInt(5500),
		},
		"three prices, equal weight": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 33, Price: big.NewInt(5000)},
				{VoteWeight: 33, Price: big.NewInt(6000)},
				{VoteWeight: 33, Price: big.NewInt(7000)},
			},
			totalWeight: math.NewInt(99),
			expected:    big.NewInt(6000),
		},
		"three prices, unequal weight": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 20, Price: big.NewInt(5000)},
				{VoteWeight: 30, Price: big.NewInt(6000)},
				{VoteWeight: 51, Price: big.NewInt(7000)},
			},
			totalWeight: math.NewInt(101),
			expected:    big.NewInt(7000),
		},
		"five prices, varied weights": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 10, Price: big.NewInt(5000)},
				{VoteWeight: 20, Price: big.NewInt(5500)},
				{VoteWeight: 30, Price: big.NewInt(6000)},
				{VoteWeight: 25, Price: big.NewInt(6500)},
				{VoteWeight: 15, Price: big.NewInt(7000)},
			},
			totalWeight: math.NewInt(100),
			expected:    big.NewInt(6000),
		},
		"prices not in order": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 30, Price: big.NewInt(6000)},
				{VoteWeight: 20, Price: big.NewInt(5000)},
				{VoteWeight: 50, Price: big.NewInt(5500)},
			},
			totalWeight: math.NewInt(100),
			expected:    big.NewInt(5500),
		},
		"even number of prices, unequal weights": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 25, Price: big.NewInt(5000)},
				{VoteWeight: 25, Price: big.NewInt(5500)},
				{VoteWeight: 30, Price: big.NewInt(6000)},
				{VoteWeight: 20, Price: big.NewInt(6500)},
			},
			totalWeight: math.NewInt(100),
			expected:    big.NewInt(5750), // Average of 5500 and 6000
		},
		"large weights": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 1000000, Price: big.NewInt(5000)},
				{VoteWeight: 2000000, Price: big.NewInt(6000)},
				{VoteWeight: 3000000, Price: big.NewInt(7000)},
			},
			totalWeight: math.NewInt(6000000),
			expected:    big.NewInt(6500),
		},
		"single price with large weight": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 1000000, Price: big.NewInt(5000)},
			},
			totalWeight: math.NewInt(1000000),
			expected:    big.NewInt(5000),
		},
		"no prices": {
			prices:      []vemath.PricePerValidator{},
			totalWeight: math.NewInt(0),
			expected:    nil,
		},
		"even total weight, two middle values": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 25, Price: big.NewInt(5000)},
				{VoteWeight: 25, Price: big.NewInt(6000)},
				{VoteWeight: 25, Price: big.NewInt(7000)},
				{VoteWeight: 25, Price: big.NewInt(8000)},
			},
			totalWeight: math.NewInt(100),
			expected:    big.NewInt(6500), // Average of 6000 and 7000
		},
		"even total weight, multiple prices with same weight": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 20, Price: big.NewInt(5000)},
				{VoteWeight: 20, Price: big.NewInt(5500)},
				{VoteWeight: 20, Price: big.NewInt(6000)},
				{VoteWeight: 20, Price: big.NewInt(6500)},
				{VoteWeight: 20, Price: big.NewInt(7000)},
			},
			totalWeight: math.NewInt(100),
			expected:    big.NewInt(6000), // Average of 5500 and 6500
		},
		"rounds down price for median with even total weight": {
			prices: []vemath.PricePerValidator{
				{VoteWeight: 25, Price: big.NewInt(1)},
				{VoteWeight: 25, Price: big.NewInt(2)},
				{VoteWeight: 25, Price: big.NewInt(3)},
				{VoteWeight: 25, Price: big.NewInt(4)},
			},
			totalWeight: math.NewInt(100),
			expected:    big.NewInt(2),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := vemath.ComputeMedian(tc.prices, tc.totalWeight)
			require.Equal(t, tc.expected, result)
		})
	}
}
