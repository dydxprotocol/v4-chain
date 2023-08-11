package keeper

import (
	"github.com/dydxprotocol/v4/x/assets/types"
)

var _ types.QueryServer = Keeper{}
