package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
)

const Name = "daemons"

// daemon errors
var (
	// Generic daemon server errors.
	ErrServerNotInitializedCorrectly = errorsmod.Register(Name, 1, "Daemon server not initialized correctly")

	// PriceFeed daemon service errors will have code 1xxx.
	ErrPriceFeedInvalidPrice = errorsmod.Register(
		Name,
		1000,
		fmt.Sprintf("Price is set to %d which is not a valid price", constants.DefaultPrice),
	)
	ErrPriceFeedLastUpdateTimeNotSet   = errorsmod.Register(Name, 1001, "LastUpdateTime is not set")
	ErrPriceFeedMarketPriceUpdateEmpty = errorsmod.Register(Name, 1002, "Market price update has length of 0")
)
