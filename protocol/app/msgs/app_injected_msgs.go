package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
)

var (
	// AppInjectedMsgSamples are msgs that are injected into the block by the proposing validator.
	// These messages are reserved for proposing validator's use only.
	AppInjectedMsgSamples = map[string]sdk.Msg{
		// clob
		"/dydxprotocol.clob.MsgProposedOperations": &clobtypes.MsgProposedOperations{
			Proposer:        "abc",
			OperationsQueue: make([]clobtypes.Operation, 0),
		},
		"/dydxprotocol.clob.MsgProposedOperationsResponse": nil,

		// perpetuals
		"/dydxprotocol.perpetuals.MsgAddPremiumVotes": &perptypes.MsgAddPremiumVotes{
			Proposer: "abc",
			Votes: []perptypes.FundingPremium{
				{PerpetualId: 0, PremiumPpm: 1_000},
			},
		},
		"/dydxprotocol.perpetuals.MsgAddPremiumVotesResponse": nil,

		// prices
		"/dydxprotocol.prices.MsgUpdateMarketPrices": &pricestypes.MsgUpdateMarketPrices{
			Proposer: "abc",
			MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
				pricestypes.NewMarketPriceUpdate(constants.MarketId0, 123_000),
			},
		},
		"/dydxprotocol.prices.MsgUpdateMarketPricesResponse": nil,
	}
)
