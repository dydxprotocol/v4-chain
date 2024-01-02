package abci_test

import (
	"cosmossdk.io/log"
	"fmt"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/abci"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCached_Mixed(t *testing.T) {
	testEvent := sdk.NewEvent("test", sdk.NewAttribute("key", "value"))
	tests := map[string]struct {
		f             func(ctx sdk.Context) error
		expectedError error
	}{
		"success": {
			f: func(ctx sdk.Context) error {
				ctx.EventManager().EmitEvent(testEvent)
				return nil
			},
		},
		"failure": {
			f: func(ctx sdk.Context) error {
				ctx.EventManager().EmitEvent(testEvent)
				return fmt.Errorf("failure")
			},
			expectedError: fmt.Errorf("failure"),
		},
		"panic": {
			f: func(ctx sdk.Context) error {
				ctx.EventManager().EmitEvent(testEvent)
				panic("panic")
			},
			expectedError: fmt.Errorf("panic"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ms := &mocks.MultiStore{}
			cms := &mocks.CacheMultiStore{}

			// Expect that the cached store is created and returned.
			ms.On("CacheMultiStore").Return(cms).Once()

			if tc.expectedError == nil {
				// For non-error cases, expect that the cache is written to the underlying store.
				cms.On("Write").Return(nil).Once()
			}

			ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())

			err := abci.RunCached(ctx, tc.f)

			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
				require.Len(t, ctx.EventManager().Events(), 0)
			} else {
				require.NoError(t, err)
				require.Equal(t, ctx.EventManager().Events(), sdk.Events{testEvent})
			}

			ms.AssertExpectations(t)
			cms.AssertExpectations(t)
		})
	}
}
