package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

var _ types.QueryServer = Keeper{}
