package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// settledUpdate is used internally in the subaccounts keeper to
// to specify changes to one or more `Subaccounts` (for example the
// result of a trade, transfer, etc).
// The subaccount is always in its settled state.
type settledUpdate struct {
	// The `Subaccount` for which this update applies to, in its settled form.
	SettledSubaccount types.Subaccount
	// A list of changes to make to any `AssetPositions` in the `Subaccount`.
	AssetUpdates []types.AssetUpdate
	// A list of changes to make to any `PerpetualPositions` in the `Subaccount`.
	PerpetualUpdates []types.PerpetualUpdate
}
