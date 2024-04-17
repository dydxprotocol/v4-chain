package constants

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	perpetualstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
)

var (
	// MsgUpdateParams replaces MsgCompleteBridge, since we delete the bridge module
	// Note: the crucial component here is that Authority is set to the module address
	// of the delaymsg module
	TestMsg1 = &perpetualstypes.MsgUpdateParams{
		Authority: types.ModuleAddress.String(),
		Params:    PerpetualsGenesisParams,
	}
	TestMsg2 = &perpetualstypes.MsgUpdateParams{
		Authority: types.ModuleAddress.String(),
		Params: perpetualstypes.Params{
			FundingRateClampFactorPpm: TestFundingRateClampFactorPpm + 1,
			PremiumVoteClampFactorPpm: TestPremiumVoteClampFactorPpm,
			MinNumVotesPerSample:      TestMinNumVotesPerSample,
		},
	}
	TestMsg3 = &perpetualstypes.MsgUpdateParams{
		Authority: types.ModuleAddress.String(),
		Params: perpetualstypes.Params{
			FundingRateClampFactorPpm: TestFundingRateClampFactorPpm + 2,
			PremiumVoteClampFactorPpm: TestPremiumVoteClampFactorPpm,
			MinNumVotesPerSample:      TestMinNumVotesPerSample,
		},
	}
	NoHandlerMsg = &testdata.TestMsg{Signers: []string{types.ModuleAddress.String()}}
)
