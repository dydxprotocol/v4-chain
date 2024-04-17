package constants

import (
	"math/big"

	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

var (
	// Get state position functions.
	GetStatePosition_ZeroPositionSize = func(
		subaccountId satypes.SubaccountId,
		clobPairId clobtypes.ClobPairId,
	) (
		statePositionSize *big.Int,
	) {
		return big.NewInt(0)
	}
)
