package ante

import (
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perpetualstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		// ------- Custom modules

		// clob
		*clobtypes.MsgProposedOperations,

		// perpetuals
		*perpetualstypes.MsgAddPremiumVotes:

		return true
	}
	return false
}
