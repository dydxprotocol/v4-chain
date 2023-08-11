package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestGetMarketPremiumsMap(t *testing.T) {
	tests := map[string]struct {
		premiumStore types.PremiumStore
		expectedMap  map[uint32]types.MarketPremiums
	}{
		"3 perpetuals from 0 to 2": {
			premiumStore: types.PremiumStore{
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{1000, -1000, 1000},
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{0, 0, 0},
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{5000, 100, -100, 0},
					},
				},
			},
			expectedMap: map[uint32]types.MarketPremiums{
				0: {
					PerpetualId: 0,
					Premiums:    []int32{1000, -1000, 1000},
				},
				1: {
					PerpetualId: 1,
					Premiums:    []int32{0, 0, 0},
				},
				2: {
					PerpetualId: 2,
					Premiums:    []int32{5000, 100, -100, 0},
				},
			},
		},
		"Some perpetuals not present; some with empty entries": {
			premiumStore: types.PremiumStore{
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{1000, -1000},
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{},
					},
					{
						PerpetualId: 5,
						Premiums:    []int32{0},
					},
				},
			},
			expectedMap: map[uint32]types.MarketPremiums{
				0: {
					PerpetualId: 0,
					Premiums:    []int32{1000, -1000},
				},
				2: {
					PerpetualId: 2,
					Premiums:    []int32{},
				},
				5: {
					PerpetualId: 5,
					Premiums:    []int32{0},
				},
			},
		},
	}

	for _, tc := range tests {
		require.Equal(
			t,
			tc.expectedMap,
			tc.premiumStore.GetMarketPremiumsMap(),
		)
	}
}

func TestNewPremiumStoreFromMarketPremiumMap(t *testing.T) {
	numPremiums := uint32(10)

	tests := map[string]struct {
		marketPremiumsMap    map[uint32]types.MarketPremiums
		numPerpetuals        uint32
		expectedPremiumStore types.PremiumStore
	}{
		"3 perpetuals from 0 to 2": {
			numPerpetuals: 3,
			marketPremiumsMap: map[uint32]types.MarketPremiums{
				2: {
					PerpetualId: 2,
					Premiums:    []int32{5000, 100, -100, 0},
				},
				0: {
					PerpetualId: 0,
					Premiums:    []int32{1000, -1000, 1000},
				},
				1: {
					PerpetualId: 1,
					Premiums:    []int32{0, 0, 0},
				},
			},
			expectedPremiumStore: types.PremiumStore{
				NumPremiums: numPremiums,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{1000, -1000, 1000},
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{0, 0, 0},
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{5000, 100, -100, 0},
					},
				},
			},
		},
		"3 perpetuals from 0 to 2, numPerpetuals = 5": {
			numPerpetuals: 5,
			marketPremiumsMap: map[uint32]types.MarketPremiums{
				2: {
					PerpetualId: 2,
					Premiums:    []int32{5000, 100, -100, 0},
				},
				0: {
					PerpetualId: 0,
					Premiums:    []int32{1000, -1000, 1000},
				},
				1: {
					PerpetualId: 1,
					Premiums:    []int32{0, 0, 0},
				},
			},
			expectedPremiumStore: types.PremiumStore{
				NumPremiums: numPremiums,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{1000, -1000, 1000},
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{0, 0, 0},
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{5000, 100, -100, 0},
					},
				},
			},
		},
		"Some perpetuals not present; some with empty entries": {
			numPerpetuals: 6,
			marketPremiumsMap: map[uint32]types.MarketPremiums{
				0: {
					PerpetualId: 0,
					Premiums:    []int32{1000, -1000},
				},
				2: {
					PerpetualId: 2,
					Premiums:    []int32{},
				},
				5: {
					PerpetualId: 5,
					Premiums:    []int32{0},
				},
			},
			expectedPremiumStore: types.PremiumStore{
				NumPremiums: numPremiums,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{1000, -1000},
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{},
					},
					{
						PerpetualId: 5,
						Premiums:    []int32{0},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		require.Equal(
			t,
			tc.expectedPremiumStore,
			*types.NewPremiumStoreFromMarketPremiumMap(
				tc.marketPremiumsMap,
				tc.numPerpetuals,
				numPremiums,
			),
		)
	}
}
