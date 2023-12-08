package clob

import clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

// TestHumanOrder is a test order with human readable price and size.
type TestHumanOrder struct {
	Order      clobtypes.Order
	HumanPrice string
	HumanSize  string
}
