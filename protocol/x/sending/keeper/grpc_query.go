package keeper

import (
	"github.com/dydxprotocol/v4/x/sending/types"
)

var _ types.QueryServer = Keeper{}
