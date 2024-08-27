package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
)

var (
	// AppInjectedMsgSamples are msgs that are injected into the block by the proposing validator.
	// These messages are reserved for proposing validator's use only.
	AppInjectedMsgSamples = map[string]sdk.Msg{
		// clob
		"/dydxprotocol.clob.MsgProposedOperations": &clobtypes.MsgProposedOperations{
			OperationsQueue: make([]clobtypes.OperationRaw, 0),
		},
		"/dydxprotocol.clob.MsgProposedOperationsResponse": nil,

		// perpetuals
		"/dydxprotocol.perpetuals.MsgAddPremiumVotes": &perptypes.MsgAddPremiumVotes{
			Votes: []perptypes.FundingPremium{
				{PerpetualId: 0, PremiumPpm: 1_000},
			},
		},
		"/dydxprotocol.perpetuals.MsgAddPremiumVotesResponse": nil,
	}
)
