package keeper_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGetIsolatedPerpetualStateTransition(t *testing.T) {
	tests := map[string]struct {
		// parameters
		subaccountQuoteQuantums   *big.Int
		isolatedPerpetualUpdate   *types.PerpetualUpdate
		isolatedPerpetualPosition *types.PerpetualPosition

		// expectation
		expectedStateTransition *types.IsolatedPerpetualPositionStateTransition
	}{
		`If perpetual update is nil, nil state transition is returned`: {
			subaccountQuoteQuantums:   big.NewInt(100_000_000), // $100
			isolatedPerpetualUpdate:   nil,
			isolatedPerpetualPosition: nil,
			expectedStateTransition:   nil,
		},
		`If perpetual update is nil and isolated position exists, nil state transition is returned`: {
			subaccountQuoteQuantums:   big.NewInt(100_000_000), // $100
			isolatedPerpetualUpdate:   nil,
			isolatedPerpetualPosition: &constants.PerpetualPosition_OneISOLong,
			expectedStateTransition:   nil,
		},
		`If perpetual update exists and isolated position is nil, state transition representing an
		isolated perpetual position being opened is returned`: {
			subaccountQuoteQuantums: big.NewInt(100_000_000), // $100
			isolatedPerpetualUpdate: &types.PerpetualUpdate{
				PerpetualId:      uint32(3),
				BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ISO
			},
			isolatedPerpetualPosition: nil,
			expectedStateTransition: &types.IsolatedPerpetualPositionStateTransition{
				PerpetualId:               uint32(3),
				QuoteQuantumsBeforeUpdate: big.NewInt(100_000_000),
				Transition:                types.Opened,
			},
		},
		`If perpetual update exists and isolated position exists, and perpetual update would close
		isolated position, state transition representing an isolated perpetual position being closed is returned`: {
			subaccountQuoteQuantums: big.NewInt(100_000_000), // $100
			isolatedPerpetualUpdate: &types.PerpetualUpdate{
				PerpetualId:      uint32(3),
				BigQuantumsDelta: new(big.Int).Neg(constants.PerpetualPosition_OneISOLong.GetBigQuantums()),
			},
			isolatedPerpetualPosition: &constants.PerpetualPosition_OneISOLong,
			expectedStateTransition: &types.IsolatedPerpetualPositionStateTransition{
				PerpetualId:               uint32(3),
				QuoteQuantumsBeforeUpdate: big.NewInt(100_000_000),
				Transition:                types.Closed,
			},
		},
		`If perpetual update exists and isolated position exists, and perpetual update would increase
		isolated position, nil state transition is returned`: {
			subaccountQuoteQuantums: big.NewInt(100_000_000), // $100
			isolatedPerpetualUpdate: &types.PerpetualUpdate{
				PerpetualId:      uint32(3),
				BigQuantumsDelta: big.NewInt(10_000_000),
			},
			isolatedPerpetualPosition: &constants.PerpetualPosition_OneISOLong,
			expectedStateTransition:   nil,
		},
		`If perpetual update exists and isolated position exists, and perpetual update would decrease
		isolated position, nil state transition is returned`: {
			subaccountQuoteQuantums: big.NewInt(100_000_000), // $100
			isolatedPerpetualUpdate: &types.PerpetualUpdate{
				PerpetualId:      uint32(3),
				BigQuantumsDelta: big.NewInt(-10_000_000),
			},
			isolatedPerpetualPosition: &constants.PerpetualPosition_OneISOLong,
			expectedStateTransition:   nil,
		},
		`If perpetual update exists and isolated position exists, and perpetual update would flip
		isolated position, nil state transition is returned`: {
			subaccountQuoteQuantums: big.NewInt(100_000_000), // $100
			isolatedPerpetualUpdate: &types.PerpetualUpdate{
				PerpetualId:      uint32(3),
				BigQuantumsDelta: big.NewInt(-10_000_000_000),
			},
			isolatedPerpetualPosition: &constants.PerpetualPosition_OneISOLong,
			expectedStateTransition:   nil,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				stateTransition := keeper.GetIsolatedPerpetualStateTransition(
					tc.subaccountQuoteQuantums,
					tc.isolatedPerpetualPosition,
					tc.isolatedPerpetualUpdate,
				)
				require.Equal(t, tc.expectedStateTransition, stateTransition)
			},
		)
	}
}
