package upgrades

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
)

// Upgrade defines a struct containing necessary fields that a MsgSoftwareUpgrade
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade struct {
	// Upgrade version name, for the upgrade handler, e.g. `v7`
	UpgradeName string

	// Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
	StoreUpgrades store.StoreUpgrades
}

// Fork defines a struct containing the requisite fields for a non-software upgrade proposal
// Hard Fork at a given height to implement.
type Fork struct {
	// Upgrade version name, for the upgrade handler, e.g. `v7`
	UpgradeName string
	// Height the upgrade occurs at.
	UpgradeHeight int64
	// Upgrade info for this fork.
	UpgradeInfo string
}
