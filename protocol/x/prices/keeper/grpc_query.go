package keeper

import (
	"github.com/dydxprotocol/v4/x/prices/types"
)

var _ types.QueryServer = Keeper{}
