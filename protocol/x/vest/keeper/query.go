package keeper

import (
	"github.com/dydxprotocol/v4/x/vest/types"
)

var _ types.QueryServer = Keeper{}
