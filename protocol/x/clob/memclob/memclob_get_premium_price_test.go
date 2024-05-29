package memclob

import (
	"math"
	"math/big"
	"testing"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestGetPremiumPrice(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders []types.MatchableOrder

		// Parameters.
		clobPair                    types.ClobPair
		indexPrice                  pricestypes.MarketPrice
		baseAtomicResolution        int32
		maxAbsPremiumVotePpm        *big.Int
		impactNotionalQuoteQuantums *big.Int

		// Expectations.
		expectedErr        error
		expectedPremiumPpm int32
		shouldPanic        bool
	}{
		`Best Bid < Index < Best Ask, premium = 0`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_010_000,    // $10_001 (Impact Ask)
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_990_000,     // $9_999 (Impact Bid)
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(1_000_000), // 100%
			indexPrice: pricestypes.MarketPrice{
				Price:    1_000_000_000, // $10_000.
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          0,
		},
		`Index < Impact Bid < Best Bid < Best Ask, positive premium`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_010_000,    // $10_001 (Impact Ask)
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     1_000_000_000, // 0.1 BTC
					Subticks:     99_990_000,    // $9_999
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_980_000,     // $9_998 (Impact Bid = $9998.2)
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(1_000_000), // 100%
			indexPrice: pricestypes.MarketPrice{
				Price:    900_000_000, // $9_000.
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          110911,                                // 9_998.2 / 9_000 - 1 = 0.110911
		},
		`Impact Bid < Best Bid < Best Ask < Impact Ask < Index, negative premium`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     2,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_020_000,    // $10_002 (Impact Ask = $10001.7)
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     1_000_000_000, // 0.1 BTC
					Subticks:     100_015_000,   // $10_001.5
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     1_000_000_000, // 0.1 BTC
					Subticks:     100_010_000,   // $10_001
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_990_000,     // $9_999
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(1_000_000), // 100%
			indexPrice: pricestypes.MarketPrice{
				Price:    1_000_300_000, // $10_003
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          -129,                                  // 10_001.7 / 10_003 - 1 = -0.000129
		},
		`Impact Bid < Best Bid < Best Ask = Impact Ask < Index, negative premium`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_010_000,    // $10_001
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_990_000,     // $9_999
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(1_000_000), // 100%
			indexPrice: pricestypes.MarketPrice{
				Price:    1_000_190_000, // $10_001.9
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          -89,                                   // 10_001 / 10_001.9 - 1 = -0.000089
		},
		`Index < Impact Bid = Best Bid < Best Ask, positive premium`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_010_000,    // $10_001
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_999_000,     // $9_999.9
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(1_000_000), // 100%
			indexPrice: pricestypes.MarketPrice{
				Price:    999_750_000, // $9_997.5
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          240,                                   // 9_999.9 / 9_997.5 - 1 = 0.000240
		},
		`Impact Bid < Index < Best Bid < Best Ask, 0 premium`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_010_000,    // $10_001
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     1_000_000_000, // 0.1 BTC
					Subticks:     99_999_000,    // $9_999.9
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 0.1 BTC
					Subticks:     99_995_000,     // $9_999.5
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_990_000,     // $9_999
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(1_000_000), // 100%
			indexPrice: pricestypes.MarketPrice{
				Price:    999_982_000, // $9_999.5
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          0,
		},
		`BestAsk < Index; Impact Ask = Infinity (low liquidity); 0 premium`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     1_000_000_000, // 0.1 BTC
					Subticks:     100_001_000,   // $10_000.1
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_998_000,     // $9_999.8
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(100_000), // 10%
			indexPrice: pricestypes.MarketPrice{
				Price:    1_000_100_000, // $10_001
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(2_000_000_000), // $2000
			expectedPremiumPpm:          0,
		},
		`Impact Bid = 0 (low liquidity); Index < Best Bid; 0 premium`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_001_000,    // $10_000.1
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     1_000_000_000, // 0.1 BTC
					Subticks:     99_998_000,    // $9_999.8
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(100_000), // 10%
			indexPrice: pricestypes.MarketPrice{
				Price:    999_950_000, // $9_999.5
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(2_000_000_000), // $2000
			expectedPremiumPpm:          0,
		},
		`Not enough liquidity on both sides, return 0`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     1_000_000_000, // 0.1 BTC
					Subticks:     100_001_000,   // $10_000.1
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     1_000_000_000, // 0.1 BTC
					Subticks:     99_998_000,    // $9_999.8
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(100_000), // 10%
			indexPrice: pricestypes.MarketPrice{
				Price:    1_000_000_000, // $10_000
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(2_000_000_000), // $2000
			expectedPremiumPpm:          0,                                     // 0%
		},
		`Orderbook is empty, return 0`: {
			placedMatchableOrders: []types.MatchableOrder{},
			clobPair:              constants.ClobPair_Btc,
			maxAbsPremiumVotePpm:  big.NewInt(100_000), // 10%
			indexPrice: pricestypes.MarketPrice{
				Price:    1_000_000_000, // $10_000
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(2_000_000_000), // $2000
			expectedPremiumPpm:          0,                                     // 0%
		},
		`Index << Impact Bid, maximum premium (clamped)`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_010_000,    // $10_001
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_999_000,     // $9_999.9
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(100_000), // 10%
			indexPrice: pricestypes.MarketPrice{
				Price:    600_000_000, // $6_000
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          100_000,
		},
		`Impact Ask << Index, minimum premium (clamped)`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_010_000,    // $10_001
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_999_000,     // $9_999.9
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(100_000), // 10%
			indexPrice: pricestypes.MarketPrice{
				Price:    1_500_000_000, // $6_000
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          -100_000,
		},
		`Index < Impact Bid < Impact Ask = Infinity (low liquidity); positive premium`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     1_000_000_000, // 0.1 BTC
					Subticks:     100_001_000,   // $10_000.1
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_998_000,     // $9_999.8
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(100_000), // 10%
			indexPrice: pricestypes.MarketPrice{
				Price:    999_900_000, // $9_999
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          80,                                    // 0.008%
		},
		`0 = Impact Bid (low liquidity) < Impact Ask; negative premium`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_001_000,    // $10_000.1
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     1_000_000_000, // 0.11 BTC
					Subticks:     99_998_000,    // $9_999.8
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(100_000), // 10%
			indexPrice: pricestypes.MarketPrice{
				Price:    1_000_100_000, // $10_001
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(5_000_000_000), // $5000
			expectedPremiumPpm:          -89,                                   // -0.0089%
		},
		"error: maxAbsPremiumVotePpm overflow int32": {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_010_000,    // $10_001 (Impact Ask)
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			maxAbsPremiumVotePpm: big.NewInt(math.MaxInt32 + 1),
			shouldPanic:          true,
		},
		"error: clob pair is not a perpetual": {
			clobPair:                    constants.ClobPair_Spot_Btc,
			maxAbsPremiumVotePpm:        big.NewInt(100_000), // 10%
			impactNotionalQuoteQuantums: big.NewInt(1000),
			expectedErr: errorsmod.Wrapf(
				types.ErrPremiumWithNonPerpetualClobPair,
				"ClobPair ID: %d",
				constants.ClobPair_Spot_Btc.Id,
			),
		},
		"error: index price is zero": {
			clobPair:                    constants.ClobPair_Btc,
			maxAbsPremiumVotePpm:        big.NewInt(100_000), // 10%
			impactNotionalQuoteQuantums: big.NewInt(1000),

			indexPrice: pricestypes.MarketPrice{
				Price: 0,
			},
			expectedErr: types.ErrZeroIndexPriceForPremiumCalculation,
		},
		`Zero impact amount, positive premium, use best bid as impact price`: {
			placedMatchableOrders: []types.MatchableOrder{
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Bob_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_SELL,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     100_010_000,    // $10_001
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
				&types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10_000_000_000, // 1 BTC
					Subticks:     99_999_000,     // $9_999.9
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
				},
			},
			clobPair:             constants.ClobPair_Btc,
			maxAbsPremiumVotePpm: big.NewInt(1_000_000), // 100%
			indexPrice: pricestypes.MarketPrice{
				Price:    999_750_000, // $9_997.5
				Exponent: -5,
			},
			// 1 baseQuantum = 10^(-10) BTC.
			baseAtomicResolution:        -10,
			impactNotionalQuoteQuantums: new(big.Int).SetUint64(0), //
			expectedPremiumPpm:          240,                       // 9_999.9 / 9_997.5 - 1 = 0.000240
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob and orderbook state.
			memclob, _ := setUpMemclobAndOrderbook(
				t,
				ctx,
				tc.placedMatchableOrders,
				nil,
				[]types.MatchableOrder{},
			)
			if len(tc.placedMatchableOrders) == 0 {
				// Create the orderbook when there are no orders.
				memclob.CreateOrderbook(tc.clobPair)
			}

			pricePremiumParams := perptypes.GetPricePremiumParams{
				IndexPrice:                  tc.indexPrice,
				BaseAtomicResolution:        tc.baseAtomicResolution,
				QuoteAtomicResolution:       lib.QuoteCurrencyAtomicResolution,
				ImpactNotionalQuoteQuantums: tc.impactNotionalQuoteQuantums,
				MaxAbsPremiumVotePpm:        tc.maxAbsPremiumVotePpm,
			}

			if tc.shouldPanic {
				require.Panics(t, func() {
					//nolint:errcheck
					memclob.GetPricePremium(
						ctx,
						tc.clobPair,
						pricePremiumParams,
					)
				})
				return
			}

			premiumPpm, err := memclob.GetPricePremium(
				ctx,
				tc.clobPair,
				pricePremiumParams,
			)

			if tc.expectedErr != nil {
				require.ErrorContains(t,
					err,
					tc.expectedErr.Error(),
				)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedPremiumPpm, premiumPpm)
		})
	}
}
