package sdk

import (
	"time"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewSdkContextWithMultistore() (
	ctx sdk.Context,
	stateStore store.CommitMultiStore,
	db *tmdb.MemDB,
) {
	db = tmdb.NewMemDB()
	stateStore = store.NewCommitMultiStore(db)
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
