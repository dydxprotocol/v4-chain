package keeper

import (
	"github.com/dydxprotocol/v4/x/perpetuals/types"
)

var _ types.QueryServer = Keeper{}
