package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

var _ types.QueryServer = Keeper{}
