package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

var _ types.QueryServer = Keeper{}
