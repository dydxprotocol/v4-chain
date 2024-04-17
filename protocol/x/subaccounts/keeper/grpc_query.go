package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

var _ types.QueryServer = Keeper{}
