package price_fetcher

import "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"

// taskLoopDefinition defines the parameters for the price fetcher's task loop:
// Which markets to query, and parameters for how to query them and process query results.
// It's mostly identical to the price fetcher's mutable state, except that the market ids
// to query for the next task loop are articulated explicitly. This struct is used to pass
// snapshot of the current parameters to the price fetcher's task loop in order to prevent
// writes to the price fetcher's mutable state from putting the task loop into an
// inconsistent state whenever the price fetcher's mutable state is updated.
//
// This was judged as a better alternative than locking the price fetcher's
// mutable state for the duration of the task loop so that it would not create any blocking
// issues for the go routine that fetches config updates from the pricefeed server. If the
// price fetcher's mutable state updates in the middle of a task loop execution, it will be
// ignored by that loop and picked up by the next one.
type taskLoopDefinition struct {
	mutableExchangeConfig *types.MutableExchangeMarketConfig
	marketExponents       map[types.MarketId]types.Exponent
	marketIds             []types.MarketId
}
