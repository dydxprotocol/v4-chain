package keeper

import (
	"github.com/dydxprotocol/v4/x/subaccounts/types"
)

var _ types.QueryServer = Keeper{}
