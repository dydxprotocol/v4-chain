package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

var _ types.QueryServer = Keeper{}
