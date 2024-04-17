package v_3_0_0

import (
	store "cosmossdk.io/store/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/upgrades"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
)

const (
	UpgradeName = "v3.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			ratelimittypes.StoreKey,
			icahosttypes.StoreKey,
		},
	},
}
