package abci_test

import (
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/abci"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCached_Success(t *testing.T) {
	ms := &mocks.MultiStore{}
	cms := &mocks.CacheMultiStore{}
	// Expect that the cached store is created and returned.
	ms.On("CacheMultiStore").Return(cms).Once()
	// Expect that the cache is written to the underlying store.
	cms.On("Write").Return(nil).Once()

	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())

	err := abci.RunCached(ctx, func(ctx sdk.Context) error { return nil })
	require.NoError(t, err)

	ms.AssertExpectations(t)
	cms.AssertExpectations(t)
}

func TestRunCached_Failure(t *testing.T) {
	ms := &mocks.MultiStore{}
	cms := &mocks.CacheMultiStore{} // We don't mock the Write method because it should not be called.

	// Expect that the cached store is created and returned.
	ms.On("CacheMultiStore").Return(cms).Once()

	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())

	// Expect that the cache is discarded. The test will fail if the cache is persisted because the
	// Write method of the CacheMultiStore is not mocked here.
	err := abci.RunCached(ctx, func(ctx sdk.Context) error { return fmt.Errorf("failure") })
	require.ErrorContains(t, err, "failure")

	ms.AssertExpectations(t)
}
