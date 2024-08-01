package types

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetVaultClobOrderClientId returns the client ID for a CLOB order where
// - 1st bit is `side-1` (subtract 1 as buy_side = 1, sell_side = 2)
//
// - next 8 bits are `layer`
func GetVaultClobOrderClientId(
	side clobtypes.Order_Side,
	layer uint8,
) uint32 {
	sideBit := uint32(side-1) << 31
	layerBits := uint32(layer) << 23

	return sideBit | layerBits
}
