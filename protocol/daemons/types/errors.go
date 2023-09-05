package types

import (
	moderrors "cosmossdk.io/errors"
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
)

const Name = "daemons"

// daemon errors
var (
	// Generic daemon server errors.
	ErrDaemonMethodNotImplemented    = moderrors.Register(Name, 1, "Daemon method not implemented")
	ErrServerNotInitializedCorrectly = moderrors.Register(Name, 2, "Daemon server not initialized correctly")

	// PriceFeed daemon service errors will have code 1xxx.
	ErrPriceFeedInvalidPrice = moderrors.Register(
		Name,
		1000,
		fmt.Sprintf("Price is set to %d which is not a valid price", constants.DefaultPrice),
	)
	ErrPriceFeedLastUpdateTimeNotSet   = moderrors.Register(Name, 1001, "LastUpdateTime is not set")
	ErrPriceFeedMarketPriceUpdateEmpty = moderrors.Register(Name, 1002, "Market price update has length of 0")
)
