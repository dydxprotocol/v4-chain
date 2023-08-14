package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
)

var _ types.QueryServer = Keeper{}
