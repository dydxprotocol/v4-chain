package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	blocktimetypes "github.com/dydxprotocol/v4/x/blocktime/types"
	bridgetypes "github.com/dydxprotocol/v4/x/bridge/types"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	perpetualstypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
)

// IsSingleAppInjectedMsg returns true if the given list of msgs contains an "app-injected msg"
// and it's the only msg in the list. Otherwise, returns false.
func IsSingleAppInjectedMsg(msgs []sdk.Msg) bool {
	return len(msgs) == 1 && IsAppInjectedMsg(msgs[0])
}

// IsAppInjectedMsg returns true if the given msg is an "app-injected msg".
// Otherwise, returns false.
func IsAppInjectedMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		*blocktimetypes.MsgIsDelayedBlock,
		*bridgetypes.MsgAcknowledgeBridge,
		*clobtypes.MsgProposedOperations,
		*perpetualstypes.MsgAddPremiumVotes,
		*pricestypes.MsgUpdateMarketPrices:
		return true
	}
	return false
}
