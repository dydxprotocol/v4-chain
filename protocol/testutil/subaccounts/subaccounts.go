package subaccounts

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func CreatePerpetualUpdate(
	id uint32,
	delta *big.Int,
) []types.PerpetualUpdate {
	return []types.PerpetualUpdate{
		{
			PerpetualId:      id,
			BigQuantumsDelta: delta,
		},
	}
}
