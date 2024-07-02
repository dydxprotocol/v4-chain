package types

// SettledUpdate is used internally in the subaccounts keeper to
// to specify changes to one or more `Subaccounts` (for example the
// result of a trade, transfer, etc).
// The subaccount is always in its settled state.
type SettledUpdate struct {
	// The `Subaccount` for which this update applies to, in its settled form.
	SettledSubaccount Subaccount
	// A list of changes to make to any `AssetPositions` in the `Subaccount`.
	AssetUpdates []AssetUpdate
	// A list of changes to make to any `PerpetualPositions` in the `Subaccount`.
	PerpetualUpdates []PerpetualUpdate
}
