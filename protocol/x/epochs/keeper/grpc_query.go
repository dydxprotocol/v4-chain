package keeper

import (
	"github.com/dydxprotocol/v4/x/epochs/types"
)

var _ types.QueryServer = Keeper{}
