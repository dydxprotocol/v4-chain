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
	// Leverage configuration for this subaccount (perpetualId -> custom imf).
	// nil means no leverage configured (use default margin requirements).
	LeverageMap map[uint32]uint32
}

func (u *SettledUpdate) GetAssetUpdates() map[uint32]AssetUpdate {
	updates := make(map[uint32]AssetUpdate)
	for _, update := range u.AssetUpdates {
		updates[update.AssetId] = update.DeepCopy()
	}
	return updates
}

func (u *SettledUpdate) GetPerpetualUpdates() map[uint32]PerpetualUpdate {
	updates := make(map[uint32]PerpetualUpdate)
	for _, update := range u.PerpetualUpdates {
		updates[update.PerpetualId] = update.DeepCopy()
	}
	return updates
}
