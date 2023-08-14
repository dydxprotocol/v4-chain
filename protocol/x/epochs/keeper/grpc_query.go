package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
)

var _ types.QueryServer = Keeper{}
