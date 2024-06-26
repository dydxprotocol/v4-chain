package deamon

import (
	"math/big"

	"cosmossdk.io/log"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
)

// PreBlockHandler is responsible for aggregating deamon data from each
// validator and writing the prices data into the store before any transactions
// are executed/finalized for a given block.
type PreBlockHandler struct { //golint:ignore
	logger log.Logger

	// keeper is the keeper for the prices module. This is utilized to write
	// deamon price data to state.
	keeper pk.Keeper
}

// NewOraclePreBlockHandler returns a new PreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewOraclePreBlockHandler(
	logger log.Logger,
	aggregateFn aggregator.AggregateFnFromContext[string, map[slinkytypes.CurrencyPair]*big.Int],
	pk pk.Keeper,
	veCodec codec.VoteExtensionCodec,
	ecCodec codec.ExtendedCommitCodec,
) *PreBlockHandler {
	va := abciaggregator.NewDefaultVoteAggregator(
		logger,
		aggregateFn,
		strategy,
	)

	return &PreBlockHandler{
		logger: logger,
		keeper: pk,
	}
}
