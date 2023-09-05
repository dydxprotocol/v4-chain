package sdk

import (
	"cosmossdk.io/store/metrics"
	dbm "github.com/cosmos/cosmos-db"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewSdkContextWithMultistore() (
	ctx sdk.Context,
	stateStore store.CommitMultiStore,
	db *dbm.MemDB,
) {
	db = dbm.NewMemDB()
	stateStore = store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	ctx = sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	ctx = ctx.WithTxBytes([]byte{1})
	return ctx, stateStore, db
}

func NewContextWithBlockHeightAndTime(
	height int64,
	time time.Time,
) (
	ctx sdk.Context,
) {
	return sdk.NewContext(nil, tmproto.Header{}, false, log.NewNopLogger()).
		WithBlockHeight(height).
		WithBlockTime(time)
}
